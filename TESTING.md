# Mira Project Testing Documentation

This document provides comprehensive information about the testing framework and practices used in the Mira project.

## Table of Contents

1. [Testing Overview](#testing-overview)
2. [Test Structure](#test-structure)
3. [Running Tests](#running-tests)
4. [Test Categories](#test-categories)
5. [Writing Tests](#writing-tests)
6. [Mocking and Test Doubles](#mocking-and-test-doubles)
7. [Performance Testing](#performance-testing)
8. [Integration Testing](#integration-testing)
9. [Coverage Analysis](#coverage-analysis)
10. [Best Practices](#best-practices)

## Testing Overview

The Mira project follows a comprehensive testing strategy that includes:

- **Unit Tests**: Testing individual components in isolation
- **Integration Tests**: Testing component interactions
- **Performance Tests**: Benchmarking critical operations
- **Security Tests**: Validating authentication and authorization
- **Controller Tests**: Testing HTTP handlers and API endpoints

### Key Features

- ✅ **Controller Layer Tests**: Complete coverage for all system controllers
- ✅ **Service Layer Tests**: Enhanced tests with edge cases and performance validation
- ✅ **Middleware Tests**: Authentication, CORS, and logging middleware
- ✅ **Utility Function Tests**: Comprehensive testing of common utilities
- ✅ **Integration Tests**: End-to-end API testing
- ✅ **Benchmark Tests**: Performance measurement and regression detection
- ✅ **Security Tests**: Authentication and authorization validation

## Test Structure

```
mira/
├── app/
│   ├── controller/
│   │   ├── system/
│   │   │   ├── user_controller_test.go
│   │   │   └── role_controller_test.go
│   │   └── auth_controller_test.go
│   ├── middleware/
│   │   ├── auth_middleware_test.go
│   │   ├── cors_test.go
│   │   └── oper_log_middleware_test.go
│   ├── security/
│   │   └── security_test.go
│   └── service/
│       └── user_service_enhanced_test.go
├── common/
│   └── utils/
│       └── utils_test.go
├── tests/
│   ├── integration/
│   │   └── api_integration_test.go
│   └── benchmark/
│       └── service_benchmark_test.go
├── test_runner.sh           # Comprehensive test runner
└── TESTING.md              # This documentation
```

## Running Tests

### Using the Test Runner Script

The project includes a comprehensive test runner script (`test_runner.sh`) that provides various testing options:

```bash
# Run all tests
./test_runner.sh

# Run specific test categories
./test_runner.sh unit          # Unit tests only
./test_runner.sh controller    # Controller tests only
./test_runner.sh integration   # Integration tests only
./test_runner.sh benchmark     # Performance benchmarks
./test_runner.sh coverage      # Tests with coverage analysis
./test_runner.sh race          # Race condition tests

# View test statistics
./test_runner.sh stats

# Clean test artifacts
./test_runner.sh clean

# Show help
./test_runner.sh help
```

### Using Go Test Commands

```bash
# Run all tests (recommended with CGO disabled to avoid linking issues)
CGO_ENABLED=0 go test ./...

# Run tests with verbose output
CGO_ENABLED=0 go test -v ./...

# Run tests with coverage
CGO_ENABLED=0 go test -cover ./...

# Generate coverage report
CGO_ENABLED=0 go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
CGO_ENABLED=0 go test -bench=. ./...

# Run tests with race detection
CGO_ENABLED=0 go test -race ./...

# Run specific test file
CGO_ENABLED=0 go test -v ./app/controller/system/user_controller_test.go

# Run specific test function
CGO_ENABLED=0 go test -run TestUserController_List ./app/controller/system

# Alternative: Use the test script for convenience
./test.sh
```

### Important: CGO Linking Issues

On some systems (especially macOS), you may encounter linking errors like:
```
ld: library not found for -lresolv
```

**Solution**: Set `CGO_ENABLED=0` when running tests:

```bash
export CGO_ENABLED=0
go test ./...
```

Or use the provided test script which handles this automatically:
```bash
./test.sh
```

This prevents the linker from trying to link with system libraries that may not be available or properly configured.

## Test Categories

### 1. Unit Tests

Unit tests focus on testing individual components in isolation:

- **Service Layer**: Business logic testing with mocked dependencies
- **Utility Functions**: Pure function testing with various inputs
- **Security Components**: Authentication and authorization logic

### 2. Controller Tests

Controller tests validate HTTP handlers and API endpoints:

- **Request/Response Validation**: Proper handling of HTTP requests
- **Authentication**: Protected endpoint access control
- **Error Handling**: Graceful error responses
- **Data Binding**: Request data validation and binding

### 3. Integration Tests

Integration tests verify component interactions:

- **API Endpoints**: End-to-end request/response cycles
- **Middleware Integration**: Authentication, CORS, logging
- **Database Integration**: Data persistence and retrieval
- **External Service Integration**: Third-party API interactions

### 4. Performance Tests

Performance tests measure and validate system performance:

- **Benchmark Tests**: Operation timing and resource usage
- **Load Testing**: Performance under concurrent load
- **Memory Profiling**: Memory allocation and garbage collection
- **Regression Testing**: Performance degradation detection

## Writing Tests

### Test File Structure

```go
package controller

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)

// Test structure
type UserControllerTestSuite struct {
    suite.Suite
    controller  *UserController
    router      *gin.Engine
    mockService *MockUserService
}

func TestUserControllerTestSuite(t *testing.T) {
    suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) SetupTest() {
    // Setup test environment
}

func (suite *UserControllerTestSuite) TestUserController_List() {
    // Test implementation
}
```

### Test Naming Conventions

- **Test Files**: `*_test.go`
- **Test Functions**: `TestFunctionName` or `TestStructName_MethodName`
- **Benchmark Functions**: `BenchmarkFunctionName`
- **Suite Tests**: `StructNameTestSuite`

### Example Controller Test

```go
func (suite *UserControllerTestSuite) TestUserController_List_Success() {
    // Arrange
    param := dto.UserListRequest{
        PageNum:  1,
        PageSize: 10,
    }

    expectedUsers := []dto.UserListResponse{
        {UserId: 1, UserName: "testuser"},
    }

    suite.mockService.On("GetUserList", param, 1, true).
        Return(expectedUsers, int64(1))

    // Act
    body, _ := json.Marshal(param)
    req, _ := http.NewRequest("GET", "/list", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    c, w := suite.createTestContext(1, "testuser")
    c.Request = req

    suite.controller.List(c)

    // Assert
    assert.Equal(suite.T(), http.StatusOK, w.Code)

    var resp response.Response
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), float64(0), resp.Code)

    suite.mockService.AssertExpectations(suite.T())
}
```

## Mocking and Test Doubles

### Mocking Services

```go
// MockUserService implements service.UserServiceInterface
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUserList(param dto.UserListRequest, userId int, withDataScope bool) ([]dto.UserListResponse, int64) {
    args := m.Called(param, userId, withDataScope)
    return args.Get(0).([]dto.UserListResponse), args.Get(1).(int64)
}

// Usage in tests
mockService := new(MockUserService)
mockService.On("GetUserList", mock.Anything, 1, true).
    Return(expectedUsers, int64(1))
```

### Test Context Creation

```go
func (suite *UserControllerTestSuite) createTestContext(userId int, userName string) (*gin.Context, *httptest.ResponseRecorder) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    userToken := &token.UserTokenResponse{
        UserId:   userId,
        UserName: userName,
    }

    c.Set(token.UserTokenKey, userToken)
    c.Request = httptest.NewRequest("GET", "/", nil)

    return c, w
}
```

## Performance Testing

### Benchmark Tests

```go
func BenchmarkUserController_List(b *testing.B) {
    controller, _, mockUserService, _, _, _ := setupUserControllerTest()

    param := dto.UserListRequest{PageNum: 1, PageSize: 10}
    expectedUsers := make([]dto.UserListResponse, 100)

    mockUserService.On("GetUserList", param, 1, true).
        Return(expectedUsers, int64(100))

    body, _ := json.Marshal(param)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)

        userToken := &token.UserTokenResponse{UserId: 1, UserName: "testuser"}
        c.Set(token.UserTokenKey, userToken)

        req, _ := http.NewRequest("GET", "/list", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        c.Request = req

        controller.List(c)
    }
}
```

### Load Testing

```go
func TestConcurrentAccess(t *testing.T) {
    const numGoroutines = 10
    const operationsPerGoroutine = 5

    done := make(chan bool, numGoroutines)

    for i := 0; i < numGoroutines; i++ {
        go func() {
            for j := 0; j < operationsPerGoroutine; j++ {
                // Perform concurrent operations
                service.GetUserByUsername("testuser")
            }
            done <- true
        }()
    }

    // Wait for completion
    for i := 0; i < numGoroutines; i++ {
        <-done
    }
}
```

## Integration Testing

### API Integration Tests

```go
func (suite *APIIntegrationTestSuite) TestAuthEndpoints() {
    suite.T().Run("GET /auth/captcha", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/auth/captcha", nil)
        w := httptest.NewRecorder()
        suite.router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var resp response.Response
        err := json.Unmarshal(w.Body.Bytes(), &resp)
        assert.NoError(t, err)
        assert.Equal(t, float64(0), resp.Code)
    })
}
```

### Security Tests

```go
func (suite *APIIntegrationTestSuite) TestSecurity() {
    suite.T().Run("Large payload protection", func(t *testing.T) {
        largePayload := strings.Repeat("a", 10*1024*1024) // 10MB
        req, _ := http.NewRequest("POST", "/system/user", strings.NewReader(largePayload))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        suite.router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
    })
}
```

## Coverage Analysis

### Generating Coverage Reports

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Show coverage percentage
go tool cover -func=coverage.out
```

### Coverage Goals

- **Service Layer**: > 80% coverage
- **Controller Layer**: > 75% coverage
- **Utility Functions**: > 90% coverage
- **Security Components**: > 85% coverage

## Best Practices

### 1. Test Organization

- **Group Related Tests**: Use test suites for complex components
- **Table-Driven Tests**: For multiple test cases with similar structure
- **Descriptive Names**: Clear, descriptive test and function names
- **AAA Pattern**: Arrange, Act, Assert structure

### 2. Test Data Management

- **Test Factories**: Create reusable test data generators
- **Cleanup**: Proper test cleanup to avoid interference
- **Isolation**: Tests should not depend on each other
- **Deterministic**: Tests should produce consistent results

### 3. Mock Usage

- **Interface-Based**: Mock interfaces, not concrete implementations
- **Specific Expectations**: Set clear expectations on mock calls
- **Verification**: Assert that expected interactions occurred
- **Realistic Behavior**: Mocks should behave like real dependencies

### 4. Error Testing

- **Edge Cases**: Test boundary conditions and invalid inputs
- **Error Paths**: Test error handling and recovery
- **Security Scenarios**: Test potential security vulnerabilities
- **Performance Edge Cases**: Test behavior with large datasets

### 5. Performance Testing

- **Baseline Measurements**: Establish performance baselines
- **Regression Detection**: Monitor for performance degradation
- **Resource Usage**: Test memory and CPU consumption
- **Concurrent Scenarios**: Test under concurrent load

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.21
    - name: Run tests
      run: ./test_runner.sh all
    - name: Upload coverage
      uses: codecov/codecov-action@v1
```

## Troubleshooting

### Common Issues

1. **Race Conditions**: Use `-race` flag to detect
2. **Test Dependencies**: Ensure tests are isolated
3. **Mock Mismatches**: Verify mock expectations
4. **Import Issues**: Check module dependencies
5. **Database State**: Clean test database between runs

### Debugging Tests

```bash
# Run tests with verbose output
go test -v ./...

# Run specific test with debugging
go test -v -run TestSpecificFunction ./path/to/package

# Show test output even on success
go test -v ./path/to/package

# Run tests with race detection
go test -race ./path/to/package
```

## Future Enhancements

- **Contract Testing**: API contract validation
- **Property-Based Testing**: Randomized test input generation
- **Visual Regression Testing**: UI component testing
- **Contract Testing**: Service-level agreements validation
- **Load Testing**: Automated load testing in CI/CD

## Contributing

When adding new tests:

1. Follow existing test patterns and naming conventions
2. Include both positive and negative test cases
3. Add appropriate mock implementations
4. Ensure tests are isolated and deterministic
5. Update documentation as needed
6. Verify test coverage meets project standards

For questions or issues with testing, please refer to the project maintainers or create an issue in the project repository.