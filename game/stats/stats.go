package stats

import "time"

type Stat struct {
	timer    time.Duration
	moveY    uint8
	moveX    uint8
	eatToken bool
}

func New(x, y uint8, t time.Duration, eat bool) Stat {
	return Stat{
		timer:    t,
		moveY:    y,
		moveX:    x,
		eatToken: eat,
	}
}
