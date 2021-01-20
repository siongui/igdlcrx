package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/siongui/instago/download"
)

func DownloadStoryFromUserIdFile(useridfile string, interval int64, verbose bool, mgr *igdl.IGDownloadManager) (err error) {
	idstrs, err := igdl.ReadNonCommentLines(useridfile)
	if err != nil {
		return
	}

	// usually there are at most 150 trays in reels_tray.
	// double the buffer to 300. 160 or 200 may be ok as well.
	c := make(chan igdl.TrayInfo, 300)

	for _, idstr := range idstrs {
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			return err
		}

		// layer = 2: also download reel mentions in story item
		c <- igdl.SetupTrayInfo(id, "", 2, false)
	}

	mgr.TrayDownloaderViaStoryAPI(c, igdl.NewTimeLimiter(interval), true, true, verbose)

	return
}

func main() {
	mgr, err := igdl.NewInstagramDownloadManager("auth-clean.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	f := flag.String("f", "users.txt", "file containing user id")
	flag.Parse()

	err = DownloadStoryFromUserIdFile(*f, 12, true, mgr)
	if err != nil {
		fmt.Println(err)
		return
	}
}
