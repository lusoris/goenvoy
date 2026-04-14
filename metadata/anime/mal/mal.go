package mal

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/golusoris/goenvoy/metadata"
)

const (
	defaultBaseURL = "https://api.myanimelist.net/v2"
	defaultAuthURL = "https://myanimelist.net/v1/oauth2"
)

// Client is a MyAnimeList API v2 client.
type Client struct {
	*metadata.BaseClient
	clientID     string
	clientSecret string
	authURL      string
	onToken      TokenCallback

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
}

// SetClientSecret sets the client secret (required for confidential clients).
func (c *Client) SetClientSecret(secret string) { c.clientSecret = secret }

// SetAuthURL overrides the default OAuth2 authorization URL.
func (c *Client) SetAuthURL(u string) { c.authURL = u }

// SetAccessToken sets a pre-existing OAuth2 access token.
func (c *Client) SetAccessToken(token string) {
	c.mu.Lock()
	c.accessToken = token
	c.mu.Unlock()
}

// SetRefreshToken sets a pre-existing OAuth2 refresh token.
func (c *Client) SetRefreshToken(token string) {
	c.mu.Lock()
	c.refreshToken = token
	c.mu.Unlock()
}

// TokenCallback is called whenever a new token pair is obtained.
type TokenCallback func(token Token)

// SetTokenCallback sets a callback invoked when tokens change.
func (c *Client) SetTokenCallback(cb TokenCallback) { c.onToken = cb }

// New creates a [Client] using the given MAL API client ID.
func New(clientID string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "mal", opts...)
	return &Client{BaseClient: bc, clientID: clientID, authURL: defaultAuthURL}
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
	u, err := url.Parse(c.BaseURL() + path)
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

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.Header.Set("X-Mal-Client-Id", c.clientID)
	}

	resp, err := c.HTTPClient().Do(req)
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

// OAuth2 with PKCE.

// PKCEChallenge holds a code verifier and its S256 challenge for the PKCE flow.
type PKCEChallenge struct {
	CodeVerifier  string
	CodeChallenge string
}

// GeneratePKCE creates a random PKCE code verifier and its S256 challenge.
func GeneratePKCE() (*PKCEChallenge, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("mal: generate random bytes: %w", err)
	}
	verifier := base64.RawURLEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return &PKCEChallenge{
		CodeVerifier:  verifier,
		CodeChallenge: challenge,
	}, nil
}

// AuthorizationURL constructs the URL the user should visit to authorize the app.
func (c *Client) AuthorizationURL(state string, pkce *PKCEChallenge) string {
	v := url.Values{
		"response_type":         {"code"},
		"client_id":             {c.clientID},
		"code_challenge":        {pkce.CodeChallenge},
		"code_challenge_method": {"S256"},
	}
	if state != "" {
		v.Set("state", state)
	}
	return c.authURL + "/authorize?" + v.Encode()
}

// ExchangeCode exchanges an authorization code for access and refresh tokens.
func (c *Client) ExchangeCode(ctx context.Context, code string, pkce *PKCEChallenge, redirectURI string) (*Token, error) {
	data := url.Values{
		"client_id":     {c.clientID},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {pkce.CodeVerifier},
	}
	if c.clientSecret != "" {
		data.Set("client_secret", c.clientSecret)
	}
	if redirectURI != "" {
		data.Set("redirect_uri", redirectURI)
	}
	return c.tokenRequest(ctx, data)
}

// RefreshToken uses the stored refresh token to obtain a new access token.
func (c *Client) RefreshToken(ctx context.Context) (*Token, error) {
	c.mu.RLock()
	rt := c.refreshToken
	c.mu.RUnlock()
	if rt == "" {
		return nil, errors.New("mal: no refresh token available")
	}
	data := url.Values{
		"client_id":     {c.clientID},
		"grant_type":    {"refresh_token"},
		"refresh_token": {rt},
	}
	if c.clientSecret != "" {
		data.Set("client_secret", c.clientSecret)
	}
	return c.tokenRequest(ctx, data)
}

func (c *Client) tokenRequest(ctx context.Context, data url.Values) (*Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.authURL+"/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("mal: create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("mal: token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("mal: read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(body, apiErr); jsonErr != nil {
			apiErr.RawBody = string(body)
		}
		return nil, apiErr
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("mal: decode token response: %w", err)
	}

	c.mu.Lock()
	c.accessToken = token.AccessToken
	c.refreshToken = token.RefreshToken
	c.mu.Unlock()
	if c.onToken != nil {
		c.onToken(token)
	}
	return &token, nil
}
