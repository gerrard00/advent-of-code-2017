package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

type command struct {
	operator string
	x        string
	y        string
}

func getCommands(input string) []command {
	commands := make([]command, 0)
	commandStringParts := strings.Split(input, ",")
	for _, commandString := range commandStringParts {
		var newCommand command

		newCommand.operator = string(commandString[0])

		singleCommandParts := strings.Split(commandString, "/")

		newCommand.x = string(singleCommandParts[0][1:])

		if len(singleCommandParts) > 1 {
			newCommand.y = singleCommandParts[1]
		}

		// fmt.Println(commandString)
		// fmt.Println(newCommand)
		commands = append(commands, newCommand)
	}

	return commands
}

type programs struct {
	official    []string
	indexByName map[string]int
}

func newPrograms(num int) *programs {
	result := programs{
		make([]string, num),
		make(map[string]int),
	}

	for i := 0; i < num; i++ {
		result.official[i] = string('a' + i)
		result.indexByName[result.official[i]] = i
	}

	return &result
}

// TODO: does his need to be a pointer?
func (p *programs) resetIndex() {
	for i := 0; i < len(p.official); i++ {
		p.indexByName[p.official[i]] = i
	}
}

func (p programs) print() {
	fmt.Println(strings.Join(p.official, ""))
	fmt.Println(p.indexByName)
}

// TODO: see if this can be a non-pointer
func (p *programs) exchange(from int, to int) {
	// p.print()
	temp := p.official[to]
	p.official[to] = p.official[from]
	p.indexByName[p.official[to]] = to
	p.official[from] = temp
	p.indexByName[p.official[from]] = from
	// // p.print()
	// fmt.Println()
}

func fn(numberOfPrograms int, commandString string) string {
	p := newPrograms(numberOfPrograms)

	commands := getCommands(commandString)

	// p.print()

	for _, com := range commands {
		// fmt.Println(com)

		switch com.operator {
		case "x":
			exchangeFrom, _ := strconv.Atoi(com.x)
			exchangeTo, _ := strconv.Atoi(com.y)
			p.exchange(exchangeFrom, exchangeTo)
		case "p":
			partnerFrom, _ := p.indexByName[com.x]
			partnerTo, _ := p.indexByName[com.y]
			p.exchange(partnerFrom, partnerTo)
		case "s":
			// TODO: this isn't the fastest way, but may still be faster than lots
			// of swaps
			swapLength, _ := strconv.Atoi(com.x)
			swapStart := len(p.official) - swapLength
			p.official = append(p.official[swapStart:], p.official[:swapStart]...)
			p.resetIndex()
		default:
			panic(fmt.Sprintf("Unknown command: %v\n", com.operator))
		}
		// p.print()
	}

	return strings.Join(p.official, "")
}

type testpair struct {
	numberOfPrograms int
	commands         string
	output           string
}

var testsOne = []testpair{
	{5, `s1,x3/4,pe/b`, "baedc"},
	// real input, real result
	// wrong jibfndocaekplhgm
	// wrong poilbkfeajgdmnhc
	{16, getRealInput(), "ionlbkfeajgdmphc"},
}

func runTest(t *testing.T, testInputs []testpair) {
	for _, pair := range testInputs {
		v := fn(pair.numberOfPrograms, pair.commands)

		if v != pair.output {
			t.Errorf("\nexpected %v\ngot      %v\n", pair.output, v)
		}
	}
}

func TestOne(t *testing.T) {
	runTest(t, testsOne)
}

func getRealInput() string {
	bytes, readFileErr := ioutil.ReadFile("sixteen_input")

	if readFileErr != nil {
		panic(readFileErr)
	}

	// d'oh! have to remember to clip the final newline
	return string(bytes[:len(bytes)-1])
}
