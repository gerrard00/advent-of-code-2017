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

func printKeys(g graph) {
	fmt.Println("Unvisited")
	fmt.Println("---------")
	for k, _ := range g {
		fmt.Println(k)
	}
	fmt.Println("---------")
}

func findRelated(targetNode *node, unvisitedNodes graph) (relatedNodes []*node, newUnvisitedNodes graph) {
	// fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	// fmt.Printf("Current node: %v unvisited count: %v\n", targetNode.id, len(unvisitedNodes))
	// printKeys(unvisitedNodes)

	if unvisitedNodes[targetNode.id] == nil {
		return nil, unvisitedNodes
	}

	relatedNodes = append(relatedNodes, targetNode)
	newUnvisitedNodes = unvisitedNodes
	delete(newUnvisitedNodes, targetNode.id)
	// fmt.Println(newUnvisitedNodes)
	// fmt.Printf("updated unvisited count: %v\n", len(newUnvisitedNodes))
	// fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n")

	var newChildRelatedNodes []*node
	for _, relatedNode := range targetNode.related {
		newChildRelatedNodes, newUnvisitedNodes = findRelated(relatedNode, newUnvisitedNodes)
		relatedNodes = append(relatedNodes, newChildRelatedNodes...)
	}

	return relatedNodes, newUnvisitedNodes
}

func fn(input string) int {
	getFirstNodeFromGraph := func(g graph) *node {
		for _, v := range g {
			return v
		}

		return nil
	}

	g := convertToGraph(input)
	unvisited := make(graph)
	// TODO: this should really just be another graph
	// var relatedNodes []*node
	for k, v := range g {
		unvisited[k] = v
	}
	groupCount := 0

	// var debug = 0

	for {
		// debug++

		// if debug > 3 {
		// 	fmt.Println("too far")
		// 	os.Exit(-1)
		// }

		targetNode := getFirstNodeFromGraph(unvisited)

		// relatedNodes, unvisited = findRelated(targetNode, unvisited)
		_, unvisited = findRelated(targetNode, unvisited)

		groupCount++

		// fmt.Printf("Group count: %v\n\n\n\n\n", groupCount)
		// printKeys(unvisited)
		if len(unvisited) == 0 {
			break
		}
	}
	return groupCount
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
		6 <-> 4, 5`, 0, 2},
	// real input, real result
	{getRealInput(), 0, 204},
}

func runTest(t *testing.T, testInputs []testpair) {
	for _, pair := range testInputs {
		v := fn(pair.graphString)

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
