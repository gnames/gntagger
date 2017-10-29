package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gnames/gnfinder"
	. "github.com/gnames/gntagger"
)

func main() {
	var path string
	dict := gnfinder.LoadDictionary()
	flag.Parse()
	var data []byte
	var err error
	switch flag.NArg() {
	case 0:
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Panic(err)
		}
	case 1:
		path = flag.Arg(0)
		data, err = ioutil.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}
	json := gnfinder.FindNamesJSON(data, &dict)
	err = ioutil.WriteFile(path+".json", json, 0644)
	if err != nil {
		log.Panic(err)
	}
	InitGUI(path)
}
