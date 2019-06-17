package dff

type DuplicateFileFinder struct {
	sortBy            int
	accessDeniedCount int
	format            int
	Option            *Option
}

func NewDuplicateFileFinder(option *Option) *DuplicateFileFinder {
	dff := DuplicateFileFinder{
		Option: option,
		sortBy: getSortValue(option.SortBy),
		format: getFormatValue(option.Format),
	}

	return &dff
}

func (d *DuplicateFileFinder) Find() ([]*UniqFile, int, error) {
	absDirs, err := isReadableDirs(d.Option.Dirs)
	if err != nil {
		return nil, 0, err
	}
	d.Option.Dirs = absDirs

	fileMap, err := collectFilesInDirs(d.Option.Dirs, d.Option.MinFileSize)
	if err != nil {
		return nil, 0, err
	}

	list, err := d.filterAndSort(fileMap)
	if err != nil {
		return nil, 0, err
	}
	return list, len(fileMap), nil
}

func (d *DuplicateFileFinder) filterAndSort(fileMap FileMap) ([]*UniqFile, error) {
	duplicateFileMap, err := findDuplicateFiles(fileMap, d.Option.MinNumOfFilesInFileGroup)
	if err != nil {
		return nil, err
	}

	for key, uniqFile := range duplicateFileMap {
		if len(uniqFile.List) < d.Option.MinNumOfFilesInFileGroup {
			delete(duplicateFileMap, key)
		}
	}
	list := getSortedValues(duplicateFileMap, d.sortBy)
	return list, nil
}

func (d *DuplicateFileFinder) Display(list []*UniqFile) {
	if len(list) < 1 {
		return
	}

	if d.format == TextFormat {
		outputFilesInTextFormat(list)
		return
	}

	outputFilesInJsonFormat(list)
}
