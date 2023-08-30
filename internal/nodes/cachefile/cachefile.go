// This is not a regular cache implementation. It creates a file for caching data.
// It contains two methods which return file handler for writing and reading.
// I made this decision to save memory because not all content returning but
// only a handlers. Than this handlers can be used for reading and writing data directly.
package cachefile

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

// GetNewFileWriter creates new file and returns its handler for writing.
// Also returns function for closing this file and error if happend.
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

// GetFileReader opens a cache file and returns its handler for reading
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
