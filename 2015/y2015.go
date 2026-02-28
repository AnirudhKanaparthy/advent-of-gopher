package y2015

import (
	"errors"

	y2015d1 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/1"
	y2015d10 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/10"
	y2015d11 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/11"
	y2015d12 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/12"
	y2015d13 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/13"
	y2015d14 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/14"
	y2015d15 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/15"
	y2015d16 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/16"
	y2015d17 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/17"
	y2015d18 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/18"
	y2015d19 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/19"
	y2015d2 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/2"
	y2015d3 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/3"
	y2015d4 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/4"
	y2015d5 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/5"
	y2015d6 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/6"
	y2015d7 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/7"
	y2015d8 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/8"
	y2015d9 "github.com/AnirudhKanaparthy/advent-of-gopher/2015/9"
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
	case 8:
		return y2015d8.MakeSolution(), nil
	case 9:
		return y2015d9.MakeSolution(), nil
	case 10:
		return y2015d10.MakeSolution(), nil
	case 11:
		return y2015d11.MakeSolution(), nil
	case 12:
		return y2015d12.MakeSolution(), nil
	case 13:
		return y2015d13.MakeSolution(), nil
	case 14:
		return y2015d14.MakeSolution(), nil
	case 15:
		return y2015d15.MakeSolution(), nil
	case 16:
		return y2015d16.MakeSolution(), nil
	case 17:
		return y2015d17.MakeSolution(), nil
	case 18:
		return y2015d18.MakeSolution(), nil
	case 19:
		return y2015d19.MakeSolution(), nil
	}
	return nil, errors.New("Invalid day provided")
}
