# GitHub Copilot Instructions for PingOne Terraform Provider

## Code Documentation Standards

### Go Documentation (godoc) Best Practices

When writing or modifying Go code, ensure all documentation follows these godoc standards:

#### Package Documentation
- Add a package comment before the `package` declaration
- Start with "Package packagename" followed by a description
- Use complete sentences with proper punctuation
- Example:
  ```go
  // Package verify provides validation functions for various data types used in the PingOne Terraform provider.
  package verify
  ```

#### Function/Method Documentation
- Add comments directly above the function/method declaration
- Start the comment with the function name
- Use complete sentences in present tense
- Document comprehensively including:
  - The overall purpose of the function
  - What is returned from the function in human terms (e.g., "a string that represents the user's email address")
  - What is required from each of the inputs in human terms, providing context as to what the various parameters do for the function
  - External dependencies that can modify the logic (e.g., environment variables that must be set by the developer)
- Example:
  ```go
  // FullIsoList returns a slice of all available ISO language codes.
  // The returned slice contains string values representing standardized ISO language codes
  // such as "en" for English and "fr-CA" for French (Canada). No parameters are required.
  func FullIsoList() []string {
  ```
  
  More comprehensive example:
  ```go
  // ValidateEmailAddress checks if the provided email address is valid and accessible.
  // It returns a boolean indicating validity and an error if validation fails.
  // The email parameter must be a non-empty string containing a properly formatted email address.
  // The timeout parameter specifies the maximum duration in seconds to wait for validation.
  // This function requires the SMTP_SERVER environment variable to be set for external validation.
  func ValidateEmailAddress(email string, timeout int) (bool, error) {
  ```

#### Type Documentation
- Document all exported types (those starting with capital letters)
- Start the comment with the type name
- Explain the purpose and usage of the type
- Example:
  ```go
  // IsoCountry represents an ISO language code with its corresponding human-readable name.
  type IsoCountry struct {
  ```

#### Field Documentation
- Document exported struct fields
- Use concise but descriptive comments
- Example:
  ```go
  type IsoCountry struct {
      // Code is the ISO language code (e.g., "en", "fr-CA")
      Code string
      // Name is the human-readable language name (e.g., "English", "French (Canada)")
      Name string
  }
  ```

#### Variable/Constant Documentation
- Document exported variables and constants
- Group related constants with a single comment block when appropriate
- Example:
  ```go
  // reservedLanguageCodes contains ISO language codes that are reserved and cannot be used for custom languages.
  var reservedLanguageCodes = []string{
  ```

## Terraform Provider Specific Guidelines

### Resource and Data Source Implementation
- Use the Plugin Framework SDK for all new resources and data sources
- Register new resources in `internal/service/<service name>/service.go`
- Use the PingOne Go SDK for all API interactions (never call endpoints directly)
- Follow the established patterns for client configuration and response parsing

### Error Handling
- Use the `framework.ParseResponse` wrapper for all SDK calls
- Implement custom error handling when appropriate using `CustomError` parameter
- Use proper retry conditions with the `Retryable` parameter
- Always check for and handle diagnostics errors

### Code Organization
- Place service-specific code in `internal/service/<service name>/`
- Use common code in `internal/service/base/` for platform components shared between services
- Follow the established directory structure for consistency

### Testing
- Write acceptance tests for all new resources and data sources
- Use the common test functions in `internal/acctest`
- Ensure tests create real configuration in PingOne (acceptance tests)
- Run `make test` for unit tests and `make testacc` for acceptance tests

## General Go Best Practices

### Code Style
- Follow standard Go formatting (use `gofmt`)
- Use meaningful variable and function names
- Keep functions focused and concise
- Use proper error handling patterns
- Leverage Go modules for dependency management

### Comments and Documentation
- Use complete sentences with proper capitalization and punctuation
- Be comprehensive and descriptive, not just concise
- Explain the "why" not just the "what" when appropriate
- Document any non-obvious behavior or edge cases
- Use present tense in function descriptions ("returns", not "will return")
- For functions, always document:
  - Overall purpose and behavior
  - Return values in human-readable terms
  - Parameter requirements and their purpose
  - External dependencies (environment variables, configuration files, etc.)

### Import Organization
- Group imports logically (standard library, third-party, local)
- Use blank lines to separate import groups
- Avoid unnecessary imports

## Code Review Considerations

When suggesting code changes or improvements:
1. Ensure all exported items are properly documented
2. Verify adherence to the established provider patterns
3. Check that error handling follows the framework standards
4. Confirm new code includes appropriate tests
5. Validate that SDK usage follows the established patterns
6. Ensure code follows Go best practices and idioms

## Additional Resources

- [Contributing Guidelines](../CONTRIBUTING.md)
- [Development Environment Setup](../contributing/development-environment.md)
- [Provider Design Guide](../contributing/provider-design.md)
- [Services Support](../contributing/services-support.md)
- [PR Checklist](../contributing/pr-checklist.md)
