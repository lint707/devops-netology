## Yandex.Cloud
variable "yc_token" {
  type = string
  description = "Yandex Cloud API key"
}
variable "yc_region" {
  type = string
  description = "Yandex Cloud Region (i.e. ru-central1-a)"
}
variable "yc_cloud_id" {
  type = string
  description = "Yandex Cloud id"
}
variable "yc_folder_id" {
  type = string
  description = "Yandex Cloud folder id"
}

#-----

provider "yandex" {
  token = var.yc_token
  cloud_id  = var.yc_cloud_id
  folder_id = var.yc_folder_id
  zone      = var.yc_region
}
# Создаем виртуалку в публичной подсети с публичным IP и образом ubuntu-2004-lts
resource "yandex_compute_instance" "publicvm" {
  name          = "publicvm"
  zone          = var.zone
  platform_id   = "standard-v2"
  resources {
    cores  = 2
    memory = 2
  }
  boot_disk {
    initialize_params {
      image_id = "fd8gfg42q4551cvt340b"
    }
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.public.id
    nat         = true
  }
}
# Создаем NAT-инстанс в публичной подсети с адресом 192.168.10.254 и образом fd80mrhj8fl2oe87o4e1
resource "yandex_compute_instance" "nat" {
  name          = "nat"
  zone          = var.zone
  platform_id   = "standard-v2"
  resources {
    cores  = 2
    memory = 2
  }
  boot_disk {
    initialize_params {
      image_id = "fd80mrhj8fl2oe87o4e1"
    }
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.public.id
    ip_address = "192.168.10.254"
    nat         = true
  }
}

# Создаем виртуалку в приватной подсети с внутренним IP и образом ubuntu-2004-lts
resource "yandex_compute_instance" "privatevm" {
  name          = "privatevm"
  zone          = var.zone
  platform_id   = "standard-v2"
  resources {
    cores  = 2
    memory = 2
  }
  boot_disk {
    initialize_params {
      image_id = "fd8gfg42q4551cvt340b"
    }
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.private.id
    ip_address = "192.168.20.10"
  }
}

