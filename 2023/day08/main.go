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
	parsed := parseInput(input)

	dirs := dirsToIndex(parsed[0])
	setup := buildSetup(parsed[2:])

	curr := "AAA"
	i := 0
	for true {
		for _, d := range dirs {
			i++
			curr = setup[curr][d]
			if curr == "ZZZ" {
				return i
			}
		}
	}

	return 0
}

func part2(input string) int {
	parsed := parseInput(input)

	dirs := dirsToIndex(parsed[0])
	setup := buildSetup(parsed[2:])

	currs := []string{}
	for s := range setup {
		if strings.HasSuffix(s, "A") {
			currs = append(currs, s)
		}
	}
	counts := []int{}
	for _, curr := range currs {
		i := 0
		done := false
		for !done {
			for _, d := range dirs {
				i++
				curr = setup[curr][d]
				if strings.HasSuffix(curr, "Z") {
					counts = append(counts, i)
					done = true
					break
				}
			}
		}
	}

	return sliceLCM(counts)
}

func sliceLCM(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	result := numbers[0]
	for _, num := range numbers[1:] {
		result = lcm(result, num)
	}
	return result
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func buildSetup(lines []string) map[string][]string {
	m := make(map[string][]string)

	for _, line := range lines {
		dir := strings.Split(line, "=")[1]
		dir = strings.Trim(dir, " ")
		dir = strings.Trim(dir, "(")
		dir = strings.Trim(dir, ")")
		dirs := strings.Split(dir, ", ")
		lr := []string{dirs[0], dirs[1]}
		m[strings.Split(line, " ")[0]] = lr
	}
	return m
}

func dirsToIndex(dirs string) []int {
	indeces := []int{}
	for _, c := range dirs {
		if c == 'L' {
			indeces = append(indeces, 0)
		} else {
			indeces = append(indeces, 1)
		}
	}
	return indeces
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
