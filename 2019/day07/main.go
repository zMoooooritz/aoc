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

type CPU struct {
	instructions       []int
	instructionPointer int

	inputs       []int
	inputPointer int

	output int
	halted bool

	pauseOnIO bool
}

func NewCPU(instructions []int, inputs []int) CPU {
	cpu := CPU{
		instructions:       instructions,
		instructionPointer: 0,
		inputs:             inputs,
		inputPointer:       0,
		output:             0,
		halted:             false,
		pauseOnIO:          false,
	}
	return cpu
}

func (c *CPU) Run() {
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
			c.halted = true
			break
		}

		c.instructionPointer += argCount + 1 // can be altered in the instructions
		if fn, ok := opCodeMap[inst.opcode]; ok {
			fn(inst.params)
		} else {
			panic("invalid opcode")
		}

		if c.pauseOnIO && inst.opcode == OUT {
			break
		}
	}
}

func (c *CPU) DoPauseOnIO(doPause bool) {
	c.pauseOnIO = doPause
}

func (c *CPU) HasHalted() bool {
	return c.halted
}

func (c *CPU) SendInput(input int) {
	c.inputs = append(c.inputs, input)
}

func (c *CPU) RecvOutput() int {
	return c.output
}

func (c *CPU) add(args []Parameter) {
	c.instructions[args[2].value] = c.evalParameter(args[0]) + c.evalParameter(args[1])
}

func (c *CPU) mul(args []Parameter) {
	c.instructions[args[2].value] = c.evalParameter(args[0]) * c.evalParameter(args[1])
}

func (c *CPU) in(args []Parameter) {
	c.instructions[args[0].value] = c.inputs[c.inputPointer]
	c.inputPointer += 1
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

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func part1(input string) int {
	instructions := parseInput(input)

	maxVal := 0
	maxPerm := []int{}
	perms := permutations([]int{0, 1, 2, 3, 4})
	for _, perm := range perms {
		curr := 0
		for iteration := range 5 {
			inst := make([]int, len(instructions))
			copy(inst, instructions)
			inputs := []int{perm[iteration], curr}
			cpu := NewCPU(inst, inputs)
			cpu.Run()
			curr = cpu.RecvOutput()
		}

		if curr > maxVal {
			maxVal = curr
			maxPerm = perm
		}
	}

	fmt.Println("MaxPerm:", maxPerm, "MaxVal:", maxVal)

	return maxVal
}

func part2(input string) int {
	instructions := parseInput(input)

	maxVal := 0
	maxPerm := []int{}
	perms := permutations([]int{5, 6, 7, 8, 9})
	for _, perm := range perms {
		cpus := [5]CPU{}
		for cpuIndex := range len(cpus) {
			inst := make([]int, len(instructions))
			copy(inst, instructions)
			cpus[cpuIndex] = NewCPU(inst, []int{perm[cpuIndex]})
			cpus[cpuIndex].DoPauseOnIO(true)
		}

		val := 0
		currCPU := 0
		for {
			cpus[currCPU].SendInput(val)
			cpus[currCPU].Run()

			if cpus[currCPU].HasHalted() {
				break
			}

			val = cpus[currCPU].RecvOutput()
			currCPU = (currCPU + 1) % len(cpus)
		}

		if val > maxVal {
			maxVal = val
			maxPerm = perm
		}
	}
	fmt.Println("MaxPerm:", maxPerm, "MaxVal:", maxVal)

	return maxVal
}

func parseInput(input string) []int {
	return cast.ToIntSliceSep(input, ",")
}
