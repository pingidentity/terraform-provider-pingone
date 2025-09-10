# Provider Design

This document provides an architectural and design overview of the PingOne Terraform Provider, outlining its structure, design principles, and development patterns.

## Overview

The PingOne Terraform Provider is a comprehensive infrastructure-as-code solution for managing PingOne identity and access management resources. Built using HashiCorp's Terraform Plugin Framework and SDK, it provides full lifecycle management of PingOne configuration through a well-structured, maintainable codebase.

## Architecture

### High-Level Architecture

The provider follows a layered architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     Terraform Core                          │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│              Provider Entry Point (main.go)                 │
│  - Version management                                       │
│  - Server factory initialization                           │
│  - Debug mode support                                      │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│           Provider Factory (internal/provider)              │
│  - Dual Framework Support (v5/v6 mux)                     │
│  - Framework Provider (Plugin Framework)                   │
│  - SDKv2 Provider (Legacy SDK)                            │
└─────────────────────────────────────────────────────────────┘
                                │
                    ┌───────────┴───────────┐
                    ▼                       ▼
┌─────────────────────────────┐   ┌─────────────────────────────┐
│    Framework Provider       │   │     SDKv2 Provider          │
│  (Plugin Framework v1.x)    │   │   (Legacy SDK v2.x)         │
│  - Modern resource impl.    │   │   - Legacy resources        │
│  - Type-safe schemas        │   │   - Gradual migration       │
│  - Better validation        │   │                             │
└─────────────────────────────┘   └─────────────────────────────┘
                    │                       │
                    └───────────┬───────────┘
                                ▼
