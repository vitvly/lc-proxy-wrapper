NIMBASE_H_PATH ?= $(CURDIR)/../nimbus-eth2/vendor/nimbus-build-system/vendor/Nim-csources-v1/c_code/
CGO_CFLAGS ?= -I$(CURDIR)/nimcache -I$(NIMBASE_H_PATH)
CGO_LDFLAGS ?= -L. -lcb

light-client-status-go: build-nim
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build

build-nim: 
	nim c --app:staticlib --header:cb.h --noMain:on --nimcache:$(CURDIR)/nimcache cb.nim

.PHONE: clean

clean:
	rm -rf nimcache libcb.a nim-test

