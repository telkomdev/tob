#!/bin/sh

check_err()
{
    if [ "$1" -ne "0" ]; then
        echo "Error # ${1} : ${2}"
        exit ${1}
    fi
}

gen_sha256_pass() 
{
    PLAIN_PASS="$1"
    OSNAME="linux"

    if [ -z "$PLAIN_PASS" ]; then
        echo "$0 requires a plain password argument"
        return 1
    fi

    echo "Generating SHA256 from '$PLAIN_PASS'"

    if [ "$(uname)" = "Darwin" ]; then
        OSNAME="darwin"
    fi

    if [ "$OSNAME" = "linux" ]; then
        printf "%s" "$PLAIN_PASS" | sha256sum | awk '{print $1}'
    else 
        printf "%s" "$PLAIN_PASS" | shasum -a 256 | awk '{print $1}'
    fi

    return 0
}

gen_sha256_pass "$@"
check_err $? "gen_sha256_pass returned an error"
