package alphabeta

import "fmt"

func (s *State) printEval(n *Node, depth uint8, score int16) {
	fmt.Println("Score eval: ", score, " - depth: ", depth)

	depthNode := s.maxDepth - depth
	for node := n; node != nil; node = node.prev {
		y, x := node.rule.GetPosition()
		fmt.Println("move ", depthNode, " on y: ", y, " - x: ", x)
		depthNode--
	}
	return
}

/*
Score eval:  -16358  - depth:  0
move  4  on y:  7  - x:  8
move  3  on y:  10  - x:  8
move  2  on y:  8  - x:  8
move  1  on y:  8  - x:  7
--
Score eval:  -16358  - depth:  0
move  4  on y:  8  - x:  9
move  3  on y:  7  - x:  6
move  2  on y:  8  - x:  8
move  1  on y:  8  - x:  7
--

*/
