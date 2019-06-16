package dff

import (
	"os"
)

const (
	SortBySize = iota + 1
	SortByTotalSize
	SortByCount
)

const (
	JsonFormat = iota + 1
	TextFormat
)

type FileDetail struct {
	dir string
	f   os.FileInfo
}

type FileMap map[string]*FileDetail

type FileMapDetail struct {
	fileMap FileMap
	dir     string
}

func NewFileMapDetail(dir string) *FileMapDetail {
	return &FileMapDetail{
		dir:     dir,
		fileMap: make(FileMap),
	}
}

type FileMapBySize map[int64][]*FileDetail

type DuplicateFileMap map[[32]byte]*UniqFile

type UniqFile struct {
	List      []string
	Size      int64
	TotalSize int64
	Count     int
}

func NewDuplicateFiles(size int64) *UniqFile {
	return &UniqFile{
		Size: size,
		List: make([]string, 0),
	}
}

// Sorting for UniqFile
type UniqFiles []*UniqFile

func (s UniqFiles) Len() int      { return len(s) }
func (s UniqFiles) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort by size
type BySize struct{ UniqFiles }

func (s BySize) Less(i, j int) bool { return s.UniqFiles[i].Size > s.UniqFiles[j].Size }

// Sort by total size
type ByTotalSize struct{ UniqFiles }

func (s ByTotalSize) Less(i, j int) bool { return s.UniqFiles[i].TotalSize > s.UniqFiles[j].TotalSize }

// Sort by
type ByCount struct{ UniqFiles }

func (s ByCount) Less(i, j int) bool { return s.UniqFiles[i].Count > s.UniqFiles[j].Count }
