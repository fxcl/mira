### **Specification: `has_perm` Middleware**

**Objective:** This document outlines the logic and testing strategy for the `has_perm.go` middleware. The primary function of this middleware is to act as a gatekeeper for API endpoints, ensuring that the incoming request is from a user who possesses the specific permission required to access the resource.

---

### **1. Functional Requirements**

1.  **Parameterized Permissions:** The middleware must be configurable with a specific permission string (e.g., `'system:user:list'`).
2.  **Super Admin Bypass:** A designated "super admin" user (identified by `UserID == 1`) must bypass all permission checks and be granted immediate access.
3.  **Permission Verification:** For non-super-admin users, the middleware must verify if the user has the required permission by checking against a central authority (the `security` service).
4.  **Access Granted:** If the user has the required permission, the request should be passed to the next handler in the Gin processing chain.
5.  **Access Denied:** If the user does not have the required permission, the request chain must be aborted, and a standardized JSON error response with `code: 601` and message `"Insufficient permissions"` must be returned to the client.

---

### **2. Non-Functional Requirements**

1.  **Modularity:** The core permission-checking logic should be delegated to the [`app/security/security.go`](app/security/security.go) module to maintain separation of concerns.
2.  **Reusability:** The middleware must be designed as a higher-order function that returns a [`gin.HandlerFunc`](https://pkg.go.dev/github.com/gin-gonic/gin#HandlerFunc), allowing it to be easily applied to any route definition.
3.  **Standardized Responses:** All error responses must conform to the project's standard error structure, generated via the [`anima/response/response.go`](anima/response/response.go) module.

---

### **3. Pseudocode**

```plaintext
MODULE Middleware.HasPermission

DEPENDENCIES:
  - GinFramework.Context
  - SecurityService
  - ResponseService

//--------------------------------------------------------------------------------
// FUNCTION HasPerm(required_permission: STRING) -> GIN_HANDLER
//
// @description:
//   A higher-order function that creates a Gin middleware handler.
//   The returned handler will check if the authenticated user has the
//   'required_permission'.
//
// @param required_permission: The permission string to validate (e.g., "system:user:list").
// @returns: A Gin handler function to be used in a route.
//--------------------------------------------------------------------------------
FUNCTION HasPerm(required_permission):

  // Define and return the actual middleware handler
  RETURN FUNCTION(context):

    // 1. Get Authenticated User
    current_user_id = SecurityService.GetAuthUserId(context)

    // 2. Handle Super Admin Edge Case
    // The user with ID 1 is a super administrator and bypasses all checks.
    IF current_user_id == 1 THEN
      context.Next() // Proceed to the main controller logic
      RETURN
    END IF

    // 3. Verify Permission for Regular Users
    user_has_permission = SecurityService.HasPerm(current_user_id, required_permission)

    // 4. Authorization Logic
    IF user_has_permission IS FALSE THEN
      // 4a. Access Denied
      error_response = ResponseService.NewError()
      error_response.SetCode(601) // 601: Insufficient Permissions
      error_response.SetMsg("Insufficient permissions")
      error_response.Json(context) // Send JSON response

      context.Abort() // Halt the request chain
      RETURN
    END IF

    // 5. Access Granted
    context.Next() // Proceed to the main controller logic

  END FUNCTION
END FUNCTION
```

---

### **4. TDD Anchors (Test Plan)**

A robust test suite should be created to validate the middleware's behavior under various conditions.

*   **`Test: HasPerm_WithSuperAdmin`**
    *   **Given:** A request context where [`security.GetAuthUserId()`](app/security/security.go:17) returns `1`.
    *   **When:** The `HasPerm` middleware is invoked with any permission string.
    *   **Then:**
        *   The `context.Next()` method must be called exactly once.
        *   The [`security.HasPerm()`](app/security/security.go:23) function must NOT be called.
        *   The response status code should not be set to an error state by this middleware.
        *   `context.Abort()` must NOT be called.

*   **`Test: HasPerm_WithSufficientPermission`**
    *   **Given:** A request context where [`security.GetAuthUserId()`](app/security/security.go:17) returns a non-super-admin ID (e.g., `101`).
    *   **And:** The [`security.HasPerm()`](app/security/security.go:23) function is mocked to return `true` for that user and permission.
    *   **When:** The `HasPerm` middleware is invoked.
    *   **Then:**
        *   The `context.Next()` method must be called exactly once.
        *   `context.Abort()` must NOT be called.

*   **`Test: HasPerm_WithInsufficientPermission`**
    *   **Given:** A request context where [`security.GetAuthUserId()`](app/security/security.go:17) returns a non-super-admin ID (e.g., `101`).
    *   **And:** The [`security.HasPerm()`](app/security/security.go:23) function is mocked to return `false`.
    *   **When:** The `HasPerm` middleware is invoked.
    *   **Then:**
        *   The `context.Next()` method must NOT be called.
        *   `context.Abort()` must be called exactly once.
        *   A JSON response must be sent with `code: 601` and `msg: "Insufficient permissions"`.

*   **`Test: HasPerm_WithUnauthenticatedUser`**
    *   **Given:** A request context where [`security.GetAuthUserId()`](app/security/security.go:17) returns a zero-value or error (simulating a user who failed authentication prior to this middleware).
    *   **And:** The [`security.HasPerm()`](app/security/security.go:23) function returns `false` for a zero-value user ID.
    *   **When:** The `HasPerm` middleware is invoked.
    *   **Then:**
        *   The behavior should be identical to the `InsufficientPermission` test case.
        *   `context.Abort()` must be called, and a `601` error returned.
