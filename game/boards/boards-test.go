package boards

/*
	This file provide more situation of board to check them in the tests.
	The schema are represent by:
	 x == p1 && o == p2 && . == token empty
*/

// GetEmptyBoard : _
func GetEmptyBoard() *[19][19]uint8 {
	var b [19][19]uint8

	return &b
}

// GetSimpleBoardP2 : _
// . . . . o
// . x x o .
// . . o . .
func GetSimpleBoardP2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][7] = 1
	b[9][8] = 1

	// set P2 on the board
	b[9][9] = 2
	b[10][8] = 2
	b[8][10] = 2
	return &b
}

// GetBoardFilledRightP2 : _
// . . . . o . . . . . . o x
// . x x o . x o x o x o x o
// . . o . . . . . . . o x .
func GetBoardFilledRightP2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][8] = 1
	b[9][7] = 1
	b[9][11] = 1
	b[9][13] = 1
	b[9][15] = 1
	b[9][17] = 1
	b[10][17] = 1
	b[8][18] = 1

	// set P2 on the board
	b[9][9] = 2
	b[10][8] = 2
	b[8][10] = 2
	b[9][12] = 2
	b[9][14] = 2
	b[9][16] = 2
	b[9][18] = 2
	b[8][17] = 2
	b[10][16] = 2
	return &b
}

// GetCaptureBoardP1_1 : _
// . o x x .
func GetCaptureBoardP1_1() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[0][17] = 1
	b[0][16] = 1

	// set P2 on the board
	b[0][15] = 2
	return &b
}

// GetCaptureBoardP2_1 : _
// . . x . . x . . x .
// . . . o x o x o . .
// . . . x o o o x . .
// x x x o o . o o x .
// . . . . o o o . . .
// . . . o . o x o x .
// . . x . . x . . x .
// . . . . . . . . . o
func GetCaptureBoardP2_1() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[10][9] = 1
	b[9][10] = 1
	b[10][12] = 1
	b[9][12] = 1
	b[7][12] = 1
	b[6][11] = 1
	b[4][12] = 1
	b[5][10] = 1
	b[4][9] = 1
	b[5][8] = 1
	b[4][6] = 1
	b[6][7] = 1
	b[7][6] = 1
	b[7][5] = 1
	b[7][4] = 1
	b[10][6] = 1

	// set P2 on the board
	b[9][9] = 2
	b[8][9] = 2
	b[9][11] = 2
	b[8][10] = 2
	b[7][11] = 2
	b[7][10] = 2
	b[5][11] = 2
	b[6][10] = 2
	b[5][9] = 2
	b[6][9] = 2
	b[5][7] = 2
	b[6][8] = 2
	b[7][7] = 2
	b[7][8] = 2
	b[8][8] = 2
	b[9][7] = 2
	b[11][13] = 2

	return &b
}

// GetBoardNoAlignP2_1 : _
// . . . o x o . . .
func GetBoardNoAlignP2_1() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][8] = 1

	// set P2 on the board
	b[9][9] = 2
	b[9][7] = 2
	return &b
}

// GetBoardNoAlignP2_2 : _
// . . x . . . . . . .
// . x . . . . . . . .
// . . . x . . . . . .
// . . . . o . . . . .
// . . . . x . . . . .
// . . . . . o o . . .
// . . . . . . x . . .
// . . . . . . . o . .
// . . . . . . . . x .
// . . . . . . . . . o
func GetBoardNoAlignP2_2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[8][8] = 1
	b[6][6] = 1
	b[4][4] = 1
	b[2][3] = 1
	b[1][1] = 1
	b[0][2] = 1

	// set P2 on the board
	b[9][9] = 2
	b[7][7] = 2
	b[5][5] = 2
	b[5][6] = 2
	b[3][4] = 2
	return &b
}

// GetBoardAlignFlankedP2 : _
// . . . . . x .
// o x x . x . o
// . . . o . o .
func GetBoardAlignFlankedP2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][9] = 1
	b[9][10] = 1
	b[9][12] = 1
	b[8][13] = 1

	// set P2 on the board
	b[9][8] = 2
	b[10][11] = 2
	b[10][13] = 2
	b[9][14] = 2
	return &b
}

// GetBoardAlignFreeP2 : _
// . x . .
// . o x .
// . o . .
func GetBoardAlignFreeP2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][9] = 1
	b[8][8] = 1

	// set P2 on the board
	b[9][8] = 2
	b[10][8] = 2
	return &b
}

// GetBoardAlignHalfP2 : _
// . . . . .
// . x . . .
// . o x . .
// . o . x .
// . . . . o
func GetBoardAlignHalfP2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][9] = 1
	b[8][8] = 1
	b[10][10] = 1

	// set P2 on the board
	b[9][8] = 2
	b[10][8] = 2
	b[11][11] = 2
	return &b
}

// GetTreeBoardP1 : _
// . . . . . . . . .
// . x . . . . . . .
// . . o o . . x o .
// . . . x o x o x .
// . . . . . . . . .
func GetTreeBoardP1() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][8] = 1
	b[9][10] = 1
	b[8][11] = 1
	b[9][12] = 1
	b[7][6] = 1

	// set P2 on the board
	b[9][9] = 2
	b[8][8] = 2
	b[9][11] = 2
	b[8][12] = 2
	b[8][7] = 2
	return &b
}

// GetWinCapturableBoardP2 : _
// . . . . . . .
// . . . . . x .
// o x x . x x o
// . . . o o o .
func GetWinCapturableBoardP2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][9] = 1
	b[9][10] = 1
	b[9][12] = 1
	b[8][13] = 1
	b[9][13] = 1

	// set P2 on the board
	b[9][8] = 2
	b[10][11] = 2
	b[10][13] = 2
	b[9][14] = 2
	b[10][12] = 2
	return &b
}

// GetWinNoCapturableBoardP2 : _
// . . . . . . .
// o x x . x x o
// . . . o . o .
func GetWinNoCapturableBoardP2() *[19][19]uint8 {
	var b [19][19]uint8

	// set P1 on the board
	b[9][9] = 1
	b[9][10] = 1
	b[9][12] = 1
	b[9][13] = 1

	// set P2 on the board
	b[9][8] = 2
	b[10][11] = 2
	b[10][13] = 2
	b[9][14] = 2
	return &b
}
