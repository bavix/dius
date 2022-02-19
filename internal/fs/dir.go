package fs

import (
	"github.com/bavix/dius/internal/wgi"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
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
		size, err := fastSize(f.path, f.i)
		f.size = &size
		f.err = err
	}

	return *f.size, f.err
}

func fastSize(path string, info os.FileInfo) (uint64, error) {
	if !info.IsDir() {
		return uint64(info.Size()), nil
	}

	newPath := filepath.Join(path, info.Name())
	files, err := ioutil.ReadDir(newPath)
	for attempts := 0; attempts < 20 && err != nil; attempts++ {
		time.Sleep(time.Millisecond * 20)
		files, err = ioutil.ReadDir(newPath)
	}

	if err != nil {
		return 0, err
	}

	var errSize error = nil
	var size uint64 = 0
	var wg wgi.WaitGroup
	for _, file := range files {
		if !file.IsDir() {
			size += uint64(file.Size())
			continue
		}

		wg.Add(1)
		go func(f os.FileInfo) {
			defer wg.Done()

			sum, err := fastSize(newPath, f)
			if err != nil {
				errSize = err
			}

			size += sum
		}(file)
	}

	wg.Wait()

	return size, errSize
}
