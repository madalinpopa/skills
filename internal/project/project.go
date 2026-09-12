package project

import (
	"fmt"
	"os"
	"path/filepath"
)

type Scope struct {
	Root    string
	Warning string
}

func Detect(start string) (Scope, error) {
	start, err := filepath.Abs(start)
	if err != nil {
		return Scope{}, err
	}
	info, err := os.Stat(start)
	if err != nil {
		return Scope{}, err
	}
	if !info.IsDir() {
		return Scope{}, fmt.Errorf("%s is not a directory", start)
	}
	if root, ok := gitRoot(start); ok {
		return Scope{Root: root}, nil
	}
	return Scope{
		Root:    start,
		Warning: "not inside a Git repository, using " + start,
	}, nil
}

func gitRoot(dir string) (string, bool) {
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
