package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gnames/gnfinder"
)

type Annotation struct{}

func (a *Annotation) NotName() string {
	return "NotName"
}

func (a *Annotation) Accepted() string {
	return "Accepted"
}

func (a *Annotation) Uninomial() string {
	return "Uninomial"
}

func (a *Annotation) Species() string {
	return "Species"
}

func (a *Annotation) Undecided() string {
	return "Undecided"
}

type Names struct {
	Data    gnfinder.OutputJSON
	Current int
}

func (n *Names) String() string {
	var current bool
	names := n.Data.Names
	str := make([]string, len(names)*4)
	for i, v := range names {
		current = (i == n.Current)
		for j, w := range nameStrings(&v, current) {
			str[i*4+j] = w
		}
	}
	return strings.Join(str, "\n")
}

func nameStrings(n *gnfinder.NameJSON, current bool) []string {
	name := make([]string, 4)
	nameString := n.Name
	if current {
		nameString = fmt.Sprintf("\033[43;30;1m%s\033[0m", nameString)
	}
	name[0] = fmt.Sprintf("Type: %s", n.Type)
	name[1] = fmt.Sprintf("Name: %s", nameString)
	name[2] = annotation(n.Annotation)
	return name
}

func annotation(a string) string {
	var color int
	annotation := Annotation{}
	switch a {
	case annotation.Accepted():
		color = 32 //green
	case annotation.NotName():
		color = 31 //red
	case annotation.Undecided():
		color = 37 //light grey
	case annotation.Species():
		color = 92 //light green
	case annotation.Uninomial():
		color = 92 //light green
	default:
		color = 33
	}
	return fmt.Sprintf("\033[%d;40;2mAnnot: %s\033[0m", color, a)
}

func NamesFromJSON() *Names {
	o := gnfinder.OutputJSON{}
	b, err := ioutil.ReadFile("../../testdata/seashells_book.json")
	if err != nil {
		log.Panicln(err)
	}
	o.FromJSON(b)
	return &Names{o, 0}
}

type Text struct {
	Original []rune
	OffsetY  int
}

func (t *Text) Markup(n *Names) string {
	markup := make([]rune, len(t.Original)+30)
	name := n.Data.Names[n.Current]
	tailLen := len(t.Original) - name.OffsetEnd
	markup = append(markup, t.Original[0:name.OffsetStart]...)
	markup = append(markup, []rune("\033[40;33;1m")...)
	markup = append(markup, t.Original[name.OffsetStart:name.OffsetEnd]...)
	markup = append(markup, []rune("\033[0m")...)
	markup = append(markup, t.Original[name.OffsetEnd:tailLen]...)
	t.OffsetY = newLinesNum(t.Original[0:name.OffsetStart])
	return string(markup)
}

func newLinesNum(rs []rune) int {
	offset := 1
	for _, v := range rs {
		if v == '\n' || v == '\v' {
			offset++
		}
	}
	return offset
}

func PrepareText() *Text {
	b, err := ioutil.ReadFile("../../testdata/seashells_book.txt")
	if err != nil {
		log.Panicln(err)
	}
	return &Text{[]rune(string(b)), 0}
}
