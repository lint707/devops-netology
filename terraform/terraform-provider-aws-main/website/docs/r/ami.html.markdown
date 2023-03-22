---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "AWS: aws_ami"
description: |-
  Creates and manages a custom Amazon Machine Image (AMI).
---

# Resource: aws_ami

The AMI resource allows the creation and management of a completely-custom
*Amazon Machine Image* (AMI).

If you just want to duplicate an existing AMI, possibly copying it to another
region, it's better to use `aws_ami_copy` instead.

If you just want to share an existing AMI with another AWS account,
it's better to use `aws_ami_launch_permission` instead.

## Example Usage

```terraform
# Create an AMI that will start a machine whose root device is backed by
# an EBS volume populated from a snapshot. We assume that such a snapshot
# already exists with the id "snap-xxxxxxxx".
resource "aws_ami" "example" {
  name                = "terraform-example"
  virtualization_type = "hvm"
  root_device_name    = "/dev/xvda"

  ebs_block_device {
    device_name = "/dev/xvda"
    snapshot_id = "snap-xxxxxxxx"
    volume_size = 8
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A region-unique name for the AMI.
* `boot_mode` - (Optional) The boot mode of the AMI. For more information, see [Boot modes](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ami-boot.html) in the Amazon Elastic Compute Cloud User Guide.
* `deprecation_time` - (Optional) The date and time to deprecate the AMI. If you specified a value for seconds, Amazon EC2 rounds the seconds to the nearest minute. Valid values: [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) (`YYYY-MM-DDTHH:MM:SSZ`)
* `description` - (Optional) A longer, human-readable description for the AMI.
* `ena_support` - (Optional) Specifies whether enhanced networking with ENA is enabled. Defaults to `false`.
* `root_device_name` - (Optional) The name of the root device (for example, `/dev/sda1`, or `/dev/xvda`).
* `virtualization_type` - (Optional) Keyword to choose what virtualization mode created instances
  will use. Can be either "paravirtual" (the default) or "hvm". The choice of virtualization type
  changes the set of further arguments that are required, as described below.
* `architecture` - (Optional) Machine architecture for created instances. Defaults to "x86_64".
* `ebs_block_device` - (Optional) Nested block describing an EBS block device that should be
  attached to created instances. The structure of this block is described below.
* `ephemeral_block_device` - (Optional) Nested block describing an ephemeral block device that
  should be attached to created instances. The structure of this block is described below.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `tpm_support` - (Optional) If the image is configured for NitroTPM support, the value is `v2.0`. For more information, see [NitroTPM](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/nitrotpm.html) in the Amazon Elastic Compute Cloud User Guide.

When `virtualization_type` is "paravirtual" the following additional arguments apply:

* `image_location` - (Required) Path to an S3 object containing an image manifest, e.g., created
  by the `ec2-upload-bundle` command in the EC2 command line tools.
* `kernel_id` - (Required) The id of the kernel image (AKI) that will be used as the paravirtual
  kernel in created instances.
* `ramdisk_id` - (Optional) The id of an initrd image (ARI) that will be used when booting the
  created instances.

When `virtualization_type` is "hvm" the following additional arguments apply:

* `sriov_net_support` - (Optional) When set to "simple" (the default), enables enhanced networking
  for created instances. No other value is supported at this time.

Nested `ebs_block_device` blocks have the following structure:

* `device_name` - (Required) The path at which the device is exposed to created instances.
* `delete_on_termination` - (Optional) Boolean controlling whether the EBS volumes created to
  support each created instance will be deleted once that instance is terminated.
* `encrypted` - (Optional) Boolean controlling whether the created EBS volumes will be encrypted. Can't be used with `snapshot_id`.
* `iops` - (Required only when `volume_type` is `io1` or `io2`) Number of I/O operations per second the
  created volumes will support.
* `snapshot_id` - (Optional) The id of an EBS snapshot that will be used to initialize the created
  EBS volumes. If set, the `volume_size` attribute must be at least as large as the referenced
  snapshot.
* `throughput` - (Optional) The throughput that the EBS volume supports, in MiB/s. Only valid for `volume_type` of `gp3`.
* `volume_size` - (Required unless `snapshot_id` is set) The size of created volumes in GiB.
  If `snapshot_id` is set and `volume_size` is omitted then the volume will have the same size
  as the selected snapshot.
* `volume_type` - (Optional) The type of EBS volume to create. Can be `standard`, `gp2`, `gp3`, `io1`, `io2`, `sc1` or `st1` (Default: `standard`).
* `kms_key_id` - (Optional) The full ARN of the AWS Key Management Service (AWS KMS) CMK to use when encrypting the snapshots of
an image during a copy operation. This parameter is only required if you want to use a non-default CMK;
if this parameter is not specified, the default CMK for EBS is used
* `outpost_arn` - (Optional) The ARN of the Outpost on which the snapshot is stored.

~> **Note:** You can specify `encrypted` or `snapshot_id` but not both.

Nested `ephemeral_block_device` blocks have the following structure:

* `device_name` - (Required) The path at which the device is exposed to created instances.
* `virtual_name` - (Required) A name for the ephemeral device, of the form "ephemeralN" where
  *N* is a volume number starting from zero.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the AMI.
* `id` - The ID of the created AMI.
* `owner_id` - The AWS account ID of the image owner.
* `root_snapshot_id` - The Snapshot ID for the root volume (for EBS-backed AMIs)
* `usage_operation` - The operation of the Amazon EC2 instance and the billing code that is associated with the AMI.
* `platform_details` - The platform details associated with the billing code of the AMI.
* `image_owner_alias` - The AWS account alias (for example, amazon, self) or the AWS account ID of the AMI owner.
* `image_type` - The type of image.
* `hypervisor` - The hypervisor type of the image.
* `owner_id` - The AWS account ID of the image owner.
* `platform` - This value is set to windows for Windows AMIs; otherwise, it is blank.
* `public` - Indicates whether the image has public launch permissions.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Timeouts

[Configuration options](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts):

* `create` - (Default `40m`)
* `update` - (Default `40m`)
* `delete` - (Default `90m`)

## Import

`aws_ami` can be imported using the ID of the AMI, e.g.,

```
$ terraform import aws_ami.example ami-12345678
```
