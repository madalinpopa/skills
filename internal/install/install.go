package install

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"time"

	"github.com/madalinpopa/skills/internal/skill"
)

const LockFile = ".skill-lock.json"

type Lock struct {
	Name      string            `json:"name"`
	Source    string            `json:"source"`
	Commit    string            `json:"commit"`
	Installed time.Time         `json:"installed"`
	Files     map[string]string `json:"files"`
}

type Desired struct {
	Dir   string
	Files map[string]skill.File
}

type Request struct {
	Name        string
	Unavailable bool
	Targets     []Desired
}

type Installer struct {
	Source  string
	Commit  string
	Backups string
	Force   bool
	DryRun  bool
	Now     func() time.Time
}

func (i Installer) Install(reqs []Request) ([]SkillPlan, error) {
	specs := make([]Spec, 0, len(reqs))
	requests := map[string]Request{}
	specByName := map[string]Spec{}
	for _, req := range reqs {
		spec, err := i.spec(req)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
		requests[req.Name] = req
		specByName[req.Name] = spec
	}
	plans := Plan(specs)
	for idx, plan := range plans {
		if plan.State == StateConflict && i.Force {
			forced, err := i.force(specByName[plan.Name])
			plans[idx] = forced
			if err != nil {
				return plans[:idx+1], fmt.Errorf("install %s: %w", plan.Name, err)
			}
			plan = forced
		}
		if i.DryRun || (plan.State != StateAdd && plan.State != StateUpdate) {
			continue
		}
		if err := i.apply(requests[plan.Name], plan); err != nil {
			plans[idx].State = StateFailed
			return plans[:idx+1], fmt.Errorf("install %s: %w", plan.Name, err)
		}
	}
	return plans, nil
}

func (i Installer) force(spec Spec) (SkillPlan, error) {
	var backups []string
	for idx := range spec.Targets {
		target := &spec.Targets[idx]
		target.Base = target.Have
		if i.DryRun {
			continue
		}
		path, err := i.backup(target.Dir)
		if err != nil {
			return SkillPlan{Name: spec.Name, State: StateFailed, Backups: backups}, err
		}
		if path != "" {
			backups = append(backups, path)
		}
	}
	plan := planSkill(spec)
	plan.Backups = backups
	return plan, nil
}

func (i Installer) spec(req Request) (Spec, error) {
	spec := Spec{Name: req.Name, Available: !req.Unavailable, Source: i.Source}
	for _, desired := range req.Targets {
		t, err := inspect(desired.Dir)
		if err != nil {
			return Spec{}, err
		}
		target := Target{
			Dir:    desired.Dir,
			Want:   hashFiles(desired.Files),
			Have:   t.files,
			Issues: append(t.issues, collisions(desired.Dir, desired.Files, t)...),
		}
		lock, err := ReadLock(desired.Dir)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return Spec{}, err
		}
		if err == nil {
			target.Source = lock.Source
			target.Base = lock.Files
			if target.Base == nil {
				target.Base = Files{}
			}
		}
		spec.Targets = append(spec.Targets, target)
	}
	return spec, nil
}

func (i Installer) apply(req Request, plan SkillPlan) error {
	stage := &staging{}
	defer stage.discard()
	for idx, desired := range req.Targets {
		lock := Lock{
			Name:      req.Name,
			Source:    i.Source,
			Commit:    i.Commit,
			Installed: i.Now(),
			Files:     hashFiles(desired.Files),
		}
		if err := stage.target(desired, plan.Targets[idx], lock); err != nil {
			return err
		}
	}
	return stage.publish()
}

func ReadLock(dir string) (Lock, error) {
	data, err := os.ReadFile(filepath.Join(dir, LockFile))
	if err != nil {
		return Lock{}, err
	}
	var lock Lock
	if err := json.Unmarshal(data, &lock); err != nil {
		return Lock{}, fmt.Errorf("%s: %w", filepath.Join(dir, LockFile), err)
	}
	return lock, nil
}

type staged struct {
	tmp   string
	final string
}

type stagedTarget struct {
	dir     string
	files   []staged
	removes []string
	lock    staged
}

type staging struct {
	targets []stagedTarget
	dirs    []string
}

func (s *staging) target(desired Desired, plan TargetPlan, lock Lock) error {
	target := stagedTarget{dir: desired.Dir}
	for _, change := range plan.Files {
		final := filepath.Join(desired.Dir, filepath.FromSlash(change.Path))
		switch change.Action {
		case ActionAdd, ActionUpdate:
			file := desired.Files[change.Path]
			written, err := s.write(final, file.Data, perm(file.Mode))
			target.files = append(target.files, written)
			if err != nil {
				s.targets = append(s.targets, target)
				return err
			}
		case ActionRemove:
			target.removes = append(target.removes, final)
		}
	}
	data, err := json.Marshal(lock, json.Deterministic(true), jsontext.WithIndent("  "))
	if err != nil {
		return err
	}
	target.lock, err = s.write(filepath.Join(desired.Dir, LockFile), append(data, '\n'), 0o644)
	s.targets = append(s.targets, target)
	return err
}

