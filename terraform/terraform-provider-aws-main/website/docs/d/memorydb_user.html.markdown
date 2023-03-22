---
subcategory: "MemoryDB for Redis"
layout: "aws"
page_title: "AWS: aws_memorydb_user"
description: |-
  Provides information about a MemoryDB User.
---

# Resource: aws_memorydb_user

Provides information about a MemoryDB User.

## Example Usage

```terraform
data "aws_memorydb_user" "example" {
  user_name = "my-user"
}
```

## Argument Reference

The following arguments are required:

* `user_name` - (Required) Name of the user.

## Attributes Reference

In addition, the following attributes are exported:

* `id` - Name of the user.
* `access_string` - The access permissions string used for this user.
* `arn` - ARN of the user.
* `authentication_mode` - Denotes the user's authentication properties.
    * `password_count` - The number of passwords belonging to the user.
    * `type` - Indicates whether the user requires a password to authenticate.
* `minimum_engine_version` - The minimum engine version supported for the user.
* `tags` - A map of tags assigned to the subnet group.
