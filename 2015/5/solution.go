package y2015d5

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
	return "<file path>"
}

func (sol Solution) Solve(part int, args []string, sb io.Writer) error {
	if len(args) < 1 {
		return errors.New("No input file provided")
	}

	inputFilePath := args[0]
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}
	text := string(data)

	var callback func(string) bool

	switch part {
	case 1:
		callback = isNiceStringPart1
	case 2:
		callback = isNiceStringPart2
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	nice := 0
	for line := range strings.SplitSeq(text, "\n") {
		if callback(line) {
			nice += 1
		}
	}

	fmt.Fprintf(sb, "Number of nice strings (%v): %v", part, nice)
	return nil
}

func isNiceStringPart1(text string) bool {
	l := len(text)
	if l <= 3 {
		return false
	}

	twiceChars := 0
	for i := range text[:l-1] {
		s := text[i : i+2]
		switch s {
		case "ab":
			fallthrough
		case "cd":
			fallthrough
		case "pq":
			fallthrough
		case "xy":
			return false
		}
		if text[i] == text[i+1] {
			twiceChars += 1
		}
	}
	if twiceChars < 1 {
		return false
	}

	vowelsCount := 0
	for _, c := range text {
		switch c {
		case 'a':
			fallthrough
		case 'e':
			fallthrough
		case 'i':
			fallthrough
		case 'o':
			fallthrough
		case 'u':
			vowelsCount += 1
		}
	}

	return vowelsCount >= 3
}

func checkCondition1(text string) bool {
	if len(text) < 4 {
		return false
	}

	pairMap := make(map[string]int)
	lastPair := ""
	for i := range text[:len(text)-1] {
		pair := text[i : i+2]

		flag := (lastPair == pair)
		lastPair = pair

		if flag {
			continue
		}

		val, ok := pairMap[pair]
		if !ok {
			pairMap[pair] = 0
			val = 0
		}

		pairMap[pair] = val + 1
		if pairMap[pair] >= 2 {
			return true
		}
	}
	return false
}

func checkCondition2(text string) bool {
	if len(text) < 3 {
		return false
	}

	lastPair := text[:2]
	for i := range text[1 : len(text)-1] {
		pair := text[i : i+2]

		if lastPair[0] == pair[1] && lastPair[1] == pair[0] {
			return true
		}

		lastPair = pair
	}
	return false
}

func isNiceStringPart2(text string) bool {
	return checkCondition1(text) && checkCondition2(text)
}
