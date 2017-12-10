package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

type registersType map[string]int

// if a > 1
type conditional struct {
	lefthand  string
	operator  string
	righthand string
}

// b inc 5 if a > 1
type command struct {
	register string
	operator string
	value    string
	cond     conditional
}

func tokenize(commandString string) command {
	newCommand := command{}

	// tokens := strings.Split(commandString, " ")

	// b inc 5 if a > 1
	_, scanErr := fmt.Sscanf(commandString, "%s %s %s if %s %s %s",
		&newCommand.register,
		&newCommand.operator,
		&newCommand.value,
		&newCommand.cond.lefthand,
		&newCommand.cond.operator,
		&newCommand.cond.righthand)

	if scanErr != nil {
		panic(scanErr)
	}

	return newCommand
}

func evaluate(variableNameOrLiteral string, registers registersType) int {
	if literalValue, err := strconv.Atoi(variableNameOrLiteral); err == nil {
		return literalValue
	}

	// TODO: should really throw if register doesn't exist
	return registers[variableNameOrLiteral]
}

func integerFunctions(operator string, operands ...int) int {
	switch operator {
	case "inc":
		return operands[0] + operands[1]
	case "dec":
		return operands[0] - operands[1]
	default:
		panic("unknown operator " + operator)
	}
}

func booleanFunctions(operator string, operands ...int) bool {
	switch operator {
	case "<":
		return operands[0] < operands[1]
	case "<=":
		return operands[0] <= operands[1]
	case "==":
		return operands[0] == operands[1]
	case "!=":
		return operands[0] != operands[1]
	case ">":
		return operands[0] > operands[1]
	case ">=":
		return operands[0] >= operands[1]
	default:
		panic("unknown operator " + operator)
	}
}

type result struct {
	highestFinal  int
	highestDuring int
}

func fn(commands string) result {
	var r result
	registers := make(registersType)

	scanner := bufio.NewScanner(strings.NewReader(commands))

	for scanner.Scan() {
		commandString := scanner.Text()
		// fmt.Println(commandString)
		currentCommand := tokenize(commandString)
		// fmt.Println(currentCommand)

		var currentValue int
		var ok bool

		if currentValue, ok = registers[currentCommand.register]; !ok {
			registers[currentCommand.register] = 0
		}

		if booleanFunctions(currentCommand.cond.operator,
			evaluate(currentCommand.cond.lefthand, registers),
			evaluate(currentCommand.cond.righthand, registers)) {
			newValue := integerFunctions(currentCommand.operator,
				currentValue, evaluate(currentCommand.value, registers))
			registers[currentCommand.register] = newValue

			if r.highestDuring < newValue {
				r.highestDuring = newValue
			}
		}
	}

	// fmt.Println(registers)

	r.highestFinal = -1

	for _, v := range registers {
		if r.highestFinal < v {
			r.highestFinal = v
		}
	}

	return r
}

type testpair struct {
	input  string
	output result
}

var testInputs = []testpair{
	{
		`b inc 5 if a > 1
a inc 1 if b < 5
c dec -10 if a >= 1
c inc -20 if c == 10`, result{1, 10}},
	// real input, real result
	{getRealInput(), result{6828, 7234}},
}

func TestOne(t *testing.T) {
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

func getRealInput() string {
	bytes, readFileErr := ioutil.ReadFile("eight_input")

	if readFileErr != nil {
		panic(readFileErr)
	}

	return string(bytes)
}
