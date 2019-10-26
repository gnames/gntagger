package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/output"
	"github.com/gnames/gntagger/annotation"
)

// Names is an object that keeps output of a name finder and the path where to
// save this data on disk
type Names struct {
	// Path to json file with names
	Path string
	// Data is a gnfinder output
	Data output.Output
}

// NewNames uses a name finder or existing information to return Names structure
// generated from a text
func NewNames(text *Text, gnt *GnTagger) *Names {
	opts := []gnfinder.Option{
		gnfinder.OptBayesThreshold(gnt.OddsLow),
		gnfinder.OptDict(dict.LoadDictionary()),
	}

	gnf := gnfinder.NewGNfinder(opts...)

	data := gnf.FindNames([]byte(string(text.Processed)))

	for i := range data.Names {
		n := &data.Names[i]
		if n.Odds != 0.0 && n.Odds < gnt.OddsHigh {
			n.Annotation = annotation.Doubtful.String()
		}
	}
	return &Names{Data: *data, Path: text.FilePath(NamesFile)}
}

// Save writes current state of names to file
func (n *Names) Save() error {
	json := n.Data.ToJSON()
	return ioutil.WriteFile(n.Path, json, 0644)
}

// NameStrings composes text to show in terminal gui
func NameStrings(n *output.Name, current bool, i int,
	total int) ([]string, error) {
	name := make([]string, 4)
	nameString := n.Name
	if current {
		nameString = fmt.Sprintf("\033[33;40;1m%s\033[0m", nameString)
	}
	name[0] = fmt.Sprintf("    %d/%d", i+1, total)
	name[1] = n.Type
	if n.Odds != 0.0 {
		name[1] = fmt.Sprintf("%s (Score: %0.2f)", name[1], math.Log10(n.Odds))
	}
	name[2] = fmt.Sprintf("Name: %s", nameString)
	ann, err := annotation.NewAnnotation(n.Annotation)
	if err != nil {
		return nil, err
	}
	name[3] = ann.Format()
	return name, nil
}

// NamesFromJSON creates gntagger's name structure from a finder output
func NamesFromJSON(path string) *Names {
	o := output.Output{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	o.FromJSON(b)
	return &Names{path, o}
}

// GetCurrentName returns currently selected name
func (n *Names) GetCurrentName() *output.Name {
	return &n.Data.Names[n.Data.Meta.CurrentName]
}

// UpdateAnnotations takes an annotation updates the current name with it.
// If needed, it propagates annotations further down the 'unseen' list.
func (n *Names) UpdateAnnotations(newAnnot annotation.Annotation, edge int,
	gnt *GnTagger) error {
	var (
		err      error
		oldAnnot annotation.Annotation
	)

	name := n.GetCurrentName()
	oldAnnot, err = annotation.NewAnnotation(name.Annotation)
	if err != nil {
		return err
	}

	name.Annotation = newAnnot.String()

	notAtEdge := n.Data.Meta.CurrentName < edge-3
	notAcceptedOrRejected := !newAnnot.In(annotation.NotName, annotation.Accepted)

	if notAtEdge || notAcceptedOrRejected {
		return nil
	}
	return tryToChangeNamesForward(n, oldAnnot, newAnnot, edge, gnt)
}

func tryToChangeNamesForward(n *Names, oldAnnot annotation.Annotation,
	newAnnot annotation.Annotation, edge int, gnt *GnTagger) error {
	for i := n.Data.Meta.CurrentName + 1; i < len(n.Data.Names); i++ {
		name := &n.Data.Names[i]
		sameName := n.GetCurrentName().Name == name.Name
		if !sameName {
			continue
		}

		if oldAnnot.In(annotation.NotAssigned, annotation.Doubtful) {
			name.Annotation = newAnnot.String()
		} else {
			nameAnnot, err := annotation.NewAnnotation(name.Annotation)
			if err != nil {
				return err
			}
			unmarkNames(name, nameAnnot, gnt)
		}
	}
	return nil
}

func unmarkNames(name *output.Name, a annotation.Annotation, gnt *GnTagger) {
	if IsDoubtful(name, gnt) {
		name.Annotation = annotation.Doubtful.String()
	} else {
		name.Annotation = annotation.NotAssigned.String()
	}
}

func IsDoubtful(n *output.Name, gnt *GnTagger) bool {
	return n.Odds != 0 && n.Odds < gnt.OddsHigh
}
