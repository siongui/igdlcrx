package main

import (
	"fmt"
	"flag"

	"github.com/siongui/instago/download"
)

func main() {
	mgr, err := igdl.NewInstagramDownloadManager("auth.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	id := flag.Int64("id", 25025320, "id of instagram user")
	flag.Parse()

	fmt.Println("Download unexpired stories (last 24 hours) of the user reel media")
	mgr.DownloadUserStory(*id)
}
