FROM golang:1.23.6-alpine3.21 as build

RUN apk add --no-cache make

WORKDIR /usr/src/logo-ls

COPY . /usr/src/logo-ls

RUN make

FROM alpine:3.21 as example

COPY --from=build /usr/src/logo-ls/bin/logo-ls /usr/local/bin/logo-ls


