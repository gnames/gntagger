package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gntagger/annotation"
)

// Names is an object that keeps output of a name finder and the path where to save this data on disk
type Names struct {
	// Path to json file with names
	Path string
	// Data is a gnfinder output
	Data gnfinder.Output
}

// NewNames uses a name finder or existing information to return Names structure generated from a text
func NewNames(text *Text, gnt *GnTagger) *Names {
	dict := gnfinder.LoadDictionary()

	opts := []gnfinder.Opt{gnfinder.WithBayesThreshold(gnt.OddsLow)}
	if gnt.Bayes {
		opts = append(opts, gnfinder.WithBayes)
	}

	data := gnfinder.FindNames(text.Processed, &dict, opts...)

	for i := range data.Names {
		n := &data.Names[i]
		if n.Odds != 0.0 && n.Odds < gnt.OddsHigh {
			n.Annotation = annotation.Doubtful.String()
		}
	}
	return &Names{Data: data, Path: text.FilePath(NamesFile)}
}

// Save writes current state of names to file
func (n *Names) Save() error {
	json := n.Data.ToJSON()
	return ioutil.WriteFile(n.Path, json, 0644)
}

// NameStrings composes text to show in terminal gui
func NameStrings(n *gnfinder.Name, current bool, i int, total int) ([]string, error) {
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
	o := gnfinder.Output{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	o.FromJSON(b)
	return &Names{path, o}
}

// GetCurrentName returns currently selected name
func (n *Names) GetCurrentName() *gnfinder.Name {
	return &n.Data.Names[n.Data.Meta.CurrentName]
}
