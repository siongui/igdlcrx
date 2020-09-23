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

build: fmt
	@echo "\033[92mCompiling Go to JavaScript ...\033[0m"
	[ -d $(CRXDIR) ] || mkdir -p $(CRXDIR)
	cp extension/manifest.json $(CRXDIR)
	cp extension/style.css $(CRXDIR)
	cd extension; gopherjs build background.go chrome.go -o $(CRXDIR)/background.js
	cd extension; gopherjs build content.go chrome.go -o $(CRXDIR)/content.js
	cd extension; gopherjs build contentfb.go chrome.go -o $(CRXDIR)/contentfb.js

pack: build
	cd $(CRXDIR); zip -r extension.zip .

fmt:
	@echo "\033[92mGo fmt source code...\033[0m"
	@go fmt extension/*.go

install:
	@echo "\033[92mInstalling GopherJS ...\033[0m"
	go get -u github.com/gopherjs/gopherjs
	@#echo "\033[92mInstalling GopherJS Bindings for Chrome ...\033[0m"
	@#go get -u github.com/fabioberger/chrome
	@echo "\033[92mInstalling github.com/siongui/godom ...\033[0m"
	go get -u github.com/siongui/godom
	@echo "\033[92mInstalling github.com/siongui/instago ...\033[0m"
	go get -u github.com/siongui/instago
