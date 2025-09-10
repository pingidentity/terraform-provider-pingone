# Internal Framework Components

This directory contains custom plan modifiers, validators, and types that extend the Terraform Plugin Framework for use in the PingOne provider.

## Plan Modifiers

### Bool Plan Modifiers

| Modifier | Purpose |
|----------|---------|
| `UnmodifiableDataLossProtection()` | Prevents modification of immutable boolean fields to protect against data loss. Forces manual resource replacement. |
| `UnmodifiableDataLossProtectionIf()` | Conditional version that prevents modification only when specified conditions are met. |
| `UnmodifiableDataLossProtectionIfFunc()` | Function-based conditional protection for boolean attributes. |

### Object Plan Modifiers

| Modifier | Purpose |
|----------|---------|
| `RequiresReplaceIfExistenceChanges()` | Forces resource replacement when an object attribute changes from set to unset or vice versa. |
| `ReplaceIfNowNull()` | Forces resource replacement when an object attribute becomes null. |

### Set Plan Modifiers

| Modifier | Purpose |
|----------|---------|
| `UnmodifiableDataLossProtection()` | Prevents modification of immutable set fields to protect against data loss. |
| `UnmodifiableDataLossProtectionIf()` | Conditional version that prevents modification only when specified conditions are met. |
| `UnmodifiableDataLossProtectionIfFunc()` | Function-based conditional protection for set attributes. |
| `UnmodifiableDataLossProtectionIfNull()` | Prevents modification only when the set is currently null. |
| `RequiresReplaceIfPreviouslyNull()` | Forces resource replacement when a set changes from null to non-null. |

### String Plan Modifiers

| Modifier | Purpose |
|----------|---------|
| `UnmodifiableDataLossProtection()` | Prevents modification of immutable string fields to protect against data loss. |
| `UnmodifiableDataLossProtectionIf()` | Conditional version that prevents modification only when specified conditions are met. |
| `UnmodifiableDataLossProtectionIfFunc()` | Function-based conditional protection for string attributes. |
| `ReplaceIfNowNull()` | Forces resource replacement when a string attribute becomes null. |

## Validators

### Bool Validators

| Validator | Purpose |
|-----------|---------|
| `BoolAtLeastOneOfMustBeTrue()` | Ensures at least one boolean attribute in a group is true. |
| `BoolMustBeValueIfPathSetToValue()` | Requires a boolean to have a specific value when another attribute matches a condition. |
| `BoolMustBeTrueIfPathSetToValue()` | Requires a boolean to be true when another attribute matches a condition. |
| `BoolMustBeFalseIfPathSetToValue()` | Requires a boolean to be false when another attribute matches a condition. |
| `BoolMustNotBeValueIfPathSetToValue()` | Prevents a boolean from having a specific value when another attribute matches a condition. |
| `BoolMustNotBeTrueIfPathSetToValue()` | Prevents a boolean from being true when another attribute matches a condition. |
| `BoolMustNotBeFalseIfPathSetToValue()` | Prevents a boolean from being false when another attribute matches a condition. |
| `ConflictsIfMatchesPathValue()` | Creates conflicts when another attribute matches a specific value. |

### Int32 Validators

| Validator | Purpose |
|-----------|---------|
| `IsDivisibleBy()` | Validates that an integer is exactly divisible by a specified denominator. |
| `IsGreaterThanPathValue()` | Validates that an integer is greater than the value of another attribute. |
| `IsGreaterThanEqualToPathValue()` | Validates that an integer is greater than or equal to the value of another attribute. |
| `IsLessThanPathValue()` | Validates that an integer is less than the value of another attribute. |
| `IsLessThanEqualToPathValue()` | Validates that an integer is less than or equal to the value of another attribute. |
| `RegexMatchesPathValue()` | Validates that an integer matches a regex pattern based on another attribute's value. |

### Int64 Validators

| Validator | Purpose |
|-----------|---------|
| `IsDivisibleBy()` | Validates that a 64-bit integer is exactly divisible by a specified denominator. |
| `IsGreaterThanPathValue()` | Validates that a 64-bit integer is greater than the value of another attribute. |
| `IsGreaterThanEqualToPathValue()` | Validates that a 64-bit integer is greater than or equal to the value of another attribute. |
| `IsLessThanPathValue()` | Validates that a 64-bit integer is less than the value of another attribute. |
| `IsLessThanEqualToPathValue()` | Validates that a 64-bit integer is less than or equal to the value of another attribute. |
| `RegexMatchesPathValue()` | Validates that a 64-bit integer matches a regex pattern based on another attribute's value. |

### List Validators

| Validator | Purpose |
|-----------|---------|
| `IsRequiredIfMatchesPathValue()` | Makes a list attribute required when another attribute matches a specific value. |

### Map Validators

| Validator | Purpose |
|-----------|---------|
| `ConflictsIfMatchesPathValue()` | Creates conflicts for map attributes when another attribute matches a specific value. |
| `IsRequiredIfMatchesPathValue()` | Makes a map attribute required when another attribute matches a specific value. |

### Object Validators

| Validator | Purpose |
|-----------|---------|
| `ConflictsIfMatchesPathValue()` | Creates conflicts for object attributes when another attribute matches a specific value. |
| `IsRequiredIfMatchesPathValue()` | Makes an object attribute required when another attribute matches a specific value. |

