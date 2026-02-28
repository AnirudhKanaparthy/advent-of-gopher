package y2015d14

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

	const totalIterations = 2503

	deers := parseAllReindeerStats(text)
	state := makeRaceState(deers)

	for t := 0; t < totalIterations; t += 1 {
		state.Step()
	}

	switch part {
	case 1:
		mDeer := ""
		mDistance := -1
		for n, d := range state.distances {
			if d > mDistance {
				mDistance = d
				mDeer = n
			}
		}
		fmt.Fprintf(w, "Max Distance: %v -> %v\n", mDeer, mDistance)
	case 2:
		mDeer := ""
		mPoints := -1
		for n, p := range state.points {
			if p > mPoints {
				mPoints = p
				mDeer = n
			}
		}
		fmt.Fprintf(w, "Max Points: %v -> %v\n", mDeer, mPoints)
	default:
		return errors.New("Expected `part` to be either 1 or 2")
	}

	return nil
}

type raceState struct {
	time       int
	stats      map[string]reindeerStat
	flying     map[string]bool
	distances  map[string]int
	nextChange map[string]int
	points     map[string]int
}

func makeRaceState(deers []reindeerStat) raceState {
	state := raceState{
		time:       0,
		stats:      make(map[string]reindeerStat),
		flying:     make(map[string]bool),
		distances:  make(map[string]int),
		nextChange: make(map[string]int),
		points:     make(map[string]int),
	}
	for _, deer := range deers {
		state.stats[deer.name] = deer
		state.distances[deer.name] = 0
		state.flying[deer.name] = true
		state.nextChange[deer.name] = deer.flyDuration
		state.points[deer.name] = 0
	}
	return state
}

func (state *raceState) Step() {
	for name, stat := range state.stats {
		if state.nextChange[name] == 0 {
			if state.flying[name] {
				state.nextChange[name] = stat.restDuration
			} else {
				state.nextChange[name] = stat.flyDuration
			}
			state.flying[name] = !(state.flying[name])
		}
		if state.flying[name] {
			state.distances[name] += stat.speed
		}
		state.nextChange[name] -= 1
	}

	type deerDist struct {
		name string
		dist int
	}

	maxDeers := make([]deerDist, 0)
	for deer, distance := range state.distances {
		if len(maxDeers) == 0 {
			maxDeers = append(maxDeers, deerDist{name: deer, dist: distance})
			continue
		}
		if distance < maxDeers[0].dist {
			continue
		}

		if distance > maxDeers[0].dist {
			maxDeers = maxDeers[:0]
		}
		maxDeers = append(maxDeers, deerDist{name: deer, dist: distance})
	}

	for _, deer := range maxDeers {
		state.points[deer.name] += 1
	}

	state.time += 1
}

type reindeerStat struct {
	name         string
	speed        int
	flyDuration  int
	restDuration int
}

func parseAllReindeerStats(text string) []reindeerStat {
	deers := make([]reindeerStat, 0)
	for line := range strings.SplitSeq(text, "\n") {
		deer, err := parseReindeer(line)
		if err != nil {
			continue
		}
		deers = append(deers, deer)
	}
	return deers
}

func parseReindeer(line string) (reindeerStat, error) {
	tokens := strings.Split(line, " ")
	if len(tokens) < 15 {
		return reindeerStat{}, errors.New("Invalid line")
	}

	name := tokens[0]

	speed, err := strconv.Atoi(tokens[3])
	if err != nil {
		return reindeerStat{}, err
	}

	flyDuration, err := strconv.Atoi(tokens[6])
	if err != nil {
		return reindeerStat{}, err
	}

	restDuration, err := strconv.Atoi(tokens[13])
	if err != nil {
		return reindeerStat{}, err
	}

	return reindeerStat{
		name:         name,
		speed:        speed,
		flyDuration:  flyDuration,
		restDuration: restDuration,
	}, nil
}
