package main

import (
	"github.com/ymohl-cl/game-builder/drivers"
	"github.com/ymohl-cl/game-builder/scripter"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/scenes/gomoku"
	"github.com/ymohl-cl/gomoku/scenes/loader"
	"github.com/ymohl-cl/gomoku/scenes/menu"
)

func main() {
	var err error
	var d drivers.VSDL
	var data *database.Data

	// init drivers sdl from game-builder
	if d, err = drivers.Init(conf.WindowWidth, conf.WindowHeight, conf.Title); err != nil {
		panic(err)
	}
	defer d.Destroy()

	// get data game
	if data, err = database.Get(); err != nil {
		panic(err)
	}

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
	var menuScene *menu.Menu
	if menuScene, err = menu.New(data, d.GetRenderer()); err != nil {
		panic(err)
	}
	if err = s.AddScene(menuScene, conf.SMenu, true); err != nil {
		panic(err)
	}

	// get and add stat scene

	// get and add gomoku scene
	var gameScene *gomoku.Gomoku
	if gameScene, err = gomoku.New(data, d.GetRenderer()); err != nil {
		panic(err)
	}
	if err = s.AddScene(gameScene, conf.SGame, false); err != nil {
		panic(err)
	}

	// run application
	s.Run(d)
}
