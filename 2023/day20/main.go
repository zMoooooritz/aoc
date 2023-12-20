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

type ModuleType int

const (
	NONE ModuleType = iota
	BROADCASTER
	FLIP_FLOP
	CONJUNCTION
)

type Module struct {
	name         string
	moduleType   ModuleType
	receivers    []string
	state        bool
	senderStates map[string]bool
}

type Network struct {
	modules map[string]Module
}

type SendInformation struct {
	src   string
	dest  string
	pulse bool
}

type Metadata struct {
	targets    map[string]bool
	counters   map[string]int
	pressCount int
}

func (n Network) getModuleByName(name string) (Module, error) {
	module, found := n.modules[name]
	if !found {
		return Module{}, fmt.Errorf("Empty Module")

	}
	return module, nil
}

var network Network
var sendQueue []SendInformation
var lowSendCount = 0
var highSendCount = 0

func part1(input string) int {
	parseInput(input)

	iterations := 1000
	for i := 0; i < iterations; i++ {
		sendQueue = []SendInformation{}
		sendQueue = append(sendQueue, SendInformation{"button", "broadcaster", false})
		for {
			queue := make([]SendInformation, len(sendQueue))
			copy(queue, sendQueue)
			sendQueue = []SendInformation{}
			for _, sI := range queue {
				handlePulse(sI, nil)
			}
			if len(sendQueue) == 0 {
				break
			}
		}
	}

	return lowSendCount * highSendCount
}

func part2(input string) int {
	parseInput(input)

	rxName := "rx"

	var finalModule Module
	for _, module := range network.modules {
		for _, recv := range module.receivers {
			if recv == rxName {
				finalModule = module
				break
			}
		}
	}

	targets := map[string]bool{}
	for label := range finalModule.senderStates {
		targets[label] = true
	}

	fmt.Println(targets)

	metadata := Metadata{targets, make(map[string]int), 1}

	for {
		sendQueue = []SendInformation{}
		sendQueue = append(sendQueue, SendInformation{"button", "broadcaster", false})
		doBreak := false
		for {
			queue := make([]SendInformation, len(sendQueue))
			copy(queue, sendQueue)
			sendQueue = []SendInformation{}
			for _, sI := range queue {
				handlePulse(sI, &metadata)
			}
			if len(metadata.targets) == 0 {
				doBreak = true
				break
			}
		}
		if doBreak {
			break
		}
		metadata.pressCount++
	}

	nums := []int{}
	for _, count := range metadata.counters {
		nums = append(nums, int(count))
	}
	return findLCM(nums)
}

func handlePulse(sI SendInformation, metadata *Metadata) {
	if sI.pulse == false {
		lowSendCount++
	} else {
		highSendCount++
	}

	module, err := network.getModuleByName(sI.dest)
	// nothing to do
	if err != nil {
		return
	}

	doSend := false
	sendPulse := false
	switch module.moduleType {
	case BROADCASTER:
		doSend = true
		sendPulse = sI.pulse
	case FLIP_FLOP:
		if sI.pulse == false {
			module.state = !module.state
			doSend = true
			sendPulse = module.state
		}
	case CONJUNCTION:
		module.senderStates[sI.src] = sI.pulse
		doSend = true
		sendPulse = false
		for _, v := range module.senderStates {
			if !v {
				sendPulse = true
				break
			}
		}

		if metadata != nil {
			if _, ok := metadata.targets[sI.dest]; ok && sendPulse == true {
				metadata.counters[sI.dest] = metadata.pressCount
				delete(metadata.targets, sI.dest)
			}
		}
	}
	if doSend {
		for _, m := range module.receivers {
			sendQueue = append(sendQueue, SendInformation{sI.dest, m, sendPulse})
		}
	}

	network.modules[sI.dest] = module
}

func findLCM(numbers []int) int {
	gcd := func(a, b int) int { //general common divisor
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}

	lcm := func(a, b int) int { //least common multiple
		return a * b / gcd(a, b)
	}

	result := int(1)
	for _, num := range numbers {
		result = lcm(result, num)
	}
	return result
}

func parseInput(input string) {
	network.modules = make(map[string]Module)

	for _, line := range strings.Split(input, "\n") {
		data := strings.Split(line, " -> ")
		module := Module{}

		if strings.HasPrefix(data[0], "%") {
			module.moduleType = FLIP_FLOP
			module.name = data[0][1:]
			module.state = false
		} else if strings.HasPrefix(data[0], "&") {
			module.moduleType = CONJUNCTION
			module.name = data[0][1:]
			module.senderStates = make(map[string]bool)
		} else {
			module.moduleType = BROADCASTER
			module.name = data[0]
		}

		receivers := []string{}
		for _, r := range strings.Split(data[1], ",") {
			receivers = append(receivers, strings.Trim(r, " "))
		}
		module.receivers = receivers

		network.modules[module.name] = module
	}

	for i, m := range network.modules {
		if m.moduleType == CONJUNCTION {
			count := 0
			for _, sm := range network.modules {
				for _, r := range sm.receivers {
					if r == m.name {
						count += 1
						network.modules[i].senderStates[sm.name] = false
						continue
					}
				}
			}
		}
	}
}
