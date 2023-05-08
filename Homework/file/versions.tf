# Указываем версию Terraform и провайдер Yandex Cloud
terraform {
  required_version = ">= 0.15"
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = ">= 0.49"
    }
  }
}

