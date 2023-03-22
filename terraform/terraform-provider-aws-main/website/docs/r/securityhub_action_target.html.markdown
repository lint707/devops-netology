---
subcategory: "Security Hub"
layout: "aws"
page_title: "AWS: aws_securityhub_action_target"
description: |-
  Creates Security Hub custom action.
---

# Resource: aws_securityhub_action_target

Creates Security Hub custom action.

## Example Usage

```terraform
resource "aws_securityhub_account" "example" {}

resource "aws_securityhub_action_target" "example" {
  depends_on  = [aws_securityhub_account.example]
  name        = "Send notification to chat"
  identifier  = "SendToChat"
  description = "This is custom action sends selected findings to chat"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The description for the custom action target.
* `identifier` - (Required) The ID for the custom action target.
* `description` - (Required) The name of the custom action target.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the Security Hub custom action target.

## Import

Security Hub custom action can be imported using the action target ARN e.g.,

```sh
$ terraform import aws_securityhub_action_target.example arn:aws:securityhub:eu-west-1:312940875350:action/custom/a
```
