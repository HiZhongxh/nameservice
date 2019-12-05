#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/nsd/${BINARY:-nsd}
ID=${ID:-0}
LOG=${LOG:-nsd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'nsd' E.g.: -e BINARY=nsd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export NSDHOME="/nsd/node${ID}/nsd"

if [ -d "$(dirname "${NSDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${NSDHOME}" "$@" | tee "${NSDHOME}/${LOG}"
else
  "${BINARY}" --home "${NSDHOME}" "$@"
fi

