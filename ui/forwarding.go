package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/adzimzf/tpot/config"
	"github.com/jroimartin/gocui"
)

func NewForwarding(list []*config.ForwardingNode) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	maxScreenX, maxScreenY = g.Size()

	g.Cursor = false
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(func(gui *gocui.Gui) error {
		x, y := g.Size()
		if _, err := g.SetView(searchInputView, 0, 0, x-1, y-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			if _, err := g.SetCurrentView(searchInputView); err != nil {
				return err
			}
		}

		x0, x1 := 0, 0
		y0, y1 := 1, 5

		for _, forwarding := range list {
			x0, x1 = 1+x1, x1+40
			if x1 > x {
				y0, y1 = y1+1, y1+5
				x0, x1 = 1, 40
			}

			viewName := forwarding.ViewName()
			if v, err := g.SetView(viewName, x0, y0, x1, y1); err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				v.Title = forwarding.Host
				v.FgColor = gocui.ColorGreen
				_, err := v.Write([]byte(fmt.Sprintf("listen: %s\nto    : %s:%s", forwarding.ListenPort, forwarding.RemoteHost, forwarding.RemotePort)))
				if err != nil {
					return err
				}
				if _, err := g.SetCurrentView(viewName); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	go updateStatus(g, list)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return
}

func updateStatus(g *gocui.Gui, list []*config.ForwardingNode) {
	for {
		time.Sleep(1 * time.Second)
		g.Update(func(gui *gocui.Gui) error {
			for _, node := range list {
				view, err := gui.View(node.ViewName())
				if err != nil {
					return err
				}
				view.FgColor = gocui.ColorRed
				if node.Status {
					view.FgColor = gocui.ColorGreen
				}
				view.Clear()
				view.Write([]byte(fmt.Sprintf("listen: %s\nto    : %s:%s\nerror : %s", node.ListenPort, node.RemoteHost, node.RemotePort, node.Error)))
			}
			return nil
		})

	}
}
