GOPATH := $(CURDIR)
export GOPATH

all: main

main:
	go install zagent
	
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


