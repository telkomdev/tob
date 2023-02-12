#!/bin/sh

check_err()
{
    if [ "$1" -ne "0" ]; then
        echo "Error # ${1} : ${2}"
        exit ${1}
    fi
}

package() 
{
    VERSION=$1
    OSNAME="linux"

    if [ -z "$VERSION" ]; then
        echo "$0 require version argument"
        return 1
    fi

    echo "packaging version $VERSION"

    echo "packaging for Apple's OSX"
    make build-osx
    tar -czvf tob-${VERSION}.darwin-amd64.tar.gz tob
    rm tob

    echo "packaging for Linux"
    make build-linux
    tar -czvf tob-${VERSION}.linux-amd64.tar.gz tob
    rm tob

    echo "packaging for Windows"
    make build-win
    zip tob-${VERSION}.win-amd64.zip tob.exe
    rm tob.exe

    echo "generate SHA256 checksum ..."

    if [ $(uname) = "Darwin" ]; then
        OSNAME="darwin"
    fi

    if [ "$OSNAME" = "linux" ]; then
        sha256sum tob-${VERSION}.darwin-amd64.tar.gz >> sha256sums.txt
        sha256sum tob-${VERSION}.linux-amd64.tar.gz >> sha256sums.txt
        sha256sum tob-${VERSION}.win-amd64.zip >> sha256sums.txt
    else 
        shasum -a 256 tob-${VERSION}.darwin-amd64.tar.gz >> sha256sums.txt
        shasum -a 256 tob-${VERSION}.linux-amd64.tar.gz >> sha256sums.txt
        shasum -a 256 tob-${VERSION}.win-amd64.zip >> sha256sums.txt
    fi

    return 0
}

package "$@"
check_err $? "package returned error"

# How to run this script

# always execute package.sh from root project folder
# ./scripts/package.sh YOUR_NEW_VERSION_HERE

# eg:
# ./scripts/package.sh 1.0.0