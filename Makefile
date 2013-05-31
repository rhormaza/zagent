GOPATH := $(CURDIR)
export GOPATH

all: agent
	
agent:
	go install search


util:
	go install util

