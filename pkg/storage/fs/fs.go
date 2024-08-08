package fs

import (
	"fmt"
	"os"
	"path"
)

type FSStorage struct {
	Dir string
}

type FSOption func(c *FSStorage)

func WithDir(dir string) FSOption {
	return func(fs *FSStorage) {
		fs.Dir = dir
	}
}

func NewFSStorage(opts ...FSOption) *FSStorage {
	s := &FSStorage{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *FSStorage) Write(b []byte, filename string) error {
	filePath := path.Join(s.Dir, filename)
	if err := os.WriteFile(filePath, b, os.ModePerm); err != nil {
		return fmt.Errorf("cannot write to %s: %w", filePath, err)
	}
	return nil
}
