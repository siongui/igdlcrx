package main

import (
	"flag"
	"fmt"

	"github.com/siongui/instago"
	"github.com/siongui/instago/download"
)

func main() {
	instago.SetUserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 142.0.0.22.109 (iPhone12,5; iOS 14_1; en_US; en-US; scale=3.00; 1242x2688; 214888322) NW/1")
	mgr, err := igdl.NewInstagramDownloadManager("auth-clean.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	f := flag.String("f", "users.txt", "file containing user id")
	flag.Parse()

	err = mgr.DownloadStoryFromUserIdFile(*f, 15, true)
	if err != nil {
		fmt.Println(err)
		return
	}
}
