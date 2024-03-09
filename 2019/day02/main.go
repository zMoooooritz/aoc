package main

import (
	_ "embed"
	"flag"
	"fmt"
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

func part1(input string) int {
	ints := parseInput(input)

	ints[1] = 12
	ints[2] = 2

	return runInstructions(ints)
}

func part2(input string) int {
	ints := parseInput(input)

	target := 19690720

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			instructions := make([]int, len(ints))
			copy(instructions, ints)
			instructions[1] = i
			instructions[2] = j
			res := runInstructions(instructions)
			if res == target {
				return 100*i + j
			}
		}
	}

	panic("no answer found")
}

func runInstructions(instructions []int) int {
	index := 0
	for index < len(instructions) {
		opcode := instructions[index]

		if opcode == 99 {
			break
		}

		fstOperator := instructions[index+1]
		sndOperator := instructions[index+2]
		target := instructions[index+3]

		switch opcode {
		case 1:
			instructions[target] = instructions[fstOperator] + instructions[sndOperator]
			break
		case 2:
			instructions[target] = instructions[fstOperator] * instructions[sndOperator]
			break
		default:
			panic("invalid opcode")
		}
		index += 4
	}

	return instructions[0]
}

func parseInput(input string) (ints []int) {
	for _, line := range strings.Split(input, ",") {
		ints = append(ints, cast.ToInt(line))
	}
	return ints
}
