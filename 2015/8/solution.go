package y2015d8

import (
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
	return "<file path>"
}

func (Solution) Solve(part int, args []string, w io.Writer) error {
	filepath := args[0]
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	text := string(data)

	switch part {
	case 1:
		solveP1(text, w)
	case 2:
		solveP2(text, w)
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	return nil
}

func solveP1(text string, w io.Writer) {
	ncode := 0
	nmem := 0
	for line := range strings.FieldsSeq(text) {
		ncode += len(line)
		nmem += countInMemorySize(line)
	}

	fmt.Fprintf(w, "%v - %v = %v\n", ncode, nmem, ncode-nmem)
}

func solveP2(text string, w io.Writer) {
	ncode := 0
	nesc := 0
	for line := range strings.FieldsSeq(text) {
		ncode += len(line)
		nesc += countEscapedSize(line)
	}

	fmt.Fprintf(w, "%v - %v = %v\n", nesc, ncode, nesc-ncode)
}

func countEscapedSize(text string) int {
	count := 1

	i := 0
	for i < len(text) {
		c := text[i]

		switch c {
		case '"':
			count += 2
		case '\\':
			count += 2
		default:
			count += 1
		}

		i += 1
	}

	count += 1
	return count
}

func countInMemorySize(text string) int {
	count := 0

	i := 0
	for i < len(text) {
		c := text[i]

		switch c {
		case '"':
		case '\\':
			i += 1
			if text[i] == 'x' {
				i += 2
			}

			count += 1
		default:
			count += 1
		}

		i += 1
	}

	return count
}
