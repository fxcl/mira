# Specification: User Validator

This document outlines the validation logic for user-related operations. Each function ensures that incoming data meets the required format and constraints before being processed by the application.

## Module: `app/validator/user_validator.go`

This module contains all validation functions related to user management.

---

### 1. `UpdateProfileValidator`

**Objective:** Validates the request payload for updating a user's own profile information.

#### Functional Requirements:

1.  The user's nickname must not be empty.
2.  The user's email must be in a valid format.
3.  The user's phone number must be in a valid format.

#### Pseudocode:

```pseudocode
FUNCTION UpdateProfileValidator(request):
  // 1. Validate Nickname
  IF request.NickName is empty THEN
    RETURN error "Please enter the user nickname"
  END IF

  // 2. Validate Email Format
  IF request.Email does not match EMAIL_REGEX THEN
    RETURN error "Incorrect email format"
  END IF

  // 3. Validate Phone Number Format
  IF request.Phonenumber does not match PHONE_REGEX THEN
    RETURN error "Incorrect mobile number format"
  END IF

  // 4. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_update_profile_fails_with_empty_nickname`
-   `test_update_profile_fails_with_invalid_email`
-   `test_update_profile_fails_with_invalid_phone`
-   `test_update_profile_succeeds_with_valid_data`

---

### 2. `UserProfileUpdatePwdValidator`

**Objective:** Validates the request for a user changing their own password.

#### Functional Requirements:

1.  The old password must be provided.
2.  The new password must be provided.

#### Pseudocode:

```pseudocode
FUNCTION UserProfileUpdatePwdValidator(request):
  // 1. Validate Old Password
  IF request.OldPassword is empty THEN
    RETURN error "Please enter the old password"
  END IF

  // 2. Validate New Password
  IF request.NewPassword is empty THEN
    RETURN error "Please enter the new password"
  END IF

  // 3. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_update_password_fails_with_empty_old_password`
-   `test_update_password_fails_with_empty_new_password`
-   `test_update_password_succeeds_with_valid_data`

---

### 3. `CreateUserValidator`

**Objective:** Validates the request to create a new user (typically by an admin).

#### Functional Requirements:

1.  A nickname must be provided.
2.  A username must be provided.
3.  A password must be provided.
4.  If a phone number is provided, it must be in a valid format.
5.  If an email is provided, it must be in a valid format.

#### Pseudocode:

```pseudocode
FUNCTION CreateUserValidator(request):
  // 1. Validate Nickname
  IF request.NickName is empty THEN
    RETURN error "Please enter the user nickname"
  END IF

  // 2. Validate Username
  IF request.UserName is empty THEN
    RETURN error "Please enter the username"
  END IF

  // 3. Validate Password
  IF request.Password is empty THEN
    RETURN error "Please enter the user password"
  END IF

  // 4. Validate Phone Number (if provided)
  IF request.Phonenumber is not empty AND request.Phonenumber does not match PHONE_REGEX THEN
    RETURN error "Incorrect mobile number format"
  END IF

  // 5. Validate Email (if provided)
  IF request.Email is not empty AND request.Email does not match EMAIL_REGEX THEN
    RETURN error "Incorrect email account format"
  END IF

  // 6. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_create_user_fails_with_empty_nickname`
-   `test_create_user_fails_with_empty_username`
-   `test_create_user_fails_with_empty_password`
-   `test_create_user_fails_with_invalid_phone`
-   `test_create_user_fails_with_invalid_email`
-   `test_create_user_succeeds_with_valid_data`
-   `test_create_user_succeeds_with_empty_optional_fields`

---

### 4. `UpdateUserValidator`

**Objective:** Validates the request to update a user's details (typically by an admin).

#### Functional Requirements:

1.  A valid `UserId` (> 0) must be provided.
2.  A nickname must be provided.
3.  If a phone number is provided, it must be in a valid format.
4.  If an email is provided, it must be in a valid format.

#### Pseudocode:

```pseudocode
FUNCTION UpdateUserValidator(request):
  // 1. Validate UserId
  IF request.UserId <= 0 THEN
    RETURN error "Parameter error"
  END IF

  // 2. Validate Nickname
  IF request.NickName is empty THEN
    RETURN error "Please enter the user nickname"
  END IF

  // 3. Validate Phone Number (if provided)
  IF request.Phonenumber is not empty AND request.Phonenumber does not match PHONE_REGEX THEN
    RETURN error "Incorrect mobile number format"
  END IF

  // 4. Validate Email (if provided)
  IF request.Email is not empty AND request.Email does not match EMAIL_REGEX THEN
    RETURN error "Incorrect email account format"
  END IF

  // 5. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_update_user_fails_with_invalid_userid`
