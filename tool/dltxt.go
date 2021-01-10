package main

import (
	"flag"
	"fmt"
	"strconv"
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
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		//fmt.Println(i)
		//continue

		err = mgr.DownloadUserStoryPostlive(i)
		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(12 * time.Second)
	}
}
