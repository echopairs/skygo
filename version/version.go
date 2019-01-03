package version

import (
	"log"
)

// 自动添加版本信息，需要在build的时候传ldflags进去，参考脚本如下
//
// #! /bin/bash
//
// VERSION=$(git describe --abbrev=0 --tags 2>/dev/null)
// REVCNT=$(git rev-list --count HEAD 2>/dev/null)
// if test "$VERSION" == ""
// then
//     VERSION="dev$REVCNT"
// else
//     DEVCNT=$(git rev-list --count $VERSION)
//     if test $REVCNT != $DEVCNT
//     then
//         VERSION="$VERSION.dev$(expr $REVCNT - $DEVCNT)"
//     fi
// fi
//
// GITCOMMIT=$(git rev-parse HEAD)
// GITBRANCH=$(git branch | cut -d" " -f2)
// BUILDTIME=$(date +%Y/%m/%d-%H:%M:%S)
//
// echo "VERSION: $VERSION"
// echo "BRANCH:  $GITBRANCH"
// echo "COMMIT:  $GITCOMMIT"
// echo "TIME:    $BUILDTIME"
//
// LDFLAGS="-s -X common/version.VERSION=$VERSION -X common/version.BUILDTIME=$BUILDTIME -X common/version.GITCOMMIT=$GITCOMMIT -X common/version.GITBRANCH=$GITBRANCH"
//
// go build -ldflags "$LDFLAGS" xxx

var (
	VERSION   = "unknown"
	BUILDTIME = "unknown"
	GITBRANCH = "unknown"
	GITCOMMIT = "unknown"

	showVersion bool
)

func Show() {
	log.Printf("VERSION: %s\n", VERSION)
	log.Printf("BUILDTIME: %s\n", BUILDTIME)
	log.Printf("GITBRANCH: %s\n", GITBRANCH)
	log.Printf("GITCOMMIT: %s\n", GITCOMMIT)
}
