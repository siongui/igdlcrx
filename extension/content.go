package main

import (
	"regexp"
	"strings"
	"time"

	. "github.com/siongui/godom"
)

func GetElementInElement(element *Object, selector string) (elm *Object, ok bool) {
	elms := element.QuerySelectorAll(selector)
	if len(elms) == 1 {
		elm = elms[0]
		ok = true
		return
	}
	return
}

func GetBestImageUrl(mediaElm *Object) string {
	img, ok := GetElementInElement(mediaElm, "img")
	if !ok {
		return ""
	}

	//src := img.Call("getAttribute", "src").String()
	//println("src: " + src)
	srcset := img.Call("getAttribute", "srcset").String()
	//println(srcset)
	srcs := strings.Split(srcset, ",")
	bestsrc := srcs[len(srcs)-1]
	s := strings.Split(bestsrc, " ")[0]
	//println("best src: " + s)
	return s
}

func GetVideoUrl(mediaElm *Object) string {
	videos := mediaElm.QuerySelectorAll("video")

	vUrl := ""
	if len(videos) != 1 {
		return vUrl
	}

	for _, source := range videos[0].QuerySelectorAll("source") {
		vUrl = source.Call("getAttribute", "src").String()
		//println(vUrl)
	}

	return vUrl
}

func ProcessArticleInRootPath(article *Object) {
	btns := article.QuerySelectorAll(".download-btn")
	if len(btns) > 0 {
		return
	}
	/*
		header, ok := GetElementInElement(article, "header")
		if !ok {
			return
		}


			userElm, ok := GetElementInElement(header, "a.sqdOP.yWX7d._8A5w5.ZIAjV")
			if !ok {
				return
			}
			username := userElm.InnerHTML()
	*/

	codetimeElm, ok := GetElementInElement(article, "div.k_Q0X.NnvRN")
	if !ok {
		return
	}
	codeElm, ok := GetElementInElement(codetimeElm, "a")
	if !ok {
		return
	}
	code := strings.TrimPrefix(codeElm.Call("getAttribute", "href").String(), "/p/")
	code = strings.TrimSuffix(code, "/")

	/*
			timeElm, ok := GetElementInElement(codetimeElm, "time")
			if !ok {
				return
			}
			time := timeElm.Call("getAttribute", "datetime").String()
			println(username + " " + code + " " + time)


		mediaElm, ok := GetElementInElement(article, "div.KL4Bh")
		if !ok {
			return
		}
		GetBestImageUrl(mediaElm)
		// TODO: how to check if video in post?
	*/

	btn := Document.CreateElement("button")
	btn.Dataset().Set("dataCode", code)
	btn.ClassList().Add("download-btn")
	btn.SetInnerHTML("Download")
	btn.AddEventListener("click", func(e Event) {
		// send code of post to background for download
		Chrome.Runtime.Call("sendMessage", "postcode:"+code)
	})
	article.Call("prepend", btn)
}

func DoRootAction() {
	println("do root action")
	articles := Document.QuerySelectorAll("article[role='presentation']")
	for _, article := range articles {
		ProcessArticleInRootPath(article)
	}
}

func DoStoryAction() {
	println("do story action")
	section, ok := GetElementInElement(Document, "section._8XqED.carul")
	if !ok {
		return
	}

	/*
		userElm, ok := GetElementInElement(section, "a.FPmhX.notranslate.R4sSg")
		if !ok {
			return
		}
		username := userElm.Call("getAttribute", "title").String()
		println(username)
	*/

	/*
		timeElm, ok := GetElementInElement(section, "time")
		if !ok {
			return
		}
		time := timeElm.Call("getAttribute", "datetime").String()
		println(time)
	*/

	mediaElm, ok := GetElementInElement(section, "div.qbCDp")
	if !ok {
		return
	}
	url1 := GetBestImageUrl(mediaElm)
	url2 := GetVideoUrl(mediaElm)

	url := ""
	if url2 == "" {
		url = url1
	} else {
		url = url2
	}
	// send code of post to background for download
	Chrome.Runtime.Call("sendMessage", url)
}

func DoUserAction() {
	println("do user action")
}

func CheckUrlAndDoAction(url string) {
	//println(time.Now().Format(time.RFC3339))
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
	// Currently this receiver do nothing meaningful.
	// Just print received URL.
	Chrome.Runtime.Get("onMessage").Call("addListener", func(message interface{}) {
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
