package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
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

type Direction int

const (
	N Direction = iota
	E
	S
	W
)

func rotateCW(dir Direction) Direction {
	rotationMap := map[Direction]Direction{
		N: E,
		E: S,
		S: W,
		W: N,
	}
	rd, ok := rotationMap[dir]
	if ok {
		return rd
	}
	panic("invalid access")
}

type Vector struct {
	x int
	y int
}

func directionToVec(dir Direction) Vector {
	vectorMap := map[Direction]Vector{
		N: {0, 1},
		E: {1, 0},
		S: {0, -1},
		W: {-1, 0},
	}
	sd, ok := vectorMap[dir]
	if ok {
		return sd
	}
	panic("invalid access")
}

type Guard struct {
	coord Coordinate
	dir   Direction
}

func (g *Guard) lookingAt() Coordinate {
	vec := directionToVec(g.dir)
	return Coordinate{
		g.coord.x + vec.x,
		g.coord.y - vec.y,
	}
}

func (g *Guard) stepForward() {
	lookingAt := g.lookingAt()
	g.coord.x = lookingAt.x
	g.coord.y = lookingAt.y
}

func (g *Guard) rotate() {
	g.dir = rotateCW(g.dir)
}

func (g *Guard) simulate() (map[Coordinate][]Direction, bool) {
	isLoop := false

	visited := map[Coordinate][]Direction{}
	visited[g.coord] = []Direction{g.dir}

	for {
		focused := g.lookingAt()
		if !inBounds(focused) {
			break
		}
		if grid[focused.y][focused.x] == obstacleSymbol {
			g.rotate()
		} else {
			g.stepForward()
		}

		if prevVisits, ok := visited[g.coord]; ok {
			if slices.Contains(prevVisits, g.dir) {
				isLoop = true
				break
			}
		}

		visited[g.coord] = append(visited[g.coord], g.dir)
	}
	return visited, isLoop
}

func inBounds(coords Coordinate) bool {
	if coords.x < 0 || coords.y < 0 || coords.y >= len(grid) || coords.x >= len(grid[0]) {
		return false
	}
	return true
}

var guardSymbol = byte('^')
var obstacleSymbol = byte('#')
var floorSymbol = byte('.')

var grid [][]byte

func part1(input string) int {
	guard := parseInput(input)

	visited, _ := guard.simulate()
	return len(visited)
}

func part2(input string) int {
	guard := parseInput(input)
	startPos := guard.coord

	ghostGuard := guard
	visited, _ := ghostGuard.simulate()

	result := 0
	for coord := range visited {
		ghostGuard = guard
		if coord != startPos {
			grid[coord.y][coord.x] = obstacleSymbol
			if _, ok := ghostGuard.simulate(); ok {
				result += 1
			}
			grid[coord.y][coord.x] = floorSymbol
		}
	}
	return result
}

func parseInput(input string) Guard {
	grid = [][]byte{}
	for _, line := range strings.Split(input, "\n") {
		grid = append(grid, []byte(line))
	}

	guard := Guard{}
	for y, row := range grid {
		for x, symbol := range row {
			if symbol == guardSymbol {
				guard = Guard{
					Coordinate{x, y},
					N,
				}
			}
		}
	}

	return guard
}
