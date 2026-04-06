resource "pingone_davinci_connector_instance" "connector-oai-github" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-github"
  }
  name = "My awesome connector-oai-github"
  property {
    name  = "actionsAddRepoAccessToSelfHostedRunnerGroupInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_repo_access_to_self_hosted_runner_group_in_org_org
  }
  property {
    name  = "actionsAddRepoAccessToSelfHostedRunnerGroupInOrg_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_repo_access_to_self_hosted_runner_group_in_org_repository_id
  }
  property {
    name  = "actionsAddRepoAccessToSelfHostedRunnerGroupInOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_repo_access_to_self_hosted_runner_group_in_org_runner_group_id
  }
  property {
    name  = "actionsAddSelectedRepoToOrgSecret_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_selected_repo_to_org_secret_org
  }
  property {
    name  = "actionsAddSelectedRepoToOrgSecret_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_selected_repo_to_org_secret_repository_id
  }
  property {
    name  = "actionsAddSelectedRepoToOrgSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_selected_repo_to_org_secret_secret_name
  }
  property {
    name  = "actionsAddSelfHostedRunnerToGroupForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_self_hosted_runner_to_group_for_org_org
  }
  property {
    name  = "actionsAddSelfHostedRunnerToGroupForOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_self_hosted_runner_to_group_for_org_runner_group_id
  }
  property {
    name  = "actionsAddSelfHostedRunnerToGroupForOrg_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_add_self_hosted_runner_to_group_for_org_runner_id
  }
  property {
    name  = "actionsCancelWorkflowRun_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_cancel_workflow_run_owner
  }
  property {
    name  = "actionsCancelWorkflowRun_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_cancel_workflow_run_repo
  }
  property {
    name  = "actionsCancelWorkflowRun_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_cancel_workflow_run_run_id
  }
  property {
    name  = "actionsCreateOrUpdateOrgSecret_actionsCreateOrUpdateOrgSecretRequest_ActionsCreateOrUpdateOrgSecretRequest_encrypted_value"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_org_secret_actions_create_or_update_org_secret_request_actions_create_or_update_org_secret_request_encrypted_value
  }
  property {
    name  = "actionsCreateOrUpdateOrgSecret_actionsCreateOrUpdateOrgSecretRequest_ActionsCreateOrUpdateOrgSecretRequest_key_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_org_secret_actions_create_or_update_org_secret_request_actions_create_or_update_org_secret_request_key_id
  }
  property {
    name  = "actionsCreateOrUpdateOrgSecret_actionsCreateOrUpdateOrgSecretRequest_ActionsCreateOrUpdateOrgSecretRequest_selected_repository_ids"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_org_secret_actions_create_or_update_org_secret_request_actions_create_or_update_org_secret_request_selected_repository_ids
  }
  property {
    name  = "actionsCreateOrUpdateOrgSecret_actionsCreateOrUpdateOrgSecretRequest_ActionsCreateOrUpdateOrgSecretRequest_visibility"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_org_secret_actions_create_or_update_org_secret_request_actions_create_or_update_org_secret_request_visibility
  }
  property {
    name  = "actionsCreateOrUpdateOrgSecret_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_org_secret_org
  }
  property {
    name  = "actionsCreateOrUpdateOrgSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_org_secret_secret_name
  }
  property {
    name  = "actionsCreateOrUpdateRepoSecret_actionsCreateOrUpdateRepoSecretRequest_ActionsCreateOrUpdateRepoSecretRequest_encrypted_value"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_repo_secret_actions_create_or_update_repo_secret_request_actions_create_or_update_repo_secret_request_encrypted_value
  }
  property {
    name  = "actionsCreateOrUpdateRepoSecret_actionsCreateOrUpdateRepoSecretRequest_ActionsCreateOrUpdateRepoSecretRequest_key_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_repo_secret_actions_create_or_update_repo_secret_request_actions_create_or_update_repo_secret_request_key_id
  }
  property {
    name  = "actionsCreateOrUpdateRepoSecret_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_repo_secret_owner
  }
  property {
    name  = "actionsCreateOrUpdateRepoSecret_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_repo_secret_repo
  }
  property {
    name  = "actionsCreateOrUpdateRepoSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_or_update_repo_secret_secret_name
  }
  property {
    name  = "actionsCreateRegistrationTokenForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_registration_token_for_org_org
  }
  property {
    name  = "actionsCreateRegistrationTokenForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_registration_token_for_repo_owner
  }
  property {
    name  = "actionsCreateRegistrationTokenForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_registration_token_for_repo_repo
  }
  property {
    name  = "actionsCreateRemoveTokenForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_remove_token_for_org_org
  }
  property {
    name  = "actionsCreateRemoveTokenForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_remove_token_for_repo_owner
  }
  property {
    name  = "actionsCreateRemoveTokenForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_remove_token_for_repo_repo
  }
  property {
    name  = "actionsCreateSelfHostedRunnerGroupForOrg_actionsCreateSelfHostedRunnerGroupForOrgRequest_ActionsCreateSelfHostedRunnerGroupForOrgRequest_allows_public_repositories"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_self_hosted_runner_group_for_org_actions_create_self_hosted_runner_group_for_org_request_actions_create_self_hosted_runner_group_for_org_request_allows_public_repositories
  }
  property {
    name  = "actionsCreateSelfHostedRunnerGroupForOrg_actionsCreateSelfHostedRunnerGroupForOrgRequest_ActionsCreateSelfHostedRunnerGroupForOrgRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_self_hosted_runner_group_for_org_actions_create_self_hosted_runner_group_for_org_request_actions_create_self_hosted_runner_group_for_org_request_name
  }
  property {
    name  = "actionsCreateSelfHostedRunnerGroupForOrg_actionsCreateSelfHostedRunnerGroupForOrgRequest_ActionsCreateSelfHostedRunnerGroupForOrgRequest_runners"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_self_hosted_runner_group_for_org_actions_create_self_hosted_runner_group_for_org_request_actions_create_self_hosted_runner_group_for_org_request_runners
  }
  property {
    name  = "actionsCreateSelfHostedRunnerGroupForOrg_actionsCreateSelfHostedRunnerGroupForOrgRequest_ActionsCreateSelfHostedRunnerGroupForOrgRequest_selected_repository_ids"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_self_hosted_runner_group_for_org_actions_create_self_hosted_runner_group_for_org_request_actions_create_self_hosted_runner_group_for_org_request_selected_repository_ids
  }
  property {
    name  = "actionsCreateSelfHostedRunnerGroupForOrg_actionsCreateSelfHostedRunnerGroupForOrgRequest_ActionsCreateSelfHostedRunnerGroupForOrgRequest_visibility"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_self_hosted_runner_group_for_org_actions_create_self_hosted_runner_group_for_org_request_actions_create_self_hosted_runner_group_for_org_request_visibility
  }
  property {
    name  = "actionsCreateSelfHostedRunnerGroupForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_self_hosted_runner_group_for_org_org
  }
  property {
    name  = "actionsCreateWorkflowDispatch_actionsCreateWorkflowDispatchRequest_ActionsCreateWorkflowDispatchRequest_inputs"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_workflow_dispatch_actions_create_workflow_dispatch_request_actions_create_workflow_dispatch_request_inputs
  }
  property {
    name  = "actionsCreateWorkflowDispatch_actionsCreateWorkflowDispatchRequest_ActionsCreateWorkflowDispatchRequest_ref"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_workflow_dispatch_actions_create_workflow_dispatch_request_actions_create_workflow_dispatch_request_ref
  }
  property {
    name  = "actionsCreateWorkflowDispatch_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_workflow_dispatch_owner
  }
  property {
    name  = "actionsCreateWorkflowDispatch_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_create_workflow_dispatch_repo
  }
  property {
    name  = "actionsDeleteArtifact_artifact_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_artifact_artifact_id
  }
  property {
    name  = "actionsDeleteArtifact_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_artifact_owner
  }
  property {
    name  = "actionsDeleteArtifact_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_artifact_repo
  }
  property {
    name  = "actionsDeleteOrgSecret_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_org_secret_org
  }
  property {
    name  = "actionsDeleteOrgSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_org_secret_secret_name
  }
  property {
    name  = "actionsDeleteRepoSecret_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_repo_secret_owner
  }
  property {
    name  = "actionsDeleteRepoSecret_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_repo_secret_repo
  }
  property {
    name  = "actionsDeleteRepoSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_repo_secret_secret_name
  }
  property {
    name  = "actionsDeleteSelfHostedRunnerFromOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_self_hosted_runner_from_org_org
  }
  property {
    name  = "actionsDeleteSelfHostedRunnerFromOrg_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_self_hosted_runner_from_org_runner_id
  }
  property {
    name  = "actionsDeleteSelfHostedRunnerFromRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_self_hosted_runner_from_repo_owner
  }
  property {
    name  = "actionsDeleteSelfHostedRunnerFromRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_self_hosted_runner_from_repo_repo
  }
  property {
    name  = "actionsDeleteSelfHostedRunnerFromRepo_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_self_hosted_runner_from_repo_runner_id
  }
  property {
    name  = "actionsDeleteSelfHostedRunnerGroupFromOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_self_hosted_runner_group_from_org_org
  }
  property {
    name  = "actionsDeleteSelfHostedRunnerGroupFromOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_self_hosted_runner_group_from_org_runner_group_id
  }
  property {
    name  = "actionsDeleteWorkflowRunLogs_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_workflow_run_logs_owner
  }
  property {
    name  = "actionsDeleteWorkflowRunLogs_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_workflow_run_logs_repo
  }
  property {
    name  = "actionsDeleteWorkflowRunLogs_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_workflow_run_logs_run_id
  }
  property {
    name  = "actionsDeleteWorkflowRun_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_workflow_run_owner
  }
  property {
    name  = "actionsDeleteWorkflowRun_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_workflow_run_repo
  }
  property {
    name  = "actionsDeleteWorkflowRun_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_delete_workflow_run_run_id
  }
  property {
    name  = "actionsDisableSelectedRepositoryGithubActionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_disable_selected_repository_github_actions_organization_org
  }
  property {
    name  = "actionsDisableSelectedRepositoryGithubActionsOrganization_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_disable_selected_repository_github_actions_organization_repository_id
  }
  property {
    name  = "actionsDisableWorkflow_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_disable_workflow_owner
  }
  property {
    name  = "actionsDisableWorkflow_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_disable_workflow_repo
  }
  property {
    name  = "actionsDownloadArtifact_archive_format"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_artifact_archive_format
  }
  property {
    name  = "actionsDownloadArtifact_artifact_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_artifact_artifact_id
  }
  property {
    name  = "actionsDownloadArtifact_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_artifact_owner
  }
  property {
    name  = "actionsDownloadArtifact_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_artifact_repo
  }
  property {
    name  = "actionsDownloadJobLogsForWorkflowRun_job_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_job_logs_for_workflow_run_job_id
  }
  property {
    name  = "actionsDownloadJobLogsForWorkflowRun_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_job_logs_for_workflow_run_owner
  }
  property {
    name  = "actionsDownloadJobLogsForWorkflowRun_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_job_logs_for_workflow_run_repo
  }
  property {
    name  = "actionsDownloadWorkflowRunLogs_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_workflow_run_logs_owner
  }
  property {
    name  = "actionsDownloadWorkflowRunLogs_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_workflow_run_logs_repo
  }
  property {
    name  = "actionsDownloadWorkflowRunLogs_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_download_workflow_run_logs_run_id
  }
  property {
    name  = "actionsEnableSelectedRepositoryGithubActionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_enable_selected_repository_github_actions_organization_org
  }
  property {
    name  = "actionsEnableSelectedRepositoryGithubActionsOrganization_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_enable_selected_repository_github_actions_organization_repository_id
  }
  property {
    name  = "actionsEnableWorkflow_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_enable_workflow_owner
  }
  property {
    name  = "actionsEnableWorkflow_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_enable_workflow_repo
  }
  property {
    name  = "actionsGetAllowedActionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_allowed_actions_organization_org
  }
  property {
    name  = "actionsGetAllowedActionsRepository_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_allowed_actions_repository_owner
  }
  property {
    name  = "actionsGetAllowedActionsRepository_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_allowed_actions_repository_repo
  }
  property {
    name  = "actionsGetArtifact_artifact_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_artifact_artifact_id
  }
  property {
    name  = "actionsGetArtifact_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_artifact_owner
  }
  property {
    name  = "actionsGetArtifact_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_artifact_repo
  }
  property {
    name  = "actionsGetGithubActionsPermissionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_github_actions_permissions_organization_org
  }
  property {
    name  = "actionsGetGithubActionsPermissionsRepository_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_github_actions_permissions_repository_owner
  }
  property {
    name  = "actionsGetGithubActionsPermissionsRepository_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_github_actions_permissions_repository_repo
  }
  property {
    name  = "actionsGetJobForWorkflowRun_job_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_job_for_workflow_run_job_id
  }
  property {
    name  = "actionsGetJobForWorkflowRun_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_job_for_workflow_run_owner
  }
  property {
    name  = "actionsGetJobForWorkflowRun_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_job_for_workflow_run_repo
  }
  property {
    name  = "actionsGetOrgPublicKey_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_org_public_key_org
  }
  property {
    name  = "actionsGetOrgSecret_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_org_secret_org
  }
  property {
    name  = "actionsGetOrgSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_org_secret_secret_name
  }
  property {
    name  = "actionsGetRepoPublicKey_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_repo_public_key_owner
  }
  property {
    name  = "actionsGetRepoPublicKey_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_repo_public_key_repo
  }
  property {
    name  = "actionsGetRepoSecret_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_repo_secret_owner
  }
  property {
    name  = "actionsGetRepoSecret_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_repo_secret_repo
  }
  property {
    name  = "actionsGetRepoSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_repo_secret_secret_name
  }
  property {
    name  = "actionsGetSelfHostedRunnerForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_self_hosted_runner_for_org_org
  }
  property {
    name  = "actionsGetSelfHostedRunnerForOrg_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_self_hosted_runner_for_org_runner_id
  }
  property {
    name  = "actionsGetSelfHostedRunnerForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_self_hosted_runner_for_repo_owner
  }
  property {
    name  = "actionsGetSelfHostedRunnerForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_self_hosted_runner_for_repo_repo
  }
  property {
    name  = "actionsGetSelfHostedRunnerForRepo_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_self_hosted_runner_for_repo_runner_id
  }
  property {
    name  = "actionsGetSelfHostedRunnerGroupForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_self_hosted_runner_group_for_org_org
  }
  property {
    name  = "actionsGetSelfHostedRunnerGroupForOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_self_hosted_runner_group_for_org_runner_group_id
  }
  property {
    name  = "actionsGetWorkflowRun_exclude_pull_requests"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_workflow_run_exclude_pull_requests
  }
  property {
    name  = "actionsGetWorkflowRun_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_workflow_run_owner
  }
  property {
    name  = "actionsGetWorkflowRun_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_workflow_run_repo
  }
  property {
    name  = "actionsGetWorkflowRun_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_workflow_run_run_id
  }
  property {
    name  = "actionsGetWorkflow_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_workflow_owner
  }
  property {
    name  = "actionsGetWorkflow_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_get_workflow_repo
  }
  property {
    name  = "actionsListArtifactsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_artifacts_for_repo_owner
  }
  property {
    name  = "actionsListArtifactsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_artifacts_for_repo_page
  }
  property {
    name  = "actionsListArtifactsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_artifacts_for_repo_per_page
  }
  property {
    name  = "actionsListArtifactsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_artifacts_for_repo_repo
  }
  property {
    name  = "actionsListJobsForWorkflowRun_filter"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_jobs_for_workflow_run_filter
  }
  property {
    name  = "actionsListJobsForWorkflowRun_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_jobs_for_workflow_run_owner
  }
  property {
    name  = "actionsListJobsForWorkflowRun_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_jobs_for_workflow_run_page
  }
  property {
    name  = "actionsListJobsForWorkflowRun_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_jobs_for_workflow_run_per_page
  }
  property {
    name  = "actionsListJobsForWorkflowRun_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_jobs_for_workflow_run_repo
  }
  property {
    name  = "actionsListJobsForWorkflowRun_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_jobs_for_workflow_run_run_id
  }
  property {
    name  = "actionsListOrgSecrets_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_org_secrets_org
  }
  property {
    name  = "actionsListOrgSecrets_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_org_secrets_page
  }
  property {
    name  = "actionsListOrgSecrets_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_org_secrets_per_page
  }
  property {
    name  = "actionsListRepoAccessToSelfHostedRunnerGroupInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_access_to_self_hosted_runner_group_in_org_org
  }
  property {
    name  = "actionsListRepoAccessToSelfHostedRunnerGroupInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_access_to_self_hosted_runner_group_in_org_page
  }
  property {
    name  = "actionsListRepoAccessToSelfHostedRunnerGroupInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_access_to_self_hosted_runner_group_in_org_per_page
  }
  property {
    name  = "actionsListRepoAccessToSelfHostedRunnerGroupInOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_access_to_self_hosted_runner_group_in_org_runner_group_id
  }
  property {
    name  = "actionsListRepoSecrets_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_secrets_owner
  }
  property {
    name  = "actionsListRepoSecrets_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_secrets_page
  }
  property {
    name  = "actionsListRepoSecrets_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_secrets_per_page
  }
  property {
    name  = "actionsListRepoSecrets_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_secrets_repo
  }
  property {
    name  = "actionsListRepoWorkflows_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_workflows_owner
  }
  property {
    name  = "actionsListRepoWorkflows_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_workflows_page
  }
  property {
    name  = "actionsListRepoWorkflows_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_workflows_per_page
  }
  property {
    name  = "actionsListRepoWorkflows_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_repo_workflows_repo
  }
  property {
    name  = "actionsListRunnerApplicationsForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_runner_applications_for_org_org
  }
  property {
    name  = "actionsListRunnerApplicationsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_runner_applications_for_repo_owner
  }
  property {
    name  = "actionsListRunnerApplicationsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_runner_applications_for_repo_repo
  }
  property {
    name  = "actionsListSelectedReposForOrgSecret_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_selected_repos_for_org_secret_org
  }
  property {
    name  = "actionsListSelectedReposForOrgSecret_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_selected_repos_for_org_secret_page
  }
  property {
    name  = "actionsListSelectedReposForOrgSecret_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_selected_repos_for_org_secret_per_page
  }
  property {
    name  = "actionsListSelectedReposForOrgSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_selected_repos_for_org_secret_secret_name
  }
  property {
    name  = "actionsListSelectedRepositoriesEnabledGithubActionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_selected_repositories_enabled_github_actions_organization_org
  }
  property {
    name  = "actionsListSelectedRepositoriesEnabledGithubActionsOrganization_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_selected_repositories_enabled_github_actions_organization_page
  }
  property {
    name  = "actionsListSelectedRepositoriesEnabledGithubActionsOrganization_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_selected_repositories_enabled_github_actions_organization_per_page
  }
  property {
    name  = "actionsListSelfHostedRunnerGroupsForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runner_groups_for_org_org
  }
  property {
    name  = "actionsListSelfHostedRunnerGroupsForOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runner_groups_for_org_page
  }
  property {
    name  = "actionsListSelfHostedRunnerGroupsForOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runner_groups_for_org_per_page
  }
  property {
    name  = "actionsListSelfHostedRunnersForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_for_org_org
  }
  property {
    name  = "actionsListSelfHostedRunnersForOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_for_org_page
  }
  property {
    name  = "actionsListSelfHostedRunnersForOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_for_org_per_page
  }
  property {
    name  = "actionsListSelfHostedRunnersForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_for_repo_owner
  }
  property {
    name  = "actionsListSelfHostedRunnersForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_for_repo_page
  }
  property {
    name  = "actionsListSelfHostedRunnersForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_for_repo_per_page
  }
  property {
    name  = "actionsListSelfHostedRunnersForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_for_repo_repo
  }
  property {
    name  = "actionsListSelfHostedRunnersInGroupForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_in_group_for_org_org
  }
  property {
    name  = "actionsListSelfHostedRunnersInGroupForOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_in_group_for_org_page
  }
  property {
    name  = "actionsListSelfHostedRunnersInGroupForOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_in_group_for_org_per_page
  }
  property {
    name  = "actionsListSelfHostedRunnersInGroupForOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_self_hosted_runners_in_group_for_org_runner_group_id
  }
  property {
    name  = "actionsListWorkflowRunArtifacts_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_run_artifacts_owner
  }
  property {
    name  = "actionsListWorkflowRunArtifacts_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_run_artifacts_page
  }
  property {
    name  = "actionsListWorkflowRunArtifacts_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_run_artifacts_per_page
  }
  property {
    name  = "actionsListWorkflowRunArtifacts_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_run_artifacts_repo
  }
  property {
    name  = "actionsListWorkflowRunArtifacts_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_run_artifacts_run_id
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_actor"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_actor
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_branch"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_branch
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_created"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_created
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_event"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_event
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_exclude_pull_requests"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_exclude_pull_requests
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_owner
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_page
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_per_page
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_repo
  }
  property {
    name  = "actionsListWorkflowRunsForRepo_status"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_for_repo_status
  }
  property {
    name  = "actionsListWorkflowRuns_actor"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_actor
  }
  property {
    name  = "actionsListWorkflowRuns_branch"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_branch
  }
  property {
    name  = "actionsListWorkflowRuns_created"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_created
  }
  property {
    name  = "actionsListWorkflowRuns_event"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_event
  }
  property {
    name  = "actionsListWorkflowRuns_exclude_pull_requests"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_exclude_pull_requests
  }
  property {
    name  = "actionsListWorkflowRuns_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_owner
  }
  property {
    name  = "actionsListWorkflowRuns_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_page
  }
  property {
    name  = "actionsListWorkflowRuns_per_page"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_per_page
  }
  property {
    name  = "actionsListWorkflowRuns_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_repo
  }
  property {
    name  = "actionsListWorkflowRuns_status"
    type  = "string"
    value = var.connector-oai-github_property_actions_list_workflow_runs_status
  }
  property {
    name  = "actionsReRunWorkflow_body"
    type  = "string"
    value = var.connector-oai-github_property_actions_re_run_workflow_body
  }
  property {
    name  = "actionsReRunWorkflow_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_re_run_workflow_owner
  }
  property {
    name  = "actionsReRunWorkflow_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_re_run_workflow_repo
  }
  property {
    name  = "actionsReRunWorkflow_run_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_re_run_workflow_run_id
  }
  property {
    name  = "actionsRemoveRepoAccessToSelfHostedRunnerGroupInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_repo_access_to_self_hosted_runner_group_in_org_org
  }
  property {
    name  = "actionsRemoveRepoAccessToSelfHostedRunnerGroupInOrg_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_repo_access_to_self_hosted_runner_group_in_org_repository_id
  }
  property {
    name  = "actionsRemoveRepoAccessToSelfHostedRunnerGroupInOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_repo_access_to_self_hosted_runner_group_in_org_runner_group_id
  }
  property {
    name  = "actionsRemoveSelectedRepoFromOrgSecret_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_selected_repo_from_org_secret_org
  }
  property {
    name  = "actionsRemoveSelectedRepoFromOrgSecret_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_selected_repo_from_org_secret_repository_id
  }
  property {
    name  = "actionsRemoveSelectedRepoFromOrgSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_selected_repo_from_org_secret_secret_name
  }
  property {
    name  = "actionsRemoveSelfHostedRunnerFromGroupForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_self_hosted_runner_from_group_for_org_org
  }
  property {
    name  = "actionsRemoveSelfHostedRunnerFromGroupForOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_self_hosted_runner_from_group_for_org_runner_group_id
  }
  property {
    name  = "actionsRemoveSelfHostedRunnerFromGroupForOrg_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_remove_self_hosted_runner_from_group_for_org_runner_id
  }
  property {
    name  = "actionsSetAllowedActionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_allowed_actions_organization_org
  }
  property {
    name  = "actionsSetAllowedActionsOrganization_selectedActions_SelectedActions_github_owned_allowed"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_allowed_actions_organization_selected_actions_selected_actions_github_owned_allowed
  }
  property {
    name  = "actionsSetAllowedActionsOrganization_selectedActions_SelectedActions_patterns_allowed"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_allowed_actions_organization_selected_actions_selected_actions_patterns_allowed
  }
  property {
    name  = "actionsSetAllowedActionsRepository_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_allowed_actions_repository_owner
  }
  property {
    name  = "actionsSetAllowedActionsRepository_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_allowed_actions_repository_repo
  }
  property {
    name  = "actionsSetAllowedActionsRepository_selectedActions_SelectedActions_github_owned_allowed"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_allowed_actions_repository_selected_actions_selected_actions_github_owned_allowed
  }
  property {
    name  = "actionsSetAllowedActionsRepository_selectedActions_SelectedActions_patterns_allowed"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_allowed_actions_repository_selected_actions_selected_actions_patterns_allowed
  }
  property {
    name  = "actionsSetGithubActionsPermissionsOrganization_actionsSetGithubActionsPermissionsOrganizationRequest_ActionsSetGithubActionsPermissionsOrganizationRequest_allowed_actions"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_github_actions_permissions_organization_actions_set_github_actions_permissions_organization_request_actions_set_github_actions_permissions_organization_request_allowed_actions
  }
  property {
    name  = "actionsSetGithubActionsPermissionsOrganization_actionsSetGithubActionsPermissionsOrganizationRequest_ActionsSetGithubActionsPermissionsOrganizationRequest_enabled_repositories"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_github_actions_permissions_organization_actions_set_github_actions_permissions_organization_request_actions_set_github_actions_permissions_organization_request_enabled_repositories
  }
  property {
    name  = "actionsSetGithubActionsPermissionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_github_actions_permissions_organization_org
  }
  property {
    name  = "actionsSetGithubActionsPermissionsRepository_actionsSetGithubActionsPermissionsRepositoryRequest_ActionsSetGithubActionsPermissionsRepositoryRequest_allowed_actions"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_github_actions_permissions_repository_actions_set_github_actions_permissions_repository_request_actions_set_github_actions_permissions_repository_request_allowed_actions
  }
  property {
    name  = "actionsSetGithubActionsPermissionsRepository_actionsSetGithubActionsPermissionsRepositoryRequest_ActionsSetGithubActionsPermissionsRepositoryRequest_enabled"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_github_actions_permissions_repository_actions_set_github_actions_permissions_repository_request_actions_set_github_actions_permissions_repository_request_enabled
  }
  property {
    name  = "actionsSetGithubActionsPermissionsRepository_owner"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_github_actions_permissions_repository_owner
  }
  property {
    name  = "actionsSetGithubActionsPermissionsRepository_repo"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_github_actions_permissions_repository_repo
  }
  property {
    name  = "actionsSetRepoAccessToSelfHostedRunnerGroupInOrg_actionsSetRepoAccessToSelfHostedRunnerGroupInOrgRequest_ActionsSetRepoAccessToSelfHostedRunnerGroupInOrgRequest_selected_repository_ids"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_repo_access_to_self_hosted_runner_group_in_org_actions_set_repo_access_to_self_hosted_runner_group_in_org_request_actions_set_repo_access_to_self_hosted_runner_group_in_org_request_selected_repository_ids
  }
  property {
    name  = "actionsSetRepoAccessToSelfHostedRunnerGroupInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_repo_access_to_self_hosted_runner_group_in_org_org
  }
  property {
    name  = "actionsSetRepoAccessToSelfHostedRunnerGroupInOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_repo_access_to_self_hosted_runner_group_in_org_runner_group_id
  }
  property {
    name  = "actionsSetSelectedReposForOrgSecret_actionsSetSelectedReposForOrgSecretRequest_ActionsSetSelectedReposForOrgSecretRequest_selected_repository_ids"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_selected_repos_for_org_secret_actions_set_selected_repos_for_org_secret_request_actions_set_selected_repos_for_org_secret_request_selected_repository_ids
  }
  property {
    name  = "actionsSetSelectedReposForOrgSecret_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_selected_repos_for_org_secret_org
  }
  property {
    name  = "actionsSetSelectedReposForOrgSecret_secret_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_selected_repos_for_org_secret_secret_name
  }
  property {
    name  = "actionsSetSelectedRepositoriesEnabledGithubActionsOrganization_actionsSetSelectedRepositoriesEnabledGithubActionsOrganizationRequest_ActionsSetSelectedRepositoriesEnabledGithubActionsOrganizationRequest_selected_repository_ids"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_selected_repositories_enabled_github_actions_organization_actions_set_selected_repositories_enabled_github_actions_organization_request_actions_set_selected_repositories_enabled_github_actions_organization_request_selected_repository_ids
  }
  property {
    name  = "actionsSetSelectedRepositoriesEnabledGithubActionsOrganization_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_selected_repositories_enabled_github_actions_organization_org
  }
  property {
    name  = "actionsSetSelfHostedRunnersInGroupForOrg_enterpriseAdminSetSelfHostedRunnersInGroupForEnterpriseRequest_EnterpriseAdminSetSelfHostedRunnersInGroupForEnterpriseRequest_runners"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_self_hosted_runners_in_group_for_org_enterprise_admin_set_self_hosted_runners_in_group_for_enterprise_request_enterprise_admin_set_self_hosted_runners_in_group_for_enterprise_request_runners
  }
  property {
    name  = "actionsSetSelfHostedRunnersInGroupForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_self_hosted_runners_in_group_for_org_org
  }
  property {
    name  = "actionsSetSelfHostedRunnersInGroupForOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_set_self_hosted_runners_in_group_for_org_runner_group_id
  }
  property {
    name  = "actionsUpdateSelfHostedRunnerGroupForOrg_actionsUpdateSelfHostedRunnerGroupForOrgRequest_ActionsUpdateSelfHostedRunnerGroupForOrgRequest_allows_public_repositories"
    type  = "string"
    value = var.connector-oai-github_property_actions_update_self_hosted_runner_group_for_org_actions_update_self_hosted_runner_group_for_org_request_actions_update_self_hosted_runner_group_for_org_request_allows_public_repositories
  }
  property {
    name  = "actionsUpdateSelfHostedRunnerGroupForOrg_actionsUpdateSelfHostedRunnerGroupForOrgRequest_ActionsUpdateSelfHostedRunnerGroupForOrgRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_actions_update_self_hosted_runner_group_for_org_actions_update_self_hosted_runner_group_for_org_request_actions_update_self_hosted_runner_group_for_org_request_name
  }
  property {
    name  = "actionsUpdateSelfHostedRunnerGroupForOrg_actionsUpdateSelfHostedRunnerGroupForOrgRequest_ActionsUpdateSelfHostedRunnerGroupForOrgRequest_visibility"
    type  = "string"
    value = var.connector-oai-github_property_actions_update_self_hosted_runner_group_for_org_actions_update_self_hosted_runner_group_for_org_request_actions_update_self_hosted_runner_group_for_org_request_visibility
  }
  property {
    name  = "actionsUpdateSelfHostedRunnerGroupForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_actions_update_self_hosted_runner_group_for_org_org
  }
  property {
    name  = "actionsUpdateSelfHostedRunnerGroupForOrg_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_actions_update_self_hosted_runner_group_for_org_runner_group_id
  }
  property {
    name  = "activityCheckRepoIsStarredByAuthenticatedUser_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_check_repo_is_starred_by_authenticated_user_owner
  }
  property {
    name  = "activityCheckRepoIsStarredByAuthenticatedUser_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_check_repo_is_starred_by_authenticated_user_repo
  }
  property {
    name  = "activityDeleteRepoSubscription_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_delete_repo_subscription_owner
  }
  property {
    name  = "activityDeleteRepoSubscription_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_delete_repo_subscription_repo
  }
  property {
    name  = "activityDeleteThreadSubscription_thread_id"
    type  = "string"
    value = var.connector-oai-github_property_activity_delete_thread_subscription_thread_id
  }
  property {
    name  = "activityGetRepoSubscription_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_get_repo_subscription_owner
  }
  property {
    name  = "activityGetRepoSubscription_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_get_repo_subscription_repo
  }
  property {
    name  = "activityGetThreadSubscriptionForAuthenticatedUser_thread_id"
    type  = "string"
    value = var.connector-oai-github_property_activity_get_thread_subscription_for_authenticated_user_thread_id
  }
  property {
    name  = "activityGetThread_thread_id"
    type  = "string"
    value = var.connector-oai-github_property_activity_get_thread_thread_id
  }
  property {
    name  = "activityListEventsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_events_for_authenticated_user_page
  }
  property {
    name  = "activityListEventsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_events_for_authenticated_user_per_page
  }
  property {
    name  = "activityListEventsForAuthenticatedUser_username"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_events_for_authenticated_user_username
  }
  property {
    name  = "activityListNotificationsForAuthenticatedUser_all"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_notifications_for_authenticated_user_all
  }
  property {
    name  = "activityListNotificationsForAuthenticatedUser_before"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_notifications_for_authenticated_user_before
  }
  property {
    name  = "activityListNotificationsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_notifications_for_authenticated_user_page
  }
  property {
    name  = "activityListNotificationsForAuthenticatedUser_participating"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_notifications_for_authenticated_user_participating
  }
  property {
    name  = "activityListNotificationsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_notifications_for_authenticated_user_per_page
  }
  property {
    name  = "activityListNotificationsForAuthenticatedUser_since"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_notifications_for_authenticated_user_since
  }
  property {
    name  = "activityListOrgEventsForAuthenticatedUser_org"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_org_events_for_authenticated_user_org
  }
  property {
    name  = "activityListOrgEventsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_org_events_for_authenticated_user_page
  }
  property {
    name  = "activityListOrgEventsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_org_events_for_authenticated_user_per_page
  }
  property {
    name  = "activityListOrgEventsForAuthenticatedUser_username"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_org_events_for_authenticated_user_username
  }
  property {
    name  = "activityListPublicEventsForRepoNetwork_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_for_repo_network_owner
  }
  property {
    name  = "activityListPublicEventsForRepoNetwork_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_for_repo_network_page
  }
  property {
    name  = "activityListPublicEventsForRepoNetwork_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_for_repo_network_per_page
  }
  property {
    name  = "activityListPublicEventsForRepoNetwork_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_for_repo_network_repo
  }
  property {
    name  = "activityListPublicEventsForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_for_user_page
  }
  property {
    name  = "activityListPublicEventsForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_for_user_per_page
  }
  property {
    name  = "activityListPublicEventsForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_for_user_username
  }
  property {
    name  = "activityListPublicEvents_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_page
  }
  property {
    name  = "activityListPublicEvents_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_events_per_page
  }
  property {
    name  = "activityListPublicOrgEvents_org"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_org_events_org
  }
  property {
    name  = "activityListPublicOrgEvents_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_org_events_page
  }
  property {
    name  = "activityListPublicOrgEvents_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_public_org_events_per_page
  }
  property {
    name  = "activityListReceivedEventsForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_received_events_for_user_page
  }
  property {
    name  = "activityListReceivedEventsForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_received_events_for_user_per_page
  }
  property {
    name  = "activityListReceivedEventsForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_received_events_for_user_username
  }
  property {
    name  = "activityListReceivedPublicEventsForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_received_public_events_for_user_page
  }
  property {
    name  = "activityListReceivedPublicEventsForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_received_public_events_for_user_per_page
  }
  property {
    name  = "activityListReceivedPublicEventsForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_received_public_events_for_user_username
  }
  property {
    name  = "activityListRepoEvents_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_events_owner
  }
  property {
    name  = "activityListRepoEvents_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_events_page
  }
  property {
    name  = "activityListRepoEvents_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_events_per_page
  }
  property {
    name  = "activityListRepoEvents_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_events_repo
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_all"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_all
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_before"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_before
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_owner
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_page
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_participating"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_participating
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_per_page
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_repo
  }
  property {
    name  = "activityListRepoNotificationsForAuthenticatedUser_since"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repo_notifications_for_authenticated_user_since
  }
  property {
    name  = "activityListReposStarredByAuthenticatedUser_direction"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_authenticated_user_direction
  }
  property {
    name  = "activityListReposStarredByAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_authenticated_user_page
  }
  property {
    name  = "activityListReposStarredByAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_authenticated_user_per_page
  }
  property {
    name  = "activityListReposStarredByAuthenticatedUser_sort"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_authenticated_user_sort
  }
  property {
    name  = "activityListReposStarredByUser_direction"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_user_direction
  }
  property {
    name  = "activityListReposStarredByUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_user_page
  }
  property {
    name  = "activityListReposStarredByUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_user_per_page
  }
  property {
    name  = "activityListReposStarredByUser_sort"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_user_sort
  }
  property {
    name  = "activityListReposStarredByUser_username"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_starred_by_user_username
  }
  property {
    name  = "activityListReposWatchedByUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_watched_by_user_page
  }
  property {
    name  = "activityListReposWatchedByUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_watched_by_user_per_page
  }
  property {
    name  = "activityListReposWatchedByUser_username"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_repos_watched_by_user_username
  }
  property {
    name  = "activityListStargazersForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_stargazers_for_repo_owner
  }
  property {
    name  = "activityListStargazersForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_stargazers_for_repo_page
  }
  property {
    name  = "activityListStargazersForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_stargazers_for_repo_per_page
  }
  property {
    name  = "activityListStargazersForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_stargazers_for_repo_repo
  }
  property {
    name  = "activityListWatchedReposForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_watched_repos_for_authenticated_user_page
  }
  property {
    name  = "activityListWatchedReposForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_watched_repos_for_authenticated_user_per_page
  }
  property {
    name  = "activityListWatchersForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_watchers_for_repo_owner
  }
  property {
    name  = "activityListWatchersForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_watchers_for_repo_page
  }
  property {
    name  = "activityListWatchersForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_watchers_for_repo_per_page
  }
  property {
    name  = "activityListWatchersForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_list_watchers_for_repo_repo
  }
  property {
    name  = "activityMarkNotificationsAsRead_activityMarkNotificationsAsReadRequest_ActivityMarkNotificationsAsReadRequest_last_read_at"
    type  = "string"
    value = var.connector-oai-github_property_activity_mark_notifications_as_read_activity_mark_notifications_as_read_request_activity_mark_notifications_as_read_request_last_read_at
  }
  property {
    name  = "activityMarkNotificationsAsRead_activityMarkNotificationsAsReadRequest_ActivityMarkNotificationsAsReadRequest_read"
    type  = "string"
    value = var.connector-oai-github_property_activity_mark_notifications_as_read_activity_mark_notifications_as_read_request_activity_mark_notifications_as_read_request_read
  }
  property {
    name  = "activityMarkRepoNotificationsAsRead_activityMarkRepoNotificationsAsReadRequest_ActivityMarkRepoNotificationsAsReadRequest_last_read_at"
    type  = "string"
    value = var.connector-oai-github_property_activity_mark_repo_notifications_as_read_activity_mark_repo_notifications_as_read_request_activity_mark_repo_notifications_as_read_request_last_read_at
  }
  property {
    name  = "activityMarkRepoNotificationsAsRead_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_mark_repo_notifications_as_read_owner
  }
  property {
    name  = "activityMarkRepoNotificationsAsRead_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_mark_repo_notifications_as_read_repo
  }
  property {
    name  = "activityMarkThreadAsRead_thread_id"
    type  = "string"
    value = var.connector-oai-github_property_activity_mark_thread_as_read_thread_id
  }
  property {
    name  = "activitySetRepoSubscription_activitySetRepoSubscriptionRequest_ActivitySetRepoSubscriptionRequest_ignored"
    type  = "string"
    value = var.connector-oai-github_property_activity_set_repo_subscription_activity_set_repo_subscription_request_activity_set_repo_subscription_request_ignored
  }
  property {
    name  = "activitySetRepoSubscription_activitySetRepoSubscriptionRequest_ActivitySetRepoSubscriptionRequest_subscribed"
    type  = "string"
    value = var.connector-oai-github_property_activity_set_repo_subscription_activity_set_repo_subscription_request_activity_set_repo_subscription_request_subscribed
  }
  property {
    name  = "activitySetRepoSubscription_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_set_repo_subscription_owner
  }
  property {
    name  = "activitySetRepoSubscription_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_set_repo_subscription_repo
  }
  property {
    name  = "activitySetThreadSubscription_activitySetThreadSubscriptionRequest_ActivitySetThreadSubscriptionRequest_ignored"
    type  = "string"
    value = var.connector-oai-github_property_activity_set_thread_subscription_activity_set_thread_subscription_request_activity_set_thread_subscription_request_ignored
  }
  property {
    name  = "activitySetThreadSubscription_thread_id"
    type  = "string"
    value = var.connector-oai-github_property_activity_set_thread_subscription_thread_id
  }
  property {
    name  = "activityStarRepoForAuthenticatedUser_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_star_repo_for_authenticated_user_owner
  }
  property {
    name  = "activityStarRepoForAuthenticatedUser_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_star_repo_for_authenticated_user_repo
  }
  property {
    name  = "activityUnstarRepoForAuthenticatedUser_owner"
    type  = "string"
    value = var.connector-oai-github_property_activity_unstar_repo_for_authenticated_user_owner
  }
  property {
    name  = "activityUnstarRepoForAuthenticatedUser_repo"
    type  = "string"
    value = var.connector-oai-github_property_activity_unstar_repo_for_authenticated_user_repo
  }
  property {
    name  = "apiVersion"
    type  = "string"
    value = var.connector-oai-github_property_api_version
  }
  property {
    name  = "appsAddRepoToInstallationForAuthenticatedUser_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_add_repo_to_installation_for_authenticated_user_installation_id
  }
  property {
    name  = "appsAddRepoToInstallationForAuthenticatedUser_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_add_repo_to_installation_for_authenticated_user_repository_id
  }
  property {
    name  = "appsCheckAuthorization_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_check_authorization_access_token
  }
  property {
    name  = "appsCheckAuthorization_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_check_authorization_client_id
  }
  property {
    name  = "appsCheckToken_appsCheckTokenRequest_AppsCheckTokenRequest_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_check_token_apps_check_token_request_apps_check_token_request_access_token
  }
  property {
    name  = "appsCheckToken_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_check_token_client_id
  }
  property {
    name  = "appsCreateContentAttachment_appsCreateContentAttachmentRequest_AppsCreateContentAttachmentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_content_attachment_apps_create_content_attachment_request_apps_create_content_attachment_request_body
  }
  property {
    name  = "appsCreateContentAttachment_appsCreateContentAttachmentRequest_AppsCreateContentAttachmentRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_content_attachment_apps_create_content_attachment_request_apps_create_content_attachment_request_title
  }
  property {
    name  = "appsCreateContentAttachment_content_reference_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_content_attachment_content_reference_id
  }
  property {
    name  = "appsCreateContentAttachment_owner"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_content_attachment_owner
  }
  property {
    name  = "appsCreateContentAttachment_repo"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_content_attachment_repo
  }
  property {
    name  = "appsCreateFromManifest_code"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_from_manifest_code
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_actions"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_actions
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_administration"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_administration
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_checks"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_checks
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_content_references"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_content_references
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_contents"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_contents
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_deployments"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_deployments
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_environments"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_environments
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_issues"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_issues
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_members"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_members
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_metadata"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_metadata
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_administration"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_administration
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_hooks"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_hooks
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_packages"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_packages
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_plan"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_plan
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_projects"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_projects
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_secrets"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_secrets
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_self_hosted_runners"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_self_hosted_runners
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_organization_user_blocking"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_organization_user_blocking
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_packages"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_packages
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_pages"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_pages
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_pull_requests"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_pull_requests
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_repository_hooks"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_repository_hooks
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_repository_projects"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_repository_projects
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_secret_scanning_alerts"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_secret_scanning_alerts
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_secrets"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_secrets
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_security_events"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_security_events
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_single_file"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_single_file
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_statuses"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_statuses
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_team_discussions"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_team_discussions
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_vulnerability_alerts"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_vulnerability_alerts
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppPermissions_workflows"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_app_permissions_workflows
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppsCreateInstallationAccessTokenRequest_repositories"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_apps_create_installation_access_token_request_repositories
  }
  property {
    name  = "appsCreateInstallationAccessToken_appsCreateInstallationAccessTokenRequest_AppsCreateInstallationAccessTokenRequest_repository_ids"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_apps_create_installation_access_token_request_apps_create_installation_access_token_request_repository_ids
  }
  property {
    name  = "appsCreateInstallationAccessToken_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_create_installation_access_token_installation_id
  }
  property {
    name  = "appsDeleteAuthorization_appsDeleteAuthorizationRequest_AppsDeleteAuthorizationRequest_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_delete_authorization_apps_delete_authorization_request_apps_delete_authorization_request_access_token
  }
  property {
    name  = "appsDeleteAuthorization_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_delete_authorization_client_id
  }
  property {
    name  = "appsDeleteInstallation_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_delete_installation_installation_id
  }
  property {
    name  = "appsDeleteToken_appsDeleteAuthorizationRequest_AppsDeleteAuthorizationRequest_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_delete_token_apps_delete_authorization_request_apps_delete_authorization_request_access_token
  }
  property {
    name  = "appsDeleteToken_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_delete_token_client_id
  }
  property {
    name  = "appsGetBySlug_app_slug"
    type  = "string"
    value = var.connector-oai-github_property_apps_get_by_slug_app_slug
  }
  property {
    name  = "appsGetInstallation_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_get_installation_installation_id
  }
  property {
    name  = "appsGetOrgInstallation_org"
    type  = "string"
    value = var.connector-oai-github_property_apps_get_org_installation_org
  }
  property {
    name  = "appsGetRepoInstallation_owner"
    type  = "string"
    value = var.connector-oai-github_property_apps_get_repo_installation_owner
  }
  property {
    name  = "appsGetRepoInstallation_repo"
    type  = "string"
    value = var.connector-oai-github_property_apps_get_repo_installation_repo
  }
  property {
    name  = "appsGetUserInstallation_username"
    type  = "string"
    value = var.connector-oai-github_property_apps_get_user_installation_username
  }
  property {
    name  = "appsListInstallationReposForAuthenticatedUser_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installation_repos_for_authenticated_user_installation_id
  }
  property {
    name  = "appsListInstallationReposForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installation_repos_for_authenticated_user_page
  }
  property {
    name  = "appsListInstallationReposForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installation_repos_for_authenticated_user_per_page
  }
  property {
    name  = "appsListInstallationsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installations_for_authenticated_user_page
  }
  property {
    name  = "appsListInstallationsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installations_for_authenticated_user_per_page
  }
  property {
    name  = "appsListInstallations_outdated"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installations_outdated
  }
  property {
    name  = "appsListInstallations_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installations_page
  }
  property {
    name  = "appsListInstallations_per_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installations_per_page
  }
  property {
    name  = "appsListInstallations_since"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_installations_since
  }
  property {
    name  = "appsListReposAccessibleToInstallation_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_repos_accessible_to_installation_page
  }
  property {
    name  = "appsListReposAccessibleToInstallation_per_page"
    type  = "string"
    value = var.connector-oai-github_property_apps_list_repos_accessible_to_installation_per_page
  }
  property {
    name  = "appsRemoveRepoFromInstallationForAuthenticatedUser_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_remove_repo_from_installation_for_authenticated_user_installation_id
  }
  property {
    name  = "appsRemoveRepoFromInstallationForAuthenticatedUser_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_remove_repo_from_installation_for_authenticated_user_repository_id
  }
  property {
    name  = "appsResetAuthorization_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_reset_authorization_access_token
  }
  property {
    name  = "appsResetAuthorization_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_reset_authorization_client_id
  }
  property {
    name  = "appsResetToken_appsCheckTokenRequest_AppsCheckTokenRequest_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_reset_token_apps_check_token_request_apps_check_token_request_access_token
  }
  property {
    name  = "appsResetToken_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_reset_token_client_id
  }
  property {
    name  = "appsRevokeAuthorizationForApplication_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_revoke_authorization_for_application_access_token
  }
  property {
    name  = "appsRevokeAuthorizationForApplication_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_revoke_authorization_for_application_client_id
  }
  property {
    name  = "appsRevokeGrantForApplication_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_revoke_grant_for_application_access_token
  }
  property {
    name  = "appsRevokeGrantForApplication_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_revoke_grant_for_application_client_id
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_actions"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_actions
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_administration"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_administration
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_checks"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_checks
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_content_references"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_content_references
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_contents"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_contents
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_deployments"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_deployments
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_environments"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_environments
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_issues"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_issues
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_members"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_members
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_metadata"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_metadata
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_administration"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_administration
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_hooks"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_hooks
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_packages"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_packages
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_plan"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_plan
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_projects"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_projects
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_secrets"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_secrets
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_self_hosted_runners"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_self_hosted_runners
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_organization_user_blocking"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_organization_user_blocking
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_packages"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_packages
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_pages"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_pages
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_pull_requests"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_pull_requests
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_repository_hooks"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_repository_hooks
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_repository_projects"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_repository_projects
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_secret_scanning_alerts"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_secret_scanning_alerts
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_secrets"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_secrets
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_security_events"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_security_events
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_single_file"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_single_file
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_statuses"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_statuses
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_team_discussions"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_team_discussions
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_vulnerability_alerts"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_vulnerability_alerts
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppPermissions_workflows"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_app_permissions_workflows
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppsScopeTokenRequest_access_token"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_apps_scope_token_request_access_token
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppsScopeTokenRequest_repositories"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_apps_scope_token_request_repositories
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppsScopeTokenRequest_repository_ids"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_apps_scope_token_request_repository_ids
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppsScopeTokenRequest_target"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_apps_scope_token_request_target
  }
  property {
    name  = "appsScopeToken_appsScopeTokenRequest_AppsScopeTokenRequest_target_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_apps_scope_token_request_apps_scope_token_request_target_id
  }
  property {
    name  = "appsScopeToken_client_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_scope_token_client_id
  }
  property {
    name  = "appsSuspendInstallation_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_suspend_installation_installation_id
  }
  property {
    name  = "appsUnsuspendInstallation_installation_id"
    type  = "string"
    value = var.connector-oai-github_property_apps_unsuspend_installation_installation_id
  }
  property {
    name  = "appsUpdateWebhookConfigForApp_appsUpdateWebhookConfigForAppRequest_AppsUpdateWebhookConfigForAppRequest_content_type"
    type  = "string"
    value = var.connector-oai-github_property_apps_update_webhook_config_for_app_apps_update_webhook_config_for_app_request_apps_update_webhook_config_for_app_request_content_type
  }
  property {
    name  = "appsUpdateWebhookConfigForApp_appsUpdateWebhookConfigForAppRequest_AppsUpdateWebhookConfigForAppRequest_secret"
    type  = "string"
    value = var.connector-oai-github_property_apps_update_webhook_config_for_app_apps_update_webhook_config_for_app_request_apps_update_webhook_config_for_app_request_secret
  }
  property {
    name  = "appsUpdateWebhookConfigForApp_appsUpdateWebhookConfigForAppRequest_AppsUpdateWebhookConfigForAppRequest_url"
    type  = "string"
    value = var.connector-oai-github_property_apps_update_webhook_config_for_app_apps_update_webhook_config_for_app_request_apps_update_webhook_config_for_app_request_url
  }
  property {
    name  = "authBearerToken"
    type  = "string"
    value = var.connector-oai-github_property_auth_bearer_token
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-github_property_base_path
  }
  property {
    name  = "checksCreateSuite_checksCreateSuiteRequest_ChecksCreateSuiteRequest_head_sha"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_suite_checks_create_suite_request_checks_create_suite_request_head_sha
  }
  property {
    name  = "checksCreateSuite_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_suite_owner
  }
  property {
    name  = "checksCreateSuite_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_suite_repo
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequestOutput_annotations"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_output_annotations
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequestOutput_images"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_output_images
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequestOutput_summary"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_output_summary
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequestOutput_text"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_output_text
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequestOutput_title"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_output_title
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_actions"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_actions
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_completed_at"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_completed_at
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_conclusion"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_conclusion
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_details_url"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_details_url
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_external_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_external_id
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_head_sha"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_head_sha
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_name
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_started_at"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_started_at
  }
  property {
    name  = "checksCreate_checksCreateRequest_ChecksCreateRequest_status"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_checks_create_request_checks_create_request_status
  }
  property {
    name  = "checksCreate_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_owner
  }
  property {
    name  = "checksCreate_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_create_repo
  }
  property {
    name  = "checksGetSuite_check_suite_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_get_suite_check_suite_id
  }
  property {
    name  = "checksGetSuite_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_get_suite_owner
  }
  property {
    name  = "checksGetSuite_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_get_suite_repo
  }
  property {
    name  = "checksGet_check_run_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_get_check_run_id
  }
  property {
    name  = "checksGet_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_get_owner
  }
  property {
    name  = "checksGet_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_get_repo
  }
  property {
    name  = "checksListAnnotations_check_run_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_annotations_check_run_id
  }
  property {
    name  = "checksListAnnotations_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_annotations_owner
  }
  property {
    name  = "checksListAnnotations_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_annotations_page
  }
  property {
    name  = "checksListAnnotations_per_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_annotations_per_page
  }
  property {
    name  = "checksListAnnotations_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_annotations_repo
  }
  property {
    name  = "checksListForRef_app_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_app_id
  }
  property {
    name  = "checksListForRef_check_name"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_check_name
  }
  property {
    name  = "checksListForRef_filter"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_filter
  }
  property {
    name  = "checksListForRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_owner
  }
  property {
    name  = "checksListForRef_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_page
  }
  property {
    name  = "checksListForRef_per_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_per_page
  }
  property {
    name  = "checksListForRef_ref"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_ref
  }
  property {
    name  = "checksListForRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_repo
  }
  property {
    name  = "checksListForRef_status"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_ref_status
  }
  property {
    name  = "checksListForSuite_check_name"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_check_name
  }
  property {
    name  = "checksListForSuite_check_suite_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_check_suite_id
  }
  property {
    name  = "checksListForSuite_filter"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_filter
  }
  property {
    name  = "checksListForSuite_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_owner
  }
  property {
    name  = "checksListForSuite_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_page
  }
  property {
    name  = "checksListForSuite_per_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_per_page
  }
  property {
    name  = "checksListForSuite_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_repo
  }
  property {
    name  = "checksListForSuite_status"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_for_suite_status
  }
  property {
    name  = "checksListSuitesForRef_app_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_suites_for_ref_app_id
  }
  property {
    name  = "checksListSuitesForRef_check_name"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_suites_for_ref_check_name
  }
  property {
    name  = "checksListSuitesForRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_suites_for_ref_owner
  }
  property {
    name  = "checksListSuitesForRef_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_suites_for_ref_page
  }
  property {
    name  = "checksListSuitesForRef_per_page"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_suites_for_ref_per_page
  }
  property {
    name  = "checksListSuitesForRef_ref"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_suites_for_ref_ref
  }
  property {
    name  = "checksListSuitesForRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_list_suites_for_ref_repo
  }
  property {
    name  = "checksRerequestSuite_check_suite_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_rerequest_suite_check_suite_id
  }
  property {
    name  = "checksRerequestSuite_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_rerequest_suite_owner
  }
  property {
    name  = "checksRerequestSuite_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_rerequest_suite_repo
  }
  property {
    name  = "checksSetSuitesPreferences_checksSetSuitesPreferencesRequest_ChecksSetSuitesPreferencesRequest_auto_trigger_checks"
    type  = "string"
    value = var.connector-oai-github_property_checks_set_suites_preferences_checks_set_suites_preferences_request_checks_set_suites_preferences_request_auto_trigger_checks
  }
  property {
    name  = "checksSetSuitesPreferences_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_set_suites_preferences_owner
  }
  property {
    name  = "checksSetSuitesPreferences_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_set_suites_preferences_repo
  }
  property {
    name  = "checksUpdate_check_run_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_check_run_id
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequestOutput_annotations"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_output_annotations
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequestOutput_images"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_output_images
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequestOutput_summary"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_output_summary
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequestOutput_text"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_output_text
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequestOutput_title"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_output_title
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_actions"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_actions
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_completed_at"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_completed_at
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_conclusion"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_conclusion
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_details_url"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_details_url
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_external_id"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_external_id
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_name
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_started_at"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_started_at
  }
  property {
    name  = "checksUpdate_checksUpdateRequest_ChecksUpdateRequest_status"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_checks_update_request_checks_update_request_status
  }
  property {
    name  = "checksUpdate_owner"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_owner
  }
  property {
    name  = "checksUpdate_repo"
    type  = "string"
    value = var.connector-oai-github_property_checks_update_repo
  }
  property {
    name  = "codeScanningDeleteAnalysis_analysis_id"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_delete_analysis_analysis_id
  }
  property {
    name  = "codeScanningDeleteAnalysis_confirm_delete"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_delete_analysis_confirm_delete
  }
  property {
    name  = "codeScanningDeleteAnalysis_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_delete_analysis_owner
  }
  property {
    name  = "codeScanningDeleteAnalysis_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_delete_analysis_repo
  }
  property {
    name  = "codeScanningGetAlert_alert_number"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_alert_alert_number
  }
  property {
    name  = "codeScanningGetAlert_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_alert_owner
  }
  property {
    name  = "codeScanningGetAlert_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_alert_repo
  }
  property {
    name  = "codeScanningGetAnalysis_analysis_id"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_analysis_analysis_id
  }
  property {
    name  = "codeScanningGetAnalysis_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_analysis_owner
  }
  property {
    name  = "codeScanningGetAnalysis_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_analysis_repo
  }
  property {
    name  = "codeScanningGetSarif_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_sarif_owner
  }
  property {
    name  = "codeScanningGetSarif_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_sarif_repo
  }
  property {
    name  = "codeScanningGetSarif_sarif_id"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_get_sarif_sarif_id
  }
  property {
    name  = "codeScanningListAlertInstances_alert_number"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alert_instances_alert_number
  }
  property {
    name  = "codeScanningListAlertInstances_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alert_instances_owner
  }
  property {
    name  = "codeScanningListAlertInstances_page"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alert_instances_page
  }
  property {
    name  = "codeScanningListAlertInstances_per_page"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alert_instances_per_page
  }
  property {
    name  = "codeScanningListAlertInstances_ref"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alert_instances_ref
  }
  property {
    name  = "codeScanningListAlertInstances_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alert_instances_repo
  }
  property {
    name  = "codeScanningListAlertsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_owner
  }
  property {
    name  = "codeScanningListAlertsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_page
  }
  property {
    name  = "codeScanningListAlertsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_per_page
  }
  property {
    name  = "codeScanningListAlertsForRepo_ref"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_ref
  }
  property {
    name  = "codeScanningListAlertsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_repo
  }
  property {
    name  = "codeScanningListAlertsForRepo_state"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_state
  }
  property {
    name  = "codeScanningListAlertsForRepo_tool_guid"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_tool_guid
  }
  property {
    name  = "codeScanningListAlertsForRepo_tool_name"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_alerts_for_repo_tool_name
  }
  property {
    name  = "codeScanningListRecentAnalyses_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_owner
  }
  property {
    name  = "codeScanningListRecentAnalyses_page"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_page
  }
  property {
    name  = "codeScanningListRecentAnalyses_per_page"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_per_page
  }
  property {
    name  = "codeScanningListRecentAnalyses_ref"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_ref
  }
  property {
    name  = "codeScanningListRecentAnalyses_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_repo
  }
  property {
    name  = "codeScanningListRecentAnalyses_sarif_id"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_sarif_id
  }
  property {
    name  = "codeScanningListRecentAnalyses_tool_guid"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_tool_guid
  }
  property {
    name  = "codeScanningListRecentAnalyses_tool_name"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_list_recent_analyses_tool_name
  }
  property {
    name  = "codeScanningUpdateAlert_alert_number"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_update_alert_alert_number
  }
  property {
    name  = "codeScanningUpdateAlert_codeScanningUpdateAlertRequest_CodeScanningUpdateAlertRequest_dismissed_reason"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_update_alert_code_scanning_update_alert_request_code_scanning_update_alert_request_dismissed_reason
  }
  property {
    name  = "codeScanningUpdateAlert_codeScanningUpdateAlertRequest_CodeScanningUpdateAlertRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_update_alert_code_scanning_update_alert_request_code_scanning_update_alert_request_state
  }
  property {
    name  = "codeScanningUpdateAlert_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_update_alert_owner
  }
  property {
    name  = "codeScanningUpdateAlert_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_update_alert_repo
  }
  property {
    name  = "codeScanningUploadSarif_codeScanningUploadSarifRequest_CodeScanningUploadSarifRequest_checkout_uri"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_code_scanning_upload_sarif_request_code_scanning_upload_sarif_request_checkout_uri
  }
  property {
    name  = "codeScanningUploadSarif_codeScanningUploadSarifRequest_CodeScanningUploadSarifRequest_commit_sha"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_code_scanning_upload_sarif_request_code_scanning_upload_sarif_request_commit_sha
  }
  property {
    name  = "codeScanningUploadSarif_codeScanningUploadSarifRequest_CodeScanningUploadSarifRequest_ref"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_code_scanning_upload_sarif_request_code_scanning_upload_sarif_request_ref
  }
  property {
    name  = "codeScanningUploadSarif_codeScanningUploadSarifRequest_CodeScanningUploadSarifRequest_sarif"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_code_scanning_upload_sarif_request_code_scanning_upload_sarif_request_sarif
  }
  property {
    name  = "codeScanningUploadSarif_codeScanningUploadSarifRequest_CodeScanningUploadSarifRequest_started_at"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_code_scanning_upload_sarif_request_code_scanning_upload_sarif_request_started_at
  }
  property {
    name  = "codeScanningUploadSarif_codeScanningUploadSarifRequest_CodeScanningUploadSarifRequest_tool_name"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_code_scanning_upload_sarif_request_code_scanning_upload_sarif_request_tool_name
  }
  property {
    name  = "codeScanningUploadSarif_owner"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_owner
  }
  property {
    name  = "codeScanningUploadSarif_repo"
    type  = "string"
    value = var.connector-oai-github_property_code_scanning_upload_sarif_repo
  }
  property {
    name  = "codesOfConductGetConductCode_key"
    type  = "string"
    value = var.connector-oai-github_property_codes_of_conduct_get_conduct_code_key
  }
  property {
    name  = "enterpriseAdminAddAuthorizedSshKey_authorized_key"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_add_authorized_ssh_key_authorized_key
  }
  property {
    name  = "enterpriseAdminAddOrgAccessToSelfHostedRunnerGroupInEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_add_org_access_to_self_hosted_runner_group_in_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminAddOrgAccessToSelfHostedRunnerGroupInEnterprise_org_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_add_org_access_to_self_hosted_runner_group_in_enterprise_org_id
  }
  property {
    name  = "enterpriseAdminAddOrgAccessToSelfHostedRunnerGroupInEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_add_org_access_to_self_hosted_runner_group_in_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminAddSelfHostedRunnerToGroupForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_add_self_hosted_runner_to_group_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminAddSelfHostedRunnerToGroupForEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_add_self_hosted_runner_to_group_for_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminAddSelfHostedRunnerToGroupForEnterprise_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_add_self_hosted_runner_to_group_for_enterprise_runner_id
  }
  property {
    name  = "enterpriseAdminCreateEnterpriseServerLicense_license"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_enterprise_server_license_license
  }
  property {
    name  = "enterpriseAdminCreateEnterpriseServerLicense_password"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_enterprise_server_license_password
  }
  property {
    name  = "enterpriseAdminCreateEnterpriseServerLicense_settings"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_enterprise_server_license_settings
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_accept"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_accept
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_enterpriseAdminCreateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_content_type"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_enterprise_admin_create_global_webhook_request_enterprise_admin_create_global_webhook_request_config_content_type
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_enterpriseAdminCreateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_insecure_ssl"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_enterprise_admin_create_global_webhook_request_enterprise_admin_create_global_webhook_request_config_insecure_ssl
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_enterpriseAdminCreateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_secret"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_enterprise_admin_create_global_webhook_request_enterprise_admin_create_global_webhook_request_config_secret
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_enterpriseAdminCreateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_url"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_enterprise_admin_create_global_webhook_request_enterprise_admin_create_global_webhook_request_config_url
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_enterpriseAdminCreateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequest_active"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_enterprise_admin_create_global_webhook_request_enterprise_admin_create_global_webhook_request_active
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_enterpriseAdminCreateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequest_events"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_enterprise_admin_create_global_webhook_request_enterprise_admin_create_global_webhook_request_events
  }
  property {
    name  = "enterpriseAdminCreateGlobalWebhook_enterpriseAdminCreateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_global_webhook_enterprise_admin_create_global_webhook_request_enterprise_admin_create_global_webhook_request_name
  }
  property {
    name  = "enterpriseAdminCreateImpersonationOAuthToken_enterpriseAdminCreateImpersonationOAuthTokenRequest_EnterpriseAdminCreateImpersonationOAuthTokenRequest_scopes"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_impersonation_oauth_token_enterprise_admin_create_impersonation_oauth_token_request_enterprise_admin_create_impersonation_oauth_token_request_scopes
  }
  property {
    name  = "enterpriseAdminCreateImpersonationOAuthToken_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_impersonation_oauth_token_username
  }
  property {
    name  = "enterpriseAdminCreateOrg_enterpriseAdminCreateOrgRequest_EnterpriseAdminCreateOrgRequest_admin"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_org_enterprise_admin_create_org_request_enterprise_admin_create_org_request_admin
  }
  property {
    name  = "enterpriseAdminCreateOrg_enterpriseAdminCreateOrgRequest_EnterpriseAdminCreateOrgRequest_login"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_org_enterprise_admin_create_org_request_enterprise_admin_create_org_request_login
  }
  property {
    name  = "enterpriseAdminCreateOrg_enterpriseAdminCreateOrgRequest_EnterpriseAdminCreateOrgRequest_profile_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_org_enterprise_admin_create_org_request_enterprise_admin_create_org_request_profile_name
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveEnvironment_enterpriseAdminCreatePreReceiveEnvironmentRequest_EnterpriseAdminCreatePreReceiveEnvironmentRequest_image_url"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_environment_enterprise_admin_create_pre_receive_environment_request_enterprise_admin_create_pre_receive_environment_request_image_url
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveEnvironment_enterpriseAdminCreatePreReceiveEnvironmentRequest_EnterpriseAdminCreatePreReceiveEnvironmentRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_environment_enterprise_admin_create_pre_receive_environment_request_enterprise_admin_create_pre_receive_environment_request_name
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveHook_enterpriseAdminCreatePreReceiveHookRequest_EnterpriseAdminCreatePreReceiveHookRequest_allow_downstream_configuration"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_hook_enterprise_admin_create_pre_receive_hook_request_enterprise_admin_create_pre_receive_hook_request_allow_downstream_configuration
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveHook_enterpriseAdminCreatePreReceiveHookRequest_EnterpriseAdminCreatePreReceiveHookRequest_enforcement"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_hook_enterprise_admin_create_pre_receive_hook_request_enterprise_admin_create_pre_receive_hook_request_enforcement
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveHook_enterpriseAdminCreatePreReceiveHookRequest_EnterpriseAdminCreatePreReceiveHookRequest_environment"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_hook_enterprise_admin_create_pre_receive_hook_request_enterprise_admin_create_pre_receive_hook_request_environment
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveHook_enterpriseAdminCreatePreReceiveHookRequest_EnterpriseAdminCreatePreReceiveHookRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_hook_enterprise_admin_create_pre_receive_hook_request_enterprise_admin_create_pre_receive_hook_request_name
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveHook_enterpriseAdminCreatePreReceiveHookRequest_EnterpriseAdminCreatePreReceiveHookRequest_script"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_hook_enterprise_admin_create_pre_receive_hook_request_enterprise_admin_create_pre_receive_hook_request_script
  }
  property {
    name  = "enterpriseAdminCreatePreReceiveHook_enterpriseAdminCreatePreReceiveHookRequest_EnterpriseAdminCreatePreReceiveHookRequest_script_repository"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_pre_receive_hook_enterprise_admin_create_pre_receive_hook_request_enterprise_admin_create_pre_receive_hook_request_script_repository
  }
  property {
    name  = "enterpriseAdminCreateRegistrationTokenForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_registration_token_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminCreateRemoveTokenForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_remove_token_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminCreateSelfHostedRunnerGroupForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_self_hosted_runner_group_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminCreateSelfHostedRunnerGroupForEnterprise_enterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_allows_public_repositories"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_self_hosted_runner_group_for_enterprise_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_allows_public_repositories
  }
  property {
    name  = "enterpriseAdminCreateSelfHostedRunnerGroupForEnterprise_enterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_self_hosted_runner_group_for_enterprise_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_name
  }
  property {
    name  = "enterpriseAdminCreateSelfHostedRunnerGroupForEnterprise_enterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_runners"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_self_hosted_runner_group_for_enterprise_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_runners
  }
  property {
    name  = "enterpriseAdminCreateSelfHostedRunnerGroupForEnterprise_enterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_selected_organization_ids"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_self_hosted_runner_group_for_enterprise_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_selected_organization_ids
  }
  property {
    name  = "enterpriseAdminCreateSelfHostedRunnerGroupForEnterprise_enterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminCreateSelfHostedRunnerGroupForEnterpriseRequest_visibility"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_self_hosted_runner_group_for_enterprise_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_enterprise_admin_create_self_hosted_runner_group_for_enterprise_request_visibility
  }
  property {
    name  = "enterpriseAdminCreateUser_enterpriseAdminCreateUserRequest_EnterpriseAdminCreateUserRequest_email"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_user_enterprise_admin_create_user_request_enterprise_admin_create_user_request_email
  }
  property {
    name  = "enterpriseAdminCreateUser_enterpriseAdminCreateUserRequest_EnterpriseAdminCreateUserRequest_login"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_create_user_enterprise_admin_create_user_request_enterprise_admin_create_user_request_login
  }
  property {
    name  = "enterpriseAdminDeleteGlobalWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_global_webhook_hook_id
  }
  property {
    name  = "enterpriseAdminDeleteImpersonationOAuthToken_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_impersonation_oauth_token_username
  }
  property {
    name  = "enterpriseAdminDeletePersonalAccessToken_token_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_personal_access_token_token_id
  }
  property {
    name  = "enterpriseAdminDeletePreReceiveEnvironment_pre_receive_environment_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_pre_receive_environment_pre_receive_environment_id
  }
  property {
    name  = "enterpriseAdminDeletePreReceiveHook_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_pre_receive_hook_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminDeletePublicKey_key_ids"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_public_key_key_ids
  }
  property {
    name  = "enterpriseAdminDeleteSelfHostedRunnerFromEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_self_hosted_runner_from_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminDeleteSelfHostedRunnerFromEnterprise_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_self_hosted_runner_from_enterprise_runner_id
  }
  property {
    name  = "enterpriseAdminDeleteSelfHostedRunnerGroupFromEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_self_hosted_runner_group_from_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminDeleteSelfHostedRunnerGroupFromEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_self_hosted_runner_group_from_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminDeleteUser_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_delete_user_username
  }
  property {
    name  = "enterpriseAdminDemoteSiteAdministrator_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_demote_site_administrator_username
  }
  property {
    name  = "enterpriseAdminDisableSelectedOrganizationGithubActionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_disable_selected_organization_github_actions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminDisableSelectedOrganizationGithubActionsEnterprise_org_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_disable_selected_organization_github_actions_enterprise_org_id
  }
  property {
    name  = "enterpriseAdminEnableOrDisableMaintenanceMode_maintenance"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_enable_or_disable_maintenance_mode_maintenance
  }
  property {
    name  = "enterpriseAdminEnableSelectedOrganizationGithubActionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_enable_selected_organization_github_actions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminEnableSelectedOrganizationGithubActionsEnterprise_org_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_enable_selected_organization_github_actions_enterprise_org_id
  }
  property {
    name  = "enterpriseAdminGetAllowedActionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_allowed_actions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminGetDownloadStatusForPreReceiveEnvironment_pre_receive_environment_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_download_status_for_pre_receive_environment_pre_receive_environment_id
  }
  property {
    name  = "enterpriseAdminGetGithubActionsPermissionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_github_actions_permissions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminGetGlobalWebhook_accept"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_global_webhook_accept
  }
  property {
    name  = "enterpriseAdminGetGlobalWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_global_webhook_hook_id
  }
  property {
    name  = "enterpriseAdminGetPreReceiveEnvironment_pre_receive_environment_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_pre_receive_environment_pre_receive_environment_id
  }
  property {
    name  = "enterpriseAdminGetPreReceiveHookForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_pre_receive_hook_for_org_org
  }
  property {
    name  = "enterpriseAdminGetPreReceiveHookForOrg_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_pre_receive_hook_for_org_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminGetPreReceiveHookForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_pre_receive_hook_for_repo_owner
  }
  property {
    name  = "enterpriseAdminGetPreReceiveHookForRepo_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_pre_receive_hook_for_repo_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminGetPreReceiveHookForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_pre_receive_hook_for_repo_repo
  }
  property {
    name  = "enterpriseAdminGetPreReceiveHook_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_pre_receive_hook_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminGetSelfHostedRunnerForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_self_hosted_runner_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminGetSelfHostedRunnerForEnterprise_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_self_hosted_runner_for_enterprise_runner_id
  }
  property {
    name  = "enterpriseAdminGetSelfHostedRunnerGroupForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_self_hosted_runner_group_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminGetSelfHostedRunnerGroupForEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_get_self_hosted_runner_group_for_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminListGlobalWebhooks_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_global_webhooks_page
  }
  property {
    name  = "enterpriseAdminListGlobalWebhooks_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_global_webhooks_per_page
  }
  property {
    name  = "enterpriseAdminListOrgAccessToSelfHostedRunnerGroupInEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_org_access_to_self_hosted_runner_group_in_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminListOrgAccessToSelfHostedRunnerGroupInEnterprise_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_org_access_to_self_hosted_runner_group_in_enterprise_page
  }
  property {
    name  = "enterpriseAdminListOrgAccessToSelfHostedRunnerGroupInEnterprise_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_org_access_to_self_hosted_runner_group_in_enterprise_per_page
  }
  property {
    name  = "enterpriseAdminListOrgAccessToSelfHostedRunnerGroupInEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_org_access_to_self_hosted_runner_group_in_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminListPersonalAccessTokens_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_personal_access_tokens_page
  }
  property {
    name  = "enterpriseAdminListPersonalAccessTokens_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_personal_access_tokens_per_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveEnvironments_direction"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_environments_direction
  }
  property {
    name  = "enterpriseAdminListPreReceiveEnvironments_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_environments_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveEnvironments_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_environments_per_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveEnvironments_sort"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_environments_sort
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForOrg_direction"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_org_direction
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_org_org
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_org_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_org_per_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForOrg_sort"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_org_sort
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForRepo_direction"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_repo_direction
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_repo_owner
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_repo_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_repo_per_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_repo_repo
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooksForRepo_sort"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_for_repo_sort
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooks_direction"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_direction
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooks_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooks_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_per_page
  }
  property {
    name  = "enterpriseAdminListPreReceiveHooks_sort"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_pre_receive_hooks_sort
  }
  property {
    name  = "enterpriseAdminListPublicKeys_direction"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_public_keys_direction
  }
  property {
    name  = "enterpriseAdminListPublicKeys_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_public_keys_page
  }
  property {
    name  = "enterpriseAdminListPublicKeys_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_public_keys_per_page
  }
  property {
    name  = "enterpriseAdminListPublicKeys_since"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_public_keys_since
  }
  property {
    name  = "enterpriseAdminListPublicKeys_sort"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_public_keys_sort
  }
  property {
    name  = "enterpriseAdminListRunnerApplicationsForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_runner_applications_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminListSelectedOrganizationsEnabledGithubActionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_selected_organizations_enabled_github_actions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminListSelectedOrganizationsEnabledGithubActionsEnterprise_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_selected_organizations_enabled_github_actions_enterprise_page
  }
  property {
    name  = "enterpriseAdminListSelectedOrganizationsEnabledGithubActionsEnterprise_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_selected_organizations_enabled_github_actions_enterprise_per_page
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnerGroupsForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runner_groups_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnerGroupsForEnterprise_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runner_groups_for_enterprise_page
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnerGroupsForEnterprise_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runner_groups_for_enterprise_per_page
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnersForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runners_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnersForEnterprise_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runners_for_enterprise_page
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnersForEnterprise_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runners_for_enterprise_per_page
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnersInGroupForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runners_in_group_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnersInGroupForEnterprise_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runners_in_group_for_enterprise_page
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnersInGroupForEnterprise_per_page"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runners_in_group_for_enterprise_per_page
  }
  property {
    name  = "enterpriseAdminListSelfHostedRunnersInGroupForEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_list_self_hosted_runners_in_group_for_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminPingGlobalWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_ping_global_webhook_hook_id
  }
  property {
    name  = "enterpriseAdminPromoteUserToBeSiteAdministrator_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_promote_user_to_be_site_administrator_username
  }
  property {
    name  = "enterpriseAdminRemoveAuthorizedSshKey_authorized_key"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_authorized_ssh_key_authorized_key
  }
  property {
    name  = "enterpriseAdminRemoveOrgAccessToSelfHostedRunnerGroupInEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_org_access_to_self_hosted_runner_group_in_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminRemoveOrgAccessToSelfHostedRunnerGroupInEnterprise_org_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_org_access_to_self_hosted_runner_group_in_enterprise_org_id
  }
  property {
    name  = "enterpriseAdminRemoveOrgAccessToSelfHostedRunnerGroupInEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_org_access_to_self_hosted_runner_group_in_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminRemovePreReceiveHookEnforcementForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_pre_receive_hook_enforcement_for_org_org
  }
  property {
    name  = "enterpriseAdminRemovePreReceiveHookEnforcementForOrg_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_pre_receive_hook_enforcement_for_org_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminRemovePreReceiveHookEnforcementForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_pre_receive_hook_enforcement_for_repo_owner
  }
  property {
    name  = "enterpriseAdminRemovePreReceiveHookEnforcementForRepo_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_pre_receive_hook_enforcement_for_repo_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminRemovePreReceiveHookEnforcementForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_pre_receive_hook_enforcement_for_repo_repo
  }
  property {
    name  = "enterpriseAdminRemoveSelfHostedRunnerFromGroupForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_self_hosted_runner_from_group_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminRemoveSelfHostedRunnerFromGroupForEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_self_hosted_runner_from_group_for_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminRemoveSelfHostedRunnerFromGroupForEnterprise_runner_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_remove_self_hosted_runner_from_group_for_enterprise_runner_id
  }
  property {
    name  = "enterpriseAdminSetAllowedActionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_allowed_actions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminSetAllowedActionsEnterprise_selectedActions_SelectedActions_github_owned_allowed"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_allowed_actions_enterprise_selected_actions_selected_actions_github_owned_allowed
  }
  property {
    name  = "enterpriseAdminSetAllowedActionsEnterprise_selectedActions_SelectedActions_patterns_allowed"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_allowed_actions_enterprise_selected_actions_selected_actions_patterns_allowed
  }
  property {
    name  = "enterpriseAdminSetAnnouncement_announcement_Announcement_announcement"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_announcement_announcement_announcement_announcement
  }
  property {
    name  = "enterpriseAdminSetAnnouncement_announcement_Announcement_expires_at"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_announcement_announcement_announcement_expires_at
  }
  property {
    name  = "enterpriseAdminSetGithubActionsPermissionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_github_actions_permissions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminSetGithubActionsPermissionsEnterprise_enterpriseAdminSetGithubActionsPermissionsEnterpriseRequest_EnterpriseAdminSetGithubActionsPermissionsEnterpriseRequest_allowed_actions"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_github_actions_permissions_enterprise_enterprise_admin_set_github_actions_permissions_enterprise_request_enterprise_admin_set_github_actions_permissions_enterprise_request_allowed_actions
  }
  property {
    name  = "enterpriseAdminSetGithubActionsPermissionsEnterprise_enterpriseAdminSetGithubActionsPermissionsEnterpriseRequest_EnterpriseAdminSetGithubActionsPermissionsEnterpriseRequest_enabled_organizations"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_github_actions_permissions_enterprise_enterprise_admin_set_github_actions_permissions_enterprise_request_enterprise_admin_set_github_actions_permissions_enterprise_request_enabled_organizations
  }
  property {
    name  = "enterpriseAdminSetOrgAccessToSelfHostedRunnerGroupInEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_org_access_to_self_hosted_runner_group_in_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminSetOrgAccessToSelfHostedRunnerGroupInEnterprise_enterpriseAdminSetOrgAccessToSelfHostedRunnerGroupInEnterpriseRequest_EnterpriseAdminSetOrgAccessToSelfHostedRunnerGroupInEnterpriseRequest_selected_organization_ids"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_org_access_to_self_hosted_runner_group_in_enterprise_enterprise_admin_set_org_access_to_self_hosted_runner_group_in_enterprise_request_enterprise_admin_set_org_access_to_self_hosted_runner_group_in_enterprise_request_selected_organization_ids
  }
  property {
    name  = "enterpriseAdminSetOrgAccessToSelfHostedRunnerGroupInEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_org_access_to_self_hosted_runner_group_in_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminSetSelectedOrganizationsEnabledGithubActionsEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_selected_organizations_enabled_github_actions_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminSetSelectedOrganizationsEnabledGithubActionsEnterprise_enterpriseAdminSetSelectedOrganizationsEnabledGithubActionsEnterpriseRequest_EnterpriseAdminSetSelectedOrganizationsEnabledGithubActionsEnterpriseRequest_selected_organization_ids"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_selected_organizations_enabled_github_actions_enterprise_enterprise_admin_set_selected_organizations_enabled_github_actions_enterprise_request_enterprise_admin_set_selected_organizations_enabled_github_actions_enterprise_request_selected_organization_ids
  }
  property {
    name  = "enterpriseAdminSetSelfHostedRunnersInGroupForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_self_hosted_runners_in_group_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminSetSelfHostedRunnersInGroupForEnterprise_enterpriseAdminSetSelfHostedRunnersInGroupForEnterpriseRequest_EnterpriseAdminSetSelfHostedRunnersInGroupForEnterpriseRequest_runners"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_self_hosted_runners_in_group_for_enterprise_enterprise_admin_set_self_hosted_runners_in_group_for_enterprise_request_enterprise_admin_set_self_hosted_runners_in_group_for_enterprise_request_runners
  }
  property {
    name  = "enterpriseAdminSetSelfHostedRunnersInGroupForEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_self_hosted_runners_in_group_for_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminSetSettings_settings"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_set_settings_settings
  }
  property {
    name  = "enterpriseAdminStartPreReceiveEnvironmentDownload_pre_receive_environment_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_start_pre_receive_environment_download_pre_receive_environment_id
  }
  property {
    name  = "enterpriseAdminSuspendUser_enterpriseAdminSuspendUserRequest_EnterpriseAdminSuspendUserRequest_reason"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_suspend_user_enterprise_admin_suspend_user_request_enterprise_admin_suspend_user_request_reason
  }
  property {
    name  = "enterpriseAdminSuspendUser_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_suspend_user_username
  }
  property {
    name  = "enterpriseAdminSyncLdapMappingForTeam_team_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_sync_ldap_mapping_for_team_team_id
  }
  property {
    name  = "enterpriseAdminSyncLdapMappingForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_sync_ldap_mapping_for_user_username
  }
  property {
    name  = "enterpriseAdminUnsuspendUser_enterpriseAdminUnsuspendUserRequest_EnterpriseAdminUnsuspendUserRequest_reason"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_unsuspend_user_enterprise_admin_unsuspend_user_request_enterprise_admin_unsuspend_user_request_reason
  }
  property {
    name  = "enterpriseAdminUnsuspendUser_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_unsuspend_user_username
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_accept"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_accept
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_enterpriseAdminUpdateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_content_type"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_enterprise_admin_update_global_webhook_request_enterprise_admin_create_global_webhook_request_config_content_type
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_enterpriseAdminUpdateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_insecure_ssl"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_enterprise_admin_update_global_webhook_request_enterprise_admin_create_global_webhook_request_config_insecure_ssl
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_enterpriseAdminUpdateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_secret"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_enterprise_admin_update_global_webhook_request_enterprise_admin_create_global_webhook_request_config_secret
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_enterpriseAdminUpdateGlobalWebhookRequest_EnterpriseAdminCreateGlobalWebhookRequestConfig_url"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_enterprise_admin_update_global_webhook_request_enterprise_admin_create_global_webhook_request_config_url
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_enterpriseAdminUpdateGlobalWebhookRequest_EnterpriseAdminUpdateGlobalWebhookRequest_active"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_enterprise_admin_update_global_webhook_request_enterprise_admin_update_global_webhook_request_active
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_enterpriseAdminUpdateGlobalWebhookRequest_EnterpriseAdminUpdateGlobalWebhookRequest_events"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_enterprise_admin_update_global_webhook_request_enterprise_admin_update_global_webhook_request_events
  }
  property {
    name  = "enterpriseAdminUpdateGlobalWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_global_webhook_hook_id
  }
  property {
    name  = "enterpriseAdminUpdateLdapMappingForTeam_enterpriseAdminUpdateLdapMappingForTeamRequest_EnterpriseAdminUpdateLdapMappingForTeamRequest_ldap_dn"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_ldap_mapping_for_team_enterprise_admin_update_ldap_mapping_for_team_request_enterprise_admin_update_ldap_mapping_for_team_request_ldap_dn
  }
  property {
    name  = "enterpriseAdminUpdateLdapMappingForTeam_team_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_ldap_mapping_for_team_team_id
  }
  property {
    name  = "enterpriseAdminUpdateLdapMappingForUser_enterpriseAdminUpdateLdapMappingForTeamRequest_EnterpriseAdminUpdateLdapMappingForTeamRequest_ldap_dn"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_ldap_mapping_for_user_enterprise_admin_update_ldap_mapping_for_team_request_enterprise_admin_update_ldap_mapping_for_team_request_ldap_dn
  }
  property {
    name  = "enterpriseAdminUpdateLdapMappingForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_ldap_mapping_for_user_username
  }
  property {
    name  = "enterpriseAdminUpdateOrgName_enterpriseAdminUpdateOrgNameRequest_EnterpriseAdminUpdateOrgNameRequest_login"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_org_name_enterprise_admin_update_org_name_request_enterprise_admin_update_org_name_request_login
  }
  property {
    name  = "enterpriseAdminUpdateOrgName_org"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_org_name_org
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveEnvironment_enterpriseAdminUpdatePreReceiveEnvironmentRequest_EnterpriseAdminUpdatePreReceiveEnvironmentRequest_image_url"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_environment_enterprise_admin_update_pre_receive_environment_request_enterprise_admin_update_pre_receive_environment_request_image_url
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveEnvironment_enterpriseAdminUpdatePreReceiveEnvironmentRequest_EnterpriseAdminUpdatePreReceiveEnvironmentRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_environment_enterprise_admin_update_pre_receive_environment_request_enterprise_admin_update_pre_receive_environment_request_name
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveEnvironment_pre_receive_environment_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_environment_pre_receive_environment_id
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForOrg_enterpriseAdminUpdatePreReceiveHookEnforcementForOrgRequest_EnterpriseAdminUpdatePreReceiveHookEnforcementForOrgRequest_allow_downstream_configuration"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_org_enterprise_admin_update_pre_receive_hook_enforcement_for_org_request_enterprise_admin_update_pre_receive_hook_enforcement_for_org_request_allow_downstream_configuration
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForOrg_enterpriseAdminUpdatePreReceiveHookEnforcementForOrgRequest_EnterpriseAdminUpdatePreReceiveHookEnforcementForOrgRequest_enforcement"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_org_enterprise_admin_update_pre_receive_hook_enforcement_for_org_request_enterprise_admin_update_pre_receive_hook_enforcement_for_org_request_enforcement
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_org_org
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForOrg_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_org_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForRepo_enterpriseAdminUpdatePreReceiveHookEnforcementForRepoRequest_EnterpriseAdminUpdatePreReceiveHookEnforcementForRepoRequest_enforcement"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_repo_enterprise_admin_update_pre_receive_hook_enforcement_for_repo_request_enterprise_admin_update_pre_receive_hook_enforcement_for_repo_request_enforcement
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_repo_owner
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForRepo_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_repo_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHookEnforcementForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enforcement_for_repo_repo
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHook_enterpriseAdminUpdatePreReceiveHookRequest_EnterpriseAdminUpdatePreReceiveHookRequest_allow_downstream_configuration"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enterprise_admin_update_pre_receive_hook_request_enterprise_admin_update_pre_receive_hook_request_allow_downstream_configuration
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHook_enterpriseAdminUpdatePreReceiveHookRequest_EnterpriseAdminUpdatePreReceiveHookRequest_enforcement"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enterprise_admin_update_pre_receive_hook_request_enterprise_admin_update_pre_receive_hook_request_enforcement
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHook_enterpriseAdminUpdatePreReceiveHookRequest_EnterpriseAdminUpdatePreReceiveHookRequest_environment"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enterprise_admin_update_pre_receive_hook_request_enterprise_admin_update_pre_receive_hook_request_environment
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHook_enterpriseAdminUpdatePreReceiveHookRequest_EnterpriseAdminUpdatePreReceiveHookRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enterprise_admin_update_pre_receive_hook_request_enterprise_admin_update_pre_receive_hook_request_name
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHook_enterpriseAdminUpdatePreReceiveHookRequest_EnterpriseAdminUpdatePreReceiveHookRequest_script"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enterprise_admin_update_pre_receive_hook_request_enterprise_admin_update_pre_receive_hook_request_script
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHook_enterpriseAdminUpdatePreReceiveHookRequest_EnterpriseAdminUpdatePreReceiveHookRequest_script_repository"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_enterprise_admin_update_pre_receive_hook_request_enterprise_admin_update_pre_receive_hook_request_script_repository
  }
  property {
    name  = "enterpriseAdminUpdatePreReceiveHook_pre_receive_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_pre_receive_hook_pre_receive_hook_id
  }
  property {
    name  = "enterpriseAdminUpdateSelfHostedRunnerGroupForEnterprise_enterprise"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_self_hosted_runner_group_for_enterprise_enterprise
  }
  property {
    name  = "enterpriseAdminUpdateSelfHostedRunnerGroupForEnterprise_enterpriseAdminUpdateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminUpdateSelfHostedRunnerGroupForEnterpriseRequest_allows_public_repositories"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_self_hosted_runner_group_for_enterprise_enterprise_admin_update_self_hosted_runner_group_for_enterprise_request_enterprise_admin_update_self_hosted_runner_group_for_enterprise_request_allows_public_repositories
  }
  property {
    name  = "enterpriseAdminUpdateSelfHostedRunnerGroupForEnterprise_enterpriseAdminUpdateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminUpdateSelfHostedRunnerGroupForEnterpriseRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_self_hosted_runner_group_for_enterprise_enterprise_admin_update_self_hosted_runner_group_for_enterprise_request_enterprise_admin_update_self_hosted_runner_group_for_enterprise_request_name
  }
  property {
    name  = "enterpriseAdminUpdateSelfHostedRunnerGroupForEnterprise_enterpriseAdminUpdateSelfHostedRunnerGroupForEnterpriseRequest_EnterpriseAdminUpdateSelfHostedRunnerGroupForEnterpriseRequest_visibility"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_self_hosted_runner_group_for_enterprise_enterprise_admin_update_self_hosted_runner_group_for_enterprise_request_enterprise_admin_update_self_hosted_runner_group_for_enterprise_request_visibility
  }
  property {
    name  = "enterpriseAdminUpdateSelfHostedRunnerGroupForEnterprise_runner_group_id"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_self_hosted_runner_group_for_enterprise_runner_group_id
  }
  property {
    name  = "enterpriseAdminUpdateUsernameForUser_enterpriseAdminUpdateUsernameForUserRequest_EnterpriseAdminUpdateUsernameForUserRequest_login"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_username_for_user_enterprise_admin_update_username_for_user_request_enterprise_admin_update_username_for_user_request_login
  }
  property {
    name  = "enterpriseAdminUpdateUsernameForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_update_username_for_user_username
  }
  property {
    name  = "enterpriseAdminUpgradeLicense_license"
    type  = "string"
    value = var.connector-oai-github_property_enterprise_admin_upgrade_license_license
  }
  property {
    name  = "gistsCheckIsStarred_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_check_is_starred_gist_id
  }
  property {
    name  = "gistsCreateComment_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_create_comment_gist_id
  }
  property {
    name  = "gistsCreateComment_gistsCreateCommentRequest_GistsCreateCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_gists_create_comment_gists_create_comment_request_gists_create_comment_request_body
  }
  property {
    name  = "gistsCreate_gistsCreateRequest_GistsCreateRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_gists_create_gists_create_request_gists_create_request_description
  }
  property {
    name  = "gistsCreate_gistsCreateRequest_GistsCreateRequest_files"
    type  = "string"
    value = var.connector-oai-github_property_gists_create_gists_create_request_gists_create_request_files
  }
  property {
    name  = "gistsDeleteComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_delete_comment_comment_id
  }
  property {
    name  = "gistsDeleteComment_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_delete_comment_gist_id
  }
  property {
    name  = "gistsDelete_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_delete_gist_id
  }
  property {
    name  = "gistsFork_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_fork_gist_id
  }
  property {
    name  = "gistsGetComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_get_comment_comment_id
  }
  property {
    name  = "gistsGetComment_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_get_comment_gist_id
  }
  property {
    name  = "gistsGetRevision_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_get_revision_gist_id
  }
  property {
    name  = "gistsGetRevision_sha"
    type  = "string"
    value = var.connector-oai-github_property_gists_get_revision_sha
  }
  property {
    name  = "gistsGet_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_get_gist_id
  }
  property {
    name  = "gistsListComments_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_comments_gist_id
  }
  property {
    name  = "gistsListComments_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_comments_page
  }
  property {
    name  = "gistsListComments_per_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_comments_per_page
  }
  property {
    name  = "gistsListCommits_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_commits_gist_id
  }
  property {
    name  = "gistsListCommits_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_commits_page
  }
  property {
    name  = "gistsListCommits_per_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_commits_per_page
  }
  property {
    name  = "gistsListForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_for_user_page
  }
  property {
    name  = "gistsListForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_for_user_per_page
  }
  property {
    name  = "gistsListForUser_since"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_for_user_since
  }
  property {
    name  = "gistsListForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_for_user_username
  }
  property {
    name  = "gistsListForks_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_forks_gist_id
  }
  property {
    name  = "gistsListForks_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_forks_page
  }
  property {
    name  = "gistsListForks_per_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_forks_per_page
  }
  property {
    name  = "gistsListPublic_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_public_page
  }
  property {
    name  = "gistsListPublic_per_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_public_per_page
  }
  property {
    name  = "gistsListPublic_since"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_public_since
  }
  property {
    name  = "gistsListStarred_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_starred_page
  }
  property {
    name  = "gistsListStarred_per_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_starred_per_page
  }
  property {
    name  = "gistsListStarred_since"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_starred_since
  }
  property {
    name  = "gistsList_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_page
  }
  property {
    name  = "gistsList_per_page"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_per_page
  }
  property {
    name  = "gistsList_since"
    type  = "string"
    value = var.connector-oai-github_property_gists_list_since
  }
  property {
    name  = "gistsStar_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_star_gist_id
  }
  property {
    name  = "gistsUnstar_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_unstar_gist_id
  }
  property {
    name  = "gistsUpdateComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_update_comment_comment_id
  }
  property {
    name  = "gistsUpdateComment_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_update_comment_gist_id
  }
  property {
    name  = "gistsUpdateComment_gistsCreateCommentRequest_GistsCreateCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_gists_update_comment_gists_create_comment_request_gists_create_comment_request_body
  }
  property {
    name  = "gistsUpdate_gist_id"
    type  = "string"
    value = var.connector-oai-github_property_gists_update_gist_id
  }
  property {
    name  = "gistsUpdate_gistsUpdateRequest_GistsUpdateRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_gists_update_gists_update_request_gists_update_request_description
  }
  property {
    name  = "gistsUpdate_gistsUpdateRequest_GistsUpdateRequest_files"
    type  = "string"
    value = var.connector-oai-github_property_gists_update_gists_update_request_gists_update_request_files
  }
  property {
    name  = "gitCreateBlob_gitCreateBlobRequest_GitCreateBlobRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_git_create_blob_git_create_blob_request_git_create_blob_request_content
  }
  property {
    name  = "gitCreateBlob_gitCreateBlobRequest_GitCreateBlobRequest_encoding"
    type  = "string"
    value = var.connector-oai-github_property_git_create_blob_git_create_blob_request_git_create_blob_request_encoding
  }
  property {
    name  = "gitCreateBlob_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_create_blob_owner
  }
  property {
    name  = "gitCreateBlob_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_create_blob_repo
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequestAuthor_date"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_author_date
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequestAuthor_email"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_author_email
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequestAuthor_name"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_author_name
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequestCommitter_date"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_committer_date
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequestCommitter_email"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_committer_email
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequestCommitter_name"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_committer_name
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequest_message"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_message
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequest_parents"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_parents
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequest_signature"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_signature
  }
  property {
    name  = "gitCreateCommit_gitCreateCommitRequest_GitCreateCommitRequest_tree"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_git_create_commit_request_git_create_commit_request_tree
  }
  property {
    name  = "gitCreateCommit_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_owner
  }
  property {
    name  = "gitCreateCommit_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_create_commit_repo
  }
  property {
    name  = "gitCreateRef_gitCreateRefRequest_GitCreateRefRequest_key"
    type  = "string"
    value = var.connector-oai-github_property_git_create_ref_git_create_ref_request_git_create_ref_request_key
  }
  property {
    name  = "gitCreateRef_gitCreateRefRequest_GitCreateRefRequest_ref"
    type  = "string"
    value = var.connector-oai-github_property_git_create_ref_git_create_ref_request_git_create_ref_request_ref
  }
  property {
    name  = "gitCreateRef_gitCreateRefRequest_GitCreateRefRequest_sha"
    type  = "string"
    value = var.connector-oai-github_property_git_create_ref_git_create_ref_request_git_create_ref_request_sha
  }
  property {
    name  = "gitCreateRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_create_ref_owner
  }
  property {
    name  = "gitCreateRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_create_ref_repo
  }
  property {
    name  = "gitCreateTag_gitCreateTagRequest_GitCreateTagRequestTagger_date"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_git_create_tag_request_git_create_tag_request_tagger_date
  }
  property {
    name  = "gitCreateTag_gitCreateTagRequest_GitCreateTagRequestTagger_email"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_git_create_tag_request_git_create_tag_request_tagger_email
  }
  property {
    name  = "gitCreateTag_gitCreateTagRequest_GitCreateTagRequestTagger_name"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_git_create_tag_request_git_create_tag_request_tagger_name
  }
  property {
    name  = "gitCreateTag_gitCreateTagRequest_GitCreateTagRequest_message"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_git_create_tag_request_git_create_tag_request_message
  }
  property {
    name  = "gitCreateTag_gitCreateTagRequest_GitCreateTagRequest_object"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_git_create_tag_request_git_create_tag_request_object
  }
  property {
    name  = "gitCreateTag_gitCreateTagRequest_GitCreateTagRequest_tag"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_git_create_tag_request_git_create_tag_request_tag
  }
  property {
    name  = "gitCreateTag_gitCreateTagRequest_GitCreateTagRequest_type"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_git_create_tag_request_git_create_tag_request_type
  }
  property {
    name  = "gitCreateTag_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_owner
  }
  property {
    name  = "gitCreateTag_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tag_repo
  }
  property {
    name  = "gitCreateTree_gitCreateTreeRequest_GitCreateTreeRequest_base_tree"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tree_git_create_tree_request_git_create_tree_request_base_tree
  }
  property {
    name  = "gitCreateTree_gitCreateTreeRequest_GitCreateTreeRequest_tree"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tree_git_create_tree_request_git_create_tree_request_tree
  }
  property {
    name  = "gitCreateTree_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tree_owner
  }
  property {
    name  = "gitCreateTree_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_create_tree_repo
  }
  property {
    name  = "gitDeleteRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_delete_ref_owner
  }
  property {
    name  = "gitDeleteRef_ref"
    type  = "string"
    value = var.connector-oai-github_property_git_delete_ref_ref
  }
  property {
    name  = "gitDeleteRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_delete_ref_repo
  }
  property {
    name  = "gitGetBlob_file_sha"
    type  = "string"
    value = var.connector-oai-github_property_git_get_blob_file_sha
  }
  property {
    name  = "gitGetBlob_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_get_blob_owner
  }
  property {
    name  = "gitGetBlob_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_get_blob_repo
  }
  property {
    name  = "gitGetCommit_commit_sha"
    type  = "string"
    value = var.connector-oai-github_property_git_get_commit_commit_sha
  }
  property {
    name  = "gitGetCommit_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_get_commit_owner
  }
  property {
    name  = "gitGetCommit_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_get_commit_repo
  }
  property {
    name  = "gitGetRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_get_ref_owner
  }
  property {
    name  = "gitGetRef_ref"
    type  = "string"
    value = var.connector-oai-github_property_git_get_ref_ref
  }
  property {
    name  = "gitGetRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_get_ref_repo
  }
  property {
    name  = "gitGetTag_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_get_tag_owner
  }
  property {
    name  = "gitGetTag_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_get_tag_repo
  }
  property {
    name  = "gitGetTag_tag_sha"
    type  = "string"
    value = var.connector-oai-github_property_git_get_tag_tag_sha
  }
  property {
    name  = "gitGetTree_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_get_tree_owner
  }
  property {
    name  = "gitGetTree_recursive"
    type  = "string"
    value = var.connector-oai-github_property_git_get_tree_recursive
  }
  property {
    name  = "gitGetTree_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_get_tree_repo
  }
  property {
    name  = "gitGetTree_tree_sha"
    type  = "string"
    value = var.connector-oai-github_property_git_get_tree_tree_sha
  }
  property {
    name  = "gitListMatchingRefs_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_list_matching_refs_owner
  }
  property {
    name  = "gitListMatchingRefs_page"
    type  = "string"
    value = var.connector-oai-github_property_git_list_matching_refs_page
  }
  property {
    name  = "gitListMatchingRefs_per_page"
    type  = "string"
    value = var.connector-oai-github_property_git_list_matching_refs_per_page
  }
  property {
    name  = "gitListMatchingRefs_ref"
    type  = "string"
    value = var.connector-oai-github_property_git_list_matching_refs_ref
  }
  property {
    name  = "gitListMatchingRefs_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_list_matching_refs_repo
  }
  property {
    name  = "gitUpdateRef_gitUpdateRefRequest_GitUpdateRefRequest_force"
    type  = "string"
    value = var.connector-oai-github_property_git_update_ref_git_update_ref_request_git_update_ref_request_force
  }
  property {
    name  = "gitUpdateRef_gitUpdateRefRequest_GitUpdateRefRequest_sha"
    type  = "string"
    value = var.connector-oai-github_property_git_update_ref_git_update_ref_request_git_update_ref_request_sha
  }
  property {
    name  = "gitUpdateRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_git_update_ref_owner
  }
  property {
    name  = "gitUpdateRef_ref"
    type  = "string"
    value = var.connector-oai-github_property_git_update_ref_ref
  }
  property {
    name  = "gitUpdateRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_git_update_ref_repo
  }
  property {
    name  = "gitignoreGetTemplate_name"
    type  = "string"
    value = var.connector-oai-github_property_gitignore_get_template_name
  }
  property {
    name  = "issuesAddAssignees_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_assignees_issue_number
  }
  property {
    name  = "issuesAddAssignees_issuesAddAssigneesRequest_IssuesAddAssigneesRequest_assignees"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_assignees_issues_add_assignees_request_issues_add_assignees_request_assignees
  }
  property {
    name  = "issuesAddAssignees_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_assignees_owner
  }
  property {
    name  = "issuesAddAssignees_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_assignees_repo
  }
  property {
    name  = "issuesAddLabels_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_labels_issue_number
  }
  property {
    name  = "issuesAddLabels_issuesAddLabelsRequest_IssuesAddLabelsRequest_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_labels_issues_add_labels_request_issues_add_labels_request_labels
  }
  property {
    name  = "issuesAddLabels_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_labels_owner
  }
  property {
    name  = "issuesAddLabels_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_add_labels_repo
  }
  property {
    name  = "issuesCheckUserCanBeAssigned_assignee"
    type  = "string"
    value = var.connector-oai-github_property_issues_check_user_can_be_assigned_assignee
  }
  property {
    name  = "issuesCheckUserCanBeAssigned_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_check_user_can_be_assigned_owner
  }
  property {
    name  = "issuesCheckUserCanBeAssigned_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_check_user_can_be_assigned_repo
  }
  property {
    name  = "issuesCreateComment_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_comment_issue_number
  }
  property {
    name  = "issuesCreateComment_issuesUpdateCommentRequest_IssuesUpdateCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_comment_issues_update_comment_request_issues_update_comment_request_body
  }
  property {
    name  = "issuesCreateComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_comment_owner
  }
  property {
    name  = "issuesCreateComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_comment_repo
  }
  property {
    name  = "issuesCreateLabel_issuesCreateLabelRequest_IssuesCreateLabelRequest_color"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_label_issues_create_label_request_issues_create_label_request_color
  }
  property {
    name  = "issuesCreateLabel_issuesCreateLabelRequest_IssuesCreateLabelRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_label_issues_create_label_request_issues_create_label_request_description
  }
  property {
    name  = "issuesCreateLabel_issuesCreateLabelRequest_IssuesCreateLabelRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_label_issues_create_label_request_issues_create_label_request_name
  }
  property {
    name  = "issuesCreateLabel_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_label_owner
  }
  property {
    name  = "issuesCreateLabel_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_label_repo
  }
  property {
    name  = "issuesCreateMilestone_issuesCreateMilestoneRequest_IssuesCreateMilestoneRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_milestone_issues_create_milestone_request_issues_create_milestone_request_description
  }
  property {
    name  = "issuesCreateMilestone_issuesCreateMilestoneRequest_IssuesCreateMilestoneRequest_due_on"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_milestone_issues_create_milestone_request_issues_create_milestone_request_due_on
  }
  property {
    name  = "issuesCreateMilestone_issuesCreateMilestoneRequest_IssuesCreateMilestoneRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_milestone_issues_create_milestone_request_issues_create_milestone_request_state
  }
  property {
    name  = "issuesCreateMilestone_issuesCreateMilestoneRequest_IssuesCreateMilestoneRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_milestone_issues_create_milestone_request_issues_create_milestone_request_title
  }
  property {
    name  = "issuesCreateMilestone_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_milestone_owner
  }
  property {
    name  = "issuesCreateMilestone_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_milestone_repo
  }
  property {
    name  = "issuesCreate_issuesCreateRequest_IssuesCreateRequest_assignee"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_issues_create_request_issues_create_request_assignee
  }
  property {
    name  = "issuesCreate_issuesCreateRequest_IssuesCreateRequest_assignees"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_issues_create_request_issues_create_request_assignees
  }
  property {
    name  = "issuesCreate_issuesCreateRequest_IssuesCreateRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_issues_create_request_issues_create_request_body
  }
  property {
    name  = "issuesCreate_issuesCreateRequest_IssuesCreateRequest_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_issues_create_request_issues_create_request_labels
  }
  property {
    name  = "issuesCreate_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_owner
  }
  property {
    name  = "issuesCreate_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_create_repo
  }
  property {
    name  = "issuesDeleteComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_comment_comment_id
  }
  property {
    name  = "issuesDeleteComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_comment_owner
  }
  property {
    name  = "issuesDeleteComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_comment_repo
  }
  property {
    name  = "issuesDeleteLabel_name"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_label_name
  }
  property {
    name  = "issuesDeleteLabel_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_label_owner
  }
  property {
    name  = "issuesDeleteLabel_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_label_repo
  }
  property {
    name  = "issuesDeleteMilestone_milestone_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_milestone_milestone_number
  }
  property {
    name  = "issuesDeleteMilestone_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_milestone_owner
  }
  property {
    name  = "issuesDeleteMilestone_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_delete_milestone_repo
  }
  property {
    name  = "issuesGetComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_comment_comment_id
  }
  property {
    name  = "issuesGetComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_comment_owner
  }
  property {
    name  = "issuesGetComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_comment_repo
  }
  property {
    name  = "issuesGetEvent_event_id"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_event_event_id
  }
  property {
    name  = "issuesGetEvent_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_event_owner
  }
  property {
    name  = "issuesGetEvent_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_event_repo
  }
  property {
    name  = "issuesGetLabel_name"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_label_name
  }
  property {
    name  = "issuesGetLabel_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_label_owner
  }
  property {
    name  = "issuesGetLabel_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_label_repo
  }
  property {
    name  = "issuesGetMilestone_milestone_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_milestone_milestone_number
  }
  property {
    name  = "issuesGetMilestone_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_milestone_owner
  }
  property {
    name  = "issuesGetMilestone_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_milestone_repo
  }
  property {
    name  = "issuesGet_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_issue_number
  }
  property {
    name  = "issuesGet_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_owner
  }
  property {
    name  = "issuesGet_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_get_repo
  }
  property {
    name  = "issuesListAssignees_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_assignees_owner
  }
  property {
    name  = "issuesListAssignees_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_assignees_page
  }
  property {
    name  = "issuesListAssignees_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_assignees_per_page
  }
  property {
    name  = "issuesListAssignees_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_assignees_repo
  }
  property {
    name  = "issuesListCommentsForRepo_direction"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_for_repo_direction
  }
  property {
    name  = "issuesListCommentsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_for_repo_owner
  }
  property {
    name  = "issuesListCommentsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_for_repo_page
  }
  property {
    name  = "issuesListCommentsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_for_repo_per_page
  }
  property {
    name  = "issuesListCommentsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_for_repo_repo
  }
  property {
    name  = "issuesListCommentsForRepo_since"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_for_repo_since
  }
  property {
    name  = "issuesListCommentsForRepo_sort"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_for_repo_sort
  }
  property {
    name  = "issuesListComments_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_issue_number
  }
  property {
    name  = "issuesListComments_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_owner
  }
  property {
    name  = "issuesListComments_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_page
  }
  property {
    name  = "issuesListComments_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_per_page
  }
  property {
    name  = "issuesListComments_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_repo
  }
  property {
    name  = "issuesListComments_since"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_comments_since
  }
  property {
    name  = "issuesListEventsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_repo_owner
  }
  property {
    name  = "issuesListEventsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_repo_page
  }
  property {
    name  = "issuesListEventsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_repo_per_page
  }
  property {
    name  = "issuesListEventsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_repo_repo
  }
  property {
    name  = "issuesListEventsForTimeline_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_timeline_issue_number
  }
  property {
    name  = "issuesListEventsForTimeline_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_timeline_owner
  }
  property {
    name  = "issuesListEventsForTimeline_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_timeline_page
  }
  property {
    name  = "issuesListEventsForTimeline_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_timeline_per_page
  }
  property {
    name  = "issuesListEventsForTimeline_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_for_timeline_repo
  }
  property {
    name  = "issuesListEvents_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_issue_number
  }
  property {
    name  = "issuesListEvents_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_owner
  }
  property {
    name  = "issuesListEvents_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_page
  }
  property {
    name  = "issuesListEvents_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_per_page
  }
  property {
    name  = "issuesListEvents_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_events_repo
  }
  property {
    name  = "issuesListForAuthenticatedUser_direction"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_direction
  }
  property {
    name  = "issuesListForAuthenticatedUser_filter"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_filter
  }
  property {
    name  = "issuesListForAuthenticatedUser_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_labels
  }
  property {
    name  = "issuesListForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_page
  }
  property {
    name  = "issuesListForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_per_page
  }
  property {
    name  = "issuesListForAuthenticatedUser_since"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_since
  }
  property {
    name  = "issuesListForAuthenticatedUser_sort"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_sort
  }
  property {
    name  = "issuesListForAuthenticatedUser_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_authenticated_user_state
  }
  property {
    name  = "issuesListForOrg_direction"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_direction
  }
  property {
    name  = "issuesListForOrg_filter"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_filter
  }
  property {
    name  = "issuesListForOrg_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_labels
  }
  property {
    name  = "issuesListForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_org
  }
  property {
    name  = "issuesListForOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_page
  }
  property {
    name  = "issuesListForOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_per_page
  }
  property {
    name  = "issuesListForOrg_since"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_since
  }
  property {
    name  = "issuesListForOrg_sort"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_sort
  }
  property {
    name  = "issuesListForOrg_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_org_state
  }
  property {
    name  = "issuesListForRepo_assignee"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_assignee
  }
  property {
    name  = "issuesListForRepo_creator"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_creator
  }
  property {
    name  = "issuesListForRepo_direction"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_direction
  }
  property {
    name  = "issuesListForRepo_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_labels
  }
  property {
    name  = "issuesListForRepo_mentioned"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_mentioned
  }
  property {
    name  = "issuesListForRepo_milestone"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_milestone
  }
  property {
    name  = "issuesListForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_owner
  }
  property {
    name  = "issuesListForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_page
  }
  property {
    name  = "issuesListForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_per_page
  }
  property {
    name  = "issuesListForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_repo
  }
  property {
    name  = "issuesListForRepo_since"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_since
  }
  property {
    name  = "issuesListForRepo_sort"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_sort
  }
  property {
    name  = "issuesListForRepo_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_for_repo_state
  }
  property {
    name  = "issuesListLabelsForMilestone_milestone_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_milestone_milestone_number
  }
  property {
    name  = "issuesListLabelsForMilestone_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_milestone_owner
  }
  property {
    name  = "issuesListLabelsForMilestone_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_milestone_page
  }
  property {
    name  = "issuesListLabelsForMilestone_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_milestone_per_page
  }
  property {
    name  = "issuesListLabelsForMilestone_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_milestone_repo
  }
  property {
    name  = "issuesListLabelsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_repo_owner
  }
  property {
    name  = "issuesListLabelsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_repo_page
  }
  property {
    name  = "issuesListLabelsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_repo_per_page
  }
  property {
    name  = "issuesListLabelsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_for_repo_repo
  }
  property {
    name  = "issuesListLabelsOnIssue_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_on_issue_issue_number
  }
  property {
    name  = "issuesListLabelsOnIssue_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_on_issue_owner
  }
  property {
    name  = "issuesListLabelsOnIssue_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_on_issue_page
  }
  property {
    name  = "issuesListLabelsOnIssue_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_on_issue_per_page
  }
  property {
    name  = "issuesListLabelsOnIssue_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels_on_issue_repo
  }
  property {
    name  = "issuesListMilestones_direction"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_milestones_direction
  }
  property {
    name  = "issuesListMilestones_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_milestones_owner
  }
  property {
    name  = "issuesListMilestones_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_milestones_page
  }
  property {
    name  = "issuesListMilestones_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_milestones_per_page
  }
  property {
    name  = "issuesListMilestones_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_milestones_repo
  }
  property {
    name  = "issuesListMilestones_sort"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_milestones_sort
  }
  property {
    name  = "issuesListMilestones_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_milestones_state
  }
  property {
    name  = "issuesList_collab"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_collab
  }
  property {
    name  = "issuesList_direction"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_direction
  }
  property {
    name  = "issuesList_filter"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_filter
  }
  property {
    name  = "issuesList_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_labels
  }
  property {
    name  = "issuesList_orgs"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_orgs
  }
  property {
    name  = "issuesList_owned"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_owned
  }
  property {
    name  = "issuesList_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_page
  }
  property {
    name  = "issuesList_per_page"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_per_page
  }
  property {
    name  = "issuesList_pulls"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_pulls
  }
  property {
    name  = "issuesList_since"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_since
  }
  property {
    name  = "issuesList_sort"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_sort
  }
  property {
    name  = "issuesList_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_list_state
  }
  property {
    name  = "issuesLock_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_lock_issue_number
  }
  property {
    name  = "issuesLock_issuesLockRequest_IssuesLockRequest_lock_reason"
    type  = "string"
    value = var.connector-oai-github_property_issues_lock_issues_lock_request_issues_lock_request_lock_reason
  }
  property {
    name  = "issuesLock_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_lock_owner
  }
  property {
    name  = "issuesLock_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_lock_repo
  }
  property {
    name  = "issuesRemoveAllLabels_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_all_labels_issue_number
  }
  property {
    name  = "issuesRemoveAllLabels_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_all_labels_owner
  }
  property {
    name  = "issuesRemoveAllLabels_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_all_labels_repo
  }
  property {
    name  = "issuesRemoveAssignees_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_assignees_issue_number
  }
  property {
    name  = "issuesRemoveAssignees_issuesRemoveAssigneesRequest_IssuesRemoveAssigneesRequest_assignees"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_assignees_issues_remove_assignees_request_issues_remove_assignees_request_assignees
  }
  property {
    name  = "issuesRemoveAssignees_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_assignees_owner
  }
  property {
    name  = "issuesRemoveAssignees_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_assignees_repo
  }
  property {
    name  = "issuesRemoveLabel_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_label_issue_number
  }
  property {
    name  = "issuesRemoveLabel_name"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_label_name
  }
  property {
    name  = "issuesRemoveLabel_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_label_owner
  }
  property {
    name  = "issuesRemoveLabel_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_remove_label_repo
  }
  property {
    name  = "issuesSetLabels_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_set_labels_issue_number
  }
  property {
    name  = "issuesSetLabels_issuesSetLabelsRequest_IssuesSetLabelsRequest_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_set_labels_issues_set_labels_request_issues_set_labels_request_labels
  }
  property {
    name  = "issuesSetLabels_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_set_labels_owner
  }
  property {
    name  = "issuesSetLabels_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_set_labels_repo
  }
  property {
    name  = "issuesUnlock_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_unlock_issue_number
  }
  property {
    name  = "issuesUnlock_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_unlock_owner
  }
  property {
    name  = "issuesUnlock_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_unlock_repo
  }
  property {
    name  = "issuesUpdateComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_comment_comment_id
  }
  property {
    name  = "issuesUpdateComment_issuesUpdateCommentRequest_IssuesUpdateCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_comment_issues_update_comment_request_issues_update_comment_request_body
  }
  property {
    name  = "issuesUpdateComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_comment_owner
  }
  property {
    name  = "issuesUpdateComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_comment_repo
  }
  property {
    name  = "issuesUpdateLabel_issuesUpdateLabelRequest_IssuesUpdateLabelRequest_color"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_label_issues_update_label_request_issues_update_label_request_color
  }
  property {
    name  = "issuesUpdateLabel_issuesUpdateLabelRequest_IssuesUpdateLabelRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_label_issues_update_label_request_issues_update_label_request_description
  }
  property {
    name  = "issuesUpdateLabel_issuesUpdateLabelRequest_IssuesUpdateLabelRequest_new_name"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_label_issues_update_label_request_issues_update_label_request_new_name
  }
  property {
    name  = "issuesUpdateLabel_name"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_label_name
  }
  property {
    name  = "issuesUpdateLabel_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_label_owner
  }
  property {
    name  = "issuesUpdateLabel_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_label_repo
  }
  property {
    name  = "issuesUpdateMilestone_issuesUpdateMilestoneRequest_IssuesUpdateMilestoneRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_milestone_issues_update_milestone_request_issues_update_milestone_request_description
  }
  property {
    name  = "issuesUpdateMilestone_issuesUpdateMilestoneRequest_IssuesUpdateMilestoneRequest_due_on"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_milestone_issues_update_milestone_request_issues_update_milestone_request_due_on
  }
  property {
    name  = "issuesUpdateMilestone_issuesUpdateMilestoneRequest_IssuesUpdateMilestoneRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_milestone_issues_update_milestone_request_issues_update_milestone_request_state
  }
  property {
    name  = "issuesUpdateMilestone_issuesUpdateMilestoneRequest_IssuesUpdateMilestoneRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_milestone_issues_update_milestone_request_issues_update_milestone_request_title
  }
  property {
    name  = "issuesUpdateMilestone_milestone_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_milestone_milestone_number
  }
  property {
    name  = "issuesUpdateMilestone_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_milestone_owner
  }
  property {
    name  = "issuesUpdateMilestone_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_milestone_repo
  }
  property {
    name  = "issuesUpdate_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_issue_number
  }
  property {
    name  = "issuesUpdate_issuesUpdateRequest_IssuesUpdateRequest_assignee"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_issues_update_request_issues_update_request_assignee
  }
  property {
    name  = "issuesUpdate_issuesUpdateRequest_IssuesUpdateRequest_assignees"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_issues_update_request_issues_update_request_assignees
  }
  property {
    name  = "issuesUpdate_issuesUpdateRequest_IssuesUpdateRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_issues_update_request_issues_update_request_body
  }
  property {
    name  = "issuesUpdate_issuesUpdateRequest_IssuesUpdateRequest_labels"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_issues_update_request_issues_update_request_labels
  }
  property {
    name  = "issuesUpdate_issuesUpdateRequest_IssuesUpdateRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_issues_update_request_issues_update_request_state
  }
  property {
    name  = "issuesUpdate_owner"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_owner
  }
  property {
    name  = "issuesUpdate_repo"
    type  = "string"
    value = var.connector-oai-github_property_issues_update_repo
  }
  property {
    name  = "licensesGetAllCommonlyUsed_featured"
    type  = "string"
    value = var.connector-oai-github_property_licenses_get_all_commonly_used_featured
  }
  property {
    name  = "licensesGetAllCommonlyUsed_page"
    type  = "string"
    value = var.connector-oai-github_property_licenses_get_all_commonly_used_page
  }
  property {
    name  = "licensesGetAllCommonlyUsed_per_page"
    type  = "string"
    value = var.connector-oai-github_property_licenses_get_all_commonly_used_per_page
  }
  property {
    name  = "licensesGetForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_licenses_get_for_repo_owner
  }
  property {
    name  = "licensesGetForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_licenses_get_for_repo_repo
  }
  property {
    name  = "licensesGet_license"
    type  = "string"
    value = var.connector-oai-github_property_licenses_get_license
  }
  property {
    name  = "markdownRenderRaw_body"
    type  = "string"
    value = var.connector-oai-github_property_markdown_render_raw_body
  }
  property {
    name  = "markdownRender_markdownRenderRequest_MarkdownRenderRequest_context"
    type  = "string"
    value = var.connector-oai-github_property_markdown_render_markdown_render_request_markdown_render_request_context
  }
  property {
    name  = "markdownRender_markdownRenderRequest_MarkdownRenderRequest_mode"
    type  = "string"
    value = var.connector-oai-github_property_markdown_render_markdown_render_request_markdown_render_request_mode
  }
  property {
    name  = "markdownRender_markdownRenderRequest_MarkdownRenderRequest_text"
    type  = "string"
    value = var.connector-oai-github_property_markdown_render_markdown_render_request_markdown_render_request_text
  }
  property {
    name  = "metaGetOctocat_s"
    type  = "string"
    value = var.connector-oai-github_property_meta_get_octocat_s
  }
  property {
    name  = "oauthAuthorizationsCreateAuthorization_oauthAuthorizationsCreateAuthorizationRequest_OauthAuthorizationsCreateAuthorizationRequest_client_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_create_authorization_oauth_authorizations_create_authorization_request_oauth_authorizations_create_authorization_request_client_id
  }
  property {
    name  = "oauthAuthorizationsCreateAuthorization_oauthAuthorizationsCreateAuthorizationRequest_OauthAuthorizationsCreateAuthorizationRequest_client_secret"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_create_authorization_oauth_authorizations_create_authorization_request_oauth_authorizations_create_authorization_request_client_secret
  }
  property {
    name  = "oauthAuthorizationsCreateAuthorization_oauthAuthorizationsCreateAuthorizationRequest_OauthAuthorizationsCreateAuthorizationRequest_fingerprint"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_create_authorization_oauth_authorizations_create_authorization_request_oauth_authorizations_create_authorization_request_fingerprint
  }
  property {
    name  = "oauthAuthorizationsCreateAuthorization_oauthAuthorizationsCreateAuthorizationRequest_OauthAuthorizationsCreateAuthorizationRequest_note"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_create_authorization_oauth_authorizations_create_authorization_request_oauth_authorizations_create_authorization_request_note
  }
  property {
    name  = "oauthAuthorizationsCreateAuthorization_oauthAuthorizationsCreateAuthorizationRequest_OauthAuthorizationsCreateAuthorizationRequest_note_url"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_create_authorization_oauth_authorizations_create_authorization_request_oauth_authorizations_create_authorization_request_note_url
  }
  property {
    name  = "oauthAuthorizationsCreateAuthorization_oauthAuthorizationsCreateAuthorizationRequest_OauthAuthorizationsCreateAuthorizationRequest_scopes"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_create_authorization_oauth_authorizations_create_authorization_request_oauth_authorizations_create_authorization_request_scopes
  }
  property {
    name  = "oauthAuthorizationsDeleteAuthorization_authorization_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_delete_authorization_authorization_id
  }
  property {
    name  = "oauthAuthorizationsDeleteGrant_grant_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_delete_grant_grant_id
  }
  property {
    name  = "oauthAuthorizationsGetAuthorization_authorization_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_authorization_authorization_id
  }
  property {
    name  = "oauthAuthorizationsGetGrant_grant_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_grant_grant_id
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprint_client_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_client_id
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprint_fingerprint"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_fingerprint
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprint_oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_client_secret"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_client_secret
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprint_oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_note"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_note
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprint_oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_note_url"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_note_url
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprint_oauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppAndFingerprintRequest_scopes"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_oauth_authorizations_get_or_create_authorization_for_app_and_fingerprint_request_scopes
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForApp_client_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_client_id
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForApp_oauthAuthorizationsGetOrCreateAuthorizationForAppRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppRequest_client_secret"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_oauth_authorizations_get_or_create_authorization_for_app_request_oauth_authorizations_get_or_create_authorization_for_app_request_client_secret
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForApp_oauthAuthorizationsGetOrCreateAuthorizationForAppRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppRequest_fingerprint"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_oauth_authorizations_get_or_create_authorization_for_app_request_oauth_authorizations_get_or_create_authorization_for_app_request_fingerprint
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForApp_oauthAuthorizationsGetOrCreateAuthorizationForAppRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppRequest_note"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_oauth_authorizations_get_or_create_authorization_for_app_request_oauth_authorizations_get_or_create_authorization_for_app_request_note
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForApp_oauthAuthorizationsGetOrCreateAuthorizationForAppRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppRequest_note_url"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_oauth_authorizations_get_or_create_authorization_for_app_request_oauth_authorizations_get_or_create_authorization_for_app_request_note_url
  }
  property {
    name  = "oauthAuthorizationsGetOrCreateAuthorizationForApp_oauthAuthorizationsGetOrCreateAuthorizationForAppRequest_OauthAuthorizationsGetOrCreateAuthorizationForAppRequest_scopes"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_get_or_create_authorization_for_app_oauth_authorizations_get_or_create_authorization_for_app_request_oauth_authorizations_get_or_create_authorization_for_app_request_scopes
  }
  property {
    name  = "oauthAuthorizationsListAuthorizations_client_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_list_authorizations_client_id
  }
  property {
    name  = "oauthAuthorizationsListAuthorizations_page"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_list_authorizations_page
  }
  property {
    name  = "oauthAuthorizationsListAuthorizations_per_page"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_list_authorizations_per_page
  }
  property {
    name  = "oauthAuthorizationsListGrants_client_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_list_grants_client_id
  }
  property {
    name  = "oauthAuthorizationsListGrants_page"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_list_grants_page
  }
  property {
    name  = "oauthAuthorizationsListGrants_per_page"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_list_grants_per_page
  }
  property {
    name  = "oauthAuthorizationsUpdateAuthorization_authorization_id"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_update_authorization_authorization_id
  }
  property {
    name  = "oauthAuthorizationsUpdateAuthorization_oauthAuthorizationsUpdateAuthorizationRequest_OauthAuthorizationsUpdateAuthorizationRequest_add_scopes"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_update_authorization_oauth_authorizations_update_authorization_request_oauth_authorizations_update_authorization_request_add_scopes
  }
  property {
    name  = "oauthAuthorizationsUpdateAuthorization_oauthAuthorizationsUpdateAuthorizationRequest_OauthAuthorizationsUpdateAuthorizationRequest_fingerprint"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_update_authorization_oauth_authorizations_update_authorization_request_oauth_authorizations_update_authorization_request_fingerprint
  }
  property {
    name  = "oauthAuthorizationsUpdateAuthorization_oauthAuthorizationsUpdateAuthorizationRequest_OauthAuthorizationsUpdateAuthorizationRequest_note"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_update_authorization_oauth_authorizations_update_authorization_request_oauth_authorizations_update_authorization_request_note
  }
  property {
    name  = "oauthAuthorizationsUpdateAuthorization_oauthAuthorizationsUpdateAuthorizationRequest_OauthAuthorizationsUpdateAuthorizationRequest_note_url"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_update_authorization_oauth_authorizations_update_authorization_request_oauth_authorizations_update_authorization_request_note_url
  }
  property {
    name  = "oauthAuthorizationsUpdateAuthorization_oauthAuthorizationsUpdateAuthorizationRequest_OauthAuthorizationsUpdateAuthorizationRequest_remove_scopes"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_update_authorization_oauth_authorizations_update_authorization_request_oauth_authorizations_update_authorization_request_remove_scopes
  }
  property {
    name  = "oauthAuthorizationsUpdateAuthorization_oauthAuthorizationsUpdateAuthorizationRequest_OauthAuthorizationsUpdateAuthorizationRequest_scopes"
    type  = "string"
    value = var.connector-oai-github_property_oauth_authorizations_update_authorization_oauth_authorizations_update_authorization_request_oauth_authorizations_update_authorization_request_scopes
  }
  property {
    name  = "orgsCheckMembershipForUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_check_membership_for_user_org
  }
  property {
    name  = "orgsCheckMembershipForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_check_membership_for_user_username
  }
  property {
    name  = "orgsCheckPublicMembershipForUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_check_public_membership_for_user_org
  }
  property {
    name  = "orgsCheckPublicMembershipForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_check_public_membership_for_user_username
  }
  property {
    name  = "orgsConvertMemberToOutsideCollaborator_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_convert_member_to_outside_collaborator_org
  }
  property {
    name  = "orgsConvertMemberToOutsideCollaborator_orgsConvertMemberToOutsideCollaboratorRequest_OrgsConvertMemberToOutsideCollaboratorRequest_async"
    type  = "string"
    value = var.connector-oai-github_property_orgs_convert_member_to_outside_collaborator_orgs_convert_member_to_outside_collaborator_request_orgs_convert_member_to_outside_collaborator_request_async
  }
  property {
    name  = "orgsConvertMemberToOutsideCollaborator_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_convert_member_to_outside_collaborator_username
  }
  property {
    name  = "orgsCreateWebhook_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_org
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequestConfig_content_type"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_config_content_type
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequestConfig_password"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_config_password
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequestConfig_secret"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_config_secret
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequestConfig_url"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_config_url
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequestConfig_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_config_username
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequest_active"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_active
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequest_events"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_events
  }
  property {
    name  = "orgsCreateWebhook_orgsCreateWebhookRequest_OrgsCreateWebhookRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_orgs_create_webhook_orgs_create_webhook_request_orgs_create_webhook_request_name
  }
  property {
    name  = "orgsDeleteWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_orgs_delete_webhook_hook_id
  }
  property {
    name  = "orgsDeleteWebhook_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_delete_webhook_org
  }
  property {
    name  = "orgsGetMembershipForAuthenticatedUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_membership_for_authenticated_user_org
  }
  property {
    name  = "orgsGetMembershipForUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_membership_for_user_org
  }
  property {
    name  = "orgsGetMembershipForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_membership_for_user_username
  }
  property {
    name  = "orgsGetWebhookConfigForOrg_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_webhook_config_for_org_hook_id
  }
  property {
    name  = "orgsGetWebhookConfigForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_webhook_config_for_org_org
  }
  property {
    name  = "orgsGetWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_webhook_hook_id
  }
  property {
    name  = "orgsGetWebhook_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_webhook_org
  }
  property {
    name  = "orgsGet_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_get_org
  }
  property {
    name  = "orgsListAppInstallations_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_app_installations_org
  }
  property {
    name  = "orgsListAppInstallations_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_app_installations_page
  }
  property {
    name  = "orgsListAppInstallations_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_app_installations_per_page
  }
  property {
    name  = "orgsListForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_for_authenticated_user_page
  }
  property {
    name  = "orgsListForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_for_authenticated_user_per_page
  }
  property {
    name  = "orgsListForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_for_user_page
  }
  property {
    name  = "orgsListForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_for_user_per_page
  }
  property {
    name  = "orgsListForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_for_user_username
  }
  property {
    name  = "orgsListMembers_filter"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_members_filter
  }
  property {
    name  = "orgsListMembers_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_members_org
  }
  property {
    name  = "orgsListMembers_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_members_page
  }
  property {
    name  = "orgsListMembers_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_members_per_page
  }
  property {
    name  = "orgsListMembers_role"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_members_role
  }
  property {
    name  = "orgsListMembershipsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_memberships_for_authenticated_user_page
  }
  property {
    name  = "orgsListMembershipsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_memberships_for_authenticated_user_per_page
  }
  property {
    name  = "orgsListMembershipsForAuthenticatedUser_state"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_memberships_for_authenticated_user_state
  }
  property {
    name  = "orgsListOutsideCollaborators_filter"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_outside_collaborators_filter
  }
  property {
    name  = "orgsListOutsideCollaborators_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_outside_collaborators_org
  }
  property {
    name  = "orgsListOutsideCollaborators_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_outside_collaborators_page
  }
  property {
    name  = "orgsListOutsideCollaborators_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_outside_collaborators_per_page
  }
  property {
    name  = "orgsListPublicMembers_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_public_members_org
  }
  property {
    name  = "orgsListPublicMembers_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_public_members_page
  }
  property {
    name  = "orgsListPublicMembers_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_public_members_per_page
  }
  property {
    name  = "orgsListWebhooks_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_webhooks_org
  }
  property {
    name  = "orgsListWebhooks_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_webhooks_page
  }
  property {
    name  = "orgsListWebhooks_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_webhooks_per_page
  }
  property {
    name  = "orgsList_per_page"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_per_page
  }
  property {
    name  = "orgsList_since"
    type  = "string"
    value = var.connector-oai-github_property_orgs_list_since
  }
  property {
    name  = "orgsPingWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_orgs_ping_webhook_hook_id
  }
  property {
    name  = "orgsPingWebhook_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_ping_webhook_org
  }
  property {
    name  = "orgsRemoveMember_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_member_org
  }
  property {
    name  = "orgsRemoveMember_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_member_username
  }
  property {
    name  = "orgsRemoveMembershipForUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_membership_for_user_org
  }
  property {
    name  = "orgsRemoveMembershipForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_membership_for_user_username
  }
  property {
    name  = "orgsRemoveOutsideCollaborator_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_outside_collaborator_org
  }
  property {
    name  = "orgsRemoveOutsideCollaborator_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_outside_collaborator_username
  }
  property {
    name  = "orgsRemovePublicMembershipForAuthenticatedUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_public_membership_for_authenticated_user_org
  }
  property {
    name  = "orgsRemovePublicMembershipForAuthenticatedUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_remove_public_membership_for_authenticated_user_username
  }
  property {
    name  = "orgsSetMembershipForUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_set_membership_for_user_org
  }
  property {
    name  = "orgsSetMembershipForUser_orgsSetMembershipForUserRequest_OrgsSetMembershipForUserRequest_role"
    type  = "string"
    value = var.connector-oai-github_property_orgs_set_membership_for_user_orgs_set_membership_for_user_request_orgs_set_membership_for_user_request_role
  }
  property {
    name  = "orgsSetMembershipForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_set_membership_for_user_username
  }
  property {
    name  = "orgsSetPublicMembershipForAuthenticatedUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_set_public_membership_for_authenticated_user_org
  }
  property {
    name  = "orgsSetPublicMembershipForAuthenticatedUser_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_set_public_membership_for_authenticated_user_username
  }
  property {
    name  = "orgsUpdateMembershipForAuthenticatedUser_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_membership_for_authenticated_user_org
  }
  property {
    name  = "orgsUpdateMembershipForAuthenticatedUser_orgsUpdateMembershipForAuthenticatedUserRequest_OrgsUpdateMembershipForAuthenticatedUserRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_membership_for_authenticated_user_orgs_update_membership_for_authenticated_user_request_orgs_update_membership_for_authenticated_user_request_state
  }
  property {
    name  = "orgsUpdateWebhookConfigForOrg_appsUpdateWebhookConfigForAppRequest_AppsUpdateWebhookConfigForAppRequest_content_type"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_config_for_org_apps_update_webhook_config_for_app_request_apps_update_webhook_config_for_app_request_content_type
  }
  property {
    name  = "orgsUpdateWebhookConfigForOrg_appsUpdateWebhookConfigForAppRequest_AppsUpdateWebhookConfigForAppRequest_secret"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_config_for_org_apps_update_webhook_config_for_app_request_apps_update_webhook_config_for_app_request_secret
  }
  property {
    name  = "orgsUpdateWebhookConfigForOrg_appsUpdateWebhookConfigForAppRequest_AppsUpdateWebhookConfigForAppRequest_url"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_config_for_org_apps_update_webhook_config_for_app_request_apps_update_webhook_config_for_app_request_url
  }
  property {
    name  = "orgsUpdateWebhookConfigForOrg_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_config_for_org_hook_id
  }
  property {
    name  = "orgsUpdateWebhookConfigForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_config_for_org_org
  }
  property {
    name  = "orgsUpdateWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_hook_id
  }
  property {
    name  = "orgsUpdateWebhook_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_org
  }
  property {
    name  = "orgsUpdateWebhook_orgsUpdateWebhookRequest_OrgsUpdateWebhookRequestConfig_content_type"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_orgs_update_webhook_request_orgs_update_webhook_request_config_content_type
  }
  property {
    name  = "orgsUpdateWebhook_orgsUpdateWebhookRequest_OrgsUpdateWebhookRequestConfig_secret"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_orgs_update_webhook_request_orgs_update_webhook_request_config_secret
  }
  property {
    name  = "orgsUpdateWebhook_orgsUpdateWebhookRequest_OrgsUpdateWebhookRequestConfig_url"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_orgs_update_webhook_request_orgs_update_webhook_request_config_url
  }
  property {
    name  = "orgsUpdateWebhook_orgsUpdateWebhookRequest_OrgsUpdateWebhookRequest_active"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_orgs_update_webhook_request_orgs_update_webhook_request_active
  }
  property {
    name  = "orgsUpdateWebhook_orgsUpdateWebhookRequest_OrgsUpdateWebhookRequest_events"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_orgs_update_webhook_request_orgs_update_webhook_request_events
  }
  property {
    name  = "orgsUpdateWebhook_orgsUpdateWebhookRequest_OrgsUpdateWebhookRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_webhook_orgs_update_webhook_request_orgs_update_webhook_request_name
  }
  property {
    name  = "orgsUpdate_org"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_org
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_billing_email"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_billing_email
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_blog"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_blog
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_company"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_company
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_default_repository_permission"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_default_repository_permission
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_description
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_email"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_email
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_has_organization_projects"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_has_organization_projects
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_has_repository_projects"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_has_repository_projects
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_location"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_location
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_members_allowed_repository_creation_type"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_members_allowed_repository_creation_type
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_members_can_create_internal_repositories"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_members_can_create_internal_repositories
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_members_can_create_pages"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_members_can_create_pages
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_members_can_create_private_repositories"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_members_can_create_private_repositories
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_members_can_create_public_repositories"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_members_can_create_public_repositories
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_members_can_create_repositories"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_members_can_create_repositories
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_name
  }
  property {
    name  = "orgsUpdate_orgsUpdateRequest_OrgsUpdateRequest_twitter_username"
    type  = "string"
    value = var.connector-oai-github_property_orgs_update_orgs_update_request_orgs_update_request_twitter_username
  }
  property {
    name  = "projectsAddCollaborator_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_add_collaborator_project_id
  }
  property {
    name  = "projectsAddCollaborator_projectsAddCollaboratorRequest_ProjectsAddCollaboratorRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_projects_add_collaborator_projects_add_collaborator_request_projects_add_collaborator_request_permission
  }
  property {
    name  = "projectsAddCollaborator_username"
    type  = "string"
    value = var.connector-oai-github_property_projects_add_collaborator_username
  }
  property {
    name  = "projectsCreateCard_column_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_card_column_id
  }
  property {
    name  = "projectsCreateCard_projectsCreateCardRequest_ProjectsCreateCardRequest_content_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_card_projects_create_card_request_projects_create_card_request_content_id
  }
  property {
    name  = "projectsCreateCard_projectsCreateCardRequest_ProjectsCreateCardRequest_content_type"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_card_projects_create_card_request_projects_create_card_request_content_type
  }
  property {
    name  = "projectsCreateCard_projectsCreateCardRequest_ProjectsCreateCardRequest_note"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_card_projects_create_card_request_projects_create_card_request_note
  }
  property {
    name  = "projectsCreateColumn_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_column_project_id
  }
  property {
    name  = "projectsCreateColumn_projectsUpdateColumnRequest_ProjectsUpdateColumnRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_column_projects_update_column_request_projects_update_column_request_name
  }
  property {
    name  = "projectsCreateForAuthenticatedUser_projectsCreateForAuthenticatedUserRequest_ProjectsCreateForAuthenticatedUserRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_authenticated_user_projects_create_for_authenticated_user_request_projects_create_for_authenticated_user_request_body
  }
  property {
    name  = "projectsCreateForAuthenticatedUser_projectsCreateForAuthenticatedUserRequest_ProjectsCreateForAuthenticatedUserRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_authenticated_user_projects_create_for_authenticated_user_request_projects_create_for_authenticated_user_request_name
  }
  property {
    name  = "projectsCreateForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_org_org
  }
  property {
    name  = "projectsCreateForOrg_projectsCreateForOrgRequest_ProjectsCreateForOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_org_projects_create_for_org_request_projects_create_for_org_request_body
  }
  property {
    name  = "projectsCreateForOrg_projectsCreateForOrgRequest_ProjectsCreateForOrgRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_org_projects_create_for_org_request_projects_create_for_org_request_name
  }
  property {
    name  = "projectsCreateForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_repo_owner
  }
  property {
    name  = "projectsCreateForRepo_projectsCreateForOrgRequest_ProjectsCreateForOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_repo_projects_create_for_org_request_projects_create_for_org_request_body
  }
  property {
    name  = "projectsCreateForRepo_projectsCreateForOrgRequest_ProjectsCreateForOrgRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_repo_projects_create_for_org_request_projects_create_for_org_request_name
  }
  property {
    name  = "projectsCreateForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_projects_create_for_repo_repo
  }
  property {
    name  = "projectsDeleteCard_card_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_delete_card_card_id
  }
  property {
    name  = "projectsDeleteColumn_column_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_delete_column_column_id
  }
  property {
    name  = "projectsDelete_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_delete_project_id
  }
  property {
    name  = "projectsGetCard_card_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_get_card_card_id
  }
  property {
    name  = "projectsGetColumn_column_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_get_column_column_id
  }
  property {
    name  = "projectsGetPermissionForUser_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_get_permission_for_user_project_id
  }
  property {
    name  = "projectsGetPermissionForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_projects_get_permission_for_user_username
  }
  property {
    name  = "projectsGet_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_get_project_id
  }
  property {
    name  = "projectsListCards_archived_state"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_cards_archived_state
  }
  property {
    name  = "projectsListCards_column_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_cards_column_id
  }
  property {
    name  = "projectsListCards_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_cards_page
  }
  property {
    name  = "projectsListCards_per_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_cards_per_page
  }
  property {
    name  = "projectsListCollaborators_affiliation"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_collaborators_affiliation
  }
  property {
    name  = "projectsListCollaborators_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_collaborators_page
  }
  property {
    name  = "projectsListCollaborators_per_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_collaborators_per_page
  }
  property {
    name  = "projectsListCollaborators_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_collaborators_project_id
  }
  property {
    name  = "projectsListColumns_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_columns_page
  }
  property {
    name  = "projectsListColumns_per_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_columns_per_page
  }
  property {
    name  = "projectsListColumns_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_columns_project_id
  }
  property {
    name  = "projectsListForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_org_org
  }
  property {
    name  = "projectsListForOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_org_page
  }
  property {
    name  = "projectsListForOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_org_per_page
  }
  property {
    name  = "projectsListForOrg_state"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_org_state
  }
  property {
    name  = "projectsListForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_repo_owner
  }
  property {
    name  = "projectsListForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_repo_page
  }
  property {
    name  = "projectsListForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_repo_per_page
  }
  property {
    name  = "projectsListForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_repo_repo
  }
  property {
    name  = "projectsListForRepo_state"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_repo_state
  }
  property {
    name  = "projectsListForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_user_page
  }
  property {
    name  = "projectsListForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_user_per_page
  }
  property {
    name  = "projectsListForUser_state"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_user_state
  }
  property {
    name  = "projectsListForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_projects_list_for_user_username
  }
  property {
    name  = "projectsMoveCard_card_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_move_card_card_id
  }
  property {
    name  = "projectsMoveCard_projectsMoveCardRequest_ProjectsMoveCardRequest_column_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_move_card_projects_move_card_request_projects_move_card_request_column_id
  }
  property {
    name  = "projectsMoveCard_projectsMoveCardRequest_ProjectsMoveCardRequest_position"
    type  = "string"
    value = var.connector-oai-github_property_projects_move_card_projects_move_card_request_projects_move_card_request_position
  }
  property {
    name  = "projectsMoveColumn_column_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_move_column_column_id
  }
  property {
    name  = "projectsMoveColumn_projectsMoveColumnRequest_ProjectsMoveColumnRequest_position"
    type  = "string"
    value = var.connector-oai-github_property_projects_move_column_projects_move_column_request_projects_move_column_request_position
  }
  property {
    name  = "projectsRemoveCollaborator_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_remove_collaborator_project_id
  }
  property {
    name  = "projectsRemoveCollaborator_username"
    type  = "string"
    value = var.connector-oai-github_property_projects_remove_collaborator_username
  }
  property {
    name  = "projectsUpdateCard_card_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_card_card_id
  }
  property {
    name  = "projectsUpdateCard_projectsUpdateCardRequest_ProjectsUpdateCardRequest_archived"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_card_projects_update_card_request_projects_update_card_request_archived
  }
  property {
    name  = "projectsUpdateCard_projectsUpdateCardRequest_ProjectsUpdateCardRequest_note"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_card_projects_update_card_request_projects_update_card_request_note
  }
  property {
    name  = "projectsUpdateColumn_column_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_column_column_id
  }
  property {
    name  = "projectsUpdateColumn_projectsUpdateColumnRequest_ProjectsUpdateColumnRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_column_projects_update_column_request_projects_update_column_request_name
  }
  property {
    name  = "projectsUpdate_project_id"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_project_id
  }
  property {
    name  = "projectsUpdate_projectsUpdateRequest_ProjectsUpdateRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_projects_update_request_projects_update_request_body
  }
  property {
    name  = "projectsUpdate_projectsUpdateRequest_ProjectsUpdateRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_projects_update_request_projects_update_request_name
  }
  property {
    name  = "projectsUpdate_projectsUpdateRequest_ProjectsUpdateRequest_organization_permission"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_projects_update_request_projects_update_request_organization_permission
  }
  property {
    name  = "projectsUpdate_projectsUpdateRequest_ProjectsUpdateRequest_private"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_projects_update_request_projects_update_request_private
  }
  property {
    name  = "projectsUpdate_projectsUpdateRequest_ProjectsUpdateRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_projects_update_projects_update_request_projects_update_request_state
  }
  property {
    name  = "pullsCheckIfMerged_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_check_if_merged_owner
  }
  property {
    name  = "pullsCheckIfMerged_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_check_if_merged_pull_number
  }
  property {
    name  = "pullsCheckIfMerged_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_check_if_merged_repo
  }
  property {
    name  = "pullsCreateReplyForReviewComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_reply_for_review_comment_comment_id
  }
  property {
    name  = "pullsCreateReplyForReviewComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_reply_for_review_comment_owner
  }
  property {
    name  = "pullsCreateReplyForReviewComment_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_reply_for_review_comment_pull_number
  }
  property {
    name  = "pullsCreateReplyForReviewComment_pullsCreateReplyForReviewCommentRequest_PullsCreateReplyForReviewCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_reply_for_review_comment_pulls_create_reply_for_review_comment_request_pulls_create_reply_for_review_comment_request_body
  }
  property {
    name  = "pullsCreateReplyForReviewComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_reply_for_review_comment_repo
  }
  property {
    name  = "pullsCreateReviewComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_owner
  }
  property {
    name  = "pullsCreateReviewComment_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pull_number
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_body
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_commit_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_commit_id
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_in_reply_to"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_in_reply_to
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_line"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_line
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_path"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_path
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_position"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_position
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_side"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_side
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_start_line"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_start_line
  }
  property {
    name  = "pullsCreateReviewComment_pullsCreateReviewCommentRequest_PullsCreateReviewCommentRequest_start_side"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_pulls_create_review_comment_request_pulls_create_review_comment_request_start_side
  }
  property {
    name  = "pullsCreateReviewComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_comment_repo
  }
  property {
    name  = "pullsCreateReview_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_owner
  }
  property {
    name  = "pullsCreateReview_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_pull_number
  }
  property {
    name  = "pullsCreateReview_pullsCreateReviewRequest_PullsCreateReviewRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_pulls_create_review_request_pulls_create_review_request_body
  }
  property {
    name  = "pullsCreateReview_pullsCreateReviewRequest_PullsCreateReviewRequest_comments"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_pulls_create_review_request_pulls_create_review_request_comments
  }
  property {
    name  = "pullsCreateReview_pullsCreateReviewRequest_PullsCreateReviewRequest_commit_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_pulls_create_review_request_pulls_create_review_request_commit_id
  }
  property {
    name  = "pullsCreateReview_pullsCreateReviewRequest_PullsCreateReviewRequest_event"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_pulls_create_review_request_pulls_create_review_request_event
  }
  property {
    name  = "pullsCreateReview_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_review_repo
  }
  property {
    name  = "pullsCreate_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_owner
  }
  property {
    name  = "pullsCreate_pullsCreateRequest_PullsCreateRequest_base"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_pulls_create_request_pulls_create_request_base
  }
  property {
    name  = "pullsCreate_pullsCreateRequest_PullsCreateRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_pulls_create_request_pulls_create_request_body
  }
  property {
    name  = "pullsCreate_pullsCreateRequest_PullsCreateRequest_draft"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_pulls_create_request_pulls_create_request_draft
  }
  property {
    name  = "pullsCreate_pullsCreateRequest_PullsCreateRequest_head"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_pulls_create_request_pulls_create_request_head
  }
  property {
    name  = "pullsCreate_pullsCreateRequest_PullsCreateRequest_issue"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_pulls_create_request_pulls_create_request_issue
  }
  property {
    name  = "pullsCreate_pullsCreateRequest_PullsCreateRequest_maintainer_can_modify"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_pulls_create_request_pulls_create_request_maintainer_can_modify
  }
  property {
    name  = "pullsCreate_pullsCreateRequest_PullsCreateRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_pulls_create_request_pulls_create_request_title
  }
  property {
    name  = "pullsCreate_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_create_repo
  }
  property {
    name  = "pullsDeletePendingReview_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_delete_pending_review_owner
  }
  property {
    name  = "pullsDeletePendingReview_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_delete_pending_review_pull_number
  }
  property {
    name  = "pullsDeletePendingReview_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_delete_pending_review_repo
  }
  property {
    name  = "pullsDeletePendingReview_review_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_delete_pending_review_review_id
  }
  property {
    name  = "pullsDeleteReviewComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_delete_review_comment_comment_id
  }
  property {
    name  = "pullsDeleteReviewComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_delete_review_comment_owner
  }
  property {
    name  = "pullsDeleteReviewComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_delete_review_comment_repo
  }
  property {
    name  = "pullsDismissReview_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_dismiss_review_owner
  }
  property {
    name  = "pullsDismissReview_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_dismiss_review_pull_number
  }
  property {
    name  = "pullsDismissReview_pullsDismissReviewRequest_PullsDismissReviewRequest_event"
    type  = "string"
    value = var.connector-oai-github_property_pulls_dismiss_review_pulls_dismiss_review_request_pulls_dismiss_review_request_event
  }
  property {
    name  = "pullsDismissReview_pullsDismissReviewRequest_PullsDismissReviewRequest_message"
    type  = "string"
    value = var.connector-oai-github_property_pulls_dismiss_review_pulls_dismiss_review_request_pulls_dismiss_review_request_message
  }
  property {
    name  = "pullsDismissReview_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_dismiss_review_repo
  }
  property {
    name  = "pullsDismissReview_review_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_dismiss_review_review_id
  }
  property {
    name  = "pullsGetReviewComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_review_comment_comment_id
  }
  property {
    name  = "pullsGetReviewComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_review_comment_owner
  }
  property {
    name  = "pullsGetReviewComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_review_comment_repo
  }
  property {
    name  = "pullsGetReview_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_review_owner
  }
  property {
    name  = "pullsGetReview_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_review_pull_number
  }
  property {
    name  = "pullsGetReview_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_review_repo
  }
  property {
    name  = "pullsGetReview_review_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_review_review_id
  }
  property {
    name  = "pullsGet_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_owner
  }
  property {
    name  = "pullsGet_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_pull_number
  }
  property {
    name  = "pullsGet_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_get_repo
  }
  property {
    name  = "pullsListCommentsForReview_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_comments_for_review_owner
  }
  property {
    name  = "pullsListCommentsForReview_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_comments_for_review_page
  }
  property {
    name  = "pullsListCommentsForReview_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_comments_for_review_per_page
  }
  property {
    name  = "pullsListCommentsForReview_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_comments_for_review_pull_number
  }
  property {
    name  = "pullsListCommentsForReview_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_comments_for_review_repo
  }
  property {
    name  = "pullsListCommentsForReview_review_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_comments_for_review_review_id
  }
  property {
    name  = "pullsListCommits_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_commits_owner
  }
  property {
    name  = "pullsListCommits_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_commits_page
  }
  property {
    name  = "pullsListCommits_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_commits_per_page
  }
  property {
    name  = "pullsListCommits_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_commits_pull_number
  }
  property {
    name  = "pullsListCommits_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_commits_repo
  }
  property {
    name  = "pullsListFiles_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_files_owner
  }
  property {
    name  = "pullsListFiles_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_files_page
  }
  property {
    name  = "pullsListFiles_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_files_per_page
  }
  property {
    name  = "pullsListFiles_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_files_pull_number
  }
  property {
    name  = "pullsListFiles_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_files_repo
  }
  property {
    name  = "pullsListRequestedReviewers_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_requested_reviewers_owner
  }
  property {
    name  = "pullsListRequestedReviewers_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_requested_reviewers_page
  }
  property {
    name  = "pullsListRequestedReviewers_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_requested_reviewers_per_page
  }
  property {
    name  = "pullsListRequestedReviewers_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_requested_reviewers_pull_number
  }
  property {
    name  = "pullsListRequestedReviewers_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_requested_reviewers_repo
  }
  property {
    name  = "pullsListReviewCommentsForRepo_direction"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_for_repo_direction
  }
  property {
    name  = "pullsListReviewCommentsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_for_repo_owner
  }
  property {
    name  = "pullsListReviewCommentsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_for_repo_page
  }
  property {
    name  = "pullsListReviewCommentsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_for_repo_per_page
  }
  property {
    name  = "pullsListReviewCommentsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_for_repo_repo
  }
  property {
    name  = "pullsListReviewCommentsForRepo_since"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_for_repo_since
  }
  property {
    name  = "pullsListReviewCommentsForRepo_sort"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_for_repo_sort
  }
  property {
    name  = "pullsListReviewComments_direction"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_direction
  }
  property {
    name  = "pullsListReviewComments_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_owner
  }
  property {
    name  = "pullsListReviewComments_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_page
  }
  property {
    name  = "pullsListReviewComments_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_per_page
  }
  property {
    name  = "pullsListReviewComments_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_pull_number
  }
  property {
    name  = "pullsListReviewComments_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_repo
  }
  property {
    name  = "pullsListReviewComments_since"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_since
  }
  property {
    name  = "pullsListReviewComments_sort"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_review_comments_sort
  }
  property {
    name  = "pullsListReviews_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_reviews_owner
  }
  property {
    name  = "pullsListReviews_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_reviews_page
  }
  property {
    name  = "pullsListReviews_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_reviews_per_page
  }
  property {
    name  = "pullsListReviews_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_reviews_pull_number
  }
  property {
    name  = "pullsListReviews_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_reviews_repo
  }
  property {
    name  = "pullsList_base"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_base
  }
  property {
    name  = "pullsList_direction"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_direction
  }
  property {
    name  = "pullsList_head"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_head
  }
  property {
    name  = "pullsList_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_owner
  }
  property {
    name  = "pullsList_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_page
  }
  property {
    name  = "pullsList_per_page"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_per_page
  }
  property {
    name  = "pullsList_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_repo
  }
  property {
    name  = "pullsList_sort"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_sort
  }
  property {
    name  = "pullsList_state"
    type  = "string"
    value = var.connector-oai-github_property_pulls_list_state
  }
  property {
    name  = "pullsMerge_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_merge_owner
  }
  property {
    name  = "pullsMerge_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_merge_pull_number
  }
  property {
    name  = "pullsMerge_pullsMergeRequest_PullsMergeRequest_commit_message"
    type  = "string"
    value = var.connector-oai-github_property_pulls_merge_pulls_merge_request_pulls_merge_request_commit_message
  }
  property {
    name  = "pullsMerge_pullsMergeRequest_PullsMergeRequest_commit_title"
    type  = "string"
    value = var.connector-oai-github_property_pulls_merge_pulls_merge_request_pulls_merge_request_commit_title
  }
  property {
    name  = "pullsMerge_pullsMergeRequest_PullsMergeRequest_merge_method"
    type  = "string"
    value = var.connector-oai-github_property_pulls_merge_pulls_merge_request_pulls_merge_request_merge_method
  }
  property {
    name  = "pullsMerge_pullsMergeRequest_PullsMergeRequest_sha"
    type  = "string"
    value = var.connector-oai-github_property_pulls_merge_pulls_merge_request_pulls_merge_request_sha
  }
  property {
    name  = "pullsMerge_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_merge_repo
  }
  property {
    name  = "pullsRemoveRequestedReviewers_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_remove_requested_reviewers_owner
  }
  property {
    name  = "pullsRemoveRequestedReviewers_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_remove_requested_reviewers_pull_number
  }
  property {
    name  = "pullsRemoveRequestedReviewers_pullsRemoveRequestedReviewersRequest_PullsRemoveRequestedReviewersRequest_reviewers"
    type  = "string"
    value = var.connector-oai-github_property_pulls_remove_requested_reviewers_pulls_remove_requested_reviewers_request_pulls_remove_requested_reviewers_request_reviewers
  }
  property {
    name  = "pullsRemoveRequestedReviewers_pullsRemoveRequestedReviewersRequest_PullsRemoveRequestedReviewersRequest_team_reviewers"
    type  = "string"
    value = var.connector-oai-github_property_pulls_remove_requested_reviewers_pulls_remove_requested_reviewers_request_pulls_remove_requested_reviewers_request_team_reviewers
  }
  property {
    name  = "pullsRemoveRequestedReviewers_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_remove_requested_reviewers_repo
  }
  property {
    name  = "pullsRequestReviewers_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_request_reviewers_owner
  }
  property {
    name  = "pullsRequestReviewers_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_request_reviewers_pull_number
  }
  property {
    name  = "pullsRequestReviewers_pullsRequestReviewersRequest_PullsRequestReviewersRequest_reviewers"
    type  = "string"
    value = var.connector-oai-github_property_pulls_request_reviewers_pulls_request_reviewers_request_pulls_request_reviewers_request_reviewers
  }
  property {
    name  = "pullsRequestReviewers_pullsRequestReviewersRequest_PullsRequestReviewersRequest_team_reviewers"
    type  = "string"
    value = var.connector-oai-github_property_pulls_request_reviewers_pulls_request_reviewers_request_pulls_request_reviewers_request_team_reviewers
  }
  property {
    name  = "pullsRequestReviewers_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_request_reviewers_repo
  }
  property {
    name  = "pullsSubmitReview_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_submit_review_owner
  }
  property {
    name  = "pullsSubmitReview_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_submit_review_pull_number
  }
  property {
    name  = "pullsSubmitReview_pullsSubmitReviewRequest_PullsSubmitReviewRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_submit_review_pulls_submit_review_request_pulls_submit_review_request_body
  }
  property {
    name  = "pullsSubmitReview_pullsSubmitReviewRequest_PullsSubmitReviewRequest_event"
    type  = "string"
    value = var.connector-oai-github_property_pulls_submit_review_pulls_submit_review_request_pulls_submit_review_request_event
  }
  property {
    name  = "pullsSubmitReview_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_submit_review_repo
  }
  property {
    name  = "pullsSubmitReview_review_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_submit_review_review_id
  }
  property {
    name  = "pullsUpdateBranch_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_branch_owner
  }
  property {
    name  = "pullsUpdateBranch_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_branch_pull_number
  }
  property {
    name  = "pullsUpdateBranch_pullsUpdateBranchRequest_PullsUpdateBranchRequest_expected_head_sha"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_branch_pulls_update_branch_request_pulls_update_branch_request_expected_head_sha
  }
  property {
    name  = "pullsUpdateBranch_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_branch_repo
  }
  property {
    name  = "pullsUpdateReviewComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_comment_comment_id
  }
  property {
    name  = "pullsUpdateReviewComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_comment_owner
  }
  property {
    name  = "pullsUpdateReviewComment_pullsUpdateReviewCommentRequest_PullsUpdateReviewCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_comment_pulls_update_review_comment_request_pulls_update_review_comment_request_body
  }
  property {
    name  = "pullsUpdateReviewComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_comment_repo
  }
  property {
    name  = "pullsUpdateReview_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_owner
  }
  property {
    name  = "pullsUpdateReview_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_pull_number
  }
  property {
    name  = "pullsUpdateReview_pullsUpdateReviewRequest_PullsUpdateReviewRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_pulls_update_review_request_pulls_update_review_request_body
  }
  property {
    name  = "pullsUpdateReview_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_repo
  }
  property {
    name  = "pullsUpdateReview_review_id"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_review_review_id
  }
  property {
    name  = "pullsUpdate_owner"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_owner
  }
  property {
    name  = "pullsUpdate_pull_number"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_pull_number
  }
  property {
    name  = "pullsUpdate_pullsUpdateRequest_PullsUpdateRequest_base"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_pulls_update_request_pulls_update_request_base
  }
  property {
    name  = "pullsUpdate_pullsUpdateRequest_PullsUpdateRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_pulls_update_request_pulls_update_request_body
  }
  property {
    name  = "pullsUpdate_pullsUpdateRequest_PullsUpdateRequest_maintainer_can_modify"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_pulls_update_request_pulls_update_request_maintainer_can_modify
  }
  property {
    name  = "pullsUpdate_pullsUpdateRequest_PullsUpdateRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_pulls_update_request_pulls_update_request_state
  }
  property {
    name  = "pullsUpdate_pullsUpdateRequest_PullsUpdateRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_pulls_update_request_pulls_update_request_title
  }
  property {
    name  = "pullsUpdate_repo"
    type  = "string"
    value = var.connector-oai-github_property_pulls_update_repo
  }
  property {
    name  = "reactionsCreateForCommitComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_commit_comment_comment_id
  }
  property {
    name  = "reactionsCreateForCommitComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_commit_comment_owner
  }
  property {
    name  = "reactionsCreateForCommitComment_reactionsCreateForCommitCommentRequest_ReactionsCreateForCommitCommentRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_commit_comment_reactions_create_for_commit_comment_request_reactions_create_for_commit_comment_request_content
  }
  property {
    name  = "reactionsCreateForCommitComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_commit_comment_repo
  }
  property {
    name  = "reactionsCreateForIssueComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_comment_comment_id
  }
  property {
    name  = "reactionsCreateForIssueComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_comment_owner
  }
  property {
    name  = "reactionsCreateForIssueComment_reactionsCreateForIssueCommentRequest_ReactionsCreateForIssueCommentRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_comment_reactions_create_for_issue_comment_request_reactions_create_for_issue_comment_request_content
  }
  property {
    name  = "reactionsCreateForIssueComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_comment_repo
  }
  property {
    name  = "reactionsCreateForIssue_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_issue_number
  }
  property {
    name  = "reactionsCreateForIssue_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_owner
  }
  property {
    name  = "reactionsCreateForIssue_reactionsCreateForIssueRequest_ReactionsCreateForIssueRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_reactions_create_for_issue_request_reactions_create_for_issue_request_content
  }
  property {
    name  = "reactionsCreateForIssue_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_issue_repo
  }
  property {
    name  = "reactionsCreateForPullRequestReviewComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_pull_request_review_comment_comment_id
  }
  property {
    name  = "reactionsCreateForPullRequestReviewComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_pull_request_review_comment_owner
  }
  property {
    name  = "reactionsCreateForPullRequestReviewComment_reactionsCreateForPullRequestReviewCommentRequest_ReactionsCreateForPullRequestReviewCommentRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_pull_request_review_comment_reactions_create_for_pull_request_review_comment_request_reactions_create_for_pull_request_review_comment_request_content
  }
  property {
    name  = "reactionsCreateForPullRequestReviewComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_pull_request_review_comment_repo
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentInOrg_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_in_org_comment_number
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_in_org_discussion_number
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_in_org_org
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentInOrg_reactionsCreateForTeamDiscussionCommentInOrgRequest_ReactionsCreateForTeamDiscussionCommentInOrgRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_in_org_reactions_create_for_team_discussion_comment_in_org_request_reactions_create_for_team_discussion_comment_in_org_request_content
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_in_org_team_slug
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentLegacy_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_legacy_comment_number
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_legacy_discussion_number
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentLegacy_reactionsCreateForTeamDiscussionCommentInOrgRequest_ReactionsCreateForTeamDiscussionCommentInOrgRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_legacy_reactions_create_for_team_discussion_comment_in_org_request_reactions_create_for_team_discussion_comment_in_org_request_content
  }
  property {
    name  = "reactionsCreateForTeamDiscussionCommentLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_comment_legacy_team_id
  }
  property {
    name  = "reactionsCreateForTeamDiscussionInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_in_org_discussion_number
  }
  property {
    name  = "reactionsCreateForTeamDiscussionInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_in_org_org
  }
  property {
    name  = "reactionsCreateForTeamDiscussionInOrg_reactionsCreateForTeamDiscussionInOrgRequest_ReactionsCreateForTeamDiscussionInOrgRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_in_org_reactions_create_for_team_discussion_in_org_request_reactions_create_for_team_discussion_in_org_request_content
  }
  property {
    name  = "reactionsCreateForTeamDiscussionInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_in_org_team_slug
  }
  property {
    name  = "reactionsCreateForTeamDiscussionLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_legacy_discussion_number
  }
  property {
    name  = "reactionsCreateForTeamDiscussionLegacy_reactionsCreateForTeamDiscussionInOrgRequest_ReactionsCreateForTeamDiscussionInOrgRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_legacy_reactions_create_for_team_discussion_in_org_request_reactions_create_for_team_discussion_in_org_request_content
  }
  property {
    name  = "reactionsCreateForTeamDiscussionLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_create_for_team_discussion_legacy_team_id
  }
  property {
    name  = "reactionsDeleteForCommitComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_commit_comment_comment_id
  }
  property {
    name  = "reactionsDeleteForCommitComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_commit_comment_owner
  }
  property {
    name  = "reactionsDeleteForCommitComment_reaction_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_commit_comment_reaction_id
  }
  property {
    name  = "reactionsDeleteForCommitComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_commit_comment_repo
  }
  property {
    name  = "reactionsDeleteForIssueComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_comment_comment_id
  }
  property {
    name  = "reactionsDeleteForIssueComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_comment_owner
  }
  property {
    name  = "reactionsDeleteForIssueComment_reaction_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_comment_reaction_id
  }
  property {
    name  = "reactionsDeleteForIssueComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_comment_repo
  }
  property {
    name  = "reactionsDeleteForIssue_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_issue_number
  }
  property {
    name  = "reactionsDeleteForIssue_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_owner
  }
  property {
    name  = "reactionsDeleteForIssue_reaction_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_reaction_id
  }
  property {
    name  = "reactionsDeleteForIssue_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_issue_repo
  }
  property {
    name  = "reactionsDeleteForPullRequestComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_pull_request_comment_comment_id
  }
  property {
    name  = "reactionsDeleteForPullRequestComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_pull_request_comment_owner
  }
  property {
    name  = "reactionsDeleteForPullRequestComment_reaction_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_pull_request_comment_reaction_id
  }
  property {
    name  = "reactionsDeleteForPullRequestComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_pull_request_comment_repo
  }
  property {
    name  = "reactionsDeleteForTeamDiscussionComment_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_comment_comment_number
  }
  property {
    name  = "reactionsDeleteForTeamDiscussionComment_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_comment_discussion_number
  }
  property {
    name  = "reactionsDeleteForTeamDiscussionComment_org"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_comment_org
  }
  property {
    name  = "reactionsDeleteForTeamDiscussionComment_reaction_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_comment_reaction_id
  }
  property {
    name  = "reactionsDeleteForTeamDiscussionComment_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_comment_team_slug
  }
  property {
    name  = "reactionsDeleteForTeamDiscussion_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_discussion_number
  }
  property {
    name  = "reactionsDeleteForTeamDiscussion_org"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_org
  }
  property {
    name  = "reactionsDeleteForTeamDiscussion_reaction_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_reaction_id
  }
  property {
    name  = "reactionsDeleteForTeamDiscussion_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_for_team_discussion_team_slug
  }
  property {
    name  = "reactionsDeleteLegacy_reaction_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_delete_legacy_reaction_id
  }
  property {
    name  = "reactionsListForCommitComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_commit_comment_comment_id
  }
  property {
    name  = "reactionsListForCommitComment_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_commit_comment_content
  }
  property {
    name  = "reactionsListForCommitComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_commit_comment_owner
  }
  property {
    name  = "reactionsListForCommitComment_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_commit_comment_page
  }
  property {
    name  = "reactionsListForCommitComment_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_commit_comment_per_page
  }
  property {
    name  = "reactionsListForCommitComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_commit_comment_repo
  }
  property {
    name  = "reactionsListForIssueComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_comment_comment_id
  }
  property {
    name  = "reactionsListForIssueComment_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_comment_content
  }
  property {
    name  = "reactionsListForIssueComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_comment_owner
  }
  property {
    name  = "reactionsListForIssueComment_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_comment_page
  }
  property {
    name  = "reactionsListForIssueComment_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_comment_per_page
  }
  property {
    name  = "reactionsListForIssueComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_comment_repo
  }
  property {
    name  = "reactionsListForIssue_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_content
  }
  property {
    name  = "reactionsListForIssue_issue_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_issue_number
  }
  property {
    name  = "reactionsListForIssue_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_owner
  }
  property {
    name  = "reactionsListForIssue_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_page
  }
  property {
    name  = "reactionsListForIssue_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_per_page
  }
  property {
    name  = "reactionsListForIssue_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_issue_repo
  }
  property {
    name  = "reactionsListForPullRequestReviewComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_pull_request_review_comment_comment_id
  }
  property {
    name  = "reactionsListForPullRequestReviewComment_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_pull_request_review_comment_content
  }
  property {
    name  = "reactionsListForPullRequestReviewComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_pull_request_review_comment_owner
  }
  property {
    name  = "reactionsListForPullRequestReviewComment_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_pull_request_review_comment_page
  }
  property {
    name  = "reactionsListForPullRequestReviewComment_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_pull_request_review_comment_per_page
  }
  property {
    name  = "reactionsListForPullRequestReviewComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_pull_request_review_comment_repo
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentInOrg_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_in_org_comment_number
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentInOrg_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_in_org_content
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_in_org_discussion_number
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_in_org_org
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_in_org_page
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_in_org_per_page
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_in_org_team_slug
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentLegacy_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_legacy_comment_number
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentLegacy_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_legacy_content
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_legacy_discussion_number
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_legacy_page
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_legacy_per_page
  }
  property {
    name  = "reactionsListForTeamDiscussionCommentLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_comment_legacy_team_id
  }
  property {
    name  = "reactionsListForTeamDiscussionInOrg_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_in_org_content
  }
  property {
    name  = "reactionsListForTeamDiscussionInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_in_org_discussion_number
  }
  property {
    name  = "reactionsListForTeamDiscussionInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_in_org_org
  }
  property {
    name  = "reactionsListForTeamDiscussionInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_in_org_page
  }
  property {
    name  = "reactionsListForTeamDiscussionInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_in_org_per_page
  }
  property {
    name  = "reactionsListForTeamDiscussionInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_in_org_team_slug
  }
  property {
    name  = "reactionsListForTeamDiscussionLegacy_content"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_legacy_content
  }
  property {
    name  = "reactionsListForTeamDiscussionLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_legacy_discussion_number
  }
  property {
    name  = "reactionsListForTeamDiscussionLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_legacy_page
  }
  property {
    name  = "reactionsListForTeamDiscussionLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_legacy_per_page
  }
  property {
    name  = "reactionsListForTeamDiscussionLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_reactions_list_for_team_discussion_legacy_team_id
  }
  property {
    name  = "reposAcceptInvitationForAuthenticatedUser_invitation_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_accept_invitation_for_authenticated_user_invitation_id
  }
  property {
    name  = "reposAddAppAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_app_access_restrictions_branch
  }
  property {
    name  = "reposAddAppAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_app_access_restrictions_owner
  }
  property {
    name  = "reposAddAppAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_app_access_restrictions_repo
  }
  property {
    name  = "reposAddAppAccessRestrictions_reposSetAppAccessRestrictionsRequest_ReposSetAppAccessRestrictionsRequest_apps"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_app_access_restrictions_repos_set_app_access_restrictions_request_repos_set_app_access_restrictions_request_apps
  }
  property {
    name  = "reposAddCollaborator_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_collaborator_owner
  }
  property {
    name  = "reposAddCollaborator_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_collaborator_repo
  }
  property {
    name  = "reposAddCollaborator_reposAddCollaboratorRequest_ReposAddCollaboratorRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_collaborator_repos_add_collaborator_request_repos_add_collaborator_request_permission
  }
  property {
    name  = "reposAddCollaborator_username"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_collaborator_username
  }
  property {
    name  = "reposAddStatusCheckContexts_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_status_check_contexts_branch
  }
  property {
    name  = "reposAddStatusCheckContexts_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_status_check_contexts_owner
  }
  property {
    name  = "reposAddStatusCheckContexts_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_status_check_contexts_repo
  }
  property {
    name  = "reposAddStatusCheckContexts_reposSetStatusCheckContextsRequest_ReposSetStatusCheckContextsRequest_contexts"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_status_check_contexts_repos_set_status_check_contexts_request_repos_set_status_check_contexts_request_contexts
  }
  property {
    name  = "reposAddTeamAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_team_access_restrictions_branch
  }
  property {
    name  = "reposAddTeamAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_team_access_restrictions_owner
  }
  property {
    name  = "reposAddTeamAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_team_access_restrictions_repo
  }
  property {
    name  = "reposAddTeamAccessRestrictions_reposSetTeamAccessRestrictionsRequest_ReposSetTeamAccessRestrictionsRequest_teams"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_team_access_restrictions_repos_set_team_access_restrictions_request_repos_set_team_access_restrictions_request_teams
  }
  property {
    name  = "reposAddUserAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_user_access_restrictions_branch
  }
  property {
    name  = "reposAddUserAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_user_access_restrictions_owner
  }
  property {
    name  = "reposAddUserAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_user_access_restrictions_repo
  }
  property {
    name  = "reposAddUserAccessRestrictions_reposSetUserAccessRestrictionsRequest_ReposSetUserAccessRestrictionsRequest_users"
    type  = "string"
    value = var.connector-oai-github_property_repos_add_user_access_restrictions_repos_set_user_access_restrictions_request_repos_set_user_access_restrictions_request_users
  }
  property {
    name  = "reposCheckCollaborator_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_check_collaborator_owner
  }
  property {
    name  = "reposCheckCollaborator_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_check_collaborator_repo
  }
  property {
    name  = "reposCheckCollaborator_username"
    type  = "string"
    value = var.connector-oai-github_property_repos_check_collaborator_username
  }
  property {
    name  = "reposCompareCommits_basehead"
    type  = "string"
    value = var.connector-oai-github_property_repos_compare_commits_basehead
  }
  property {
    name  = "reposCompareCommits_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_compare_commits_owner
  }
  property {
    name  = "reposCompareCommits_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_compare_commits_repo
  }
  property {
    name  = "reposCreateCommitComment_commit_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_comment_commit_sha
  }
  property {
    name  = "reposCreateCommitComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_comment_owner
  }
  property {
    name  = "reposCreateCommitComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_comment_repo
  }
  property {
    name  = "reposCreateCommitComment_reposCreateCommitCommentRequest_ReposCreateCommitCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_comment_repos_create_commit_comment_request_repos_create_commit_comment_request_body
  }
  property {
    name  = "reposCreateCommitComment_reposCreateCommitCommentRequest_ReposCreateCommitCommentRequest_line"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_comment_repos_create_commit_comment_request_repos_create_commit_comment_request_line
  }
  property {
    name  = "reposCreateCommitComment_reposCreateCommitCommentRequest_ReposCreateCommitCommentRequest_path"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_comment_repos_create_commit_comment_request_repos_create_commit_comment_request_path
  }
  property {
    name  = "reposCreateCommitComment_reposCreateCommitCommentRequest_ReposCreateCommitCommentRequest_position"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_comment_repos_create_commit_comment_request_repos_create_commit_comment_request_position
  }
  property {
    name  = "reposCreateCommitSignatureProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_signature_protection_branch
  }
  property {
    name  = "reposCreateCommitSignatureProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_signature_protection_owner
  }
  property {
    name  = "reposCreateCommitSignatureProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_signature_protection_repo
  }
  property {
    name  = "reposCreateCommitStatus_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_status_owner
  }
  property {
    name  = "reposCreateCommitStatus_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_status_repo
  }
  property {
    name  = "reposCreateCommitStatus_reposCreateCommitStatusRequest_ReposCreateCommitStatusRequest_context"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_status_repos_create_commit_status_request_repos_create_commit_status_request_context
  }
  property {
    name  = "reposCreateCommitStatus_reposCreateCommitStatusRequest_ReposCreateCommitStatusRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_status_repos_create_commit_status_request_repos_create_commit_status_request_description
  }
  property {
    name  = "reposCreateCommitStatus_reposCreateCommitStatusRequest_ReposCreateCommitStatusRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_status_repos_create_commit_status_request_repos_create_commit_status_request_state
  }
  property {
    name  = "reposCreateCommitStatus_reposCreateCommitStatusRequest_ReposCreateCommitStatusRequest_target_url"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_status_repos_create_commit_status_request_repos_create_commit_status_request_target_url
  }
  property {
    name  = "reposCreateCommitStatus_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_commit_status_sha
  }
  property {
    name  = "reposCreateDeployKey_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deploy_key_owner
  }
  property {
    name  = "reposCreateDeployKey_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deploy_key_repo
  }
  property {
    name  = "reposCreateDeployKey_reposCreateDeployKeyRequest_ReposCreateDeployKeyRequest_key"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deploy_key_repos_create_deploy_key_request_repos_create_deploy_key_request_key
  }
  property {
    name  = "reposCreateDeployKey_reposCreateDeployKeyRequest_ReposCreateDeployKeyRequest_read_only"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deploy_key_repos_create_deploy_key_request_repos_create_deploy_key_request_read_only
  }
  property {
    name  = "reposCreateDeployKey_reposCreateDeployKeyRequest_ReposCreateDeployKeyRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deploy_key_repos_create_deploy_key_request_repos_create_deploy_key_request_title
  }
  property {
    name  = "reposCreateDeploymentStatus_deployment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_deployment_id
  }
  property {
    name  = "reposCreateDeploymentStatus_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_owner
  }
  property {
    name  = "reposCreateDeploymentStatus_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repo
  }
  property {
    name  = "reposCreateDeploymentStatus_reposCreateDeploymentStatusRequest_ReposCreateDeploymentStatusRequest_auto_inactive"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repos_create_deployment_status_request_repos_create_deployment_status_request_auto_inactive
  }
  property {
    name  = "reposCreateDeploymentStatus_reposCreateDeploymentStatusRequest_ReposCreateDeploymentStatusRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repos_create_deployment_status_request_repos_create_deployment_status_request_description
  }
  property {
    name  = "reposCreateDeploymentStatus_reposCreateDeploymentStatusRequest_ReposCreateDeploymentStatusRequest_environment"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repos_create_deployment_status_request_repos_create_deployment_status_request_environment
  }
  property {
    name  = "reposCreateDeploymentStatus_reposCreateDeploymentStatusRequest_ReposCreateDeploymentStatusRequest_environment_url"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repos_create_deployment_status_request_repos_create_deployment_status_request_environment_url
  }
  property {
    name  = "reposCreateDeploymentStatus_reposCreateDeploymentStatusRequest_ReposCreateDeploymentStatusRequest_log_url"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repos_create_deployment_status_request_repos_create_deployment_status_request_log_url
  }
  property {
    name  = "reposCreateDeploymentStatus_reposCreateDeploymentStatusRequest_ReposCreateDeploymentStatusRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repos_create_deployment_status_request_repos_create_deployment_status_request_state
  }
  property {
    name  = "reposCreateDeploymentStatus_reposCreateDeploymentStatusRequest_ReposCreateDeploymentStatusRequest_target_url"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_status_repos_create_deployment_status_request_repos_create_deployment_status_request_target_url
  }
  property {
    name  = "reposCreateDeployment_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_owner
  }
  property {
    name  = "reposCreateDeployment_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repo
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_auto_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_auto_merge
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_description
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_environment"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_environment
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_production_environment"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_production_environment
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_ref
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_required_contexts"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_required_contexts
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_task"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_task
  }
  property {
    name  = "reposCreateDeployment_reposCreateDeploymentRequest_ReposCreateDeploymentRequest_transient_environment"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_deployment_repos_create_deployment_request_repos_create_deployment_request_transient_environment
  }
  property {
    name  = "reposCreateDispatchEvent_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_dispatch_event_owner
  }
  property {
    name  = "reposCreateDispatchEvent_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_dispatch_event_repo
  }
  property {
    name  = "reposCreateDispatchEvent_reposCreateDispatchEventRequest_ReposCreateDispatchEventRequest_client_payload"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_dispatch_event_repos_create_dispatch_event_request_repos_create_dispatch_event_request_client_payload
  }
  property {
    name  = "reposCreateDispatchEvent_reposCreateDispatchEventRequest_ReposCreateDispatchEventRequest_event_type"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_dispatch_event_repos_create_dispatch_event_request_repos_create_dispatch_event_request_event_type
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_allow_merge_commit"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_allow_merge_commit
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_allow_rebase_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_allow_rebase_merge
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_allow_squash_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_allow_squash_merge
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_auto_init"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_auto_init
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_delete_branch_on_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_delete_branch_on_merge
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_description
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_gitignore_template"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_gitignore_template
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_has_downloads"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_has_downloads
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_has_issues"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_has_issues
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_has_projects"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_has_projects
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_has_wiki"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_has_wiki
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_homepage"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_homepage
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_is_template"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_is_template
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_license_template"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_license_template
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_name
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_private"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_private
  }
  property {
    name  = "reposCreateForAuthenticatedUser_reposCreateForAuthenticatedUserRequest_ReposCreateForAuthenticatedUserRequest_team_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_for_authenticated_user_repos_create_for_authenticated_user_request_repos_create_for_authenticated_user_request_team_id
  }
  property {
    name  = "reposCreateFork_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_fork_owner
  }
  property {
    name  = "reposCreateFork_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_fork_repo
  }
  property {
    name  = "reposCreateFork_reposCreateForkRequest_ReposCreateForkRequest_organization"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_fork_repos_create_fork_request_repos_create_fork_request_organization
  }
  property {
    name  = "reposCreateInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_org
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_allow_merge_commit"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_allow_merge_commit
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_allow_rebase_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_allow_rebase_merge
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_allow_squash_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_allow_squash_merge
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_auto_init"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_auto_init
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_delete_branch_on_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_delete_branch_on_merge
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_description
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_gitignore_template"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_gitignore_template
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_has_issues"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_has_issues
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_has_projects"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_has_projects
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_has_wiki"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_has_wiki
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_homepage"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_homepage
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_is_template"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_is_template
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_license_template"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_license_template
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_name
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_private"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_private
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_team_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_team_id
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_use_squash_pr_title_as_default"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_use_squash_pr_title_as_default
  }
  property {
    name  = "reposCreateInOrg_reposCreateInOrgRequest_ReposCreateInOrgRequest_visibility"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_in_org_repos_create_in_org_request_repos_create_in_org_request_visibility
  }
  property {
    name  = "reposCreateOrUpdateFileContents_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_owner
  }
  property {
    name  = "reposCreateOrUpdateFileContents_path"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_path
  }
  property {
    name  = "reposCreateOrUpdateFileContents_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repo
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequestAuthor_date"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_author_date
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequestAuthor_email"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_author_email
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequestAuthor_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_author_name
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequestCommitter_date"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_committer_date
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequestCommitter_email"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_committer_email
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequestCommitter_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_committer_name
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequest_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_branch
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequest_content"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_content
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequest_message"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_message
  }
  property {
    name  = "reposCreateOrUpdateFileContents_reposCreateOrUpdateFileContentsRequest_ReposCreateOrUpdateFileContentsRequest_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_or_update_file_contents_repos_create_or_update_file_contents_request_repos_create_or_update_file_contents_request_sha
  }
  property {
    name  = "reposCreatePagesSite_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_pages_site_owner
  }
  property {
    name  = "reposCreatePagesSite_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_pages_site_repo
  }
  property {
    name  = "reposCreatePagesSite_reposCreatePagesSiteRequest_ReposCreatePagesSiteRequestSource_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_pages_site_repos_create_pages_site_request_repos_create_pages_site_request_source_branch
  }
  property {
    name  = "reposCreatePagesSite_reposCreatePagesSiteRequest_ReposCreatePagesSiteRequestSource_path"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_pages_site_repos_create_pages_site_request_repos_create_pages_site_request_source_path
  }
  property {
    name  = "reposCreateRelease_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_owner
  }
  property {
    name  = "reposCreateRelease_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_repo
  }
  property {
    name  = "reposCreateRelease_reposCreateReleaseRequest_ReposCreateReleaseRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_repos_create_release_request_repos_create_release_request_body
  }
  property {
    name  = "reposCreateRelease_reposCreateReleaseRequest_ReposCreateReleaseRequest_draft"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_repos_create_release_request_repos_create_release_request_draft
  }
  property {
    name  = "reposCreateRelease_reposCreateReleaseRequest_ReposCreateReleaseRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_repos_create_release_request_repos_create_release_request_name
  }
  property {
    name  = "reposCreateRelease_reposCreateReleaseRequest_ReposCreateReleaseRequest_prerelease"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_repos_create_release_request_repos_create_release_request_prerelease
  }
  property {
    name  = "reposCreateRelease_reposCreateReleaseRequest_ReposCreateReleaseRequest_tag_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_repos_create_release_request_repos_create_release_request_tag_name
  }
  property {
    name  = "reposCreateRelease_reposCreateReleaseRequest_ReposCreateReleaseRequest_target_commitish"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_release_repos_create_release_request_repos_create_release_request_target_commitish
  }
  property {
    name  = "reposCreateUsingTemplate_reposCreateUsingTemplateRequest_ReposCreateUsingTemplateRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_using_template_repos_create_using_template_request_repos_create_using_template_request_description
  }
  property {
    name  = "reposCreateUsingTemplate_reposCreateUsingTemplateRequest_ReposCreateUsingTemplateRequest_include_all_branches"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_using_template_repos_create_using_template_request_repos_create_using_template_request_include_all_branches
  }
  property {
    name  = "reposCreateUsingTemplate_reposCreateUsingTemplateRequest_ReposCreateUsingTemplateRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_using_template_repos_create_using_template_request_repos_create_using_template_request_name
  }
  property {
    name  = "reposCreateUsingTemplate_reposCreateUsingTemplateRequest_ReposCreateUsingTemplateRequest_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_using_template_repos_create_using_template_request_repos_create_using_template_request_owner
  }
  property {
    name  = "reposCreateUsingTemplate_reposCreateUsingTemplateRequest_ReposCreateUsingTemplateRequest_private"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_using_template_repos_create_using_template_request_repos_create_using_template_request_private
  }
  property {
    name  = "reposCreateUsingTemplate_template_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_using_template_template_owner
  }
  property {
    name  = "reposCreateUsingTemplate_template_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_using_template_template_repo
  }
  property {
    name  = "reposCreateWebhook_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_owner
  }
  property {
    name  = "reposCreateWebhook_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repo
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequestConfig_content_type"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_config_content_type
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequestConfig_digest"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_config_digest
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequestConfig_secret"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_config_secret
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequestConfig_token"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_config_token
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequestConfig_url"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_config_url
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequest_active"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_active
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequest_events"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_events
  }
  property {
    name  = "reposCreateWebhook_reposCreateWebhookRequest_ReposCreateWebhookRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_create_webhook_repos_create_webhook_request_repos_create_webhook_request_name
  }
  property {
    name  = "reposDeclineInvitationForAuthenticatedUser_invitation_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_decline_invitation_for_authenticated_user_invitation_id
  }
  property {
    name  = "reposDeleteAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_access_restrictions_branch
  }
  property {
    name  = "reposDeleteAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_access_restrictions_owner
  }
  property {
    name  = "reposDeleteAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_access_restrictions_repo
  }
  property {
    name  = "reposDeleteAdminBranchProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_admin_branch_protection_branch
  }
  property {
    name  = "reposDeleteAdminBranchProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_admin_branch_protection_owner
  }
  property {
    name  = "reposDeleteAdminBranchProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_admin_branch_protection_repo
  }
  property {
    name  = "reposDeleteBranchProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_branch_protection_branch
  }
  property {
    name  = "reposDeleteBranchProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_branch_protection_owner
  }
  property {
    name  = "reposDeleteBranchProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_branch_protection_repo
  }
  property {
    name  = "reposDeleteCommitComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_commit_comment_comment_id
  }
  property {
    name  = "reposDeleteCommitComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_commit_comment_owner
  }
  property {
    name  = "reposDeleteCommitComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_commit_comment_repo
  }
  property {
    name  = "reposDeleteCommitSignatureProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_commit_signature_protection_branch
  }
  property {
    name  = "reposDeleteCommitSignatureProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_commit_signature_protection_owner
  }
  property {
    name  = "reposDeleteCommitSignatureProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_commit_signature_protection_repo
  }
  property {
    name  = "reposDeleteDeployKey_key_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_deploy_key_key_id
  }
  property {
    name  = "reposDeleteDeployKey_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_deploy_key_owner
  }
  property {
    name  = "reposDeleteDeployKey_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_deploy_key_repo
  }
  property {
    name  = "reposDeleteDeployment_deployment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_deployment_deployment_id
  }
  property {
    name  = "reposDeleteDeployment_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_deployment_owner
  }
  property {
    name  = "reposDeleteDeployment_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_deployment_repo
  }
  property {
    name  = "reposDeleteFile_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_owner
  }
  property {
    name  = "reposDeleteFile_path"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_path
  }
  property {
    name  = "reposDeleteFile_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repo
  }
  property {
    name  = "reposDeleteFile_reposDeleteFileRequest_ReposDeleteFileRequestAuthor_email"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repos_delete_file_request_repos_delete_file_request_author_email
  }
  property {
    name  = "reposDeleteFile_reposDeleteFileRequest_ReposDeleteFileRequestAuthor_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repos_delete_file_request_repos_delete_file_request_author_name
  }
  property {
    name  = "reposDeleteFile_reposDeleteFileRequest_ReposDeleteFileRequestCommitter_email"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repos_delete_file_request_repos_delete_file_request_committer_email
  }
  property {
    name  = "reposDeleteFile_reposDeleteFileRequest_ReposDeleteFileRequestCommitter_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repos_delete_file_request_repos_delete_file_request_committer_name
  }
  property {
    name  = "reposDeleteFile_reposDeleteFileRequest_ReposDeleteFileRequest_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repos_delete_file_request_repos_delete_file_request_branch
  }
  property {
    name  = "reposDeleteFile_reposDeleteFileRequest_ReposDeleteFileRequest_message"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repos_delete_file_request_repos_delete_file_request_message
  }
  property {
    name  = "reposDeleteFile_reposDeleteFileRequest_ReposDeleteFileRequest_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_file_repos_delete_file_request_repos_delete_file_request_sha
  }
  property {
    name  = "reposDeleteInvitation_invitation_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_invitation_invitation_id
  }
  property {
    name  = "reposDeleteInvitation_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_invitation_owner
  }
  property {
    name  = "reposDeleteInvitation_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_invitation_repo
  }
  property {
    name  = "reposDeletePagesSite_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_pages_site_owner
  }
  property {
    name  = "reposDeletePagesSite_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_pages_site_repo
  }
  property {
    name  = "reposDeletePullRequestReviewProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_pull_request_review_protection_branch
  }
  property {
    name  = "reposDeletePullRequestReviewProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_pull_request_review_protection_owner
  }
  property {
    name  = "reposDeletePullRequestReviewProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_pull_request_review_protection_repo
  }
  property {
    name  = "reposDeleteReleaseAsset_asset_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_release_asset_asset_id
  }
  property {
    name  = "reposDeleteReleaseAsset_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_release_asset_owner
  }
  property {
    name  = "reposDeleteReleaseAsset_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_release_asset_repo
  }
  property {
    name  = "reposDeleteRelease_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_release_owner
  }
  property {
    name  = "reposDeleteRelease_release_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_release_release_id
  }
  property {
    name  = "reposDeleteRelease_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_release_repo
  }
  property {
    name  = "reposDeleteWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_webhook_hook_id
  }
  property {
    name  = "reposDeleteWebhook_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_webhook_owner
  }
  property {
    name  = "reposDeleteWebhook_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_webhook_repo
  }
  property {
    name  = "reposDelete_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_owner
  }
  property {
    name  = "reposDelete_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_delete_repo
  }
  property {
    name  = "reposDownloadTarballArchive_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_download_tarball_archive_owner
  }
  property {
    name  = "reposDownloadTarballArchive_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_download_tarball_archive_ref
  }
  property {
    name  = "reposDownloadTarballArchive_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_download_tarball_archive_repo
  }
  property {
    name  = "reposDownloadZipballArchive_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_download_zipball_archive_owner
  }
  property {
    name  = "reposDownloadZipballArchive_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_download_zipball_archive_ref
  }
  property {
    name  = "reposDownloadZipballArchive_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_download_zipball_archive_repo
  }
  property {
    name  = "reposGetAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_access_restrictions_branch
  }
  property {
    name  = "reposGetAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_access_restrictions_owner
  }
  property {
    name  = "reposGetAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_access_restrictions_repo
  }
  property {
    name  = "reposGetAdminBranchProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_admin_branch_protection_branch
  }
  property {
    name  = "reposGetAdminBranchProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_admin_branch_protection_owner
  }
  property {
    name  = "reposGetAdminBranchProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_admin_branch_protection_repo
  }
  property {
    name  = "reposGetAllStatusCheckContexts_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_all_status_check_contexts_branch
  }
  property {
    name  = "reposGetAllStatusCheckContexts_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_all_status_check_contexts_owner
  }
  property {
    name  = "reposGetAllStatusCheckContexts_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_all_status_check_contexts_repo
  }
  property {
    name  = "reposGetAllTopics_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_all_topics_owner
  }
  property {
    name  = "reposGetAllTopics_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_all_topics_page
  }
  property {
    name  = "reposGetAllTopics_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_all_topics_per_page
  }
  property {
    name  = "reposGetAllTopics_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_all_topics_repo
  }
  property {
    name  = "reposGetAppsWithAccessToProtectedBranch_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_apps_with_access_to_protected_branch_branch
  }
  property {
    name  = "reposGetAppsWithAccessToProtectedBranch_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_apps_with_access_to_protected_branch_owner
  }
  property {
    name  = "reposGetAppsWithAccessToProtectedBranch_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_apps_with_access_to_protected_branch_repo
  }
  property {
    name  = "reposGetBranchProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_branch_protection_branch
  }
  property {
    name  = "reposGetBranchProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_branch_protection_owner
  }
  property {
    name  = "reposGetBranchProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_branch_protection_repo
  }
  property {
    name  = "reposGetBranch_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_branch_branch
  }
  property {
    name  = "reposGetBranch_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_branch_owner
  }
  property {
    name  = "reposGetBranch_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_branch_repo
  }
  property {
    name  = "reposGetCodeFrequencyStats_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_code_frequency_stats_owner
  }
  property {
    name  = "reposGetCodeFrequencyStats_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_code_frequency_stats_repo
  }
  property {
    name  = "reposGetCollaboratorPermissionLevel_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_collaborator_permission_level_owner
  }
  property {
    name  = "reposGetCollaboratorPermissionLevel_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_collaborator_permission_level_repo
  }
  property {
    name  = "reposGetCollaboratorPermissionLevel_username"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_collaborator_permission_level_username
  }
  property {
    name  = "reposGetCombinedStatusForRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_combined_status_for_ref_owner
  }
  property {
    name  = "reposGetCombinedStatusForRef_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_combined_status_for_ref_page
  }
  property {
    name  = "reposGetCombinedStatusForRef_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_combined_status_for_ref_per_page
  }
  property {
    name  = "reposGetCombinedStatusForRef_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_combined_status_for_ref_ref
  }
  property {
    name  = "reposGetCombinedStatusForRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_combined_status_for_ref_repo
  }
  property {
    name  = "reposGetCommitActivityStats_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_activity_stats_owner
  }
  property {
    name  = "reposGetCommitActivityStats_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_activity_stats_repo
  }
  property {
    name  = "reposGetCommitComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_comment_comment_id
  }
  property {
    name  = "reposGetCommitComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_comment_owner
  }
  property {
    name  = "reposGetCommitComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_comment_repo
  }
  property {
    name  = "reposGetCommitSignatureProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_signature_protection_branch
  }
  property {
    name  = "reposGetCommitSignatureProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_signature_protection_owner
  }
  property {
    name  = "reposGetCommitSignatureProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_signature_protection_repo
  }
  property {
    name  = "reposGetCommit_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_owner
  }
  property {
    name  = "reposGetCommit_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_page
  }
  property {
    name  = "reposGetCommit_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_per_page
  }
  property {
    name  = "reposGetCommit_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_ref
  }
  property {
    name  = "reposGetCommit_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_commit_repo
  }
  property {
    name  = "reposGetContent_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_content_owner
  }
  property {
    name  = "reposGetContent_path"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_content_path
  }
  property {
    name  = "reposGetContent_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_content_ref
  }
  property {
    name  = "reposGetContent_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_content_repo
  }
  property {
    name  = "reposGetContributorsStats_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_contributors_stats_owner
  }
  property {
    name  = "reposGetContributorsStats_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_contributors_stats_repo
  }
  property {
    name  = "reposGetDeployKey_key_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deploy_key_key_id
  }
  property {
    name  = "reposGetDeployKey_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deploy_key_owner
  }
  property {
    name  = "reposGetDeployKey_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deploy_key_repo
  }
  property {
    name  = "reposGetDeploymentStatus_deployment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deployment_status_deployment_id
  }
  property {
    name  = "reposGetDeploymentStatus_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deployment_status_owner
  }
  property {
    name  = "reposGetDeploymentStatus_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deployment_status_repo
  }
  property {
    name  = "reposGetDeploymentStatus_status_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deployment_status_status_id
  }
  property {
    name  = "reposGetDeployment_deployment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deployment_deployment_id
  }
  property {
    name  = "reposGetDeployment_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deployment_owner
  }
  property {
    name  = "reposGetDeployment_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_deployment_repo
  }
  property {
    name  = "reposGetLatestPagesBuild_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_latest_pages_build_owner
  }
  property {
    name  = "reposGetLatestPagesBuild_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_latest_pages_build_repo
  }
  property {
    name  = "reposGetLatestRelease_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_latest_release_owner
  }
  property {
    name  = "reposGetLatestRelease_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_latest_release_repo
  }
  property {
    name  = "reposGetPagesBuild_build_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pages_build_build_id
  }
  property {
    name  = "reposGetPagesBuild_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pages_build_owner
  }
  property {
    name  = "reposGetPagesBuild_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pages_build_repo
  }
  property {
    name  = "reposGetPages_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pages_owner
  }
  property {
    name  = "reposGetPages_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pages_repo
  }
  property {
    name  = "reposGetParticipationStats_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_participation_stats_owner
  }
  property {
    name  = "reposGetParticipationStats_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_participation_stats_repo
  }
  property {
    name  = "reposGetPullRequestReviewProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pull_request_review_protection_branch
  }
  property {
    name  = "reposGetPullRequestReviewProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pull_request_review_protection_owner
  }
  property {
    name  = "reposGetPullRequestReviewProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_pull_request_review_protection_repo
  }
  property {
    name  = "reposGetPunchCardStats_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_punch_card_stats_owner
  }
  property {
    name  = "reposGetPunchCardStats_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_punch_card_stats_repo
  }
  property {
    name  = "reposGetReadmeInDirectory_dir"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_readme_in_directory_dir
  }
  property {
    name  = "reposGetReadmeInDirectory_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_readme_in_directory_owner
  }
  property {
    name  = "reposGetReadmeInDirectory_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_readme_in_directory_ref
  }
  property {
    name  = "reposGetReadmeInDirectory_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_readme_in_directory_repo
  }
  property {
    name  = "reposGetReadme_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_readme_owner
  }
  property {
    name  = "reposGetReadme_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_readme_ref
  }
  property {
    name  = "reposGetReadme_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_readme_repo
  }
  property {
    name  = "reposGetReleaseAsset_asset_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_asset_asset_id
  }
  property {
    name  = "reposGetReleaseAsset_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_asset_owner
  }
  property {
    name  = "reposGetReleaseAsset_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_asset_repo
  }
  property {
    name  = "reposGetReleaseByTag_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_by_tag_owner
  }
  property {
    name  = "reposGetReleaseByTag_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_by_tag_repo
  }
  property {
    name  = "reposGetReleaseByTag_tag"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_by_tag_tag
  }
  property {
    name  = "reposGetRelease_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_owner
  }
  property {
    name  = "reposGetRelease_release_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_release_id
  }
  property {
    name  = "reposGetRelease_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_release_repo
  }
  property {
    name  = "reposGetStatusChecksProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_status_checks_protection_branch
  }
  property {
    name  = "reposGetStatusChecksProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_status_checks_protection_owner
  }
  property {
    name  = "reposGetStatusChecksProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_status_checks_protection_repo
  }
  property {
    name  = "reposGetTeamsWithAccessToProtectedBranch_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_teams_with_access_to_protected_branch_branch
  }
  property {
    name  = "reposGetTeamsWithAccessToProtectedBranch_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_teams_with_access_to_protected_branch_owner
  }
  property {
    name  = "reposGetTeamsWithAccessToProtectedBranch_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_teams_with_access_to_protected_branch_repo
  }
  property {
    name  = "reposGetUsersWithAccessToProtectedBranch_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_users_with_access_to_protected_branch_branch
  }
  property {
    name  = "reposGetUsersWithAccessToProtectedBranch_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_users_with_access_to_protected_branch_owner
  }
  property {
    name  = "reposGetUsersWithAccessToProtectedBranch_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_users_with_access_to_protected_branch_repo
  }
  property {
    name  = "reposGetWebhookConfigForRepo_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_webhook_config_for_repo_hook_id
  }
  property {
    name  = "reposGetWebhookConfigForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_webhook_config_for_repo_owner
  }
  property {
    name  = "reposGetWebhookConfigForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_webhook_config_for_repo_repo
  }
  property {
    name  = "reposGetWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_webhook_hook_id
  }
  property {
    name  = "reposGetWebhook_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_webhook_owner
  }
  property {
    name  = "reposGetWebhook_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_webhook_repo
  }
  property {
    name  = "reposGet_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_owner
  }
  property {
    name  = "reposGet_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_get_repo
  }
  property {
    name  = "reposListBranchesForHeadCommit_commit_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches_for_head_commit_commit_sha
  }
  property {
    name  = "reposListBranchesForHeadCommit_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches_for_head_commit_owner
  }
  property {
    name  = "reposListBranchesForHeadCommit_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches_for_head_commit_repo
  }
  property {
    name  = "reposListBranches__protected"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches__protected
  }
  property {
    name  = "reposListBranches_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches_owner
  }
  property {
    name  = "reposListBranches_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches_page
  }
  property {
    name  = "reposListBranches_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches_per_page
  }
  property {
    name  = "reposListBranches_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_branches_repo
  }
  property {
    name  = "reposListCollaborators_affiliation"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_collaborators_affiliation
  }
  property {
    name  = "reposListCollaborators_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_collaborators_owner
  }
  property {
    name  = "reposListCollaborators_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_collaborators_page
  }
  property {
    name  = "reposListCollaborators_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_collaborators_per_page
  }
  property {
    name  = "reposListCollaborators_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_collaborators_repo
  }
  property {
    name  = "reposListCommentsForCommit_commit_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_comments_for_commit_commit_sha
  }
  property {
    name  = "reposListCommentsForCommit_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_comments_for_commit_owner
  }
  property {
    name  = "reposListCommentsForCommit_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_comments_for_commit_page
  }
  property {
    name  = "reposListCommentsForCommit_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_comments_for_commit_per_page
  }
  property {
    name  = "reposListCommentsForCommit_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_comments_for_commit_repo
  }
  property {
    name  = "reposListCommitCommentsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_comments_for_repo_owner
  }
  property {
    name  = "reposListCommitCommentsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_comments_for_repo_page
  }
  property {
    name  = "reposListCommitCommentsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_comments_for_repo_per_page
  }
  property {
    name  = "reposListCommitCommentsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_comments_for_repo_repo
  }
  property {
    name  = "reposListCommitStatusesForRef_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_statuses_for_ref_owner
  }
  property {
    name  = "reposListCommitStatusesForRef_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_statuses_for_ref_page
  }
  property {
    name  = "reposListCommitStatusesForRef_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_statuses_for_ref_per_page
  }
  property {
    name  = "reposListCommitStatusesForRef_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_statuses_for_ref_ref
  }
  property {
    name  = "reposListCommitStatusesForRef_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commit_statuses_for_ref_repo
  }
  property {
    name  = "reposListCommits_author"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_author
  }
  property {
    name  = "reposListCommits_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_owner
  }
  property {
    name  = "reposListCommits_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_page
  }
  property {
    name  = "reposListCommits_path"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_path
  }
  property {
    name  = "reposListCommits_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_per_page
  }
  property {
    name  = "reposListCommits_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_repo
  }
  property {
    name  = "reposListCommits_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_sha
  }
  property {
    name  = "reposListCommits_since"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_since
  }
  property {
    name  = "reposListCommits_until"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_commits_until
  }
  property {
    name  = "reposListContributors_anon"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_contributors_anon
  }
  property {
    name  = "reposListContributors_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_contributors_owner
  }
  property {
    name  = "reposListContributors_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_contributors_page
  }
  property {
    name  = "reposListContributors_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_contributors_per_page
  }
  property {
    name  = "reposListContributors_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_contributors_repo
  }
  property {
    name  = "reposListDeployKeys_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deploy_keys_owner
  }
  property {
    name  = "reposListDeployKeys_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deploy_keys_page
  }
  property {
    name  = "reposListDeployKeys_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deploy_keys_per_page
  }
  property {
    name  = "reposListDeployKeys_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deploy_keys_repo
  }
  property {
    name  = "reposListDeploymentStatuses_deployment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployment_statuses_deployment_id
  }
  property {
    name  = "reposListDeploymentStatuses_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployment_statuses_owner
  }
  property {
    name  = "reposListDeploymentStatuses_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployment_statuses_page
  }
  property {
    name  = "reposListDeploymentStatuses_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployment_statuses_per_page
  }
  property {
    name  = "reposListDeploymentStatuses_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployment_statuses_repo
  }
  property {
    name  = "reposListDeployments_environment"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_environment
  }
  property {
    name  = "reposListDeployments_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_owner
  }
  property {
    name  = "reposListDeployments_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_page
  }
  property {
    name  = "reposListDeployments_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_per_page
  }
  property {
    name  = "reposListDeployments_ref"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_ref
  }
  property {
    name  = "reposListDeployments_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_repo
  }
  property {
    name  = "reposListDeployments_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_sha
  }
  property {
    name  = "reposListDeployments_task"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_deployments_task
  }
  property {
    name  = "reposListForAuthenticatedUser_affiliation"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_affiliation
  }
  property {
    name  = "reposListForAuthenticatedUser_before"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_before
  }
  property {
    name  = "reposListForAuthenticatedUser_direction"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_direction
  }
  property {
    name  = "reposListForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_page
  }
  property {
    name  = "reposListForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_per_page
  }
  property {
    name  = "reposListForAuthenticatedUser_since"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_since
  }
  property {
    name  = "reposListForAuthenticatedUser_sort"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_sort
  }
  property {
    name  = "reposListForAuthenticatedUser_type"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_type
  }
  property {
    name  = "reposListForAuthenticatedUser_visibility"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_authenticated_user_visibility
  }
  property {
    name  = "reposListForOrg_direction"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_org_direction
  }
  property {
    name  = "reposListForOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_org_org
  }
  property {
    name  = "reposListForOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_org_page
  }
  property {
    name  = "reposListForOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_org_per_page
  }
  property {
    name  = "reposListForOrg_sort"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_org_sort
  }
  property {
    name  = "reposListForOrg_type"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_org_type
  }
  property {
    name  = "reposListForUser_direction"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_user_direction
  }
  property {
    name  = "reposListForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_user_page
  }
  property {
    name  = "reposListForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_user_per_page
  }
  property {
    name  = "reposListForUser_sort"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_user_sort
  }
  property {
    name  = "reposListForUser_type"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_user_type
  }
  property {
    name  = "reposListForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_for_user_username
  }
  property {
    name  = "reposListForks_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_forks_owner
  }
  property {
    name  = "reposListForks_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_forks_page
  }
  property {
    name  = "reposListForks_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_forks_per_page
  }
  property {
    name  = "reposListForks_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_forks_repo
  }
  property {
    name  = "reposListForks_sort"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_forks_sort
  }
  property {
    name  = "reposListInvitationsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_invitations_for_authenticated_user_page
  }
  property {
    name  = "reposListInvitationsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_invitations_for_authenticated_user_per_page
  }
  property {
    name  = "reposListInvitations_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_invitations_owner
  }
  property {
    name  = "reposListInvitations_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_invitations_page
  }
  property {
    name  = "reposListInvitations_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_invitations_per_page
  }
  property {
    name  = "reposListInvitations_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_invitations_repo
  }
  property {
    name  = "reposListLanguages_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_languages_owner
  }
  property {
    name  = "reposListLanguages_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_languages_repo
  }
  property {
    name  = "reposListPagesBuilds_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pages_builds_owner
  }
  property {
    name  = "reposListPagesBuilds_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pages_builds_page
  }
  property {
    name  = "reposListPagesBuilds_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pages_builds_per_page
  }
  property {
    name  = "reposListPagesBuilds_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pages_builds_repo
  }
  property {
    name  = "reposListPublic_since"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_public_since
  }
  property {
    name  = "reposListPublic_visibility"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_public_visibility
  }
  property {
    name  = "reposListPullRequestsAssociatedWithCommit_commit_sha"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pull_requests_associated_with_commit_commit_sha
  }
  property {
    name  = "reposListPullRequestsAssociatedWithCommit_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pull_requests_associated_with_commit_owner
  }
  property {
    name  = "reposListPullRequestsAssociatedWithCommit_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pull_requests_associated_with_commit_page
  }
  property {
    name  = "reposListPullRequestsAssociatedWithCommit_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pull_requests_associated_with_commit_per_page
  }
  property {
    name  = "reposListPullRequestsAssociatedWithCommit_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_pull_requests_associated_with_commit_repo
  }
  property {
    name  = "reposListReleaseAssets_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_release_assets_owner
  }
  property {
    name  = "reposListReleaseAssets_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_release_assets_page
  }
  property {
    name  = "reposListReleaseAssets_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_release_assets_per_page
  }
  property {
    name  = "reposListReleaseAssets_release_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_release_assets_release_id
  }
  property {
    name  = "reposListReleaseAssets_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_release_assets_repo
  }
  property {
    name  = "reposListReleases_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_releases_owner
  }
  property {
    name  = "reposListReleases_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_releases_page
  }
  property {
    name  = "reposListReleases_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_releases_per_page
  }
  property {
    name  = "reposListReleases_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_releases_repo
  }
  property {
    name  = "reposListTags_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_tags_owner
  }
  property {
    name  = "reposListTags_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_tags_page
  }
  property {
    name  = "reposListTags_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_tags_per_page
  }
  property {
    name  = "reposListTags_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_tags_repo
  }
  property {
    name  = "reposListTeams_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_teams_owner
  }
  property {
    name  = "reposListTeams_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_teams_page
  }
  property {
    name  = "reposListTeams_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_teams_per_page
  }
  property {
    name  = "reposListTeams_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_teams_repo
  }
  property {
    name  = "reposListWebhooks_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_webhooks_owner
  }
  property {
    name  = "reposListWebhooks_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_webhooks_page
  }
  property {
    name  = "reposListWebhooks_per_page"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_webhooks_per_page
  }
  property {
    name  = "reposListWebhooks_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_list_webhooks_repo
  }
  property {
    name  = "reposMerge_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_merge_owner
  }
  property {
    name  = "reposMerge_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_merge_repo
  }
  property {
    name  = "reposMerge_reposMergeRequest_ReposMergeRequest_base"
    type  = "string"
    value = var.connector-oai-github_property_repos_merge_repos_merge_request_repos_merge_request_base
  }
  property {
    name  = "reposMerge_reposMergeRequest_ReposMergeRequest_commit_message"
    type  = "string"
    value = var.connector-oai-github_property_repos_merge_repos_merge_request_repos_merge_request_commit_message
  }
  property {
    name  = "reposMerge_reposMergeRequest_ReposMergeRequest_head"
    type  = "string"
    value = var.connector-oai-github_property_repos_merge_repos_merge_request_repos_merge_request_head
  }
  property {
    name  = "reposPingWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_ping_webhook_hook_id
  }
  property {
    name  = "reposPingWebhook_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_ping_webhook_owner
  }
  property {
    name  = "reposPingWebhook_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_ping_webhook_repo
  }
  property {
    name  = "reposRemoveAppAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_app_access_restrictions_branch
  }
  property {
    name  = "reposRemoveAppAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_app_access_restrictions_owner
  }
  property {
    name  = "reposRemoveAppAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_app_access_restrictions_repo
  }
  property {
    name  = "reposRemoveAppAccessRestrictions_reposSetAppAccessRestrictionsRequest_ReposSetAppAccessRestrictionsRequest_apps"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_app_access_restrictions_repos_set_app_access_restrictions_request_repos_set_app_access_restrictions_request_apps
  }
  property {
    name  = "reposRemoveCollaborator_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_collaborator_owner
  }
  property {
    name  = "reposRemoveCollaborator_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_collaborator_repo
  }
  property {
    name  = "reposRemoveCollaborator_username"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_collaborator_username
  }
  property {
    name  = "reposRemoveStatusCheckContexts_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_status_check_contexts_branch
  }
  property {
    name  = "reposRemoveStatusCheckContexts_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_status_check_contexts_owner
  }
  property {
    name  = "reposRemoveStatusCheckContexts_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_status_check_contexts_repo
  }
  property {
    name  = "reposRemoveStatusCheckContexts_reposSetStatusCheckContextsRequest_ReposSetStatusCheckContextsRequest_contexts"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_status_check_contexts_repos_set_status_check_contexts_request_repos_set_status_check_contexts_request_contexts
  }
  property {
    name  = "reposRemoveStatusCheckProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_status_check_protection_branch
  }
  property {
    name  = "reposRemoveStatusCheckProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_status_check_protection_owner
  }
  property {
    name  = "reposRemoveStatusCheckProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_status_check_protection_repo
  }
  property {
    name  = "reposRemoveTeamAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_team_access_restrictions_branch
  }
  property {
    name  = "reposRemoveTeamAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_team_access_restrictions_owner
  }
  property {
    name  = "reposRemoveTeamAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_team_access_restrictions_repo
  }
  property {
    name  = "reposRemoveTeamAccessRestrictions_reposSetTeamAccessRestrictionsRequest_ReposSetTeamAccessRestrictionsRequest_teams"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_team_access_restrictions_repos_set_team_access_restrictions_request_repos_set_team_access_restrictions_request_teams
  }
  property {
    name  = "reposRemoveUserAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_user_access_restrictions_branch
  }
  property {
    name  = "reposRemoveUserAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_user_access_restrictions_owner
  }
  property {
    name  = "reposRemoveUserAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_user_access_restrictions_repo
  }
  property {
    name  = "reposRemoveUserAccessRestrictions_reposSetUserAccessRestrictionsRequest_ReposSetUserAccessRestrictionsRequest_users"
    type  = "string"
    value = var.connector-oai-github_property_repos_remove_user_access_restrictions_repos_set_user_access_restrictions_request_repos_set_user_access_restrictions_request_users
  }
  property {
    name  = "reposRenameBranch_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_rename_branch_branch
  }
  property {
    name  = "reposRenameBranch_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_rename_branch_owner
  }
  property {
    name  = "reposRenameBranch_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_rename_branch_repo
  }
  property {
    name  = "reposRenameBranch_reposRenameBranchRequest_ReposRenameBranchRequest_new_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_rename_branch_repos_rename_branch_request_repos_rename_branch_request_new_name
  }
  property {
    name  = "reposReplaceAllTopics_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_replace_all_topics_owner
  }
  property {
    name  = "reposReplaceAllTopics_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_replace_all_topics_repo
  }
  property {
    name  = "reposReplaceAllTopics_reposReplaceAllTopicsRequest_ReposReplaceAllTopicsRequest_names"
    type  = "string"
    value = var.connector-oai-github_property_repos_replace_all_topics_repos_replace_all_topics_request_repos_replace_all_topics_request_names
  }
  property {
    name  = "reposRequestPagesBuild_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_request_pages_build_owner
  }
  property {
    name  = "reposRequestPagesBuild_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_request_pages_build_repo
  }
  property {
    name  = "reposSetAdminBranchProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_admin_branch_protection_branch
  }
  property {
    name  = "reposSetAdminBranchProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_admin_branch_protection_owner
  }
  property {
    name  = "reposSetAdminBranchProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_admin_branch_protection_repo
  }
  property {
    name  = "reposSetAppAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_app_access_restrictions_branch
  }
  property {
    name  = "reposSetAppAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_app_access_restrictions_owner
  }
  property {
    name  = "reposSetAppAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_app_access_restrictions_repo
  }
  property {
    name  = "reposSetAppAccessRestrictions_reposSetAppAccessRestrictionsRequest_ReposSetAppAccessRestrictionsRequest_apps"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_app_access_restrictions_repos_set_app_access_restrictions_request_repos_set_app_access_restrictions_request_apps
  }
  property {
    name  = "reposSetStatusCheckContexts_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_status_check_contexts_branch
  }
  property {
    name  = "reposSetStatusCheckContexts_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_status_check_contexts_owner
  }
  property {
    name  = "reposSetStatusCheckContexts_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_status_check_contexts_repo
  }
  property {
    name  = "reposSetStatusCheckContexts_reposSetStatusCheckContextsRequest_ReposSetStatusCheckContextsRequest_contexts"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_status_check_contexts_repos_set_status_check_contexts_request_repos_set_status_check_contexts_request_contexts
  }
  property {
    name  = "reposSetTeamAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_team_access_restrictions_branch
  }
  property {
    name  = "reposSetTeamAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_team_access_restrictions_owner
  }
  property {
    name  = "reposSetTeamAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_team_access_restrictions_repo
  }
  property {
    name  = "reposSetTeamAccessRestrictions_reposSetTeamAccessRestrictionsRequest_ReposSetTeamAccessRestrictionsRequest_teams"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_team_access_restrictions_repos_set_team_access_restrictions_request_repos_set_team_access_restrictions_request_teams
  }
  property {
    name  = "reposSetUserAccessRestrictions_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_user_access_restrictions_branch
  }
  property {
    name  = "reposSetUserAccessRestrictions_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_user_access_restrictions_owner
  }
  property {
    name  = "reposSetUserAccessRestrictions_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_user_access_restrictions_repo
  }
  property {
    name  = "reposSetUserAccessRestrictions_reposSetUserAccessRestrictionsRequest_ReposSetUserAccessRestrictionsRequest_users"
    type  = "string"
    value = var.connector-oai-github_property_repos_set_user_access_restrictions_repos_set_user_access_restrictions_request_repos_set_user_access_restrictions_request_users
  }
  property {
    name  = "reposTestPushWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_test_push_webhook_hook_id
  }
  property {
    name  = "reposTestPushWebhook_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_test_push_webhook_owner
  }
  property {
    name  = "reposTestPushWebhook_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_test_push_webhook_repo
  }
  property {
    name  = "reposTransfer_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_transfer_owner
  }
  property {
    name  = "reposTransfer_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_transfer_repo
  }
  property {
    name  = "reposTransfer_reposTransferRequest_ReposTransferRequest_new_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_transfer_repos_transfer_request_repos_transfer_request_new_owner
  }
  property {
    name  = "reposTransfer_reposTransferRequest_ReposTransferRequest_team_ids"
    type  = "string"
    value = var.connector-oai-github_property_repos_transfer_repos_transfer_request_repos_transfer_request_team_ids
  }
  property {
    name  = "reposUpdateBranchProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_branch
  }
  property {
    name  = "reposUpdateBranchProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_owner
  }
  property {
    name  = "reposUpdateBranchProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repo
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviewsDismissalRestrictions_apps"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_dismissal_restrictions_apps
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviewsDismissalRestrictions_teams"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_dismissal_restrictions_teams
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviewsDismissalRestrictions_users"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_dismissal_restrictions_users
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviews_dismiss_stale_reviews"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_dismiss_stale_reviews
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviews_require_code_owner_reviews"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_require_code_owner_reviews
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviews_required_approving_review_count"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_required_approving_review_count
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredStatusChecks_checks"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_status_checks_checks
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredStatusChecks_contexts"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_status_checks_contexts
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRequiredStatusChecks_strict"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_status_checks_strict
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRestrictions_apps"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_restrictions_apps
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRestrictions_teams"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_restrictions_teams
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequestRestrictions_users"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_restrictions_users
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequest_allow_deletions"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_allow_deletions
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequest_allow_force_pushes"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_allow_force_pushes
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequest_block_creations"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_block_creations
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequest_contexts"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_contexts
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequest_enforce_admins"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_enforce_admins
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequest_required_conversation_resolution"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_conversation_resolution
  }
  property {
    name  = "reposUpdateBranchProtection_reposUpdateBranchProtectionRequest_ReposUpdateBranchProtectionRequest_required_linear_history"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_branch_protection_repos_update_branch_protection_request_repos_update_branch_protection_request_required_linear_history
  }
  property {
    name  = "reposUpdateCommitComment_comment_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_commit_comment_comment_id
  }
  property {
    name  = "reposUpdateCommitComment_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_commit_comment_owner
  }
  property {
    name  = "reposUpdateCommitComment_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_commit_comment_repo
  }
  property {
    name  = "reposUpdateCommitComment_reposUpdateCommitCommentRequest_ReposUpdateCommitCommentRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_commit_comment_repos_update_commit_comment_request_repos_update_commit_comment_request_body
  }
  property {
    name  = "reposUpdateInformationAboutPagesSite_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_information_about_pages_site_owner
  }
  property {
    name  = "reposUpdateInformationAboutPagesSite_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_information_about_pages_site_repo
  }
  property {
    name  = "reposUpdateInformationAboutPagesSite_reposUpdateInformationAboutPagesSiteRequest_ReposUpdateInformationAboutPagesSiteRequestSource_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_information_about_pages_site_repos_update_information_about_pages_site_request_repos_update_information_about_pages_site_request_source_branch
  }
  property {
    name  = "reposUpdateInformationAboutPagesSite_reposUpdateInformationAboutPagesSiteRequest_ReposUpdateInformationAboutPagesSiteRequestSource_path"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_information_about_pages_site_repos_update_information_about_pages_site_request_repos_update_information_about_pages_site_request_source_path
  }
  property {
    name  = "reposUpdateInformationAboutPagesSite_reposUpdateInformationAboutPagesSiteRequest_ReposUpdateInformationAboutPagesSiteRequest_cname"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_information_about_pages_site_repos_update_information_about_pages_site_request_repos_update_information_about_pages_site_request_cname
  }
  property {
    name  = "reposUpdateInformationAboutPagesSite_reposUpdateInformationAboutPagesSiteRequest_ReposUpdateInformationAboutPagesSiteRequest_https_enforced"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_information_about_pages_site_repos_update_information_about_pages_site_request_repos_update_information_about_pages_site_request_https_enforced
  }
  property {
    name  = "reposUpdateInformationAboutPagesSite_reposUpdateInformationAboutPagesSiteRequest_ReposUpdateInformationAboutPagesSiteRequest_public"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_information_about_pages_site_repos_update_information_about_pages_site_request_repos_update_information_about_pages_site_request_public
  }
  property {
    name  = "reposUpdateInvitation_invitation_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_invitation_invitation_id
  }
  property {
    name  = "reposUpdateInvitation_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_invitation_owner
  }
  property {
    name  = "reposUpdateInvitation_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_invitation_repo
  }
  property {
    name  = "reposUpdateInvitation_reposUpdateInvitationRequest_ReposUpdateInvitationRequest_permissions"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_invitation_repos_update_invitation_request_repos_update_invitation_request_permissions
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_branch
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_owner
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_repo
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_reposUpdatePullRequestReviewProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviewsDismissalRestrictions_apps"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_repos_update_pull_request_review_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_dismissal_restrictions_apps
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_reposUpdatePullRequestReviewProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviewsDismissalRestrictions_teams"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_repos_update_pull_request_review_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_dismissal_restrictions_teams
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_reposUpdatePullRequestReviewProtectionRequest_ReposUpdateBranchProtectionRequestRequiredPullRequestReviewsDismissalRestrictions_users"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_repos_update_pull_request_review_protection_request_repos_update_branch_protection_request_required_pull_request_reviews_dismissal_restrictions_users
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_reposUpdatePullRequestReviewProtectionRequest_ReposUpdatePullRequestReviewProtectionRequest_dismiss_stale_reviews"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_repos_update_pull_request_review_protection_request_repos_update_pull_request_review_protection_request_dismiss_stale_reviews
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_reposUpdatePullRequestReviewProtectionRequest_ReposUpdatePullRequestReviewProtectionRequest_require_code_owner_reviews"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_repos_update_pull_request_review_protection_request_repos_update_pull_request_review_protection_request_require_code_owner_reviews
  }
  property {
    name  = "reposUpdatePullRequestReviewProtection_reposUpdatePullRequestReviewProtectionRequest_ReposUpdatePullRequestReviewProtectionRequest_required_approving_review_count"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_pull_request_review_protection_repos_update_pull_request_review_protection_request_repos_update_pull_request_review_protection_request_required_approving_review_count
  }
  property {
    name  = "reposUpdateReleaseAsset_asset_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_asset_asset_id
  }
  property {
    name  = "reposUpdateReleaseAsset_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_asset_owner
  }
  property {
    name  = "reposUpdateReleaseAsset_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_asset_repo
  }
  property {
    name  = "reposUpdateReleaseAsset_reposUpdateReleaseAssetRequest_ReposUpdateReleaseAssetRequest_label"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_asset_repos_update_release_asset_request_repos_update_release_asset_request_label
  }
  property {
    name  = "reposUpdateReleaseAsset_reposUpdateReleaseAssetRequest_ReposUpdateReleaseAssetRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_asset_repos_update_release_asset_request_repos_update_release_asset_request_name
  }
  property {
    name  = "reposUpdateReleaseAsset_reposUpdateReleaseAssetRequest_ReposUpdateReleaseAssetRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_asset_repos_update_release_asset_request_repos_update_release_asset_request_state
  }
  property {
    name  = "reposUpdateRelease_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_owner
  }
  property {
    name  = "reposUpdateRelease_release_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_release_id
  }
  property {
    name  = "reposUpdateRelease_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_repo
  }
  property {
    name  = "reposUpdateRelease_reposUpdateReleaseRequest_ReposUpdateReleaseRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_repos_update_release_request_repos_update_release_request_body
  }
  property {
    name  = "reposUpdateRelease_reposUpdateReleaseRequest_ReposUpdateReleaseRequest_draft"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_repos_update_release_request_repos_update_release_request_draft
  }
  property {
    name  = "reposUpdateRelease_reposUpdateReleaseRequest_ReposUpdateReleaseRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_repos_update_release_request_repos_update_release_request_name
  }
  property {
    name  = "reposUpdateRelease_reposUpdateReleaseRequest_ReposUpdateReleaseRequest_prerelease"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_repos_update_release_request_repos_update_release_request_prerelease
  }
  property {
    name  = "reposUpdateRelease_reposUpdateReleaseRequest_ReposUpdateReleaseRequest_tag_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_repos_update_release_request_repos_update_release_request_tag_name
  }
  property {
    name  = "reposUpdateRelease_reposUpdateReleaseRequest_ReposUpdateReleaseRequest_target_commitish"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_release_repos_update_release_request_repos_update_release_request_target_commitish
  }
  property {
    name  = "reposUpdateStatusCheckProtection_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_status_check_protection_branch
  }
  property {
    name  = "reposUpdateStatusCheckProtection_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_status_check_protection_owner
  }
  property {
    name  = "reposUpdateStatusCheckProtection_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_status_check_protection_repo
  }
  property {
    name  = "reposUpdateStatusCheckProtection_reposUpdateStatusCheckProtectionRequest_ReposUpdateStatusCheckProtectionRequest_contexts"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_status_check_protection_repos_update_status_check_protection_request_repos_update_status_check_protection_request_contexts
  }
  property {
    name  = "reposUpdateStatusCheckProtection_reposUpdateStatusCheckProtectionRequest_ReposUpdateStatusCheckProtectionRequest_strict"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_status_check_protection_repos_update_status_check_protection_request_repos_update_status_check_protection_request_strict
  }
  property {
    name  = "reposUpdateWebhookConfigForRepo_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_config_for_repo_hook_id
  }
  property {
    name  = "reposUpdateWebhookConfigForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_config_for_repo_owner
  }
  property {
    name  = "reposUpdateWebhookConfigForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_config_for_repo_repo
  }
  property {
    name  = "reposUpdateWebhookConfigForRepo_reposUpdateWebhookConfigForRepoRequest_ReposUpdateWebhookConfigForRepoRequest_content_type"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_config_for_repo_repos_update_webhook_config_for_repo_request_repos_update_webhook_config_for_repo_request_content_type
  }
  property {
    name  = "reposUpdateWebhookConfigForRepo_reposUpdateWebhookConfigForRepoRequest_ReposUpdateWebhookConfigForRepoRequest_secret"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_config_for_repo_repos_update_webhook_config_for_repo_request_repos_update_webhook_config_for_repo_request_secret
  }
  property {
    name  = "reposUpdateWebhookConfigForRepo_reposUpdateWebhookConfigForRepoRequest_ReposUpdateWebhookConfigForRepoRequest_url"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_config_for_repo_repos_update_webhook_config_for_repo_request_repos_update_webhook_config_for_repo_request_url
  }
  property {
    name  = "reposUpdateWebhook_hook_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_hook_id
  }
  property {
    name  = "reposUpdateWebhook_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_owner
  }
  property {
    name  = "reposUpdateWebhook_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repo
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequestConfig_address"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_config_address
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequestConfig_content_type"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_config_content_type
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequestConfig_room"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_config_room
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequestConfig_secret"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_config_secret
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequestConfig_url"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_config_url
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequest_active"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_active
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequest_add_events"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_add_events
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequest_events"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_events
  }
  property {
    name  = "reposUpdateWebhook_reposUpdateWebhookRequest_ReposUpdateWebhookRequest_remove_events"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_webhook_repos_update_webhook_request_repos_update_webhook_request_remove_events
  }
  property {
    name  = "reposUpdate_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_owner
  }
  property {
    name  = "reposUpdate_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repo
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_allow_forking"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_allow_forking
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_allow_merge_commit"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_allow_merge_commit
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_allow_rebase_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_allow_rebase_merge
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_allow_squash_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_allow_squash_merge
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_allow_update_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_allow_update_branch
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_archived"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_archived
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_default_branch"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_default_branch
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_delete_branch_on_merge"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_delete_branch_on_merge
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_description
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_has_issues"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_has_issues
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_has_projects"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_has_projects
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_has_wiki"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_has_wiki
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_homepage"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_homepage
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_is_template"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_is_template
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_name
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_private"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_private
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_use_squash_pr_title_as_default"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_use_squash_pr_title_as_default
  }
  property {
    name  = "reposUpdate_reposUpdateRequest_ReposUpdateRequest_visibility"
    type  = "string"
    value = var.connector-oai-github_property_repos_update_repos_update_request_repos_update_request_visibility
  }
  property {
    name  = "reposUploadReleaseAsset_body"
    type  = "string"
    value = var.connector-oai-github_property_repos_upload_release_asset_body
  }
  property {
    name  = "reposUploadReleaseAsset_label"
    type  = "string"
    value = var.connector-oai-github_property_repos_upload_release_asset_label
  }
  property {
    name  = "reposUploadReleaseAsset_name"
    type  = "string"
    value = var.connector-oai-github_property_repos_upload_release_asset_name
  }
  property {
    name  = "reposUploadReleaseAsset_owner"
    type  = "string"
    value = var.connector-oai-github_property_repos_upload_release_asset_owner
  }
  property {
    name  = "reposUploadReleaseAsset_release_id"
    type  = "string"
    value = var.connector-oai-github_property_repos_upload_release_asset_release_id
  }
  property {
    name  = "reposUploadReleaseAsset_repo"
    type  = "string"
    value = var.connector-oai-github_property_repos_upload_release_asset_repo
  }
  property {
    name  = "searchCode_order"
    type  = "string"
    value = var.connector-oai-github_property_search_code_order
  }
  property {
    name  = "searchCode_page"
    type  = "string"
    value = var.connector-oai-github_property_search_code_page
  }
  property {
    name  = "searchCode_per_page"
    type  = "string"
    value = var.connector-oai-github_property_search_code_per_page
  }
  property {
    name  = "searchCode_q"
    type  = "string"
    value = var.connector-oai-github_property_search_code_q
  }
  property {
    name  = "searchCode_sort"
    type  = "string"
    value = var.connector-oai-github_property_search_code_sort
  }
  property {
    name  = "searchCommits_order"
    type  = "string"
    value = var.connector-oai-github_property_search_commits_order
  }
  property {
    name  = "searchCommits_page"
    type  = "string"
    value = var.connector-oai-github_property_search_commits_page
  }
  property {
    name  = "searchCommits_per_page"
    type  = "string"
    value = var.connector-oai-github_property_search_commits_per_page
  }
  property {
    name  = "searchCommits_q"
    type  = "string"
    value = var.connector-oai-github_property_search_commits_q
  }
  property {
    name  = "searchCommits_sort"
    type  = "string"
    value = var.connector-oai-github_property_search_commits_sort
  }
  property {
    name  = "searchIssuesAndPullRequests_order"
    type  = "string"
    value = var.connector-oai-github_property_search_issues_and_pull_requests_order
  }
  property {
    name  = "searchIssuesAndPullRequests_page"
    type  = "string"
    value = var.connector-oai-github_property_search_issues_and_pull_requests_page
  }
  property {
    name  = "searchIssuesAndPullRequests_per_page"
    type  = "string"
    value = var.connector-oai-github_property_search_issues_and_pull_requests_per_page
  }
  property {
    name  = "searchIssuesAndPullRequests_q"
    type  = "string"
    value = var.connector-oai-github_property_search_issues_and_pull_requests_q
  }
  property {
    name  = "searchIssuesAndPullRequests_sort"
    type  = "string"
    value = var.connector-oai-github_property_search_issues_and_pull_requests_sort
  }
  property {
    name  = "searchLabels_order"
    type  = "string"
    value = var.connector-oai-github_property_search_labels_order
  }
  property {
    name  = "searchLabels_page"
    type  = "string"
    value = var.connector-oai-github_property_search_labels_page
  }
  property {
    name  = "searchLabels_per_page"
    type  = "string"
    value = var.connector-oai-github_property_search_labels_per_page
  }
  property {
    name  = "searchLabels_q"
    type  = "string"
    value = var.connector-oai-github_property_search_labels_q
  }
  property {
    name  = "searchLabels_repository_id"
    type  = "string"
    value = var.connector-oai-github_property_search_labels_repository_id
  }
  property {
    name  = "searchLabels_sort"
    type  = "string"
    value = var.connector-oai-github_property_search_labels_sort
  }
  property {
    name  = "searchRepos_order"
    type  = "string"
    value = var.connector-oai-github_property_search_repos_order
  }
  property {
    name  = "searchRepos_page"
    type  = "string"
    value = var.connector-oai-github_property_search_repos_page
  }
  property {
    name  = "searchRepos_per_page"
    type  = "string"
    value = var.connector-oai-github_property_search_repos_per_page
  }
  property {
    name  = "searchRepos_q"
    type  = "string"
    value = var.connector-oai-github_property_search_repos_q
  }
  property {
    name  = "searchRepos_sort"
    type  = "string"
    value = var.connector-oai-github_property_search_repos_sort
  }
  property {
    name  = "searchTopics_page"
    type  = "string"
    value = var.connector-oai-github_property_search_topics_page
  }
  property {
    name  = "searchTopics_per_page"
    type  = "string"
    value = var.connector-oai-github_property_search_topics_per_page
  }
  property {
    name  = "searchTopics_q"
    type  = "string"
    value = var.connector-oai-github_property_search_topics_q
  }
  property {
    name  = "searchUsers_order"
    type  = "string"
    value = var.connector-oai-github_property_search_users_order
  }
  property {
    name  = "searchUsers_page"
    type  = "string"
    value = var.connector-oai-github_property_search_users_page
  }
  property {
    name  = "searchUsers_per_page"
    type  = "string"
    value = var.connector-oai-github_property_search_users_per_page
  }
  property {
    name  = "searchUsers_q"
    type  = "string"
    value = var.connector-oai-github_property_search_users_q
  }
  property {
    name  = "searchUsers_sort"
    type  = "string"
    value = var.connector-oai-github_property_search_users_sort
  }
  property {
    name  = "secretScanningGetAlert_alert_number"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_get_alert_alert_number
  }
  property {
    name  = "secretScanningGetAlert_owner"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_get_alert_owner
  }
  property {
    name  = "secretScanningGetAlert_repo"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_get_alert_repo
  }
  property {
    name  = "secretScanningListAlertsForRepo_owner"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_list_alerts_for_repo_owner
  }
  property {
    name  = "secretScanningListAlertsForRepo_page"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_list_alerts_for_repo_page
  }
  property {
    name  = "secretScanningListAlertsForRepo_per_page"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_list_alerts_for_repo_per_page
  }
  property {
    name  = "secretScanningListAlertsForRepo_repo"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_list_alerts_for_repo_repo
  }
  property {
    name  = "secretScanningListAlertsForRepo_resolution"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_list_alerts_for_repo_resolution
  }
  property {
    name  = "secretScanningListAlertsForRepo_secret_type"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_list_alerts_for_repo_secret_type
  }
  property {
    name  = "secretScanningListAlertsForRepo_state"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_list_alerts_for_repo_state
  }
  property {
    name  = "secretScanningUpdateAlert_alert_number"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_update_alert_alert_number
  }
  property {
    name  = "secretScanningUpdateAlert_owner"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_update_alert_owner
  }
  property {
    name  = "secretScanningUpdateAlert_repo"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_update_alert_repo
  }
  property {
    name  = "secretScanningUpdateAlert_secretScanningUpdateAlertRequest_SecretScanningUpdateAlertRequest_resolution"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_update_alert_secret_scanning_update_alert_request_secret_scanning_update_alert_request_resolution
  }
  property {
    name  = "secretScanningUpdateAlert_secretScanningUpdateAlertRequest_SecretScanningUpdateAlertRequest_state"
    type  = "string"
    value = var.connector-oai-github_property_secret_scanning_update_alert_secret_scanning_update_alert_request_secret_scanning_update_alert_request_state
  }
  property {
    name  = "teamsAddMemberLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_member_legacy_team_id
  }
  property {
    name  = "teamsAddMemberLegacy_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_member_legacy_username
  }
  property {
    name  = "teamsAddOrUpdateMembershipForUserInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_membership_for_user_in_org_org
  }
  property {
    name  = "teamsAddOrUpdateMembershipForUserInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_membership_for_user_in_org_team_slug
  }
  property {
    name  = "teamsAddOrUpdateMembershipForUserInOrg_teamsAddOrUpdateMembershipForUserInOrgRequest_TeamsAddOrUpdateMembershipForUserInOrgRequest_role"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_membership_for_user_in_org_teams_add_or_update_membership_for_user_in_org_request_teams_add_or_update_membership_for_user_in_org_request_role
  }
  property {
    name  = "teamsAddOrUpdateMembershipForUserInOrg_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_membership_for_user_in_org_username
  }
  property {
    name  = "teamsAddOrUpdateMembershipForUserLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_membership_for_user_legacy_team_id
  }
  property {
    name  = "teamsAddOrUpdateMembershipForUserLegacy_teamsAddOrUpdateMembershipForUserInOrgRequest_TeamsAddOrUpdateMembershipForUserInOrgRequest_role"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_membership_for_user_legacy_teams_add_or_update_membership_for_user_in_org_request_teams_add_or_update_membership_for_user_in_org_request_role
  }
  property {
    name  = "teamsAddOrUpdateMembershipForUserLegacy_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_membership_for_user_legacy_username
  }
  property {
    name  = "teamsAddOrUpdateProjectPermissionsInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_project_permissions_in_org_org
  }
  property {
    name  = "teamsAddOrUpdateProjectPermissionsInOrg_project_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_project_permissions_in_org_project_id
  }
  property {
    name  = "teamsAddOrUpdateProjectPermissionsInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_project_permissions_in_org_team_slug
  }
  property {
    name  = "teamsAddOrUpdateProjectPermissionsInOrg_teamsAddOrUpdateProjectPermissionsInOrgRequest_TeamsAddOrUpdateProjectPermissionsInOrgRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_project_permissions_in_org_teams_add_or_update_project_permissions_in_org_request_teams_add_or_update_project_permissions_in_org_request_permission
  }
  property {
    name  = "teamsAddOrUpdateProjectPermissionsLegacy_project_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_project_permissions_legacy_project_id
  }
  property {
    name  = "teamsAddOrUpdateProjectPermissionsLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_project_permissions_legacy_team_id
  }
  property {
    name  = "teamsAddOrUpdateProjectPermissionsLegacy_teamsAddOrUpdateProjectPermissionsLegacyRequest_TeamsAddOrUpdateProjectPermissionsLegacyRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_project_permissions_legacy_teams_add_or_update_project_permissions_legacy_request_teams_add_or_update_project_permissions_legacy_request_permission
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_in_org_org
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsInOrg_owner"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_in_org_owner
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsInOrg_repo"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_in_org_repo
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_in_org_team_slug
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsInOrg_teamsAddOrUpdateRepoPermissionsInOrgRequest_TeamsAddOrUpdateRepoPermissionsInOrgRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_in_org_teams_add_or_update_repo_permissions_in_org_request_teams_add_or_update_repo_permissions_in_org_request_permission
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsLegacy_owner"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_legacy_owner
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsLegacy_repo"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_legacy_repo
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_legacy_team_id
  }
  property {
    name  = "teamsAddOrUpdateRepoPermissionsLegacy_teamsAddOrUpdateRepoPermissionsLegacyRequest_TeamsAddOrUpdateRepoPermissionsLegacyRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_teams_add_or_update_repo_permissions_legacy_teams_add_or_update_repo_permissions_legacy_request_teams_add_or_update_repo_permissions_legacy_request_permission
  }
  property {
    name  = "teamsCheckPermissionsForProjectInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_project_in_org_org
  }
  property {
    name  = "teamsCheckPermissionsForProjectInOrg_project_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_project_in_org_project_id
  }
  property {
    name  = "teamsCheckPermissionsForProjectInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_project_in_org_team_slug
  }
  property {
    name  = "teamsCheckPermissionsForProjectLegacy_project_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_project_legacy_project_id
  }
  property {
    name  = "teamsCheckPermissionsForProjectLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_project_legacy_team_id
  }
  property {
    name  = "teamsCheckPermissionsForRepoInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_repo_in_org_org
  }
  property {
    name  = "teamsCheckPermissionsForRepoInOrg_owner"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_repo_in_org_owner
  }
  property {
    name  = "teamsCheckPermissionsForRepoInOrg_repo"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_repo_in_org_repo
  }
  property {
    name  = "teamsCheckPermissionsForRepoInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_repo_in_org_team_slug
  }
  property {
    name  = "teamsCheckPermissionsForRepoLegacy_owner"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_repo_legacy_owner
  }
  property {
    name  = "teamsCheckPermissionsForRepoLegacy_repo"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_repo_legacy_repo
  }
  property {
    name  = "teamsCheckPermissionsForRepoLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_check_permissions_for_repo_legacy_team_id
  }
  property {
    name  = "teamsCreateDiscussionCommentInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_comment_in_org_discussion_number
  }
  property {
    name  = "teamsCreateDiscussionCommentInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_comment_in_org_org
  }
  property {
    name  = "teamsCreateDiscussionCommentInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_comment_in_org_team_slug
  }
  property {
    name  = "teamsCreateDiscussionCommentInOrg_teamsCreateDiscussionCommentInOrgRequest_TeamsCreateDiscussionCommentInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_comment_in_org_teams_create_discussion_comment_in_org_request_teams_create_discussion_comment_in_org_request_body
  }
  property {
    name  = "teamsCreateDiscussionCommentLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_comment_legacy_discussion_number
  }
  property {
    name  = "teamsCreateDiscussionCommentLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_comment_legacy_team_id
  }
  property {
    name  = "teamsCreateDiscussionCommentLegacy_teamsCreateDiscussionCommentInOrgRequest_TeamsCreateDiscussionCommentInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_comment_legacy_teams_create_discussion_comment_in_org_request_teams_create_discussion_comment_in_org_request_body
  }
  property {
    name  = "teamsCreateDiscussionInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_in_org_org
  }
  property {
    name  = "teamsCreateDiscussionInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_in_org_team_slug
  }
  property {
    name  = "teamsCreateDiscussionInOrg_teamsCreateDiscussionInOrgRequest_TeamsCreateDiscussionInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_in_org_teams_create_discussion_in_org_request_teams_create_discussion_in_org_request_body
  }
  property {
    name  = "teamsCreateDiscussionInOrg_teamsCreateDiscussionInOrgRequest_TeamsCreateDiscussionInOrgRequest_private"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_in_org_teams_create_discussion_in_org_request_teams_create_discussion_in_org_request_private
  }
  property {
    name  = "teamsCreateDiscussionInOrg_teamsCreateDiscussionInOrgRequest_TeamsCreateDiscussionInOrgRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_in_org_teams_create_discussion_in_org_request_teams_create_discussion_in_org_request_title
  }
  property {
    name  = "teamsCreateDiscussionLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_legacy_team_id
  }
  property {
    name  = "teamsCreateDiscussionLegacy_teamsCreateDiscussionInOrgRequest_TeamsCreateDiscussionInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_legacy_teams_create_discussion_in_org_request_teams_create_discussion_in_org_request_body
  }
  property {
    name  = "teamsCreateDiscussionLegacy_teamsCreateDiscussionInOrgRequest_TeamsCreateDiscussionInOrgRequest_private"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_legacy_teams_create_discussion_in_org_request_teams_create_discussion_in_org_request_private
  }
  property {
    name  = "teamsCreateDiscussionLegacy_teamsCreateDiscussionInOrgRequest_TeamsCreateDiscussionInOrgRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_discussion_legacy_teams_create_discussion_in_org_request_teams_create_discussion_in_org_request_title
  }
  property {
    name  = "teamsCreate_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_org
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_description
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_ldap_dn"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_ldap_dn
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_maintainers"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_maintainers
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_name
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_parent_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_parent_team_id
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_permission
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_privacy"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_privacy
  }
  property {
    name  = "teamsCreate_teamsCreateRequest_TeamsCreateRequest_repo_names"
    type  = "string"
    value = var.connector-oai-github_property_teams_create_teams_create_request_teams_create_request_repo_names
  }
  property {
    name  = "teamsDeleteDiscussionCommentInOrg_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_comment_in_org_comment_number
  }
  property {
    name  = "teamsDeleteDiscussionCommentInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_comment_in_org_discussion_number
  }
  property {
    name  = "teamsDeleteDiscussionCommentInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_comment_in_org_org
  }
  property {
    name  = "teamsDeleteDiscussionCommentInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_comment_in_org_team_slug
  }
  property {
    name  = "teamsDeleteDiscussionCommentLegacy_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_comment_legacy_comment_number
  }
  property {
    name  = "teamsDeleteDiscussionCommentLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_comment_legacy_discussion_number
  }
  property {
    name  = "teamsDeleteDiscussionCommentLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_comment_legacy_team_id
  }
  property {
    name  = "teamsDeleteDiscussionInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_in_org_discussion_number
  }
  property {
    name  = "teamsDeleteDiscussionInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_in_org_org
  }
  property {
    name  = "teamsDeleteDiscussionInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_in_org_team_slug
  }
  property {
    name  = "teamsDeleteDiscussionLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_legacy_discussion_number
  }
  property {
    name  = "teamsDeleteDiscussionLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_discussion_legacy_team_id
  }
  property {
    name  = "teamsDeleteInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_in_org_org
  }
  property {
    name  = "teamsDeleteInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_in_org_team_slug
  }
  property {
    name  = "teamsDeleteLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_delete_legacy_team_id
  }
  property {
    name  = "teamsGetByName_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_by_name_org
  }
  property {
    name  = "teamsGetByName_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_by_name_team_slug
  }
  property {
    name  = "teamsGetDiscussionCommentInOrg_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_comment_in_org_comment_number
  }
  property {
    name  = "teamsGetDiscussionCommentInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_comment_in_org_discussion_number
  }
  property {
    name  = "teamsGetDiscussionCommentInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_comment_in_org_org
  }
  property {
    name  = "teamsGetDiscussionCommentInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_comment_in_org_team_slug
  }
  property {
    name  = "teamsGetDiscussionCommentLegacy_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_comment_legacy_comment_number
  }
  property {
    name  = "teamsGetDiscussionCommentLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_comment_legacy_discussion_number
  }
  property {
    name  = "teamsGetDiscussionCommentLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_comment_legacy_team_id
  }
  property {
    name  = "teamsGetDiscussionInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_in_org_discussion_number
  }
  property {
    name  = "teamsGetDiscussionInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_in_org_org
  }
  property {
    name  = "teamsGetDiscussionInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_in_org_team_slug
  }
  property {
    name  = "teamsGetDiscussionLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_legacy_discussion_number
  }
  property {
    name  = "teamsGetDiscussionLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_discussion_legacy_team_id
  }
  property {
    name  = "teamsGetLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_legacy_team_id
  }
  property {
    name  = "teamsGetMemberLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_member_legacy_team_id
  }
  property {
    name  = "teamsGetMemberLegacy_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_member_legacy_username
  }
  property {
    name  = "teamsGetMembershipForUserInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_membership_for_user_in_org_org
  }
  property {
    name  = "teamsGetMembershipForUserInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_membership_for_user_in_org_team_slug
  }
  property {
    name  = "teamsGetMembershipForUserInOrg_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_membership_for_user_in_org_username
  }
  property {
    name  = "teamsGetMembershipForUserLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_membership_for_user_legacy_team_id
  }
  property {
    name  = "teamsGetMembershipForUserLegacy_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_get_membership_for_user_legacy_username
  }
  property {
    name  = "teamsListChildInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_child_in_org_org
  }
  property {
    name  = "teamsListChildInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_child_in_org_page
  }
  property {
    name  = "teamsListChildInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_child_in_org_per_page
  }
  property {
    name  = "teamsListChildInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_child_in_org_team_slug
  }
  property {
    name  = "teamsListChildLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_child_legacy_page
  }
  property {
    name  = "teamsListChildLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_child_legacy_per_page
  }
  property {
    name  = "teamsListChildLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_child_legacy_team_id
  }
  property {
    name  = "teamsListDiscussionCommentsInOrg_direction"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_in_org_direction
  }
  property {
    name  = "teamsListDiscussionCommentsInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_in_org_discussion_number
  }
  property {
    name  = "teamsListDiscussionCommentsInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_in_org_org
  }
  property {
    name  = "teamsListDiscussionCommentsInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_in_org_page
  }
  property {
    name  = "teamsListDiscussionCommentsInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_in_org_per_page
  }
  property {
    name  = "teamsListDiscussionCommentsInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_in_org_team_slug
  }
  property {
    name  = "teamsListDiscussionCommentsLegacy_direction"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_legacy_direction
  }
  property {
    name  = "teamsListDiscussionCommentsLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_legacy_discussion_number
  }
  property {
    name  = "teamsListDiscussionCommentsLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_legacy_page
  }
  property {
    name  = "teamsListDiscussionCommentsLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_legacy_per_page
  }
  property {
    name  = "teamsListDiscussionCommentsLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussion_comments_legacy_team_id
  }
  property {
    name  = "teamsListDiscussionsInOrg_direction"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_in_org_direction
  }
  property {
    name  = "teamsListDiscussionsInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_in_org_org
  }
  property {
    name  = "teamsListDiscussionsInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_in_org_page
  }
  property {
    name  = "teamsListDiscussionsInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_in_org_per_page
  }
  property {
    name  = "teamsListDiscussionsInOrg_pinned"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_in_org_pinned
  }
  property {
    name  = "teamsListDiscussionsInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_in_org_team_slug
  }
  property {
    name  = "teamsListDiscussionsLegacy_direction"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_legacy_direction
  }
  property {
    name  = "teamsListDiscussionsLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_legacy_page
  }
  property {
    name  = "teamsListDiscussionsLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_legacy_per_page
  }
  property {
    name  = "teamsListDiscussionsLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_discussions_legacy_team_id
  }
  property {
    name  = "teamsListForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_for_authenticated_user_page
  }
  property {
    name  = "teamsListForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_for_authenticated_user_per_page
  }
  property {
    name  = "teamsListMembersInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_in_org_org
  }
  property {
    name  = "teamsListMembersInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_in_org_page
  }
  property {
    name  = "teamsListMembersInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_in_org_per_page
  }
  property {
    name  = "teamsListMembersInOrg_role"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_in_org_role
  }
  property {
    name  = "teamsListMembersInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_in_org_team_slug
  }
  property {
    name  = "teamsListMembersLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_legacy_page
  }
  property {
    name  = "teamsListMembersLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_legacy_per_page
  }
  property {
    name  = "teamsListMembersLegacy_role"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_legacy_role
  }
  property {
    name  = "teamsListMembersLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_members_legacy_team_id
  }
  property {
    name  = "teamsListProjectsInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_projects_in_org_org
  }
  property {
    name  = "teamsListProjectsInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_projects_in_org_page
  }
  property {
    name  = "teamsListProjectsInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_projects_in_org_per_page
  }
  property {
    name  = "teamsListProjectsInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_projects_in_org_team_slug
  }
  property {
    name  = "teamsListProjectsLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_projects_legacy_page
  }
  property {
    name  = "teamsListProjectsLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_projects_legacy_per_page
  }
  property {
    name  = "teamsListProjectsLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_projects_legacy_team_id
  }
  property {
    name  = "teamsListReposInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_repos_in_org_org
  }
  property {
    name  = "teamsListReposInOrg_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_repos_in_org_page
  }
  property {
    name  = "teamsListReposInOrg_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_repos_in_org_per_page
  }
  property {
    name  = "teamsListReposInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_repos_in_org_team_slug
  }
  property {
    name  = "teamsListReposLegacy_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_repos_legacy_page
  }
  property {
    name  = "teamsListReposLegacy_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_repos_legacy_per_page
  }
  property {
    name  = "teamsListReposLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_repos_legacy_team_id
  }
  property {
    name  = "teamsList_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_org
  }
  property {
    name  = "teamsList_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_page
  }
  property {
    name  = "teamsList_per_page"
    type  = "string"
    value = var.connector-oai-github_property_teams_list_per_page
  }
  property {
    name  = "teamsRemoveMemberLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_member_legacy_team_id
  }
  property {
    name  = "teamsRemoveMemberLegacy_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_member_legacy_username
  }
  property {
    name  = "teamsRemoveMembershipForUserInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_membership_for_user_in_org_org
  }
  property {
    name  = "teamsRemoveMembershipForUserInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_membership_for_user_in_org_team_slug
  }
  property {
    name  = "teamsRemoveMembershipForUserInOrg_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_membership_for_user_in_org_username
  }
  property {
    name  = "teamsRemoveMembershipForUserLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_membership_for_user_legacy_team_id
  }
  property {
    name  = "teamsRemoveMembershipForUserLegacy_username"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_membership_for_user_legacy_username
  }
  property {
    name  = "teamsRemoveProjectInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_project_in_org_org
  }
  property {
    name  = "teamsRemoveProjectInOrg_project_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_project_in_org_project_id
  }
  property {
    name  = "teamsRemoveProjectInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_project_in_org_team_slug
  }
  property {
    name  = "teamsRemoveProjectLegacy_project_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_project_legacy_project_id
  }
  property {
    name  = "teamsRemoveProjectLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_project_legacy_team_id
  }
  property {
    name  = "teamsRemoveRepoInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_repo_in_org_org
  }
  property {
    name  = "teamsRemoveRepoInOrg_owner"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_repo_in_org_owner
  }
  property {
    name  = "teamsRemoveRepoInOrg_repo"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_repo_in_org_repo
  }
  property {
    name  = "teamsRemoveRepoInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_repo_in_org_team_slug
  }
  property {
    name  = "teamsRemoveRepoLegacy_owner"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_repo_legacy_owner
  }
  property {
    name  = "teamsRemoveRepoLegacy_repo"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_repo_legacy_repo
  }
  property {
    name  = "teamsRemoveRepoLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_remove_repo_legacy_team_id
  }
  property {
    name  = "teamsUpdateDiscussionCommentInOrg_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_in_org_comment_number
  }
  property {
    name  = "teamsUpdateDiscussionCommentInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_in_org_discussion_number
  }
  property {
    name  = "teamsUpdateDiscussionCommentInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_in_org_org
  }
  property {
    name  = "teamsUpdateDiscussionCommentInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_in_org_team_slug
  }
  property {
    name  = "teamsUpdateDiscussionCommentInOrg_teamsCreateDiscussionCommentInOrgRequest_TeamsCreateDiscussionCommentInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_in_org_teams_create_discussion_comment_in_org_request_teams_create_discussion_comment_in_org_request_body
  }
  property {
    name  = "teamsUpdateDiscussionCommentLegacy_comment_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_legacy_comment_number
  }
  property {
    name  = "teamsUpdateDiscussionCommentLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_legacy_discussion_number
  }
  property {
    name  = "teamsUpdateDiscussionCommentLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_legacy_team_id
  }
  property {
    name  = "teamsUpdateDiscussionCommentLegacy_teamsCreateDiscussionCommentInOrgRequest_TeamsCreateDiscussionCommentInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_comment_legacy_teams_create_discussion_comment_in_org_request_teams_create_discussion_comment_in_org_request_body
  }
  property {
    name  = "teamsUpdateDiscussionInOrg_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_in_org_discussion_number
  }
  property {
    name  = "teamsUpdateDiscussionInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_in_org_org
  }
  property {
    name  = "teamsUpdateDiscussionInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_in_org_team_slug
  }
  property {
    name  = "teamsUpdateDiscussionInOrg_teamsUpdateDiscussionInOrgRequest_TeamsUpdateDiscussionInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_in_org_teams_update_discussion_in_org_request_teams_update_discussion_in_org_request_body
  }
  property {
    name  = "teamsUpdateDiscussionInOrg_teamsUpdateDiscussionInOrgRequest_TeamsUpdateDiscussionInOrgRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_in_org_teams_update_discussion_in_org_request_teams_update_discussion_in_org_request_title
  }
  property {
    name  = "teamsUpdateDiscussionLegacy_discussion_number"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_legacy_discussion_number
  }
  property {
    name  = "teamsUpdateDiscussionLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_legacy_team_id
  }
  property {
    name  = "teamsUpdateDiscussionLegacy_teamsUpdateDiscussionInOrgRequest_TeamsUpdateDiscussionInOrgRequest_body"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_legacy_teams_update_discussion_in_org_request_teams_update_discussion_in_org_request_body
  }
  property {
    name  = "teamsUpdateDiscussionLegacy_teamsUpdateDiscussionInOrgRequest_TeamsUpdateDiscussionInOrgRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_discussion_legacy_teams_update_discussion_in_org_request_teams_update_discussion_in_org_request_title
  }
  property {
    name  = "teamsUpdateInOrg_org"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_in_org_org
  }
  property {
    name  = "teamsUpdateInOrg_team_slug"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_in_org_team_slug
  }
  property {
    name  = "teamsUpdateInOrg_teamsUpdateInOrgRequest_TeamsUpdateInOrgRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_in_org_teams_update_in_org_request_teams_update_in_org_request_description
  }
  property {
    name  = "teamsUpdateInOrg_teamsUpdateInOrgRequest_TeamsUpdateInOrgRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_in_org_teams_update_in_org_request_teams_update_in_org_request_name
  }
  property {
    name  = "teamsUpdateInOrg_teamsUpdateInOrgRequest_TeamsUpdateInOrgRequest_parent_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_in_org_teams_update_in_org_request_teams_update_in_org_request_parent_team_id
  }
  property {
    name  = "teamsUpdateInOrg_teamsUpdateInOrgRequest_TeamsUpdateInOrgRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_in_org_teams_update_in_org_request_teams_update_in_org_request_permission
  }
  property {
    name  = "teamsUpdateInOrg_teamsUpdateInOrgRequest_TeamsUpdateInOrgRequest_privacy"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_in_org_teams_update_in_org_request_teams_update_in_org_request_privacy
  }
  property {
    name  = "teamsUpdateLegacy_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_legacy_team_id
  }
  property {
    name  = "teamsUpdateLegacy_teamsUpdateLegacyRequest_TeamsUpdateLegacyRequest_description"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_legacy_teams_update_legacy_request_teams_update_legacy_request_description
  }
  property {
    name  = "teamsUpdateLegacy_teamsUpdateLegacyRequest_TeamsUpdateLegacyRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_legacy_teams_update_legacy_request_teams_update_legacy_request_name
  }
  property {
    name  = "teamsUpdateLegacy_teamsUpdateLegacyRequest_TeamsUpdateLegacyRequest_parent_team_id"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_legacy_teams_update_legacy_request_teams_update_legacy_request_parent_team_id
  }
  property {
    name  = "teamsUpdateLegacy_teamsUpdateLegacyRequest_TeamsUpdateLegacyRequest_permission"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_legacy_teams_update_legacy_request_teams_update_legacy_request_permission
  }
  property {
    name  = "teamsUpdateLegacy_teamsUpdateLegacyRequest_TeamsUpdateLegacyRequest_privacy"
    type  = "string"
    value = var.connector-oai-github_property_teams_update_legacy_teams_update_legacy_request_teams_update_legacy_request_privacy
  }
  property {
    name  = "usersAddEmailForAuthenticatedUser_usersAddEmailForAuthenticatedUserRequest_UsersAddEmailForAuthenticatedUserRequest_emails"
    type  = "string"
    value = var.connector-oai-github_property_users_add_email_for_authenticated_user_users_add_email_for_authenticated_user_request_users_add_email_for_authenticated_user_request_emails
  }
  property {
    name  = "usersCheckFollowingForUser_target_user"
    type  = "string"
    value = var.connector-oai-github_property_users_check_following_for_user_target_user
  }
  property {
    name  = "usersCheckFollowingForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_users_check_following_for_user_username
  }
  property {
    name  = "usersCheckPersonIsFollowedByAuthenticated_username"
    type  = "string"
    value = var.connector-oai-github_property_users_check_person_is_followed_by_authenticated_username
  }
  property {
    name  = "usersCreateGpgKeyForAuthenticatedUser_usersCreateGpgKeyForAuthenticatedUserRequest_UsersCreateGpgKeyForAuthenticatedUserRequest_armored_public_key"
    type  = "string"
    value = var.connector-oai-github_property_users_create_gpg_key_for_authenticated_user_users_create_gpg_key_for_authenticated_user_request_users_create_gpg_key_for_authenticated_user_request_armored_public_key
  }
  property {
    name  = "usersCreatePublicSshKeyForAuthenticatedUser_usersCreatePublicSshKeyForAuthenticatedUserRequest_UsersCreatePublicSshKeyForAuthenticatedUserRequest_key"
    type  = "string"
    value = var.connector-oai-github_property_users_create_public_ssh_key_for_authenticated_user_users_create_public_ssh_key_for_authenticated_user_request_users_create_public_ssh_key_for_authenticated_user_request_key
  }
  property {
    name  = "usersCreatePublicSshKeyForAuthenticatedUser_usersCreatePublicSshKeyForAuthenticatedUserRequest_UsersCreatePublicSshKeyForAuthenticatedUserRequest_title"
    type  = "string"
    value = var.connector-oai-github_property_users_create_public_ssh_key_for_authenticated_user_users_create_public_ssh_key_for_authenticated_user_request_users_create_public_ssh_key_for_authenticated_user_request_title
  }
  property {
    name  = "usersDeleteEmailForAuthenticatedUser_usersDeleteEmailForAuthenticatedUserRequest_UsersDeleteEmailForAuthenticatedUserRequest_emails"
    type  = "string"
    value = var.connector-oai-github_property_users_delete_email_for_authenticated_user_users_delete_email_for_authenticated_user_request_users_delete_email_for_authenticated_user_request_emails
  }
  property {
    name  = "usersDeleteGpgKeyForAuthenticatedUser_gpg_key_id"
    type  = "string"
    value = var.connector-oai-github_property_users_delete_gpg_key_for_authenticated_user_gpg_key_id
  }
  property {
    name  = "usersDeletePublicSshKeyForAuthenticatedUser_key_id"
    type  = "string"
    value = var.connector-oai-github_property_users_delete_public_ssh_key_for_authenticated_user_key_id
  }
  property {
    name  = "usersFollow_username"
    type  = "string"
    value = var.connector-oai-github_property_users_follow_username
  }
  property {
    name  = "usersGetByUsername_username"
    type  = "string"
    value = var.connector-oai-github_property_users_get_by_username_username
  }
  property {
    name  = "usersGetContextForUser_subject_id"
    type  = "string"
    value = var.connector-oai-github_property_users_get_context_for_user_subject_id
  }
  property {
    name  = "usersGetContextForUser_subject_type"
    type  = "string"
    value = var.connector-oai-github_property_users_get_context_for_user_subject_type
  }
  property {
    name  = "usersGetContextForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_users_get_context_for_user_username
  }
  property {
    name  = "usersGetGpgKeyForAuthenticatedUser_gpg_key_id"
    type  = "string"
    value = var.connector-oai-github_property_users_get_gpg_key_for_authenticated_user_gpg_key_id
  }
  property {
    name  = "usersGetPublicSshKeyForAuthenticatedUser_key_id"
    type  = "string"
    value = var.connector-oai-github_property_users_get_public_ssh_key_for_authenticated_user_key_id
  }
  property {
    name  = "usersListEmailsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_emails_for_authenticated_user_page
  }
  property {
    name  = "usersListEmailsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_emails_for_authenticated_user_per_page
  }
  property {
    name  = "usersListFollowedByAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_followed_by_authenticated_user_page
  }
  property {
    name  = "usersListFollowedByAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_followed_by_authenticated_user_per_page
  }
  property {
    name  = "usersListFollowersForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_followers_for_authenticated_user_page
  }
  property {
    name  = "usersListFollowersForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_followers_for_authenticated_user_per_page
  }
  property {
    name  = "usersListFollowersForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_followers_for_user_page
  }
  property {
    name  = "usersListFollowersForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_followers_for_user_per_page
  }
  property {
    name  = "usersListFollowersForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_users_list_followers_for_user_username
  }
  property {
    name  = "usersListFollowingForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_following_for_user_page
  }
  property {
    name  = "usersListFollowingForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_following_for_user_per_page
  }
  property {
    name  = "usersListFollowingForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_users_list_following_for_user_username
  }
  property {
    name  = "usersListGpgKeysForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_gpg_keys_for_authenticated_user_page
  }
  property {
    name  = "usersListGpgKeysForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_gpg_keys_for_authenticated_user_per_page
  }
  property {
    name  = "usersListGpgKeysForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_gpg_keys_for_user_page
  }
  property {
    name  = "usersListGpgKeysForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_gpg_keys_for_user_per_page
  }
  property {
    name  = "usersListGpgKeysForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_users_list_gpg_keys_for_user_username
  }
  property {
    name  = "usersListPublicEmailsForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_public_emails_for_authenticated_user_page
  }
  property {
    name  = "usersListPublicEmailsForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_public_emails_for_authenticated_user_per_page
  }
  property {
    name  = "usersListPublicKeysForUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_public_keys_for_user_page
  }
  property {
    name  = "usersListPublicKeysForUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_public_keys_for_user_per_page
  }
  property {
    name  = "usersListPublicKeysForUser_username"
    type  = "string"
    value = var.connector-oai-github_property_users_list_public_keys_for_user_username
  }
  property {
    name  = "usersListPublicSshKeysForAuthenticatedUser_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_public_ssh_keys_for_authenticated_user_page
  }
  property {
    name  = "usersListPublicSshKeysForAuthenticatedUser_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_public_ssh_keys_for_authenticated_user_per_page
  }
  property {
    name  = "usersList_per_page"
    type  = "string"
    value = var.connector-oai-github_property_users_list_per_page
  }
  property {
    name  = "usersList_since"
    type  = "string"
    value = var.connector-oai-github_property_users_list_since
  }
  property {
    name  = "usersUnfollow_username"
    type  = "string"
    value = var.connector-oai-github_property_users_unfollow_username
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_bio"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_bio
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_blog"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_blog
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_company"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_company
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_email"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_email
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_hireable"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_hireable
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_location"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_location
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_name"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_name
  }
  property {
    name  = "usersUpdateAuthenticated_usersUpdateAuthenticatedRequest_UsersUpdateAuthenticatedRequest_twitter_username"
    type  = "string"
    value = var.connector-oai-github_property_users_update_authenticated_users_update_authenticated_request_users_update_authenticated_request_twitter_username
  }
}
