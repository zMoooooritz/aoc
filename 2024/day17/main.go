package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"slices"
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

type CPU struct {
	regA int
	regB int
	regC int

	instructionPointer int

	instructions []int
	output       []int
}

func (c *CPU) fetchVal(index int) int {
	if index < 4 {
		return index
	} else if index == 4 {
		return c.regA
	} else if index == 5 {
		return c.regB
	} else if index == 6 {
		return c.regC
	} else {
		panic("invalid")
	}
}

func (c *CPU) run() {
	opMap := map[int]op{
		0: c.adv,
		1: c.bxl,
		2: c.bst,
		3: c.jnz,
		4: c.bxc,
		5: c.out,
		6: c.bdv,
		7: c.cdv,
	}

	for c.instructionPointer < len(c.instructions) {
		oldInstPointer := c.instructionPointer
		opId := c.instructions[c.instructionPointer]
		opMap[opId](c.instructions[c.instructionPointer+1])
		if oldInstPointer == c.instructionPointer {
			c.gotoNextInstruction()
		}
	}
}

func (c *CPU) gotoNextInstruction() {
	c.instructionPointer += 2
}

func (c *CPU) div(operand int) int {
	numerator := c.regA
	denominator := int(math.Pow(2, float64(c.fetchVal(operand))))
	return int(numerator / denominator)
}

func (c *CPU) adv(operand int) {
	c.regA = c.div(operand)
}

func (c *CPU) bxl(operand int) {
	c.regB = c.regB ^ operand
}

func (c *CPU) bst(operand int) {
	c.regB = c.fetchVal(operand) % 8
}

func (c *CPU) jnz(operand int) {
	if c.regA == 0 {
		return
	}
	c.instructionPointer = operand
}

func (c *CPU) bxc(operand int) {
	_ = operand
	c.regB = c.regB ^ c.regC
}

func (c *CPU) out(operand int) {
	c.output = append(c.output, c.fetchVal(operand)%8)
}

func (c *CPU) bdv(operand int) {
	c.regB = c.div(operand)
}

func (c *CPU) cdv(operand int) {
	c.regC = c.div(operand)
}

type op func(int)

func part1(input string) int {
	cpu := parseInput(input)

	cpu.run()

	strOut := []string{}
	for _, o := range cpu.output {
		strOut = append(strOut, cast.ToString(o))
	}
	fmt.Println(strings.Join(strOut[:], ","))

	return 0
}

func copyCPU(cpu CPU) CPU {
	newCPU := cpu
	newCPU.instructions = make([]int, len(cpu.instructions))
	copy(newCPU.instructions, cpu.instructions)
	newCPU.output = []int{}
	return newCPU
}

func part2(input string) int {
	baseCPU := parseInput(input)

	seed := 0
	for iteration := len(baseCPU.instructions) - 1; iteration >= 0; iteration -= 1 {
		seed <<= 3
		for {
			cpu := copyCPU(baseCPU)
			cpu.regA = seed
			cpu.run()
			if !slices.Equal(cpu.output, baseCPU.instructions[iteration:]) {
				seed += 1
			} else {
				break
			}
		}
	}
	return seed
}

func parseInput(input string) CPU {
	cpu := CPU{}
	lines := strings.Split(input, "\n")
	cpu.regA = cast.ToInt(strings.TrimPrefix(lines[0], "Register A: "))
	cpu.regB = cast.ToInt(strings.TrimPrefix(lines[1], "Register B: "))
	cpu.regC = cast.ToInt(strings.TrimPrefix(lines[2], "Register C: "))
	cpu.instructions = cast.ToIntSliceSep(strings.TrimPrefix(lines[4], "Program: "), ",")
	return cpu
}
