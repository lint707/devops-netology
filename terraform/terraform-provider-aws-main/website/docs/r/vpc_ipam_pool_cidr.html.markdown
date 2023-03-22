---
subcategory: "VPC IPAM (IP Address Manager)"
layout: "aws"
page_title: "AWS: aws_vpc_ipam_pool_cidr"
description: |-
  Provisions a CIDR from an IPAM address pool.
---

# Resource: aws_vpc_ipam_pool_cidr

Provisions a CIDR from an IPAM address pool.

~> **NOTE:** Provisioning Public IPv4 or Public IPv6 require [steps outside the scope of this resource](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html#prepare-for-byoip). The resource accepts `message` and `signature` as part of the `cidr_authorization_context` attribute but those must be generated ahead of time. Public IPv6 CIDRs that are provisioned into a Pool with `publicly_advertisable = true` and all public IPv4 CIDRs also require creating a Route Origin Authorization (ROA) object in your Regional Internet Registry (RIR).

~> **NOTE:** In order to deprovision CIDRs all Allocations must be released. Allocations created by a VPC take up to 30 minutes to be released. However, for IPAM to properly manage the removal of allocation records created by VPCs and other resources, you must [grant it permissions](https://docs.aws.amazon.com/vpc/latest/ipam/choose-single-user-or-orgs-ipam.html) in
either a single account or organizationally. If you are unable to deprovision a cidr after waiting over 30 minutes, you may be missing the Service Linked Role.

## Example Usage

Basic usage:

```terraform
data "aws_region" "current" {}

resource "aws_vpc_ipam" "example" {
  operating_regions {
    region_name = data.aws_region.current.name
  }
}

resource "aws_vpc_ipam_pool" "example" {
  address_family = "ipv4"
  ipam_scope_id  = aws_vpc_ipam.example.private_default_scope_id
  locale         = data.aws_region.current.name
}

resource "aws_vpc_ipam_pool_cidr" "example" {
  ipam_pool_id = aws_vpc_ipam_pool.example.id
  cidr         = "172.2.0.0/16"
}
```

Provision Public IPv6 Pool CIDRs:

```terraform
data "aws_region" "current" {}

resource "aws_vpc_ipam" "example" {
  operating_regions {
    region_name = data.aws_region.current.name
  }
}

resource "aws_vpc_ipam_pool" "ipv6_test_public" {
  address_family = "ipv6"
  ipam_scope_id  = aws_vpc_ipam.example.public_default_scope_id
  locale         = "us-east-1"
  description    = "public ipv6"
  advertisable   = false
  aws_service    = "ec2"
}

resource "aws_vpc_ipam_pool_cidr" "ipv6_test_public" {
  ipam_pool_id = aws_vpc_ipam_pool.ipv6_test_public.id
  cidr         = var.ipv6_cidr
  cidr_authorization_context {
    message   = var.message
    signature = var.signature
  }
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Optional) The CIDR you want to assign to the pool.
* `cidr_authorization_context` - (Optional) A signed document that proves that you are authorized to bring the specified IP address range to Amazon using BYOIP. This is not stored in the state file. See [cidr_authorization_context](#cidr_authorization_context) for more information.
* `ipam_pool_id` - (Required) The ID of the pool to which you want to assign a CIDR.

### cidr_authorization_context

* `message` - (Optional) The plain-text authorization message for the prefix and account.
* `signature` - (Optional) The signed authorization message for the prefix and account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IPAM Pool Cidr concatenated with the IPAM Pool ID.

## Import

IPAMs can be imported using the `<cidr>_<ipam-pool-id>`, e.g.

```
$ terraform import aws_vpc_ipam_pool_cidr.example 172.2.0.0/24_ipam-pool-0e634f5a1517cccdc
```
