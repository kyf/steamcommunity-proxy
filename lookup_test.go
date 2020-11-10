package main

import(
    "testing"
    "os"
    "fmt"
    "log"
)

func TestDnsLookup(t *testing.T){
    logger := log.New(
        os.Stdout,
        "[dns]",
        log.LstdFlags,
    )
    ip, err := DnsLookUp(domainName, dnsList)
    if err != nil {
        t.Errorf("%s", err.Error())
        return
    }
    println(fmt.Sprintf("%v", ip))
}
