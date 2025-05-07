FROM golang:1.24.3-alpine AS build-env

WORKDIR /build

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY *.go go.mod go.sum /build/
COPY static/ /build/static/

RUN go version
RUN go build -ldflags="-X 'main.buildTime=$(date -R -u)'"

FROM scratch

COPY --from=build-env /build/arnested.dk /arnested.dk

ENTRYPOINT ["/arnested.dk"]
