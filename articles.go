package kavoshrank

import (
	"context"
	"fmt"
	"net/url"
)

type ArticleService struct {
	client *Client
}

type GenerateArticleRequest struct {
	ProjectID           int64           `json:"project_id,omitempty"`
	Topic               string          `json:"topic"`
	WordCount           int             `json:"word_count,omitempty"`
	TargetWordCount     int             `json:"target_word_count,omitempty"`
	CategoryID          int             `json:"category_id,omitempty"`
	CategoryName        string          `json:"category_name,omitempty"`
	Model               string          `json:"model,omitempty"`
	Status              string          `json:"status,omitempty"`
	PublishStatus       string          `json:"publish_status,omitempty"`
	UseSERP             bool            `json:"use_serp,omitempty"`
	IncludeVideo        bool            `json:"include_video,omitempty"`
	ExtraInstructions   string          `json:"extra_instructions,omitempty"`
	FeaturedImageType   string          `json:"featured_image_type,omitempty"`
	FeaturedImagePrompt string          `json:"featured_image_prompt,omitempty"`
	ShareSocial         bool            `json:"share_social,omitempty"`
	SocialChannels      map[string]bool `json:"social_channels,omitempty"`
	Async               bool            `json:"async,omitempty"`
}

type ArticleResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Excerpt   string `json:"excerpt,omitempty"`
	MetaDesc  string `json:"meta_description,omitempty"`
	Status    string `json:"status,omitempty"`
	JobID     string `json:"job_id,omitempty"`
	AudioURL  string `json:"audio_url,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type QueuedJobResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type JobStatusResponse struct {
	JobID    string `json:"job_id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
}

type ActiveJobResponse struct {
	HasActiveJob bool   `json:"has_active_job"`
	JobID        string `json:"job_id,omitempty"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
}

type CategoryItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// Generate generates an AI SEO article.
func (s *ArticleService) Generate(ctx context.Context, req *GenerateArticleRequest) (*ArticleResponse, error) {
	path := "/api/v1/articles/generate"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/generate", s.client.ProjectID())
	}

	var resp struct {
		Data *ArticleResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "POST", path, nil, req, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	// In case response is flat
	var flat ArticleResponse
	if err := s.client.execute(ctx, "POST", path, nil, req, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}

// GenerateAsync queues an asynchronous article generation job.
func (s *ArticleService) GenerateAsync(ctx context.Context, req *GenerateArticleRequest) (*QueuedJobResponse, error) {
	req.Async = true
	path := "/api/v1/articles/generate"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/generate", s.client.ProjectID())
	}

	var resp struct {
		Data *QueuedJobResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "POST", path, nil, req, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat QueuedJobResponse
	if err := s.client.execute(ctx, "POST", path, nil, req, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}

// GetJob retrieves the current status of an asynchronous generation job.
func (s *ArticleService) GetJob(ctx context.Context, jobID string) (*JobStatusResponse, error) {
	path := fmt.Sprintf("/api/v1/articles/jobs/%s", jobID)
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/jobs/%s", s.client.ProjectID(), jobID)
	}

	var resp struct {
		Data *JobStatusResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", path, nil, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat JobStatusResponse
	if err := s.client.execute(ctx, "GET", path, nil, nil, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}

// GetActiveJob checks whether there is an active job running in the project.
func (s *ArticleService) GetActiveJob(ctx context.Context) (*ActiveJobResponse, error) {
	path := "/api/v1/articles/active-job"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/active-job", s.client.ProjectID())
	}

	var resp struct {
		Data *ActiveJobResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", path, nil, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat ActiveJobResponse
	if err := s.client.execute(ctx, "GET", path, nil, nil, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}

// DismissActiveJob dismisses the ongoing active job.
func (s *ArticleService) DismissActiveJob(ctx context.Context) error {
	path := "/api/v1/articles/active-job/dismiss"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/active-job/dismiss", s.client.ProjectID())
	}

	return s.client.execute(ctx, "POST", path, nil, nil, nil)
}

// Publish publishes an article by ID.
func (s *ArticleService) Publish(ctx context.Context, articleID int64, params map[string]any) error {
	path := fmt.Sprintf("/api/v1/articles/%d/publish", articleID)
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/%d/publish", s.client.ProjectID(), articleID)
	}

	return s.client.execute(ctx, "POST", path, nil, params, nil)
}

// SuggestTitles suggests SEO titles for a given keyword.
func (s *ArticleService) SuggestTitles(ctx context.Context, keyword string) ([]string, error) {
	path := "/api/v1/articles/title-suggestions"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/title-suggestions", s.client.ProjectID())
	}

	var resp struct {
		Data []string `json:"data"`
	}
	if err := s.client.execute(ctx, "POST", path, nil, map[string]string{"keyword": keyword}, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GoogleSuggest queries autocomplete search suggestions from Google.
func (s *ArticleService) GoogleSuggest(ctx context.Context, keyword string) ([]string, error) {
	path := "/api/v1/articles/google-suggest"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/google-suggest", s.client.ProjectID())
	}

	query := url.Values{}
	query.Set("keyword", keyword)

	var resp struct {
		Data []string `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", path, query, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Categories returns available article categories.
func (s *ArticleService) Categories(ctx context.Context) ([]CategoryItem, error) {
	path := "/api/v1/articles/categories"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/articles/categories", s.client.ProjectID())
	}

	var resp struct {
		Data []CategoryItem `json:"data"`
	}
	if err := s.client.execute(ctx, "GET", path, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
