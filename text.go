package gntagger

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"unicode"

	jsoniter "github.com/json-iterator/go"
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
	// Cleaned text in bytes
	ProcessedBytes []byte
	// Path to the text file
	Path string
	// Files is a map that contains names of the files created by gntagger
	Files map[FileType]string
	// TextMeta describes provides metainformation about text:
	// Checksum, Githash, Timestamp
	TextMeta
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
	printable := PrintableBytes(t.Raw)
	t.ProcessedBytes = Wrap(printable, width)
	t.Processed = []rune(string(t.ProcessedBytes))
}

func (t *Text) FilePath(f FileType) string {
	return filepath.Join(t.Path, t.Files[f])
}

func PrintableBytes(b []byte) []byte {
	init := make([]byte, 0, len(b))
	target := bytes.NewBuffer(init)
	source := bytes.NewReader(b)

	for {
		c, _, err := source.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}

		if unicode.IsPrint(c) || unicode.IsSpace(c) {
			target.WriteRune(c)
		}
	}
	return target.Bytes()
}

func Wrap(b []byte, width int) []byte {
	if width <= 0 {
		width = 5
	}
	var word bytes.Buffer
	var endReached bool
	init := make([]byte, 0, len(b))
	target := bytes.NewBuffer(init)
	source := bytes.NewReader(b)
	lineCursor := 0
	wordCursor := 0

	for {
		if endReached {
			break
		}
		c, _, err := source.ReadRune()
		if err == io.EOF {
			endReached = true
		} else if err != nil {
			log.Panic(err)
		}

		switch {
		case unicode.IsSpace(c) || endReached:
			if wordCursor == 0 {
				lineCursor++
			} else if lineCursor+wordCursor >= width {
				target.WriteRune(rune('\n'))
				lineCursor = 0
			}

			if c == rune('\n') {
				lineCursor = 0
			} else {
				lineCursor += wordCursor + 1
			}

			_, err := word.WriteTo(target)
			if err != nil {
				log.Panic(err)
			}

			target.WriteRune(c)
			wordCursor = 0
		default:
			word.WriteRune(c)
			wordCursor++
		}
	}
	return target.Bytes()
}

func preparePath(path string) string {
	if path == "" {
		return "./gntagger_input"
	}
	d, f := filepath.Split(path)
	dir := filepath.Join(d, f+"_gntagger")
	return dir
}

func prepareFilesAndText(t *Text, w int, bayes *bool) *Names {
	exist, err := fileExistsAlready(t)
	if err != nil {
		log.Panic(err)
	}

	if exist {
		processedTextFromFile(t)
		return NamesFromJSON(t.FilePath(NamesFile))
	} else {
		t.Process(w)
		names = NewNames(t, bayes)
		createFilesGently(t, names)
		return names
	}
}
