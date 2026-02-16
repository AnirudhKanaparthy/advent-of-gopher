package y2015d2

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
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

func (sol Solution) Solve(part int, args []string, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("No file path provided")
	}

	inputFilePath := args[0]

	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}
	text := string(data)

	totalRibbonNeeded := 0
	totalWrapNeeded := 0
	for line := range strings.SplitSeq(text, "\n") {
		if len(line) == 0 {
			continue
		}

		box, err := parseIntoBox(line)
		if err != nil {
			return err
		}
		totalWrapNeeded += box.giftWrapNeeded()
		totalRibbonNeeded += box.ribbonNeeded()
	}

	switch part {
	case 1:
		fmt.Fprintf(w, "Total wrap needed: %v", totalWrapNeeded)
	case 2:
		fmt.Fprintf(w, "Total ribbon needed: %v", totalRibbonNeeded)
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	return nil
}

type box struct {
	l int
	w int
	h int
}

func (b *box) giftWrapNeeded() int {
	p := b.l * b.w
	q := b.w * b.h
	r := b.h * b.l

	s := min(min(p, q), r)

	return 2*p + 2*q + 2*r + s
}

func (b *box) ribbonNeededV1() int {
	s := []int{b.l, b.h, b.w}
	slices.Sort(s)
	return 2*(s[0]+s[1]) + (b.l * b.h * b.w)
}

func (b *box) ribbonNeeded() int {
	// A more fun approach

	p := min(b.l, b.h)
	q := min(b.h, b.w)
	r := min(b.w, b.l)

	m := max(p, q)
	n := max(q, r)
	o := max(r, p)

	t := (p + q + r + m + n + o) / 3
	return 2*t + (b.l * b.w * b.h)
}

func parseIntoBox(text string) (box, error) {
	if len(text) == 0 {
		return box{}, errors.New("empty input")
	}

	dimsRaw := strings.Split(text, "x")

	dims := [3]int{0, 0, 0}
	for i, dimRaw := range dimsRaw {
		dim, err := strconv.Atoi(dimRaw)
		if err != nil {
			return box{}, errors.New("wrong box format")
		}
		dims[i] = dim
	}

	return box{dims[0], dims[1], dims[2]}, nil
}
