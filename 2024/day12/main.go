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

func (c Coordinate) rotateCW() Coordinate {
	for index, dir := range dirs {
		if dir == c {
			rot := dirs[(index+1)%4]
			return rot
		}
	}
	panic("invalid coordinate")
}

func (c Coordinate) rotateCCW() Coordinate {
	for index, dir := range dirs {
		if dir == c {
			rot := dirs[(index+3)%4]
			return rot
		}
	}
	panic("invalid coordinate")
}

func (c Coordinate) flip() Coordinate {
	return Coordinate{-c.x, -c.y}
}

var dirs = []Coordinate{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
var grid = [][]byte{}
var height = 0
var width = 0

type Region struct {
	plots []Coordinate
}

func (r *Region) area() int {
	return len(r.plots)
}

func (r *Region) perimeter() int {
	per := 0
	for _, p := range r.plots {
		options := 4
		for _, dir := range dirs {
			neighbour := Coordinate{p.x + dir.x, p.y + dir.y}
			if slices.Contains(r.plots, neighbour) {
				options -= 1
			}
		}
		per += options
	}

	return per
}

func (r *Region) lineSweepPerimeterHori(dir Coordinate) int {
	per := 0
	for y := 0; y < height; y += 1 {
		maxX := 0
		minX := 10000
		borderPlots := map[Coordinate]struct{}{}
		for _, p := range r.plots {
			if p.y == y {
				neighbour := Coordinate{p.x + dir.x, p.y + dir.y}
				if !slices.Contains(r.plots, neighbour) {
					maxX = max(maxX, p.x)
					minX = min(minX, p.x)
					borderPlots[p] = struct{}{}
				}
			}
		}
		empty := true
		for x := minX; x < maxX+1; x += 1 {
			if _, ok := borderPlots[Coordinate{x, y}]; ok {
				if empty {
					per += 1
					empty = false
				}
			} else {
				empty = true
			}
		}
	}
	return per
}

func (r *Region) lineSweepPerimeterVert(dir Coordinate) int {
	per := 0
	for x := 0; x < width; x += 1 {
		maxY := 0
		minY := 10000
		borderPlots := map[Coordinate]struct{}{}
		for _, p := range r.plots {
			if p.x == x {
				neighbour := Coordinate{p.x + dir.x, p.y + dir.y}
				if !slices.Contains(r.plots, neighbour) {
					maxY = max(maxY, p.y)
					minY = min(minY, p.y)
					borderPlots[p] = struct{}{}
				}
			}
		}
		empty := true
		for y := minY; y < maxY+1; y += 1 {
			if _, ok := borderPlots[Coordinate{x, y}]; ok {
				if empty {
					per += 1
					empty = false
				}
			} else {
				empty = true
			}
		}
	}
	return per
}

func (r *Region) advancedPerimeter() int {
	per := 0

	per += r.lineSweepPerimeterHori(Coordinate{0, -1})
	per += r.lineSweepPerimeterHori(Coordinate{0, 1})
	per += r.lineSweepPerimeterVert(Coordinate{-1, 0})
	per += r.lineSweepPerimeterVert(Coordinate{1, 0})

	return per
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

	region.plots = append(region.plots, start)
	activeNodes := []Coordinate{start}
	seen[start] = struct{}{}

	for len(activeNodes) > 0 {
		var node Coordinate
		node, activeNodes = activeNodes[0], activeNodes[1:]

		nodeId := grid[node.y][node.x]
		if nodeId != plotId {
			continue
		} else {
			if !slices.Contains(region.plots, node) {
				region.plots = append(region.plots, node)
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

func constructRegions() []Region {
	regions := []Region{}
	seen := map[Coordinate]struct{}{}

	for y, row := range grid {
		for x := range row {
			coord := Coordinate{x, y}
			if _, ok := seen[coord]; ok {
				continue
			}
			region := buildRegion(coord)
			regions = append(regions, region)
			for _, plot := range region.plots {
				seen[plot] = struct{}{}
			}
		}
	}

	return regions
}

func part1(input string) int {
	parseInput(input)

	result := 0
	for _, region := range constructRegions() {
		result += region.area() * region.perimeter()
	}

	return result
}

func part2(input string) int {
	parseInput(input)

	result := 0
	for _, region := range constructRegions() {
		result += region.area() * region.advancedPerimeter()
	}

	return result
}

func parseInput(input string) {
	grid = [][]byte{}
	for _, line := range strings.Split(input, "\n") {
		lineData := []byte{}
		for _, chr := range line {
			lineData = append(lineData, byte(chr))
		}
		grid = append(grid, lineData)
	}

	height = len(grid)
	width = len(grid[0])
}
