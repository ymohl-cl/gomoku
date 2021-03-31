package main

import (
	goui "github.com/ymohl-cl/go-ui"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/scenes/gomoku"
	"github.com/ymohl-cl/gomoku/scenes/loader"
	"github.com/ymohl-cl/gomoku/scenes/menu"
)

func main() {
	var err error
	var gb goui.GameBuilder
	var data *database.Data

	// init drivers sdl from go-ui
	if gb, err = goui.New(goui.ConfigUI{
		Window: goui.Window{
			Title:  conf.Title,
			Width:  conf.WindowWidth,
			Height: conf.WindowHeight,
		},
	}); err != nil {
		panic(err)
	}
	defer func() {
		if err := gb.Close(); err != nil {
			panic(err)
		}
	}()

	// get data game
	if data, err = database.Get(); err != nil {
		panic(err)
	}

	s := gb.Script()
	r := gb.Renderer().Driver()

	// get and add loader scene
	var loaderScene *loader.Load
	if loaderScene, err = loader.New(nil, r); err != nil {
		panic(err)
	}
	if err = s.SetLoader(loaderScene); err != nil {
		panic(err)
	}

	// get and add menu scene
	var menuScene *menu.Menu
	if menuScene, err = menu.New(data, r); err != nil {
		panic(err)
	}
	if err = s.AddScene("menu", menuScene); err != nil {
		panic(err)
	}

	// get and add gomoku scene
	var gameScene *gomoku.Gomoku
	if gameScene, err = gomoku.New(data, r); err != nil {
		panic(err)
	}
	if err = s.AddScene("game", gameScene); err != nil {
		panic(err)
	}

	// run application
	if err = gb.Run("menu"); err != nil {
		panic(err)
	}
	return
}
