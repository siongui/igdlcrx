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

	f := flag.String("f", "users.txt", "file containing user id")
	flag.Parse()

	ids, err := igdl.ReadNonCommentLines(*f)
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr.DownloadEmptyIds(ids)
}
