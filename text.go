package gntagger

import (
	"path/filepath"
	"github.com/gnames/gnfinder"
	"log"
	"io/ioutil"
	"regexp"
	"os"
	"fmt"
)

type Text struct {
	// Path to the text file
	Path string
	// Content of the file converted to runes
	Original []rune
}

func (t *Text) Markup(n *Names) string {
	var markup []rune
	name := n.Data.Names[n.Data.Meta.CurrentName]
	markup = append(markup, t.Original[0:name.OffsetStart]...)
	markup = append(markup, []rune("\033[40;33;1m")...)
	markup = append(markup, t.Original[name.OffsetStart:name.OffsetEnd]...)
	markup = append(markup, []rune("\033[0m")...)
	markup = append(markup, t.Original[name.OffsetEnd:len(t.Original)]...)
	return string(markup)
}

func PrepareText(path string) *Text {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	return &Text{path, []rune(string(b))}
}

func prepareData(text []byte, path string) (*Text, *Names, error) {
	dir, file := prepareFilepaths(path)
	cleanData := sanitizeText(text)
	textPath, jsonPath := createFilesGently(dir, file, cleanData)
	txt := PrepareText(textPath)
	names := NamesFromJSON(jsonPath)
	return txt, names, nil
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
	dir := filepath.Join(d, f+"_input")
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
			return false, fmt.Errorf("path %s exists but it is a dir", path)
		}
	} else if err != nil {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return false, err
		}
		return false, nil
	} else if !stat.IsDir() {
		return false, fmt.Errorf("path %s exists but it is not a dir", dir)
	} else {
		return false, nil
	}
}
