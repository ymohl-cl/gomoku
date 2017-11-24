package alignment

import rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"

// InfoWin keep info of raw when a align size == 5.
// Need to check if the win align is capturable or not
type InfoWin struct {
	mask       [11]uint8
	dirY       int8
	dirX       int8
	capturable bool
}

// GetMask return mask from info win
func (i InfoWin) GetMask() *[11]uint8 {
	return &i.mask
}

// GetDirection : _
func (i InfoWin) GetDirection() (int8, int8) {
	return i.dirY, i.dirX
}

// IsCapturable : _
func (i InfoWin) IsCapturable() bool {
	return i.capturable
}

// Alignment describe the alignments to five available spot and size align > 2
type Alignment struct {
	Size    uint8
	Style   uint8
	IsThree bool
	iWin    *InfoWin
}

// GetInfosWin : _
func (a Alignment) GetInfosWin() *InfoWin {
	return a.iWin
}

// NewInfosWin : _
func (a *Alignment) NewInfosWin(mask *[11]uint8, dirY, dirX int8, capturable bool) {
	a.iWin = new(InfoWin)
	copy(a.iWin.mask[:], (*mask)[:])
	a.iWin.dirY = dirY
	a.iWin.dirX = dirX
	a.iWin.capturable = capturable
}

/*
// GetSize : _
func (a *Align) GetSize() int8 {
	return int8(a.Size)
}

// IsStyle : _
func (a *Align) IsStyle(style uint8) bool {
	if a.Style&style != 0 {
		return true
	}
	return false
}
*/

// Copy on the 'a' alignment, size and style from 'src' alignment parameter
func (a *Alignment) Copy(src *Alignment) {
	a.Size = src.Size
	a.Style = 0
	a.Style |= src.Style
}

// IsBetter return true if 'a' is a better alignment than 'compare'
func (a Alignment) IsBetter(compare *Alignment) bool {
	if a.Size > compare.Size {
		return true
	} else if a.Size == compare.Size && a.Style <= compare.Style {
		return true
	}
	return false
}

func (a *Alignment) setStyleByMask(mask *[11]uint8, index int, player uint8) {
	var left, right bool

	valueLeft := (*mask)[(index+0)-1]
	valueRight := (*mask)[(index+4)+1]
	if valueLeft == rdef.Empty || valueLeft == player {
		left = true
	}
	if valueRight == rdef.Empty || valueRight == player {
		right = true
	}
	if left && right {
		a.Style |= rdef.AlignFree
	} else if left || right {
		a.Style |= rdef.AlignHalf
	} else {
		a.Style |= rdef.AlignFlanked
	}
}

// New create alignment after scan it
func New(mask *[11]uint8, player uint8) *Alignment {
	var ret, tmp Alignment

	for i := 1; i <= 5; i++ {
		tmp.Size = 0
		tmp.Style = 0
		available := uint8(0)
		for c := 0; c < 5; c++ {
			value := (*mask)[i+c]
			if value == player {
				tmp.Size++
			} else if value == rdef.Empty {
				available++
			} else {
				break
			}
		}
		if available+tmp.Size == 5 && tmp.Size >= ret.Size {
			tmp.setStyleByMask(mask, i, player)
			if !ret.IsBetter(&tmp) {
				ret.Copy(&tmp)
			}
		}
	}

	return &ret
}

// AnalyzeThree check the alignment to define a free-three and record it on isThree attribute
func (a *Alignment) AnalyzeThree(mask *[11]uint8) bool {
	player := mask[5]
	left := [4]uint8{mask[5-1], mask[5-2], mask[5-3], mask[5-4]}
	right := [4]uint8{mask[5+1], mask[5+2], mask[5+3], mask[5+4]}

	switch {
	// [0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0]
	case left[0] == player && left[1] == player && left[2] == rdef.Empty &&
		right[0] == rdef.Empty:
		fallthrough
	// [0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0]
	case left[0] == player && left[1] == rdef.Empty && right[0] == player &&
		right[1] == rdef.Empty:
		fallthrough
	// [0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0]
	case left[0] == rdef.Empty && right[0] == player && right[1] == player &&
		right[2] == rdef.Empty:
		fallthrough
	// [0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0]
	case (left[3] == rdef.Empty || left[3] == player) && left[2] == player &&
		left[1] == player && left[0] == rdef.Empty &&
		(right[0] == rdef.Empty || right[0] == player):
		fallthrough
	// [0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0]
	case (left[1] == rdef.Empty || left[1] == player) && left[0] == player &&
		right[0] == rdef.Empty && right[1] == player &&
		(right[1] == rdef.Empty || right[1] == player):
		fallthrough
	// [0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0]
	case (left[0] == rdef.Empty || left[0] == player) && right[0] == player &&
		right[1] == rdef.Empty && right[2] == player &&
		(right[3] == rdef.Empty || right[3] == player):
		fallthrough
	// [0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0]
	case (left[0] == rdef.Empty || left[0] == player) && right[0] == rdef.Empty &&
		right[1] == player && right[2] == player &&
		(right[3] == rdef.Empty || right[3] == player):
		fallthrough
	// [0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0]
	case (left[2] == rdef.Empty || left[2] == player) && left[1] == player &&
		left[0] == rdef.Empty && right[0] == player &&
		(right[1] == rdef.Empty || right[1] == player):
		fallthrough
	// [0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0
	case (left[3] == rdef.Empty || left[3] == player) && left[2] == player &&
		left[1] == rdef.Empty && left[0] == player &&
		(right[0] == rdef.Empty || right[0] == player):
		a.IsThree = true
		return true
	default:
		//default
	}

	return false
}

// IsConsecutive : _
func (a Alignment) IsConsecutive(mask *[11]uint8) bool {
	var size uint8
	player := mask[5]

	for i := 4; i >= 0; i-- {
		if (*mask)[i] == player {
			size++
		} else {
			break
		}
	}

	for i := 5; i <= 10; i++ {
		if (*mask)[i] == player {
			size++
		} else {
			break
		}
	}
	if size >= a.Size {
		return true
	}
	return false
}
