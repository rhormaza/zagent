GOPATH := $(CURDIR)
export GOPATH

all: main

main:
	go install zagent

zrouter:
	go build zrouter

zstatus:
	go build zstatus

zsearch:
	go build zsearch

zutil:
	go build zutil

log4go:
	go build log4go

zjson:
	go build zjson

ztcpserver:
	go build ztcpserver


