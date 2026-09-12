package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/madalinpopa/skills/internal/install"
)

func (a app) find(installed []install.Installation, names, agents []string, global bool) ([]install.Installation, error) {
	chosen, err := install.Find(installed, names)
	missing, notInstalled := errors.AsType[install.NotInstalledError](err)
	if err == nil || global || !notInstalled {
		return chosen, err
	}
	targets, err := install.Targets(a.cfg, a.home, true, agents)
	if err != nil {
		return nil, missing
	}
	home, err := install.Scan(targets)
	if err != nil {
		return nil, missing
	}
	var elsewhere, nowhere []string
	for _, name := range missing.Names {
		if _, err := install.Find(home, []string{name}); err == nil {
			elsewhere = append(elsewhere, name)
		} else {
			nowhere = append(nowhere, name)
		}
	}
	if len(elsewhere) == 0 {
		return nil, missing
	}
	hint := fmt.Sprintf("%s: only installed globally, use --global", strings.Join(elsewhere, ", "))
	if len(nowhere) > 0 {
		return nil, fmt.Errorf("%w; %s", install.NotInstalledError{Names: nowhere}, hint)
	}
	return nil, errors.New(hint)
}
