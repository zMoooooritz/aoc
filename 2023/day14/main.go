package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/zMoooooritz/advent-of-code/util"
	"golang.org/x/exp/slices"
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

type Direction int8

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

func part1(input string) int {
	data := parseInput(input)

	tiltBoard(data, NORTH)

	return calculateForce(data)
}

var states = make(map[string]int)

func part2(input string) int {
	data := parseInput(input)

	loopStart := 0
	loopSize := 0
	cycles := 1000000000
	for i := 0; i < cycles; i++ {
		tiltBoard(data, NORTH)
		tiltBoard(data, WEST)
		tiltBoard(data, SOUTH)
		tiltBoard(data, EAST)

		currState := stringify(data)
		if _, ok := states[currState]; ok {
			loopSize = i - states[currState]
			loopStart = states[currState]
			break
		} else {
			states[currState] = i
		}
	}
	step := (cycles - (loopStart + 1)) % loopSize

	for i := 0; i < step; i++ {
		tiltBoard(data, NORTH)
		tiltBoard(data, WEST)
		tiltBoard(data, SOUTH)
		tiltBoard(data, EAST)
	}

	return calculateForce(data)
}

func tiltBoard(data [][]byte, dir Direction) {
	for i := 0; i < len(data); i++ {
		line := getLineAsString(data, i, dir)
		ln := []byte(line)
		moveStones(ln)
		setLineInData(data, string(ln), i, dir)
	}
}

func calculateForce(data [][]byte) int {
	result := 0
	for y, line := range data {
		for _, char := range line {
			if char == 'O' {
				result += len(data) - y
			}
		}
	}
	return result
}

func getLineAsString(data [][]byte, index int, dir Direction) string {
	switch dir {
	case NORTH:
		return getForwardColAsString(data, index)
	case EAST:
		return getBackwardRowAsString(data, index)
	case SOUTH:
		return getBackwardColAsString(data, index)
	case WEST:
		return getForwardRowAsString(data, index)
	}
	return ""
}

func moveStones(data []byte) {
	for i := 1; i < len(data); i++ {
		if data[i] != 'O' {
			continue
		}
		for j := i; j > 0; j-- {
			if data[j-1] == '.' {
				data[j-1] = 'O'
				data[j] = '.'
			} else {
				break
			}
		}
	}
}

func setLineInData(data [][]byte, line string, index int, dir Direction) {
	switch dir {
	case NORTH:
		setForwardColInData(data, line, index)
	case EAST:
		setBackwardRowInData(data, line, index)
	case SOUTH:
		setBackwardColInData(data, line, index)
	case WEST:
		setForwardRowInData(data, line, index)
	}
}

func stringify(data [][]byte) string {
	bytes := []byte{}
	for _, line := range data {
		bytes = append(bytes, line...)
	}
	return string(bytes)
}

func setForwardRowInData(data [][]byte, row string, index int) {
	data[index] = []byte(row)
}

func setBackwardRowInData(data [][]byte, row string, index int) {
	rw := []byte(row)
	slices.Reverse(rw)
	data[index] = rw
}

func setForwardColInData(data [][]byte, col string, index int) {
	cl := []byte(col)
	for i, char := range cl {
		data[i][index] = char
	}
}

func setBackwardColInData(data [][]byte, col string, index int) {
	cl := []byte(col)
	slices.Reverse(cl)
	for i, char := range cl {
		data[i][index] = char
	}
}

func getForwardRowAsString(data [][]byte, index int) string {
	return string(data[index])
}

func getBackwardRowAsString(data [][]byte, index int) string {
	row := data[index]
	slices.Reverse(row)
	return string(row)
}

func getForwardColAsString(data [][]byte, index int) string {
	col := []byte{}
	for _, row := range data {
		col = append(col, row[index])
	}
	return string(col)
}

func getBackwardColAsString(data [][]byte, index int) string {
	col := []byte{}
	for _, row := range data {
		col = append(col, row[index])
	}
	slices.Reverse(col)
	return string(col)
}

func parseInput(input string) (ans [][]byte) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, []byte(line))
	}

	return ans
}
