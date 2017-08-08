package loader

import (
	"time"

	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/gomoku/conf"
)

/*
** Endpoint action from objects click
 */

func (l *Load) addLoadingBar() {
	var b *block.Block
	var err error
	var loop bool

	loop = true
	for loop {
		select {
		case <-l.closer:
			loop = false
		default:
			if b, err = l.lastLoadBlock.Clone(l.renderer); err != nil {
				panic(err)
			}
			x, y := l.lastLoadBlock.GetPosition()
			if x+conf.LoadBlockWidth*2 > conf.WindowWidth-conf.LoadBlockWidth {
				l.resetLoadingBlock()
			} else {
				b.UpdatePosition(x+conf.LoadBlockWidth, y)
				l.m.Lock()
				l.layers[layerLoadingBar] = append(l.layers[layerLoadingBar], b)
				l.m.Unlock()
				l.lastLoadBlock = b
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	l.resetLoadingBlock()
}

func (l *Load) resetLoadingBlock() {
	l.m.Lock()
	l.lastLoadBlock = l.layers[layerLoadingBar][0].(*block.Block)
	del := l.layers[layerLoadingBar][1:]
	l.layers[layerLoadingBar] = l.layers[layerLoadingBar][:1]
	l.m.Unlock()
	go clearLoadingBar(del)
}

func clearLoadingBar(sl []objects.Object) {
	var err error

	for _, v := range sl {
		if err = v.Close(); err != nil {
			panic(err)
		}
	}
}
