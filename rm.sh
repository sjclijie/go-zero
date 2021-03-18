#!/bin/bash


if [ -e /Users/jaylee/Workspace/go-zero/test/iamAppModel.go ]; then
  rm -rf /Users/jaylee/Workspace/go-zero/test/iamAppModel.go
fi

if [ -e /Users/jaylee/Workspace/go-zero/test/vars.go ]; then
  rm -rf /Users/jaylee/Workspace/go-zero/test/vars.go
fi

if [ -d /Users/jaylee/Workspace/go-zero/example/bookstore/api/internal ]; then
  rm -rf /Users/jaylee/Workspace/go-zero/example/bookstore/api/internal/
fi
