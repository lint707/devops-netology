---
subcategory: "AppSync"
layout: "aws"
page_title: "AWS: aws_appsync_api_cache"
description: |-
  Provides an AppSync API Cache.
---

# Resource: aws_appsync_api_cache

Provides an AppSync API Cache.

## Example Usage

```terraform
resource "aws_appsync_graphql_api" "example" {
  authentication_type = "API_KEY"
  name                = "example"
}

resource "aws_appsync_api_cache" "example" {
  api_id               = aws_appsync_graphql_api.example.id
  api_caching_behavior = "FULL_REQUEST_CACHING"
  type                 = "LARGE"
  ttl                  = 900
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) The GraphQL API ID.
* `api_caching_behavior` - (Required) Caching behavior. Valid values are `FULL_REQUEST_CACHING` and `PER_RESOLVER_CACHING`.
* `type` - (Required) The cache instance type. Valid values are `SMALL`, `MEDIUM`, `LARGE`, `XLARGE`, `LARGE_2X`, `LARGE_4X`, `LARGE_8X`, `LARGE_12X`, `T2_SMALL`, `T2_MEDIUM`, `R4_LARGE`, `R4_XLARGE`, `R4_2XLARGE`, `R4_4XLARGE`, `R4_8XLARGE`.
* `ttl` - (Required) TTL in seconds for cache entries.
* `at_rest_encryption_enabled` - (Optional) At-rest encryption flag for cache. You cannot update this setting after creation.
* `transit_encryption_enabled` - (Optional) Transit encryption flag when connecting to cache. You cannot update this setting after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The AppSync API ID.

## Import

`aws_appsync_api_cache` can be imported using the AppSync API ID, e.g.,

```
$ terraform import aws_appsync_api_cache.example xxxxx
```
