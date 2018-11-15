
all:clean  csbox rboxtool  fboxtool

csbox:
	go build -ldflags "-s -w" -o $@ ../csbox/*.go

rboxtool:
	go build -ldflags "-s -w" -o $@ ../rboxtool/*.go

fboxtool:
	go build -ldflags "-s -w" -o $@ ../fboxtool/*.go

csboxWebAssembly:
	GOARCH=wasm GOOS=js go build -o ../csbox/static/js/csbox.wasm ../csbox/webassembly/*.go

clean:
	rm -rf ./csbox ./rboxtool ./fboxtool

installrbox:
	echo ${GOPATH}
	rm -rf /soft/golib/src/github.com/rlds/rbox/base
	cp -r ../base /soft/golib/src/github.com/rlds/rbox

installfbox:
	echo ${GOPATH}
	rm -rf /soft/golib/src/github.com/rlds/rbox/fbox
	cp -r ../fbox /soft/golib/src/github.com/rlds/rbox

linux:
	GOOS=linux go build -ldflags "-s -w" -o csbox ../csbox/*.go
