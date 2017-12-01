package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/gnames/gntagger"
)

var (
	githash    = "n/a"
	gittag     = "n/a"
	buildstamp = "n/a"
)

func main() {
	flag.Usage = func() {
		help :=
			"Usage of %s:\n\n" +
				"%s file_with_names\n\nor\n\n" +
				"cat file | %s\n\n" +
				"If the input is a file, the results will be placed " +
				"next to the file under a directory with the name " +
				"\"your_file_name_gntagger\".\n" +
				"If the input compes from a pipe the results will be located in " +
				"a ./gntagger_input directory.\n\n" +
				"Documentation: https://godoc.org/github.com/gnames/gntagger\n\n"

		fmt.Fprintf(os.Stderr, help, os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var (
		path string
		data []byte
		err  error
	)

	versionFlag := flag.Bool("version", false, "Print version")
	bayesFlag := flag.Bool("bayes", false, "Use bayes name-finding")
	flag.Parse()

	if *versionFlag {
		fmt.Printf(" Git commit hash: %s\n Git tag: %s\n UTC Build Time: %s\n\n",
			githash, gittag, buildstamp)
		os.Exit(0)
	}

	switch flag.NArg() {
	case 0:
		if ok := checkStdin(); ok {
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Panic(err)
			}
		} else {
			flag.Usage()
			os.Exit(0)
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

	text := NewText(data, path, githash)
	WarningIfPreviousData(text)
	InitGUI(text, bayesFlag)
	defer infoOnExit(text)
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		log.Panic(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func infoOnExit(t *Text) {
	path, err := filepath.Abs(t.Path)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("\n Files with curated names and the text are located in"+
		"\n\n %s \n\n",
		path,
	)
	for _, e := range t.Errors() {
		fmt.Printf("\n %s \n\n", e.Error())
	}
}
