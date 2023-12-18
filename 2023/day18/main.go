package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/zMoooooritz/advent-of-code/cast"
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

type Instruction struct {
	dir   Direction
	count int
	color string
}

type Vertex struct {
	X int
	Y int
}

func part1(input string) int {
	ins := parseInput(input)

	sX, sY, mX, mY := calculateDims(ins)

	grid := make([][]byte, mY+1)
	for i := range grid {
		grid[i] = make([]byte, mX+1)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	drawPath(grid, ins, sX, sY)

	fillInterior(grid, sX+1, sY+1)

	return calculateArea(grid)
}

func calculateDims(ins []Instruction) (int, int, int, int) {
	minX, maxX, minY, maxY, currX, currY := 0, 0, 0, 0, 0, 0
	for _, i := range ins {
		switch i.dir {
		case NORTH:
			currY -= i.count
		case EAST:
			currX += i.count
		case SOUTH:
			currY += i.count
		case WEST:
			currX -= i.count
		}
		minX = maths.MinInt(minX, currX)
		maxX = maths.MaxInt(maxX, currX)
		minY = maths.MinInt(minY, currY)
		maxY = maths.MaxInt(maxY, currY)
	}
	return -minX, -minY, maxX - minX, maxY - minY
}

func drawPath(grid [][]byte, ins []Instruction, sX, sY int) {
	currX, currY := sX, sY
	grid[currY][currX] = '#'
	for _, i := range ins {
		dX, dY := directionToOffsets(i.dir)
		for j := 0; j < i.count; j++ {
			currX += dX
			currY += dY
			grid[currY][currX] = '#'
		}
	}
}

func fillInterior(grid [][]byte, x, y int) {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[0]) || grid[y][x] != '.' {
		return
	}
	grid[y][x] = 'X'

	fillInterior(grid, x+1, y)
	fillInterior(grid, x-1, y)
	fillInterior(grid, x, y+1)
	fillInterior(grid, x, y-1)
}

func calculateArea(grid [][]byte) int {
	area := 0
	for _, line := range grid {
		for _, c := range line {
			if c == '#' || c == 'X' {
				area++
			}
		}
	}
	return area
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

func part2(input string) int {
	ins := parseInput(input)
	fixInstructions(ins)
	return areaCalculation(ins)
}

func fixInstructions(ins []Instruction) {
	for i := range ins {
		count, _ := strconv.ParseInt(ins[i].color[1:6], 16, 32)
		ins[i].count = int(count)
		switch ins[i].color[6] {
		case '0':
			ins[i].dir = EAST
		case '1':
			ins[i].dir = SOUTH
		case '2':
			ins[i].dir = WEST
		case '3':
			ins[i].dir = NORTH
		}
	}
}

func areaCalculation(ins []Instruction) int {
	vertices := instructionsToVertices(ins)

	// shoelace formula
	area := 0
	prevIndex := len(vertices) - 1
	for currIndex := range vertices {
		area += vertices[prevIndex].X*vertices[currIndex].Y - vertices[prevIndex].Y*vertices[currIndex].X
		prevIndex = currIndex
	}

	// count the perimeter accordingly
	border := 0
	for i := range ins {
		border += ins[i].count
	}

	return area/2 + border/2 + 1
}

func instructionsToVertices(ins []Instruction) []Vertex {
	vertices := []Vertex{}
	currX, currY := 0, 0
	for _, in := range ins {
		dX, dY := directionToOffsets(in.dir)
		currX += dX * in.count
		currY += dY * in.count
		vertices = append(vertices, Vertex{currX, currY})
	}
	return vertices
}

func parseInput(input string) (ans []Instruction) {
	for _, line := range strings.Split(input, "\n") {
		data := strings.Split(line, " ")
		dir := NORTH
		switch data[0] {
		case "U":
			dir = NORTH
		case "R":
			dir = EAST
		case "D":
			dir = SOUTH
		case "L":
			dir = WEST
		}
		count := cast.ToInt(data[1])
		color := strings.Trim(strings.Trim(data[2], "("), ")")
		ans = append(ans, Instruction{dir, count, color})
	}
	return ans
}
