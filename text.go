package gntagger

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"unicode"

	"runtime"

	jsoniter "github.com/json-iterator/go"
)

// FileType describes types created by gntagger during name finding and
// curation.
type FileType int

const (
	// InputFile contains text used for name-finding.
	InputFile FileType = iota
	// NamesFile contains JSON output with names and metadata.
	NamesFile
	// MetaFile creates meta-information used for various purposes.
	MetaFile
)

// TextMeta used for creating the content of the MetaFile
type TextMeta struct {
	// Checksum is a hash calculated from the content.
	Checksum string `json:"text_checksum"`
	// GNtaggerVersion is the version of gntagger that generated output.
	GNtaggerVersion string `json:"gntagger_version"`
	// Timestamp of the last save.
	Timestamp string `json:"save_timestamp"`
}

// ToJSON converts meta-information into JSON format
func (t *TextMeta) ToJSON() []byte {
	json, err := jsoniter.Marshal(t)
	if err != nil {
		log.Panic(err)
	}
	return json
}

// Text contains text of the input and its metadata
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
	// Checksum, GNtaggerVersion, Timestamp
	TextMeta
	// errors that accumulated during the process. They will be shown on exit.
	errors map[string]error
}

// NewText creates new Text object
func NewText(data []byte, path string, GNtaggerVersion string) *Text {
	path = preparePath(path)
	checksum := fmt.Sprintf("%x", sha1.Sum(data))
	meta := TextMeta{
		Checksum:        checksum,
		Timestamp:       timestamp(),
		GNtaggerVersion: GNtaggerVersion,
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

// Errors returns list of errors that happened during execution of the
// gntagger.
func (t *Text) Errors() []error {
	res := make([]error, 0, len(t.errors))
	for _, v := range t.errors {
		res = append(res, v)
	}
	return res
}

// AddError adds a new arror to the Text's error collection.
func (t *Text) AddError(err error) {
	if t.errors == nil {
		t.errors = make(map[string]error)
	}
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(err.Error())))
	if _, ok := t.errors[hash]; !ok {
		t.errors[hash] = err
	}
}

// Process removes all non-printable characters from Text and wraps
// its lines making sure that all scientific names are visible.
func (t *Text) Process(width int) {
	if runtime.GOOS == "windows" {
		width = 100
	}
	printable := printableBytes(t.Raw)
	t.ProcessedBytes = wrap(printable, width)
	t.Processed = []rune(string(t.ProcessedBytes))
}

// FilePath returns a file path for a given FileType.
func (t *Text) FilePath(f FileType) string {
	return filepath.Join(t.Path, t.Files[f])
}

// PrepareFilesAndText creates files if needed and returns list of names
// for curation.
func PrepareFilesAndText(t *Text, w int, gnt *GnTagger) *Names {
	exist, err := fileExistsAlready(t)
	if err != nil {
		log.Panic(err)
	}

	if exist {
		processedTextFromFile(t)
		return NamesFromJSON(t.FilePath(NamesFile))
	}
	t.Process(w)
	names := NewNames(t, gnt)
	createFilesGently(t, names)
	return names
}

func printableBytes(b []byte) []byte {
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

		if (unicode.IsPrint(c) || unicode.IsSpace(c)) && c != rune('\r') && c != rune('\v') {
			target.WriteRune(c)
		}
	}
	return target.Bytes()
}

func wrap(b []byte, width int) []byte {
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
