package ruler

//Token values and size of board
const (
	TokenEmpty         = 0
	TokenP1            = 1
	TokenP2            = 2
	sizeY              = 19
	sizeX              = 19
	threeHorizontal    = 1 << 0
	threeVertical      = 1 << 1
	threeDiagonalLeft  = 1 << 2
	threeDiagonalRight = 1 << 3
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

// Rules struct contain all Rules checks informations
type Rules struct {
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
	IsWin         bool
	CapturableWin []*Capture
	PositionWin   []*Capture
	MessageWin    string

	// Winneable
	NbToken uint8
	NbLine  uint8
}

// New : Return a new instance and set default values of Rules struct
func New() *Rules {
	return &Rules{}
}

// isOnTheBoard : check if the (pX, pY) point are inside the board
func isOnTheBoard(pX, pY int8) bool {
	if pX < 0 || pY < 0 || pX >= 19 || pY >= 19 {
		return false
	}
	return true
}

// isEmpty : check if the (pX, pY) point are empty position
func isEmpty(b *[][]uint8, pX, pY int8) bool {
	if (*b)[uint8(pY)][uint8(pX)] == TokenEmpty {
		return true
	}
	return false
}

// IsAvailablePosition : check if the position is available
func (r *Rules) isAvailablePosition(b *[][]uint8, pX, pY int8) bool {

	if !isOnTheBoard(pX, pY) {
		r.MovedStr = "Index out of board"
		return false
	}
	if !isEmpty(b, pX, pY) {
		r.MovedStr = "Index already used"
		return false
	}

	//Check around posX/posY
	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {

			x := pX + xi
			y := pY + yi

			if (xi == 0 && yi == 0) || !isOnTheBoard(x, y) {
				continue
			}

			if !isEmpty(b, x, y) {
				r.IsMoved = true
				return true
			}
		}
	}

	r.MovedStr = "Not neighborhood"
	return false
}

// CheckRules : check all rules (captures/doubleThree/win conditions ..), more details on following functions
func (r *Rules) CheckRules(board *[][]uint8, pX, pY int8, player uint8, nbCaps uint8) {
	var maskThree uint8

	if !r.isAvailablePosition(board, pX, pY) {
		return
	}

	//Check around posX/posY
	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {
			x := pX + xi
			y := pY + yi
			if (xi == 0 && yi == 0) || !isOnTheBoard(x, y) {
				continue
			}
			if !isEmpty(board, x, y) {
				r.CheckCapture(board, pX, pY, xi, yi, player)
			}
			r.CheckDoubleThree(board, pX, pY, xi, yi, player, &maskThree)
		}
	}

	if nbCaps+r.NbCaps >= 5 {
		r.IsWin = true
		r.MessageWin = "Win by capture: " + string(nbCaps+r.NbCaps)
		return
	}

	r.CheckWinner(board, pX, pY, player)

	if r.IsWin == false && r.IsCaptured == false && r.NbThree >= 2 {
		r.IsMoved = false
		r.MovedStr = "DoubleThree"
	}

	r.CheckWinneable(board, pX, pY, player)
	return
}

// GetCaptures : return the slice of captured points
func (r *Rules) GetCaptures() []*Capture {
	return r.caps
}

// GetNbrCaptures : return the numbers of captured points (len of slice)
func (r Rules) GetNbrCaptures() int32 {
	return int32(len(r.caps))
}

