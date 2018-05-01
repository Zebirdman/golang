package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if checkArgs(os.Args) {
		input := " "
		r := bufio.NewReader(os.Stdin)
		for input != "exit" {
			line, _ := r.ReadString('\n')
			input = line[:len(line)-1]
		}
		fmt.Print("Exiting Calculator...\n")
	}
}

func checkArgs(a []string) bool {
	if len(a) > 1 {
		if a[1] == "--help" && len(a) == 2 {
			fmt.Print("work in progress\n")
			return false
		}
		fmt.Print("calc: invalid options '")
		for index, value := range a[1:] {
			if value != "--help" {
				fmt.Printf("%s", value)
			}
			if index != (len(a) - 2) {
				fmt.Printf(" ")
			}
		}
		fmt.Print("'\nTry 'calc --help' for more information.\n")
		return false
	}
	fmt.Print("Simple Calculator: Author Ben Futterleib\n>")
	return true
}
