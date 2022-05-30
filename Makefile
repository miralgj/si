BINDIR      := $(CURDIR)/bin
BINNAME     ?= si

SHELL      = /usr/bin/env bash

GOBIN         = $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN         = $(shell go env GOPATH)/bin
endif

.PHONY: build
build: build-linux build-linux-armv7 build-linux-mips build-mac build-mac-arm build-windows

.PHONY: build-windows
build-windows: export GOOS=windows
build-windows: export GOARCH=amd64
build-windows: export GO111MODULE=on
build-windows: export GOPROXY=$(MOD_PROXY_URL)
build-windows:
	go build -v -o bin/$(BINNAME)-windows-amd64 cmd/si/main.go

.PHONY: build-linux
build-linux: export GOOS=linux
build-linux: export GOARCH=amd64
build-linux: export CGO_ENABLED=0
build-linux: export GO111MODULE=on
build-linux: export GOPROXY=$(MOD_PROXY_URL)
build-linux:
	go build -v -o bin/$(BINNAME)-linux-amd64 cmd/si/main.go

.PHONY: build-linux-mips
build-linux-mips: export GOOS=linux
build-linux-mips: export GOARCH=mips64le
build-linux-mips: export CGO_ENABLED=0
build-linux-mips: export GO111MODULE=on
build-linux-mips: export GOPROXY=$(MOD_PROXY_URL)
build-linux-mips:
	go build -v -o bin/$(BINNAME)-linux-mips64 cmd/si/main.go

.PHONY: build-linux-armv7
build-linux-armv7: export GOOS=linux
build-linux-armv7: export GOARCH=arm
build-linux-armv7: export GOARM=7
build-linux-armv7: export CGO_ENABLED=0
build-linux-armv7: export GO111MODULE=on
build-linux-armv7: export GOPROXY=$(MOD_PROXY_URL)
build-linux-armv7:
	go build -v -o bin/$(BINNAME)-linux-armv7 cmd/si/main.go

.PHONY: build-mac
build-mac: export GOOS=darwin
build-mac: export GOARCH=amd64
build-mac: export CGO_ENABLED=0
build-mac: export GO111MODULE=on
build-mac: export GOPROXY=$(MOD_PROXY_URL)
build-mac:
	go build -v -o bin/$(BINNAME)-darwin-amd64 cmd/si/main.go

.PHONY: build-mac-arm
build-mac-arm: export GOOS=darwin
build-mac-arm: export GOARCH=arm64
build-mac-arm: export CGO_ENABLED=0
build-mac-arm: export GO111MODULE=on
build-mac-arm: export GOPROXY=$(MOD_PROXY_URL)
build-mac-arm:
	go build -v -o bin/$(BINNAME)-darwin-arm64 cmd/si/main.go

.PHONY: clean
clean:
	rm -f $(BINDIR)/*

.PHONY: dist
dist:

.PHONY: checksum
checksum:
	cd $(BINDIR) ; \
	for f in $$(ls -1 si-*-{arm64,armv7,amd64,mips64} 2>/dev/null) ; do \
		echo "Creating $${f}.sha256sum" ; \
		shasum -a 256 "$${f}" > "$${f}.sha256sum" ; \
	done
