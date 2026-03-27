package gitforge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type httpClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func newHTTPClient(baseURL, token string, client *http.Client) *httpClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &httpClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		httpClient: client,
	}
}

func (c *httpClient) get(ctx context.Context, path string, query url.Values) (json.RawMessage, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *httpClient) post(ctx context.Context, path string, body any) (json.RawMessage, error) {
	return c.mutate(ctx, http.MethodPost, path, body)
}

func (c *httpClient) patch(ctx context.Context, path string, body any) (json.RawMessage, error) {
	return c.mutate(ctx, http.MethodPatch, path, body)
}

func (c *httpClient) del(ctx context.Context, path string, query url.Values) error {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req)
	return err
}

func (c *httpClient) getRaw(ctx context.Context, path string, query url.Values) ([]byte, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp.StatusCode, raw)
	}

	return raw, nil
}

func (c *httpClient) delWithBody(ctx context.Context, path string, body any) (json.RawMessage, error) {
	return c.mutate(ctx, http.MethodDelete, path, body)
}

func (c *httpClient) mutate(ctx context.Context, method, path string, body any) (json.RawMessage, error) {
	u := c.baseURL + path
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req)
}

func (c *httpClient) do(req *http.Request) (json.RawMessage, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil, nil
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp.StatusCode, raw)
	}

	return json.RawMessage(raw), nil
}

func (c *httpClient) parseError(status int, raw []byte) error {
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		return &GitForgeError{StatusCode: status, Code: "unknown", Message: fmt.Sprintf("HTTP %d", status)}
	}

	code := stringFromMap(body, "code")
	if code == "" {
		code = stringFromMap(body, "error")
	}
	if code == "" {
		code = "unknown"
	}
	message := stringFromMap(body, "message")
	if message == "" {
		message = fmt.Sprintf("HTTP %d", status)
	}

	if status == 409 && code == "branch_moved" {
		if sha := stringFromMap(body, "currentSha"); sha != "" {
			return &RefUpdateError{
				GitForgeError: &GitForgeError{StatusCode: 409, Code: code, Message: message},
				CurrentSHA:    sha,
			}
		}
	}

	return &GitForgeError{StatusCode: status, Code: code, Message: message}
}

func stringFromMap(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
