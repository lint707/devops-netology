---
subcategory: "Data Pipeline"
layout: "aws"
page_title: "AWS: aws_datapipeline_pipeline"
description: |-
  Provides details about a specific DataPipeline.
---

# Source: aws_datapipeline_pipeline

Provides details about a specific DataPipeline Pipeline.

## Example Usage

```terraform
data "aws_datapipeline_pipeline" "example" {
  pipeline_id = "pipelineID"
}
```

## Argument Reference

The following arguments are required:

* `pipeline_id` - (Required) ID of the pipeline.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Name of Pipeline.
* `description` - Description of Pipeline.
* `tags` - A map of tags assigned to the resource.

