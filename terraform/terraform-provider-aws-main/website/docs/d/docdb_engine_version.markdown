---
subcategory: "DocDB (DocumentDB)"
layout: "aws"
page_title: "AWS: aws_docdb_engine_version"
description: |-
  Information about a DocumentDB engine version.
---

# Data Source: aws_docdb_engine_version

Information about a DocumentDB engine version.

## Example Usage

```terraform
data "aws_docdb_engine_version" "test" {
  version = "3.6.0"
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Optional) DB engine. (Default: `docdb`)
* `parameter_group_family` - (Optional) The name of a specific DB parameter group family. An example parameter group family is `docdb3.6`.
* `preferred_versions` - (Optional) Ordered list of preferred engine versions. The first match in this list will be returned. If no preferred matches are found and the original search returned more than one result, an error is returned. If both the `version` and `preferred_versions` arguments are not configured, the data source will return the default version for the engine.
* `version` - (Optional) Version of the DB engine. For example, `3.6.0`. If `version` and `preferred_versions` are not set, the data source will provide information for the AWS-defined default version. If both the `version` and `preferred_versions` arguments are not configured, the data source will return the default version for the engine.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `engine_description` - The description of the database engine.
* `exportable_log_types` - Set of log types that the database engine has available for export to CloudWatch Logs.
* `supports_log_exports_to_cloudwatch` - Indicates whether the engine version supports exporting the log types specified by `exportable_log_types` to CloudWatch Logs.
* `valid_upgrade_targets` - A set of engine versions that this database engine version can be upgraded to.
* `version_description` - The description of the database engine version.
