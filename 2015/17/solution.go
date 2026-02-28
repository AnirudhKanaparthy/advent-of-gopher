package y2015d17

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
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

	switch part {
	case 1:
	case 2:
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	containers := parseContainerSizes(text)
	// containers := []int{20, 15, 10, 5, 5}
	combs, ways := countCombinations(containers, 150)
	fmt.Fprintf(w, "All possible combinations: %v\n", combs)
	fmt.Fprintf(w, "Ways: %v\n", ways)

	return nil
}

func parseContainerSizes(text string) []int {
	containers := make([]int, 0)
	for line := range strings.FieldsSeq(text) {
		val, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		containers = append(containers, val)
	}
	return containers
}

func countCombinations(containers []int, total int) (int, map[int]int) {
	ways := make(map[int]int)

	var helper func([]int, int, int) int
	helper = func(containers []int, total int, using int) int {
		if total < 0 || len(containers) == 0 {
			return 0
		}
		if total == 0 {
			ways[using] += 1
			return 1
		}
		if len(containers) == 1 {
			if containers[0] == total {
				ways[using+1] += 1
				return 1
			}
			return 0
		}

		return helper(containers[1:], total-containers[0], using+1) +
			helper(containers[1:], total, using)
	}
	return helper(containers, total, 0), ways
}
