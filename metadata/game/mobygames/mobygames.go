package mobygames

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.mobygames.com/v1"

// Client is a MobyGames API client.
type Client struct {
	*metadata.BaseClient
	apiKey string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("mobygames: %s: %s", e.Status, e.Body)
}

// New creates a MobyGames [Client] with the given API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "mobygames", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}

func (c *Client) get(ctx context.Context, endpoint string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_key", c.apiKey)

	u := c.BaseURL() + "/" + endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("mobygames: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("mobygames: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mobygames: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("mobygames: decode response: %w", err)
	}
	return nil
}

// GetGenres returns the list of all genres.
func (c *Client) GetGenres(ctx context.Context) ([]Genre, error) {
	var result []Genre
	if err := c.get(ctx, "genres", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetPlatforms returns the list of all platforms.
func (c *Client) GetPlatforms(ctx context.Context) ([]Platform, error) {
	var result []Platform
	if err := c.get(ctx, "platforms", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetGroups returns a paginated list of game groups.
func (c *Client) GetGroups(ctx context.Context, limit, offset int) ([]Group, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	var result []Group
	if err := c.get(ctx, "groups", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SearchGames searches for games by title.
func (c *Client) SearchGames(ctx context.Context, title string, offset, limit int) ([]Game, error) {
	params := url.Values{}
	params.Set("title", title)
	params.Set("offset", strconv.Itoa(offset))
	params.Set("limit", strconv.Itoa(limit))
	params.Set("format", "normal")
	var result struct {
		Games []Game `json:"games"`
	}
	if err := c.get(ctx, "games", params, &result); err != nil {
		return nil, err
	}
	return result.Games, nil
}

// GetGame returns a single game by ID.
func (c *Client) GetGame(ctx context.Context, gameID int) (*Game, error) {
	params := url.Values{}
	params.Set("format", "normal")
	var game Game
	if err := c.get(ctx, "games/"+strconv.Itoa(gameID), params, &game); err != nil {
		return nil, err
	}
	return &game, nil
}

// GetRecentGames returns recently updated games.
func (c *Client) GetRecentGames(ctx context.Context, age, limit, offset int) ([]Game, error) {
	params := url.Values{}
	params.Set("age", strconv.Itoa(age))
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("format", "normal")
	var result struct {
		Games []Game `json:"games"`
	}
	if err := c.get(ctx, "games/recent", params, &result); err != nil {
		return nil, err
	}
	return result.Games, nil
}

// GetRandomGames returns random games.
func (c *Client) GetRandomGames(ctx context.Context, limit int) ([]Game, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("format", "normal")
	var result struct {
		Games []Game `json:"games"`
	}
	if err := c.get(ctx, "games/random", params, &result); err != nil {
		return nil, err
	}
	return result.Games, nil
}

// GetGamePlatforms returns platforms for a specific game.
func (c *Client) GetGamePlatforms(ctx context.Context, gameID int) ([]PlatformDetail, error) {
	var result struct {
		Platforms []PlatformDetail `json:"platforms"`
	}
	if err := c.get(ctx, "games/"+strconv.Itoa(gameID)+"/platforms", nil, &result); err != nil {
		return nil, err
	}
	return result.Platforms, nil
}

// GetGameScreenshots returns screenshots for a game on a specific platform.
func (c *Client) GetGameScreenshots(ctx context.Context, gameID, platformID int) ([]Screenshot, error) {
	endpoint := "games/" + strconv.Itoa(gameID) + "/platforms/" + strconv.Itoa(platformID) + "/screenshots"
	var result struct {
		Screenshots []Screenshot `json:"screenshots"`
	}
	if err := c.get(ctx, endpoint, nil, &result); err != nil {
		return nil, err
	}
	return result.Screenshots, nil
}

// GetGameCovers returns cover art for a game on a specific platform.
func (c *Client) GetGameCovers(ctx context.Context, gameID, platformID int) ([]CoverGroup, error) {
	endpoint := "games/" + strconv.Itoa(gameID) + "/platforms/" + strconv.Itoa(platformID) + "/covers"
	var result struct {
		CoverGroups []CoverGroup `json:"cover_groups"`
	}
	if err := c.get(ctx, endpoint, nil, &result); err != nil {
		return nil, err
	}
	return result.CoverGroups, nil
}
