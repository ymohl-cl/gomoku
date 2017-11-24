package ruler

import "fmt"

// PrintCaptures : _
func (r *Rules) PrintCaptures() {
	for _, capture := range r.captures {
		fmt.Println("(y - x): ", capture.Y, " - ", capture.X)
	}
}

// PrintThrees : _
func (r *Rules) PrintThrees() {
	fmt.Println("nb threes ", r.NumberThree)
}

// PrintAlignments : _
func (r *Rules) PrintAlignments() {
	fmt.Println("nb alignment ", len(r.aligns))
	for _, a := range r.aligns {
		a.Print()
	}
}

// PrintMove player, and spot position
func (r *Rules) PrintMove() {
	fmt.Println("player: ", r.player, " (y: ", r.y, " - x: ", r.x, ")")
}
