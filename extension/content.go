package main

import (
	"regexp"

	"github.com/fabioberger/chrome"
	. "github.com/siongui/godom"
)

func IsStoryUrl(url string) bool {
	re := regexp.MustCompile(`^https:\/\/www\.instagram\.com\/stories\/[a-zA-Z_.]+\/\d+\/$`)
	return re.MatchString(url)
}

func IsRootUrl(url string) bool {
	re := regexp.MustCompile(`^https:\/\/www\.instagram\.com\/$`)
	return re.MatchString(url)
}

func IsUserUrl(url string) bool {
	re := regexp.MustCompile(`^https:\/\/www\.instagram\.com\/[a-zA-Z_.]+\/$`)
	return re.MatchString(url)
}

func findStoryLink(url string) {
	sec := Document.QuerySelector("section._8XqED.carul")
	user := sec.QuerySelector("a.FPmhX.notranslate.R4sSg")
	println(user.InnerHTML())
}

func main() {
	c := chrome.NewChrome()

	c.Runtime.OnMessage(func(message interface{}, sender chrome.MessageSender, sendResponse func(interface{})) {
		url := message.(string)
		if IsStoryUrl(url) {
			println("Receive Story URL: " + url)
		}
		if IsRootUrl(url) {
			println("Receive Root URL: " + url)
		}
		if IsUserUrl(url) {
			println("Receive User URL: " + url)
		}
	})

	println(Window.Location().Href())
}
