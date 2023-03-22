# Provider
provider "yandex" {
  token     = var.yc_token
  cloud_id  = var.yc_cloud_id
  zone      = var.yc_region
}