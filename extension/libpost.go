package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/siongui/instago"
)

// Decode JSON data returned by Instagram post API
type postInfo struct {
	GraphQL struct {
		ShortcodeMedia instago.IGMedia `json:"shortcode_media"`
	} `json:"graphql"`
}

func GetPostInfo(code string) (em instago.IGMedia, err error) {
	url := instago.CodeToUrl(code) + "?__a=1"

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("resp.StatusCode = " + strconv.Itoa(resp.StatusCode))
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	pi := postInfo{}
	err = json.Unmarshal(b, &pi)
	if err != nil {
		return
	}
	em = pi.GraphQL.ShortcodeMedia
	return
}
