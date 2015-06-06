package main

import (
    "fmt"
    "net"
    "net/http"
    "log"
    "time"
    "io/ioutil"
)

var waitTime int = 10
var host string = "fbrx.noip.me"
var ipResolverHost string = "http://whatismyip.akamai.com/" //http://v4.ident.me/

func main() {
    var currentIp, dnsIp string
    for{
        fmt.Println(time.Now())
        currentIp = determineCurrentIp()
        dnsIp = resolveIpFromDns(host)
        fmt.Printf("current ip: %v - ip from dns: %v\n", currentIp, dnsIp)
        if(currentIp == dnsIp){
            fmt.Println("no change detected...")
        }else{
            fmt.Println("deteced change...triggering update")
            updateIp(currentIp)
        }
        time.Sleep((time.Duration(waitTime) * time.Second))
    }
}


func resolveIpFromDns(host string) string {
    var addrs, err = net.LookupHost(host)
    fmt.Println("looking up host: " + host)
    if(err != nil){
        log.Fatal(err)
    }

    fmt.Println("resolved to " + addrs[0])
    return addrs[0]
}

func determineCurrentIp() string {
    resp, err := http.Get(ipResolverHost)
    if(err != nil){
        log.Fatal(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if(err != nil){
        log.Fatal(err)
    }
    fmt.Println("determined current ip: " + string(body))

    return string(body)
}

func updateIp(newIp string) string {
    fmt.Printf("...updated to %v...\n")
    return newIp
}
