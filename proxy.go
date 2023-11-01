package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

/*
#include <stdlib.h>
#include "verifproxy.h"

typedef void (*callback_type)(char *, int);
void goCallback_cgo(char *, int);

*/
import "C"

type Web3UrlType struct {
	Kind    string `toml:"kind"`
	Web3Url string `toml:"web3Url"`
}
type Config struct {
	Eth2Network      string `toml:"network"`
	TrustedBlockRoot string `toml:"trusted-block-root"`
	// Web3Url          Web3UrlType `toml:"web3-url"`
	Web3Url    string `toml:"web3-url"`
	RpcAddress string `toml:"rpc-address"`
	RpcPort    uint16 `toml:"rpc-port"`
	LogLevel   string `toml:"log-level"`
}

const (
	ProxyInitialised int = iota
	OptimisticHeader
	FinalizedHeader
)

type ProxyEvent struct {
	EventType int
	Msg       string
}

var proxyEventChan chan *ProxyEvent
var once sync.Once

//export goCallback
func goCallback(json *C.char, cbType int) {
	//C.free(unsafe.Pointer(json))
	//fmt.Println("### goCallback " + goStr)
	if proxyEventChan != nil {
		once.Do(func() {
			proxyEventChan <- &ProxyEvent{ProxyInitialised, ""}
		})
		goStr := C.GoString(json)
		if cbType == 0 { // finalized header
			proxyEventChan <- &ProxyEvent{FinalizedHeader, goStr}
		} else if cbType == 1 { // optimistic header
			proxyEventChan <- &ProxyEvent{OptimisticHeader, goStr}
		}
	}
}

var nimContextPtr unsafe.Pointer

func StartLightClient(ctx context.Context, cfg *Config, proxyEventCh chan *ProxyEvent) {
	proxyEventChan = proxyEventCh
	cb := (C.callback_type)(unsafe.Pointer(C.goCallback_cgo))

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		jsonBytes, _ := json.Marshal(cfg)
		jsonStr := string(jsonBytes)
		fmt.Println("### jsonStr: ", jsonStr)
		configCStr := C.CString(jsonStr)
		nimContextPtr = unsafe.Pointer(C.startLightClientProxy(configCStr, cb))
		fmt.Println("ptr: %p", nimContextPtr)
		fmt.Println("inside go-func after startLcViaJson")
	}()

}
