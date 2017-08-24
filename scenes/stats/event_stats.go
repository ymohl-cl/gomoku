package stats

import (
	"github.com/ymohl-cl/gomoku/conf"
)

/*
** Endpoint action from objects click
 */

// Return to menu
func (s *Stats) Return(values ...interface{}) {
	var err error

	go func() {
		if err = s.switcher(conf.SMenu, true); err != nil {
			panic(err)
		}
	}()
}
