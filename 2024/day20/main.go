package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/zMoooooritz/advent-of-code/ds/pq"
	"github.com/zMoooooritz/advent-of-code/ds/spcl"
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

var grid = [][]byte{}
var width = 0
var height = 0

func isValidCoord(coord spcl.Coordinate) bool {
	return coord.X >= 0 && coord.X < width && coord.Y >= 0 && coord.Y < height
}

func dijkstra(start spcl.Coordinate) map[spcl.Coordinate]int {
	distMap := map[spcl.Coordinate]int{}

	pq := pq.PriorityQueue[spcl.Coordinate]{}
	seen := map[spcl.Coordinate]struct{}{}

	pq.Insert(start, 0)

	for !pq.IsEmpty() {
		node, distance := pq.DeleteMin()

		if _, ok := seen[node]; ok {
			continue
		}
		seen[node] = struct{}{}
		distMap[node] = distance

		for _, n := range node.CardinalNeighbours() {
			if grid[n.Y][n.X] == '#' {
				continue
			}
			pq.Insert(n, distance+1)
		}

	}
	return distMap
}

func findBasicShortcuts(distMap map[spcl.Coordinate]int) map[int]int {
	cutMap := map[int]int{}
	for coord, dist := range distMap {
		for _, dir := range spcl.CARDINAL_DIRS {
			border, target := coord, coord
			doubleDir := spcl.Coordinate(dir)
			doubleDir.Mul(2)
			border.Add(dir)
			target.Add(spcl.Vector(doubleDir))

			if !isValidCoord(border) || !isValidCoord(target) {
				continue
			}

			if grid[border.Y][border.X] != '#' {
				continue
			}

			if d, ok := distMap[target]; ok {
				if d < dist {
					distDiff := dist - d - 2
					cutMap[distDiff] += 1
				}
			}
		}
	}
	return cutMap
}

func hammingDist(start, end spcl.Coordinate) int {
	return maths.AbsInt(start.X-end.X) + maths.AbsInt(start.Y-end.Y)
}

func jumpOptions(coord spcl.Coordinate) []spcl.Coordinate {
	jumpDist := 20

	coords := []spcl.Coordinate{}
	for x := -jumpDist; x <= jumpDist; x++ {
		for y := -jumpDist; y <= jumpDist; y++ {
			jumpCoord := spcl.Coordinate{X: coord.X + x, Y: coord.Y + y}
			if hammingDist(coord, jumpCoord) > jumpDist {
				continue
			}
			if isValidCoord(jumpCoord) {
				coords = append(coords, jumpCoord)
			}
		}
	}

	return coords
}

func findAdvancedShortcuts(distMap map[spcl.Coordinate]int) map[int]int {
	cutMap := map[int]int{}
	for coord, dist := range distMap {

		targets := jumpOptions(coord)
		for _, target := range targets {
			if d, ok := distMap[target]; ok {
				if d < dist {
					distDiff := dist - d - hammingDist(coord, target)
					cutMap[distDiff] += 1
				}
			}

		}
	}
	return cutMap
}

func part1(input string) int {
	start, end := parseInput(input)

	_ = start

	distMap := dijkstra(end)
	cutMap := findBasicShortcuts(distMap)

	cheatOptions := 0
	for dist, count := range cutMap {
		if dist >= 100 {
			cheatOptions += count
		}
	}

	return cheatOptions
}

func part2(input string) int {
	start, end := parseInput(input)

	_ = start

	distMap := dijkstra(end)
	cutMap := findAdvancedShortcuts(distMap)

	cheatOptions := 0
	for dist, count := range cutMap {
		if dist >= 100 {
			cheatOptions += count
		}
	}

	return cheatOptions
}

func parseInput(input string) (spcl.Coordinate, spcl.Coordinate) {
	grid = [][]byte{}
	for _, row := range strings.Split(input, "\n") {
		grid = append(grid, []byte(row))
	}

	height = len(grid)
	width = len(grid[0])

	start, end := spcl.Coordinate{}, spcl.Coordinate{}
	for y, row := range grid {
		for x, c := range row {
			if c == 'S' {
				start = spcl.Coordinate{X: x, Y: y}
			}
			if c == 'E' {
				end = spcl.Coordinate{X: x, Y: y}
			}
		}
	}
	return start, end
}
