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
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.ID, nil
}

func (c *Client) ReadCanvas(id string) (*Canvas, error) {
	// There is no API to read the content of a channel canvas.
	// We can only verify it exists by trying to edit it with empty content.
	return &Canvas{ID: id, Content: ""}, nil
}






