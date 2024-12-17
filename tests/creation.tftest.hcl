provider "cloudfoundry" {
  api_url = "https://api.fr.cloud.gov"
}
provider "cloudfoundry-community" {
  api_url = "https://api.fr.cloud.gov"
}

variables {
  cf_space_prefix = "glr-cg-ci-tests"
  ci_server_token = "fake-gdg-server-token"
}

run "test-system-creation" {
  override_resource {
    target = cloudfoundry_app.gitlab-runner-manager
  }

  assert {
    condition     = cloudfoundry_app.gitlab-runner-manager.id == output.manager_app_id && output.manager_app_id != null
    error_message = "Runner manager should have an ID"
  }

  assert {
    condition     = module.egress_proxy.app_id == output.egress_app_id && output.egress_app_id != null
    error_message = "Egress app should have an ID"
  }

  assert {
    condition     = module.manager_space.space_id == output.manager_space_id && output.manager_space_id != null
    error_message = "Manager space was created"
  }

  assert {
    condition     = module.worker_space.space_id == output.worker_space_id && output.worker_space_id != null
    error_message = "Worker space was created"
  }

  assert {
    condition     = module.egress_space.space_id == output.egress_space_id && output.egress_space_id != null
    error_message = "Egress space was created"
  }

  assert {
    condition     = module.object_store_instance.bucket_id == output.object_cache_service_id && output.object_cache_service_id != null
    error_message = "S3 cacheing bucket exists"
  }

  assert {
    condition     = cloudfoundry_app.gitlab-runner-manager.environment.CF_USERNAME == output.service_account_username
    error_message = "Service account username is passed to the manager as CF_USERNAME"
  }

  assert {
    condition     = cloudfoundry_app.gitlab-runner-manager.environment.CF_PASSWORD == local.sa_cf_password
    error_message = "Service account password is passed to the manager as CF_PASSWORD"
  }

  assert {
    condition     = length(cloudfoundry_network_policy.egress_routing.policy) == 2
    error_message = "Egress routing policy should have two entries"
  }

  assert {
    condition = alltrue([
      for p in cloudfoundry_network_policy.egress_routing.policy : p.source_app == output.manager_app_id
    ])
    error_message = "Egress routing only allows routing from the manager app"
  }

  assert {
    condition = alltrue([
      for p in cloudfoundry_network_policy.egress_routing.policy : p.destination_app == output.egress_app_id
    ])
    error_message = "Egress routing only allows routing to the egress app"
  }

  assert {
    condition = alltrue([
      for p in cloudfoundry_network_policy.egress_routing.policy : contains(["8080", "61443"], p.port)
    ])
    error_message = "Egress routing only allows ports 8080 or 61443"
  }

  assert {
    condition     = cloudfoundry_service_instance.runner_service_account.space == output.worker_space_id
    error_message = "The service account user is created on the worker space"
  }

  assert {
    condition = alltrue([
      cloudfoundry_space_role.service-account-egress-role.username == output.service_account_username,
      cloudfoundry_space_role.service-account-egress-role.type == "space_developer",
      cloudfoundry_space_role.service-account-egress-role.space == output.egress_space_id
    ])
    error_message = "Service account is granted space_developer on the egress space"
  }
}
