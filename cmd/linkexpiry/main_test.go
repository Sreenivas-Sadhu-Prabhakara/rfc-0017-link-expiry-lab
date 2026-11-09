package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"example.com/rfc-go-labs/lab"
)

func TestSampleReports(t *testing.T) {
	data, err := os.ReadFile("../../examples/saturation.json")
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run(nil, bytes.NewReader(data), &out); err != nil {
		t.Fatal(err)
	}
	var report lab.Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Steps) != 7 || report.Steps[4].Outcome != "delayed: table full" {
		t.Fatalf("unexpected sample trace: %+v", report)
	}
	out.Reset()
	if err := run([]string{"-html"}, bytes.NewReader(data), &out); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"30-second idle expiry", "delayed: table full", "synthetic event script"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("HTML missing %q", want)
		}
	}
}

func TestInvalidInput(t *testing.T) {
	for _, input := range []string{
		`{"capacity":1,"extra":true}`,
		`{"capacity":1,"events":[{"action":"begin","link":"a"}]} {}`,
		`{"capacity":1,"events":[{"action":"release","link":"a"}]}`,
		strings.Repeat(" ", (1<<20)+1),
	} {
		var out bytes.Buffer
		if err := run(nil, strings.NewReader(input), &out); err == nil || out.Len() != 0 {
			t.Fatalf("got %d output bytes and error %v for invalid input (%d bytes)", out.Len(), err, len(input))
		}
	}
}
