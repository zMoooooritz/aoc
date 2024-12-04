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

type Direction int

const (
	N Direction = iota
	NE
	E
	SE
	S
	SW
	W
	NW
)

var directions = []Direction{
	N, NE, E, SE, S, SW, W, NW,
}

var diagDirections = []Direction{
	NE, SE, SW, NW,
}

type Coordinate struct {
	x int
	y int
}

type StepDirection struct {
	x int
	y int
}

func dirToStep(dir Direction) StepDirection {
	stepMap := map[Direction]StepDirection{
		N:  {0, 1},
		NE: {1, 1},
		E:  {1, 0},
		SE: {1, -1},
		S:  {0, -1},
		SW: {-1, -1},
		W:  {-1, 0},
		NW: {-1, 1},
	}
	sd, ok := stepMap[dir]
	if ok {
		return sd
	}
	return StepDirection{0, 0}
}

func toTargetPos(start Coordinate, step StepDirection, steps int) Coordinate {
	return Coordinate{
		start.x + steps*step.x,
		start.y - steps*step.y,
	}
}

func searchInGrid(grid [][]byte, searchTerm string, start Coordinate, direction Direction) bool {
	stepDirection := dirToStep(direction)
	endPos := toTargetPos(start, stepDirection, len(searchTerm)-1)
	gridHeight := len(grid)
	gridWidth := len(grid[0])

	if endPos.x < 0 || endPos.y < 0 || endPos.x >= gridWidth || endPos.y >= gridHeight {
		return false
	}

	for index, char := range searchTerm {
		pos := toTargetPos(start, stepDirection, index)
		if byte(char) != grid[pos.y][pos.x] {
			return false
		}
	}
	return true
}

func part1(input string) int {
	grid := parseInput(input)
	searchTerm := "XMAS"

	result := 0
	for y := range grid {
		for x := range grid {
			for _, direction := range directions {
				found := searchInGrid(grid, searchTerm, Coordinate{x, y}, direction)
				if found {
					result += 1
				}
			}
		}
	}

	return result
}

func part2(input string) int {
	grid := parseInput(input)
	searchTerm := "MAS"

	result := 0
	foundCenters := map[Coordinate]StepDirection{}
	ignoreCenters := map[Coordinate]struct{}{}
	for y := range grid {
		for x := range grid {
			for _, direction := range diagDirections {
				start := Coordinate{x, y}
				found := searchInGrid(grid, searchTerm, Coordinate{x, y}, direction)
				if found {
					stepDir := dirToStep(direction)
					centerCoord := toTargetPos(start, stepDir, 1)
					if prevStepDir, ok := foundCenters[centerCoord]; ok {
						if prevStepDir.x != -stepDir.x || prevStepDir.y != -stepDir.y {
							if _, ok := ignoreCenters[centerCoord]; !ok {
								result += 1
								ignoreCenters[centerCoord] = struct{}{}
							}
						}
					} else {
						foundCenters[centerCoord] = stepDir
					}
				}
			}
		}
	}

	return result
}

func parseInput(input string) (parsed [][]byte) {
	for _, line := range strings.Split(input, "\n") {
		parsed = append(parsed, []byte(line))
	}
	return parsed
}
