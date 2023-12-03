package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"unicode"

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
		f, l := firstAndLastDigit(line)
		value, _ := strconv.Atoi(strconv.Itoa(f) + strconv.Itoa(l))
		result += value
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)
	result := 0
	for _, line := range parsed {
		f, l := firstAndLastDigit(digitReplace(line))
		value, _ := strconv.Atoi(strconv.Itoa(f) + strconv.Itoa(l))
		result += value
	}

	return result
}

func digitReplace(str string) string {
	strs := [10]string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	newStr := ""
	for i := 0; i < len(str); i++ {
		for si, s := range strs {
			if strings.HasPrefix(str[i:], s) {
				newStr += strconv.Itoa(si)
				break
			}
		}
		newStr += string(str[i])
	}
	return newStr
}

func firstAndLastDigit(str string) (int, int) {
	first := -1
	last := -1
	for _, char := range str {
		if unicode.IsDigit(char) {
			if first == -1 {
				first = int(char - '0')
				last = first
			} else {
				last = int(char - '0')
			}
		}
	}
	return first, last
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
