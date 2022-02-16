package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"log"
	"os"
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
		if file.IsDir() {
			wg.Add(1)

			go func(f os.FileInfo) {
				defer wg.Done()
				res, err := fastSize(newPath, f)
				if err != nil {
					errSize = err
				}

				size += res
			}(file)

			continue
		}

		size += uint64(file.Size())
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

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func(f os.FileInfo) {
			defer wg.Done()
			bytes, err := fastSize(path, f)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%6s %12s\n", strings.ReplaceAll(humanize.Bytes(bytes), " ", ""), f.Name())
		}(file)
	}

	wg.Wait()
}
