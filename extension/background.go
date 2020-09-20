package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/siongui/instago"
)

func GetStoryFilenameUrl(storyinfo string) (filename, url string) {
	sss := strings.Split(storyinfo, ",")
	if len(sss) != 3 {
		return
	}
	url = sss[2]

	urlnoq, err := instago.StripQueryString(url)
	if err != nil {
		return
	}
	eee := strings.Split(urlnoq, ".")
	ext := eee[len(eee)-1]

	timestamp := sss[1]
	t, _ := time.Parse(time.RFC3339, timestamp)
	loc := time.FixedZone("UTC+8", +8*60*60)

	filename = sss[0] + "-story-" + t.In(loc).Format(time.RFC3339) + "-" + strconv.FormatInt(t.Unix(), 10) + "." + ext
	// chrome.downloads does not allow ":" in filename
	filename = strings.Replace(filename, ":", "-", -1)
	return
}

func DownloadPost(code string) {
	em, err := instago.GetPostInfoNoLogin(code)
	if err != nil {
		println(err.Error())
		return
	}
	println(em)
}

func main() {
	// Currently do nothing meaningful
	Chrome.Tabs.Get("onUpdated").Call("addListener", func(tabId int, changeInfo map[string]interface{}) {
		if _, ok := changeInfo["url"]; !ok {
			return
		}

		url := changeInfo["url"].(string)

		queryInfo := make(map[string]interface{})
		queryInfo["active"] = true
		queryInfo["currentWindow"] = true

		Chrome.Tabs.Call("query", queryInfo, func(tabs []map[string]interface{}) {
			if len(tabs) == 1 {
				Chrome.Tabs.Call("sendMessage", tabs[0]["id"], url)
			}
		})
	})

	// Receive code of post from content.
	// Call chrome.downloads API to download files
	Chrome.Runtime.Get("onMessage").Call("addListener", func(message interface{}) {
		msg := message.(string)
		if strings.HasPrefix(msg, "postcode:") {
			code := strings.TrimPrefix(msg, "postcode:")
			/*
				createProperties := make(map[string]string)
				createProperties["url"] = "https://www.instagram.com/p/" + code + "/"
				Chrome.Tabs.Call("create", createProperties)
			*/
			go DownloadPost(code)
			return
		}
		if strings.HasPrefix(msg, "storyinfo:") {
			storyinfo := strings.TrimPrefix(msg, "storyinfo:")
			options := make(map[string]string)
			filename, url := GetStoryFilenameUrl(storyinfo)
			options["url"] = url
			options["filename"] = filename
			//println(filename)
			//println(url)
			Chrome.Downloads.Call("download", options)
			return
		}
		println("Received msg from content: " + msg)
	})
}
