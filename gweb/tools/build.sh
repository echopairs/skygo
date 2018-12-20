#!/usr/bin/env bash

script=$(readlink -f "$0")
route=$(dirname "$script")


VERSION=$(git describe --abbrev=0 --tags 2>/dev/null)
GITCOMMIT=$(git rev-parse --short HEAD)
GITBRANCE=$(git branch | grep "*" | cut -d" " -f2)
BUILDTIME=$(date +%Y-%m-%d-%H:%M:%S)

if test -z "$VERSION"
then
    VERSION=$GITCOMMIT
fi

echo "VERSION:" $VERSION
echo "COMMIT:"  $GITCOMMIT
echo "BRANCH:"  $GITBRANCE
echo "TIME:"    $BUILDTIME

cd ${route}/../cmd

LDFLAGS="-s -X github.com/echopairs/skygo/version.VERSION=\"$VERSION\"
    -X github.com/echopairs/skygo/version.GITBRANCH=\"$GITBRANCE\"
    -X github.com/echopairs/skygo/version.GITCOMMIT=\"$GITCOMMIT\"
    -X github.com/echopairs/skygo/version.BUILDTIME=$BUILDTIME"

echo $LDFLAGS

for e in "./cloudweb"
do
    echo "building $e ..."
    name=${e##*/}
    go build -ldflags "$LDFLAGS" -o bin/$name $e
done

# https://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
