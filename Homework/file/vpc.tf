# Создаем пустую VPC с названием vpc
resource "yandex_vpc_network" "vpc" {
  name = "vpc"
}

# Создаем публичную подсеть с названием public и сетью 192.168.10.0/24
resource "yandex_vpc_subnet" "public" {
  name           = "public"
  zone           = var.zone
  network_id     = yandex_vpc_network.vpc.id
  v4_cidr_blocks = ["192.168.10.0/24"]
}
# Создаем таблицу маршрутизации
resource "yandex_vpc_route_table" "private-rt" {
  name      = "private-rt"
  network_id = yandex_vpc_network.vpc.id
# Добавляем статический маршрут для 0.0.0.0/0 через NAT-инстанс
  static_route {
    destination_prefix = "0.0.0.0/0"
    next_hop_address   = "192.168.10.254"
  }
}
# Создаем приватную подсеть с названием private, сетью 192.168.20.0/24 и без автоматического выделения публичных IP
resource "yandex_vpc_subnet" "private" {
  name           = "private"
  zone           = var.zone
  network_id     = yandex_vpc_network.vpc.id
  v4_cidr_blocks = ["192.168.20.0/24"]
  route_table_id = yandex_vpc_route_table.private-rt.id # привязка таблицы маршрутизации
}


