package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	y2015 "github.com/AnirudhKanaparthy/advent-of-gopher/2015"
	"github.com/AnirudhKanaparthy/advent-of-gopher/common"
)

func MakeSolution(year int, day int) (common.Solver, error) {
	switch year {
	case 2015:
		return y2015.MakeSolution(day)
	}
	return nil, errors.New("Invalid Year provided")
}

func prFatalln(f string, v ...any) {
	fmt.Printf(f, v...)
	fmt.Printf("\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 4 {
		prFatalln("Usage: %v <year> <day> <part> [solution specific args]...", os.Args[0])
	}

	year, err := strconv.Atoi(os.Args[1])
	if err != nil || year < 2015 || 2015 < year {
		prFatalln("Invalid Year")
	}
	day, err := strconv.Atoi(os.Args[2])
	if err != nil || day < 1 || 31 < day {
		prFatalln("Invalid Day")
	}

	part, err := strconv.Atoi(os.Args[3])
	if err != nil || part < 1 {
		prFatalln("Invalid Part")
	}

	solution, err := MakeSolution(year, day)
	if err != nil {
		prFatalln("%s", err.Error())
	}

	args := os.Args[4:]
	var sb strings.Builder
	err = solution.Solve(part, args, &sb)
	if err != nil {
		fmt.Printf("Usage: %v %v %v %v %v\n", os.Args[0], year, day, part, solution.ArgsString(part, args))
		prFatalln("%s", err.Error())
	}
	fmt.Printf("Solution:\n%v\n", sb.String())
}
