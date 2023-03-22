---
subcategory: "Network Manager"
layout: "aws"
page_title: "AWS: aws_networkmanager_attachment_accepter"
description: |-
  Terraform resource for managing an AWS NetworkManager Attachment Accepter.
---

# Resource: aws_networkmanager_attachment_accepter

Terraform resource for managing an AWS NetworkManager Attachment Accepter.

## Example Usage

### Basic Usage

```terraform
resource "aws_networkmanager_attachment_accepter" "test" {
  attachment_id   = aws_networkmanager_vpc_attachment.vpc.id
  attachment_type = aws_networkmanager_vpc_attachment.vpc.attachment_type
}
```

## Argument Reference

The following arguments are required:

* `attachment_id` - (Required) The ID of the attachment.
* `attachment_type` - The type of attachment. Valid values can be found in the [AWS Documentation](https://docs.aws.amazon.com/networkmanager/latest/APIReference/API_ListAttachments.html#API_ListAttachments_RequestSyntax)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attachment_policy_rule_number` - The policy rule number associated with the attachment.
* `core_network_arn` - The ARN of a core network.
* `core_network_id` - The id of a core network.
* `edge_location` - The Region where the edge is located.
* `owner_account_id` - The ID of the attachment account owner.
* `resource_arn` - The attachment resource ARN.
* `segment_name` - The name of the segment attachment.
* `state` - The state of the attachment.
