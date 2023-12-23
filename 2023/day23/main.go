package main

import (
	"container/list"
	_ "embed"
	"flag"
	"fmt"
	"image"
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

type Tile int

const (
	NONE Tile = iota
	PATH
	SLOPE_UP
	SLOPE_RIGHT
	SLOPE_DOWN
	SLOPE_LEFT
	FOREST
	START
	END
)

type State struct {
	Pos      image.Point
	Path     map[image.Point]struct{}
	Distance int
}

func part1(input string) int {
	tiles := parseInput(input)
	// fmt.Println(tiles)

	return findLongestPath(tiles)
}

// topologische sortierung
// longest path mit dp

func findLongestPath(tiles [][]Tile) int {
	queue := list.New()
	grid := map[image.Point]bool{}

	var start, end image.Point
	for y, line := range tiles {
		for x, v := range line {
			grid[image.Point{x, y}] = true
			if v == START {
				start = image.Point{x, y}
			}
			if v == END {
				end = image.Point{x, y}
			}
		}
	}

	dists := make(map[image.Point]int)
	dists[start] = 0

	queue.PushBack(State{start, make(map[image.Point]struct{}), 0})

	for queue.Len() > 0 {
		element := queue.Remove(queue.Back()).(State)

		tile := tiles[element.Pos.Y][element.Pos.X]
		if tile == END {
			continue
		}

		pa := image.Point{0, 1}
		pb := image.Point{0, -1}
		pc := image.Point{1, 0}
		pd := image.Point{-1, 0}
		directions := []image.Point{pa, pb, pc, pd}
		if isSlopeTile(tile) {
			directions = []image.Point{slopeToDirection(tile)}
		}

		for _, dir := range directions {
			newPos := element.Pos.Add(dir)
			if _, ok := grid[newPos]; !ok {
				continue
			}
			newTile := tiles[newPos.Y][newPos.X]
			if newTile == FOREST {
				continue
			}

			if _, ok := element.Path[newPos]; ok {
				continue
			}

			newDist := dists[element.Pos] + 1
			if _, ok := dists[newPos]; !ok || newDist > dists[newPos] { // TODO
				dists[newPos] = newDist

				newPath := make(map[image.Point]struct{})
				for k := range element.Path {
					newPath[k] = struct{}{}
				}
				newPath[newPos] = struct{}{}

				queue.PushFront(State{newPos, newPath, newDist})
			}
		}
	}
	return dists[end]
}

func part2(input string) int {
	return 0
}

func isSlopeTile(tile Tile) bool {
	return tile == SLOPE_UP || tile == SLOPE_RIGHT || tile == SLOPE_DOWN || tile == SLOPE_LEFT
}

func slopeToDirection(slope Tile) image.Point {
	switch slope {
	case SLOPE_UP:
		return image.Point{0, -1}
	case SLOPE_RIGHT:
		return image.Point{1, 0}
	case SLOPE_DOWN:
		return image.Point{0, 1}
	case SLOPE_LEFT:
		return image.Point{-1, 0}
	}
	panic("invalid slope tile")
}

func parseInput(input string) [][]Tile {
	tiles := [][]Tile{}
	for _, line := range strings.Split(input, "\n") {
		row := []Tile{}
		for _, c := range line {
			t := NONE
			switch c {
			case '#':
				t = FOREST
			case '.':
				t = PATH
			case '<':
				t = SLOPE_LEFT
			case '^':
				t = SLOPE_UP
			case '>':
				t = SLOPE_RIGHT
			case 'v':
				t = SLOPE_DOWN
			}
			row = append(row, t)
		}
		tiles = append(tiles, row)
	}
	for x, c := range tiles[0] {
		if c == PATH {
			tiles[0][x] = START
			break
		}
	}
	lastLineIndex := len(tiles) - 1
	for x := len(tiles[lastLineIndex]) - 1; x >= 0; x-- {
		if tiles[lastLineIndex][x] == PATH {
			tiles[lastLineIndex][x] = END
			break
		}
	}
	return tiles
}
