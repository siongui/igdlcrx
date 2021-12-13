package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/siongui/instago/download"
)

func main() {
	mgr, err := igdl.NewInstagramDownloadManager("auth.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = mgr.LoadCollectionList()
	if err != nil {
		fmt.Println(err)
		return
	}

	var collectionName string
	flag.StringVar(&collectionName, "collection", "Download", "name of the collection")
	flag.Parse()

	sleepInterval := 12 // seconds

	for {
		fmt.Println("Download saved post in collection:", collectionName)
		err = mgr.DownloadSavedCollectionPosts(collectionName)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("=========================")
		fmt.Print(time.Now().Format(time.RFC3339))
		fmt.Println(": sleep ", sleepInterval, " second")
		fmt.Println("=========================")
		time.Sleep(time.Duration(sleepInterval) * time.Second)
	}
}
