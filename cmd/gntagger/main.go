package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gnames/gnfinder"
	. "github.com/gnames/gntagger"
)

func main() {
	var path string
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
		fmt.Printf("Please enter the path to your text file\n")
		os.Exit(1)
	}
	txt, names := prepareData(data)
	InitGUI(txt, names)
}

func prepareData(data []byte) (*Text, *Names) {
	dict := gnfinder.LoadDictionary()
	cleanData := sanitizeText(data)
	dir, err := ioutil.TempDir("", "gntagger")
	if err != nil {
		log.Panic(err)
	}
	textPath := filepath.Join(dir, "text.txt")
	jsonPath := textPath + ".json"
	err = ioutil.WriteFile(textPath, cleanData, 0644)
	if err != nil {
		log.Panic(err)
	}

	json := gnfinder.FindNamesJSON(cleanData, &dict)
	err = ioutil.WriteFile(jsonPath, json, 0644)
	if err != nil {
		log.Panic(err)
	}
	txt := PrepareText(textPath)
	names := NamesFromJSON(jsonPath)
	return txt, names
}

func sanitizeText(b []byte) []byte {
	re := regexp.MustCompile("[^[:print:]\n]")
	return []byte(re.ReplaceAllString(string(b), ""))
}
