#!/bin/sh

NAME="pm-creds"
SOURCE_DIR="pm-creds"
TMP_DIR="tmp"
DIST_DIR="dist"

BUILD_FOR="darwin/arm64 darwin/amd64 windows/amd64 windows/386 linux/arm64 linux/amd64 linux/386"

mkdir -p ${TMP_DIR} ${DIST_DIR}

cd ${SOURCE_DIR}
for BUILD in ${BUILD_FOR}; do
	GOOS=${BUILD%/*}
	GOARCH=${BUILD#*/}
	FILENAME=${NAME}
	ZIPNAME="${NAME}-${GOOS}-${GOARCH}.zip"
	echo "Building for ${BUILD} as ${FILENAME}"
	if [ "${GOOS}" = "windows" ]; then
		FILENAME="${FILENAME}.exe"
	fi
	
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../${TMP_DIR}/${FILENAME}
	zip -j ../${DIST_DIR}/${ZIPNAME} ../${TMP_DIR}/${FILENAME}
	rm ../${TMP_DIR}/${FILENAME}
done
