GO ?= go


.PHONY: all clean server utils

all: server utils

clean:
	$(RM) -rf build/output/*

server: cmd/goblog-server/main.go
	export CGO_ENABLED=0 && $(GO) build -o build/output/$@ $<

utils: cmd/goblog-utils/main.go
	export CGO_ENABLED=0 && $(GO) build -o build/output/$@ $<
