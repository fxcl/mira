package curl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Client *http.Client
}

// HTTP request parameters
type RequestParam struct {
	Url     string
	Method  string
	Header  map[string]interface{}
	Query   map[string]interface{}
	Json    map[string]interface{}
	Form    map[string]interface{}
	Body    string
	Context context.Context
}

func DefaultClient() *Request {
	return &Request{
		Client: http.DefaultClient,
	}
}

// Initialize the client
func NewClient(client *http.Client) *Request {
	return &Request{
		Client: client,
	}
}

// Send request
func (r *Request) Send(requestParam *RequestParam) (string, error) {
	log.Printf("Sending request: %+v\n", requestParam)

	var (
		request *http.Request
		err     error
	)

	// Create a new HTTP request
	request, err = createRequest(requestParam)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set request header information
	if requestParam.Header != nil {
		for key, value := range requestParam.Header {
			request.Header.Set(key, fmt.Sprint(value))
		}
	}

	// Associate the request with the provided context
	if requestParam.Context != nil {
		request = request.WithContext(requestParam.Context)
	}

	// Send HTTP request
	result, err := r.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	defer result.Body.Close()

	// Read the response body content and add it to the buffer
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, result.Body); err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return buffer.String(), nil
}

// Create request
func createRequest(requestParam *RequestParam) (*http.Request, error) {
	switch strings.ToLower(requestParam.Method) {
	case "get":
		return getRequest(requestParam)
	case "post":
		return postRequest(requestParam)
	default:
		return getRequest(requestParam)
	}
}

// get request
func getRequest(requestParam *RequestParam) (*http.Request, error) {
	// Parse URL
	url, err := url.Parse(requestParam.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	query := url.Query()
	for key, value := range requestParam.Query {
		query.Set(key, fmt.Sprint(value))
	}
	url.RawQuery = query.Encode()

	// Update the URL of the request parameters
	requestParam.Url = url.String()

	return http.NewRequest("GET", requestParam.Url, nil)
}

// post request
func postRequest(requestParam *RequestParam) (*http.Request, error) {
	var body io.Reader

	// Pass parameters in Json format
	if requestParam.Json != nil {
		// Serialize json to a byte array
		jsonData, _ := json.Marshal(requestParam.Json)
		body = bytes.NewBuffer(jsonData)
	}

	// Pass parameters in Form format
	if requestParam.Form != nil {
		// Create form data
		formData := url.Values{}
		for key, value := range requestParam.Form {
			formData.Add(key, fmt.Sprint(value))
		}
		body = strings.NewReader(formData.Encode())
	}

	// Pass parameters in Body
	if requestParam.Body != "" {
		body = strings.NewReader(requestParam.Body)
	}

	return http.NewRequest("POST", requestParam.Url, body)
}
