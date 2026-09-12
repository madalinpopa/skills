package install

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/madalinpopa/skills/internal/skill"
)

type Installation struct {
	Name    string
	Targets []Destination
	Issues  []Issue
}

func Scan(dests []Destination) ([]Installation, error) {
	byName := map[string]*Installation{}
	for _, dest := range dests {
		entries, err := os.ReadDir(dest.Dir)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			dir := filepath.Join(dest.Dir, entry.Name())
			info, err := os.Stat(dir)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, err
			}
			if !info.IsDir() {
				continue
			}
			_, err = ReadLock(dir)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			issue, damaged := lockIssue(err)
			if err != nil && !damaged {
				return nil, err
			}
			found, ok := byName[entry.Name()]
			if !ok {
				found = &Installation{Name: entry.Name()}
				byName[entry.Name()] = found
			}
			if damaged {
				found.Issues = append(found.Issues, issue)
			}
			found.Targets = append(found.Targets, dest)
		}
	}
	installed := make([]Installation, 0, len(byName))
	for _, name := range slices.Sorted(maps.Keys(byName)) {
		installed = append(installed, *byName[name])
	}
	return installed, nil
}

func Describe(inst Installation) (string, error) {
	dir := filepath.Join(inst.Targets[0].Dir, inst.Name)
	root, err := os.OpenRoot(dir)
	if err != nil {
		return "", err
	}
	data, err := root.ReadFile("SKILL.md")
	if err = errors.Join(err, root.Close()); err != nil {
		return "", err
	}
	meta, err := skill.Describe(data)
	if err != nil {
		return "", fmt.Errorf("%s: %w", filepath.Join(dir, "SKILL.md"), err)
	}
	return meta.Description, nil
}

func Find(installed []Installation, names []string) ([]Installation, error) {
	if len(names) == 0 {
		return installed, nil
	}
	byName := map[string]Installation{}
	for _, inst := range installed {
		byName[inst.Name] = inst
	}
	var found []Installation
	var missing []string
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		inst, ok := byName[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		found = append(found, inst)
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("not installed: %s", strings.Join(missing, ", "))
	}
	return found, nil
}

func UpdateRequests(storeDir string, catalog []skill.Skill, installed []Installation) ([]Request, error) {
	byName := map[string]skill.Skill{}
	for _, s := range catalog {
		byName[s.Name] = s
	}
	reqs := make([]Request, 0, len(installed))
	for _, inst := range installed {
		s, ok := byName[inst.Name]
		if !ok {
			reqs = append(reqs, Request{Name: inst.Name, Unavailable: true})
			continue
		}
		req, err := request(storeDir, s, inst.Targets)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}
