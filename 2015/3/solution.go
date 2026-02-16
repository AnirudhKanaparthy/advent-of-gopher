package y2015d3

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
		return errors.New("No input file provided")
	}

	inputFilePath := args[0]
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}
	text := string(data)

	var result string
	switch part {
	case 1:
		result = fmt.Sprintf("Houses that receive at least one gift: %v", giftDelivery(text))
	case 2:
		result = fmt.Sprintf("Houses that receive at least one gift with robot: %v", giftDeliveryWithRobot(text))
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	fmt.Fprintf(w, "%s", result)
	return nil
}

type vec2i struct {
	x int
	y int
}

func giftDelivery(instructions string) int {
	visited := make(map[vec2i]bool)
	pos := vec2i{0, 0}
	visited[pos] = true
	for _, c := range instructions {
		switch c {
		case '^':
			pos.y += 1
		case 'v':
			pos.y -= 1
		case '<':
			pos.x -= 1
		case '>':
			pos.x += 1
		default:
			continue
		}
		visited[pos] = true
	}
	return len(visited)
}

func giftDeliveryWithRobot(instructions string) int {
	visited := make(map[vec2i]bool)
	santa := vec2i{0, 0}
	roboSanta := vec2i{0, 0}
	visited[santa] = true
	for i, c := range instructions {
		ptr := &santa
		if i%2 == 0 {
			ptr = &roboSanta
		}

		switch c {
		case '^':
			ptr.y += 1
		case 'v':
			ptr.y -= 1
		case '<':
			ptr.x -= 1
		case '>':
			ptr.x += 1
		default:
			continue
		}
		visited[*ptr] = true
	}
	return len(visited)
}
