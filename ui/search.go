package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type search struct {
	g *gocui.Gui
}

func (s *search) register(data []string) error {
	return s.g.SetKeybinding("search_input", gocui.KeyTab, gocui.ModNone, func(gui *gocui.Gui, inputV *gocui.View) error {
		resultV, err := gui.View("search_result")
		if err != nil {
			return err
		}
		resultV.Overwrite = true
		resultV.Clear()
		keyword := inputV.Buffer()
		fmt.Fprint(resultV, formatResult(lookup(keyword, data), keyword, arrowPos{}))
		return err
	})
}

func newSearch(g *gocui.Gui) *search {
	return &search{g: g}
}
