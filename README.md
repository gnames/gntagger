# gntagger [![Doc Status][doc-img]][doc]
Takes output from gnfinder and uses it to find false positives and false negatives in the result

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
