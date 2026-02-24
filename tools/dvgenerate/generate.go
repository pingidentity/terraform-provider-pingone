//go:generate go run ./cmd/generate/generate.go

package dvgenerate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"github.com/pingidentity/terraform-provider-pingone/dvgenerate/internal"
	"github.com/samir-gandhi/davinci-client-go/davinci"
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
	dir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	// Go up two levels to the provider root directory (tools/dvgenerate -> tools -> root)
	dir = filepath.Dir(filepath.Dir(dir))
	fmt.Println("base directory:", dir)

	GenerateReferenceTemplate(dir)
	GenerateConnectorHCLExamples(dir)
}

func GenerateReferenceTemplate(baseDirectory string) {

	fileNameDirectory := fmt.Sprintf("%s/templates/guides", baseDirectory)

	conns, err := readConnectors()
	if err != nil {
		panic(err)
	}

	err = writeConnectorTemplate(fileNameDirectory, internal.ConnectorReferenceTmpl, conns, true)
	if err != nil {
		panic(err)
	}
}

func GenerateConnectorHCLExamples(baseDirectory string) {

	conns, err := readConnectors()
	if err != nil {
		panic(err)
	}

	fileNameDirectory := fmt.Sprintf("%s/examples/davinci-connector-instances", baseDirectory)
	fmt.Println("connector examples directory:", fileNameDirectory)
	for _, conn := range conns {
		err := writeConnector(fileNameDirectory, internal.ConnectorTmpl, conn, false)
		if err != nil {
			panic(err)
		}
	}
}

func writeConnectorTemplate(fileNameDirectory, templateString string, conns connectionByName, overwrite bool) error {
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

func readConnectors() (connectionByName, error) {

	connectorSchema := []davinci.Connector{}
	err := davinci.Unmarshal(internal.ConnectorSchemaBytes, &connectorSchema, davinci.ExportCmpOpts{
		IgnoreEnvironmentMetadata: true,
	})
	if err != nil {
		return nil, err
	}

	connectorList := make(connectionByName, 0)
	for _, connectorSchemaItem := range connectorSchema {

		connectorDocData := connectorDocData{}

		if v := connectorSchemaItem.ConnectorID; v != nil {
			connectorDocData.ConnectorId = *v
		} else {
			connectorDocData.ConnectorId = "No value"
		}

		if v := connectorSchemaItem.Name; v != nil {
			connectorDocData.ConnectorName = *v
		} else {
			connectorDocData.ConnectorName = "No name"
		}

		if connectorSchemaItem.Properties != nil {
			connectorProperties := make(connectionPropertyByName, 0)
			for propertyName, property := range connectorSchemaItem.Properties {
				propertyType := "string"
				if v := property.Type; v != nil {
					propertyType = *v
				}

				propertyType = rewritePropertyType(propertyType)

				description := property.Info

				if description != nil && strings.TrimSpace(*description) != "" && !strings.HasSuffix(strings.TrimSpace(*description), ".") {
					descriptionTemp := fmt.Sprintf("%s.", *description)
					description = &descriptionTemp
				}

				connectorProperty := connectorDocPropertyData{
					Name:               propertyName,
					Type:               &propertyType,
					Description:        description,
					ConsoleDisplayName: property.DisplayName,
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

type connectionPropertyByName []connectorDocPropertyData

func (a connectionPropertyByName) Len() int           { return len(a) }
func (a connectionPropertyByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a connectionPropertyByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type connectionByName []connectorDocData

func (a connectionByName) Len() int { return len(a) }
func (a connectionByName) Less(i, j int) bool {
	return fmt.Sprintf("%s%s", a[i].ConnectorName, a[i].ConnectorId) < fmt.Sprintf("%s%s", a[j].ConnectorName, a[j].ConnectorId)
}
func (a connectionByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

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
