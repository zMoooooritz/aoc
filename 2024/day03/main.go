package main

import (
	_ "embed"
	"flag"
	"fmt"
	"regexp"
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

type Operation struct {
	raw       string
	operation string
	operand1  int
	operand2  int
}

func (o *Operation) evaluate() int {
	return o.operand1 * o.operand2
}

func part1(input string) int {
	re := regexp.MustCompile(`(mul)\((\d{1,3}),(\d{1,3})\)`)
	matches := re.FindAllStringSubmatch(input, -1)

	result := 0
	for _, match := range matches {
		operation := Operation{
			match[0],
			match[1],
			cast.ToInt(match[2]),
			cast.ToInt(match[3]),
		}
		result += operation.evaluate()
	}

	return result
}

func part2(input string) int {
	re := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)|do\(\)|don't\(\)`)
	matches := re.FindAllStringSubmatch(input, -1)

	operations := []Operation{}

	for _, match := range matches {
		if strings.HasPrefix(match[0], "mul") {
			mulre := regexp.MustCompile(`(mul)\((\d{1,3}),(\d{1,3})\)`)
			mtch := mulre.FindAllStringSubmatch(match[0], -1)[0]
			operations = append(operations, Operation{
				mtch[0],
				mtch[1],
				cast.ToInt(mtch[2]),
				cast.ToInt(mtch[3]),
			})
		} else if strings.HasPrefix(match[0], "don't") {
			operations = append(operations, Operation{
				match[0],
				"don't",
				0,
				0,
			})
		} else {
			operations = append(operations, Operation{
				match[0],
				"do",
				0,
				0,
			})
		}
	}

	isEnabled := true
	result := 0
	for _, operation := range operations {
		if operation.operation == "do" {
			isEnabled = true
		} else if operation.operation == "don't" {
			isEnabled = false
		} else if operation.operation == "mul" {
			if isEnabled {
				result += operation.evaluate()
			}
		}
	}
	return result
}
