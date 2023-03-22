---
subcategory: "AppSync"
layout: "aws"
page_title: "AWS: aws_appsync_resolver"
description: |-
  Provides an AppSync Resolver.
---

# Resource: aws_appsync_resolver

Provides an AppSync Resolver.

## Example Usage

```terraform
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = "tf-example"

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = "tf_example"
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

# UNIT type resolver (default)
resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF

  caching_config {
    caching_keys = [
      "$context.identity.sub",
      "$context.arguments.id",
    ]
    ttl = 60
  }
}

# PIPELINE type resolver
resource "aws_appsync_resolver" "Mutation_pipelineTest" {
  type              = "Mutation"
  api_id            = aws_appsync_graphql_api.test.id
  field             = "pipelineTest"
  request_template  = "{}"
  response_template = "$util.toJson($ctx.result)"
  kind              = "PIPELINE"
  pipeline_config {
    functions = [
      aws_appsync_function.test1.function_id,
      aws_appsync_function.test2.function_id,
      aws_appsync_function.test3.function_id,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) The API ID for the GraphQL API.
* `type` - (Required) The type name from the schema defined in the GraphQL API.
* `field` - (Required) The field name from the schema defined in the GraphQL API.
* `request_template` - (Optional) The request mapping template for UNIT resolver or 'before mapping template' for PIPELINE resolver. Required for non-Lambda resolvers.
* `response_template` - (Optional) The response mapping template for UNIT resolver or 'after mapping template' for PIPELINE resolver. Required for non-Lambda resolvers.
* `data_source` - (Optional) The data source name.
* `max_batch_size` - (Optional) The maximum batching size for a resolver. Valid values are between `0` and `2000`.
* `kind`  - (Optional) The resolver type. Valid values are `UNIT` and `PIPELINE`.
* `sync_config` - (Optional) Describes a Sync configuration for a resolver. See [Sync Config](#sync-config).
* `pipeline_config` - (Optional) The PipelineConfig.
    * `functions` - (Required) The list of Function ID.
* `caching_config` - (Optional) The CachingConfig.
    * `caching_keys` - (Optional) The list of caching key.
    * `ttl` - (Optional) The TTL in seconds.

### Sync Config

The following arguments are supported:

* `conflict_detection` - (Optional) The Conflict Detection strategy to use. Valid values are `NONE` and `VERSION`.
* `conflict_handler` - (Optional) The Conflict Resolution strategy to perform in the event of a conflict. Valid values are `NONE`, `OPTIMISTIC_CONCURRENCY`, `AUTOMERGE`, and `LAMBDA`.
* `lambda_conflict_handler_config` - (Optional) The Lambda Conflict Handler Config when configuring `LAMBDA` as the Conflict Handler. See [Lambda Conflict Handler Config](#lambda-conflict-handler-config).

#### Lambda Conflict Handler Config

The following arguments are supported:

* `lambda_conflict_handler_arn` - (Optional) The Amazon Resource Name (ARN) for the Lambda function to use as the Conflict Handler.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN

## Import

`aws_appsync_resolver` can be imported with their `api_id`, a hyphen, `type`, a hypen and `field` e.g.,

```
$ terraform import aws_appsync_resolver.example abcdef123456-exampleType-exampleField
```
