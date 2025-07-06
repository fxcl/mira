# Department Validator Specification

This document specifies the validation logic for department-related operations. The validators ensure that incoming data for creating and updating departments is consistent and valid before being processed by the application's services.

## File Location

[`app/validator/dept_validator.go`](app/validator/dept_validator.go)

---

## Functions

### `CreateDeptValidator`

Validates the request payload when creating a new department.

#### Signature

```go
func CreateDeptValidator(param dto.CreateDeptRequest) error
```

#### Parameters

-   `param` ([`dto.CreateDeptRequest`](app/dto/dept_request.go)): The request object containing the details for the new department.

#### Validation Logic

1.  **Check `ParentId`**:
    -   **Condition**: The `ParentId` field is less than or equal to `0`.
    -   **Error**: If the condition is met, it returns an error: `"please select the parent department"`.

2.  **Check `DeptName`**:
    -   **Condition**: The `DeptName` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the department name"`.

#### TDD Anchors
- It should return an error if `ParentId` is not positive.
- It should return an error if `DeptName` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` if all validation checks pass.
-   Returns an `error` object with a descriptive message if any validation check fails.

---

### `UpdateDeptValidator`

Validates the request payload when updating an existing department.

#### Signature

```go
func UpdateDeptValidator(param dto.UpdateDeptRequest) error
```

#### Parameters

-   `param` ([`dto.UpdateDeptRequest`](app/dto/dept_request.go)): The request object containing the details for the department to be updated.

#### Validation Logic

1.  **Check `DeptId`**:
    -   **Condition**: The `DeptId` field is less than or equal to `0`.
    -   **Error**: If the condition is met, it returns an error: `"parameter error"`.

2.  **Check `ParentId`**:
    -   **Condition**: The `DeptId` is not `100` AND the `ParentId` is less than or equal to `0`. This implies that only the root department (ID 100) can have a non-positive `ParentId`.
    -   **Error**: If the condition is met, it returns an error: `"please select the parent department"`.

3.  **Check `DeptName`**:
    -   **Condition**: The `DeptName` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the department name"`.

4.  **Check Self-Reference**:
    -   **Condition**: The `DeptId` is the same as the `ParentId`.
    -   **Error**: If the condition is met, it returns a formatted error: `"failed to modify menu " + param.DeptName + ", the parent department cannot be itself"`.

#### TDD Anchors
- It should return an error if `DeptId` is not positive.
- It should return an error if `ParentId` is not positive (and `DeptId` is not 100).
- It should return an error if `DeptName` is empty.
- It should return an error if `DeptId` and `ParentId` are the same.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` if all validation checks pass.
-   Returns an `error` object with a descriptive message if any validation check fails.
