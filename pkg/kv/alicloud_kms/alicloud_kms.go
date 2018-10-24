package alicloud_kms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"

	"github.com/jetstack/vault-unsealer/pkg/kv"
)

type alicloudKms struct {
	kmsClient 	*kms.Client
	store      	kv.Service
	keyId 		string
}

var _ kv.Service = &alicloudKms{}

func New(store kv.Service, keyId, regionId, alicloudAccessKey, alicloudSecretKey string) (kv.Service, error) {
	client, err := kms.NewClientWithAccessKey(regionId, alicloudAccessKey, alicloudSecretKey)

	if err != nil {
		return nil, fmt.Errorf("error creating alicloud kms client: %s", err.Error())
	}

	return &alicloudKms{
		store:   store,
		kmsClient:     client,
		keyId: keyId,
	}, nil
}

func (a *alicloudKms) encrypt(s []byte) ([]byte, error) {
	requestStruct := kms.CreateEncryptRequest()
	requestStruct.KeyId = a.keyId
	requestStruct.RpcRequest.SetScheme("HTTPS")
	requestStruct.Plaintext = string(s[:])

	genresp, err := a.kmsClient.Encrypt(requestStruct)

	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %s", err.Error())
	}

	return []byte(genresp.CiphertextBlob), nil
}

func (a *alicloudKms) decrypt(s []byte) ([]byte, error) {
	derequestStruct := kms.CreateDecryptRequest()
	derequestStruct.CiphertextBlob = string(s[:])
	derequestStruct.RpcRequest.SetScheme("HTTPS")

	degenresp, err := a.kmsClient.Decrypt(derequestStruct)

	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %s", err.Error())
	}

	return []byte(degenresp.Plaintext), nil
}

func (a *alicloudKms) Get(key string) ([]byte, error) {
	cipherText, err := a.store.Get(key)

	if err != nil {
		return nil, err
	}

	return a.decrypt(cipherText)
}

func (a *alicloudKms) Set(key string, val []byte) error {
	cipherText, err := a.encrypt(val)

	if err != nil {
		return err
	}

	return a.store.Set(key, cipherText)
}

func (a *alicloudKms) Test(key string) error {
	// TODO: Implement a test if a Set is likely to work, Alicloud doesn't seemt to provide a dry-run on the parameter store
	return nil
}
