package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
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

type Equation struct {
	a int
	b int
	c int
}

type EquationSystem struct {
	fst Equation
	snd Equation
}

func (es *EquationSystem) solve() int {
	detD := es.fst.a*es.snd.b - es.fst.b*es.snd.a
	detDx := es.fst.c*es.snd.b - es.fst.b*es.snd.c
	detDy := es.fst.a*es.snd.c - es.fst.c*es.snd.a

	if detD != 0 {
		x := float64(detDx) / float64(detD)
		y := float64(detDy) / float64(detD)
		if x == float64(int(x)) && y == float64(int(y)) {
			return 3*int(x) + int(y)
		}
	}
	return 0
}

func part1(input string) int {
	systems := parseInput(input)

	result := 0
	for _, system := range systems {
		result += system.solve()
	}

	return result
}

func part2(input string) int {
	systems := parseInput(input)

	offset := 10000000000000

	result := 0
	for _, system := range systems {
		system.fst.c += offset
		system.snd.c += offset
		result += system.solve()
	}

	return result
}

func parseEquationSystem(input string) EquationSystem {
	// Split the input string into lines
	lines := strings.Split(input, "\n")

	// Extract button offsets
	buttonA := strings.Split(strings.TrimPrefix(lines[0], "Button A: "), ", ")
	buttonB := strings.Split(strings.TrimPrefix(lines[1], "Button B: "), ", ")

	buttonA_X, _ := strconv.Atoi(strings.TrimPrefix(buttonA[0], "X"))
	buttonA_Y, _ := strconv.Atoi(strings.TrimPrefix(buttonA[1], "Y"))

	buttonB_X, _ := strconv.Atoi(strings.TrimPrefix(buttonB[0], "X"))
	buttonB_Y, _ := strconv.Atoi(strings.TrimPrefix(buttonB[1], "Y"))

	// Extract prize values
	prize := strings.Split(strings.TrimPrefix(lines[2], "Prize: "), ", ")
	prize_X, _ := strconv.Atoi(strings.TrimPrefix(prize[0], "X="))
	prize_Y, _ := strconv.Atoi(strings.TrimPrefix(prize[1], "Y="))

	// Create the equations
	eq1 := Equation{a: buttonA_X, b: buttonB_X, c: prize_X} // X-offset for Button A
	eq2 := Equation{a: buttonA_Y, b: buttonB_Y, c: prize_Y} // Y-offset for Button A

	// Return the system of equations
	return EquationSystem{fst: eq1, snd: eq2}
}

func parseInput(input string) (data []EquationSystem) {
	systems := []EquationSystem{}
	for _, esd := range strings.Split(input, "\n\n") {
		systems = append(systems, parseEquationSystem(esd))
	}
	return systems
}
