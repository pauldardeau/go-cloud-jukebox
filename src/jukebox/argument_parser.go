package jukebox

import (
   "fmt"
   "strconv"
   "strings"
)

type ArgumentParser struct {
   dict_all_reserved_words map[string]string
   dict_bool_options map[string]string
   dict_int_options map[string]string
   dict_string_options map[string]string
   dict_commands map[string]string
   list_commands []string
}


const TYPE_BOOL = "bool"
const TYPE_INT = "int"
const TYPE_STRING = "string"


func NewArgumentParser() *ArgumentParser {
   var arg_parser ArgumentParser
   arg_parser.dict_all_reserved_words = make(map[string]string)
   arg_parser.dict_bool_options = make(map[string]string)
   arg_parser.dict_int_options = make(map[string]string)
   arg_parser.dict_string_options = make(map[string]string)
   arg_parser.dict_commands = make(map[string]string)
   arg_parser.list_commands = []string{}
   return &arg_parser
}

func (ap *ArgumentParser) addOption(o string,
                                    option_type string,
                                    help string) {
    ap.dict_all_reserved_words[o] = option_type

    if (option_type == TYPE_BOOL) {
        ap.dict_bool_options[o] = help
    } else if (option_type == TYPE_INT) {
        ap.dict_int_options[o] = help
    } else if (option_type == TYPE_STRING) {
        ap.dict_string_options[o] = help
    }
}

func (ap *ArgumentParser) AddOptionalBoolFlag(flag string, help string) {
    ap.addOption(flag, TYPE_BOOL, help)
}

func (ap *ArgumentParser) AddOptionalIntArgument(arg string, help string) {
    ap.addOption(arg, TYPE_INT, help)
}

func (ap *ArgumentParser) AddOptionalStringArgument(arg string, help string) {
    ap.addOption(arg, TYPE_STRING, help)
}

func (ap *ArgumentParser) AddRequiredArgument(arg string, help string) {
    ap.dict_commands[arg] = help
    ap.list_commands = append(ap.list_commands, arg)
}

func (ap *ArgumentParser) ParseArgs(args []string) map[string]interface{} {

    dict_args := make(map[string]interface{})

    num_args := len(args)
    working := true
    i := 0
    commands_found := 0

    if num_args == 0 {
        working = false
    }

    for {
        if ! working {
            break
        }

        arg := args[i]

        dict_value, ok := ap.dict_all_reserved_words[arg]

        if ok {
            arg_type := dict_value
            arg = arg[2:]
            if arg_type == TYPE_BOOL {
                fmt.Printf("adding key=%s value=true\n", arg)
                dict_args[arg] = true
            } else if arg_type == TYPE_INT {
                i++
                if i < num_args {
                    next_arg := args[i]
                    int_value, int_err := strconv.Atoi(next_arg)
                    if int_err == nil {
                        fmt.Printf("adding key=%s value=%d\n", arg, int_value)
                        dict_args[arg] = int_value
                    }
                } else {
                    // missing int value
                }
            } else if arg_type == TYPE_STRING {
                i++
                if i < num_args {
                    next_arg := args[i]
                    fmt.Printf("adding key=%s value=%s\n", arg, next_arg)
                    dict_args[arg] = next_arg
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
                if commands_found < len(ap.list_commands) {
                    command_name := ap.list_commands[commands_found]
                    fmt.Printf("adding key=%s value=%s\n", command_name, arg)
                    dict_args[command_name] = arg
                    commands_found++
                } else {
                    // unrecognized command
                }
            }
        }

        i++
        if i >= num_args {
            working = false
        }
    }

    return dict_args
}

