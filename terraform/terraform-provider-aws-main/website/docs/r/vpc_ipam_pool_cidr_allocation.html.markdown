---
subcategory: "VPC IPAM (IP Address Manager)"
layout: "aws"
page_title: "AWS: aws_vpc_ipam_pool_cidr_allocation"
description: |-
  Allocates (reserves) a CIDR from an IPAM address pool, preventing usage by IPAM.
---

# Resource: aws_vpc_ipam_pool_cidr_allocation

Allocates (reserves) a CIDR from an IPAM address pool, preventing usage by IPAM. Only works for private IPv4.

## Example Usage

Basic usage:

```terraform
data "aws_region" "current" {}

resource "aws_vpc_ipam_pool_cidr_allocation" "example" {
  ipam_pool_id = aws_vpc_ipam_pool.example.id
  cidr         = "172.2.0.0/24"
  depends_on = [
    aws_vpc_ipam_pool_cidr.example
  ]
}

resource "aws_vpc_ipam_pool_cidr" "example" {
  ipam_pool_id = aws_vpc_ipam_pool.example.id
  cidr         = "172.2.0.0/16"
}

resource "aws_vpc_ipam_pool" "example" {
  address_family = "ipv4"
  ipam_scope_id  = aws_vpc_ipam.example.private_default_scope_id
  locale         = data.aws_region.current.name
}

resource "aws_vpc_ipam" "example" {
  operating_regions {
    region_name = data.aws_region.current.name
  }
}
```

With the `disallowed_cidrs` attribute:

```terraform
data "aws_region" "current" {}

resource "aws_vpc_ipam_pool_cidr_allocation" "example" {
  ipam_pool_id  = aws_vpc_ipam_pool.example.id
  netmaskLength = 28

  disallowed_cidrs = [
    "172.2.0.0/28"
  ]

  depends_on = [
    aws_vpc_ipam_pool_cidr.example
  ]
}

resource "aws_vpc_ipam_pool_cidr" "example" {
  ipam_pool_id = aws_vpc_ipam_pool.example.id
  cidr         = "172.2.0.0/16"
}

resource "aws_vpc_ipam_pool" "example" {
  address_family = "ipv4"
  ipam_scope_id  = aws_vpc_ipam.example.private_default_scope_id
  locale         = data.aws_region.current.name
}

resource "aws_vpc_ipam" "example" {
  operating_regions {
    region_name = data.aws_region.current.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Optional) The CIDR you want to assign to the pool.
* `description` - (Optional) The description for the allocation.
* `disallowed_cidrs` - (Optional) Exclude a particular CIDR range from being returned by the pool.
* `ipam_pool_id` - (Required) The ID of the pool to which you want to assign a CIDR.
* `netmask_length` - (Optional) The netmask length of the CIDR you would like to allocate to the IPAM pool. Valid Values: `0-32`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the allocation.
* `resource_id` - The ID of the resource.
* `resource_owner` - The owner of the resource.
* `resource_type` - The type of the resource.

## Import

IPAMs can be imported using the `allocation id`, e.g.

```
$ terraform import aws_vpc_ipam_pool_cidr_allocation.example
```
