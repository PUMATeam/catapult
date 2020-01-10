#!/bin/bash

test -z "$(gofmt -l ../../)"
SUCCESS=$?

if [ $SUCCESS -eq 0 ]; then
  exit 0
else
  gofmt -l ../../
fi

