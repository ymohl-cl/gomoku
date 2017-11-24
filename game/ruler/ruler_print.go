package ruler

import "fmt"

func (r *Rules) PrintCaptures() {
	for _, capture := range r.captures {
		fmt.Println("(y - x): ", capture.Y, " - ", capture.X)
	}
}

func (r *Rules) PrintThrees() {
	fmt.Println("nb threes ", r.NumberThree)
}

func (r *Rules) PrintAlignments() {
	fmt.Println("nb alignment ", len(r.aligns))
	for _, a := range r.aligns {
		a.Print()
	}
}
