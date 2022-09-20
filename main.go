package main

import (
	"bytes"
	"fmt"
	"os"
	"unsafe"

	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
)

/*
#include <stdlib.h>
#include "lightclient.h"

typedef void (*callback_type)(char *);
void goCallback_cgo(char *);

*/
import "C"

//export goCallback
func goCallback(json *C.char) {
	goStr := C.GoString(json)
	//C.free(unsafe.Pointer(json))
	fmt.Println("### goCallback " + goStr)
}

func main() {
	fmt.Println("vim-go")
	cb := (C.callback_type)(unsafe.Pointer(C.goCallback_cgo))
	C.testEcho()
	C.setOptimisticHeaderCallback(cb)
	C.setFinalizedHeaderCallback(cb)
	fmt.Println("vim-go 2")
	type Config struct {
		Eth2Network      string `toml:"network"`
		TrustedBlockRoot string `toml:"trusted-block-root"`
	}

	var testConfig = Config{
		Eth2Network:      "mainnet",
		TrustedBlockRoot: "0x60d6462ea001c078ec95c4cccb7982e82503b99bd5fd91a5b5ed06a0d736fa6f",
		//Eth2Network:      "prater",
		//TrustedBlockRoot: "0x017e4563ebf7fed67cff819c63d8da397b4ed0452a3bbd7cae13476abc5020e4",
	}

	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(testConfig)
	if err != nil {
		return
	}
	tomlFileName := "config.toml"
	f, err := os.Create(tomlFileName)
	if err != nil {
		return
	}
	//defer f.Close()
	f.WriteString(buffer.String())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		configCStr := C.CString(tomlFileName)
		C.startLc(configCStr)
	}()
	fmt.Println("vim-go 3")

	<-signals
	//C.invokeHeaderCallback()
	//C.testEcho()

}
