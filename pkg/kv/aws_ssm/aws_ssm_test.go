package aws_ssm

import (
	"os"
	"testing"
)

func TestAWSIntegration(t *testing.T) {
	region := os.Getenv("AWS_REGION")

	if region == "" {
		t.Skip("Skip AWS integration tests: not environment variable 'AWS_REGION' specified")
	}

	payloadKey := "test123"
	payloadValue := "payload123"

	a, err := New("test-integration-")
	if err != nil {
		t.Errorf("Unexpected error creating SSM kv: %s", err)
	}

	err = a.Set(payloadKey, []byte(payloadValue))
	if err != nil {
		t.Errorf("Unexpected error storing value in SSM kv: %s", err)
	}

	out, err := a.Get("test123")
	if err != nil {
		t.Errorf("Unexpected error storing value in SSM kv: %s", err)
	}

	if exp, act := payloadValue, string(out); exp != act {
		t.Errorf("Unexpected decrypt output: exp=%s act=%s", exp, act)
	}

	_, err = a.Get("test-not-existing")
	if err == nil {
		t.Errorf("Expected error getting a non existing key")
	}

}
