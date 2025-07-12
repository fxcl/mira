package ipaddress

import (
	"encoding/json"
	"fmt"
	"net"

	"mira/common/curl"

	"github.com/mileusna/useragent"
)

// IP address
type IpAddress struct {
	Ip         string `json:"ip"`
	Pro        string `json:"pro"`
	ProCode    string `json:"proCode"`
	City       string `json:"city"`
	CityCode   string `json:"cityCode"`
	Region     string `json:"region"`
	RegionCode string `json:"regionCode"`
	Addr       string `json:"addr"`
	Browser    string `json:"browser"`
	Os         string `json:"os"`
}

type Client struct {
	Curl *curl.Request
}

func NewClient(c *curl.Request) *Client {
	return &Client{Curl: c}
}

// GetAddress gets the address based on the IP
func (c *Client) GetAddress(ip string, userAgent string) (*IpAddress, error) {
	var ipAddress IpAddress

	// Parse userAgent
	userAgentData := useragent.Parse(userAgent)
	ipAddress.Browser = userAgentData.Name
	ipAddress.Os = userAgentData.OS

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		ipAddress.Ip = ip
		ipAddress.Addr = "Unknown Address"
		return &ipAddress, nil
	}

	if parsedIP.IsLoopback() || parsedIP.IsPrivate() {
		ipAddress.Ip = ip
		ipAddress.Addr = "Intranet Address"
		return &ipAddress, nil
	}

	body, err := c.Curl.Send(&curl.RequestParam{
		Url: "http://whois.pconline.com.cn/ipJson.jsp",
		Query: map[string]interface{}{
			"ip":   ip,
			"json": true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ip address info: %w", err)
	}

	if err := json.Unmarshal([]byte(body), &ipAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ip address info: %w", err)
	}

	return &ipAddress, nil
}

// GetAddress gets the address based on the IP
func GetAddress(ip string, userAgent string) (*IpAddress, error) {
	return NewClient(curl.DefaultClient()).GetAddress(ip, userAgent)
}
