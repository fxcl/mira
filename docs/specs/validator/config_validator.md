# Configuration Validator Specification

This document outlines the validation logic for configuration-related operations. These validators ensure that the incoming data for creating and updating system configurations meets the required criteria before being processed.

## File Location

[`app/validator/config_validator.go`](app/validator/config_validator.go)

---

## Functions

### `CreateConfigValidator`

Validates the request payload when creating a new configuration parameter.

#### Signature

```go
func CreateConfigValidator(param dto.CreateConfigRequest) error
```

#### Parameters

-   `param` ([`dto.CreateConfigRequest`](app/dto/config_request.go)): The request object containing the details for the new configuration.

#### Validation Logic

1.  **Check `ConfigName`**:
    -   **Condition**: The `ConfigName` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the parameter name"`.

2.  **Check `ConfigKey`**:
    -   **Condition**: The `ConfigKey` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the parameter key"`.

3.  **Check `ConfigValue`**:
    -   **Condition**: The `ConfigValue` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the parameter value"`.

#### TDD Anchors
- It should return an error if `ConfigName` is empty.
- It should return an error if `ConfigKey` is empty.
- It should return an error if `ConfigValue` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` if all validation checks pass.
-   Returns an `error` object with a descriptive message if any validation check fails.

---

### `UpdateConfigValidator`

Validates the request payload when updating an existing configuration parameter.

#### Signature

```go
func UpdateConfigValidator(param dto.UpdateConfigRequest) error
```

#### Parameters

-   `param` ([`dto.UpdateConfigRequest`](app/dto/config_request.go)): The request object containing the details for the configuration to be updated.

#### Validation Logic

1.  **Check `ConfigId`**:
    -   **Condition**: The `ConfigId` field is less than or equal to `0`.
    -   **Error**: If the condition is met, it returns an error: `"parameter error"`.

2.  **Check `ConfigName`**:
    -   **Condition**: The `ConfigName` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the parameter name"`.

3.  **Check `ConfigKey`**:
    -   **Condition**: The `ConfigKey` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the parameter key"`.

4.  **Check `ConfigValue`**:
    -   **Condition**: The `ConfigValue` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the parameter value"`.

#### TDD Anchors
- It should return an error if `ConfigId` is not positive.
- It should return an error if `ConfigName` is empty.
- It should return an error if `ConfigKey` is empty.
- It should return an error if `ConfigValue` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` if all validation checks pass.
-   Returns an `error` object with a descriptive message if any validation check fails.
