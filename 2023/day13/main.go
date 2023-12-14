package main

import (
	_ "embed"
	"flag"
	"fmt"
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
	blocks := parseInput(input)

	result := 0
	for _, rows := range blocks {
		cols := tanspose(rows)

		result += 100*findReflection(rows) + findReflection(cols)
	}

	return result
}

func part2(input string) int {
	blocks := parseInput(input)

	result := 0
	for _, rows := range blocks {
		cols := tanspose(rows)

		for i := 0; i < len(rows)-1; i++ {
			if countReflectionMistakes(rows, i) == 1 {
				result += 100 * (i + 1)
			}
		}
		for i := 0; i < len(cols)-1; i++ {
			if countReflectionMistakes(cols, i) == 1 {
				result += i + 1
			}
		}
	}

	return result
}

func tanspose(data []string) []string {
	trans := []string{}

	for i := 0; i < len(data[0]); i++ {
		trans = append(trans, "")
	}

	for _, r := range data {
		for j := range trans {
			trans[j] += string(r[j])
		}
	}
	return trans
}

func findReflection(data []string) int {
	for i := 0; i < len(data)-1; i++ {
		isReflection := true
		for j := 0; i-j >= 0 && i+j+1 < len(data); j++ {
			if data[i-j] != data[i+j+1] {
				isReflection = false
			}
		}
		if isReflection {
			return i + 1
		}
	}
	return 0
}

func countReflectionMistakes(data []string, index int) int {
	count := 0
	for j := 0; index-j >= 0 && index+j+1 < len(data); j++ {
		line := data[index-j]
		reflection := data[index+j+1]
		for k := 0; k < len(line); k++ {
			if line[k] != reflection[k] {
				count++
			}
		}
	}
	return count
}

func switchAtIndex(in string, i int) string {
	if in[i] == '#' {
		return replaceAtIndex(in, '.', i)
	} else {
		return replaceAtIndex(in, '#', i)
	}
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func parseInput(input string) [][]string {
	ans := []string{}
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}

	blocks := [][]string{}
	block := []string{}

	for _, line := range ans {
		if len(line) == 0 {
			blocks = append(blocks, block)
			block = []string{}
		} else {
			block = append(block, line)
		}
	}
	blocks = append(blocks, block)
	return blocks
}
