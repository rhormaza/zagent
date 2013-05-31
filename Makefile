GOPATH := $(CURDIR)
export GOPATH

all: log4go util search main

main:
	go build -o zagent main
	mv zagent bin/
	
search:
	go build search

util:
	go build util

log4go:
	go build log4go

tcpserver:
	go build tcpserver
