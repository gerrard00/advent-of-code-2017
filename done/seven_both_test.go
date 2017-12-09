package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

type node struct {
	name            string
	weight          int
	correctedWeight int
	totalWeight     int
	children        []*node
}

func (n node) getTotalWeight() (int, *node) {
	var badStack *node = nil
	if n.totalWeight > 0 {
		return n.totalWeight, badStack
	}

	sum := n.weight

	childWeightDistribution := make(map[int][]*node)

	for _, child := range n.children {
		var currentChildTotalWeight int
		var currentBadStack *node

		currentChildTotalWeight, currentBadStack = child.getTotalWeight()
		sum += currentChildTotalWeight
		childWeightDistribution[currentChildTotalWeight] =
			append(childWeightDistribution[currentChildTotalWeight], child)

		if badStack == nil && currentBadStack != nil {
			badStack = currentBadStack
		}
	}

	if badStack == nil && len(childWeightDistribution) > 1 {
		var goodWeight int
		var badWeight int

		for k, v := range childWeightDistribution {
			if len(v) == 1 {
				badWeight = k
				// fmt.Printf("distro containing bad guy: %v\n", childWeightDistribution)
				badStack = v[0]
			} else {
				goodWeight = k
			}
		}

		if badStack != nil {
			badStack.correctedWeight = badStack.weight + (goodWeight - badWeight)
			// fmt.Printf("Current: %v Good: %v Bad: %v\n",
			// 	badStack.weight, goodWeight, badWeight)
		}
	}

	// fmt.Printf("distribution: %v\n", childWeightDistribution)
	return sum, badStack
}

func convertToGraph(input string) node {
	scanner := bufio.NewScanner(strings.NewReader(input))

	// map of pending children to parent nodes
	unrootedNodes := make(map[string]*node)
	pendingChildren := make(map[string]*node)

	for scanner.Scan() {
		newNode := node{}
		currentLine := scanner.Text()
		fmt.Sscanf(currentLine,
			"%s (%d)",
			&newNode.name, &newNode.weight)

		arrowIndex := strings.Index(currentLine, "->")
		if arrowIndex > -1 {
			childNames := strings.Split(currentLine[arrowIndex+3:], ", ")
			for _, childName := range childNames {
				pendingChildren[childName] = &newNode
			}
		}
		unrootedNodes[newNode.name] = &newNode
	}

	for k, v := range pendingChildren {
		v.children = append(v.children, unrootedNodes[k])
		delete(unrootedNodes, k)
	}

	if len(unrootedNodes) > 1 {
		panic("Multiple unrooted nodes found")
	}

	// TODO: hacky, is there another way?
	for _, v := range unrootedNodes {
		return *v
	}

	return node{}
}

func printGraph(graph *node, depth int) {
	prefix := strings.Repeat("\t", depth)
	fmt.Printf("%v%v\n", prefix, graph.name)

	newDepth := depth + 1
	for _, v := range graph.children {
		printGraph(v, newDepth)
	}
}

func findBadGuy(graph node) int {
	// totalWeight, badStack := graph.getTotalWeight()
	_, badStack := graph.getTotalWeight()

	// fmt.Printf("Total weight: %v Bad stack: %v Corrected Weight: %v\n",
	// 	totalWeight, badStack.name, badStack.correctedWeight)

	return badStack.correctedWeight
}

func fn(input string) int {
	root := convertToGraph(input)

	// printGraph(&root, 0)
	return findBadGuy(root)
}

type testpair struct {
	input  string
	output int
}

var testsOne = []testpair{
	{
		`pbga (66)
xhth (57)
ebii (61)
havc (66)
ktlj (57)
fwft (72) -> ktlj, cntj, xhth
qoyq (66)
padx (45) -> pbga, havc, qoyq
tknk (41) -> ugml, padx, fwft
jptl (61)
ugml (68) -> gyxo, ebii, jptl
gyxo (61)
cntj (57)`, 60},
	// real input, real result
	{getRealInput(), 391},
}

func runTest(t *testing.T, testInputs []testpair) {
	for _, pair := range testInputs {
		v := fn(pair.input)

		if v != pair.output {
			t.Error(
				"For", pair.input,
				"expected", pair.output,
				"got", v,
			)
		}
	}
}

func TestOne(t *testing.T) {
	runTest(t, testsOne)
}

func getRealInput() string {
	bytes, readFileErr := ioutil.ReadFile("seven_one_input")

	if readFileErr != nil {
		panic(readFileErr)
	}

	return string(bytes)
}
