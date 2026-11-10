package kavoshrank

import (
	"context"
	"fmt"
)

type CopilotService struct {
	client *Client
}

type CopilotConsultRequest struct {
	ProjectID       int64            `json:"project_id,omitempty"`
	Question        string           `json:"question"`
	CurrentProduct  map[string]any   `json:"current_product,omitempty"`
	RelatedProducts []map[string]any `json:"related_products,omitempty"`
	ShopperID       string           `json:"shopper_id,omitempty"`
}

type CopilotConsultResponse struct {
	Reply           string           `json:"reply"`
	Recommendations []map[string]any `json:"recommendations,omitempty"`
	Confidence      float64          `json:"confidence,omitempty"`
}

type CopilotSettings struct {
	IsEnabled        bool   `json:"is_enabled"`
	WidgetPosition   string `json:"widget_position,omitempty"`
	PrimaryColor     string `json:"primary_color,omitempty"`
	WelcomeMessage   string `json:"welcome_message,omitempty"`
	DailyCapRequests int    `json:"daily_cap_requests,omitempty"`
}

// Consult asks the Storefront AI Copilot for customer pre-purchase guidance.
func (s *CopilotService) Consult(ctx context.Context, req *CopilotConsultRequest) (*CopilotConsultResponse, error) {
	path := "/api/v1/plugin/copilot/consult"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/copilot/consult", s.client.ProjectID())
	}

	var resp struct {
		Data *CopilotConsultResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "POST", path, nil, req, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat CopilotConsultResponse
	if err := s.client.execute(ctx, "POST", path, nil, req, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}

// GetSettings retrieves current Storefront Copilot settings.
func (s *CopilotService) GetSettings(ctx context.Context) (*CopilotSettings, error) {
	path := "/api/v1/plugin/copilot/settings"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/copilot/settings", s.client.ProjectID())
	}

	var resp struct {
		Data *CopilotSettings `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", path, nil, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat CopilotSettings
	if err := s.client.execute(ctx, "GET", path, nil, nil, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}
