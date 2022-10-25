package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"unsafe"

	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
)

/*
#include <stdlib.h>
#include "lcproxy.h"

typedef void (*callback_type)(char *);
void goCallback_cgo(char *);

*/
import "C"

type Web3UrlType struct {
	Kind    string `toml:"kind"`
	Web3Url string `toml:"web3Url"`
}
type Config struct {
	Eth2Network      string      `toml:"network"`
	TrustedBlockRoot string      `toml:"trusted-block-root"`
	Web3Url          Web3UrlType `toml:"web3-url"`
	RpcAddress       string      `toml:"rpc-address"`
	RpcPort          uint16      `toml:"rpc-port"`
}

type BeaconBlockHeader struct {
	Slot          uint64 `json:"slot"`
	ProposerIndex uint64 `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
}

//export goCallback
func goCallback(json *C.char) {
	goStr := C.GoString(json)
	//C.free(unsafe.Pointer(json))
	fmt.Println("### goCallback " + goStr)
	// var hdr BeaconBlockHeader
	// err := json.NewDecoder([]byte(goStr)).Decode(&hdr)
	// if err != nil {
	// 	fmt.Println("### goCallback json parse error: " + err)
	// }
	// fmt.Println("Unmarshal result: " + hdr)
}

func StartLightClient(ctx context.Context, cfg *Config) {
	fmt.Println("vim-go")
	cb := (C.callback_type)(unsafe.Pointer(C.goCallback_cgo))
	C.testEcho()
	C.setOptimisticHeaderCallback(cb)
	C.setFinalizedHeaderCallback(cb)
	fmt.Println("vim-go 2")
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(cfg)
	if err != nil {
		return
	}
	tomlFileName := "config.toml"
	f, err := os.Create(tomlFileName)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(buffer.String())

	go func() {
		configCStr := C.CString(tomlFileName)
		C.startLc(configCStr)
		fmt.Println("inside go-func after startLc")
	}()
	go func() {
		fmt.Println("Before range ctx.Done()")
		for range ctx.Done() {
			fmt.Println("inside go-func ctx.Done()")
			C.quit()
		}
	}()
	fmt.Println("vim-go 3")

}

func main() {
	var testConfig = Config{
		Eth2Network:      "mainnet",
		TrustedBlockRoot: "0x6fd2f9fdc616a0755f3f88ac58e4a7871788c3128261a18bf9645eee7042eb53",
		Web3Url:          Web3UrlType{"HttpUrl", "https://mainnet.infura.io/v3/800c641949d64d768a5070a1b0511938"},
		RpcAddress:       "127.0.0.1",
		RpcPort:          8545,
		//Eth2Network:      "prater",
		//TrustedBlockRoot: "0x017e4563ebf7fed67cff819c63d8da397b4ed0452a3bbd7cae13476abc5020e4",
	}
	ctx, cancel := context.WithCancel(context.Background())
	StartLightClient(ctx, &testConfig)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Before range signals")
	// time.Sleep(8 * time.Second)
	// cancel()
	for range signals {
		fmt.Println("Signal caught, exiting")
		cancel()
	}
	fmt.Println("Exiting")

}
