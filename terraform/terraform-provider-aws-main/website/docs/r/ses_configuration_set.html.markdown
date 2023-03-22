---
subcategory: "SES (Simple Email)"
layout: "aws"
page_title: "AWS: aws_ses_configuration_set"
description: |-
  Provides an SES configuration set
---

# Resource: aws_ses_configuration_set

Provides an SES configuration set resource.

## Example Usage

```terraform
resource "aws_ses_configuration_set" "test" {
  name = "some-configuration-set-test"
}
```

### Require TLS Connections

```terraform
resource "aws_ses_configuration_set" "test" {
  name = "some-configuration-set-test"

  delivery_options {
    tls_policy = "Require"
  }
}
```

## Argument Reference

The following argument is required:

* `name` - (Required) Name of the configuration set.

The following argument is optional:

* `delivery_options` - (Optional) Whether messages that use the configuration set are required to use TLS. See below.
* `reputation_metrics_enabled` - (Optional) Whether or not Amazon SES publishes reputation metrics for the configuration set, such as bounce and complaint rates, to Amazon CloudWatch. The default value is `false`.
* `sending_enabled` - (Optional) Whether email sending is enabled or disabled for the configuration set. The default value is `true`.
* `tracking_options` - (Optional) Domain that is used to redirect email recipients to an Amazon SES-operated domain. See below. **NOTE:** This functionality is best effort.

### delivery_options

* `tls_policy` - (Optional) Whether messages that use the configuration set are required to use Transport Layer Security (TLS). If the value is `Require`, messages are only delivered if a TLS connection can be established. If the value is `Optional`, messages can be delivered in plain text if a TLS connection can't be established. Valid values: `Require` or `Optional`. Defaults to `Optional`.

### tracking_options

* `custom_redirect_domain` - (Optional) Custom subdomain that is used to redirect email recipients to the Amazon SES event tracking domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - SES configuration set ARN.
* `id` - SES configuration set name.
* `last_fresh_start` - Date and time at which the reputation metrics for the configuration set were last reset. Resetting these metrics is known as a fresh start.

## Import

SES Configuration Sets can be imported using their `name`, e.g.,

```
$ terraform import aws_ses_configuration_set.test some-configuration-set-test
```
