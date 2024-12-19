package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

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

var options = []string{}

var optionCount = map[string]int{}

func countOptions(pattern string) int {

	if count, ok := optionCount[pattern]; ok {
		return count
	}

	if len(pattern) == 0 {
		return 1
	}
	total := 0
	for _, opt := range options {
		if strings.HasPrefix(pattern, opt) {
			total += countOptions(strings.TrimPrefix(pattern, opt))
		}
	}
	optionCount[pattern] = total
	return total
}

func part1(input string) int {
	targets := parseInput(input)

	result := 0
	for _, t := range targets {
		if countOptions(t) > 0 {
			result += 1
		}
	}

	return result
}

func part2(input string) int {
	targets := parseInput(input)

	result := 0
	for _, t := range targets {
		result += countOptions(t)
	}

	return result
}

func parseInput(input string) []string {
	targets := []string{}

	lines := strings.Split(input, "\n")
	for _, pattern := range strings.Split(lines[0], ",") {
		options = append(options, strings.TrimSpace(pattern))
	}

	for _, line := range lines[2:] {
		targets = append(targets, line)
	}
	return targets
}
