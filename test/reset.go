package main

import(
	"fmt"
)

func reSet(src *string,dest string,deft string){
	if dest != "" {
		*src = dest
	}else{
		*src = deft
	}
}

func main(){
	var src,dest string
	dest = "dest"
	deft := "deft"

	reSet(&src,dest,deft)
	fmt.Println(src)
}