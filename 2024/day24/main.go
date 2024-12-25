package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
	"sort"
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

type Operation int

const (
	AND Operation = iota
	OR
	XOR
)

type Gate struct {
	fstIn string
	sndIn string
	out   string
	op    Operation
}

var initialWires = map[string]bool{}
var gates = []Gate{}

func part1(input string) int {
	parseInput(input)

	wires := run()

	return wiresToInt(wires, 'z')
}

func run() map[string]bool {
	activeWires := []string{}
	wires := map[string]bool{}
	for k, v := range initialWires {
		activeWires = append(activeWires, k)
		wires[k] = v
	}

	outWires := []string{}
	for _, g := range gates {
		if strings.HasPrefix(g.out, "z") {
			outWires = append(outWires, g.out)
		}
	}

	for {
		if containsAll(activeWires, outWires) {
			break
		}

		for _, g := range gates {
			if slices.Contains(activeWires, g.fstIn) && slices.Contains(activeWires, g.sndIn) {
				out := false
				if g.op == AND {
					out = wires[g.fstIn] && wires[g.sndIn]
				} else if g.op == OR {
					out = wires[g.fstIn] || wires[g.sndIn]
				} else {
					out = wires[g.fstIn] != wires[g.sndIn]
				}
				wires[g.out] = out
				activeWires = append(activeWires, g.out)
			}
		}
	}

	return wires
}

func wiresToInt(wireStates map[string]bool, prefix byte) int {
	wireNames := []string{}
	for n := range wireStates {
		if byte(n[0]) == prefix {
			wireNames = append(wireNames, n)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(wireNames)))
	result := 0
	for _, w := range wireNames {
		val := 0
		if wireStates[w] {
			val = 1
		}
		result = result*2 + val
	}
	return result
}

func containsAll(a, b []string) bool {
	for _, x := range b {
		if !slices.Contains(a, x) {
			return false
		}
	}
	return true
}

func findGateOut(a, b string, op Operation) string {
	for _, g := range gates {
		if g.op != op {
			continue
		}
		if g.fstIn == a && g.sndIn == b || g.fstIn == b && g.sndIn == a {
			return g.out
		}
	}
	return ""
}

func part2(input string) int {
	parseInput(input)

	swapped := []string{}
	carry := ""
	for i := range 45 {
		n := fmt.Sprintf("%02d", i)

		m1 := findGateOut("x"+n, "y"+n, XOR)
		n1 := findGateOut("x"+n, "y"+n, AND)

		var r1, z1, c1 string

		if carry != "" {
			r1 = findGateOut(carry, m1, AND)
			if r1 == "" {
				m1, n1 = n1, m1
				swapped = append(swapped, m1, n1)
				r1 = findGateOut(carry, m1, AND)
			}

			z1 = findGateOut(carry, m1, XOR)

			if strings.HasPrefix(m1, "z") {
				m1, z1 = z1, m1
				swapped = append(swapped, m1, z1)
			}
			if strings.HasPrefix(n1, "z") {
				n1, z1 = z1, n1
				swapped = append(swapped, n1, z1)
			}
			if strings.HasPrefix(r1, "z") {
				r1, z1 = z1, r1
				swapped = append(swapped, r1, z1)
			}

			c1 = findGateOut(r1, n1, OR)
		}

		if strings.HasPrefix(c1, "z") && c1 != "z45" {
			c1, z1 = z1, c1
			swapped = append(swapped, c1, z1)
		}

		if carry == "" {
			carry = n1
		} else {
			carry = c1
		}
	}

	slices.Sort(sort.StringSlice(swapped))
	fmt.Println(strings.Join(swapped, ","))

	return 0
}

func parseInput(input string) {
	initialWires = map[string]bool{}
	gates = []Gate{}

	wireInput := true
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			wireInput = false
			continue
		}
		if wireInput {
			splt := strings.Split(line, ": ")
			initialWires[splt[0]] = splt[1] == "1"
		} else {
			splt := strings.Split(line, " -> ")
			spl := strings.Split(splt[0], " ")
			op := AND
			if spl[1] == "AND" {
				op = AND
			} else if spl[1] == "OR" {
				op = OR
			} else {
				op = XOR
			}
			gates = append(gates, Gate{spl[0], spl[2], splt[1], op})
		}
	}
}
