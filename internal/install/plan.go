package install

import (
	"maps"
	"path"
	"slices"
	"strings"
)

type Mode int8

const (
	ModeFile Mode = iota
	ModeExecutable
	ModeUnknown
)

type Entry struct {
	Hash string
	Mode Mode
}

type Files map[string]Entry

type Action string

const (
	ActionUnchanged Action = "unchanged"
	ActionAdd       Action = "add"
	ActionUpdate    Action = "update"
	ActionRemove    Action = "remove"
	ActionKeep      Action = "keep"
	ActionConflict  Action = "conflict"
)

type State string

const (
	StateUnchanged   State = "unchanged"
	StateAdd         State = "add"
	StateUpdate      State = "update"
	StateConflict    State = "conflict"
	StateRemove      State = "remove"
	StateUnavailable State = "unavailable"
	StateForeign     State = "foreign"
	StateUnsupported State = "unsupported"
	StateFailed      State = "failed"
)

type Issue struct {
	Path   string
	Reason string
}

type FileChange struct {
	Path   string
	Action Action
}

type Target struct {
	Dir    string
	Source string
	Want   Files
	Have   Files
	Base   Files
	Issues []Issue
}

type Spec struct {
	Name      string
	Available bool
	Source    string
	Targets   []Target
}

type TargetPlan struct {
	Dir   string
	Files []FileChange
}

type SkillPlan struct {
	Name      string
	State     State
	Source    string
	Managed   bool
	Forced    bool
	Conflicts []string
	Backups   []string
	Issues    []Issue
	Targets   []TargetPlan
}

func Plan(specs []Spec) []SkillPlan {
	plans := make([]SkillPlan, 0, len(specs))
	for _, spec := range specs {
		plans = append(plans, planSkill(spec))
	}
	sortByName(plans)
	return plans
}

func sortByName(plans []SkillPlan) {
	slices.SortFunc(plans, func(a, b SkillPlan) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func PlanFiles(want, have, base Files) []FileChange {
	var changes []FileChange
	for _, p := range paths(want, have, base) {
		changes = append(changes, FileChange{Path: p, Action: fileAction(p, want, have, base)})
	}
	return changes
}

func paths(sets ...Files) []string {
	union := map[string]struct{}{}
	for _, set := range sets {
		for p := range set {
			union[p] = struct{}{}
		}
	}
	return slices.Sorted(maps.Keys(union))
}

func fileAction(p string, want, have, base Files) Action {
	w, wanted := want[p]
	h, present := have[p]
	b, tracked := base[p]
	switch {
	case wanted && !present:
		return ActionAdd
	case wanted && h == w:
		return ActionUnchanged
	case wanted && tracked && h == b:
		return ActionUpdate
	case wanted:
		return ActionConflict
	case !present:
		return ActionUnchanged
	case tracked && h == b:
		return ActionRemove
	case tracked:
		return ActionConflict
	default:
		return ActionKeep
	}
}

func planSkill(spec Spec) SkillPlan {
	plan := SkillPlan{Name: spec.Name}
	if !spec.Available {
		plan.State = StateUnavailable
		return plan
	}
	managed, missing, changed := false, false, false
	for _, target := range spec.Targets {
		if target.Source != "" && target.Source != spec.Source {
			plan.State = StateForeign
			plan.Source = target.Source
			return plan
		}
		if target.Base == nil {
			missing = true
		} else {
			managed = true
		}
		plan.Issues = append(plan.Issues, target.Issues...)
		files := PlanFiles(target.Want, target.Have, target.Base)
		for _, change := range files {
			switch change.Action {
			case ActionConflict:
				plan.Conflicts = append(plan.Conflicts, path.Join(target.Dir, change.Path))
			case ActionAdd, ActionUpdate, ActionRemove:
				changed = true
			}
		}
		plan.Targets = append(plan.Targets, TargetPlan{Dir: target.Dir, Files: files})
	}
	plan.Managed = managed
	switch {
	case len(plan.Issues) > 0:
		plan.State = StateUnsupported
		plan.Conflicts = nil
		plan.Targets = nil
	case len(plan.Conflicts) > 0:
		plan.State = StateConflict
	case missing:
		plan.State = StateAdd
	case changed:
		plan.State = StateUpdate
	default:
		plan.State = StateUnchanged
	}
	return plan
}
