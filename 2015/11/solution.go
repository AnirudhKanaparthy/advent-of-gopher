package y2015d11

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Solution struct{}

func MakeSolution() *Solution {
	return &Solution{}
}

func (sol Solution) ArgsString(int, []string) string {
	return "<sequence>"
}

func (Solution) Solve(part int, args []string, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("No sequence provided")
	}
	seq := args[0]
	if len(seq) != 8 {
		return errors.New("Sequence provided must be of length 8")
	}

	toCheck := []stringChecker{
		checkNoiol,
		checkPairs2,
		checkInc3,
	}

	for {
		seq = incString(seq)
		if checkAll(seq, toCheck) {
			fmt.Fprintf(w, "Valid sequence: %v", seq)
			break
		}
	}

	return nil
}

func incChar(b byte) byte {
	return (((b - 'a') + 1) % ('z' - 'a' + 1)) + 'a'
}

func incString(text string) string {
	// Probably not efficient
	if len(text) == 0 {
		return ""
	}

	last := len(text) - 1
	front := text[:last]
	if text[last] == 'z' {
		front = incString(front)
	}
	back := incChar(text[last])
	return fmt.Sprintf("%s%c", front, back)
}

type stringChecker func(string) bool

func checkInc3(text string) bool {
	if len(text) < 3 {
		return false
	}
	if text[0]+1 == text[1] && text[1]+1 == text[2] {
		return true
	}
	return checkInc3(text[1:])
}

func checkNoiol(text string) bool {
	return !strings.ContainsAny(text, "iol")
}

func checkPairs2(text string) bool {
	if len(text) < 2 {
		return false
	}
	count := 0
	i := 0
	for i < len(text)-1 {
		if text[i] == text[i+1] {
			count += 1
			i += 1
		}
		i += 1

		if count >= 2 {
			return true
		}
	}
	return false
}

func checkAll(text string, checkers []stringChecker) bool {
	for _, check := range checkers {
		if !check(text) {
			return false
		}
	}
	return true
}
