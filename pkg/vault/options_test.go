package vault

import (
	"testing"

	aggregator "github.com/appscode/go/util/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestVaultOptions_Validate(t *testing.T) {
	testData := []struct {
		testName    string
		opts        *VaultOptions
		expectedErr error
	}{
		{
			"secret threshold is zero, validation failed",
			&VaultOptions{
				SecretShares:    1,
				SecretThreshold: 0,
			},
			errors.New("secret threshold must be positive"),
		},
		{
			"secret threshold > secret shares, validation failed",
			&VaultOptions{
				SecretShares:    1,
				SecretThreshold: 2,
			},
			errors.New("secret threshold must be less than or equal to secret shares"),
		},
		{
			"validation successful",
			&VaultOptions{
				SecretShares:    10,
				SecretThreshold: 2,
			},
			nil,
		},
	}

	for _, test := range testData {
		t.Run(test.testName, func(t *testing.T) {
			errs := test.opts.Validate()
			if test.expectedErr != nil {
				assert.EqualError(t, aggregator.NewAggregate(errs), test.expectedErr.Error())
			} else {
				assert.Nil(t, errs)
			}
		})
	}
}
