package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type MemosClient struct {
	BaseUrl string
	Token   string
	Client  *http.Client
}

func NewMemosClient(baseUrl string, token string) *MemosClient {
	return &MemosClient{
		BaseUrl: baseUrl,
		Token:   token,
		Client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

// GetMemos retrieve a list of memos from url
func (m *MemosClient) GetMemos(pageSize int, pageToken string, filter string, tag string) (*ListMemosResponse, error) {
	reqUrl, err := m.buildUrl(pageSize, pageToken, filter, tag)
	if err != nil {
		return nil, fmt.Errorf("could not build url: %w", err)
	}

	req, err := m.createRequest(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return m.processResponse(resp)
}

func (m *MemosClient) buildUrl(pageSize int, pageToken string, filter string, tag string) (string, error) {
	baseUrl := "https://" + m.BaseUrl + "/api/v1/memos"
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/memo", baseUrl))
	if err != nil {
		return "", err
	}

	q := u.Query()
	if pageSize > 0 {
		q.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}

	filterParts := []string{}
	if filter != "" {
		filterParts = append(filterParts, fmt.Sprintf("tag == '%s'", tag))
	}

	if len(filterParts) > 0 {
		q.Set("filter", strings.Join(filterParts, " && "))
	}

	if tag != "" {
		filterParts = append(filterParts, tag)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil

}

func (m *MemosClient) createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.Token)

	return req, nil
}

func (m *MemosClient) processResponse(resp *http.Response) (*ListMemosResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var listResp ListMemosResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		fmt.Printf("Failed to unmarshal. Response body: %s\n", string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &listResp, nil
}
