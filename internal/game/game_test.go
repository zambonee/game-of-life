package game

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	cells := []Coordinates{{X: 1, Y: 2}, {X: -3, Y: 4}, {X: 1, Y: 2}}
	state := New(cells)

	if len(state.cells) != 2 {
		t.Fatalf("New() created %d cells, want 2", len(state.cells))
	}
}

func TestGetNeighbors(t *testing.T) {
	center := Coordinates{X: 0, Y: 0}
	state := New([]Coordinates{
		center,
		{X: -1, Y: -1},
		{X: 0, Y: -1},
		{X: 1, Y: -1},
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: -1, Y: 1},
		{X: 0, Y: 1},
	})
	state.cells[Coordinates{X: 1, Y: 1}] = emergent

	got := state.getNeighbors(center)
	if len(got.liveNeighbors) != 7 {
		t.Errorf("getNeighbors() found %d live neighbors, want 7", len(got.liveNeighbors))
	}
	if len(got.deadNeighbors) != 1 {
		t.Errorf("getNeighbors() found %d dead neighbors, want 1", len(got.deadNeighbors))
	}
	if got.deadNeighbors[0] != (Coordinates{X: 1, Y: 1}) {
		t.Errorf("getNeighbors() dead neighbor = %v, want (1, 1)", got.deadNeighbors[0])
	}
}

func TestGetNeighborsAtInt64Boundaries(t *testing.T) {
	tests := []struct {
		name              string
		center            Coordinates
		wrappedCoordinate Coordinates
		wantDeadNeighbors int
	}{
		{
			name:              "minimum x",
			center:            Coordinates{X: math.MinInt64, Y: 0},
			wrappedCoordinate: Coordinates{X: math.MaxInt64, Y: 0},
			wantDeadNeighbors: 5,
		},
		{
			name:              "maximum x",
			center:            Coordinates{X: math.MaxInt64, Y: 0},
			wrappedCoordinate: Coordinates{X: math.MinInt64, Y: 0},
			wantDeadNeighbors: 5,
		},
		{
			name:              "minimum y",
			center:            Coordinates{X: 0, Y: math.MinInt64},
			wrappedCoordinate: Coordinates{X: 0, Y: math.MaxInt64},
			wantDeadNeighbors: 5,
		},
		{
			name:              "maximum y",
			center:            Coordinates{X: 0, Y: math.MaxInt64},
			wrappedCoordinate: Coordinates{X: 0, Y: math.MinInt64},
			wantDeadNeighbors: 5,
		},
		{
			name:              "minimum corner",
			center:            Coordinates{X: math.MinInt64, Y: math.MinInt64},
			wrappedCoordinate: Coordinates{X: math.MaxInt64, Y: math.MaxInt64},
			wantDeadNeighbors: 3,
		},
		{
			name:              "maximum corner",
			center:            Coordinates{X: math.MaxInt64, Y: math.MaxInt64},
			wrappedCoordinate: Coordinates{X: math.MinInt64, Y: math.MinInt64},
			wantDeadNeighbors: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := New([]Coordinates{tt.center, tt.wrappedCoordinate})
			got := state.getNeighbors(tt.center)

			if len(got.liveNeighbors) != 0 {
				t.Errorf("getNeighbors() found %d live neighbors across the boundary, want 0", len(got.liveNeighbors))
			}
			if len(got.deadNeighbors) != tt.wantDeadNeighbors {
				t.Errorf("getNeighbors() found %d dead neighbors, want %d", len(got.deadNeighbors), tt.wantDeadNeighbors)
			}
		})
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
	for cell, cellState := range state.cells {
		if cellState != alive {
			t.Errorf("cell %v has state %v after tick, want alive", cell, cellState)
		}
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
