package gntagger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// ShowWarningIfPreviousData takes *Text pointer. It warns if previously created
// data exist and backups the old data.
func ShowWarningIfPreviousData(text *Text) {
	var warning string
	if old := previousDataChecksums(text); old.Checksum != "" {
		if text.Checksum != old.Checksum {
			warning = "\nYour input file has changed."
		} else if text.GNtaggerVersion != old.GNtaggerVersion {
			warning = "\nYour gntagger is updated."
		}
		if warning != "" {
			warning += " Older data are backed up for comparison."
			moveOldFiles(text, old.Timestamp)
			showWarning(warning)
		}
	}
}

func moveOldFiles(text *Text, timestamp string) {
	for ft, f := range text.Files {
		newPath := filepath.Join(text.Path, timestamp+"_"+f)
		err := os.Rename(text.FilePath(ft), newPath)
		if err != nil {
			log.Panic(err)
		}
	}
}

func showWarning(warning string) {
	warning += "\n\nPress \"Enter\" to continue."
	fmt.Println(warning)
	fmt.Scanln()
}

func previousDataChecksums(text *Text) *TextMeta {
	var old TextMeta
	for f := range text.Files {
		_, err := os.Stat(text.FilePath(f))
		if err != nil && os.IsNotExist(err) {
			return &TextMeta{}
		}
	}

	meta, err := ioutil.ReadFile(text.FilePath(MetaFile))
	if err != nil {
		log.Panic(err)
	}
	r := bytes.NewReader(meta)
	d := jsoniter.NewDecoder(r)
	err = d.Decode(&old)
	if err != nil {
		log.Panic(err)
	}
	return &old
}

func timestamp() string {
	t := time.Now()
	return t.Format("20060102_150405")
}

func fileExistsAlready(text *Text) (bool, error) {
	filePath := text.FilePath(NamesFile)
	if stat, err := os.Stat(text.Path); err == nil && stat.IsDir() {
		if stat, err := os.Stat(filePath); err == nil && !stat.IsDir() {
			return true, nil
		} else if err != nil {
			return false, nil
		} else {
			return false, fmt.Errorf("Path %s exists but it is a dir", filePath)
		}
	} else if err != nil {
		err := os.Mkdir(text.Path, 0755)
		return false, err
	} else if !stat.IsDir() {
		return false, fmt.Errorf("Path %s exists but it is not a dir", text.Path)
	}

	return false, nil
}

func createFilesGently(t *Text, names *Names) {
	err := ioutil.WriteFile(t.FilePath(InputFile),
		t.ProcessedBytes, 0644)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(names.Path, names.Data.ToJSON(), 0644)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(t.FilePath(MetaFile), t.TextMeta.ToJSON(), 0644)
	if err != nil {
		log.Panic(err)
	}
	processedTextFromFile(t)
}

func processedTextFromFile(t *Text) {
	// by some reason the text does not show correctly
	// in GUI unless it is taken from its file
	txt, err := ioutil.ReadFile(t.FilePath(InputFile))
	if err != nil {
		log.Panic(err)
	}
	t.Processed = []rune(string(txt))
}
