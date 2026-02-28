package y2015d9

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
	return "<file path>"
}

func (Solution) Solve(part int, args []string, w io.Writer) error {
	filepath := args[0]
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	text := string(data)

	g, err := makeGraph(text)
	if err != nil {
		return fmt.Errorf("Cannot make graph, error while parsing: %v", err)
	}

	var result graphPath
	switch part {
	case 1:
		fmt.Fprintf(w, "Shortest path:\n")
		result, err = g.shortestPath([]string{})
		if err != nil {
			return err
		}
	case 2:
		fmt.Fprintf(w, "Longest path:\n")
		result, err = g.longestPath([]string{})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	for _, p := range result.path[:len(result.path)-1] {
		fmt.Fprintf(w, "%v -> ", p)
	}
	fmt.Fprintf(w, "%v : %v\n", result.path[len(result.path)-1], result.cost)
	return nil
}

type graph map[string]map[string]int

func (g graph) destinations(src string) map[string]int {
	vM, ok := g[src]
	if !ok {
		vM = make(map[string]int)
	}

	for d, tM := range g {
		v, ok := tM[src]
		if ok {
			vM[d] = v
		}
	}
	return vM
}

func identity[T any](a T, b T) (T, T) {
	// 4 copies, 1 extra copy (if you had used swap)
	return a, b
}

func (g graph) set(a string, b string, val int) {
	if strings.Compare(a, b) > 0 {
		b, a = identity(a, b)
	}

	dests, ok := g[a]
	if !ok {
		g[a] = make(map[string]int)
		dests = g[a]
	}

	dests[b] = val
}

func makeGraph(text string) (graph, error) {
	g := make(graph)

	for line := range strings.SplitSeq(text, "\n") {
		tokens := strings.Fields(line)
		if len(tokens) < 5 {
			return graph{}, errors.New("Invalid distance description")
		}

		a := tokens[0]
		b := tokens[2]
		vt := tokens[4]

		v, err := strconv.Atoi(vt)
		if err != nil {
			return graph{}, errors.New("Invalid distance value")
		}

		g.set(a, b, v)
	}
	return g, nil
}

type graphPath struct {
	path []string
	cost int
}

func (g graph) nodeCount() int {
	return len(g) + 1
}

func (g graph) shortestPath(visited []string) (graphPath, error) {
	if g.nodeCount() == len(visited) {
		return graphPath{path: []string{}, cost: 0}, nil
	}

	nextLoc := ""
	shortest := graphPath{path: nil, cost: 99999999}

	if len(visited) == 0 {
		for start := range g {
			p, err := g.shortestPath([]string{start})
			if err != nil {
				return graphPath{}, err
			}

			if p.cost < shortest.cost {
				nextLoc = start
				shortest = p
			}
		}

		t := append([]string{}, nextLoc)
		t = append(t, shortest.path...)
		shortest.path = t

		return shortest, nil
	}

	start := visited[len(visited)-1]
	for dest, c := range g.destinations(start) {
		if slices.Index(visited, dest) != -1 {
			continue
		}

		visited = append(visited, dest) // Visit

		p, err := g.shortestPath(visited)
		if err != nil {
			return graphPath{}, err
		}

		visited = visited[:len(visited)-1] // Undo

		p.cost += c
		if p.cost < shortest.cost {
			nextLoc = dest
			shortest = p
		}
	}

	t := append([]string{}, nextLoc)
	t = append(t, shortest.path...)
	shortest.path = t

	return shortest, nil
}

func (g graph) longestPath(visited []string) (graphPath, error) {
	if g.nodeCount() == len(visited) {
		return graphPath{path: []string{}, cost: 0}, nil
	}

	nextLoc := ""
	longest := graphPath{path: nil, cost: -1}

	if len(visited) == 0 {
		for start := range g {
			p, err := g.longestPath([]string{start})
			if err != nil {
				return graphPath{}, err
			}

			if p.cost > longest.cost {
				nextLoc = start
				longest = p
			}
		}

		t := append([]string{}, nextLoc)
		t = append(t, longest.path...)
		longest.path = t

		return longest, nil
	}

	start := visited[len(visited)-1]
	for dest, c := range g.destinations(start) {
		if slices.Index(visited, dest) != -1 {
			continue
		}

		visited = append(visited, dest) // Visit

		p, err := g.longestPath(visited)
		if err != nil {
			return graphPath{}, err
		}

		visited = visited[:len(visited)-1] // Undo

		p.cost += c
		if p.cost > longest.cost {
			nextLoc = dest
			longest = p
		}
	}

	t := append([]string{}, nextLoc)
	t = append(t, longest.path...)
	longest.path = t

	return longest, nil
}
