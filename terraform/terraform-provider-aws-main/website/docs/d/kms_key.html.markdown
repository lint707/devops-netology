---
subcategory: "KMS (Key Management)"
layout: "aws"
page_title: "AWS: aws_kms_key"
description: |-
  Get information on a AWS Key Management Service (KMS) Key
---

# aws_kms_key

Use this data source to get detailed information about
the specified KMS Key with flexible key id input.
This can be useful to reference key alias
without having to hard code the ARN as input.

## Example Usage

```terraform
data "aws_kms_key" "by_alias" {
  key_id = "alias/my-key"
}

data "aws_kms_key" "by_id" {
  key_id = "1234abcd-12ab-34cd-56ef-1234567890ab"
}

data "aws_kms_key" "by_alias_arn" {
  key_id = "arn:aws:kms:us-east-1:111122223333:alias/my-key"
}

data "aws_kms_key" "by_key_arn" {
  key_id = "arn:aws:kms:us-east-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
}
```

## Argument Reference

* `key_id` - (Required) Key identifier which can be one of the following format:
    * Key ID. E.g: `1234abcd-12ab-34cd-56ef-1234567890ab`
    * Key ARN. E.g.: `arn:aws:kms:us-east-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab`
    * Alias name. E.g.: `alias/my-key`
    * Alias ARN: E.g.: `arn:aws:kms:us-east-1:111122223333:alias/my-key`
* `grant_tokens` - (Optional) List of grant tokens

## Attributes Reference

* `id`: The globally unique identifier for the key
* `arn`: The Amazon Resource Name (ARN) of the key
* `aws_account_id`: The twelve-digit account ID of the AWS account that owns the key
* `creation_date`: The date and time when the key was created
* `deletion_date`: The date and time after which AWS KMS deletes the key. This value is present only when `key_state` is `PendingDeletion`, otherwise this value is 0
* `description`: The description of the key.
* `enabled`: Specifies whether the key is enabled. When `key_state` is `Enabled` this value is true, otherwise it is false
* `expiration_model`: Specifies whether the Key's key material expires. This value is present only when `origin` is `EXTERNAL`, otherwise this value is empty
* `key_manager`: The key's manager
* `key_state`: The state of the key
* `key_usage`: Specifies the intended use of the key
* `customer_master_key_spec`: Specifies whether the key contains a symmetric key or an asymmetric key pair and the encryption algorithms or signing algorithms that the key supports
* `multi_region`: Indicates whether the KMS key is a multi-Region (`true`) or regional (`false`) key.
* `multi_region_configuration`: Lists the primary and replica keys in same multi-Region key. Present only when the value of `multi_region` is `true`.
* `origin`: When this value is `AWS_KMS`, AWS KMS created the key material. When this value is `EXTERNAL`, the key material was imported from your existing key management infrastructure or the CMK lacks key material
* `valid_to`: The time at which the imported key material expires. This value is present only when `origin` is `EXTERNAL` and whose `expiration_model` is `KEY_MATERIAL_EXPIRES`, otherwise this value is 0

The `multi_region_configuration` object supports the following:

* `multi_region_key_type`: Indicates whether the KMS key is a `PRIMARY` or `REPLICA` key.
* `primary_key`: The key ARN and Region of the primary key. This is the current KMS key if it is the primary key.
* `replica_keys`: The key ARNs and Regions of all replica keys. Includes the current KMS key if it is a replica key.

The `primary_key` and `replica_keys` objects support the following:

* `arn`: The key ARN of a primary or replica key of a multi-Region key.
* `region`: The AWS Region of a primary or replica key in a multi-Region key.
