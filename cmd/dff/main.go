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
	t := time.Now()
	duplicateFileFinder := dff.NewDuplicateFileFinder(option)
	duplicateList, scanned, err := duplicateFileFinder.Find()
	if err != nil {
		log.Error(err)
		return
	}
	duplicateFileFinder.Display(duplicateList)

	// Logging
	log.WithFields(log.Fields{
		"number_of_files_scanned":               scanned,
		"duplicate_group_count":                 len(duplicateList),
		"minimum_number_of_files_in_file_group": duplicateFileFinder.Option.MinNumOfFilesInFileGroup,
		"min_file_size":                         duplicateFileFinder.Option.MinFileSize,
		"running_time(sec)":                     time.Since(t).Seconds(),
	}).Info("result")
}

func init() {

	// Set flags
	fs := pflag.NewFlagSet("dff", pflag.ContinueOnError)
	dirs := fs.StringArrayP("dir", "d", []string{}, "Target directories")
	minNumOfFilesInFileGroup := fs.IntP("min-count", "c", 2, "Minimum number of files in file group")
	minFileSize := fs.Int64P("min-size", "s", 0, "Minimum file size (Byte)")
	verbose := fs.BoolP("verbose", "v", false, "Verbose")
	sortBy := fs.StringP("sort", "r", "total", "Sort by [size | total | count]")
	format := fs.StringP("format", "f", "text", "Output format [json | text]")
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
