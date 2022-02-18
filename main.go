package main

import (
	"fmt"
	"github.com/bavix/dius/internal/wgi"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func shortPath(base, path string) string {
	newPath, _ := filepath.Rel(base, path)
	if newPath == "." {
		return newPath
	}

	if filepath.ToSlash(base) == filepath.Dir(path) {
		return filepath.Base(path)
	}

	return path
}

func main() {
	pwd, _ := os.Getwd()
	path := pwd
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var total uint64
	var wg wgi.WaitGroup
	for _, file := range files {
		wg.Add(1)

		go func(f os.FileInfo) {
			defer wg.Done()
			bytes, err := fastSize(path, f)
			if err != nil {
				color.Red(err.Error())
				return
			}

			total += bytes
			line := fmt.Sprintf(
				"%-7s %s\n",
				strings.ReplaceAll(humanize.IBytes(bytes), " ", ""),
				shortPath(pwd, filepath.Join(path, f.Name())))

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

	color.Green(
		fmt.Sprintf("%-7s %s\n",
			strings.ReplaceAll(humanize.IBytes(total), " ", ""),
			shortPath(pwd, path)))
}
