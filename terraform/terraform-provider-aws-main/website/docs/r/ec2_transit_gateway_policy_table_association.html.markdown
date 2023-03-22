---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "AWS: aws_ec2_transit_gateway_policy_table_association"
description: |-
  Manages an EC2 Transit Gateway Policy Table association
---

# Resource: aws_ec2_transit_gateway_policy_table_association

Manages an EC2 Transit Gateway Policy Table association.

## Example Usage

```terraform
resource "aws_ec2_transit_gateway_policy_table_association" "example" {
  transit_gateway_attachment_id   = aws_networkmanager_transit_gateway_peering.example.transit_gateway_peering_attachment_id
  transit_gateway_policy_table_id = aws_ec2_transit_gateway_policy_table.example.id
}
```

## Argument Reference

The following arguments are supported:

* `transit_gateway_attachment_id` - (Required) Identifier of EC2 Transit Gateway Attachment.
* `transit_gateway_policy_table_id` - (Required) Identifier of EC2 Transit Gateway Policy Table.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - EC2 Transit Gateway Policy Table identifier combined with EC2 Transit Gateway Attachment identifier
* `resource_id` - Identifier of the resource
* `resource_type` - Type of the resource

## Import

`aws_ec2_transit_gateway_policy_table_association` can be imported by using the EC2 Transit Gateway Policy Table identifier, an underscore, and the EC2 Transit Gateway Attachment identifier, e.g.,

```
$ terraform import aws_ec2_transit_gateway_policy_table_association.example tgw-rtb-12345678_tgw-attach-87654321
```
