package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"strconv"
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

type Equation struct {
	value    int
	operands []int
}

func part1(input string) int {
	equations := parseInput(input)

	result := 0
	for _, equation := range equations {
		operatorCount := len(equation.operands) - 1
		combinations := int(math.Pow(2, float64(operatorCount)))
		for c := range combinations {
			value := equation.operands[0]
			for index, o := range equation.operands[1:] {
				v := c
				v >>= index
				v %= 2

				if v == 0 {
					value += o
				} else {
					value *= o
				}
			}
			if value == equation.value {
				result += value
				break
			}
		}
	}

	return result
}

func concatInts(a, b int) int {
	aStr := strconv.Itoa(a)
	bStr := strconv.Itoa(b)

	concatenated := aStr + bStr

	result, _ := strconv.Atoi(concatenated)

	return result
}

func part2(input string) int {
	equations := parseInput(input)

	result := 0
	for _, equation := range equations {
		operatorCount := len(equation.operands) - 1
		combinations := int(math.Pow(3, float64(operatorCount)))
		for c := range combinations {
			value := equation.operands[0]
			for index, o := range equation.operands[1:] {
				v := c
				if index > 0 {
					v /= int(math.Pow(3, float64(index)))
				}
				v %= 3

				if v == 0 {
					value += o
				} else if v == 1 {
					value *= o
				} else {
					value = concatInts(value, o)
				}
			}
			if value == equation.value {
				result += value
				break
			}
		}
	}

	return result
}

func parseInput(input string) []Equation {
	equations := []Equation{}
	for _, line := range strings.Split(input, "\n") {
		s := strings.Split(line, ":")
		equations = append(equations, Equation{
			cast.ToInt(s[0]),
			cast.ToIntSlice(s[1]),
		})
	}
	return equations
}
