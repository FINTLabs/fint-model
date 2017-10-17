FROM golang

RUN go-wrapper download github.com/mitchellh/gox
RUN go-wrapper install github.com/mitchellh/gox

WORKDIR /go/src/app/vendor/github.com/FINTprosjektet/fint-model

COPY . .

RUN go-wrapper download
