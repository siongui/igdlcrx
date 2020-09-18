package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/fabioberger/chrome"
	. "github.com/siongui/godom"
)

func DoRootAction() {
	println("do root action")
	articles := Document.QuerySelectorAll("article[role='presentation']")
	for _, article := range articles {
		header := article.QuerySelector("header")
		userElm := header.QuerySelector("a.sqdOP.yWX7d._8A5w5.ZIAjV")
		username := userElm.InnerHTML()

		codetimeElm := article.QuerySelector("div.k_Q0X.NnvRN")
		codeElm := codetimeElm.QuerySelector("a")
		code := strings.TrimPrefix(codeElm.Call("getAttribute", "href").String(), "/p/")
		code = strings.TrimSuffix(code, "/")
		timeElm := codetimeElm.QuerySelector("time")
		time := timeElm.Call("getAttribute", "datetime").String()

		println(username + " " + code + " " + time)

		mediaElm := article.QuerySelector("div.KL4Bh")
		imgs := mediaElm.QuerySelectorAll("img")
		if len(imgs) == 1 {
			//src := imgs[0].Call("getAttribute", "src").String()
			//println("src: " + src)
			srcset := imgs[0].Call("getAttribute", "srcset").String()
			//println(srcset)
			srcs := strings.Split(srcset, ",")
			bestsrc := srcs[len(srcs)-1]
			s := strings.Split(bestsrc, " ")[0]
			println("best src: " + s)
		}
	}
}

func DoStoryAction() {
	println("do story action")
}

func DoUserAction() {
	println("do user action")
}

func CheckUrlAndDoAction(url string) {
	println(time.Now().Format(time.RFC3339))
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

func main() {
	c := chrome.NewChrome()

	// Currently this receiver do nothing meaningful.
	// Just print received URL.
	c.Runtime.OnMessage(func(message interface{}, sender chrome.MessageSender, sendResponse func(interface{})) {
		url := message.(string)
		//CheckUrlAndDoAction(url)
		println("Received URL from background: " + url)
	})

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				CheckUrlAndDoAction(Window.Location().Href())
			}
		}
	}()

}
