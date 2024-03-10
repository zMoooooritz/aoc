package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"strings"

	"github.com/zMoooooritz/advent-of-code/cast"
	"github.com/zMoooooritz/advent-of-code/util"
)

//go:embed input.txt
var input string

type OpCode int

const (
	NULLCODE OpCode = iota
	ADD
	MUL
	IN
	OUT
	JIT
	JIF
	LT
	EQ
	HALT = 99
)

var opArgCount = map[OpCode]int{
	ADD:  3,
	MUL:  3,
	IN:   1,
	OUT:  1,
	JIT:  2,
	JIF:  2,
	LT:   3,
	EQ:   3,
	HALT: 0,
}

type ParameterMode int

const (
	POSITION ParameterMode = iota
	IMMEDIATE
)

type Parameter struct {
	value int
	mode  ParameterMode
}

type Instruction struct {
	opcode      OpCode
	firstParam  Parameter
	secondParam Parameter
	thridParam  Parameter
}

var instructions []int

func init() {
	// do this in init (not main) so test file has same input
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
}

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

	if part == 1 {
		ans := part1(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) int {
	instructions = parseInput(input)

	return executeInstructions()
}

func part2(input string) int {
	instructions = parseInput(input)

	return executeInstructions()
}

func executeInstructions() int {
	instructionPointer := 0
	for instructionPointer < len(instructions) {
		// fmt.Println(instructions)
		op := parseInstruction(instructionPointer)
		opcode := op.opcode
		argCount := operationArgumentCount(opcode)
		// fmt.Println("instPointer:", instructionPointer, "op:", op)

		if opcode == HALT {
			break
		}

		instructionPointer += argCount + 1
		switch opcode {
		case ADD:
			instructions[op.thridParam.value] = evaluateParameter(op.firstParam) + evaluateParameter(op.secondParam)
			break
		case MUL:
			instructions[op.thridParam.value] = evaluateParameter(op.firstParam) * evaluateParameter(op.secondParam)
			break
		case IN:
			instructions[op.firstParam.value] = 5
		case OUT:
			fmt.Println(evaluateParameter(op.firstParam))
		case JIT:
			if evaluateParameter(op.firstParam) != 0 {
				instructionPointer = evaluateParameter(op.secondParam)
			}
		case JIF:
			if evaluateParameter(op.firstParam) == 0 {
				instructionPointer = evaluateParameter(op.secondParam)
			}
		case LT:
			if evaluateParameter(op.firstParam) < evaluateParameter(op.secondParam) {
				instructions[op.thridParam.value] = 1
			} else {
				instructions[op.thridParam.value] = 0
			}
		case EQ:
			if evaluateParameter(op.firstParam) == evaluateParameter(op.secondParam) {
				instructions[op.thridParam.value] = 1
			} else {
				instructions[op.thridParam.value] = 0
			}
		default:
			panic("invalid opcode")
		}
	}

	return instructions[0]
}

func parseInstruction(instructionPointer int) Instruction {
	operation := Instruction{}
	instruction := instructions[instructionPointer]
	opcode := OpCode(instruction % 100)
	argCount := operationArgumentCount(opcode)

	parameterModes := instruction / 100

	operation.opcode = opcode
	counter := 1
	if argCount >= 1 {
		operation.firstParam = Parameter{instructions[instructionPointer+counter], ParameterMode((parameterModes / int(math.Pow10(counter-1))) % 10)}
	}
	counter++
	if argCount >= 2 {
		operation.secondParam = Parameter{instructions[instructionPointer+counter], ParameterMode((parameterModes / int(math.Pow10(counter-1))) % 10)}
	}
	counter++
	if argCount >= 3 {
		operation.thridParam = Parameter{instructions[instructionPointer+counter], ParameterMode((parameterModes / int(math.Pow10(counter-1))) % 10)}
	}

	return operation
}

func operationArgumentCount(opCode OpCode) int {
	if count, ok := opArgCount[opCode]; ok {
		return count
	}
	panic("invalid opcode")
}

func evaluateParameter(param Parameter) int {
	if param.mode == POSITION {
		return instructions[param.value]
	} else {
		return param.value
	}
}

func parseInput(input string) (ints []int) {
	for _, line := range strings.Split(input, ",") {
		ints = append(ints, cast.ToInt(line))
	}
	return ints
}