// CheckCapture : Check a capture case in a line
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token
func (r *Rules) CheckCapture(board *[][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) {

	x := posX + xi
	y := posY + yi

	//No check if same player
	if (*board)[y][x] == currentPlayer {
		return
	}
	//No check if part of a capture line is out of Board
	if !isOnTheBoard(posX+(2*xi), posY+(2*yi)) || isEmpty(board, posX+(2*xi), posY+(2*yi)) {
		return
	}
	if !isOnTheBoard(posX+(3*xi), posY+(3*yi)) || isEmpty(board, posX+(3*xi), posY+(3*yi)) {
		return
	}
	//Capture check
	if currentPlayer != (*board)[uint8(posY+(2*yi))][uint8(posX+(2*xi))] {
		if currentPlayer == (*board)[uint8(posY+(3*yi))][uint8(posX+(3*xi))] {
			r.IsCaptured = true
			r.caps = append(r.caps, NewCapture(posX+xi, posY+yi))
			r.caps = append(r.caps, NewCapture(posX+(xi*2), posY+(yi*2)))
			r.NbCaps++
			return
		}
	}
	return
}

// IsThree : Check a "three case" in a line and record it
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token
func (r *Rules) IsThree(board *[][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) bool {
	var nbMe uint8

	r.PositionWin = nil
	//Check line
	for i := int8(-1); i <= 2; i++ {
		x := posX + (xi * i)
		y := posY + (yi * i)
		if !isOnTheBoard(x, y) {
			return false
		}
		if (*board)[y][x] == currentPlayer {
			nbMe++
			r.PositionWin = append(r.PositionWin, &Capture{X: x, Y: y})
			if nbMe > 2 {
				return false
			}
		} else if !isEmpty(board, x, y) {
			return false
		}
	}

	//Check Extremites
	if nbMe == 2 {
		if !isOnTheBoard(posX+(-2*xi), posY+(-2*yi)) {
			if isEmpty(board, posX+(-1*xi), posY+(-1*yi)) {
				return false
			}
		} else if !isEmpty(board, posX+(-2*xi), posY+(-2*yi)) {
			return false
		}
		if !isOnTheBoard(posX+(3*xi), posY+(3*yi)) {
			if isEmpty(board, posX+(2*xi), posY+(2*yi)) {
				return false
			}
		} else if !isEmpty(board, posX+(3*xi), posY+(3*yi)) {
			return false
		}
		return true
	}
	return false
}

// CheckDoubleThree : Call IsThree twice to record a "three" case
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token
func (r *Rules) CheckDoubleThree(board *[][]uint8, posX, posY, xi, yi int8, currentPlayer uint8, maskThree *uint8) {
	var flag bool

	if yi == 1 && xi == 1 {
		if *maskThree&threeDiagonalLeft != 0 {
			return
		}
	} else if yi == 1 && xi == 0 {
		if *maskThree&threeVertical != 0 {
			return
		}
	} else if yi == 1 && xi == -1 {
		if *maskThree&threeDiagonalRight != 0 {
			return
		}
	} else if yi == 0 && xi == 1 {
		if *maskThree&threeHorizontal != 0 {
			return
		}
	}

	if r.IsThree(board, posX, posY, xi, yi, currentPlayer) {
		r.NbThree++
		flag = true
	} else if r.IsThree(board, posX+xi, posY+yi, xi, yi, currentPlayer) {
		r.NbThree++
		flag = true
	}

	if !flag {
		return
	}

	r.PositionWin = append(r.PositionWin, &Capture{X: posX, Y: posY})
	r.CheckCaptureAfterMove(board, posX, posY, currentPlayer)

	if yi == -1 && xi == -1 {
		*maskThree |= threeDiagonalLeft
	} else if yi == -1 && xi == 0 {
		*maskThree |= threeVertical
	} else if yi == -1 && xi == 1 {
		*maskThree |= threeDiagonalRight
	} else if yi == 0 && xi == -1 {
		*maskThree |= threeHorizontal
	}
}

// isWin : Check if there is 5 current player token on board = win
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token
func (r *Rules) isWin(board *[][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) bool {
	var nbMe uint8

	r.PositionWin = nil
	for i := int8(-4); i <= 4; i++ {
		x := posX + (xi * i)
		y := posY + (yi * i)
		if !isOnTheBoard(x, y) {
			continue
		}
		if (x == posX && y == posY) || (*board)[y][x] == currentPlayer {
			nbMe++
			r.PositionWin = append(r.PositionWin, &Capture{X: x, Y: y})
			if nbMe == 5 {
				return true
			}
		} else {
			r.PositionWin = nil
			nbMe = 0
		}
	}

	return false
}

// CheckWinner : Call isWin for vertical, horizontal and diagonals lines around point
// board is actual board
// posX and poxY are where player want to play position
// currentPlayer is the actual player token
func (r *Rules) CheckWinner(board *[][]uint8, posX, posY int8, currentPlayer uint8) {

	if r.isWin(board, posX, posY, 1, -1, currentPlayer) {
		r.IsWin = true
	} else if r.isWin(board, posX, posY, 1, 0, currentPlayer) {
		r.IsWin = true
	} else if r.isWin(board, posX, posY, 1, 1, currentPlayer) {
		r.IsWin = true
	} else if r.isWin(board, posX, posY, 0, 1, currentPlayer) {
		r.IsWin = true
	}

	if r.IsWin {
		r.MessageWin = "Win by align five token"
		r.CapturableWin = nil
		r.CheckCaptureAfterMove(board, posX, posY, currentPlayer)
	}
}

func (r *Rules) CheckWinneable(board *[][]uint8, posX, posY int8, currentPlayer uint8) {
	if nb := r.isWinneable(board, posX, posY, 1, -1, currentPlayer); nb >= r.NbToken {
		if nb > r.NbToken {
			r.NbToken = nb
			r.NbLine = 1
		} else {
			r.NbLine++
		}
	}

	if nb := r.isWinneable(board, posX, posY, 1, 0, currentPlayer); nb >= r.NbToken {
		if nb > r.NbToken {
			r.NbToken = nb
			r.NbLine = 1
		} else {
			r.NbLine++
		}
	}

	if nb := r.isWinneable(board, posX, posY, 1, 1, currentPlayer); nb >= r.NbToken {
		if nb > r.NbToken {
			r.NbToken = nb
			r.NbLine = 1
		} else {
			r.NbLine++
		}
	}

	if nb := r.isWinneable(board, posX, posY, 0, 1, currentPlayer); nb >= r.NbToken {
		if nb > r.NbToken {
			r.NbToken = nb
			r.NbLine = 1
		} else {
			r.NbLine++
		}
	}
}

func (r *Rules) isWinneable(board *[][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) uint8 {
	var nbToken uint8
	var nbAvailable uint8

	for i := int8(-4); i <= 4; i++ {
		x := posX + (xi * i)
		y := posY + (yi * i)
		if !isOnTheBoard(x, y) {
			continue
		}
		if (x == posX && y == posY) || (*board)[y][x] == currentPlayer {
			nbToken++
			nbAvailable++
		} else if (*board)[y][x] == TokenEmpty {
			nbAvailable++
		} else {
			nbAvailable = 0
			nbToken = 0
		}

		if nbAvailable >= 5 {
			return nbToken
		}
	}

	return 0
}

func (r *Rules) CheckCaptureAfterMove(board *[][]uint8, posX, posY int8, currentPlayer uint8) {
	(*board)[uint8(posY)][uint8(posX)] = currentPlayer
	player := uint8(TokenP1)
	if currentPlayer == TokenP1 {
		player = TokenP2
	}
	for _, cap := range r.PositionWin {
		for yi := int8(-1); yi <= 1; yi++ {
			for xi := int8(-1); xi <= 1; xi++ {
				x := cap.X + xi
				y := cap.Y + yi
				if (xi == 0 && yi == 0) || !isOnTheBoard(x, y) {
					continue
				}
				//Can we put in this posX/posY
				if (*board)[y][x] == currentPlayer {
					rFake := New()
					rFake.CheckCapture(board, cap.X+2*xi, cap.Y+2*yi, -xi, -yi, player)
					if rFake.IsCaptured == true {
						r.CapturableWin = append(r.CapturableWin, &Capture{X: x, Y: y})
					}
				}
			}
		}
	}
	(*board)[uint8(posY)][uint8(posX)] = TokenEmpty
}
