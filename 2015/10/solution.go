package y2015d10

import (
	"fmt"
	"io"
	"strings"
)

type Solution struct{}

func MakeSolution() *Solution {
	return &Solution{}
}

func (sol Solution) ArgsString(int, []string) string {
	return "<first sequence>"
}

func (Solution) Solve(part int, args []string, w io.Writer) error {
	var times int
	switch part {
	case 1:
		times = 40
	case 2:
		times = 50
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	cur := args[0]
	for i := 0; i < times; i += 1 {
		cur = nextSeq(cur)
		fmt.Fprintf(w, "%2v: %v\n", i+1, len(cur))
	}

	fmt.Fprintf(w, "Answer: %v\n", len(cur))
	return nil
}

func nextSeq(text string) string {
	var sb strings.Builder
	count := 1
	for i := 1; i < len(text); i += 1 {
		if text[i] != text[i-1] {
			fmt.Fprintf(&sb, "%v%c", count, text[i-1])
			count = 0
		}
		count += 1
	}
	fmt.Fprintf(&sb, "%v%c", count, text[len(text)-1])

	return sb.String()
}
