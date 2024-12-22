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

func calcNextSecretNumber(secretNumber int) int {
	secretNumber ^= secretNumber << 6
	secretNumber &= 0xFFFFFF
	secretNumber ^= secretNumber >> 5
	secretNumber &= 0xFFFFFF
	secretNumber ^= secretNumber << 11
	secretNumber &= 0xFFFFFF
	return secretNumber
}

func part1(input string) int {
	values := parseInput(input)

	result := 0
	results := map[int]int{}

	iterationCount := 2000
	for _, value := range values {
		curr := value
		for range iterationCount {
			value = calcNextSecretNumber(value)
		}
		results[curr] = value
		result += value
	}

	return result
}

type Sequence struct {
	fst int
	snd int
	trd int
	frt int
}

func part2(input string) int {
	values := parseInput(input)

	prices := map[Sequence]int{}

	iterationCount := 2000
	for _, value := range values {

		priceList := []int{value % 10}
		for range iterationCount {
			value = calcNextSecretNumber(value)
			priceList = append(priceList, value%10)
		}

		priceChanges := []int{}
		for index := range iterationCount - 1 {
			priceChanges = append(priceChanges, priceList[index+1]-priceList[index])
		}

		seen := map[Sequence]struct{}{}
		for index := range iterationCount - 4 {
			cSeq := Sequence{priceChanges[index], priceChanges[index+1], priceChanges[index+2], priceChanges[index+3]}
			if _, ok := seen[cSeq]; !ok {
				seen[cSeq] = struct{}{}
				prices[cSeq] += priceList[index+4]
			}
		}
	}

	maxSeq := Sequence{}
	maxVal := 0

	for seq, val := range prices {
		if val > maxVal {
			maxVal = val
			maxSeq = seq
		}
	}

	fmt.Println("Seq:", maxSeq, "Val:", maxVal)

	return maxVal
}

func parseInput(input string) []int {
	values := []int{}
	for _, line := range strings.Split(input, "\n") {
		values = append(values, cast.ToInt(line))
	}
	return values
}
