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
var username, password, noipHost string

func main() {
    username = "user"
    password = "pass"
    noipHost = "host"

    var currentIp, dnsIp string
    for{
        fmt.Println(time.Now())
        currentIp = determineCurrentIp()
        dnsIp = resolveIpFromDns(host)
        if(currentIp == dnsIp){
            fmt.Printf("no change detected...going to sleep for %v seconds\n", waitTime)
        }else{
            fmt.Println("deteced change...triggering update")
            updateIp(currentIp)
        }
        time.Sleep((time.Duration(waitTime) * time.Second))
    }
}


func resolveIpFromDns(host string) string {
    var addrs, err = net.LookupHost(host)
    if(err != nil){
        log.Fatal(err)
    }

    fmt.Printf("resolved host %v to ip: %v\n", host, addrs[0])
    return addrs[0]
}

func determineCurrentIp() string {
    curIp := callUrlAndGetResponse(ipResolverHost)
    fmt.Printf("determined current ip from %v: %v\n", ipResolverHost, curIp)
    return curIp
}

func updateIp(newIp string) string {
    url := "http://" + username + ":" + password + "@dynupdate.no-ip.com/nic/update?hostname=" + noipHost + "&myip=" + newIp
    fmt.Printf("...updated to %v using url %v...\n", newIp, url)
    return newIp
}

func callUrlAndGetResponse(url string) string{
    resp, err := http.Get(ipResolverHost)
    if(err != nil){
        log.Fatal(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if(err != nil){
        log.Fatal(err)
    }

    return string(body)
}
