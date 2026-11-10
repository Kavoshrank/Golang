package kavoshrank

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_Options(t *testing.T) {
	client := NewClient(
		WithBaseURL("https://api.example.com"),
		WithSiteToken("site-token-xyz"),
		WithBearerToken("bearer-jwt-123"),
		WithProjectID(77),
		WithTimeout(15*time.Second),
		WithHeader("X-Custom", "val"),
	)

	if client.BaseURL() != "https://api.example.com" {
		t.Fatalf("unexpected BaseURL: %s", client.BaseURL())
	}
	if client.ProjectID() != 77 {
		t.Fatalf("unexpected ProjectID: %d", client.ProjectID())
	}
	if client.cfg.siteToken != "site-token-xyz" {
		t.Fatalf("unexpected siteToken: %s", client.cfg.siteToken)
	}
	if client.cfg.bearerToken != "bearer-jwt-123" {
		t.Fatalf("unexpected bearerToken: %s", client.cfg.bearerToken)
	}
	if client.cfg.headers["X-Custom"] != "val" {
		t.Fatalf("unexpected custom header: %s", client.cfg.headers["X-Custom"])
	}
}

func TestArticleService_Generate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/projects/5/articles/generate" {
			t.Errorf("expected path /api/v1/projects/5/articles/generate, got %s", r.URL.Path)
		}
		if r.Header.Get("X-KavoshRank-Secret-Token") != "my-site-token" {
			t.Errorf("missing site token header")
		}

		var req GenerateArticleRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Topic != "Golang Microservices" {
			t.Errorf("unexpected topic: %s", req.Topic)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": ArticleResponse{
				ID:    101,
				Title: "Golang Microservices Guide",
			},
		})
	}))
	defer ts.Close()

	client := NewClient(
		WithBaseURL(ts.URL),
		WithSiteToken("my-site-token"),
		WithProjectID(5),
	)

	resp, err := client.Articles.Generate(context.Background(), &GenerateArticleRequest{
		Topic:     "Golang Microservices",
		WordCount: 1500,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != 101 || resp.Title != "Golang Microservices Guide" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestAudioService_Generate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": GenerateAudioResponse{
				AudioURL: "https://cdn.example.com/podcast.mp3",
				PostID:   42,
			},
		})
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL))
	resp, err := client.Audio.Generate(context.Background(), &GenerateAudioRequest{
		PostID:  42,
		Title:   "Audio Post",
		Content: "Full text content",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AudioURL != "https://cdn.example.com/podcast.mp3" {
		t.Fatalf("unexpected audio url: %s", resp.AudioURL)
	}
}

func TestCopilotService_Consult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": CopilotConsultResponse{
				Reply: "The Pro model includes wireless charging.",
			},
		})
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL))
	resp, err := client.Copilot.Consult(context.Background(), &CopilotConsultRequest{
		Question: "Does this include wireless charging?",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Reply, "wireless charging") {
		t.Fatalf("unexpected reply: %s", resp.Reply)
	}
}

func TestReviewService_GenerateDigest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": ReviewDigestResponse{
				ProductID: 88,
				Status:    "ok",
				Consensus: "positive",
			},
		})
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL))
	resp, err := client.Reviews.GenerateDigest(context.Background(), &ReviewDigestRequest{
		ProductID: 88,
		Reviews: []RawReview{
			{ID: 1, Content: "Great product", Rating: 5},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "ok" || resp.Consensus != "positive" {
		t.Fatalf("unexpected digest: %+v", resp)
	}
}

func TestErrors_Handling(t *testing.T) {
	// 401 error
	ts401 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid site token"})
	}))
	defer ts401.Close()

	client401 := NewClient(WithBaseURL(ts401.URL))
	_, err := client401.Health.Check(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*AuthenticationError); !ok {
		t.Fatalf("expected *AuthenticationError, got %T: %v", err, err)
	}

	// 400 error
	ts400 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "validation error",
			"errors":  map[string]string{"topic": "required"},
		})
	}))
	defer ts400.Close()

	client400 := NewClient(WithBaseURL(ts400.URL))
	_, err = client400.Articles.Generate(context.Background(), &GenerateArticleRequest{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	vErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
	if vErr.Errors["topic"] != "required" {
		t.Fatalf("unexpected validation field: %v", vErr.Errors)
	}

	// 429 error
	ts429 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "rate limited"})
	}))
	defer ts429.Close()

	client429 := NewClient(WithBaseURL(ts429.URL))
	_, err = client429.Health.Check(context.Background())
	if err == nil {
		t.Fatal("expected rate limit error, got nil")
	}
	rlErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("expected *RateLimitError, got %T: %v", err, err)
	}
	if rlErr.RetryAfter != 30 {
		t.Fatalf("expected retry after 30, got %d", rlErr.RetryAfter)
	}
}

func TestWebhookReceiver(t *testing.T) {
	receiver := NewWebhookReceiver("secret-token-123")

	if !receiver.VerifyToken("secret-token-123") {
		t.Fatal("expected token to verify")
	}
	if receiver.VerifyToken("wrong-token") {
		t.Fatal("expected wrong token to fail")
	}

	dispatched := false
	receiver.On("publish_article", func(event *WebhookEvent) (any, error) {
		dispatched = true
		return map[string]any{"id": event.Payload["article_id"]}, nil
	})

	handler := receiver.HTTPHandler()

	// Test unauthorized call
	reqInvalid := httptest.NewRequest("POST", "/webhook", strings.NewReader(`{"action":"ping"}`))
	reqInvalid.Header.Set("X-KavoshRank-Secret-Token", "wrong-token")
	recInvalid := httptest.NewRecorder()
	handler(recInvalid, reqInvalid)
	if recInvalid.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recInvalid.Code)
	}

	// Test authorized dispatch
	reqValid := httptest.NewRequest("POST", "/webhook", strings.NewReader(`{"action":"publish_article","article_id":456}`))
	reqValid.Header.Set("X-KavoshRank-Secret-Token", "secret-token-123")
	recValid := httptest.NewRecorder()
	handler(recValid, reqValid)

	if recValid.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recValid.Code)
	}
	if !dispatched {
		t.Fatal("expected webhook action handler to be dispatched")
	}
}
