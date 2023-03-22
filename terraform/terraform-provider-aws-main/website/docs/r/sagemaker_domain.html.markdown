---
subcategory: "SageMaker"
layout: "aws"
page_title: "AWS: aws_sagemaker_domain"
description: |-
  Provides a SageMaker Domain resource.
---

# Resource: aws_sagemaker_domain

Provides a SageMaker Domain resource.

## Example Usage

### Basic usage

```terraform
resource "aws_sagemaker_domain" "example" {
  domain_name = "example"
  auth_mode   = "IAM"
  vpc_id      = aws_vpc.test.id
  subnet_ids  = [aws_subnet.test.id]

  default_user_settings {
    execution_role = aws_iam_role.test.arn
  }
}

resource "aws_iam_role" "example" {
  name               = "example"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.example.json
}

data "aws_iam_policy_document" "example" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["sagemaker.amazonaws.com"]
    }
  }
}
```

### Using Custom Images

```terraform
resource "aws_sagemaker_image" "test" {
  image_name = "example"
  role_arn   = aws_iam_role.test.arn
}

resource "aws_sagemaker_app_image_config" "test" {
  app_image_config_name = "example"

  kernel_gateway_image_config {
    kernel_spec {
      name = "example"
    }
  }
}

resource "aws_sagemaker_image_version" "test" {
  image_name = aws_sagemaker_image.test.id
  base_image = "base-image"
}

resource "aws_sagemaker_domain" "test" {
  domain_name = "example"
  auth_mode   = "IAM"
  vpc_id      = aws_vpc.test.id
  subnet_ids  = [aws_subnet.test.id]

  default_user_settings {
    execution_role = aws_iam_role.test.arn

    kernel_gateway_app_settings {
      custom_image {
        app_image_config_name = aws_sagemaker_app_image_config.test.app_image_config_name
        image_name            = aws_sagemaker_image_version.test.image_name
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) The domain name.
* `auth_mode` - (Required) The mode of authentication that members use to access the domain. Valid values are `IAM` and `SSO`.
* `vpc_id` - (Required) The ID of the Amazon Virtual Private Cloud (VPC) that Studio uses for communication.
* `subnet_ids` - (Required) The VPC subnets that Studio uses for communication.
* `default_user_settings` - (Required) The default user settings. See [Default User Settings](#default-user-settings) below.
* `retention_policy` - (Optional) The retention policy for this domain, which specifies whether resources will be retained after the Domain is deleted. By default, all resources are retained. See [Retention Policy](#retention-policy) below.
* `kms_key_id` - (Optional) The AWS KMS customer managed CMK used to encrypt the EFS volume attached to the domain.
* `app_network_access_type` - (Optional) Specifies the VPC used for non-EFS traffic. The default value is `PublicInternetOnly`. Valid values are `PublicInternetOnly` and `VpcOnly`.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

### Default User Settings

* `execution_role` - (Required) The execution role ARN for the user.
* `security_groups` - (Optional) The security groups.
* `sharing_settings` - (Optional) The sharing settings. See [Sharing Settings](#sharing-settings) below.
* `tensor_board_app_settings` - (Optional) The TensorBoard app settings. See [TensorBoard App Settings](#tensorboard-app-settings) below.
* `jupyter_server_app_settings` - (Optional) The Jupyter server's app settings. See [Jupyter Server App Settings](#jupyter-server-app-settings) below.
* `kernel_gateway_app_settings` - (Optional) The kernel gateway app settings. See [Kernel Gateway App Settings](#kernel-gateway-app-settings) below.

#### Sharing Settings

* `notebook_output_option` - (Optional) Whether to include the notebook cell output when sharing the notebook. The default is `Disabled`. Valid values are `Allowed` and `Disabled`.
* `s3_kms_key_id` - (Optional) When `notebook_output_option` is Allowed, the AWS Key Management Service (KMS) encryption key ID used to encrypt the notebook cell output in the Amazon S3 bucket.
* `s3_output_path` - (Optional) When `notebook_output_option` is Allowed, the Amazon S3 bucket used to save the notebook cell output.

#### TensorBoard App Settings

* `default_resource_spec` - (Optional) The default instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance. see [Default Resource Spec](#default-resource-spec) below.

#### Kernel Gateway App Settings

* `default_resource_spec` - (Optional) The default instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance. see [Default Resource Spec](#default-resource-spec) below.
* `custom_image` - (Optional) A list of custom SageMaker images that are configured to run as a KernelGateway app. see [Custom Image](#custom-image) below.
* `lifecycle_config_arns` - (Optional) The Amazon Resource Name (ARN) of the Lifecycle Configurations.

#### Jupyter Server App Settings

* `default_resource_spec` - (Optional) The default instance type and the Amazon Resource Name (ARN) of the SageMaker image created on the instance. see [Default Resource Spec](#default-resource-spec) below.
* `lifecycle_config_arns` - (Optional) The Amazon Resource Name (ARN) of the Lifecycle Configurations.

##### Default Resource Spec

* `instance_type` - (Optional) The instance type that the image version runs on.. For valid values see [SageMaker Instance Types](https://docs.aws.amazon.com/sagemaker/latest/dg/notebooks-available-instance-types.html).
* `lifecycle_config_arn` - (Optional) The Amazon Resource Name (ARN) of the Lifecycle Configuration attached to the Resource.
* `sagemaker_image_arn` - (Optional) The ARN of the SageMaker image that the image version belongs to.
* `sagemaker_image_version_arn` - (Optional) The ARN of the image version created on the instance.

##### Custom Image

* `app_image_config_name` - (Required) The name of the App Image Config.
* `image_name` - (Required) The name of the Custom Image.
* `image_version_number` - (Optional) The version number of the Custom Image.

### Retention Policy

* `home_efs_file_system` - (Optional) The retention policy for data stored on an Amazon Elastic File System (EFS) volume. Valid values are `Retain` or `Delete`.  Default value is `Retain`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Domain.
* `arn` - The Amazon Resource Name (ARN) assigned by AWS to this Domain.
* `url` - The domain's URL.
* `single_sign_on_managed_application_instance_id` - The SSO managed application instance ID.
* `home_efs_file_system_id` - The ID of the Amazon Elastic File System (EFS) managed by this Domain.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

SageMaker Code Domains can be imported using the `id`, e.g.,

```
$ terraform import aws_sagemaker_domain.test_domain d-8jgsjtilstu8
```
