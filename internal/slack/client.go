package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is the client used to interact with the Slack API.
type Client struct {
	httpClient *http.Client
	token      string
	workspace  string
	baseURL    *url.URL
}

// NewClient creates a new Slack client.
func NewClient(token, workspace string) *Client {
	baseURL, _ := url.Parse(fmt.Sprintf("https://%s.slack.com", workspace))
	return &Client{
		httpClient: &http.Client{},
		token:      token,
		workspace:  workspace,
		baseURL:    baseURL,
	}
}

// newTestClient creates a new client for testing purposes.
func newTestClient(serverURL string) *Client {
	baseURL, _ := url.Parse(serverURL)
	return &Client{
		httpClient: &http.Client{},
		token:      "test-token",
		workspace:  "test-workspace",
		baseURL:    baseURL,
	}
}

// Canvas represents a Slack Canvas.
type Canvas struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// FileInfo represents the metadata of a Slack file (canvas).
type FileInfo struct {
	ID       string `json:"id"`
	Group    string `json:"group"`
	IsPublic bool   `json:"is_public"`
	Shares   struct {
		Private map[string][]string `json:"private"`
	} `json:"shares"`
}

type apiResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func (c *Client) CreateCanvas(content, channelID string) (string, error) {
	rel, err := url.Parse("/api/conversations.canvases.create")
	if err != nil {
		return "", err
	}
	url := c.baseURL.ResolveReference(rel)

	values := map[string]interface{}{"channel_id": channelID, "document_content": map[string]string{"type": "markdown", "markdown": content}}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("Slack API Response: %s\n", string(body))

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	if !apiResp.Ok {
		return "", fmt.Errorf("slack API error: %s", apiResp.Error)
	}

	var result struct {
		CanvasID string `json:"canvas_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.CanvasID, nil
}

func (c *Client) ReadCanvas(id string) (*Canvas, error) {
	// There is no API to read the content of a channel canvas.
	// We can only verify it exists by trying to edit it with empty content.
	return &Canvas{ID: id, Content: ""}, nil
}

func (c *Client) CreateUserCanvas(content, channelID string, private bool, userIDs []string) (string, error) {
	rel, err := url.Parse("/api/canvases.create")
	if err != nil {
		return "", err
	}
	url := c.baseURL.ResolveReference(rel)

	values := map[string]interface{}{
		"document_content": map[string]string{"type": "markdown", "markdown": content},
	}
	if channelID != "" {
		values["channel_id"] = channelID
	}
	if private {
		values["private"] = private
		values["user_ids"] = userIDs
	}

	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("Slack API Response: %s\n", string(body))

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	if !apiResp.Ok {
		return "", fmt.Errorf("slack API error: %s", apiResp.Error)
	}

	var result struct {
		CanvasID string `json:"canvas_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.CanvasID, nil
}

func (c *Client) UpdateUserCanvas(id, content string) error {
	rel, err := url.Parse("/api/canvases.edit")
	if err != nil {
		return err
	}
	url := c.baseURL.ResolveReference(rel)

	values := map[string]interface{}{"canvas_id": id, "changes": []map[string]interface{}{{"operation": "replace", "document_content": map[string]string{"type": "markdown", "markdown": content}}}}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Slack API Response: %s\n", string(body))

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return err
	}

	if !apiResp.Ok {
		return fmt.Errorf("slack API error: %s", apiResp.Error)
	}

	return nil
}

func (c *Client) DeleteUserCanvas(id string) error {
	rel, err := url.Parse("/api/canvases.delete")
	if err != nil {
		return err
	}
	url := c.baseURL.ResolveReference(rel)

	values := map[string]interface{}{"canvas_id": id}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Slack API Response: %s\n", string(body))

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return err
	}

	if !apiResp.Ok {
		return fmt.Errorf("slack API error: %s", apiResp.Error)
	}

	return nil
}

func (c *Client) GetFileInfo(fileID string) (*FileInfo, error) {
	rel, err := url.Parse(fmt.Sprintf("/api/files.info?file=%s", fileID))
	if err != nil {
		return nil, err
	}
	url := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Slack API Response: %s\n", string(body))

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Ok {
		return nil, fmt.Errorf("slack API error: %s", apiResp.Error)
	}

	var result struct {
		File FileInfo `json:"file"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result.File, nil
}
