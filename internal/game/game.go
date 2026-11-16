package game

import (
	"math"
)

type Coordinates struct {
	X int64
	Y int64
}

type State struct {
	// TODO: first pass is holding state in map. Consider the space complexity of holding 2 game states instead.
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
	nextCells := make(map[Coordinates]struct{})
	for c := range s.cells {
		neighbors := s.getNeighbors(c)

		// If an "alive" cell had less than 2 or more than 3 alive neighbors (in any of the 8 surrounding cells), it becomes dead.
		if len(neighbors.liveNeighbors) < 2 || len(neighbors.liveNeighbors) > 3 {
			continue
		}
		nextCells[c] = struct{}{}

		// If a "dead" cell had *exactly* 3 alive neighbors, it becomes alive.
		// TODO: find better efficiencies for this lookup
		for _, n := range neighbors.deadNeighbors {
			neighbors := s.getNeighbors(n)
			if len(neighbors.liveNeighbors) == 3 {
				nextCells[n] = struct{}{}
			}
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
