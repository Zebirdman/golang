package main

import (
	"bytes"
	"fmt"
)

type path []byte
type numma int

func main() {
	pathName := path("/usr/bin/hi")

	pathName.Truncate()
	fmt.Printf("%s\n", pathName)
}

func (p *path) Truncate() {
	i := bytes.LastIndex(*p, []byte("/"))
	if i >= 0 {
		*p = (*p)[0:i]
	}
}
