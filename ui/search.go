package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type search struct {
	g     *gocui.Gui
	hosts []string
}

const (
	autoCompleteChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()_-+:;\"/,.<>|\\{}|[] "
	searchInputView   = "search_input"
	searchResultView  = "search_result"
)

func (s *search) register() error {
	for _, r := range autoCompleteChars {
		err := s.g.SetKeybinding(searchInputView, r, gocui.ModNone, s.handleType(r))
		if err != nil {
			return err
		}
	}

	keyBackSpaceList := []gocui.Key{gocui.KeyBackspace, gocui.KeyBackspace2}
	for _, key := range keyBackSpaceList {
		err := s.g.SetKeybinding(searchInputView, key, gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {

			if len(v.Buffer()) < 1 {
				return nil
			}

			text := v.Buffer()[:len(v.Buffer())-2]
			v.Clear()
			newText := strings.Trim(text, "\n")
			if err := v.SetCursor(len(newText), 0); err != nil {
				return err
			}
			if _, err := fmt.Fprint(v, newText); err != nil {
				return err
			}
			return s.updateResult(text, gui)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *search) handleType(c rune) func(gui *gocui.Gui, inputV *gocui.View) error {
	return func(gui *gocui.Gui, v *gocui.View) error {
		text := v.Buffer()
		v.Clear()
		newText := strings.Trim(text, "\n") + string(c)
		if err := v.SetCursor(len(newText), 0); err != nil {
			return err
		}
		if _, err := fmt.Fprint(v, newText); err != nil {
			return err
		}
		return s.updateResult(newText, gui)
	}
}

func (s *search) updateResult(keyword string, gui *gocui.Gui) error {
	resultV, err := gui.View(searchResultView)
	if err != nil {
		return err
	}
	resultV.Overwrite = true
	resultV.Clear()
	_, err = fmt.Fprint(resultV, formatResult(lookup(keyword, s.hosts), keyword, arrowPos{}))
	return err
}

func newSearch(g *gocui.Gui, hosts []string) *search {
	return &search{g: g, hosts: hosts}
}
