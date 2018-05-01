package main

import (
	"errors"
	"fmt"
)

const (
	appName = "dimager"
)

type option struct {
	Name    string   // how the argument will be passed ie '-v'
	Arg     bool     // will it require an operand?
	Operand string   // storage for the operand once processed
	Coargs  []string // codependant arguments list
	Enabled bool     // enables the option
}

var (
	// option none, used for when no arguments are passed
	none = &option{"", false, "", nil, false}
	// triggers the help page
	help = &option{"--help", false, "", nil, false}
	// option for specifying a docker host
	host = &option{"-h", true, "", []string{"-p"}, false}
	// option for docker cert path when using tls
	path = &option{"-c", true, "", []string{"-h"}, false}
)

// string array for cycling through arg checks
var cmdOptions = [3]*option{help, host, path}

func checkArgs(a []string) (*option, error) {
	var err error
	var option = none
	if len(a) > 1 {
		args := a[1:]

		for index, argValue := range args {

			for _, op := range cmdOptions {
				// first check name of argument matches an option name
				if argValue == op.Name {
					if op.Arg == true {
						err = checkOperand(index, args, op)
						if err != nil {
							option = op
						}
					}
					op.Enabled = true
				}
			}
		}
		// if argument doesnt match return and report

		return option, err

	}
	return option, errors.New("missing operand")
}

func checkOperand(i int, a []string, o *option) error {
	var err error
	//fmt.Printf("Index = %d , Length = %d , Arg = %s\n", i, len(a), a[i])
	if i+2 > len(a) {
		err = fmt.Errorf("missing argument")
		// if net arg is an option then we know the operand is missing
	} else {
		for _, op := range cmdOptions {
			if a[1] == op.Name && op.Name != "--help" {
				err = fmt.Errorf("missing argument")
			} else {
				op.Operand = a[1]
			}
		}
	}
	return err
}

func errorPage(s string, e error) {
	if s == "" {
		fmt.Printf("%s: %s\n", appName, e)
	} else {
		fmt.Printf("%s: %s '%s'\n", appName, e, s)
	}
	fmt.Printf("Try '%s --help' for more information\n", appName)
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
