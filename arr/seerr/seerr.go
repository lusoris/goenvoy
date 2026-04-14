package seerr

import (
	"context"
	"fmt"
	"net/url"

	"github.com/golusoris/goenvoy/arr/v2"
)

// Client is a Seerr API client.
type Client struct {
	base *arr.BaseClient
}

// New creates a Seerr [Client] for the instance at baseURL.
func New(baseURL, apiKey string, opts ...arr.Option) (*Client, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{base: base}, nil
}

// GetStatus returns the current Seerr server status.
func (c *Client) GetStatus(ctx context.Context) (*StatusResponse, error) {
	var out StatusResponse
	if err := c.base.Get(ctx, "/api/v1/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMe returns the currently authenticated user.
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var out User
	if err := c.base.Get(ctx, "/api/v1/auth/me", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Search performs a multi-search for movies, TV shows, and people.
func (c *Client) Search(ctx context.Context, query string, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/search?query=%s&page=%d", url.QueryEscape(query), page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverMovies returns a paginated list of discovered movies.
func (c *Client) DiscoverMovies(ctx context.Context, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/discover/movies?page=%d", page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverTV returns a paginated list of discovered TV shows.
func (c *Client) DiscoverTV(ctx context.Context, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/discover/tv?page=%d", page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverTrending returns currently trending movies and TV shows.
func (c *Client) DiscoverTrending(ctx context.Context, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/discover/trending?page=%d", page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverUpcomingMovies returns upcoming movies.
func (c *Client) DiscoverUpcomingMovies(ctx context.Context, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/discover/movies/upcoming?page=%d", page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverUpcomingTV returns upcoming TV shows.
func (c *Client) DiscoverUpcomingTV(ctx context.Context, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/discover/tv/upcoming?page=%d", page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// pagedRequests is the list-requests response shape.
type pagedRequests struct {
	PageInfo PageInfo       `json:"pageInfo"`
	Results  []MediaRequest `json:"results,omitempty"`
}

// GetRequests returns all media requests with pagination and optional filter.
func (c *Client) GetRequests(ctx context.Context, take, skip int, filter string) ([]MediaRequest, *PageInfo, error) {
	var out pagedRequests
	path := fmt.Sprintf("/api/v1/request?take=%d&skip=%d", take, skip)
	if filter != "" {
		path += "&filter=" + url.QueryEscape(filter)
	}
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Results, &out.PageInfo, nil
}

// GetRequest returns a single media request by its ID.
func (c *Client) GetRequest(ctx context.Context, id int) (*MediaRequest, error) {
	var out MediaRequest
	path := fmt.Sprintf("/api/v1/request/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRequest creates a new media request.
func (c *Client) CreateRequest(ctx context.Context, body *CreateRequestBody) (*MediaRequest, error) {
	var out MediaRequest
	if err := c.base.Post(ctx, "/api/v1/request", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRequest removes a media request by its ID.
func (c *Client) DeleteRequest(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/request/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// ApproveRequest approves a pending request.
func (c *Client) ApproveRequest(ctx context.Context, id int) (*MediaRequest, error) {
	var out MediaRequest
	path := fmt.Sprintf("/api/v1/request/%d/approve", id)
	if err := c.base.Post(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeclineRequest declines a pending request.
func (c *Client) DeclineRequest(ctx context.Context, id int) (*MediaRequest, error) {
	var out MediaRequest
	path := fmt.Sprintf("/api/v1/request/%d/decline", id)
	if err := c.base.Post(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RetryRequest retries a failed request.
func (c *Client) RetryRequest(ctx context.Context, id int) (*MediaRequest, error) {
	var out MediaRequest
	path := fmt.Sprintf("/api/v1/request/%d/retry", id)
	if err := c.base.Post(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRequestCount returns aggregate request count statistics.
func (c *Client) GetRequestCount(ctx context.Context) (*RequestCount, error) {
	var out RequestCount
	if err := c.base.Get(ctx, "/api/v1/request/count", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovie returns full details for a movie by its TMDB ID.
func (c *Client) GetMovie(ctx context.Context, tmdbID int) (*MovieDetails, error) {
	var out MovieDetails
	path := fmt.Sprintf("/api/v1/movie/%d", tmdbID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieRecommendations returns recommended movies for a given TMDB ID.
func (c *Client) GetMovieRecommendations(ctx context.Context, tmdbID, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/movie/%d/recommendations?page=%d", tmdbID, page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieSimilar returns similar movies for a given TMDB ID.
func (c *Client) GetMovieSimilar(ctx context.Context, tmdbID, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/movie/%d/similar?page=%d", tmdbID, page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTV returns full details for a TV show by its TMDB ID.
func (c *Client) GetTV(ctx context.Context, tmdbID int) (*TvDetails, error) {
	var out TvDetails
	path := fmt.Sprintf("/api/v1/tv/%d", tmdbID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeason returns details and episodes for a specific season.
func (c *Client) GetTVSeason(ctx context.Context, tmdbID, seasonNumber int) (*Season, error) {
	var out Season
	path := fmt.Sprintf("/api/v1/tv/%d/season/%d", tmdbID, seasonNumber)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVRecommendations returns recommended TV shows for a given TMDB ID.
func (c *Client) GetTVRecommendations(ctx context.Context, tmdbID, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/tv/%d/recommendations?page=%d", tmdbID, page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSimilar returns similar TV shows for a given TMDB ID.
func (c *Client) GetTVSimilar(ctx context.Context, tmdbID, page int) (*SearchResults, error) {
	var out SearchResults
	path := fmt.Sprintf("/api/v1/tv/%d/similar?page=%d", tmdbID, page)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// pagedMedia is the list-media response shape.
type pagedMedia struct {
	PageInfo PageInfo    `json:"pageInfo"`
	Results  []MediaInfo `json:"results,omitempty"`
}

// GetMedia returns a paginated list of all media items.
func (c *Client) GetMedia(ctx context.Context, take, skip int) ([]MediaInfo, *PageInfo, error) {
	var out pagedMedia
	path := fmt.Sprintf("/api/v1/media?take=%d&skip=%d", take, skip)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Results, &out.PageInfo, nil
}

// DeleteMedia removes a media item by its ID.
func (c *Client) DeleteMedia(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/media/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// UpdateMediaStatus updates a media item's availability status.
func (c *Client) UpdateMediaStatus(ctx context.Context, id int, status string, is4k bool) (*MediaInfo, error) {
	var out MediaInfo
	path := fmt.Sprintf("/api/v1/media/%d/%s", id, status)
	body := map[string]bool{"is4k": is4k}
	if err := c.base.Post(ctx, path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// pagedUsers is the list-users response shape.
type pagedUsers struct {
	PageInfo PageInfo `json:"pageInfo"`
	Results  []User   `json:"results,omitempty"`
}

// GetUsers returns a paginated list of users.
func (c *Client) GetUsers(ctx context.Context, take, skip int) ([]User, *PageInfo, error) {
	var out pagedUsers
	path := fmt.Sprintf("/api/v1/user?take=%d&skip=%d", take, skip)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Results, &out.PageInfo, nil
}

// GetUser returns a single user by their ID.
func (c *Client) GetUser(ctx context.Context, id int) (*User, error) {
	var out User
	path := fmt.Sprintf("/api/v1/user/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteUser removes a user by their ID.
func (c *Client) DeleteUser(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/user/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetUserQuota returns the request quota for a specific user.
func (c *Client) GetUserQuota(ctx context.Context, userID int) (*UserQuota, error) {
	var out UserQuota
	path := fmt.Sprintf("/api/v1/user/%d/quota", userID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// pagedIssues is the list-issues response shape.
type pagedIssues struct {
	PageInfo PageInfo `json:"pageInfo"`
	Results  []Issue  `json:"results,omitempty"`
}

// GetIssues returns a paginated list of issues.
func (c *Client) GetIssues(ctx context.Context, take, skip int, filter string) ([]Issue, *PageInfo, error) {
	var out pagedIssues
	path := fmt.Sprintf("/api/v1/issue?take=%d&skip=%d", take, skip)
	if filter != "" {
		path += "&filter=" + url.QueryEscape(filter)
	}
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Results, &out.PageInfo, nil
}

// GetIssue returns a single issue by its ID.
func (c *Client) GetIssue(ctx context.Context, id int) (*Issue, error) {
	var out Issue
	path := fmt.Sprintf("/api/v1/issue/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateIssue creates a new issue report.
func (c *Client) CreateIssue(ctx context.Context, body *CreateIssueBody) (*Issue, error) {
	var out Issue
	if err := c.base.Post(ctx, "/api/v1/issue", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteIssue removes an issue by its ID.
func (c *Client) DeleteIssue(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/issue/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// AddIssueComment adds a comment to an existing issue.
func (c *Client) AddIssueComment(ctx context.Context, issueID int, message string) (*Issue, error) {
	var out Issue
	path := fmt.Sprintf("/api/v1/issue/%d/comment", issueID)
	body := map[string]string{"message": message}
	if err := c.base.Post(ctx, path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResolveIssue marks an issue as resolved.
func (c *Client) ResolveIssue(ctx context.Context, id int) (*Issue, error) {
	var out Issue
	path := fmt.Sprintf("/api/v1/issue/%d/resolved", id)
	if err := c.base.Post(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReopenIssue reopens a resolved issue.
func (c *Client) ReopenIssue(ctx context.Context, id int) (*Issue, error) {
	var out Issue
	path := fmt.Sprintf("/api/v1/issue/%d/open", id)
	if err := c.base.Post(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIssueCount returns aggregate issue count statistics.
func (c *Client) GetIssueCount(ctx context.Context) (*IssueCount, error) {
	var out IssueCount
	if err := c.base.Get(ctx, "/api/v1/issue/count", &out); err != nil {
		return nil, err
	}
	return &out, nil
}
