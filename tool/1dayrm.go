package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/siongui/instago/download"
)

func MovePost(fromdir, todir string) (err error) {
	err = filepath.Walk(fromdir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, e)
			return e
		}
		if info.Mode().IsRegular() {
			if strings.Contains(info.Name(), "-post-") {
				e2 := igdl.MoveFileToDir(todir, path, info)
				if e2 != nil {
					return e2
				}
			}
		}

		return nil
	})
	return
}

func main() {
	root := flag.String("root", "Instagram", "dir of downloaded files")
	todir := flag.String("todir", "~/Pictures", "dir to which files are moved")
	flag.Parse()

	igdl.MoveExpiredStory(*root, *todir)
	MovePost(*root, *todir)
}
