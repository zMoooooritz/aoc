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

type Lens struct {
	label  string
	focLen int
}

type Box struct {
	lenses []Lens
}

func (b *Box) addLens(lens Lens) {
	index := b.indexOfLens(lens)
	if index != -1 {
		b.lenses[index] = lens
	} else {
		b.lenses = append(b.lenses, lens)
	}
}

func (b *Box) removeLens(lens Lens) {
	index := b.indexOfLens(lens)
	if index != -1 {
		b.lenses = append(b.lenses[:index], b.lenses[index+1:]...)
	}
}

func (b *Box) indexOfLens(lens Lens) int {
	for i, l := range b.lenses {
		if l.label == lens.label {
			return i
		}
	}
	return -1
}

func part1(input string) int {
	parsed := parseInput(input)

	result := 0
	for _, str := range parsed {
		result += hashString(str)
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)

	boxes := []Box{}
	for i := 0; i < 256; i++ {
		boxes = append(boxes, Box{})
	}

	for _, str := range parsed {
		data := []string{}
		sign := "-"
		focLen := 0
		if strings.Contains(str, "=") {
			data = strings.Split(str, "=")
			sign = "="
			focLen = cast.ToInt(data[1])
		} else {
			data = strings.Split(str, "-")
		}
		label := data[0]
		hash := hashString(label)
		lens := Lens{label, focLen}
		if sign == "=" {
			boxes[hash].addLens(lens)
		} else {
			boxes[hash].removeLens(lens)
		}
	}

	result := 0
	for boxIndex, box := range boxes {
		for lensIndex, l := range box.lenses {
			result += (boxIndex + 1) * (lensIndex + 1) * l.focLen
		}
	}
	return result
}

func hashString(str string) int {
	hash := 0
	for _, c := range str {
		hash = (hash + int(c)) * 17 % 256
	}
	return hash
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, ",") {
		ans = append(ans, strings.Trim(line, "\n"))
	}
	return ans
}
