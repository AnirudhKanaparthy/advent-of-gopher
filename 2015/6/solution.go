package y2015d6

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
	if len(args) < 1 {
		return errors.New("No input file provided")
	}

	inputFilePath := args[0]

	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}
	text := string(data)

	const (
		gridWidth  int = 1000
		gridHeight int = 1000
	)

	source := strings.Trim(text, " \t\n")

	switch part {
	case 1:
		result, err := solveBool(source, gridWidth, gridHeight)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Number of lights on: %v", result)
	case 2:
		result, err := solveInt(source, gridWidth, gridHeight)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Total brightness: %v", result)
	default:
		return fmt.Errorf("Expected part to be 1 or 2, got %v", part)
	}

	return nil
}

func solveBool(source string, gridWidth int, gridHeight int) (int, error) {
	grid := makeGrid[bool](gridWidth, gridHeight)
	for _, line := range strings.Split(source, "\n") {
		line = strings.Trim(line, " \t\n")

		ctx := makeLexerCtx(line)
		instruction, err := parseInstruction(&ctx)
		if err != nil {
			return -1, err
		}

		err = doBool(&grid, instruction)
		if err != nil {
			return -1, fmt.Errorf("Cannot perform instruction: %v", err.Error())
		}
	}
	return countBool(&grid), nil
}

func solveInt(source string, gridWidth int, gridHeight int) (int, error) {
	grid := makeGrid[int](gridWidth, gridHeight)
	for _, line := range strings.Split(source, "\n") {
		line = strings.Trim(line, " \t\n")

		ctx := makeLexerCtx(line)
		instruction, err := parseInstruction(&ctx)
		if err != nil {
			return -1, err
		}

		err = doInt(&grid, instruction)
		if err != nil {
			return -1, fmt.Errorf("Cannot perform instruction: %v", err.Error())
		}
	}
	return countInt(&grid), nil
}

type grid[T bool | int] struct {
	grid   []T
	width  int
	height int
}

func makeGrid[T bool | int](width int, height int) grid[T] {
	return grid[T]{
		make([]T, width*height),
		width,
		height,
	}
}

func (g grid[T]) index(x int, y int) (int, error) {
	if x < 0 || g.width <= x {
		return -1, fmt.Errorf("x (%v) is out of range (0, %v)", x, g.width)
	}
	if y < 0 || g.height <= y {
		return -1, fmt.Errorf("y (%v) is out of range (0, %v)", y, g.height)
	}
	return y*g.height + x, nil
}

func (g grid[T]) at(x int, y int) (T, error) {
	idx, err := g.index(x, y)
	if err != nil {
		var d T
		return d, err
	}
	return g.grid[idx], nil
}

func (g *grid[T]) set(x int, y int, v T) error {
	idx, err := g.index(x, y)
	if err != nil {
		return err
	}
	g.grid[idx] = v
	return nil
}

func (g *grid[T]) setRange(gridRange pointRange, callback func(v T) T) error {
	fromX := min(gridRange.from.x, gridRange.to.x)
	fromY := min(gridRange.from.y, gridRange.to.y)

	toX := max(gridRange.from.x, gridRange.to.x)
	toY := max(gridRange.from.y, gridRange.to.y)

	for i := fromX; i < toX+1; i += 1 {
		for j := fromY; j < toY+1; j += 1 {
			val, err := g.at(i, j)
			if err != nil {
				return err
			}
			if err := g.set(i, j, callback(val)); err != nil {
				return err
			}
		}
	}
	return nil
}

func doBool(g *grid[bool], instruction instruction) error {
	switch instruction.instrutionType {
	case instructionTurnOn:
		return g.setRange(instruction.gridRange, func(bool) bool { return true })
	case instructionTurnOff:
		return g.setRange(instruction.gridRange, func(bool) bool { return false })
	case instructionToggle:
		return g.setRange(instruction.gridRange, func(v bool) bool { return !v })
	default:
		return errors.New("Invalid instruction")
	}
}

func doInt(g *grid[int], instruction instruction) error {
	switch instruction.instrutionType {
	case instructionTurnOn:
		return g.setRange(instruction.gridRange, func(v int) int { return v + 1 })
	case instructionTurnOff:
		return g.setRange(instruction.gridRange, func(v int) int {
			if v > 0 {
				return v - 1
			}
			return v
		})
	case instructionToggle:
		return g.setRange(instruction.gridRange, func(v int) int { return v + 2 })
	default:
		return errors.New("Invalid instruction")
	}
}

func countBool(g *grid[bool]) int {
	count := 0
	for _, isOn := range g.grid {
		if isOn {
			count += 1
		}
	}
	return count
}

func countInt(g *grid[int]) int {
	count := 0
	for _, b := range g.grid {
		count += b
	}
	return count
}

type tokenType int

const (
	tokenOff tokenType = iota
	tokenOn
	tokenTurn
	tokenToggle
	tokenThrough
	tokenNumber
	tokenComma
)

