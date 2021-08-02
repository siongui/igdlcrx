package main

import (
	"flag"
	"fmt"

	"github.com/siongui/instago/download"
)

func main() {
	id := flag.String("id", "25025320", "id of instagram user")
	clean := flag.Bool("clean", true, "if use clean account")
	flag.Parse()

	accountJson := "auth.json"
	if *clean {
		accountJson = "auth-clean.json"
	}

	mgr, err := igdl.NewInstagramDownloadManager(accountJson)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *clean {
		fmt.Println("Download unexpired stories (last 24 hours) of the user (no postlive) using clean account")
	} else {
		fmt.Println("Download unexpired stories (last 24 hours) of the user (no postlive)")
	}
	mgr.DownloadUserStory(*id)
}
