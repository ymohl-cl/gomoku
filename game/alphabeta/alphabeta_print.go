package alphabeta

import (
	"fmt"

	"github.com/ymohl-cl/gomoku/game/boards"
)

func (i *InfoPlayer) Print() {
	fmt.Println("nb capture: ", i.totalCapture)
}

func (s *State) Print() {
	fmt.Println("maxDepth: ", s.maxDepth, " - player: ", s.currentPlayer)
	boards.Print(s.board)
	fmt.Print("p1 ")
	s.infoP1.Print()
	fmt.Print("p2 ")
	s.infoP2.Print()
}
