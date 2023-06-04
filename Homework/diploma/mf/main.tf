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

resource "yandex_kubernetes_cluster" "kub-diploma" {
  name        = "kub-diploma"
  network_id = yandex_vpc_network.internal.id
  master {
    regional {
      region = "ru-central1"

      location {
        zone      = yandex_vpc_subnet.internal-a.zone
        subnet_id = yandex_vpc_subnet.internal-a.id
      }

      location {
        zone      = yandex_vpc_subnet.internal-b.zone
        subnet_id = yandex_vpc_subnet.internal-b.id
      }

      location {
        zone      = yandex_vpc_subnet.internal-c.zone
        subnet_id = yandex_vpc_subnet.internal-c.id
      }
    }
    version   = "1.22"
    public_ip = true
  }
  release_channel = "RAPID"
  node_service_account_id = yandex_iam_service_account.docker.id
  service_account_id      = yandex_iam_service_account.instances.id
}
resource "yandex_kubernetes_node_group" "diploma-group-auto" {
  cluster_id  = yandex_kubernetes_cluster.kub-diploma.id
  name        = "diploma-group-auto"
  version     = "1.22"

  instance_template {
    platform_id = "standard-v2"
    nat         = true

    resources {
      core_fraction = 20 
      memory        = 2
      cores         = 2
    }

    boot_disk {
      type = "network-hdd"
      size = 64
    }

    scheduling_policy {
      preemptible = false
    }
  }

  scale_policy {
    fixed_scale {
      size = 3
    }
  }

  allocation_policy {
    location {
      zone = "ru-central1-a"
    }

    location {
      zone = "ru-central1-b"
    }

    location {
      zone = "ru-central1-c"
    }
  }

  maintenance_policy {
    auto_upgrade = false
    auto_repair  = true
  }
}
