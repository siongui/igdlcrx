package main

import (
	"flag"
	"fmt"

	"github.com/siongui/instago/download"
)

func main() {
	mgr, err := igdl.NewInstagramDownloadManager("auth.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	root := flag.String("root", "Instagram", "dir of downloaded files")
	todir := flag.String("todir", "~/Pictures", "dir to which files are moved")
	flag.Parse()

	mgr.CheckZero(*root, *todir)
}
