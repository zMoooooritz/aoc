package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
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

func applyRule(val int) []int {
	sVal := strconv.Itoa(val)
	if val == 0 {
		return []int{1}
	} else if len(sVal)%2 == 0 {
		fVal, _ := strconv.Atoi(sVal[:len(sVal)/2])
		bVal, _ := strconv.Atoi(sVal[len(sVal)/2:])
		return []int{fVal, bVal}
	} else {
		return []int{val * 2024}
	}
}

func part1(input string) int {
	config := parseInput(input)
	fmt.Println(config)

	for i := 0; i < 25; i++ {
		newConfig := []int{}
		for _, c := range config {
			newConfig = append(newConfig, applyRule(c)...)
		}
		config = newConfig
	}

	return len(config)
}

func part2(input string) int {
	config := parseInput(input)
	fmt.Println(config)

	configMap := map[int]int{}
	for _, c := range config {
		if _, ok := configMap[c]; ok {
			configMap[c] += 1
		} else {
			configMap[c] = 1
		}
	}

	for i := 0; i < 75; i++ {
		newConfigMap := map[int]int{}
		for k, v := range configMap {
			for _, nk := range applyRule(k) {
				if _, ok := newConfigMap[nk]; ok {
					newConfigMap[nk] += v
				} else {
					newConfigMap[nk] = v
				}
			}
		}
		configMap = newConfigMap
	}

	result := 0
	for _, v := range configMap {
		result += v
	}
	return result
}

func parseInput(input string) (config []int) {
	return cast.ToIntSlice(input)
}
