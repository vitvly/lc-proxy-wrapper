include Makefile.vars

CGO_CFLAGS_TEST ?= -I$(CURDIR)/nimcache -I$(NIMBASE_H_PATH)
CGO_LDFLAGS_TEST ?= -L. -lcb

VERIFPROXY_OBJS = $(shell find $(NIMBUS_ETH1_PATH)/nimcache/libverifproxy -name "*.o")
LIBNATPMP_OBJS = $(shell find $(NIMBUS_ETH1_PATH)/vendor/nim-nat-traversal/vendor/libnatpmp-upstream -name "*.o")
LIBMINIUPNPC_OBJS = $(shell find $(NIMBUS_ETH1_PATH)/vendor/nim-nat-traversal/vendor/miniupnp/miniupnpc -name "*.o")
LIBBACKTRACE_OBJS = $(shell find $(NIMBUS_ETH1_PATH)/vendor/nim-libbacktrace/vendor/libbacktrace-upstream -name "*.o")
LIBBACKTRACE_WRAPPER_OBJS = $(NIMBUS_ETH1_PATH)/vendor/nim-libbacktrace/libbacktrace_wrapper.o
#ALL_OBJS = $(VERIFPROXY_OBJS) $(LIBNATPMP_OBJS) $(LIBMINIUPNPC_OBJS) $(LIBBACKTRACE_WRAPPER_OBJS) $(LIBBACKTRACE_OBJS)

build-light-client-go:
	#$(warning $(OBJS))
	$(shell echo ${ALL_OBJS} | tr ' ' '\n' > objs.lst)
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -x -ldflags $(LDFLAGS)

build-light-client-go-exe:
	#$(warning $(OBJS))
	$(shell echo ${ALL_OBJS} | tr ' ' '\n' > objs.lst)
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -x -ldflags $(LDFLAGS) -o verif-proxy-wrapper ./main 

build-verif-proxy-wrapper:
	CGO_CFLAGS="$(CGO_CFLAGS)" go build -x -v -ldflags '-v "-extldflags=-Wl,-rpath,. -lverifproxy -L$(VERIF_PROXY_OUT_PATH)"' 

build-verif-proxy-wrapper-exe:
	CGO_CFLAGS="$(CGO_CFLAGS)" go build -x -v -ldflags '-v "-extldflags=-Wl,-rpath,. -lverifproxy -L$(VERIF_PROXY_OUT_PATH)"' -o verif-proxy-wrapper ./main 

build-go: build-nim
	CGO_CFLAGS="$(CGO_CFLAGS_TEST)" CGO_LDFLAGS="$(CGO_LDFLAGS_TEST)" go build

# build-verif-proxy-wrapper: 
# 	nim c --app:staticlib --header:cb.h --noMain:on --nimcache:$(CURDIR)/nimcache cb.nim

light-client-status-go: build-nim

.PHONY: clean

clean:
	rm -rf nimcache libcb.a verif-proxy-wrapper

