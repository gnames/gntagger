// Copyright © 2019 Dmitry Mozzherin <dmozzherin@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gnames/gntagger"
	"github.com/gnames/gntagger/termui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	build   string
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gntagger",
	Short: "finds scientific names in text and helps to annotates result",
	Long: `gntagger finds scientific names in text and provides
tools to curate and annotate found names.

To find names in a utf-8 encoded site:

gntagger ./your_file.txt

To find names from STDOUT:

cat ./your_file.txt | gntagger

If the input is a file, the results will be placed next to
the file under a directory with the name
your_file_gntagger.

If the input compes from a pipe the results will be located in
a ./gntagger_input directory.

Documentation: https://godoc.org/github.com/gnames/gntagger
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var err error
		var path string

		versionFlag(cmd)

		gnt := gntagger.NewGnTagger()

		switch len(args) {
		case 0:
			if ok := checkStdin(); ok {
				data, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					log.Panic(err)
				}
			} else {
				cmd.Help()
				os.Exit(0)
			}
		case 1:
			path = args[0]
			data, err = ioutil.ReadFile(path)
			if err != nil {
				log.Panic(err)
			}
		default:
			cmd.Help()
			os.Exit(2)
		}

		text := gntagger.NewText(data, path, version)
		gntagger.ShowWarningIfPreviousData(text)
		termui.InitGUI(text, gnt)
		defer infoOnExit(text)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string, b string) {
	version = v
	build = b
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gntagger.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "V", false, "show version and build timestamp.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gntagger" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gntagger")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		log.Panic(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func infoOnExit(t *gntagger.Text) {
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

func versionFlag(cmd *cobra.Command) {
	ver, err := cmd.Flags().GetBool("version")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if ver {
		fmt.Printf("\nversion: %s\n\nbuild:   %s\n\n", version, build)
		os.Exit(0)
	}
}