┌─────────────────────────────────────────────────────────────┐
│               Service Layer (internal/service)              │
│  ├── authorize/    - PingOne Authorize resources           │
│  ├── base/         - Core/foundational resources           │
│  ├── credentials/  - PingOne Credentials resources         │
│  ├── mfa/          - Multi-factor authentication           │
│  ├── risk/         - PingOne Risk Management               │
│  ├── sso/          - Single Sign-On resources             │
│  └── verify/       - PingOne Verify resources             │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│           Framework & Utilities (internal/framework)        │
│  - Type conversion utilities                               │
│  - Custom validators                                       │
│  - Plan modifiers                                         │
│  - Common resource patterns                               │
│  - Schema helpers                                         │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│              Client Layer (internal/client)                 │
│  - Configuration management                                │
│  - Authentication handling                                 │
│  - API client initialization                              │
│  - Global options support                                 │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│            PingOne Go SDK                                  │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                    PingOne APIs                             │
│  - Management API                                          │
│  - Service-specific APIs                                   │
│  - Multi-region support                                    │
└─────────────────────────────────────────────────────────────┘
```

### Core Components

#### 1. Provider Entry Point (`main.go`)
- **Purpose**: Application entry point and server initialization
- **Key Features**:
  - Version management for releases
  - Debug mode support for development
  - Terraform Plugin Protocol v6 server setup
  - Provider server factory delegation

#### 2. Provider Factory (`internal/provider/factory.go`)
- **Purpose**: Manages dual framework support and provider multiplexing
- **Key Features**:
  - **Dual Framework Support**: Seamlessly combines Plugin Framework v1.x and SDK v2.x providers
  - **Protocol Upgrade**: Automatically upgrades SDK v2 (protocol v5) to protocol v6
  - **Multiplexing**: Routes requests to appropriate provider implementation
  - **Migration Strategy**: Enables gradual migration from SDK v2 to Plugin Framework

#### 3. Framework Provider (`internal/provider/framework/`)
- **Purpose**: Modern provider implementation using Terraform Plugin Framework
- **Key Features**:
  - Type-safe schema definitions
  - Enhanced validation and plan modification
  - Improved error handling and diagnostics
  - Resource and data source registration
  - Provider configuration management

#### 4. Service Layer (`internal/service/`)
Organized by PingOne service domains:

- **`authorize/`**: PingOne Authorize (fine-grained authorization)
- **`base/`**: Core foundational resources (environments, organizations, licenses)
- **`credentials/`**: PingOne Credentials (digital credentials and wallets)
- **`mfa/`**: Multi-factor authentication resources
- **`risk/`**: PingOne Risk Management (risk assessment and policies)
- **`sso/`**: Single Sign-On resources (applications, users, groups, policies)
- **`verify/`**: PingOne Verify (identity verification and proofing)

Each service module follows a consistent pattern:
```
service/
├── service.go                    # Resource/DataSource registration
├── resource_*.go                 # Resource implementations
├── data_source_*.go             # Data source implementations
├── *_test.go                    # Acceptance tests
├── utils_*.go                   # Service-specific utilities
└── sweep.go                     # Test cleanup utilities
```

**Note**: This document describes the patterns for manually implemented resources and data sources. Considerations and patterns specific to generated code implementations are not covered in this design document.

#### 5. Framework Utilities (`internal/framework/`)
Provides reusable components and patterns:

- **Type Conversion**: Utilities for converting between Terraform types and Go SDK types
- **Custom Validators**: Domain-specific validation logic
- **Plan Modifiers**: Custom plan modification behaviors
- **Schema Helpers**: Common schema patterns and attribute definitions
- **Custom Types**: PingOne-specific type implementations (ResourceID, etc.)

#### 6. Client Layer (`internal/client/`)
- **Configuration Management**: Provider configuration parsing and validation
- **Authentication**: OAuth2 client credentials and access token management
- **API Client Initialization**: PingOne SDK client setup with proper configuration
- **Global Options**: Provider-wide behavioral configurations
- **Multi-Region Support**: Automatic endpoint selection based on region

### Design Principles

#### 1. **Separation of Concerns**
- **Service Isolation**: Each PingOne service has its own module with clear boundaries
- **Layer Separation**: Distinct layers for provider logic, business logic, and API interaction
- **Single Responsibility**: Each component has a clearly defined purpose

#### 2. **Dual Framework Strategy**
- **Progressive Migration**: Gradual transition from SDK v2 to Plugin Framework
- **Backward Compatibility**: Existing resources continue to work during migration
- **Protocol Unification**: Both frameworks present a unified interface via protocol v6

#### 3. **Type Safety**
- **Framework Types**: Leverages Plugin Framework's type-safe attribute system
- **Custom Types**: Implements domain-specific types (ResourceID, enum types)
- **Validation**: Comprehensive validation at multiple layers

#### 4. **Consistency**
- **Naming Conventions**: Standardized resource and data source naming
- **File Organization**: Consistent directory and file structure across services
- **Code Patterns**: Reusable patterns for common operations

#### 5. **Testability**
- **Acceptance Tests**: Comprehensive integration testing against real PingOne APIs
- **Test Utilities**: Shared testing infrastructure and helpers
- **Test Isolation**: Independent test environments and cleanup procedures

#### 6. **Code Generation Strategy**
- **Migration to Generated Code**: The provider is transitioning from manually created and maintained resources and data sources to using generated code over time
- **Generated Code Identification**: Generated code files are marked with a comment at the top of the file indicating they are generated
- **Implementation Approach**: New resources and data sources will prioritize generated implementations, while existing manual implementations will be migrated incrementally
- **Maintenance Benefits**: Generated code reduces maintenance overhead, ensures consistency across resources, and accelerates development of new features
- **Quality Assurance**: Generated code follows the same testing and validation standards as manually written code

## Framework Utilities Reference

The `internal/framework/` package provides a comprehensive set of utilities and patterns for implementing consistent, well-documented, and properly validated resources and data sources. These utilities should be used when defining schemas, descriptions, and validators to ensure uniformity across the provider.

### Schema Description Utilities

The `internal/framework/` package provides schema description utilities that standardize how resource and data source documentation is generated. These utilities ensure consistent formatting between plain text and Markdown descriptions, automatically handle common documentation patterns like allowed values and default values, and provide chainable methods for building comprehensive attribute documentation.

**Why use these utilities:**
- Ensures consistent documentation formatting across all resources
- Automatically generates both plain text and Markdown descriptions
- Provides standardized messaging for common validation scenarios
- Supports documentation of complex validation relationships between attributes
- Enables automated generation of provider documentation

### Common Schema Attributes

The framework provides pre-built schema attributes for common patterns used throughout the provider. These attributes include standard ID fields with appropriate custom types, link attributes for foreign key references with built-in validation, data source return attributes for computed results, and filter attributes for SCIM and data filtering operations.

**Why use common attributes:**
- Ensures consistency in schema definitions across resources
- Provides pre-configured validation and plan modifiers
- Reduces boilerplate code in resource implementations
- Standardizes attribute naming and behavior patterns
- Includes appropriate custom types for type safety

### Custom Validators

The framework includes specialized validators organized by data type (`stringvalidator`, `boolvalidator`, `int32validator`, etc.) that handle PingOne-specific validation scenarios. These validators support conditional validation based on other field values, content validation for specific formats, cross-field validation for conflicting or mutually exclusive attributes, and domain-specific validation rules.

**Why use framework validators:**
- Provides consistent validation behavior across resources
- Handles complex validation scenarios specific to PingOne
- Includes proper error messaging for validation failures
- Supports conditional validation based on resource state
- Reduces code duplication in validation logic

### Custom Plan Modifiers

The framework provides plan modifiers that handle PingOne-specific lifecycle management scenarios, including data loss protection for immutable fields, conditional plan modification based on resource state, and specialized replacement behaviors for sensitive operations.

**Why use framework plan modifiers:**
- Ensures consistent lifecycle management across resources
- Provides appropriate warnings for potentially destructive operations
- Handles edge cases specific to PingOne resource management
- Standardizes immutability patterns across the provider
- Includes proper user messaging for plan modifications

### Custom Types

The framework includes custom types specifically designed for PingOne resources, including `ResourceIDType` for PingOne UUID validation and `DaVinci` types for DaVinci service-specific requirements. These types provide enhanced validation, consistent serialization/deserialization, and improved type safety throughout the provider.

**Why use custom types:**
- Provides compile-time type safety for PingOne resource identifiers
- Includes built-in validation for PingOne-specific formats
- Ensures consistent handling of resource references
- Enables better error messages for type-related issues
- Supports automatic conversion between Terraform and SDK types

### Implementation Guidelines

When implementing resources and data sources, developers should leverage the framework utilities to ensure consistency and reduce implementation overhead. The utilities handle common patterns, provide standardized validation, and ensure proper documentation generation. Refer to existing implementations in `internal/service/` for examples of proper usage patterns.

## SDK Utilities Reference

The `internal/sdk` package provides essential utilities for making API calls to PingOne services with consistent error handling, retry logic, and response parsing. These utilities abstract the complexity of error handling across multiple PingOne service APIs and provide standardized patterns for CRUD operations.

**Why use SDK utilities:**
- Ensures consistent error handling across all PingOne service APIs
- Provides automatic retry logic for transient failures and permission propagation delays
- Standardizes error message formatting for better user experience
- Handles multi-service API differences transparently
- Includes built-in timeout and retry configuration
- Supports custom error handling for domain-specific business logic

### Core Functions

The package provides `ParseResponse` and `ParseResponseWithCustomTimeout` functions that serve as the primary interface for executing SDK API calls. These functions handle error parsing, retry logic, response validation, and diagnostic generation uniformly across all PingOne services.

### Error Handling

The SDK utilities include built-in error handlers for common scenarios such as resource not found warnings (useful for read operations where external deletion should not fail) and invalid value errors with enhanced messaging. Custom error handlers can be created for domain-specific business logic requirements.

### Retry Logic

The package provides configurable retry conditions including default retry behavior for create/read operations that may encounter permission propagation delays. Custom retry logic can be implemented for handling specific race conditions or service-specific transient failures.

### Multi-Service Support

The SDK utilities automatically handle different PingOne service APIs (Management, MFA, Authorize, Credentials, Risk, Verify) with uniform error marshaling and retry logic, eliminating the need for service-specific implementations.

### Integration Patterns

Developers should use the SDK utilities for all API interactions following established CRUD operation patterns. The utilities integrate with the framework's environment validation and provide consistent diagnostic handling across all resource types.

## Validation Utilities Reference

The `internal/verify` package provides validation functions, regular expressions, and validators for common data formats and PingOne-specific requirements. These utilities ensure consistent validation across all resources and data sources and provide a centralized location for validation logic.

**Why use validation utilities:**
- Ensures consistent validation behavior across all resources
- Provides pre-defined regular expressions for common data formats
- Includes PingOne-specific validation for resource identifiers and formats
- Supports both SDKv2 and Plugin Framework validation patterns
- Centralizes validation logic for maintainability
- Includes comprehensive locale and language validation support

### Regular Expressions

The package includes pre-defined regular expressions for resource identifiers (PingOne UUIDs and DaVinci IDs), network and protocol formats (IPv4, IPv6, URLs, domains), and other common formats (timestamps, color codes, country codes). These patterns are available in both partial match and full string match variants.

### Validation Functions

The package provides validation functions compatible with both SDKv2 and Plugin Framework implementations. These include resource ID validators for PingOne and DaVinci resources, and framework validators that return appropriate validator types for Plugin Framework usage.

### Locale and Language Utilities

The package includes comprehensive ISO language code support with functions for generating language validators, retrieving complete ISO language lists, and managing reserved language codes specific to PingOne.

### OIDC Attribute Validation

Specialized utilities for OpenID Connect attribute validation include functions for retrieving lists of illegal and overrideable OIDC attribute names, along with formatted strings for documentation purposes.

### Integration Patterns

Developers should use the validation utilities consistently across Framework and SDKv2 resources, apply appropriate validators in acceptance tests using the provided regex patterns, and leverage the centralized validation logic rather than creating custom patterns.

## Utility Functions Reference

The `internal/utils` package provides helper functions for common data transformations, type conversions, and utility operations used throughout the provider. These utilities standardize common operations and reduce code duplication across implementations.

**Why use utility functions:**
- Standardizes common data transformation operations
- Provides type-safe conversions between different data types
- Includes specialized functions for working with SDK enum values
- Offers cryptographically secure random generation for testing
- Supports schema composition and reusability patterns
- Reduces code duplication across resource implementations

### Enum Utilities

The package provides functions for converting SDK enum values to formats compatible with Terraform validators and framework descriptions. These utilities handle the complexity of enum serialization and provide consistent interfaces for working with SDK-provided enum values.

### Type Conversion Utilities

Conversion functions are available for transforming slices between different types (string, int, int32, int64) to the `any` type required by various Terraform framework interfaces. These utilities provide type-safe conversions for use in validators, descriptions, and other framework operations.

### JSON Utilities

The package includes functions for performing semantic comparison of JSON content, which is useful for state comparison where formatting differences should not trigger changes. These utilities handle JSON unmarshaling and deep comparison automatically.

### String Utilities

Cryptographically secure random string generation functions are provided for test data creation and other scenarios requiring secure random values. The utilities support custom character sets and various random value generation patterns.

### Schema Utilities

Functions for merging schema attribute maps enable schema composition and reusability patterns. These utilities support optional overwrite behavior and facilitate the creation of base schemas that can be extended with additional attributes.

### Integration Patterns

Developers should leverage these utilities for consistent data transformations, use enum utilities when working with SDK-provided enum values in validators and descriptions, apply random generation functions for secure test data creation, and utilize schema merging for creating reusable schema components.

## Provider Documentation

The PingOne Terraform Provider uses an automated documentation generation system that creates comprehensive documentation from multiple sources. Understanding the relationship between templates, examples, and generated documentation is essential for maintaining high-quality provider documentation.

### Documentation Architecture

The provider documentation system consists of three main components that work together to generate the final documentation published at `/docs`:

#### Templates (`/templates`)
Template files define the structure and content for documentation pages. These templates use Go template syntax and are organized by type and service subcategory:

- **Structure**: `/templates/{type}/{service-subcategory}/`
- **Types**: `data-sources/`, `resources/`, `guides/`
- **Service Subcategories**: Each subdirectory corresponds to a service package in `internal/service/`
  - `authorize/` - PingOne Authorize resources
  - `base/` - Core foundational resources
  - `credentials/` - PingOne Credentials resources
  - `mfa/` - Multi-factor authentication resources
  - `risk/` - PingOne Risk Management resources
  - `sso/` - Single Sign-On resources
  - `verify/` - PingOne Verify resources

#### Examples (`/examples`)
Example configurations demonstrate real-world usage patterns and are referenced by templates during documentation generation:

- **Structure**: `/examples/{type}/{service-subcategory}/`
- **Types**: `data-sources/`, `resources/`, `guides/`, `provider/`
- **Service Alignment**: Subcategories mirror the service package structure
- **Content**: Complete Terraform configuration examples showing practical usage

#### Generated Documentation (`/docs`)
The final documentation is generated by running `make generate`, which processes templates and examples to create comprehensive documentation:

- **Generation Process**: Combines template content with example configurations
- **Output Structure**: `/docs/{type}/{service-subcategory}/`
- **Content**: Complete documentation pages with descriptions, schemas, and examples

### Documentation Generation Process

#### Make Target
```bash
make generate
```

This command processes all templates and examples to generate the complete documentation set in `/docs`. The generation process:

1. **Reads Templates**: Processes Go template files from `/templates`
2. **Incorporates Examples**: Includes relevant example configurations from `/examples`
3. **Generates Schema Documentation**: Extracts schema information from resource and data source implementations
4. **Creates Final Documentation**: Outputs complete documentation pages to `/docs`

#### Template Processing
Templates use Go template syntax to include dynamic content:
- **Schema Information**: Automatically extracted from resource implementations
- **Example Configurations**: Referenced from corresponding example files
- **Cross-References**: Links to related resources and data sources
- **Metadata**: Provider version information and generation timestamps

### Service Package Alignment

The documentation structure directly mirrors the service package organization:

#### Service Package → Documentation Mapping
- `internal/service/authorize/` → `templates/*/authorize/` → `examples/*/authorize/` → `docs/*/authorize/`
- `internal/service/base/` → `templates/*/base/` → `examples/*/base/` → `docs/*/base/`
- `internal/service/credentials/` → `templates/*/credentials/` → `examples/*/credentials/` → `docs/*/credentials/`
- `internal/service/mfa/` → `templates/*/mfa/` → `examples/*/mfa/` → `docs/*/mfa/`
- `internal/service/risk/` → `templates/*/risk/` → `examples/*/risk/` → `docs/*/risk/`
- `internal/service/sso/` → `templates/*/sso/` → `examples/*/sso/` → `docs/*/sso/`
- `internal/service/verify/` → `templates/*/verify/` → `examples/*/verify/` → `docs/*/verify/`

This alignment ensures that documentation organization matches the provider's internal structure, making it easier for developers to locate relevant documentation and examples.

