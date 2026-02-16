package common

import "io"

type Solver interface {
	ArgsString(part int, args []string) string
	Solve(part int, args []string, w io.Writer) error
}
