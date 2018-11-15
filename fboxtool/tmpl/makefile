
all:{{.BoxConf.Name}}

{{.BoxConf.Name}}:
	go build -ldflags "-s -w" -o $@ ./*.go

clean:
	rm -rf ./{{.BoxConf.Name}}
