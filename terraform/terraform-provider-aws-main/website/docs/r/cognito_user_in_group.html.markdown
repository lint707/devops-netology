---
subcategory: "Cognito IDP (Identity Provider)"
layout: "aws"
page_title: "AWS: aws_cognito_user_in_group"
description: |-
  Adds the specified user to the specified group.
---

# Resource: aws_cognito_user_in_group

Adds the specified user to the specified group.

## Example Usage

```terraform
resource "aws_cognito_user_pool" "example" {
  name = "example"

  password_policy {
    temporary_password_validity_days = 7
    minimum_length                   = 6
    require_uppercase                = false
    require_symbols                  = false
    require_numbers                  = false
  }
}

resource "aws_cognito_user" "example" {
  user_pool_id = aws_cognito_user_pool.test.id
  username     = "example"
}

resource "aws_cognito_user_group" "example" {
  user_pool_id = aws_cognito_user_pool.test.id
  name         = "example"
}

resource "aws_cognito_user_in_group" "example" {
  user_pool_id = aws_cognito_user_pool.example.id
  group_name   = aws_cognito_user_group.example.name
  username     = aws_cognito_user.example.username
}
```

## Argument Reference

The following arguments are required:

* `user_pool_id` - (Required) The user pool ID of the user and group.
* `group_name` - (Required) The name of the group to which the user is to be added.
* `username` - (Required) The username of the user to be added to the group.

## Attributes Reference

No additional attributes are exported.
