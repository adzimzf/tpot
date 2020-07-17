package ui

import (
	"strings"

	"github.com/jroimartin/gocui"
)

type keyEnterBinding struct {
	g *gocui.Gui
}

func (k *keyEnterBinding) register(result *string) error {
	return k.g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		v, err := gui.View("search_result")
		if err != nil {
			return err
		}
		*result = k.findResult(v.Buffer())
		return gocui.ErrQuit
	})
}

func (k *keyEnterBinding) findResult(s string) string {
	for _, st := range strings.Split(s, "\n") {
		if strings.Contains(st, ">") {
			for _, s2 := range strings.Split(st, "|") {
				if strings.Contains(s2, ">") {
					return cleanText(strings.Replace(s2, ">", "", -1))
				}
			}
		}
	}
	return ""
}

func newKeyEnterBinding(g *gocui.Gui) *keyEnterBinding {
	return &keyEnterBinding{g: g}
}
