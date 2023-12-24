package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/zMoooooritz/advent-of-code/cast"
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

type Vector3 struct {
	X int
	Y int
	Z int
}

type Vector2 struct {
	X float64
	Y float64
}

type HailStone struct {
	Pos Vector3
	Vel Vector3
}

// there is probably a nicer way to deal with the floating point imprecision
func part1(input string) int {
	hailStones := parseInput(input)
	_ = hailStones

	testAreaMin := 200000000000000.0
	testAreaMax := 400000000000000.0
	factor := 10000000000.0

	count := 0
	for i := range hailStones {
		for j := i + 1; j < len(hailStones); j++ {
			p1Curr := Vector2{float64(hailStones[i].Pos.X), float64(hailStones[i].Pos.Y)}
			p1Next := Vector2{float64(hailStones[i].Pos.X) + factor*float64(hailStones[i].Vel.X), float64(hailStones[i].Pos.Y) + factor*float64(hailStones[i].Vel.Y)}
			p2Curr := Vector2{float64(hailStones[j].Pos.X), float64(hailStones[j].Pos.Y)}
			p2Next := Vector2{float64(hailStones[j].Pos.X) + factor*float64(hailStones[j].Vel.X), float64(hailStones[j].Pos.Y) + factor*float64(hailStones[j].Vel.Y)}

			denom := (p1Curr.X-p1Next.X)*(p2Curr.Y-p2Next.Y) - (p1Curr.Y-p1Next.Y)*(p2Curr.X-p2Next.X)

			// no intersection
			if denom == 0 {
				continue
			}

			pX := ((p1Curr.X*p1Next.Y-p1Curr.Y*p1Next.X)*(p2Curr.X-p2Next.X) - (p1Curr.X-p1Next.X)*(p2Curr.X*p2Next.Y-p2Curr.Y*p2Next.X)) / denom
			pY := ((p1Curr.X*p1Next.Y-p1Curr.Y*p1Next.X)*(p2Curr.Y-p2Next.Y) - (p1Curr.Y-p1Next.Y)*(p2Curr.X*p2Next.Y-p2Curr.Y*p2Next.X)) / denom

			// not inside test area
			if !(testAreaMin <= pX && pX <= testAreaMax && testAreaMin <= pY && pY <= testAreaMax) {
				continue
			}

			p1InPast := ((p1Curr.X < pX) != (p1Curr.X < p1Next.X)) || ((p1Curr.Y < pY) != (p1Curr.Y < p1Next.Y))
			p2InPast := ((p2Curr.X < pX) != (p2Curr.X < p2Next.X)) || ((p2Curr.Y < pY) != (p2Curr.Y < p2Next.Y))

			if p1InPast || p2InPast {
				continue
			}

			count++
		}
	}

	return count
}

func part2(input string) int {
	return 0
}

func parseInput(input string) []HailStone {
	hailStones := []HailStone{}
	for _, line := range strings.Split(input, "\n") {
		data := strings.Split(line, "@")
		hailStones = append(hailStones, HailStone{toVector3(data[0]), toVector3(data[1])})
	}
	return hailStones
}

func toVector3(input string) Vector3 {
	data := strings.Split(input, ",")
	intData := []int{}
	for _, d := range data {
		intData = append(intData, cast.ToInt(strings.Trim(d, " ")))
	}
	return Vector3{intData[0], intData[1], intData[2]}
}
