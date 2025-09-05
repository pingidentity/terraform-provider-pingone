package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
	headerLines = 10
)

/*
  This script is used to verify that any beta resources or data sources in the provider have the '//go:build beta' build tag.
  It follows these steps:
  1. Finds all service_beta.go files in the provider
  2. Extracts resource and data source names from BetaResources() and BetaDataSources() functions
  3. Checks if the corresponding resource or data source file has the required build tag
  4. Reports any files that are missing the tag
*/

func main() {
	errors := 0

	fmt.Println("Checking for beta resources and data sources without the required build tag...")

	// Walk through the internal/service directory
	err := filepath.Walk("./internal/service", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, "service_beta.go") {
			serviceDir := filepath.Dir(path)
			serviceName := filepath.Base(serviceDir)
			fmt.Printf("\nChecking %s service (%s)...\n", serviceName, path)

			// Parse the service_beta.go file
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				fmt.Printf("  Error parsing %s: %v\n", path, err)
				return nil
			}

			betaResources := []string{}
			betaDataSources := []string{}

			// Find the BetaResources and BetaDataSources functions
			for _, decl := range node.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok {
					if funcDecl.Name.Name == "BetaResources" {
						if funcDecl.Body != nil {
							betaResources = extractResourceNames(funcDecl.Body, "Resource")
						}
					}

					if funcDecl.Name.Name == "BetaDataSources" {
						if funcDecl.Body != nil {
							betaDataSources = extractResourceNames(funcDecl.Body, "DataSource")
						}
					}
				}
			}

			// Check resources for the beta build tag
			for _, resource := range betaResources {
				snakeCase := camelToSnake(resource)
				resourceFile := filepath.Join(serviceDir, fmt.Sprintf("resource_%s.go", snakeCase))

				if _, err := os.Stat(resourceFile); os.IsNotExist(err) {
					fmt.Printf("  %sError: Could not find resource file for %s at %s%s\n",
						colorRed, resource, resourceFile, colorReset)
					errors++
					continue
				}

				// Check if the file has the beta build tag
				if !hasBetaBuildTag(resourceFile) {
					fmt.Printf("  %sError: Resource file %s is missing //go:build beta tag%s\n",
						colorRed, resourceFile, colorReset)
					errors++
				} else {
					fmt.Printf("  %s✓ %s has the correct build tag%s\n",
						colorGreen, resourceFile, colorReset)
				}
			}

			// Check data sources for the beta build tag
			for _, dataSource := range betaDataSources {
				snakeCase := camelToSnake(dataSource)
				dataSourceFile := filepath.Join(serviceDir, fmt.Sprintf("data_source_%s.go", snakeCase))

				if _, err := os.Stat(dataSourceFile); os.IsNotExist(err) {
					fmt.Printf("  %sError: Could not find data source file for %s at %s%s\n",
						colorRed, dataSource, dataSourceFile, colorReset)
					errors++
					continue
				}

				// Check if the file has the beta build tag
				if !hasBetaBuildTag(dataSourceFile) {
					fmt.Printf("  %sError: Data source file %s is missing //go:build beta tag%s\n",
						colorRed, dataSourceFile, colorReset)
					errors++
				} else {
					fmt.Printf("  %s✓ %s has the correct build tag%s\n",
						colorGreen, dataSourceFile, colorReset)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("%sError walking the path: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	if errors > 0 {
		fmt.Printf("\n%sFound %d error(s). Some beta resources or data sources are missing the required //go:build beta tag%s\n",
			colorRed, errors, colorReset)
		os.Exit(1)
	} else {
		fmt.Printf("\n%sAll beta resources and data sources have the correct build tag!%s\n",
			colorGreen, colorReset)
		os.Exit(0)
	}
}

func camelToSnake(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result += "_"
		}
		result += strings.ToLower(string(r))
	}
	return result
}

// hasBetaBuildTag checks if a file has the //go:build beta tag
func hasBetaBuildTag(filePath string) bool {
	cleanPath := filepath.Clean(filePath)

	// Ensure the path is within the project directory
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		fmt.Printf("  %sError resolving absolute path for %s: %v%s\n", colorRed, cleanPath, err, colorReset)
		return false
	}
	rootDir, err := filepath.Abs(".")
	if err != nil {
		fmt.Printf("  %sError resolving project root directory: %v%s\n", colorRed, err, colorReset)
		return false
	}
	if !strings.HasPrefix(absPath, rootDir) {
		fmt.Printf("  %sError: File path %s is outside project directory%s\n", colorRed, cleanPath, colorReset)
		return false
	}

	content, err := os.ReadFile(cleanPath) // #nosec G304 - path is confirmed to be within the project directory above
	if err != nil {
		fmt.Printf("  %sError reading file %s: %v%s\n", colorRed, cleanPath, err, colorReset)
		return false
	}

	lines := strings.Split(string(content), "\n")
	for i := 0; i < min(headerLines, len(lines)); i++ {
		if strings.TrimSpace(lines[i]) == "//go:build beta" {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func extractResourceNames(body *ast.BlockStmt, suffix string) []string {
	var names []string

	for _, stmt := range body.List {
		if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
			if len(returnStmt.Results) > 0 {
				// Handle different return patterns
				switch res := returnStmt.Results[0].(type) {
				case *ast.CompositeLit:
					// Direct return of array literal []func(){...}
					for _, elt := range res.Elts {
						if ident, ok := elt.(*ast.Ident); ok {
							name := ident.Name
							if strings.HasPrefix(name, "New") && strings.HasSuffix(name, suffix) {
								resourceName := name[3 : len(name)-len(suffix)] // Remove "New" and suffix
								names = append(names, resourceName)
							}
						}
					}
				case *ast.Ident:
					// Return of a variable (probably empty array)
					// Skip these as they likely don't contain anything
				}
			}
		} else if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			// Handle pattern like: resources := []func(){...}
			for _, rhs := range assignStmt.Rhs {
				if compLit, ok := rhs.(*ast.CompositeLit); ok {
					for _, elt := range compLit.Elts {
						if ident, ok := elt.(*ast.Ident); ok {
							name := ident.Name
							if strings.HasPrefix(name, "New") && strings.HasSuffix(name, suffix) {
								resourceName := name[3 : len(name)-len(suffix)] // Remove "New" and suffix
								names = append(names, resourceName)
							}
						}
					}
				}
			}
		}
	}

	return names
}
