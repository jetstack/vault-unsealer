package aws_kms

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Options struct {
	KmsKeyID string

	// TODO: should make it auto generated
	SsmKeyPrefix string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.KmsKeyID, "aws.kms-key-id", o.KmsKeyID, "The ID or ARN of the AWS KMS key to encrypt values")
	fs.StringVar(&o.SsmKeyPrefix, "aws.ssm-key-prefix", o.SsmKeyPrefix, "The Key Prefix for SSM Parameter store")
}

func (o *Options) Validate() []error {
	var errs []error
	if o.KmsKeyID == "" {
		errs = append(errs, errors.New("aws kms key id must be non-empty"))
	}
	return errs
}

func (o *Options) Apply() error {
	return nil
}
