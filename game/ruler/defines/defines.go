package defines

// outOfBoard is used for the mask of raw (see ruler.CheckRules)
const (
	Empty      = 0
	Player1    = uint8(1)
	Player2    = uint8(2)
	OutOfBoard = 3

	SizeY = 19
	SizeX = 19

	AlignFree    = 1 << 0
	AlignHalf    = 1 << 1
	AlignFlanked = 2 << 2
)

// GetOtherPlayer provide the value of player's opponent provided on parameter
func GetOtherPlayer(p uint8) uint8 {
	if p == Player1 {
		return Player2
	}
	return Player1
}

// IsOnTheBoard : check if the (y, x) point are inside the board
func IsOnTheBoard(y, x int8) bool {
	if x < 0 || y < 0 || x >= 19 || y >= 19 {
		return false
	}
	return true
}

// IsEmpty : check if the (pX, pY) point are Empty position
func IsEmpty(b *[19][19]uint8, y, x int8) bool {
	if (*b)[uint8(y)][uint8(x)] == Empty {
		return true
	}
	return false
}
