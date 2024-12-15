package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
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

type Robot struct {
	pos  Coordinate
	velo Coordinate
}

func (r *Robot) step() {
	r.pos.x = (r.pos.x + r.velo.x + width) % width
	r.pos.y = (r.pos.y + r.velo.y + height) % height
}

type Region struct {
	fields []Coordinate
}

func (r *Region) area() int {
	return len(r.fields)
}

var dirs = []Coordinate{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
var width = 0
var height = 0

func part1(input string) int {
	width = 101
	height = 103

	robots := parseInput(input)

	iterations := 100
	for range iterations {
		for index, r := range robots {
			r.step()
			robots[index] = r
		}
	}

	nw, ne, sw, se := 0, 0, 0, 0
	for _, r := range robots {
		if r.pos.x < width/2 && r.pos.y < height/2 {
			nw += 1
		} else if r.pos.x > width/2 && r.pos.y < height/2 {
			ne += 1
		} else if r.pos.x < width/2 && r.pos.y > height/2 {
			sw += 1
		} else if r.pos.x > width/2 && r.pos.y > height/2 {
			se += 1
		}
	}

	return nw * ne * sw * se
}

var grid = [][]byte{}

func buildGrid(robots []Robot) {
	presence := map[Coordinate]int{}

	for _, r := range robots {
		if _, ok := presence[r.pos]; ok {
			presence[r.pos] += 1
		} else {
			presence[r.pos] = 1
		}
	}

	for y := range height {
		for x := range width {
			coord := Coordinate{x, y}
			if _, ok := presence[coord]; ok {
				grid[y][x] = byte(presence[coord] + '0')
			} else {
				grid[y][x] = '.'
			}

		}
	}
}

func neighbours(coord Coordinate) []Coordinate {
	n := []Coordinate{}

	for _, dir := range dirs {
		newCoord := Coordinate{coord.x + dir.x, coord.y + dir.y}
		if newCoord.x >= 0 && newCoord.y >= 0 && newCoord.x < width && newCoord.y < height {
			n = append(n, newCoord)
		}
	}
	return n
}

func buildRegion(start Coordinate) Region {
	seen := map[Coordinate]struct{}{}
	plotId := grid[start.y][start.x]
	region := Region{}

	region.fields = append(region.fields, start)
	activeNodes := []Coordinate{start}
	seen[start] = struct{}{}

	for len(activeNodes) > 0 {
		var node Coordinate
		node, activeNodes = activeNodes[0], activeNodes[1:]

		nodeId := grid[node.y][node.x]
		if nodeId != plotId {
			continue
		} else {
			if !slices.Contains(region.fields, node) {
				region.fields = append(region.fields, node)
			}
		}

		for _, n := range neighbours(node) {
			if _, ok := seen[n]; ok {
				continue
			}

			activeNodes = append(activeNodes, n)
			seen[n] = struct{}{}
		}
	}

	return region
}

func visualizeGrid() {
	for _, row := range grid {
		strRow := []string{}
		for _, c := range row {
			strRow = append(strRow, string(c))
		}
		fmt.Println(strings.Join(strRow, ""))
	}
}

func part2(input string) int {
	width = 101
	height = 103

	grid = make([][]byte, height)
	for h := range height {
		grid[h] = make([]byte, width)
	}

	robots := parseInput(input)

	iterations := 10000
	result := 0
	for iteration := range iterations {
		for index, r := range robots {
			r.step()
			robots[index] = r
		}

		buildGrid(robots)
		foundBigRegion := false
		for x := range width {
			for y := range height {
				if grid[y][x] == '.' {
					continue
				}
				region := buildRegion(Coordinate{x, y})
				if region.area() > 25 {
					foundBigRegion = true
				}
			}
		}

		if foundBigRegion {
			visualizeGrid()
			result = iteration + 1
			break
		}
	}

	return result
}

func parseInput(input string) []Robot {
	robots := []Robot{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimLeft(line, "p=")
		spl := strings.Split(line, " v=")
		spl1 := strings.Split(spl[0], ",")
		spl2 := strings.Split(spl[1], ",")

		robots = append(robots, Robot{
			Coordinate{cast.ToInt(spl1[0]), cast.ToInt(spl1[1])},
			Coordinate{cast.ToInt(spl2[0]), cast.ToInt(spl2[1])},
		})
	}
	return robots
}
