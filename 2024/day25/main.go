package main

import (
	_ "embed"
	"flag"
	"fmt"
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

type Profile struct {
	heightMap []int
	width     int
	height    int
	isKey     bool
}

func part1(input string) int {
	profiles := parseInput(input)

	fmt.Println(profiles)

	fits := 0
	for i := 0; i < len(profiles); i++ {
		for j := i + 1; j < len(profiles); j++ {
			fst, snd := profiles[i], profiles[j]
			if fst.isKey == snd.isKey {
				continue
			}
			if fst.width != snd.width || fst.height != snd.height {
				continue
			}
			isValid := true
			for x := range fst.width {
				if fst.heightMap[x]+snd.heightMap[x]+2 > fst.height {
					isValid = false
					break
				}
			}
			if isValid {
				fits += 1
			}
		}
	}

	return fits
}

func part2(input string) int {
	fmt.Println("AOC 2024 DONE")
	return 0
}

func parseInput(input string) []Profile {
	profiles := [][]string{}
	profile := []string{}
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			profiles = append(profiles, profile)
			profile = []string{}
			continue
		}

		profile = append(profile, line)
	}
	profiles = append(profiles, profile)

	pStructs := []Profile{}
	for _, profile := range profiles {
		pStructs = append(pStructs, toProfile(profile))
	}

	return pStructs
}

func toProfile(profile []string) Profile {
	isKey := profile[0][0] == '#'
	height := len(profile)
	width := len(profile[0])

	heightMap := []int{}
	for x := range width {
		count := 0
		for y := range height {
			if profile[y][x] == '#' {
				count += 1
			}
		}
		heightMap = append(heightMap, count-1)
	}
	return Profile{
		heightMap,
		width,
		height,
		isKey,
	}
}
