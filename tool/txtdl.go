package main

import (
	"flag"
	"fmt"
	"strconv"

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

	for _, id := range ids {
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Download unexpired stories (last 24 hours) of the user reel media", n)
		mgr.DownloadUserReelMedia(n)
		return
	}
}
