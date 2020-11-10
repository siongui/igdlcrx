package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/siongui/instago"
	igdl "github.com/siongui/instago/download"
)

var usernameId map[string]string
var idUserTray map[string]instago.UserTray

type storyHandler struct {
	mgr *instago.IGApiManager
}

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/alive/" {
		fmt.Fprintf(w, "ok")
	}
}

func (sh *storyHandler) storyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url path: " + r.URL.Path)

	sss := strings.Split(r.URL.Path, "/")
	if len(sss) != 5 {
		return
	}

	username := sss[2]
	storyid := sss[3]
	//fmt.Println("username:", username, ", story id:", storyid)

	// get id from username
	id, ok := usernameId[username]
	if !ok {
		fmt.Println(username, "id not in cache")
		id2, err := sh.mgr.GetIdFromWebStoryUrl("https://www.instagram.com" + r.URL.Path)
		if err == nil {
			id = id2
			usernameId[username] = id
		} else {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("username:", username, ", id: ", id, ", story id:", storyid)

	// get user story tray if not exist
	tray, ok := idUserTray[id]
	if !ok {
		ut, err := sh.mgr.GetUserStory(id)
		if err == nil {
			tray = ut
			idUserTray[id] = tray
		} else {
			fmt.Println("fail to get user story tray:", err)
			return
		}
	}

	// get story item from item id
	item := instago.IGItem{}
	for _, itm := range tray.Reel.Items {
		if strings.HasPrefix(itm.Id, storyid) {
			item = itm
		}
	}
	// story item does not exist, return
	if item.GetTimestamp() == 0 {
		fmt.Println("story item not found")
		return
	}

	_, err := igdl.GetStoryItem(item, username)
	if err != nil {
		fmt.Println("fail to download story item:", err)
		return
	}
}

func main() {
	usernameId = make(map[string]string)
	idUserTray = make(map[string]instago.UserTray)

	mgr, err := instago.NewInstagramApiManager("auth.json")
	if err != nil {
		log.Fatal(err)
	}

	url, err := mgr.GetGetWebFeedReelsTrayUrl()
	if err != nil {
		log.Fatal(err)
	}

	rms, err := mgr.GetWebFeedReelsTray(url)
	if err != nil {
		log.Fatal(err)
	}

	for i, rm := range rms {
		fmt.Println(i, rm.User.Username, rm.User.Id)
		usernameId[rm.User.Username] = rm.User.Id
	}

	myStoryHandler := &storyHandler{mgr: mgr}

	http.HandleFunc("/alive/", aliveHandler)
	http.HandleFunc("/stories/", myStoryHandler.storyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
