GO ?= go


.PHONY: all clean

all: server migration

clean:
	$(RM) -rf build/output/*

server: cmd/goblog-server/main.go
	$(GO) build -o build/output/$@ $<

migration: cmd/goblog-utils/main.go
	$(GO) build -o build/output/$@ $<
