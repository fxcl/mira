# Specification: Auth Validator

**Module:** `app/validator/auth_validator.go`

## Overview

This module provides validation functions for authentication-related operations, such as user registration and login. It ensures that the data provided by the client meets the required format and constraints before being processed by the application's services.

## Dependencies

- `mira/app/dto`: For request data transfer objects (`RegisterRequest`, `LoginRequest`).
- `errors`: Standard Go library for error creation.

---

## Function: `RegisterValidator`

### Signature

```go
func RegisterValidator(param dto.RegisterRequest) error
```

### Description

Validates the data provided for user registration. It checks for empty fields, password confirmation, and length constraints on username and password.

### Parameters

- `param` ([`dto.RegisterRequest`](app/dto/auth_request.go:5)): An object containing the user's registration details.
  - `Username` (string): The desired username.
  - `Password` (string): The desired password.
  - `ConfirmPassword` (string): The re-typed password for confirmation.

### Return Value

- `error`: An error object describing the validation failure, or `nil` if validation is successful.

### Pseudocode Logic

```pseudocode
FUNCTION RegisterValidator(request: RegisterRequest): Error | Nil
  // 1. Ensure username is not an empty string.
  IF request.Username is empty THEN
    RETURN Error("username cannot be empty")
  END IF

  // 2. Ensure password is not an empty string.
  IF request.Password is empty THEN
    RETURN Error("password cannot be empty")
  END IF

  // 3. Ensure the password and confirmation password match.
  IF request.ConfirmPassword is not equal to request.Password THEN
    RETURN Error("passwords do not match")
  END IF

  // 4. Validate username length.
  usernameLength = length of request.Username
  IF usernameLength < 2 OR usernameLength > 20 THEN
    RETURN Error("username length must be between 2 and 20 characters")
  END IF

  // 5. Validate password length.
  passwordLength = length of request.Password
  IF passwordLength < 5 OR passwordLength > 20 THEN
    RETURN Error("password length must be between 5 and 20 characters")
  END IF

  // 6. If all checks pass, return nil.
  RETURN Nil
END FUNCTION
```

### TDD Anchors

- **TestEmptyUsername**: Should return an error if the username is empty.
- **TestEmptyPassword**: Should return an error if the password is empty.
- **TestMismatchedPasswords**: Should return an error if `Password` and `ConfirmPassword` do not match.
- **TestUsernameTooShort**: Should return an error if the username has fewer than 2 characters.
- **TestUsernameTooLong**: Should return an error if the username has more than 20 characters.
- **TestPasswordTooShort**: Should return an error if the password has fewer than 5 characters.
- **TestPasswordTooLong**: Should return an error if the password has more than 20 characters.
- **TestValidRegistration**: Should return `nil` for a valid registration request.

---

## Function: `LoginValidator`

### Signature

```go
func LoginValidator(param dto.LoginRequest) error
```

### Description

Validates the data provided for user login. It checks that the username and password fields are not empty.

### Parameters

- `param` ([`dto.LoginRequest`](app/dto/auth_request.go:11)): An object containing the user's login credentials.
  - `Username` (string): The user's username.
  - `Password` (string): The user's password.

### Return Value

- `error`: An error object describing the validation failure, or `nil` if validation is successful.

### Pseudocode Logic

```pseudocode
FUNCTION LoginValidator(request: LoginRequest): Error | Nil
  // 1. Ensure username is not an empty string.
  IF request.Username is empty THEN
    RETURN Error("username cannot be empty")
  END IF

  // 2. Ensure password is not an empty string.
  IF request.Password is empty THEN
    RETURN Error("password cannot be empty")
  END IF

  // 3. If all checks pass, return nil.
  RETURN Nil
END FUNCTION
```

### TDD Anchors

- **TestLoginEmptyUsername**: Should return an error if the username is empty.
- **TestLoginEmptyPassword**: Should return an error if the password is empty.
- **TestValidLogin**: Should return `nil` for a valid login request.
