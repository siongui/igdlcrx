# cannot use relative path in GOROOT, otherwise 6g not found. For example,
#   export GOROOT=../go  (=> 6g not found)
# it is also not allowed to use relative path in GOPATH
ifndef TRAVIS
export GOROOT=$(realpath ../paligo/go)
export GOPATH=$(realpath ../paligo)
export PATH := $(GOROOT)/bin:$(GOPATH)/bin:$(PATH)
endif

CRXDIR=$(CURDIR)/crx
ZIPFILE=$(CRXDIR)/extension.zip
LOCAL_LIBBGDIR=extension/libbackground
LIBBGDIR=$(GOPATH)/src/github.com/siongui/igdlcrx/$(LOCAL_LIBBGDIR)

build: fmt
	@echo "\033[92mCompiling Go to JavaScript ...\033[0m"
	[ -d $(CRXDIR) ] || mkdir -p $(CRXDIR)
	cp extension/manifest.json $(CRXDIR)
	cp extension/style.css $(CRXDIR)
	cp extension/request.js $(CRXDIR)
	cd extension; gopherjs build background.go chrome.go -o $(CRXDIR)/background.js
	cd extension; gopherjs build content.go chrome.go -o $(CRXDIR)/content.js
	cd extension; gopherjs build contentfb.go chrome.go -o $(CRXDIR)/contentfb.js

copylocallib:
	[ -d $(LIBBGDIR) ] || mkdir -p $(LIBBGDIR)
	cp $(LOCAL_LIBBGDIR)/*.go $(LIBBGDIR)

pack: build
	cd $(CRXDIR); zip -r extension.zip .

localhost: fmt
	@echo "\033[92mlocalhost Server Running ...\033[0m"
	@go run localhost/server.go

userstory2layer: fmt
	@echo "\033[92mDownload user $(id) unexpired stories and stories of reel mentions...\033[0m"
	@go run localhost/userstory2layer.go -id=$(id)

fmt:
	@echo "\033[92mGo fmt source code...\033[0m"
	@go fmt extension/*.go
	@go fmt $(LOCAL_LIBBGDIR)/*.go
	@go fmt localhost/*.go

install:
	@echo "\033[92mInstalling GopherJS ...\033[0m"
	go get -u github.com/gopherjs/gopherjs
	@#echo "\033[92mInstalling GopherJS Bindings for Chrome ...\033[0m"
	@#go get -u github.com/fabioberger/chrome
	@echo "\033[92mInstalling github.com/siongui/godom ...\033[0m"
	go get -u github.com/siongui/godom
	@echo "\033[92mInstalling github.com/siongui/instago ...\033[0m"
	go get -u github.com/siongui/instago
	go get -u github.com/siongui/instago/download
	@echo "\033[92mInstalling github.com/extension/libbackground ...\033[0m"
	go get -u github.com/siongui/igdlcrx/extension/libbackground
