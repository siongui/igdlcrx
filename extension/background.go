package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/siongui/instago"
)

var mgr = instago.NewApiManager(nil, nil)
var usernameId map[string]string

func DownloadFBStory(fbstory string) {
	sss := strings.Split(fbstory, ",")
	if len(sss) != 2 {
		println("facebook story message not correct")
		return
	}
	username := sss[0]
	// FIXME: handle blob url
	url := sss[1]

	//println(username)
	//println(url)

	filename := username + "-facebook-story." + GetStoryExt(url)

	options := make(map[string]string)
	options["url"] = url
	options["filename"] = Rename(filename)
	//println(filename)
	//println(url)
	Chrome.Downloads.Call("download", options)
}

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

func Rename(s string) string {
	// chrome.downloads does not allow ":" in filename
	return strings.Replace(s, ":", "_", -1)
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
		id2, err := mgr.GetIdFromWebStoryUrl(storyurl)
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
	filename = Rename(filename)
	return
}

func DownloadIGMedia(em instago.IGMedia) (err error) {
	urls, err := em.GetMediaUrls()
	if err != nil {
		return
	}

	for index, url := range urls {
		var taggedusers []instago.IGTaggedUser
		if len(urls) == 1 {
			taggedusers = em.EdgeMediaToTaggedUser.GetIdUsernamePairs()
		} else {
			taggedusers = em.EdgeSidecarToChildren.Edges[index].Node.EdgeMediaToTaggedUser.GetIdUsernamePairs()
		}

		// prevent panic in instago.BuildFilename method
		_, err = instago.StripQueryString(url)
		if err != nil {
			return
		}

		filename := instago.GetPostFilename(
			em.GetUsername(),
			em.GetUserId(),
			em.GetPostCode(),
			url,
			em.GetTimestamp(),
			taggedusers)
		if index > 0 {
			filename = instago.AppendIndexToFilename(filename, index)
		}

		options := make(map[string]string)
		options["url"] = url
		options["filename"] = Rename(filename)
		//println(filename)
		//println(url)
		Chrome.Downloads.Call("download", options)

	}
	return
}

func DownloadPost(code string) {
	em, err := instago.GetPostInfoNoLogin(code)
	if err != nil {
		println(err.Error())
		return
	}
	err = DownloadIGMedia(em)
	if err != nil {
		println(err.Error())
		return
	}
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
		if strings.HasPrefix(msg, "fbstory:") {
			fbstory := strings.TrimPrefix(msg, "fbstory:")
			go DownloadFBStory(fbstory)
			return
		}
		println("Received msg from content: " + msg)
	})

	storyQH, u1, u2, err := mgr.GetWebQueryHash()
	if err == nil {
		println(storyQH)
		println(u1)
		println(u2)
	} else {
		println(err.Error())
	}
}
