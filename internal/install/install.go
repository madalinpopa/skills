package install

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
			if err != nil {
				return nil, fmt.Errorf("install %s: %w", plan.Name, err)
			}
			plans[idx] = forced
			plan = forced
		}
		if i.DryRun || (plan.State != StateAdd && plan.State != StateUpdate) {
			continue
		}
		if err := i.apply(requests[plan.Name], plan); err != nil {
			return nil, fmt.Errorf("install %s: %w", plan.Name, err)
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
			return SkillPlan{}, err
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
		have, err := hashDir(desired.Dir)
		if err != nil {
			return Spec{}, err
		}
		target := Target{Dir: desired.Dir, Want: hashFiles(desired.Files), Have: have}
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

type staging struct {
	files   []staged
	removes []string
	dirs    []string
}

func (s *staging) target(desired Desired, plan TargetPlan, lock Lock) error {
	for _, change := range plan.Files {
		final := filepath.Join(desired.Dir, filepath.FromSlash(change.Path))
		switch change.Action {
		case ActionAdd, ActionUpdate:
			file := desired.Files[change.Path]
			if err := s.write(final, file.Data, perm(file.Mode)); err != nil {
				return err
			}
		case ActionRemove:
			s.removes = append(s.removes, final)
		}
	}
	data, err := json.Marshal(lock, json.Deterministic(true), jsontext.WithIndent("  "))
	if err != nil {
		return err
	}
	return s.write(filepath.Join(desired.Dir, LockFile), append(data, '\n'), 0o644)
}

func (s *staging) write(final string, data []byte, mode fs.FileMode) error {
	if err := s.mkdir(filepath.Dir(final)); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(final), "."+filepath.Base(final)+".*.tmp")
	if err != nil {
		return err
	}
	s.files = append(s.files, staged{tmp: tmp.Name(), final: final})
	_, writeErr := tmp.Write(data)
	if err := errors.Join(writeErr, tmp.Close()); err != nil {
		return err
	}
	return os.Chmod(tmp.Name(), mode)
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
	for _, f := range s.files {
		if err := os.Rename(f.tmp, f.final); err != nil {
			return err
		}
	}
	s.files = nil
	s.dirs = nil
	for _, path := range s.removes {
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (s *staging) discard() {
	for _, f := range s.files {
		_ = os.Remove(f.tmp)
	}
	for _, dir := range s.dirs {
		_ = os.Remove(dir)
	}
}

func hashDir(dir string) (Files, error) {
	root, err := os.OpenRoot(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return Files{}, nil
	}
	if err != nil {
		return nil, err
	}
	files, err := hashTree(root.FS())
	return files, errors.Join(err, root.Close())
}

func hashTree(fsys fs.FS) (Files, error) {
	files := Files{}
	err := fs.WalkDir(fsys, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.Type().IsRegular() || entry.Name() == LockFile {
			return nil
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		files[p] = hash(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
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
