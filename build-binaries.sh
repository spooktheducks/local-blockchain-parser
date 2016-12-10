#!/bin/bash

mkdir -p dist

echo Building Linux binary...
GOOS=linux GOARCH=amd64 go build github.com/WikiLeaksFreedomForce/local-blockchain-parser
mv local-blockchain-parser dist/local-blockchain-parser-linuxamd64

echo Building OSX binary...
GOOS=darwin GOARCH=amd64 go build github.com/WikiLeaksFreedomForce/local-blockchain-parser
mv local-blockchain-parser dist/local-blockchain-parser-osxamd64

echo Building Windows amd64 binary...
GOOS=windows GOARCH=amd64 go build github.com/WikiLeaksFreedomForce/local-blockchain-parser
mv local-blockchain-parser.exe dist/local-blockchain-parser-windowsamd64.exe

echo Building Windows 386 binary...
GOOS=windows GOARCH=386 go build github.com/WikiLeaksFreedomForce/local-blockchain-parser
mv local-blockchain-parser.exe dist/local-blockchain-parser-windowsamd64.exe