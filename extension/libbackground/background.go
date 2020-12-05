package libbackground

import (
	"errors"
	"strings"

	"github.com/siongui/instago"
)

var mgr = instago.NewApiManager(nil, nil)
var usernameId = make(map[string]string)
var idUserTray = make(map[string]instago.UserTray)

func GetStoryIdFromStoryUrl(storyurl string) string {
	sss := strings.Split(storyurl, "/")
	return sss[len(sss)-2]
}

func GetUsernameFromStoryUrl(storyurl string) string {
	sss := strings.Split(storyurl, "/")
	return sss[len(sss)-3]
}

func ResetVariables() {
	usernameId = make(map[string]string)
	idUserTray = make(map[string]instago.UserTray)
}

func SetUsernameId(username, id string) {
	usernameId[username] = id
}

func GetIdFromUsername(username, storyurl string) (id string, err error) {
	id, ok := usernameId[username]
	if !ok {
		id2, err := mgr.GetIdFromWebStoryUrl(storyurl)
		if err == nil {
			id = id2
			usernameId[username] = id
		}
	}
	return
}

func GetStoryItem(id, storyurl string) (item instago.IGItem, err error) {
	// get user story tray if not exist
	tray, ok := idUserTray[id]
	if !ok {
		ut, err2 := mgr.GetUserStory(id)
		if err2 == nil {
			tray = ut
			idUserTray[id] = tray
		} else {
			err = err2
			return
		}
	}

	// get story item from item id
	for _, itm := range tray.Reel.Items {
		//println(GetStoryIdFromStoryUrl(storyurl))
		//println(itm.Id)
		if strings.HasPrefix(itm.Id, GetStoryIdFromStoryUrl(storyurl)) {
			item = itm
		}
	}

	// check if story item is found
	if item.GetTimestamp() == 0 {
		err = errors.New("item not found: " + storyurl)
	}
	return
}
