#!/bin/bash

OS_VERSION="fedora-31"
PASSWORD="fedora"
OS_IMG="os.img"
if test -f "$FILE"; then
  echo "$FILE exists, nothing to do"
  exit 0
fi

VIRT_BUILDER=$(which virt-builder)
if [ $? -eq 0 ]; then
  ${VIRT_BUILDER} ${OS_VERSION} --size 6G -o ${OS_IMG} --root-password password:${PASSWORD} --install "tar"
else
  echo virt-builder not found
  exit 1
fi
