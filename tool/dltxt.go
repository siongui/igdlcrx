package main

import (
	"flag"
	"fmt"

	"github.com/siongui/instago/download"
)

func main() {
	mgr, err := igdl.NewInstagramDownloadManager("auth-clean.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	f := flag.String("f", "users.txt", "file containing user id")
	flag.Parse()

	err = mgr.DownloadStoryFromUserIdFile(*f, 12, true)
	if err != nil {
		fmt.Println(err)
		return
	}
}
