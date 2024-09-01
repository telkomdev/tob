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

    echo "packaging tob-http-agent version $VERSION"

    echo "packaging for Apple's OSX"
    make build-http-agent-osx
    tar -czvf tob-http-agent-${VERSION}.darwin-amd64.tar.gz tob-http-agent
    rm tob-http-agent

    echo "packaging for Apple's OSX with Apple Chip"
    make build-http-agent-osx-arm
    tar -czvf tob-http-agent-${VERSION}.darwin-arm64.tar.gz tob-http-agent
    rm tob-http-agent

    echo "packaging for Linux"
    make build-http-agent-linux
    tar -czvf tob-http-agent-${VERSION}.linux-amd64.tar.gz tob-http-agent
    rm tob-http-agent

    echo "generate SHA256 checksum ..."

    if [ $(uname) = "Darwin" ]; then
        OSNAME="darwin"
    fi

    if [ "$OSNAME" = "linux" ]; then
        sha256sum tob-http-agent-${VERSION}.darwin-amd64.tar.gz >> tob-http-agent-sha256sums.txt
        sha256sum tob-http-agent-${VERSION}.darwin-arm64.tar.gz >> tob-http-agent-sha256sums.txt
        sha256sum tob-http-agent-${VERSION}.linux-amd64.tar.gz >> tob-http-agent-sha256sums.txt
    else 
        shasum -a 256 tob-http-agent-${VERSION}.darwin-amd64.tar.gz >> tob-http-agent-sha256sums.txt
        shasum -a 256 tob-http-agent-${VERSION}.darwin-arm64.tar.gz >> tob-http-agent-sha256sums.txt
        shasum -a 256 tob-http-agent-${VERSION}.linux-amd64.tar.gz >> tob-http-agent-sha256sums.txt
    fi

    return 0
}

package "$@"
check_err $? "package_http_agent returned error"

# How to run this script

# always execute package_http_agent.sh from root project folder
# ./scripts/package_http_agent.sh YOUR_NEW_VERSION_HERE

# eg:
# ./scripts/package_http_agent.sh 1.0.0