package ui

import (
	"github.com/jroimartin/gocui"
	"strings"
)

type navDir int

const (
	NavUp navDir = iota
	NavDown
	NavLeft
	NavRight
)

var NavKeyMap = map[navDir]gocui.Key{
	NavUp:    gocui.KeyArrowUp,
	NavDown:  gocui.KeyArrowDown,
	NavLeft:  gocui.KeyArrowLeft,
	NavRight: gocui.KeyArrowRight,
}

type ArrowNav struct {
	g *gocui.Gui
	s *search
}

func (a *ArrowNav) newHandler(dir navDir) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		resultV, err := g.View(searchResultView)
		if err != nil {
			return err
		}
		inputV, err := g.View(searchInputView)
		if err != nil {
			return err
		}

		text := cleanText(resultV.Buffer())
		keyword := strings.TrimSpace(inputV.Buffer())
		pos, data := findArrowPos(text)
		nextPos := a.nextPos(pos, findMaxXY(text), dir)
		resultV.Clear()

		_, err = resultV.Write([]byte(formatResult(lookup(keyword, data), keyword, nextPos)))
		return err
	}
}

// nextPos decide the next arrow position
func (a *ArrowNav) nextPos(cp arrowPos, m maxXY, dir navDir) arrowPos {
	switch dir {
	case NavUp:
		if p := cp.Y - 1; p > -1 {
			cp.Y = p
		}
	case NavDown:
		if p := cp.Y + 1; p < m.Y {
			cp.Y = p
		}
	case NavLeft:
		if p := cp.X - 1; p > -1 {
			cp.X = p
		}
	case NavRight:
		if p := cp.X + 1; p < m.X {
			cp.X = p
		}
	}
	return cp
}

func (a *ArrowNav) registerArrowNav() error {
	for dir, key := range NavKeyMap {
		if err := a.g.SetKeybinding(searchInputView, key, gocui.ModNone, a.newHandler(dir)); err != nil {
			return err
		}
	}
	return nil
}

func NewArrowNav(g *gocui.Gui) *ArrowNav {
	return &ArrowNav{g: g}
}
