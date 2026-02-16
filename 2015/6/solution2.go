package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFilePath = "input.txt"
const (
	gridWidth  int = 1000
	gridHeight int = 1000
)

type Grid struct {
	grid   []int
	width  int
	height int
}

func MakeGrid(width int, height int) Grid {
	return Grid{
		make([]int, width*height),
		width,
		height,
	}
}

func (g *Grid) At(x int, y int) (int, error) {
	if x < 0 || g.width <= x {
		return -1, fmt.Errorf("x (%v) is out of range (0, %v)", x, g.width)
	}
	if y < 0 || g.height <= y {
		return -1, fmt.Errorf("y (%v) is out of range (0, %v)", y, g.height)
	}
	return g.grid[y*g.height+x], nil
}

func (g *Grid) Set(x int, y int, v int) error {
	if x < 0 || g.width <= x {
		return fmt.Errorf("x (%v) is out of range (0, %v)", x, g.width)
	}
	if y < 0 || g.height <= y {
		return fmt.Errorf("y (%v) is out of range (0, %v)", y, g.height)
	}
	g.grid[y*g.height+x] = v
	return nil
}

func (g *Grid) Toggle(x int, y int) error {
	val, err := g.At(x, y)
	if err != nil {
		return err
	}
	val += 2
	return g.Set(x, y, val)
}

func (g *Grid) Increment(x int, y int) error {
	val, err := g.At(x, y)
	if err != nil {
		return err
	}
	val += 1
	return g.Set(x, y, val)
}

func (g *Grid) Decrement(x int, y int) error {
	val, err := g.At(x, y)
	if err != nil {
		return err
	}
	if val > 0 {
		val -= 1
	}
	return g.Set(x, y, val)
}

