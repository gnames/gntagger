# gntagger [![Doc Status][doc-img]][doc]

gntagger not only finds scientific names in a document. It also allows user to
go through each found name, see it in a context of a text, and then accept or
reject the found name.

We made this program so we can improve on quality of name-finding algorithm,
but it is useful for anybody who needs to extract scientific names from a book
or a scientific paper. The program works for MS Windows, Mac and Linux and it
runs from a command line interface -- CMD in case of windows, or a terminal
under Mac and Linux.

gntagger allows to curate 4000 names spread over 600 pages in about 2 hours. It
is significantly faster than curaion made in a text editor or pdf viewer.

[![Ascii Cast][asciicast-img]][asciicast]

## Installation

The  program is just an executable file that runs from a command line. Download
latest zip or tar file for your operating system, extract the file and place it
somehere in your `PATH`, so it is visible by your system.

gntagger works with plain texts, so if you need to find names in a PDF file,
first you need to convert it to text.

## Conversion of PDF to text

### Linux

You can just use `less` command.

```
less paper1.pdf | gntagger
```

Another option is pdftotext from xpdf package.

### Mac

Use `xpdf` package:

```bash
brew install Caskroom/cask/xquartz
brew install xpdf
pdftotext -layout doc.pdf doc.txt
```

### Windows

Download [Xpdf tools][xpdf-tools], unzip them, and use pdftotext.exe

```
pdftotext.exe -layout doc.pdf doc.txt
```

## Usage

To find out version

```bash
gntagger -version
```

To get names from a file (processed text and list of names will be saved in the
same directory as the text file)

```bash
gntagger file_with_names.txt

#using bayes algorithm

gntagger -bayes file_with_names.txt

# on windows
gntagger.exe  file_with_names.txt
```

To get names from stanard input

```
# linux

less file.pdf | gntagger
less file.pdf | gntagger -bayes

# mac

pdftotext -layout file.pdf | gntagger
pdftotext -layout file.pdf | gntagger -bayes
```

Note that -layout flag for pdftotext tries to preserve original structure of
the text, as it was in the original PDF. It significantly increases chances for
finding names that are split between the end and the start of two lines.

## User Interface

The user interface of the program consists of 2 panels. The left panel
contains detected scientific names, with a "current name" located in the middle
of the screen and highlighted. The left panel contains the full text, where
the "current name" is highlighted and aligned with the "current name" in the
left panel.

The program is designed to move though the names quickly. Navigate to the
next/previous name in the left panel using Right/Left arrow keys. All names
have an empty annotation at the beginning. Pressing Right Arrow key
automatically "accepts" found name if the annotation is empty. Other keys
allow to annotate the "current name" differently:

* Space: rejects a name with "NotName" annotation

* 'y':   re-accepts mistakenly rejected name with "Accepted" annotation

* 'u':   marks a name as "Uninomial"

* 'g':   marks a name as "Genus"

* 's':  marks a name as "Species"

* Ctrl-C: saves curation and exits application

* Ctrl-S: saves curations made so far

**Current names are saved to clipboard automatically**, so it is easy to paste
them into a browser, speadsheet, database, or text editor.

The program autosaves results of curation. If the program crashes, or exited
the user can continue curation at the last point instead of starting from
scratch.

## Development

### Running tests

Install ginkgo, a BDD testing framefork for Go.

```bash
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```

To run tests go to root directory of the project and run

```bash
ginkgo

#or

go test
```


### Build executable

```bash
go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` \
                   -X main.githash=`git rev-parse HEAD | cut -c1-7` \
                   -X main.gittag=`git describe --tags`" \
         -o gntagger -a cmd/gntagger/main.go
```


[asciicast-img]: https://asciinema.org/a/wNfIt2TfZiyrAwJZKhuq5DkHV.png
[asciicast]: https://asciinema.org/a/wNfIt2TfZiyrAwJZKhuq5DkHV
[doc-img]: https://godoc.org/github.com/gnames/gntagger?status.png
[doc]: https://godoc.org/github.com/gnames/gntagger
[xpdf-tools]: https://www.xpdfreader.com/download.html
