package y2015

import (
	"errors"

	y2015d1 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/1"
	y2015d2 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/2"
	y2015d3 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/3"
	y2015d4 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/4"
	y2015d5 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/5"
	y2015d6 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/6"
	y2015d7 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/7"
	"github.com/AnirudhKanaparthy/advent-of-gopher/common"
)

func MakeSolution(day int) (common.Solver, error) {
	switch day {
	case 1:
		return y2015d1.MakeSolution(), nil
	case 2:
		return y2015d2.MakeSolution(), nil
	case 3:
		return y2015d3.MakeSolution(), nil
	case 4:
		return y2015d4.MakeSolution(), nil
	case 5:
		return y2015d5.MakeSolution(), nil
	case 6:
		return y2015d6.MakeSolution(), nil
	case 7:
		return y2015d7.MakeSolution(), nil
	}
	return nil, errors.New("Invalid day provided")
}
