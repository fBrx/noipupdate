package main

import (
    "fmt"
//    "net"
    "net/http"
    "log"
    "time"
    "io/ioutil"
    "flag"
    "strings"
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
    var currentIp, lastIp string
    currentIp = "not set"
    lastIp = "not set"
    for{
        fmt.Println(time.Now())
        lastIp = currentIp
        currentIp = determineCurrentIp()
        if(currentIp == lastIp){
            fmt.Printf("no change detected (current ip is %v)...going to sleep for %v seconds\n", currentIp, interval)
        }else{
            fmt.Printf("detected change (%v to %v)...triggering update\n", lastIp, currentIp)
            updateIp(currentIp)
        }
        time.Sleep((time.Duration(interval) * time.Second))
    }
}

func determineCurrentIp() string {
    curIp := callUrlAndGetResponse(ipResolver)
    fmt.Printf("determined current ip from %v: %v\n", ipResolver, curIp)
    return curIp
}

func updateIp(newIp string) {
    url := "http://" + username + ":" + password + "@dynupdate.no-ip.com/nic/update?hostname=" + host + "&myip=" + newIp
    fmt.Printf("updating to %v using url %v...\n", newIp, url)
    resp := callUrlAndGetResponse(url)
    fmt.Println("received response: " + resp)
    if(strings.HasPrefix(resp, "good") || strings.HasPrefix(resp, "nochg")){
        fmt.Printf("response was goot...continueing after %v seconds\n", interval)
    }else{
        log.Fatal("response indicated error, please fix")
    }
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
