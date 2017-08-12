package controller

import (
	"strconv"
	"time"
)

type Control Gomoku;Menu {
	notice *text.Text
	renderer *sdl.Renderer
}

func (c *Control) SetNotice(str string) {
	idSDL := c.notice.NewIDSDL()
	if c.notice.IsInit() == true {
		c.notice.Close()
	}
	c.notice.UpdateText(str, c.renderer)
	if err := c.notice.Init(c.renderer); err != nil {
		panic(errors.New(objects.ErrorRenderer))
	}
	time.Sleep(3 * time.Second)
	if c.notice.GetIDSDL() == idSDL {
		c.notice.Close()
	}
}

// TimeToString transform duration to string with format 00:00:00
func TimeToString(d time.Duration) string {
	var str string

	// get hour, min and sec
	hour := int(d.Hours())
	min := int(d.Minutes())
	sec := int(d.Seconds())

	sec -= (min * 60)
	min -= (hour * 60)

	// if hour is on only one char, add 0
	if hour < 10 {
		str += "0"
	}
	// add hour to string
	str += strconv.Itoa(hour) + ":"

	// if min is on only one char, add 0
	if min < 10 {
		str += "0"
	}
	// add min to string
	str += strconv.Itoa(min) + ":"

	// if sec is on only one char, add 0
	if sec < 10 {
		str += "0"
	}
	// add min to string
	str += strconv.Itoa(sec)

	return str
}
