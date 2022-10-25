NIMBASE_H_PATH ?= $(CURDIR)/../nimbus-eth1/vendor/nimbus-build-system/vendor/Nim-csources-v1/c_code/
CGO_CFLAGS_TEST ?= -I$(CURDIR)/nimcache -I$(NIMBASE_H_PATH)
CGO_LDFLAGS_TEST ?= -L. -lcb
CGO_CFLAGS ?= -I$(CURDIR)/../nimbus-eth1/nimcache/liblcproxy -I$(NIMBASE_H_PATH)
LCPROXY_OBJS = $(shell find $(CURDIR)/../nimbus-eth1/nimcache/liblcproxy -name "*.o")
LIBNATPMP_OBJS = $(shell find $(CURDIR)/../nimbus-eth1/vendor/nim-nat-traversal/vendor/libnatpmp-upstream -name "*.o")
LIBMINIUPNPC_OBJS = $(shell find $(CURDIR)/../nimbus-eth1/vendor/nim-nat-traversal/vendor/miniupnp/miniupnpc -name "*.o")
LIBBACKTRACE_OBJS = $(shell find $(CURDIR)/../nimbus-eth1/vendor/nim-libbacktrace/vendor/libbacktrace-upstream -name "*.o")
LIBBACKTRACE_WRAPPER_OBJS = $(CURDIR)/../nimbus-eth1/vendor/nim-libbacktrace/libbacktrace_wrapper.o
ALL_OBJS = $(LCPROXY_OBJS) $(LIBNATPMP_OBJS) $(LIBMINIUPNPC_OBJS) $(LIBBACKTRACE_WRAPPER_OBJS) $(LIBBACKTRACE_OBJS)
CGO_LDFLAGS = -v -t -whyload 

build-light-client-go:
	#$(warning $(OBJS))
	$(shell echo ${ALL_OBJS} | tr ' ' '\n' > objs.lst)
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -x -ldflags '-v "-extldflags=-filelist objs.lst"'

build-go: build-nim
	CGO_CFLAGS="$(CGO_CFLAGS_TEST)" CGO_LDFLAGS="$(CGO_LDFLAGS_TEST)" go build

build-lc-proxy-wrapper: 
	nim c --app:staticlib --header:cb.h --noMain:on --nimcache:$(CURDIR)/nimcache cb.nim

light-client-status-go: build-nim

.PHONY: clean

clean:
	rm -rf nimcache libcb.a lc-proxy-wrapper

