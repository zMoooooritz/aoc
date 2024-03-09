package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
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

func part1(input string) int {
	minVal, maxVal := parseInput(input)
	count := 0
	for val := minVal; val <= maxVal; val++ {
		if isValidPassword(strconv.Itoa(val)) {
			count++
		}
	}

	return count
}

func part2(input string) int {
	minVal, maxVal := parseInput(input)
	count := 0
	for val := minVal; val <= maxVal; val++ {
		if isValidPassword2(strconv.Itoa(val)) {
			count++
		}
	}

	return count
}

func isValidPassword(password string) bool {
	if len(password) != 6 {
		return false
	}

	previous := 0
	hasDuplicate := false
	for _, c := range password {
		val := int(c - '0')
		if val < previous {
			return false
		}
		if val == previous {
			hasDuplicate = true
		}
		previous = val
	}
	return hasDuplicate
}

func isValidPassword2(password string) bool {
	if len(password) != 6 {
		return false
	}

	repeatCount := 0
	previous := 0
	hasDuplicate := false
	for _, c := range password {
		val := int(c - '0')
		if val < previous {
			return false
		}
		if val == previous {
			repeatCount += 1
		} else {
			if repeatCount == 2 {
				hasDuplicate = true
			}
			repeatCount = 1
		}
		previous = val
	}
	if repeatCount == 2 {
		hasDuplicate = true
	}
	return hasDuplicate
}

func parseInput(input string) (int, int) {
	data := strings.Split(input, "-")
	minVal, _ := strconv.Atoi(data[0])
	maxVal, _ := strconv.Atoi(data[1])
	return minVal, maxVal
}
