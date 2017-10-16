FROM golang

RUN go-wrapper download github.com/mitchellh/gox
RUN go-wrapper install github.com/mitchellh/gox

WORKDIR /go/src/app/vendor/github.com/FINTprosjektet/fint-model

COPY . .

RUN go-wrapper download
ARG VERSION=0.0.0
RUN gox -output="./build/{{.Dir}}-{{.OS}}" -verbose -rebuild -osarch="darwin/amd64 windows/amd64" -ldflags="-X main.Version=${VERSION}"
