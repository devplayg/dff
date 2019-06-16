package dff

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type DuplicateFileFinder struct {
	sortBy            int
	accessDeniedCount int
	format            int
	option            *Option
}

func NewDuplicateFileFinder(option *Option) *DuplicateFileFinder {
	dff := DuplicateFileFinder{
		option: option,
		sortBy: getSortValue(option.SortBy),
		format: getFormatValue(option.Format),
	}

	return &dff
}

func (d *DuplicateFileFinder) Start(t time.Time) error {
	absDirs, err := isReadableDirs(d.option.Dirs)
	if err != nil {
		return err
	}
	d.option.Dirs = absDirs

	fileMap, err := collectFilesInDirs(d.option.Dirs, d.option.MinFileSize)
	if err != nil {
		return err
	}

	duplicateFileMap, err := findDuplicateFiles(fileMap, d.option.MinNumOfFilesInFileGroup)
	if err != nil {
		return err
	}

	duplicateFileGroupCount := d.displayDuplicateFileGroups(duplicateFileMap)

	log.WithFields(log.Fields{
		"number_of_files_scanned":               len(fileMap),
		"duplicate_group_count":                 duplicateFileGroupCount,
		"minimum_number_of_files_in_file_group": d.option.MinNumOfFilesInFileGroup,
		"min_file_size":                         d.option.MinFileSize,
		"running_time(sec)":                     time.Since(t).Seconds(),
	}).Info("result")
	return nil
}

func (d *DuplicateFileFinder) displayDuplicateFileGroups(duplicateFileMap DuplicateFileMap) int {
	for key, uniqFile := range duplicateFileMap {
		if len(uniqFile.List) < d.option.MinNumOfFilesInFileGroup {
			delete(duplicateFileMap, key)
		}
	}
	list := getSortedValues(duplicateFileMap, d.sortBy)
	if len(list) < 1 {
		return 0
	}

	if d.format == TextFormat {
		outputFilesInTextFormat(list)
		return len(list)
	}
	outputFilesInJsonFormat(list)
	return len(list)
}
