#!/bin/bash

APP_NAME="mflow"
VERSION=$(git describe --abbrev=0 --tags)


build(){
	mkdir -p "release/${APP_NAME}-${VERSION}"
        cp config.json "release/${APP_NAME}-${VERSION}"
        cp oracle.json "release/${APP_NAME}-${VERSION}"
        go build -o "release/${APP_NAME}-${VERSION}/${APP_NAME}" -v
        pushd release
        tar -czvf "${APP_NAME}-${VERSION}.tar.gz" "./${APP_NAME}-${VERSION}"
        rm -fr "${APP_NAME}-${VERSION}"
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
	go get
}

case "$1" in
   build) 
	   clean
	   build 
	;;
   run)  run;;
   clean) clean;;
   getdeps) getdeps;;
   *) echo "usage $0 build|run|getdeps|clean|test" >&2
      exit 1
    ;;
esac

