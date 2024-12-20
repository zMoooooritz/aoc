package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/zMoooooritz/advent-of-code/ds/pq"
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

var grid = [][]byte{}
var width = 0
var height = 0

func isValidCoord(coord spcl.Coordinate) bool {
	return coord.X >= 0 && coord.X < width && coord.Y >= 0 && coord.Y < height
}

type State struct {
	Pos spcl.Coordinate
	Dir spcl.Vector
}

func dijkstra(start, end spcl.Coordinate) int {
	pq := pq.PriorityQueue[State]{}
	seen := map[State]struct{}{}

	pq.Insert(State{start, spcl.Vector{X: 1, Y: 0}}, 0)

	for !pq.IsEmpty() {
		state, distance := pq.DeleteMin()

		if state.Pos == end {
			return distance
		}

		if _, ok := seen[state]; ok {
			continue
		}
		seen[state] = struct{}{}

		straight := state.Pos
		straight.Add(state.Dir)
		if isValidCoord(straight) && grid[straight.Y][straight.X] != '#' {
			pq.Insert(State{straight, state.Dir}, distance+1)
		}

		cwDir := state.Dir
		cwDir.RotateCW()
		cw := state.Pos
		cw.Add(cwDir)
		if isValidCoord(cw) && grid[cw.Y][cw.X] != '#' {
			pq.Insert(State{cw, cwDir}, distance+1001)
		}

		ccwDir := state.Dir
		ccwDir.RotateCCW()
		ccw := state.Pos
		ccw.Add(ccwDir)
		if isValidCoord(ccw) && grid[ccw.Y][ccw.X] != '#' {
			pq.Insert(State{ccw, ccwDir}, distance+1001)
		}
	}
	return -1
}

func dijkstraExpanded(start, end spcl.Coordinate) (int, map[State][]State, []spcl.Vector) {
	pq := pq.PriorityQueue[State]{}
	seen := map[State]struct{}{}

	minDist := 999999999
	endDirs := []spcl.Vector{}
	cameFrom := map[State][]State{}

	minStateCost := map[State]int{}

	startState := State{start, spcl.Vector{X: 1, Y: 0}}
	minStateCost[startState] = 0
	pq.Insert(startState, 0)

	for !pq.IsEmpty() {
		state, distance := pq.DeleteMin()

		if state.Pos == end {
			if distance > minDist {
				break
			}
			minDist = distance
			endDirs = append(endDirs, state.Dir)
		}

		if _, ok := seen[state]; ok {
			continue
		}
		seen[state] = struct{}{}

		straight := state.Pos
		straight.Add(state.Dir)
		if isValidCoord(straight) && grid[straight.Y][straight.X] != '#' {
			nextState := State{straight, state.Dir}
			nextCost := distance + 1
			cost, ok := minStateCost[nextState]
			if !ok || nextCost <= cost {
				minStateCost[nextState] = nextCost
				cameFrom[nextState] = append(cameFrom[nextState], state)
				pq.Insert(nextState, nextCost)
			}
		}

		cwDir := state.Dir
		cwDir.RotateCW()
		cw := state.Pos
		cw.Add(cwDir)
		if isValidCoord(cw) && grid[cw.Y][cw.X] != '#' {
			nextState := State{cw, cwDir}
			nextCost := distance + 1001
			cost, ok := minStateCost[nextState]
			if !ok || nextCost <= cost {
				minStateCost[nextState] = nextCost
				cameFrom[nextState] = append(cameFrom[nextState], state)
				pq.Insert(nextState, nextCost)
			}
		}

		ccwDir := state.Dir
		ccwDir.RotateCCW()
		ccw := state.Pos
		ccw.Add(ccwDir)
		if isValidCoord(ccw) && grid[ccw.Y][ccw.X] != '#' {
			nextState := State{ccw, ccwDir}
			nextCost := distance + 1001
			cost, ok := minStateCost[nextState]
			if !ok || nextCost <= cost {
				minStateCost[nextState] = nextCost
				cameFrom[nextState] = append(cameFrom[nextState], state)
				pq.Insert(nextState, nextCost)
			}
		}
	}
	return minDist, cameFrom, endDirs
}

func countShortestPathTiles(cameFrom map[State][]State, end spcl.Coordinate, endDirs []spcl.Vector) int {
	queue := []map[State][]spcl.Coordinate{}
	paths := [][]spcl.Coordinate{}

	for _, endDir := range endDirs {
		pathMap := map[State][]spcl.Coordinate{
			{end, endDir}: {end},
		}
		queue = append(queue, pathMap)
	}

	for len(queue) > 0 {
		var state map[State][]spcl.Coordinate
		state, queue = queue[0], queue[1:]
		var currState State
		var path []spcl.Coordinate
		for s, p := range state {
			currState = s
			path = p
			break
		}

		if _, ok := cameFrom[currState]; !ok {
			paths = append(paths, path)
		}

		prevs := cameFrom[currState]
		for _, prev := range prevs {
			pathMap := map[State][]spcl.Coordinate{
				prev: append([]spcl.Coordinate{prev.Pos}, path...),
			}
			queue = append(queue, pathMap)
		}
	}

	visited := map[spcl.Coordinate]struct{}{}
	for _, path := range paths {
		for _, node := range path {
			visited[node] = struct{}{}
		}
	}

	return len(visited)
}

func part1(input string) int {
	start, end := parseInput(input)
	_, _ = start, end

	dist := dijkstra(start, end)

	return dist
}

func part2(input string) int {
	start, end := parseInput(input)
	_, _ = start, end

	_, cameFrom, endDirs := dijkstraExpanded(start, end)

	tiles := countShortestPathTiles(cameFrom, end, endDirs)

	return tiles
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
