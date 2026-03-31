resource "pingone_davinci_connector_instance" "pingOneLDAPConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneLDAPConnector"
  }
  name = "My awesome pingOneLDAPConnector"
  property {
    name  = "attributeName"
    type  = "string"
    value = var.pingoneldapconnector_property_attribute_name
  }
  property {
    name  = "attributes"
    type  = "string"
    value = var.pingoneldapconnector_property_attributes
  }
  property {
    name  = "baseDN"
    type  = "string"
    value = var.pingoneldapconnector_property_base_dn
  }
  property {
    name  = "bindDN"
    type  = "string"
    value = var.pingoneldapconnector_property_bind_dn
  }
  property {
    name  = "bindPassword"
    type  = "string"
    value = var.pingoneldapconnector_property_bind_password
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.pingoneldapconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.pingoneldapconnector_property_client_secret
  }
  property {
    name  = "controls"
    type  = "string"
    value = var.pingoneldapconnector_property_controls
  }
  property {
    name  = "dn"
    type  = "string"
    value = var.pingoneldapconnector_property_dn
  }
  property {
    name  = "entryAttributes"
    type  = "string"
    value = var.pingoneldapconnector_property_entry_attributes
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingoneldapconnector_property_env_id
  }
  property {
    name  = "filter"
    type  = "string"
    value = var.pingoneldapconnector_property_filter
  }
  property {
    name  = "gatewayId"
    type  = "string"
    value = var.pingoneldapconnector_property_gateway_id
  }
  property {
    name  = "jsonAttributes"
    type  = "string"
    value = var.pingoneldapconnector_property_json_attributes
  }
  property {
    name  = "ldapUrl"
    type  = "string"
    value = var.pingoneldapconnector_property_ldap_url
  }
  property {
    name  = "modifications"
    type  = "string"
    value = var.pingoneldapconnector_property_modifications
  }
  property {
    name  = "newDn"
    type  = "string"
    value = var.pingoneldapconnector_property_new_dn
  }
  property {
    name  = "newPassword"
    type  = "string"
    value = var.pingoneldapconnector_property_new_password
  }
  property {
    name  = "oldPassword"
    type  = "string"
    value = var.pingoneldapconnector_property_old_password
  }
  property {
    name  = "password"
    type  = "string"
    value = var.pingoneldapconnector_property_password
  }
  property {
    name  = "proxyAuthzDn"
    type  = "string"
    value = var.pingoneldapconnector_property_proxy_authz_dn
  }
  property {
    name  = "proxyAuthzUser"
    type  = "string"
    value = var.pingoneldapconnector_property_proxy_authz_user
  }
  property {
    name  = "proxyAuthzUsername"
    type  = "string"
    value = var.pingoneldapconnector_property_proxy_authz_username
  }
  property {
    name  = "region"
    type  = "string"
    value = var.pingoneldapconnector_property_region
  }
  property {
    name  = "retrieveOperationalAttributes"
    type  = "string"
    value = var.pingoneldapconnector_property_retrieve_operational_attributes
  }
  property {
    name  = "scope"
    type  = "string"
    value = var.pingoneldapconnector_property_scope
  }
  property {
    name  = "sizeLimit"
    type  = "string"
    value = var.pingoneldapconnector_property_size_limit
  }
  property {
    name  = "timeLimit"
    type  = "string"
    value = var.pingoneldapconnector_property_time_limit
  }
  property {
    name  = "typesOnly"
    type  = "string"
    value = var.pingoneldapconnector_property_types_only
  }
  property {
    name  = "useJsonAttributes"
    type  = "string"
    value = var.pingoneldapconnector_property_use_json_attributes
  }
}
