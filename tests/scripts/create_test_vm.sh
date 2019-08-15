#!/bin/bash

OS_IMG="os.img"
if test -f "$FILE"; then
    echo "$FILE exists, nothing to do"
    exit 0
fi

if [ $# -eq 0 ]
  then
    URL="https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-1503.qcow2"
    echo "No image URL supplied, downloading ${URL}..."
else
    URL=$1
    echo "Downloading URL ${URL}..."
fi

curl "${URL}" -o os.img
