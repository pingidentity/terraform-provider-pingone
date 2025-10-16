# Acceptance Testing Strategy

This document outlines the comprehensive testing strategy for the terraform-provider-pingone, based on analysis of existing test patterns and Terraform provider best practices. All new resources and data sources should follow these testing patterns to ensure consistency, reliability, and maintainability.

## Overview

The provider uses HashiCorp's terraform-plugin-testing framework with acceptance tests that interact with real PingOne environments. Tests are organized by service and follow consistent naming and structure patterns.

## Test Organization

### File Structure
- **Resource tests**: `resource_<name>_test.go`
- **Data source tests**: `data_source_<name>_test.go` 
- **Test helpers**: Located in `internal/acctest/service/<service>/`

### Package Structure
- Tests are placed in `<service>_test` packages (e.g., `sso_test`)
- Import test helpers from `internal/acctest/service/<service>`
- Use shared environment configurations from `internal/acctest`

## Core Testing Patterns

### 1. Standard Test Functions

Every resource should implement these core test functions:

#### **Removal Drift Tests**
Tests resource behavior when the underlying PingOne resource is deleted outside Terraform:
```go
func TestAcc<Resource>_RemovalDrift(t *testing.T) {
    // Test resource removal detection
    // Test environment removal detection
}
```

#### **New Environment Tests**
Tests resource creation in a freshly provisioned environment:
```go
func TestAcc<Resource>_NewEnv(t *testing.T) {
    // Test resource creation in new environment
    // Important for root-level resources (populations, etc.)
}
```

#### **Full Schema Tests**
Tests complete resource lifecycle with minimal and maximal configurations:
```go
func TestAcc<Resource>_Full(t *testing.T) {
    // Test minimal → maximal → minimal transitions
    // Test import functionality
    // Test all schema attributes
}
```

### 2. Schema Testing Strategy

#### **Minimal Schema Testing**
- Test with only required fields set
- Validate API defaults are properly handled
- Use `TestCheckNoResourceAttr` for optional fields that should be absent

```go
minimalCheck := resource.ComposeTestCheckFunc(
    resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
    resource.TestCheckResourceAttr(resourceFullName, "name", name),
    resource.TestCheckNoResourceAttr(resourceFullName, "description"),
    resource.TestCheckNoResourceAttr(resourceFullName, "optional_field"),
)
```

#### **Maximal Schema Testing**
- Test with all fields (required and optional) set
- Validate complex nested objects and arrays
- Test all supported enum values and configurations

```go
fullCheck := resource.ComposeTestCheckFunc(
    resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
    resource.TestCheckResourceAttr(resourceFullName, "name", name),
    resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
    resource.TestCheckResourceAttr(resourceFullName, "complex_field.nested_field", "expected_value"),
    resource.TestCheckResourceAttr(resourceFullName, "array_field.#", "2"),
)
```

#### **Transition Testing**
Test that resources can successfully transition between different configurations:
- Minimal → Maximal → Minimal
- Different enum values
- Adding/removing optional blocks
- Updating immutable fields (should force replacement)

#### **Backward Compatibility Testing**
When deprecating schema fields, implement comprehensive backward compatibility tests:
- **Dual Functionality**: Test that both deprecated and new fields work simultaneously
- **Gradual Migration**: Validate that users can migrate from old to new fields incrementally
- **Deprecation Warnings**: Ensure deprecated fields generate appropriate warnings
- **Legacy Support**: Test that existing configurations continue to work unchanged

```go
// Example: Testing deprecated field alongside new field
func TestAcc<Resource>_BackwardCompatibility_DeprecatedField(t *testing.T) {
    // Test old field only (legacy behavior)
    // Test new field only (current behavior)  
    // Test both fields together (migration period)
    // Test transition from old to new field
}
```

### 3. Pre-Check Functions

The `internal/acctest` package provides comprehensive pre-check functions to validate test requirements before execution. These functions ensure tests only run when appropriate dependencies and environment variables are configured.

#### **Essential Pre-Checks**

**Basic Client Authentication**
```go
acctest.PreCheckClient(t)
```
- **Purpose**: Validates core PingOne client credentials
- **Required Environment Variables**: 
  - `PINGONE_CLIENT_ID`
  - `PINGONE_CLIENT_SECRET` 
  - `PINGONE_ENVIRONMENT_ID`
  - `PINGONE_REGION_CODE`
- **When to Use**: Every acceptance test

