# gntagger [![Doc Status][doc-img]][doc]

Finds names using gnfinder package internally and assists in interactive
curation of found scientific names

The project is in its early stages. It is already quite useful for tagging
`false positives`. We are working on adding fast methods for detecting
`false negatives`.

[![Ascii Cast][asciicast-img]][asciicast]

## Conversion of PDF to text

For modern linux environments you can use `less` command. For mac there is an
`xpdf` package:

```bash
brew install Caskroom/cask/xquartz
brew install xpdf
pdftotext doc.pdf doc.txt
```



[asciicast-img]: https://asciinema.org/a/wNfIt2TfZiyrAwJZKhuq5DkHV.png
[asciicast]: https://asciinema.org/a/wNfIt2TfZiyrAwJZKhuq5DkHV
[doc-img]: https://godoc.org/github.com/gnames/gntagger?status.png
[doc]: https://godoc.org/github.com/gnames/gntagger
