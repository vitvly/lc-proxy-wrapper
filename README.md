# Build Nim
nim c --app:staticlib --header:cb.h --noMain:on --nimcache:$HOME/c/nim-test/nimcache cb.nim

# Build go
go build
./nim-test
