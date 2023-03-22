---
subcategory: "Cognito IDP (Identity Provider)"
layout: "aws"
page_title: "AWS: aws_cognito_resource_server"
side_bar_current: "docs-aws-resource-cognito-resource-server"
description: |-
  Provides a Cognito Resource Server.
---

# Resource: aws_cognito_resource_server

Provides a Cognito Resource Server.

## Example Usage

### Create a basic resource server

```terraform
resource "aws_cognito_user_pool" "pool" {
  name = "pool"
}

resource "aws_cognito_resource_server" "resource" {
  identifier = "https://example.com"
  name       = "example"

  user_pool_id = aws_cognito_user_pool.pool.id
}
```

### Create a resource server with sample-scope

```terraform
resource "aws_cognito_user_pool" "pool" {
  name = "pool"
}

resource "aws_cognito_resource_server" "resource" {
  identifier = "https://example.com"
  name       = "example"

  scope {
    scope_name        = "sample-scope"
    scope_description = "a Sample Scope Description"
  }

  user_pool_id = aws_cognito_user_pool.pool.id
}
```

## Argument Reference

The following arguments are supported:

* `identifier` - (Required) An identifier for the resource server.
* `name` - (Required) A name for the resource server.
* `scope` - (Optional) A list of [Authorization Scope](#authorization-scope).

### Authorization Scope

* `scope_name` - (Required) The scope name.
* `scope_description` - (Required) The scope description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scope_identifiers` - A list of all scopes configured for this resource server in the format identifier/scope_name.

## Import

`aws_cognito_resource_server` can be imported using their User Pool ID and Identifier, e.g.,

```
$ terraform import aws_cognito_resource_server.example us-west-2_abc123:https://example.com
```
