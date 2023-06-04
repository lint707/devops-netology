resource "yandex_vpc_network" "internal" {
  name = "internal"
}

resource "yandex_vpc_subnet" "internal-a" {
  name           = "internal-a"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network.internal.id
  v4_cidr_blocks = ["10.0.0.0/24"]
}

resource "yandex_vpc_subnet" "internal-b" {
  name           = "internal-b"
  zone           = "ru-central1-b"
  network_id     = yandex_vpc_network.internal.id
  v4_cidr_blocks = ["10.1.0.0/24"]
}

resource "yandex_vpc_subnet" "internal-c" {
  name           = "internal-c"
  zone           = "ru-central1-c"
  network_id     = yandex_vpc_network.internal.id
  v4_cidr_blocks = ["10.2.0.0/24"]
}
