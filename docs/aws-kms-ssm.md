# Unsealing Vault using encrypted keys from AWS SSM

vault-unsealer can use AWS KMS and SSM to store and retrieve encrypted Vault keys.

vault-unsealer can be used to initialize a vault and auto store its keys in SSM.
[See Initializing a Vault for more info](#initializing-a-vault)

## Pre-existing Vault keys

1. Create or use a KMS key in the region you want. AWS Console Encryption Keys Page.

2. Note the alias name of the key, for example: `alias/vault`.

3. Note the key UUID at the end of the ARN, for example: `arn:aws:kms:<region>:<aws-account-id>:key/<key-uuid>`.

4. Add the following IAM permissions to the IAM user/role which vault-unsealer will be using.


<b>
:warning: WARNING<br />

    If running in Kubernetes, its strongly suggested that use a project such as kube2iam
    to limit access to a specific pod for accessing the keys via a separate IAM role.
</b>

Using the policy below on nodes role will allow any pod to access the root token and unseal keys stored in SSM with KMS.

Full IAM policy (Change the place holders `<region>`, `<aws-account-id>` and `<key-uuid>`)

```javascript
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VaultUnsealerReadSSMParameters",
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameters"
            ],
            "Resource": [
                "arn:aws:ssm:<region>:<aws-account-id>:parameter/*"
            ]
        },
        {
            "Sid": "VaultUnsealerWriteSSMParameters",
            "Effect": "Allow",
            "Action": [
                "ssm:PutParameter",
                "ssm:DeleteParameter"
            ],
            "Resource": [
                "arn:aws:ssm:<region>:<aws-account-id>:parameter/vault-unsealer-*"
            ]
        },
        {
            "Sid": "VaultUnsealerGetKMS",
            "Effect": "Allow",
            "Action": [
                "kms:Get*",
                "kms:ListKeys",
                "kms:ListAliases"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Sid": "VaultUnsealerEncryptDecryptKms",
            "Effect": "Allow",
            "Action": [
                "kms:DescribeKey",
                "kms:Encrypt",
                "kms:Decrypt"
            ],
            "Resource": [
                "arn:aws:kms:<region>:<aws-account-id>:key/<key-uuid>"
            ]
        }
    ]
}
```

## Setting up an existing vault (already initialized)

Using the vault root token and unseal keys you can setup the SSM Parameter store with the following commands:

1.Export keys and prefix for SSM:

```bash
    export PREFIX=<your-prefix>- \
    KMS_KEY_ID=<kms-key-id> \
    ROOT_KEY=<vault-root-token> \
    UNSEAL0=<vault-unseal-key-1> \
    UNSEAL1=<vault-unseal-key-2> \
    UNSEAL2=<vault-unseal-key-3> \
    UNSEAL3=<vault-unseal-key-4> \
    UNSEAL4=<vault-unseal-key-5>
```

2.Encrypt and Put SSM Parameters to AWS:

```bash
$ mkdir -p /tmp/vault
  echo $ROOT_KEY > /tmp/vault/root-key
  echo $UNSEAL0 > /tmp/vault/unseal0
  echo $UNSEAL1 > /tmp/vault/unseal1
  echo $UNSEAL2 > /tmp/vault/unseal2
  echo $UNSEAL3 > /tmp/vault/unseal3
  echo $UNSEAL4 > /tmp/vault/unseal4

  echo "Encrypting Vault root token"
  aws kms encrypt --key-id $KMS_KEY_ID \
    --plaintext fileb:///tmp/vault/root-key \
    --output text \
    --query CiphertextBlob > /tmp/vault/root.enc

  echo "Creating a new SSM paramter key ${PREFIX}vault-root for Vault root token"
  aws ssm put-parameter --name "${PREFIX}vault-unseal-root" \
    --value "$(cat /tmp/vault/root.enc )" \
    --type String

  for i in {0..4}; do
    echo "Encrypting unseal${i} key"
    aws kms encrypt --key-id $KMS_KEY_ID \
      --encryption-context "Tool=vault-unsealer" \
      --plaintext fileb:///tmp/vault/unseal${i} \
      --output text \
      --query CiphertextBlob > /tmp/vault/unseal${i}.enc

    echo "Creating a new SSM paramter key ${PREFIX}vault-unseal-${i}"
    aws ssm put-parameter --name "${PREFIX}vault-unseal-${i}" \
      --value "$(cat /tmp/vault/unseal${i}.enc)" \
      --type String
  done

  rm -rf /tmp/vault
```

## Initializing a Vault

If your vault is not yet initialized you can initialized it using the parameter store as follow:

```bash
export AWS_REGION=<region>
```

Run the command aws kms list-aliases to get a list of the kms keys you need, you must use the alias name.

```javascript
{
    "Aliases": [
        {
            "AliasName": "alias/MyKmsKey",
            "AliasArn": "arn:aws:kms:us-west-2:123456789012:alias/myKMSKey",
            "TargetKeyId": "4e4ad123-20cf-4ffe-a55f-edd96ca41bef"
        }
    ]
}
```

Note: alias key must have the prefix alias/. you may use any number of secrets shares and threshold for your needs. (default is 1 secret share and 1 threshold)

```bash
vault-unsealer init \
  --mode aws-kms-ssm \
  --aws-kms-key-id <alias/kms-alias-key> \
  --aws-ssm-key-prefix <your-prefix>- \
  --secret-shares 5 \
  --secret-threshold 3
```

```console
INFO[0015] root token stored in key store key=vault-root
```

This will create 6 keys in the AWS SSM:<br/>
* vault-root
* vault-unsealer-0
* vault-unsealer-1
* vault-unsealer-2
* vault-unsealer-3
* vault-unsealer-4

## Test vault-unsealer on AWS with KMS and SSM

```bash
vault-unsealer  unseal \
  --mode aws-kms-ssm \
  --aws-kms-key-id  "alias/`<key-alias>`" \
  --aws-ssm-key-prefix `<ssm-key-prefix>`- \
  --secret-shares 5 \
  --secret-threshold 3
```

```console
INFO[0000] checking if vault is sealed...
INFO[0000] vault sealed: true
INFO[0002] successfully unsealed vault
```

## Getting the Root token and Unseal keys from SSM and Decrypt with KMS.

Since you've generated the the keys and root token using vault-unsealer you can use the following commands to get the unseal keys and root token:

Make sure to change the region / AWS CLI profile to your needs.

Note that you don't have to pass the KMS key id since it's in the metadata of the stored parameter on SSM.

```bash
export AWS_DEFAULT_REGION=us-west-2
export AWS_PROFILE=dev
echo "Fetching Vault unseal keys and root token from AWS..."

aws ssm get-parameters --names kubernetes-vault-root | jq -r '.Parameters[].Value'  | base64 -D > /tmp/root-token
ROOT_TOKEN=$(aws kms decrypt --ciphertext-blob fileb:///tmp/root-token --encryption-context Tool=vault-unsealer | jq -r '.Plaintext' | base64 -D)

for i in {0..4}; do
aws ssm get-parameters --names kubernetes-vault-unseal-${i} | jq -r '.Parameters[].Value'  | base64 -D > /tmp/unseal-${i}
echo "Unseal Key $((i+1)): $(aws kms decrypt --ciphertext-blob fileb:///tmp/unseal-${i} --encryption-context Tool=vault-unsealer | jq -r '.Plaintext' | base64 -D)"
done

echo "Initial Root token: ${ROOT_TOKEN}"
rm /tmp/unseal*
rm /tmp/root-token
echo "Done."
```