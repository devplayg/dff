package main

import (
	"fmt"
	"github.com/devplayg/dff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var fs *pflag.FlagSet
var version = "1.0.7"
var duplicateFileFinder *dff.DuplicateFileFinder

func main() {
	err := duplicateFileFinder.Start(time.Now())
	if err != nil {
		log.Error(err)
		return
	}
}

func init() {
	// Get flag set
	fs = pflag.NewFlagSet("dff", pflag.ContinueOnError)

	// Get arguments
	dirs := fs.StringArrayP("dir", "d", []string{}, "Target directories")
	minNumOfFilesInFileGroup := fs.IntP("min-count", "c", 2, "Minimum number of files in file group")
	minFileSize := fs.Int64P("min-size", "s", 1e6, "Minimum file size (Byte)")
	verbose := fs.BoolP("verbose", "v", false, "Verbose")
	sortBy := fs.StringP("sort", "r", "total", "Sort by [size | total | count]")
	format := fs.StringP("format", "f", "json", "Output format [json | text]")
	fs.Usage = printHelp
	_ = fs.Parse(os.Args[1:])

	// Check target directories
	if len(*dirs) < 1 {
		printHelp()
		os.Exit(0)
	}

	dff.InitLogger(*verbose)

	duplicateFileFinder = dff.NewDuplicateFileFinder(*dirs, *minNumOfFilesInFileGroup, *minFileSize, *sortBy, *format)
}

func printHelp() {
	fmt.Printf("Duplicate file finder v%s\n\n", version)
	fs.PrintDefaults()
}