**Test Stability Control**
```go
acctest.PreCheckNoTestAccFlaky(t)  // Standard tests
acctest.PreCheckTestAccFlaky(t)    // Flaky test scenarios
```
- **Purpose**: Control execution of potentially unstable tests
- **Environment Variable**: `TESTACC_FLAKY=true`
- **When to Use**: 
  - `PreCheckNoTestAccFlaky`: Most standard tests (skips when flaky flag is set)
  - `PreCheckTestAccFlaky`: Tests known to be flaky or under development

#### **Environment and Infrastructure Pre-Checks**

**New Environment Creation**
```go
acctest.PreCheckNewEnvironment(t)
```
- **Purpose**: Validates ability to create new PingOne environments
- **Required Environment Variables**: `PINGONE_LICENSE_ID`
- **When to Use**: Tests that create ephemeral environments

**Organization Access**
```go
acctest.PreCheckOrganisationName(t)   // Name-based access
acctest.PreCheckOrganisationID(t)     // ID-based access
```
- **Required Environment Variables**: 
  - `PINGONE_ORGANIZATION_NAME`
  - `PINGONE_ORGANIZATION_ID`
- **When to Use**: Tests requiring organization-level operations

**Regional Restrictions**
```go
acctest.PreCheckRegionSupportsWorkforce(t)
```
- **Purpose**: Skips tests in regions where workforce environments aren't supported
- **Automatic Detection**: Based on `PINGONE_REGION_CODE` (skips CA, SG)
- **When to Use**: Tests using workforce-specific features

#### **Feature Flag Pre-Checks**

**Standard Feature Flag Control**
```go
acctest.PreCheckNoFeatureFlag(t)     // Standard tests
acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)  // Feature-specific tests
```
- **Purpose**: Control test execution based on feature flags
- **Environment Variable**: `FEATURE_FLAG`
- **Available Flags**: `DAVINCI`
- **When to Use**:
  - `PreCheckNoFeatureFlag`: Standard tests that should run without feature flags
  - `PreCheckFeatureFlag`: Tests requiring specific feature enablement

#### **Certificate and Key Management Pre-Checks**

**PKCS#12 Certificates**
```go
acctest.PreCheckPKCS12Key(t)                    // Encrypted PKCS#12
acctest.PreCheckPKCS12UnencryptedKey(t)         // Unencrypted PKCS#12
```
- **Required Environment Variables**:
  - `PINGONE_KEY_PKCS12`, `PINGONE_KEY_PKCS12_PASSWORD` (encrypted)
  - `PINGONE_KEY_PKCS12_UNENCRYPTED` (unencrypted)

**Certificate Signing and Formats**
```go
acctest.PreCheckPKCS12WithCSR(t)               // Certificate signing requests
acctest.PreCheckPKCS12CSRResponse(t)           // CSR responses
acctest.PreCheckAPNSPKCS8Key(t)                // Apple Push Notification keys
acctest.PreCheckPKCS7Cert(t)                   // PKCS#7 certificates
acctest.PreCheckPEMCert(t)                     // PEM format certificates
```
- **When to Use**: Tests involving certificate management, signing, or specific formats

#### **Third-Party Integration Pre-Checks**

**Google/Firebase Integration**
```go
acctest.PreCheckGoogleJSONKey(t)               // Google Play integrity
acctest.PreCheckGoogleFirebaseCredentials(t)   // Firebase messaging
```
- **Required Environment Variables**: 
  - `PINGONE_GOOGLE_JSON_KEY`
  - `PINGONE_GOOGLE_FIREBASE_CREDENTIALS`

**SMS Provider Integration**
```go
acctest.PreCheckTwilio(t)                      // Twilio SMS
acctest.PreCheckSyniverse(t)                   // Syniverse SMS
```
- **Twilio Variables**: `PINGONE_TWILIO_SID`, `PINGONE_TWILIO_AUTH_TOKEN`, `PINGONE_TWILIO_NUMBER`
- **Syniverse Variables**: `PINGONE_SYNIVERSE_AUTH_TOKEN`, `PINGONE_SYNIVERSE_NUMBER`
- **Optional Skip**: `PINGONE_TWILIO_TEST_SKIP=true`, `PINGONE_SYNIVERSE_TEST_SKIP=true`

**Custom Domain SSL**
```go
acctest.PreCheckCustomDomainSSL(t)
```
- **Required Environment Variables**: 
  - `PINGONE_DOMAIN_CERTIFICATE_PEM`
  - `PINGONE_DOMAIN_INTERMEDIATE_CERTIFICATE_PEM`
  - `PINGONE_DOMAIN_KEY_PEM`

