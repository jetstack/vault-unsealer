terraform {
  backend "s3" {
    bucket  = "vault-store"
    key     = "terraform.vault.state"
    region  = "us-east-1"
  }
}

provider "aws" {
  region = "eu-west-1"
}

resource "aws_kms_key" "vault_kms_key" {
  description             = "Vault dev kms key"
  tags {
    Environment = "${var.environment}"
    Terraform   = "true"
  }
}

resource "aws_kms_alias" "vault_kms_key_alias" {
  name          = "alias/${var.environment}-vault"
  target_key_id = "${aws_kms_key.vault_kms_key.key_id}"
}

resource "aws_ssm_parameter" "vault-unseal-0" {
  name        = "${var.environment}-vault-unseal-0"
  description = "Vault Unseal key number 0"
  type        = "SecureString"
  value       = "${var.vault_unseal_key0}"
  key_id  = "${aws_kms_key.vault_kms_key.key_id}"

  tags {
    Environment = "${var.environment}"
    Terraform   = "true"
  }
}
resource "aws_ssm_parameter" "vault-unseal-1" {
  name        = "${var.environment}-vault-unseal-1"
  description = "Vault Unseal key number 1"
  type        = "SecureString"
  value       = "${var.vault_unseal_key1}"
  key_id  = "${aws_kms_key.vault_kms_key.key_id}"

  tags {
    Environment = "${var.environment}"
    Terraform   = "true"
  }
}

resource "aws_ssm_parameter" "vault-unseal-2" {
  name        = "${var.environment}-vault-unseal-2"
  description = "Vault Unseal key number 2"
  type        = "SecureString"
  value       = "${var.vault_unseal_key2}"
  key_id  = "${aws_kms_key.vault_kms_key.key_id}"

  tags {
    Environment = "${var.environment}"
    Terraform   = "true"
  }
}

resource "aws_ssm_parameter" "vault-unseal-3" {
  name        = "${var.environment}-vault-unseal-3"
  description = "Vault Unseal key number 3"
  type        = "SecureString"
  value       = "${var.vault_unseal_key3}"
  key_id  = "${aws_kms_key.vault_kms_key.key_id}"

  tags {
    Environment = "${var.environment}"
    Terraform   = "true"
  }
}