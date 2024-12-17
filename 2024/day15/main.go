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

func (c *Coordinate) add(d Coordinate) {
	c.x += d.x
	c.y += d.y
}

var dirs = []Coordinate{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type Warehouse struct {
	grid [][]byte

	height int
	width  int

	robotCoord Coordinate
}

func (w *Warehouse) isValidCoord(coord Coordinate) bool {
	return coord.x >= 0 && coord.x < w.width && coord.y >= 0 && coord.y < w.height
}

func (w *Warehouse) doMoveRobot(move Coordinate) {
	invalidCoord := Coordinate{-1, -1}
	curr := w.robotCoord
	free := invalidCoord
	for {
		curr.add(move)
		if !w.isValidCoord(curr) {
			break
		}
		if w.grid[curr.y][curr.x] == '#' {
			break
		}
		if w.grid[curr.y][curr.x] == '.' {
			free = curr
			break
		}
	}

	if free != invalidCoord {
		w.grid[free.y][free.x] = 'O'
		w.grid[w.robotCoord.y][w.robotCoord.x] = '.'
		w.robotCoord.add(move)
		w.grid[w.robotCoord.y][w.robotCoord.x] = '@'
	}
}

func (w *Warehouse) calcCoordSum() int {
	result := 0
	for y, row := range w.grid {
		for x, c := range row {
			if c == 'O' {
				result += 100*y + x
			}
		}
	}
	return result
}

func (w *Warehouse) printGrid() {
	for _, row := range w.grid {
		fmt.Println(string(row))
	}
}

func part1(input string) int {
	warehouse, moves := parseInput(input)

	for _, move := range moves {
		warehouse.doMoveRobot(move)
	}

	warehouse.printGrid()

	return warehouse.calcCoordSum()
}

func part2(input string) int {

	return 0
}

func parseInput(input string) (Warehouse, []Coordinate) {
	w := Warehouse{}

	w.grid = [][]byte{}
	currLine := 0
	for y, line := range strings.Split(input, "\n") {
		currLine = y
		if line == "" {
			break
		}
		lineData := []byte{}
		for x, chr := range line {
			lineData = append(lineData, byte(chr))
			if chr == '@' {
				w.robotCoord = Coordinate{x, y}
			}
		}
		w.grid = append(w.grid, lineData)
	}

	currLine += 1
	moves := []Coordinate{}
	for _, line := range strings.Split(input, "\n")[currLine:] {
		for _, c := range line {
			switch c {
			case '<':
				moves = append(moves, dirs[3])
			case 'v':
				moves = append(moves, dirs[2])
			case '>':
				moves = append(moves, dirs[1])
			case '^':
				moves = append(moves, dirs[0])

			}
		}
	}

	w.height = len(w.grid)
	w.width = len(w.grid[0])

	return w, moves
}
