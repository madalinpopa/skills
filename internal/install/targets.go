package install

import (
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/madalinpopa/skills/internal/config"
	"github.com/madalinpopa/skills/internal/skill"
)

type Variant string

const (
	VariantClaude Variant = "claude"
	VariantAgents Variant = "agents"
)

type Destination struct {
	Dir     string
	Variant Variant
}

var ErrUnknownAgent = errors.New("unknown agent")

func Targets(cfg config.Config, root string, global bool, agents []string) ([]Destination, error) {
	if len(agents) == 0 {
		agents = cfg.Defaults.Agents
	}
	byDir := map[string]Destination{}
	for _, name := range agents {
		agent, ok := cfg.Agents[name]
		if !ok {
			available := strings.Join(slices.Sorted(maps.Keys(cfg.Agents)), ", ")
			return nil, fmt.Errorf("%w %q; available: %s", ErrUnknownAgent, name, available)
		}
		dest := Destination{Dir: filepath.Join(root, agent.Project), Variant: variantOf(name)}
		if global {
			dest.Dir = agent.GlobalPath(root)
		}
		if existing, ok := byDir[dest.Dir]; ok && existing.Variant != dest.Variant {
			return nil, fmt.Errorf("agent %q shares %s with an agent that needs a different format", name, dest.Dir)
		}
		byDir[dest.Dir] = dest
	}
	dests := slices.Collect(maps.Values(byDir))
	slices.SortFunc(dests, func(a, b Destination) int {
		return strings.Compare(a.Dir, b.Dir)
	})
	return dests, nil
}

func Select(catalog []skill.Skill, names []string) ([]skill.Skill, error) {
	byName := map[string]skill.Skill{}
	for _, s := range catalog {
		byName[s.Name] = s
	}
	var selected []skill.Skill
	var missing []string
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		s, ok := byName[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		selected = append(selected, s)
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("not published in the store: %s", strings.Join(missing, ", "))
	}
	return selected, nil
}

func variantOf(agent string) Variant {
	if agent == "claude" {
		return VariantClaude
	}
	return VariantAgents
}
