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

	"github.com/pingidentity/terraform-provider-pingone/dvgenerate/internal"
)

type connectorDocData struct {
	ConnectorName string                     `json:"name,omitempty"`
	ConnectorId   string                     `json:"connectorId,omitempty"`
	RawProperties map[string]any             `json:"properties,omitempty"`
	Properties    []connectorDocPropertyData `json:"-"`
}

type connectorDocPropertyData struct {
	Name               string
	Type               *string
	Description        *string
	ConsoleDisplayName *string
	Value              *string
}

func Generate(input []byte) {
	conns, err := readConnectors(input)
	if err != nil {
		fmt.Printf("Warning: %s\n", err)
		return
	}

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
		err := writeConnector(fileNameDirectory, internal.ConnectorTmpl, conn, true)
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

func readConnectors(input []byte) ([]connectorDocData, error) {

	connectorList := []connectorDocData{}

	if len(input) == 0 {
		return nil, fmt.Errorf("no connector schema input provided. Please provide schema via -file flag or stdin.  Skipping generation")
	}

	err := json.Unmarshal(input, &connectorList)
	if err != nil {
		return nil, err
	}

	for i := range connectorList {
		if connectorList[i].ConnectorId == "" {
			connectorList[i].ConnectorId = "No value"
		}

		if connectorList[i].ConnectorName == "" {
			connectorList[i].ConnectorName = "No name"
		}

		connectorList[i].ConnectorName = sanitizeForTemplate(connectorList[i].ConnectorName)

		if connectorList[i].RawProperties != nil {
			connectorProperties := make([]connectorDocPropertyData, 0)
			for propertyName, property := range connectorList[i].RawProperties {

				propertyMap, ok := property.(map[string]any)
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

				description = sanitizeForTemplate(description)

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

				displayName = sanitizeForTemplate(displayName)

				connectorProperty := connectorDocPropertyData{
					Name:               propertyName,
					Type:               &propertyType,
					Description:        descriptionPtr,
					ConsoleDisplayName: &displayName,
				}

				exampleFound := false
				if v, ok := internal.ExampleValues[connectorList[i].ConnectorId][propertyName]; ok {
					connectorProperty.Value = &v.Value

					if v.OverridingType != nil {
						connectorProperty.Type = v.OverridingType
					}

					exampleFound = true
				}

				if !exampleFound {
					defaultValue := fmt.Sprintf("var.%s_property_%s", strings.ToLower(connectorList[i].ConnectorId), camelToSnake(propertyName))
					connectorProperty.Value = &defaultValue
				}

				connectorProperties = append(connectorProperties, connectorProperty)
			}

			slices.SortFunc(connectorProperties, func(i, j connectorDocPropertyData) int {
				return strings.Compare(i.Name, j.Name)
			})
			connectorList[i].Properties = connectorProperties
		}
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
		if unicode.IsUpper(r) {
			if i > 0 && camel[i-1] != '_' && !unicode.IsUpper(rune(camel[i-1])) {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else {
			buf.WriteRune('_')
		}
	}

	// Return the contents of the buffer as a string
	return buf.String()
}

func sanitizeForTemplate(input string) string {
	// Replaces "{{" with "{{ "{{" }}" and "}}" with "{{ "}}" }}"
	// Use placeholders to avoid double replacement
	output := strings.ReplaceAll(input, "{{", "___OPEN_BRACE___")
	output = strings.ReplaceAll(output, "}}", "___CLOSE_BRACE___")

	output = strings.ReplaceAll(output, "___OPEN_BRACE___", "{{ \"{{\" }}")
	output = strings.ReplaceAll(output, "___CLOSE_BRACE___", "{{ \"}}\" }}")
	return output
}
