FROM golang:1.7.3 as builder

WORKDIR /go/src/github.com/jetstack-experimental/vault-unsealer

COPY ./ /go/src/github.com/jetstack-experimental/vault-unsealer/

RUN make go_verify

RUN make go_build

FROM alpine:3.6

COPY --from=builder /go/src/github.com/jetstack-experimental/vault-unsealer/vault-unsealer_linux_amd64 /usr/local/bin/vault-unsealer


 