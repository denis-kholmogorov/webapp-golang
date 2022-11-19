FROM golang:alpine AS builder

WORKDIR /build

COPY ./web /usr/local/go/src/web
COPY ./github.com/gin-gonic/gin /usr/local/go/src/github.com/gin-gonic/gin
COPY ./github.com/jackc/pgx /usr/local/go/src/github.com/jackc/pgx
COPY ./github.com/georgysavva/scany /usr/local/go/src/github.com/georgysavva/scany


RUN ls -la

RUN go build -o webgin /usr/local/go/src/web/webgin.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/webgin /build/webgin

RUN go run webgin