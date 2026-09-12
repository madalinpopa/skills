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
	Name        string
	Description string
	Targets     []Destination
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
			if !entry.IsDir() {
				continue
			}
			dir := filepath.Join(dest.Dir, entry.Name())
			_, err = ReadLock(dir)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, err
			}
			found, ok := byName[entry.Name()]
			if !ok {
				var meta skill.Metadata
				meta, err = describe(dir)
				if err != nil {
					return nil, err
				}
				found = &Installation{Name: entry.Name(), Description: meta.Description}
				byName[entry.Name()] = found
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

func describe(dir string) (skill.Metadata, error) {
	file := filepath.Join(dir, "SKILL.md")
	data, err := os.ReadFile(file)
	if err != nil {
		return skill.Metadata{}, err
	}
	meta, err := skill.Describe(data)
	if err != nil {
		return skill.Metadata{}, fmt.Errorf("%s: %w", file, err)
	}
	return meta, nil
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
