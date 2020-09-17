package main

import (
	"github.com/fabioberger/chrome"
)

func main() {
	c := chrome.NewChrome()
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
}
