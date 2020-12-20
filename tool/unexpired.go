package main

import (
	"fmt"

	"github.com/siongui/instago/download"
)

func main() {
	mgr, err := igdl.NewInstagramDownloadManager("auth.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = mgr.LoadCleanDownloadManager("auth-clean.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr.DownloadUnexpiredStoryOfAllFollowingUsers(2)
}
