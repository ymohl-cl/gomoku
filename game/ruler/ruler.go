package ruler

import "fmt"

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

func (r *Rules) SetCapture(board [][]uint8, posX, posY int8, xi, yi int8) {
	fmt.Println("SetCapture")
	if !checkOnTheBoard(posX+(2*xi), posY+(2*yi)) || checkPosition(board, posX+(2*xi), posY+(2*yi)) {
		return
	}
	if !checkOnTheBoard(posX+(3*xi), posY+(3*yi)) || checkPosition(board, posX+(3*xi), posY+(3*yi)) {
		return
	}

	fmt.Println("...")
	if board[uint8(posY)][uint8(posX)] != board[uint8(posY+(2*yi))][uint8(posX+(2*xi))] {
		fmt.Println("neigh... ok x: ", xi*2, "y: ", yi*2)
		if board[uint8(posY)][uint8(posX)] == board[uint8(posY+(3*yi))][uint8(posX+(3*xi))] {
			fmt.Println("capture ok")
			r.caps = append(r.caps, NewCapture(posX+xi, posY+yi))
			r.caps = append(r.caps, NewCapture(posX+(xi*2), posY+(yi*2)))
			return
		}
	}
}

func (r *Rules) GetCaptures() []Capture {
	return r.caps
}

func (r *Rules) CheckCapture(board [][]uint8, posX, posY int8, currentPlayer uint8) {
	for yi := int8(-1); yi <= 1; yi++ {
		for xi := int8(-1); xi <= 1; xi++ {
			x := posX + xi
			y := posY + yi
			if (xi == 0 && yi == 0) || !checkOnTheBoard(x, y) {
				continue
			}
			if !checkPosition(board, x, y) && board[y][x] != currentPlayer {
				r.SetCapture(board, posX, posY, xi, yi)
			}
		}
	}
	return
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
