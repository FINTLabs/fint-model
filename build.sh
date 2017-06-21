#!/usr/bin/env sh

rm -rf build
CGO_ENABLED=0 GOOS=windows go build -o ./build/windows/fint-cli.exe
CGO_ENABLED=0 GOOS=darwin go build -o ./build/mac/fint-cli
CGO_ENABLED=0 GOOS=linux go build -o ./build/linux/fint-cli

zip build/windows.zip ./build/windows/fint-cli.exe
zip build/linux.zip ./build/linux/fint-cli
zip build/mac.zip ./build/mac/fint-cli