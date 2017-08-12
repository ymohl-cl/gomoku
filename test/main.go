package main

import "fmt"

// Prout : _
type Prout struct {
	Name string
}

// IFace : _
type IFace interface {
	Printer()
	ChangeName()
}

// Toto : _
type Toto struct {
	*Prout
}

// Tata : _
type Tata struct {
	*Prout
}

// Printer Test template
func (p Prout) Printer() {
	fmt.Println(p.Name)
}

// ChangeName : _
func (p *Prout) ChangeName() {
	p.Name = "Yolo"
}

func main() {
	to := Toto{Prout: new(Prout)}
	ta := Tata{Prout: new(Prout)}
	var io IFace
	var ia IFace

	io = to
	ia = ta
	//	pr := Prout{age: 10, sexe: "H"}
	//	to.Name = "toto"
	//	ta.Name = "tata"
	io.Printer()
	ia.Printer()
	io.ChangeName()
	ia.ChangeName()
	io.Printer()
	ia.Printer()
	to.Printer()
	to.Name = "Machin"
	to.Printer()
}
