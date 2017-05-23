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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type kvCfg struct {
	googleCloudStorageBucket string
	googleCloudStoragePrefix string
	googleCloudKMSProject    string
	googleCloudKMSLocation   string
	googleCloudKMSKeyRing    string
	googleCloudKMSCryptoKey  string
}

var kvConfig kvCfg
var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vault-unsealer",
	Short: "Automates initialisation and unsealing of Hashicorp Vault.",
	Long: `This is a CLI tool to help automate the setup and management of
Hashicorp Vault.

It will continuously attempt to unseal the target Vault instance, by retrieving
unseal keys from a Google Cloud KMS keyring.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// Google Cloud KMS flags
	RootCmd.Flags().StringVar(&kvConfig.googleCloudKMSProject, "google-cloud-kms-project", "", "The Google Cloud KMS project to use")
	RootCmd.Flags().StringVar(&kvConfig.googleCloudKMSLocation, "google-cloud-kms-location", "", "The Google Cloud KMS location to use (eg. 'global', 'europe-west1')")
	RootCmd.Flags().StringVar(&kvConfig.googleCloudKMSKeyRing, "google-cloud-kms-key-ring", "", "The name of the Google Cloud KMS key ring to use")
	RootCmd.Flags().StringVar(&kvConfig.googleCloudKMSCryptoKey, "google-cloud-kms-crypto-key", "", "The name of the Google Cloud KMS crypt key to use")

	// Google Cloud Storage flags
	RootCmd.Flags().StringVar(&kvConfig.googleCloudStorageBucket, "google-cloud-storage-bucket", "", "The name of the Google Cloud Storage bucket to store values in")
	RootCmd.Flags().StringVar(&kvConfig.googleCloudStoragePrefix, "google-cloud-storage-prefix", "", "The prefix to use for values store in Google Cloud Storage")
}
