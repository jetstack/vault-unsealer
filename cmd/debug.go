package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug vault unseal keys",
}

var debugGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get unseal keys by chosen mode",
	Example: "vault-unsealer get vault-unseal-0 vault-unseal-1 vault-unseal-2",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig.BindPFlag(cfgInitRootToken, cmd.PersistentFlags().Lookup(cfgInitRootToken))
		appConfig.BindPFlag(cfgStoreRootToken, cmd.PersistentFlags().Lookup(cfgStoreRootToken))
		appConfig.BindPFlag(cfgOverwriteExisting, cmd.PersistentFlags().Lookup(cfgOverwriteExisting))

		if len(args) == 0 {
			logrus.Fatal("Provide at least one key ID to get")
		}

		store, err := kvStoreForConfig(appConfig)
		if err != nil {
			logrus.Fatalf("error creating kv store: %s", err.Error())
		}

		var result *multierror.Error
		for _, a := range args {
			k, err := store.Get(a)
			if err != nil {
				result = multierror.Append(result, fmt.Errorf("failed to get key '%s': %s", a, err))
				continue
			}

			fmt.Printf("%s: %s", a, k)
		}

		if result != nil {
			logrus.Fatalf("failed to get key(s): %s", result.ErrorOrNil())
		}

		os.Exit(0)
	},
}

var debugSetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Set unseal keys by chosen mode using a ':' separated name-key pair",
	Example: "vault-unsealer set vault-unseal-0:my-secret-share-0 vault-unseal-1:my-secret-share-1 vault-unseal-2:my-secret-share=2",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig.BindPFlag(cfgInitRootToken, cmd.PersistentFlags().Lookup(cfgInitRootToken))
		appConfig.BindPFlag(cfgStoreRootToken, cmd.PersistentFlags().Lookup(cfgStoreRootToken))
		appConfig.BindPFlag(cfgOverwriteExisting, cmd.PersistentFlags().Lookup(cfgOverwriteExisting))

		if len(args) == 0 {
			logrus.Fatal("Provide at least one name:key pair to set")
		}

		store, err := kvStoreForConfig(appConfig)
		if err != nil {
			logrus.Fatalf("error creating kv store: %s", err.Error())
		}

		var result *multierror.Error
		for _, a := range args {
			pair := strings.Split(a, ":")
			if len(pair) != 2 {
				err := fmt.Errorf("failed to parse name:key pair '%s', expecting single separator ':'", a)
				result = multierror.Append(result, err)
				continue
			}

			if err := store.Set(pair[0], []byte(pair[1])); err != nil {
				err := fmt.Errorf("failed to store key '%s' with value '%s': %s", pair[0], pair[1], err)
				result = multierror.Append(result, err)
			}
		}

		if result != nil {
			logrus.Fatalf("failed to set key(s): %s", result.ErrorOrNil())
		}

		fmt.Print("Key(s) set successfully.\n")
		os.Exit(0)
	},
}

var debugTestCmd = &cobra.Command{
	Use:     "test",
	Short:   "Test unseal keys by chosen mode",
	Example: "vault-unsealer test vault-unseal-0 vault-unseal-1 vault-unseal-2",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig.BindPFlag(cfgInitRootToken, cmd.PersistentFlags().Lookup(cfgInitRootToken))
		appConfig.BindPFlag(cfgStoreRootToken, cmd.PersistentFlags().Lookup(cfgStoreRootToken))
		appConfig.BindPFlag(cfgOverwriteExisting, cmd.PersistentFlags().Lookup(cfgOverwriteExisting))

		if len(args) == 0 {
			logrus.Fatal("Provide at least one key ID to test")
		}

		store, err := kvStoreForConfig(appConfig)
		if err != nil {
			logrus.Fatalf("error creating kv store: %s", err.Error())
		}

		var result *multierror.Error
		for _, a := range args {
			if err := store.Test(a); err != nil {
				result = multierror.Append(result, fmt.Errorf("failed to test key '%s': %s", a, err))
			}
		}

		if result != nil {
			logrus.Fatalf("failed to test key(s): %s", result.ErrorOrNil())
		}

		fmt.Print("Key(s) tested successfully.\n")
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(debugCmd)
	debugCmd.AddCommand(debugGetCmd)
	debugCmd.AddCommand(debugSetCmd)
	debugCmd.AddCommand(debugTestCmd)
}
