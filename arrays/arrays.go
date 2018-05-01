package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type arrayBuilder struct {
	name string
	*bufio.Reader
	arr []int
}

func main() {
	array1 := &arrayBuilder{"array 1", makeRead(), make([]int, 0, 10)}
	array2 := &arrayBuilder{"array 2", makeRead(), make([]int, 0, 10)}
	array1.create()
	array2.create()
	array1.display()
	array2.display()
	array1.showSlice()
	array2.showSlice()
}

func (a *arrayBuilder) create() {
	for {
		fmt.Printf("Enter a number for %s, or -1 to exit\n", a.name)
		text, _ := a.ReadString('\n')
		input, _ := strconv.Atoi(text[0 : len(text)-1])
		if input != -1 {
			a.arr = append(a.arr, input)
			fmt.Printf("You added %d\n", input)
		} else {
			break
		}
	}
}

func (a *arrayBuilder) display() {
	fmt.Printf("%s is %d long\n", a.name, len(a.arr))
	for index, value := range a.arr {
		fmt.Printf("Index %d contains %d\n", index, value)
	}
}

func (a *arrayBuilder) showSlice() {
	fmt.Printf("length = %d Capacity = %d\n", len(a.arr), cap(a.arr))
}

func makeRead() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}
