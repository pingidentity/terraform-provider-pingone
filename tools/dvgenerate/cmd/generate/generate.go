package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pingidentity/terraform-provider-pingone/dvgenerate"
)

func main() {

	var jsonFile string
	flag.StringVar(&jsonFile, "file", "", "The path to the JSON file containing the connector schema.")
	flag.Parse()

	var input []byte
	var err error

	if jsonFile != "" {
		input, err = os.ReadFile(jsonFile)
		if err != nil {
			panic(fmt.Errorf("error reading specified file %s: %w", jsonFile, err))
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				panic(fmt.Errorf("error reading stdin: %w", err))
			}
		}
	}

	if len(input) == 0 {
		defaultFile := "internal/connector-schema.json"

		if _, err := os.Stat(defaultFile); err == nil {
			input, err = os.ReadFile(defaultFile)
			if err != nil {
				panic(fmt.Errorf("error reading default file %s: %w", defaultFile, err))
			}
		}
	}

	dvgenerate.Generate(input)
}
