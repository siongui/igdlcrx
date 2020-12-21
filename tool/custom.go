package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/siongui/instago/download"
)

func isLatestReelMediaDownloaded(username string, latestReelMedia int64) bool {
	utimes, err := igdl.GetReelMediaUnixTimesInUserStoryDir(username)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("In isLatestReelMediaDownloaded", err)
		}
		return false
	}

	lrm := strconv.FormatInt(latestReelMedia, 10)
	for _, utime := range utimes {
		if lrm == utime {
			return true
		}
	}
	return false
}

func DownloadEmptyIds(m *igdl.IGDownloadManager, emptyids []string) {
	trays, err := m.GetMultipleReelsMedia(emptyids)
	if err != nil {
		log.Println(err)
		return
	}

	for _, tray := range trays {
		for _, item := range tray.Items {
			username := tray.User.GetUsername()
			id := tray.User.GetUserId()
			_, err = igdl.GetStoryItem(item, username)
			if err != nil {
				igdl.PrintUsernameIdMsg(username, id, err)
			}
		}
	}
}

func DownloadReelsTray(m *igdl.IGDownloadManager, interval1 int, ignoreMuted, verbose bool) {
	rt, err := m.GetReelsTray()
	if err != nil {
		log.Println(err)
		return
	}

	go igdl.PrintLiveBroadcasts(rt.Broadcasts)

	emptyids := []string{}
	for _, tray := range rt.Trays {
		username := tray.GetUsername()
		id := tray.Id
		//items := tray.GetItems()

		if ignoreMuted && tray.Muted {
			if verbose {
				igdl.PrintUsernameIdMsg(username, id, " is muted && ignoreMuted set. no download")
			}
			continue
		}

		if isLatestReelMediaDownloaded(username, tray.LatestReelMedia) {
			if verbose {
				igdl.PrintUsernameIdMsg(username, id, " all downloaded")
			}
			continue
		}

		if tray.HasBestiesMedia {
			igdl.PrintUsernameIdMsg(username, id, "has close friend (besties) story item(s)")
		}

		if verbose {
			igdl.UsernameIdColorPrint(username, id)
			fmt.Println(" has undownloaded items")
		}

		items := tray.GetItems()
		if len(items) > 0 {
			for _, item := range items {
				_, err = igdl.GetStoryItem(item, username)
				if err != nil {
					igdl.PrintUsernameIdMsg(username, id, err)
				}
			}
		} else {
			emptyids = append(emptyids, strconv.FormatInt(id, 10))
		}

		if len(emptyids) > 20 {
			break
		}
	}

	DownloadEmptyIds(m, emptyids)
	//	igdl.PrintMsgSleep(interval1, "DownloadStoryAndPostLiveForever: ")
}

func main() {
	mgr, err := igdl.NewInstagramDownloadManager("auth.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	typ := flag.String("downloadtype", "timeline", "Download 1) timeline 2) story 3) highlight 4) saved posts")
	outputdir := flag.String("outputdir", "Instagram", "dir to save post and story")
	//datadir := flag.String("datadir", "Data", "dir to save data")
	flag.Parse()

	igdl.SetOutputDir(*outputdir)
	//igdl.SetDataDir(*datadir)

	switch *typ {
	case "timeline":
		fmt.Println("Download timeline")
		mgr.DownloadTimeline(1)
	case "story":
		fmt.Println("Download Stories and Post lives")
		DownloadReelsTray(mgr, 25, true, true)
	case "highlight":
		fmt.Println("Download all story highlights of all following users")
		mgr.DownloadStoryHighlights()
	case "saved":
		fmt.Println("Download all saved posts")
		mgr.DownloadSavedPosts(-1, false)
	default:
		fmt.Println("You have to choose a download type")
	}
}
