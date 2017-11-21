package gntagger

import (
	"crypto/sha1"
	"fmt"
	"log"
	"path/filepath"
	"regexp"

	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/go-wordwrap"
)

type FileType int

const (
	InputFile FileType = iota
	NamesFile
	MetaFile
)

type TextMeta struct {
	// Checksum is a hash calculated from the content
	Checksum string `json:"text_checksum"`
	// Githash is a hash of a git commit. It would be 'n/a' if not given
	Githash string `json:"app_git_commit"`
	// Timestamp of the last save
	Timestamp string `json:"save_timestamp"`
}

func (t *TextMeta) ToJSON() []byte {
	json, err := jsoniter.Marshal(t)
	if err != nil {
		log.Panic(err)
	}
	return json
}

type Text struct {
	// Raw text, as it was given by a user
	Raw []byte
	// Cleaned text after removing non-printable characters and wrapping
	// according to the width of a user's terminal
	Processed []rune
	// Path to the text file
	Path string
	// Files is a map that contains names of the files created by gntagger
	Files map[FileType]string
	// TextMeta describes provides metainformation about text:
	// Checksum, Githash, Timestamp
	TextMeta
	// width of the user's terminal window
	width uint
}

func NewText(data []byte, path string, githash string) *Text {
	path = preparePath(path)
	checksum := fmt.Sprintf("%x", sha1.Sum(data))
	meta := TextMeta{
		Checksum:  checksum,
		Githash:   githash,
		Timestamp: timestamp(),
	}

	text := &Text{
		Raw:      data,
		Path:     path,
		TextMeta: meta,
		Files: map[FileType]string{
			InputFile: "input.txt",
			NamesFile: "names.json",
			MetaFile:  "meta.json",
		},
	}
	return text
}

func (t *Text) Process(width int) {
	cleanData := sanitizeText(t.Raw)
	alignedText := wordwrap.WrapString(cleanData, uint(width))
	t.Processed = []rune(alignedText)
}

func (t *Text) FilePath(f FileType) string {
	return filepath.Join(t.Path, t.Files[f])
}

func sanitizeText(b []byte) string {
	re := regexp.MustCompile("[^[:print:]\n]")
	return re.ReplaceAllString(string(b), "")
}

func preparePath(path string) string {
	if path == "" {
		return "./gntagger_input"
	}
	d, f := filepath.Split(path)
	dir := filepath.Join(d, f+"_gntagger")
	return dir
}
