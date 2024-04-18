package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/matthiasbruns/gnutella-go/cache"
	"github.com/matthiasbruns/gnutella-go/gnutella"
	"github.com/matthiasbruns/gnutella-go/whatsmyip"
)

var ipCache []*net.TCPAddr
var remoteAddr *net.TCPAddr
var listenPort uint16 = 45309

func init() {
	fmt.Println("Starting Gnutella CLI")
}

func main() {
	ctx := context.Background()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ignore := false
	whatsmyip.DiscoverPublicIP(func(ip string, err error) {
		if ignore {
			return
		}

		ignore = true

		defer wg.Done()
		fmt.Println("Public IP:", ip)

		remoteAddr = &net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: int(listenPort),
		}
	})
	wg.Wait()

	// get ips
	ipCache = cache.QueryBazookaVendors(ctx)
	fmt.Println("IPs in cache:", len(ipCache))

	listeningChan := make(chan *net.TCPAddr)
	// get remote ip
	go func() {
		err := gnutella.StartListen("", listenPort, listeningChan)
		if err != nil {
			panic(err)
		}
	}()

	// try handshake until success
	for _, ip := range ipCache {
		successChan := make(chan bool)

		_, err := gnutella.InitTCPHandshake(ctx, remoteAddr, ip, successChan)
		if err != nil {
			fmt.Println("Error while initiating handshake with", ip, err)
			continue
		}

		success := <-successChan

		if !success {
			fmt.Println("Handshake failed with", ip)
		}
	}

	for {
	}
}
