package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
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

var cache = make(map[string]int)

func part1(input string) int {
	parsed := parseInput(input)

	result := 0
	for _, line := range parsed {
		data := strings.Split(line, " ")
		springs := data[0]
		groups := cast.ToIntSliceSep(data[1], ",")

		result += recursion(springs, groups)
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)

	result := 0
	for _, line := range parsed {
		data := strings.Split(line, " ")
		springs := repeatTimes(data[0], 5, "?")
		groups := cast.ToIntSliceSep(repeatTimes(data[1], 5, ","), ",")

		result += recursion(springs, groups)
	}

	return result
}

func repeatTimes(str string, count int, sep string) string {
	newStr := ""
	for i := 0; i < count; i++ {
		newStr += str + sep
	}
	return newStr[:len(newStr)-1]
}

func recursion(springs string, groups []int) int {
	key := springs
	for _, group := range groups {
		key += strconv.Itoa(group) + ","
	}
	if v, ok := cache[key]; ok {
		return v
	}

	if len(springs) == 0 {
		if len(groups) == 0 {
			return 1
		} else {
			return 0
		}
	}

	if strings.HasPrefix(springs, "?") {
		return recursion(strings.Replace(springs, "?", ".", 1), groups) +
			recursion(strings.Replace(springs, "?", "#", 1), groups)
	}

	if strings.HasPrefix(springs, ".") {
		res := recursion(strings.TrimPrefix(springs, "."), groups)
		cache[key] = res
		return res
	}

	if strings.HasPrefix(springs, "#") {
		if len(groups) == 0 || len(springs) < groups[0] {
			cache[key] = 0
			return 0
		}
		if strings.Contains(springs[0:groups[0]], ".") {
			cache[key] = 0
			return 0
		}

		if len(groups) > 1 {
			if len(springs) < groups[0]+1 || string(springs[groups[0]]) == "#" {
				cache[key] = 0
				return 0
			}
			res := recursion(springs[groups[0]+1:], groups[1:])
			cache[key] = res
			return res
		} else {
			res := recursion(springs[groups[0]:], groups[1:])
			cache[key] = res
			return res
		}
	}

	return 0
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
