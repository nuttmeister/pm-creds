#!/bin/sh

NAME="pm-creds"
SOURCE_DIR="pm-creds"
DIST_DIR="dist"

BUILD_FOR="darwin/arm64 darwin/amd64 windows/amd64 windows/386 linux/arm64 linux/amd64 linux/386"

mkdir -p ${SOURCE_DIR} && cd ${SOURCE_DIR}

for BUILD in ${BUILD_FOR}; do
	GOOS=${BUILD%/*}
	GOARCH=${BUILD#*/}
	FILENAME="${NAME}-${GOOS}-${GOARCH}"
	echo "Building for ${BUILD} as ${FILENAME}"
	if [ "${GOOS}" = "windows" ]; then
		FILENAME="${FILENAME}.exe"
	fi
	
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../${DIST_DIR}/${FILENAME}
done
