package ipaddress

import (
	"testing"
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

func TestClient_GetAddress_ExternalIP_MockServer(t *testing.T) {
	t.Skip("Skipping external IP test - requires proper mocking setup")
}
