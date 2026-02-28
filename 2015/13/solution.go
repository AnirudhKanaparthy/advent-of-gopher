package y2015d13

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

func (Solution) Solve(part int, args []string, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("No input file provided")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	text := string(data)

	rules := parse(text)
	switch part {
	case 1:
	case 2:
		addYourself(rules)
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	table := makeTable(rules)

	maxPerm := []string{}
	maxHappiness := -1
	for _, perm := range makePerms(table) {
		h := happiness(perm, rules)
		if h > maxHappiness {
			maxHappiness = h
			maxPerm = perm
		}
	}

	fmt.Fprintf(w, "Maximizing table order: %v -> %v\n", maxPerm, maxHappiness)
	return nil
}

func addYourself(rules map[string][]rule) {
	rules["You"] = make([]rule, 0)
}

func makePerms[T any](arr []T) [][]T {
	// TODO: Need to optimize this
	// As of now we aren't taking into
	// account he circlular nature of the table
	// [a b c] == [b c a]
	allPerms := make([][]T, 0)
	if len(arr) == 1 {
		return [][]T{arr}
	}

	for i, head := range arr {
		cpy := slices.Clone(arr)
		cpy = slices.Delete(cpy, i, i+1)

		res := makePerms(cpy)
		for j := range res {
			res[j] = slices.Insert(res[j], 0, head)
		}
		allPerms = append(allPerms, res...)
	}
	return allPerms
}

func getKeys[T comparable, E any](m map[T]E) []T {
	ks := make([]T, 0)
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func makeTable(rules map[string][]rule) []string {
	names := make(map[string]bool)
	for k := range rules {
		names[k] = true
	}
	for _, v := range rules {
		for _, r := range v {
			names[r.affectee] = true
		}
	}
	return getKeys(names)
}

func makeNeighbours(table []string) map[string][2]string {
	n := make(map[string][2]string)
	for i := range table {
		n[table[i]] = [2]string{
			table[(i+len(table)-1)%len(table)],
			table[(i+1)%len(table)],
		}
	}
	return n
}

func happiness(table []string, rules map[string][]rule) int {
	nb := makeNeighbours(table)

	score := 0
	for a, n := range nb {
		for _, r := range rules[a] {
			if r.affectee == n[0] || r.affectee == n[1] {
				score += r.amount
			}
		}
	}

	return score
}

type rule struct {
	affectee string
	amount   int
}

func parseRule(line string) (string, rule, error) {
	// Alice would lose 2 happiness units by sitting next to Bob.
	line = strings.Trim(line, " \n\t.")
	tokens := strings.Split(line, " ")
	if len(tokens) < 11 {
		return "", rule{}, errors.New("Invalid rule string")
	}
	amount, err := strconv.Atoi(tokens[3])
	if err != nil {
		return "", rule{}, err
	}

	var sign int
	switch tokens[2] {
	case "lose":
		sign = -1
	case "gain":
		sign = +1
	default:
		return "", rule{}, errors.New("Invalid rule string")
	}

	return tokens[0], rule{
		affectee: tokens[10],
		amount:   amount * sign,
	}, nil
}

func parse(source string) map[string][]rule {
	rules := make(map[string][]rule)
	for line := range strings.SplitSeq(source, "\n") {
		a, r, err := parseRule(line)
		if err != nil {
			continue
		}
		_, ok := rules[a]
		if !ok {
			rules[a] = make([]rule, 0)
		}
		rules[a] = append(rules[a], r)
	}
	return rules
}
