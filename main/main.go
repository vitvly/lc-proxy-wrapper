package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	proxy "github.com/vitvly/lc-proxy-wrapper"
)

func main() {
	var testConfig = proxy.Config{
		Eth2Network:      "mainnet",
		TrustedBlockRoot: "0x0216f0250965ceb8a54d5220d27ab0776b51695edf85e987f9ea012a5b6f6f40",
		//Web3Url:          Web3UrlType{"HttpUrl", "https://mainnet.infura.io/v3/800c641949d64d768a5070a1b0511938"},
		Web3Url:    "https://mainnet.infura.io/v3/800c641949d64d768a5070a1b0511938",
		RpcAddress: "127.0.0.1",
		RpcPort:    8545,
		LogLevel:   "INFO",
		//Eth2Network:      "prater",
		//TrustedBlockRoot: "0x017e4563ebf7fed67cff819c63d8da397b4ed0452a3bbd7cae13476abc5020e4",
	}
	proxyEventCh, err := proxy.StartVerifProxy(&testConfig)
	if err != nil {
		fmt.Println("Error when starting proxy:", err)
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Before range signals")
	for {
		select {
		case <-signals:
			fmt.Println("Signal caught, exiting")
			break
		case ev := <-proxyEventCh:
			fmt.Println("event caught", ev.EventType)

		}
	}

}
