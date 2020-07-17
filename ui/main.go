package ui

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jroimartin/gocui"
)

var maxScreenX, maxScreenY int

func GetSelectedHost(hosts []string) string {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	maxScreenX, maxScreenY = g.Size()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	if err := newLayout(g).register(hosts); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err = NewArrowNav(g).registerArrowNav(); err != nil {
		log.Panicln(err)
	}

	s := newSearch(g)
	if err := s.register(hosts); err != nil {
		log.Panicln(err)
	}

	var result string
	if err := newKeyEnterBinding(g).register(&result); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return result

}

func lookup(keyword string, data []string) map[string]stringResult {
	res := make(map[string]stringResult, len(data))
	for _, datum := range data {
		if strings.Contains(datum, strings.TrimSpace(keyword)) {
			res[datum] = stringResult{
				FormatedData: datum,
			}
		}
	}
	return res
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

type stringResult struct {
	CharMatch    int
	FormatedData string
}

func formatResult(d map[string]stringResult, keyword string, ap arrowPos) string {
	screenMaxY := maxScreenY - 3
	var res string
	var y, x int
	newList := make([]string, screenMaxY)
	for _, key := range sortKey(d) {
		prefix := "   "
		formattedHost := colorizeSelectedWord(d[key].FormatedData, keyword)
		if y == ap.Y && x == ap.X {
			prefix = arrowString
			formattedHost = fmt.Sprintf("\u001B[33;1m%s\u001B[0m", d[key].FormatedData)
		}
		newList[y] += fmt.Sprintf("%s%-60s|", prefix, formattedHost)
		y++
		if y >= screenMaxY {
			x++
			y = 0
		}
	}
	for _, s := range newList {
		res += s + "\n"
	}
	return res
}

type arrowPos struct {
	X, Y int
}

func sortKey(d map[string]stringResult) []string {
	var res []string
	for s := range d {
		res = append(res, s)
	}
	sort.Strings(res)
	return res
}

const arrowString = "\u001B[33;1m" + " > " + "\u001B[0m"

func cleanText(s string) string {
	chars := []string{" ", "\u001B[33;1m", "\u001B[0m", "\u001B[37;7m", "\u001B[0m"}
	for _, c := range chars {
		s = strings.Replace(s, c, "", -1)
	}
	return s
}

func colorizeSelectedWord(text, keyword string) string {
	key := strings.TrimSpace(keyword)
	return strings.Replace(text,
		strings.TrimSpace(key),
		fmt.Sprintf("\u001B[37;7m%s\u001B[0m", key), 1)
}

func findArrowPos(res string) (ap arrowPos, data []string) {
	for i, s := range strings.Split(res, "\n") {
		for j, q := range strings.Split(s, "|") {
			if strings.Contains(q, ">") {
				ap.X = j
				ap.Y = i
			}
			data = append(data, strings.Trim(strings.TrimSpace(q), ">"))
		}
	}
	return
}

type maxXY struct {
	X, Y int
}

func findMaxXY(s string) (m maxXY) {
	return maxXY{
		X: len(strings.Split(s, "|")),
		Y: len(strings.Split(strings.Trim(s, "\n"), "\n")),
	}

}

func debug(s string) {
	f, err := os.OpenFile("debug.txt", os.O_APPEND|os.O_RDWR, 777)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(fmt.Sprintf("========\n%s\n======\n", s))
	if err != nil {
		log.Fatal(err)
	}
	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}

}
