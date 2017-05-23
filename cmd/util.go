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

	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv"
	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv/cloudkms"
	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv/gcs"
)

func kvStoreForFlags(cfg kvCfg) (kv.Service, error) {
	g, err := gcs.New(cfg.googleCloudStorageBucket, cfg.googleCloudStoragePrefix)

	if err != nil {
		return nil, fmt.Errorf("error creating google cloud storage kv store: %s", err.Error())
	}

	kms, err := cloudkms.New(g,
		cfg.googleCloudKMSProject,
		cfg.googleCloudKMSLocation,
		cfg.googleCloudKMSKeyRing,
		cfg.googleCloudKMSCryptoKey)

	if err != nil {
		return nil, fmt.Errorf("error creating google cloud kms kv store: %s", err.Error())
	}

	return kms, nil
}
