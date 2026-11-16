package game

import (
	"math"
)

type Coordinates struct {
	X int64
	Y int64
}

type State struct {
	cells map[Coordinates]struct{}
}

func New(c []Coordinates) *State {
	state := &State{
		cells: make(map[Coordinates]struct{}, len(c)),
	}
	for _, v := range c {
		state.cells[v] = struct{}{}
	}
	return state
}

// TickForward progresses the game state to the next generation
func (s *State) TickForward() {
	neighborCounts := make(map[Coordinates]int, len(s.cells))
	for cell := range s.cells {
		for _, yOffset := range []int64{-1, 0, 1} {
			y, ok := offsetCoordinate(cell.Y, yOffset)
			if !ok {
				continue
			}
			for _, xOffset := range []int64{-1, 0, 1} {
				if xOffset == 0 && yOffset == 0 {
					continue
				}
				x, ok := offsetCoordinate(cell.X, xOffset)
				if !ok {
					continue
				}
				neighborCounts[Coordinates{X: x, Y: y}]++
			}
		}
	}

	nextCells := make(map[Coordinates]struct{})
	for cell, count := range neighborCounts {
		_, alive := s.cells[cell]
		if count == 3 || (alive && count == 2) {
			nextCells[cell] = struct{}{}
		}
	}

	s.cells = nextCells
}

// GetState returns the coordinates of live cells at the current game state
func (s *State) GetState() []Coordinates {
	result := make([]Coordinates, 0, len(s.cells))
	for c := range s.cells {
		result = append(result, c)
	}
	return result
}

type getNeighborsResponse struct {
	deadNeighbors []Coordinates
	liveNeighbors []Coordinates
}

func (s *State) getNeighbors(c Coordinates) getNeighborsResponse {
	result := getNeighborsResponse{}
	for _, yOffset := range []int64{-1, 0, 1} {
		y, ok := offsetCoordinate(c.Y, yOffset)
		if !ok {
			continue
		}
		for _, xOffset := range []int64{-1, 0, 1} {
			if yOffset == 0 && xOffset == 0 {
				continue
			}
			x, ok := offsetCoordinate(c.X, xOffset)
			if !ok {
				continue
			}
			neighbor := Coordinates{X: x, Y: y}
			if _, exists := s.cells[neighbor]; exists {
				result.liveNeighbors = append(result.liveNeighbors, neighbor)
				continue
			}
			result.deadNeighbors = append(result.deadNeighbors, neighbor)
		}
	}
	return result
}

func offsetCoordinate(value, offset int64) (int64, bool) {
	if (offset == -1 && value == math.MinInt64) || (offset == 1 && value == math.MaxInt64) {
		return 0, false
	}
	return value + offset, true
}
