# cannot use relative path in GOROOT, otherwise 6g not found. For example,
#   export GOROOT=../go  (=> 6g not found)
# it is also not allowed to use relative path in GOPATH
export GOROOT=$(realpath ../../go)
export GOPATH=$(realpath ..)
export PATH := $(GOROOT)/bin:$(GOPATH)/bin:$(PATH)

IGDIR=../Instagram

### example Makefile for simonedegiacomi/gphotosuploader
### Don't forget to prepare auth file of gphotosuploader
upload:
	while [ true ] ; do \
		gphotosuploader --upload ${IGDIR} --maxConcurrent 10 ; \
		sleep 10 ; \
	done

zero:
	go run zero.go -root=${IGDIR}

install_gphotosuploader:
	go get github.com/simonedegiacomi/gphotosuploader
