package main

import (
	_ "embed"
	"flag"
	"fmt"
	"slices"
	"sort"
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

type Link struct {
	fst string
	snd string
}

type Interconnection struct {
	fst string
	snd string
	trd string
}

var adjacencies = map[string][]string{}

func part1(input string) int {
	parseInput(input)

	interConns := map[Interconnection]struct{}{}
	for node, conn := range adjacencies {
		connCount := len(conn)
		for i := 0; i < connCount; i++ {
			for j := i + 1; j < connCount; j++ {
				conn1, conn2 := conn[i], conn[j]
				if slices.Contains(adjacencies[conn1], conn2) && slices.Contains(adjacencies[conn2], conn1) {
					tmp := []string{node, conn1, conn2}
					sort.Sort(sort.StringSlice(tmp))
					interConns[Interconnection{tmp[0], tmp[1], tmp[2]}] = struct{}{}
				}
			}
		}
	}

	filtered := []Interconnection{}
	for interConn := range interConns {
		if strings.HasPrefix(interConn.fst, "t") || strings.HasPrefix(interConn.snd, "t") || strings.HasPrefix(interConn.trd, "t") {
			filtered = append(filtered, interConn)
		}
	}

	return len(filtered)
}

// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm (not an idea of my own ;))
func bronKerbosch(r, p, x []string, cliques *[][]string) {
	if len(p) == 0 && len(x) == 0 {
		clique := make([]string, len(r))
		copy(clique, r)
		*cliques = append(*cliques, clique)
		return
	}

	pCopy := make([]string, len(p))
	copy(pCopy, p)

	for _, v := range pCopy {
		rNew := append(r, v)

		pNew := intersect(p, adjacencies[v])
		xNew := intersect(x, adjacencies[v])

		bronKerbosch(rNew, pNew, xNew, cliques)

		p = remove(p, v)
		x = remove(x, v)
	}
}

func intersect(a, b []string) []string {
	set := map[string]bool{}
	for _, v := range b {
		set[v] = true
	}
	intersection := []string{}
	for _, v := range a {
		if set[v] {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func remove(slice []string, elem string) []string {
	result := []string{}
	for _, v := range slice {
		if v != elem {
			result = append(result, v)
		}
	}
	return result
}

func part2(input string) int {
	parseInput(input)

	r := []string{}
	p := []string{}
	for k := range adjacencies {
		p = append(p, k)
	}
	x := []string{}

	cliques := [][]string{}
	bronKerbosch(r, p, x, &cliques)

	maxClique := []string{}
	maxSize := 0
	for _, clique := range cliques {
		if len(clique) > maxSize {
			maxSize = len(clique)
			maxClique = clique
		}
	}

	sort.Sort(sort.StringSlice(maxClique))
	fmt.Println(strings.Join(maxClique, ","))
	return maxSize
}

func parseInput(input string) {
	links := []Link{}
	for _, line := range strings.Split(input, "\n") {
		splt := strings.Split(line, "-")
		links = append(links, Link{splt[0], splt[1]})
	}

	adjacencies = map[string][]string{}

	for _, link := range links {
		if _, ok := adjacencies[link.fst]; !ok {
			adjacencies[link.fst] = []string{}
		}
		adjacencies[link.fst] = append(adjacencies[link.fst], link.snd)
		if _, ok := adjacencies[link.snd]; !ok {
			adjacencies[link.snd] = []string{}
		}
		adjacencies[link.snd] = append(adjacencies[link.snd], link.fst)
	}
}
