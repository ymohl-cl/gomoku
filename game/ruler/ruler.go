package ruler

// outOfBoard is used on mask of raw (see CheckRules)
const (
	empty      = 0
	Player1    = uint8(1)
	Player2    = uint8(2)
	outOfBoard = 3

	sizeY = 19
	sizeX = 19

	alignFree    = 1 << 0
	alignHalf    = 1 << 1
	alignFlanked = 2 << 2

	noNeighbourMessage     = "this spot don't have a neighbour token"
	outOfBoardMessage      = "this spot is out of board"
	doubleThreeMessage     = "this spot compose one double three"
	spotAlreadyUsedMessage = "this spot is already used by other token"
	winByCaptureMessage    = "congratulation, winner by capture"
	winByAlignmentMessage  = "congratulation, winner by alignment"
)

// Spot on the board define by Y and X
type Spot struct {
	Y int8
	X int8
}

// InfoWin keep info of raw when a align size == 5.
// Need to check if the win align is capturable or not
type InfoWin struct {
	mask [11]uint8
	dirY int8
	dirX int8
}

// Align describe the alignments to five available spot and size align > 2
type Align struct {
	size  uint8
	style uint8
	iWin  *InfoWin
}

// Rules struct contain all Rules checks informations
type Rules struct {
	/* Private attribute */
	player   uint8
	y        int8
	x        int8
	captures []*Spot
	aligns   []Align

	/* Public attribute */
	Info          string
	NumberThree   uint8
	NumberCapture uint8
	Movable       bool
	Win           bool
}

// GetOtherPlayer provide the value of player's opponent provided on parameter
func GetOtherPlayer(p uint8) uint8 {
	if p == Player1 {
		return Player2
	}
	return Player1
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
	if (*b)[uint8(y)][uint8(x)] == empty {
		return true
	}
	return false
}

// positionIsCapturable : Check if spot defined by posY and posX is capturable
func positionIsCapturable(b *[19][19]uint8, posY, posX int8, player uint8) bool {
	opponent := GetOtherPlayer(player)

	for dirY := int8(-1); dirY <= 1; dirY++ {
		for dirX := int8(-1); dirX <= 1; dirX++ {
			if dirY == 0 && dirX == 0 {
				continue
			}

			if !isOnTheBoard(posY+dirY, posX+dirX) || !isOnTheBoard(posY+dirY*2, posX+dirX*2) || !isOnTheBoard(posY+dirY*-1, posX+dirX*-1) {
				continue
			}

			if (*b)[posY+dirY][posX+dirX] == player {
				if ((*b)[posY+dirY*2][posX+dirX*2] == opponent && (*b)[posY+dirY*-1][posX+dirX*-1] == empty) ||
					((*b)[posY+dirY*2][posX+dirX*2] == empty && (*b)[posY+dirY*-1][posX+dirX*-1] == opponent) {
					return true
				}
			}
		}
	}
	return false
}

// newInfoWin : _
func (a *Align) newInfoWin(mask *[11]uint8, dirY, dirX int8) {
	a.iWin = new(InfoWin)
	copy(a.iWin.mask[:], (*mask)[:])
	a.iWin.dirY = dirY
	a.iWin.dirX = dirX
}

// alignIsCapturable : browse all point of alignment and check if the spot is capturable
func (a *Align) isCapturable(b *[19][19]uint8, player uint8, posY, posX int8) bool {
	i := a.iWin
	if i == nil {
		return false
	}

	for index := 5; index >= 0 && i.mask[index] == player; index-- {
		y := posY + i.dirY*int8(5-index)
		x := posX + i.dirX*int8(5-index)
		if positionIsCapturable(b, y, x, player) {
			return true
		}
	}

	for index := 6; index < 11 && i.mask[index] == player; index++ {
		y := posY + i.dirY*int8(5-index)
		x := posX + i.dirX*int8(5-index)
		if positionIsCapturable(b, y, x, player) {
			return true
		}
	}

	return false
}

// New return a new instance with one player and one position (y, x)
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

// ApplyMove write on the board the spot of move and delete captured spot
func (r *Rules) ApplyMove(b *[19][19]uint8) {
	(*b)[r.y][r.x] = r.player
	for _, capture := range r.captures {
		(*b)[capture.Y][capture.X] = empty
	}
}

// RestoreMove delete on the board the sport of move and restore deleted captured spot
func (r *Rules) RestoreMove(b *[19][19]uint8) {
	opponent := GetOtherPlayer(r.player)

	(*b)[r.y][r.x] = empty
	for _, capture := range r.captures {
		(*b)[capture.Y][capture.X] = opponent
	}
}

// GetCaptures : return the slice of captured points
func (r *Rules) GetCaptures() []*Spot {
	return r.captures
}

// GetPlayer provide index of player
func (r Rules) GetPlayer() uint8 {
	return r.player
}

func (r Rules) GetPosition() (int8, int8) {
	return r.y, r.x
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

// addCapture : _
func (r *Rules) addCapture(y, x int8) {
	r.captures = append(r.captures, &Spot{Y: y, X: x})
}

// isAvailablePosition : check if the position is available
func (r *Rules) isAvailablePosition(b *[19][19]uint8) bool {

	if !isOnTheBoard(r.y, r.x) {
		r.Info = outOfBoardMessage
		return false
	}
	if !isEmpty(b, r.y, r.x) {
		r.Info = spotAlreadyUsedMessage
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
				r.Movable = true
				return true
			}
		}
	}

	r.Info = noNeighbourMessage
	return false
}

