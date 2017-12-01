package ruler

import (
	"github.com/ymohl-cl/gomoku/game/ruler/alignment"
	// rdef is the ruler defines
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

// outOfBoard is used on mask of raw (see CheckRules)
const (
	noNeighbourMessage     = "this spot don't have a neighbour token"
	outOfBoardMessage      = "this spot is out of board"
	doubleThreeMessage     = "this spot compose one double three"
	spotAlreadyUsedMessage = "this spot is already used by other token"
	winByCaptureMessage    = "congratulation, winner by capture"
	winByAlignmentMessage  = "congratulation, winner by alignment"
)

var mask [11]uint8

// Spot on the board define by Y and X
type Spot struct {
	Y int8
	X int8
}

// Rules struct contain all Rules checks informations
type Rules struct {
	/* Private attribute */
	player   uint8
	y        int8
	x        int8
	captures []*Spot
	aligns   []*alignment.Alignment

	/* Public attribute */
	Info          string
	NumberThree   uint8
	NumberCapture uint8
	Movable       bool
	Win           bool
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
		(*b)[capture.Y][capture.X] = rdef.Empty
	}
}

// RestoreMove delete on the board the sport of move and restore deleted captured spot
func (r *Rules) RestoreMove(b *[19][19]uint8) {
	opponent := rdef.GetOtherPlayer(r.player)

	(*b)[r.y][r.x] = rdef.Empty
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

// GetPosition : _
func (r Rules) GetPosition() (int8, int8) {
	return r.y, r.x
}

// GetNumberAlignment return the number of alignement
func (r Rules) GetNumberAlignment() int16 {
	return int16(len(r.aligns))
}

// GetBetterAlignment return the alignment with the most value to the evaluation (size and style)
func (r Rules) GetBetterAlignment() *alignment.Alignment {
	var save *alignment.Alignment

	for _, a := range r.aligns {
		if a.IsBetter(save) {
			save = a
		}
	}
	return save
}

// IsMyPosition check the position on the board is equal at player on the rule
func (r Rules) IsMyPosition(b *[19][19]uint8) bool {
	if (*b)[r.y][r.x] == r.player {
		return true
	}
	return false
}

// IsAvailablePositionLight : check if the position is available
func IsAvailablePositionLight(b *[19][19]uint8, vy, vx int8) bool {

	if !rdef.IsOnTheBoard(vy, vx) {
		return false
	}
	if !rdef.IsEmpty(b, vy, vx) {
		return false
	}

	//Check around posX/posY
	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {

			x := vx + xi
			y := vy + yi

			if (xi == 0 && yi == 0) || !rdef.IsOnTheBoard(y, x) {
				continue
			}

			if !rdef.IsEmpty(b, y, x) {
				return true
			}
		}
	}

	return false
}

// isAvailablePosition : check if the position is available
func (r *Rules) isAvailablePosition(b *[19][19]uint8) bool {

	if !rdef.IsOnTheBoard(r.y, r.x) {
		r.Info = outOfBoardMessage
		return false
	}
	if !rdef.IsEmpty(b, r.y, r.x) {
		r.Info = spotAlreadyUsedMessage
		return false
	}

	//Check around posX/posY
	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {

			x := r.x + xi
			y := r.y + yi

			if (xi == 0 && yi == 0) || !rdef.IsOnTheBoard(y, x) {
				continue
			}

			if !rdef.IsEmpty(b, y, x) {
				r.Movable = true
				return true
			}
		}
	}

	r.Info = noNeighbourMessage
	return false
}