#### **Pre-Check Usage Patterns**

**Standard Test Pattern**
```go
PreCheck: func() {
    acctest.PreCheckNoTestAccFlaky(t)
    acctest.PreCheckClient(t)
    acctest.PreCheckNoFeatureFlag(t)
},
```

**New Environment Test Pattern**
```go
PreCheck: func() {
    acctest.PreCheckNoTestAccFlaky(t)
    acctest.PreCheckClient(t)
    acctest.PreCheckNewEnvironment(t)
    acctest.PreCheckNoFeatureFlag(t)
},
```

**Integration Test Pattern**
```go
PreCheck: func() {
    acctest.PreCheckNoTestAccFlaky(t)
    acctest.PreCheckClient(t)
    acctest.PreCheckGoogleJSONKey(t)          // Third-party integration
    acctest.PreCheckNoFeatureFlag(t)
},
```

**Feature-Specific Test Pattern**
```go
PreCheck: func() {
    acctest.PreCheckNoTestAccFlaky(t)
    acctest.PreCheckClient(t)
    acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)
},
```

### 4. Environment Strategy

#### **Shared Environment Usage**
Most tests should use the shared environment for efficiency:
```go
PreCheck: func() {
    acctest.PreCheckNoTestAccFlaky(t)
    acctest.PreCheckClient(t)
    acctest.PreCheckNoFeatureFlag(t)
},
```

Use `acctest.GenericSandboxEnvironment()` in test configurations.

#### **Ephemeral Environment Usage**
Use dedicated environments for tests that:
- Modify shared/global configuration (e.g., enabling languages)
- Test root-level resources (populations) to validate race conditions
- Require specific environment settings or feature flags

```go
PreCheck: func() {
    acctest.PreCheckNoTestAccFlaky(t)
    acctest.PreCheckClient(t)
    acctest.PreCheckNewEnvironment(t)
    acctest.PreCheckNoFeatureFlag(t)
},
```

### 5. Import Testing

All resources must test import functionality:
```go
{
    ResourceName: resourceFullName,
    ImportStateIdFunc: func() resource.ImportStateIdFunc {
        return func(s *terraform.State) (string, error) {
            rs, ok := s.RootModule().Resources[resourceFullName]
            if !ok {
                return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
            }
            return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
        }
    }(),
    ImportState:       true,
    ImportStateVerify: true,
    ImportStateVerifyIgnore: []string{
        "sensitive_field",
        "computed_field",
    },
}
```

### 6. Error Testing

#### **API Error Validation**
Test expected API error conditions:
```go
{
    Config:      testAccConfig_InvalidConfiguration(resourceName),
    ExpectError: regexp.MustCompile("Expected error message pattern"),
}
```

#### **Import Error Testing**
Test invalid import scenarios:
```go
{
    ImportState:   true,
    ImportStateId: "invalid-id-format",
    ExpectError:   regexp.MustCompile("Unexpected Import Identifier"),
}
```

#### **Configuration Validation**
Test Terraform configuration validation:
```go
{
    Config:      testAccConfig_ConflictingAttributes(resourceName),
    ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
}
```

## Data Source Testing

### Standard Data Source Tests

#### **By Name Lookup**
```go
func TestAcc<DataSource>DataSource_ByNameFull(t *testing.T) {
    // Test data source lookup by name
    // Use TestCheckResourceAttrPair to compare with resource
}
```

#### **By ID Lookup**
```go
func TestAcc<DataSource>DataSource_ByIDFull(t *testing.T) {
    // Test data source lookup by ID
    // Use TestCheckResourceAttrPair to compare with resource
}
```

#### **Not Found Scenarios**
```go
func TestAcc<DataSource>DataSource_NotFound(t *testing.T) {
    // Test data source behavior with non-existent resources
    // Should return appropriate error
}
```

## Advanced Testing Patterns

### 1. Complex Resource Types

For resources with multiple types/variants (e.g., applications with different OIDC types):
- Test each type's minimal and maximal schema separately
- Test transitions between compatible types
- Test error conditions for incompatible type changes

### 2. Regional Functionality

For features available only in specific regions:
- Use appropriate environment configurations
- Test both supported and unsupported regions
- Validate proper error handling for unsupported regions

### 3. Feature Flag Dependencies

For functionality behind feature flags:
- Use `acctest.PreCheckFeatureFlag(t)` when feature is required
- Use `acctest.PreCheckNoFeatureFlag(t)` for standard tests
- Test both enabled and disabled states where applicable

