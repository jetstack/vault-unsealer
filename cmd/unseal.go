// Copyright © 2017 Jetstack Ltd. <james@jetstack.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/vault"
)

type unsealCfg struct {
	unsealPeriod time.Duration
}

var unsealConfig unsealCfg

// unsealCmd represents the unseal command
var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		store, err := kvStoreForConfig(appConfig)

		if err != nil {
			logrus.Fatalf("error creating kv store: %s", err.Error())
		}

		cl, err := api.NewClient(api.DefaultConfig())

		if err != nil {
			logrus.Fatalf("error connecting to vault: %s", err.Error())
		}

		v, err := vault.New("vault", store, cl)

		if err != nil {
			logrus.Fatalf("error creating vault helper: %s", err.Error())
		}

		for {
			func() {
				logrus.Infof("checking if vault is sealed...")
				sealed, err := v.Sealed()
				if err != nil {
					logrus.Errorf("error checking if vault is sealed: %s", err.Error())
					return
				}

				logrus.Infof("vault sealed: %t", sealed)

				// If vault is not sealed, we stop here and wait another 30 seconds
				if !sealed {
					return
				}

				if err = v.Unseal(); err != nil {
					logrus.Errorf("error unsealing vault: %s", err.Error())
					return
				}

				logrus.Infof("successfully unsealed vault")
			}()
			// wait 30 seconds before trying again
			time.Sleep(time.Second * 30)
		}
	},
}

func init() {
	unsealCmd.Flags().DurationVar(&unsealConfig.unsealPeriod, "unseal-period", time.Second*30, "How often to attempt to unseal the vault instance")

	RootCmd.AddCommand(unsealCmd)
}
