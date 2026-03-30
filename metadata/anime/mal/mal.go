package mal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://api.myanimelist.net/v2"
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "goenvoy/0.0.1"
)

// Option configures a [Client].
type Option func(*Client)

// WithHTTPClient sets a custom [http.Client].
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithTimeout overrides the default HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(cl *Client) { cl.httpClient.Timeout = d }
}

// WithUserAgent sets the User-Agent header for all requests.
func WithUserAgent(ua string) Option {
	return func(cl *Client) { cl.userAgent = ua }
}

// WithBaseURL overrides the default API base URL.
func WithBaseURL(u string) Option {
	return func(cl *Client) { cl.rawBaseURL = u }
}

// Client is a MyAnimeList API v2 client.
type Client struct {
	clientID   string
	rawBaseURL string
	httpClient *http.Client
	userAgent  string
}

// New creates a [Client] using the given MAL API client ID.
func New(clientID string, opts ...Option) *Client {
	c := &Client{
		clientID:   clientID,
		rawBaseURL: defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		userAgent:  defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Err        string `json:"error"`
	Message    string `json:"message"`
	// RawBody holds the raw response body when the error response could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("mal: HTTP %d: %s: %s", e.StatusCode, e.Err, e.Message)
	}
	if e.Err != "" {
		return fmt.Sprintf("mal: HTTP %d: %s", e.StatusCode, e.Err)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("mal: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("mal: HTTP %d", e.StatusCode)
}

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	u, err := url.Parse(c.rawBaseURL + path)
	if err != nil {
		return fmt.Errorf("mal: parse URL: %w", err)
	}

	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("mal: create request: %w", err)
	}

	req.Header.Set("X-MAL-CLIENT-ID", c.clientID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mal: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mal: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.RawBody = string(body)
		}
		return apiErr
	}

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("mal: decode response: %w", err)
		}
	}
	return nil
}

// fieldsParam returns a url.Values with the fields parameter set from the given list.
func fieldsParam(fields []string) url.Values {
	v := url.Values{}
	if len(fields) > 0 {
		v.Set("fields", strings.Join(fields, ","))
	}
	return v
}

// paginatedParams builds url.Values with fields, limit, and offset.
func paginatedParams(fields []string, limit, offset int) url.Values {
	v := fieldsParam(fields)
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		v.Set("offset", strconv.Itoa(offset))
	}
	return v
}

