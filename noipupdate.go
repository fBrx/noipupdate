package main

import (
    "fmt"
    "net"
    "net/http"
    "log"
    "time"
    "io/ioutil"
    "flag"
)

var ipResolver string
var username, password, host, checkHost string
var interval int
var verbose bool

func initArgs() {
    flag.StringVar(&host, "host", "", "the no-ip hostname to update")
    flag.StringVar(&checkHost, "checkHost", "", "the hostname to check for the current dns entry. will default to host if not set")
    flag.StringVar(&username, "username", "", "the username of your no-ip account")
    flag.StringVar(&password, "password", "", "the password of your no-ip account")
    flag.IntVar(&interval, "interval", 300, "the interval (in seconds) in which to perform update checks")
    flag.BoolVar(&verbose, "v", false, "enable verbose logging")
    flag.StringVar(&ipResolver, "ipResolver", "http://whatismyip.akamai.com/", "the url to check for the current ip. response must only contain ip as plain text.")

    flag.Parse()

    if(checkHost == ""){
        checkHost = host
    }

}

func main() {
    initArgs()

    var currentIp, dnsIp string
    for{
        fmt.Println(time.Now())
        currentIp = determineCurrentIp()
        dnsIp = resolveIpFromDns(checkHost)
        if(currentIp == dnsIp){
            fmt.Printf("no change detected...going to sleep for %v seconds\n", interval)
        }else{
            fmt.Println("deteced change...triggering update")
            updateIp(currentIp)
        }
        time.Sleep((time.Duration(interval) * time.Second))
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
    curIp := callUrlAndGetResponse(ipResolver)
    fmt.Printf("determined current ip from %v: %v\n", ipResolver, curIp)
    return curIp
}

func updateIp(newIp string) string {
    url := "http://" + username + ":" + password + "@dynupdate.no-ip.com/nic/update?hostname=" + host + "&myip=" + newIp
    fmt.Printf("updating to %v using url %v...\n", newIp, url)
//    resp := callUrlAndGetResponse(url)
//    fmt.Println("received response: " + resp)
    return newIp
}

func callUrlAndGetResponse(url string) string{
    resp, err := http.Get(url)
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
