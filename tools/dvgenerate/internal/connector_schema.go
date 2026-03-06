package internal

import (
	_ "embed"
)

//go:embed connector-schema.json
var ConnectorSchemaBytes []byte
