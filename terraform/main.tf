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

resource "yandex_compute_image" "image" {
  name       = "image"
  source_family = "ubuntu-2004-lts"
}

resource "yandex_compute_instance" "ntvm" {
  name = "ntvm"
  zone = "ru-central1-a"
  allow_stopping_for_update = true

  resources {
    cores = 2
    memory = 2
    core_fraction = 20
  }

  boot_disk {
    initialize_params {
      image_id = "${yandex_compute_image.image.id}"
    }
  }

  network_interface {
	    subnet_id = "${yandex_vpc_subnet.internal-a.id}"
	    nat = true
	  }
}
