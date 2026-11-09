package lab

import (
	"fmt"
	"regexp"
	"sort"
)

const IdleSeconds = 30

type Scenario struct {
	Capacity int     `json:"capacity"`
	Events   []Event `json:"events"`
}

type Event struct {
	At     int    `json:"at_seconds"`
	Action string `json:"action"`
	Link   string `json:"link"`
}

type Step struct {
	Event
	Outcome string   `json:"outcome"`
	Expired []string `json:"expired"`
	Active  []string `json:"active"`
	Idle    []string `json:"idle"`
}

type Report struct {
	Capacity    int    `json:"capacity"`
	IdleSeconds int    `json:"idle_seconds"`
	Steps       []Step `json:"steps"`
}

type entry struct {
	active    bool
	idleSince int
}

var linkName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)

// Replay advances a bounded table through a script in logical seconds.
func Replay(s Scenario) (Report, error) {
	if s.Capacity < 1 || s.Capacity > 64 {
		return Report{}, fmt.Errorf("capacity must be between 1 and 64")
	}
	if len(s.Events) == 0 || len(s.Events) > 500 {
		return Report{}, fmt.Errorf("expected 1 to 500 events")
	}
	report := Report{Capacity: s.Capacity, IdleSeconds: IdleSeconds, Steps: make([]Step, 0, len(s.Events))}
	table := make(map[string]entry, s.Capacity)
	lastAt := 0
	for i, event := range s.Events {
		if event.At < lastAt {
			return Report{}, fmt.Errorf("event %d: time must be nonnegative and nondecreasing", i+1)
		}
		lastAt = event.At
		if !linkName.MatchString(event.Link) {
			return Report{}, fmt.Errorf("event %d: invalid link name", i+1)
		}
		step := Step{Event: event, Expired: []string{}, Active: []string{}, Idle: []string{}}
		for name, e := range table {
			if !e.active && event.At-e.idleSince >= IdleSeconds {
				delete(table, name)
				step.Expired = append(step.Expired, name)
			}
		}
		e, found := table[event.Link]
		switch event.Action {
		case "begin":
			switch {
			case e.active:
				step.Outcome = "delayed: link busy"
			case !found && len(table) == s.Capacity:
				step.Outcome = "delayed: table full"
			case found:
				table[event.Link] = entry{active: true}
				step.Outcome = "started: reused idle entry"
			default:
				table[event.Link] = entry{active: true}
				step.Outcome = "started: inserted entry"
			}
		case "release":
			if !e.active {
				return Report{}, fmt.Errorf("event %d: release requires an active link", i+1)
			}
			table[event.Link] = entry{idleSince: event.At}
			step.Outcome = "released: entry idle"
		default:
			return Report{}, fmt.Errorf("event %d: unknown action %q", i+1, event.Action)
		}
		for name, e := range table {
			if e.active {
				step.Active = append(step.Active, name)
			} else {
				step.Idle = append(step.Idle, name)
			}
		}
		sort.Strings(step.Expired)
		sort.Strings(step.Active)
		sort.Strings(step.Idle)
		report.Steps = append(report.Steps, step)
	}
	return report, nil
}
