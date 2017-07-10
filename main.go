package main

import (
	"flag"

	"github.com/jetstack-experimental/vault-unsealer/cmd"
)

func main() {
	flag.Parse()
	cmd.Execute()
}
