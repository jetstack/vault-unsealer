package aws_ssm

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gitlab.jetstack.net/jetstack-experimental/vault-unsealer/pkg/kv"
)

type awsSSM struct {
	ssmService *ssm.SSM

	keyPrefix string
}

var _ kv.Service = &awsSSM{}

func New(keyPrefix string) (kv.Service, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &awsSSM{
		ssmService: ssm.New(sess),
		keyPrefix:  keyPrefix,
	}, nil
}

func (a *awsSSM) Get(key string) ([]byte, error) {
	out, err := a.ssmService.GetParameters(&ssm.GetParametersInput{
		Names: []*string{
			aws.String(a.name(key)),
		},
	})
	if err != nil {
		return []byte{}, err
	}

	if len(out.Parameters) < 1 {
		return []byte{}, fmt.Errorf("key '%s' not found", key)
	}

	return base64.StdEncoding.DecodeString(*out.Parameters[0].Value)
}

func (a *awsSSM) name(key string) string {
	return fmt.Sprintf("%s%s", a.keyPrefix, key)
}

func (a *awsSSM) Set(key string, val []byte) error {
	_, err := a.ssmService.PutParameter(&ssm.PutParameterInput{
		Description: aws.String("vault-unsealer"),
		Name:        aws.String(a.name(key)),
		Overwrite:   aws.Bool(true), // TODO: Potentally dangerous: Overwriting of the seal key == loosing vault backend
		Value:       aws.String(base64.StdEncoding.EncodeToString(val)),
		Type:        aws.String("String"),
	})
	return err
}

func (g *awsSSM) Test(key string) error {
	// TODO: Implement a test if a Set is likely to work, AWS doesn't seemt to provide a dry-run on the parameter store
	return nil
}