func (g *Grid) DoRange(gridRange Range, callback func(int) int) error {
	fromX := min(gridRange.from.x, gridRange.to.x)
	fromY := min(gridRange.from.y, gridRange.to.y)

	toX := max(gridRange.from.x, gridRange.to.x)
	toY := max(gridRange.from.y, gridRange.to.y)

	for i := fromX; i < toX+1; i += 1 {
		for j := fromY; j < toY+1; j += 1 {
			v, err := g.At(i, j)
			if err != nil {
				return err
			}
			if err := g.Set(i, j, callback(v)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Grid) SetRange(v int, gridRange Range) error {
	return g.DoRange(gridRange, func(int) int { return v })
}

func (g *Grid) ToggleRange(gridRange Range) error {
	return g.DoRange(gridRange, func(v int) int { return v + 2 })
}

func (g *Grid) IncrementRange(gridRange Range) error {
	return g.DoRange(gridRange, func(v int) int { return v + 1 })
}

func (g *Grid) DecrementRange(gridRange Range) error {
	return g.DoRange(gridRange, func(v int) int {
		if v > 0 {
			return v - 1
		}
		return v
	})
}

func (g *Grid) Do(instruction Instruction) error {
	switch instruction.instrutionType {
	case InstructionTurnOn:
		return g.IncrementRange(instruction.gridRange)
	case InstructionTurnOff:
		return g.DecrementRange(instruction.gridRange)
	case InstructionToggle:
		return g.ToggleRange(instruction.gridRange)
	default:
		return errors.New("Invalid Instruction")
	}
}

func (g Grid) Count() int {
	count := 0
	for _, b := range g.grid {
		count += b
	}
	return count
}

func (g Grid) String() string {
	var sb strings.Builder
	for i := 0; i < g.width; i += 1 {
		for j := 0; j < g.height; j += 1 {
			b, err := g.At(i, j)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err.Error())
				sb.WriteString(". ")
			} else if b > 0 {
				sb.WriteString("# ")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

type TokenType int

const (
	TokenOff TokenType = iota
	TokenOn
	TokenTurn
	TokenToggle
	TokenThrough
	TokenNumber
	TokenComma
)

func (tt TokenType) String() string {
	switch tt {
	case TokenOff:
		return "TokenOff"
	case TokenOn:
		return "TokenOn"
	case TokenTurn:
		return "TokenTurn"
	case TokenToggle:
		return "TokenToggle"
	case TokenThrough:
		return "TokenThrough"
	case TokenNumber:
		return "TokenNumber"
	case TokenComma:
		return "TokenComma"
	default:
		return "TokenUnknown"
	}
}

type Token struct {
	tokenType TokenType
	text      string
}

type LexerCtx struct {
	source string
	pos    int
}

func MakeLexerCtx(source string) LexerCtx {
	return LexerCtx{source, 0}
}

func (ctx *LexerCtx) NextChar() byte {
	if ctx.pos >= len(ctx.source) {
		return 0
	}

	c := ctx.source[ctx.pos]
	ctx.pos += 1
	return c
}

func (ctx *LexerCtx) StepBack() {
	if ctx.pos <= 0 {
		return
	}
	ctx.pos -= 1
}

func IsAlphabet(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (ctx *LexerCtx) NextToken() (Token, error) {
	c := ctx.NextChar()

	for c == ' ' || c == '\t' {
		c = ctx.NextChar()
	}

	switch c {
	case ',':
		return Token{TokenComma, ","}, nil
	case 0:
		return Token{}, errors.New("end")
	}

	if IsAlphabet(c) {
		word := ""
		for IsAlphabet(c) {
			word += string(c)
			c = ctx.NextChar()
		}
		if c != 0 {
			ctx.StepBack()
		}

		t := Token{}
		switch word {
		case "turn":
			t.tokenType = TokenTurn
		case "toggle":
			t.tokenType = TokenToggle
		case "through":
			t.tokenType = TokenThrough
		case "on":
			t.tokenType = TokenOn
		case "off":
			t.tokenType = TokenOff
		default:
			return Token{}, errors.New("invalid keyword")
		}

		t.text = word
		return t, nil
	}

	if IsDigit(c) {
		word := ""
		for IsDigit(c) {
			word += string(c)
			c = ctx.NextChar()
		}
		if c != 0 {
			ctx.StepBack()
		}

		return Token{TokenNumber, word}, nil
	}

	return Token{}, errors.New("unknown token")
}

// Instruction Grammar
//
// Instruction := Command [CommandArg] Range
// Command     := "turn" | "toggle"
// CommandArg  := "on" | "off"
// Range       := Number "," Number "through" Number "," Number
// Number      := [0-9]+

type Point struct {
	x, y int
}

type Range struct {
	from Point
	to   Point
}

type InstructionType int

const (
	InstructionUnknown InstructionType = iota
	InstructionTurnOn
	InstructionTurnOff
	InstructionToggle
)

type Instruction struct {
	instrutionType InstructionType
	gridRange      Range
}

func ParseInstructionType(lexer *LexerCtx) (InstructionType, error) {
	curToken, err := lexer.NextToken()
	if err != nil {
		return InstructionUnknown, err
	}

	switch curToken.tokenType {
	case TokenTurn:
		curToken, err = lexer.NextToken()
		if err != nil {
			return InstructionUnknown, err
		}
		switch curToken.tokenType {
		case TokenOn:
			return InstructionTurnOn, nil
		case TokenOff:
			return InstructionTurnOff, nil
		default:
			return InstructionUnknown, errors.New("You can either turn 'on' or 'off'")
		}
	case TokenToggle:
		return InstructionToggle, nil
	default:
		return InstructionUnknown, errors.New("Instruction is neither 'turn' or 'toggle'")
	}
}

func ParseNumber(lexer *LexerCtx) (int, error) {
	curToken, err := lexer.NextToken()
	if err != nil {
		return -1, err
	}
	if curToken.tokenType != TokenNumber {
		return -1, fmt.Errorf("Expected number, got: '%v'", curToken.text)
	}
	x, err := strconv.Atoi(curToken.text)
	if err != nil {
		return -1, fmt.Errorf("Expected number, got: '%v'", curToken.text)
	}
	return x, nil
}

func ParsePoint(lexer *LexerCtx) (Point, error) {
	x, err := ParseNumber(lexer)
	if err != nil {
		return Point{}, err
	}

	curToken, err := lexer.NextToken()
	if curToken.tokenType != TokenComma {
		return Point{}, fmt.Errorf("Expected ',', got '%v'", curToken.text)
	}

	y, err := ParseNumber(lexer)
	if err != nil {
		return Point{}, err
	}

	return Point{x: x, y: y}, nil
}

func ParseRange(lexer *LexerCtx) (Range, error) {
	pointX, err := ParsePoint(lexer)
	if err != nil {
		return Range{}, nil
	}

	curToken, err := lexer.NextToken()
	if curToken.tokenType != TokenThrough {
		return Range{}, fmt.Errorf("Expected 'through', got '%v'", curToken.text)
	}

	pointY, err := ParsePoint(lexer)
	if err != nil {
		return Range{}, nil
	}

	return Range{from: pointX, to: pointY}, nil
}

func ParseInstruction(lexer *LexerCtx) (Instruction, error) {
	instructionType, err := ParseInstructionType(lexer)
	if err != nil {
		return Instruction{}, err
	}
	gridRange, err := ParseRange(lexer)
	if err != nil {
		return Instruction{}, err
	}
	return Instruction{
		instrutionType: instructionType,
		gridRange:      gridRange,
	}, nil
}

func main() {
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		panic(err)
	}

	source := strings.Trim(string(data), " \t\n")
	grid := MakeGrid(gridWidth, gridHeight)

	for i, line := range strings.Split(source, "\n") {
		line = strings.Trim(line, " \t\n")
		fmt.Printf("Instruction %v: %v\n", i, line)

		ctx := MakeLexerCtx(line)
		instruction, err := ParseInstruction(&ctx)
		if err != nil {
			panic(err)
		}

		err = grid.Do(instruction)
		if err != nil {
			panic(fmt.Sprintf("Cannot perform instruction: %v", err.Error()))
		} else {
			fmt.Println("Successfully performed instruction")
		}
		fmt.Printf("\n")
	}
	fmt.Println(grid.String())
	fmt.Printf("Brightness: %v\n", grid.Count())
}
