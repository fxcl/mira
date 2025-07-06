# Specification: Menu Validator

This document outlines the validation logic for menu creation and update operations, as implemented in [`app/validator/menu_validator.go`](app/validator/menu_validator.go).

## Modules

-   [`menu_validator`](app/validator/menu_validator.go)

## Dependencies

-   [`dto.CreateMenuRequest`](app/dto/menu_request.go)
-   [`dto.UpdateMenuRequest`](app/dto/menu_request.go)
-   [`constant`](common/types/constant/constant.go)
-   [`utils`](common/utils/utils.go)

---

## Function: `CreateMenuValidator`

Validates the request payload for creating a new menu.

### Signature

```pseudocode
FUNCTION CreateMenuValidator(request: CreateMenuRequest) -> ERROR
```

### TDD Anchors

-   **AssertsErrorOnEmptyMenuName**: Test fails if `MenuName` is an empty string.
-   **AssertsErrorOnEmptyPathForDirectoryOrMenu**: Test fails if `MenuType` is `DIRECTORY` or `MENU` and `Path` is empty.
-   **AssertsErrorOnInvalidExternalLinkPath**: Test fails if `IsFrame` is true and `Path` does not begin with `http` or `https`.
-   **AssertsSuccessOnValidRequest**: Test passes with a valid request object.

### Logic

```pseudocode
BEGIN FUNCTION CreateMenuValidator(request)
    // 1. Validate Menu Name
    IF request.MenuName IS EMPTY
        RETURN ERROR "please enter the menu name"
    END IF

    // 2. Validate Path for specific menu types
    IF (request.MenuType IS "DIRECTORY" OR request.MenuType IS "MENU") AND request.Path IS EMPTY
        RETURN ERROR "please enter the route address"
    END IF

    // 3. Validate external link path format
    IF request.IsFrame IS "YES" AND request.Path DOES NOT START WITH "http"
        RETURN ERROR "failed to add menu {request.MenuName}, the address must start with http(s)://"
    END IF

    // 4. Return success
    RETURN NIL
END FUNCTION
```

---

## Function: `UpdateMenuValidator`

Validates the request payload for updating an existing menu.

### Signature

```pseudocode
FUNCTION UpdateMenuValidator(request: UpdateMenuRequest) -> ERROR
```

### TDD Anchors

-   **AssertsErrorOnInvalidMenuID**: Test fails if `MenuId` is zero or negative.
-   **AssertsErrorOnEmptyMenuName**: Test fails if `MenuName` is an empty string.
-   **AssertsErrorOnEmptyPathForDirectoryOrMenu**: Test fails if `MenuType` is `DIRECTORY` or `MENU` and `Path` is empty.
-   **AssertsErrorOnInvalidExternalLinkPath**: Test fails if `IsFrame` is true and `Path` does not begin with `http` or `https`.
-   **AssertsErrorOnParentIdSameAsMenuId**: Test fails if `ParentId` is the same as `MenuId`.
-   **AssertsSuccessOnValidRequest**: Test passes with a valid request object.

### Logic

```pseudocode
BEGIN FUNCTION UpdateMenuValidator(request)
    // 1. Validate Menu ID
    IF request.MenuId <= 0
        RETURN ERROR "parameter error"
    END IF

    // 2. Validate Menu Name
    IF request.MenuName IS EMPTY
        RETURN ERROR "please enter the menu name"
    END IF

    // 3. Validate Path for specific menu types
    IF (request.MenuType IS "DIRECTORY" OR request.MenuType IS "MENU") AND request.Path IS EMPTY
        RETURN ERROR "please enter the route address"
    END IF

    // 4. Validate external link path format
    IF request.IsFrame IS "YES" AND request.Path DOES NOT START WITH "http"
        RETURN ERROR "failed to modify menu {request.MenuName}, the address must start with http(s)://"
    END IF

    // 5. Validate parent menu selection
    IF request.MenuId IS EQUAL TO request.ParentId
        RETURN ERROR "failed to modify menu {request.MenuName}, the parent menu cannot be itself"
    END IF

    // 6. Return success
    RETURN NIL
END FUNCTION
