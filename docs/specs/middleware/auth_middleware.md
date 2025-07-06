# Spec: Auth Middleware

This document outlines the functional specification, pseudocode, and testing anchors for the authentication middleware.

**File:** [`app/middleware/auth_middleware.go`](app/middleware/auth_middleware.go)

## 1. Functional Requirements

1.  **Authentication Check**: The middleware must verify if a user is authenticated for the incoming request.
2.  **Unauthenticated Access**: If no authenticated user is found, the middleware must immediately terminate the request and return a `401 Unauthorized` error.
3.  **Token Refresh**: The middleware must check the expiration time of the user's token. If the token is set to expire within a predefined threshold (20 minutes), it should automatically refresh the token.
4.  **User Status Verification**: The middleware must check the status of the authenticated user.
5.  **Disabled User Handling**: If the user's status is not "normal" (e.g., disabled, locked), the middleware must terminate the request and return a custom `601` error.
6.  **Request Continuation**: If the user is authenticated, their status is normal, and the token is valid, the middleware must pass the request to the next handler in the chain.

## 2. Constraints & Edge Cases

1.  **Dependency on `security.GetAuthUser`**: The middleware's logic is entirely dependent on the `security.GetAuthUser` function correctly retrieving user data from the request context. If this function fails or returns invalid data, the middleware's behavior is undefined.
2.  **Token Refresh Failure**: The specification does not define behavior if `token.RefreshToken` fails. The current implementation proceeds with the request even if the refresh fails.
3.  **Time Synchronization**: The token expiration check relies on the server's clock being synchronized. Significant clock skew could lead to premature or delayed token refreshes.
4.  **Concurrent Requests**: If multiple requests arrive simultaneously with a near-expiry token, the refresh logic might be triggered multiple times. The `token.RefreshToken` implementation should be idempotent or handle this gracefully.

## 3. Pseudocode

```plaintext
MODULE AuthMiddleware

  DEPENDENCIES
    - ginContext: Framework context for request/response handling.
    - security: Module to retrieve authenticated user information.
    - response: Module for creating standardized JSON error responses.
    - token: Module for handling token operations (e.g., refresh).
    - constants: Application-defined constants (e.g., user status).
    - time: System time library.

  FUNCTION AuthMiddleware(): gin.HandlerFunc
    // This function is a factory that returns the actual middleware handler.
    RETURN FUNCTION (context)
      // 1. Retrieve authenticated user from the request context.
      authUser = security.GetAuthUser(context)

      // 2. Handle unauthenticated users.
      IF authUser IS NULL THEN
        errorResponse = response.NewError()
        errorResponse.SetCode(401)
        errorResponse.SetMsg("Not logged in")
        errorResponse.Json(context)
        context.Abort() // Stop processing further handlers.
        RETURN
      END IF

      // 3. Check if the token needs refreshing.
      // Define the time window for proactive refresh.
      refreshThreshold = time.Now() + 20 * time.Minute
      IF authUser.ExpireTime < refreshThreshold THEN
        // The token is about to expire, so refresh it.
        token.RefreshToken(context, authUser.UserTokenResponse)
      END IF

      // 4. Check if the user account is active.
      IF authUser.Status IS NOT constants.NORMAL_STATUS THEN
        errorResponse = response.NewError()
        errorResponse.SetCode(601) // Custom error code for disabled user.
        errorResponse.SetMsg("User is disabled")
        errorResponse.Json(context)
        context.Abort() // Stop processing further handlers.
        RETURN
      END IF

      // 5. If all checks pass, proceed to the next handler.
      context.Next()
    END FUNCTION
  END FUNCTION

END MODULE
```

## 4. TDD Anchors

```go
package middleware_test

import (
    "testing"
    // ... other necessary imports
)

// Test case for when no user is authenticated.
func TestAuthMiddleware_UnauthenticatedRequest(t *testing.T) {
    // SETUP: Create a request context without an authenticated user.
    // EXPECT: Response status code should be 401.
    // EXPECT: Response body should contain "Not logged in".
}

// Test case for when an authenticated user is disabled or locked.
func TestAuthMiddleware_DisabledUser(t *testing.T) {
    // SETUP: Create a request context with an authenticated user whose status is NOT 'normal'.
    // EXPECT: Response status code should be 601.
    // EXPECT: Response body should contain "User is disabled".
}

// Test case for when a token is close to expiring.
func TestAuthMiddleware_TokenNearingExpiry(t *testing.T) {
    // SETUP: Create a request context with a user whose token expires in < 20 minutes.
    // SETUP: Mock the `token.RefreshToken` function.
    // ACTION: Execute the middleware.
    // EXPECT: The mocked `token.RefreshToken` function to be called exactly once.
    // EXPECT: The request to be passed to the next handler (context.Next() is called).
}

// Test case for when a token is not close to expiring.
func TestAuthMiddleware_TokenNotNearingExpiry(t *testing.T) {
    // SETUP: Create a request context with a user whose token expires in > 20 minutes.
    // SETUP: Mock the `token.RefreshToken` function.
    // ACTION: Execute the middleware.
    // EXPECT: The mocked `token.RefreshToken` function NOT to be called.
    // EXPECT: The request to be passed to the next handler.
}

// Test case for the standard happy path.
func TestAuthMiddleware_ValidUserAndToken(t *testing.T) {
    // SETUP: Create a request context with a normal-status user whose token is not expiring soon.
    // ACTION: Execute the middleware.
    // EXPECT: The request to be passed to the next handler.
    // EXPECT: The response status code to be set by the final handler (e.g., 200 OK), not the middleware.
}
