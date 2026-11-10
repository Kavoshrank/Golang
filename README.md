# KavoshRank Go SDK

Official Go SDK for integrating with the KavoshRank AI SEO, Content Automation, Audio Studio, and Storefront Copilot Master Engine.

## Requirements

- Go 1.21 or later

## Installation

```bash
go get github.com/Kavoshrank/Golang
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kavoshrank "github.com/Kavoshrank/Golang"
)

func main() {
	client := kavoshrank.NewClient(
		kavoshrank.WithSiteToken("YOUR_SITE_SECRET_TOKEN"), // or WithBearerToken("JWT")
		kavoshrank.WithBaseURL(kavoshrank.DefaultBaseURL),
		kavoshrank.WithProjectID(10), // optional default project ID
		kavoshrank.WithTimeout(30*time.Second),
	)

	ctx := context.Background()

	// 1. Generate an AI SEO Article
	article, err := client.Articles.Generate(ctx, &kavoshrank.GenerateArticleRequest{
		Topic:     "Kubernetes Production Best Practices",
		WordCount: 1500,
		UseSERP:   true,
	})
	if err != nil {
		log.Fatalf("generate failed: %v", err)
	}
	fmt.Printf("Article ID: %d, Title: %s\n", article.ID, article.Title)
}
```

### 1. Articles & SEO Generation

```go
// Asynchronous generation
queued, err := client.Articles.GenerateAsync(ctx, &kavoshrank.GenerateArticleRequest{
    Topic: "PostgreSQL High Availability",
})

// Query async job status
job, err := client.Articles.GetJob(ctx, queued.JobID)

// Title suggestions
titles, err := client.Articles.SuggestTitles(ctx, "postgresql ha")

// Google autocomplete suggestions
suggestions, err := client.Articles.GoogleSuggest(ctx, "postgresql ha")

// Publish an article
err = client.Articles.Publish(ctx, article.ID, nil)
```

### 2. Audio Studio & Narration

```go
audio, err := client.Audio.Generate(ctx, &kavoshrank.GenerateAudioRequest{
    PostID:  101,
    Title:   "Kubernetes Production Best Practices",
    Content: "Narration script content for the article...",
    Voice:   "Puck",
    Theme:   "glassmorphic", // glassmorphic, minimal, waveform
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Audio URL:", audio.AudioURL)
```

### 3. Storefront AI Copilot

```go
consult, err := client.Copilot.Consult(ctx, &kavoshrank.CopilotConsultRequest{
    Question: "Is this graphics card compatible with a 650W power supply?",
    CurrentProduct: map[string]any{
        "id":    55,
        "title": "GeForce RTX 4070 Dual Fan",
        "price": 599,
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Assistant reply:", consult.Reply)
```

### 4. Review Sentiment Digest

```go
digest, err := client.Reviews.GenerateDigest(ctx, &kavoshrank.ReviewDigestRequest{
    ProductID: 55,
    Reviews: []kavoshrank.RawReview{
        {ID: 1, Rating: 5, Content: "Extremely quiet and runs cold."},
        {ID: 2, Rating: 4, Content: "Great efficiency, slightly large for ITX cases."},
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Consensus:", digest.Consensus)
fmt.Println("Summary:", digest.Summary)
```

### 5. Inbound Site Events

```go
_, err := client.Events.Send(ctx, "order_placed", map[string]any{
    "order_id": 10024,
    "total":    149.99,
})
```

---

## Receiving Webhooks in Go

Use `WebhookReceiver` to handle and verify incoming actions from the KavoshRank Master Engine:

```go
package main

import (
	"fmt"
	"net/http"

	kavoshrank "github.com/Kavoshrank/Golang"
)

func main() {
	receiver := kavoshrank.NewWebhookReceiver("YOUR_SITE_SECRET_TOKEN")

	receiver.On("publish_article", func(event *kavoshrank.WebhookEvent) (any, error) {
		title := event.Payload["title"]
		fmt.Printf("Received publish command for: %v\n", title)

		// Save article to your application database
		return map[string]any{"status": "published"}, nil
	})

	receiver.On("audit", func(event *kavoshrank.WebhookEvent) (any, error) {
		return map[string]any{"healthy": true}, nil
	})

	http.HandleFunc("/api/kavoshrank/webhook", receiver.HTTPHandler())
	http.ListenAndServe(":8080", nil)
}
```

---

## Error Handling

The SDK maps HTTP responses into typed Go errors:

- `*kavoshrank.AuthenticationError`: 401 Unauthorized / 403 Forbidden.
- `*kavoshrank.ValidationError`: 400 Bad Request / 422 Unprocessable (check `.Errors` field map).
- `*kavoshrank.RateLimitError`: 429 Too Many Requests (check `.RetryAfter` in seconds).
- `*kavoshrank.APIError`: Other HTTP errors.

```go
_, err := client.Articles.Generate(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *kavoshrank.ValidationError:
        fmt.Printf("Validation error: %s (fields: %v)\n", e.Message, e.Errors)
    case *kavoshrank.RateLimitError:
        fmt.Printf("Rate limited. Retry after %d seconds\n", e.RetryAfter)
    case *kavoshrank.AuthenticationError:
        fmt.Printf("Authentication failed: %s\n", e.Message)
    default:
        fmt.Printf("API error: %v\n", err)
    }
}
```

## Running Tests

```bash
cd packages/go
go test -v ./...
```
