package gntagger

import (
	"io/ioutil"
	"log"

	"github.com/gnames/gnfinder"
)

func NamesFromJSON() gnfinder.OutputJSON {
	o := gnfinder.OutputJSON{}
	b, err := ioutil.ReadFile("../../testdata/seashells_book.json")
	if err != nil {
		log.Panicln(err)
	}
	o.FromJSON(b)
	return o
}
