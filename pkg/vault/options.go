package vault

import (
	"github.com/pkg/errors"
	aws "github.com/soter/vault-unsealer/pkg/kv/aws_kms"
	google "github.com/soter/vault-unsealer/pkg/kv/cloudkms"
	"github.com/spf13/pflag"
)

// That configures the vault API
type VaultOptions struct {
	KeyPrefix string

	// how many key parts exist
	SecretShares int
	// how many of these parts are needed to unseal vault  (secretThreshold <= secretShares)
	SecretThreshold int

	// should the root token be stored in the keyStore
	StoreRootToken bool

	// overwrite existing tokens
	OverwriteExisting bool

	Google *google.Options
	Aws    *aws.Options
}

func NewVaultOptions() *VaultOptions {
	return &VaultOptions{
		KeyPrefix:       "vault",
		SecretThreshold: 3,
		SecretShares:    5,
		StoreRootToken:  true,
	}
}

func (o *VaultOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.StoreRootToken, "store-root-token", o.StoreRootToken, "should the root token be stored in the key store")
	fs.BoolVar(&o.OverwriteExisting, "overwrite-existing", o.OverwriteExisting, "overwrite existing unseal keys and root tokens, possibly dangerous!")
	fs.IntVar(&o.SecretShares, "secret-shares", o.SecretShares, "Total count of secret shares that exist")
	fs.IntVar(&o.SecretThreshold, "secret-threshold", o.SecretThreshold, "Minimum required secret shares to unseal")
}

func (o *VaultOptions) Validate() []error {
	var errs []error
	if o.SecretThreshold <= 0 {
		errs = append(errs, errors.New("secret threshold must be positive"))
	}
	if o.SecretShares <= 0 {
		errs = append(errs, errors.New("secret shares must be positive"))
	}
	if o.SecretThreshold > o.SecretShares {
		errs = append(errs, errors.New("secret threshold must be less than or equal to secret shares"))
	}
	return errs
}

func (o *VaultOptions) Apply() error {
	return nil
}
