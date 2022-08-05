package main

import (
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -Inimcache
#cgo CFLAGS: -I../nimbus-eth2/vendor/nimbus-build-system/vendor/Nim-csources-v1/c_code/
#cgo LDFLAGS: -L. -lcb
#include <stdlib.h>
#include "cb.h"

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
	C.setHeaderCallback((C.callback_type)(unsafe.Pointer(C.goCallback_cgo)))

	C.invokeHeaderCallback()
	C.testEcho()

}
