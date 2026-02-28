package y2015d18

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Solution struct{}

func MakeSolution() *Solution {
	return &Solution{}
}

func (sol Solution) ArgsString(int, []string) string {
	return "<filepath>"
}

func (Solution) Solve(part int, args []string, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("No input file provided")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	text := string(data)

	const steps = 100

	grid, err := ParseLights(text)
	if err != nil {
		return err
	}

	var Stepper func(*Grid[bool])
	switch part {
	case 1:
		Stepper = StepLightsP1
	case 2:
		h := grid.height - 1
		w := grid.width - 1

		grid.Set(0, 0, true)
		grid.Set(w, 0, true)
		grid.Set(0, h, true)
		grid.Set(w, h, true)

		Stepper = StepLightsP2
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	before := CountOnLights(&grid)

	for i := 0; i < steps; i += 1 {
		Stepper(&grid)
	}

	after := CountOnLights(&grid)

	fmt.Fprintf(w, "Before: %v, After: %v\n", before, after)

	return nil
}

type Grid[T any] struct {
	data   []T
	width  int
	height int
}

func (g *Grid[T]) index(x int, y int) int {
	if x < 0 || x >= g.width {
		return -1
	}
	if y < 0 || y >= g.height {
		return -1
	}
	return y*g.width + x
}

func (g *Grid[T]) At(x int, y int) T {
	idx := g.index(x, y)
	return g.data[idx]
}

func (g *Grid[T]) Set(x int, y int, v T) {
	idx := g.index(x, y)
	g.data[idx] = v
}

func MakeGrid[T any](width int, height int) Grid[T] {
	return Grid[T]{
		data:   make([]T, width*height),
		width:  width,
		height: height,
	}
}

type Point struct {
	x int
	y int
}

var Neighbours = [8]Point{
	{-1, -1},
	{+0, -1},
	{+1, -1},

	{-1, +0},
	{+1, +0},

	{-1, +1},
	{+0, +1},
	{+1, +1},
}

func FprintLights(w io.Writer, g *Grid[bool]) {
	for j := 0; j < g.height; j += 1 {
		for i := 0; i < g.width; i += 1 {
			switch g.At(i, j) {
			case true:
				fmt.Fprintf(w, "#")
			case false:
				fmt.Fprintf(w, ".")
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func ParseLights(source string) (Grid[bool], error) {
	height := 0
	data := make([]bool, 0, len(source))
	for line := range strings.SplitSeq(source, "\n") {
		for col, c := range line {
			switch c {
			case '#':
				data = append(data, true)
			case '.':
				data = append(data, false)
			default:
				return Grid[bool]{}, fmt.Errorf("Error while parsing, (line %v, column %v)", height+1, col)
			}
		}
		height += 1
	}
	width := strings.Index(source, "\n")
	return Grid[bool]{
		data:   data,
		height: height,
		width:  width,
	}, nil
}

func LightsNeighbours(g *Grid[bool]) Grid[uint8] {
	nbs := MakeGrid[uint8](g.width, g.height)
	for j := 0; j < g.height; j += 1 {
		for i := 0; i < g.width; i += 1 {

			var val uint8 = 0
			for k := range Neighbours {
				at := Neighbours[k]
				x := i + at.x
				if x < 0 || x >= g.width {
					continue
				}
				y := j + at.y
				if y < 0 || y >= g.height {
					continue
				}

				if g.At(x, y) {
					val |= 1 << k
				}
			}
			nbs.Set(i, j, val)

		}
	}
	return nbs
}

func LightsCountNeighbours(g *Grid[bool]) Grid[int8] {
	nbs := MakeGrid[int8](g.width, g.height)
	for j := 0; j < g.height; j += 1 {
		for i := 0; i < g.width; i += 1 {

			var val int8 = 0
			for k := range Neighbours {
				at := Neighbours[k]
				x := i + at.x
				if x < 0 || x >= g.width {
					continue
				}
				y := j + at.y
				if y < 0 || y >= g.height {
					continue
				}

				if g.At(x, y) {
					val += 1
				}
			}
			nbs.Set(i, j, val)

		}
	}
	return nbs
}

func StepLightsP1(g *Grid[bool]) {
	nbs := LightsCountNeighbours(g)

	for j := 0; j < g.height; j += 1 {
		for i := 0; i < g.width; i += 1 {
			n := nbs.At(i, j)
			if g.At(i, j) {
				g.Set(i, j, (n == 2) || (n == 3))
			} else {
				g.Set(i, j, n == 3)
			}
		}
	}
}

func StepLightsP2(g *Grid[bool]) {
	nbs := LightsCountNeighbours(g)

	for j := 0; j < g.height; j += 1 {
		for i := 0; i < g.width; i += 1 {
			if (i == 0 || i == g.width-1) &&
				(j == 0 || j == g.height-1) {
				continue
			}
			n := nbs.At(i, j)
			if g.At(i, j) {
				g.Set(i, j, (n == 2) || (n == 3))
			} else {
				g.Set(i, j, n == 3)
			}
		}
	}
}

func CountOnLights(g *Grid[bool]) int {
	count := 0
	for j := 0; j < g.height; j += 1 {
		for i := 0; i < g.width; i += 1 {
			if g.At(i, j) {
				count += 1
			}
		}
	}
	return count
}
