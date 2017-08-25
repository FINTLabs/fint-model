FROM golang:latest

RUN go-wrapper download github.com/mitchellh/gox
RUN go-wrapper install github.com/mitchellh/gox

WORKDIR /go/src/app/vendor/github.com/FINTprosjektet/fint-model

COPY . .

RUN go-wrapper download
RUN gox -output="./build/{{.Dir}}-{{.OS}}" -rebuild -osarch="windows/amd64"
RUN curl --upload-file build/fint-model-windows.exe https://transfer.sh/fint-model-windows.exe
