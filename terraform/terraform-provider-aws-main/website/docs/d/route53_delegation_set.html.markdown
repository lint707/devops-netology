---
subcategory: "Route 53"
layout: "aws"
page_title: "AWS: aws_route53_delegation_set"
description: |-
    Provides details about a specific Route 53 Delegation Set
---

# Data Source: aws_route53_delegation_set

`aws_route53_delegation_set` provides details about a specific Route 53 Delegation Set.

This data source allows to find a list of name servers associated with a specific delegation set.

## Example Usage

The following example shows how to get a delegation set from its id.

```terraform
data "aws_route53_delegation_set" "dset" {
  id = "MQWGHCBFAKEID"
}
```

## Argument Reference


* `id` - (Required) The Hosted Zone id of the desired delegation set.

The following attribute is additionally exported:

* `arn` - The Amazon Resource Name (ARN) of the Delegation Set.
* `caller_reference` - Caller Reference of the delegation set.
* `name_servers` - The list of DNS name servers for the delegation set.
