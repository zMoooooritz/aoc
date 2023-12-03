package main

import (
	_ "embed"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/zMoooooritz/advent-of-code/maths"
	"github.com/zMoooooritz/advent-of-code/util"
)

//go:embed input.txt
var input string

type Pull struct {
	Red   int
	Green int
	Blue  int
}

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

func part1(input string) int {
	parsed := parseInput(input)
	bag := Pull{12, 13, 14}
	result := 0

	for _, p := range parsed {
		gameId, pulls := parseGame(p)
		foundError := false
		for _, pull := range pulls {
			if isInvalidPull(pull, bag) {
				foundError = true
			}
		}
		if !foundError {
			result += gameId
		}

	}
	return result
}

func part2(input string) int {
	parsed := parseInput(input)
	result := 0

	for _, p := range parsed {
		_, pulls := parseGame(p)
		merge := mergePulls(pulls)
		result += merge.Red * merge.Green * merge.Blue
	}
	return result
}

func isInvalidPull(pull Pull, bag Pull) bool {
	return pull.Red > bag.Red || pull.Green > bag.Green || pull.Blue > bag.Blue
}

func mergePulls(pulls []Pull) Pull {
	merge := Pull{0, 0, 0}
	for _, pull := range pulls {
		merge.Red = maths.MaxInt(merge.Red, pull.Red)
		merge.Green = maths.MaxInt(merge.Green, pull.Green)
		merge.Blue = maths.MaxInt(merge.Blue, pull.Blue)
	}
	return merge
}

func parseGame(str string) (int, []Pull) {
	var pulls []Pull
	split := strings.Split(str, ":")
	gameStr := split[0]
	gameId, _ := strconv.Atoi(strings.Split(gameStr, " ")[1])
	for _, pull := range strings.Split(split[1], ";") {
		pulls = append(pulls, parsePull(pull))
	}
	return gameId, pulls
}

func parsePull(str string) Pull {
	pull := Pull{0, 0, 0}
	pv := reflect.ValueOf(&pull)
	v := pv.Elem()
	typeOf := v.Type()
	for _, s := range strings.Split(str, ",") {
		for i := 0; i < v.NumField(); i++ {
			color := strings.ToLower(typeOf.Field(i).Name)
			if !strings.Contains(s, color) {
				continue
			}
			s = strings.Trim(s, color)
			s = strings.Trim(s, " ")
			value, _ := strconv.Atoi(s)
			v.Field(i).SetInt(int64(value))
		}
	}
	return pull
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
