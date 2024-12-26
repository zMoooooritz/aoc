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

type Layer struct {
	index int
	data  [][]int
}

func (l *Layer) countOcc(target int) int {
	count := 0
	for _, row := range l.data {
		for _, value := range row {
			if value == target {
				count += 1
			}
		}
	}
	return count
}

func part1(input string) int {
	test := false

	width := 25
	height := 6
	if test {
		width = 3
		height = 2
	}

	layers := parseInput(input, width, height)

	minCount := 99999
	minIndex := 0
	for idx, l := range layers {
		cnt := l.countOcc(0)
		if cnt < minCount {
			minCount = cnt
			minIndex = idx
		}
	}

	l := layers[minIndex]
	return l.countOcc(1) * l.countOcc(2)
}

func overlayLayers(layers []Layer) Layer {
	width := len(layers[0].data[0])
	height := len(layers[0].data)

	data := [][]int{}
	for y := range height {
		row := []int{}
		for x := range width {
			for _, l := range layers {
				if l.data[y][x] == 2 {
					continue
				}
				row = append(row, l.data[y][x])
				break
			}
		}
		data = append(data, row)
	}
	return Layer{index: -1, data: data}
}

func part2(input string) int {
	test := false

	width := 25
	height := 6
	if test {
		width = 2
		height = 2
	}

	layers := parseInput(input, width, height)
	layer := overlayLayers(layers)

	for _, row := range layer.data {
		fmt.Println(row)
	}

	return 0
}

func parseInput(input string, width, height int) []Layer {
	layerSize := width * height
	layerCount := len(input) / layerSize

	layers := []Layer{}
	for index := range layerCount {
		data := input[index*layerSize : (index+1)*layerSize]
		result := [][]int{}
		row := []int{}
		for i, c := range data {
			row = append(row, int(c-'0'))

			if (i+1)%width == 0 {
				result = append(result, row)
				row = []int{}
			}
		}
		layers = append(layers, Layer{index, result})
	}
	return layers
}
