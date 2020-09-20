package main

import (
	"regexp"
	"strings"
	"time"

	. "github.com/siongui/godom"
)

var debug = true

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
		if debug {
			println("cannot find img element in GetBestImageUrl")
		}
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

	for i, source := range videos[0].QuerySelectorAll("source") {
		vUrl = source.Call("getAttribute", "src").String()
		if i == 0 {
			break
		}
		//println(vUrl)
	}

	return vUrl
}

func ProcessArticleInRootPath(article *Object) {
	btns := article.QuerySelectorAll(".download-timeline-post-btn")
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
	btn.ClassList().Add("download-timeline-post-btn")
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
		if debug {
			println("cannot find section in DoStoryAction")
		}
		return
	}

	btns := section.QuerySelectorAll(".download-story-btn")
	if len(btns) > 0 {
		if debug {
			println("story download button exist. exit.")
		}
		return
	}

	userElm, ok := GetElementInElement(section, "a.FPmhX.notranslate.R4sSg")
	if !ok {
		if debug {
			println("cannot find userElm in DoStoryAction")
		}
		return
	}
	username := userElm.Call("getAttribute", "title").String()
	if debug {
		println("story username: " + username)
	}

	timeElm, ok := GetElementInElement(section, "time")
	if !ok {
		if debug {
			println("cannot find timeElm in DoStoryAction")
		}
		return
	}
	time := timeElm.Call("getAttribute", "datetime").String()
	if debug {
		println("story timestamp: " + time)
	}

	mediaElm, ok := GetElementInElement(section, "div.qbCDp")
	if !ok {
		if debug {
			println("cannot find mediaElm in DoStoryAction")
		}
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
	if debug {
		println("story url: " + url)
	}

	btn := Document.CreateElement("button")
	btn.Dataset().Set("dataUrl", url)
	btn.ClassList().Add("download-story-btn")
	btn.SetInnerHTML("Download")
	btn.AddEventListener("click", func(e Event) {
		// send code of post to background for download
		Chrome.Runtime.Call("sendMessage", "storyinfo:"+username+","+time+","+url)
	})
	controlElm, ok := GetElementInElement(section, "div.GHEPc")
	if !ok {
		return
	}
	controlElm.AppendChild(btn)
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
	re := regexp.MustCompile(`^https:\/\/www\.instagram\.com\/stories\/[a-zA-Z\d_.]+\/\d+\/$`)
	return re.MatchString(url)
}

func IsRootUrl(url string) bool {
	re := regexp.MustCompile(`^https:\/\/www\.instagram\.com\/$`)
	return re.MatchString(url)
}

func IsUserUrl(url string) bool {
	re := regexp.MustCompile(`^https:\/\/www\.instagram\.com\/[a-zA-Z\d_.]+\/$`)
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
