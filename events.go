package kavoshrank

import (
	"context"
)

type EventService struct {
	client *Client
}

type InboundEventRequest struct {
	Event   string         `json:"event"`
	Payload map[string]any `json:"payload"`
}

type InboundEventResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// Send dispatches an inbound site event to KavoshRank (leads, comments, forms, orders).
func (s *EventService) Send(ctx context.Context, eventType string, payload map[string]any) (*InboundEventResponse, error) {
	req := InboundEventRequest{
		Event:   eventType,
		Payload: payload,
	}

	var resp struct {
		Data *InboundEventResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "POST", "/api/v1/plugin/events", nil, req, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat InboundEventResponse
	if err := s.client.execute(ctx, "POST", "/api/v1/plugin/events", nil, req, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}
