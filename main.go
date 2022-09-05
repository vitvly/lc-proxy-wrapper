package main

import (
	"bytes"
	"fmt"
	"os"

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
	//cb := (C.callback_type)(unsafe.Pointer(C.goCallback_cgo))
	C.testEcho()
	// C.setOptimisticHeaderCallback(cb)
	// C.setFinalizedHeaderCallback(cb)
	fmt.Println("vim-go 2")
	type Config struct {
		Network          string
		TrustedBlockRoot string
	}

	var testConfig = Config{
		Network:          "prater",
		TrustedBlockRoot: "0x017e4563ebf7fed67cff819c63d8da397b4ed0452a3bbd7cae13476abc5020e4",
	}

	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(testConfig)
	if err != nil {
		return
	}
	var configStr = buffer.String()
	configCStr := C.CString(configStr)
	fmt.Println("configStr: ", configStr)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		C.startLc(configCStr)
	}()
	fmt.Println("vim-go 3")

	<-signals
	//C.invokeHeaderCallback()
	//C.testEcho()

}
