package gntagger

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gntagger/annotation"
)

type Names struct {
	// Path to json file with names
	Path string
	// Data is a gnfinder output
	Data gnfinder.Output
}

func NewNames(text *Text, bayes *bool) *Names {
	dict := gnfinder.LoadDictionary()

	opts := []gnfinder.Opt{gnfinder.WithBayesThreshold(1)}
	if *bayes {
		opts = append(opts, gnfinder.WithBayes)
	}

	data := gnfinder.FindNames(text.Processed, &dict, opts...)

	for i, _ := range data.Names {
		n := &data.Names[i]
		if n.Odds < 100 {
			n.Annotation = annotation.Doubtful.String()
		}
	}
	return &Names{Data: data, Path: text.FilePath(NamesFile)}
}

func (n *Names) Save() error {
	json := n.Data.ToJSON()
	return ioutil.WriteFile(n.Path, json, 0644)
}

func nameStrings(n *gnfinder.Name, current bool, i int, total int) ([]string, error) {
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
	annotation, err := annotation.NewAnnotation(n.Annotation)
	if err != nil {
		return nil, err
	}
	name[3] = annotation.Format()
	return name, nil
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

func (names *Names) GetCurrentName() *gnfinder.Name {
	return &names.Data.Names[names.Data.Meta.CurrentName]
}
