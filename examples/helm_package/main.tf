terraform {
  required_providers {
    artifacthub = {
      version = "0.1.0"
      source  = "arldka.cloud/dev/artifacthub"
    }
  }
}

variable "repo_name" {
  type    = string
  default = "artifact-hub"
}

variable "name" {
  type    = string
  default = "artifact-hub"
}

data "artifacthub_helm_package" "default" {
  repo_name = var.repo_name
  name      = var.name
  version   = "v1.1.1"
}

# Returns the package_id
output "package_id" {
  value = data.artifacthub_helm_package.default.id
}