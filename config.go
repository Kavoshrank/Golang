package kavoshrank

import (
	"net/http"
	"time"
)

const (
	DefaultBaseURL = "https://backend.kavoshrank.ir"
	DefaultTimeout = 30 * time.Second
)

// Option is a functional option for configuring the Client.
type Option func(*config)

type config struct {
	baseURL     string
	siteToken   string
	bearerToken string
	projectID   int64
	httpClient  *http.Client
	timeout     time.Duration
	headers     map[string]string
}

func defaultConfig() *config {
	return &config{
		baseURL: DefaultBaseURL,
		timeout: DefaultTimeout,
		headers: make(map[string]string),
	}
}

// WithBaseURL sets the base URL for the KavoshRank API.
func WithBaseURL(url string) Option {
	return func(c *config) {
		if url != "" {
			c.baseURL = url
		}
	}
}

// WithSiteToken sets the site secret token used by WordPress and custom CMS connectors.
func WithSiteToken(token string) Option {
	return func(c *config) {
		c.siteToken = token
	}
}

// WithBearerToken sets the JWT or API token for user and project API calls.
func WithBearerToken(token string) Option {
	return func(c *config) {
		c.bearerToken = token
	}
}

// WithProjectID sets the default project context ID for project-scoped routes.
func WithProjectID(id int64) Option {
	return func(c *config) {
		c.projectID = id
	}
}

// WithHTTPClient allows passing a custom net/http Client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		if d > 0 {
			c.timeout = d
		}
	}
}

// WithHeader adds a custom HTTP header to all outbound requests.
func WithHeader(key, value string) Option {
	return func(c *config) {
		c.headers[key] = value
	}
}
