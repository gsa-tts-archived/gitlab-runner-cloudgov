terraform {
  required_version = "~> 1.5"
  required_providers {
    cloudfoundry = {
      source  = "cloudfoundry/cloudfoundry"
      version = "1.2.0"
    }
    cloudfoundry-community = {
      source  = "cloudfoundry-community/cloudfoundry"
      version = "0.53.1"
    }
  }
}
provider "cloudfoundry" {}

provider "cloudfoundry-community" {
  api_url  = "https://api.fr.cloud.gov"
  user     = var.cf_user
  password = var.cf_password
}

module "sandbox-runner" {
  source = "../"

  cf_org_manager          = var.cf_org_manager
  cf_community_user       = var.cf_user
  cf_space_prefix         = var.cf_space_prefix
  ci_server_url           = "gsa-0.gitlab-dedicated.us"
  ci_server_token         = var.ci_server_token
  docker_hub_user         = var.docker_hub_user
  docker_hub_token        = var.docker_hub_token
  manager_instances       = 1
  developer_emails        = var.developer_emails
  worker_disk_size        = var.worker_disk_size
  program_technologies    = var.program_technologies
  worker_egress_allowlist = var.worker_egress_allowlist
  allow_ssh               = var.allow_ssh
}
