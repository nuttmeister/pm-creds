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
	codesign --force --timestamp --options runtime --sign "${AC_APPID}" pm-creds
	codesign --verify --verbose pm-creds
	ditto -c -k --keepParent --rsrc pm-creds ${FILENAME}.zip
	xcrun altool --notarize-app --primary-bundle-id "se.execit.pm-creds.zip" -u "${AC_USERNAME}" -p "${AC_PASSWORD}" -t osx -f ${FILENAME}.zip
elif [ "${GOOS}" = "windows" ]; then
	# npm install --save-dev signcode
	# ./node_modules/signcode/cli.js sign ${FILENAME}.exe \
	# 	--cert cert-win.p12 \
	# 	--password ${WIN_PASSWORD} \
	# 	--name 'pm-creds' \
	# 	--url 'https://github.com/nuttmeister/pm-creds'
	# signcode verify ${FILENAME}.exe
	ditto -c -k --keepParent --rsrc pm-creds.exe ${FILENAME}.zip
else
	ditto -c -k --keepParent --rsrc pm-creds ${FILENAME}.zip
fi

rm -rf pm-creds

cd ../..
