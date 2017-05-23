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

package gcs

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"

	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv"
)

type gcsStorage struct {
	cl     *storage.Client
	bucket string
	prefix string
}

func New(bucket, prefix string) (kv.Service, error) {
	cl, err := storage.NewClient(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error creating gcs client: %s", err.Error())
	}

	return &gcsStorage{cl, bucket, prefix}, nil
}

func (g *gcsStorage) Set(key, val string) error {
	ctx := context.Background()
	n := objectNameWithPrefix(g.prefix, key)
	if _, err := g.cl.Bucket(g.bucket).Object(n).NewWriter(ctx).Write([]byte(val)); err != nil {
		return fmt.Errorf("error writing key '%s' to gcs bucket '%s'", n, g.bucket)
	}
	return nil
}

func (g *gcsStorage) Get(key string) (string, error) {
	ctx := context.Background()
	n := objectNameWithPrefix(g.prefix, key)

	r, err := g.cl.Bucket(g.bucket).Object(n).NewReader(ctx)

	if err != nil {
		return "", fmt.Errorf("error getting object for key '%s': %s", n, err.Error())
	}

	b, err := ioutil.ReadAll(r)

	if err != nil {
		return "", fmt.Errorf("error reading object with key '%s': %s", n, err.Error())
	}

	return string(b), nil
}

func objectNameWithPrefix(prefix, key string) string {
	return fmt.Sprintf("%s%s", prefix, key)
}
