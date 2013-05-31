GOPATH := $(CURDIR)
export GOPATH

all: log4go util search
	
search:
	go install search

util:
	go build util

log4go:
	go build log4go
