package main

import (
	"os"

	logs "github.com/appscode/go/log/golog"
	"github.com/soter/vault-unsealer/commands"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := commands.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
