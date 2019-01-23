package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appConfig *viper.Viper
var cfgFile string

const cfgSecretShares = "secret-shares"
const cfgSecretThreshold = "secret-threshold"

const cfgMode = "mode"
const cfgModeValueAlicloudKMSOSS = "alicloud-kms-ssm"
const cfgModeValueAWSKMSSSM = "aws-kms-ssm"
const cfgModeValueGoogleCloudKMSGCS = "google-cloud-kms-gcs"
const cfgModeValueLocal = "local"

const cfgGoogleCloudKMSProject = "google-cloud-kms-project"
const cfgGoogleCloudKMSLocation = "google-cloud-kms-location"
const cfgGoogleCloudKMSKeyRing = "google-cloud-kms-key-ring"
const cfgGoogleCloudKMSCryptoKey = "google-cloud-kms-crypto-key"

const cfgGoogleCloudStorageBucket = "google-cloud-storage-bucket"
const cfgGoogleCloudStoragePrefix = "google-cloud-storage-prefix"

const cfgAWSKMSKeyID = "aws-kms-key-id"
const cfgAWSSSMKeyPrefix = "aws-ssm-key-prefix"

const cfgAlicloudKMSRegion = "alicloud-kms-region"
const cfgAlicloudKMSKeyID = "alicloud-kms-id"

const cfgAlicloudStorageEndpoint = "alicloud-storage-endpoint"
const cfgAlicloudStorageBucket = "alicloud-storage-bucket"
const cfgAlicloudStoragePrefix = "alicloud-storage-prefix"

const cfgLocalKeyDir = "local-key-dir"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vault-unsealer",
	Short: "Automates initialisation and unsealing of Hashicorp Vault.",
	Long: `This is a CLI tool to help automate the setup and management of
Hashicorp Vault.

It will continuously attempt to unseal the target Vault instance, by retrieving
unseal keys from a Google Cloud, AWS KMS keyring or local in path
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

func configIntVar(key string, defaultValue int, description string) {
	RootCmd.PersistentFlags().Int(key, defaultValue, description)
	appConfig.BindPFlag(key, RootCmd.PersistentFlags().Lookup(key))
}

func configStringVar(key, defaultValue, description string) {
	RootCmd.PersistentFlags().String(key, defaultValue, description)
	appConfig.BindPFlag(key, RootCmd.PersistentFlags().Lookup(key))
}

func init() {
	appConfig = viper.New()
	appConfig.SetEnvPrefix("vault_unsealer")
	replacer := strings.NewReplacer("-", "_")
	appConfig.SetEnvKeyReplacer(replacer)
	appConfig.AutomaticEnv()

	// SelectMode
	configStringVar(
		cfgMode,
		cfgModeValueGoogleCloudKMSGCS,
		fmt.Sprintf("Select the mode to use '%s' => Google Cloud Storage with encryption using Google KMS; '%s' => AWS SSM parameter store using AWS KMS encryption; '%s' => Alicloud KMS parameter store using Alicloud KMS encryption; %s => Use local keys in path", cfgModeValueGoogleCloudKMSGCS, cfgModeValueAWSKMSSSM, cfgModeValueAlicloudKMSOSS, cfgModeValueLocal),
	)

	// Secret config
	configIntVar(cfgSecretShares, 1, "Total count of secret shares that exist")
	configIntVar(cfgSecretThreshold, 1, "Minimum required secret shares to unseal")

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

	// Alicloud KMS flags
	configStringVar(cfgAlicloudKMSRegion, "", "The Region string of the Alicloud KMS to encrypt values")
	configStringVar(cfgAlicloudKMSKeyID, "", "The Key ID string of the Alicloud KMS key to encrypt values")

	// Alicloud Storage flags
	configStringVar(cfgAlicloudStorageEndpoint, "", "The endpoint of the Alicloud Bucket")
	configStringVar(cfgAlicloudStorageBucket, "", "The Region string of the Alicloud KMS to encrypt values")
	configStringVar(cfgAlicloudStoragePrefix, "", "The prefix to use for values store in Alicloud Storage")

	configStringVar("local-key-dir", "", "Directory of key shares in path")
}
