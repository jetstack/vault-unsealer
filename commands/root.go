package commands

import (
	"flag"

	"github.com/golang/glog"
	"github.com/soter/vault-unsealer/pkg/worker"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	opts := worker.NewWorkerOptions()

	cmd := &cobra.Command{
		Use:   "vault-unsealer",
		Short: "Automates initialisation and unsealing of Hashicorp Vault.",

		Run: func(cmd *cobra.Command, args []string) {
			if errs := opts.Validate(); errs != nil {
				glog.Fatal(errs)
			}
			if err := opts.Run(); err != nil {
				glog.Fatal(err)
			}
		},
	}

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})

	opts.AddFlags(cmd.Flags())

	return cmd
}
