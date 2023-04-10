package jukebox

import (
	"fmt"
	"strconv"
	"strings"
)

type ArgumentParser struct {
	dictAllReservedWords map[string]string
	dictBoolOptions      map[string]string
	dictIntOptions       map[string]string
	dictStringOptions    map[string]string
	dictCommands         map[string]string
	listCommands         []string
	debugMode            bool
}

const TypeBool = "bool"
const TypeInt = "int"
const TypeString = "string"

func NewArgumentParser(debugMode bool) *ArgumentParser {
	var argParser ArgumentParser
	argParser.dictAllReservedWords = make(map[string]string)
	argParser.dictBoolOptions = make(map[string]string)
	argParser.dictIntOptions = make(map[string]string)
	argParser.dictStringOptions = make(map[string]string)
	argParser.dictCommands = make(map[string]string)
	argParser.listCommands = []string{}
	argParser.debugMode = debugMode
	return &argParser
}

func (ap *ArgumentParser) addOption(o string,
	optionType string,
	help string) {
	ap.dictAllReservedWords[o] = optionType

	if optionType == TypeBool {
		ap.dictBoolOptions[o] = help
	} else if optionType == TypeInt {
		ap.dictIntOptions[o] = help
	} else if optionType == TypeString {
		ap.dictStringOptions[o] = help
	}
}

func (ap *ArgumentParser) AddOptionalBoolFlag(flag string, help string) {
	ap.addOption(flag, TypeBool, help)
}

func (ap *ArgumentParser) AddOptionalIntArgument(arg string, help string) {
	ap.addOption(arg, TypeInt, help)
}

func (ap *ArgumentParser) AddOptionalStringArgument(arg string, help string) {
	ap.addOption(arg, TypeString, help)
}

func (ap *ArgumentParser) AddRequiredArgument(arg string, help string) {
	ap.dictCommands[arg] = help
	ap.listCommands = append(ap.listCommands, arg)
}

func (ap *ArgumentParser) ParseArgs(args []string) *PropertySet {

	ps := NewPropertySet()

	numArgs := len(args)
	working := true
	i := 0
	commandsFound := 0

	if numArgs == 0 {
		working = false
	}

	for {
		if !working {
			break
		}

		arg := args[i]

		dictValue, ok := ap.dictAllReservedWords[arg]

		if ok {
			argType := dictValue
			arg = arg[2:]
			if argType == TypeBool {
				if ap.debugMode {
					fmt.Printf("adding key=%s value=true\n", arg)
				}
				ps.Add(arg, NewBoolPropertyValue(true))
			} else if argType == TypeInt {
				i++
				if i < numArgs {
					nextArg := args[i]
					intValue, intErr := strconv.Atoi(nextArg)
					if intErr == nil {
						if ap.debugMode {
							fmt.Printf("adding key=%s value=%d\n", arg, intValue)
						}
						ps.Add(arg, NewIntPropertyValue(intValue))
					}
				} else {
					// missing int value
				}
			} else if argType == TypeString {
				i++
				if i < numArgs {
					nextArg := args[i]
					if ap.debugMode {
						fmt.Printf("adding key=%s value=%s\n", arg, nextArg)
					}
					ps.Add(arg, NewStringPropertyValue(nextArg))
				} else {
					// missing string value
				}
			} else {
				// unrecognized type
			}
		} else {
			if strings.HasPrefix(arg, "--") {
				// unrecognized option
			} else {
				if commandsFound < len(ap.listCommands) {
					commandName := ap.listCommands[commandsFound]
					if ap.debugMode {
						fmt.Printf("adding key=%s value=%s\n", commandName, arg)
					}
					ps.Add(commandName, NewStringPropertyValue(arg))
					commandsFound++
				} else {
					// unrecognized command
				}
			}
		}

		i++
		if i >= numArgs {
			working = false
		}
	}

	return ps
}
