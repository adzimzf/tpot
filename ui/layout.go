package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type layout struct {
	s *search
	g *gocui.Gui
}

func (l *layout) register(data []string) error {
	l.g.SetManagerFunc(func(gui *gocui.Gui) error {
		maxX, maxY := l.g.Size()
		if v, err := l.g.SetView(searchInputView, 0, maxY-3, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Editable = true
			if _, err := l.g.SetCurrentView(searchInputView); err != nil {
				return err
			}
		}
		if v, err := l.g.SetView(searchResultView, 0, 1, maxX-1, maxY-3); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Type to Search or Arrow to Navigate"
			v.Editable = true
			fmt.Fprintln(v, formatResult(lookup("", data), "", arrowPos{}))
		}
		return nil
	})
	return nil

}

func newLayout(g *gocui.Gui) *layout {
	return &layout{
		g: g,
	}
}
