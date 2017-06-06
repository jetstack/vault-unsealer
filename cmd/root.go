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
	"github.com/spf13/viper"
)

var appConfig *viper.Viper
var cfgFile string

const cfgMode = "mode"
const cfgModeValueAWSKMSSSM = "aws-kms-ssm"
const cfgModeValueGoogleCloudKMSGCS = "google-cloud-kms-gcs"

const cfgGoogleCloudKMSProject = "google-cloud-kms-project"
const cfgGoogleCloudKMSLocation = "google-cloud-kms-location"
const cfgGoogleCloudKMSKeyRing = "google-cloud-kms-key-ring"
const cfgGoogleCloudKMSCryptoKey = "google-cloud-kms-crypto-key"

const cfgGoogleCloudStorageBucket = "google-cloud-storage-bucket"
const cfgGoogleCloudStoragePrefix = "google-cloud-storage-prefix"

const cfgAWSKMSKeyID = "aws-kms-key-id"
const cfgAWSSSMKeyPrefix = "aws-ssm-key-prefix"

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

func configStringVar(key, defaultValue, description string) {
	RootCmd.PersistentFlags().String(key, defaultValue, description)
	appConfig.BindPFlag(key, RootCmd.PersistentFlags().Lookup(key))
}

func init() {
	appConfig = viper.New()
	appConfig.SetEnvPrefix("vault_unsealer")
	appConfig.AutomaticEnv()

	// SelectMode
	configStringVar(
		cfgMode,
		cfgModeValueGoogleCloudKMSGCS,
		fmt.Sprintf("Select the mode to use '%s' => Google Cloud Storage with encryption using Google KMS; '%s' => AWS SSM parameter store using AWS KMS encryption", cfgModeValueGoogleCloudKMSGCS, cfgModeValueAWSKMSSSM),
	)

	// Google Cloud KMS flags
	configStringVar(cfgGoogleCloudKMSProject, "", "The Google Cloud KMS project to use")
	configStringVar(cfgGoogleCloudKMSLocation, "", "The Google Cloud KMS location to use (eg. 'global', 'europe-west1')")
	configStringVar(cfgGoogleCloudKMSKeyRing, "", "The name of the Google Cloud KMS key ring to use")
	configStringVar(cfgGoogleCloudKMSCryptoKey, "", "The name of the Google Cloud KMS crypt key to use")

	// Google Cloud Storage flags
	configStringVar(cfgGoogleCloudStorageBucket, "", "The name of the Google Cloud Storage bucket to store values in")
	configStringVar(cfgGoogleCloudStoragePrefix, "", "The prefix to use for values store in Google Cloud Storage")

	// AWS KMS Storage flags
	configStringVar("aws-kms-key-id", "", "The ID or ARN of the AWS KMS key to encrypt values")

	// AWS SSM Parameter Storage flags
	configStringVar("aws-ssm-key-prefix", "", "The Key Prefix for SSM Parameter store")

}
