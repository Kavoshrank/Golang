package kavoshrank

import (
	"context"
)

type HealthService struct {
	client *Client
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

type PluginUpdateInfo struct {
	Name            string            `json:"name"`
	Slug            string            `json:"slug"`
	Version         string            `json:"version"`
	DownloadURL     string            `json:"download_url"`
	Package         string            `json:"package"`
	RequiresPHP     string            `json:"requires_php,omitempty"`
	RequiresAtLeast string            `json:"requires_at_least,omitempty"`
	TestedUpTo      string            `json:"tested_up_to,omitempty"`
	LastUpdated     string            `json:"last_updated,omitempty"`
	Sections        map[string]string `json:"sections,omitempty"`
}

// Check checks KavoshRank system health.
func (s *HealthService) Check(ctx context.Context) (*HealthResponse, error) {
	var resp struct {
		Data *HealthResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", "/api/v1/health", nil, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat HealthResponse
	if err := s.client.execute(ctx, "GET", "/api/v1/health", nil, nil, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}

// GetUpdateInfo retrieves the latest published plugin/connector release record.
func (s *HealthService) GetUpdateInfo(ctx context.Context) (*PluginUpdateInfo, error) {
	var resp struct {
		Data *PluginUpdateInfo `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", "/api/v1/plugin/update-info", nil, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat PluginUpdateInfo
	if err := s.client.execute(ctx, "GET", "/api/v1/plugin/update-info", nil, nil, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}
