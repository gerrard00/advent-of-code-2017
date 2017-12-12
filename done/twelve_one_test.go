package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

type node struct {
	id      int
	related map[int]*node
}

type graph map[int]*node

func convertToGraph(input string) graph {
	scanner := bufio.NewScanner(strings.NewReader(input))

	// 2 <-> 0, 3, 4
	graph := make(graph)
	// map of pending children to parent nodes
	type pendingRelationship struct {
		toId     int
		fromNode *node
	}

	pendingChildren := make([]pendingRelationship, 0)

	for scanner.Scan() {
		newNode := node{}
		currentLine := scanner.Text()
		fmt.Sscanf(currentLine,
			"%d <-> ",
			&newNode.id)

		newNode.related = make(map[int]*node)
		arrowIndex := strings.Index(currentLine, "<->")
		if arrowIndex > -1 {
			childNames := strings.Split(currentLine[arrowIndex+4:], ", ")
			for _, childName := range childNames {
				newId, _ := strconv.Atoi(childName)
				pendingChildren =
					append(pendingChildren, pendingRelationship{newId, &newNode})
			}
		}
		graph[newNode.id] = &newNode
	}

	for _, entry := range pendingChildren {
		// fmt.Printf("%v related to %v\n", entry.fromNode.id, entry.toId)
		// TODO: should really skip assignments if they exist
		entry.fromNode.related[entry.toId] = graph[entry.toId]
		graph[entry.toId].related[entry.fromNode.id] = entry.fromNode
	}

	return graph
}

func printGraph(g graph) {
	for _, n := range g {
		fmt.Printf("%v\n", n.id)
		for k, _ := range n.related {
			fmt.Printf("\t->%v\n", k)
		}
	}
}

func findRelated(targetNode *node, visitedNodeIds map[int]bool) (relatedNodes []*node, newVisitedNodeIds map[int]bool) {
	// fmt.Printf("curent visited: %v\n", visitedNodeIds)
	if visitedNodeIds[targetNode.id] {
		return nil, visitedNodeIds
	}

	relatedNodes = append(relatedNodes, targetNode)
	if visitedNodeIds == nil {
		newVisitedNodeIds = make(map[int]bool)
	} else {
		newVisitedNodeIds = visitedNodeIds
	}

	newVisitedNodeIds[targetNode.id] = true

	var newChildRelatedNodes []*node
	for _, relatedNode := range targetNode.related {
		newChildRelatedNodes, newVisitedNodeIds = findRelated(relatedNode, newVisitedNodeIds)
		relatedNodes = append(relatedNodes, newChildRelatedNodes...)
	}

	return relatedNodes, newVisitedNodeIds
}

func fn(input string, targetID int) int {
	g := convertToGraph(input)

	targetNode := g[targetID]

	relatedNodes, _ := findRelated(targetNode, nil)
	// printGraph(g)

	return len(relatedNodes)
}

type testpair struct {
	graphString string
	targetID    int
	output      int
}

var testsOne = []testpair{
	{
		`0 <-> 2
		1 <-> 1
		2 <-> 0, 3, 4
		3 <-> 2, 4
		4 <-> 2, 3, 6
		5 <-> 6
		6 <-> 4, 5`, 0, 6},
	// real input, real result
	{getRealInput(), 0, 378},
}

func runTest(t *testing.T, testInputs []testpair) {
	for _, pair := range testInputs {
		v := fn(pair.graphString, pair.targetID)

		if v != pair.output {
			t.Error(
				"For", pair.graphString,
				"target", pair.targetID,
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
	bytes, readFileErr := ioutil.ReadFile("twelve_input")

	if readFileErr != nil {
		panic(readFileErr)
	}

	return string(bytes)
}
