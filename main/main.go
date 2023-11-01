package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	proxy "github.com/vitvly/lc-proxy-wrapper"
)

func main() {
	var testConfig = proxy.Config{
		Eth2Network:      "mainnet",
		TrustedBlockRoot: "0x9a6364b5702e5854172e175c9236a743743e8a9f139d8e1ec0a668eebc0f176b",
		//Web3Url:          Web3UrlType{"HttpUrl", "https://mainnet.infura.io/v3/800c641949d64d768a5070a1b0511938"},
		Web3Url:    "https://mainnet.infura.io/v3/800c641949d64d768a5070a1b0511938",
		RpcAddress: "127.0.0.1",
		RpcPort:    8545,
		LogLevel:   "INFO",
		//Eth2Network:      "prater",
		//TrustedBlockRoot: "0x017e4563ebf7fed67cff819c63d8da397b4ed0452a3bbd7cae13476abc5020e4",
	}
	ctx, cancel := context.WithCancel(context.Background())
	proxyEventCh := make(chan *proxy.ProxyEvent, 10)
	proxy.StartLightClient(ctx, &testConfig, proxyEventCh)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Before range signals")
	for {
		select {
		case <-signals:
			fmt.Println("Signal caught, exiting")
			cancel()
			break
		case ev := <-proxyEventCh:
			fmt.Println("event caught", ev.EventType)

		}
	}

}
