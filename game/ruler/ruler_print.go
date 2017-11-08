package ruler

import "fmt"

func (r *Rules) printCaptures() {
	for _, capture := range r.captures {
		fmt.Println("(y - x): ", capture.Y, " - ", capture.X)
	}
}

func (r *Rules) printThrees() {
	fmt.Println("nb threes ", r.NumberThree)
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
		if a.iWin != nil {
			if r.Win == false {
				fmt.Print(", capturable")
			} else {
				fmt.Print(", no capturable")
			}
		}
		fmt.Println("]")
	}
}
