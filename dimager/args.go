package main

import (
	"fmt"
)

type option struct {
	Name    string   // how the argument will be passed ie '-v'
	reqArg  bool     // will it require an argument?
	Operand string   // storage for the operand once processed
	Coargs  []string // codependant arguments list
	Enabled bool     // enables the option
}

var (
	// application name
	appName = "Your app"
	// option none, used for when no arguments are passed
	none = &option{"none", false, "", nil, false}
	// storage for a bad option passed
	invalid = &option{"", false, "", nil, false}
	// triggers the help page
	help = &option{"help", false, "", nil, false}
	// array for holding cmd line cmd line options
	cmdOptions = make([]*option, 0, 10)
	// holds the help page string
	appHelpPage = ""
	// debug options
	debug = false
)

func initArgs(n, h string) {
	appHelpPage = h
	appName = n
	debugArguments()
}

func newOption(n string, r bool, c []string) *option {
	no := &option{n, r, "", c, false}
	cmdOptions = append(cmdOptions, no)
	return no
}

func debugArguments() {
	if debug {
		// show enabled options and Arguments
		fmt.Printf("Enabled arguments:\n")
		for _, opt := range cmdOptions {
			if opt.Enabled {
				fmt.Printf("Name: %s  Operand: %s\n", opt.Name, opt.Operand)
			}
		}
	}
}

func needHelp(a []string) bool {
	if len(a) > 1 {
		for _, value := range a {
			if value == "--help" {
				help.Enabled = true
				return true
			}
		}
	}
	return false
}

func checkArgs(allArgs []string) (*option, error) {
	if !needHelp(allArgs) && len(allArgs) > 1 {
		a := allArgs[1:]
		for i := 0; i < len(a); i++ {
			argbytes := []byte(a[i])
			for x := 0; x < len(argbytes); x++ {
				if len(argbytes) == 1 {
					if argbytes[0] == '-' {
						invalid.Operand = string(argbytes[0])
						return invalid, fmt.Errorf("inappropriate")
					}
					invalid.Operand = string(argbytes[0])
					return invalid, fmt.Errorf("invalid option")
				}
				if x == 0 {
					x = 1
				}
				valOpt, err := checkCommand(argbytes[x])
				if err != nil {
					return invalid, fmt.Errorf("unknown command")
				}
				// if argument required check that it is valid
				if valOpt.reqArg {
					if x+1 != len(argbytes) {
						valOpt.Operand = string(argbytes[x+1])
						return valOpt, fmt.Errorf("invalid parameter")
					}
					if i+1 == len(a) {
						return valOpt, fmt.Errorf("requires an argument")
					}
					// TODO: need to account fo someone putting in a "" string
					if nb := []byte(a[i+1]); nb[0] == '-' {
						valOpt.Operand = string(nb[0])
						return valOpt, fmt.Errorf("invalid parameter")
					}
					valOpt.Operand = string(a[i+1])
					i++
				}
			}
		}
		return nil, nil
	}
	return none, fmt.Errorf("missing arguments")
}

func checkCommand(b byte) (*option, error) {
	for _, cmd := range cmdOptions {
		if string(b) == cmd.Name {
			cmd.Enabled = true
			return cmd, nil
		}
	}
	invalid.Operand = string(b)
	return invalid, fmt.Errorf("invalid command")
}

func showErrors(o *option, e error) {
	if help.Enabled {
		fmt.Printf(appHelpPage)
		return
	}
	switch e.Error() {
	case "missing arguments":
		fmt.Printf("%s: %s\n", appName, e)
		break
	case "invalid parameter":
		fmt.Printf("%s: command '%s' %s '%s'\n", appName, o.Name, e, o.Operand)
		break
	case "requires an argument":
		fmt.Printf("%s: option '%s' %s\n", appName, o.Name, e)
		break
	default:
		fmt.Printf("%s: %s '%s'\n", appName, e, o.Operand)
		break
	}
	fmt.Printf("Try '%s --help' for more information\n", appName)
}
