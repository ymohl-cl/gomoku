package main

import (
	"fmt"

	"github.com/ymohl-cl/game-builder/builder"
	"github.com/ymohl-cl/game-builder/drivers"
	"github.com/ymohl-cl/game-builder/scripter"
	"github.com/ymohl-cl/gomoku/conf"
)

func main() {
	d, err := drivers.Init(conf.WindowWidth, conf.WindowHeight, conf.Title)
	defer d.Destroy()
	if err != nil {
		fmt.Println(err)
	}

	s := scripter.New(nil)
	builder.Run(d, s)
}
