#!/bin/sh
set -e
GOOS=${1}
GOARCH=${2}

FRIENDLY=${GOOS}
if [ ${FRIENDLY} = "darwin" ]; then
	FRIENDLY=macOS
fi

FILENAME="pm-creds-${TAG}-${FRIENDLY}-${GOOS}"

GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../dist/${FILENAME}/pm-creds
if [ "${GOOS}" = "darwin" ]; then
	codesign --force -s "${AC_APPID}" --timestamp ../dist/${FILENAME}/pm-creds
	codesign -v ../dist/${FILENAME}/pm-creds
fi

zip -j ../dist/${FILENAME}.zip ../dist/${FILENAME}/pm-creds

if [ "${GOOS}" = "darwin" ]; then
	xcrun altool --notarize-app --primary-bundle-id "se.execit.pm-creds.zip" -u "${AC_USERNAME}" -p "${AC_PASSWORD}" -t osx -f ../dist/${FILENAME}.zip
fi
