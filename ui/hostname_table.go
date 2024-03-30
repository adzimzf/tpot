package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"log"
	"os"
	"sort"
	"strings"
	"time"
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

// findMatchPositions return the position and matches char of keyword in the string
// eg:
//
//	 keyword: this
//	 str: this is word
//	 return:
//			matches: this
//			positions: [0,3]
func findMatchPositions(keyword, str string) (matches []string, positions [][]int) {
	for pos, char := range str {
		foundPos := -1
		for _, keywordChar := range keyword {
			if char == keywordChar {
				foundPos = pos
				break
			}
		}
		if foundPos != -1 {
			matches = append(matches, string(char))
			positions = append(positions, []int{foundPos, foundPos + 1})
			str = str[:foundPos] + " " + str[foundPos+1:]
		}

	}
	return matches, positions
}

// lookup the keyword in the datum (list of host)
func lookup(keyword string, datum []string) []stringResult {

	// empty string is not needed, but the datum might contain it.
	var tmpDatum []string
	for _, s := range datum {
		if s != "" {
			tmpDatum = append(tmpDatum, s)
		}
	}
	datum = tmpDatum

	var res []stringResult
	matchStrings := fuzzy.RankFindNormalizedFold(keyword, datum)

	if keyword == "" {
		sort.Slice(matchStrings, func(i, j int) bool {
			return matchStrings[i].Target < matchStrings[j].Target
		})
	} else {
		sort.Slice(matchStrings, func(i, j int) bool {
			return matchStrings[i].Distance < matchStrings[j].Distance
		})
	}

	for _, matchString := range matchStrings {
		_, pos := findMatchPositions(keyword, matchString.Target)
		res = append(res, stringResult{
			MatchPositions: pos,
			Keyword:        keyword,
			Data:           matchString.Target,
		})
	}
	return res
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

type stringResult struct {
	MatchPositions [][]int
	Keyword        string
	Data           string
}

// colorizeMatchChars the matching string
func (s stringResult) colorizeMatchChars() string {
	str := strings.Builder{}
	colorized := false
	for i := 0; i < len(s.Data); i++ {
		for _, position := range s.MatchPositions {
			if i == position[0] { // start to colorizeMatchChars
				colorized = true
			}
			if i == position[1] { // reset colorizeMatchChars
				colorized = false
			}
		}
		if colorized {
			str.WriteString("\u001B[37;7m")
			str.WriteRune(rune(s.Data[i]))
			str.WriteString("\u001B[0m")
		} else {
			str.WriteRune(rune(s.Data[i]))
		}

	}
	return str.String()
}

// divChars the space before pipe
func (s stringResult) divChars() string {
	const maxSpace = 60
	str := strings.Builder{}
	for i := 0; i < maxSpace-len(s.Data); i++ {
		str.WriteString(" ")
	}
	str.WriteRune(dividerChar)
	return str.String()
}

// formatResult colorize & create table to be shown as a string
// d is a list of node
// keyword is a keyword to be colorize
// ap is the current arrow position
func formatResult(hosts []stringResult, keyword string, ap arrowPos) string {
	screenMaxY := maxScreenY - 3
	var res string
	var y, x int
	newList := make([]string, screenMaxY)

	for _, key := range hosts {
		formattedHost := strings.Builder{}
		// if the host is selected, the color and prefix will be different
		if y == ap.Y && x == ap.X {
			formattedHost.WriteString(arrowColorized)
			formattedHost.WriteString("\u001B[33;1m")
			formattedHost.WriteString(key.Data)
			formattedHost.WriteString("\u001B[0m")
		} else {
			formattedHost.WriteString("   ")
			formattedHost.WriteString(key.colorizeMatchChars())
		}
		formattedHost.WriteString(key.divChars())

		newList[y] += formattedHost.String()
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

// cleanText clear the text from color character
func cleanText(s string) string {
	chars := []string{" ", "\u001B[33;1m", "\u001B[0m", "\u001B[37;7m", "\u001B[0m"}
	for _, c := range chars {
		s = strings.Replace(s, c, "", -1)
	}
	return s
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

func Debug(i ...interface{}) {
	go func() {
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
	}()
}
