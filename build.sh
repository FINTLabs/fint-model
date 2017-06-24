#!/usr/bin/env sh


rm -rf build

mkdir build
cd build
mkdir windows
mkdir mac
mkdir linux

cd ..

CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=i686-w64-mingw32-gcc go build -o ./build/windows/fint-model.exe
GOOS=darwin go build -o ./build/mac/fint-model
CGO_ENABLED=0 GOOS=linux go build -o ./build/linux/fint-model

zip build/windows.zip ./build/windows/fint-model.exe
zip build/linux.zip ./build/linux/fint-model
zip build/mac.zip ./build/mac/fint-model