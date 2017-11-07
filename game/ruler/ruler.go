package ruler

import "strconv"

//Token values and size of board
const (
	TokenEmpty = 0
	TokenP1    = 1
	TokenP2    = 2
	outOfBoard = 3

	sizeY = 19
	sizeX = 19

	alignFree    = 1 << 0
	alignHalf    = 1 << 1
	alignFlanked = 2 << 2
)

// Capture struct is a couple (X, Y) of a captured point
type Capture struct {
	X int8
	Y int8
}

// NewCapture : Return a new instance and set (X, Y) values of a captured point
func NewCapture(posX, posY int8) *Capture {
	return &Capture{X: posX, Y: posY}
}

type Align struct {
	size       uint8
	style      uint8
	capturable bool
}

// Rules struct contain all Rules checks informations
type Rules struct {
	player uint8
	y      int8
	x      int8

	//Move
	IsMoved  bool
	MovedStr string

	//Capture
	IsCaptured bool
	caps       []*Capture
	NbCaps     uint8

	//Double-Three
	NbThree uint8

	//Win
	IsWin bool
	//	CapturableWin []*Capture
	//	PositionWin   []*Capture
	MessageWin string

	// Winneable
	//	NbToken uint8
	//	NbLine  uint8
	aligns []Align
}

// New return a new instance and with player and one position (y, x)
func New(p uint8, y, x int8) *Rules {
	return &Rules{
		player: p,
		y:      y,
		x:      x,
	}
}

// Init a rules with a player and one position (y, x)
func (r *Rules) Init(p uint8, y, x int8) {
	r.player = p
	r.y = y
	r.x = x
}

// getOtherPlayer provide the value of player's opponent provided on parameter
func getOtherPlayer(p uint8) uint8 {
	if p == TokenP1 {
		return TokenP2
	}
	return TokenP1
}

// isOnTheBoard : check if the (y, x) point are inside the board
func isOnTheBoard(y, x int8) bool {
	if x < 0 || y < 0 || x >= 19 || y >= 19 {
		return false
	}
	return true
}

// isEmpty : check if the (pX, pY) point are empty position
func isEmpty(b *[19][19]uint8, y, x int8) bool {
	if (*b)[uint8(y)][uint8(x)] == TokenEmpty {
		return true
	}
	return false
}

// IsAvailablePosition : check if the position is available
func (r *Rules) isAvailablePosition(b *[19][19]uint8) bool {

	if !isOnTheBoard(r.y, r.x) {
		r.MovedStr = "Index out of board"
		return false
	}
	if !isEmpty(b, r.y, r.x) {
		r.MovedStr = "Index already used"
		return false
	}

	//Check around posX/posY
	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {

			x := r.x + xi
			y := r.y + yi

			if (xi == 0 && yi == 0) || !isOnTheBoard(y, x) {
				continue
			}

			if !isEmpty(b, y, x) {
				r.IsMoved = true
				return true
			}
		}
	}

	r.MovedStr = "Not neighborhood"
	return false
}

func (r *Rules) getMaskFromBoard(b *[19][19]uint8, dirY, dirX int8, mask *[11]uint8) {
	tx := uint8(r.x + (dirX * 5))
	ty := uint8(r.y + (dirY * 5))
	otx := uint8(r.x + (dirX * -5))
	oty := uint8(r.y + (dirY * -5))

	(*mask)[5] = r.player

	for i := int8(0); i < 5; i++ {

		if isOnTheBoard(int8(ty), int8(tx)) {
			(*mask)[i] = (*b)[ty][tx]
		} else {
			(*mask)[i] = outOfBoard
		}

		if isOnTheBoard(int8(oty), int8(otx)) {
			(*mask)[10-i] = (*b)[oty][otx]
		} else {
			(*mask)[10-i] = outOfBoard
		}

		if dirY != 0 {
			ty++
			oty--
		}
		if dirX != 0 {
			if dirX < 0 {
				tx++
				otx--

			} else {
				tx--
				otx++
			}
		}
	}
}

func (r *Rules) analyzeCapture(mask *[11]uint8, dirY, dirX int8) {
	cible := getOtherPlayer(r.player)

	if (*mask)[4] == cible && (*mask)[3] == cible && (*mask)[2] == r.player {
		r.NbCaps++
		r.caps = append(r.caps, NewCapture(r.x+dirX, r.y+dirY))
		r.caps = append(r.caps, NewCapture(r.x+(dirX*2), r.y+(dirY*2)))
	}

	if (*mask)[6] == cible && (*mask)[7] == cible && (*mask)[8] == r.player {
		r.NbCaps++
		r.caps = append(r.caps, NewCapture(r.x+(dirX*-1), r.y+(dirY*-1)))
		r.caps = append(r.caps, NewCapture(r.x+(dirX*-2), r.y+(dirY*-2)))
	}
}