### Schema Validators

| Validator | Purpose |
|-----------|---------|
| `ConflictsIfMatchesPathValue()` | Creates conflicts at the schema level when attributes match specific values. |
| `IsRequiredIfMatchesPathValue()` | Makes attributes required at the schema level based on other attribute values. |
| `IsRequiredIfRegexMatchesPathValue()` | Makes attributes required when another attribute's value matches a regex pattern. |
| `RegexMatchesPathValue()` | Validates that schema-level values match regex patterns based on other attributes. |
| `ShouldBeDefinedValueIfPathMatchesValue()` | Suggests that an attribute should be defined when another attribute matches a value. |

### Set Validators

| Validator | Purpose |
|-----------|---------|
| `ConflictsIfMatchesPathValue()` | Creates conflicts for set attributes when another attribute matches a specific value. |
| `IsRequiredIfMatchesPathValue()` | Makes a set attribute required when another attribute matches a specific value. |

### String Validators

| Validator | Purpose |
|-----------|---------|
| `ConflictsIfMatchesPathValue()` | Creates conflicts for string attributes when another attribute matches a specific value. |
| `IsRequiredIfMatchesPathValue()` | Makes a string attribute required when another attribute matches a specific value. |
| `IsRequiredIfRegexMatchesPathValue()` | Makes a string attribute required when another attribute's value matches a regex pattern. |
| `RegexMatchesPathValue()` | Validates that a string matches a regex pattern based on another attribute's value. |
| `ShouldBeDefinedValueIfPathMatchesValue()` | Suggests that a string should be defined when another attribute matches a value. |
| `ShouldNotContain()` | Validates that a string does not contain specific substrings. |
| `StringIsBase64Encoded()` | Validates that a string contains valid base64-encoded content. |
| `IsB64ContentType()` | Validates that a string is a valid base64-encoded content with proper content type prefix. |

## Custom Types

### PingOne Types

| Type | Purpose |
|------|---------|
| `ResourceIDType` | Custom string type with PingOne resource ID validation. Ensures values conform to PingOne's resource identifier format. |
| `ResourceIDValue` | Value type for PingOne resource identifiers with built-in validation and conversion methods. |

### DaVinci Types

| Type | Purpose |
|------|---------|
| `ResourceIDType` | Custom string type with DaVinci resource ID validation for DaVinci-specific identifiers. |
| `ResourceIDValue` | Value type for DaVinci resource identifiers with built-in validation and conversion methods. |

## Internal Utilities

### String Defaults

| Utility | Purpose |
|---------|---------|
| `StaticValueUnknownable` | Provides static default values for string attributes that can handle unknown states during planning. |

### Schema Utilities

| Utility | Purpose |
|---------|---------|
| `SchemaAttributeDescription` | Utilities for creating consistent attribute descriptions across the provider. |
| `SchemaCommon` | Common schema definitions and utilities shared across resources. |

## Usage Examples

### Using Plan Modifiers

```go
// Protect an immutable field
"immutable_field": schema.StringAttribute{
    Description: "This field cannot be changed after creation",
    Required:    true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.UnmodifiableDataLossProtection(),
    },
},

// Force replacement when existence changes
"optional_object": schema.SingleNestedAttribute{
    Optional: true,
    PlanModifiers: []planmodifier.Object{
        objectplanmodifier.RequiresReplaceIfExistenceChanges(),
    },
    // ... attribute definitions
},
```

### Using Validators

```go
// Validate base64 encoding
"certificate": schema.StringAttribute{
    Description: "Base64-encoded certificate",
    Required:    true,
    Validators: []validator.String{
        stringvalidator.StringIsBase64Encoded(),
    },
},

// Conditional requirements
"password": schema.StringAttribute{
    Optional: true,
    Validators: []validator.String{
        stringvalidator.IsRequiredIfMatchesPathValue(
            types.StringValue("PASSWORD"),
            path.MatchRelative().AtParent().AtName("auth_type"),
        ),
    },
},

// Ensure at least one boolean is true
"enable_sso": schema.BoolAttribute{
    Optional: true,
    Default:  booldefault.StaticBool(false),
    Validators: []validator.Bool{
        boolvalidator.BoolAtLeastOneOfMustBeTrue(
            false, // default value
            path.MatchRelative().AtParent().AtName("enable_mfa"),
            path.MatchRelative().AtParent().AtName("enable_risk"),
        ),
    },
},
```

### Using Custom Types

```go
// PingOne resource ID with validation
"environment_id": schema.StringAttribute{
    Description: "The ID of the PingOne environment",
    Required:    true,
    CustomType:  pingonetypes.ResourceIDType{},
},

// DaVinci resource ID with validation
"flow_id": schema.StringAttribute{
    Description: "The ID of the DaVinci flow",
    Required:    true,
    CustomType:  davincitypes.ResourceIDType{},
},
```

## Contributing

When adding new validators or plan modifiers to this framework:

1. Follow the existing naming conventions
2. Add comprehensive godoc comments explaining the purpose and usage
3. Include both Description() and MarkdownDescription() methods for user-facing documentation
4. Add unit tests for the new functionality
5. Update this README with the new component
6. Follow the established patterns for error messages and validation logic

All components should be designed to work seamlessly with the Terraform Plugin Framework and provide clear, actionable error messages to users.
