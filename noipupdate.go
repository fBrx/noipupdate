package main

import (
    "fmt"
    "net"
)

func main() {
    resolveIp("fbrx.noip.me")
}


func resolveIp(host string) string {
    var addrs, err = net.LookupHost(host)
    fmt.Println("looking up host: " + host)
    if(err == nil){
        fmt.Println("resolved to " + addrs[0])
        return addrs[0]
    }else{
        fmt.Println(err)
        return "";
    }
}
