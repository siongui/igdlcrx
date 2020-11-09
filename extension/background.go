package main

import (
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/siongui/instago"
)

var mgr = instago.NewApiManager(nil, nil)
var usernameId map[string]string
var idUserTray map[string]instago.UserTray

func DownloadFBPhoto(fbphoto string) {
	sss := strings.Split(fbphoto, ",,,")
	if len(sss) != 2 {
		println("facebook photo message not correct")
		println(fbphoto)
		return
	}
	username := sss[0]
	url := sss[1]

	urlnoq, _ := instago.StripQueryString(url)
	filename := username + "-facebook-photo-" + path.Base(urlnoq)

	options := make(map[string]string)
	options["url"] = url
	options["filename"] = Rename(filename)
	//println(filename)
	//println(url)
	Chrome.Downloads.Call("download", options)
}

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

	ext := "jpg"
	if strings.HasPrefix(url, "blob:") {
		ext = "mp4"
	}

	filename := username + "-facebook-story." + ext

	options := make(map[string]string)
	options["url"] = url
	options["filename"] = Rename(filename)
	//println(filename)
	//println(url)
	Chrome.Downloads.Call("download", options)
}

// To be removed
func GetTimeStr(timestamp string) (string, string) {
	t, _ := time.Parse(time.RFC3339, timestamp)
	loc := time.FixedZone("UTC+8", +8*60*60)
	return t.In(loc).Format(time.RFC3339), strconv.FormatInt(t.Unix(), 10)
}

func GetStoryId(storyurl string) string {
	sss := strings.Split(storyurl, "/")
	return sss[len(sss)-2]
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
	//timestamp := sss[1]
	mediaUrl = sss[2]
	storyurl := sss[3]

	// get id from username
	id, ok := usernameId[username]
	if !ok {
		id2, err := mgr.GetIdFromWebStoryUrl(storyurl)
		if err == nil {
			id = id2
			usernameId[username] = id
		} else {
			println(err.Error())
			filename = ""
			return
		}
	}

	// get user story tray if not exist
	tray, ok := idUserTray[id]
	if !ok {
		ut, err := mgr.GetUserStory(id)
		if err == nil {
			tray = ut
			idUserTray[id] = tray
		} else {
			println(err.Error())
			filename = ""
			return
		}
	}

	// get story item from item id
	item := instago.IGItem{}
	for _, itm := range tray.Reel.Items {
		//println(GetStoryId(storyurl))
		//println(itm.Id)
		if strings.HasPrefix(itm.Id, GetStoryId(storyurl)) {
			item = itm
		}
	}
	// story item does not exist, return
	if item.GetTimestamp() == 0 {
		println("story item not found")
		filename = ""
		return
	}

	filename = instago.BuildStoryFilename(mediaUrl, item.GetUsername(), item.GetUserId(), item.GetTimestamp())

	var appendIdUsernames []instago.IGTaggedUser
	for _, rm := range item.ReelMentions {
		pair := instago.IGTaggedUser{Id: rm.GetUserId(), Username: rm.GetUsername()}
		appendIdUsernames = append(appendIdUsernames, pair)
	}
	filename = instago.AppendTaggedUsersToFilename(username, id, filename, appendIdUsernames)

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
	if filename == "" {
		return
	}
	options["url"] = url
	options["filename"] = filename
	//println(filename)
	//println(url)
	Chrome.Downloads.Call("download", options)
}

func main() {
	usernameId = make(map[string]string)
	idUserTray = make(map[string]instago.UserTray)

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
		if strings.HasPrefix(msg, "fbphoto:") {
			fbphoto := strings.TrimPrefix(msg, "fbphoto:")
			go DownloadFBPhoto(fbphoto)
			return
		}

		if msg == "pageReload" {
			idUserTray = make(map[string]instago.UserTray)
			println("page reloaded")
			return
		}

		println("Received msg from content: " + msg)
	})

	/*
		storyQH, u1, u2, err := mgr.GetWebQueryHash()
		if err == nil {
			println(storyQH)
			println(u1)
			println(u2)
		} else {
			println(err.Error())
		}
	*/

	rturl, err := mgr.GetGetWebFeedReelsTrayUrl()
	if err != nil {
		println(err.Error())
		return
	} else {
		println(rturl)
	}

	rms, err := mgr.GetWebFeedReelsTray(rturl)
	if err != nil {
		println(err.Error())
		return
	}

	for _, rm := range rms {
		usernameId[rm.User.Username] = rm.User.Id
	}
}
