//go:generate go run ./cmd/generate/generate.go

package dvgenerate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/dvgenerate/internal"
)

type connectorDocData struct {
	ConnectorName string
	ConnectorId   string
	Properties    []connectorDocPropertyData
}

type connectorDocPropertyData struct {
	Name               string
	Type               *string
	Description        *string
	ConsoleDisplayName *string
	Value              *string
}

func Generate() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	baseDirectory := absDir
	if strings.HasSuffix(baseDirectory, "tools/dvgenerate") {
		baseDirectory = filepath.Dir(filepath.Dir(baseDirectory))
	}

	conns, err := readConnectors()
	if err != nil {
		panic(err)
	}

	GenerateReferenceTemplate(baseDirectory, conns)
	GenerateConnectorHCLExamples(baseDirectory, conns)
}

func GenerateReferenceTemplate(baseDirectory string, conns []connectorDocData) {

	fileNameDirectory := fmt.Sprintf("%s/templates/guides", baseDirectory)

	err := writeConnectorTemplate(fileNameDirectory, internal.ConnectorReferenceTmpl, conns, true)
	if err != nil {
		panic(err)
	}
}

func GenerateConnectorHCLExamples(baseDirectory string, conns []connectorDocData) {

	fileNameDirectory := fmt.Sprintf("%s/examples/davinci-connector-instances", baseDirectory)
	for _, conn := range conns {
		err := writeConnector(fileNameDirectory, internal.ConnectorTmpl, conn, false)
		if err != nil {
			panic(err)
		}
	}
}

func writeConnectorTemplate(fileNameDirectory, templateString string, conns []connectorDocData, overwrite bool) error {
	fileName := fmt.Sprintf("%s/connector-instance-reference.md.tmpl", fileNameDirectory)

	t, err := template.New("ConnectorReferenceTemplate").Parse(templateString)
	if err != nil {
		return err
	}

	return writeTemplateFile(t, fileName, overwrite, conns)
}

func writeConnector(fileNameDirectory, templateString string, conn connectorDocData, overwrite bool) error {
	fileName := fmt.Sprintf("%s/%s.tf", fileNameDirectory, conn.ConnectorId)

	t, err := template.New(fmt.Sprintf("Connector-%s", conn.ConnectorId)).Parse(templateString)
	if err != nil {
		return err
	}

	return writeTemplateFile(t, fileName, overwrite, conn)
}

func writeTemplateFile(t *template.Template, fileName string, overwrite bool, data any) error {
	// Check if the file exists
	if _, err := os.Stat(fileName); err == nil {
		if !overwrite {
			return fmt.Errorf("file %s already exists and overwrite is set to false", fileName)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file exists: %v", err)
	}

	fileName = filepath.Clean(fileName)
	outputFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = t.Execute(outputFile, data)
	if err != nil {
		return err
	}

	return nil
}

type connectorSchemaWrapper struct {
	pingone.DaVinciConnectorDetailsResponse
	ConnectorID *string `json:"connectorId,omitempty"`
	Name        *string `json:"name,omitempty"`
}

func (c *connectorSchemaWrapper) UnmarshalJSON(data []byte) error {
	type localConnectorSchema struct {
		ConnectorID *string                `json:"connectorId,omitempty"`
		Name        *string                `json:"name,omitempty"`
		Properties  map[string]interface{} `json:"properties,omitempty"`
	}

	var aux localConnectorSchema
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.ConnectorID = aux.ConnectorID
	c.Name = aux.Name
	c.Properties = aux.Properties

	return nil
}

func readConnectors() ([]connectorDocData, error) {

	connectorSchema := []connectorSchemaWrapper{}
	err := json.Unmarshal(internal.ConnectorSchemaBytes, &connectorSchema)
	if err != nil {
		return nil, err
	}

	connectorList := make([]connectorDocData, 0)
	for _, connectorSchemaItem := range connectorSchema {

		connectorDocData := connectorDocData{}

		if v := connectorSchemaItem.ConnectorID; v != nil && *v != "" {
			connectorDocData.ConnectorId = *v
		} else {
			connectorDocData.ConnectorId = "No value"
		}

		if v := connectorSchemaItem.Name; v != nil && *v != "" {
			connectorDocData.ConnectorName = *v
		} else {
			connectorDocData.ConnectorName = "No name"
		}

		if connectorSchemaItem.Properties != nil {
			connectorProperties := make([]connectorDocPropertyData, 0)
			for propertyName, property := range connectorSchemaItem.Properties {

				propertyMap, ok := property.(map[string]interface{})
				if !ok {
					continue
				}

				propertyType := "string"
				if v, ok := propertyMap["type"].(string); ok && v != "" {
					propertyType = v
				}

				propertyType = rewritePropertyType(propertyType)

				description := ""
				if v, ok := propertyMap["info"].(string); ok {
					description = v
				}

				if strings.TrimSpace(description) != "" && !strings.HasSuffix(strings.TrimSpace(description), ".") {
					description = fmt.Sprintf("%s.", description)
				}

				var descriptionPtr *string
				if description != "" {
					descriptionPtr = &description
				}

				var displayName string
				if v, ok := propertyMap["displayName"].(string); ok {
					displayName = v
				}

				connectorProperty := connectorDocPropertyData{
					Name:               propertyName,
					Type:               &propertyType,
					Description:        descriptionPtr,
					ConsoleDisplayName: &displayName,
				}

				exampleFound := false
				if v, ok := internal.ExampleValues[connectorDocData.ConnectorId][propertyName]; ok {
					connectorProperty.Value = &v.Value

					if v.OverridingType != nil {
						connectorProperty.Type = v.OverridingType
					}

					exampleFound = true
				}

				if !exampleFound {
					defaultValue := fmt.Sprintf("var.%s_property_%s", strings.ToLower(connectorDocData.ConnectorId), camelToSnake(propertyName))
					connectorProperty.Value = &defaultValue
				}

				connectorProperties = append(connectorProperties, connectorProperty)
			}

			slices.SortFunc(connectorProperties, func(i, j connectorDocPropertyData) int {
				return strings.Compare(i.Name, j.Name)
			})
			connectorDocData.Properties = connectorProperties
		}

		connectorList = append(connectorList, connectorDocData)
	}

	slices.SortFunc(connectorList, func(i, j connectorDocData) int {
		return strings.Compare(i.ConnectorName, j.ConnectorName)
	})

	return connectorList, nil
}

func rewritePropertyType(dvSchemaPropertyType string) (propertyType string) {
	switch dvSchemaPropertyType {
	case "string", "boolean", "number", "":
		propertyType = dvSchemaPropertyType
	default:
		propertyType = "json"
	}

	return
}

func camelToSnake(camel string) string {
	// A buffer to build the output string
	var buf bytes.Buffer

	// Loop through each rune in the string
	for i, r := range camel {
		// If the rune is an uppercase letter and it's not the first character,
		// write an underscore to the buffer
		if unicode.IsUpper(r) && i > 0 {
			buf.WriteRune('_')
		}
		// Write the lowercase version of the current rune to the buffer
		buf.WriteRune(unicode.ToLower(r))
	}

	// Return the contents of the buffer as a string
	return buf.String()
}
