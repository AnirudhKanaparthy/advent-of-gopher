package y2015d19

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"regexp"
	"slices"
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

	parts := strings.Split(text, "\n\n")
	rulesSrc := parts[0]

	rules, err := ParseRules(rulesSrc)
	if err != nil {
		return err
	}
	initial := strings.Trim(parts[1], " \t\n")

	switch part {
	case 1:
		unique := make(map[string]int)
		for _, rule := range rules {
			possible, err := rule.Apply(initial)
			if err != nil {
				return err
			}
			maps.Copy(unique, possible)
		}
		fmt.Fprintf(w, "Number of unique molecules: %v\n", len(unique))
	case 2:
		steps, err := StepsNeededV2(rules, initial)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Steps needed: %v\n", steps)
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	return nil
}

type Rule struct {
	from string
	to   string
}

func ParseRules(source string) ([]Rule, error) {
	rules := make([]Rule, 0)
	i := 0
	for line := range strings.SplitSeq(source, "\n") {
		tokens := strings.Split(line, "=>")
		if len(tokens) < 2 {
			return rules, fmt.Errorf("Invalid rule at %v", i)
		}

		from := strings.Trim(tokens[0], " \t\n")
		to := strings.Trim(tokens[1], " \t\n")

		rules = append(rules, Rule{from: from, to: to})
	}
	return rules, nil
}

func (r Rule) Apply(molecule string) (map[string]int, error) {
	exp, err := regexp.Compile(r.from)
	if err != nil {
		return nil, err
	}

	possible := make(map[string]int, 0)
	matches := exp.FindAllStringIndex(molecule, -1)

	for _, m := range matches {
		var sb strings.Builder
		sb.WriteString(molecule[:m[0]])
		sb.WriteString(r.to)
		sb.WriteString(molecule[m[1]:])

		possible[sb.String()] += 1
	}

	return possible, nil
}

func PrintfJourney(journey []string) {
	for _, p := range journey[:len(journey)-1] {
		fmt.Printf("%v => ", p)
	}
	fmt.Printf("%v", journey[len(journey)-1])
}

func StepsNeeded(rules []Rule, target string) (int, error) {
	shortest := math.Inf(+1)

	count := 0
	var helper func([]string, string, float64) (int, error)
	helper = func(journey []string, current string, path float64) (int, error) {
		if len(current) > len(target) {
			count += 1
			PrintfJourney(journey)
			fmt.Printf("[len]\n")
			return 0, nil
		}
		if current == target {
			shortest = min(shortest, path)
			count += 1

			PrintfJourney(journey)
			fmt.Printf("[short](%v)\n", shortest)
			return 1, nil
		}

		unique := make(map[string]int)
		for _, r := range rules {
			possible, err := r.Apply(current)
			if err != nil {
				return -1, err
			}
			maps.Copy(unique, possible)
		}
		totalCount := 0
		for p := range unique {
			cpy := slices.Clone(journey)
			cpy = append(cpy, p)
			c, err := helper(cpy, p, path+1)
			if err != nil {
				return -1, err
			}
			totalCount += c
		}

		PrintfJourney(journey)
		fmt.Printf("[store]\n")
		return totalCount, nil
	}

	total, err := helper([]string{"e"}, "e", 0.0)
	if err != nil {
		return -1, err
	}
	_ = total
	return int(shortest), nil
}

func StepsNeededV2(rules []Rule, target string) (int, error) {
	cache := make(map[string]float64)
	cacheHits := 0
	moleculesChecked := 0

	var helper func(string) (float64, error)
	helper = func(molecule string) (float64, error) {
		moleculesChecked += 1
		if moleculesChecked%100000 == 0 {
			fmt.Printf("Info:\n\tMolecules checked: %v\n\t       Cache hits: %v\n\t Cache hits ratio: %.3f\n", moleculesChecked, cacheHits, float64(cacheHits)/float64(moleculesChecked))
		}

		if len(molecule) > len(target) {
			return math.Inf(+1), nil
		}
		if molecule == target {
			return 0.0, nil
		}

		c, ok := cache[molecule]
		if ok {
			cacheHits += 1
			return c, nil
		}

		uniqueMolecules := make(map[string]int)
		for _, r := range rules {
			newMolecules, err := r.Apply(molecule)
			if err != nil {
				return -1.0, err
			}
			maps.Copy(uniqueMolecules, newMolecules)
		}

		smallestCost := math.Inf(+1)
		// Iterate through all the possible molecules and DFS
		for m := range uniqueMolecules {
			cost, err := helper(m)
			if err != nil {
				return -1.0, err
			}
			smallestCost = min(smallestCost, cost+1.0)
		}

		cache[molecule] = smallestCost
		return smallestCost, nil
	}

	cost, err := helper("e")
	if err != nil {
		return -1, err
	}
	fmt.Printf("Info:\n\tMolecules checked: %v\n\t       Cache hits: %v\n\t Cache hits ratio: %.3f\n", moleculesChecked, cacheHits, float64(cacheHits)/float64(moleculesChecked))
	return int(cost), nil
}
