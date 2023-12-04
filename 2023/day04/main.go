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

func part1(input string) int {
	parsed := parseInput(input)
	result := 0
	for _, line := range parsed {
		parts := strings.Split(line, ":")
		// card := parts[0]
		win := cast.ToIntSlice(strings.Split(parts[1], "|")[0])
		have := cast.ToIntSlice(strings.Split(parts[1], "|")[1])
		hits := len(intersect(win, have))
		if hits > 0 {
			result += 1 << (hits - 1)
		}
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)
	cardCounts := []int{}
	for range parsed {
		cardCounts = append(cardCounts, 1)
	}

	for i, line := range parsed {
		parts := strings.Split(line, ":")
		// card := parts[0]
		win := cast.ToIntSlice(strings.Split(parts[1], "|")[0])
		have := cast.ToIntSlice(strings.Split(parts[1], "|")[1])
		hits := len(intersect(win, have))
		for j := 0; j < hits; j++ {
			cardCounts[i+j+1] += cardCounts[i]
		}
	}

	return maths.SumIntSlice(cardCounts)
}

func intersect(l1 []int, l2 []int) []int {
	var result []int

	unique := make(map[int]bool)

	for _, num := range l1 {
		unique[num] = true
	}

	for _, num := range l2 {
		if unique[num] {
			result = append(result, num)
			delete(unique, num)
		}
	}
	return result
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
