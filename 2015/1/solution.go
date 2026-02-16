package y2015d1

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Solution struct{}

func MakeSolution() *Solution {
	return &Solution{}
}

func (sol Solution) ArgsString(int, []string) string {
	return "<filepath>"
}

func (sol Solution) Solve(part int, args []string, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("Need at least one argument")
	}

	inputFilePath := args[0]

	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}
	contents := string(data)

	switch part {
	case 1:
		fmt.Fprintf(w, "Floor position: %v", floorPos(contents))
	case 2:
		fmt.Fprintf(w, "Basement position: %v", basementPos(contents))
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}
	return nil
}

func basementPos(instructions string) int {
	floor := 0
	for i, c := range instructions {
		switch c {
		case '(':
			floor += 1
		case ')':
			floor -= 1
		}

		if floor == -1 {
			return i + 1
		}
	}
	return floor
}

func floorPos(instructions string) int {
	floor := 0
	for _, c := range instructions {
		switch c {
		case '(':
			floor += 1
		case ')':
			floor -= 1
		}
	}
	return floor
}
