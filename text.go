package gntagger

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"crypto/sha1"
	"encoding/base64"

	"github.com/gnames/gnfinder"
	"github.com/mitchellh/go-wordwrap"
)

const shaFileName = "sha.txt"

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

func computeSha(input []byte) (string, error) {
	hasher := sha1.New()
	if _, err := hasher.Write(input); err != nil {
		return "", err
	}
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

func checkIfShaEquals(dir string) (bool, error) {
	hasher := sha1.New()
	hasher.Write(textData)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	shaExisting, err := ioutil.ReadFile(filepath.Join(dir, shaFileName))
	if err != nil {
		return false, err
	}

	return string(shaExisting) == sha, nil
}

func PrepareText(path string) *Text {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	return &Text{path, []rune(string(b))}
}

func prepareData(text []byte, dir string, file string, width int,
	bayes *bool) (*Text, *Names, error) {
	cleanData := sanitizeText(text)
	alignedText := wordwrap.WrapString(string(cleanData), uint(width))
	textPath, jsonPath, err :=
		createFilesGently(dir, file, []byte(alignedText), bayes)
	if err != nil {
		return nil, nil, err
	}
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

func jsonFileName(fileName string) string {
	return fileName + ".json"
}

func createFilesGently(dir string, file string, text []byte,
	bayes *bool) (string, string, error) {
	dict := gnfinder.LoadDictionary()
	textPath := filepath.Join(dir, file)
	jsonPath := jsonFileName(textPath)

	err := ioutil.WriteFile(textPath, text, 0644)
	if err != nil {
		log.Panic(err)
	}

	var opts []gnfinder.Opt
	if *bayes {
		opts = []gnfinder.Opt{gnfinder.WithBayes}
	}

	json := gnfinder.FindNamesJSON(text, &dict, opts...)
	err = ioutil.WriteFile(jsonPath, json, 0644)
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	sha, err := computeSha(textData)
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	err = ioutil.WriteFile(filepath.Join(dir, shaFileName),
		[]byte(sha), 0644)
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	return textPath, jsonPath, nil
}

func checkIfPathExists(dir string, file string) bool {
	var (
		err    error
		exists bool
	)
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	} else if err != nil {
		log.Panic(err)
		os.Exit(2)
	} else if !stat.IsDir() {
		log.Panicf("Path %s exists. It should be a dir", dir)
		os.Exit(2)
	}

	dataFilePath := filepath.Join(dir, file)
	stat, err = os.Stat(dataFilePath)
	if os.IsNotExist(err) {
		exists = false
	} else if err != nil {
		log.Panic(err)
		os.Exit(2)
	} else if stat.IsDir() {
		log.Panicf("Path %s exists. It should not be a dir", dir)
		os.Exit(2)
	} else {
		exists = true
	}
	return exists
}
