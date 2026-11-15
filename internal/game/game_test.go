package game

import "testing"

func TestNew(t *testing.T) {
	cells := []coordinates{{X: 1, Y: 2}, {X: -3, Y: 4}, {X: 1, Y: 2}}
	state := New(cells)

	if len(state.Cells) != 2 {
		t.Fatalf("New() created %d cells, want 2", len(state.Cells))
	}
}

func TestGetNeighbors(t *testing.T) {
	center := coordinates{X: 0, Y: 0}
	state := New([]coordinates{
		center,
		{X: -1, Y: -1},
		{X: 0, Y: -1},
		{X: 1, Y: -1},
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: -1, Y: 1},
		{X: 0, Y: 1},
	})
	state.Cells[coordinates{X: 1, Y: 1}] = emergent

	got := state.getNeighbors(center)
	if len(got.liveNeighbors) != 7 {
		t.Errorf("getNeighbors() found %d live neighbors, want 7", len(got.liveNeighbors))
	}
	if len(got.deadNeighbors) != 1 {
		t.Errorf("getNeighbors() found %d dead neighbors, want 1", len(got.deadNeighbors))
	}
	if got.deadNeighbors[0] != (coordinates{X: 1, Y: 1}) {
		t.Errorf("getNeighbors() dead neighbor = %v, want (1, 1)", got.deadNeighbors[0])
	}
}

func TestTickForward(t *testing.T) {
	tests := []struct {
		name  string
		start []coordinates
		want  []coordinates
	}{
		{
			name:  "underpopulation",
			start: []coordinates{{X: 0, Y: 0}},
		},
		{
			name: "stable block",
			start: []coordinates{
				{X: 0, Y: 0}, {X: 1, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1},
			},
			want: []coordinates{
				{X: 0, Y: 0}, {X: 1, Y: 0},
				{X: 0, Y: 1}, {X: 1, Y: 1},
			},
		},
		{
			name: "blinker",
			start: []coordinates{
				{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0},
			},
			want: []coordinates{
				{X: 0, Y: -1}, {X: 0, Y: 0}, {X: 0, Y: 1},
			},
		},
		{
			name: "overpopulation",
			start: []coordinates{
				{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1},
				{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0},
				{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1},
			},
			want: []coordinates{
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
	start := []coordinates{
		{X: 1, Y: 0},
		{X: 2, Y: 1},
		{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2},
	}
	firstState := []coordinates{
		{X: 0, Y: 1}, {X: 2, Y: 1},
		{X: 1, Y: 2}, {X: 2, Y: 2},
		{X: 1, Y: 3},
	}
	lastState := []coordinates{
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

func assertLiveCells(t *testing.T, state *State, want []coordinates) {
	t.Helper()

	wantSet := make(map[coordinates]struct{}, len(want))
	for _, cell := range want {
		wantSet[cell] = struct{}{}
	}
	if len(state.Cells) != len(wantSet) {
		t.Fatalf("got %d live cells, want %d: %v", len(state.Cells), len(wantSet), state.Cells)
	}
	for cell, cellState := range state.Cells {
		if cellState != alive {
			t.Errorf("cell %v has state %v after tick, want alive", cell, cellState)
		}
		if _, ok := wantSet[cell]; !ok {
			t.Errorf("unexpected live cell %v", cell)
		}
	}
	for cell := range wantSet {
		if _, ok := state.Cells[cell]; !ok {
			t.Errorf("missing live cell %v", cell)
		}
	}
}
