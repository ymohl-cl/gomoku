package ruler

import "fmt"

//Token values and size of board
const (
	TokenEmpty = 0
	TokenP1    = 1
	TokenP2    = 2
	outOfBoard = 3

	sizeY = 19
	sizeX = 19

	threeHorizontal    = 1 << 0
	threeVertical      = 1 << 1
	threeDiagonalLeft  = 1 << 2
	threeDiagonalRight = 1 << 3

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

// New : Return a new instance and set default values of Rules struct
func New(p uint8, py, px int8) *Rules {
	return &Rules{
		player: p,
		y:      py,
		x:      px,
	}
}

// Init : _
func (r *Rules) Init(p uint8, py, px int8) {
	r.player = p
	r.y = py
	r.x = px
}

// getOtherPlayer provide the value of opponent player provided on parameter
func getOtherPlayer(p uint8) uint8 {
	if p == TokenP1 {
		return TokenP2
	}
	return TokenP1
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

func (r *Rules) Print() {
	if r.player == TokenP1 {
		fmt.Println("Player: P1")
	} else if r.player == TokenP2 {
		fmt.Println("Player: P2")
	} else {
		fmt.Println("Player not define")
	}

	fmt.Println("Can move: ", r.IsMoved, " | details: ", r.MovedStr)

	fmt.Println("Captured: ", r.IsCaptured, " | Nomber: ", r.NbCaps)
	for _, cap := range r.caps {
		fmt.Println("\t capture: ", cap.Y, "-", cap.X)
	}

	fmt.Println("Number tree: ", r.NbThree)

	fmt.Println("IsWin: ", r.IsWin, " - details: ", r.MessageWin)

	fmt.Println("Number align: ", len(r.aligns))
	for _, align := range r.aligns {
		fmt.Print("align size: ", align.size, " - Capturable: ", align.capturable)
		if align.style&alignFree != 0 {
			fmt.Println(" | style free ")
		} else if align.style&alignHalf != 0 {
			fmt.Println(" | style halfFree ")
		} else if align.style&alignFlanked != 0 {
			fmt.Println(" | style flanked ")
		} else {
			fmt.Println(" | style not define")
		}
	}
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

func (r *Rules) getMaskFromBoard(b *[][]uint8, pX, pY, xi, yi int8, mask *[11]uint8) {
	tx := uint8(pX + (xi * 5))
	ty := uint8(pY + (yi * 5))
	otx := uint8(pX + (xi * -5))
	oty := uint8(pY + (yi * -5))

	(*mask)[5] = r.player

	for i := int8(0); i < 5; i++ {

		if isOnTheBoard(int8(tx), int8(ty)) {
			(*mask)[i] = (*b)[ty][tx]
		} else {
			(*mask)[i] = outOfBoard
		}

		if isOnTheBoard(int8(otx), int8(oty)) {
			(*mask)[10-i] = (*b)[oty][otx]
		} else {
			(*mask)[10-i] = outOfBoard
		}

		if yi != 0 {
			ty++
			oty--
		}
		if xi != 0 {
			if xi < 0 {
				tx++
				otx--

			} else {
				tx--
				otx++
			}
		}
	}
}

func (r *Rules) analyzeCapture(mask *[11]uint8, player uint8, pX, pY, xi, yi int8) {
	cible := getOtherPlayer(player)

	if (*mask)[4] == cible && (*mask)[3] == cible && (*mask)[2] == player {
		r.NbCaps++

		r.caps = append(r.caps, NewCapture(pX+xi, pY+yi))
		r.caps = append(r.caps, NewCapture(pX+(xi*2), pY+(yi*2)))
	}

	if (*mask)[6] == cible && (*mask)[7] == cible && (*mask)[8] == player {
		r.NbCaps++
		r.caps = append(r.caps, NewCapture(pX+(xi-1), pY+(yi*-1)))
		r.caps = append(r.caps, NewCapture(pX+(xi*-2), pY+(yi*-2)))
	}
}

func (r *Rules) positionIsCapturable(b *[][]uint8, pX, pY int8, player uint8) bool {
	opponent := getOtherPlayer(player)

	for y := int8(-1); y <= 1; y++ {
		for x := int8(-1); x <= 1; x++ {
			if x == 0 && y == 0 {
				continue
			}

			xi := pX + x
			yi := pY + y

			if !isOnTheBoard(xi, yi) || !isOnTheBoard(xi*2, yi*2) || !isOnTheBoard(xi*-1, yi*-1) {
				continue
			}

			if (*b)[yi][xi] == player {
				if ((*b)[yi*2][xi*2] == opponent && (*b)[yi*-1][xi*-1] == TokenEmpty) ||
					((*b)[yi*2][xi*2] == TokenEmpty && (*b)[yi*-1][xi*-1] == opponent) {
					return true
				}
			}
		}
	}
	return false
}

func (r *Rules) alignIsCapturable(b *[][]uint8, mask *[11]uint8, pX, pY, xi, yi int8, player uint8) bool {
	var index int

	for index = 5; index >= 0 && (*mask)[index] == player; index-- {
	}
	index += 1

	for i := 0; i < 5; i++ {
		tx := pX + xi*int8(index)
		ty := pY + yi*int8(index)
		if r.positionIsCapturable(b, tx, ty, player) {
			return true
		}
	}
	return false
}

func (r *Rules) analyzeAlign(mask *[11]uint8, player uint8, b *[][]uint8, pX, pY, xi, yi int8) {
	a := Align{size: 1}
	availablePosition := 1
	var lastIndex int
	var left, right bool

	for i := 4; i >= 0; i-- {
		lastIndex = i
		if (*mask)[i] == player {
			a.size++
			availablePosition++
		} else if (*mask)[i] == TokenEmpty {
			availablePosition++
		} else {
			break
		}
	}
	if (*mask)[lastIndex] == TokenEmpty ||
		((*mask)[lastIndex] != player && (*mask)[lastIndex+1] == TokenEmpty) {
		left = true
	}

	for i := 6; i <= 10; i++ {
		lastIndex = i
		if (*mask)[i] == player {
			a.size++
			availablePosition++
		} else if (*mask)[i] == TokenEmpty {
			availablePosition++
		} else {
			break
		}
	}
	if (*mask)[lastIndex] == TokenEmpty ||
		((*mask)[lastIndex] != player && (*mask)[lastIndex-1] == TokenEmpty) {
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

	if a.size == 5 && r.alignIsCapturable(b, mask, pX, pY, xi, yi, player) {
		a.capturable = true
	}
	r.aligns = append(r.aligns, a)
}

func (r *Rules) analyzeMoveCondition() bool {
	if r.NbThree >= 2 && r.NbCaps == 0 {
		fmt.Println("Double tree")
		return false
	}

	return true
}

func (r *Rules) analyzeWinCondition(nbCaps uint8) {
	totalCaps := r.NbCaps + nbCaps

	if totalCaps >= 5 {
		r.IsWin = true
		r.MessageWin = "Great, win by capture: " + string(totalCaps)
	}

	for _, align := range r.aligns {
		if align.size >= 5 && !align.capturable {
			r.IsWin = true
			r.MessageWin = "Great, win by align: " + string(align.size)
		}
	}
}

func (r Rules) PrintMask(m *[11]uint8) {
	for index, v := range *m {
		fmt.Print("[ ", index, ": ", v, " ] ")
	}
	fmt.Println(" || ")
}

// CheckRules : check all rules (captures/doubleThree/win conditions ..), more details on following functions
func (r *Rules) CheckRules(board *[][]uint8, pX, pY int8, player uint8, nbCaps uint8) {
	var stop bool
	var mask [11]uint8

	if !r.isAvailablePosition(board, pX, pY) {
		return
	}

	//	fmt.Println("START NEW RULE")

	//Check around posX/posY
	for yi := int8(-1); yi <= 0; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {
			if xi == 0 && yi == 0 {
				stop = true
				break
			}

			r.getMaskFromBoard(board, pX, pY, xi, yi, &mask)

			//			r.PrintMask(&mask)
			r.analyzeCapture(&mask, player, pX, pY, xi, yi)
			r.analyzeAlign(&mask, player, board, pX, pY, xi, yi)
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

// GetNbrCaptures : return the numbers of captured points (len of slice)
/*
func (r Rules) GetNbrCaptures() int32 {
	return int32(len(r.caps))
}
*/
// CheckCapture : Check a capture case in a line
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token
/*
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
*/
// IsThree : Check a "three case" in a line and record it
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token
/*
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
*/
// CheckDoubleThree : Call IsThree twice to record a "three" case
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token
/*
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
*/
// isWin : Check if there is 5 current player token on board = win
// board is actual board
// posX and poxY are where player want to play position
// xi and yi are steps of posX and posY respectively for the direction check
// currentPlayer is the actual player token

/*
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
*/
// CheckWinner : Call isWin for vertical, horizontal and diagonals lines around point
// board is actual board
// posX and poxY are where player want to play position
// currentPlayer is the actual player token
/*
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
*/
/*
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
*/
/*
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
*/
/*
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
*/
