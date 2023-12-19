package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strconv"
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

type Category int

const (
	X Category = iota
	M
	A
	S
)

type Comparator int

const (
	GT Comparator = iota
	LT
)

type Part struct {
	categories []int
}

func (p Part) sum() int {
	sum := 0
	for _, rating := range p.categories {
		sum += rating
	}
	return sum
}

type Range struct {
	start int
	end   int
}

type RangePart struct {
	categorieRanges [][]Range
}

func (rp RangePart) availableCombinations() int {
	result := 1
	for _, pr := range rp.categorieRanges {
		rating := 0
		for _, r := range pr {
			rating += r.end - r.start + 1
		}
		result *= rating
	}
	return result
}

type Rule struct {
	category    Category
	comparator  Comparator
	value       int
	destination string
}

type Workflow struct {
	name  string
	rules []Rule
}

func (workflow Workflow) processPart(part Part) string {
	for _, r := range workflow.rules {
		value := part.categories[r.category]

		if r.comparator == LT && value < r.value {
			return r.destination
		}
		if r.comparator == GT && value > r.value {
			return r.destination
		}
	}
	panic("invalid workflow")
}

var workflows []Workflow

func part1(input string) int {
	_, parts := parseInput(input)

	result := 0
	for _, part := range parts {
		currWorflowName := "in"
		for true {
			wf := getWorkflowByName(currWorflowName)
			currWorflowName = wf.processPart(part)
			if currWorflowName == "A" {
				result += part.sum()
				break
			}
			if currWorflowName == "R" {
				break
			}
		}
	}

	return result
}

func part2(input string) int {
	parseInput(input)
	workflows = flattenWorkflows(workflows)
	fullRange := []Range{{1, 4000}}
	rp := RangePart{}
	rp.categorieRanges = [][]Range{fullRange, fullRange, fullRange, fullRange}
	return calculateValidCombinations("in", rp)
}

func calculateValidCombinations(workflowName string, rp RangePart) int {
	if workflowName == "A" {
		return rp.availableCombinations()
	}
	if workflowName == "R" {
		return 0
	}

	wf := getWorkflowByName(workflowName)

	// deep-copy slices of RangePart for recursion
	rpL, rpR := rp, rp
	ranges := rp.categorieRanges
	rpL.categorieRanges = make([][]Range, len(ranges))
	copy(rpL.categorieRanges, ranges)
	rpR.categorieRanges = make([][]Range, len(ranges))
	copy(rpR.categorieRanges, ranges)

	fstRule := wf.rules[0]
	sndRule := wf.rules[1]
	lR, rR := splitRanges(ranges[fstRule.category], fstRule)
	rpL.categorieRanges[fstRule.category] = lR
	rpR.categorieRanges[fstRule.category] = rR

	return calculateValidCombinations(fstRule.destination, rpL) + calculateValidCombinations(sndRule.destination, rpR)
}

func splitRanges(ranges []Range, rule Rule) ([]Range, []Range) {
	lR, rR := []Range{}, []Range{}
	splitValue := rule.value
	if rule.comparator == GT {
		splitValue += 1
	}

	for _, r := range ranges {
		if splitValue <= r.start {
			rR = append(rR, r)
		} else if r.end < splitValue {
			lR = append(lR, r)
		} else {
			lR = append(lR, Range{r.start, splitValue - 1})
			rR = append(rR, Range{splitValue, r.end})
		}
	}

	if rule.comparator == LT {
		return lR, rR
	} else {
		return rR, lR
	}
}

func flattenWorkflows(workflows []Workflow) []Workflow {
	flattened := []Workflow{}
	for _, wf := range workflows {
		for r := range wf.rules {
			if wf.rules[r].destination != "A" && wf.rules[r].destination != "R" {
				wf.rules[r].destination += "1"
			}
		}

		if wf.name == "in" {
			flattened = append(flattened, wf)
			continue
		}

		i := 0
		for ; i < len(wf.rules)-2; i++ {
			firstRule := wf.rules[i]
			secondRule := Rule{X, GT, 0, wf.name + strconv.Itoa(i+2)}
			flattened = append(flattened, Workflow{wf.name + strconv.Itoa(i+1), []Rule{firstRule, secondRule}})
		}
		lastRules := wf.rules[len(wf.rules)-2:]
		flattened = append(flattened, Workflow{wf.name + strconv.Itoa(i+1), lastRules})
	}
	return flattened
}

func getWorkflowByName(name string) Workflow {
	for _, wf := range workflows {
		if wf.name == name {
			return wf
		}
	}
	fmt.Println(name)
	panic("no workflow found")
}

func parseInput(input string) ([]Workflow, []Part) {
	data := strings.Split(input, "\n\n")
	wf, pt := data[0], data[1]

	workflows = []Workflow{}

	for _, line := range strings.Split(wf, "\n") {
		line = strings.Trim(line, "}")
		split := strings.Split(line, "{")
		name := split[0]
		rules := []Rule{}
		for _, r := range strings.Split(split[1], ",") {
			if strings.Contains(r, ":") {
				colonsplit := strings.Split(r, ":")
				if strings.Contains(colonsplit[0], "<") {
					splt := strings.Split(colonsplit[0], "<")
					rules = append(rules, Rule{strToCategory(splt[0]), LT, cast.ToInt(splt[1]), colonsplit[1]})
				} else {
					splt := strings.Split(colonsplit[0], ">")
					rules = append(rules, Rule{strToCategory(splt[0]), GT, cast.ToInt(splt[1]), colonsplit[1]})
				}
			} else {
				rules = append(rules, Rule{X, GT, 0, r})
			}
		}
		workflows = append(workflows, Workflow{name, rules})
	}

	parts := []Part{}
	for _, line := range strings.Split(pt, "\n") {
		line = strings.Trim(line, "{")
		line = strings.Trim(line, "}")

		var part Part
		part.categories = []int{0, 0, 0, 0}
		for _, a := range strings.Split(line, ",") {
			splt := strings.Split(a, "=")
			value := cast.ToInt(splt[1])
			category := strToCategory(splt[0])
			part.categories[category] = value
		}
		parts = append(parts, part)
	}

	return workflows, parts
}

func strToCategory(in string) Category {
	if in == "x" {
		return X
	}
	if in == "m" {
		return M
	}
	if in == "a" {
		return A
	}
	if in == "s" {
		return S
	}
	panic("invalid in value")
}
