package aws_kms

import (
	"testing"

	aggregator "github.com/appscode/go/util/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func getValidationError() []error {
	var errs []error
	errs = append(errs, errors.New("aws kms key id must be non-empty"))
	return errs
}

func TestOptions_Validate(t *testing.T) {
	testData := []struct {
		testName    string
		opts        *Options
		expectedErr error
	}{
		{
			"aws key id provided, validation successful",
			&Options{
				"test-key",
				"",
			},
			nil,
		},
		{
			"aws key id not provided, validation failed",
			&Options{
				"",
				"",
			},
			aggregator.NewAggregate(getValidationError()),
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
