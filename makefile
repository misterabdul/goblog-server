GO ?= go


.PHONY: all clean

all: server

clean:
	$(RM) -rf build/output/*

server: cmd/goblog-server/main.go
	$(GO) build -o build/output/$@ $<