-   `test_update_user_fails_with_empty_nickname`
-   `test_update_user_fails_with_invalid_phone`
-   `test_update_user_fails_with_invalid_email`
-   `test_update_user_succeeds_with_valid_data`

---

### 5. `RemoveUserValidator`

**Objective:** Validates the request to remove one or more users.

#### Functional Requirements:

1.  The super administrator (user ID 1) cannot be deleted.
2.  The currently authenticated user cannot delete themselves.

#### Pseudocode:

```pseudocode
FUNCTION RemoveUserValidator(userIds_to_delete, authenticated_user_id):
  // 1. Check for Super Admin deletion
  IF userIds_to_delete contains 1 THEN
    RETURN error "The super administrator cannot be deleted"
  END IF

  // 2. Check for self-deletion
  IF userIds_to_delete contains authenticated_user_id THEN
    RETURN error "The current user cannot be deleted"
  END IF

  // 3. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_remove_user_fails_when_deleting_super_admin`
-   `test_remove_user_fails_when_deleting_self`
-   `test_remove_user_succeeds_for_valid_users`

---

### 6. `ChangeUserStatusValidator`

**Objective:** Validates the request to change a user's account status (e.g., active/inactive).

#### Functional Requirements:

1.  A valid `UserId` (> 0) must be provided.
2.  A status value must be provided.

#### Pseudocode:

```pseudocode
FUNCTION ChangeUserStatusValidator(request):
  // 1. Validate UserId
  IF request.UserId <= 0 THEN
    RETURN error "Parameter error"
  END IF

  // 2. Validate Status
  IF request.Status is empty THEN
    RETURN error "Please select a status"
  END IF

  // 3. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_change_status_fails_with_invalid_userid`
-   `test_change_status_fails_with_empty_status`
-   `test_change_status_succeeds_with_valid_data`

---

### 7. `ResetUserPwdValidator`

**Objective:** Validates the request to reset a user's password (typically by an admin).

#### Functional Requirements:

1.  A valid `UserId` (> 0) must be provided.
2.  A new password must be provided.

#### Pseudocode:

```pseudocode
FUNCTION ResetUserPwdValidator(request):
  // 1. Validate UserId
  IF request.UserId <= 0 THEN
    RETURN error "Parameter error"
  END IF

  // 2. Validate Password
  IF request.Password is empty THEN
    RETURN error "Please enter the user password"
  END IF

  // 3. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_reset_password_fails_with_invalid_userid`
-   `test_reset_password_fails_with_empty_password`
-   `test_reset_password_succeeds_with_valid_data`

---

### 8. `ImportUserValidator`

**Objective:** Validates user data being imported from an external source (e.g., a spreadsheet). This is similar to `CreateUserValidator` but may have different constraints (e.g., password is not required on import).

#### Functional Requirements:

1.  A nickname must be provided.
2.  A username must be provided.
3.  If a phone number is provided, it must be in a valid format.
4.  If an email is provided, it must be in a valid format.

#### Pseudocode:

```pseudocode
FUNCTION ImportUserValidator(request):
  // 1. Validate Nickname
  IF request.NickName is empty THEN
    RETURN error "Please enter the user nickname"
  END IF

  // 2. Validate Username
  IF request.UserName is empty THEN
    RETURN error "Please enter the username"
  END IF

  // 3. Validate Phone Number (if provided)
  IF request.Phonenumber is not empty AND request.Phonenumber does not match PHONE_REGEX THEN
    RETURN error "Incorrect mobile number format"
  END IF

  // 4. Validate Email (if provided)
  IF request.Email is not empty AND request.Email does not match EMAIL_REGEX THEN
    RETURN error "Incorrect email account format"
  END IF

  // 5. Success
  RETURN nil
END FUNCTION
```

#### TDD Anchors:

-   `test_import_user_fails_with_empty_nickname`
-   `test_import_user_fails_with_empty_username`
-   `test_import_user_fails_with_invalid_phone`
-   `test_import_user_fails_with_invalid_email`
-   `test_import_user_succeeds_with_valid_data`
