# Vault-unsealer

This project aims to make it easier to automate the secure unsealing of a Vault
server.

## Usage

```
This is a CLI tool to help automate the setup and management of
Hashicorp Vault.

It will continuously attempt to unseal the target Vault instance, by retrieving
unseal keys from a Google Cloud KMS keyring.

Usage:
  vault-unsealer [command]

Available Commands:
  init        Initialise the target Vault instance
  unseal      A brief description of your command

Flags:
      --google-cloud-kms-crypto-key string   The name of the Google Cloud KMS crypt key to use
      --google-cloud-kms-key-ring string     The name of the Google Cloud KMS key ring to use
      --google-cloud-kms-location string     The Google Cloud KMS location to use (eg. 'global', 'europe-west1')
      --google-cloud-kms-project string      The Google Cloud KMS project to use
      --google-cloud-storage-bucket string   The name of the Google Cloud Storage bucket to store values in
      --google-cloud-storage-prefix string   The prefix to use for values store in Google Cloud Storage
  -h, --help                                 help for vault-unsealer

Use "vault-unsealer [command] --help" for more information about a command.
```