---
subcategory: "CloudWatch Logs"
layout: "aws"
page_title: "AWS: aws_cloudwatch_log_stream"
description: |-
  Provides a CloudWatch Log Stream resource.
---

# Resource: aws_cloudwatch_log_stream

Provides a CloudWatch Log Stream resource.

## Example Usage

```terraform
resource "aws_cloudwatch_log_group" "yada" {
  name = "Yada"
}

resource "aws_cloudwatch_log_stream" "foo" {
  name           = "SampleLogStream1234"
  log_group_name = aws_cloudwatch_log_group.yada.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the log stream. Must not be longer than 512 characters and must not contain `:`
* `log_group_name` - (Required) The name of the log group under which the log stream is to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) specifying the log stream.

## Import

Cloudwatch Log Stream can be imported using the stream's `log_group_name` and `name`, e.g.,

```
$ terraform import aws_cloudwatch_log_stream.foo Yada:SampleLogStream1234
```
