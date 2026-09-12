package skill

import (
	"errors"
	"io/fs"
	"slices"
	"strings"
)

func Catalog(fsys fs.FS) ([]Skill, error) {
	all, err := Discover(fsys)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	published := slices.DeleteFunc(all, func(s Skill) bool {
		return s.Status != Published
	})
	slices.SortFunc(published, func(a, b Skill) int {
		return strings.Compare(a.Name, b.Name)
	})
	return published, nil
}
