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

	"github.com/spf13/viper"

	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv"
	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv/aws_kms"
	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv/aws_ssm"
	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv/cloudkms"
	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv/gcs"
)

func kvStoreForConfig(cfg *viper.Viper) (kv.Service, error) {

	if cfg.GetString(cfgMode) == cfgModeValueGoogleCloudKMSGCS {
		g, err := gcs.New(
			cfg.GetString(cfgGoogleCloudStorageBucket),
			cfg.GetString(cfgGoogleCloudStoragePrefix),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud storage kv store: %s", err.Error())
		}

		kms, err := cloudkms.New(g,
			cfg.GetString(cfgGoogleCloudKMSProject),
			cfg.GetString(cfgGoogleCloudKMSLocation),
			cfg.GetString(cfgGoogleCloudKMSKeyRing),
			cfg.GetString(cfgGoogleCloudKMSCryptoKey),
		)

		if err != nil {
			return nil, fmt.Errorf("error creating google cloud kms kv store: %s", err.Error())
		}

		return kms, nil
	}

	if cfg.GetString(cfgMode) == cfgModeValueAWSKMSSSM {
		ssm, err := aws_ssm.New(cfg.GetString(cfgAWSSSMKeyPrefix))
		if err != nil {
			return nil, fmt.Errorf("error creating AWS SSM kv store: %s", err.Error())
		}

		kms, err := aws_kms.New(ssm, cfg.GetString(cfgAWSKMSKeyID))
		if err != nil {
			return nil, fmt.Errorf("error creating AWS KMS ID kv store: %s", err.Error())
		}

		return kms, nil
	}

	return nil, fmt.Errorf("Unsupported backend mode: '%s'", cfg.GetString(cfgMode))
}
