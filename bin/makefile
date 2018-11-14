
all:clean  ShowServer

ShowServer:
	go build -ldflags "-s -w" -o $@ ../webserver/*.go

buildbox:
	go build -ldflags "-s -w" -o $@ ../boxtool/*.go

clean:
	rm -rf ./ShowServer ./buildbox

init:
	mkdir ./log

installmbox:
	echo ${GOPATH}
	rm -rf ${GOPATH}/src/github.com/rlds/rbox/base
	cp -r ../base ${GOPATH}/src/github.com/rlds/rbox

linux:
	GOOS=linux go build -ldflags "-s -w" -o ShowServer ../webserver/*.go
