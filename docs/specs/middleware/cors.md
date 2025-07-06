# Spec: CORS Middleware

**Module:** [`app/middleware/cors.go`](app/middleware/cors.go)

## 1. Functional Requirements

1.  **Enable Cross-Origin Requests**: The middleware must add the necessary HTTP headers to allow web clients from different origins to make requests to the API.
2.  **Handle Preflight Requests**: It must correctly handle HTTP `OPTIONS` requests (preflight requests) by returning a successful status without passing the request to other handlers.
3.  **Configurable Origin**: The allowed origin should be permissive (`*`) by default but designed to be easily changed to a specific domain for production environments.
4.  **Support Standard Methods and Headers**: The middleware should explicitly allow common HTTP methods (e.g., `GET`, `POST`, `PUT`, `DELETE`) and headers (e.g., `Content-Type`, `Authorization`) required for typical API interactions.
5.  **Allow Credentials**: It must support credentialed requests (e.g., cookies, authorization headers) from cross-origin clients.

## 2. Non-Functional Requirements

1.  **Performance**: The middleware should have minimal performance overhead, as it runs on every applicable request.
2.  **Security**: While permissive by default, the implementation should highlight where to tighten security (e.g., replacing `*` with a specific domain) to prevent security risks like Cross-Site Request Forgery (CSRF).
3.  **Integration**: It must integrate seamlessly as a standard [`gin.HandlerFunc`](https://pkg.go.dev/github.com/gin-gonic/gin#HandlerFunc) into the Gin framework's middleware chain.

## 3. Edge Cases

1.  **No Origin Header**: If a request does not include an `Origin` header, the middleware should not add any CORS headers and simply pass the request to the next handler. This is typical for same-origin requests or server-to-server communication.
2.  **Non-Preflight OPTIONS Request**: While unlikely, if an `OPTIONS` request is received that is not a CORS preflight request, it will still be handled by the preflight logic, which is an acceptable behavior.

## 4. Pseudocode

```plaintext
// Module: middleware.cors

FUNCTION CorsMiddleware(context AS HttpContext):
  // This function acts as a factory and returns the actual middleware handler.
  RETURN FUNCTION(context AS HttpContext):
    // 1. Extract request details
    request_method = context.Request.Method
    request_origin = context.Request.Header.Get("Origin")

    // 2. Add CORS headers if the request is from a browser (has an Origin header)
    IF request_origin IS NOT EMPTY THEN
      SET response_header "Access-Control-Allow-Origin" to "*" // TDD Anchor 2.1
      SET response_header "Access-Control-Allow-Methods" to "POST, GET, OPTIONS, PUT, DELETE, UPDATE"
      SET response_header "Access-Control-Allow-Headers" to "Origin, X-Requested-With, Content-Type, Accept, Authorization"
      SET response_header "Access-Control-Expose-Headers" to "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type"
      SET response_header "Access-Control-Allow-Credentials" to "true" // TDD Anchor 2.2
    END IF

    // 3. Handle preflight (OPTIONS) requests
    IF request_method IS "OPTIONS" THEN
      // Abort the request with a "No Content" status.
      // This prevents it from hitting other handlers down the chain.
      context.AbortWithStatus(204) // TDD Anchor 1.1
      RETURN // Stop further processing
    END IF

    // 4. For all other requests, proceed to the next handler in the chain.
    context.Next() // TDD Anchor 2.3, TDD Anchor 3.1
END FUNCTION
```

## 5. TDD Anchors

1.  **Test Case: Preflight Request Handling**
    *   **Given**: An `OPTIONS` request is sent with an `Origin` header (e.g., `http://example.com`).
    *   **When**: The `CorsMiddleware` processes the request.
    *   **Then**:
        *   **1.1**: The HTTP response status code MUST be `204 No Content`.
        *   The `Access-Control-Allow-Origin` response header MUST be `*`.
        *   The `Access-Control-Allow-Methods` response header MUST be present and contain the allowed methods.
        *   The next handler in the chain MUST NOT be called.

2.  **Test Case: Actual Cross-Origin Request (Non-OPTIONS)**
    *   **Given**: A `GET` request is sent with an `Origin` header (e.g., `http://example.com`).
    *   **When**: The `CorsMiddleware` processes the request.
    *   **Then**:
        *   **2.1**: The `Access-Control-Allow-Origin` response header MUST be `*`.
        *   **2.2**: The `Access-Control-Allow-Credentials` response header MUST be `true`.
        *   **2.3**: The request MUST be passed to the next handler in the chain (`context.Next()` must be called).

3.  **Test Case: Same-Origin or Server-to-Server Request**
    *   **Given**: A `GET` request is sent *without* an `Origin` header.
    *   **When**: The `CorsMiddleware` processes the request.
    *   **Then**:
        *   No `Access-Control-*` headers should be set on the response.
        *   **3.1**: The request MUST be passed to the next handler in the chain (`context.Next()` must be called).
