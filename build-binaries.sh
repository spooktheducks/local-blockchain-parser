#!/bin/bash

mkdir -p dist

echo Building Linux binary...
GOROOT_FINAL=/usr/local/go GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH github.com/spooktheducks/local-blockchain-parser
mv local-blockchain-parser dist/local-blockchain-parser-linuxamd64

echo Building OSX binary...
GOROOT_FINAL=/usr/local/go GOOS=darwin GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH github.com/spooktheducks/local-blockchain-parser
mv local-blockchain-parser dist/local-blockchain-parser-osxamd64

echo Building Windows amd64 binary...
GOROOT_FINAL=/usr/local/go GOOS=windows GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH github.com/spooktheducks/local-blockchain-parser
mv local-blockchain-parser.exe dist/local-blockchain-parser-windowsamd64.exe

echo Building Windows 386 binary...
GOROOT_FINAL=/usr/local/go GOOS=windows GOARCH=386 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH github.com/spooktheducks/local-blockchain-parser
mv local-blockchain-parser.exe dist/local-blockchain-parser-windows386.exe

