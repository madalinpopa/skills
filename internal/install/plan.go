package install

import (
	"maps"
	"path"
	"slices"
	"strings"
)

type Files map[string]string

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
)

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
	Conflicts []string
	Backups   []string
	Targets   []TargetPlan
}

func Plan(specs []Spec) []SkillPlan {
	plans := make([]SkillPlan, 0, len(specs))
	for _, spec := range specs {
		plans = append(plans, planSkill(spec))
	}
	slices.SortFunc(plans, func(a, b SkillPlan) int {
		return strings.Compare(a.Name, b.Name)
	})
	return plans
}

func PlanFiles(want, have, base Files) []FileChange {
	union := map[string]struct{}{}
	for p := range maps.Keys(want) {
		union[p] = struct{}{}
	}
	for p := range maps.Keys(have) {
		union[p] = struct{}{}
	}
	for p := range maps.Keys(base) {
		union[p] = struct{}{}
	}
	var changes []FileChange
	for _, p := range slices.Sorted(maps.Keys(union)) {
		changes = append(changes, FileChange{Path: p, Action: fileAction(p, want, have, base)})
	}
	return changes
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
	installed := false
	changed := false
	for _, target := range spec.Targets {
		if target.Source != "" && target.Source != spec.Source {
			plan.State = StateForeign
			plan.Source = target.Source
			return plan
		}
		if target.Base != nil {
			installed = true
		}
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
	switch {
	case len(plan.Conflicts) > 0:
		plan.State = StateConflict
	case !installed:
		plan.State = StateAdd
	case changed:
		plan.State = StateUpdate
	default:
		plan.State = StateUnchanged
	}
	return plan
}
