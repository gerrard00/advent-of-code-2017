package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type registersType map[string]int

const lastSoundRegisterName = "last_sound"

type argument struct {
	value        string
	isLiteral    bool
	literalValue int
}

// b inc 5 if a > 1
type command struct {
	operator  string
	arguments []argument
}

func convertStringToCommands(commandString string) []command {
	// jgz a -2
	commands := make([]command, 0)

	scanner := bufio.NewScanner(strings.NewReader(commandString))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 {
			continue
		}

		newCommand := command{}

		tokens := strings.Split(line, " ")

		newCommand.operator = tokens[0]
		newCommand.arguments = make([]argument, 0)

		for tokenIndex := 1; tokenIndex < len(tokens); tokenIndex++ {
			newArgument := argument{
				value: tokens[tokenIndex],
			}

			if literalValue, conversionError := strconv.Atoi(newArgument.value); conversionError == nil {
				newArgument.literalValue = literalValue
				newArgument.isLiteral = true
			} else {
				newArgument.isLiteral = false
			}

			newCommand.arguments = append(newCommand.arguments, newArgument)
		}

		commands = append(commands, newCommand)

	}
	return commands
}

func printRegisters(r registersType) {
	fmt.Printf("last\t%v\n", r[lastSoundRegisterName])

	sortedKeys := make([]string, 0)

	for key := range r {
		if key != lastSoundRegisterName {
			sortedKeys = append(sortedKeys, key)
		}
	}

	sort.Strings(sortedKeys)

	for _, registerName := range sortedKeys {
		fmt.Printf("%v\t%v\n", registerName, r[registerName])
	}
	fmt.Println()
}

func fn(input string) int {
	var registers registersType

	getValueOfArgument := func(arg argument) int {
		if arg.isLiteral {
			return arg.literalValue
		}

		return registers[arg.value]
	}

	commands := convertStringToCommands(input)
	fmt.Println(commands)

	registers = make(registersType)
	registers[lastSoundRegisterName] = -1

	for executionPointer := 0; executionPointer >= 0 && executionPointer < len(commands); executionPointer++ {
		currentCommand := commands[executionPointer]

		switch currentCommand.operator {
		case "set":
			targetRegister := currentCommand.arguments[0].value
			newValue := getValueOfArgument(currentCommand.arguments[1])
			registers[targetRegister] = newValue

			fmt.Printf("\tset %v = %v\n", targetRegister, newValue)
			// default:
		case "add":
			basisToAddTo := registers[currentCommand.arguments[0].value]
			valueToAdd := getValueOfArgument(currentCommand.arguments[1])

			fmt.Printf("\t%v + %v\n", basisToAddTo, valueToAdd)
			registers[currentCommand.arguments[0].value] = basisToAddTo + valueToAdd
		case "mul":
			basisToMultiply := registers[currentCommand.arguments[0].value]
			valueToMultiply := getValueOfArgument(currentCommand.arguments[1])

			fmt.Printf("\t%v * %v\n", basisToMultiply, valueToMultiply)
			registers[currentCommand.arguments[0].value] = basisToMultiply * valueToMultiply
		case "mod":
			dividend := registers[currentCommand.arguments[0].value]
			divisor := getValueOfArgument(currentCommand.arguments[1])

			fmt.Printf("\t%v %% %v\n", dividend, divisor)
			registers[currentCommand.arguments[0].value] = dividend % divisor
		case "snd":
			frequencyToPlay := getValueOfArgument(currentCommand.arguments[0])
			fmt.Printf(">>>>>>>>>>>>>>>%v<<<<<<<<<<<<<<<\n", frequencyToPlay)
			registers[lastSoundRegisterName] = frequencyToPlay
		case "rcv":
			shouldAct := getValueOfArgument(currentCommand.arguments[0])

			if shouldAct != 0 {
				valToRecover := registers[lastSoundRegisterName]
				if valToRecover > 0 {
					// TODO: what does recover mean?
					fmt.Printf("Recovered %v\n", valToRecover)
					return valToRecover
				}
			}
		case "jgz":
			shouldAct := getValueOfArgument(currentCommand.arguments[0])

			fmt.Printf("\tjump? %v => %v\n", currentCommand.arguments[0].value, shouldAct)
			if shouldAct != 0 {

				fmt.Printf("\tjump from %v to ", executionPointer)
				targetInstruction := executionPointer + getValueOfArgument(currentCommand.arguments[1])

				// subtract one so we can let the for-loop after though increment
				// the expression pointer to our actual target
				executionPointer = targetInstruction - 1
				fmt.Printf("%v\n", targetInstruction)
			}
		default:
			panic(fmt.Sprintf("Unknown operator %v\n", currentCommand.operator))
		}

		printRegisters(registers)
	}

	return registers[lastSoundRegisterName]
}

type testpair struct {
	input  string
	output int
}

var testInputs = []testpair{
	{
		`
set a 1
add a 2
mul a a
mod a 5
snd a
set a 0
rcv a
jgz a -1
set a 1
jgz a -2
		`, 4},
	// real input, real result
	{getRealInput(), 4601},
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
	bytes, readFileErr := ioutil.ReadFile("eighteen_input")

	if readFileErr != nil {
		panic(readFileErr)
	}

	return string(bytes)
}
