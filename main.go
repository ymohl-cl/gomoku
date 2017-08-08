package main

import (
	"github.com/ymohl-cl/game-builder/drivers"
	"github.com/ymohl-cl/game-builder/scripter"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/scenes/loader"
)

func main() {
	var err error
	var d drivers.VSDL

	// init drivers sdl from game-builder
	if d, err = drivers.Init(conf.WindowWidth, conf.WindowHeight, conf.Title); err != nil {
		panic(err)
	}
	defer d.Destroy()

	// get new scripter application
	s := scripter.New()

	// get and add loader scene
	var loaderScene *loader.Load
	if loaderScene, err = loader.New(nil, d.GetRenderer()); err != nil {
		panic(err)
	}
	if err = s.AddLoader(loaderScene); err != nil {
		panic(err)
	}

	// get and add menu scene

	// get and add stat scene

	// get and add gomoku scene

	// run application
	s.Run(d)
}
