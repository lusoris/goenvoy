package kitsu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/lusoris/goenvoy/metadata"
)

const (
	defaultBaseURL   = "https://kitsu.io/api/edge"
	defaultAuthURL   = "https://kitsu.io/api/oauth/token"
	jsonAPIMediaType = "application/vnd.api+json"
)

// Client is a Kitsu API client. Authentication is not required for public reads.
type Client struct {
	*metadata.BaseClient
	authURL string
	onToken TokenCallback

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
}

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

// New creates a Kitsu [Client].
func New(opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "kitsu", opts...)
	return &Client{BaseClient: bc, authURL: defaultAuthURL}
}


// APIError is returned when the Kitsu API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return "kitsu: HTTP " + strconv.Itoa(e.StatusCode) + ": " + e.Body
	}
	return "kitsu: HTTP " + e.Status
}

// jsonAPIResource is the JSON:API single-resource envelope.
type jsonAPIResource[T any] struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes T      `json:"attributes"`
	} `json:"data"`
}

// jsonAPICollection is the JSON:API collection envelope.
type jsonAPICollection[T any] struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes T      `json:"attributes"`
	} `json:"data"`
	Links PageLinks `json:"links"`
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	u := c.BaseURL() + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("kitsu: create request: %w", err)
	}

	req.Header.Set("Accept", jsonAPIMediaType)
	req.Header.Set("User-Agent", c.UserAgent())

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("kitsu: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("kitsu: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	return body, nil
}

func getResource[T any](c *Client, ctx context.Context, path string) (*T, error) {
	body, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var envelope jsonAPIResource[T]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("kitsu: decode response: %w", err)
	}

	result := envelope.Data.Attributes
	setID(&result, envelope.Data.ID)

	return &result, nil
}

func getCollection[T any](c *Client, ctx context.Context, path string) ([]T, error) {
	body, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var envelope jsonAPICollection[T]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("kitsu: decode response: %w", err)
	}

	items := make([]T, len(envelope.Data))
	for i, d := range envelope.Data {
		items[i] = d.Attributes
		setID(&items[i], d.ID)
	}

	return items, nil
}

// setID sets the ID field on types that embed it via the json:"id" tag.
// This is necessary because JSON:API places id outside attributes.
func setID(v any, id string) {
	type idSetter interface{ setID(string) }
	if s, ok := v.(idSetter); ok {
		s.setID(id)
	}
}

func (a *Anime) setID(id string)      { a.ID = id }
func (m *Manga) setID(id string)      { m.ID = id }
func (e *Episode) setID(id string)    { e.ID = id }
func (ch *Character) setID(id string) { ch.ID = id }
func (ca *Category) setID(id string)  { ca.ID = id }
func (u *User) setID(id string)       { u.ID = id }

// GetAnime fetches a single anime by its Kitsu ID.
func (c *Client) GetAnime(ctx context.Context, id int64) (*Anime, error) {
	return getResource[Anime](c, ctx, "/anime/"+strconv.FormatInt(id, 10))
}

// SearchAnime searches for anime by text query.
func (c *Client) SearchAnime(ctx context.Context, query string, limit, offset int) ([]Anime, error) {
	path := "/anime?filter%5Btext%5D=" + url.QueryEscape(query) +
		"&page%5Blimit%5D=" + strconv.Itoa(limit) +
		"&page%5Boffset%5D=" + strconv.Itoa(offset)
	return getCollection[Anime](c, ctx, path)
}

// TrendingAnime returns the currently trending anime.
func (c *Client) TrendingAnime(ctx context.Context) ([]Anime, error) {
	return getCollection[Anime](c, ctx, "/trending/anime")
}

// GetAnimeEpisodes returns episodes for the given anime ID.
func (c *Client) GetAnimeEpisodes(ctx context.Context, animeID int64, limit, offset int) ([]Episode, error) {
	path := "/anime/" + strconv.FormatInt(animeID, 10) + "/episodes" +
		"?page%5Blimit%5D=" + strconv.Itoa(limit) +
		"&page%5Boffset%5D=" + strconv.Itoa(offset)
	return getCollection[Episode](c, ctx, path)
}

