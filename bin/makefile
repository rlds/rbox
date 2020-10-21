
all:clean  csbox rboxtool  fboxtool boxapi

csbox:
	go build -ldflags "-s -w" -o $@ ../csbox/*.go

rboxtool:
	go build -ldflags "-s -w" -o $@ ../rboxtool/*.go

fboxtool:
	go build -ldflags "-s -w" -o $@ ../fboxtool/*.go

boxapi:
	go build -ldflags "-s -w" -o $@ ../boxapi/*.go

csboxWebAssembly:
	GOARCH=wasm GOOS=js go build -o ../csbox/static/js/csbox.wasm ../csbox/webassembly/*.go

clean:
	rm -rf ./csbox ./rboxtool ./fboxtool ./boxapi

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
	GOOS=linux go build -ldflags "-s -w" -o boxapi ../boxapi/*.go
