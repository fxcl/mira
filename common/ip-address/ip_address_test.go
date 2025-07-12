package ipaddress

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"mira/common/curl"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func TestGetAddress_InternalIP(t *testing.T) {
	testCases := []struct {
		name      string
		ip        string
		userAgent string
		expected  string
	}{
		{
			name:      "Localhost IPv4",
			ip:        "127.0.0.1",
			userAgent: "test-agent",
			expected:  "Intranet Address",
		},
		{
			name:      "Localhost IPv6",
			ip:        "::1",
			userAgent: "test-agent",
			expected:  "Intranet Address",
		},
		{
			name:      "Private Network IP",
			ip:        "192.168.1.1",
			userAgent: "test-agent",
			expected:  "Intranet Address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := GetAddress(tc.ip, tc.userAgent)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if addr.Addr != tc.expected {
				t.Errorf("Expected address '%s', got '%s'", tc.expected, addr.Addr)
			}
			if addr.Ip != tc.ip {
				t.Errorf("Expected IP '%s', got '%s'", tc.ip, addr.Ip)
			}
		})
	}
}

func TestClient_GetAddress_ExternalIP(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the GBK response from the real API
		gbkWriter := transform.NewWriter(w, simplifiedchinese.GBK.NewEncoder())
		fmt.Fprintln(gbkWriter, `{"ip":"8.8.8.8","addr":"Google DNS"}`)
	}))
	defer server.Close()

	// Create a new client with the mock server's client
	curlClient := curl.NewClient(server.Client())
	client := NewClient(curlClient)

	// Get the address
	addr, err := client.GetAddress("8.8.8.8", "test-agent")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the address
	if addr.Addr != "Google DNS" {
		t.Errorf("Expected address 'Google DNS', got '%s'", addr.Addr)
	}
}
