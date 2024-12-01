package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
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

func part1(input string) int {
	l1, l2 := parseInput(input)

	slices.Sort(l1)
	slices.Sort(l2)

	diff := 0
	for index := range l1 {
		diff += maths.AbsInt(l1[index] - l2[index])
	}

	return diff
}

func part2(input string) int {
	l1, l2 := parseInput(input)

	count := map[int]int{}

	for _, v := range l2 {
		if _, ok := count[v]; ok {
			count[v] += 1
		} else {
			count[v] = 1
		}
	}

	result := 0
	for _, v := range l1 {
		if _, ok := count[v]; ok {
			result += v * count[v]
		}
	}

	return result
}

func parseInput(input string) (l1, l2 []int) {
	for _, line := range strings.Split(input, "\n") {
		vals := strings.Split(line, "   ")
		l1 = append(l1, cast.ToInt(vals[0]))
		l2 = append(l2, cast.ToInt(vals[1]))
	}

	return l1, l2
}
