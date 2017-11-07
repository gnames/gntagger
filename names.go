package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"

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

func (a *Annotation) Genus() string {
	return "Genus"
}

func (a *Annotation) Species() string {
	return "Species"
}

func (a *Annotation) Doubtful() string {
	return "Doubtful"
}

type Names struct {
	// Path to json file with names
	Path string
	// Data is a gnfinder output
	Data gnfinder.OutputJSON
}

func (n *Names) Save() error {
	json := n.Data.ToJSON()
	return ioutil.WriteFile(n.Path, json, 0644)
}

func nameStrings(n *gnfinder.NameJSON, current bool, i int, total int) []string {
	name := make([]string, 4)
	nameString := n.Name
	if current {
		nameString = fmt.Sprintf("\033[43;30;1m%s\033[0m", nameString)
	}
	name[0] = fmt.Sprintf("    %d/%d", i, total)
	name[1] = fmt.Sprintf("Type: %s", n.Type)
	name[2] = fmt.Sprintf("Name: %s", nameString)
	name[3] = annotation(n.Annotation)
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
	case annotation.Doubtful():
		color = 37 //light grey
	case annotation.Species():
		color = 35 //magenta
	case annotation.Genus():
		color = 35 //magenta
	case annotation.Uninomial():
		color = 35 //magenta
	default:
		color = 33
	}
	return fmt.Sprintf("\033[%d;40;2mAnnot: %s\033[0m", color, a)
}

func NamesFromJSON(path string) *Names {
	o := gnfinder.OutputJSON{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	o.FromJSON(b)
	return &Names{path, o}
}

type Text struct {
	// Path to the text file
	Path string
	// Content of the file converted to runes
	Original []rune
	// Vertical offset of the text to align it with name
	OffsetY int
}

func (t *Text) Markup(n *Names) string {
	var markup []rune
	name := n.Data.Names[n.Data.Meta.CurrentName]
	markup = append(markup, t.Original[0:name.OffsetStart]...)
	markup = append(markup, []rune("\033[40;33;1m")...)
	markup = append(markup, t.Original[name.OffsetStart:name.OffsetEnd]...)
	markup = append(markup, []rune("\033[0m")...)
	markup = append(markup, t.Original[name.OffsetEnd:len(t.Original)]...)
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

func PrepareText(path string) *Text {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	return &Text{path, []rune(string(b)), 0}
}

func (names *Names) currentName() *gnfinder.NameJSON {
	return &names.Data.Names[names.Data.Meta.CurrentName]
}