// analyzeWinCondition : Check conditions to accept the win move
func (r *Rules) analyzeWinCondition(b *[19][19]uint8, nbCaptured uint8) {
	totalCaps := r.NumberCapture + nbCaptured

	if totalCaps >= 5 {
		r.Win = true
		r.Info = winByCaptureMessage
		return
	}

	for _, align := range r.aligns {
		if i := align.GetInfosWin(); i != nil {
			if !i.IsCapturable() {
				r.Win = true
				r.Info = winByAlignmentMessage
			}
			// no win because capture is possible
			// r.ApplyMove(b)
			// bool := align.isCapturable(b, r.player, r.y, r.x)
			// r.RestoreMove(b)
		}
	}
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

// analyzeAlignments : Check three alignment and win alignment
func (r *Rules) analyzeAlignments(dirY, dirX int8) {
	if alignment.AnalyzeThree(&mask) {
		r.NumberThree++
	}

	if a := alignment.AnalyzeWin(&mask); a != nil {
		a.NewInfosWin(&mask, dirY, dirX, false)
		r.aligns = append(r.aligns, a)
	}

	return
}

// addCapture : _
func (r *Rules) addCapture(y, x int8) {
	r.captures = append(r.captures, &Spot{Y: y, X: x})
}

// analyzeCapture : Check and records captured spot
func (r *Rules) analyzeCapture(dirY, dirX int8) {
	cible := rdef.GetOtherPlayer(r.player)

	if mask[4] == cible && mask[3] == cible && mask[2] == r.player {
		// simulate capture to the check alignment
		mask[4] = rdef.Empty
		mask[3] = rdef.Empty

		r.NumberCapture++
		r.addCapture(r.y+dirY, r.x+dirX)
		r.addCapture(r.y+(dirY*2), r.x+(dirX*2))
	}

	if mask[6] == cible && mask[7] == cible && mask[8] == r.player {
		// simulate capture to the check alignment
		mask[6] = rdef.Empty
		mask[7] = rdef.Empty

		r.NumberCapture++
		r.addCapture(r.y+(dirY*-1), r.x+(dirX*-1))
		r.addCapture(r.y+(dirY*-2), r.x+(dirX*-2))
	}
}

// getMaskFromBoard : provide one slice of 11 spots. index 5 is a spot of move
// it's more simplest and fast to applies all rules.
// example [3 0 1 1 0 2 2 0 0 0 0]
//            < < < =| |= > > >
func (r *Rules) getMaskFromBoard(b *[19][19]uint8, dirY, dirX int8) {
	mask[5] = r.player

	for i := int8(1); i <= 5; i++ {
		leftY := uint8(r.y + (dirY * i))
		leftX := uint8(r.x + (dirX * i))
		rightY := uint8(r.y + (dirY * -i))
		rightX := uint8(r.x + (dirX * -i))

		if rdef.IsOnTheBoard(int8(leftY), int8(leftX)) {
			mask[5-i] = (*b)[leftY][leftX]
		} else {
			mask[5-i] = rdef.OutOfBoard
		}

		if rdef.IsOnTheBoard(int8(rightY), int8(rightX)) {
			mask[5+i] = (*b)[rightY][rightX]
		} else {
			mask[5+i] = rdef.OutOfBoard
		}
	}
}

// CheckRules : check all rules (captures/doubleThree/win conditions ..)
// create 4 masks to check all direction
func (r *Rules) CheckRules(board *[19][19]uint8, nbCaps uint8) {
	var stop bool
	//	var mask [11]uint8

	//	if !r.isAvailablePosition(board) {
	//		return
	//	}
	r.Movable = true
	// check around move position. y and x represent one direction
	for y := int8(-1); y <= 0; y++ {
		for x := int8(-1); x <= 1; x++ {
			if x == 0 && y == 0 {
				// all direction are checked so break
				stop = true
				break
			}

			// create mask to the direction
			r.getMaskFromBoard(board, y, x)

			// record the capturable spots
			r.analyzeCapture(y, x)
			// record informations alignment
			r.analyzeAlignments(y, x)
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

// getAlignment : Check alignment type and records when there are enought available spot
func (r *Rules) getAlignment(dirY, dirX int8) {
	var a *alignment.Alignment

	a = alignment.New(&mask, r.player)
	if a.Size < 2 {
		return
	}

	if a.IsThree {
		r.NumberThree++
	}

	r.aligns = append(r.aligns, a)
	return
}

// UpdateAlignments define aligns with a new state baord
func (r *Rules) UpdateAlignments(board *[19][19]uint8) {
	r.aligns = nil
	r.NumberThree = 0
	// check around move position. y and x represent one direction
	for y := int8(-1); y <= 0; y++ {
		for x := int8(-1); x <= 1; x++ {
			if x == 0 && y == 0 {
				// all direction are checked so break
				return
			}

			// create mask to the direction
			r.getMaskFromBoard(board, y, x)

			// record informations alinment
			r.getAlignment(y, x)
		}
	}
}

/*
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
				if ((*b)[posY+dirY*2][posX+dirX*2] == opponent && (*b)[posY+dirY*-1][posX+dirX*-1] == Empty) ||
					((*b)[posY+dirY*2][posX+dirX*2] == Empty && (*b)[posY+dirY*-1][posX+dirX*-1] == opponent) {
					return true
				}
			}
		}
	}
	return false
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

*/
