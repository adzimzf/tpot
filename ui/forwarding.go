package ui

import (
	"fmt"
	"log"
	"sync"
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

// remoteHostMap stores the index char to show first
var remoteHostMap = sync.Map{}

type remoteHostMapper struct {
	index int
}

// refreshRemoteHost determine the remoteHost to be shown in the UI
// if the number of character is greater than maximum Box, then
// every 1 second it'll move 1 character to the left
func refreshRemoteHost(remoteHost string) string {
	const maxChar = 30

	if len(remoteHost) <= maxChar {
		return remoteHost
	}

	load, ok := remoteHostMap.Load(remoteHost)
	if !ok {
		// start show from the first character
		// when first loading
		remoteHostMap.Store(remoteHost, remoteHostMapper{
			index: 1,
		})
		return remoteHost
	}

	m, ok := load.(remoteHostMapper)
	if !ok {
		return remoteHost
	}

	if m.index >= len(remoteHost) {
		// restart the index once it reaches the maximum character
		remoteHostMap.Store(remoteHost, remoteHostMapper{
			index: 1,
		})
		return remoteHost
	}

	const maxShowNext = 25
	if len(remoteHost[m.index:]) <= maxShowNext {
		x := remoteHost[m.index:] + "  " + remoteHost[0:maxChar+2-len(remoteHost[m.index:])-4]
		m.index++
		remoteHostMap.Store(remoteHost, m)
		return x
	}

	m.index++
	remoteHostMap.Store(remoteHost, m)
	return remoteHost[m.index-1:]
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
				remoteHost := refreshRemoteHost(node.RemoteHost)
				view.Write([]byte(fmt.Sprintf("listen: %s\nto    : %s:%s\nerror : %s", node.ListenPort, remoteHost, node.RemotePort, node.Error)))
			}
			return nil
		})
	}
}
