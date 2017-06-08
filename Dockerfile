FROM alpine:3.6

COPY vault-unsealer_linux_amd64 /usr/local/bin/vault-unsealer

CMD ["/usr/local/bin/vault-unsealer"]
