package main

import (
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/siongui/igdlcrx/extension/libbackground"
	"github.com/siongui/instago"
)

var mgr = instago.NewApiManager(nil, nil)
var isLocalhostAlive = false

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
	//storyurl := sss[3]

	/*
		item, err := libbackground.GetStoryItemFromStoryUrl(storyurl)
		if err != nil {
			println(err.Error())
			filename = ""
			return
		}

		filename = instago.BuildStoryFilename(mediaUrl, item.GetUsername(), item.GetUserId(), item.GetTimestamp())

		var appendIdUsernames []instago.IGTaggedUser
		for _, rm := range item.ReelMentions {
			pair := instago.IGTaggedUser{Id: rm.GetUserId(), Username: rm.GetUsername()}
			appendIdUsernames = append(appendIdUsernames, pair)
		}
		filename = instago.AppendTaggedUsersToFilename(username, item.GetUserId(), filename, appendIdUsernames)
	*/

	filename = username + "-" + timestamp + "-" + mediaUrl

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

func SendMessageToContentScript(msg string) {
	queryInfo := make(map[string]interface{})
	queryInfo["active"] = true
	queryInfo["currentWindow"] = true

	Chrome.Tabs.Call("query", queryInfo, func(tabs []map[string]interface{}) {
		if len(tabs) == 1 {
			Chrome.Tabs.Call("sendMessage", tabs[0]["id"], msg)
		}
	})
}

func IsLocalhostAlive() bool {
	// check if localhost server is alive
	resp, err := http.Get("http://localhost:8080/alive/")
	if err != nil {
		isLocalhostAlive = false
		return isLocalhostAlive
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		isLocalhostAlive = false
		return isLocalhostAlive
	}

	if string(body) == "ok" {
		isLocalhostAlive = true
		return isLocalhostAlive
	}

	isLocalhostAlive = false
	return isLocalhostAlive
}

func SendReelMentionsToContentScript(storyUrl string) {
	item, err := libbackground.GetStoryItemFromStoryUrl(storyUrl)
	if err != nil {
		println(err.Error())
		return
	}

	msg := ""
	for _, rm := range item.ReelMentions {
		if rm.IsPublic() {
			msg += rm.GetUserId() + ":" + rm.GetUsername() + ":" + rm.DisplayType + ":" + "public;"
		} else {
			msg += rm.GetUserId() + ":" + rm.GetUsername() + ":" + rm.DisplayType + ":" + "private;"
		}
	}
	SendMessageToContentScript("reelmentions:" + msg)
}

func main() {
	// Currently do nothing meaningful
	Chrome.Tabs.Get("onUpdated").Call("addListener", func(tabId int, changeInfo map[string]interface{}) {
		if _, ok := changeInfo["url"]; !ok {
			return
		}

		url := changeInfo["url"].(string)
		SendMessageToContentScript(url)
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
		if strings.HasPrefix(msg, "visitStoryUrl:") {
			storyUrl := strings.TrimPrefix(msg, "visitStoryUrl:")
			go SendReelMentionsToContentScript(storyUrl)
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
		if strings.HasPrefix(msg, "localhost:") {
			path := strings.TrimPrefix(msg, "localhost:")
			go func() {
				http.Get("http://localhost:8080" + path)
			}()
			return
		}

		if msg == "pageReload" {
			libbackground.ResetVariables()
			println("page reloaded")
			return
		}

		if msg == "isLocalhostAlive" {
			go func() {
				if IsLocalhostAlive() {
					SendMessageToContentScript("localhostIsAlive")
				}
			}()
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

	/*
		// get web url of reels tray
		rturl, err := mgr.GetGetWebFeedReelsTrayUrl()
		if err != nil {
			println(err.Error())
			return
		} else {
			println(rturl)
		}

		// get web reels tray
		rms, err := mgr.GetWebFeedReelsTray(rturl)
		if err != nil {
			println(err.Error())
			return
		}

		// set id - username pairs via data of reels tray
		for _, rm := range rms {
			libbackground.SetUsernameId(rm.User.Username, rm.User.Id)
		}
	*/
}
