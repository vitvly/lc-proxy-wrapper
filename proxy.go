package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"github.com/vitvly/lc-proxy-wrapper/types"
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

var proxyEventChan chan *types.ProxyEvent

var nimContextPtr unsafe.Pointer

//export goCallback
func goCallback(json *C.char, cbType int) {
	//C.free(unsafe.Pointer(json))
	//fmt.Println("### goCallback " + goStr)
	var goStr string
	if json != nil {
		goStr = C.GoString(json)
	}
	if proxyEventChan != nil {
		if cbType == 0 { // finalized header
			proxyEventChan <- &types.ProxyEvent{types.FinalizedHeader, goStr}
		} else if cbType == 1 { // optimistic header
			proxyEventChan <- &types.ProxyEvent{types.OptimisticHeader, goStr}
		} else if cbType == 2 { // stopped
			proxyEventChan <- &types.ProxyEvent{types.Stopped, goStr}
			close(proxyEventChan)
			nimContextPtr = nil
		} else if cbType == 3 { // error
			proxyEventChan <- &types.ProxyEvent{types.Error, goStr}
			close(proxyEventChan)
			nimContextPtr = nil
		}
	}
}

func StartVerifProxy(cfg *Config, proxyEventCh chan *types.ProxyEvent) error {
	if nimContextPtr != nil {
		// Other instance running
		return errors.New("Nimbux proxy already (still) running")
	}
	proxyEventChan = proxyEventCh
	cb := (C.callback_type)(unsafe.Pointer(C.goCallback_cgo))

	jsonBytes, _ := json.Marshal(cfg)
	jsonStr := string(jsonBytes)
	fmt.Println("### jsonStr: ", jsonStr)
	configCStr := C.CString(jsonStr)
	nimContextPtr = unsafe.Pointer(C.startVerifProxy(configCStr, cb))
	fmt.Println("ptr: %p", nimContextPtr)
	fmt.Println("inside go-func after startLcViaJson")

	return nil

}

func StopVerifProxy() {
	C.stopVerifProxy((*C.struct_VerifProxyContext)(nimContextPtr))
}
