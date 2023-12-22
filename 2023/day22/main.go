package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"sort"
	"strings"

	"github.com/zMoooooritz/advent-of-code/cast"
	"github.com/zMoooooritz/advent-of-code/maths"
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

type Corner struct {
	X int
	Y int
	Z int
}

type Brick struct {
	FstCorner   Corner
	SndCorner   Corner
	Index       int
	Setteled    bool
	SupportedBy []int
	Supports    []int
}

func (b *Brick) distanceToGround() int {
	return maths.MinInt(b.FstCorner.Z, b.SndCorner.Z) - 1
}

func (b *Brick) height() int {
	return maths.MaxInt(b.FstCorner.Z, b.SndCorner.Z) - maths.MinInt(b.FstCorner.Z, b.SndCorner.Z) + 1
}

func (b *Brick) moveDown(distance int) {
	b.FstCorner.Z -= distance
	b.SndCorner.Z -= distance
}

func (b *Brick) toFootprint() image.Rectangle {
	return cornersToRectangle(b.FstCorner, b.SndCorner)
}

func cornersToRectangle(fstCorner, sndCorner Corner) image.Rectangle {
	minX := maths.MinInt(fstCorner.X, sndCorner.X)
	minY := maths.MinInt(fstCorner.Y, sndCorner.Y)
	maxX := maths.MaxInt(fstCorner.X, sndCorner.X)
	maxY := maths.MaxInt(fstCorner.Y, sndCorner.Y)
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

type Space struct {
	Bricks []Brick
}

func (s *Space) getLowestUnsettledBrick() (Brick, error) {
	minDist := int(1e9)
	var minDistBrick Brick
	for _, b := range s.Bricks {
		if b.Setteled {
			continue
		}

		dist := b.distanceToGround()
		if dist < minDist {
			minDist = dist
			minDistBrick = b
		}
	}
	if minDist == int(1e9) {
		return minDistBrick, fmt.Errorf("all bricks are setteled")
	}
	return minDistBrick, nil
}

func (s *Space) getSupportingBricks(fstCorner, sndCorner Corner) ([]Brick, int) {
	bricks := []Brick{}
	highestPoint := 0
	for _, b := range s.Bricks {
		if !b.Setteled {
			continue
		}
		brickHighestPoint := maths.MaxInt(b.FstCorner.Z, b.SndCorner.Z)
		if brickHighestPoint < highestPoint {
			continue
		}

		areaRect := cornersToRectangle(fstCorner, sndCorner)
		brickRect := b.toFootprint()

		if areaRect.Overlaps(brickRect) {
			if brickHighestPoint == highestPoint {
				bricks = append(bricks, b)
			} else {
				bricks = []Brick{b}
			}

			highestPoint = brickHighestPoint
		}
	}
	return bricks, highestPoint
}

func (s *Space) supportCount(index int) int {
	brick := s.Bricks[index]
	dist := brick.distanceToGround()
	brickRect := brick.toFootprint()

	// fmt.Println(brick, brick.distanceToGround(), brick.height(), dist)

	count := 0
	for _, b := range s.Bricks {
		if b.height()+b.distanceToGround() != dist {
			continue
		}

		// fmt.Println("Check brick: ", b)
		if brickRect.Overlaps(b.toFootprint()) {
			// fmt.Println("Support von: ", b, " an: ", brick)
			count++
		}
	}
	// fmt.Println("======================")
	return count
}

func (s *Space) makeBricksFall() {
	lowestBrick, err := s.getLowestUnsettledBrick()

	for err == nil {
		supportingBricks, highestPoint := s.getSupportingBricks(lowestBrick.FstCorner, lowestBrick.SndCorner)
		supportIndex := []int{}
		for _, b := range supportingBricks {
			supportIndex = append(supportIndex, b.Index)
		}
		dist := lowestBrick.distanceToGround()
		lowestBrick.moveDown(dist - highestPoint)
		lowestBrick.Setteled = true
		lowestBrick.SupportedBy = supportIndex

		s.Bricks[lowestBrick.Index] = lowestBrick
		lowestBrick, err = s.getLowestUnsettledBrick()
	}

	for i, brick := range s.Bricks {
		for _, supportedBy := range brick.SupportedBy {
			s.Bricks[supportedBy].Supports = append(s.Bricks[supportedBy].Supports, i)
		}
	}
}

func (s *Space) fallingCountOnRemove(brick Brick) int {
	falling := make(map[int]bool)
	falling[brick.Index] = true

	currBricks := []Brick{brick}
	for len(currBricks) > 0 {
		nextBricks := []Brick{}
		for _, b := range currBricks {
			for _, sIndex := range b.Supports {
				nextBricks = append(nextBricks, s.Bricks[sIndex])
			}
			isFalling := true
			for _, sIndex := range b.SupportedBy {
				if _, ok := falling[sIndex]; !ok {
					isFalling = false
				}
			}
			if isFalling {
				falling[b.Index] = true
			}
		}

		currBricks = nextBricks
	}
	return len(falling) - 1
}

func part1(input string) int {
	bricks := parseInput(input)

	space := Space{bricks}
	space.makeBricksFall()

	count := 0
	for _, brick := range space.Bricks {
		isBrickRequired := false
		for _, supportsIndex := range brick.Supports {
			if len(space.Bricks[supportsIndex].SupportedBy) == 1 {
				isBrickRequired = true
				break
			}
		}
		if !isBrickRequired || len(brick.Supports) == 0 {
			count++
		}
	}
	return count
}

func part2(input string) int {
	bricks := parseInput(input)

	space := Space{bricks}
	space.makeBricksFall()

	result := 0
	for _, b := range space.Bricks {
		result += space.fallingCountOnRemove(b)
	}
	return result
}

func parseInput(input string) (data []Brick) {
	for i, line := range strings.Split(input, "\n") {
		splt := strings.Split(line, "~")
		data = append(data, Brick{toCorner(splt[0]), toCorner(splt[1]), i, false, []int{}, []int{}})
	}
	sort.SliceStable(data, func(i, j int) bool {
		return maths.MinInt(data[i].FstCorner.Z, data[i].SndCorner.Z) < maths.MinInt(data[j].FstCorner.Z, data[j].FstCorner.Z)
	})
	for i := range data {
		data[i].Index = i
	}
	return data
}

func toCorner(input string) Corner {
	coords := strings.Split(input, ",")
	return Corner{cast.ToInt(coords[0]), cast.ToInt(coords[1]), cast.ToInt(coords[2])}
}
