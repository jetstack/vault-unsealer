FROM golang:1.10.3 as builder

ARG CI_COMMIT_TAG
ARG CI_COMMIT_SHA
ARG CI_DATE

COPY . /go/src/github.com/jetstack/vault-unsealer
WORKDIR /go/src/github.com/jetstack/vault-unsealer

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-w -X main.version=${CI_COMMIT_TAG} -X main.commit=${CI_COMMIT_SHA} -X main.date=${CI_DATE}" -o vault-unsealer

FROM alpine:3.6
RUN apk add --update ca-certificates
COPY --from=builder /go/src/github.com/jetstack/vault-unsealer/vault-unsealer /usr/local/bin/vault-unsealer

ENTRYPOINT ["/usr/local/bin/vault-unsealer"]