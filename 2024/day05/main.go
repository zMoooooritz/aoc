package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
	"sort"
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

type BeforeConds map[int][]int

type UpdateSequence []int

func (us UpdateSequence) middleNumber() int {
	return us[len(us)/2]
}

func (us UpdateSequence) isValid(beforeConds BeforeConds) bool {
	seenNumbers := []int{}

	for _, n := range us {
		if conds, ok := beforeConds[n]; ok {
			for _, cond := range conds {
				if slices.Contains(seenNumbers, cond) {
					return false
				}
			}
		}
		seenNumbers = append(seenNumbers, n)
	}
	return true
}

func part1(input string) int {
	beforeConds, sequences := parseInput(input)

	result := 0
	for _, sequence := range sequences {
		if sequence.isValid(beforeConds) {
			result += sequence.middleNumber()
		}
	}

	return result
}

func part2(input string) int {
	beforeConds, sequences := parseInput(input)

	result := 0
	for _, sequence := range sequences {
		if !sequence.isValid(beforeConds) {
			sort.SliceStable(sequence, func(i, j int) bool {
				if cond, ok := beforeConds[sequence[i]]; ok {
					if slices.Contains(cond, sequence[j]) {
						return true
					}
				}
				return false
			})

			result += sequence.middleNumber()
		}
	}

	return result
}

func parseInput(input string) (BeforeConds, []UpdateSequence) {
	beforeConds := BeforeConds{}
	sequences := []UpdateSequence{}

	for _, line := range strings.Split(input, "\n") {
		if strings.Contains(line, "|") {
			splt := cast.ToIntSliceSep(line, "|")
			src := splt[0]
			dst := splt[1]
			if _, ok := beforeConds[src]; ok {
				beforeConds[src] = append(beforeConds[src], dst)
			} else {
				beforeConds[src] = []int{dst}
			}
		} else if strings.Contains(line, ",") {
			splt := cast.ToIntSliceSep(line, ",")
			sequences = append(sequences, UpdateSequence(splt))
		}
	}
	return beforeConds, sequences
}
