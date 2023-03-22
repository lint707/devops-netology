terraform {
  backend "remote" {
    organization = "hashicorp-v2"

    workspaces {
      name = "terraform-provider-aws-repository"
    }
  }

  required_providers {
    github = {
      source  = "integrations/github"
      version = "4.30.0"
    }
  }

  required_version = ">= 0.13.5"
}

provider "github" {
  owner = "hashicorp"
}
