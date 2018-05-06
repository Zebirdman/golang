package main

import (
	"fmt"
)

const (
	appName = "dimager"
)

type option struct {
	Name    string   // how the argument will be passed ie '-v'
	reqArg  bool     // will it require an argument?
	Operand string   // storage for the operand once processed
	Coargs  []string // codependant arguments list
	Enabled bool     // enables the option
}

var (
	// option none, used for when no arguments are passed
	none = &option{"none", false, "", nil, false}
	// storage for a bad option passed
	invalid = &option{"", false, "", nil, false}
	// triggers the help page
	help = &option{"help", false, "", nil, false}

	// option for specifying a docker host
	host = &option{"h", true, "", []string{"-p"}, false}
	// option for docker cert path when using tls
	path = &option{"c", true, "", []string{"-h"}, false}
	// option for specifying that an image tag prefix should be added
	addP = &option{"p", true, "", nil, false}
	// option to specify that existing image tag prefix will be replaced
	repP = &option{"r", true, "", nil, false}
	// option to specify that the existing image tags be removed after rename
	clean = &option{"d", false, "", nil, false}
	// option for verbose output from operation
	verb = &option{"v", false, "", nil, false}
)

// string array for cycling through arg checks
var cmdOptions = [6]*option{host, path, addP, repP, clean, verb}

func needsHelp(a []string) bool {
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
	if !needsHelp(allArgs) {
		a := allArgs[1:]
		for i := 0; i < len(a); i++ {
			argbytes := []byte(a[i])
			for x := 0; x < len(argbytes); x++ {
				//fmt.Print(argbytes)
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
				if valOpt.reqArg {
					if x+1 != len(argbytes) {
						valOpt.Operand = string(argbytes[x+1])
						return valOpt, fmt.Errorf("invalid parameter")
					}
					if i+1 == len(a) {
						return valOpt, fmt.Errorf("requires an argument")
					}
					if nb := []byte(a[i+1]); nb[0] == '-' {
						valOpt.Operand = string(nb[0])
						return valOpt, fmt.Errorf("invalid parameter")
					}
					valOpt.Operand = string(a[i+1])
					i++
				}
			}
		}
	}
	return none, nil
}

func showErrors(o *option, e error) {
	if help.Enabled {
		helpPage()
	} else {
		switch e.Error() {
		case "invalid option":
			fmt.Printf("%s: %s '%s'\n", appName, e, o.Operand)
			fmt.Printf("Try '%s --help' for more information\n", appName)
			break
		case "inappropriate":
			fmt.Printf("%s: %s '%s'\n", appName, e, o.Operand)
			fmt.Printf("Try '%s --help' for more information\n", appName)
			break
		case "unknown command":
			fmt.Printf("%s: %s '%s'\n", appName, e, o.Operand)
			fmt.Printf("Try '%s --help' for more information\n", appName)
			break
		case "invalid parameter":
			fmt.Printf("%s: command '%s' %s '%s'\n", appName, o.Name, e, o.Operand)
			fmt.Printf("Try '%s --help' for more information\n", appName)
			break
		case "requires an argument":
			fmt.Printf("%s: option '%s' %s\n", appName, o.Name, e)
			fmt.Printf("Try '%s --help' for more information\n", appName)
			break
		}
	}
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

func helpPage() {
	fmt.Printf(`dimager: allows for the easy renaming of docker image tags prefix's
  usefull if we want to retag images to use with a pivate registry
  author: Ben Futterleib

Usage: dimager [OPTION]... [-s] SCRIPT_NAME (1st form)
  or: dimager [OPTION]... [-h] DOCKER_HOST [-p] DOCKER_CERT_PATH (2nd form)
In the first form specify a path to a script containing exports for the
  DOCKER_HOST, DOCKER_CERT_PATH, DOCKER_TLS_VERIFY
environment variables, must be valid executable script
In the second form pass the DOCKER_HOST and DOCKER_CERT_PATH values using
  the given flags, DOCKER_TLS_VERIFY will be set automatically

Arguments:
  -h 	specify the connection for the docker host
  -p 	specify the path to the directory holding the client certs for docker
  -s 	give a path to a script containing the relevant export env variables
  --help 	display this help and exit
`)
}
