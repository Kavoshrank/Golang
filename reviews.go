package kavoshrank

import (
	"context"
	"fmt"
)

type ReviewService struct {
	client *Client
}

type RawReview struct {
	ID        int64  `json:"id"`
	Rating    int    `json:"rating,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at,omitempty"`
}

type ReviewDigestRequest struct {
	ProjectID    int64       `json:"project_id,omitempty"`
	ProductID    int64       `json:"product_id"`
	ProductTitle string      `json:"product_title,omitempty"`
	Reviews      []RawReview `json:"reviews"`
	ForceRefresh bool        `json:"force_refresh,omitempty"`
}

type DigestPoint struct {
	Text      string  `json:"text"`
	ReviewIDs []int64 `json:"review_ids,omitempty"`
}

type ReviewDigestResponse struct {
	ProductID            int64         `json:"product_id"`
	Status               string        `json:"status"` // "ok" or "insufficient_reviews"
	TotalReviewsAnalyzed int           `json:"total_reviews_analyzed"`
	Consensus            string        `json:"consensus,omitempty"` // "positive", "mixed", "negative"
	Summary              string        `json:"summary,omitempty"`
	Pros                 []DigestPoint `json:"pros,omitempty"`
	Cons                 []DigestPoint `json:"cons,omitempty"`
}

// GenerateDigest generates an automated AI review sentiment digest.
func (s *ReviewService) GenerateDigest(ctx context.Context, req *ReviewDigestRequest) (*ReviewDigestResponse, error) {
	path := "/api/v1/plugin/reviews/digest"
	if s.client.ProjectID() > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/woo/review-digest", s.client.ProjectID())
	}

	var resp struct {
		Data *ReviewDigestResponse `json:"data"`
	}
	if err := s.client.execute(ctx, "POST", path, nil, req, &resp); err != nil {
		return nil, err
	}
	if resp.Data != nil {
		return resp.Data, nil
	}

	var flat ReviewDigestResponse
	if err := s.client.execute(ctx, "POST", path, nil, req, &flat); err != nil {
		return nil, err
	}
	return &flat, nil
}
