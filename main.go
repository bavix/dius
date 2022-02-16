package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func fastSize(path string, info os.FileInfo) (uint64, error) {
	if !info.IsDir() {
		return uint64(info.Size()), nil
	}

	newPath := path + "/" + info.Name()
	files, err := ioutil.ReadDir(newPath)
	if err != nil {
		return 0, err
	}

	var size uint64 = 0
	var errSize error = nil
	var wg sync.WaitGroup
	for _, file := range files {
		if !file.IsDir() {
			size += uint64(file.Size())

			continue
		}

		wg.Add(1)

		go func(f os.FileInfo) {
			defer wg.Done()

			err := filepath.Walk(newPath+"/"+f.Name(),
				func(_ string, fw fs.FileInfo, err error) error {
					size += uint64(fw.Size())
					return err
				})

			if err != nil {
				errSize = err
			}
		}(file)
	}

	wg.Wait()

	return size, errSize
}

func main() {
	path, _ := os.Getwd()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	sort.SliceStable(files, func(i, j int) bool {
		return strings.Compare(files[i].Name(), files[j].Name()) == -1
	})

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func(f os.FileInfo) {
			defer wg.Done()
			bytes, err := fastSize(path, f)
			if err == nil {
				line := fmt.Sprintf("%-7s %s\n", strings.ReplaceAll(humanize.Bytes(bytes), " ", ""), f.Name())
				if f.IsDir() {
					color.Blue(line)
				} else if strings.HasPrefix(f.Name(), ".") {
					color.Cyan(line)
				} else {
					color.White(line)
				}
			} else {
				color.Red(err.Error())
			}
		}(file)
	}

	wg.Wait()
}
