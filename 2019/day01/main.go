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
	data := parseInput(input)

	result := 0
	for _, line := range data {
		result += (line / 3) - 2
	}
	return result
}

func part2(input string) int {
	data := parseInput(input)

	result := 0
	for _, line := range data {
		result += requiredFuel(line)
	}
	return result
}

func requiredFuel(mass int) int {
	reqFuel := (mass / 3) - 2
	if reqFuel <= 0 {
		return 0
	}
	return reqFuel + requiredFuel(reqFuel)
}

func parseInput(input string) (data []int) {
	for _, line := range strings.Split(input, "\n") {
		data = append(data, cast.ToInt(line))
	}
	return data
}