func (tt tokenType) string() string {
	switch tt {
	case tokenOff:
		return "tokenOff"
	case tokenOn:
		return "tokenOn"
	case tokenTurn:
		return "tokenTurn"
	case tokenToggle:
		return "tokenToggle"
	case tokenThrough:
		return "tokenThrough"
	case tokenNumber:
		return "tokenNumber"
	case tokenComma:
		return "tokenComma"
	default:
		return "tokenUnknown"
	}
}

type token struct {
	tokenType tokenType
	text      string
}

type lexerCtx struct {
	source string
	pos    int
}

func makeLexerCtx(source string) lexerCtx {
	return lexerCtx{source, 0}
}

func (ctx *lexerCtx) nextChar() byte {
	if ctx.pos >= len(ctx.source) {
		return 0
	}

	c := ctx.source[ctx.pos]
	ctx.pos += 1
	return c
}

func (ctx *lexerCtx) stepBack() {
	if ctx.pos <= 0 {
		return
	}
	ctx.pos -= 1
}

func isAlphabet(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (ctx *lexerCtx) nextToken() (token, error) {
	c := ctx.nextChar()

	for c == ' ' || c == '\t' {
		c = ctx.nextChar()
	}

	switch c {
	case ',':
		return token{tokenComma, ","}, nil
	case 0:
		return token{}, errors.New("end")
	}

	if isAlphabet(c) {
		word := ""
		for isAlphabet(c) {
			word += string(c)
			c = ctx.nextChar()
		}
		if c != 0 {
			ctx.stepBack()
		}

		t := token{}
		switch word {
		case "turn":
			t.tokenType = tokenTurn
		case "toggle":
			t.tokenType = tokenToggle
		case "through":
			t.tokenType = tokenThrough
		case "on":
			t.tokenType = tokenOn
		case "off":
			t.tokenType = tokenOff
		default:
			return token{}, errors.New("invalid keyword")
		}

		t.text = word
		return t, nil
	}

	if isDigit(c) {
		word := ""
		for isDigit(c) {
			word += string(c)
			c = ctx.nextChar()
		}
		if c != 0 {
			ctx.stepBack()
		}

		return token{tokenNumber, word}, nil
	}

	return token{}, errors.New("unknown token")
}

// instruction Grammar
//
// instruction := Command [CommandArg] pointRange
// Command     := "turn" | "toggle"
// CommandArg  := "on" | "off"
// pointRange       := Number "," Number "through" Number "," Number
// Number      := [0-9]+

type point struct {
	x, y int
}

type pointRange struct {
	from point
	to   point
}

type instructionType int

const (
	instructionUnknown instructionType = iota
	instructionTurnOn
	instructionTurnOff
	instructionToggle
)

type instruction struct {
	instrutionType instructionType
	gridRange      pointRange
}

func parseInstructionType(lexer *lexerCtx) (instructionType, error) {
	curToken, err := lexer.nextToken()
	if err != nil {
		return instructionUnknown, err
	}

	switch curToken.tokenType {
	case tokenTurn:
		curToken, err = lexer.nextToken()
		if err != nil {
			return instructionUnknown, err
		}
		switch curToken.tokenType {
		case tokenOn:
			return instructionTurnOn, nil
		case tokenOff:
			return instructionTurnOff, nil
		default:
			return instructionUnknown, errors.New("You can either turn 'on' or 'off'")
		}
	case tokenToggle:
		return instructionToggle, nil
	default:
		return instructionUnknown, errors.New("instruction is neither 'turn' or 'toggle'")
	}
}

func parseNumber(lexer *lexerCtx) (int, error) {
	curToken, err := lexer.nextToken()
	if err != nil {
		return -1, err
	}
	if curToken.tokenType != tokenNumber {
		return -1, fmt.Errorf("Expected number, got: '%v'", curToken.text)
	}
	x, err := strconv.Atoi(curToken.text)
	if err != nil {
		return -1, fmt.Errorf("Expected number, got: '%v'", curToken.text)
	}
	return x, nil
}

func parsePoint(lexer *lexerCtx) (point, error) {
	x, err := parseNumber(lexer)
	if err != nil {
		return point{}, err
	}

	curToken, err := lexer.nextToken()
	if curToken.tokenType != tokenComma {
		return point{}, fmt.Errorf("Expected ',', got '%v'", curToken.text)
	}

	y, err := parseNumber(lexer)
	if err != nil {
		return point{}, err
	}

	return point{x: x, y: y}, nil
}

func parseRange(lexer *lexerCtx) (pointRange, error) {
	pointX, err := parsePoint(lexer)
	if err != nil {
		return pointRange{}, nil
	}

	curToken, err := lexer.nextToken()
	if curToken.tokenType != tokenThrough {
		return pointRange{}, fmt.Errorf("Expected 'through', got '%v'", curToken.text)
	}

	pointY, err := parsePoint(lexer)
	if err != nil {
		return pointRange{}, nil
	}

	return pointRange{from: pointX, to: pointY}, nil
}

func parseInstruction(lexer *lexerCtx) (instruction, error) {
	instructionType, err := parseInstructionType(lexer)
	if err != nil {
		return instruction{}, err
	}
	gridRange, err := parseRange(lexer)
	if err != nil {
		return instruction{}, err
	}
	return instruction{
		instrutionType: instructionType,
		gridRange:      gridRange,
	}, nil
}
