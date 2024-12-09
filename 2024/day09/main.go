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

func part1(input string) int {
	data := parseInput(input)

	fwd := 0
	bwd := len(data) - 1
	blkPos := 0

	result := 0
	for fwd <= bwd {
		fwdVal := data[fwd]
		if fwdVal > 0 {
			if fwd%2 == 0 {
				fileID := fwd / 2
				result += fileID * blkPos
				blkPos += 1
				data[fwd] -= 1
			} else {
				bwdVal := data[bwd]
				if bwdVal > 0 {
					if bwd%2 == 0 {
						fileID := bwd / 2
						result += fileID * blkPos
						blkPos += 1
						data[bwd] -= 1
						data[fwd] -= 1
					} else {
						data[bwd] = 0
						bwd -= 1
					}
				} else {
					bwd -= 1
				}
			}
		} else {
			fwd += 1
		}
	}

	return result
}

type DiskMap struct {
	files []int
}

func createDiskMap(blocks []int) DiskMap {
	diskMap := DiskMap{files: []int{}}
	idx := 0

	for i := 0; i < len(blocks); i++ {
		if i%2 != 0 {
			for j := 0; j < blocks[i]; j++ {
				diskMap.files = append(diskMap.files, 0)
			}
		} else {
			for j := 0; j < blocks[i]; j++ {
				diskMap.files = append(diskMap.files, idx)
			}
			idx++
		}
	}
	return diskMap
}

func (d *DiskMap) sort(idxShift int) {
	left := idxShift
	right := len(d.files) - 1

	for left < right {
		for left < right && d.files[left] != 0 {
			left++
		}
		for left < right && d.files[right] == 0 {
			right--
		}
		if left < right {
			d.files[left], d.files[right] = d.files[right], d.files[left]
		}
	}
}

func (d *DiskMap) sortFiles(idxShift int) {
	maxID := 0
	for _, file := range d.files {
		if file > maxID {
			maxID = file
		}
	}

	for id := maxID; id > 0; id-- {
		fileStart, fileEnd := -1, -1

		for i := idxShift; i < len(d.files); i++ {
			if d.files[i] == id {
				if fileStart == -1 {
					fileStart = i
				}
				fileEnd = i
			}
		}

		if fileStart == -1 {
			continue
		}

		fileLength := fileEnd - fileStart + 1

		spaceStart, spaceLength := -1, 0
		for i := idxShift; i < fileStart; i++ {
			if d.files[i] == 0 {
				if spaceStart == -1 {
					spaceStart = i
				}
				spaceLength++

				if spaceLength >= fileLength {
					break
				}
			} else {
				spaceStart, spaceLength = -1, 0
			}
		}

		if spaceStart != -1 && spaceLength >= fileLength {
			for i := 0; i < fileLength; i++ {
				d.files[spaceStart+i] = id
			}

			for i := fileStart; i <= fileEnd; i++ {
				d.files[i] = 0
			}
		}
	}
}

func part2(input string) int {
	data := parseInput(input)
	diskMap := createDiskMap(data)

	diskMap.sortFiles(data[0])

	result := 0
	for i, f := range diskMap.files {
		result += i * f
	}
	return result
}

func parseInput(input string) []int {
	data := []int{}
	for _, chr := range input {
		data = append(data, cast.ToInt(string(chr)))
	}
	return data
}
