package lab

import (
	"reflect"
	"testing"
)

func TestSaturationAndExpiry(t *testing.T) {
	report, err := Replay(Scenario{Capacity: 2, Events: []Event{
		{At: 0, Action: "begin", Link: "a"},
		{At: 1, Action: "begin", Link: "b"},
		{At: 2, Action: "release", Link: "b"},
		{At: 31, Action: "begin", Link: "c"},
		{At: 32, Action: "begin", Link: "c"},
		{At: 33, Action: "begin", Link: "b"},
		{At: 40, Action: "release", Link: "a"},
		{At: 70, Action: "begin", Link: "b"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	checks := []struct {
		index   int
		outcome string
		expired []string
		active  []string
		idle    []string
	}{
		{3, "delayed: table full", []string{}, []string{"a"}, []string{"b"}},
		{4, "started: inserted entry", []string{"b"}, []string{"a", "c"}, []string{}},
		{5, "delayed: table full", []string{}, []string{"a", "c"}, []string{}},
		{7, "started: inserted entry", []string{"a"}, []string{"b", "c"}, []string{}},
	}
	for _, want := range checks {
		step := report.Steps[want.index]
		if step.Outcome != want.outcome || !reflect.DeepEqual(step.Expired, want.expired) ||
			!reflect.DeepEqual(step.Active, want.active) || !reflect.DeepEqual(step.Idle, want.idle) {
			t.Errorf("step %d = %+v, want %+v", want.index, step, want)
		}
	}
}

func TestReuseAndBusyLink(t *testing.T) {
	report, err := Replay(Scenario{Capacity: 1, Events: []Event{
		{At: 0, Action: "begin", Link: "a"},
		{At: 1, Action: "release", Link: "a"},
		{At: 30, Action: "begin", Link: "a"},
		{At: 100, Action: "begin", Link: "a"},
		{At: 100, Action: "release", Link: "a"},
		{At: 130, Action: "begin", Link: "a"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"started: inserted entry", "released: entry idle", "started: reused idle entry",
		"delayed: link busy", "released: entry idle", "started: inserted entry"}
	for i, step := range report.Steps {
		if step.Outcome != want[i] {
			t.Errorf("step %d outcome = %q, want %q", i, step.Outcome, want[i])
		}
	}
}

func TestInvalidScenario(t *testing.T) {
	for name, s := range map[string]Scenario{
		"capacity": {Capacity: 0, Events: []Event{{0, "begin", "a"}}},
		"empty":    {Capacity: 1},
		"negative": {Capacity: 1, Events: []Event{{-1, "begin", "a"}}},
		"reverse":  {Capacity: 1, Events: []Event{{2, "begin", "a"}, {1, "release", "a"}}},
		"bad link": {Capacity: 1, Events: []Event{{0, "begin", "<a>"}}},
		"bad verb": {Capacity: 1, Events: []Event{{0, "retry", "a"}}},
		"release":  {Capacity: 1, Events: []Event{{0, "release", "a"}}},
	} {
		t.Run(name, func(t *testing.T) {
			report, err := Replay(s)
			if err == nil || len(report.Steps) != 0 {
				t.Fatalf("Replay = %+v, %v, want error and no steps", report, err)
			}
		})
	}
}
