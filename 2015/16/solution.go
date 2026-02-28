package y2015d16

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
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

	var matcher func(aunt, map[string]int) int
	switch part {
	case 1:
		matcher = matchAuntP1
	case 2:
		matcher = matchAuntP2
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	details := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	aunts := parseAllAunts(text)

	type auntScore struct {
		id    int
		score int
	}

	maxAunts := []auntScore{}
	scores := matchAllAunts(aunts, details, matcher)
	for i, s := range scores {
		if len(maxAunts) == 0 {
			maxAunts = append(maxAunts, auntScore{aunts[i].id, s})
		}

		if s < maxAunts[0].score {
			continue
		}

		if s > maxAunts[0].score {
			maxAunts = maxAunts[:0]
		}
		maxAunts = append(maxAunts, auntScore{aunts[i].id, s})
	}

	for _, s := range maxAunts {
		fmt.Fprintf(w, "It's Aunt Sue %v -> %v\n", s.id, s.score)
	}
	return nil
}

type aunt struct {
	id      int
	details map[string]int
}

func matchAuntP1(a aunt, details map[string]int) int {
	count := 0
	for k, v := range a.details {
		if details[k] == v {
			count += 1
		}
	}
	return count
}

func matchAuntP2(a aunt, details map[string]int) int {
	count := 0

	greaterThan := map[string]bool{
		"cats":  true,
		"trees": true,
	}
	fewerThan := map[string]bool{
		"pomeranians": true,
		"goldfish":    true,
	}

	for k, v := range a.details {
		if greaterThan[k] {
			if details[k] < v {
				count += 1
			}
			continue
		}

		if fewerThan[k] {
			if details[k] > v {
				count += 1
			}
			continue
		}

		if details[k] == v {
			count += 1
		}
	}
	return count
}

func matchAllAunts(aunts []aunt, details map[string]int, matcher func(aunt, map[string]int) int) []int {
	scores := make([]int, 0)
	for _, a := range aunts {
		scores = append(scores, matcher(a, details))
	}
	return scores
}

func parseAunt(line string) (aunt, error) {
	i := strings.Index(line, ":")
	if i == -1 {
		return aunt{}, errors.New("Invalid aunt description")
	}
	nextLine, found := strings.CutPrefix(line[:i], "Sue ")
	if !found {
		return aunt{}, errors.New("Invalid aunt description")
	}
	id, err := strconv.Atoi(nextLine)
	if err != nil {
		return aunt{}, err
	}
	_ = id

	keyValsExp, err := regexp.Compile(`([a-z]+)\s*:\s*(\d+)`)
	if err != nil {
		return aunt{}, err
	}

	items := make(map[string]int)
	for _, m := range keyValsExp.FindAllStringSubmatch(line[i:], -1) {
		items[m[1]], err = strconv.Atoi(m[2])
		if err != nil {
			return aunt{}, err
		}
	}
	return aunt{id: id, details: items}, nil
}

func parseAllAunts(text string) []aunt {
	aunts := make([]aunt, 0)
	for line := range strings.SplitSeq(text, "\n") {
		a, err := parseAunt(line)
		if err != nil {
			continue
		}
		aunts = append(aunts, a)
	}
	return aunts
}
