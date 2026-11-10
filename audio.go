package kavoshrank

import (
	"context"
	"fmt"
	"net/url"
)

type AudioService struct {
	client *Client
}

type GenerateAudioRequest struct {
	ProjectID  int64  `json:"project_id,omitempty"`
	PostID     int    `json:"post_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Voice      string `json:"voice,omitempty"`
	Theme      string `json:"theme,omitempty"`
	AutoInject bool   `json:"auto_inject,omitempty"`
	Async      bool   `json:"async,omitempty"`
}

type GenerateAudioResponse struct {
	AudioURL        string `json:"audio_url"`
	PostID          int    `json:"post_id"`
	AudioMinutes    int    `json:"audio_minutes,omitempty"`
	TokensDeducted  int64  `json:"tokens_deducted,omitempty"`
	RemainingTokens int64  `json:"remaining_tokens,omitempty"`
}

type AudioArticleItem struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	AudioURL  string `json:"audio_url,omitempty"`
	HasAudio  bool   `json:"has_audio"`
	Duration  int    `json:"duration,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type AudioSettings struct {
	PlayerTheme      string  `json:"player_theme"`
	SpeakerVoice     string  `json:"speaker_voice"`
	SpeechRate       float64 `json:"speech_rate"`
	AutoInsertPlayer bool    `json:"auto_insert_player"`
}

// Generate creates audio narration/podcast from article content.
func (s *AudioService) Generate(ctx context.Context, req *GenerateAudioRequest) (*GenerateAudioResponse, error) {
	path := "/api/v1/audio/generate"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/audio/generate", s.client.ProjectID())
	}

	var resp struct {
		Data *GenerateAudioResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "POST", path, nil, req, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat GenerateAudioResponse
	if err := s.client.execute(ctx, "POST", path, nil, req, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}

// ListArticles returns articles eligible for audio conversion.
func (s *AudioService) ListArticles(ctx context.Context, page, limit int) ([]AudioArticleItem, error) {
	path := "/api/v1/audio/articles"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/audio/articles", s.client.ProjectID())
	}

	query := url.Values{}
	if page > 0 {
		query.Set("page", fmt.Sprintf("%d", page))
	}
	if limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", limit))
	}

	var resp struct {
		Data []AudioArticleItem `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", path, query, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
