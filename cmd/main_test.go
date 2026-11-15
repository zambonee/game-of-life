package main

import (
	"strings"
	"testing"

	"github.com/zambonee/game-of-life/internal/game"
)

func TestParseLife106(t *testing.T) {
	input := "#Life 1.06\n\n-9223372036854775808 9223372036854775807\n0 -1\n"

	got, err := parseLife106(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseLife106() error = %v", err)
	}
	want := []game.Coordinates{
		{X: -9223372036854775808, Y: 9223372036854775807},
		{X: 0, Y: -1},
	}
	if len(got) != len(want) {
		t.Fatalf("parseLife106() returned %d cells, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("parseLife106()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestParseLife106RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "missing header", input: "0 1\n"},
		{name: "missing coordinate", input: "#Life 1.06\n0\n"},
		{name: "invalid integer", input: "#Life 1.06\nnope 1\n"},
		{name: "integer overflow", input: "#Life 1.06\n9223372036854775808 0\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseLife106(strings.NewReader(tt.input)); err == nil {
				t.Fatal("parseLife106() error = nil, want an error")
			}
		})
	}
}

func TestRunAdvancesTenGenerations(t *testing.T) {
	input := "#Life 1.06\n1 0\n2 1\n0 2\n1 2\n2 2\n"
	var output strings.Builder

	if err := run(strings.NewReader(input), &output); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	want := "#Life 1.06\n2 4\n3 5\n4 3\n4 4\n4 5\n"
	if got := output.String(); got != want {
		t.Errorf("run() output:\n%s\nwant:\n%s", got, want)
	}
}
