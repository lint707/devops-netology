---
subcategory: "Connect"
layout: "aws"
page_title: "AWS: aws_connect_user_hierarchy_structure"
description: |-
  Provides details about a specific Amazon Connect User Hierarchy Structure
---

# Resource: aws_connect_user_hierarchy_structure

Provides an Amazon Connect User Hierarchy Structure resource. For more information see
[Amazon Connect: Getting Started](https://docs.aws.amazon.com/connect/latest/adminguide/amazon-connect-get-started.html)

## Example Usage

### Basic

```terraform
resource "aws_connect_user_hierarchy_structure" "example" {
  instance_id = "aaaaaaaa-bbbb-cccc-dddd-111111111111"

  hierarchy_structure {
    level_one {
      name = "levelone"
    }
  }
}
```

### With Five Levels

```terraform
resource "aws_connect_user_hierarchy_structure" "example" {
  instance_id = "aaaaaaaa-bbbb-cccc-dddd-111111111111"

  hierarchy_structure {
    level_one {
      name = "levelone"
    }

    level_two {
      name = "leveltwo"
    }

    level_three {
      name = "levelthree"
    }

    level_four {
      name = "levelfour"
    }

    level_five {
      name = "levelfive"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `hierarchy_structure` - (Required) A block that defines the hierarchy structure's levels. The `hierarchy_structure` block is documented below.
* `instance_id` - (Required) Specifies the identifier of the hosting Amazon Connect Instance.

A `hierarchy_structure` block supports the following arguments:

* `level_one` - (Optional) A block that defines the details of level one. The level block is documented below.
* `level_two` - (Optional) A block that defines the details of level two. The level block is documented below.
* `level_three` - (Optional) A block that defines the details of level three. The level block is documented below.
* `level_four` - (Optional) A block that defines the details of level four. The level block is documented below.
* `level_five` - (Optional) A block that defines the details of level five. The level block is documented below.

Each level block supports the following arguments:

* `name` - (Required) The name of the user hierarchy level. Must not be more than 50 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hierarchy_structure` - In addition to the arguments defined initially, there are attributes added to the levels created. These additional attributes are documented below.
* `id` - The identifier of the hosting Amazon Connect Instance.

A level block supports the following additional attributes:

* `arn` -  The Amazon Resource Name (ARN) of the hierarchy level.
* `id` -  The identifier of the hierarchy level.

## Import

Amazon Connect User Hierarchy Structures can be imported using the `instance_id`, e.g.,

```
$ terraform import aws_connect_user_hierarchy_structure.example f1288a1f-6193-445a-b47e-af739b2
```
