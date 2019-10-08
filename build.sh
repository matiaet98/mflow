#!/bin/bash

APP_NAME="mflow"
VERSION=$(git describe --abbrev=0 --tags)


build(){
	mkdir -p release/mflow
        cp config.json release/mflow/
        cp oracle.json release/mflow/
        go build -o release/mflow/$APP_NAME -v
        pushd release
        tar -czvf "${APP_NAME}-${VERSION}.tar.gz" ./mflow
        rm -fr mflow
        popd
}

test(){
	go test -v -cover
}

clean(){ 
	go clean
	rm -f release/*
}

run(){
	go run main.go
}

getdeps(){
	go mod vendor
	go mod download
	go mod verify
	go mod tidy
	go mod graph
}

case "$1" in
   build) build ;;
   run)  run;;
   clean) clean;;
   getdeps) getdeps;;
   *) echo "usage $0 build|run|getdeps|clean|test" >&2
      exit 1
    ;;
esac

