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
	IsWin      bool
	MessageWin string
}

// New : Return a new instance and set default values of Rules struct
func New() *Rules {
	return &Rules{IsMoved: false, MovedStr: "", IsCaptured: false, NbThree: 0, IsWin: false}
}

func (r Rules) NewCopy() *Rules {
	return &Rules{
		IsMoved:    r.IsMoved,
		MovedStr:   r.MovedStr,
		IsCaptured: r.IsCaptured,
		caps:       r.caps,
		NbThree:    r.NbThree,
		IsWin:      r.IsWin,
	}
}

func (r *Rules) Copy(ref *Rules) {
	r.IsMoved = ref.IsMoved
	r.MovedStr = ref.MovedStr
	r.IsCaptured = ref.IsCaptured
	r.caps = ref.caps
	r.NbCaps = ref.NbCaps
	r.NbThree = ref.NbThree
	r.IsWin = ref.IsWin
}

func (r *Rules) Clean() {
	r.IsMoved = false
	r.MovedStr = ""
	r.IsCaptured = false
	r.NbCaps = 0
	r.caps = nil
	r.NbThree = 0
	r.IsWin = false
}

// checkPosition : check if the (posX, posY) point are inside the board
// posX and poxY are where player want to play position
func checkOnTheBoard(posX, posY int8) bool {
	if posX < 0 || posY < 0 || posX >= 19 || posY >= 19 {
		return false
	}
	return true
}

// checkPosition : check if the (posX, posY) point are empty
// board is actual board
// posX and poxY are where player want to play position
func checkPosition(board *[][]uint8, posX, posY int8) bool {
	if (*board)[uint8(posY)][uint8(posX)] == TokenEmpty {
		return true
	}
	return false
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
	if !checkOnTheBoard(posX+(2*xi), posY+(2*yi)) || checkPosition(board, posX+(2*xi), posY+(2*yi)) {
		return
	}
	if !checkOnTheBoard(posX+(3*xi), posY+(3*yi)) || checkPosition(board, posX+(3*xi), posY+(3*yi)) {
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
func (r Rules) IsThree(board *[][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) bool {
	var nbMe uint8

	//Check line
	for i := int8(-1); i <= 2; i++ {
		x := posX + (xi * i)
		y := posY + (yi * i)
		if !checkOnTheBoard(x, y) {
			return false
		}
		if (*board)[y][x] == currentPlayer {
			nbMe++
			if nbMe > 2 {
				return false
			}
		} else if !checkPosition(board, x, y) {
			return false
		}
	}

	//Check Extremites
	if nbMe == 2 {
		if !checkOnTheBoard(posX+(-2*xi), posY+(-2*yi)) {
			if checkPosition(board, posX+(-1*xi), posY+(-1*yi)) {
				return false
			}
		} else if !checkPosition(board, posX+(-2*xi), posY+(-2*yi)) {
			return false
		}
		if !checkOnTheBoard(posX+(3*xi), posY+(3*yi)) {
			if checkPosition(board, posX+(2*xi), posY+(2*yi)) {
				return false
			}
		} else if !checkPosition(board, posX+(3*xi), posY+(3*yi)) {
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
func (r Rules) isWin(board *[][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) bool {
	var nbMe uint8

	for i := int8(-4); i <= 4; i++ {
		x := posX + (xi * i)
		y := posY + (yi * i)
		if !checkOnTheBoard(x, y) {
			continue
		}
		if (x == posX && y == posY) || (*board)[y][x] == currentPlayer {
			nbMe++
			if nbMe == 5 {
				return true
			}
		} else {
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
	}
}

/*function negamax(node, depth, α, β, color)
02     if depth = 0 or node is a terminal node
03         return color * the heuristic value of node

04     childNodes := GenerateMoves(node)
05     childNodes := OrderMoves(childNodes)
06     bestValue := −∞
07     foreach child in childNodes
08         v := −negamax(child, depth − 1, −β, −α, −color)
09         bestValue := max( bestValue, v )
10         α := max( α, v )
11         if α ≥ β
12             break
13     return bestValue
*/

// CheckRules : Execute all Check (moves/captures/doubleThree/Win)
// board is actual board
// posX and poxY are where player want to play position
// currentPlayer is the actual player token
func (r *Rules) CheckRules(board *[][]uint8, posX, posY int8, currentPlayer uint8, nbCaps uint8) {
	var maskThree uint8

	//Basic init posX/posY check
	if !checkOnTheBoard(posX, posY) {
		r.MovedStr = "Index not allowed"
		return
	}
	if !checkPosition(board, posX, posY) {
		r.MovedStr = "Already used"
		return
	}

	//Check around posX/posY
	//	wg := new(sync.WaitGroup)
	//	m := new(sync.Mutex)
	//	for i := int8(-1); i <= 1; i++ {
	//		for j := int8(-1); j <= 1; j++ {
	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {
			//			wg.Add(1)
			//			go func(xi, yi int8) {
			//				defer wg.Done()
			x := posX + xi
			y := posY + yi
			if (xi == 0 && yi == 0) || !checkOnTheBoard(x, y) {
				//				return
				continue
			}
			//Can we put in this posX/posY
			if !checkPosition(board, x, y) {
				//We can put here
				r.IsMoved = true
				r.CheckCapture(board, posX, posY, xi, yi, currentPlayer)
			}
			//				r.CheckDoubleThree(board, posX, posY, xi, yi, currentPlayer, &maskThree, m)
			r.CheckDoubleThree(board, posX, posY, xi, yi, currentPlayer, &maskThree)
			//			}(j, i)
		}
	}
	//	wg.Wait()
	// Check if this move is a winning
	if r.IsMoved == false {
		r.MovedStr = "Not neighborhood"
		return
	}

	if nbCaps+r.NbCaps >= 5 {
		r.IsWin = true
		r.MessageWin = "Win by capture: " + string(nbCaps+r.NbCaps)
		return
	}
	r.CheckWinner(board, posX, posY, currentPlayer)

	if r.IsWin == false && r.IsCaptured == false && r.NbThree >= 2 {
		r.IsMoved = false
		r.MovedStr = "DoubleThree"
	}
	return
}
