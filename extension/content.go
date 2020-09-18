package main

import (
	"regexp"
	"strings"
	"time"

	. "github.com/siongui/godom"
)

func GetBestImageUrl(imgs []*Object) string {
	if len(imgs) != 1 {
		return ""
	}

	//src := imgs[0].Call("getAttribute", "src").String()
	//println("src: " + src)
	srcset := imgs[0].Call("getAttribute", "srcset").String()
	//println(srcset)
	srcs := strings.Split(srcset, ",")
	bestsrc := srcs[len(srcs)-1]
	s := strings.Split(bestsrc, " ")[0]
	println("best src: " + s)
	return s
}

func GetVideoUrl(videos []*Object) string {
	vUrl := ""
	if len(videos) != 1 {
		return vUrl
	}

	for _, source := range videos[0].QuerySelectorAll("source") {
		vUrl = source.Call("getAttribute", "src").String()
		println(vUrl)
	}

	return vUrl
}

func ProcessArticleInRootPath(article *Object) {
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
	GetBestImageUrl(imgs)

	// send code of post to background for download
	Chrome.Runtime.Call("sendMessage", code)
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
	sections := Document.QuerySelectorAll("section._8XqED.carul")
	if len(sections) != 1 {
		println("cannot find story element: section._8XqED.carul")
		return
	}

	userElm := sections[0].QuerySelector("a.FPmhX.notranslate.R4sSg")
	username := userElm.Call("getAttribute", "title").String()
	println(username)

	timeElm := sections[0].QuerySelector("time")
	time := timeElm.Call("getAttribute", "datetime").String()
	println(time)

	mediaElm := sections[0].QuerySelector("div.qbCDp")
	imgs := mediaElm.QuerySelectorAll("img")
	url1 := GetBestImageUrl(imgs)
	videos := mediaElm.QuerySelectorAll("video")
	url2 := GetVideoUrl(videos)

	url := ""
	if url2 == "" {
		url = url1
	} else {
		url = url2
	}
	// send code of post to background for download
	Window.Get("chrome").Get("runtime").Call("sendMessage", url)
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