func (s *staging) write(final string, data []byte, mode fs.FileMode) (staged, error) {
	if err := s.mkdir(filepath.Dir(final)); err != nil {
		return staged{}, err
	}
	tmp, err := os.CreateTemp(filepath.Dir(final), "."+filepath.Base(final)+".*.tmp")
	if err != nil {
		return staged{}, err
	}
	file := staged{tmp: tmp.Name(), final: final}
	_, writeErr := tmp.Write(data)
	if err := errors.Join(writeErr, tmp.Close()); err != nil {
		return file, err
	}
	return file, os.Chmod(tmp.Name(), mode)
}

func (s *staging) mkdir(dir string) error {
	var missing []string
	for d := dir; ; d = filepath.Dir(d) {
		if _, err := os.Stat(d); err == nil {
			break
		}
		missing = append(missing, d)
		if filepath.Dir(d) == d {
			break
		}
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	s.dirs = append(s.dirs, missing...)
	return nil
}

func (s *staging) publish() error {
	for _, target := range s.targets {
		if err := target.publish(); err != nil {
			return fmt.Errorf("publish %s: %w", target.dir, err)
		}
	}
	return nil
}

func (t stagedTarget) publish() error {
	for _, f := range t.files {
		if err := os.Rename(f.tmp, f.final); err != nil {
			return err
		}
	}
	for _, path := range t.removes {
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return os.Rename(t.lock.tmp, t.lock.final)
}

func (s *staging) discard() {
	for _, target := range s.targets {
		for _, f := range target.files {
			_ = os.Remove(f.tmp)
		}
		_ = os.Remove(target.lock.tmp)
	}
	for _, dir := range s.dirs {
		_ = os.Remove(dir)
	}
}

func Check(dir string) ([]Issue, error) {
	t, err := inspect(dir)
	return t.issues, err
}

type tree struct {
	files  Files
	dirs   map[string]bool
	issues []Issue
}

func inspect(dir string) (tree, error) {
	t := tree{files: Files{}, dirs: map[string]bool{}}
	info, err := os.Lstat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return t, nil
	}
	if err != nil {
		return tree{}, err
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		t.issues = []Issue{{Path: dir, Reason: "is a symlink"}}
		return t, nil
	}
	if !info.IsDir() {
		t.issues = []Issue{{Path: dir, Reason: "is a file where a directory is expected"}}
		return t, nil
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return tree{}, err
	}
	fsys := root.FS()
	err = fs.WalkDir(fsys, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil || p == "." {
			return err
		}
		full := filepath.Join(dir, filepath.FromSlash(p))
		switch {
		case entry.Type()&fs.ModeSymlink != 0:
			t.issues = append(t.issues, Issue{Path: full, Reason: "is a symlink"})
		case entry.IsDir():
			t.dirs[p] = true
		case !entry.Type().IsRegular():
			t.issues = append(t.issues, Issue{Path: full, Reason: "is not a regular file"})
		case entry.Name() != LockFile:
			data, err := fs.ReadFile(fsys, p)
			if err != nil {
				return err
			}
			t.files[p] = hash(data)
		}
		return nil
	})
	if err = errors.Join(err, root.Close()); err != nil {
		return tree{}, err
	}
	return t, nil
}

func collisions(dir string, want map[string]skill.File, t tree) []Issue {
	var issues []Issue
	seen := map[string]bool{}
	for _, p := range slices.Sorted(maps.Keys(want)) {
		if t.dirs[p] {
			issues = append(issues, Issue{Path: filepath.Join(dir, filepath.FromSlash(p)), Reason: "is a directory where a file is expected"})
		}
		for parent := path.Dir(p); parent != "."; parent = path.Dir(parent) {
			if _, ok := t.files[parent]; ok && !seen[parent] {
				seen[parent] = true
				issues = append(issues, Issue{Path: filepath.Join(dir, filepath.FromSlash(parent)), Reason: "is a file where a directory is expected"})
			}
		}
	}
	return issues
}

func hashDir(dir string) (Files, error) {
	t, err := inspect(dir)
	return t.files, err
}

func hashFiles(files map[string]skill.File) Files {
	hashes := make(Files, len(files))
	for path, file := range files {
		hashes[path] = hash(file.Data)
	}
	return hashes
}

func hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func perm(mode fs.FileMode) fs.FileMode {
	if mode&0o111 != 0 {
		return 0o755
	}
	return 0o644
}

func Requests(storeDir string, skills []skill.Skill, dests []Destination) ([]Request, error) {
	reqs := make([]Request, 0, len(skills))
	for _, s := range skills {
		req, err := request(storeDir, s, dests)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func request(storeDir string, s skill.Skill, dests []Destination) (Request, error) {
	req := Request{Name: s.Name}
	for _, dest := range dests {
		files, err := skill.Render(os.DirFS(filepath.Join(storeDir, filepath.FromSlash(s.Dir))), string(dest.Variant))
		if err != nil {
			return Request{}, fmt.Errorf("%s: %w", s.Name, err)
		}
		req.Targets = append(req.Targets, Desired{Dir: filepath.Join(dest.Dir, s.Name), Files: files})
	}
	return req, nil
}
