package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
	"strings"

	"github.com/zMoooooritz/advent-of-code/cast"
	"github.com/zMoooooritz/advent-of-code/ds/spcl"
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

var grid [][]byte
var width int
var height int

func isValidCoord(coord spcl.Coordinate) bool {
	return coord.X >= 0 && coord.X < width && coord.Y >= 0 && coord.Y < height
}

func dijkstra(start, end spcl.Coordinate) int {
	next := []spcl.Coordinate{start}

	dist := map[spcl.Coordinate]int{}

	currDist := 0
	for len(next) > 0 {
		active := make([]spcl.Coordinate, len(next))
		copy(active, next)
		next = []spcl.Coordinate{}

		for len(active) > 0 {
			var node spcl.Coordinate
			node, active = active[0], active[1:]
			dist[node] = currDist

			neighbours := node.CardinalNeighbours()
			for _, n := range neighbours {
				if !isValidCoord(n) {
					continue
				}
				if _, ok := dist[n]; ok {
					continue
				}
				if grid[n.Y][n.X] == '#' {
					continue
				}
				if !slices.Contains(next, n) {
					next = append(next, n)
				}

			}
		}

		currDist += 1
	}

	if d, ok := dist[end]; ok {
		return d
	}
	return -1
}

func part1(input string) int {
	coords := parseInput(input)

	width = 71
	height = 71

	drawCount := 1024

	grid = make([][]byte, height)
	for i := range height {
		grid[i] = slices.Repeat([]byte{'.'}, width)
	}

	for index, coord := range coords {
		if index < drawCount {
			grid[coord.Y][coord.X] = '#'
		} else {
			break
		}
	}

	start := spcl.Coordinate{X: 0, Y: 0}
	end := spcl.Coordinate{X: width - 1, Y: height - 1}
	dist := dijkstra(start, end)

	return dist
}

var coords []spcl.Coordinate

func binSearch(start, end int) int {
	mid := (start + end) / 2

	grid = make([][]byte, height)
	for i := range height {
		grid[i] = slices.Repeat([]byte{'.'}, width)
	}

	for index, coord := range coords {
		if index < mid {
			grid[coord.Y][coord.X] = '#'
		} else {
			break
		}
	}
	startCoord := spcl.Coordinate{X: 0, Y: 0}
	endCoord := spcl.Coordinate{X: width - 1, Y: height - 1}
	midResult := dijkstra(startCoord, endCoord)

	if start == end-1 {
		return mid
	} else if midResult == -1 {
		return binSearch(start, mid-1)
	} else {
		return binSearch(mid, end)
	}
}

func part2(input string) int {
	coords = parseInput(input)

	width = 71
	height = 71

	breaking := binSearch(0, len(coords)-1)
	fmt.Println(coords[breaking])

	return 0
}

func parseInput(input string) []spcl.Coordinate {
	coords := []spcl.Coordinate{}
	for _, line := range strings.Split(input, "\n") {
		spl := strings.Split(line, ",")
		coords = append(coords, spcl.Coordinate{X: cast.ToInt(spl[0]), Y: cast.ToInt(spl[1])})
	}
	return coords
}
