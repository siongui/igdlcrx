package main

import (
	"fmt"
	"log"
	"net/http"
)

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/alive/" {
		fmt.Fprintf(w, "ok")
	}
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "url path: %s!", r.URL.Path)
}

func main() {
	http.HandleFunc("/alive/", aliveHandler)
	http.HandleFunc("/stories/", storyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
