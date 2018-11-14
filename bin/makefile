
all:clean  ShowServer

GetTidLogs:
	go build -ldflags "-s -w" -o $@ ../getlog/*.go

ShowServer:
	go build -ldflags "-s -w" -o $@ ../webserver/*.go

buildbox:
	go build -ldflags "-s -w" -o $@ ../boxtool/*.go

clean:
	rm -rf ./ShowServer ./GetTidLogs ./buildbox

init:
	mkdir ./log

installmbox:
	echo ${GOPATH}
	rm -rf ${GOPATH}/src/rlds/mbox
	cp -r ../mbox ${GOPATH}/src/rlds/

linux:
	GOOS=linux go build -ldflags "-s -w" -o ShowServer ../webserver/*.go
