package main

import (
	"../bworker"
)

func main() {
	var wk Woker
	bworker.Run(setConf(), &wk)
}
