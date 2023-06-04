# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# The following variable is used to configure the provider's authentication
# token. You don't need to provide a token on the command line to apply changes,
# though: using the remote backend, Terraform will execute remotely in Terraform
# Cloud where your token is already securely stored in your workspace!

provider "yandex" {
  token = var.yc_token
  cloud_id  = var.yc_cloud_id
  folder_id = var.yc_folder_id
  zone      = var.yc_region
} 
