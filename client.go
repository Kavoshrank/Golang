package kavoshrank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client is the main entry point to the KavoshRank Go SDK.
type Client struct {
	cfg *config

	Articles *ArticleService
	Audio    *AudioService
	Copilot  *CopilotService
	Reviews  *ReviewService
	Events   *EventService
	Health   *HealthService
}

// NewClient initializes a new KavoshRank API client with functional options.
func NewClient(opts ...Option) *Client {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.httpClient == nil {
		cfg.httpClient = &http.Client{
			Timeout: cfg.timeout,
		}
	}

	c := &Client{cfg: cfg}
	c.Articles = &ArticleService{client: c}
	c.Audio = &AudioService{client: c}
	c.Copilot = &CopilotService{client: c}
	c.Reviews = &ReviewService{client: c}
	c.Events = &EventService{client: c}
	c.Health = &HealthService{client: c}

	return c
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return strings.TrimRight(c.cfg.baseURL, "/")
}

// ProjectID returns the configured project ID.
func (c *Client) ProjectID() int64 {
	return c.cfg.projectID
}

// execute performs the HTTP request and handles standard error mapping.
func (c *Client) execute(ctx context.Context, method, path string, query url.Values, body any, target any) error {
	fullURL := fmt.Sprintf("%s/%s", c.BaseURL(), strings.TrimLeft(path, "/"))
	if len(query) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, query.Encode())
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.cfg.siteToken != "" {
		req.Header.Set("X-KavoshRank-Secret-Token", c.cfg.siteToken)
		req.Header.Set("X-TSW-Secret-Token", c.cfg.siteToken)
	}

	if c.cfg.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.bearerToken)
	}

	for k, v := range c.cfg.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.cfg.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleHTTPError(resp.StatusCode, resp.Header, respBody)
	}

	if target != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, target); err != nil {
			return fmt.Errorf("unmarshal response (status %d): %w", resp.StatusCode, err)
		}
	}

	return nil
}

func (c *Client) handleHTTPError(statusCode int, header http.Header, body []byte) error {
	var errPayload struct {
		Message string            `json:"message"`
		Error   string            `json:"error"`
		Errors  map[string]string `json:"errors"`
	}

	msg := string(body)
	if err := json.Unmarshal(body, &errPayload); err == nil {
		if errPayload.Message != "" {
			msg = errPayload.Message
		} else if errPayload.Error != "" {
			msg = errPayload.Error
		}
	}

	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthenticationError{StatusCode: statusCode, Message: msg}
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return &ValidationError{StatusCode: statusCode, Message: msg, Errors: errPayload.Errors}
	case http.StatusTooManyRequests:
		retryAfter, _ := strconv.Atoi(header.Get("Retry-After"))
		return &RateLimitError{StatusCode: statusCode, Message: msg, RetryAfter: retryAfter}
	default:
		return &APIError{StatusCode: statusCode, Message: msg, Errors: errPayload.Errors}
	}
}
