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
	data := padData(parsed)

	start := -1
	result := 0

	for j := 1; j < len(data)-1; j++ {
		for i := 1; i < len(data[0]); i++ {
			if unicode.IsDigit(rune(data[j][i])) && start == -1 {
				start = i
			}

			if !unicode.IsDigit(rune(data[j][i])) && start != -1 {
				if isAdjacendPart(data, start, j, i-start) {
					number, _ := strconv.Atoi(data[j][start:i])
					result += number
				}
				start = -1
			}
		}
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)
	data := padData(parsed)
	result := 0

	for j := 1; j < len(data)-1; j++ {
		for i := 1; i < len(data[0])-1; i++ {
			if data[j][i] == '*' {
				result += calculateGearRatio(data, i, j)
			}
		}
	}
	return result
}

func padData(input []string) []string {
	emptyLine := strings.Repeat(".", len(input[0]))
	data := []string{emptyLine}
	data = append(data, input...)
	data = append(data, emptyLine)

	for j := 0; j < len(data); j++ {
		data[j] = "." + data[j] + "."
	}
	return data
}

func isAdjacendPart(data []string, x, y, count int) bool {
	for j := y - 1; j <= y+1; j++ {
		for i := x - 1; i <= x+count; i++ {
			char := data[j][i]
			if !unicode.IsDigit(rune(char)) && char != '.' {
				return true
			}
		}
	}
	return false
}

func calculateGearRatio(data []string, x, y int) int {
	values := map[int]bool{}

	for j := y - 1; j <= y+1; j++ {
		for i := x - 1; i <= x+1; i++ {
			result := extractNumber(data[j], i)
			values[result] = true
		}
	}
	delete(values, 0)

	ratio := 0
	if len(values) == 2 {
		ratio = 1
		for k := range values {
			ratio *= k
		}
	}

	return ratio
}

func extractNumber(str string, x int) int {
	start := x
	end := x

	for unicode.IsDigit(rune(str[start])) {
		start -= 1
	}
	for unicode.IsDigit(rune(str[end])) {
		end += 1
	}

	if start == end {
		return 0
	}
	val, err := strconv.Atoi(str[start+1 : end])
	if err != nil {
		return 0
	}
	return val
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
