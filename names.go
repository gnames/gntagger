package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gnames/gnfinder"
)

type AnnotationId int

const (
	AnnotationNotAssigned AnnotationId = iota
	AnnotationNotName
	AnnotationAccepted
	AnnotationUninomial
	AnnotationGenus
	AnnotationSpecies
	AnnotationDoubtful
)

type Names struct {
	// Path to json file with names
	Path string
	// Data is a gnfinder output
	Data gnfinder.OutputJSON
}

var annotationNames = []string{
	"",
	"NotName",
	"Accepted",
	"Uninomial",
	"Genus",
	"Species",
	"Doubtful",
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
	name[0] = fmt.Sprintf("    %d/%d", i+1, total)
	name[1] = fmt.Sprintf("Type: %s", n.Type)
	name[2] = fmt.Sprintf("Name: %s", nameString)
	name[3] = annotationOfName(n.Annotation).formatString()
	return name
}

func annotationOfName(annotationName string) AnnotationId {
	if annotationName == "NotAssigned" {
		return AnnotationNotAssigned
	}
	for idx, an := range annotationNames {
		if an == annotationName {
			return AnnotationId(idx)
		}
	}
	return AnnotationNotAssigned
}

func (annotation AnnotationId) color() int {
	switch annotation {
	case AnnotationAccepted:
		return 32 //green
	case AnnotationNotName:
		return 31 //red
	case AnnotationDoubtful:
		return 37 //light grey
	case AnnotationSpecies:
		return 35 //magenta
	case AnnotationGenus:
		return 35 //magenta
	case AnnotationUninomial:
		return 35 //magenta
	case AnnotationNotAssigned:
	default:
		return 33
	}
	return 33
}

func (annotation AnnotationId) name() string {
	return annotationNames[int(annotation)]
}

func (annotation AnnotationId) formatString() string {
	color := annotation.color()
	name := annotation.name()
	return fmt.Sprintf("\033[%d;40;2mAnnot: %s\033[0m", color, name)
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
