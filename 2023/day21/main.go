package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"strings"

	"github.com/zMoooooritz/advent-of-code/ds/pq"
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
	START
	GARDEN
	ROCK
)

type State struct {
	Pos image.Point
}

func part1(input string) int {
	tiles := parseInput(input)

	return bfs(tiles, 64)
}

func bfs(tiles [][]Tile, steps int) int {
	grid := map[image.Point]bool{}

	start := State{}
	for y, line := range tiles {
		for x, v := range line {
			grid[image.Point{x, y}] = true
			if v == START {
				start = State{image.Point{x, y}}
			}
		}
	}

	pq := pq.PriorityQueue[State]{}
	seen := make(map[State]int)
	pq.Insert(start, 0)

	for !pq.IsEmpty() {
		state, distance := pq.DeleteMin()

		if _, ok := seen[state]; ok {
			continue
		}
		seen[state] = distance

		pa := image.Point{0, 1}
		pb := image.Point{0, -1}
		pc := image.Point{1, 0}
		pd := image.Point{-1, 0}
		directions := []image.Point{pa, pb, pc, pd}

		for _, dir := range directions {
			newPos := state.Pos.Add(dir)
			if _, ok := grid[newPos]; !ok {
				continue
			}
			tile := tiles[newPos.Y][newPos.X]
			if tile != GARDEN {
				continue
			}
			if !(distance < steps) {
				continue
			}

			pq.Insert(State{newPos}, distance+1)
		}
	}

	reachablePlots := 0
	for _, v := range seen {
		if v%2 == steps%2 {
			reachablePlots++
		}
	}

	return reachablePlots
}

func part2(input string) int {
	return 0
}

func parseInput(input string) (tiles [][]Tile) {
	for _, line := range strings.Split(input, "\n") {
		row := []Tile{}
		for _, c := range line {
			switch c {
			case 'S':
				row = append(row, START)
			case '#':
				row = append(row, ROCK)
			case '.':
				row = append(row, GARDEN)
			default:
				row = append(row, NONE)
			}
		}
		tiles = append(tiles, row)

	}
	return tiles
}
