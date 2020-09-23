package main

import (
	"regexp"
	"time"

	. "github.com/siongui/godom"
	"github.com/siongui/instago"
)

func DoFacebookPhotoAction(url string) {
	println("photo url: " + url)
}

func DoFacebookStoryAction(url string) {
	println("story url: " + url)
	storyElm, ok := GetElementInElement(Document, "div[data-pagelet='Stories']")
	if !ok {
		return
	}

	imgUrl := ""
	// FIXME: not the imgElm we want
	imgElm, ok := GetElementInElement(storyElm, "img")
	if ok {
		imgUrl = imgElm.GetAttribute("src")
	}
	println(imgUrl)
}

func IsFacebookPhotoUrl(url string) bool {
	re1 := regexp.MustCompile(`^https:\/\/www\.facebook\.com\/photo\/\?fbid=\d+&set=[a-z\d.]+$`)
	re2 := regexp.MustCompile(`^https:\/\/www\.facebook\.com\/[a-zA-Z\d.]+\/photos\/[a-zA-Z\d.]+\/[a-zA-Z\d.]+\/?$`)
	return re1.MatchString(url) || re2.MatchString(url)
}

func IsFacebookStoryUrl(url string) bool {
	urlnoq, _ := instago.StripQueryString(url)
	re := regexp.MustCompile(`^https:\/\/www\.facebook\.com\/stories\/\d+\/[a-zA-Z\d=]+\/$`)
	return re.MatchString(urlnoq) || (url == "https://www.facebook.com/stories")
}

func CheckFacebookUrl(url string) {
	if IsFacebookPhotoUrl(url) {
		DoFacebookPhotoAction(url)
	}
	if IsFacebookStoryUrl(url) {
		DoFacebookStoryAction(url)
	}
}

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				CheckFacebookUrl(Window.Location().Href())
			}
		}
	}()
}
