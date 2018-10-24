package alicloud_oss

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"

	"github.com/jetstack/vault-unsealer/pkg/kv"
)

type ossStorage struct {
	cl     		*oss.Client
	endpoint 	string
	bucket 		string
	prefix 		string
}

func New(endpoint, bucket, prefix, alicloudAccessKey, alicloudSecretKey string) (kv.Service, error) {
	cl, err := oss.New(endpoint, alicloudAccessKey, alicloudSecretKey)

	lsRes, _ := cl.ListBuckets()
	for _, bucket := range lsRes.Buckets {
		fmt.Println("Buckets:", bucket.Name)

	}


	if err != nil {
		return nil, fmt.Errorf("error creating oss client: %s", err.Error())
	}

	return &ossStorage{cl, endpoint, bucket, prefix}, nil
}

func (o *ossStorage) Set(key string, val []byte) error {
	n := objectNameWithPrefix(o.prefix, key)
	bucket, err := o.cl.Bucket(o.bucket)
	if err != nil {
		return fmt.Errorf("error writing key '%s' to oss bucket '%s'", n, o.bucket)
	}

	err = bucket.PutObject(n, bytes.NewReader(val))
	if err != nil {
		return fmt.Errorf("error writing key '%s' to oss bucket '%s'", n, o.bucket)
	}

	return nil
}

func (o *ossStorage) Get(key string) ([]byte, error) {
	n := objectNameWithPrefix(o.prefix, key)
	b := new(bytes.Buffer)

	bucket, err := o.cl.Bucket(o.bucket)
	if err != nil {
		return nil, fmt.Errorf("error writing key '%s' to oss bucket '%s'", n, o.bucket)
	}

	lsRes, _ := o.cl.ListBuckets()
	for _, bucket := range lsRes.Buckets {
		fmt.Println("Buckets:", bucket.Name)

	}


	body, err := bucket.GetObject(n)
	if err != nil {
		if err != nil {
			return nil, kv.NewNotFoundError("error getting object for key '%s': %s", n, err.Error())
		}
		//return nil, fmt.Errorf("error writing key '%s' to oss bucket '%s'", n, o.bucket)
	}

	io.Copy(b, body)
	body.Close()

	return b.Bytes(), nil
}

func objectNameWithPrefix(prefix, key string) string {
	return fmt.Sprintf("%s%s", prefix, key)
}

func (o *ossStorage) Test(key string) error {
	// TODO: Implement me properly
	return nil
}
