package main

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/siongui/instago"
)

type IGStoryUser struct {
	User struct {
		Id            string `json:"id"`
		ProfilePicUrl string `json:"profile_pic_url"`
		Username      string `json:"username"`
	} `json:"user"`
}

func IsStoryUrl(url string) bool {
	re := regexp.MustCompile(`^https:\/\/www\.instagram\.com\/stories\/[a-zA-Z\d_.]+\/\d+\/$`)
	return re.MatchString(url)
}

func GetUserFromStoryUrl(url string) (user IGStoryUser, err error) {
	if !IsStoryUrl(url) {
		err = errors.New(url + " is not a valid story url")
		return
	}

	jsonurl := url + "?__a=1"
	b, err := instago.GetHTTPResponseNoLogin(jsonurl)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &user)
	return
}

func GetId(url string) (id string, err error) {
	user, err := GetUserFromStoryUrl(url)
	id = user.User.Id
	return
}
