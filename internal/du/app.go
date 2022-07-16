package du

import (
	"context"
	"fmt"
	"github.com/bavix/dius/internal/fs"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		files, err = ioutil.ReadDir(path)
		for attempts := 0; err != nil; attempts++ {
			select {
			case <-ctx.Done():
				log.Fatal(err)
			default:
				time.Sleep(time.Millisecond * time.Duration(attempts*10))
				files, err = ioutil.ReadDir(path)
			}
		}
	}

	var total uint64
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
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

			atomic.AddUint64(&total, bytes)
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

	if !pathIsFile {
		color.Green(
			fmt.Sprintf("%-8s %s\n",
				strings.ReplaceAll(humanize.IBytes(total), " ", ""),
				shortPath(pwd, path)))
	}
}
