package dff

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/minio/highwayhash"
	log "github.com/sirupsen/logrus"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func init() {

}

func getSortValue(val string) int {
	val = strings.TrimSpace(strings.ToLower(val))
	if val == "size" {
		return SortBySize
	}
	if val == "count" {
		return SortByCount
	}

	return SortByTotalSize
}

func getFormatValue(val string) int {
	val = strings.TrimSpace(strings.ToLower(val))
	if val == "text" {
		return TextFormat
	}
	return JsonFormat
}

func isReadableDirs(dirs []string) ([]string, error) {
	absDirs := make([]string, 0)
	for _, dir := range dirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}

		err = isValidDir(absDir)
		if err != nil {
			return nil, err
		}
		absDirs = append(absDirs, absDir)
	}
	return absDirs, nil
}

func isValidDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		return err
	}
	return nil
}

func generateHashKeyOfFile(path string, highwayhash hash.Hash) ([32]byte, error) {
	hash, err := getHighwayFileHash(highwayhash, path)
	if err != nil {
		return [32]byte{}, err
	}

	var key [32]byte
	copy(key[:], hash)

	return key, nil
}

func getHighwayFileHash(highwayHash hash.Hash, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	highwayHash.Reset()
	if _, err = io.Copy(highwayHash, file); err != nil {
		return nil, err
	}

	checksum := highwayHash.Sum(nil)
	return checksum, nil
}

func collectFilesInDirs(dirs []string, minFileSize int64) (FileMap, error) {
	ch := make(chan *FileMapDetail, len(dirs))
	for _, dir := range dirs {
		go searchDir(dir, minFileSize, ch)
	}

	fileMap := make(FileMap)
	for i := 0; i < len(dirs); i++ {
		filMapDetail := <-ch // Receive file map from goroutine
		for path, fileDetail := range filMapDetail.fileMap {
			fileMap[path] = fileDetail
		}
		log.Debugf("[%s] is merged into file map", filMapDetail.dir)
	}
	return fileMap, nil
}

func searchDir(dir string, minFileSize int64, ch chan *FileMapDetail) error {
	log.Infof("collecting files in [%s]", dir)
	fileMapDetail := NewFileMapDetail(dir)
	defer func() {
		ch <- fileMapDetail
	}()

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Error(err)
		}
		if err == nil && !f.IsDir() && f != nil && f.Mode().IsRegular() && f.Size() >= minFileSize {
			fileMapDetail.fileMap[path] = &FileDetail{filepath.Dir(path), f}
		}
		return nil
	})
	return err
}

func findDuplicateFiles(fileMap FileMap, minNumOfFilesInFileGroup int) (DuplicateFileMap, error) {
	fileMapBySize := classifyFilesBySize(fileMap)

	highwayHash, err := generateHighwayHashKey()
	if err != nil {
		return nil, err
	}

	duplicateFileMap := make(DuplicateFileMap)
	for _, list := range fileMapBySize {
		if len(list) < minNumOfFilesInFileGroup {
			continue
		}
		updateDuplicateFileMap(duplicateFileMap, list, highwayHash)
	}
	return duplicateFileMap, nil
}

func classifyFilesBySize(fileMap FileMap) FileMapBySize {
	fileMapBySize := make(FileMapBySize)
	for _, fileDetail := range fileMap {
		if _, ok := fileMapBySize[fileDetail.f.Size()]; !ok {
			fileMapBySize[fileDetail.f.Size()] = make([]*FileDetail, 0)
		}
		fileMapBySize[fileDetail.f.Size()] = append(fileMapBySize[fileDetail.f.Size()], fileDetail)
	}
	return fileMapBySize
}

func generateHighwayHashKey() (hash.Hash, error) {
	key := sha256.Sum256([]byte("Duplicate File Finder"))
	highwayhash, err := highwayhash.New(key[:])
	if err != nil {
		return nil, err
	}
	return highwayhash, err
}

func updateDuplicateFileMap(duplicateFileMap DuplicateFileMap, list []*FileDetail, highwayHash hash.Hash) {
	for _, fileDetail := range list {
		path := filepath.Join(fileDetail.dir, fileDetail.f.Name())
		key, err := generateHashKeyOfFile(path, highwayHash)
		if err != nil {
			log.Error(err)
			continue
		}

		if _, ok := duplicateFileMap[key]; !ok {
			duplicateFileMap[key] = NewDuplicateFiles(fileDetail.f.Size())
		}
		duplicateFileMap[key].List = append(duplicateFileMap[key].List, path)
		duplicateFileMap[key].TotalSize += fileDetail.f.Size()
		duplicateFileMap[key].Count++
	}
}

func outputFilesInTextFormat(list []*UniqFile) {
	for _, uniqFile := range list {
		fmt.Printf("total_size=%s, size=%d, count=%d, \n", ByteCountDecimal(uniqFile.TotalSize), uniqFile.Size, len(uniqFile.List))
		for _, path := range uniqFile.List {
			fmt.Printf("\t%s\n", path)
		}
	}
}

func outputFilesInJsonFormat(list []*UniqFile) {
	b, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(string(b))
}

func getSortedValues(duplicateFileMap DuplicateFileMap, sortBy int) []*UniqFile {
	list := make([]*UniqFile, 0, len(duplicateFileMap))
	for _, v := range duplicateFileMap {
		list = append(list, v)
	}

	if sortBy == SortByCount {
		sort.Sort(ByCount{list})
		return list
	}

	if sortBy == SortBySize {
		sort.Sort(BySize{list})
		return list
	}

	sort.Sort(ByTotalSize{list})
	return list
}

func InitLogger(verbose bool) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
