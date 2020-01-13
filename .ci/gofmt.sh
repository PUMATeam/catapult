#!/bin/bash

test -z "$(gofmt -l $GITHUB_WORKSPACE)"
SUCCESS=$?

if [ $SUCCESS -eq 0 ]; then
  exit 0
else
  echo Unformatted files:
  echo "$(gofmt -l $GITHUB_WORKSPACE)"
  exit $SUCCESS
fi

