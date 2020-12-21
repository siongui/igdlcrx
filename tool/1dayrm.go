package main

import (
	"flag"

	"github.com/siongui/instago/download"
)

func main() {
	root := flag.String("root", "Instagram", "dir of downloaded files")
	todir := flag.String("todir", "~/Pictures", "dir to which files are moved")
	flag.Parse()

	igdl.MoveExpiredStory(*root, *todir)
}