// GetManga fetches a single manga by its Kitsu ID.
func (c *Client) GetManga(ctx context.Context, id int64) (*Manga, error) {
	return getResource[Manga](c, ctx, "/manga/"+strconv.FormatInt(id, 10))
}

// SearchManga searches for manga by text query.
func (c *Client) SearchManga(ctx context.Context, query string, limit, offset int) ([]Manga, error) {
	path := "/manga?filter%5Btext%5D=" + url.QueryEscape(query) +
		"&page%5Blimit%5D=" + strconv.Itoa(limit) +
		"&page%5Boffset%5D=" + strconv.Itoa(offset)
	return getCollection[Manga](c, ctx, path)
}

// TrendingManga returns the currently trending manga.
func (c *Client) TrendingManga(ctx context.Context) ([]Manga, error) {
	return getCollection[Manga](c, ctx, "/trending/manga")
}

// GetCharacter fetches a single character by its Kitsu ID.
func (c *Client) GetCharacter(ctx context.Context, id int64) (*Character, error) {
	return getResource[Character](c, ctx, "/characters/"+strconv.FormatInt(id, 10))
}

// SearchCharacters searches for characters by name.
func (c *Client) SearchCharacters(ctx context.Context, name string, limit, offset int) ([]Character, error) {
	path := "/characters?filter%5Bname%5D=" + url.QueryEscape(name) +
		"&page%5Blimit%5D=" + strconv.Itoa(limit) +
		"&page%5Boffset%5D=" + strconv.Itoa(offset)
	return getCollection[Character](c, ctx, path)
}

// GetCategory fetches a single category by its Kitsu ID.
func (c *Client) GetCategory(ctx context.Context, id int64) (*Category, error) {
	return getResource[Category](c, ctx, "/categories/"+strconv.FormatInt(id, 10))
}

// GetCategories returns categories, optionally filtered by slug.
func (c *Client) GetCategories(ctx context.Context, limit, offset int) ([]Category, error) {
	path := "/categories?page%5Blimit%5D=" + strconv.Itoa(limit) +
		"&page%5Boffset%5D=" + strconv.Itoa(offset)
	return getCollection[Category](c, ctx, path)
}

// GetUser fetches a single user by their Kitsu ID.
func (c *Client) GetUser(ctx context.Context, id int64) (*User, error) {
	return getResource[User](c, ctx, "/users/"+strconv.FormatInt(id, 10))
}

// SearchUsers searches for users by name.
func (c *Client) SearchUsers(ctx context.Context, query string, limit, offset int) ([]User, error) {
	path := "/users?filter%5Bquery%5D=" + url.QueryEscape(query) +
		"&page%5Blimit%5D=" + strconv.Itoa(limit) +
		"&page%5Boffset%5D=" + strconv.Itoa(offset)
	return getCollection[User](c, ctx, path)
}

// OAuth2.

// Authenticate performs the resource owner password grant to obtain tokens.
func (c *Client) Authenticate(ctx context.Context, username, password string) (*Token, error) {
	data := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	return c.tokenRequest(ctx, data)
}

// RefreshToken uses the stored refresh token to obtain a new access token.
func (c *Client) RefreshToken(ctx context.Context) (*Token, error) {
	c.mu.RLock()
	rt := c.refreshToken
	c.mu.RUnlock()
	if rt == "" {
		return nil, errors.New("kitsu: no refresh token available")
	}
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {rt},
	}
	return c.tokenRequest(ctx, data)
}

func (c *Client) tokenRequest(ctx context.Context, data url.Values) (*Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.authURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("kitsu: create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("kitsu: token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("kitsu: read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("kitsu: decode token response: %w", err)
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
