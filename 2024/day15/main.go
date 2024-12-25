package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

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

type Warehouse struct {
	grid [][]byte

	height int
	width  int

	robotCoord spcl.Coordinate
}

func (w *Warehouse) getByteAt(coord spcl.Coordinate) byte {
	return w.grid[coord.Y][coord.X]
}

func (w *Warehouse) isValidCoord(coord spcl.Coordinate) bool {
	return coord.X >= 0 && coord.X < w.width && coord.Y >= 0 && coord.Y < w.height
}

func (w *Warehouse) doMoveRobot(move spcl.Vector) {
	invalidCoord := spcl.Coordinate{X: -1, Y: -1}
	curr := w.robotCoord
	free := invalidCoord
	for {
		curr.Add(move)
		if !w.isValidCoord(curr) {
			break
		}
		if w.grid[curr.Y][curr.X] == '#' {
			break
		}
		if w.grid[curr.Y][curr.X] == '.' {
			free = curr
			break
		}
	}

	if free != invalidCoord {
		w.grid[free.Y][free.X] = 'O'
		w.grid[w.robotCoord.Y][w.robotCoord.X] = '.'
		w.robotCoord.Add(move)
		w.grid[w.robotCoord.Y][w.robotCoord.X] = '@'
	}
}

func (w *Warehouse) buildConnectedBlock(coord spcl.Coordinate, dir spcl.Vector) map[spcl.Coordinate]struct{} {
	block := map[spcl.Coordinate]struct{}{}
	if dir.Y == 0 {
		curr := coord
		expected := 'X'
		if dir.X > 0 {
			expected = '['
		} else {
			expected = ']'
		}
		for {
			curr.Add(dir)
			if !w.isValidCoord(curr) {
				break
			}
			if w.getByteAt(curr) == '#' || w.getByteAt(curr) == '.' {
				break
			}
			if w.getByteAt(curr) == byte(expected) {
				if expected == '[' {
					expected = ']'
				} else {
					expected = '['
				}
			} else {
				break
			}
			block[curr] = struct{}{}
		}
	} else {
		curr := coord
		curr.Add(dir)
		if !w.isValidCoord(curr) {
			return block
		}

		if w.grid[curr.Y][curr.X] == '[' {
			other := spcl.Coordinate{X: curr.X + 1, Y: curr.Y}
			block[curr] = struct{}{}
			block[other] = struct{}{}
			for k, v := range w.buildConnectedBlock(curr, dir) {
				block[k] = v
			}
			for k, v := range w.buildConnectedBlock(other, dir) {
				block[k] = v
			}
		} else if w.grid[curr.Y][curr.X] == ']' {
			other := spcl.Coordinate{X: curr.X - 1, Y: curr.Y}
			block[curr] = struct{}{}
			block[other] = struct{}{}
			for k, v := range w.buildConnectedBlock(curr, dir) {
				block[k] = v
			}
			for k, v := range w.buildConnectedBlock(other, dir) {
				block[k] = v
			}
		} else {
			return block
		}
	}
	return block
}

func (w *Warehouse) calcCoordSum() int {
	result := 0
	for y, row := range w.grid {
		for x, c := range row {
			if c == 'O' || c == '[' {
				result += 100*y + x
			}
		}
	}
	return result
}

func (w *Warehouse) scaleUp() {
	newGrid := [][]byte{}

	for _, row := range w.grid {
		newRow := []byte{}
		for _, c := range row {
			switch c {
			case '#':
				newRow = append(newRow, '#', '#')
			case 'O':
				newRow = append(newRow, '[', ']')
			case '.':
				newRow = append(newRow, '.', '.')
			case '@':
				newRow = append(newRow, '@', '.')
			}
		}
		newGrid = append(newGrid, newRow)
	}

	w.grid = newGrid
	w.width *= 2

	w.robotCoord.X *= 2
}

func (w *Warehouse) canMoveBlock(block map[spcl.Coordinate]struct{}, move spcl.Vector) bool {
	isPossible := true
	for k := range block {
		curr := k
		curr.Add(move)
		if w.getByteAt(curr) == '#' {
			isPossible = false
			break
		}
	}
	return isPossible
}

func (w *Warehouse) moveBlock(block map[spcl.Coordinate]struct{}, move spcl.Vector) {
	oldData := map[spcl.Coordinate]byte{}
	for k := range block {
		oldData[spcl.Coordinate{X: k.X, Y: k.Y}] = w.grid[k.Y][k.X]
		w.grid[k.Y][k.X] = '.'
	}

	for k := range block {
		curr := k
		curr.Add(move)
		w.grid[curr.Y][curr.X] = oldData[k]
	}
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

	return warehouse.calcCoordSum()
}

func part2(input string) int {
	warehouse, moves := parseInput(input)
	_ = moves

	warehouse.scaleUp()

	for _, move := range moves {
		block := warehouse.buildConnectedBlock(warehouse.robotCoord, move)

		if len(block) != 0 && warehouse.canMoveBlock(block, move) {
			warehouse.moveBlock(block, move)
		}

		robotCoord := warehouse.robotCoord
		robotCoord.Add(move)
		if warehouse.getByteAt(robotCoord) == '.' {
			warehouse.grid[warehouse.robotCoord.Y][warehouse.robotCoord.X] = '.'
			warehouse.robotCoord.Add(move)
			warehouse.grid[warehouse.robotCoord.Y][warehouse.robotCoord.X] = '@'
		}
	}

	return warehouse.calcCoordSum()
}

func parseInput(input string) (Warehouse, []spcl.Vector) {
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
				w.robotCoord = spcl.Coordinate{X: x, Y: y}
			}
		}
		w.grid = append(w.grid, lineData)
	}

	currLine += 1
	moves := []spcl.Vector{}

	for _, line := range strings.Split(input, "\n")[currLine:] {
		for _, c := range line {
			switch c {
			case '<':
				moves = append(moves, spcl.CARDINAL_DIRS[3])
			case 'v':
				moves = append(moves, spcl.CARDINAL_DIRS[2])
			case '>':
				moves = append(moves, spcl.CARDINAL_DIRS[1])
			case '^':
				moves = append(moves, spcl.CARDINAL_DIRS[0])

			}
		}
	}

	w.height = len(w.grid)
	w.width = len(w.grid[0])

	return w, moves
}