### 4. Immutable Field Testing

For resources with immutable fields:
- Test that changes force resource replacement
- Validate error messages for unsupported updates
- Test workarounds for resources that associate customer data

### 5. Schema Deprecation and Backward Compatibility

When deprecating schema fields or changing resource behavior:

#### **Deprecation Strategy**
- **Simultaneous Support**: Both deprecated and new functionality must work during the deprecation period
- **Gradual Migration**: Users should be able to migrate incrementally without breaking changes
- **Clear Warnings**: Deprecated fields should generate helpful deprecation warnings
- **Documentation**: Update both code comments and user documentation

#### **Required Deprecation Tests**
```go
func TestAcc<Resource>_Deprecation_<FieldName>(t *testing.T) {
    // Test deprecated field works (legacy configuration)
    // Test new field works (current configuration)
    // Test both fields work together (migration period)
    // Test migration path from deprecated to new field
    // Validate deprecation warnings are generated
}
```

#### **Backward Compatibility Validation**
- Existing user configurations must continue to work unchanged
- No functional regression during deprecation period
- Smooth migration path from old to new schema
- Proper handling of edge cases during transition

### 6. Collection Type Testing

For ambiguous collection fields:
- Test ordering for `list` types
- Test uniqueness for `set` types
- Validate proper handling of duplicates
- Test empty collections

## Test Configuration Patterns

### Configuration Helpers

Use consistent configuration helper functions:
```go
func testAcc<Resource>Config_Minimal(resourceName, name string) string {
    return fmt.Sprintf(`
%[1]s

resource "pingone_<resource>" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  // Only required fields
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAcc<Resource>Config_Full(resourceName, name string) string {
    return fmt.Sprintf(`
%[1]s

resource "pingone_<resource>" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test description"
  // All fields including optional
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
```

### Test Helper Functions

Implement standard helper functions in `internal/acctest/service/<service>/`:
```go
func <Resource>_CheckDestroy(s *terraform.State) error
func <Resource>_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc
func <Resource>_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, resourceID string)
```

## Quality Standards

### Test Coverage Requirements

All resources and data sources must have:
- [ ] **Removal drift detection**
- [ ] **Minimal schema validation** with API defaults
- [ ] **Maximal schema validation** with all optional fields
- [ ] **Schema transition testing** (minimal ↔ maximal)
- [ ] **Import functionality testing**
- [ ] **Error condition testing**
- [ ] **Data source lookup testing** (if applicable)
- [ ] **Backward compatibility testing** (when deprecating fields)

#### **Additional Requirements for Schema Changes**
When modifying existing resources:
- [ ] **Deprecation testing** for any removed or changed fields
- [ ] **Migration path validation** from old to new schema
- [ ] **Dual functionality testing** during deprecation periods
- [ ] **Warning validation** for deprecated field usage

### Test Reliability

- Use `t.Parallel()` for parallel execution
- Implement proper cleanup in CheckDestroy functions
- Handle race conditions for environment-level resources
- Use appropriate timeouts for resource operations
- Validate test stability across multiple runs

### Maintenance Considerations

- Keep test configurations DRY with helper functions
- Use descriptive test and variable names
- Document complex test scenarios
- Regular review and update of test patterns
- Ensure tests remain valid as APIs evolve

## Examples

### Complete Resource Test Structure
```go
func TestAcc<Resource>_RemovalDrift(t *testing.T) { /* ... */ }
func TestAcc<Resource>_NewEnv(t *testing.T) { /* ... */ }
func TestAcc<Resource>_Full(t *testing.T) {
    // Minimal test step
    // Maximal test step  
    // Transition tests
    // Import test
    // Error tests
}
func TestAcc<Resource>_<SpecificScenario>(t *testing.T) { /* ... */ }

// For resources with deprecated fields
func TestAcc<Resource>_Deprecation_<FieldName>(t *testing.T) {
    // Legacy configuration test
    // New configuration test
    // Migration path test
    // Dual functionality test
}
```

### Complete Data Source Test Structure
```go
func TestAcc<DataSource>DataSource_ByNameFull(t *testing.T) { /* ... */ }
func TestAcc<DataSource>DataSource_ByIDFull(t *testing.T) { /* ... */ }
func TestAcc<DataSource>DataSource_NotFound(t *testing.T) { /* ... */ }
```

This comprehensive testing strategy ensures that all provider functionality is thoroughly validated, providing confidence in the reliability and correctness of the terraform-provider-pingone.