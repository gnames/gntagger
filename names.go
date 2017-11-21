package gntagger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"

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
	Data gnfinder.Output
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

func nameStrings(n *gnfinder.Name, current bool, i int,
	total int) ([]string, error) {
	name := make([]string, 4)
	nameString := n.Name
	if current {
		nameString = fmt.Sprintf("\033[43;30;1m%s\033[0m", nameString)
	}
	name[0] = fmt.Sprintf("    %d/%d", i+1, total)
	if n.Odds != 0.0 {
		name[1] = fmt.Sprintf("Log Odds: %0.2f", math.Log10(n.Odds))
	} else {
		name[1] = fmt.Sprintf("Type: %s", n.Type)
	}
	name[2] = fmt.Sprintf("Name: %s", nameString)
	annotation, err := annotationOfName(n.Annotation)
	if err != nil {
		return nil, err
	}
	name[3] = annotation.formatString()
	return name, nil
}

func annotationOfName(annotationName string) (AnnotationId, error) {
	if annotationName == "NotAssigned" {
		return AnnotationNotAssigned, nil
	}
	for idx, an := range annotationNames {
		if an == annotationName {
			return AnnotationId(idx), nil
		}
	}
	errMsg :=
		fmt.Sprintf("annotation name %s isn't supported", annotationName)
	return -1, errors.New(errMsg)
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
	o := gnfinder.Output{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	o.FromJSON(b)
	return &Names{path, o}
}

func (names *Names) currentName() *gnfinder.Name {
	return &names.Data.Names[names.Data.Meta.CurrentName]
}
