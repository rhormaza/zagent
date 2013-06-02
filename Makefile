GOPATH := $(CURDIR)
export GOPATH

all: log4go zutil zjson zsearchlog zrouter main

main:
	go build -o zagent main
	mv zagent bin/
	
zsearchlog:
	go build zsearchlog

zutil:
	go build zutil

log4go:
	go build log4go

zjson:
	go build zjson

zrouter:
	go build zrouter

ztcpserver:
	go build ztcpserver
