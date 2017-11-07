package ruler

import "fmt"

func (r *Rules) printCaptures() {
	for _, cap := range r.caps {
		fmt.Println("(y - x): ", cap.Y, " - ", cap.X)
	}
}

func (r *Rules) printThrees() {
	fmt.Println("nb threes ", r.NbThree)
}

func (r *Rules) printAlignments() {
	fmt.Println("nb alignment ", len(r.aligns))
	for _, a := range r.aligns {
		fmt.Print("[size: ", a.size)
		if a.style&alignFree != 0 {
			fmt.Print(", style: free")
		} else if a.style&alignHalf != 0 {
			fmt.Print(", style: half")
		} else if a.style&alignFlanked != 0 {
			fmt.Print(", style: flanked")
		} else {
			fmt.Print(", style: not define")
		}
		if a.size >= 5 {
			if a.capturable {
				fmt.Print(", capturable")
			} else {
				fmt.Print(", no capturable")
			}
		}
		fmt.Println("]")
	}
}

func (r *Rules) print() {
	if r.player == TokenP1 {
		fmt.Println("Player: P1")
	} else if r.player == TokenP2 {
		fmt.Println("Player: P2")
	} else {
		fmt.Println("Player not define")
	}

	fmt.Println("Can move: ", r.IsMoved, " | details: ", r.MovedStr)

	r.printCaptures()

	fmt.Println("Number tree: ", r.NbThree)

	fmt.Println("IsWin: ", r.IsWin, " - details: ", r.MessageWin)

	r.printAlignments()
}
