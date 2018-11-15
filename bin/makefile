
all:clean  ShowServer

csbox:
	go build -ldflags "-s -w" -o $@ ../csbox/*.go

rboxtool:
	go build -ldflags "-s -w" -o $@ ../rboxtool/*.go

fboxtool:
	go build -ldflags "-s -w" -o $@ ../fboxtool/*.go


clean:
	rm -rf ./csbox ./rboxtool ./fboxtool

init:
	mkdir ./log

installrbox:
	echo ${GOPATH}
	rm -rf /soft/golib/src/github.com/rlds/rbox/base
	cp -r ../base /soft/golib/src/github.com/rlds/rbox

installfbox:
	echo ${GOPATH}
	rm -rf /soft/golib/src/github.com/rlds/rbox/fbox
	cp -r ../fbox /soft/golib/src/github.com/rlds/rbox

linux:
	GOOS=linux go build -ldflags "-s -w" -o ShowServer ../showserver/*.go
