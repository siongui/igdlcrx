package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/siongui/instago"
)

var usernameId map[string]string

func GetStoryExt(mediaUrl string) string {
	urlnoq, _ := instago.StripQueryString(mediaUrl)
	eee := strings.Split(urlnoq, ".")
	return eee[len(eee)-1]
}

func GetTimeStr(timestamp string) string {
	t, _ := time.Parse(time.RFC3339, timestamp)
	loc := time.FixedZone("UTC+8", +8*60*60)
	return t.In(loc).Format(time.RFC3339) + "-" + strconv.FormatInt(t.Unix(), 10)
}

func GetStoryFilenameUrl(storyinfo string) (filename, mediaUrl string) {
	sss := strings.Split(storyinfo, ",")
	if len(sss) != 4 {
		return
	}
	username := sss[0]
	timestamp := sss[1]
	mediaUrl = sss[2]
	storyurl := sss[3]

	id, ok := usernameId[username]
	if !ok {
		id2, err := GetId(storyurl)
		if err == nil {
			id = id2
			usernameId[username] = id
		} else {
			println(err.Error())
			id = ""
		}
	}

	ext := GetStoryExt(mediaUrl)

	filename = username + "-" + id + "-story-" + GetTimeStr(timestamp) + "." + ext
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

func DownloadStory(storyinfo string) {
	options := make(map[string]string)
	filename, url := GetStoryFilenameUrl(storyinfo)
	options["url"] = url
	options["filename"] = filename
	//println(filename)
	//println(url)
	Chrome.Downloads.Call("download", options)
}

func main() {
	usernameId = make(map[string]string)

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
			go DownloadStory(storyinfo)
			return
		}
		println("Received msg from content: " + msg)
	})
}
