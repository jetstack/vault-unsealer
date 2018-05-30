package worker

import (
	"time"

	"github.com/pkg/errors"
	aws "github.com/soter/vault-unsealer/pkg/kv/aws_kms"
	google "github.com/soter/vault-unsealer/pkg/kv/cloudkms"
	"github.com/soter/vault-unsealer/pkg/vault"
	"github.com/spf13/pflag"
)

const (
	ModeGoogleCloudKmsGCS = "google-cloud-kms-gcs"
	ModeAwsKmsSsm         = "aws-kms-ssm"
)

type WorkerOptions struct {
	// retry period to try initializing and unsealing
	ReTryPeriod time.Duration

	// Select the mode to use
	// 	- 'google-cloud-kms-gcs' => Google Cloud Storage with encryption using Google KMS
	// 	- 'aws-kms-ssm' => AWS SSM parameter store using AWS KMS encryption
	Mode string

	// ca cert file for vault api client, if vault used a self signed certificate
	CaCertFile string

	// If InSecureTLS true, then it will skip tls verification when communicating with vault server
	InSecureTLS bool

	Vault  *vault.VaultOptions
	Google *google.Options
	Aws    *aws.Options
}

func NewWorkerOptions() *WorkerOptions {
	return &WorkerOptions{
		ReTryPeriod: 10 * time.Second,
		Vault:       vault.NewVaultOptions(),
		Google:      google.NewOptions(),
		Aws:         aws.NewOptions(),
	}
}

func (o *WorkerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Mode, "mode", o.Mode, "Select the mode to use 'google-cloud-kms-gcs' => Google Cloud Storage with encryption using Google KMS; 'aws-kms-ssm' => AWS SSM parameter store using AWS KMS")
	fs.DurationVar(&o.ReTryPeriod, "retry-period", o.ReTryPeriod, "How often to attempt to unseal the vault instance")
	fs.StringVar(&o.CaCertFile, "ca-cert-file", o.Mode, "Path to the ca cert file that will be used to verify self signed vault server certificate")
	fs.BoolVar(&o.InSecureTLS, "insecure-tls", o.InSecureTLS, "To skip tls verification when communicating with vault server")

	o.Vault.AddFlags(fs)
	o.Google.AddFlags(fs)
	o.Aws.AddFlags(fs)
}

func (o *WorkerOptions) Validate() []error {
	var errs []error
	if o.Mode != ModeGoogleCloudKmsGCS && o.Mode != ModeAwsKmsSsm {
		errs = append(errs, errors.New("invalid mode"))
	}

	errs = append(errs, o.Vault.Validate()...)

	if o.Mode == ModeGoogleCloudKmsGCS {
		errs = append(errs, o.Google.Validate()...)
	}
	if o.Mode == ModeAwsKmsSsm {
		errs = append(errs, o.Aws.Validate()...)
	}

	return errs
}

func (o *WorkerOptions) Apply() error {
	return nil
}
