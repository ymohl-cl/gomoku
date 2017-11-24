package alignment

import (
	"fmt"

	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

// Print : _
func (a Alignment) Print() {
	fmt.Print("[size: ", a.Size)
	if a.Style&rdef.AlignFree != 0 {
		fmt.Print(", style: free")
	} else if a.Style&rdef.AlignHalf != 0 {
		fmt.Print(", style: half")
	} else if a.Style&rdef.AlignFlanked != 0 {
		fmt.Print(", style: flanked")
	} else {
		fmt.Print(", style: not define")
	}
	if a.iWin != nil {
		if a.iWin.IsCapturable() == false {
			fmt.Print(", no capturable")
		} else {
			fmt.Print(", capturable")
		}
	}
	fmt.Println("]")
}
