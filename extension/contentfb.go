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

func SendMessage(username, url string) {
	println("download " + username + " " + url)
	// send code of post to background for download
	Chrome.Runtime.Call("sendMessage", "fbstory:"+username+","+url)
}

func DownloadFacebookStoryButton(username, url string) {
	//println(username + " " + url)
	btns := Document.QuerySelectorAll(".fb-story-dl-button")
	if len(btns) == 0 {
		container, ok := GetElementInElement(Document, "i.hu5pjgll.m6k467ps.sp_F8MGkiW7YGw_1_5x.sx_acc4fa")
		if !ok {
			println("cannot find story button container")
			return
		}
		btn := Document.CreateElement("button")
		btn.ClassList().Add("fb-story-dl-button")
		btn.Dataset().Set("username", username)
		btn.Dataset().Set("url", url)
		btn.SetInnerHTML("Download")
		btn.AddEventListener("click", func(e Event) {
			e.StopPropagation()
			e.PreventDefault()
			SendMessage(e.Target().Dataset().Get("username").String(),
				e.Target().Dataset().Get("url").String())
		})
		container.ParentNode().ParentNode().ParentNode().AppendChild(btn)
	} else {
		btn := btns[0]
		btn.Dataset().Set("username", username)
		btn.Dataset().Set("url", url)
	}
}

func DoFacebookStoryAction(url string) {
	//println("story url: " + url)
	/*
		storyElm, ok := GetElementInElement(Document, "div[data-pagelet='Stories']")
		if !ok {
			println("cannot find story element")
			return
		}
	*/
	storyElm := Document

	// try to find story username
	username := ""
	userElm, ok := GetElementInElement(storyElm, "img.q9iuea42.qs4al1v0.eprw1yos.a4d05b8z.sibfvsnu.px9q9ucb.j2ut9x2k.p4hiznlx.a8c37x1j.qypqp5cg.bixrwtb6.q676j6op")
	if !ok {
		println("cannot find element containing username")
		return
	}
	if userElm.HasAttribute("alt") {
		username = userElm.GetAttribute("alt")
	}
	if username == "" {
		println("cannot get facebook story username")
		return
	}
	//println(username)

	// try to find story image (if exist)
	imgUrl := ""
	imgElm, ok := GetElementInElement(storyElm, "img.g5ia77u1.arfg74bv.n00je7tq.pmk7jnqg.j9ispegn.rk01pc8j.ke6wolob.k4urcfbm.du4w35lb")
	if ok {
		if imgElm.HasAttribute("src") {
			imgUrl = imgElm.GetAttribute("src")
		}
	}
	if imgUrl != "" {
		// download and return
		//println(imgUrl)
		DownloadFacebookStoryButton(username, imgUrl)
		return
	}

	// story image not exist. find story video.
	videoUrl := ""
	videoElm, ok := GetElementInElement(storyElm, "video.k4urcfbm.datstx6m.a8c37x1j")
	if ok {
		if videoElm.HasAttribute("src") {
			videoUrl = videoElm.GetAttribute("src")
		}

		if videoUrl != "" {
			// download and return
			//println(videoUrl)
			DownloadFacebookStoryButton(username, videoUrl)
			return
		}
	}
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
