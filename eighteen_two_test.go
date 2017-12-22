package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
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

func printRegisters(r registersType, printer printfFunc) {
	printer("last\t%v\n", r[lastSoundRegisterName])

	sortedKeys := make([]string, 0)

	for key := range r {
		if key != lastSoundRegisterName {
			sortedKeys = append(sortedKeys, key)
		}
	}

	sort.Strings(sortedKeys)

	for _, registerName := range sortedKeys {
		printer("%v\t%v\n", registerName, r[registerName])
	}
	fmt.Println()
}

type printfFunc (func(format string, a ...interface{}))

func getWrappedPrintf(programID int) printfFunc {
	return func(format string, a ...interface{}) {
		modifiedFormat := fmt.Sprintf("PROG-%v\t %v", programID, format)
		fmt.Printf(modifiedFormat, a...)
	}
}

func doIt(programId int, commands []command,
	waitGroup *sync.WaitGroup,
	sendChan chan int, receiveChan chan int, sendCounterChan chan int) {

	var registers registersType
	wrappedPrintf := getWrappedPrintf(programId)

	getValueOfArgument := func(arg argument) int {
		if arg.isLiteral {
			return arg.literalValue
		}

		return registers[arg.value]
	}

	registers = make(registersType)
	registers[lastSoundRegisterName] = -1
	registers["p"] = programId

	for executionPointer := 0; executionPointer >= 0 && executionPointer < len(commands); executionPointer++ {
		currentCommand := commands[executionPointer]

		switch currentCommand.operator {
		case "set":
			targetRegister := currentCommand.arguments[0].value
			newValue := getValueOfArgument(currentCommand.arguments[1])
			registers[targetRegister] = newValue

			wrappedPrintf("\tset %v = %v\n", targetRegister, newValue)
			// default:
		case "add":
			basisToAddTo := registers[currentCommand.arguments[0].value]
			valueToAdd := getValueOfArgument(currentCommand.arguments[1])

			wrappedPrintf("\t%v + %v\n", basisToAddTo, valueToAdd)
			registers[currentCommand.arguments[0].value] = basisToAddTo + valueToAdd
		case "mul":
			basisToMultiply := registers[currentCommand.arguments[0].value]
			valueToMultiply := getValueOfArgument(currentCommand.arguments[1])

			wrappedPrintf("\t%v * %v\n", basisToMultiply, valueToMultiply)
			registers[currentCommand.arguments[0].value] = basisToMultiply * valueToMultiply
		case "mod":
			dividend := registers[currentCommand.arguments[0].value]
			divisor := getValueOfArgument(currentCommand.arguments[1])

			wrappedPrintf("\t%v %% %v\n", dividend, divisor)
			registers[currentCommand.arguments[0].value] = dividend % divisor
		case "snd":
			frequencyToPlay := getValueOfArgument(currentCommand.arguments[0])
			registers[lastSoundRegisterName] = frequencyToPlay
			go func() {
				wrappedPrintf("send %v\n", frequencyToPlay)
				sendChan <- frequencyToPlay
				if sendCounterChan != nil {
					wrappedPrintf("send counter %v\n", frequencyToPlay)
					sendCounterChan <- frequencyToPlay
				}
			}()
		case "rcv":
			targetRegister := currentCommand.arguments[0].value
			wrappedPrintf("\t\tabout to receive into %v\n", targetRegister)

			// hacky receive deadlock preventer
			mightBeHung := true
			go func(kill *bool) {
				time.Sleep(100 * time.Millisecond)

				if *kill {
					receiveChan <- -100

				}
			}(&mightBeHung)

			valueReceived := <-receiveChan
			// hacky receive deadlock preventer
			mightBeHung = false

			if valueReceived == -100 {
				wrappedPrintf("Program %v is deadlocked\n", programId)
				waitGroup.Done()
				return
			}
			wrappedPrintf("\t\treceived into %v\n", targetRegister)
			registers[targetRegister] = valueReceived
		case "jgz":
			shouldAct := getValueOfArgument(currentCommand.arguments[0])

			wrappedPrintf("\tjump? %v => %v\n", currentCommand.arguments[0].value, shouldAct)
			if shouldAct != 0 {

				wrappedPrintf("\tjump from %v to ", executionPointer)
				targetInstruction := executionPointer + getValueOfArgument(currentCommand.arguments[1])

				// subtract one so we can let the for-loop after though increment
				// the expression pointer to our actual target
				executionPointer = targetInstruction - 1
				wrappedPrintf("%v\n", targetInstruction)
			}
		default:
			panic(fmt.Sprintf("Unknown operator %v\n", currentCommand.operator))
		}

		// printRegisters(registers, wrappedPrintf)
	}
	waitGroup.Done()
	close(sendChan)
	return
}

func fn(input string) int {
	var commands []command

	commands = convertStringToCommands(input)
	// fmt.Println(commands)

	var waitGroup sync.WaitGroup

	zeroChan := make(chan int)
	oneChan := make(chan int)

	sendCounterChan := make(chan int)
	sendCounter := 0
	// finalResulChan := make(chan int)

	go func() {
		for s := range sendCounterChan {
			sendCounter++
			fmt.Printf("<<<<-%v %v->>>>\n", s, sendCounter)
		}
		// finalResulChan <- sendCounter
	}()

	waitGroup.Add(1)
	go doIt(0, commands, &waitGroup, zeroChan, oneChan, nil)
	waitGroup.Add(1)
	go doIt(1, commands, &waitGroup, oneChan, zeroChan, sendCounterChan)

	waitGroup.Wait()
	close(sendCounterChan)
	// finalResult := <-finalResulChan

	return sendCounter
}

type testpair struct {
	input  string
	output int
}

var testInputs = []testpair{
	// {
	// 	`
	// snd 1
	// snd 2
	// snd p
	// rcv a
	// rcv b
	// rcv c
	// rcv d
	// 	`, 3},
	// real input, real result
	{getRealInput(), -1},
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
