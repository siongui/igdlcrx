package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/siongui/instago"
)

var mgr instago.IGApiManager
var usernameId map[string]string

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/alive/" {
		fmt.Fprintf(w, "ok")
	}
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
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
		id2, err := mgr.GetIdFromWebStoryUrl("https://www.instagram.com" + r.URL.Path)
		if err == nil {
			id = id2
			usernameId[username] = id
		} else {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("username:", username, ", id: ", id, ", story id:", storyid)
}

func main() {
	usernameId = make(map[string]string)

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

	http.HandleFunc("/alive/", aliveHandler)
	http.HandleFunc("/stories/", storyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
