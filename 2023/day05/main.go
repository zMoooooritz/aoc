package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"strings"
	"unicode"

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

type MappingStage struct {
	mappings []Mapping
}

func (ms *MappingStage) add(m Mapping) {
	ms.mappings = append(ms.mappings, m)
}

func (ms *MappingStage) apply(val int, reverse bool) int {
	for _, m := range ms.mappings {
		var err error
		var nVal int
		if reverse {
			nVal, err = m.revApply(val)
		} else {
			nVal, err = m.apply(val)
		}
		if err == nil {
			return nVal
		}
	}
	return val
}

type Mapping struct {
	src int
	dst int
	cnt int
}

func (m *Mapping) apply(val int) (int, error) {
	if m.src <= val && val < m.src+m.cnt {
		return m.dst + val - m.src, nil
	}
	return 0, errors.New("unable to map value")
}

func (m *Mapping) revApply(val int) (int, error) {
	if m.dst <= val && val < m.dst+m.cnt {
		return m.src + val - m.dst, nil
	}
	return 0, errors.New("unable to map value")
}

type Range struct {
	start int
	end   int
}

func part1(input string) int {
	parsed := parseInput(input)

	seeds := cast.ToIntSlice(strings.Split(parsed[0], ":")[1])
	mappingStages := buildStages(parsed[1:])

	end := 0
	for i, s := range seeds {
		val := s
		for _, ms := range mappingStages {
			val = ms.apply(val, false)
		}
		if i == 0 {
			end = val
		}

		if val < end {
			end = val
		}
	}

	return end
}

// stupid brute force solutions
// the proper solution probably works on ranges and divides/merges them according to the mapping
// currently not in the mood to fiddle around with ranges/intervals :/
func part2(input string) int {
	parsed := parseInput(input)

	seedIn := cast.ToIntSlice(strings.Split(parsed[0], ":")[1])
	seeds := []Range{}
	for i := 0; i < len(seedIn); i += 2 {
		seeds = append(seeds, Range{seedIn[i], seedIn[i] + seedIn[i+1]})
	}

	mappingStages := buildStages(parsed[1:])
	end := 0
	for i, s := range seeds {
		for j := s.start; j < s.end; j++ {
			val := j
			for _, ms := range mappingStages {
				val = ms.apply(val, false)
			}
			if i == 0 && j == seeds[0].start {
				end = val
			}

			if val < end {
				end = val
			}

		}
	}
	return end
}

func buildStages(data []string) []MappingStage {
	stages := []MappingStage{}
	mappingStage := MappingStage{}
	stageActive := false

	for _, line := range data {
		if line == "" {
			continue
		}
		if unicode.IsDigit(rune(line[0])) {
			vals := cast.ToIntSlice(line)
			mappingStage.add(Mapping{vals[1], vals[0], vals[2]})
			stageActive = true
		} else {
			if stageActive {
				stages = append(stages, mappingStage)
				mappingStage = MappingStage{}
				stageActive = false
			}
		}
	}
	stages = append(stages, mappingStage)
	return stages
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
