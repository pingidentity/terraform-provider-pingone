resource "pingone_davinci_connector_instance" "connector-oai-datadogapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-datadogapi"
  }
  name = "My awesome connector-oai-datadogapi"
  property {
    name  = "apiV1ApplicationKeyPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_application_key_post_accept
  }
  property {
    name  = "apiV1ApplicationKeyPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_application_key_post_body
  }
  property {
    name  = "apiV1ApplicationKeyPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_application_key_post_content_type
  }
  property {
    name  = "apiV1CheckRunPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_check_run_post_accept
  }
  property {
    name  = "apiV1CheckRunPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_check_run_post_body
  }
  property {
    name  = "apiV1CheckRunPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_check_run_post_content_type
  }
  property {
    name  = "apiV1DailyCustomReportsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_daily_custom_reports_get_accept
  }
  property {
    name  = "apiV1DailyCustomReportsGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_daily_custom_reports_get_page_number
  }
  property {
    name  = "apiV1DailyCustomReportsGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_daily_custom_reports_get_page_size
  }
  property {
    name  = "apiV1DailyCustomReportsGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_daily_custom_reports_get_sort
  }
  property {
    name  = "apiV1DailyCustomReportsGet_sort_dir"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_daily_custom_reports_get_sort_dir
  }
  property {
    name  = "apiV1DailyCustomReportsReportIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_daily_custom_reports_report_id_get_accept
  }
  property {
    name  = "apiV1DailyCustomReportsReportIdGet_report_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_daily_custom_reports_report_id_get_report_id
  }
  property {
    name  = "apiV1DashboardDashboardIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_delete_accept
  }
  property {
    name  = "apiV1DashboardDashboardIdDelete_dashboard_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_delete_dashboard_id
  }
  property {
    name  = "apiV1DashboardDashboardIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_get_accept
  }
  property {
    name  = "apiV1DashboardDashboardIdGet_dashboard_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_get_dashboard_id
  }
  property {
    name  = "apiV1DashboardDashboardIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_put_accept
  }
  property {
    name  = "apiV1DashboardDashboardIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_put_body
  }
  property {
    name  = "apiV1DashboardDashboardIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_put_content_type
  }
  property {
    name  = "apiV1DashboardDashboardIdPut_dashboard_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_dashboard_id_put_dashboard_id
  }
  property {
    name  = "apiV1DashboardDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_delete_accept
  }
  property {
    name  = "apiV1DashboardDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_delete_content_type
  }
  property {
    name  = "apiV1DashboardGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_get_accept
  }
  property {
    name  = "apiV1DashboardGet_filter_deleted"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_get_filter_deleted
  }
  property {
    name  = "apiV1DashboardGet_filter_shared"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_get_filter_shared
  }
  property {
    name  = "apiV1DashboardListsManualGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_get_accept
  }
  property {
    name  = "apiV1DashboardListsManualListIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_delete_accept
  }
  property {
    name  = "apiV1DashboardListsManualListIdDelete_list_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_delete_list_id
  }
  property {
    name  = "apiV1DashboardListsManualListIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_get_accept
  }
  property {
    name  = "apiV1DashboardListsManualListIdGet_list_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_get_list_id
  }
  property {
    name  = "apiV1DashboardListsManualListIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_put_accept
  }
  property {
    name  = "apiV1DashboardListsManualListIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_put_body
  }
  property {
    name  = "apiV1DashboardListsManualListIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_put_content_type
  }
  property {
    name  = "apiV1DashboardListsManualListIdPut_list_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_list_id_put_list_id
  }
  property {
    name  = "apiV1DashboardListsManualPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_post_accept
  }
  property {
    name  = "apiV1DashboardListsManualPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_post_body
  }
  property {
    name  = "apiV1DashboardListsManualPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_lists_manual_post_content_type
  }
  property {
    name  = "apiV1DashboardPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_patch_accept
  }
  property {
    name  = "apiV1DashboardPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_patch_body
  }
  property {
    name  = "apiV1DashboardPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_patch_content_type
  }
  property {
    name  = "apiV1DashboardPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_post_accept
  }
  property {
    name  = "apiV1DashboardPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_post_body
  }
  property {
    name  = "apiV1DashboardPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_post_content_type
  }
  property {
    name  = "apiV1DashboardPublicPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_post_accept
  }
  property {
    name  = "apiV1DashboardPublicPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_post_body
  }
  property {
    name  = "apiV1DashboardPublicPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_post_content_type
  }
  property {
    name  = "apiV1DashboardPublicTokenDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_delete_accept
  }
  property {
    name  = "apiV1DashboardPublicTokenDelete_token"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_delete_token
  }
  property {
    name  = "apiV1DashboardPublicTokenGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_get_accept
  }
  property {
    name  = "apiV1DashboardPublicTokenGet_token"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_get_token
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_delete_accept
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_delete_content_type
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationDelete_token"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_delete_token
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_get_accept
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_get_page_number
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_get_page_size
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationGet_token"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_get_token
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_post_accept
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_post_body
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_post_content_type
  }
  property {
    name  = "apiV1DashboardPublicTokenInvitationPost_token"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_invitation_post_token
  }
  property {
    name  = "apiV1DashboardPublicTokenPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_put_accept
  }
  property {
    name  = "apiV1DashboardPublicTokenPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_put_body
  }
  property {
    name  = "apiV1DashboardPublicTokenPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_put_content_type
  }
  property {
    name  = "apiV1DashboardPublicTokenPut_token"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_dashboard_public_token_put_token
  }
  property {
    name  = "apiV1DistributionPointsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_distribution_points_post_accept
  }
  property {
    name  = "apiV1DistributionPointsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_distribution_points_post_body
  }
  property {
    name  = "apiV1DistributionPointsPost_content_encoding"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_distribution_points_post_content_encoding
  }
  property {
    name  = "apiV1DistributionPointsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_distribution_points_post_content_type
  }
  property {
    name  = "apiV1DowntimeCancelByScopePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_cancel_by_scope_post_accept
  }
  property {
    name  = "apiV1DowntimeCancelByScopePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_cancel_by_scope_post_body
  }
  property {
    name  = "apiV1DowntimeCancelByScopePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_cancel_by_scope_post_content_type
  }
  property {
    name  = "apiV1DowntimeDowntimeIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_delete_accept
  }
  property {
    name  = "apiV1DowntimeDowntimeIdDelete_downtime_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_delete_downtime_id
  }
  property {
    name  = "apiV1DowntimeDowntimeIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_get_accept
  }
  property {
    name  = "apiV1DowntimeDowntimeIdGet_downtime_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_get_downtime_id
  }
  property {
    name  = "apiV1DowntimeDowntimeIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_put_accept
  }
  property {
    name  = "apiV1DowntimeDowntimeIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_put_body
  }
  property {
    name  = "apiV1DowntimeDowntimeIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_put_content_type
  }
  property {
    name  = "apiV1DowntimeDowntimeIdPut_downtime_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_downtime_id_put_downtime_id
  }
  property {
    name  = "apiV1DowntimeGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_get_accept
  }
  property {
    name  = "apiV1DowntimeGet_current_only"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_get_current_only
  }
  property {
    name  = "apiV1DowntimePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_post_accept
  }
  property {
    name  = "apiV1DowntimePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_post_body
  }
  property {
    name  = "apiV1DowntimePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_downtime_post_content_type
  }
  property {
    name  = "apiV1EventsEventIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_event_id_get_accept
  }
  property {
    name  = "apiV1EventsEventIdGet_event_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_event_id_get_event_id
  }
  property {
    name  = "apiV1EventsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_accept
  }
  property {
    name  = "apiV1EventsGet_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_end
  }
  property {
    name  = "apiV1EventsGet_exclude_aggregate"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_exclude_aggregate
  }
  property {
    name  = "apiV1EventsGet_page"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_page
  }
  property {
    name  = "apiV1EventsGet_priority"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_priority
  }
  property {
    name  = "apiV1EventsGet_sources"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_sources
  }
  property {
    name  = "apiV1EventsGet_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_start
  }
  property {
    name  = "apiV1EventsGet_tags"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_tags
  }
  property {
    name  = "apiV1EventsGet_unaggregated"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_get_unaggregated
  }
  property {
    name  = "apiV1EventsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_post_accept
  }
  property {
    name  = "apiV1EventsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_post_body
  }
  property {
    name  = "apiV1EventsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_events_post_content_type
  }
  property {
    name  = "apiV1GraphEmbedEmbedIdEnableGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_embed_id_enable_get_accept
  }
  property {
    name  = "apiV1GraphEmbedEmbedIdEnableGet_embed_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_embed_id_enable_get_embed_id
  }
  property {
    name  = "apiV1GraphEmbedEmbedIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_embed_id_get_accept
  }
  property {
    name  = "apiV1GraphEmbedEmbedIdGet_embed_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_embed_id_get_embed_id
  }
  property {
    name  = "apiV1GraphEmbedEmbedIdRevokeGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_embed_id_revoke_get_accept
  }
  property {
    name  = "apiV1GraphEmbedEmbedIdRevokeGet_embed_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_embed_id_revoke_get_embed_id
  }
  property {
    name  = "apiV1GraphEmbedGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_get_accept
  }
  property {
    name  = "apiV1GraphEmbedPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_post_accept
  }
  property {
    name  = "apiV1GraphEmbedPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_post_body
  }
  property {
    name  = "apiV1GraphEmbedPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_embed_post_content_type
  }
  property {
    name  = "apiV1GraphSnapshotGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_accept
  }
  property {
    name  = "apiV1GraphSnapshotGet_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_end
  }
  property {
    name  = "apiV1GraphSnapshotGet_event_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_event_query
  }
  property {
    name  = "apiV1GraphSnapshotGet_graph_def"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_graph_def
  }
  property {
    name  = "apiV1GraphSnapshotGet_height"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_height
  }
  property {
    name  = "apiV1GraphSnapshotGet_metric_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_metric_query
  }
  property {
    name  = "apiV1GraphSnapshotGet_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_start
  }
  property {
    name  = "apiV1GraphSnapshotGet_title"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_title
  }
  property {
    name  = "apiV1GraphSnapshotGet_width"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_graph_snapshot_get_width
  }
  property {
    name  = "apiV1HostHostNameMutePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_host_host_name_mute_post_accept
  }
  property {
    name  = "apiV1HostHostNameMutePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_host_host_name_mute_post_body
  }
  property {
    name  = "apiV1HostHostNameMutePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_host_host_name_mute_post_content_type
  }
  property {
    name  = "apiV1HostHostNameMutePost_host_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_host_host_name_mute_post_host_name
  }
  property {
    name  = "apiV1HostHostNameUnmutePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_host_host_name_unmute_post_accept
  }
  property {
    name  = "apiV1HostHostNameUnmutePost_host_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_host_host_name_unmute_post_host_name
  }
  property {
    name  = "apiV1HostsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_accept
  }
  property {
    name  = "apiV1HostsGet_count"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_count
  }
  property {
    name  = "apiV1HostsGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_filter
  }
  property {
    name  = "apiV1HostsGet_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_from
  }
  property {
    name  = "apiV1HostsGet_include_hosts_metadata"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_include_hosts_metadata
  }
  property {
    name  = "apiV1HostsGet_include_muted_hosts_data"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_include_muted_hosts_data
  }
  property {
    name  = "apiV1HostsGet_sort_dir"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_sort_dir
  }
  property {
    name  = "apiV1HostsGet_sort_field"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_sort_field
  }
  property {
    name  = "apiV1HostsGet_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_get_start
  }
  property {
    name  = "apiV1HostsTotalsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_totals_get_accept
  }
  property {
    name  = "apiV1HostsTotalsGet_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_hosts_totals_get_from
  }
  property {
    name  = "apiV1IntegrationAwsAvailableNamespaceRulesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_available_namespace_rules_get_accept
  }
  property {
    name  = "apiV1IntegrationAwsDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_delete_accept
  }
  property {
    name  = "apiV1IntegrationAwsDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_delete_content_type
  }
  property {
    name  = "apiV1IntegrationAwsFilteringDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_filtering_delete_accept
  }
  property {
    name  = "apiV1IntegrationAwsFilteringDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_filtering_delete_content_type
  }
  property {
    name  = "apiV1IntegrationAwsFilteringGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_filtering_get_accept
  }
  property {
    name  = "apiV1IntegrationAwsFilteringGet_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_filtering_get_account_id
  }
  property {
    name  = "apiV1IntegrationAwsFilteringPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_filtering_post_accept
  }
  property {
    name  = "apiV1IntegrationAwsFilteringPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_filtering_post_body
  }
  property {
    name  = "apiV1IntegrationAwsFilteringPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_filtering_post_content_type
  }
  property {
    name  = "apiV1IntegrationAwsGenerateNewExternalIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_generate_new_external_id_put_accept
  }
  property {
    name  = "apiV1IntegrationAwsGenerateNewExternalIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_generate_new_external_id_put_body
  }
  property {
    name  = "apiV1IntegrationAwsGenerateNewExternalIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_generate_new_external_id_put_content_type
  }
  property {
    name  = "apiV1IntegrationAwsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_get_accept
  }
  property {
    name  = "apiV1IntegrationAwsGet_access_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_get_access_key_id
  }
  property {
    name  = "apiV1IntegrationAwsGet_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_get_account_id
  }
  property {
    name  = "apiV1IntegrationAwsGet_role_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_get_role_name
  }
  property {
    name  = "apiV1IntegrationAwsLogsCheckAsyncPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_check_async_post_accept
  }
  property {
    name  = "apiV1IntegrationAwsLogsCheckAsyncPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_check_async_post_body
  }
  property {
    name  = "apiV1IntegrationAwsLogsCheckAsyncPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_check_async_post_content_type
  }
  property {
    name  = "apiV1IntegrationAwsLogsDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_delete_accept
  }
  property {
    name  = "apiV1IntegrationAwsLogsDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_delete_content_type
  }
  property {
    name  = "apiV1IntegrationAwsLogsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_get_accept
  }
  property {
    name  = "apiV1IntegrationAwsLogsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_post_accept
  }
  property {
    name  = "apiV1IntegrationAwsLogsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_post_body
  }
  property {
    name  = "apiV1IntegrationAwsLogsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_post_content_type
  }
  property {
    name  = "apiV1IntegrationAwsLogsServicesAsyncPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_services_async_post_accept
  }
  property {
    name  = "apiV1IntegrationAwsLogsServicesAsyncPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_services_async_post_body
  }
  property {
    name  = "apiV1IntegrationAwsLogsServicesAsyncPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_services_async_post_content_type
  }
  property {
    name  = "apiV1IntegrationAwsLogsServicesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_services_get_accept
  }
  property {
    name  = "apiV1IntegrationAwsLogsServicesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_services_post_accept
  }
  property {
    name  = "apiV1IntegrationAwsLogsServicesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_services_post_body
  }
  property {
    name  = "apiV1IntegrationAwsLogsServicesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_logs_services_post_content_type
  }
  property {
    name  = "apiV1IntegrationAwsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_post_accept
  }
  property {
    name  = "apiV1IntegrationAwsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_post_body
  }
  property {
    name  = "apiV1IntegrationAwsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_post_content_type
  }
  property {
    name  = "apiV1IntegrationAwsPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_put_accept
  }
  property {
    name  = "apiV1IntegrationAwsPut_access_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_put_access_key_id
  }
  property {
    name  = "apiV1IntegrationAwsPut_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_put_account_id
  }
  property {
    name  = "apiV1IntegrationAwsPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_put_body
  }
  property {
    name  = "apiV1IntegrationAwsPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_put_content_type
  }
  property {
    name  = "apiV1IntegrationAwsPut_role_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_aws_put_role_name
  }
  property {
    name  = "apiV1IntegrationAzureDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_delete_accept
  }
  property {
    name  = "apiV1IntegrationAzureDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_delete_content_type
  }
  property {
    name  = "apiV1IntegrationAzureGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_get_accept
  }
  property {
    name  = "apiV1IntegrationAzureHostFiltersPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_host_filters_post_accept
  }
  property {
    name  = "apiV1IntegrationAzureHostFiltersPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_host_filters_post_body
  }
  property {
    name  = "apiV1IntegrationAzureHostFiltersPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_host_filters_post_content_type
  }
  property {
    name  = "apiV1IntegrationAzurePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_post_accept
  }
  property {
    name  = "apiV1IntegrationAzurePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_post_body
  }
  property {
    name  = "apiV1IntegrationAzurePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_post_content_type
  }
  property {
    name  = "apiV1IntegrationAzurePut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_put_accept
  }
  property {
    name  = "apiV1IntegrationAzurePut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_put_body
  }
  property {
    name  = "apiV1IntegrationAzurePut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_azure_put_content_type
  }
  property {
    name  = "apiV1IntegrationGcpDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_delete_accept
  }
  property {
    name  = "apiV1IntegrationGcpDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_delete_content_type
  }
  property {
    name  = "apiV1IntegrationGcpGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_get_accept
  }
  property {
    name  = "apiV1IntegrationGcpPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_post_accept
  }
  property {
    name  = "apiV1IntegrationGcpPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_post_body
  }
  property {
    name  = "apiV1IntegrationGcpPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_post_content_type
  }
  property {
    name  = "apiV1IntegrationGcpPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_put_accept
  }
  property {
    name  = "apiV1IntegrationGcpPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_put_body
  }
  property {
    name  = "apiV1IntegrationGcpPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_gcp_put_content_type
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_post_accept
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_post_body
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_post_content_type
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNameDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_delete_accept
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNameDelete_service_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_delete_service_name
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNameGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_get_accept
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNameGet_service_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_get_service_name
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNamePut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_put_accept
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNamePut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_put_body
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNamePut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_put_content_type
  }
  property {
    name  = "apiV1IntegrationPagerdutyConfigurationServicesServiceNamePut_service_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_pagerduty_configuration_services_service_name_put_service_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNameDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_delete_accept
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNameDelete_account_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_delete_account_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNameDelete_channel_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_delete_channel_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNameGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_get_accept
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNameGet_account_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_get_account_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNameGet_channel_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_get_channel_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNamePatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_patch_accept
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNamePatch_account_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_patch_account_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNamePatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_patch_body
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNamePatch_channel_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_patch_channel_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsChannelNamePatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_channel_name_patch_content_type
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_get_accept
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsGet_account_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_get_account_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_post_accept
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsPost_account_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_post_account_name
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_post_body
  }
  property {
    name  = "apiV1IntegrationSlackConfigurationAccountsAccountNameChannelsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_configuration_accounts_account_name_channels_post_content_type
  }
  property {
    name  = "apiV1IntegrationSlackDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_delete_accept
  }
  property {
    name  = "apiV1IntegrationSlackGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_get_accept
  }
  property {
    name  = "apiV1IntegrationSlackPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_post_accept
  }
  property {
    name  = "apiV1IntegrationSlackPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_post_body
  }
  property {
    name  = "apiV1IntegrationSlackPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_post_content_type
  }
  property {
    name  = "apiV1IntegrationSlackPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_put_accept
  }
  property {
    name  = "apiV1IntegrationSlackPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_put_body
  }
  property {
    name  = "apiV1IntegrationSlackPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_slack_put_content_type
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNameDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_delete_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNameDelete_custom_variable_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_delete_custom_variable_name
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNameGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_get_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNameGet_custom_variable_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_get_custom_variable_name
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNamePut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_put_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNamePut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_put_body
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNamePut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_put_content_type
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesCustomVariableNamePut_custom_variable_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_custom_variable_name_put_custom_variable_name
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_post_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_post_body
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationCustomVariablesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_custom_variables_post_content_type
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_post_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_post_body
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_post_content_type
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNameDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_delete_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNameDelete_webhook_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_delete_webhook_name
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNameGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_get_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNameGet_webhook_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_get_webhook_name
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNamePut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_put_accept
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNamePut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_put_body
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNamePut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_put_content_type
  }
  property {
    name  = "apiV1IntegrationWebhooksConfigurationWebhooksWebhookNamePut_webhook_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_integration_webhooks_configuration_webhooks_webhook_name_put_webhook_name
  }
  property {
    name  = "apiV1LogsConfigIndexOrderGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_index_order_get_accept
  }
  property {
    name  = "apiV1LogsConfigIndexOrderPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_index_order_put_accept
  }
  property {
    name  = "apiV1LogsConfigIndexOrderPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_index_order_put_body
  }
  property {
    name  = "apiV1LogsConfigIndexOrderPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_index_order_put_content_type
  }
  property {
    name  = "apiV1LogsConfigIndexesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_get_accept
  }
  property {
    name  = "apiV1LogsConfigIndexesNameGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_name_get_accept
  }
  property {
    name  = "apiV1LogsConfigIndexesNameGet_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_name_get_name
  }
  property {
    name  = "apiV1LogsConfigIndexesNamePut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_name_put_accept
  }
  property {
    name  = "apiV1LogsConfigIndexesNamePut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_name_put_body
  }
  property {
    name  = "apiV1LogsConfigIndexesNamePut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_name_put_content_type
  }
  property {
    name  = "apiV1LogsConfigIndexesNamePut_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_name_put_name
  }
  property {
    name  = "apiV1LogsConfigIndexesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_post_accept
  }
  property {
    name  = "apiV1LogsConfigIndexesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_post_body
  }
  property {
    name  = "apiV1LogsConfigIndexesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_indexes_post_content_type
  }
  property {
    name  = "apiV1LogsConfigPipelineOrderGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipeline_order_get_accept
  }
  property {
    name  = "apiV1LogsConfigPipelineOrderPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipeline_order_put_accept
  }
  property {
    name  = "apiV1LogsConfigPipelineOrderPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipeline_order_put_body
  }
  property {
    name  = "apiV1LogsConfigPipelineOrderPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipeline_order_put_content_type
  }
  property {
    name  = "apiV1LogsConfigPipelinesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_get_accept
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_delete_accept
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdDelete_pipeline_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_delete_pipeline_id
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_get_accept
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdGet_pipeline_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_get_pipeline_id
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_put_accept
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_put_body
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_put_content_type
  }
  property {
    name  = "apiV1LogsConfigPipelinesPipelineIdPut_pipeline_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_pipeline_id_put_pipeline_id
  }
  property {
    name  = "apiV1LogsConfigPipelinesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_post_accept
  }
  property {
    name  = "apiV1LogsConfigPipelinesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_post_body
  }
  property {
    name  = "apiV1LogsConfigPipelinesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_logs_config_pipelines_post_content_type
  }
  property {
    name  = "apiV1MetricsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_get_accept
  }
  property {
    name  = "apiV1MetricsGet_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_get_from
  }
  property {
    name  = "apiV1MetricsGet_host"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_get_host
  }
  property {
    name  = "apiV1MetricsGet_tag_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_get_tag_filter
  }
  property {
    name  = "apiV1MetricsMetricNameGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_metric_name_get_accept
  }
  property {
    name  = "apiV1MetricsMetricNameGet_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_metric_name_get_metric_name
  }
  property {
    name  = "apiV1MetricsMetricNamePut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_metric_name_put_accept
  }
  property {
    name  = "apiV1MetricsMetricNamePut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_metric_name_put_body
  }
  property {
    name  = "apiV1MetricsMetricNamePut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_metric_name_put_content_type
  }
  property {
    name  = "apiV1MetricsMetricNamePut_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_metrics_metric_name_put_metric_name
  }
  property {
    name  = "apiV1MonitorCanDeleteGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_can_delete_get_accept
  }
  property {
    name  = "apiV1MonitorCanDeleteGet_monitor_ids"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_can_delete_get_monitor_ids
  }
  property {
    name  = "apiV1MonitorGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_accept
  }
  property {
    name  = "apiV1MonitorGet_group_states"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_group_states
  }
  property {
    name  = "apiV1MonitorGet_id_offset"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_id_offset
  }
  property {
    name  = "apiV1MonitorGet_monitor_tags"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_monitor_tags
  }
  property {
    name  = "apiV1MonitorGet_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_name
  }
  property {
    name  = "apiV1MonitorGet_page"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_page
  }
  property {
    name  = "apiV1MonitorGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_page_size
  }
  property {
    name  = "apiV1MonitorGet_tags"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_tags
  }
  property {
    name  = "apiV1MonitorGet_with_downtimes"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_get_with_downtimes
  }
  property {
    name  = "apiV1MonitorGroupsSearchGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_groups_search_get_accept
  }
  property {
    name  = "apiV1MonitorGroupsSearchGet_page"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_groups_search_get_page
  }
  property {
    name  = "apiV1MonitorGroupsSearchGet_per_page"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_groups_search_get_per_page
  }
  property {
    name  = "apiV1MonitorGroupsSearchGet_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_groups_search_get_query
  }
  property {
    name  = "apiV1MonitorGroupsSearchGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_groups_search_get_sort
  }
  property {
    name  = "apiV1MonitorMonitorIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_delete_accept
  }
  property {
    name  = "apiV1MonitorMonitorIdDelete_force"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_delete_force
  }
  property {
    name  = "apiV1MonitorMonitorIdDelete_monitor_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_delete_monitor_id
  }
  property {
    name  = "apiV1MonitorMonitorIdDowntimesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_downtimes_get_accept
  }
  property {
    name  = "apiV1MonitorMonitorIdDowntimesGet_monitor_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_downtimes_get_monitor_id
  }
  property {
    name  = "apiV1MonitorMonitorIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_get_accept
  }
  property {
    name  = "apiV1MonitorMonitorIdGet_group_states"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_get_group_states
  }
  property {
    name  = "apiV1MonitorMonitorIdGet_monitor_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_get_monitor_id
  }
  property {
    name  = "apiV1MonitorMonitorIdMutePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_mute_post_accept
  }
  property {
    name  = "apiV1MonitorMonitorIdMutePost_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_mute_post_end
  }
  property {
    name  = "apiV1MonitorMonitorIdMutePost_monitor_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_mute_post_monitor_id
  }
  property {
    name  = "apiV1MonitorMonitorIdMutePost_scope"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_mute_post_scope
  }
  property {
    name  = "apiV1MonitorMonitorIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_put_accept
  }
  property {
    name  = "apiV1MonitorMonitorIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_put_body
  }
  property {
    name  = "apiV1MonitorMonitorIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_put_content_type
  }
  property {
    name  = "apiV1MonitorMonitorIdPut_monitor_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_put_monitor_id
  }
  property {
    name  = "apiV1MonitorMonitorIdUnmutePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_unmute_post_accept
  }
  property {
    name  = "apiV1MonitorMonitorIdUnmutePost_all_scopes"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_unmute_post_all_scopes
  }
  property {
    name  = "apiV1MonitorMonitorIdUnmutePost_monitor_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_unmute_post_monitor_id
  }
  property {
    name  = "apiV1MonitorMonitorIdUnmutePost_scope"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_unmute_post_scope
  }
  property {
    name  = "apiV1MonitorMonitorIdValidatePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_validate_post_accept
  }
  property {
    name  = "apiV1MonitorMonitorIdValidatePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_validate_post_body
  }
  property {
    name  = "apiV1MonitorMonitorIdValidatePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_validate_post_content_type
  }
  property {
    name  = "apiV1MonitorMonitorIdValidatePost_monitor_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_monitor_id_validate_post_monitor_id
  }
  property {
    name  = "apiV1MonitorPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_post_accept
  }
  property {
    name  = "apiV1MonitorPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_post_body
  }
  property {
    name  = "apiV1MonitorPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_post_content_type
  }
  property {
    name  = "apiV1MonitorSearchGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_search_get_accept
  }
  property {
    name  = "apiV1MonitorSearchGet_page"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_search_get_page
  }
  property {
    name  = "apiV1MonitorSearchGet_per_page"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_search_get_per_page
  }
  property {
    name  = "apiV1MonitorSearchGet_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_search_get_query
  }
  property {
    name  = "apiV1MonitorSearchGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_search_get_sort
  }
  property {
    name  = "apiV1MonitorValidatePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_validate_post_accept
  }
  property {
    name  = "apiV1MonitorValidatePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_validate_post_body
  }
  property {
    name  = "apiV1MonitorValidatePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monitor_validate_post_content_type
  }
  property {
    name  = "apiV1MonthlyCustomReportsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monthly_custom_reports_get_accept
  }
  property {
    name  = "apiV1MonthlyCustomReportsGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monthly_custom_reports_get_page_number
  }
  property {
    name  = "apiV1MonthlyCustomReportsGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monthly_custom_reports_get_page_size
  }
  property {
    name  = "apiV1MonthlyCustomReportsGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monthly_custom_reports_get_sort
  }
  property {
    name  = "apiV1MonthlyCustomReportsGet_sort_dir"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monthly_custom_reports_get_sort_dir
  }
  property {
    name  = "apiV1MonthlyCustomReportsReportIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monthly_custom_reports_report_id_get_accept
  }
  property {
    name  = "apiV1MonthlyCustomReportsReportIdGet_report_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_monthly_custom_reports_report_id_get_report_id
  }
  property {
    name  = "apiV1NotebooksGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_accept
  }
  property {
    name  = "apiV1NotebooksGet_author_handle"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_author_handle
  }
  property {
    name  = "apiV1NotebooksGet_count"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_count
  }
  property {
    name  = "apiV1NotebooksGet_exclude_author_handle"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_exclude_author_handle
  }
  property {
    name  = "apiV1NotebooksGet_include_cells"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_include_cells
  }
  property {
    name  = "apiV1NotebooksGet_is_template"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_is_template
  }
  property {
    name  = "apiV1NotebooksGet_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_query
  }
  property {
    name  = "apiV1NotebooksGet_sort_dir"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_sort_dir
  }
  property {
    name  = "apiV1NotebooksGet_sort_field"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_sort_field
  }
  property {
    name  = "apiV1NotebooksGet_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_start
  }
  property {
    name  = "apiV1NotebooksGet_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_get_type
  }
  property {
    name  = "apiV1NotebooksNotebookIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_delete_accept
  }
  property {
    name  = "apiV1NotebooksNotebookIdDelete_notebook_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_delete_notebook_id
  }
  property {
    name  = "apiV1NotebooksNotebookIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_get_accept
  }
  property {
    name  = "apiV1NotebooksNotebookIdGet_notebook_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_get_notebook_id
  }
  property {
    name  = "apiV1NotebooksNotebookIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_put_accept
  }
  property {
    name  = "apiV1NotebooksNotebookIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_put_body
  }
  property {
    name  = "apiV1NotebooksNotebookIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_put_content_type
  }
  property {
    name  = "apiV1NotebooksNotebookIdPut_notebook_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_notebook_id_put_notebook_id
  }
  property {
    name  = "apiV1NotebooksPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_post_accept
  }
  property {
    name  = "apiV1NotebooksPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_post_body
  }
  property {
    name  = "apiV1NotebooksPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_notebooks_post_content_type
  }
  property {
    name  = "apiV1OrgGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_get_accept
  }
  property {
    name  = "apiV1OrgPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_post_accept
  }
  property {
    name  = "apiV1OrgPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_post_body
  }
  property {
    name  = "apiV1OrgPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_post_content_type
  }
  property {
    name  = "apiV1OrgPublicIdDowngradePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_downgrade_post_accept
  }
  property {
    name  = "apiV1OrgPublicIdDowngradePost_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_downgrade_post_public_id
  }
  property {
    name  = "apiV1OrgPublicIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_get_accept
  }
  property {
    name  = "apiV1OrgPublicIdGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_get_public_id
  }
  property {
    name  = "apiV1OrgPublicIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_put_accept
  }
  property {
    name  = "apiV1OrgPublicIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_put_body
  }
  property {
    name  = "apiV1OrgPublicIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_put_content_type
  }
  property {
    name  = "apiV1OrgPublicIdPut_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_org_public_id_put_public_id
  }
  property {
    name  = "apiV1QueryGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_query_get_accept
  }
  property {
    name  = "apiV1QueryGet_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_query_get_from
  }
  property {
    name  = "apiV1QueryGet_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_query_get_query
  }
  property {
    name  = "apiV1QueryGet_to"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_query_get_to
  }
  property {
    name  = "apiV1SearchGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_search_get_accept
  }
  property {
    name  = "apiV1SearchGet_q"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_search_get_q
  }
  property {
    name  = "apiV1SecurityAnalyticsSignalsSignalIdAddToIncidentPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_security_analytics_signals_signal_id_add_to_incident_patch_accept
  }
  property {
    name  = "apiV1SecurityAnalyticsSignalsSignalIdAddToIncidentPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_security_analytics_signals_signal_id_add_to_incident_patch_body
  }
  property {
    name  = "apiV1SecurityAnalyticsSignalsSignalIdAddToIncidentPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_security_analytics_signals_signal_id_add_to_incident_patch_content_type
  }
  property {
    name  = "apiV1SecurityAnalyticsSignalsSignalIdAddToIncidentPatch_signal_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_security_analytics_signals_signal_id_add_to_incident_patch_signal_id
  }
  property {
    name  = "apiV1ServiceDependenciesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_get_accept
  }
  property {
    name  = "apiV1ServiceDependenciesGet_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_get_end
  }
  property {
    name  = "apiV1ServiceDependenciesGet_env"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_get_env
  }
  property {
    name  = "apiV1ServiceDependenciesGet_primary_tag"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_get_primary_tag
  }
  property {
    name  = "apiV1ServiceDependenciesGet_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_get_start
  }
  property {
    name  = "apiV1ServiceDependenciesServiceGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_service_get_accept
  }
  property {
    name  = "apiV1ServiceDependenciesServiceGet_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_service_get_end
  }
  property {
    name  = "apiV1ServiceDependenciesServiceGet_env"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_service_get_env
  }
  property {
    name  = "apiV1ServiceDependenciesServiceGet_primary_tag"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_service_get_primary_tag
  }
  property {
    name  = "apiV1ServiceDependenciesServiceGet_service"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_service_get_service
  }
  property {
    name  = "apiV1ServiceDependenciesServiceGet_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_service_dependencies_service_get_start
  }
  property {
    name  = "apiV1SloBulkDeletePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_bulk_delete_post_accept
  }
  property {
    name  = "apiV1SloBulkDeletePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_bulk_delete_post_body
  }
  property {
    name  = "apiV1SloBulkDeletePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_bulk_delete_post_content_type
  }
  property {
    name  = "apiV1SloCanDeleteGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_can_delete_get_accept
  }
  property {
    name  = "apiV1SloCanDeleteGet_ids"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_can_delete_get_ids
  }
  property {
    name  = "apiV1SloCorrectionGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_get_accept
  }
  property {
    name  = "apiV1SloCorrectionPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_post_accept
  }
  property {
    name  = "apiV1SloCorrectionPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_post_body
  }
  property {
    name  = "apiV1SloCorrectionPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_post_content_type
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_delete_accept
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdDelete_slo_correction_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_delete_slo_correction_id
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_get_accept
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdGet_slo_correction_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_get_slo_correction_id
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_patch_accept
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_patch_body
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_patch_content_type
  }
  property {
    name  = "apiV1SloCorrectionSloCorrectionIdPatch_slo_correction_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_correction_slo_correction_id_patch_slo_correction_id
  }
  property {
    name  = "apiV1SloGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_get_accept
  }
  property {
    name  = "apiV1SloGet_ids"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_get_ids
  }
  property {
    name  = "apiV1SloGet_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_get_limit
  }
  property {
    name  = "apiV1SloGet_metrics_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_get_metrics_query
  }
  property {
    name  = "apiV1SloGet_offset"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_get_offset
  }
  property {
    name  = "apiV1SloGet_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_get_query
  }
  property {
    name  = "apiV1SloGet_tags_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_get_tags_query
  }
  property {
    name  = "apiV1SloPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_post_accept
  }
  property {
    name  = "apiV1SloPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_post_body
  }
  property {
    name  = "apiV1SloPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_post_content_type
  }
  property {
    name  = "apiV1SloSearchGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_search_get_accept
  }
  property {
    name  = "apiV1SloSearchGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_search_get_page_number
  }
  property {
    name  = "apiV1SloSearchGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_search_get_page_size
  }
  property {
    name  = "apiV1SloSearchGet_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_search_get_query
  }
  property {
    name  = "apiV1SloSloIdCorrectionsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_corrections_get_accept
  }
  property {
    name  = "apiV1SloSloIdCorrectionsGet_slo_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_corrections_get_slo_id
  }
  property {
    name  = "apiV1SloSloIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_delete_accept
  }
  property {
    name  = "apiV1SloSloIdDelete_force"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_delete_force
  }
  property {
    name  = "apiV1SloSloIdDelete_slo_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_delete_slo_id
  }
  property {
    name  = "apiV1SloSloIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_get_accept
  }
  property {
    name  = "apiV1SloSloIdGet_slo_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_get_slo_id
  }
  property {
    name  = "apiV1SloSloIdGet_with_configured_alert_ids"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_get_with_configured_alert_ids
  }
  property {
    name  = "apiV1SloSloIdHistoryGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_history_get_accept
  }
  property {
    name  = "apiV1SloSloIdHistoryGet_apply_correction"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_history_get_apply_correction
  }
  property {
    name  = "apiV1SloSloIdHistoryGet_from_ts"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_history_get_from_ts
  }
  property {
    name  = "apiV1SloSloIdHistoryGet_slo_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_history_get_slo_id
  }
  property {
    name  = "apiV1SloSloIdHistoryGet_target"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_history_get_target
  }
  property {
    name  = "apiV1SloSloIdHistoryGet_to_ts"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_history_get_to_ts
  }
  property {
    name  = "apiV1SloSloIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_put_accept
  }
  property {
    name  = "apiV1SloSloIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_put_body
  }
  property {
    name  = "apiV1SloSloIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_put_content_type
  }
  property {
    name  = "apiV1SloSloIdPut_slo_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_slo_slo_id_put_slo_id
  }
  property {
    name  = "apiV1SyntheticsCiBatchBatchIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_ci_batch_batch_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsCiBatchBatchIdGet_batch_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_ci_batch_batch_id_get_batch_id
  }
  property {
    name  = "apiV1SyntheticsLocationsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_locations_get_accept
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_delete_accept
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdDelete_location_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_delete_location_id
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdGet_location_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_get_location_id
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_put_accept
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_put_body
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_put_content_type
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsLocationIdPut_location_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_location_id_put_location_id
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_post_accept
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_post_body
  }
  property {
    name  = "apiV1SyntheticsPrivateLocationsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_private_locations_post_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsApiPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_post_accept
  }
  property {
    name  = "apiV1SyntheticsTestsApiPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_post_body
  }
  property {
    name  = "apiV1SyntheticsTestsApiPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_post_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsApiPublicIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_public_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsApiPublicIdGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_public_id_get_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsApiPublicIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_public_id_put_accept
  }
  property {
    name  = "apiV1SyntheticsTestsApiPublicIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_public_id_put_body
  }
  property {
    name  = "apiV1SyntheticsTestsApiPublicIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_public_id_put_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsApiPublicIdPut_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_api_public_id_put_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_post_accept
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_post_body
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_post_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_get_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_put_accept
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_put_body
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_put_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdPut_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_put_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsGet_from_ts"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_get_from_ts
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsGet_probe_dc"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_get_probe_dc
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_get_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsGet_to_ts"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_get_to_ts
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsResultIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_result_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsResultIdGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_result_id_get_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsBrowserPublicIdResultsResultIdGet_result_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_browser_public_id_results_result_id_get_result_id
  }
  property {
    name  = "apiV1SyntheticsTestsDeletePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_delete_post_accept
  }
  property {
    name  = "apiV1SyntheticsTestsDeletePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_delete_post_body
  }
  property {
    name  = "apiV1SyntheticsTestsDeletePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_delete_post_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_post_accept
  }
  property {
    name  = "apiV1SyntheticsTestsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_post_body
  }
  property {
    name  = "apiV1SyntheticsTestsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_post_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_get_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_put_accept
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_put_body
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_put_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdPut_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_put_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsGet_from_ts"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_get_from_ts
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsGet_probe_dc"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_get_probe_dc
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_get_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsGet_to_ts"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_get_to_ts
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsResultIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_result_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsResultIdGet_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_result_id_get_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdResultsResultIdGet_result_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_results_result_id_get_result_id
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdStatusPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_status_put_accept
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdStatusPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_status_put_body
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdStatusPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_status_put_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsPublicIdStatusPut_public_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_public_id_status_put_public_id
  }
  property {
    name  = "apiV1SyntheticsTestsTriggerCiPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_trigger_ci_post_accept
  }
  property {
    name  = "apiV1SyntheticsTestsTriggerCiPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_trigger_ci_post_body
  }
  property {
    name  = "apiV1SyntheticsTestsTriggerCiPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_trigger_ci_post_content_type
  }
  property {
    name  = "apiV1SyntheticsTestsTriggerPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_trigger_post_accept
  }
  property {
    name  = "apiV1SyntheticsTestsTriggerPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_trigger_post_body
  }
  property {
    name  = "apiV1SyntheticsTestsTriggerPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_tests_trigger_post_content_type
  }
  property {
    name  = "apiV1SyntheticsVariablesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_get_accept
  }
  property {
    name  = "apiV1SyntheticsVariablesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_post_accept
  }
  property {
    name  = "apiV1SyntheticsVariablesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_post_body
  }
  property {
    name  = "apiV1SyntheticsVariablesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_post_content_type
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_delete_accept
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdDelete_variable_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_delete_variable_id
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_get_accept
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdGet_variable_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_get_variable_id
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_put_accept
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_put_body
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_put_content_type
  }
  property {
    name  = "apiV1SyntheticsVariablesVariableIdPut_variable_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_synthetics_variables_variable_id_put_variable_id
  }
  property {
    name  = "apiV1TagsHostsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_get_accept
  }
  property {
    name  = "apiV1TagsHostsGet_source"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_get_source
  }
  property {
    name  = "apiV1TagsHostsHostNameDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_delete_accept
  }
  property {
    name  = "apiV1TagsHostsHostNameDelete_host_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_delete_host_name
  }
  property {
    name  = "apiV1TagsHostsHostNameDelete_source"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_delete_source
  }
  property {
    name  = "apiV1TagsHostsHostNameGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_get_accept
  }
  property {
    name  = "apiV1TagsHostsHostNameGet_host_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_get_host_name
  }
  property {
    name  = "apiV1TagsHostsHostNameGet_source"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_get_source
  }
  property {
    name  = "apiV1TagsHostsHostNamePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_post_accept
  }
  property {
    name  = "apiV1TagsHostsHostNamePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_post_body
  }
  property {
    name  = "apiV1TagsHostsHostNamePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_post_content_type
  }
  property {
    name  = "apiV1TagsHostsHostNamePost_host_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_post_host_name
  }
  property {
    name  = "apiV1TagsHostsHostNamePost_source"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_post_source
  }
  property {
    name  = "apiV1TagsHostsHostNamePut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_put_accept
  }
  property {
    name  = "apiV1TagsHostsHostNamePut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_put_body
  }
  property {
    name  = "apiV1TagsHostsHostNamePut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_put_content_type
  }
  property {
    name  = "apiV1TagsHostsHostNamePut_host_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_put_host_name
  }
  property {
    name  = "apiV1TagsHostsHostNamePut_source"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_tags_hosts_host_name_put_source
  }
  property {
    name  = "apiV1UsageAnalyzedLogsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_analyzed_logs_get_accept
  }
  property {
    name  = "apiV1UsageAnalyzedLogsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_analyzed_logs_get_end_hr
  }
  property {
    name  = "apiV1UsageAnalyzedLogsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_analyzed_logs_get_start_hr
  }
  property {
    name  = "apiV1UsageAttributionGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_accept
  }
  property {
    name  = "apiV1UsageAttributionGet_end_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_end_month
  }
  property {
    name  = "apiV1UsageAttributionGet_fields"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_fields
  }
  property {
    name  = "apiV1UsageAttributionGet_include_descendants"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_include_descendants
  }
  property {
    name  = "apiV1UsageAttributionGet_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_limit
  }
  property {
    name  = "apiV1UsageAttributionGet_offset"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_offset
  }
  property {
    name  = "apiV1UsageAttributionGet_sort_direction"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_sort_direction
  }
  property {
    name  = "apiV1UsageAttributionGet_sort_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_sort_name
  }
  property {
    name  = "apiV1UsageAttributionGet_start_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_attribution_get_start_month
  }
  property {
    name  = "apiV1UsageAuditLogsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_audit_logs_get_accept
  }
  property {
    name  = "apiV1UsageAuditLogsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_audit_logs_get_end_hr
  }
  property {
    name  = "apiV1UsageAuditLogsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_audit_logs_get_start_hr
  }
  property {
    name  = "apiV1UsageAwsLambdaGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_aws_lambda_get_accept
  }
  property {
    name  = "apiV1UsageAwsLambdaGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_aws_lambda_get_end_hr
  }
  property {
    name  = "apiV1UsageAwsLambdaGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_aws_lambda_get_start_hr
  }
  property {
    name  = "apiV1UsageBillableSummaryGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_billable_summary_get_accept
  }
  property {
    name  = "apiV1UsageBillableSummaryGet_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_billable_summary_get_month
  }
  property {
    name  = "apiV1UsageCiAppGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_ci_app_get_accept
  }
  property {
    name  = "apiV1UsageCiAppGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_ci_app_get_end_hr
  }
  property {
    name  = "apiV1UsageCiAppGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_ci_app_get_start_hr
  }
  property {
    name  = "apiV1UsageCspmGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_cspm_get_accept
  }
  property {
    name  = "apiV1UsageCspmGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_cspm_get_end_hr
  }
  property {
    name  = "apiV1UsageCspmGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_cspm_get_start_hr
  }
  property {
    name  = "apiV1UsageCwsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_cws_get_accept
  }
  property {
    name  = "apiV1UsageCwsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_cws_get_end_hr
  }
  property {
    name  = "apiV1UsageCwsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_cws_get_start_hr
  }
  property {
    name  = "apiV1UsageDbmGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_dbm_get_accept
  }
  property {
    name  = "apiV1UsageDbmGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_dbm_get_end_hr
  }
  property {
    name  = "apiV1UsageDbmGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_dbm_get_start_hr
  }
  property {
    name  = "apiV1UsageFargateGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_fargate_get_accept
  }
  property {
    name  = "apiV1UsageFargateGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_fargate_get_end_hr
  }
  property {
    name  = "apiV1UsageFargateGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_fargate_get_start_hr
  }
  property {
    name  = "apiV1UsageHostsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hosts_get_accept
  }
  property {
    name  = "apiV1UsageHostsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hosts_get_end_hr
  }
  property {
    name  = "apiV1UsageHostsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hosts_get_start_hr
  }
  property {
    name  = "apiV1UsageHourlyAttributionGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hourly_attribution_get_accept
  }
  property {
    name  = "apiV1UsageHourlyAttributionGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hourly_attribution_get_end_hr
  }
  property {
    name  = "apiV1UsageHourlyAttributionGet_include_descendants"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hourly_attribution_get_include_descendants
  }
  property {
    name  = "apiV1UsageHourlyAttributionGet_next_record_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hourly_attribution_get_next_record_id
  }
  property {
    name  = "apiV1UsageHourlyAttributionGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hourly_attribution_get_start_hr
  }
  property {
    name  = "apiV1UsageHourlyAttributionGet_tag_breakdown_keys"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hourly_attribution_get_tag_breakdown_keys
  }
  property {
    name  = "apiV1UsageHourlyAttributionGet_usage_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_hourly_attribution_get_usage_type
  }
  property {
    name  = "apiV1UsageIncidentManagementGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_incident_management_get_accept
  }
  property {
    name  = "apiV1UsageIncidentManagementGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_incident_management_get_end_hr
  }
  property {
    name  = "apiV1UsageIncidentManagementGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_incident_management_get_start_hr
  }
  property {
    name  = "apiV1UsageIndexedSpansGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_indexed_spans_get_accept
  }
  property {
    name  = "apiV1UsageIndexedSpansGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_indexed_spans_get_end_hr
  }
  property {
    name  = "apiV1UsageIndexedSpansGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_indexed_spans_get_start_hr
  }
  property {
    name  = "apiV1UsageIngestedSpansGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_ingested_spans_get_accept
  }
  property {
    name  = "apiV1UsageIngestedSpansGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_ingested_spans_get_end_hr
  }
  property {
    name  = "apiV1UsageIngestedSpansGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_ingested_spans_get_start_hr
  }
  property {
    name  = "apiV1UsageIotGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_iot_get_accept
  }
  property {
    name  = "apiV1UsageIotGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_iot_get_end_hr
  }
  property {
    name  = "apiV1UsageIotGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_iot_get_start_hr
  }
  property {
    name  = "apiV1UsageLogsByIndexGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_by_index_get_accept
  }
  property {
    name  = "apiV1UsageLogsByIndexGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_by_index_get_end_hr
  }
  property {
    name  = "apiV1UsageLogsByIndexGet_index_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_by_index_get_index_name
  }
  property {
    name  = "apiV1UsageLogsByIndexGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_by_index_get_start_hr
  }
  property {
    name  = "apiV1UsageLogsByRetentionGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_by_retention_get_accept
  }
  property {
    name  = "apiV1UsageLogsByRetentionGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_by_retention_get_end_hr
  }
  property {
    name  = "apiV1UsageLogsByRetentionGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_by_retention_get_start_hr
  }
  property {
    name  = "apiV1UsageLogsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_get_accept
  }
  property {
    name  = "apiV1UsageLogsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_get_end_hr
  }
  property {
    name  = "apiV1UsageLogsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_logs_get_start_hr
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_accept
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_end_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_end_month
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_fields"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_fields
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_include_descendants"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_include_descendants
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_next_record_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_next_record_id
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_sort_direction"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_sort_direction
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_sort_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_sort_name
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_start_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_start_month
  }
  property {
    name  = "apiV1UsageMonthlyAttributionGet_tag_breakdown_keys"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_monthly_attribution_get_tag_breakdown_keys
  }
  property {
    name  = "apiV1UsageNetworkFlowsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_network_flows_get_accept
  }
  property {
    name  = "apiV1UsageNetworkFlowsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_network_flows_get_end_hr
  }
  property {
    name  = "apiV1UsageNetworkFlowsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_network_flows_get_start_hr
  }
  property {
    name  = "apiV1UsageNetworkHostsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_network_hosts_get_accept
  }
  property {
    name  = "apiV1UsageNetworkHostsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_network_hosts_get_end_hr
  }
  property {
    name  = "apiV1UsageNetworkHostsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_network_hosts_get_start_hr
  }
  property {
    name  = "apiV1UsageOnlineArchiveGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_online_archive_get_accept
  }
  property {
    name  = "apiV1UsageOnlineArchiveGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_online_archive_get_end_hr
  }
  property {
    name  = "apiV1UsageOnlineArchiveGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_online_archive_get_start_hr
  }
  property {
    name  = "apiV1UsageProfilingGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_profiling_get_accept
  }
  property {
    name  = "apiV1UsageProfilingGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_profiling_get_end_hr
  }
  property {
    name  = "apiV1UsageProfilingGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_profiling_get_start_hr
  }
  property {
    name  = "apiV1UsageRumGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_rum_get_accept
  }
  property {
    name  = "apiV1UsageRumGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_rum_get_end_hr
  }
  property {
    name  = "apiV1UsageRumGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_rum_get_start_hr
  }
  property {
    name  = "apiV1UsageRumSessionsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_rum_sessions_get_accept
  }
  property {
    name  = "apiV1UsageRumSessionsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_rum_sessions_get_end_hr
  }
  property {
    name  = "apiV1UsageRumSessionsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_rum_sessions_get_start_hr
  }
  property {
    name  = "apiV1UsageRumSessionsGet_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_rum_sessions_get_type
  }
  property {
    name  = "apiV1UsageSdsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_sds_get_accept
  }
  property {
    name  = "apiV1UsageSdsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_sds_get_end_hr
  }
  property {
    name  = "apiV1UsageSdsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_sds_get_start_hr
  }
  property {
    name  = "apiV1UsageSnmpGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_snmp_get_accept
  }
  property {
    name  = "apiV1UsageSnmpGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_snmp_get_end_hr
  }
  property {
    name  = "apiV1UsageSnmpGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_snmp_get_start_hr
  }
  property {
    name  = "apiV1UsageSummaryGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_summary_get_accept
  }
  property {
    name  = "apiV1UsageSummaryGet_end_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_summary_get_end_month
  }
  property {
    name  = "apiV1UsageSummaryGet_include_org_details"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_summary_get_include_org_details
  }
  property {
    name  = "apiV1UsageSummaryGet_start_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_summary_get_start_month
  }
  property {
    name  = "apiV1UsageSyntheticsApiGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_api_get_accept
  }
  property {
    name  = "apiV1UsageSyntheticsApiGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_api_get_end_hr
  }
  property {
    name  = "apiV1UsageSyntheticsApiGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_api_get_start_hr
  }
  property {
    name  = "apiV1UsageSyntheticsBrowserGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_browser_get_accept
  }
  property {
    name  = "apiV1UsageSyntheticsBrowserGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_browser_get_end_hr
  }
  property {
    name  = "apiV1UsageSyntheticsBrowserGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_browser_get_start_hr
  }
  property {
    name  = "apiV1UsageSyntheticsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_get_accept
  }
  property {
    name  = "apiV1UsageSyntheticsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_get_end_hr
  }
  property {
    name  = "apiV1UsageSyntheticsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_synthetics_get_start_hr
  }
  property {
    name  = "apiV1UsageTimeseriesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_timeseries_get_accept
  }
  property {
    name  = "apiV1UsageTimeseriesGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_timeseries_get_end_hr
  }
  property {
    name  = "apiV1UsageTimeseriesGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_timeseries_get_start_hr
  }
  property {
    name  = "apiV1UsageTopAvgMetricsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_top_avg_metrics_get_accept
  }
  property {
    name  = "apiV1UsageTopAvgMetricsGet_day"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_top_avg_metrics_get_day
  }
  property {
    name  = "apiV1UsageTopAvgMetricsGet_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_top_avg_metrics_get_limit
  }
  property {
    name  = "apiV1UsageTopAvgMetricsGet_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_top_avg_metrics_get_month
  }
  property {
    name  = "apiV1UsageTopAvgMetricsGet_names"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_top_avg_metrics_get_names
  }
  property {
    name  = "apiV1UsageTopAvgMetricsGet_next_record_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_usage_top_avg_metrics_get_next_record_id
  }
  property {
    name  = "apiV1ValidateGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v1_validate_get_accept
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_delete_accept
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdDelete_api_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_delete_api_key_id
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_get_accept
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdGet_api_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_get_api_key_id
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_get_include
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_patch_accept
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdPatch_api_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_patch_api_key_id
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_patch_body
  }
  property {
    name  = "apiV2ApiKeysApiKeyIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_api_key_id_patch_content_type
  }
  property {
    name  = "apiV2ApiKeysGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_accept
  }
  property {
    name  = "apiV2ApiKeysGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_filter
  }
  property {
    name  = "apiV2ApiKeysGet_filter_created_at_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_filter_created_at_end
  }
  property {
    name  = "apiV2ApiKeysGet_filter_created_at_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_filter_created_at_start
  }
  property {
    name  = "apiV2ApiKeysGet_filter_modified_at_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_filter_modified_at_end
  }
  property {
    name  = "apiV2ApiKeysGet_filter_modified_at_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_filter_modified_at_start
  }
  property {
    name  = "apiV2ApiKeysGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_include
  }
  property {
    name  = "apiV2ApiKeysGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_page_number
  }
  property {
    name  = "apiV2ApiKeysGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_page_size
  }
  property {
    name  = "apiV2ApiKeysGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_get_sort
  }
  property {
    name  = "apiV2ApiKeysPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_post_accept
  }
  property {
    name  = "apiV2ApiKeysPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_post_body
  }
  property {
    name  = "apiV2ApiKeysPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_api_keys_post_content_type
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_delete_accept
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdDelete_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_delete_app_key_id
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_get_accept
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdGet_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_get_app_key_id
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_get_include
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_patch_accept
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdPatch_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_patch_app_key_id
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_patch_body
  }
  property {
    name  = "apiV2ApplicationKeysAppKeyIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_app_key_id_patch_content_type
  }
  property {
    name  = "apiV2ApplicationKeysGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_get_accept
  }
  property {
    name  = "apiV2ApplicationKeysGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_get_filter
  }
  property {
    name  = "apiV2ApplicationKeysGet_filter_created_at_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_get_filter_created_at_end
  }
  property {
    name  = "apiV2ApplicationKeysGet_filter_created_at_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_get_filter_created_at_start
  }
  property {
    name  = "apiV2ApplicationKeysGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_get_page_number
  }
  property {
    name  = "apiV2ApplicationKeysGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_get_page_size
  }
  property {
    name  = "apiV2ApplicationKeysGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_application_keys_get_sort
  }
  property {
    name  = "apiV2AuditEventsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_get_accept
  }
  property {
    name  = "apiV2AuditEventsGet_filter_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_get_filter_from
  }
  property {
    name  = "apiV2AuditEventsGet_filter_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_get_filter_query
  }
  property {
    name  = "apiV2AuditEventsGet_filter_to"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_get_filter_to
  }
  property {
    name  = "apiV2AuditEventsGet_page_cursor"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_get_page_cursor
  }
  property {
    name  = "apiV2AuditEventsGet_page_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_get_page_limit
  }
  property {
    name  = "apiV2AuditEventsGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_get_sort
  }
  property {
    name  = "apiV2AuditEventsSearchPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_search_post_accept
  }
  property {
    name  = "apiV2AuditEventsSearchPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_search_post_body
  }
  property {
    name  = "apiV2AuditEventsSearchPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_audit_events_search_post_content_type
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_delete_accept
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdDelete_authn_mapping_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_delete_authn_mapping_id
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_get_accept
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdGet_authn_mapping_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_get_authn_mapping_id
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_patch_accept
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdPatch_authn_mapping_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_patch_authn_mapping_id
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_patch_body
  }
  property {
    name  = "apiV2AuthnMappingsAuthnMappingIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_authn_mapping_id_patch_content_type
  }
  property {
    name  = "apiV2AuthnMappingsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_get_accept
  }
  property {
    name  = "apiV2AuthnMappingsGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_get_filter
  }
  property {
    name  = "apiV2AuthnMappingsGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_get_page_number
  }
  property {
    name  = "apiV2AuthnMappingsGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_get_page_size
  }
  property {
    name  = "apiV2AuthnMappingsGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_get_sort
  }
  property {
    name  = "apiV2AuthnMappingsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_post_accept
  }
  property {
    name  = "apiV2AuthnMappingsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_post_body
  }
  property {
    name  = "apiV2AuthnMappingsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_authn_mappings_post_content_type
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_delete_accept
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdDelete_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_delete_app_key_id
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_get_accept
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdGet_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_get_app_key_id
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_patch_accept
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdPatch_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_patch_app_key_id
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_patch_body
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysAppKeyIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_app_key_id_patch_content_type
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_get_accept
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_get_filter
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysGet_filter_created_at_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_get_filter_created_at_end
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysGet_filter_created_at_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_get_filter_created_at_start
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_get_page_number
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_get_page_size
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_get_sort
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_post_accept
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_post_body
  }
  property {
    name  = "apiV2CurrentUserApplicationKeysPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_current_user_application_keys_post_content_type
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_delete_accept
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_delete_content_type
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsDelete_dashboard_list_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_delete_dashboard_list_id
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_get_accept
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsGet_dashboard_list_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_get_dashboard_list_id
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_post_accept
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_post_body
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_post_content_type
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPost_dashboard_list_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_post_dashboard_list_id
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_put_accept
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_put_body
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_put_content_type
  }
  property {
    name  = "apiV2DashboardListsManualDashboardListIdDashboardsPut_dashboard_list_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_dashboard_lists_manual_dashboard_list_id_dashboards_put_dashboard_list_id
  }
  property {
    name  = "apiV2EventsSearchPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_events_search_post_accept
  }
  property {
    name  = "apiV2EventsSearchPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_events_search_post_body
  }
  property {
    name  = "apiV2EventsSearchPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_events_search_post_content_type
  }
  property {
    name  = "apiV2IncidentsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_get_accept
  }
  property {
    name  = "apiV2IncidentsGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_get_include
  }
  property {
    name  = "apiV2IncidentsGet_page_offset"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_get_page_offset
  }
  property {
    name  = "apiV2IncidentsGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_get_page_size
  }
  property {
    name  = "apiV2IncidentsIncidentIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_delete_accept
  }
  property {
    name  = "apiV2IncidentsIncidentIdDelete_incident_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_delete_incident_id
  }
  property {
    name  = "apiV2IncidentsIncidentIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_get_accept
  }
  property {
    name  = "apiV2IncidentsIncidentIdGet_incident_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_get_incident_id
  }
  property {
    name  = "apiV2IncidentsIncidentIdGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_get_include
  }
  property {
    name  = "apiV2IncidentsIncidentIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_patch_accept
  }
  property {
    name  = "apiV2IncidentsIncidentIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_patch_body
  }
  property {
    name  = "apiV2IncidentsIncidentIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_patch_content_type
  }
  property {
    name  = "apiV2IncidentsIncidentIdPatch_incident_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_patch_incident_id
  }
  property {
    name  = "apiV2IncidentsIncidentIdPatch_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_incident_id_patch_include
  }
  property {
    name  = "apiV2IncidentsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_post_accept
  }
  property {
    name  = "apiV2IncidentsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_post_body
  }
  property {
    name  = "apiV2IncidentsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_incidents_post_content_type
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_get_accept
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_delete_accept
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdDelete_integration_service_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_delete_integration_service_id
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_get_accept
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdGet_integration_service_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_get_integration_service_id
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_patch_accept
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_patch_body
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_patch_content_type
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesIntegrationServiceIdPatch_integration_service_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_integration_service_id_patch_integration_service_id
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_post_accept
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_post_body
  }
  property {
    name  = "apiV2IntegrationOpsgenieServicesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_integration_opsgenie_services_post_content_type
  }
  property {
    name  = "apiV2LogsAnalyticsAggregatePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_analytics_aggregate_post_accept
  }
  property {
    name  = "apiV2LogsAnalyticsAggregatePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_analytics_aggregate_post_body
  }
  property {
    name  = "apiV2LogsAnalyticsAggregatePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_analytics_aggregate_post_content_type
  }
  property {
    name  = "apiV2LogsConfigArchiveOrderGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archive_order_get_accept
  }
  property {
    name  = "apiV2LogsConfigArchiveOrderPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archive_order_put_accept
  }
  property {
    name  = "apiV2LogsConfigArchiveOrderPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archive_order_put_body
  }
  property {
    name  = "apiV2LogsConfigArchiveOrderPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archive_order_put_content_type
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_delete_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdDelete_archive_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_delete_archive_id
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_get_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdGet_archive_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_get_archive_id
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_put_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdPut_archive_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_put_archive_id
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_put_body
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_put_content_type
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_delete_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersDelete_archive_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_delete_archive_id
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_delete_content_type
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_get_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersGet_archive_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_get_archive_id
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_post_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersPost_archive_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_post_archive_id
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_post_body
  }
  property {
    name  = "apiV2LogsConfigArchivesArchiveIdReadersPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_archive_id_readers_post_content_type
  }
  property {
    name  = "apiV2LogsConfigArchivesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_get_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_post_accept
  }
  property {
    name  = "apiV2LogsConfigArchivesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_post_body
  }
  property {
    name  = "apiV2LogsConfigArchivesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_archives_post_content_type
  }
  property {
    name  = "apiV2LogsConfigMetricsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_get_accept
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_delete_accept
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdDelete_metric_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_delete_metric_id
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_get_accept
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdGet_metric_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_get_metric_id
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_patch_accept
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_patch_body
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_patch_content_type
  }
  property {
    name  = "apiV2LogsConfigMetricsMetricIdPatch_metric_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_metric_id_patch_metric_id
  }
  property {
    name  = "apiV2LogsConfigMetricsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_post_accept
  }
  property {
    name  = "apiV2LogsConfigMetricsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_post_body
  }
  property {
    name  = "apiV2LogsConfigMetricsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_metrics_post_content_type
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_get_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_get_page_number
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_get_page_size
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_post_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_post_body
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_post_content_type
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_delete_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdDelete_restriction_query_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_delete_restriction_query_id
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_get_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdGet_restriction_query_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_get_restriction_query_id
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_patch_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_patch_body
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_patch_content_type
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdPatch_restriction_query_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_patch_restriction_query_id
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_delete_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_delete_content_type
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesDelete_restriction_query_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_delete_restriction_query_id
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_get_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_get_page_number
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_get_page_size
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesGet_restriction_query_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_get_restriction_query_id
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_post_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_post_body
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_post_content_type
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRestrictionQueryIdRolesPost_restriction_query_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_restriction_query_id_roles_post_restriction_query_id
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRoleRoleIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_role_role_id_get_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesRoleRoleIdGet_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_role_role_id_get_role_id
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesUserUserIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_user_user_id_get_accept
  }
  property {
    name  = "apiV2LogsConfigRestrictionQueriesUserUserIdGet_user_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_config_restriction_queries_user_user_id_get_user_id
  }
  property {
    name  = "apiV2LogsEventsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_accept
  }
  property {
    name  = "apiV2LogsEventsGet_filter_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_filter_from
  }
  property {
    name  = "apiV2LogsEventsGet_filter_index"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_filter_index
  }
  property {
    name  = "apiV2LogsEventsGet_filter_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_filter_query
  }
  property {
    name  = "apiV2LogsEventsGet_filter_to"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_filter_to
  }
  property {
    name  = "apiV2LogsEventsGet_page_cursor"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_page_cursor
  }
  property {
    name  = "apiV2LogsEventsGet_page_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_page_limit
  }
  property {
    name  = "apiV2LogsEventsGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_get_sort
  }
  property {
    name  = "apiV2LogsEventsSearchPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_search_post_accept
  }
  property {
    name  = "apiV2LogsEventsSearchPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_search_post_body
  }
  property {
    name  = "apiV2LogsEventsSearchPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_events_search_post_content_type
  }
  property {
    name  = "apiV2LogsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_post_body
  }
  property {
    name  = "apiV2LogsPost_content_encoding"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_post_content_encoding
  }
  property {
    name  = "apiV2LogsPost_ddtags"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_logs_post_ddtags
  }
  property {
    name  = "apiV2MetricsConfigBulkTagsDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_config_bulk_tags_delete_accept
  }
  property {
    name  = "apiV2MetricsConfigBulkTagsDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_config_bulk_tags_delete_content_type
  }
  property {
    name  = "apiV2MetricsConfigBulkTagsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_config_bulk_tags_post_accept
  }
  property {
    name  = "apiV2MetricsConfigBulkTagsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_config_bulk_tags_post_body
  }
  property {
    name  = "apiV2MetricsConfigBulkTagsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_config_bulk_tags_post_content_type
  }
  property {
    name  = "apiV2MetricsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_get_accept
  }
  property {
    name  = "apiV2MetricsGet_filter_configured"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_get_filter_configured
  }
  property {
    name  = "apiV2MetricsGet_filter_include_percentiles"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_get_filter_include_percentiles
  }
  property {
    name  = "apiV2MetricsGet_filter_metric_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_get_filter_metric_type
  }
  property {
    name  = "apiV2MetricsGet_filter_tags"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_get_filter_tags
  }
  property {
    name  = "apiV2MetricsGet_filter_tags_configured"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_get_filter_tags_configured
  }
  property {
    name  = "apiV2MetricsGet_window_seconds"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_get_window_seconds
  }
  property {
    name  = "apiV2MetricsMetricNameAllTagsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_all_tags_get_accept
  }
  property {
    name  = "apiV2MetricsMetricNameAllTagsGet_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_all_tags_get_metric_name
  }
  property {
    name  = "apiV2MetricsMetricNameEstimateGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_estimate_get_accept
  }
  property {
    name  = "apiV2MetricsMetricNameEstimateGet_filter_groups"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_estimate_get_filter_groups
  }
  property {
    name  = "apiV2MetricsMetricNameEstimateGet_filter_hours_ago"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_estimate_get_filter_hours_ago
  }
  property {
    name  = "apiV2MetricsMetricNameEstimateGet_filter_num_aggregations"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_estimate_get_filter_num_aggregations
  }
  property {
    name  = "apiV2MetricsMetricNameEstimateGet_filter_pct"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_estimate_get_filter_pct
  }
  property {
    name  = "apiV2MetricsMetricNameEstimateGet_filter_timespan_h"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_estimate_get_filter_timespan_h
  }
  property {
    name  = "apiV2MetricsMetricNameEstimateGet_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_estimate_get_metric_name
  }
  property {
    name  = "apiV2MetricsMetricNameTagsDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_delete_accept
  }
  property {
    name  = "apiV2MetricsMetricNameTagsDelete_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_delete_metric_name
  }
  property {
    name  = "apiV2MetricsMetricNameTagsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_get_accept
  }
  property {
    name  = "apiV2MetricsMetricNameTagsGet_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_get_metric_name
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_patch_accept
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_patch_body
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_patch_content_type
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPatch_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_patch_metric_name
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_post_accept
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_post_body
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_post_content_type
  }
  property {
    name  = "apiV2MetricsMetricNameTagsPost_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_tags_post_metric_name
  }
  property {
    name  = "apiV2MetricsMetricNameVolumesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_volumes_get_accept
  }
  property {
    name  = "apiV2MetricsMetricNameVolumesGet_metric_name"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_metrics_metric_name_volumes_get_metric_name
  }
  property {
    name  = "apiV2PermissionsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_permissions_get_accept
  }
  property {
    name  = "apiV2ProcessesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_processes_get_accept
  }
  property {
    name  = "apiV2ProcessesGet_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_processes_get_from
  }
  property {
    name  = "apiV2ProcessesGet_page_cursor"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_processes_get_page_cursor
  }
  property {
    name  = "apiV2ProcessesGet_page_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_processes_get_page_limit
  }
  property {
    name  = "apiV2ProcessesGet_search"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_processes_get_search
  }
  property {
    name  = "apiV2ProcessesGet_tags"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_processes_get_tags
  }
  property {
    name  = "apiV2ProcessesGet_to"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_processes_get_to
  }
  property {
    name  = "apiV2RolesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_get_accept
  }
  property {
    name  = "apiV2RolesGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_get_filter
  }
  property {
    name  = "apiV2RolesGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_get_page_number
  }
  property {
    name  = "apiV2RolesGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_get_page_size
  }
  property {
    name  = "apiV2RolesGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_get_sort
  }
  property {
    name  = "apiV2RolesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_post_accept
  }
  property {
    name  = "apiV2RolesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_post_body
  }
  property {
    name  = "apiV2RolesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_post_content_type
  }
  property {
    name  = "apiV2RolesRoleIdClonePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_clone_post_accept
  }
  property {
    name  = "apiV2RolesRoleIdClonePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_clone_post_body
  }
  property {
    name  = "apiV2RolesRoleIdClonePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_clone_post_content_type
  }
  property {
    name  = "apiV2RolesRoleIdClonePost_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_clone_post_role_id
  }
  property {
    name  = "apiV2RolesRoleIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_delete_accept
  }
  property {
    name  = "apiV2RolesRoleIdDelete_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_delete_role_id
  }
  property {
    name  = "apiV2RolesRoleIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_get_accept
  }
  property {
    name  = "apiV2RolesRoleIdGet_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_get_role_id
  }
  property {
    name  = "apiV2RolesRoleIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_patch_accept
  }
  property {
    name  = "apiV2RolesRoleIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_patch_body
  }
  property {
    name  = "apiV2RolesRoleIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_patch_content_type
  }
  property {
    name  = "apiV2RolesRoleIdPatch_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_patch_role_id
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_delete_accept
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_delete_content_type
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsDelete_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_delete_role_id
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_get_accept
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsGet_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_get_role_id
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_post_accept
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_post_body
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_post_content_type
  }
  property {
    name  = "apiV2RolesRoleIdPermissionsPost_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_permissions_post_role_id
  }
  property {
    name  = "apiV2RolesRoleIdUsersDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_delete_accept
  }
  property {
    name  = "apiV2RolesRoleIdUsersDelete_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_delete_content_type
  }
  property {
    name  = "apiV2RolesRoleIdUsersDelete_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_delete_role_id
  }
  property {
    name  = "apiV2RolesRoleIdUsersGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_get_accept
  }
  property {
    name  = "apiV2RolesRoleIdUsersGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_get_filter
  }
  property {
    name  = "apiV2RolesRoleIdUsersGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_get_page_number
  }
  property {
    name  = "apiV2RolesRoleIdUsersGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_get_page_size
  }
  property {
    name  = "apiV2RolesRoleIdUsersGet_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_get_role_id
  }
  property {
    name  = "apiV2RolesRoleIdUsersGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_get_sort
  }
  property {
    name  = "apiV2RolesRoleIdUsersPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_post_accept
  }
  property {
    name  = "apiV2RolesRoleIdUsersPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_post_body
  }
  property {
    name  = "apiV2RolesRoleIdUsersPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_post_content_type
  }
  property {
    name  = "apiV2RolesRoleIdUsersPost_role_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_roles_role_id_users_post_role_id
  }
  property {
    name  = "apiV2RumAnalyticsAggregatePost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_analytics_aggregate_post_accept
  }
  property {
    name  = "apiV2RumAnalyticsAggregatePost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_analytics_aggregate_post_body
  }
  property {
    name  = "apiV2RumAnalyticsAggregatePost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_analytics_aggregate_post_content_type
  }
  property {
    name  = "apiV2RumEventsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_get_accept
  }
  property {
    name  = "apiV2RumEventsGet_filter_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_get_filter_from
  }
  property {
    name  = "apiV2RumEventsGet_filter_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_get_filter_query
  }
  property {
    name  = "apiV2RumEventsGet_filter_to"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_get_filter_to
  }
  property {
    name  = "apiV2RumEventsGet_page_cursor"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_get_page_cursor
  }
  property {
    name  = "apiV2RumEventsGet_page_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_get_page_limit
  }
  property {
    name  = "apiV2RumEventsGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_get_sort
  }
  property {
    name  = "apiV2RumEventsSearchPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_search_post_accept
  }
  property {
    name  = "apiV2RumEventsSearchPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_search_post_body
  }
  property {
    name  = "apiV2RumEventsSearchPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_rum_events_search_post_content_type
  }
  property {
    name  = "apiV2SamlConfigurationsIdpMetadataPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_saml_configurations_idp_metadata_post_accept
  }
  property {
    name  = "apiV2SamlConfigurationsIdpMetadataPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_saml_configurations_idp_metadata_post_content_type
  }
  property {
    name  = "apiV2SamlConfigurationsIdpMetadataPost_idp_file"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_saml_configurations_idp_metadata_post_idp_file
  }
  property {
    name  = "apiV2SecurityCloudWorkloadPolicyDownloadGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_cloud_workload_policy_download_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_delete_accept
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdDelete_agent_rule_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_delete_agent_rule_id
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdGet_agent_rule_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_get_agent_rule_id
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_patch_accept
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdPatch_agent_rule_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_patch_agent_rule_id
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_patch_body
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesAgentRuleIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_agent_rule_id_patch_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_post_accept
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_post_body
  }
  property {
    name  = "apiV2SecurityMonitoringCloudWorkloadSecurityAgentRulesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_cloud_workload_security_agent_rules_post_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_post_accept
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_post_body
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_post_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_delete_accept
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdDelete_security_filter_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_delete_security_filter_id
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdGet_security_filter_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_get_security_filter_id
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_patch_accept
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_patch_body
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_patch_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringConfigurationSecurityFiltersSecurityFilterIdPatch_security_filter_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_configuration_security_filters_security_filter_id_patch_security_filter_id
  }
  property {
    name  = "apiV2SecurityMonitoringRulesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringRulesGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_get_page_number
  }
  property {
    name  = "apiV2SecurityMonitoringRulesGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_get_page_size
  }
  property {
    name  = "apiV2SecurityMonitoringRulesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_post_accept
  }
  property {
    name  = "apiV2SecurityMonitoringRulesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_post_body
  }
  property {
    name  = "apiV2SecurityMonitoringRulesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_post_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_delete_accept
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdDelete_rule_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_delete_rule_id
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdGet_rule_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_get_rule_id
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdPut_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_put_accept
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdPut_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_put_body
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdPut_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_put_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringRulesRuleIdPut_rule_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_rules_rule_id_put_rule_id
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_get_accept
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsGet_filter_from"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_get_filter_from
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsGet_filter_query"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_get_filter_query
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsGet_filter_to"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_get_filter_to
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsGet_page_cursor"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_get_page_cursor
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsGet_page_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_get_page_limit
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_get_sort
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSearchPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_search_post_accept
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSearchPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_search_post_body
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSearchPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_search_post_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdAssigneePatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_assignee_patch_accept
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdAssigneePatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_assignee_patch_body
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdAssigneePatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_assignee_patch_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdAssigneePatch_signal_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_assignee_patch_signal_id
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdIncidentsPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_incidents_patch_accept
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdIncidentsPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_incidents_patch_body
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdIncidentsPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_incidents_patch_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdIncidentsPatch_signal_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_incidents_patch_signal_id
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdStatePatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_state_patch_accept
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdStatePatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_state_patch_body
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdStatePatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_state_patch_content_type
  }
  property {
    name  = "apiV2SecurityMonitoringSignalsSignalIdStatePatch_signal_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_security_monitoring_signals_signal_id_state_patch_signal_id
  }
  property {
    name  = "apiV2SeriesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_series_post_accept
  }
  property {
    name  = "apiV2SeriesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_series_post_body
  }
  property {
    name  = "apiV2SeriesPost_content_encoding"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_series_post_content_encoding
  }
  property {
    name  = "apiV2SeriesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_series_post_content_type
  }
  property {
    name  = "apiV2ServiceAccountsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_post_accept
  }
  property {
    name  = "apiV2ServiceAccountsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_post_body
  }
  property {
    name  = "apiV2ServiceAccountsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_post_content_type
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_delete_accept
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdDelete_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_delete_app_key_id
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdDelete_service_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_delete_service_account_id
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_get_accept
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdGet_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_get_app_key_id
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdGet_service_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_get_service_account_id
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_patch_accept
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdPatch_app_key_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_patch_app_key_id
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_patch_body
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_patch_content_type
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysAppKeyIdPatch_service_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_app_key_id_patch_service_account_id
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_accept
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_filter
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_filter_created_at_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_filter_created_at_end
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_filter_created_at_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_filter_created_at_start
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_page_number
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_page_size
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_service_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_service_account_id
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_get_sort
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_post_accept
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_post_body
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_post_content_type
  }
  property {
    name  = "apiV2ServiceAccountsServiceAccountIdApplicationKeysPost_service_account_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_service_accounts_service_account_id_application_keys_post_service_account_id
  }
  property {
    name  = "apiV2ServicesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_get_accept
  }
  property {
    name  = "apiV2ServicesGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_get_filter
  }
  property {
    name  = "apiV2ServicesGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_get_include
  }
  property {
    name  = "apiV2ServicesGet_page_offset"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_get_page_offset
  }
  property {
    name  = "apiV2ServicesGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_get_page_size
  }
  property {
    name  = "apiV2ServicesPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_post_accept
  }
  property {
    name  = "apiV2ServicesPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_post_body
  }
  property {
    name  = "apiV2ServicesPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_post_content_type
  }
  property {
    name  = "apiV2ServicesServiceIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_delete_accept
  }
  property {
    name  = "apiV2ServicesServiceIdDelete_service_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_delete_service_id
  }
  property {
    name  = "apiV2ServicesServiceIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_get_accept
  }
  property {
    name  = "apiV2ServicesServiceIdGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_get_include
  }
  property {
    name  = "apiV2ServicesServiceIdGet_service_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_get_service_id
  }
  property {
    name  = "apiV2ServicesServiceIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_patch_accept
  }
  property {
    name  = "apiV2ServicesServiceIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_patch_body
  }
  property {
    name  = "apiV2ServicesServiceIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_patch_content_type
  }
  property {
    name  = "apiV2ServicesServiceIdPatch_service_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_services_service_id_patch_service_id
  }
  property {
    name  = "apiV2TeamsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_get_accept
  }
  property {
    name  = "apiV2TeamsGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_get_filter
  }
  property {
    name  = "apiV2TeamsGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_get_include
  }
  property {
    name  = "apiV2TeamsGet_page_offset"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_get_page_offset
  }
  property {
    name  = "apiV2TeamsGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_get_page_size
  }
  property {
    name  = "apiV2TeamsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_post_accept
  }
  property {
    name  = "apiV2TeamsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_post_body
  }
  property {
    name  = "apiV2TeamsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_post_content_type
  }
  property {
    name  = "apiV2TeamsTeamIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_delete_accept
  }
  property {
    name  = "apiV2TeamsTeamIdDelete_team_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_delete_team_id
  }
  property {
    name  = "apiV2TeamsTeamIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_get_accept
  }
  property {
    name  = "apiV2TeamsTeamIdGet_include"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_get_include
  }
  property {
    name  = "apiV2TeamsTeamIdGet_team_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_get_team_id
  }
  property {
    name  = "apiV2TeamsTeamIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_patch_accept
  }
  property {
    name  = "apiV2TeamsTeamIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_patch_body
  }
  property {
    name  = "apiV2TeamsTeamIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_patch_content_type
  }
  property {
    name  = "apiV2TeamsTeamIdPatch_team_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_teams_team_id_patch_team_id
  }
  property {
    name  = "apiV2UsageApplicationSecurityGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_application_security_get_accept
  }
  property {
    name  = "apiV2UsageApplicationSecurityGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_application_security_get_end_hr
  }
  property {
    name  = "apiV2UsageApplicationSecurityGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_application_security_get_start_hr
  }
  property {
    name  = "apiV2UsageCostByOrgGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_cost_by_org_get_accept
  }
  property {
    name  = "apiV2UsageCostByOrgGet_end_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_cost_by_org_get_end_month
  }
  property {
    name  = "apiV2UsageCostByOrgGet_start_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_cost_by_org_get_start_month
  }
  property {
    name  = "apiV2UsageEstimatedCostGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_estimated_cost_get_accept
  }
  property {
    name  = "apiV2UsageEstimatedCostGet_end_date"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_estimated_cost_get_end_date
  }
  property {
    name  = "apiV2UsageEstimatedCostGet_end_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_estimated_cost_get_end_month
  }
  property {
    name  = "apiV2UsageEstimatedCostGet_start_date"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_estimated_cost_get_start_date
  }
  property {
    name  = "apiV2UsageEstimatedCostGet_start_month"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_estimated_cost_get_start_month
  }
  property {
    name  = "apiV2UsageEstimatedCostGet_view"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_estimated_cost_get_view
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_accept
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_filter_include_descendants"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_filter_include_descendants
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_filter_product_families"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_filter_product_families
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_filter_timestamp_end"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_filter_timestamp_end
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_filter_timestamp_start"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_filter_timestamp_start
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_filter_versions"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_filter_versions
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_page_limit"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_page_limit
  }
  property {
    name  = "apiV2UsageHourlyUsageGet_page_next_record_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_hourly_usage_get_page_next_record_id
  }
  property {
    name  = "apiV2UsageLambdaTracedInvocationsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_lambda_traced_invocations_get_accept
  }
  property {
    name  = "apiV2UsageLambdaTracedInvocationsGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_lambda_traced_invocations_get_end_hr
  }
  property {
    name  = "apiV2UsageLambdaTracedInvocationsGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_lambda_traced_invocations_get_start_hr
  }
  property {
    name  = "apiV2UsageObservabilityPipelinesGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_observability_pipelines_get_accept
  }
  property {
    name  = "apiV2UsageObservabilityPipelinesGet_end_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_observability_pipelines_get_end_hr
  }
  property {
    name  = "apiV2UsageObservabilityPipelinesGet_start_hr"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_usage_observability_pipelines_get_start_hr
  }
  property {
    name  = "apiV2UserInvitationsPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_user_invitations_post_accept
  }
  property {
    name  = "apiV2UserInvitationsPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_user_invitations_post_body
  }
  property {
    name  = "apiV2UserInvitationsPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_user_invitations_post_content_type
  }
  property {
    name  = "apiV2UserInvitationsUserInvitationUuidGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_user_invitations_user_invitation_uuid_get_accept
  }
  property {
    name  = "apiV2UserInvitationsUserInvitationUuidGet_user_invitation_uuid"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_user_invitations_user_invitation_uuid_get_user_invitation_uuid
  }
  property {
    name  = "apiV2UsersGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_get_accept
  }
  property {
    name  = "apiV2UsersGet_filter"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_get_filter
  }
  property {
    name  = "apiV2UsersGet_filter_status"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_get_filter_status
  }
  property {
    name  = "apiV2UsersGet_page_number"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_get_page_number
  }
  property {
    name  = "apiV2UsersGet_page_size"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_get_page_size
  }
  property {
    name  = "apiV2UsersGet_sort"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_get_sort
  }
  property {
    name  = "apiV2UsersGet_sort_dir"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_get_sort_dir
  }
  property {
    name  = "apiV2UsersPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_post_accept
  }
  property {
    name  = "apiV2UsersPost_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_post_body
  }
  property {
    name  = "apiV2UsersPost_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_post_content_type
  }
  property {
    name  = "apiV2UsersUserIdDelete_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_delete_accept
  }
  property {
    name  = "apiV2UsersUserIdDelete_user_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_delete_user_id
  }
  property {
    name  = "apiV2UsersUserIdGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_get_accept
  }
  property {
    name  = "apiV2UsersUserIdGet_user_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_get_user_id
  }
  property {
    name  = "apiV2UsersUserIdOrgsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_orgs_get_accept
  }
  property {
    name  = "apiV2UsersUserIdOrgsGet_user_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_orgs_get_user_id
  }
  property {
    name  = "apiV2UsersUserIdPatch_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_patch_accept
  }
  property {
    name  = "apiV2UsersUserIdPatch_body"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_patch_body
  }
  property {
    name  = "apiV2UsersUserIdPatch_content_type"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_patch_content_type
  }
  property {
    name  = "apiV2UsersUserIdPatch_user_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_patch_user_id
  }
  property {
    name  = "apiV2UsersUserIdPermissionsGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_permissions_get_accept
  }
  property {
    name  = "apiV2UsersUserIdPermissionsGet_user_id"
    type  = "string"
    value = var.connector-oai-datadogapi_property_api_v2_users_user_id_permissions_get_user_id
  }
  property {
    name  = "authApiKey"
    type  = "string"
    value = var.connector-oai-datadogapi_property_auth_api_key
  }
  property {
    name  = "authApplicationKey"
    type  = "string"
    value = var.connector-oai-datadogapi_property_auth_application_key
  }
  property {
    name  = "basePath"
    type  = "string"
    value = var.connector-oai-datadogapi_property_base_path
  }
  property {
    name  = "rootGet_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_root_get_accept
  }
  property {
    name  = "v1MonitorMuteAllPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_v1_monitor_mute_all_post_accept
  }
  property {
    name  = "v1MonitorUnmuteAllPost_accept"
    type  = "string"
    value = var.connector-oai-datadogapi_property_v1_monitor_unmute_all_post_accept
  }
}
