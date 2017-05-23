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

package cloudkms

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"

	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv"
)

// googleKms is an implementation of the kv.Service interface, that encrypts
// and decrypts data using Google Cloud KMS before storing into another kv
// backend.
type googleKms struct {
	svc     *cloudkms.Service
	store   kv.Service
	keyPath string
}

var _ kv.Service = &googleKms{}

func New(store kv.Service, project, location, keyring, cryptoKey string) (kv.Service, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)

	if err != nil {
		return nil, fmt.Errorf("error creating google client: %s", err.Error())
	}

	kmsService, err := cloudkms.New(client)

	if err != nil {
		return nil, fmt.Errorf("error creating google kms service client: %s", err.Error())
	}

	return &googleKms{
		store:   store,
		svc:     kmsService,
		keyPath: fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", project, location, keyring, cryptoKey),
	}, nil
}

func (g *googleKms) encrypt(s string) (string, error) {
	resp, err := g.svc.Projects.Locations.KeyRings.CryptoKeys.Encrypt(g.keyPath, &cloudkms.EncryptRequest{
		Plaintext: s,
	}).Do()

	if err != nil {
		return "", fmt.Errorf("error encrypting data: %s", err.Error())
	}

	return resp.Ciphertext, nil
}

func (g *googleKms) decrypt(s string) (string, error) {
	resp, err := g.svc.Projects.Locations.KeyRings.CryptoKeys.Decrypt(g.keyPath, &cloudkms.DecryptRequest{
		Ciphertext: s,
	}).Do()

	if err != nil {
		return "", fmt.Errorf("error decrypting data: %s", err.Error())
	}

	return resp.Plaintext, nil
}

func (g *googleKms) Get(key string) (string, error) {
	cipherText, err := g.store.Get(key)

	if err != nil {
		return "", err
	}

	return g.decrypt(cipherText)
}

func (g *googleKms) Set(key, val string) error {
	cipherText, err := g.encrypt(val)

	if err != nil {
		return err
	}

	return g.store.Set(key, cipherText)
}
