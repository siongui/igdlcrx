package main

import (
	"github.com/fabioberger/chrome"
)

func main() {
	c := chrome.NewChrome()

	// Currently do nothing meaningful
	c.Tabs.OnUpdated(func(tabId int, changeInfo chrome.Object, tab chrome.Tab) {
		if _, ok := changeInfo["url"]; ok {
			url := changeInfo["url"].(string)

			queryInfo := make(map[string]interface{})
			queryInfo["active"] = true
			queryInfo["currentWindow"] = true

			c.Tabs.Query(queryInfo, func(tabs []chrome.Tab) {
				if len(tabs) == 1 {
					c.Tabs.SendMessage(tabs[0].Id, url, nil)
				}
			})
		}
	})

	// Receive code of post from content.
	// Call chrome.downloads API to download files
	c.Runtime.OnMessage(func(message interface{}, sender chrome.MessageSender, sendResponse func(interface{})) {
		msg := message.(string)
		println("Received msg from content: " + msg)
	})
}
