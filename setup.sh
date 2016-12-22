#!/bin/bash

set -e

PROJECT_NAMESPACE="${GOPATH}/src/github.com/WikiLeaksFreedomForce/"
PROJECT="local-blockchain-parser"

# You should be using ssh, update your origin after cloning
REPOSITORY="https://github.com/WikiLeaksFreedomForce/local-blockchain-parser.git"

if [ -z ${GOPATH+x} ]
	then
		echo "please install go and setup gopath, see readme"
		exit 1
	else
		go get github.com/tools/godep
		mkdir -p ${PROJECT_NAMESPACE}
		cd ${PROJECT_NAMESPACE}
		git clone ${REPOSITORY} ${PROJECT}
		cd ${PROJECT_NAMESPACE}${PROJECT}
		godep go build
		echo "build success"
fi
