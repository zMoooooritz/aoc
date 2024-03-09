package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
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

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

type Instruction struct {
	Direction Direction
	Distance  int
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

func part1(input string) int {
	first, second := parseInput(input)

	firstVisited, secondVisited := visitedCoordinates(first), visitedCoordinates(second)
	intersections := intersect(firstVisited, secondVisited)

	startingPoint := image.Point{0, 0}
	dist := int(^uint(0) >> 1)
	for _, intersection := range intersections {
		if intersection == startingPoint {
			continue
		}
		hamDist := maths.AbsInt(intersection.X) + maths.AbsInt(intersection.Y)
		dist = maths.MinInt(dist, hamDist)
	}
	return dist
}

func part2(input string) int {
	first, second := parseInput(input)

	firstVisited, secondVisited := visitedCoordinates(first), visitedCoordinates(second)
	intersections := intersect(firstVisited, secondVisited)

	dists := make(map[image.Point]int)
	visited := make(map[image.Point]bool)
	for _, intersection := range intersections {
		dists[intersection] = 0
	}

	for i, coord := range firstVisited {
		if contains(intersections, coord) && !visited[coord] {
			dists[coord] += i
			visited[coord] = true
		}
	}

	visited = make(map[image.Point]bool)
	for i, coord := range secondVisited {
		if contains(intersections, coord) && !visited[coord] {
			dists[coord] += i
			visited[coord] = true
		}
	}

	startingPoint := image.Point{0, 0}
	dist := int(^uint(0) >> 1)
	for k, v := range dists {
		if k == startingPoint {
			continue
		}
		dist = maths.MinInt(dist, v)
	}

	return dist
}

func contains(slice []image.Point, point image.Point) bool {
	for _, p := range slice {
		if p == point {
			return true
		}
	}
	return false
}

func visitedCoordinates(instructions []Instruction) []image.Point {
	currentPosition := image.Point{0, 0}
	positions := []image.Point{}
	positions = append(positions, currentPosition)
	for _, inst := range instructions {
		dir := directionToPoint(inst.Direction)
		for dist := 0; dist < inst.Distance; dist++ {
			currentPosition = currentPosition.Add(dir)
			positions = append(positions, currentPosition)
		}
	}
	return positions
}

func intersect(l1 []image.Point, l2 []image.Point) []image.Point {
	var result []image.Point

	unique := make(map[image.Point]bool)

	for _, p := range l1 {
		unique[p] = true
	}

	for _, p := range l2 {
		if unique[p] {
			result = append(result, p)
			delete(unique, p)
		}
	}
	return result
}

func parseInput(input string) (first, second []Instruction) {
	instructions := [][]Instruction{}
	for _, data := range strings.Split(input, "\n") {
		inst := []Instruction{}
		for _, i := range strings.Split(data, ",") {
			dir := toDirection(rune(i[0]))
			dist := cast.ToInt(i[1:])

			inst = append(inst, Instruction{dir, dist})
		}
		instructions = append(instructions, inst)
	}
	return instructions[0], instructions[1]
}

func directionToPoint(dir Direction) image.Point {
	switch dir {
	case NORTH:
		return image.Point{0, -1}
	case EAST:
		return image.Point{1, 0}
	case SOUTH:
		return image.Point{0, 1}
	case WEST:
		return image.Point{-1, 0}
	}
	panic("invalid direction")
}

func toDirection(d rune) Direction {
	switch d {
	case 'U':
		return NORTH
	case 'R':
		return EAST
	case 'D':
		return SOUTH
	case 'L':
		return WEST
	}
	panic("invalid direction")
}