// analyzeCapture : Check and records captured spot
func (r *Rules) analyzeCapture(mask *[11]uint8, dirY, dirX int8) {
	cible := GetOtherPlayer(r.player)

	if (*mask)[4] == cible && (*mask)[3] == cible && (*mask)[2] == r.player {
		r.NumberCapture++
		r.addCapture(r.y+dirY, r.x+dirX)
		r.addCapture(r.y+(dirY*2), r.x+(dirX*2))
	}

	if (*mask)[6] == cible && (*mask)[7] == cible && (*mask)[8] == r.player {
		r.NumberCapture++
		r.addCapture(r.y+(dirY*-1), r.x+(dirX*-1))
		r.addCapture(r.y+(dirY*-2), r.x+(dirX*-2))
	}
}

// analyzeAlign : Check alignment type and records when there are enought available spot
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
		} else if (*mask)[i] == empty {
			availablePosition++
		} else {
			break
		}
	}
	if (*mask)[lastIndex] == empty ||
		((*mask)[lastIndex] != r.player && (*mask)[lastIndex+1] == empty) {
		left = true
	}

	for i := 6; i <= 10; i++ {
		lastIndex = i
		if (*mask)[i] == r.player {
			a.size++
			availablePosition++
		} else if (*mask)[i] == empty {
			availablePosition++
		} else {
			break
		}
	}
	if (*mask)[lastIndex] == empty ||
		((*mask)[lastIndex] != r.player && (*mask)[lastIndex-1] == empty) {
		right = true
	}

	if availablePosition < 5 || a.size == 1 {
		return
	} else if availablePosition == 5 {
		a.style |= alignFlanked
	} else if left == true && right == true {
		a.style |= alignFree
		if a.size == 3 {
			r.NumberThree++
		}
	} else {
		a.style |= alignHalf
	}

	if a.size == 5 {
		a.newInfoWin(mask, dirY, dirX)
	}
	r.aligns = append(r.aligns, a)
}

// analyzeMoveCondition : Check conditions to accept the move
func (r *Rules) analyzeMoveCondition() bool {
	if r.NumberThree >= 2 && r.NumberCapture == 0 {
		r.Movable = false
		r.Info = doubleThreeMessage
		return false
	}

	return true
}

// analyzeWinCondition : Check conditions to accept the win move
func (r *Rules) analyzeWinCondition(b *[19][19]uint8, nbCaptured uint8) {
	totalCaps := r.NumberCapture + nbCaptured

	if totalCaps >= 5 {
		r.Win = true
		r.Info = winByCaptureMessage
	}

	for _, align := range r.aligns {
		if align.size >= 5 {
			r.ApplyMove(b)
			ret := align.isCapturable(b, r.player, r.y, r.x)
			r.RestoreMove(b)

			if ret == false {
				r.Win = true
				r.Info = winByAlignmentMessage
			}
		}
	}
}

// getMaskFromBoard : provide one slice of 11 spots. index 5 is a spot of move
// it's more simplest and fast to applies all rules.
// example [3 0 1 1 0 2 2 0 0 0 0]
//            < < < =| |= > > >
func (r *Rules) getMaskFromBoard(b *[19][19]uint8, dirY, dirX int8, mask *[11]uint8) {
	(*mask)[5] = r.player

	for i := int8(1); i <= 5; i++ {
		leftY := uint8(r.y + (dirY * i))
		leftX := uint8(r.x + (dirX * i))
		rightY := uint8(r.y + (dirY * -i))
		rightX := uint8(r.x + (dirX * -i))

		if isOnTheBoard(int8(leftY), int8(leftX)) {
			(*mask)[5-i] = (*b)[leftY][leftX]
		} else {
			(*mask)[5-i] = outOfBoard
		}

		if isOnTheBoard(int8(rightY), int8(rightX)) {
			(*mask)[5+i] = (*b)[rightY][rightX]
		} else {
			(*mask)[5+i] = outOfBoard
		}
	}
}

// CheckRules : check all rules (captures/doubleThree/win conditions ..)
// create 4 masks to check all direction
func (r *Rules) CheckRules(board *[19][19]uint8, nbCaps uint8) {
	var stop bool
	var mask [11]uint8

	if !r.isAvailablePosition(board) {
		return
	}

	// check around move position. y and x represent one direction
	for y := int8(-1); y <= 0; y++ {
		for x := int8(-1); x <= 1; x++ {
			if x == 0 && y == 0 {
				// all direction are checked so break
				stop = true
				break
			}

			// create mask to the direction
			r.getMaskFromBoard(board, y, x, &mask)

			// record the capturable spots
			r.analyzeCapture(&mask, y, x)
			// record informations alinment
			r.analyzeAlign(&mask, board, y, x)
		}

		if stop == true {
			break
		}
	}

	// check condition to move
	if !r.analyzeMoveCondition() {
		return
	}

	// check if win exist
	r.analyzeWinCondition(board, nbCaps)
	return
}
