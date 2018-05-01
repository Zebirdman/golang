package main

import "fmt"

func main() {
	dad := &Parent{&Person{"John", 33}, "male"}
	mum := &Parent{&Person{"Judy", 30}, "female"}
	son := &Child{&Person{"Brett", 12}, "male"}
	daughter := &Child{&Person{"Clair", 10}, "female"}

	fam := &Family{dad, mum, son, daughter}
	fam.getDetails()
	dad.getDetails()
	fam.Mum.getDetails()
}

type Person struct {
	Name string
	Age  int
}

type Family struct {
	Dad      *Parent
	Mum      *Parent
	Son      *Child
	Daughter *Child
}

type Parent struct {
	*Person
	Sex string
}

type Child struct {
	*Person
	Sex string
}

func (p *Family) getDetails() {
	fmt.Printf("%s is %d years old\n", p.Dad.Name, p.Dad.Age)
	fmt.Printf("%s is %d years old\n", p.Mum.Name, p.Mum.Age)
	fmt.Printf("%s is %d years old\n", p.Son.Name, p.Son.Age)
	fmt.Printf("%s is %d years old\n", p.Daughter.Name, p.Daughter.Age)
}

func (p *Person) getDetails() {
	fmt.Printf("%s is %d years old\n", p.Name, p.Age)
}
