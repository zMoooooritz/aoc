package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
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

func predsAsList(node string) []string {
	preds := []string{}
	if pred, ok := predList[node]; ok {
		preds = append(preds, pred)
		preds = append(preds, predsAsList(pred)...)
		return preds
	}
	return preds

}

var predList = map[string]string{}

func part1(input string) int {
	parseInput(input)

	result := 0
	for k := range predList {
		result += len(predsAsList(k))
	}
	return result
}

func part2(input string) int {
	parseInput(input)

	youPreds := predsAsList("YOU")
	slices.Reverse(youPreds)
	sanPreds := predsAsList("SAN")
	slices.Reverse(sanPreds)

	for len(youPreds) > 0 && len(sanPreds) > 0 && youPreds[0] == sanPreds[0] {
		youPreds = youPreds[1:]
		sanPreds = sanPreds[1:]
	}

	return len(youPreds) + len(sanPreds)
}

func parseInput(input string) map[string]string {
	predList = map[string]string{}
	for _, line := range strings.Split(input, "\n") {
		splt := strings.Split(line, ")")
		predList[splt[1]] = splt[0]
	}
	return predList
}
