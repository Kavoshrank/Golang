package kavoshrank

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// WebhookEvent represents an incoming event or action triggered by KavoshRank.
type WebhookEvent struct {
	Action  string         `json:"action"`
	Payload map[string]any `json:"payload"`
	Header  http.Header    `json:"-"`
}

// ActionHandler is a callback function for processing incoming webhook actions.
type ActionHandler func(event *WebhookEvent) (any, error)

// WebhookReceiver processes and verifies incoming webhook requests from KavoshRank.
type WebhookReceiver struct {
	secretToken string
	handlers    map[string]ActionHandler
}

// NewWebhookReceiver creates a new webhook receiver with the given secret token.
func NewWebhookReceiver(secretToken string) *WebhookReceiver {
	return &WebhookReceiver{
		secretToken: strings.TrimSpace(secretToken),
		handlers:    make(map[string]ActionHandler),
	}
}

// On registers a handler callback for an action (e.g. "publish_article", "audit", "ping").
func (r *WebhookReceiver) On(action string, handler ActionHandler) {
	r.handlers[action] = handler
}

// VerifyToken securely checks if the provided token matches the configured secret token.
func (r *WebhookReceiver) VerifyToken(token string) bool {
	if r.secretToken == "" || token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(r.secretToken), []byte(strings.TrimSpace(token))) == 1
}

// ExtractToken extracts the secret token from standard headers.
func (r *WebhookReceiver) ExtractToken(req *http.Request) string {
	candidates := []string{
		"X-KavoshRank-Secret-Token",
		"X-TSW-Secret-Token",
	}

	for _, header := range candidates {
		if val := req.Header.Get(header); val != "" {
			return val
		}
	}

	auth := req.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}

	return ""
}

// HTTPHandler returns an http.HandlerFunc that verifies and processes incoming webhooks.
func (r *WebhookReceiver) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := r.ExtractToken(req)
		if !r.VerifyToken(token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized: invalid or missing KavoshRank secret token",
			})
			return
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to read request body",
			})
			return
		}

		var payload map[string]any
		if len(bodyBytes) > 0 {
			_ = json.Unmarshal(bodyBytes, &payload)
		}
		if payload == nil {
			payload = make(map[string]any)
		}

		action := "ping"
		if a, ok := payload["action"].(string); ok && a != "" {
			action = a
		} else if e, ok := payload["event"].(string); ok && e != "" {
			action = e
		}

		event := &WebhookEvent{
			Action:  action,
			Payload: payload,
			Header:  req.Header,
		}

		w.Header().Set("Content-Type", "application/json")

		if handler, exists := r.handlers[action]; exists {
			res, err := handler(event)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"action":  action,
					"error":   err.Error(),
				})
				return
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"action":  action,
				"result":  res,
			})
			return
		}

		// Default ping handler
		if action == "ping" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"action":  "ping",
				"message": "pong",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"action":  action,
			"message": "no handler registered for action",
		})
	}
}
