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

type Position struct {
	X int
	Y int
}

func part1(input string) int {
	parsed := parseInput(input)

	data := [][]rune{}
	start := Position{}
	for y, line := range parsed {
		data = append(data, []rune(line))
		x := strings.Index(line, "S")
		if x != -1 {
			start = Position{x, y}
		}
	}

	curr := startPosition(data, start)
	prev := start

	dist := 1
	for {
		tmp := curr
		curr = nextPosition(data, prev, curr)
		prev = tmp
		dist += 1

		if curr == start {
			break
		}
	}

	return dist / 2
}

// disgusting code
// should have invested more time in the beginning parsing the data properly
// to allow working with the directions more easily
func part2(input string) int {
	parsed := parseInput(input)

	data := [][]rune{}
	processed := [][]int{}
	zeroes := []int{}
	start := Position{}
	for range parsed[0] {
		zeroes = append(zeroes, 0)
	}
	for y, line := range parsed {
		data = append(data, []rune(line))
		z := make([]int, len(zeroes))
		copy(z, zeroes)
		processed = append(processed, z)
		x := strings.Index(line, "S")
		if x != -1 {
			start = Position{x, y}
		}
	}

	replaceStart(data, start)
	curr := startPosition(data, start)
	processed[curr.Y][curr.X] = 1
	prev := start

	for {
		tmp := curr
		curr = nextPosition(data, prev, curr)
		prev = tmp
		processed[curr.Y][curr.X] = 1

		if curr == start {
			break
		}
	}

	result := 0
	for y, line := range processed {
		inside := false
		for x, c := range line {

			symbol := data[y][x]
			if c == 1 {
				if symbol == '|' || symbol == 'F' || symbol == '7' {
					inside = !inside
				}
				continue
			}

			if inside {
				result += 1
				processed[y][x] = 2
			}
		}
	}

	return result
}

func startPosition(data [][]rune, start Position) Position {
	if start.X > 0 {
		left := data[start.Y][start.X-1]
		if left == '-' || left == 'L' || left == 'F' {
			return Position{start.X - 1, start.Y}
		}
	}
	if start.X < len(data[0])-1 {
		right := data[start.Y][start.X+1]
		if right == '-' || right == 'J' || right == '7' {
			return Position{start.X + 1, start.Y}
		}
	}
	if start.Y > 0 {
		above := data[start.Y-1][start.X]
		if above == '|' || above == '7' || above == 'F' {
			return Position{start.X, start.Y - 1}
		}
	}
	if start.X < len(data)-1 {
		below := data[start.Y+1][start.X]
		if below == '|' || below == 'L' || below == 'J' {
			return Position{start.X, start.Y + 1}
		}
	}
	return Position{}
}

// does not work if s is at the border
// too lazy to implement this here
func replaceStart(data [][]rune, start Position) {
	above := data[start.Y-1][start.X]
	below := data[start.Y+1][start.X]
	right := data[start.Y][start.X+1]
	left := data[start.Y][start.X-1]

	if above == '|' || above == '7' || above == 'F' {
		if right == '-' || right == 'J' || right == '7' {
			data[start.Y][start.X] = 'L'
		} else if left == '-' || left == 'L' || left == 'F' {
			data[start.Y][start.X] = 'J'
		} else {
			data[start.Y][start.X] = '|'
		}
	} else if below == '|' || below == 'L' || below == 'J' {
		if right == '-' || right == 'J' || right == '7' {
			data[start.Y][start.X] = 'F'
		} else {
			data[start.Y][start.X] = '7'
		}
	} else {
		data[start.Y][start.X] = '-'
	}
}

func nextPosition(data [][]rune, prev, curr Position) Position {
	next := curr
	if data[curr.Y][curr.X] == '|' {
		if prev.Y < curr.Y {
			next.Y = next.Y + 1
		} else {
			next.Y = next.Y - 1
		}
	}
	if data[curr.Y][curr.X] == '-' {
		if prev.X < curr.X {
			next.X = next.X + 1
		} else {
			next.X = next.X - 1
		}
	}
	if data[curr.Y][curr.X] == 'L' {
		if prev.Y < curr.Y {
			next.X = next.X + 1
		} else {
			next.Y = next.Y - 1
		}
	}
	if data[curr.Y][curr.X] == 'J' {
		if prev.Y < curr.Y {
			next.X = next.X - 1
		} else {
			next.Y = next.Y - 1
		}
	}
	if data[curr.Y][curr.X] == '7' {
		if prev.Y > curr.Y {
			next.X = next.X - 1
		} else {
			next.Y = next.Y + 1
		}
	}
	if data[curr.Y][curr.X] == 'F' {
		if prev.Y > curr.Y {
			next.X = next.X + 1
		} else {
			next.Y = next.Y + 1
		}
	}
	return next
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
