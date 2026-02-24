package internal

import (
	_ "embed"
)

//go:embed connector_schema/connector-schema.json
var ConnectorSchemaBytes []byte
