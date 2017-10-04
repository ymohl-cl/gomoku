package ruler

const (
	TokenEmpty = 0
	TokenP1    = 1
	TokenP2    = 2
	sizeY      = 19
	sizeX      = 19
)

// Rules is a implementation interface to the rule of gomoku
type Rules struct {
	caps []Capture
}

type Capture struct {
	X int8
	Y int8
}

func NewCapture(posX, posY int8) Capture {
	return Capture{X: posX, Y: posY}
}

func (r Rules) GetNbrCaptures() int32 {
	return int32(len(r.caps))
}

func New() *Rules {
	return &Rules{}
}

func checkOnTheBoard(posX, posY int8) bool {
	if posX < 0 || posY < 0 || posX >= 19 || posY >= 19 {
		return false
	}
	return true
}

func checkPosition(board [][]uint8, posX, posY int8) bool {
	if board[uint8(posY)][uint8(posX)] == TokenEmpty {
		return true
	}
	return false
}

func (r Rules) CheckMove(board [][]uint8, posX, posY int8) (bool, string) {
	if !checkOnTheBoard(posX, posY) {
		return false, "Index not allowed"
	}
	if !checkPosition(board, posX, posY) {
		return false, "Already used"
	}

	for y := posY - 1; y <= posY+1; y++ {
		for x := posX - 1; x <= posX+1; x++ {
			if (posX == x && posY == y) || !checkOnTheBoard(x, y) {
				continue
			}
			if !checkPosition(board, x, y) {
				return true, "Cool"
			}
		}
	}
	return false, "Not neighborhood"
}

func (r *Rules) SetCapture(board [][]uint8, posX, posY int8, xi, yi int8, currentPlayer uint8) bool {
	if !checkOnTheBoard(posX+(2*xi), posY+(2*yi)) || checkPosition(board, posX+(2*xi), posY+(2*yi)) {
		return false
	}
	if !checkOnTheBoard(posX+(3*xi), posY+(3*yi)) || checkPosition(board, posX+(3*xi), posY+(3*yi)) {
		return false
	}
	if currentPlayer != board[uint8(posY+(2*yi))][uint8(posX+(2*xi))] {
		if currentPlayer == board[uint8(posY+(3*yi))][uint8(posX+(3*xi))] {
			r.caps = append(r.caps, NewCapture(posX+xi, posY+yi))
			r.caps = append(r.caps, NewCapture(posX+(xi*2), posY+(yi*2)))
			return true
		}
	}
	return false
}

func (r *Rules) GetCaptures() []Capture {
	return r.caps
}

func (r *Rules) CheckCapture(board [][]uint8, posX, posY int8, currentPlayer uint8) bool {
	var captured bool

	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {
			x := posX + xi
			y := posY + yi
			if (xi == 0 && yi == 0) || !checkOnTheBoard(x, y) {
				continue
			}
			if !checkPosition(board, x, y) && board[y][x] != currentPlayer {
				if r.SetCapture(board, posX, posY, xi, yi, currentPlayer) {
					captured = true
				}
			}
		}
	}
	return captured
}

func (r Rules) IsThree(board [][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) bool {
	var nbMe uint8

	for i := int8(-1); i <= 2; i++ {
		x := posX + (xi * i)
		y := posY + (yi * i)
		if !checkOnTheBoard(x, y) {
			return false
		}
		if board[y][x] == currentPlayer {
			nbMe++
			if nbMe > 2 {
				return false
			}
		} else if !checkPosition(board, x, y) {
			return false
		}
	}

	if nbMe == 2 {
		// check extremites
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

func (r Rules) CheckDoubleThree(board [][]uint8, posX, posY int8, currentPlayer uint8) bool {
	nbThree := 0

	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {
			if xi == 0 && yi == 0 {
				continue
			}
			if r.IsThree(board, posX, posY, xi, yi, currentPlayer) {
				nbThree++
			} else if r.IsThree(board, posX+xi, posY+yi, xi, yi, currentPlayer) {
				nbThree++
			}
			if nbThree >= 2 {
				return true
			}
		}
	}
	return false
}

func (r Rules) isWin(board [][]uint8, posX, posY, xi, yi int8, currentPlayer uint8) bool {
	var nbMe uint8

	for i := int8(-4); i <= 4; i++ {
		x := posX + (xi * i)
		y := posY + (yi * i)
		if !checkOnTheBoard(x, y) {
			return false
		}
		if board[y][x] == currentPlayer {
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

func (r Rules) CheckWinner(board [][]uint8, posX, posY int8, currentPlayer uint8) bool {

	if r.isWin(board, posX, posY, 1, -1, currentPlayer) {
		return true
	} else if r.isWin(board, posX, posY, 1, 0, currentPlayer) {
		return true
	} else if r.isWin(board, posX, posY, 1, 1, currentPlayer) {
		return true
	} else if r.isWin(board, posX, posY, 0, 1, currentPlayer) {
		return true
	}
	return false
}

/*
func (r Rules) CheckIsFinish(board [][]uint8, posX, posY uint8) bool {
	// ex: posY = 0 && posX = 0

}*/

// Capture
// Free-three forbidden
// Double-three
// Game-ending
// Alignement

// Fonctionnalites:
//- Le joueur qui a choisi ou obtenu par un tirage au sort les pions noirs (et que l'on appelle Noir par convention) joue toujours le premier en plaçant son premier pion sur l'intersection centrale du damier