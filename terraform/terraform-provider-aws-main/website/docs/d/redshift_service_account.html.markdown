---
subcategory: "Redshift"
layout: "aws"
page_title: "AWS: aws_redshift_service_account"
description: |-
  Get AWS Redshift Service Account for storing audit data in S3.
---

# Data Source: aws_redshift_service_account

Use this data source to get the Account ID of the [AWS Redshift Service Account](http://docs.aws.amazon.com/redshift/latest/mgmt/db-auditing.html#db-auditing-enable-logging)
in a given region for the purpose of allowing Redshift to store audit data in S3.

## Example Usage

```terraform
data "aws_redshift_service_account" "main" {}

resource "aws_s3_bucket" "bucket" {
  bucket        = "tf-redshift-logging-test-bucket"
  force_destroy = true
}

resource "aws_s3_bucket_policy" "allow_audit_logging" {
  bucket = aws_s3_bucket.bucket.id
  policy = <<EOF
{
	"Version": "2008-10-17",
	"Statement": [
		{
            "Sid": "Put bucket policy needed for audit logging",
            "Effect": "Allow",
            "Principal": {
		        "AWS": "${data.aws_redshift_service_account.main.arn}"
            },
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::tf-redshift-logging-test-bucket/*"
        },
        {
            "Sid": "Get bucket policy needed for audit logging ",
            "Effect": "Allow",
            "Principal": {
		        "AWS": "${data.aws_redshift_service_account.main.arn}"
            },
            "Action": "s3:GetBucketAcl",
            "Resource": "arn:aws:s3:::tf-redshift-logging-test-bucket"
        }
	]
}
EOF
}
```

## Argument Reference

* `region` - (Optional) Name of the region whose AWS Redshift account ID is desired.
Defaults to the region from the AWS provider configuration.

## Attributes Reference

* `id` - The ID of the AWS Redshift service account in the selected region.
* `arn` - The ARN of the AWS Redshift service account in the selected region.
