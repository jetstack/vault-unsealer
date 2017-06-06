FROM alpine:3.5

RUN apk add --update ca-certificates

COPY vault-unsealer_linux_amd64 /usr/local/bin/vault-unsealer

CMD ["/usr/local/bin/vault-unsealer"]
