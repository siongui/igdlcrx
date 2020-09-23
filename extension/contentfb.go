package main

import (
	"regexp"
	"time"

	. "github.com/siongui/godom"
	"github.com/siongui/instago"
)

func IsFacebookPhotoUrl(url string) bool {
	re := regexp.MustCompile(`^https:\/\/www\.facebook\.com\/photo\/\?fbid=\d+&set=[a-z\d.]+$`)
	return re.MatchString(url)
}

func IsFacebookStoryUrl(url string) bool {
	urlnoq, _ := instago.StripQueryString(url)
	re := regexp.MustCompile(`^https:\/\/www\.facebook\.com\/stories\/\d+\/[a-zA-Z\d=]+\/$`)
	return re.MatchString(urlnoq)
}

func CheckFacebookUrl(url string) {
	if IsFacebookPhotoUrl(url) {
		println("photo url: " + url)
	}
	if IsFacebookStoryUrl(url) {
		println("story url: " + url)
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
