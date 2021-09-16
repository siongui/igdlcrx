# cannot use relative path in GOROOT, otherwise 6g not found. For example,
#   export GOROOT=../go  (=> 6g not found)
# it is also not allowed to use relative path in GOPATH
ifndef GITHUB_ACTIONS
export GOROOT=$(realpath ../../paligo/go)
export PATH := $(GOROOT)/bin:$(PATH)
endif

userstory: fmt
	@echo "\033[92mDownload user unexpired stories ...\033[0m"
	@go run userstory.go -clean=false -id=$(id)

cleanstory: fmt
	@echo "\033[92mDownload user unexpired stories using clean account ...\033[0m"
	@go run userstory.go -clean=true -id=$(id)

fmt:
	@echo "\033[92mGo fmt source code...\033[0m"
	@go fmt *.go

modinit:
	go mod init github.com/siongui/igdlcrx/tool

modtidy:
	#go list -m all
	go mod tidy
