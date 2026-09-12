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
)

func (i Installer) backup(dir string) (string, error) {
	have, err := hashDir(dir)
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	identity := strings.TrimPrefix(strings.TrimPrefix(abs, filepath.VolumeName(abs)), string(filepath.Separator))
	dest := filepath.Join(i.Backups, i.Now().UTC().Format("20060102-150405"), identity)
	if err = os.CopyFS(dest, os.DirFS(dir)); err != nil {
		return "", fmt.Errorf("backup %s: %w", dir, err)
	}
	copied, err := hashDir(dest)
	if err != nil {
		return "", err
	}
	if !maps.Equal(copied, have) {
		return "", fmt.Errorf("backup %s: copy does not match the original", dest)
	}
	return dest, nil
}

func (i Installer) Remove(installed []Installation) ([]SkillPlan, error) {
	plans := make([]SkillPlan, 0, len(installed))
	for _, inst := range installed {
		plan, err := i.remove(inst)
		if err != nil {
			return nil, fmt.Errorf("remove %s: %w", inst.Name, err)
		}
		plans = append(plans, plan)
	}
	slices.SortFunc(plans, func(a, b SkillPlan) int {
		return strings.Compare(a.Name, b.Name)
	})
	return plans, nil
}

func (i Installer) remove(inst Installation) (SkillPlan, error) {
	plan := SkillPlan{Name: inst.Name, State: StateRemove}
	dirs := make([]string, 0, len(inst.Targets))
	for _, target := range inst.Targets {
		dir := filepath.Join(target.Dir, inst.Name)
		have, err := hashDir(dir)
		if err != nil {
			return SkillPlan{}, err
		}
		lock, err := ReadLock(dir)
		if err != nil {
			return SkillPlan{}, err
		}
		for _, p := range slices.Sorted(maps.Keys(have)) {
			if lock.Files[p] != have[p] {
				plan.Conflicts = append(plan.Conflicts, filepath.Join(dir, filepath.FromSlash(p)))
			}
		}
		dirs = append(dirs, dir)
	}
	if len(plan.Conflicts) > 0 && !i.Force {
		plan.State = StateConflict
		return plan, nil
	}
	for _, dir := range dirs {
		path, err := i.backup(dir)
		if err != nil {
			return SkillPlan{}, err
		}
		plan.Backups = append(plan.Backups, path)
	}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return SkillPlan{}, err
		}
	}
	return plan, nil
}
