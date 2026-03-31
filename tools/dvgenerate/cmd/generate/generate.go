package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pingidentity/terraform-provider-pingone/dvgenerate"
)

func main() {

	var jsonFile string
	flag.StringVar(&jsonFile, "file", "", "The path to the JSON file containing the connector schema.")
	flag.Parse()

	if jsonFile == "" {
		fmt.Println("Error: The -file flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	input, err := os.ReadFile(jsonFile)
	if err != nil {
		panic(fmt.Errorf("error reading specified file %s: %w", jsonFile, err))
	}

	dvgenerate.Generate(input)
}
