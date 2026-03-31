resource "pingone_davinci_connector_instance" "connector-oai-launchdarklyrestapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-launchdarklyrestapi"
  }
  name = "My awesome connector-oai-launchdarklyrestapi"
  property {
    name  = "associateRepositoriesAndProjects_associateRepositoriesAndProjectsRequest_AssociateRepositoriesAndProjectsRequest_mappings"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_associate_repositories_and_projects_associate_repositories_and_projects_request_associate_repositories_and_projects_request_mappings
  }
  property {
    name  = "authApiKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_auth_api_key
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_base_path
  }
  property {
    name  = "copyFeatureFlag_copyFeatureFlagRequest_CopyFeatureFlagRequestSource_currentVersion"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_copy_feature_flag_request_copy_feature_flag_request_source_current_version
  }
  property {
    name  = "copyFeatureFlag_copyFeatureFlagRequest_CopyFeatureFlagRequestSource_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_copy_feature_flag_request_copy_feature_flag_request_source_key
  }
  property {
    name  = "copyFeatureFlag_copyFeatureFlagRequest_CopyFeatureFlagRequestTarget_currentVersion"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_copy_feature_flag_request_copy_feature_flag_request_target_current_version
  }
  property {
    name  = "copyFeatureFlag_copyFeatureFlagRequest_CopyFeatureFlagRequestTarget_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_copy_feature_flag_request_copy_feature_flag_request_target_key
  }
  property {
    name  = "copyFeatureFlag_copyFeatureFlagRequest_CopyFeatureFlagRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_copy_feature_flag_request_copy_feature_flag_request_comment
  }
  property {
    name  = "copyFeatureFlag_copyFeatureFlagRequest_CopyFeatureFlagRequest_excludedActions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_copy_feature_flag_request_copy_feature_flag_request_excluded_actions
  }
  property {
    name  = "copyFeatureFlag_copyFeatureFlagRequest_CopyFeatureFlagRequest_includedActions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_copy_feature_flag_request_copy_feature_flag_request_included_actions
  }
  property {
    name  = "copyFeatureFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_feature_flag_key
  }
  property {
    name  = "copyFeatureFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_copy_feature_flag_project_key
  }
  property {
    name  = "createBigSegmentExport_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_export_environment_key
  }
  property {
    name  = "createBigSegmentExport_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_export_project_key
  }
  property {
    name  = "createBigSegmentExport_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_export_segment_key
  }
  property {
    name  = "createBigSegmentImport_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_import_environment_key
  }
  property {
    name  = "createBigSegmentImport_file"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_import_file
  }
  property {
    name  = "createBigSegmentImport_mode"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_import_mode
  }
  property {
    name  = "createBigSegmentImport_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_import_project_key
  }
  property {
    name  = "createBigSegmentImport_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_import_segment_key
  }
  property {
    name  = "createBigSegmentStoreIntegration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_config"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_store_integration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_config
  }
  property {
    name  = "createBigSegmentStoreIntegration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_store_integration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_name
  }
  property {
    name  = "createBigSegmentStoreIntegration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_on"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_store_integration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_on
  }
  property {
    name  = "createBigSegmentStoreIntegration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_store_integration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_tags
  }
  property {
    name  = "createBigSegmentStoreIntegration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_store_integration_environment_key
  }
  property {
    name  = "createBigSegmentStoreIntegration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_store_integration_integration_key
  }
  property {
    name  = "createBigSegmentStoreIntegration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_big_segment_store_integration_project_key
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_applicationKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_application_key
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_applicationKind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_application_kind
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_applicationName"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_application_name
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_deploymentMetadata"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_deployment_metadata
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_environmentKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_environment_key
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_eventMetadata"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_event_metadata
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_eventTime"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_event_time
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_eventType"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_event_type
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_projectKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_project_key
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_version"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_version
  }
  property {
    name  = "createDeploymentEvent_createDeploymentEventRequest_CreateDeploymentEventRequest_versionName"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_deployment_event_create_deployment_event_request_create_deployment_event_request_version_name
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_canReshuffleTraffic"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_can_reshuffle_traffic
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_flags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_flags
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_hypothesis"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_hypothesis
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_metrics"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_metrics
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_primaryFunnelKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_primary_funnel_key
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_primarySingleMetricKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_primary_single_metric_key
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_randomizationUnit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_randomization_unit
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequestIteration_treatments"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_iteration_treatments
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_description
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_key
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequest_maintainerId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_maintainer_id
  }
  property {
    name  = "createExperiment_createExperimentRequest_CreateExperimentRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_create_experiment_request_create_experiment_request_name
  }
  property {
    name  = "createExperiment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_environment_key
  }
  property {
    name  = "createExperiment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_experiment_project_key
  }
  property {
    name  = "createFlagLink_createFlagLinkRequest_CreateFlagLinkRequest_deepLink"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_create_flag_link_request_create_flag_link_request_deep_link
  }
  property {
    name  = "createFlagLink_createFlagLinkRequest_CreateFlagLinkRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_create_flag_link_request_create_flag_link_request_description
  }
  property {
    name  = "createFlagLink_createFlagLinkRequest_CreateFlagLinkRequest_integrationKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_create_flag_link_request_create_flag_link_request_integration_key
  }
  property {
    name  = "createFlagLink_createFlagLinkRequest_CreateFlagLinkRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_create_flag_link_request_create_flag_link_request_key
  }
  property {
    name  = "createFlagLink_createFlagLinkRequest_CreateFlagLinkRequest_metadata"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_create_flag_link_request_create_flag_link_request_metadata
  }
  property {
    name  = "createFlagLink_createFlagLinkRequest_CreateFlagLinkRequest_timestamp"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_create_flag_link_request_create_flag_link_request_timestamp
  }
  property {
    name  = "createFlagLink_createFlagLinkRequest_CreateFlagLinkRequest_title"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_create_flag_link_request_create_flag_link_request_title
  }
  property {
    name  = "createFlagLink_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_feature_flag_key
  }
  property {
    name  = "createFlagLink_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_flag_link_project_key
  }
  property {
    name  = "createInsightGroup_createInsightGroupRequest_CreateInsightGroupRequest_applicationKeys"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_insight_group_create_insight_group_request_create_insight_group_request_application_keys
  }
  property {
    name  = "createInsightGroup_createInsightGroupRequest_CreateInsightGroupRequest_environmentKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_insight_group_create_insight_group_request_create_insight_group_request_environment_key
  }
  property {
    name  = "createInsightGroup_createInsightGroupRequest_CreateInsightGroupRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_insight_group_create_insight_group_request_create_insight_group_request_key
  }
  property {
    name  = "createInsightGroup_createInsightGroupRequest_CreateInsightGroupRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_insight_group_create_insight_group_request_create_insight_group_request_name
  }
  property {
    name  = "createInsightGroup_createInsightGroupRequest_CreateInsightGroupRequest_projectKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_insight_group_create_insight_group_request_create_insight_group_request_project_key
  }
  property {
    name  = "createIntegrationDeliveryConfiguration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_config"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_integration_delivery_configuration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_config
  }
  property {
    name  = "createIntegrationDeliveryConfiguration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_integration_delivery_configuration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_name
  }
  property {
    name  = "createIntegrationDeliveryConfiguration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_on"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_integration_delivery_configuration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_on
  }
  property {
    name  = "createIntegrationDeliveryConfiguration_createBigSegmentStoreIntegrationRequest_CreateBigSegmentStoreIntegrationRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_integration_delivery_configuration_create_big_segment_store_integration_request_create_big_segment_store_integration_request_tags
  }
  property {
    name  = "createIntegrationDeliveryConfiguration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_integration_delivery_configuration_environment_key
  }
  property {
    name  = "createIntegrationDeliveryConfiguration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_integration_delivery_configuration_integration_key
  }
  property {
    name  = "createIntegrationDeliveryConfiguration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_integration_delivery_configuration_project_key
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_canReshuffleTraffic"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_can_reshuffle_traffic
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_flags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_flags
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_hypothesis"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_hypothesis
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_metrics"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_metrics
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_primaryFunnelKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_primary_funnel_key
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_primarySingleMetricKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_primary_single_metric_key
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_randomizationUnit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_randomization_unit
  }
  property {
    name  = "createIteration_createIterationRequest_CreateIterationRequest_treatments"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_create_iteration_request_create_iteration_request_treatments
  }
  property {
    name  = "createIteration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_environment_key
  }
  property {
    name  = "createIteration_experiment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_experiment_key
  }
  property {
    name  = "createIteration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_iteration_project_key
  }
  property {
    name  = "createMetricGroup_createMetricGroupRequest_CreateMetricGroupRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_create_metric_group_request_create_metric_group_request_description
  }
  property {
    name  = "createMetricGroup_createMetricGroupRequest_CreateMetricGroupRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_create_metric_group_request_create_metric_group_request_key
  }
  property {
    name  = "createMetricGroup_createMetricGroupRequest_CreateMetricGroupRequest_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_create_metric_group_request_create_metric_group_request_kind
  }
  property {
    name  = "createMetricGroup_createMetricGroupRequest_CreateMetricGroupRequest_maintainerId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_create_metric_group_request_create_metric_group_request_maintainer_id
  }
  property {
    name  = "createMetricGroup_createMetricGroupRequest_CreateMetricGroupRequest_metrics"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_create_metric_group_request_create_metric_group_request_metrics
  }
  property {
    name  = "createMetricGroup_createMetricGroupRequest_CreateMetricGroupRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_create_metric_group_request_create_metric_group_request_name
  }
  property {
    name  = "createMetricGroup_createMetricGroupRequest_CreateMetricGroupRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_create_metric_group_request_create_metric_group_request_tags
  }
  property {
    name  = "createMetricGroup_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_metric_group_project_key
  }
  property {
    name  = "createOAuth2Client_createOAuth2ClientRequest_CreateOAuth2ClientRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_oauth2_client_create_oauth2_client_request_create_oauth2_client_request_description
  }
  property {
    name  = "createOAuth2Client_createOAuth2ClientRequest_CreateOAuth2ClientRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_oauth2_client_create_oauth2_client_request_create_oauth2_client_request_name
  }
  property {
    name  = "createOAuth2Client_createOAuth2ClientRequest_CreateOAuth2ClientRequest_redirectUri"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_oauth2_client_create_oauth2_client_request_create_oauth2_client_request_redirect_uri
  }
  property {
    name  = "createSubscription_createSubscriptionRequest_CreateSubscriptionRequest_apiKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_create_subscription_request_create_subscription_request_api_key
  }
  property {
    name  = "createSubscription_createSubscriptionRequest_CreateSubscriptionRequest_config"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_create_subscription_request_create_subscription_request_config
  }
  property {
    name  = "createSubscription_createSubscriptionRequest_CreateSubscriptionRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_create_subscription_request_create_subscription_request_name
  }
  property {
    name  = "createSubscription_createSubscriptionRequest_CreateSubscriptionRequest_on"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_create_subscription_request_create_subscription_request_on
  }
  property {
    name  = "createSubscription_createSubscriptionRequest_CreateSubscriptionRequest_statements"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_create_subscription_request_create_subscription_request_statements
  }
  property {
    name  = "createSubscription_createSubscriptionRequest_CreateSubscriptionRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_create_subscription_request_create_subscription_request_tags
  }
  property {
    name  = "createSubscription_createSubscriptionRequest_CreateSubscriptionRequest_url"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_create_subscription_request_create_subscription_request_url
  }
  property {
    name  = "createSubscription_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_subscription_integration_key
  }
  property {
    name  = "createTriggerWorkflow_createTriggerWorkflowRequest_CreateTriggerWorkflowRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_trigger_workflow_create_trigger_workflow_request_create_trigger_workflow_request_comment
  }
  property {
    name  = "createTriggerWorkflow_createTriggerWorkflowRequest_CreateTriggerWorkflowRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_trigger_workflow_create_trigger_workflow_request_create_trigger_workflow_request_instructions
  }
  property {
    name  = "createTriggerWorkflow_createTriggerWorkflowRequest_CreateTriggerWorkflowRequest_integrationKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_trigger_workflow_create_trigger_workflow_request_create_trigger_workflow_request_integration_key
  }
  property {
    name  = "createTriggerWorkflow_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_trigger_workflow_environment_key
  }
  property {
    name  = "createTriggerWorkflow_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_trigger_workflow_feature_flag_key
  }
  property {
    name  = "createTriggerWorkflow_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_trigger_workflow_project_key
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_description
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_environmentKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_environment_key
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_flagKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_flag_key
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_key
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_name
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_projectKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_project_key
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_stages"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_stages
  }
  property {
    name  = "createWorkflowTemplate_createWorkflowTemplateRequest_CreateWorkflowTemplateRequest_workflowId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_create_workflow_template_create_workflow_template_request_create_workflow_template_request_workflow_id
  }
  property {
    name  = "deleteApplicationVersion_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_application_version_application_key
  }
  property {
    name  = "deleteApplicationVersion_version_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_application_version_version_key
  }
  property {
    name  = "deleteApplication_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_application_application_key
  }
  property {
    name  = "deleteApprovalRequestForFlag_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_approval_request_for_flag_environment_key
  }
  property {
    name  = "deleteApprovalRequestForFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_approval_request_for_flag_feature_flag_key
  }
  property {
    name  = "deleteApprovalRequestForFlag_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_approval_request_for_flag_id
  }
  property {
    name  = "deleteApprovalRequestForFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_approval_request_for_flag_project_key
  }
  property {
    name  = "deleteApprovalRequest_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_approval_request_id
  }
  property {
    name  = "deleteBigSegmentStoreIntegration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_big_segment_store_integration_environment_key
  }
  property {
    name  = "deleteBigSegmentStoreIntegration_integration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_big_segment_store_integration_integration_id
  }
  property {
    name  = "deleteBigSegmentStoreIntegration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_big_segment_store_integration_integration_key
  }
  property {
    name  = "deleteBigSegmentStoreIntegration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_big_segment_store_integration_project_key
  }
  property {
    name  = "deleteBranches_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_branches_repo
  }
  property {
    name  = "deleteBranches_request_body"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_branches_request_body
  }
  property {
    name  = "deleteContextInstances_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_context_instances_environment_key
  }
  property {
    name  = "deleteContextInstances_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_context_instances_id
  }
  property {
    name  = "deleteContextInstances_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_context_instances_project_key
  }
  property {
    name  = "deleteCustomRole_custom_role_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_custom_role_custom_role_key
  }
  property {
    name  = "deleteDestination_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_destination_environment_key
  }
  property {
    name  = "deleteDestination_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_destination_id
  }
  property {
    name  = "deleteDestination_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_destination_project_key
  }
  property {
    name  = "deleteEnvironment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_environment_environment_key
  }
  property {
    name  = "deleteEnvironment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_environment_project_key
  }
  property {
    name  = "deleteFeatureFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_feature_flag_feature_flag_key
  }
  property {
    name  = "deleteFeatureFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_feature_flag_project_key
  }
  property {
    name  = "deleteFlagConfigScheduledChanges_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_config_scheduled_changes_environment_key
  }
  property {
    name  = "deleteFlagConfigScheduledChanges_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_config_scheduled_changes_feature_flag_key
  }
  property {
    name  = "deleteFlagConfigScheduledChanges_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_config_scheduled_changes_id
  }
  property {
    name  = "deleteFlagConfigScheduledChanges_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_config_scheduled_changes_project_key
  }
  property {
    name  = "deleteFlagFollowers_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_followers_environment_key
  }
  property {
    name  = "deleteFlagFollowers_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_followers_feature_flag_key
  }
  property {
    name  = "deleteFlagFollowers_member_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_followers_member_id
  }
  property {
    name  = "deleteFlagFollowers_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_followers_project_key
  }
  property {
    name  = "deleteFlagLink_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_link_feature_flag_key
  }
  property {
    name  = "deleteFlagLink_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_link_id
  }
  property {
    name  = "deleteFlagLink_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_flag_link_project_key
  }
  property {
    name  = "deleteInsightGroup_insight_group_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_insight_group_insight_group_key
  }
  property {
    name  = "deleteIntegrationDeliveryConfiguration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_integration_delivery_configuration_environment_key
  }
  property {
    name  = "deleteIntegrationDeliveryConfiguration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_integration_delivery_configuration_id
  }
  property {
    name  = "deleteIntegrationDeliveryConfiguration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_integration_delivery_configuration_integration_key
  }
  property {
    name  = "deleteIntegrationDeliveryConfiguration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_integration_delivery_configuration_project_key
  }
  property {
    name  = "deleteMember_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_member_id
  }
  property {
    name  = "deleteMetricGroup_metric_group_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_metric_group_metric_group_key
  }
  property {
    name  = "deleteMetricGroup_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_metric_group_project_key
  }
  property {
    name  = "deleteMetric_metric_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_metric_metric_key
  }
  property {
    name  = "deleteMetric_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_metric_project_key
  }
  property {
    name  = "deleteOAuthClient_client_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_oauth_client_client_id
  }
  property {
    name  = "deleteProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_project_project_key
  }
  property {
    name  = "deleteRelayAutoConfig_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_relay_auto_config_id
  }
  property {
    name  = "deleteReleasePipeline_pipeline_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_release_pipeline_pipeline_key
  }
  property {
    name  = "deleteReleasePipeline_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_release_pipeline_project_key
  }
  property {
    name  = "deleteRepositoryProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_repository_project_project_key
  }
  property {
    name  = "deleteRepositoryProject_repository_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_repository_project_repository_key
  }
  property {
    name  = "deleteRepository_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_repository_repo
  }
  property {
    name  = "deleteSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_segment_environment_key
  }
  property {
    name  = "deleteSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_segment_project_key
  }
  property {
    name  = "deleteSegment_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_segment_segment_key
  }
  property {
    name  = "deleteSubscription_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_subscription_id
  }
  property {
    name  = "deleteSubscription_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_subscription_integration_key
  }
  property {
    name  = "deleteTeam_team_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_team_team_key
  }
  property {
    name  = "deleteToken_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_token_id
  }
  property {
    name  = "deleteTriggerWorkflow_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_trigger_workflow_environment_key
  }
  property {
    name  = "deleteTriggerWorkflow_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_trigger_workflow_feature_flag_key
  }
  property {
    name  = "deleteTriggerWorkflow_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_trigger_workflow_id
  }
  property {
    name  = "deleteTriggerWorkflow_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_trigger_workflow_project_key
  }
  property {
    name  = "deleteUser_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_user_environment_key
  }
  property {
    name  = "deleteUser_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_user_project_key
  }
  property {
    name  = "deleteUser_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_user_user_key
  }
  property {
    name  = "deleteWebhook_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_webhook_id
  }
  property {
    name  = "deleteWorkflowTemplate_template_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_workflow_template_template_key
  }
  property {
    name  = "deleteWorkflow_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_workflow_environment_key
  }
  property {
    name  = "deleteWorkflow_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_workflow_feature_flag_key
  }
  property {
    name  = "deleteWorkflow_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_workflow_project_key
  }
  property {
    name  = "deleteWorkflow_workflow_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_delete_workflow_workflow_id
  }
  property {
    name  = "evaluateContextInstance_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_evaluate_context_instance_environment_key
  }
  property {
    name  = "evaluateContextInstance_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_evaluate_context_instance_filter
  }
  property {
    name  = "evaluateContextInstance_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_evaluate_context_instance_limit
  }
  property {
    name  = "evaluateContextInstance_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_evaluate_context_instance_offset
  }
  property {
    name  = "evaluateContextInstance_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_evaluate_context_instance_project_key
  }
  property {
    name  = "evaluateContextInstance_request_body"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_evaluate_context_instance_request_body
  }
  property {
    name  = "evaluateContextInstance_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_evaluate_context_instance_sort
  }
  property {
    name  = "getAllReleasePipelines_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_all_release_pipelines_filter
  }
  property {
    name  = "getAllReleasePipelines_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_all_release_pipelines_limit
  }
  property {
    name  = "getAllReleasePipelines_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_all_release_pipelines_offset
  }
  property {
    name  = "getAllReleasePipelines_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_all_release_pipelines_project_key
  }
  property {
    name  = "getApplicationVersions_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_application_versions_application_key
  }
  property {
    name  = "getApplicationVersions_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_application_versions_filter
  }
  property {
    name  = "getApplicationVersions_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_application_versions_limit
  }
  property {
    name  = "getApplicationVersions_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_application_versions_offset
  }
  property {
    name  = "getApplicationVersions_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_application_versions_sort
  }
  property {
    name  = "getApplication_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_application_application_key
  }
  property {
    name  = "getApplication_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_application_expand
  }
  property {
    name  = "getApplications_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_applications_expand
  }
  property {
    name  = "getApplications_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_applications_filter
  }
  property {
    name  = "getApplications_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_applications_limit
  }
  property {
    name  = "getApplications_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_applications_offset
  }
  property {
    name  = "getApplications_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_applications_sort
  }
  property {
    name  = "getApprovalForFlag_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_for_flag_environment_key
  }
  property {
    name  = "getApprovalForFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_for_flag_feature_flag_key
  }
  property {
    name  = "getApprovalForFlag_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_for_flag_id
  }
  property {
    name  = "getApprovalForFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_for_flag_project_key
  }
  property {
    name  = "getApprovalRequest_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_request_expand
  }
  property {
    name  = "getApprovalRequest_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_request_id
  }
  property {
    name  = "getApprovalRequests_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_requests_expand
  }
  property {
    name  = "getApprovalRequests_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_requests_filter
  }
  property {
    name  = "getApprovalRequests_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_requests_limit
  }
  property {
    name  = "getApprovalRequests_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approval_requests_offset
  }
  property {
    name  = "getApprovalsForFlag_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approvals_for_flag_environment_key
  }
  property {
    name  = "getApprovalsForFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approvals_for_flag_feature_flag_key
  }
  property {
    name  = "getApprovalsForFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_approvals_for_flag_project_key
  }
  property {
    name  = "getAuditLogEntries_after"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_audit_log_entries_after
  }
  property {
    name  = "getAuditLogEntries_before"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_audit_log_entries_before
  }
  property {
    name  = "getAuditLogEntries_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_audit_log_entries_limit
  }
  property {
    name  = "getAuditLogEntries_q"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_audit_log_entries_q
  }
  property {
    name  = "getAuditLogEntries_spec"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_audit_log_entries_spec
  }
  property {
    name  = "getAuditLogEntry_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_audit_log_entry_id
  }
  property {
    name  = "getBigSegmentExport_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_export_environment_key
  }
  property {
    name  = "getBigSegmentExport_export_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_export_export_id
  }
  property {
    name  = "getBigSegmentExport_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_export_project_key
  }
  property {
    name  = "getBigSegmentExport_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_export_segment_key
  }
  property {
    name  = "getBigSegmentImport_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_import_environment_key
  }
  property {
    name  = "getBigSegmentImport_import_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_import_import_id
  }
  property {
    name  = "getBigSegmentImport_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_import_project_key
  }
  property {
    name  = "getBigSegmentImport_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_import_segment_key
  }
  property {
    name  = "getBigSegmentStoreIntegration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_store_integration_environment_key
  }
  property {
    name  = "getBigSegmentStoreIntegration_integration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_store_integration_integration_id
  }
  property {
    name  = "getBigSegmentStoreIntegration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_store_integration_integration_key
  }
  property {
    name  = "getBigSegmentStoreIntegration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_big_segment_store_integration_project_key
  }
  property {
    name  = "getBranch_branch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_branch_branch
  }
  property {
    name  = "getBranch_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_branch_flag_key
  }
  property {
    name  = "getBranch_proj_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_branch_proj_key
  }
  property {
    name  = "getBranch_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_branch_repo
  }
  property {
    name  = "getBranches_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_branches_repo
  }
  property {
    name  = "getContextAttributeNames_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_attribute_names_environment_key
  }
  property {
    name  = "getContextAttributeNames_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_attribute_names_filter
  }
  property {
    name  = "getContextAttributeNames_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_attribute_names_project_key
  }
  property {
    name  = "getContextAttributeValues_attribute_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_attribute_values_attribute_name
  }
  property {
    name  = "getContextAttributeValues_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_attribute_values_environment_key
  }
  property {
    name  = "getContextAttributeValues_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_attribute_values_filter
  }
  property {
    name  = "getContextAttributeValues_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_attribute_values_project_key
  }
  property {
    name  = "getContextInstanceSegmentsMembershipByEnv_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instance_segments_membership_by_env_environment_key
  }
  property {
    name  = "getContextInstanceSegmentsMembershipByEnv_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instance_segments_membership_by_env_project_key
  }
  property {
    name  = "getContextInstanceSegmentsMembershipByEnv_request_body"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instance_segments_membership_by_env_request_body
  }
  property {
    name  = "getContextInstances_continuation_token"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_continuation_token
  }
  property {
    name  = "getContextInstances_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_environment_key
  }
  property {
    name  = "getContextInstances_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_filter
  }
  property {
    name  = "getContextInstances_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_id
  }
  property {
    name  = "getContextInstances_include_total_count"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_include_total_count
  }
  property {
    name  = "getContextInstances_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_limit
  }
  property {
    name  = "getContextInstances_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_project_key
  }
  property {
    name  = "getContextInstances_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_instances_sort
  }
  property {
    name  = "getContextKindsByProjectKey_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_context_kinds_by_project_key_project_key
  }
  property {
    name  = "getContexts_continuation_token"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_continuation_token
  }
  property {
    name  = "getContexts_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_environment_key
  }
  property {
    name  = "getContexts_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_filter
  }
  property {
    name  = "getContexts_include_total_count"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_include_total_count
  }
  property {
    name  = "getContexts_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_key
  }
  property {
    name  = "getContexts_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_kind
  }
  property {
    name  = "getContexts_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_limit
  }
  property {
    name  = "getContexts_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_project_key
  }
  property {
    name  = "getContexts_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_contexts_sort
  }
  property {
    name  = "getCustomRole_custom_role_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_custom_role_custom_role_key
  }
  property {
    name  = "getCustomRoles_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_custom_roles_limit
  }
  property {
    name  = "getCustomRoles_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_custom_roles_offset
  }
  property {
    name  = "getCustomWorkflow_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_custom_workflow_environment_key
  }
  property {
    name  = "getCustomWorkflow_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_custom_workflow_feature_flag_key
  }
  property {
    name  = "getCustomWorkflow_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_custom_workflow_project_key
  }
  property {
    name  = "getCustomWorkflow_workflow_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_custom_workflow_workflow_id
  }
  property {
    name  = "getDataExportEventsUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_data_export_events_usage_from
  }
  property {
    name  = "getDataExportEventsUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_data_export_events_usage_to
  }
  property {
    name  = "getDependentFlagsByEnv_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_dependent_flags_by_env_environment_key
  }
  property {
    name  = "getDependentFlagsByEnv_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_dependent_flags_by_env_feature_flag_key
  }
  property {
    name  = "getDependentFlagsByEnv_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_dependent_flags_by_env_project_key
  }
  property {
    name  = "getDependentFlags_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_dependent_flags_feature_flag_key
  }
  property {
    name  = "getDependentFlags_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_dependent_flags_project_key
  }
  property {
    name  = "getDeploymentFrequencyChart_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_application_key
  }
  property {
    name  = "getDeploymentFrequencyChart_bucket_ms"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_bucket_ms
  }
  property {
    name  = "getDeploymentFrequencyChart_bucket_type"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_bucket_type
  }
  property {
    name  = "getDeploymentFrequencyChart_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_environment_key
  }
  property {
    name  = "getDeploymentFrequencyChart_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_expand
  }
  property {
    name  = "getDeploymentFrequencyChart_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_from
  }
  property {
    name  = "getDeploymentFrequencyChart_group_by"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_group_by
  }
  property {
    name  = "getDeploymentFrequencyChart_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_project_key
  }
  property {
    name  = "getDeploymentFrequencyChart_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_frequency_chart_to
  }
  property {
    name  = "getDeployment_deployment_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_deployment_id
  }
  property {
    name  = "getDeployment_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployment_expand
  }
  property {
    name  = "getDeployments_after"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_after
  }
  property {
    name  = "getDeployments_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_application_key
  }
  property {
    name  = "getDeployments_before"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_before
  }
  property {
    name  = "getDeployments_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_environment_key
  }
  property {
    name  = "getDeployments_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_expand
  }
  property {
    name  = "getDeployments_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_from
  }
  property {
    name  = "getDeployments_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_kind
  }
  property {
    name  = "getDeployments_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_limit
  }
  property {
    name  = "getDeployments_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_project_key
  }
  property {
    name  = "getDeployments_status"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_status
  }
  property {
    name  = "getDeployments_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_deployments_to
  }
  property {
    name  = "getDestination_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_destination_environment_key
  }
  property {
    name  = "getDestination_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_destination_id
  }
  property {
    name  = "getDestination_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_destination_project_key
  }
  property {
    name  = "getEnvironment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_environment_environment_key
  }
  property {
    name  = "getEnvironment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_environment_project_key
  }
  property {
    name  = "getEnvironmentsByProject_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_environments_by_project_filter
  }
  property {
    name  = "getEnvironmentsByProject_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_environments_by_project_limit
  }
  property {
    name  = "getEnvironmentsByProject_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_environments_by_project_offset
  }
  property {
    name  = "getEnvironmentsByProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_environments_by_project_project_key
  }
  property {
    name  = "getEnvironmentsByProject_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_environments_by_project_sort
  }
  property {
    name  = "getEvaluationsUsage_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_evaluations_usage_environment_key
  }
  property {
    name  = "getEvaluationsUsage_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_evaluations_usage_feature_flag_key
  }
  property {
    name  = "getEvaluationsUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_evaluations_usage_from
  }
  property {
    name  = "getEvaluationsUsage_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_evaluations_usage_project_key
  }
  property {
    name  = "getEvaluationsUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_evaluations_usage_to
  }
  property {
    name  = "getEvaluationsUsage_tz"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_evaluations_usage_tz
  }
  property {
    name  = "getEventsUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_events_usage_from
  }
  property {
    name  = "getEventsUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_events_usage_to
  }
  property {
    name  = "getEventsUsage_type"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_events_usage_type
  }
  property {
    name  = "getExperimentResultsForMetricGroup_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_for_metric_group_environment_key
  }
  property {
    name  = "getExperimentResultsForMetricGroup_experiment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_for_metric_group_experiment_key
  }
  property {
    name  = "getExperimentResultsForMetricGroup_iteration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_for_metric_group_iteration_id
  }
  property {
    name  = "getExperimentResultsForMetricGroup_metric_group_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_for_metric_group_metric_group_key
  }
  property {
    name  = "getExperimentResultsForMetricGroup_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_for_metric_group_project_key
  }
  property {
    name  = "getExperimentResults_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_environment_key
  }
  property {
    name  = "getExperimentResults_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_expand
  }
  property {
    name  = "getExperimentResults_experiment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_experiment_key
  }
  property {
    name  = "getExperimentResults_iteration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_iteration_id
  }
  property {
    name  = "getExperimentResults_metric_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_metric_key
  }
  property {
    name  = "getExperimentResults_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_results_project_key
  }
  property {
    name  = "getExperiment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_environment_key
  }
  property {
    name  = "getExperiment_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_expand
  }
  property {
    name  = "getExperiment_experiment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_experiment_key
  }
  property {
    name  = "getExperiment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiment_project_key
  }
  property {
    name  = "getExperimentationKeysUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experimentation_keys_usage_from
  }
  property {
    name  = "getExperimentationKeysUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experimentation_keys_usage_to
  }
  property {
    name  = "getExperimentationSettings_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experimentation_settings_project_key
  }
  property {
    name  = "getExperimentationUnitsUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experimentation_units_usage_from
  }
  property {
    name  = "getExperimentationUnitsUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experimentation_units_usage_to
  }
  property {
    name  = "getExperiments_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiments_environment_key
  }
  property {
    name  = "getExperiments_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiments_expand
  }
  property {
    name  = "getExperiments_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiments_filter
  }
  property {
    name  = "getExperiments_lifecycle_state"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiments_lifecycle_state
  }
  property {
    name  = "getExperiments_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiments_limit
  }
  property {
    name  = "getExperiments_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiments_offset
  }
  property {
    name  = "getExperiments_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_experiments_project_key
  }
  property {
    name  = "getExpiringContextTargets_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_context_targets_environment_key
  }
  property {
    name  = "getExpiringContextTargets_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_context_targets_feature_flag_key
  }
  property {
    name  = "getExpiringContextTargets_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_context_targets_project_key
  }
  property {
    name  = "getExpiringFlagsForUser_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_flags_for_user_environment_key
  }
  property {
    name  = "getExpiringFlagsForUser_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_flags_for_user_project_key
  }
  property {
    name  = "getExpiringFlagsForUser_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_flags_for_user_user_key
  }
  property {
    name  = "getExpiringTargetsForSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_targets_for_segment_environment_key
  }
  property {
    name  = "getExpiringTargetsForSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_targets_for_segment_project_key
  }
  property {
    name  = "getExpiringTargetsForSegment_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_targets_for_segment_segment_key
  }
  property {
    name  = "getExpiringUserTargetsForSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_user_targets_for_segment_environment_key
  }
  property {
    name  = "getExpiringUserTargetsForSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_user_targets_for_segment_project_key
  }
  property {
    name  = "getExpiringUserTargetsForSegment_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_user_targets_for_segment_segment_key
  }
  property {
    name  = "getExpiringUserTargets_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_user_targets_environment_key
  }
  property {
    name  = "getExpiringUserTargets_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_user_targets_feature_flag_key
  }
  property {
    name  = "getExpiringUserTargets_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_expiring_user_targets_project_key
  }
  property {
    name  = "getExtinctions_branch_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_extinctions_branch_name
  }
  property {
    name  = "getExtinctions_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_extinctions_flag_key
  }
  property {
    name  = "getExtinctions_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_extinctions_from
  }
  property {
    name  = "getExtinctions_proj_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_extinctions_proj_key
  }
  property {
    name  = "getExtinctions_repo_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_extinctions_repo_name
  }
  property {
    name  = "getExtinctions_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_extinctions_to
  }
  property {
    name  = "getFeatureFlagScheduledChange_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_scheduled_change_environment_key
  }
  property {
    name  = "getFeatureFlagScheduledChange_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_scheduled_change_feature_flag_key
  }
  property {
    name  = "getFeatureFlagScheduledChange_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_scheduled_change_id
  }
  property {
    name  = "getFeatureFlagScheduledChange_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_scheduled_change_project_key
  }
  property {
    name  = "getFeatureFlagStatusAcrossEnvironments_env"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_status_across_environments_env
  }
  property {
    name  = "getFeatureFlagStatusAcrossEnvironments_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_status_across_environments_feature_flag_key
  }
  property {
    name  = "getFeatureFlagStatusAcrossEnvironments_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_status_across_environments_project_key
  }
  property {
    name  = "getFeatureFlagStatus_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_status_environment_key
  }
  property {
    name  = "getFeatureFlagStatus_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_status_feature_flag_key
  }
  property {
    name  = "getFeatureFlagStatus_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_status_project_key
  }
  property {
    name  = "getFeatureFlagStatuses_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_statuses_environment_key
  }
  property {
    name  = "getFeatureFlagStatuses_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_statuses_project_key
  }
  property {
    name  = "getFeatureFlag_env"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_env
  }
  property {
    name  = "getFeatureFlag_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_expand
  }
  property {
    name  = "getFeatureFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_feature_flag_key
  }
  property {
    name  = "getFeatureFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flag_project_key
  }
  property {
    name  = "getFeatureFlags_archived"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_archived
  }
  property {
    name  = "getFeatureFlags_compare"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_compare
  }
  property {
    name  = "getFeatureFlags_env"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_env
  }
  property {
    name  = "getFeatureFlags_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_expand
  }
  property {
    name  = "getFeatureFlags_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_filter
  }
  property {
    name  = "getFeatureFlags_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_limit
  }
  property {
    name  = "getFeatureFlags_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_offset
  }
  property {
    name  = "getFeatureFlags_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_project_key
  }
  property {
    name  = "getFeatureFlags_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_sort
  }
  property {
    name  = "getFeatureFlags_summary"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_summary
  }
  property {
    name  = "getFeatureFlags_tag"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_feature_flags_tag
  }
  property {
    name  = "getFlagConfigScheduledChanges_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_config_scheduled_changes_environment_key
  }
  property {
    name  = "getFlagConfigScheduledChanges_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_config_scheduled_changes_feature_flag_key
  }
  property {
    name  = "getFlagConfigScheduledChanges_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_config_scheduled_changes_project_key
  }
  property {
    name  = "getFlagDefaultsByProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_defaults_by_project_project_key
  }
  property {
    name  = "getFlagEvents_after"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_after
  }
  property {
    name  = "getFlagEvents_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_application_key
  }
  property {
    name  = "getFlagEvents_before"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_before
  }
  property {
    name  = "getFlagEvents_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_environment_key
  }
  property {
    name  = "getFlagEvents_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_expand
  }
  property {
    name  = "getFlagEvents_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_from
  }
  property {
    name  = "getFlagEvents_global"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_global
  }
  property {
    name  = "getFlagEvents_has_experiments"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_has_experiments
  }
  property {
    name  = "getFlagEvents_impact_size"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_impact_size
  }
  property {
    name  = "getFlagEvents_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_limit
  }
  property {
    name  = "getFlagEvents_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_project_key
  }
  property {
    name  = "getFlagEvents_query"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_query
  }
  property {
    name  = "getFlagEvents_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_events_to
  }
  property {
    name  = "getFlagFollowers_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_followers_environment_key
  }
  property {
    name  = "getFlagFollowers_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_followers_feature_flag_key
  }
  property {
    name  = "getFlagFollowers_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_followers_project_key
  }
  property {
    name  = "getFlagLinks_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_links_feature_flag_key
  }
  property {
    name  = "getFlagLinks_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_links_project_key
  }
  property {
    name  = "getFlagStatusChart_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_status_chart_application_key
  }
  property {
    name  = "getFlagStatusChart_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_status_chart_environment_key
  }
  property {
    name  = "getFlagStatusChart_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_flag_status_chart_project_key
  }
  property {
    name  = "getFollowersByProjEnv_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_followers_by_proj_env_environment_key
  }
  property {
    name  = "getFollowersByProjEnv_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_followers_by_proj_env_project_key
  }
  property {
    name  = "getInsightGroup_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insight_group_expand
  }
  property {
    name  = "getInsightGroup_insight_group_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insight_group_insight_group_key
  }
  property {
    name  = "getInsightGroups_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insight_groups_expand
  }
  property {
    name  = "getInsightGroups_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insight_groups_limit
  }
  property {
    name  = "getInsightGroups_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insight_groups_offset
  }
  property {
    name  = "getInsightGroups_query"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insight_groups_query
  }
  property {
    name  = "getInsightGroups_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insight_groups_sort
  }
  property {
    name  = "getInsightsRepositories_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insights_repositories_expand
  }
  property {
    name  = "getInsightsScores_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insights_scores_application_key
  }
  property {
    name  = "getInsightsScores_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insights_scores_environment_key
  }
  property {
    name  = "getInsightsScores_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_insights_scores_project_key
  }
  property {
    name  = "getIntegrationDeliveryConfigurationByEnvironment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_integration_delivery_configuration_by_environment_environment_key
  }
  property {
    name  = "getIntegrationDeliveryConfigurationByEnvironment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_integration_delivery_configuration_by_environment_project_key
  }
  property {
    name  = "getIntegrationDeliveryConfigurationById_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_integration_delivery_configuration_by_id_environment_key
  }
  property {
    name  = "getIntegrationDeliveryConfigurationById_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_integration_delivery_configuration_by_id_id
  }
  property {
    name  = "getIntegrationDeliveryConfigurationById_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_integration_delivery_configuration_by_id_integration_key
  }
  property {
    name  = "getIntegrationDeliveryConfigurationById_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_integration_delivery_configuration_by_id_project_key
  }
  property {
    name  = "getLeadTimeChart_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_application_key
  }
  property {
    name  = "getLeadTimeChart_bucket_ms"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_bucket_ms
  }
  property {
    name  = "getLeadTimeChart_bucket_type"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_bucket_type
  }
  property {
    name  = "getLeadTimeChart_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_environment_key
  }
  property {
    name  = "getLeadTimeChart_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_expand
  }
  property {
    name  = "getLeadTimeChart_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_from
  }
  property {
    name  = "getLeadTimeChart_group_by"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_group_by
  }
  property {
    name  = "getLeadTimeChart_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_project_key
  }
  property {
    name  = "getLeadTimeChart_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_lead_time_chart_to
  }
  property {
    name  = "getLegacyExperimentResults_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_legacy_experiment_results_environment_key
  }
  property {
    name  = "getLegacyExperimentResults_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_legacy_experiment_results_feature_flag_key
  }
  property {
    name  = "getLegacyExperimentResults_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_legacy_experiment_results_from
  }
  property {
    name  = "getLegacyExperimentResults_metric_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_legacy_experiment_results_metric_key
  }
  property {
    name  = "getLegacyExperimentResults_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_legacy_experiment_results_project_key
  }
  property {
    name  = "getLegacyExperimentResults_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_legacy_experiment_results_to
  }
  property {
    name  = "getMauSdksByType_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_sdks_by_type_from
  }
  property {
    name  = "getMauSdksByType_sdktype"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_sdks_by_type_sdktype
  }
  property {
    name  = "getMauSdksByType_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_sdks_by_type_to
  }
  property {
    name  = "getMauUsageByCategory_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_by_category_from
  }
  property {
    name  = "getMauUsageByCategory_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_by_category_to
  }
  property {
    name  = "getMauUsage_anonymous"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_anonymous
  }
  property {
    name  = "getMauUsage_environment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_environment
  }
  property {
    name  = "getMauUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_from
  }
  property {
    name  = "getMauUsage_groupby"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_groupby
  }
  property {
    name  = "getMauUsage_project"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_project
  }
  property {
    name  = "getMauUsage_sdk"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_sdk
  }
  property {
    name  = "getMauUsage_sdktype"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_sdktype
  }
  property {
    name  = "getMauUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_mau_usage_to
  }
  property {
    name  = "getMember_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_member_id
  }
  property {
    name  = "getMembers_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_members_filter
  }
  property {
    name  = "getMembers_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_members_limit
  }
  property {
    name  = "getMembers_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_members_offset
  }
  property {
    name  = "getMembers_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_members_sort
  }
  property {
    name  = "getMetricGroup_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_group_expand
  }
  property {
    name  = "getMetricGroup_metric_group_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_group_metric_group_key
  }
  property {
    name  = "getMetricGroup_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_group_project_key
  }
  property {
    name  = "getMetricGroups_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_groups_expand
  }
  property {
    name  = "getMetricGroups_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_groups_project_key
  }
  property {
    name  = "getMetric_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_expand
  }
  property {
    name  = "getMetric_metric_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_metric_key
  }
  property {
    name  = "getMetric_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_project_key
  }
  property {
    name  = "getMetric_version_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metric_version_id
  }
  property {
    name  = "getMetrics_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metrics_expand
  }
  property {
    name  = "getMetrics_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_metrics_project_key
  }
  property {
    name  = "getOAuthClientById_client_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_oauth_client_by_id_client_id
  }
  property {
    name  = "getProject_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_project_expand
  }
  property {
    name  = "getProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_project_project_key
  }
  property {
    name  = "getProjects_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_projects_expand
  }
  property {
    name  = "getProjects_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_projects_filter
  }
  property {
    name  = "getProjects_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_projects_limit
  }
  property {
    name  = "getProjects_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_projects_offset
  }
  property {
    name  = "getProjects_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_projects_sort
  }
  property {
    name  = "getPullRequests_after"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_after
  }
  property {
    name  = "getPullRequests_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_application_key
  }
  property {
    name  = "getPullRequests_before"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_before
  }
  property {
    name  = "getPullRequests_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_environment_key
  }
  property {
    name  = "getPullRequests_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_expand
  }
  property {
    name  = "getPullRequests_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_from
  }
  property {
    name  = "getPullRequests_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_limit
  }
  property {
    name  = "getPullRequests_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_project_key
  }
  property {
    name  = "getPullRequests_query"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_query
  }
  property {
    name  = "getPullRequests_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_sort
  }
  property {
    name  = "getPullRequests_status"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_status
  }
  property {
    name  = "getPullRequests_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_pull_requests_to
  }
  property {
    name  = "getRelayProxyConfig_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_relay_proxy_config_id
  }
  property {
    name  = "getReleaseByFlagKey_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_by_flag_key_flag_key
  }
  property {
    name  = "getReleaseByFlagKey_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_by_flag_key_project_key
  }
  property {
    name  = "getReleaseFrequencyChart_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_application_key
  }
  property {
    name  = "getReleaseFrequencyChart_bucket_ms"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_bucket_ms
  }
  property {
    name  = "getReleaseFrequencyChart_bucket_type"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_bucket_type
  }
  property {
    name  = "getReleaseFrequencyChart_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_environment_key
  }
  property {
    name  = "getReleaseFrequencyChart_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_expand
  }
  property {
    name  = "getReleaseFrequencyChart_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_from
  }
  property {
    name  = "getReleaseFrequencyChart_global"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_global
  }
  property {
    name  = "getReleaseFrequencyChart_group_by"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_group_by
  }
  property {
    name  = "getReleaseFrequencyChart_has_experiments"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_has_experiments
  }
  property {
    name  = "getReleaseFrequencyChart_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_project_key
  }
  property {
    name  = "getReleaseFrequencyChart_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_frequency_chart_to
  }
  property {
    name  = "getReleasePipelineByKey_pipeline_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_pipeline_by_key_pipeline_key
  }
  property {
    name  = "getReleasePipelineByKey_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_release_pipeline_by_key_project_key
  }
  property {
    name  = "getRepositories_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_repositories_flag_key
  }
  property {
    name  = "getRepositories_proj_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_repositories_proj_key
  }
  property {
    name  = "getRepositories_with_branches"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_repositories_with_branches
  }
  property {
    name  = "getRepositories_with_references_for_default_branch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_repositories_with_references_for_default_branch
  }
  property {
    name  = "getRepository_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_repository_repo
  }
  property {
    name  = "getSearchUsers_after"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_after
  }
  property {
    name  = "getSearchUsers_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_environment_key
  }
  property {
    name  = "getSearchUsers_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_filter
  }
  property {
    name  = "getSearchUsers_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_limit
  }
  property {
    name  = "getSearchUsers_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_offset
  }
  property {
    name  = "getSearchUsers_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_project_key
  }
  property {
    name  = "getSearchUsers_q"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_q
  }
  property {
    name  = "getSearchUsers_search_after"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_search_after
  }
  property {
    name  = "getSearchUsers_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_search_users_sort
  }
  property {
    name  = "getSegmentMembershipForContext_context_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_context_context_key
  }
  property {
    name  = "getSegmentMembershipForContext_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_context_environment_key
  }
  property {
    name  = "getSegmentMembershipForContext_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_context_project_key
  }
  property {
    name  = "getSegmentMembershipForContext_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_context_segment_key
  }
  property {
    name  = "getSegmentMembershipForUser_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_user_environment_key
  }
  property {
    name  = "getSegmentMembershipForUser_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_user_project_key
  }
  property {
    name  = "getSegmentMembershipForUser_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_user_segment_key
  }
  property {
    name  = "getSegmentMembershipForUser_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_membership_for_user_user_key
  }
  property {
    name  = "getSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_environment_key
  }
  property {
    name  = "getSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_project_key
  }
  property {
    name  = "getSegment_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segment_segment_key
  }
  property {
    name  = "getSegments_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segments_environment_key
  }
  property {
    name  = "getSegments_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segments_filter
  }
  property {
    name  = "getSegments_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segments_limit
  }
  property {
    name  = "getSegments_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segments_offset
  }
  property {
    name  = "getSegments_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segments_project_key
  }
  property {
    name  = "getSegments_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_segments_sort
  }
  property {
    name  = "getServiceConnectionUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_service_connection_usage_from
  }
  property {
    name  = "getServiceConnectionUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_service_connection_usage_to
  }
  property {
    name  = "getStaleFlagsChart_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stale_flags_chart_application_key
  }
  property {
    name  = "getStaleFlagsChart_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stale_flags_chart_environment_key
  }
  property {
    name  = "getStaleFlagsChart_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stale_flags_chart_expand
  }
  property {
    name  = "getStaleFlagsChart_group_by"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stale_flags_chart_group_by
  }
  property {
    name  = "getStaleFlagsChart_maintainer_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stale_flags_chart_maintainer_id
  }
  property {
    name  = "getStaleFlagsChart_maintainer_team_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stale_flags_chart_maintainer_team_key
  }
  property {
    name  = "getStaleFlagsChart_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stale_flags_chart_project_key
  }
  property {
    name  = "getStatistics_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_statistics_flag_key
  }
  property {
    name  = "getStatistics_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_statistics_project_key
  }
  property {
    name  = "getStreamUsageBySdkVersion_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_by_sdk_version_from
  }
  property {
    name  = "getStreamUsageBySdkVersion_sdk"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_by_sdk_version_sdk
  }
  property {
    name  = "getStreamUsageBySdkVersion_source"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_by_sdk_version_source
  }
  property {
    name  = "getStreamUsageBySdkVersion_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_by_sdk_version_to
  }
  property {
    name  = "getStreamUsageBySdkVersion_tz"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_by_sdk_version_tz
  }
  property {
    name  = "getStreamUsageBySdkVersion_version"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_by_sdk_version_version
  }
  property {
    name  = "getStreamUsageSdkversion_source"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_sdkversion_source
  }
  property {
    name  = "getStreamUsage_from"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_from
  }
  property {
    name  = "getStreamUsage_source"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_source
  }
  property {
    name  = "getStreamUsage_to"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_to
  }
  property {
    name  = "getStreamUsage_tz"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_stream_usage_tz
  }
  property {
    name  = "getSubscriptionByID_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_subscription_by_id_id
  }
  property {
    name  = "getSubscriptionByID_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_subscription_by_id_integration_key
  }
  property {
    name  = "getSubscriptions_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_subscriptions_integration_key
  }
  property {
    name  = "getTags_archived"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_tags_archived
  }
  property {
    name  = "getTags_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_tags_kind
  }
  property {
    name  = "getTags_pre"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_tags_pre
  }
  property {
    name  = "getTeamMaintainers_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_maintainers_limit
  }
  property {
    name  = "getTeamMaintainers_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_maintainers_offset
  }
  property {
    name  = "getTeamMaintainers_team_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_maintainers_team_key
  }
  property {
    name  = "getTeamRoles_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_roles_limit
  }
  property {
    name  = "getTeamRoles_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_roles_offset
  }
  property {
    name  = "getTeamRoles_team_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_roles_team_key
  }
  property {
    name  = "getTeam_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_expand
  }
  property {
    name  = "getTeam_team_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_team_team_key
  }
  property {
    name  = "getTeams_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_teams_expand
  }
  property {
    name  = "getTeams_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_teams_filter
  }
  property {
    name  = "getTeams_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_teams_limit
  }
  property {
    name  = "getTeams_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_teams_offset
  }
  property {
    name  = "getToken_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_token_id
  }
  property {
    name  = "getTokens_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_tokens_limit
  }
  property {
    name  = "getTokens_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_tokens_offset
  }
  property {
    name  = "getTokens_show_all"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_tokens_show_all
  }
  property {
    name  = "getTriggerWorkflowById_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_trigger_workflow_by_id_environment_key
  }
  property {
    name  = "getTriggerWorkflowById_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_trigger_workflow_by_id_feature_flag_key
  }
  property {
    name  = "getTriggerWorkflowById_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_trigger_workflow_by_id_id
  }
  property {
    name  = "getTriggerWorkflowById_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_trigger_workflow_by_id_project_key
  }
  property {
    name  = "getTriggerWorkflows_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_trigger_workflows_environment_key
  }
  property {
    name  = "getTriggerWorkflows_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_trigger_workflows_feature_flag_key
  }
  property {
    name  = "getTriggerWorkflows_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_trigger_workflows_project_key
  }
  property {
    name  = "getUserAttributeNames_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_attribute_names_environment_key
  }
  property {
    name  = "getUserAttributeNames_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_attribute_names_project_key
  }
  property {
    name  = "getUserFlagSetting_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_flag_setting_environment_key
  }
  property {
    name  = "getUserFlagSetting_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_flag_setting_feature_flag_key
  }
  property {
    name  = "getUserFlagSetting_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_flag_setting_project_key
  }
  property {
    name  = "getUserFlagSetting_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_flag_setting_user_key
  }
  property {
    name  = "getUserFlagSettings_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_flag_settings_environment_key
  }
  property {
    name  = "getUserFlagSettings_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_flag_settings_project_key
  }
  property {
    name  = "getUserFlagSettings_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_flag_settings_user_key
  }
  property {
    name  = "getUser_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_environment_key
  }
  property {
    name  = "getUser_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_project_key
  }
  property {
    name  = "getUser_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_user_user_key
  }
  property {
    name  = "getUsers_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_users_environment_key
  }
  property {
    name  = "getUsers_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_users_limit
  }
  property {
    name  = "getUsers_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_users_project_key
  }
  property {
    name  = "getUsers_search_after"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_users_search_after
  }
  property {
    name  = "getWebhook_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_webhook_id
  }
  property {
    name  = "getWorkflowTemplates_search"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflow_templates_search
  }
  property {
    name  = "getWorkflowTemplates_summary"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflow_templates_summary
  }
  property {
    name  = "getWorkflows_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflows_environment_key
  }
  property {
    name  = "getWorkflows_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflows_feature_flag_key
  }
  property {
    name  = "getWorkflows_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflows_limit
  }
  property {
    name  = "getWorkflows_offset"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflows_offset
  }
  property {
    name  = "getWorkflows_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflows_project_key
  }
  property {
    name  = "getWorkflows_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflows_sort
  }
  property {
    name  = "getWorkflows_status"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_get_workflows_status
  }
  property {
    name  = "patchApplicationVersion_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_application_version_application_key
  }
  property {
    name  = "patchApplicationVersion_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_application_version_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchApplicationVersion_version_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_application_version_version_key
  }
  property {
    name  = "patchApplication_application_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_application_application_key
  }
  property {
    name  = "patchApplication_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_application_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchBigSegmentStoreIntegration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_big_segment_store_integration_environment_key
  }
  property {
    name  = "patchBigSegmentStoreIntegration_integration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_big_segment_store_integration_integration_id
  }
  property {
    name  = "patchBigSegmentStoreIntegration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_big_segment_store_integration_integration_key
  }
  property {
    name  = "patchBigSegmentStoreIntegration_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_big_segment_store_integration_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchBigSegmentStoreIntegration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_big_segment_store_integration_project_key
  }
  property {
    name  = "patchCustomRole_custom_role_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_custom_role_custom_role_key
  }
  property {
    name  = "patchCustomRole_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_custom_role_patch_relay_auto_config_request_patch_relay_auto_config_request_comment
  }
  property {
    name  = "patchCustomRole_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_patch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_custom_role_patch_relay_auto_config_request_patch_relay_auto_config_request_patch
  }
  property {
    name  = "patchDestination_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_destination_environment_key
  }
  property {
    name  = "patchDestination_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_destination_id
  }
  property {
    name  = "patchDestination_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_destination_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchDestination_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_destination_project_key
  }
  property {
    name  = "patchEnvironment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_environment_environment_key
  }
  property {
    name  = "patchEnvironment_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_environment_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchEnvironment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_environment_project_key
  }
  property {
    name  = "patchExperiment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_experiment_environment_key
  }
  property {
    name  = "patchExperiment_experiment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_experiment_experiment_key
  }
  property {
    name  = "patchExperiment_patchExperimentRequest_PatchExperimentRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_experiment_patch_experiment_request_patch_experiment_request_comment
  }
  property {
    name  = "patchExperiment_patchExperimentRequest_PatchExperimentRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_experiment_patch_experiment_request_patch_experiment_request_instructions
  }
  property {
    name  = "patchExperiment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_experiment_project_key
  }
  property {
    name  = "patchExpiringFlagsForUser_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_flags_for_user_environment_key
  }
  property {
    name  = "patchExpiringFlagsForUser_patchExpiringFlagsForUserRequest_PatchExpiringFlagsForUserRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_flags_for_user_patch_expiring_flags_for_user_request_patch_expiring_flags_for_user_request_comment
  }
  property {
    name  = "patchExpiringFlagsForUser_patchExpiringFlagsForUserRequest_PatchExpiringFlagsForUserRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_flags_for_user_patch_expiring_flags_for_user_request_patch_expiring_flags_for_user_request_instructions
  }
  property {
    name  = "patchExpiringFlagsForUser_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_flags_for_user_project_key
  }
  property {
    name  = "patchExpiringFlagsForUser_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_flags_for_user_user_key
  }
  property {
    name  = "patchExpiringTargetsForSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_for_segment_environment_key
  }
  property {
    name  = "patchExpiringTargetsForSegment_patchExpiringTargetsForSegmentRequest_PatchExpiringTargetsForSegmentRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_for_segment_patch_expiring_targets_for_segment_request_patch_expiring_targets_for_segment_request_comment
  }
  property {
    name  = "patchExpiringTargetsForSegment_patchExpiringTargetsForSegmentRequest_PatchExpiringTargetsForSegmentRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_for_segment_patch_expiring_targets_for_segment_request_patch_expiring_targets_for_segment_request_instructions
  }
  property {
    name  = "patchExpiringTargetsForSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_for_segment_project_key
  }
  property {
    name  = "patchExpiringTargetsForSegment_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_for_segment_segment_key
  }
  property {
    name  = "patchExpiringTargets_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_environment_key
  }
  property {
    name  = "patchExpiringTargets_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_feature_flag_key
  }
  property {
    name  = "patchExpiringTargets_patchExpiringTargetsRequest_PatchExpiringTargetsRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_patch_expiring_targets_request_patch_expiring_targets_request_comment
  }
  property {
    name  = "patchExpiringTargets_patchExpiringTargetsRequest_PatchExpiringTargetsRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_patch_expiring_targets_request_patch_expiring_targets_request_instructions
  }
  property {
    name  = "patchExpiringTargets_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_targets_project_key
  }
  property {
    name  = "patchExpiringUserTargetsForSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_for_segment_environment_key
  }
  property {
    name  = "patchExpiringUserTargetsForSegment_patchExpiringUserTargetsForSegmentRequest_PatchExpiringUserTargetsForSegmentRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_for_segment_patch_expiring_user_targets_for_segment_request_patch_expiring_user_targets_for_segment_request_comment
  }
  property {
    name  = "patchExpiringUserTargetsForSegment_patchExpiringUserTargetsForSegmentRequest_PatchExpiringUserTargetsForSegmentRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_for_segment_patch_expiring_user_targets_for_segment_request_patch_expiring_user_targets_for_segment_request_instructions
  }
  property {
    name  = "patchExpiringUserTargetsForSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_for_segment_project_key
  }
  property {
    name  = "patchExpiringUserTargetsForSegment_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_for_segment_segment_key
  }
  property {
    name  = "patchExpiringUserTargets_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_environment_key
  }
  property {
    name  = "patchExpiringUserTargets_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_feature_flag_key
  }
  property {
    name  = "patchExpiringUserTargets_patchExpiringTargetsRequest_PatchExpiringTargetsRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_patch_expiring_targets_request_patch_expiring_targets_request_comment
  }
  property {
    name  = "patchExpiringUserTargets_patchExpiringTargetsRequest_PatchExpiringTargetsRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_patch_expiring_targets_request_patch_expiring_targets_request_instructions
  }
  property {
    name  = "patchExpiringUserTargets_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_expiring_user_targets_project_key
  }
  property {
    name  = "patchFeatureFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_feature_flag_feature_flag_key
  }
  property {
    name  = "patchFeatureFlag_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_feature_flag_patch_relay_auto_config_request_patch_relay_auto_config_request_comment
  }
  property {
    name  = "patchFeatureFlag_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_patch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_feature_flag_patch_relay_auto_config_request_patch_relay_auto_config_request_patch
  }
  property {
    name  = "patchFeatureFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_feature_flag_project_key
  }
  property {
    name  = "patchFlagConfigScheduledChange_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_config_scheduled_change_environment_key
  }
  property {
    name  = "patchFlagConfigScheduledChange_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_config_scheduled_change_feature_flag_key
  }
  property {
    name  = "patchFlagConfigScheduledChange_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_config_scheduled_change_id
  }
  property {
    name  = "patchFlagConfigScheduledChange_ignore_conflicts"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_config_scheduled_change_ignore_conflicts
  }
  property {
    name  = "patchFlagConfigScheduledChange_patchFlagConfigScheduledChangeRequest_PatchFlagConfigScheduledChangeRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_config_scheduled_change_patch_flag_config_scheduled_change_request_patch_flag_config_scheduled_change_request_comment
  }
  property {
    name  = "patchFlagConfigScheduledChange_patchFlagConfigScheduledChangeRequest_PatchFlagConfigScheduledChangeRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_config_scheduled_change_patch_flag_config_scheduled_change_request_patch_flag_config_scheduled_change_request_instructions
  }
  property {
    name  = "patchFlagConfigScheduledChange_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_config_scheduled_change_project_key
  }
  property {
    name  = "patchFlagDefaultsByProject_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_defaults_by_project_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchFlagDefaultsByProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_flag_defaults_by_project_project_key
  }
  property {
    name  = "patchInsightGroup_insight_group_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_insight_group_insight_group_key
  }
  property {
    name  = "patchInsightGroup_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_insight_group_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchIntegrationDeliveryConfiguration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_integration_delivery_configuration_environment_key
  }
  property {
    name  = "patchIntegrationDeliveryConfiguration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_integration_delivery_configuration_id
  }
  property {
    name  = "patchIntegrationDeliveryConfiguration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_integration_delivery_configuration_integration_key
  }
  property {
    name  = "patchIntegrationDeliveryConfiguration_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_integration_delivery_configuration_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchIntegrationDeliveryConfiguration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_integration_delivery_configuration_project_key
  }
  property {
    name  = "patchMember_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_member_id
  }
  property {
    name  = "patchMember_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_member_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchMembers_patchMembersRequest_PatchMembersRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_members_patch_members_request_patch_members_request_comment
  }
  property {
    name  = "patchMembers_patchMembersRequest_PatchMembersRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_members_patch_members_request_patch_members_request_instructions
  }
  property {
    name  = "patchMetricGroup_metric_group_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_metric_group_metric_group_key
  }
  property {
    name  = "patchMetricGroup_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_metric_group_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchMetricGroup_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_metric_group_project_key
  }
  property {
    name  = "patchMetric_metric_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_metric_metric_key
  }
  property {
    name  = "patchMetric_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_metric_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchMetric_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_metric_project_key
  }
  property {
    name  = "patchOAuthClient_client_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_oauth_client_client_id
  }
  property {
    name  = "patchOAuthClient_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_oauth_client_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchProject_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_project_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_project_project_key
  }
  property {
    name  = "patchRelayAutoConfig_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_relay_auto_config_id
  }
  property {
    name  = "patchRelayAutoConfig_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_relay_auto_config_patch_relay_auto_config_request_patch_relay_auto_config_request_comment
  }
  property {
    name  = "patchRelayAutoConfig_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_patch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_relay_auto_config_patch_relay_auto_config_request_patch_relay_auto_config_request_patch
  }
  property {
    name  = "patchReleaseByFlagKey_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_release_by_flag_key_flag_key
  }
  property {
    name  = "patchReleaseByFlagKey_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_release_by_flag_key_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchReleaseByFlagKey_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_release_by_flag_key_project_key
  }
  property {
    name  = "patchReleasePipeline_pipeline_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_release_pipeline_pipeline_key
  }
  property {
    name  = "patchReleasePipeline_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_release_pipeline_project_key
  }
  property {
    name  = "patchRepository_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_repository_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchRepository_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_repository_repo
  }
  property {
    name  = "patchSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_segment_environment_key
  }
  property {
    name  = "patchSegment_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_segment_patch_relay_auto_config_request_patch_relay_auto_config_request_comment
  }
  property {
    name  = "patchSegment_patchRelayAutoConfigRequest_PatchRelayAutoConfigRequest_patch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_segment_patch_relay_auto_config_request_patch_relay_auto_config_request_patch
  }
  property {
    name  = "patchSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_segment_project_key
  }
  property {
    name  = "patchSegment_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_segment_segment_key
  }
  property {
    name  = "patchTeam_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_team_expand
  }
  property {
    name  = "patchTeam_patchTeamsRequest_PatchTeamsRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_team_patch_teams_request_patch_teams_request_comment
  }
  property {
    name  = "patchTeam_patchTeamsRequest_PatchTeamsRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_team_patch_teams_request_patch_teams_request_instructions
  }
  property {
    name  = "patchTeam_team_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_team_team_key
  }
  property {
    name  = "patchTeams_patchTeamsRequest_PatchTeamsRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_teams_patch_teams_request_patch_teams_request_comment
  }
  property {
    name  = "patchTeams_patchTeamsRequest_PatchTeamsRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_teams_patch_teams_request_patch_teams_request_instructions
  }
  property {
    name  = "patchToken_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_token_id
  }
  property {
    name  = "patchToken_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_token_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "patchTriggerWorkflow_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_trigger_workflow_environment_key
  }
  property {
    name  = "patchTriggerWorkflow_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_trigger_workflow_feature_flag_key
  }
  property {
    name  = "patchTriggerWorkflow_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_trigger_workflow_id
  }
  property {
    name  = "patchTriggerWorkflow_patchTriggerWorkflowRequest_PatchTriggerWorkflowRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_trigger_workflow_patch_trigger_workflow_request_patch_trigger_workflow_request_comment
  }
  property {
    name  = "patchTriggerWorkflow_patchTriggerWorkflowRequest_PatchTriggerWorkflowRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_trigger_workflow_patch_trigger_workflow_request_patch_trigger_workflow_request_instructions
  }
  property {
    name  = "patchTriggerWorkflow_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_trigger_workflow_project_key
  }
  property {
    name  = "patchWebhook_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_webhook_id
  }
  property {
    name  = "patchWebhook_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_patch_webhook_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "postApprovalRequestApplyForFlag_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_apply_for_flag_environment_key
  }
  property {
    name  = "postApprovalRequestApplyForFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_apply_for_flag_feature_flag_key
  }
  property {
    name  = "postApprovalRequestApplyForFlag_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_apply_for_flag_id
  }
  property {
    name  = "postApprovalRequestApplyForFlag_postApprovalRequestApplyRequest_PostApprovalRequestApplyRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_apply_for_flag_post_approval_request_apply_request_post_approval_request_apply_request_comment
  }
  property {
    name  = "postApprovalRequestApplyForFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_apply_for_flag_project_key
  }
  property {
    name  = "postApprovalRequestApply_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_apply_id
  }
  property {
    name  = "postApprovalRequestApply_postApprovalRequestApplyRequest_PostApprovalRequestApplyRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_apply_post_approval_request_apply_request_post_approval_request_apply_request_comment
  }
  property {
    name  = "postApprovalRequestForFlag_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_environment_key
  }
  property {
    name  = "postApprovalRequestForFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_feature_flag_key
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_comment
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_description
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_executionDate"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_execution_date
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_instructions
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_integrationConfig"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_integration_config
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_notifyMemberIds"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_notify_member_ids
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_notifyTeamKeys"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_notify_team_keys
  }
  property {
    name  = "postApprovalRequestForFlag_postApprovalRequestForFlagRequest_PostApprovalRequestForFlagRequest_operatingOnId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_post_approval_request_for_flag_request_post_approval_request_for_flag_request_operating_on_id
  }
  property {
    name  = "postApprovalRequestForFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_for_flag_project_key
  }
  property {
    name  = "postApprovalRequestReviewForFlag_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_for_flag_environment_key
  }
  property {
    name  = "postApprovalRequestReviewForFlag_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_for_flag_feature_flag_key
  }
  property {
    name  = "postApprovalRequestReviewForFlag_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_for_flag_id
  }
  property {
    name  = "postApprovalRequestReviewForFlag_postApprovalRequestReviewRequest_PostApprovalRequestReviewRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_for_flag_post_approval_request_review_request_post_approval_request_review_request_comment
  }
  property {
    name  = "postApprovalRequestReviewForFlag_postApprovalRequestReviewRequest_PostApprovalRequestReviewRequest_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_for_flag_post_approval_request_review_request_post_approval_request_review_request_kind
  }
  property {
    name  = "postApprovalRequestReviewForFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_for_flag_project_key
  }
  property {
    name  = "postApprovalRequestReview_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_id
  }
  property {
    name  = "postApprovalRequestReview_postApprovalRequestReviewRequest_PostApprovalRequestReviewRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_post_approval_request_review_request_post_approval_request_review_request_comment
  }
  property {
    name  = "postApprovalRequestReview_postApprovalRequestReviewRequest_PostApprovalRequestReviewRequest_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_review_post_approval_request_review_request_post_approval_request_review_request_kind
  }
  property {
    name  = "postApprovalRequest_postApprovalRequestRequest_PostApprovalRequestRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_post_approval_request_request_post_approval_request_request_comment
  }
  property {
    name  = "postApprovalRequest_postApprovalRequestRequest_PostApprovalRequestRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_post_approval_request_request_post_approval_request_request_description
  }
  property {
    name  = "postApprovalRequest_postApprovalRequestRequest_PostApprovalRequestRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_post_approval_request_request_post_approval_request_request_instructions
  }
  property {
    name  = "postApprovalRequest_postApprovalRequestRequest_PostApprovalRequestRequest_integrationConfig"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_post_approval_request_request_post_approval_request_request_integration_config
  }
  property {
    name  = "postApprovalRequest_postApprovalRequestRequest_PostApprovalRequestRequest_notifyMemberIds"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_post_approval_request_request_post_approval_request_request_notify_member_ids
  }
  property {
    name  = "postApprovalRequest_postApprovalRequestRequest_PostApprovalRequestRequest_notifyTeamKeys"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_post_approval_request_request_post_approval_request_request_notify_team_keys
  }
  property {
    name  = "postApprovalRequest_postApprovalRequestRequest_PostApprovalRequestRequest_resourceId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_approval_request_post_approval_request_request_post_approval_request_request_resource_id
  }
  property {
    name  = "postCustomRole_postCustomRoleRequest_PostCustomRoleRequest_basePermissions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_custom_role_post_custom_role_request_post_custom_role_request_base_permissions
  }
  property {
    name  = "postCustomRole_postCustomRoleRequest_PostCustomRoleRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_custom_role_post_custom_role_request_post_custom_role_request_description
  }
  property {
    name  = "postCustomRole_postCustomRoleRequest_PostCustomRoleRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_custom_role_post_custom_role_request_post_custom_role_request_key
  }
  property {
    name  = "postCustomRole_postCustomRoleRequest_PostCustomRoleRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_custom_role_post_custom_role_request_post_custom_role_request_name
  }
  property {
    name  = "postCustomRole_postCustomRoleRequest_PostCustomRoleRequest_policy"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_custom_role_post_custom_role_request_post_custom_role_request_policy
  }
  property {
    name  = "postDestination_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_destination_environment_key
  }
  property {
    name  = "postDestination_postDestinationRequest_PostDestinationRequest_config"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_destination_post_destination_request_post_destination_request_config
  }
  property {
    name  = "postDestination_postDestinationRequest_PostDestinationRequest_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_destination_post_destination_request_post_destination_request_kind
  }
  property {
    name  = "postDestination_postDestinationRequest_PostDestinationRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_destination_post_destination_request_post_destination_request_name
  }
  property {
    name  = "postDestination_postDestinationRequest_PostDestinationRequest_on"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_destination_post_destination_request_post_destination_request_on
  }
  property {
    name  = "postDestination_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_destination_project_key
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInnerSource_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_source_key
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInnerSource_version"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_source_version
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_color"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_color
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_confirmChanges"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_confirm_changes
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_critical"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_critical
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_defaultTrackEvents"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_default_track_events
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_defaultTtl"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_default_ttl
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_key
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_name
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_requireComments"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_require_comments
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_secureMode"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_secure_mode
  }
  property {
    name  = "postEnvironment_postProjectRequestEnvironmentsInner_PostProjectRequestEnvironmentsInner_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_post_project_request_environments_inner_post_project_request_environments_inner_tags
  }
  property {
    name  = "postEnvironment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_environment_project_key
  }
  property {
    name  = "postExtinction_branch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_extinction_branch
  }
  property {
    name  = "postExtinction_get_extinctions200_response_items_value_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_extinction_get_extinctions200_response_items_value_inner
  }
  property {
    name  = "postExtinction_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_extinction_repo
  }
  property {
    name  = "postFeatureFlag_clone"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_clone
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequestClientSideAvailability_usingEnvironmentId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_client_side_availability_using_environment_id
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequestClientSideAvailability_usingMobileKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_client_side_availability_using_mobile_key
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequestDefaults_offVariation"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_defaults_off_variation
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequestDefaults_onVariation"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_defaults_on_variation
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequestMigrationSettings_contextKind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_migration_settings_context_kind
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequestMigrationSettings_stageCount"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_migration_settings_stage_count
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_customProperties"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_custom_properties
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_description
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_includeInSnippet"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_include_in_snippet
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_key
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_maintainerId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_maintainer_id
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_maintainerTeamKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_maintainer_team_key
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_name
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_purpose"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_purpose
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_tags
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_temporary"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_temporary
  }
  property {
    name  = "postFeatureFlag_postFeatureFlagRequest_PostFeatureFlagRequest_variations"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_post_feature_flag_request_post_feature_flag_request_variations
  }
  property {
    name  = "postFeatureFlag_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_feature_flag_project_key
  }
  property {
    name  = "postFlagConfigScheduledChanges_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_config_scheduled_changes_environment_key
  }
  property {
    name  = "postFlagConfigScheduledChanges_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_config_scheduled_changes_feature_flag_key
  }
  property {
    name  = "postFlagConfigScheduledChanges_ignore_conflicts"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_config_scheduled_changes_ignore_conflicts
  }
  property {
    name  = "postFlagConfigScheduledChanges_postFlagConfigScheduledChangesRequest_PostFlagConfigScheduledChangesRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_config_scheduled_changes_post_flag_config_scheduled_changes_request_post_flag_config_scheduled_changes_request_comment
  }
  property {
    name  = "postFlagConfigScheduledChanges_postFlagConfigScheduledChangesRequest_PostFlagConfigScheduledChangesRequest_executionDate"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_config_scheduled_changes_post_flag_config_scheduled_changes_request_post_flag_config_scheduled_changes_request_execution_date
  }
  property {
    name  = "postFlagConfigScheduledChanges_postFlagConfigScheduledChangesRequest_PostFlagConfigScheduledChangesRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_config_scheduled_changes_post_flag_config_scheduled_changes_request_post_flag_config_scheduled_changes_request_instructions
  }
  property {
    name  = "postFlagConfigScheduledChanges_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_config_scheduled_changes_project_key
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_environment_key
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_feature_flag_key
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequestSource_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_source_key
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequestSource_version"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_source_version
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_comment
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_description
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequest_excludedActions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_excluded_actions
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequest_includedActions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_included_actions
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequest_notifyMemberIds"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_notify_member_ids
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_postFlagCopyConfigApprovalRequestRequest_PostFlagCopyConfigApprovalRequestRequest_notifyTeamKeys"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_post_flag_copy_config_approval_request_request_post_flag_copy_config_approval_request_request_notify_team_keys
  }
  property {
    name  = "postFlagCopyConfigApprovalRequest_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_flag_copy_config_approval_request_project_key
  }
  property {
    name  = "postMemberTeams_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_member_teams_id
  }
  property {
    name  = "postMemberTeams_postMemberTeamsRequest_PostMemberTeamsRequest_teamKeys"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_member_teams_post_member_teams_request_post_member_teams_request_team_keys
  }
  property {
    name  = "postMembers_post_members_request_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_members_post_members_request_inner
  }
  property {
    name  = "postMetric_postMetricRequest_GetFeatureFlags200ResponseItemsInnerExperimentsItemsInnerMetricEventDefault_disabled"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_get_feature_flags200_response_items_inner_experiments_items_inner_metric_event_default_disabled
  }
  property {
    name  = "postMetric_postMetricRequest_GetFeatureFlags200ResponseItemsInnerExperimentsItemsInnerMetricEventDefault_value"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_get_feature_flags200_response_items_inner_experiments_items_inner_metric_event_default_value
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_analysisType"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_analysis_type
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_description
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_eventKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_event_key
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_isActive"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_is_active
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_isNumeric"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_is_numeric
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_key
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_kind
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_name
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_percentileValue"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_percentile_value
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_randomizationUnits"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_randomization_units
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_selector"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_selector
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_successCriteria"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_success_criteria
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_tags
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_unit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_unit
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_unitAggregationType"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_unit_aggregation_type
  }
  property {
    name  = "postMetric_postMetricRequest_PostMetricRequest_urls"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_post_metric_request_post_metric_request_urls
  }
  property {
    name  = "postMetric_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_metric_project_key
  }
  property {
    name  = "postMigrationSafetyIssues_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_migration_safety_issues_environment_key
  }
  property {
    name  = "postMigrationSafetyIssues_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_migration_safety_issues_flag_key
  }
  property {
    name  = "postMigrationSafetyIssues_postMigrationSafetyIssuesRequest_PostMigrationSafetyIssuesRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_migration_safety_issues_post_migration_safety_issues_request_post_migration_safety_issues_request_comment
  }
  property {
    name  = "postMigrationSafetyIssues_postMigrationSafetyIssuesRequest_PostMigrationSafetyIssuesRequest_instructions"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_migration_safety_issues_post_migration_safety_issues_request_post_migration_safety_issues_request_instructions
  }
  property {
    name  = "postMigrationSafetyIssues_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_migration_safety_issues_project_key
  }
  property {
    name  = "postProject_postProjectRequest_PostProjectRequestDefaultClientSideAvailability_usingEnvironmentId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_project_post_project_request_post_project_request_default_client_side_availability_using_environment_id
  }
  property {
    name  = "postProject_postProjectRequest_PostProjectRequestDefaultClientSideAvailability_usingMobileKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_project_post_project_request_post_project_request_default_client_side_availability_using_mobile_key
  }
  property {
    name  = "postProject_postProjectRequest_PostProjectRequest_environments"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_project_post_project_request_post_project_request_environments
  }
  property {
    name  = "postProject_postProjectRequest_PostProjectRequest_includeInSnippetByDefault"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_project_post_project_request_post_project_request_include_in_snippet_by_default
  }
  property {
    name  = "postProject_postProjectRequest_PostProjectRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_project_post_project_request_post_project_request_key
  }
  property {
    name  = "postProject_postProjectRequest_PostProjectRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_project_post_project_request_post_project_request_name
  }
  property {
    name  = "postProject_postProjectRequest_PostProjectRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_project_post_project_request_post_project_request_tags
  }
  property {
    name  = "postRelayAutoConfig_postRelayAutoConfigRequest_PostRelayAutoConfigRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_relay_auto_config_post_relay_auto_config_request_post_relay_auto_config_request_name
  }
  property {
    name  = "postRelayAutoConfig_postRelayAutoConfigRequest_PostRelayAutoConfigRequest_policy"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_relay_auto_config_post_relay_auto_config_request_post_relay_auto_config_request_policy
  }
  property {
    name  = "postReleasePipeline_postReleasePipelineRequest_PostReleasePipelineRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_release_pipeline_post_release_pipeline_request_post_release_pipeline_request_description
  }
  property {
    name  = "postReleasePipeline_postReleasePipelineRequest_PostReleasePipelineRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_release_pipeline_post_release_pipeline_request_post_release_pipeline_request_key
  }
  property {
    name  = "postReleasePipeline_postReleasePipelineRequest_PostReleasePipelineRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_release_pipeline_post_release_pipeline_request_post_release_pipeline_request_name
  }
  property {
    name  = "postReleasePipeline_postReleasePipelineRequest_PostReleasePipelineRequest_phases"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_release_pipeline_post_release_pipeline_request_post_release_pipeline_request_phases
  }
  property {
    name  = "postReleasePipeline_postReleasePipelineRequest_PostReleasePipelineRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_release_pipeline_post_release_pipeline_request_post_release_pipeline_request_tags
  }
  property {
    name  = "postReleasePipeline_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_release_pipeline_project_key
  }
  property {
    name  = "postRepository_postRepositoryRequest_PostRepositoryRequest_commitUrlTemplate"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_repository_post_repository_request_post_repository_request_commit_url_template
  }
  property {
    name  = "postRepository_postRepositoryRequest_PostRepositoryRequest_defaultBranch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_repository_post_repository_request_post_repository_request_default_branch
  }
  property {
    name  = "postRepository_postRepositoryRequest_PostRepositoryRequest_hunkUrlTemplate"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_repository_post_repository_request_post_repository_request_hunk_url_template
  }
  property {
    name  = "postRepository_postRepositoryRequest_PostRepositoryRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_repository_post_repository_request_post_repository_request_name
  }
  property {
    name  = "postRepository_postRepositoryRequest_PostRepositoryRequest_sourceLink"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_repository_post_repository_request_post_repository_request_source_link
  }
  property {
    name  = "postRepository_postRepositoryRequest_PostRepositoryRequest_type"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_repository_post_repository_request_post_repository_request_type
  }
  property {
    name  = "postSegment_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_environment_key
  }
  property {
    name  = "postSegment_postSegmentRequest_PostSegmentRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_post_segment_request_post_segment_request_description
  }
  property {
    name  = "postSegment_postSegmentRequest_PostSegmentRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_post_segment_request_post_segment_request_key
  }
  property {
    name  = "postSegment_postSegmentRequest_PostSegmentRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_post_segment_request_post_segment_request_name
  }
  property {
    name  = "postSegment_postSegmentRequest_PostSegmentRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_post_segment_request_post_segment_request_tags
  }
  property {
    name  = "postSegment_postSegmentRequest_PostSegmentRequest_unbounded"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_post_segment_request_post_segment_request_unbounded
  }
  property {
    name  = "postSegment_postSegmentRequest_PostSegmentRequest_unboundedContextKind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_post_segment_request_post_segment_request_unbounded_context_kind
  }
  property {
    name  = "postSegment_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_segment_project_key
  }
  property {
    name  = "postTeamMembers_file"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_members_file
  }
  property {
    name  = "postTeamMembers_team_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_members_team_key
  }
  property {
    name  = "postTeam_expand"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_expand
  }
  property {
    name  = "postTeam_postTeamRequest_PostTeamRequest_customRoleKeys"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_post_team_request_post_team_request_custom_role_keys
  }
  property {
    name  = "postTeam_postTeamRequest_PostTeamRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_post_team_request_post_team_request_description
  }
  property {
    name  = "postTeam_postTeamRequest_PostTeamRequest_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_post_team_request_post_team_request_key
  }
  property {
    name  = "postTeam_postTeamRequest_PostTeamRequest_memberIDs"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_post_team_request_post_team_request_member_ids
  }
  property {
    name  = "postTeam_postTeamRequest_PostTeamRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_post_team_request_post_team_request_name
  }
  property {
    name  = "postTeam_postTeamRequest_PostTeamRequest_permissionGrants"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_team_post_team_request_post_team_request_permission_grants
  }
  property {
    name  = "postToken_postTokenRequest_PostTokenRequest_customRoleIds"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_token_post_token_request_post_token_request_custom_role_ids
  }
  property {
    name  = "postToken_postTokenRequest_PostTokenRequest_defaultApiVersion"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_token_post_token_request_post_token_request_default_api_version
  }
  property {
    name  = "postToken_postTokenRequest_PostTokenRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_token_post_token_request_post_token_request_description
  }
  property {
    name  = "postToken_postTokenRequest_PostTokenRequest_inlineRole"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_token_post_token_request_post_token_request_inline_role
  }
  property {
    name  = "postToken_postTokenRequest_PostTokenRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_token_post_token_request_post_token_request_name
  }
  property {
    name  = "postToken_postTokenRequest_PostTokenRequest_role"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_token_post_token_request_post_token_request_role
  }
  property {
    name  = "postToken_postTokenRequest_PostTokenRequest_serviceToken"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_token_post_token_request_post_token_request_service_token
  }
  property {
    name  = "postWebhook_postWebhookRequest_PostWebhookRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_webhook_post_webhook_request_post_webhook_request_name
  }
  property {
    name  = "postWebhook_postWebhookRequest_PostWebhookRequest_on"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_webhook_post_webhook_request_post_webhook_request_on
  }
  property {
    name  = "postWebhook_postWebhookRequest_PostWebhookRequest_secret"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_webhook_post_webhook_request_post_webhook_request_secret
  }
  property {
    name  = "postWebhook_postWebhookRequest_PostWebhookRequest_sign"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_webhook_post_webhook_request_post_webhook_request_sign
  }
  property {
    name  = "postWebhook_postWebhookRequest_PostWebhookRequest_statements"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_webhook_post_webhook_request_post_webhook_request_statements
  }
  property {
    name  = "postWebhook_postWebhookRequest_PostWebhookRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_webhook_post_webhook_request_post_webhook_request_tags
  }
  property {
    name  = "postWebhook_postWebhookRequest_PostWebhookRequest_url"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_webhook_post_webhook_request_post_webhook_request_url
  }
  property {
    name  = "postWorkflow_dry_run"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_dry_run
  }
  property {
    name  = "postWorkflow_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_environment_key
  }
  property {
    name  = "postWorkflow_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_feature_flag_key
  }
  property {
    name  = "postWorkflow_postWorkflowRequest_PostWorkflowRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_post_workflow_request_post_workflow_request_description
  }
  property {
    name  = "postWorkflow_postWorkflowRequest_PostWorkflowRequest_maintainerId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_post_workflow_request_post_workflow_request_maintainer_id
  }
  property {
    name  = "postWorkflow_postWorkflowRequest_PostWorkflowRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_post_workflow_request_post_workflow_request_name
  }
  property {
    name  = "postWorkflow_postWorkflowRequest_PostWorkflowRequest_stages"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_post_workflow_request_post_workflow_request_stages
  }
  property {
    name  = "postWorkflow_postWorkflowRequest_PostWorkflowRequest_templateKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_post_workflow_request_post_workflow_request_template_key
  }
  property {
    name  = "postWorkflow_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_project_key
  }
  property {
    name  = "postWorkflow_template_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_post_workflow_template_key
  }
  property {
    name  = "putBranch_branch"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_branch
  }
  property {
    name  = "putBranch_putBranchRequest_PutBranchRequest_commitTime"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_put_branch_request_put_branch_request_commit_time
  }
  property {
    name  = "putBranch_putBranchRequest_PutBranchRequest_head"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_put_branch_request_put_branch_request_head
  }
  property {
    name  = "putBranch_putBranchRequest_PutBranchRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_put_branch_request_put_branch_request_name
  }
  property {
    name  = "putBranch_putBranchRequest_PutBranchRequest_references"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_put_branch_request_put_branch_request_references
  }
  property {
    name  = "putBranch_putBranchRequest_PutBranchRequest_syncTime"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_put_branch_request_put_branch_request_sync_time
  }
  property {
    name  = "putBranch_putBranchRequest_PutBranchRequest_updateSequenceId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_put_branch_request_put_branch_request_update_sequence_id
  }
  property {
    name  = "putBranch_repo"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_branch_repo
  }
  property {
    name  = "putContextFlagSetting_context_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_flag_setting_context_key
  }
  property {
    name  = "putContextFlagSetting_context_kind"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_flag_setting_context_kind
  }
  property {
    name  = "putContextFlagSetting_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_flag_setting_environment_key
  }
  property {
    name  = "putContextFlagSetting_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_flag_setting_feature_flag_key
  }
  property {
    name  = "putContextFlagSetting_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_flag_setting_project_key
  }
  property {
    name  = "putContextFlagSetting_putContextFlagSettingRequest_PutContextFlagSettingRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_flag_setting_put_context_flag_setting_request_put_context_flag_setting_request_comment
  }
  property {
    name  = "putContextFlagSetting_putContextFlagSettingRequest_PutContextFlagSettingRequest_setting"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_flag_setting_put_context_flag_setting_request_put_context_flag_setting_request_setting
  }
  property {
    name  = "putContextKind_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_kind_key
  }
  property {
    name  = "putContextKind_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_kind_project_key
  }
  property {
    name  = "putContextKind_putContextKindRequest_PutContextKindRequest_archived"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_kind_put_context_kind_request_put_context_kind_request_archived
  }
  property {
    name  = "putContextKind_putContextKindRequest_PutContextKindRequest_description"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_kind_put_context_kind_request_put_context_kind_request_description
  }
  property {
    name  = "putContextKind_putContextKindRequest_PutContextKindRequest_hideInTargeting"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_kind_put_context_kind_request_put_context_kind_request_hide_in_targeting
  }
  property {
    name  = "putContextKind_putContextKindRequest_PutContextKindRequest_name"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_kind_put_context_kind_request_put_context_kind_request_name
  }
  property {
    name  = "putContextKind_putContextKindRequest_PutContextKindRequest_version"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_context_kind_put_context_kind_request_put_context_kind_request_version
  }
  property {
    name  = "putExperimentationSettings_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_experimentation_settings_project_key
  }
  property {
    name  = "putExperimentationSettings_putExperimentationSettingsRequest_PutExperimentationSettingsRequest_randomizationUnits"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_experimentation_settings_put_experimentation_settings_request_put_experimentation_settings_request_randomization_units
  }
  property {
    name  = "putFlagDefaultsByProject_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_project_key
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestBooleanDefaults_falseDescription"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_boolean_defaults_false_description
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestBooleanDefaults_falseDisplayName"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_boolean_defaults_false_display_name
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestBooleanDefaults_offVariation"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_boolean_defaults_off_variation
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestBooleanDefaults_onVariation"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_boolean_defaults_on_variation
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestBooleanDefaults_trueDescription"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_boolean_defaults_true_description
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestBooleanDefaults_trueDisplayName"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_boolean_defaults_true_display_name
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestDefaultClientSideAvailability_usingEnvironmentId"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_default_client_side_availability_using_environment_id
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequestDefaultClientSideAvailability_usingMobileKey"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_default_client_side_availability_using_mobile_key
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequest_tags"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_tags
  }
  property {
    name  = "putFlagDefaultsByProject_putFlagDefaultsByProjectRequest_PutFlagDefaultsByProjectRequest_temporary"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_defaults_by_project_put_flag_defaults_by_project_request_put_flag_defaults_by_project_request_temporary
  }
  property {
    name  = "putFlagFollowers_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_followers_environment_key
  }
  property {
    name  = "putFlagFollowers_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_followers_feature_flag_key
  }
  property {
    name  = "putFlagFollowers_member_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_followers_member_id
  }
  property {
    name  = "putFlagFollowers_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_followers_project_key
  }
  property {
    name  = "putFlagSetting_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_setting_environment_key
  }
  property {
    name  = "putFlagSetting_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_setting_feature_flag_key
  }
  property {
    name  = "putFlagSetting_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_setting_project_key
  }
  property {
    name  = "putFlagSetting_putContextFlagSettingRequest_PutContextFlagSettingRequest_comment"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_setting_put_context_flag_setting_request_put_context_flag_setting_request_comment
  }
  property {
    name  = "putFlagSetting_putContextFlagSettingRequest_PutContextFlagSettingRequest_setting"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_setting_put_context_flag_setting_request_put_context_flag_setting_request_setting
  }
  property {
    name  = "putFlagSetting_user_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_put_flag_setting_user_key
  }
  property {
    name  = "resetEnvironmentMobileKey_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_environment_mobile_key_environment_key
  }
  property {
    name  = "resetEnvironmentMobileKey_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_environment_mobile_key_project_key
  }
  property {
    name  = "resetEnvironmentSDKKey_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_environment_sdkkey_environment_key
  }
  property {
    name  = "resetEnvironmentSDKKey_expiry"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_environment_sdkkey_expiry
  }
  property {
    name  = "resetEnvironmentSDKKey_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_environment_sdkkey_project_key
  }
  property {
    name  = "resetRelayAutoConfig_expiry"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_relay_auto_config_expiry
  }
  property {
    name  = "resetRelayAutoConfig_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_relay_auto_config_id
  }
  property {
    name  = "resetToken_expiry"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_token_expiry
  }
  property {
    name  = "resetToken_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_reset_token_id
  }
  property {
    name  = "searchContextInstances_continuation_token"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_continuation_token
  }
  property {
    name  = "searchContextInstances_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_environment_key
  }
  property {
    name  = "searchContextInstances_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_filter
  }
  property {
    name  = "searchContextInstances_include_total_count"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_include_total_count
  }
  property {
    name  = "searchContextInstances_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_limit
  }
  property {
    name  = "searchContextInstances_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_project_key
  }
  property {
    name  = "searchContextInstances_searchContextInstancesRequest_SearchContextInstancesRequest_continuationToken"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_search_context_instances_request_search_context_instances_request_continuation_token
  }
  property {
    name  = "searchContextInstances_searchContextInstancesRequest_SearchContextInstancesRequest_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_search_context_instances_request_search_context_instances_request_filter
  }
  property {
    name  = "searchContextInstances_searchContextInstancesRequest_SearchContextInstancesRequest_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_search_context_instances_request_search_context_instances_request_limit
  }
  property {
    name  = "searchContextInstances_searchContextInstancesRequest_SearchContextInstancesRequest_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_search_context_instances_request_search_context_instances_request_sort
  }
  property {
    name  = "searchContextInstances_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_context_instances_sort
  }
  property {
    name  = "searchContexts_continuation_token"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_continuation_token
  }
  property {
    name  = "searchContexts_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_environment_key
  }
  property {
    name  = "searchContexts_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_filter
  }
  property {
    name  = "searchContexts_include_total_count"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_include_total_count
  }
  property {
    name  = "searchContexts_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_limit
  }
  property {
    name  = "searchContexts_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_project_key
  }
  property {
    name  = "searchContexts_searchContextsRequest_SearchContextsRequest_continuationToken"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_search_contexts_request_search_contexts_request_continuation_token
  }
  property {
    name  = "searchContexts_searchContextsRequest_SearchContextsRequest_filter"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_search_contexts_request_search_contexts_request_filter
  }
  property {
    name  = "searchContexts_searchContextsRequest_SearchContextsRequest_limit"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_search_contexts_request_search_contexts_request_limit
  }
  property {
    name  = "searchContexts_searchContextsRequest_SearchContextsRequest_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_search_contexts_request_search_contexts_request_sort
  }
  property {
    name  = "searchContexts_sort"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_search_contexts_sort
  }
  property {
    name  = "updateBigSegmentContextTargets_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_context_targets_environment_key
  }
  property {
    name  = "updateBigSegmentContextTargets_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_context_targets_project_key
  }
  property {
    name  = "updateBigSegmentContextTargets_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_context_targets_segment_key
  }
  property {
    name  = "updateBigSegmentContextTargets_updateBigSegmentContextTargetsRequest_UpdateBigSegmentContextTargetsRequestIncluded_add"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_context_targets_update_big_segment_context_targets_request_update_big_segment_context_targets_request_included_add
  }
  property {
    name  = "updateBigSegmentContextTargets_updateBigSegmentContextTargetsRequest_UpdateBigSegmentContextTargetsRequestIncluded_remove"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_context_targets_update_big_segment_context_targets_request_update_big_segment_context_targets_request_included_remove
  }
  property {
    name  = "updateBigSegmentTargets_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_targets_environment_key
  }
  property {
    name  = "updateBigSegmentTargets_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_targets_project_key
  }
  property {
    name  = "updateBigSegmentTargets_segment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_targets_segment_key
  }
  property {
    name  = "updateBigSegmentTargets_updateBigSegmentContextTargetsRequest_UpdateBigSegmentContextTargetsRequestIncluded_add"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_targets_update_big_segment_context_targets_request_update_big_segment_context_targets_request_included_add
  }
  property {
    name  = "updateBigSegmentTargets_updateBigSegmentContextTargetsRequest_UpdateBigSegmentContextTargetsRequestIncluded_remove"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_big_segment_targets_update_big_segment_context_targets_request_update_big_segment_context_targets_request_included_remove
  }
  property {
    name  = "updateDeployment_deployment_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_deployment_deployment_id
  }
  property {
    name  = "updateDeployment_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_deployment_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "updateFlagLink_feature_flag_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_flag_link_feature_flag_key
  }
  property {
    name  = "updateFlagLink_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_flag_link_id
  }
  property {
    name  = "updateFlagLink_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_flag_link_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "updateFlagLink_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_flag_link_project_key
  }
  property {
    name  = "updateSubscription_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_subscription_id
  }
  property {
    name  = "updateSubscription_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_subscription_integration_key
  }
  property {
    name  = "updateSubscription_patch_relay_auto_config_request_patch_inner"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_update_subscription_patch_relay_auto_config_request_patch_inner
  }
  property {
    name  = "validateIntegrationDeliveryConfiguration_environment_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_validate_integration_delivery_configuration_environment_key
  }
  property {
    name  = "validateIntegrationDeliveryConfiguration_id"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_validate_integration_delivery_configuration_id
  }
  property {
    name  = "validateIntegrationDeliveryConfiguration_integration_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_validate_integration_delivery_configuration_integration_key
  }
  property {
    name  = "validateIntegrationDeliveryConfiguration_project_key"
    type  = "string"
    value = var.connector-oai-launchdarklyrestapi_property_validate_integration_delivery_configuration_project_key
  }
}
