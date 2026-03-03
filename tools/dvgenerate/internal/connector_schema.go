package internal

import (
	"fmt"
	"os"
)

var ConnectorSchemaBytes []byte

func init() {
	var err error
	ConnectorSchemaBytes, err = os.ReadFile("../../bin/connector-schema.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to read connector schema from ../../bin/connector-schema.json: %v", err))
	}
}
