package main

import (
	_ "embed"
	"flag"
	"fmt"
	"sort"
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

type Game struct {
	Hand  string
	Bid   int
	Value int
}

func part1(input string) int {
	parsed := parseInput(input)
	games := []Game{}

	for _, line := range parsed {
		s := strings.Split(line, " ")
		hand, bid := s[0], s[1]

		singleCardValue := singleCardValue(hand, false)
		clusteredCards := clusterHand(hand)
		value := calculateHandValue(singleCardValue, clusteredCards)

		games = append(games, Game{hand, cast.ToInt(bid), value})
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].Value < games[j].Value
	})

	result := 0
	for i, g := range games {
		result += (i + 1) * g.Bid
	}

	return result
}

func part2(input string) int {
	parsed := parseInput(input)
	games := []Game{}

	for _, line := range parsed {
		s := strings.Split(line, " ")
		hand, bid := s[0], s[1]

		singleCardValue := singleCardValue(hand, true)
		clusteredCards := clusterJokerHand(hand)
		value := calculateHandValue(singleCardValue, clusteredCards)

		games = append(games, Game{hand, cast.ToInt(bid), value})
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].Value < games[j].Value
	})

	result := 0
	for i, g := range games {
		result += (i + 1) * g.Bid
	}

	return result
}

func calculateHandValue(sgv int, cc []int) int {
	value := sgv
	cv := 0
	if cc[0] == 5 {
		cv = 6
	} else if cc[0] == 4 {
		cv = 5
	} else if cc[0] == 3 && cc[1] == 2 {
		cv = 4
	} else if cc[0] == 3 {
		cv = 3
	} else if cc[0] == 2 && cc[1] == 2 {
		cv = 2
	} else if cc[0] == 2 {
		cv = 1
	}
	value += cv * 10000000000

	return value
}

func clusterJokerHand(hand string) []int {
	cluster := []int{0}
	m := make(map[rune]int)
	jCount := 0
	for _, char := range hand {
		if char == 'J' {
			jCount += 1
		} else {
			m[char] += 1
		}
	}
	for _, v := range m {
		cluster = append(cluster, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(cluster)))
	cluster[0] += jCount
	return cluster
}

func clusterHand(hand string) []int {
	cluster := []int{}
	m := make(map[rune]int)
	for _, char := range hand {
		m[char] += 1
	}
	for _, v := range m {
		cluster = append(cluster, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(cluster)))
	return cluster
}

func singleCardValue(hand string, joker bool) int {
	val := 0
	for _, c := range hand {
		val = val*100 + cardToInt(c, joker)
	}
	return val
}

func cardToInt(card rune, joker bool) int {
	if unicode.IsDigit(card) {
		return int(card) - '0'
	}
	m := make(map[rune]int)
	m['T'] = 10
	if joker {
		m['J'] = 1
	} else {
		m['J'] = 11
	}
	m['Q'] = 12
	m['K'] = 13
	m['A'] = 14
	return m[card]
}

func parseInput(input string) (ans []string) {
	for _, line := range strings.Split(input, "\n") {
		ans = append(ans, line)
	}
	return ans
}
