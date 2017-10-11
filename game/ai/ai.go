package ai

type AI struct {
}

func New() *AI {
	return &AI{}
}

func (a AI) Play(b [][]uint8) (uint8, uint8) {
	return 9, 10
}
