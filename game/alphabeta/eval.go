package alphabeta

import "fmt"

func maxWeight(v1, v2 int8) int8 {
	if v1 >= v2 {
		return v1
	}
	return v2
}

func eval(n *Node, depth uint8) int8 {
	if n.rule.Win {
		return -127 + int8(maxDepth-depth)
	}

	nb := 0
	for node := n; node != nil; node = node.prev {
		nb++
	}
	fmt.Println("Number node: ", nb)
	return -int8(60 + n.rule.NumberCapture)
}

/*func eval(s *State, stape uint8, r *ruler.Rules, base *State, maxDepth uint8) int8 {
	ret := int8(63)

	// Var Caps
	var basicCurrent int8
	var basicOther int8
	if s.player == ruler.Player2 {
		basicCurrent = int8(s.nbCapsOther - base.nbCapsOther)
		basicOther = int8(s.nbCapsCurrent - base.nbCapsCurrent)
	} else {
		basicCurrent = int8(s.nbCapsOther - base.nbCapsCurrent)
		basicOther = int8(s.nbCapsCurrent - base.nbCapsOther)
	}

	//Win
	if r.Win {
		if basicCurrent > 0 {
			return -127 + int8(maxDepth-stape)
		} else {
			return -117 + int8(maxDepth-stape)
		}

	}

	//Caps
	ret += (basicCurrent*7 - basicOther*6)

	//Three
	treeCurrent := s.nbTOther
	treeOther := s.nbTCurrent

	//	if len(r.CapturableWin) > 0 {
	//		ret += int8(treeCurrent - treeOther*3)
	//	} else {
	ret += int8(treeCurrent*3 - treeOther)
	//	}

	//Align
	AlignCurrent := s.nbAlignOther
	AlignOther := s.nbAlignCurrent

	ret += int8(AlignCurrent*4 - AlignOther)

	return -ret
}*/
