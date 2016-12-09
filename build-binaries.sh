#!/bin/bash

mkdir -p dist/windows dist/osx dist/linux

GOOS=linux GOARCH=amd64 go build github.com/WikiLeaksFreedomForce/local-blockchain-parser
mv local-blockchain-parser dist/linux

GOOS=darwin GOARCH=amd64 go build github.com/WikiLeaksFreedomForce/local-blockchain-parser
mv local-blockchain-parser dist/osx

GOOS=windows GOARCH=amd64 go build github.com/WikiLeaksFreedomForce/local-blockchain-parser
mv local-blockchain-parser.exe dist/windows