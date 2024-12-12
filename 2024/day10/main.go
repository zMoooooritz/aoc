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

type Coordinate struct {
	x int
	y int
}

var grid = [][]int{}
var height = 0
var width = 0

func neighbours(coord Coordinate) []Coordinate {
	n := []Coordinate{}
	dirs := []Coordinate{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for _, dir := range dirs {
		newCoord := Coordinate{coord.x + dir.x, coord.y + dir.y}
		if newCoord.x >= 0 && newCoord.y >= 0 && newCoord.x < width && newCoord.y < height {
			n = append(n, newCoord)
		}
	}
	return n
}

func bfs(start Coordinate, p1 bool) int {
	foundDest := map[Coordinate]int{}
	activeNodes := []Coordinate{start}

	for len(activeNodes) > 0 {
		var node Coordinate
		node, activeNodes = activeNodes[0], activeNodes[1:]

		nodeVal := grid[node.y][node.x]
		if nodeVal == 9 {
			if _, ok := foundDest[node]; ok {
				foundDest[node] += 1
			} else {
				foundDest[node] = 1
			}
			continue
		}
		for _, n := range neighbours(node) {
			newNodeVal := grid[n.y][n.x]
			if nodeVal+1 == newNodeVal {
				activeNodes = append(activeNodes, n)
			}
		}
	}

	rating := 0
	for _, v := range foundDest {
		rating += v
	}
	if p1 {
		return len(foundDest)
	} else {
		return rating

	}
}

func part1(input string) int {
	startCoords := parseInput(input)

	result := 0
	for _, coord := range startCoords {
		result += bfs(coord, true)
	}

	return result
}

func part2(input string) int {
	startCoords := parseInput(input)

	result := 0
	for _, coord := range startCoords {
		result += bfs(coord, false)
	}

	return result
}

func parseInput(input string) []Coordinate {
	grid = [][]int{}
	for _, line := range strings.Split(input, "\n") {
		lineData := []int{}
		for _, chr := range line {
			lineData = append(lineData, cast.ToInt(string(chr)))
		}
		grid = append(grid, lineData)
	}

	height = len(grid)
	width = len(grid[0])

	startCoords := []Coordinate{}
	for y, row := range grid {
		for x, val := range row {
			if val == 0 {
				startCoords = append(startCoords, Coordinate{x, y})
			}
		}
	}

	return startCoords
}
