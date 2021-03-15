#!/bin/sh
set -e
GOOS=${1}
GOARCH=${2}

FRIENDLY=${GOOS}
if [ ${FRIENDLY} = "darwin" ]; then
	FRIENDLY=macOS
fi

FILENAME="pm-creds-${TAG}-${FRIENDLY}-${GOOS}"

# Build
mkdir -p dist/${FILENAME}
cd pm-creds
GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../dist/${FILENAME}/pm-creds

# Notarize and/or zip.
cd ../dist/${FILENAME}

if [ "${GOOS}" = "darwin" ]; then
	codesign --deep --force -s "${AC_APPID}" --timestamp pm-creds
	codesign --verify --verbose ${FILENAME}/pm-creds
	ditto -c -k --keepParent --rsrc pm-creds ${FILENAME}.zip
	xcrun altool --notarize-app --primary-bundle-id "se.execit.pm-creds.zip" -u "${AC_USERNAME}" -p "${AC_PASSWORD}" -t osx -f ${FILENAME}.zip
else
	ditto -c -k --keepParent --rsrc pm-creds ${FILENAME}.zip
fi

rm -rf pm-creds

cd ../..
