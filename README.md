# Vault-unsealer

This project aims to make it easier to automate the secure initializing and unsealing of a Vault
server.

## Usage

```
Automates initialisation and unsealing of Hashicorp Vault.

Usage:
  vault-unsealer [flags]

Flags:
      --alsologtostderr                  log to standard error as well as files
      --aws.kms-key-id string            The ID or ARN of the AWS KMS key to encrypt values
      --aws.ssm-key-prefix string        The Key Prefix for SSM Parameter store
      --ca-cert-file string              Path to the ca cert file that will be used to verify self signed vault server certificate
      --google.kms-crypto-key string     The name of the Google Cloud KMS crypto key to use
      --google.kms-key-ring string       The name of the Google Cloud KMS key ring to use
      --google.kms-location string       The Google Cloud KMS location to use (eg. 'global', 'europe-west1')
      --google.kms-project string        The Google Cloud KMS project to use
      --google.storage-bucket string     The name of the Google Cloud Storage bucket to store values in
      --google.storage-prefix string     The prefix to use for values store in Google Cloud Storage
  -h, --help                             help for vault-unsealer
      --insecure-tls                     To skip tls verification when communicating with vault server
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --mode string                      Select the mode to use 'google-cloud-kms-gcs' => Google Cloud Storage with encryption using Google KMS; 'aws-kms-ssm' => AWS SSM parameter store using AWS KMS
      --overwrite-existing               overwrite existing unseal keys and root tokens, possibly dangerous!
      --retry-period duration            How often to attempt to unseal the vault instance (default 10s)
      --secret-shares int                Total count of secret shares that exist (default 5)
      --secret-threshold int             Minimum required secret shares to unseal (default 3)
      --stderrthreshold severity         logs at or above this threshold go to stderr
      --store-root-token                 should the root token be stored in the key store (default true)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

```
