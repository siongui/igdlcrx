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

	if !img.HasAttribute("srcset") {
		if debug {
			println("cannot find img srcset attribute in GetBestImageUrl")
		}
		return ""
	}

	srcset := img.GetAttribute("srcset")
	if debug {
		println(srcset)
	}
	srcs := strings.Split(srcset, ",")
	if len(srcs) == 0 {
		return ""
	}
	bestsrc := srcs[len(srcs)-1]
	s := strings.Split(bestsrc, " ")[0]
	if debug {
		println("best src: " + s)
	}
	return s
}

func GetVideoUrl(mediaElm *Object) string {
	videos := mediaElm.QuerySelectorAll("video")

	vUrl := ""
	if len(videos) != 1 {
		return vUrl
	}

	for i, source := range videos[0].QuerySelectorAll("source") {
		if !source.HasAttribute("src") {
			if debug {
				println("video source src attr not exist")
			}
			return ""
		}
		vUrl = source.GetAttribute("src")
		if debug {
			println(vUrl)
		}
		if i == 0 {
			break
		}
	}

	if debug {
		println("video url: " + vUrl)
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
	if !codeElm.HasAttribute("href") {
		return
	}
	code := strings.TrimPrefix(codeElm.GetAttribute("href"), "/p/")
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
	//println("do story action")

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

	if !userElm.HasAttribute("title") {
		if debug {
			println("cannot find userElm title attribute in DoStoryAction")
		}
		return
	}
	username := userElm.GetAttribute("title")
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

	if !timeElm.HasAttribute("datetime") {
		if debug {
			println("cannot find timeElm datetime attribute in DoStoryAction")
		}
		return
	}

	timestamp := timeElm.GetAttribute("datetime")
	if debug {
		println("story timestamp: " + timestamp)
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

	mediaUrl := ""
	if url2 == "" {
		mediaUrl = url1
	} else {
		mediaUrl = url2
	}
	if mediaUrl == "" {
		if debug {
			println("mediaUrl is empty string in DoStoryAction")
		}
		return
	}

	if debug {
		println("story mediaUrl: " + mediaUrl)
	}

	btn := Document.CreateElement("button")
	btn.Dataset().Set("dataMediaUrl", mediaUrl)
	btn.ClassList().Add("download-story-btn")
	btn.SetInnerHTML("Download")
	btn.AddEventListener("click", func(e Event) {
		// send code of post to background for download
		Chrome.Runtime.Call("sendMessage", "storyinfo:"+username+","+timestamp+","+mediaUrl+","+Window.Location().Href())
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
