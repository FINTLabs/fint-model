FROM golang:latest

WORKDIR /go/src/app/vendor/github.com/FINTprosjektet/fint-model
COPY . .

RUN go-wrapper download github.com/mitchellh/gox
RUN go-wrapper install github.com/mitchellh/gox
RUN go-wrapper download
RUN gox -output="./build/{{.OS}}/{{.Dir}}" -rebuild -osarch="windows/amd64"
RUN curl --upload-file build/windows/fint-model.exe https://transfer.sh/fint-model.exe
