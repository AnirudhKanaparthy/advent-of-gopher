package y2015d7

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
	return "<file path>"
}

func (Solution) Solve(part int, args []string, w io.Writer) error {
	filepath := args[0]
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	source := string(data)

	expressions := make(map[int][]string)
	i := 0
	for line := range strings.SplitSeq(source, "\n") {
		expressions[i] = strings.Fields(line)
		i += 1
	}

	wires := make(map[string]uint16)

	for len(expressions) != 0 {
		evalLoop(expressions, wires)
	}

	valA, ok := wires["a"]
	if !ok {
		return errors.New("Cannot find wire with name 'a'")
	}

	switch part {
	case 1:
	case 2:
		wires = make(map[string]uint16)
		wires["b"] = valA

		expressions = make(map[int][]string)
		i := 0
		for line := range strings.SplitSeq(source, "\n") {
			expressions[i] = strings.Fields(line)
			i += 1
		}

		for len(expressions) != 0 {
			evalLoop(expressions, wires)
		}
		valA, ok = wires["a"]
		if !ok {
			return errors.New("Cannot find wire with name 'a' the second time")
		}
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	fmt.Fprintf(w, "Value for a: %v\n", valA)

	return nil
}

func evalNumber(tokens *[]string, wires map[string]uint16) (uint16, error) {
	valTxt := (*tokens)[0]
	val, err := strconv.Atoi(valTxt)
	if err != nil {
		val, ok := wires[valTxt]
		if !ok {
			return 0, errors.New("Cannot be evaluated")
		}

		*tokens = (*tokens)[1:]
		return val, nil
	}
	*tokens = (*tokens)[1:]
	return uint16(val), nil
}

func evalUnary(tokens *[]string, wires map[string]uint16) (uint16, error) {
	if len(*tokens) < 2 {
		return 0, errors.New("Invalid expression")
	}
	if (*tokens)[0] != "NOT" {
		return 0, errors.New("Invalid unary expression")
	}
	t := *tokens
	*tokens = (*tokens)[1:]

	val, err := evalNumber(tokens, wires)
	if err != nil {
		*tokens = t
		return val, err
	}

	return ^val, nil
}

func evalBinary(tokens *[]string, wires map[string]uint16) (uint16, error) {
	if len(*tokens) < 3 {
		return 0, errors.New("Invalid expression")
	}
	valA, err := evalNumber(tokens, wires)
	if err != nil {
		return valA, err
	}

	t := *tokens
	op := (*tokens)[0]
	*tokens = (*tokens)[1:]

	valB, err := evalNumber(tokens, wires)
	if err != nil {
		*tokens = t
		return valB, err
	}

	var res uint16
	switch op {
	case "AND":
		res = valA & valB
	case "OR":
		res = valA | valB
	case "RSHIFT":
		res = valA >> valB
	case "LSHIFT":
		res = valA << valB
	default:
		*tokens = t
		return 0, errors.New("Cannot be evaluated")
	}

	return res, nil
}

func ParseExpression(tokens *[]string, wires map[string]uint16) (uint16, error) {
	t := *tokens

	val, err := evalBinary(tokens, wires)
	if err == nil {
		return val, err
	}
	*tokens = t

	val, err = evalUnary(tokens, wires)
	if err == nil {
		return val, err
	}
	*tokens = t

	val, err = evalNumber(tokens, wires)
	if err == nil {
		return val, err
	}
	*tokens = t

	return val, errors.New("Invalid expression")
}

func evalStatement(tokens *[]string, wires map[string]uint16) error {
	t := *tokens

	val, err := ParseExpression(tokens, wires)
	if err != nil {
		*tokens = t
		return err
	}
	if len(*tokens) < 2 {
		*tokens = t
		return errors.New("Expression is not being assigned")
	}
	if (*tokens)[0] != "->" {
		*tokens = t
		return errors.New("Invalid statement")
	}

	wires[(*tokens)[1]] = val
	return nil
}

func evalLoop(expressions map[int][]string, wires map[string]uint16) {
	toDelete := make([]int, 0)
	for i, expression := range expressions {
		if evalStatement(&expression, wires) == nil {
			toDelete = append(toDelete, i)
		}
	}
	for _, i := range toDelete {
		delete(expressions, i)
	}
}
