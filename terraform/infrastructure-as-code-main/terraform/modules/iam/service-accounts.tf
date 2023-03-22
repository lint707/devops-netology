
variable project_name { default = "" }
variable folder_id { default = "" }

resource "yandex_iam_service_account" "resource" {
  name        = "${var.project_name}-resource"
  description = "service account to manage VMs+k8s"
}

resource "yandex_resourcemanager_folder_iam_binding" "editor" {
  folder_id = "${var.folder_id}"

  role = "editor"

  members = [
    "serviceAccount:${yandex_iam_service_account.resource.id}",
  ]
}

output "id" {
  value       = "${yandex_iam_service_account.resource.id}" 
  description = "id_resource"
}
