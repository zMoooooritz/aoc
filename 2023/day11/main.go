package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

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

type Galaxy struct {
	Index int
	X     int
	Y     int
}

func part1(input string) int {
	return calculateDistances(parseInput(input), 2)
}

func part2(input string) int {
	return calculateDistances(parseInput(input), 1000000)
}

func calculateDistances(data []string, expansionTime int) int {
	galaxies, spaceCols, spaceRows := parseGalaxies(data)

	result := 0
	for _, g1 := range galaxies {
		for _, g2 := range galaxies {
			if g2.Index <= g1.Index {
				continue
			}
			result += SpaceHamDistance(g1, g2, spaceRows, spaceCols, expansionTime)
		}
	}
	return result
}

func parseGalaxies(data []string) ([]Galaxy, []int, []int) {
	spaceRows := []int{}
	emptyCols := []bool{}
	spaceCols := []int{}
	galaxies := []Galaxy{}
	galaxyCounter := 0

	for range data[0] {
		emptyCols = append(emptyCols, true)
	}

	for y, line := range data {
		emptyRow := true
		for x, char := range line {
			if char == '#' {
				galaxies = append(galaxies, Galaxy{galaxyCounter, x, y})
				emptyCols[x] = false
				emptyRow = false
				galaxyCounter++
			}
		}
		if emptyRow {
			spaceRows = append(spaceRows, y)
		}
	}

	for i, e := range emptyCols {
		if e {
			spaceCols = append(spaceCols, i)
		}
	}

	return galaxies, spaceCols, spaceRows
}

func SpaceHamDistance(g1, g2 Galaxy, spaceRows, spaceCols []int, expansionTime int) int {
	ham := maths.AbsInt(g1.X-g2.X) + maths.AbsInt(g1.Y-g2.Y)
	spaceX := RangeCount(maths.MinInt(g1.X, g2.X), maths.MaxInt(g1.X, g2.X), spaceCols)
	spaceY := RangeCount(maths.MinInt(g1.Y, g2.Y), maths.MaxInt(g1.Y, g2.Y), spaceRows)
	return ham + (expansionTime-1)*(spaceX+spaceY)
}

func RangeCount(minV, maxV int, data []int) int {
	count := 0
	for _, v := range data {
		if minV <= v && v <= maxV {
			count++
		}
	}
	return count
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
