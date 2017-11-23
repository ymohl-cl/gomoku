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
