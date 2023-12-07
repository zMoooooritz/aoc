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
	times := cast.ToIntSlice(strings.Trim(strings.Split(parsed[0], ":")[1], " "))
	dists := cast.ToIntSlice(strings.Trim(strings.Split(parsed[1], ":")[1], " "))

	result := 1
	for i := 0; i < len(times); i++ {
		time, dist := times[i], dists[i]
		result *= time - 2*binSearch(time, dist) + 1
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)
	time := cast.ToInt(strings.ReplaceAll(strings.Split(parsed[0], ":")[1], " ", ""))
	dist := cast.ToInt(strings.ReplaceAll(strings.Split(parsed[1], ":")[1], " ", ""))

	return time - 2*binSearch(time, dist) + 1
}

func binSearch(time, dist int) int {
	l, r := 0, time
	for l < r {
		mid := (l + r) / 2
		d := travelDist(mid, time)
		if d <= dist {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return l
}

func travelDist(holdTime, fullTime int) int {
	return (fullTime - holdTime) * holdTime
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
