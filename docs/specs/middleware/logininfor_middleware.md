# Specification: `LogininforMiddleware`

**File:** [`app/middleware/logininfor_middleware.go`](/app/middleware/logininfor_middleware.go)

## 1. Overview

The `LogininforMiddleware` is a Gin middleware designed to automatically record login attempts. It intercepts requests to login endpoints, captures user details, IP information, and the outcome of the login attempt (success or failure), and persists this information as a login log entry.

It works by wrapping the login controller. It reads the request, prepares a log entry, passes control to the actual login logic, and then inspects the response to determine the final status of the login attempt before saving the log.

## 2. Dependencies

| Module/Package | Usage |
| :--- | :--- |
| `bytes` | Used for creating and managing byte buffers to cache the request and response bodies. |
| `encoding/json` | Used to unmarshal the captured JSON response body to determine the login outcome. |
| `io` | Used for I/O operations, specifically to create a `NopCloser` for the cached request body. |
| [`mira/anima/datetime`](/anima/datetime) | Provides custom datetime types for timestamping the login attempt. |
| [`mira/anima/response`](/anima/response) | Defines the standard API response structure, used to parse the captured response. |
| [`mira/app/dto`](/app/dto) | Contains Data Transfer Objects for login requests (`LoginRequest`) and saving log data (`SaveLogininforRequest`). |
| [`mira/app/service`](/app/service) | Contains the `LogininforService` used to persist the final log entry to the database. |
| [`mira/common/ip-address`](/common/ip-address) | A utility to parse the client's IP address and User-Agent to determine location, browser, and OS. |
| [`mira/common/response-writer`](/common/response-writer) | Provides a custom `ResponseWriter` that allows capturing the response body. |
| [`mira/common/types/constant`](/common/types/constant) | Defines system-wide constants like `NORMAL_STATUS` and `EXCEPTION_STATUS`. |
| `time` | Standard library for getting the current time. |
| `github.com/gin-gonic/gin` | The web framework used. The middleware is a `gin.HandlerFunc`. |

## 3. Pseudocode

### Module: `logininfor_middleware`

#### FUNCTION `LogininforMiddleware()`: gin.HandlerFunc

This function returns a Gin handler that performs the core logging logic.

```pseudocode
BEGIN FUNCTION LogininforMiddleware() RETURNS gin.HandlerFunc
    RETURN new gin.HandlerFunc(ctx) WHERE
        // -- Phase 1: Before Handler Execution --

        // 1. Cache Request Body
        // The request body can only be read once. Cache it to allow multiple reads
        // by the binding logic and the downstream handler.
        bodyBytes <- ctx.GetRawData()
        ctx.Request.Body <- RECREATE_READABLE_BODY_FROM(bodyBytes)

        // 2. Prepare Response Capturing
        // Replace the default response writer with a custom one that buffers the response body.
        customWriter <- CREATE_CUSTOM_RESPONSE_WRITER(original_writer = ctx.Writer)

        // 3. Bind and Validate Login Parameters
        DECLARE loginParams AS dto.LoginRequest
        err <- BIND_JSON_BODY(ctx, &loginParams)

        // TDD Anchor: Test with invalid/malformed login data.
        IF err IS NOT NULL THEN
            SEND_ERROR_RESPONSE(ctx, code=400, message=err.Error())
            ABORT_REQUEST_CHAIN(ctx)
            RETURN
        END IF

        // 4. Restore Request Body
        // Binding consumes the body, so restore it for the next handler in the chain.
        ctx.Request.Body <- RECREATE_READABLE_BODY_FROM(bodyBytes)

        // 5. Gather Contextual Information
        ipInfo <- GET_IP_ADDRESS_DETAILS(ip=ctx.ClientIP(), userAgent=ctx.Request.UserAgent())

        // 6. Initialize Login Record
        loginRecord <- CREATE dto.SaveLogininforRequest WITH {
            UserName:      loginParams.Username,
            Ipaddr:        ipInfo.Ip,
            LoginLocation: ipInfo.Addr,
            Browser:       ipInfo.Browser,
            Os:            ipInfo.Os,
            Status:        constant.NORMAL_STATUS, // Assume success initially
            LoginTime:     CURRENT_TIMESTAMP()
        }

        // 7. Process Request
        // Set the custom writer and pass control to the next handler (e.g., login controller).
        ctx.Writer <- customWriter
        ctx.Next()

        // -- Phase 2: After Handler Execution --

        // 8. Analyze Response and Finalize Record
        DECLARE apiResponse AS response.Response
        UNMARSHAL_JSON(customWriter.Body.Bytes(), &apiResponse)

        // TDD Anchor: Test with a successful login response (code 200).
        // TDD Anchor: Test with a failed login response (e.g., code 500, 401).
        IF apiResponse.Code IS NOT 200 THEN
            loginRecord.Status <- constant.EXCEPTION_STATUS
        END IF
        loginRecord.Msg <- apiResponse.Msg

        // 9. Persist Login Record
        logininforService <- CREATE service.LogininforService
        // TDD Anchor: Verify that this service method is called with the correct loginRecord data.
        logininforService.CreateSysLogininfor(loginRecord)
    END
END FUNCTION
```

## 4. Edge Cases & Constraints

1.  **Consumed Request Body:** The primary constraint is that `ctx.Request.Body` is an `io.ReadCloser` and can only be read once. The middleware correctly handles this by caching the body in a byte slice (`bodyBytes`) and creating a new `io.NopCloser` from it each time it needs to be read.
2.  **Non-JSON Login Requests:** The middleware assumes the login request has a JSON body that can be bound to `dto.LoginRequest`. If the content type is different, the binding will fail, and the middleware will correctly abort with a 400 error.
3.  **Unstructured API Responses:** The logic assumes the downstream handler returns a standard JSON response that matches the `response.Response` struct. If the response is not valid JSON or has a different structure, `json.Unmarshal` will not populate `apiResponse` correctly, potentially leading to an inaccurate log message or status.
4.  **Handler Panics:** If a downstream handler panics, this middleware will not complete its post-execution phase, and the login attempt will not be logged. A separate recovery middleware should be used to handle panics.
