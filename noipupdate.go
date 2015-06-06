package main

import (
    "fmt"
    "net"
    "log"
    "time"
)

func main() {
    var ip, newIp string
    var waitTime int = 10
    var host string = "fbrx.noip.me"
    for{
        fmt.Println(time.Now())
        newIp = resolveIp(host)
        fmt.Printf("current ip: %v - ip from dns: %v\n", ip, newIp)
        if(ip == newIp){
            fmt.Println("no change detected...")
        }else{
            fmt.Println("deteced change...triggering update")
            ip = updateIp(newIp)
        }
        time.Sleep((time.Duration(waitTime) * time.Second))
    }
}


func resolveIp(host string) string {
    var addrs, err = net.LookupHost(host)
    fmt.Println("looking up host: " + host)
    if(err != nil){
        log.Fatal(err)
    }

    fmt.Println("resolved to " + addrs[0])
    return addrs[0]

}

func updateIp(newIp string) string {
    fmt.Println("...updated...")
    return newIp
}
