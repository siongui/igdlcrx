package main

import (
	"regexp"

	"github.com/fabioberger/chrome"
	. "github.com/siongui/godom"
)

func DoRootAction() {
	println("do root action")
}

func DoStoryAction() {
	println("do story action")
}

func DoUserAction() {
	println("do user action")
}

func CheckUrlAndDoAction(url string) {
	if IsRootUrl(url) {
		DoRootAction()
	}
	if IsStoryUrl(url) {
		DoStoryAction()
	}
	if IsUserUrl(url) {
		DoUserAction()
	}
}

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
		CheckUrlAndDoAction(url)
	})

	CheckUrlAndDoAction(Window.Location().Href())
}
