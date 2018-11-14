
package main

import(
    "fmt"
    "net"
)

func main(){
    l, _ := net.Listen("tcp", ":0") // listen on localhost
    port := l.Addr().(*net.TCPAddr).Port
    ip := l.Addr().(*net.TCPAddr).IP
    fmt.Println("--",ip, port)
}



