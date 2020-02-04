#!/bin/bash

APP_NAME="mflow"
VERSION=$(git describe --abbrev=0 --tags)


build_first(){
	mkdir -p "release/${APP_NAME}-${VERSION}"
        cp config.json "release/${APP_NAME}-${VERSION}"
        cp oracle.json "release/${APP_NAME}-${VERSION}"
        go build -o "release/${APP_NAME}-${VERSION}/${APP_NAME}" -v
        pushd release
        tar -czvf "${APP_NAME}-${VERSION}.tar.gz" "./${APP_NAME}-${VERSION}"
        rm -fr "${APP_NAME}-${VERSION}"
        popd
}

build(){
	mkdir -p "release/${APP_NAME}-${VERSION}"
        go build -o "release/${APP_NAME}-${VERSION}/${APP_NAME}" -v
        pushd "release/${APP_NAME}-${VERSION}"
        tar -czvf "${APP_NAME}-${VERSION}.tar.gz" "${APP_NAME}"
        rm -fr "${APP_NAME}"
        popd
}


release(){
   build
   pushd release/${APP_NAME}-${VERSION}
   md5sum ${APP_NAME}-${VERSION}.tar.gz | awk '{print $1}' > ${APP_NAME}-${VERSION}.tar.gz.md5
   sha1sum ${APP_NAME}-${VERSION}.tar.gz | awk '{print $1}' > ${APP_NAME}-${VERSION}.tar.gz.sha1
   echo "Usuario sua: "
   read username
   echo "Password sua: "
   read password
   curl --noproxy '*' -v -k -u "$username:$password" --upload-file ${APP_NAME}-${VERSION}.tar.gz "https://nexus.cloudint.afip.gob.ar/nexus/repository/fisca-infraestructura-raw/$APP_NAME/$VERSION/${APP_NAME}-${VERSION}.tar.gz"
   curl --noproxy '*' -v -k -u "$username:$password" --upload-file ${APP_NAME}-${VERSION}.tar.gz.md5 "https://nexus.cloudint.afip.gob.ar/nexus/repository/fisca-seleccion-casos-raw/$APP_NAME/$VERSION/${APP_NAME}-${VERSION}.tar.gz.md5"
   curl --noproxy '*' -v -k -u "$username:$password" --upload-file ${APP_NAME}-${VERSION}.tar.gz.sha1 "https://nexus.cloudint.afip.gob.ar/nexus/repository/fisca-seleccion-casos-raw/$APP_NAME/$VERSION/${APP_NAME}-${VERSION}.tar.gz.sha1"
   echo "Release: https://nexus.cloudint.afip.gob.ar/nexus/repository/fisca-infraestructura-raw/$APP_NAME/$VERSION/${APP_NAME}-${VERSION}.tar.gz"
}


test(){
	go test -v -cover
}

clean(){ 
	go clean
	rm -fr release/*
}

run(){
	go run main.go
}

getdeps(){
	go get
}

case "$1" in
   build_first) 
	   clean
	   build_first
	;;
   build) 
	   clean
	   build 
	;;
   release) release;;
   run)  run;;
   clean) clean;;
   getdeps) getdeps;;
   test) test;;
   *) echo "usage $0 build_first|build|release|run|getdeps|clean|test" >&2
      exit 1
    ;;
esac

