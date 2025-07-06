# Dictionary Validator Specification

This document outlines the validation logic for dictionary-related operations, including dictionary types and dictionary data. These validators ensure that incoming data for creating and updating dictionary entries is valid before being processed.

## File Location

[`app/validator/dict_validator.go`](app/validator/dict_validator.go)

---

## Dictionary Type Validators

### `CreateDictTypeValidator`

Validates the request payload for creating a new dictionary type.

#### Signature

```go
func CreateDictTypeValidator(param dto.CreateDictTypeRequest) error
```

#### Parameters

-   `param` ([`dto.CreateDictTypeRequest`](app/dto/dict_request.go)): The request object containing the details for the new dictionary type.

#### Validation Logic

1.  **Check `DictName`**:
    -   **Condition**: The `DictName` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the dictionary name"`.

2.  **Check `DictType`**:
    -   **Condition**: The `DictType` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the dictionary type"`.

#### TDD Anchors
- It should return an error if `DictName` is empty.
- It should return an error if `DictType` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` on successful validation.
-   Returns an `error` with a descriptive message if validation fails.

---

### `UpdateDictTypeValidator`

Validates the request payload for updating an existing dictionary type.

#### Signature

```go
func UpdateDictTypeValidator(param dto.UpdateDictTypeRequest) error
```

#### Parameters

-   `param` ([`dto.UpdateDictTypeRequest`](app/dto/dict_request.go)): The request object for updating the dictionary type.

#### Validation Logic

1.  **Check `DictId`**:
    -   **Condition**: The `DictId` is less than or equal to `0`.
    -   **Error**: Returns an error: `"parameter error"`.

2.  **Check `DictName`**:
    -   **Condition**: The `DictName` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the dictionary name"`.

3.  **Check `DictType`**:
    -   **Condition**: The `DictType` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the dictionary type"`.

#### TDD Anchors
- It should return an error if `DictId` is not positive.
- It should return an error if `DictName` is empty.
- It should return an error if `DictType` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` on successful validation.
-   Returns an `error` with a descriptive message if validation fails.

---

## Dictionary Data Validators

### `CreateDictDataValidator`

Validates the request payload for creating new dictionary data.

#### Signature

```go
func CreateDictDataValidator(param dto.CreateDictDataRequest) error
```

#### Parameters

-   `param` ([`dto.CreateDictDataRequest`](app/dto/dict_request.go)): The request object for creating new dictionary data.

#### Validation Logic

1.  **Check `DictLabel`**:
    -   **Condition**: The `DictLabel` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the data label"`.

2.  **Check `DictValue`**:
    -   **Condition**: The `DictValue` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the data key value"`.

#### TDD Anchors
- It should return an error if `DictLabel` is empty.
- It should return an error if `DictValue` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` on successful validation.
-   Returns an `error` with a descriptive message if validation fails.

---

### `UpdateDictDataValidator`

Validates the request payload for updating existing dictionary data.

#### Signature

```go
func UpdateDictDataValidator(param dto.UpdateDictDataRequest) error
```

#### Parameters

-   `param` ([`dto.UpdateDictDataRequest`](app/dto/dict_request.go)): The request object for updating dictionary data.

#### Validation Logic

1.  **Check `DictCode`**:
    -   **Condition**: The `DictCode` is less than or equal to `0`.
    -   **Error**: Returns an error: `"parameter error"`.

2.  **Check `DictLabel`**:
    -   **Condition**: The `DictLabel` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the data label"`.

3.  **Check `DictValue`**:
    -   **Condition**: The `DictValue` field is an empty string (`""`).
    -   **Error**: Returns an error: `"please enter the data key value"`.

#### TDD Anchors
- It should return an error if `DictCode` is not positive.
- It should return an error if `DictLabel` is empty.
- It should return an error if `DictValue` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` on successful validation.
-   Returns an `error` with a descriptive message if validation fails.
