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
	flag.Usage = func() {
		help :=  "Usage of %s:\n\n%s file_with_names\n\nor\n\n" +
		"cat file | %s\n\n" +
		"If the input is a file, the results will be placed " +
		"next to the file under a directory with the same name " +
		"as the file.\n" +
		"If the input compes from a pipe the results will be located in " +
		"a ./gntagger_input directory.\n\n" +
		"Documentation: https://godoc.org/github.com/gnames/gntagger\n\n"

		fmt.Fprintf(os.Stderr, help, os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

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

	txt, names := prepareData(data, path)
	if err != nil {
		log.Panic(err)
	}
	InitGUI(txt, names)
}

func prepareData(text []byte, path string) (*Text, *Names) {
	dir, file := prepareFilepaths(path)
	cleanData := sanitizeText(text)
	textPath, jsonPath := createFilesGently(dir, file, cleanData)
	txt := PrepareText(textPath)
	names := NamesFromJSON(jsonPath)
	return txt, names
}

func sanitizeText(b []byte) []byte {
	re := regexp.MustCompile("[^[:print:]\n]")
	return []byte(re.ReplaceAllString(string(b), ""))
}

func prepareFilepaths(path string) (string, string) {
	file := "input.txt"
	if path == "" {
		return "./gntagger_input", file
	}
	d, f := filepath.Split(path)
	ext := filepath.Ext(f)
	dir := filepath.Join(d, f[0:len(f)-len(ext)])
	return dir, file
}

func createFilesGently(dir string, file string, text []byte) (string, string) {
	dict := gnfinder.LoadDictionary()
	textPath := filepath.Join(dir, file)
	jsonPath := textPath + ".json"

	ok, err := fileExistsAlready(dir, file)
	if err != nil {
		log.Panic(err)
	}
	if ok {
		return textPath, jsonPath
	}

	err = ioutil.WriteFile(textPath, text, 0644)
	if err != nil {
		log.Panic(err)
	}

	json := gnfinder.FindNamesJSON(text, &dict)
	err = ioutil.WriteFile(jsonPath, json, 0644)
	if err != nil {
		log.Panic(err)
	}
	return textPath, jsonPath
}

func fileExistsAlready(dir string, file string) (bool, error) {
	path := filepath.Join(dir, file)
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			return true, nil
		} else if err != nil {
			return false, err
		} else {
			return false, fmt.Errorf("Path %s exists but it is a dir", path)
		}
	} else if err != nil {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return false, err
		}
		return false, nil
	} else if !stat.IsDir() {
		return false, fmt.Errorf("Path %s exists but it is not a dir", dir)
	} else {
		return false, nil
	}
}