func (r *Rules) positionIsCapturable(b *[19][19]uint8, posY, posX int8) bool {
	opponent := getOtherPlayer(r.player)

	for dirY := int8(-1); dirY <= 1; dirY++ {
		for dirX := int8(-1); dirX <= 1; dirX++ {
			if dirY == 0 && dirX == 0 {
				continue
			}

			if !isOnTheBoard(posY+dirY, posX+dirX) || !isOnTheBoard(posY+dirY*2, posX+dirX*2) || !isOnTheBoard(posY+dirY*-1, posX+dirX*-1) {
				continue
			}

			if (*b)[posY+dirY][posX+dirX] == r.player {
				if ((*b)[posY+dirY*2][posX+dirX*2] == opponent && (*b)[posY+dirY*-1][posX+dirX*-1] == TokenEmpty) ||
					((*b)[posY+dirY*2][posX+dirX*2] == TokenEmpty && (*b)[posY+dirY*-1][posX+dirX*-1] == opponent) {
					return true
				}
			}
		}
	}
	return false
}

func (r *Rules) alignIsCapturable(b *[19][19]uint8, mask *[11]uint8, dirY, dirX int8) bool {
	for index := 5; index >= 0 && (*mask)[index] == r.player; index-- {
		y := r.y + dirY*int8(5-index)
		x := r.x + dirX*int8(5-index)
		if r.positionIsCapturable(b, y, x) {
			return true
		}
	}

	for index := 6; index < 11 && (*mask)[index] == r.player; index++ {
		y := r.y + dirY*int8(5-index)
		x := r.x + dirX*int8(5-index)
		if r.positionIsCapturable(b, y, x) {
			return true
		}
	}

	return false
}

func (r *Rules) analyzeAlign(mask *[11]uint8, b *[19][19]uint8, dirY, dirX int8) {
	a := Align{size: 1}
	availablePosition := 1
	var lastIndex int
	var left, right bool

	for i := 4; i >= 0; i-- {
		lastIndex = i
		if (*mask)[i] == r.player {
			a.size++
			availablePosition++
		} else if (*mask)[i] == TokenEmpty {
			availablePosition++
		} else {
			break
		}
	}
	if (*mask)[lastIndex] == TokenEmpty ||
		((*mask)[lastIndex] != r.player && (*mask)[lastIndex+1] == TokenEmpty) {
		left = true
	}

	for i := 6; i <= 10; i++ {
		lastIndex = i
		if (*mask)[i] == r.player {
			a.size++
			availablePosition++
		} else if (*mask)[i] == TokenEmpty {
			availablePosition++
		} else {
			break
		}
	}
	if (*mask)[lastIndex] == TokenEmpty ||
		((*mask)[lastIndex] != r.player && (*mask)[lastIndex-1] == TokenEmpty) {
		right = true
	}

	if availablePosition < 5 || a.size == 1 {
		return
	} else if availablePosition == 5 {
		a.style |= alignFlanked
	} else if left == true && right == true {
		a.style |= alignFree
		if a.size == 3 {
			r.NbThree++
		}
	} else {
		a.style |= alignHalf
	}

	if a.size == 5 && r.alignIsCapturable(b, mask, dirY, dirX) {
		a.capturable = true
	}
	r.aligns = append(r.aligns, a)
}

func (r *Rules) analyzeMoveCondition() bool {
	if r.NbThree >= 2 && r.NbCaps == 0 {
		r.IsMoved = false
		r.MovedStr = "Double three"
		return false
	}

	return true
}

func (r *Rules) analyzeWinCondition(nbCaps uint8) {
	totalCaps := r.NbCaps + nbCaps

	if totalCaps >= 5 {
		r.IsWin = true
		r.MessageWin = "Great, win by capture: " + strconv.FormatUint(uint64(totalCaps), 10)
	}

	for _, align := range r.aligns {
		if align.size >= 5 && !align.capturable {
			r.IsWin = true
			r.MessageWin = "Great, win by align: " + strconv.FormatUint(uint64(align.size), 10)
		}
	}
}

// CheckRules : check all rules (captures/doubleThree/win conditions ..), more details on following functions
func (r *Rules) CheckRules(board *[19][19]uint8, pX, pY int8, player uint8, nbCaps uint8) {
	var stop bool
	var mask [11]uint8

	if !r.isAvailablePosition(board) {
		return
	}

	//Check around posX/posY
	for yi := int8(-1); yi <= 0; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {
			if xi == 0 && yi == 0 {
				stop = true
				break
			}

			r.getMaskFromBoard(board, yi, xi, &mask)

			r.analyzeCapture(&mask, yi, xi)
			r.analyzeAlign(&mask, board, yi, xi)
		}

		if stop == true {
			break
		}
	}

	if !r.analyzeMoveCondition() {
		return
	}

	r.analyzeWinCondition(nbCaps)
	return
}

// GetCaptures : return the slice of captured points
func (r *Rules) GetCaptures() []*Capture {
	return r.caps
}

// GetMaxAlign return the most alignment
func (r Rules) GetMaxAlign() uint8 {
	nbToken := uint8(0)

	for _, align := range r.aligns {
		if align.size > nbToken {
			nbToken = align.size
		}
	}
	return nbToken
}
