package ui

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

// maxScreenX maximum screen wide to show table UI
// maxScreenY maximum screen high to show table UI
var maxScreenX, maxScreenY int

const (
	// dividerChar is a character to create table
	dividerChar = '│'

	// arrowColorized is a character ( > ) to indicate the current selected item
	arrowColorized = "\u001B[33;1m" + " > " + "\u001B[0m"
)

// GetSelectedHost will prompt user an table UI, and let the user
// select node list by typing or moving with an arrow
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

	s := newSearch(g, hosts)
	if err := s.register(); err != nil {
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

func lookup(keyword string, datum []string) map[string]stringResult {
	res := make(map[string]stringResult, len(datum))
	for _, data := range datum {
		if data == "" {
			continue
		}
		keyword = strings.TrimSpace(keyword)
		if strings.Contains(data, keyword) {
			res[data] = stringResult{
				FormattedData: data,
			}
		}
	}
	return res
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

type stringResult struct {
	CharMatch     int
	FormattedData string
}

// formatResult colorize & create table to be shown as a string
// d is a list of node
// keyword is a keyword to be colorize
// ap is the current arrow position
func formatResult(d map[string]stringResult, keyword string, ap arrowPos) string {
	screenMaxY := maxScreenY - 3
	var res string
	var y, x int
	newList := make([]string, screenMaxY)
	for _, key := range sortKey(d) {
		prefix := "   "
		formattedHost := colorizeSelectedWord(d[key].FormattedData, keyword)
		if y == ap.Y && x == ap.X {
			prefix = arrowColorized
			formattedHost = fmt.Sprintf("\u001B[33;1m%s\u001B[0m", d[key].FormattedData)
		}
		newList[y] += fmt.Sprintf("%s%-60s%s", prefix, formattedHost, string(dividerChar))
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

// arrowPos contain the X and Y of array position in the table list
type arrowPos struct {
	X, Y int
}

// sortKey sort the table item from A-Z to improve readability
func sortKey(d map[string]stringResult) []string {
	var res []string
	for s := range d {
		res = append(res, s)
	}
	sort.Strings(res)
	return res
}

// cleanText clear the text from color character
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
		for j, q := range strings.Split(s, string(dividerChar)) {
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
		X: len(strings.Split(s, string(dividerChar))),
		Y: len(strings.Split(strings.Trim(s, "\n"), "\n")),
	}

}

func debug(i ...interface{}) {
	s := time.Now().String() + "\n"
	for _, i1 := range i {
		s += fmt.Sprintf("%v", i1)
	}
	f, err := os.OpenFile(".debug", os.O_APPEND|os.O_RDWR, 777)
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
