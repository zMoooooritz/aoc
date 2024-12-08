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

type Coordinate struct {
	x int
	y int
}

func coordsToAntiNodesBasic(a, b Coordinate, width, height int) []Coordinate {
	diff := Coordinate{a.x - b.x, a.y - b.y}

	aAnti := Coordinate{a.x + diff.x, a.y + diff.y}
	bAnti := Coordinate{b.x - diff.x, b.y - diff.y}

	antiNodes := []Coordinate{}
	if aAnti.x >= 0 && aAnti.y >= 0 && aAnti.x < width && aAnti.y < height {
		antiNodes = append(antiNodes, aAnti)
	}
	if bAnti.x >= 0 && bAnti.y >= 0 && bAnti.x < width && bAnti.y < height {
		antiNodes = append(antiNodes, bAnti)
	}

	return antiNodes
}

func coordsToAntiNodesAdvanced(a, b Coordinate, width, height int) []Coordinate {
	diff := Coordinate{a.x - b.x, a.y - b.y}

	antiNodes := []Coordinate{}
	for index := 0; ; index += 1 {
		node := Coordinate{a.x + index*diff.x, a.y + index*diff.y}
		if node.x < 0 || node.y < 0 || node.x >= width || node.y >= height {
			break
		}
		antiNodes = append(antiNodes, node)
	}

	for index := 0; ; index += 1 {
		node := Coordinate{b.x - index*diff.x, b.y - index*diff.y}
		if node.x < 0 || node.y < 0 || node.x >= width || node.y >= height {
			break
		}
		antiNodes = append(antiNodes, node)
	}

	return antiNodes
}

func solve(grid [][]byte, antiNodeFunc func(Coordinate, Coordinate, int, int) []Coordinate) int {
	height := len(grid)
	width := len(grid[0])

	frequencyMap := map[byte][]Coordinate{}
	for y, row := range grid {
		for x, c := range row {
			if c > 46 { // 46 == .
				if _, ok := frequencyMap[c]; ok {
					frequencyMap[c] = append(frequencyMap[c], Coordinate{x, y})
				} else {
					frequencyMap[c] = []Coordinate{{x, y}}
				}
			}
		}
	}

	antiNodeMap := map[Coordinate]struct{}{}

	for _, locs := range frequencyMap {
		for i := 0; i < len(locs); i++ {
			for j := i + 1; j < len(locs); j++ {
				antiNodes := antiNodeFunc(locs[i], locs[j], width, height)
				for _, antiNode := range antiNodes {
					antiNodeMap[antiNode] = struct{}{}
				}
			}
		}
	}

	return len(antiNodeMap)
}

func part1(input string) int {
	grid := parseInput(input)
	return solve(grid, coordsToAntiNodesBasic)
}

func part2(input string) int {
	grid := parseInput(input)
	return solve(grid, coordsToAntiNodesAdvanced)
}

func parseInput(input string) [][]byte {
	grid := [][]byte{}
	for _, line := range strings.Split(input, "\n") {
		grid = append(grid, []byte(line))
	}
	return grid
}
