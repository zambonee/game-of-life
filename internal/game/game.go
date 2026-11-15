package game

import "math"

type Coordinates struct {
	X int64
	Y int64
}

type cellState uint8

const (
	alive    cellState = iota // alive now and alive next generation (default state for new game state)
	dead                      // alive now but dead next generation
	emergent                  // dead now (does not exist in grid as live) but alive next generation
)

type State struct {
	// TODO: first pass is holding state in map. Consider the space complexity of holding 2 game states instead.
	Cells map[Coordinates]cellState
}

func New(c []Coordinates) *State {
	state := &State{
		Cells: make(map[Coordinates]cellState, len(c)),
	}
	for _, v := range c {
		state.Cells[v] = alive
	}
	return state
}

func (s *State) TickForward() {
	for c := range s.Cells {
		neighbors := s.getNeighbors(c)

		// If an "alive" cell had less than 2 or more than 3 alive neighbors (in any of the 8 surrounding cells), it becomes dead.
		if len(neighbors.liveNeighbors) < 2 || len(neighbors.liveNeighbors) > 3 {
			s.Cells[c] = dead
		}

		// If a "dead" cell had *exactly* 3 alive neighbors, it becomes alive.
		// TODO: find better efficiencies for this lookup
		for _, n := range neighbors.deadNeighbors {
			neighbors := s.getNeighbors(n)
			if len(neighbors.liveNeighbors) == 3 {
				s.Cells[n] = emergent
			}
		}
	}
	s.cleanState()
}

func (s *State) GetState() []Coordinates {
	result := make([]Coordinates, 0, len(s.Cells))
	for c := range s.Cells {
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
			cellState, exists := s.Cells[neighbor]
			if exists && cellState != emergent {
				result.liveNeighbors = append(result.liveNeighbors, neighbor)
				continue
			}
			result.deadNeighbors = append(result.deadNeighbors, neighbor)
		}
	}
	return result
}

func offsetCoordinate(value, offset int64) (int64, bool) {
	if offset == -1 && value == math.MinInt64 || offset == 1 && value == math.MaxInt64 {
		return 0, false
	}
	return value + offset, true
}

func (s *State) cleanState() {
	for cell, state := range s.Cells {
		switch state {
		case dead:
			delete(s.Cells, cell)
		case emergent:
			s.Cells[cell] = alive
		}
	}
}
