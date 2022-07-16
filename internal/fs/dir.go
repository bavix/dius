package fs

import (
	"os"
	"path/filepath"
	"sync/atomic"
)

type FileInterface interface {
	Size() (uint64, error)
	Name() string
	IsDir() bool
}

type file struct {
	i    os.FileInfo
	path string
	size *uint64
	err  error
}

func New(path string, info os.FileInfo) FileInterface {
	return &file{path: path, i: info}
}

func (f *file) Name() string {
	return f.i.Name()
}

func (f *file) IsDir() bool {
	return f.i.IsDir()
}

func (f *file) Size() (uint64, error) {
	if f.size == nil {
		size, err := dirSize(f.path, f.i)
		f.size = &size
		f.err = err
	}

	return *f.size, f.err
}

func dirSize(path string, info os.FileInfo) (uint64, error) {
	if !info.IsDir() {
		return uint64(info.Size()), nil
	}

	var size uint64 = 0
	err := filepath.Walk(filepath.Join(path, info.Name()),
		func(path string, info os.FileInfo, err error) error {
			if err == nil {
				atomic.AddUint64(&size, uint64(info.Size()))
			}

			return err
		})

	if err != nil {
		return 0, err
	}

	return size, nil
}
