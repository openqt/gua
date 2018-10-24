.PHONY: all build version

.EXPORT_ALL_VARIABLES:

APPVER=v1.0.0
CGO_ENABLED=0
GO111MODULE=on

GITVER=$(shell git rev-parse --short HEAD)
GOVER=$(shell go version)
BUILDTIME=$(shell date +%FT%T%z)

HTTPS_PROXY=socks5://127.0.0.1:1080/

all: build version

build:
	go-bindata -pkg yi --prefix yi -o yi/bindata.go yi/data.json
	go build -v -ldflags "-X 'github.com/openqt/gua/yi.AppVersion=${APPVER}' -X 'github.com/openqt/gua/yi.GoVersion=${GOVER}' -X 'github.com/openqt/gua/yi.GitVersion=${GITVER}' -X 'github.com/openqt/gua/yi.BuildTime=${BUILDTIME}'"

version:
	./gua version

update:
	go mod vendor -v
