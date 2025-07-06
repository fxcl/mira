# Post Validator Specification

This document outlines the validation logic for post-related operations. The validators ensure that the incoming data for creating and updating posts meets the required criteria before being processed by the application's services.

## File Location

[`app/validator/post_validator.go`](app/validator/post_validator.go)

---

## Functions

### `CreatePostValidator`

Validates the request payload when creating a new post.

#### Signature

```go
func CreatePostValidator(param dto.CreatePostRequest) error
```

#### Parameters

-   `param` ([`dto.CreatePostRequest`](app/dto/post_request.go)): The request object containing the details for the new post.

#### Validation Logic

1.  **Check `PostCode`**:
    -   **Condition**: The `PostCode` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the post code"`.

2.  **Check `PostName`**:
    -   **Condition**: The `PostName` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the post name"`.

#### TDD Anchors
- It should return an error if `PostCode` is empty.
- It should return an error if `PostName` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` if all validation checks pass.
-   Returns an `error` object with a descriptive message if any validation check fails.

---

### `UpdatePostValidator`

Validates the request payload when updating an existing post.

#### Signature

```go
func UpdatePostValidator(param dto.UpdatePostRequest) error
```

#### Parameters

-   `param` ([`dto.UpdatePostRequest`](app/dto/post_request.go)): The request object containing the details for the post to be updated.

#### Validation Logic

1.  **Check `PostId`**:
    -   **Condition**: The `PostId` field is less than or equal to `0`.
    -   **Error**: If the condition is met, it returns an error: `"parameter error"`.

2.  **Check `PostCode`**:
    -   **Condition**: The `PostCode` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the post code"`.

3.  **Check `PostName`**:
    -   **Condition**: The `PostName` field is an empty string (`""`).
    -   **Error**: If the condition is met, it returns an error: `"please enter the post name"`.

#### TDD Anchors
- It should return an error if `PostId` is not positive.
- It should return an error if `PostCode` is empty.
- It should return an error if `PostName` is empty.
- It should return `nil` if all fields are valid.

#### Return Value

-   Returns `nil` if all validation checks pass.
-   Returns an `error` object with a descriptive message if any validation check fails.
