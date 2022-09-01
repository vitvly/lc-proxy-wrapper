package main

import (
	"fmt"
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
	C.startLightClient()
	fmt.Println("vim-go 3")

	//C.invokeHeaderCallback()
	//C.testEcho()

}
