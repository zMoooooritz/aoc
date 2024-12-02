package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/zMoooooritz/advent-of-code/cast"
	"github.com/zMoooooritz/advent-of-code/maths"
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

func isSafe(vals []int) bool {
	if len(vals) <= 1 {
		return true
	}

	incrementing := vals[0] < vals[1]
	isSafe := true

	for index := range len(vals) - 1 {
		if incrementing != (vals[index] < vals[index+1]) {
			isSafe = false
			break
		}
		step := maths.AbsInt(vals[index] - vals[index+1])
		if step < 1 || step > 3 {
			isSafe = false
			break
		}
	}
	return isSafe
}

func isAnySave(vals []int) bool {
	if isSafe(vals) {
		return true
	}

	for index := range vals {
		modifiedVals := make([]int, 0, len(vals)-1)
		modifiedVals = append(modifiedVals, vals[:index]...)
		modifiedVals = append(modifiedVals, vals[index+1:]...)
		if isSafe(modifiedVals) {
			return true
		}
	}
	return false
}

func part1(input string) int {
	parsed := parseInput(input)
	safeCount := 0

	for _, report := range parsed {
		if isSafe(report) {
			safeCount += 1
		}
	}

	return safeCount
}

func part2(input string) int {
	parsed := parseInput(input)
	safeCount := 0

	for _, report := range parsed {
		if isAnySave(report) {
			safeCount += 1
		}
	}

	return safeCount
}

func parseInput(input string) (parsed [][]int) {
	for _, line := range strings.Split(input, "\n") {
		row := []int{}
		for _, element := range strings.Split(line, " ") {
			row = append(row, cast.ToInt(element))
		}
		parsed = append(parsed, row)
	}
	return parsed
}
