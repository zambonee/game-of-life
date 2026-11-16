package game

import (
	"testing"
)

func TestNew(t *testing.T) {
	cells := []Coordinates{{X: 1, Y: 2}, {X: -3, Y: 4}, {X: 1, Y: 2}}
	state := New(cells)

	if len(state.cells) != 2 {
		t.Fatalf("New() created %d cells, want 2", len(state.cells))
	}
}

func TestTickForward(t *testing.T) {
	tests := []struct {
		name  string
		start []Coordinates
		want  []Coordinates
	}{
		{
			name:  "underpopulation",
			start: []Coordinates{{X: 0, Y: 0}},
		},
		{
			name: "birth even when neighbors die",
			start: []Coordinates{
				{X: -1, Y: -1}, {X: 1, Y: -1}, {X: -1, Y: 1},
			},
			want: []Coordinates{{X: 0, Y: 0}},
		},
		{
			name: "stable block",
			start: []Coordinates{
				{X: 0, Y: 0}, {X: 1, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1},
			},
			want: []Coordinates{
				{X: 0, Y: 0}, {X: 1, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1},
			},
		},
		{
			name: "blinker",
			start: []Coordinates{
				{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0},
			},
			want: []Coordinates{
				{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1},
			},
		},
		{
			name: "overpopulation",
			start: []Coordinates{
				{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1},
				{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0},
				{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1},
			},
			want: []Coordinates{
				{X: -1, Y: -1}, {X: 1, Y: -1},
				{X: -1, Y: 1}, {X: 1, Y: 1},
				{X: 0, Y: -2}, {X: 2, Y: 0}, {X: -2, Y: 0}, {X: 0, Y: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := New(tt.start)
			state.TickForward()
			assertLiveCells(t, state, tt.want)
		})
	}
}

func TestTickForwardTwice(t *testing.T) {
	start := []Coordinates{
		{X: 1, Y: 0},
		{X: 2, Y: 1},
		{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2},
	}
	firstState := []Coordinates{
		{X: 0, Y: 1}, {X: 2, Y: 1},
		{X: 1, Y: 2}, {X: 2, Y: 2},
		{X: 1, Y: 3},
	}
	lastState := []Coordinates{
		{X: 2, Y: 1},
		{X: 0, Y: 2}, {X: 2, Y: 2},
		{X: 1, Y: 3}, {X: 2, Y: 3},
	}
	state := New(start)

	state.TickForward()
	assertLiveCells(t, state, firstState)

	state.TickForward()
	assertLiveCells(t, state, lastState)
}

func assertLiveCells(t *testing.T, state *State, want []Coordinates) {
	t.Helper()

	wantSet := make(map[Coordinates]struct{}, len(want))
	for _, cell := range want {
		wantSet[cell] = struct{}{}
	}
	if len(state.cells) != len(wantSet) {
		t.Fatalf("got %d live cells, want %d: %v", len(state.cells), len(wantSet), state.cells)
	}
	for cell := range state.cells {
		if _, ok := wantSet[cell]; !ok {
			t.Errorf("unexpected live cell %v", cell)
		}
	}
	for cell := range wantSet {
		if _, ok := state.cells[cell]; !ok {
			t.Errorf("missing live cell %v", cell)
		}
	}
}
