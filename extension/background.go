package main

func main() {
	// Currently do nothing meaningful
	Chrome.Tabs.Get("onUpdated").Call("addListener", func(tabId int, changeInfo map[string]interface{}) {
		if _, ok := changeInfo["url"]; !ok {
			return
		}

		url := changeInfo["url"].(string)

		queryInfo := make(map[string]interface{})
		queryInfo["active"] = true
		queryInfo["currentWindow"] = true

		Chrome.Tabs.Call("query", queryInfo, func(tabs []map[string]interface{}) {
			if len(tabs) == 1 {
				Chrome.Tabs.Call("sendMessage", tabs[0]["id"], url)
			}
		})
	})

	// Receive code of post from content.
	// Call chrome.downloads API to download files
	Chrome.Runtime.Get("onMessage").Call("addListener", func(message interface{}) {
		msg := message.(string)
		println("Received msg from content: " + msg)
	})
}
