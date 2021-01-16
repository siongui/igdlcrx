package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

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

	ids, err := igdl.ReadNonCommentLines(*f)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, id := range ids {
		id = strings.TrimSpace(id)

		err = mgr.DownloadUserStoryLayer(id, 2, 12)
		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(12 * time.Second)
	}
}
