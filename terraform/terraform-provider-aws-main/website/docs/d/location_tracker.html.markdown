---
subcategory: "Location"
layout: "aws"
page_title: "AWS: aws_location_tracker"
description: |-
    Retrieve information about a Location Service Tracker.
---

# Data Source: aws_location_tracker

Retrieve information about a Location Service Tracker.

## Example Usage

```terraform
data "aws_location_tracker" "example" {
  tracker_name = "example"
}
```

## Argument Reference

* `tracker_name` - (Required) The name of the tracker resource.

## Attribute Reference

* `create_time` - The timestamp for when the tracker resource was created in ISO 8601 format.
* `description` - The optional description for the tracker resource.
* `kms_key_id` - A key identifier for an AWS KMS customer managed key assigned to the Amazon Location resource.
* `position_filtering` - The position filtering method of the tracker resource.
* `tags` - Key-value map of resource tags for the tracker.
* `tracker_arn` - The Amazon Resource Name (ARN) for the tracker resource. Used when you need to specify a resource across all AWS.
* `update_time` - The timestamp for when the tracker resource was last updated in ISO 8601 format.
