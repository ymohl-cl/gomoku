package alignment

import (
	"fmt"
	"testing"

	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

func ExampleInfoWin_GetMask() {
	i := new(InfoWin)
	i.mask = [11]uint8{2, 2, 2, 2, 1, 0, 0, 1, 3, 3, 3}
	fmt.Println(i.GetMask())
	// Output:
	//&[2 2 2 2 1 0 0 1 3 3 3]
}

func ExampleInfoWin_GetDirection() {
	i := new(InfoWin)

	i.dirY = 5
	i.dirX = 2
	y, x := i.GetDirection()
	fmt.Println("(", y, " - ", x, ")")
	// Output:
	//( 5  -  2 )
}

func TestIsCapturable(t *testing.T) {
	i := new(InfoWin)

	// test: 0
	if i.IsCapturable() == true {
		t.Error(t.Name() + " > test: 0")
	}

	i.capturable = true

	// test : 1
	if i.IsCapturable() == false {
		t.Error(t.Name() + " > test: 1")
	}
}

func ExampleAlignment_NewInfosWin() {
	var a Alignment

	mask := [11]uint8{2, 2, 2, 2, 1, 0, 0, 1, 3, 3, 3}
	a.NewInfosWin(&mask, 5, 2, false)
	i := a.GetInfosWin()
	y, x := i.GetDirection()
	fmt.Println(i.GetMask())
	fmt.Println("(", y, " - ", x, ")")
	// Output:
	//&[2 2 2 2 1 0 0 1 3 3 3]
	//( 5  -  2 )
}

func TestCopy(t *testing.T) {
	var dst Alignment
	src := Alignment{Size: 2, Style: rdef.AlignFree}

	// test: 0 > src and dst match
	dst.Copy(&src)
	if src.Size != dst.Size || src.Style != dst.Style {
		t.Error(t.Name() + " > test: 0")
	}
}

func TestIsBetter(t *testing.T) {
	var one, two Alignment

	// test: 0 > one is better by size
	one = Alignment{Size: 2, Style: rdef.AlignFree}
	two = Alignment{Size: 1, Style: rdef.AlignFree}
	if !one.IsBetter(&two) {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > one is better by style
	one = Alignment{Size: 2, Style: rdef.AlignFree}
	two = Alignment{Size: 2, Style: rdef.AlignHalf}
	if !one.IsBetter(&two) {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > two is better by size
	one = Alignment{Size: 2, Style: rdef.AlignFree}
	two = Alignment{Size: 3, Style: rdef.AlignHalf}
	if one.IsBetter(&two) {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3 > one is better by equality
	one = Alignment{Size: 3, Style: rdef.AlignHalf}
	two = Alignment{Size: 3, Style: rdef.AlignHalf}
	if !one.IsBetter(&two) {
		t.Error(t.Name() + " > test: 3")
	}

	// test: 4 > reverse test 1
	one = Alignment{Size: 3, Style: rdef.AlignHalf}
	two = Alignment{Size: 3, Style: rdef.AlignFree}
	if one.IsBetter(&two) {
		t.Error(t.Name() + " > test: 4")
	}
}

func TestSetStyleByMask(t *testing.T) {
	var mask [11]uint8
	var a Alignment

	// test: 0 > free style
	mask = [11]uint8{0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0}
	a.setStyleByMask(&mask, 2, rdef.Player1)
	if a.Style&rdef.AlignFlanked != 0 || a.Style&rdef.AlignFree == 0 || a.Style&rdef.AlignHalf != 0 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > half-free style
	a.Style = 0
	mask = [11]uint8{0, 0, 0, 0, 1, 1, 1, 2, 0, 0, 0}
	a.setStyleByMask(&mask, 2, rdef.Player1)
	if a.Style&rdef.AlignFlanked != 0 || a.Style&rdef.AlignFree != 0 || a.Style&rdef.AlignHalf == 0 {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > flanked style
	a.Style = 0
	mask = [11]uint8{3, 3, 0, 0, 1, 1, 1, 2, 0, 0, 0}
	a.setStyleByMask(&mask, 2, rdef.Player1)
	if a.Style&rdef.AlignFlanked == 0 || a.Style&rdef.AlignFree != 0 || a.Style&rdef.AlignHalf != 0 {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3 > flanked style
	a.Style = 0
	mask = [11]uint8{3, 3, 1, 1, 0, 1, 0, 1, 0, 0, 0}
	a.setStyleByMask(&mask, 2, rdef.Player1)
	if a.Style&rdef.AlignFlanked != 0 || a.Style&rdef.AlignFree != 0 || a.Style&rdef.AlignHalf == 0 {
		t.Error(t.Name() + " > test: 3")
	}

	// test: 4 > free style
	a.Style = 0
	mask = [11]uint8{3, 3, 1, 1, 0, 1, 0, 1, 0, 0, 0}
	a.setStyleByMask(&mask, 3, rdef.Player1)
	if a.Style&rdef.AlignFlanked != 0 || a.Style&rdef.AlignFree == 0 || a.Style&rdef.AlignHalf != 0 {
		t.Error(t.Name() + " > test: 4")
	}
}

func TestNew(t *testing.T) {
	var mask [11]uint8
	var a *Alignment

	// test: 0 > all spots set to the player
	mask = [11]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	a = New(&mask, rdef.Player1)
	if a.Size != 5 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 0 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 1 > check extremities out of scope
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	a = New(&mask, rdef.Player1)
	if a.Size != 1 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 1 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 2 > check extremities on the scope
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 2 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 2 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 3 > check three alignment without salt
	mask = [11]uint8{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 3 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 4 > check three alignment with salt
	mask = [11]uint8{0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 4 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 5 > check three alignment with salt. start mask by 'out of board'
	mask = [11]uint8{3, 3, 1, 1, 0, 1, 0, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 5 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 6 > check alignment cut by ennemy spot
	mask = [11]uint8{3, 3, 1, 1, 2, 1, 0, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 1 || a.Style&rdef.AlignHalf == 0 {
		t.Error(t.Name()+" > test: 6 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 7 > reverse test 6
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 2, 1, 1, 3, 3}
	a = New(&mask, rdef.Player1)
	if a.Size != 1 || a.Style&rdef.AlignHalf == 0 {
		t.Error(t.Name()+" > test: 7 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 8 > check with not enought space to alignment.
	mask = [11]uint8{0, 1, 1, 2, 0, 1, 0, 0, 2, 3, 3}
	a = New(&mask, rdef.Player1)
	if a.Size != 0 {
		t.Error(t.Name()+" > test: 8 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 9 > 3 token aligned with salt
	mask = [11]uint8{0, 0, 1, 0, 0, 1, 1, 1, 2, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 9 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 10 > 2 test vrac
	mask = [11]uint8{3, 1, 1, 0, 0, 1, 1, 2, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 3 || a.Style&rdef.AlignHalf == 0 {
		t.Error(t.Name()+" > test: 10 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 11 > test vrac
	mask = [11]uint8{3, 3, 2, 0, 2, 1, 0, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.Size != 2 || a.Style&rdef.AlignHalf == 0 {
		t.Error(t.Name()+" > test: 11 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}

	// test: 12 > test vrac
	mask = [11]uint8{0, 0, 0, 0, 1, 1, 1, 0, 3, 3, 3}
	a = New(&mask, rdef.Player1)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 {
		t.Error(t.Name()+" > test: 12 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
}

func TestAnalyzeThree(t *testing.T) {
	var mask [11]uint8
	var a *Alignment

	// test: 0
	mask = [11]uint8{0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 0 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 1
	mask = [11]uint8{0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 1 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 2
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 2 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 3
	mask = [11]uint8{0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 3 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 4
	mask = [11]uint8{0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 4 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 5
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 5 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 6
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 6 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 7
	mask = [11]uint8{0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 7 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 8
	mask = [11]uint8{0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 8 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 9
	mask = [11]uint8{0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree == true {
		t.Error(t.Name()+" > test: 9 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
	// test: 10
	mask = [11]uint8{0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	a.AnalyzeThree(&mask)
	if a.Size != 3 || a.Style&rdef.AlignFree == 0 || a.IsThree != true {
		t.Error(t.Name()+" > test: 10 > (Size: ", a.Size, " - Style: ", a.Style, ")")
	}
}

func TestIsConsecutive(t *testing.T) {
	var mask [11]uint8
	var a *Alignment

	// test: 0
	mask = [11]uint8{0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if !a.IsConsecutive(&mask) {
		t.Error(t.Name() + " > test: 0")

	}

	// test: 1
	mask = [11]uint8{0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.IsConsecutive(&mask) {
		t.Error(t.Name() + " > test: 1")

	}
	// test: 2
	mask = [11]uint8{0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0}
	a = New(&mask, rdef.Player1)
	if a.IsConsecutive(&mask) {
		t.Error(t.Name() + " > test: 2")

	}

	// test: 3
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0}
	a = New(&mask, rdef.Player1)
	if a.IsConsecutive(&mask) {
		t.Error(t.Name() + " > test: 3")

	}

	// test: 4
	mask = [11]uint8{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0}
	a = New(&mask, rdef.Player1)
	if !a.IsConsecutive(&mask) {
		t.Error(t.Name() + " > test: 4")

	}
}
