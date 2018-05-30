FROM ubuntu

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY vault-unsealer /usr/local/bin/vault-unsealer

ENTRYPOINT ["/usr/local/bin/vault-unsealer"]
