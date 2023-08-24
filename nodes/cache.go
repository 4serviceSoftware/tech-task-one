package nodes

import (
	"errors"
	"io"
	"os"
)

type Cache struct {
	repo Repository
}

const (
	cacheDir      = "./.cache"
	cacheFileName = "./.cache/nodescache"
)

func NewCache(r Repository) *Cache {
	return &Cache{repo: r}
}

func (c *Cache) Put() error {
	service := NewService(c.repo)
	err := os.MkdirAll(cacheDir, 0744)
	if err != nil {
		return err
	}
	f, err := os.Create(cacheFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	err = service.WriteJsonNodesTree(f, 0)
	if err != nil {
		os.Remove(cacheFileName)
		return err
	}
	return nil
}

func (c *Cache) Get(w io.Writer) error {
	f, err := os.Open(cacheFileName)
	if err != nil {
		return err
	}
	fi, err := os.Stat(cacheFileName)
	if err != nil {
		return err
	}
	if fi.Size() <= 0 {
		return errors.New("Nodes Cache: Empty cache file.")
	}
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	return nil
}
