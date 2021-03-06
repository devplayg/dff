# Duplicate File Finder

[![Build Status](https://travis-ci.com/devplayg/dff.svg?branch=master)](https://travis-ci.com/devplayg/dff)
[![Go Report Card](https://goreportcard.com/badge/github.com/devplayg/dff)](https://goreportcard.com/report/github.com/devplayg/dff)
[![Release](https://img.shields.io/github/release/devplayg/dff.svg)](https://github.com/devplayg/dff/releases)

finds duplicate files in directories

*Powered by [HighwayHash](https://github.com/google/highwayhash) hash algorithm*

    Duplicate file finder
    
      -d, --dir stringArray   Target directories
      -f, --format string     Output format [json | text] (default "text")
      -c, --min-count int     Minimum number of files in file group (default 2)
      -s, --min-size int      Minimum file size (Byte) (default 0)
      -r, --sort string       Sort by [size | total | count] (default "total")
      -v, --verbose           Verbose


### Find duplicate files in a specific directory 

    dff -d /dir
    
### Find duplicate files in specific directories

    dff -d /dir1 -d /dir2 -d /dir3
    
### Find duplicate files if there are 10 or more identical files (Default: 2)

    dff -d /dir -c 10
    
### Find duplicate files of 2 MB or more (Default: 1 MB)

    dff -d /dir -s 2000000 
    
### Output format

JSON (Default)
    
    dff -d /dir -f json

Text
    
    dff -d /dir -f text

### Sort

Sort by file size sum

    dff -d /dir -r total

Sort by file size
 
    dff -d /dir -r size
    
Sort by duplicate count    
    
    dff -d /dir -r count
    
### Example

```go
option = &dff.Option{
    Dirs:                     []string{"/path/to/dir1", "/path/to/dir2"},
    MinNumOfFilesInFileGroup: 3,
    MinFileSize:              10000000,
    SortBy:                   "total",
    Format:                   "json",
}
duplicateFileFinder := dff.NewDuplicateFileFinder(option)
duplicateList, scanned, err := duplicateFileFinder.Find()
if err != nil {
    log.Error(err)
    return
}
```