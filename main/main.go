package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	proxy "github.com/vitvly/lc-proxy-wrapper"
	types "github.com/vitvly/lc-proxy-wrapper/types"
)

func main() {
	var testConfig = proxy.Config{
		Eth2Network:      "mainnet",
		TrustedBlockRoot: "0xe4350d9d7f53d3558a208cbc4bc93231169082b61f25114ccd6460200b1feabe",
		//Web3Url:          Web3UrlType{"HttpUrl", "https://mainnet.infura.io/v3/800c641949d64d768a5070a1b0511938"},
		Web3Url:    "https://mainnet.infura.io/v3/800c641949d64d768a5070a1b0511938",
		RpcAddress: "127.0.0.1",
		RpcPort:    8545,
		LogLevel:   "INFO",
		//Eth2Network:      "prater",
		//TrustedBlockRoot: "0x017e4563ebf7fed67cff819c63d8da397b4ed0452a3bbd7cae13476abc5020e4",
	}
	ctx, cancel := context.WithCancel(context.Background())
	proxyEventCh := make(chan *types.ProxyEvent, 10)
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
