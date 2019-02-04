FROM golang AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/app/vendor/github.com/FINTprosjektet/fint-model
ARG VERSION=0.0.0
COPY . .
RUN go install -v -ldflags "-X main.Version=${VERSION}"
RUN /go/bin/fint-model --version

FROM gcr.io/distroless/static
VOLUME [ "/src" ]
WORKDIR /src
COPY --from=builder /go/bin/fint-model /usr/bin/fint-model
ENTRYPOINT [ "/usr/bin/fint-model" ]
