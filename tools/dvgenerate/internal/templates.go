package internal

import _ "embed"

//go:embed templates/connector.tmpl
var ConnectorTmpl string

//go:embed templates/connector_reference.tmpl
var ConnectorReferenceTmpl string
