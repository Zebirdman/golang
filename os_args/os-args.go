package main

import (
	"fmt"
	"os"
)

func main() {
	var arg string
	path := os.Args[0]

	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	fmt.Print(path)
	fmt.Print("\n")
	switch arg {
	case "run":
		fmt.Printf("You have run\n")
	case "stop":
		fmt.Print("You have stopped\n")
	case "--help":
		fmt.Print("help page")
	default:
		fmt.Print("You did neither\n")
	}
}
