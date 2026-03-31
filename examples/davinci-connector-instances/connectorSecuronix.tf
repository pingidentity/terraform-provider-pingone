resource "pingone_davinci_connector_instance" "connectorSecuronix" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSecuronix"
  }
  name = "My awesome connectorSecuronix"
  property {
    name  = "allowChildCases"
    type  = "string"
    value = var.connectorsecuronix_property_allow_child_cases
  }
  property {
    name  = "city"
    type  = "string"
    value = var.connectorsecuronix_property_city
  }
  property {
    name  = "companycode"
    type  = "string"
    value = var.connectorsecuronix_property_companycode
  }
  property {
    name  = "costcentername"
    type  = "string"
    value = var.connectorsecuronix_property_costcentername
  }
  property {
    name  = "country"
    type  = "string"
    value = var.connectorsecuronix_property_country
  }
  property {
    name  = "countrycode"
    type  = "string"
    value = var.connectorsecuronix_property_countrycode
  }
  property {
    name  = "department"
    type  = "string"
    value = var.connectorsecuronix_property_department
  }
  property {
    name  = "division"
    type  = "string"
    value = var.connectorsecuronix_property_division
  }
  property {
    name  = "domainName"
    type  = "string"
    value = var.connectorsecuronix_property_domain_name
  }
  property {
    name  = "employeeId"
    type  = "string"
    value = var.connectorsecuronix_property_employee_id
  }
  property {
    name  = "employeetype"
    type  = "string"
    value = var.connectorsecuronix_property_employeetype
  }
  property {
    name  = "employeetypedescription"
    type  = "string"
    value = var.connectorsecuronix_property_employeetypedescription
  }
  property {
    name  = "eventTimeFrom"
    type  = "string"
    value = var.connectorsecuronix_property_event_time_from
  }
  property {
    name  = "eventTimeTo"
    type  = "string"
    value = var.connectorsecuronix_property_event_time_to
  }
  property {
    name  = "firstname"
    type  = "string"
    value = var.connectorsecuronix_property_firstname
  }
  property {
    name  = "from"
    type  = "string"
    value = var.connectorsecuronix_property_from
  }
  property {
    name  = "hiredate"
    type  = "string"
    value = var.connectorsecuronix_property_hiredate
  }
  property {
    name  = "incidentId"
    type  = "string"
    value = var.connectorsecuronix_property_incident_id
  }
  property {
    name  = "ipfrom"
    type  = "string"
    value = var.connectorsecuronix_property_ipfrom
  }
  property {
    name  = "ipto"
    type  = "string"
    value = var.connectorsecuronix_property_ipto
  }
  property {
    name  = "jobcode"
    type  = "string"
    value = var.connectorsecuronix_property_jobcode
  }
  property {
    name  = "key"
    type  = "string"
    value = var.connectorsecuronix_property_key
  }
  property {
    name  = "lanid"
    type  = "string"
    value = var.connectorsecuronix_property_lanid
  }
  property {
    name  = "lastname"
    type  = "string"
    value = var.connectorsecuronix_property_lastname
  }
  property {
    name  = "latitude"
    type  = "string"
    value = var.connectorsecuronix_property_latitude
  }
  property {
    name  = "location"
    type  = "string"
    value = var.connectorsecuronix_property_location
  }
  property {
    name  = "longitude"
    type  = "string"
    value = var.connectorsecuronix_property_longitude
  }
  property {
    name  = "lookupname"
    type  = "string"
    value = var.connectorsecuronix_property_lookupname
  }
  property {
    name  = "manageremployeeid"
    type  = "string"
    value = var.connectorsecuronix_property_manageremployeeid
  }
  property {
    name  = "max"
    type  = "string"
    value = var.connectorsecuronix_property_max
  }
  property {
    name  = "offset"
    type  = "string"
    value = var.connectorsecuronix_property_offset
  }
  property {
    name  = "prettyJson"
    type  = "string"
    value = var.connectorsecuronix_property_pretty_json
  }
  property {
    name  = "queryId"
    type  = "string"
    value = var.connectorsecuronix_property_query_id
  }
  property {
    name  = "rangeType"
    type  = "string"
    value = var.connectorsecuronix_property_range_type
  }
  property {
    name  = "region"
    type  = "string"
    value = var.connectorsecuronix_property_region
  }
  property {
    name  = "sort"
    type  = "string"
    value = var.connectorsecuronix_property_sort
  }
  property {
    name  = "source"
    type  = "string"
    value = var.connectorsecuronix_property_source
  }
  property {
    name  = "status"
    type  = "string"
    value = var.connectorsecuronix_property_status
  }
  property {
    name  = "statusdescription"
    type  = "string"
    value = var.connectorsecuronix_property_statusdescription
  }
  property {
    name  = "timeZone"
    type  = "string"
    value = var.connectorsecuronix_property_time_zone
  }
  property {
    name  = "title"
    type  = "string"
    value = var.connectorsecuronix_property_title
  }
  property {
    name  = "to"
    type  = "string"
    value = var.connectorsecuronix_property_to
  }
  property {
    name  = "token"
    type  = "string"
    value = var.connectorsecuronix_property_token
  }
  property {
    name  = "tpiAddress"
    type  = "string"
    value = var.connectorsecuronix_property_tpi_address
  }
  property {
    name  = "tpiCategory"
    type  = "string"
    value = var.connectorsecuronix_property_tpi_category
  }
  property {
    name  = "tpiCriticality"
    type  = "string"
    value = var.connectorsecuronix_property_tpi_criticality
  }
  property {
    name  = "tpiDate"
    type  = "string"
    value = var.connectorsecuronix_property_tpi_date
  }
  property {
    name  = "tpiDomain"
    type  = "string"
    value = var.connectorsecuronix_property_tpi_domain
  }
  property {
    name  = "tpiSrc"
    type  = "string"
    value = var.connectorsecuronix_property_tpi_src
  }
  property {
    name  = "tpiType"
    type  = "string"
    value = var.connectorsecuronix_property_tpi_type
  }
  property {
    name  = "userid"
    type  = "string"
    value = var.connectorsecuronix_property_userid
  }
  property {
    name  = "valueUCustomfield11"
    type  = "string"
    value = var.connectorsecuronix_property_value_ucustomfield11
  }
  property {
    name  = "valueUCustomfield4"
    type  = "string"
    value = var.connectorsecuronix_property_value_ucustomfield4
  }
  property {
    name  = "violator"
    type  = "string"
    value = var.connectorsecuronix_property_violator
  }
  property {
    name  = "workemail"
    type  = "string"
    value = var.connectorsecuronix_property_workemail
  }
  property {
    name  = "workphone"
    type  = "string"
    value = var.connectorsecuronix_property_workphone
  }
}
