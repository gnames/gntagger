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
	case 1:
		path = flag.Arg(0)
		data, err = ioutil.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
	default:
		fmt.Printf("Please enter the path to your text file\n")
		os.Exit(1)
	}
	json := gnfinder.FindNamesJSON(data, &dict)
	jsonPath := path + ".json"
	if _, err = os.Stat(jsonPath); os.IsNotExist(err) {
		err = ioutil.WriteFile(jsonPath, json, 0644)
		if err != nil {
			log.Panic(err)
		}
	}
	InitGUI(path)
}
