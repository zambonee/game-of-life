package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/zambonee/game-of-life/internal/game"
)

const generationsToRun = 10

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(input io.Reader, output io.Writer) error {
	cells, err := parseLife106(input)
	if err != nil {
		return err
	}

	state := game.New(cells)
	for range generationsToRun {
		state.TickForward()
	}

	cells = state.GetState()
	writer := bufio.NewWriter(output)

	if err := write(writer, cells); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush output: %w", err)
	}

	return nil
}

func parseLife106(input io.Reader) ([]game.Coordinates, error) {
	scanner := bufio.NewScanner(input)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("read input: %w", err)
		}
		return nil, fmt.Errorf("missing Life 1.06 header")
	}
	if strings.TrimSpace(scanner.Text()) != "#Life 1.06" {
		return nil, fmt.Errorf("line 1: expected #Life 1.06 header")
	}

	var cells []game.Coordinates
	lineNumber := 1
	for scanner.Scan() {
		lineNumber++
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected two coordinates", lineNumber)
		}

		x, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid x coordinate %q: %w", lineNumber, fields[0], err)
		}
		y, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid y coordinate %q: %w", lineNumber, fields[1], err)
		}
		cells = append(cells, game.Coordinates{X: x, Y: y})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}

	return cells, nil
}

func write(writer io.Writer, cells []game.Coordinates) error {
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].X == cells[j].X {
			return cells[i].Y < cells[j].Y
		}
		return cells[i].X < cells[j].X
	})

	if _, err := fmt.Fprintln(writer, "#Life 1.06"); err != nil {
		return fmt.Errorf("write Life 1.06 header: %w", err)
	}
	for _, cell := range cells {
		if _, err := fmt.Fprintf(writer, "%d %d\n", cell.X, cell.Y); err != nil {
			return fmt.Errorf("write cell %d %d: %w", cell.X, cell.Y, err)
		}
	}
	return nil
}
