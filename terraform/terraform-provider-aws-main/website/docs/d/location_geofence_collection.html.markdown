---
subcategory: "Location"
layout: "aws"
page_title: "AWS: aws_location_geofence_collection"
description: |-
    Retrieve information about a Location Service Geofence Collection.
---

# Data Source: aws_location_geofence_collection

Retrieve information about a Location Service Geofence Collection.

## Example Usage

### Basic Usage

```terraform
data "aws_location_geofence_collection" "example" {
  collection_name = "example"
}
```

## Argument Reference

The following arguments are required:

* `collection_name` - (Required) The name of the geofence collection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `collection_arn` - The Amazon Resource Name (ARN) for the geofence collection resource. Used when you need to specify a resource across all AWS.
* `create_time` - The timestamp for when the geofence collection resource was created in ISO 8601 format.
* `description` - The optional description of the geofence collection resource.
* `kms_key_id` - A key identifier for an AWS KMS customer managed key assigned to the Amazon Location resource.
* `tags` - Key-value map of resource tags for the geofence collection.
* `update_time` - The timestamp for when the geofence collection resource was last updated in ISO 8601 format.
