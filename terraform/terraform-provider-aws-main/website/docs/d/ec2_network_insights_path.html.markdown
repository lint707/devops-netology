---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_ec2_network_insights_path"
description: |-
    Provides details about a specific Network Insights Path.
---

# Data Source: aws_ec2_network_insights_path

`aws_ec2_network_insights_path` provides details about a specific Network Insights Path.

## Example Usage

```terraform
data "aws_ec2_network_insights_path" "example" {
  network_insights_path_id = aws_ec2_network_insights_path.example.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
Network Insights Paths. The given filters must match exactly one Network Insights Path
whose data will be exported as attributes.

* `network_insights_path_id` - (Optional) The ID of the Network Insights Path to select.
* `filter` - (Optional) Configuration block(s) for filtering. Detailed below.

### filter Configuration Block

The following arguments are supported by the `filter` configuration block:

* `name` - (Required) The name of the filter field. Valid values can be found in the EC2 [`DescribeNetworkInsightsPaths`](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeNetworkInsightsPaths.html) API Reference.
* `values` - (Required) Set of values that are accepted for the given filter field. Results will be selected if any given value matches.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the selected Network Insights Path.
* `destination` - The AWS resource that is the destination of the path.
* `destination_ip` - The IP address of the AWS resource that is the destination of the path.
* `destination_port` - The destination port.
* `protocol` - The protocol.
* `source` - The AWS resource that is the source of the path.
* `source_ip` - The IP address of the AWS resource that is the source of the path.
* `tags` - A map of tags assigned to the resource.
