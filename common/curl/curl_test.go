package curl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest_Send_Get(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the method is GET
		if r.Method != "GET" {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		// Check the query parameter
		if r.URL.Query().Get("param1") != "value1" {
			t.Errorf("Expected query param 'param1' to be 'value1', got '%s'", r.URL.Query().Get("param1"))
		}
		fmt.Fprintln(w, "Hello, client")
	}))
	defer server.Close()

	// Create a new client
	client := DefaultClient()

	// Create request parameters
	requestParam := &RequestParam{
		Url:    server.URL,
		Method: "GET",
		Query: map[string]interface{}{
			"param1": "value1",
		},
	}

	// Send the request
	body, err := client.Send(requestParam)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response body
	expectedBody := "Hello, client\n"
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestRequest_Send_Post_Json(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the method is POST
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Check the content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}
		fmt.Fprintln(w, "Hello, client")
	}))
	defer server.Close()

	// Create a new client
	client := DefaultClient()

	// Create request parameters
	requestParam := &RequestParam{
		Url:    server.URL,
		Method: "POST",
		Json: map[string]interface{}{
			"param1": "value1",
		},
	}

	// Send the request
	body, err := client.Send(requestParam)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response body
	expectedBody := "Hello, client\n"
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestRequest_Send_Post_Form(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the method is POST
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		// Check the content type
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type 'application/x-www-form-urlencoded', got '%s'", r.Header.Get("Content-Type"))
		}
		fmt.Fprintln(w, "Hello, client")
	}))
	defer server.Close()

	// Create a new client
	client := DefaultClient()

	// Create request parameters
	requestParam := &RequestParam{
		Url:    server.URL,
		Method: "POST",
		Form: map[string]interface{}{
			"param1": "value1",
		},
	}

	// Send the request
	body, err := client.Send(requestParam)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response body
	expectedBody := "Hello, client\n"
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}
