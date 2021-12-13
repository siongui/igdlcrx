# cannot use relative path in GOROOT, otherwise 6g not found. For example,
#   export GOROOT=../go  (=> 6g not found)
# it is also not allowed to use relative path in GOPATH
ifndef GOROOT
export GOROOT=$(realpath ../go)
export GOPATH=$(CURDIR)/gphoto
export PATH := $(GOROOT)/bin:$(GOPATH)/bin:$(PATH)
endif

#export GO111MODULE=off
IGDIR=Instagram
CRXDIR=$(CURDIR)/crx
ZIPFILE=$(CRXDIR)/extension.zip

build: fmt
	@echo "\033[92mCompiling Go to JavaScript ...\033[0m"
	[ -d $(CRXDIR) ] || mkdir -p $(CRXDIR)
	cp extension/manifest.json $(CRXDIR)
	cp extension/style.css $(CRXDIR)
	cp extension/request.js $(CRXDIR)
	cd extension; gopherjs build background.go chrome.go -o $(CRXDIR)/background.js
	cd extension; gopherjs build content.go chrome.go -o $(CRXDIR)/content.js
	cd extension; gopherjs build contentfb.go chrome.go -o $(CRXDIR)/contentfb.js

pack: build
	cd $(CRXDIR); zip -r extension.zip .

localhost: fmt
	@echo "\033[92mlocalhost Server Running ...\033[0m"
	@go run localhost/server.go

2layer: fmt
	@echo "\033[92mDownload user $(id) unexpired stories and stories of reel mentions...\033[0m"
	@go run tool/userstory2layer.go -id=$(id)

allpost: fmt
	@echo "\033[92mDownload all user posts ...\033[0m"
	@go run tool/allposts.go -id=$(id)

txtdl: fmt
	@echo "\033[92mDownload unexpired stories in users.txt ...\033[0m"
	@go run tool/txtdl.go -f=users.txt

dltxt: fmt
	@echo "\033[92mDownload unexpired stories in users.txt ...\033[0m"
	@go run tool/dltxt.go -f=users.txt

reelstray: fmt
	@echo "\033[92mDownload stories in reels tray...\033[0m"
	@go run tool/download.go -downloadtype=story -outputdir=${IGDIR}

reelstray2: fmt
	@echo "\033[92mDownload stories in reels tray...\033[0m"
	@go run tool/download.go -downloadtype=story2 -outputdir=${IGDIR}

allunexpired: fmt
	@echo "\033[92mDownload all unexpired stories ...\033[0m"
	@go run tool/unexpired.go

userstory: fmt
	@echo "\033[92mDownload user unexpired stories ...\033[0m"
	@go run tool/userstory.go -clean=false -id=$(id)

cleanstory: fmt
	@echo "\033[92mDownload user unexpired stories using clean account ...\033[0m"
	@go run tool/userstory.go -clean=true -id=$(id)

1dayrm:
	@echo "\033[92mRemove files older than one day ...\033[0m"
	@go run tool/1dayrm.go -root=${IGDIR} -todir=${HOME}/Pictures/

rmemptydir: fmt
	@echo "\033[92mRemove empty dir ...\033[0m"
	@go run tool/rmemptydir.go -root=${IGDIR}

savedcollection: fmt
	@echo "\033[92mDownload saved posts in collection ...\033[0m"
	@go run tool/savedcollection.go -collection=Download

fmt:
	@echo "\033[92mGo fmt source code...\033[0m"
	@go fmt extension/*.go
	@go fmt extension/libbackground/*.go
	@go fmt localhost/*.go
	@go fmt tool/*.go

modinit:
	go mod init github.com/siongui/igdlcrx

modtidy:
	#go list -m all
	go mod tidy

install_gopherjs:
	@echo "\033[92mInstalling GopherJS ...\033[0m"
	go get -u github.com/gopherjs/gopherjs
