package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	kavoshrank "github.com/Kavoshrank/Golang"
)

func main() {
	// 1. Initialize Client
	client := kavoshrank.NewClient(
		kavoshrank.WithSiteToken("YOUR_SITE_SECRET_TOKEN"),
		kavoshrank.WithBaseURL(kavoshrank.DefaultBaseURL),
		kavoshrank.WithTimeout(30*time.Second),
	)

	ctx := context.Background()

	// 2. Generate an AI SEO Article
	article, err := client.Articles.Generate(ctx, &kavoshrank.GenerateArticleRequest{
		Topic:     "Modern Cloud Native Go Development in 2026",
		WordCount: 1500,
		UseSERP:   true,
	})
	if err != nil {
		log.Printf("article generation error: %v", err)
	} else {
		fmt.Printf("Generated article ID: %d, Title: %s\n", article.ID, article.Title)
	}

	// 3. Generate Audio Narration / Podcast
	audio, err := client.Audio.Generate(ctx, &kavoshrank.GenerateAudioRequest{
		PostID:  101,
		Title:   "Modern Cloud Native Go Development in 2026",
		Content: "Full article narration content...",
		Voice:   "Puck",
		Theme:   "glassmorphic",
	})
	if err != nil {
		log.Printf("audio generation error: %v", err)
	} else {
		fmt.Printf("Podcast audio URL: %s\n", audio.AudioURL)
	}

	// 4. Consult AI Storefront Copilot
	reply, err := client.Copilot.Consult(ctx, &kavoshrank.CopilotConsultRequest{
		Question: "Which plan supports multi-region deployment?",
	})
	if err != nil {
		log.Printf("copilot error: %v", err)
	} else {
		fmt.Printf("Copilot response: %s\n", reply.Reply)
	}

	// 5. Incoming Webhook Receiver Example
	receiver := kavoshrank.NewWebhookReceiver("YOUR_SITE_SECRET_TOKEN")
	receiver.On("publish_article", func(event *kavoshrank.WebhookEvent) (any, error) {
		fmt.Printf("Received publish command for article: %v\n", event.Payload["title"])
		// Save to your local database or service
		return map[string]any{"published": true}, nil
	})

	http.HandleFunc("/kavoshrank/webhook", receiver.HTTPHandler())
	fmt.Println("Webhook receiver ready at :8080/kavoshrank/webhook")
	// _ = http.ListenAndServe(":8080", nil)
}
