package y2015d15

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

	var evalFunc func(map[string]int, map[string]ingredient) int
	switch part {
	case 1:
		evalFunc = evaluateP1
	case 2:
		evalFunc = evaluateP2
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	const totalIngredients = 100

	ingredients := parseAllIngredients(text)

	items := make([]string, len(ingredients))
	for i := range ingredients {
		items[i] = ingredients[i].name
	}

	combs := createCombs(items, totalIngredients)

	mIngrs := make(map[string]ingredient)
	for _, igr := range ingredients {
		mIngrs[igr.name] = igr
	}

	maxScore := -1
	maxComb := map[string]int{}
	for _, c := range combs {
		score := evalFunc(c, mIngrs)
		if score > maxScore {
			maxScore = score
			maxComb = c
		}
	}

	fmt.Fprintln(w, maxScore, maxComb)

	return nil
}

func createCombs[T comparable](items []T, total int) []map[T]int {
	if len(items) == 1 {
		return []map[T]int{
			{items[0]: total},
		}
	}
	comb := make([]map[T]int, 0)
	for i := 0; i < total; i += 1 {
		c := createCombs(items[1:], total-i)
		for j := range c {
			c[j][items[0]] = i
		}
		comb = append(comb, c...)
	}
	return comb
}

func evaluateP1(comb map[string]int, allIngredients map[string]ingredient) int {
	cap := 0
	dur := 0
	fla := 0
	tex := 0

	for _, igr := range allIngredients {
		cap += igr.capacity * comb[igr.name]
		dur += igr.durability * comb[igr.name]
		fla += igr.flavor * comb[igr.name]
		tex += igr.texture * comb[igr.name]
	}

	cap = max(0, cap)
	dur = max(0, dur)
	fla = max(0, fla)
	tex = max(0, tex)

	return cap * dur * fla * tex
}

func evaluateP2(comb map[string]int, allIngredients map[string]ingredient) int {
	cap := 0
	dur := 0
	fla := 0
	tex := 0
	cal := 0

	for _, igr := range allIngredients {
		n := comb[igr.name]

		cap += n * igr.capacity
		dur += n * igr.durability
		fla += n * igr.flavor
		tex += n * igr.texture
		cal += n * igr.calories
	}

	cap = max(0, cap)
	dur = max(0, dur)
	fla = max(0, fla)
	tex = max(0, tex)

	if cal == 500 {
		return cap * dur * fla * tex
	} else {
		return 0
	}
}

type ingredient struct {
	name       string
	capacity   int
	durability int
	flavor     int
	texture    int
	calories   int
}

func (i ingredient) string() string {
	return fmt.Sprintf(
		`{
  "name"       : %v,
  "capacity"   : %v,
  "durability" : %v,
  "flavor"     : %v,
  "texture"    : %v,
  "calories"   : %v
}`, i.name, i.capacity, i.durability, i.flavor, i.texture, i.calories)
}

func parseIngredient(line string) (ingredient, error) {
	tokens := strings.Split(line, " ")
	name := strings.Trim(tokens[0], ": \n")

	capacity, err := strconv.Atoi(strings.Trim(tokens[2], ", "))
	if err != nil {
		return ingredient{}, nil
	}

	durability, err := strconv.Atoi(strings.Trim(tokens[4], ", "))
	if err != nil {
		return ingredient{}, nil
	}

	flavor, err := strconv.Atoi(strings.Trim(tokens[6], ", \n"))
	if err != nil {
		return ingredient{}, nil
	}

	texture, err := strconv.Atoi(strings.Trim(tokens[8], ", "))
	if err != nil {
		return ingredient{}, nil
	}

	calories, err := strconv.Atoi(strings.Trim(tokens[10], ", \n"))
	if err != nil {
		return ingredient{}, nil
	}

	return ingredient{
		name:       name,
		capacity:   capacity,
		durability: durability,
		flavor:     flavor,
		texture:    texture,
		calories:   calories,
	}, nil
}

func parseAllIngredients(source string) []ingredient {
	igrs := make([]ingredient, 0)
	for line := range strings.SplitSeq(source, "\n") {
		igr, err := parseIngredient(line)
		if err != nil {
			continue
		}
		igrs = append(igrs, igr)
	}
	return igrs
}
