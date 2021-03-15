#!/bin/sh
set -e
GOOS=${1}
GOARCH=${2}
TAG=$(git tag --points-at HEAD)

FRIENDLY=${GOOS}
if [ ${FRIENDLY} = "darwin" ]; then
	FRIENDLY=macOS
fi

FILENAME="pm-creds-${TAG}-${FRIENDLY}-${GOARCH}"

# Build
mkdir -p dist/${FILENAME}
cd pm-creds
GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../dist/${FILENAME}/pm-creds

# Notarize and/or zip.
cd ../dist/${FILENAME}

if [ "${GOOS}" = "darwin" ]; then
	codesign --deep --force -s "${AC_APPID}" --timestamp --options runtime pm-creds
	codesign --verify --verbose pm-creds
	ditto -c -k --keepParent --rsrc pm-creds ${FILENAME}.zip
	xcrun altool --notarize-app --primary-bundle-id "se.execit.pm-creds.zip" -u "${AC_USERNAME}" -p "${AC_PASSWORD}" -t osx -f ${FILENAME}.zip
else
	ditto -c -k --keepParent --rsrc pm-creds ${FILENAME}.zip
fi

rm -rf pm-creds

cd ../..
