package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

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

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

func part1(input string) int {
	return startTrace(parseInput(input), -1, 0, EAST)
}

func part2(input string) int {
	parsed := parseInput(input)
	maxVal := 0
	for i := range parsed {
		r, l := startTrace(parsed, -1, i, EAST), startTrace(parsed, len(parsed[0]), i, WEST)
		maxVal = maths.MaxInt(maxVal, r, l)
	}
	for i := range parsed[0] {
		d, u := startTrace(parsed, i, -1, SOUTH), startTrace(parsed, i, len(parsed), NORTH)
		maxVal = maths.MaxInt(maxVal, d, u)
	}
	return maxVal
}

func startTrace(parsed [][]byte, x, y int, dir Direction) int {
	energized := make([][]byte, len(parsed))
	for i := range energized {
		energized[i] = make([]byte, len(parsed[0]))
		for j := range energized[i] {
			energized[i][j] = ' '
		}
	}
	directions := make([][][]Direction, len(parsed))
	for i := range directions {
		directions[i] = make([][]Direction, len(parsed[0]))
		for j := range directions[i] {
			directions[i][j] = []Direction{}
		}
	}

	traceRay(parsed, energized, directions, x, y, dir)

	return countEnergizedCells(energized)
}

func countEnergizedCells(energized [][]byte) int {
	result := 0
	for _, line := range energized {
		for _, c := range line {
			if c == '#' {
				result++
			}
		}
	}
	return result
}

func traceRay(data [][]byte, energized [][]byte, directions [][][]Direction, x, y int, dir Direction) {
	xo, yo := directionToOffsets(dir)
	x += xo
	y += yo

	if x < 0 || x > len(data[0])-1 || y < 0 || y > len(data)-1 {
		return
	}
	for _, d := range directions[y][x] {
		if dir == d {
			return
		}
	}

	directions[y][x] = append(directions[y][x], dir)
	energized[y][x] = '#'
	dirs := nextDirections(dir, data[y][x])

	for _, d := range dirs {
		traceRay(data, energized, directions, x, y, d)
	}
}

func nextDirections(dir Direction, tile byte) []Direction {
	if tile == '.' {
		return []Direction{dir}
	}

	switch dir {
	case NORTH:
		switch tile {
		case '|':
			return []Direction{NORTH}
		case '-':
			return []Direction{EAST, WEST}
		case '/':
			return []Direction{EAST}
		case '\\':
			return []Direction{WEST}
		}
	case SOUTH:
		switch tile {
		case '|':
			return []Direction{SOUTH}
		case '-':
			return []Direction{EAST, WEST}
		case '/':
			return []Direction{WEST}
		case '\\':
			return []Direction{EAST}
		}
	case EAST:
		switch tile {
		case '|':
			return []Direction{NORTH, SOUTH}
		case '-':
			return []Direction{EAST}
		case '/':
			return []Direction{NORTH}
		case '\\':
			return []Direction{SOUTH}
		}
	case WEST:
		switch tile {
		case '|':
			return []Direction{NORTH, SOUTH}
		case '-':
			return []Direction{WEST}
		case '/':
			return []Direction{SOUTH}
		case '\\':
			return []Direction{NORTH}
		}
	}
	return []Direction{}
}

func directionToOffsets(dir Direction) (int, int) {
	switch dir {
	case NORTH:
		return 0, -1
	case EAST:
		return 1, 0
	case SOUTH:
		return 0, 1
	case WEST:
		return -1, 0
	}
	return 0, 0
}

func parseInput(input string) (ans [][]byte) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, []byte(line))
	}
	return ans
}
