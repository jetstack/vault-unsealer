package main

import (
	"flag"

	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/cmd"
)

func main() {
	flag.Parse()
	cmd.Execute()
}
