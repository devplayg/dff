package main

import (
	"fmt"
	"github.com/devplayg/dff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var version = "1.0.8"
var option *dff.Option

func main() {
	duplicateFileFinder := dff.NewDuplicateFileFinder(option)
	err := duplicateFileFinder.Start(time.Now())
	if err != nil {
		log.Error(err)
		return
	}
}

func init() {

	// Set flags
	fs := pflag.NewFlagSet("dff", pflag.ContinueOnError)
	dirs := fs.StringArrayP("dir", "d", []string{}, "Target directories")
	minNumOfFilesInFileGroup := fs.IntP("min-count", "c", 2, "Minimum number of files in file group")
	minFileSize := fs.Int64P("min-size", "s", 1e6, "Minimum file size (Byte)")
	verbose := fs.BoolP("verbose", "v", false, "Verbose")
	sortBy := fs.StringP("sort", "r", "total", "Sort by [size | total | count]")
	format := fs.StringP("format", "f", "json", "Output format [json | text]")
	fs.Usage = func() {
		fmt.Printf("Duplicate file finder v%s\n\n", version)
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[1:])

	// Check target directories
	if len(*dirs) < 1 {
		fs.Usage()
		os.Exit(0)
	}

	// Initialize Logger
	dff.InitLogger(*verbose)

	// Set options
	option = &dff.Option{
		Dirs:                     *dirs,
		MinNumOfFilesInFileGroup: *minNumOfFilesInFileGroup,
		MinFileSize:              *minFileSize,
		SortBy:                   *sortBy,
		Format:                   *format,
	}
}
