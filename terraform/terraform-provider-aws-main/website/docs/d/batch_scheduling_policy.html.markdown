---
subcategory: "Batch"
layout: "aws"
page_title: "AWS: aws_batch_scheduling_policy"
description: |-
    Provides details about a Batch Scheduling Policy
---

# Data Source: aws_batch_scheduling_policy

The Batch Scheduling Policy data source allows access to details of a specific Scheduling Policy within AWS Batch.

## Example Usage

```terraform
data "aws_batch_scheduling_policy" "test" {
  arn = "arn:aws:batch:us-east-1:012345678910:scheduling-policy/example"
}
```

## Argument Reference

The following arguments are supported:

* `arn` - (Required) The Amazon Resource Name (ARN) of the scheduling policy.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

* `fairshare_policy` - A fairshare policy block specifies the `compute_reservation`, `share_delay_seconds`, and `share_distribution` of the scheduling policy. The `fairshare_policy` block is documented below.
* `name` - Specifies the name of the scheduling policy.
* `tags` - Key-value map of resource tags

A `fairshare_policy` block supports the following arguments:

* `compute_reservation` - A value used to reserve some of the available maximum vCPU for fair share identifiers that have not yet been used. For more information, see [FairsharePolicy](https://docs.aws.amazon.com/batch/latest/APIReference/API_FairsharePolicy.html).
* `share_delay_seconds` - The time period to use to calculate a fair share percentage for each fair share identifier in use, in seconds. For more information, see [FairsharePolicy](https://docs.aws.amazon.com/batch/latest/APIReference/API_FairsharePolicy.html).
* `share_distribution` - One or more share distribution blocks which define the weights for the fair share identifiers for the fair share policy. For more information, see [FairsharePolicy](https://docs.aws.amazon.com/batch/latest/APIReference/API_FairsharePolicy.html). The `share_distribution` block is documented below.

A `share_distribution` block supports the following arguments:

* `share_identifier` - A fair share identifier or fair share identifier prefix. For more information, see [ShareAttributes](https://docs.aws.amazon.com/batch/latest/APIReference/API_ShareAttributes.html).
* `weight_factor` - The weight factor for the fair share identifier. For more information, see [ShareAttributes](https://docs.aws.amazon.com/batch/latest/APIReference/API_ShareAttributes.html).
