package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"math"
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

type State struct {
	Pos image.Point
	Dir image.Point
}

func part1(input string) int {
	parsed := parseInput(input)

	return dijkstra(parsed, 1, 3)
}

func part2(input string) int {
	parsed := parseInput(input)

	return dijkstra(parsed, 4, 10)
}

func dijkstra(data [][]int, minStepsPerDirection, maxStepsPerDirection int) int {
	grid := map[image.Point]int{}
	end := image.Point{0, 0}

	for y, line := range data {
		for x, v := range line {
			grid[image.Point{x, y}] = v
			end = image.Point{x, y}
		}
	}

	pq := pq.PriorityQueue[State]{}
	seen := map[State]bool{}

	pq.Insert(State{image.Point{0, 0}, image.Point{1, 0}}, 0)
	pq.Insert(State{image.Point{0, 0}, image.Point{0, 1}}, 0)

	for !pq.IsEmpty() {
		state, distance := pq.DeleteMin()

		if state.Pos == end {
			return distance
		}

		if _, ok := seen[state]; ok {
			continue
		}
		seen[state] = true

		for i := -maxStepsPerDirection; i <= maxStepsPerDirection; i++ {
			newPos := state.Pos.Add(state.Dir.Mul(i))
			if _, ok := grid[newPos]; !ok || (-minStepsPerDirection < i && i < minStepsPerDirection) {
				continue
			}

			h := 0
			s := int(math.Copysign(1, float64(i)))
			for j := s; j != i+s; j += s {
				h += grid[state.Pos.Add(state.Dir.Mul(j))]
			}
			pq.Insert(State{newPos, image.Point{state.Dir.Y, state.Dir.X}}, distance+h)
		}
	}
	return -1
}

func parseInput(input string) (ans [][]int) {
	for _, line := range strings.Split(input, "\n") {
		ln := []int{}
		for _, c := range line {
			ln = append(ln, int(c-'0'))
		}
		ans = append(ans, ln)
	}
	return ans
}
