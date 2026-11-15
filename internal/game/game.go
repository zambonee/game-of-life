package game

type coordinates struct {
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
	Cells map[coordinates]cellState
}

func New(c []coordinates) *State {
	state := &State{
		Cells: make(map[coordinates]cellState, len(c)),
	}
	for _, v := range c {
		state.Cells[v] = alive
	}
	return state
}

func (s *State) TickForward() {
	// TODO: consider behavior on int64 overflow

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

// TODO: return something consumable for the caller
func (s *State) GetState() map[coordinates]cellState {
	return s.Cells
}

type getNeighborsResponse struct {
	deadNeighbors []coordinates
	liveNeighbors []coordinates
}

func (s *State) getNeighbors(c coordinates) getNeighborsResponse {
	result := getNeighborsResponse{}
	// range through rather than explicitly repeat coordinates for readability
	for y := c.Y - 1; y <= c.Y+1; y++ {
		for x := c.X - 1; x <= c.X+1; x++ {
			if y == c.Y && x == c.X {
				continue
			}
			cellState, exists := s.Cells[coordinates{Y: y, X: x}]
			if exists && cellState != emergent {
				result.liveNeighbors = append(result.liveNeighbors, coordinates{Y: y, X: x})
				continue
			}
			result.deadNeighbors = append(result.deadNeighbors, coordinates{Y: y, X: x})
		}
	}
	return result
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
