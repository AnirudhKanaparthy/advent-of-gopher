package y2015d12

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
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
		text, err = pruneRed(text)
		if err != nil {
			return err
		}
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	re, err := regexp.Compile(`[-+]?\d+(\.\d+)?(e\d+)?`)
	if err != nil {
		return err
	}
	sum := 0
	for _, m := range re.FindAll([]byte(text), -1) {
		v, err := strconv.Atoi(string(m))
		if err != nil {
			return err
		}
		sum += v
	}

	fmt.Fprintf(w, "Sum: %v\n", sum)

	return nil
}

func pruneRed(text string) (string, error) {
	se, err := regexp.Compile(`:\s*"red"`)
	if err != nil {
		return "", err
	}

	for {
		m := se.FindIndex([]byte(text))
		if m == nil {
			break
		}

		i := m[0]
		s := 1
		for i >= 0 {
			switch text[i] {
			case '}':
				s += 1
			case '{':
				s -= 1
			}
			if s == 0 {
				break
			}
			i -= 1
		}

		start := i
		s = 0
		for i < len(text) {
			switch text[i] {
			case '{':
				s += 1
			case '}':
				s -= 1
			}
			i += 1
			if s == 0 {
				break
			}
		}
		text = text[:start] + text[i:]
	}
	return text, nil
}
