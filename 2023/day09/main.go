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
	parsed := parseInput(input)

	result := 0
	for _, line := range parsed {
		diffs := createDiffs(cast.ToIntSlice(line))
		appendPreds(diffs)
		result += diffs[0][len(diffs[0])-1]
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)

	result := 0
	for _, line := range parsed {
		diffs := createDiffs(cast.ToIntSlice(line))
		prependPreds(diffs)
		result += diffs[0][0]
	}

	return result
}

func appendPreds(diffs [][]int) {
	diffs[len(diffs)-1] = append(diffs[len(diffs)-1], 0)
	for i := len(diffs) - 2; i >= 0; i-- {
		curr := diffs[i]
		prev := diffs[i+1]
		diffs[i] = append(curr, curr[len(curr)-1]+prev[len(prev)-1])
	}
}

func prependPreds(diffs [][]int) {
	diffs[len(diffs)-1] = append(diffs[len(diffs)-1], 0)
	for i := len(diffs) - 2; i >= 0; i-- {
		newDiff := []int{
			diffs[i][0] - diffs[i+1][0],
		}
		newDiff = append(newDiff, diffs[i]...)
		diffs[i] = newDiff
	}
}

func createDiffs(series []int) [][]int {
	diffs := [][]int{}
	diffs = append(diffs, series)
	for {
		lastDiff := diffs[len(diffs)-1]
		if allZero(lastDiff) {
			return diffs
		}

		diff := []int{}
		for i := 0; i < len(lastDiff)-1; i++ {
			diff = append(diff, lastDiff[i+1]-lastDiff[i])
		}
		diffs = append(diffs, diff)
	}
}

func allZero(data []int) bool {
	for _, d := range data {
		if d != 0 {
			return false
		}
	}
	return true
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
