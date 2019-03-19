package main

import (
	"github.com/jetstack/vault-unsealer/cmd"
)

var (
	version string = "dev"
	commit  string = "unknown"
	date    string = "unknown"
)

func main() {
	cmd.Version.Version = version
	cmd.Version.Commit = commit
	cmd.Version.BuildDate = date
	cmd.Execute()
}
