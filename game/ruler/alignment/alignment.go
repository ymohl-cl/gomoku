package alignment

import rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"

// Spot on the board define by Y and X
type Spot struct {
	Y int8
	X int8
}

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

// IsStyle check is the style parameters match with the style alignment
func (a *Alignment) IsStyle(style uint8) bool {
	if a.Style&style != 0 {
		return true
	}
	return false
}

// Copy on the 'a' alignment, size and style from 'src' alignment parameter
func (a *Alignment) Copy(src *Alignment) {
	a.Size = src.Size
	a.Style = src.Style
	a.IsThree = src.IsThree
	a.iWin = src.iWin
}

// clear the alignment
func (a *Alignment) clear() {
	a.Size = 0
	a.Style = 0
	a.IsThree = false
	a.iWin = nil
}

// IsBetter return true if 'a' is a better alignment than 'compare'
func (a Alignment) IsBetter(compare *Alignment) bool {
	if compare == nil {
		return true
	}

	if a.Size > compare.Size {
		return true
	} else if a.Size < compare.Size {
		return false
	}

	// alignments have a same size
	if a.IsThree {
		return true
	} else if compare.IsThree {
		return false
	} else if a.Style <= compare.Style {
		return true
	}
	return false
}

func (a *Alignment) setStyleByMask(mask *[11]uint8, index int, player uint8) {
	var left, right bool

	vLeft := (*mask)[(index + 0)]
	vRight := (*mask)[(index + 4)]
	vLeftOut := (*mask)[(index+0)-1]
	vRightOut := (*mask)[(index+4)+1]
	if vLeftOut != player && vLeftOut != rdef.Empty && vRightOut != player && vRightOut != rdef.Empty {
		left = false
		right = false
	} else {
		if vLeft == rdef.Empty || (vLeft == player && vLeftOut == rdef.Empty) {
			left = true
		}
		if vRight == rdef.Empty || (vRight == player && vRightOut == rdef.Empty) {
			right = true
		}
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
	var ret Alignment

	for i := 1; i <= 5; i++ {
		size := uint8(0)
		available := uint8(0)

		for c := 0; c < 5; c++ {
			value := mask[i+c]
			if value == player {
				size++
			} else if value == rdef.Empty {
				available++
			} else {
				break
			}
		}

		if available+size == 5 && size >= ret.Size {
			tmp := Alignment{Size: size}
			tmp.setStyleByMask(mask, i, player)
			if tmp.Size == 3 && tmp.Style&rdef.AlignFree != 0 {
				if AnalyzeThree(mask) {
					tmp.IsThree = true
				}
			}
			if !ret.IsBetter(&tmp) {
				ret.Copy(&tmp)
			}
		}
	}

	return &ret
}

// AnalyzeWin alignment and return it if it's a case
func AnalyzeWin(mask *[11]uint8) *Alignment {
	size := uint8(1)
	player := mask[5]

	for i := 4; i >= 0; i-- {
		if mask[i] == player {
			size++
		} else {
			break
		}
	}
	for i := 6; i < 11; i++ {
		if mask[i] == player {
			size++
		} else {
			break
		}
	}
	if size >= 5 {
		return &Alignment{Size: size}
	}
	return nil
}

// AnalyzeThree check the alignment to define a free-three and record it on isThree attribute
func AnalyzeThree(mask *[11]uint8) bool {
	player := mask[5]

	switch {
	// [0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0]
	case mask[5-1] == player && mask[5-2] == player && mask[5-3] == rdef.Empty &&
		mask[5+1] == rdef.Empty && (mask[5-4] == rdef.Empty || mask[5+2] == rdef.Empty):
		fallthrough
	// [0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0]
	case mask[5-1] == player && mask[5-2] == rdef.Empty && mask[5+1] == player &&
		mask[5+2] == rdef.Empty && (mask[5+3] == rdef.Empty || mask[5-3] == rdef.Empty):
		fallthrough
	// [0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0]
	case mask[5-1] == rdef.Empty && mask[5+1] == player && mask[5+2] == player &&
		mask[5+3] == rdef.Empty && (mask[5+4] == rdef.Empty || mask[5-2] == rdef.Empty):
		fallthrough
	// [0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0]
	case mask[5-4] == rdef.Empty && mask[5-3] == player && mask[5-2] == player &&
		mask[5-1] == rdef.Empty && mask[5+1] == rdef.Empty:
		fallthrough
	// [0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0]
	case mask[5-2] == rdef.Empty && mask[5-1] == player && mask[5+1] == rdef.Empty &&
		mask[5+2] == player && mask[5+3] == rdef.Empty:
		fallthrough
	// [0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0]
	case mask[5-1] == rdef.Empty && mask[5+1] == player && mask[5+2] == rdef.Empty &&
		mask[5+3] == player && mask[5+4] == rdef.Empty:
		fallthrough
	// [0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0]
	case mask[5-1] == rdef.Empty && mask[5+1] == rdef.Empty && mask[5+2] == player &&
		mask[5+3] == player && mask[5+4] == rdef.Empty:
		fallthrough
	// [0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0]
	case mask[5-3] == rdef.Empty && mask[5-2] == player && mask[5-1] == rdef.Empty &&
		mask[5+1] == player && mask[5+2] == rdef.Empty:
		fallthrough
	// [0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0
	case mask[5-4] == rdef.Empty && mask[5-3] == player && mask[5-2] == rdef.Empty &&
		mask[5-1] == player && mask[5+1] == rdef.Empty:
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

// positionIsCapturable : Check if spot defined by posY and posX is capturable
func positionIsCapturable(b *[19][19]uint8, posY, posX int8, player uint8) (bool, int8, int8) {
	opponent := rdef.GetOtherPlayer(player)

	for dirY := int8(-1); dirY <= 1; dirY++ {
		for dirX := int8(-1); dirX <= 1; dirX++ {
			if dirY == 0 && dirX == 0 {
				continue
			}

			if !rdef.IsOnTheBoard(posY+dirY, posX+dirX) || !rdef.IsOnTheBoard(posY+dirY*2, posX+dirX*2) || !rdef.IsOnTheBoard(posY+dirY*-1, posX+dirX*-1) {
				continue
			}

			if (*b)[posY+dirY][posX+dirX] == player {
				if ((*b)[posY+dirY*2][posX+dirX*2] == opponent && (*b)[posY+dirY*-1][posX+dirX*-1] == rdef.Empty) ||
					((*b)[posY+dirY*2][posX+dirX*2] == rdef.Empty && (*b)[posY+dirY*-1][posX+dirX*-1] == opponent) {
					if (*b)[posY+dirY*-1][posX+dirX*-1] == rdef.Empty {
						return true, posY + dirY*-1, posX + dirX*-1
					}
					return true, posY + dirY*2, posX + dirX*2
				}
			}
		}
	}
	return false, 0, 0
}

// alignIsCapturable : browse all point of alignment and check if the spot is capturable
func (a *Alignment) IsCapturable(b *[19][19]uint8, player uint8, posY, posX int8) []*Spot {
	var spots []*Spot
	i := a.iWin
	if i == nil {
		return nil
	}

	for index := 5; index >= 0 && i.mask[index] == player; index-- {
		y := posY + i.dirY*int8(5-index)
		x := posX + i.dirX*int8(5-index)
		if ok, py, px := positionIsCapturable(b, y, x, player); ok {
			i.capturable = true
			spots = append(spots, &Spot{py, px})
		}
	}

	for index := 6; index < 11 && i.mask[index] == player; index++ {
		y := posY + i.dirY*int8(5-index)
		x := posX + i.dirX*int8(5-index)
		if ok, py, px := positionIsCapturable(b, y, x, player); ok {
			i.capturable = true
			spots = append(spots, &Spot{py, px})
		}
	}
	return spots
}

func (a *Alignment) GetCaptureStatus() bool {
	if a.iWin == nil {
		return false
	}
	return a.iWin.capturable
}
