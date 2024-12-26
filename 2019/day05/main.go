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
	opcode OpCode
	params []Parameter
}

func NewCPU(instructions []int) CPU {
	cpu := CPU{
		instructions:       instructions,
		output:             0,
		instructionPointer: 0,
	}
	return cpu
}

type CPU struct {
	instructions []int
	output       int

	instructionPointer int
}

func (c *CPU) run() {
	opCodeMap := map[OpCode]func([]Parameter){
		ADD: c.add,
		MUL: c.mul,
		IN:  c.in,
		OUT: c.out,
		JIT: c.jit,
		JIF: c.jif,
		LT:  c.lt,
		EQ:  c.eq,
	}

	for c.instructionPointer < len(c.instructions) {
		inst := c.parseCurrentInstruction()
		argCount := c.operationArgumentCount(inst.opcode)

		if inst.opcode == HALT {
			break
		}

		c.instructionPointer += argCount + 1 // can be altered in the instructions
		if fn, ok := opCodeMap[inst.opcode]; ok {
			fn(inst.params)
		} else {
			panic("invalid opcode")
		}
	}
}

func (c *CPU) add(args []Parameter) {
	c.instructions[args[2].value] = c.evalParameter(args[0]) + c.evalParameter(args[1])
}

func (c *CPU) mul(args []Parameter) {
	c.instructions[args[2].value] = c.evalParameter(args[0]) * c.evalParameter(args[1])
}

func (c *CPU) in(args []Parameter) {
	c.instructions[args[0].value] = 5
}

func (c *CPU) out(args []Parameter) {
	c.output = c.evalParameter(args[0])
}

func (c *CPU) jit(args []Parameter) {
	if c.evalParameter(args[0]) != 0 {
		c.instructionPointer = c.evalParameter(args[1])
	}
}

func (c *CPU) jif(args []Parameter) {
	if c.evalParameter(args[0]) == 0 {
		c.instructionPointer = c.evalParameter(args[1])
	}

}

func (c *CPU) lt(args []Parameter) {
	if c.evalParameter(args[0]) < c.evalParameter(args[1]) {
		c.instructions[args[2].value] = 1
	} else {
		c.instructions[args[2].value] = 0
	}
}

func (c *CPU) eq(args []Parameter) {
	if c.evalParameter(args[0]) == c.evalParameter(args[1]) {
		c.instructions[args[2].value] = 1
	} else {
		c.instructions[args[2].value] = 0
	}
}

func (c *CPU) parseCurrentInstruction() Instruction {
	instruction := Instruction{}
	inst := c.instructions[c.instructionPointer]
	opcode := OpCode(inst % 100)
	argCount := c.operationArgumentCount(opcode)

	parameterModes := inst / 100

	instruction.opcode = opcode

	fmt.Println("OpCode:", opcode)

	for index := range argCount {
		instruction.params = append(instruction.params, Parameter{c.instructions[c.instructionPointer+index+1], ParameterMode((parameterModes / int(math.Pow10(index))) % 10)})
	}
	return instruction
}

func (c *CPU) operationArgumentCount(opCode OpCode) int {
	if count, ok := opArgCount[opCode]; ok {
		return count
	}
	panic("invalid opcode")
}

func (c *CPU) evalParameter(param Parameter) int {
	if param.mode == POSITION {
		return c.instructions[param.value]
	} else {
		return param.value
	}
}

func (c *CPU) fetchOutput() int {
	return c.output
}

func (c *CPU) fetchFirstRegister() int {
	return c.instructions[0]
}

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
	instructions := parseInput(input)

	cpu := NewCPU(instructions)
	cpu.run()

	test := true

	result := cpu.fetchOutput()
	if test {
		result = cpu.fetchFirstRegister()
	}
	return result
}

func part2(input string) int {
	instructions := parseInput(input)

	cpu := NewCPU(instructions)
	cpu.run()

	test := true

	result := cpu.fetchOutput()
	if test {
		result = cpu.fetchFirstRegister()
	}
	return result
}

func parseInput(input string) (ints []int) {
	for _, line := range strings.Split(input, ",") {
		ints = append(ints, cast.ToInt(line))
	}
	return ints
}
