package nodes

import (
	"errors"
	"os"
	"path/filepath"
)

type CacheFile struct {
	fileName string
}

func NewCacheFile(fileName string) *CacheFile {
	return &CacheFile{fileName: fileName}
}

func (c *CacheFile) GetNewFileWriter() (*os.File, func(), error) {
	if _, err := os.Stat(filepath.Dir(c.fileName)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(c.fileName), 0744)
		if err != nil {
			return nil, nil, err
		}
	}
	f, err := os.Create(c.fileName)
	if err != nil {
		return nil, nil, err
	}
	return f, func() { f.Close() }, nil
}

func (c *CacheFile) GetFileReader() (*os.File, error) {
	f, err := os.Open(c.fileName)
	if err != nil {
		return nil, err
	}
	fi, err := os.Stat(c.fileName)
	if err != nil {
		return nil, err
	}
	if fi.Size() <= 0 {
		return nil, errors.New("Nodes Cache: Empty cache file.")
	}
	return f, nil
}
