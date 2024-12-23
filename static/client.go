package static

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// StaticScanClient defines the interface for interacting with the scanning service.
type StaticScanClient interface {
	UploadApp(filePath string) (*UploadResult, error)
	TriggerScan(hash string) ([]byte, error)
	GetScanLogs(hash string) (*ScanLogsResult, error)
	GetJsonReport(hash string) ([]byte, error)
	GetPdfReport(hash string) ([]byte, error)
	DeleteScan(hash string) (*DeleteResult, error)
}

// apiClient implements the StaticScanClient interface.
type apiClient struct {
	baseURL  string // Base URL for the scanning service API
	apiToken string // Authorization token for API access
}

// NewClient creates and returns a new instance of StaticScanClient.
func NewClient(baseURL, apiToken string) StaticScanClient {
	return &apiClient{
		baseURL:  baseURL,
		apiToken: apiToken,
	}
}

// UploadApp uploads a file to the scanning service.
// Returns the upload result containing the file's unique hash.
func (c *apiClient) UploadApp(filePath string) (*UploadResult, error) {
	// Prepare multipart form data
	buf := &bytes.Buffer{}
	mpw := multipart.NewWriter(buf)

	// Open the file to be uploaded
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Add the file to the form data
	fWriter, err := mpw.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = io.Copy(fWriter, file); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Close the multipart writer
	if err := mpw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/upload", buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", mpw.FormDataContentType())

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed: %s", resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response JSON
	result := &UploadResult{}
	if err := json.Unmarshal(respBody, result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return result, nil
}

// TriggerScan initiates a scan for the uploaded file using its hash.
// Returns the scan response upon completion.
func (c *apiClient) TriggerScan(hash string) ([]byte, error) {
	data := url.Values{}
	data.Set("hash", hash)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/scan", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scan trigger failed: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// GetScanLogs fetches logs for the ongoing or completed scan using its hash.
func (c *apiClient) GetScanLogs(hash string) (*ScanLogsResult, error) {
	data := url.Values{}
	data.Set("hash", hash)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/scan_logs", c.baseURL), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get scan logs: %s", resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response JSON
	result := &ScanLogsResult{}
	if err := json.Unmarshal(respBody, result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return result, nil
}

// GetJsonReport fetches the JSON report for a completed scan using its hash.
func (c *apiClient) GetJsonReport(hash string) ([]byte, error) {
	return c.getReport(hash, "/api/v1/report_json")
}

// GetPdfReport fetches the PDF report for a completed scan using its hash.
func (c *apiClient) GetPdfReport(hash string) ([]byte, error) {
	return c.getReport(hash, "/api/v1/download_pdf")
}

// getReport is a helper function to fetch reports from the specified endpoint.
func (c *apiClient) getReport(hash, endpoint string) ([]byte, error) {
	data := url.Values{}
	data.Set("hash", hash)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, c.baseURL+endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch report: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// DeleteScan deletes the scan data using its hash.
// Returns the result of the deletion.
func (c *apiClient) DeleteScan(hash string) (*DeleteResult, error) {
	data := url.Values{}
	data.Set("hash", hash)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/delete_scan", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to delete scan: %s", resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response JSON
	result := &DeleteResult{}
	if err := json.Unmarshal(respBody, result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return result, nil
}
