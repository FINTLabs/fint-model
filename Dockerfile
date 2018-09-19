FROM golang:alpine AS builder
WORKDIR /go/src/app/vendor/github.com/FINTprosjektet/fint-model
ARG VERSION=0.0.0
COPY . .
RUN go install -v -ldflags "-X main.Version=${VERSION}"
RUN /go/bin/fint-model --version

FROM alpine
RUN apk --update add ca-certificates
COPY --from=builder /go/bin/fint-model /usr/bin/fint-model
WORKDIR /src
VOLUME [ "/src" ]
ENTRYPOINT [ "/usr/bin/fint-model" ]
