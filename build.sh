#!/usr/bin/env sh


rm -rf build

gox -output="./build/{{.OS}}/{{.Dir}}" -rebuild -osarch="darwin/amd64 windows/amd64"