package du

import (
	"fmt"
	"github.com/bavix/dius/internal/fs"
	"github.com/bavix/dius/internal/wgi"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

func Execute(_ *cobra.Command, args []string) {
	pwd, _ := os.Getwd()
	path := pwd
	if len(args) > 0 {
		path = args[0]
	}

	var files []os.FileInfo
	pathInfo, err := os.Stat(path)
	pathExists := !os.IsNotExist(err)
	pathIsFile := pathExists && pathInfo.Mode().IsRegular()
	if err != nil {
		log.Fatal(err)
	}

	if pathIsFile {
		files = append(files, pathInfo)
	} else {
		files, err = ioutil.ReadDir(path)
		for attempts := 0; attempts < 10 && err != nil; attempts++ {
			time.Sleep(time.Millisecond * 50)
			files, err = ioutil.ReadDir(path)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	var total uint64
	var wg wgi.WaitGroup
	for _, file := range files {
		wg.Add(1)

		fi := fs.New(path, file)
		go func() {
			defer wg.Done()
			bytes, err := fi.Size()
			if err != nil {
				color.Red(err.Error())
				return
			}

			filename := filepath.Join(path, fi.Name())
			if pathIsFile {
				filename = path
			}

			total += bytes
			line := fmt.Sprintf(
				"%-8s %s\n",
				strings.ReplaceAll(humanize.IBytes(bytes), " ", ""),
				shortPath(pwd, filename))

			if fi.IsDir() {
				color.Blue(line)
			} else if strings.HasPrefix(fi.Name(), ".") {
				color.Cyan(line)
			} else {
				color.White(line)
			}
		}()
	}

	wg.Wait()

	if pathIsFile {
		color.Green(
			fmt.Sprintf("%-8s %s\n",
				strings.ReplaceAll(humanize.IBytes(total), " ", ""),
				shortPath(pwd, path)))
	}
}