// GetAnime returns details for the anime with the given ID.
func (c *Client) GetAnime(ctx context.Context, animeID int, fields []string) (*Anime, error) {
	var out Anime
	path := fmt.Sprintf("/anime/%d", animeID)
	if err := c.get(ctx, path, fieldsParam(fields), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchAnime searches for anime by query string.
func (c *Client) SearchAnime(ctx context.Context, query string, fields []string, limit, offset int) ([]Anime, *Paging, error) {
	p := paginatedParams(fields, limit, offset)
	p.Set("q", query)
	var out animeListResponse
	if err := c.get(ctx, "/anime", p, &out); err != nil {
		return nil, nil, err
	}
	anime := make([]Anime, len(out.Data))
	for i := range out.Data {
		anime[i] = out.Data[i].Node
	}
	return anime, &out.Paging, nil
}

// AnimeRanking returns anime ranked by the given ranking type.
func (c *Client) AnimeRanking(ctx context.Context, rankingType string, fields []string, limit, offset int) ([]AnimeRanked, *Paging, error) {
	p := paginatedParams(fields, limit, offset)
	p.Set("ranking_type", rankingType)
	var out animeRankingResponse
	if err := c.get(ctx, "/anime/ranking", p, &out); err != nil {
		return nil, nil, err
	}
	ranked := make([]AnimeRanked, len(out.Data))
	for i := range out.Data {
		ranked[i] = AnimeRanked{Anime: out.Data[i].Node, Ranking: out.Data[i].Ranking}
	}
	return ranked, &out.Paging, nil
}

// AnimeRanked pairs an anime with its ranking position.
type AnimeRanked struct {
	Anime   Anime
	Ranking Ranking
}

// SeasonalAnime returns anime for a given year and season.
func (c *Client) SeasonalAnime(ctx context.Context, year int, season string, fields []string, sort string, limit, offset int) ([]Anime, *Paging, error) {
	p := paginatedParams(fields, limit, offset)
	if sort != "" {
		p.Set("sort", sort)
	}
	path := fmt.Sprintf("/anime/season/%d/%s", year, season)
	var out animeListResponse
	if err := c.get(ctx, path, p, &out); err != nil {
		return nil, nil, err
	}
	anime := make([]Anime, len(out.Data))
	for i := range out.Data {
		anime[i] = out.Data[i].Node
	}
	return anime, &out.Paging, nil
}

// GetManga returns details for the manga with the given ID.
func (c *Client) GetManga(ctx context.Context, mangaID int, fields []string) (*Manga, error) {
	var out Manga
	path := fmt.Sprintf("/manga/%d", mangaID)
	if err := c.get(ctx, path, fieldsParam(fields), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchManga searches for manga by query string.
func (c *Client) SearchManga(ctx context.Context, query string, fields []string, limit, offset int) ([]Manga, *Paging, error) {
	p := paginatedParams(fields, limit, offset)
	p.Set("q", query)
	var out mangaListResponse
	if err := c.get(ctx, "/manga", p, &out); err != nil {
		return nil, nil, err
	}
	manga := make([]Manga, len(out.Data))
	for i := range out.Data {
		manga[i] = out.Data[i].Node
	}
	return manga, &out.Paging, nil
}

// MangaRanking returns manga ranked by the given ranking type.
func (c *Client) MangaRanking(ctx context.Context, rankingType string, fields []string, limit, offset int) ([]MangaRanked, *Paging, error) {
	p := paginatedParams(fields, limit, offset)
	p.Set("ranking_type", rankingType)
	var out mangaRankingResponse
	if err := c.get(ctx, "/manga/ranking", p, &out); err != nil {
		return nil, nil, err
	}
	ranked := make([]MangaRanked, len(out.Data))
	for i := range out.Data {
		ranked[i] = MangaRanked{Manga: out.Data[i].Node, Ranking: out.Data[i].Ranking}
	}
	return ranked, &out.Paging, nil
}

// MangaRanked pairs a manga with its ranking position.
type MangaRanked struct {
	Manga   Manga
	Ranking Ranking
}

// ForumBoards returns the forum board list.
func (c *Client) ForumBoards(ctx context.Context) ([]ForumCategory, error) {
	var out forumBoardsResponse
	if err := c.get(ctx, "/forum/boards", nil, &out); err != nil {
		return nil, err
	}
	return out.Categories, nil
}

// ForumTopicDetail returns posts and an optional poll for a topic.
func (c *Client) ForumTopicDetail(ctx context.Context, topicID, limit, offset int) (*ForumTopicDetail, *Paging, error) {
	p := url.Values{}
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		p.Set("offset", strconv.Itoa(offset))
	}
	path := fmt.Sprintf("/forum/topic/%d", topicID)
	var out forumTopicDetailResponse
	if err := c.get(ctx, path, p, &out); err != nil {
		return nil, nil, err
	}
	return &out.Data, &out.Paging, nil
}

// ForumTopics searches for forum topics. At minimum, query must be non-empty.
func (c *Client) ForumTopics(ctx context.Context, query string, boardID, subboardID, limit, offset int) ([]ForumTopic, *Paging, error) {
	p := url.Values{}
	if query != "" {
		p.Set("q", query)
	}
	if boardID > 0 {
		p.Set("board_id", strconv.Itoa(boardID))
	}
	if subboardID > 0 {
		p.Set("subboard_id", strconv.Itoa(subboardID))
	}
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		p.Set("offset", strconv.Itoa(offset))
	}
	var out forumTopicsResponse
	if err := c.get(ctx, "/forum/topics", p, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Paging, nil
}
