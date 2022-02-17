package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func fastSize(path string, info os.FileInfo) (uint64, error) {
	if !info.IsDir() {
		return uint64(info.Size()), nil
	}

	newPath := filepath.Join(path, info.Name())
	files, err := ioutil.ReadDir(newPath)
	if err != nil {
		return 0, err
	}

	var errSize error = nil
	var size uint64 = 0
	var wg sync.WaitGroup
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

func main() {
	path, _ := os.Getwd()
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func(f os.FileInfo) {
			defer wg.Done()
			bytes, err := fastSize(path, f)
			if err != nil {
				color.Red(err.Error())
				return
			}

			line := fmt.Sprintf("%-7s %s\n", strings.ReplaceAll(humanize.Bytes(bytes), " ", ""), f.Name())
			if f.IsDir() {
				color.Blue(line)
			} else if strings.HasPrefix(f.Name(), ".") {
				color.Cyan(line)
			} else {
				color.White(line)
			}
		}(file)
	}

	wg.Wait()
}
