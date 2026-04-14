package rawg

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

const defaultBaseURL = "https://api.rawg.io/api"

// Client is a RAWG Video Games Database API client.
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
	return fmt.Sprintf("rawg: %s: %s", e.Status, e.Body)
}

// New creates a RAWG [Client] with the given API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "rawg", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}

func (c *Client) get(ctx context.Context, endpoint string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("key", c.apiKey)

	u := c.BaseURL() + "/" + endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("rawg: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("rawg: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("rawg: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("rawg: decode response: %w", err)
	}
	return nil
}

// SearchGames searches for games by query string.
func (c *Client) SearchGames(ctx context.Context, query string, page, pageSize int) (*PagedResult[GameListItem], error) {
	params := url.Values{}
	params.Set("search", query)
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[GameListItem]
	if err := c.get(ctx, "games", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGame returns a single game by ID.
func (c *Client) GetGame(ctx context.Context, id int) (*Game, error) {
	var game Game
	if err := c.get(ctx, "games/"+strconv.Itoa(id), nil, &game); err != nil {
		return nil, err
	}
	return &game, nil
}

// GetGameBySlug returns a single game by slug.
func (c *Client) GetGameBySlug(ctx context.Context, slug string) (*Game, error) {
	var game Game
	if err := c.get(ctx, "games/"+url.PathEscape(slug), nil, &game); err != nil {
		return nil, err
	}
	return &game, nil
}

// GetGameScreenshots returns screenshots for a game.
func (c *Client) GetGameScreenshots(ctx context.Context, gameID, page, pageSize int) (*PagedResult[Screenshot], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[Screenshot]
	if err := c.get(ctx, "games/"+strconv.Itoa(gameID)+"/screenshots", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGameTrailers returns trailers for a game.
func (c *Client) GetGameTrailers(ctx context.Context, gameID int) (*PagedResult[Trailer], error) {
	var result PagedResult[Trailer]
	if err := c.get(ctx, "games/"+strconv.Itoa(gameID)+"/movies", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGameAdditions returns DLC and editions for a game.
func (c *Client) GetGameAdditions(ctx context.Context, gameID, page, pageSize int) (*PagedResult[GameListItem], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[GameListItem]
	if err := c.get(ctx, "games/"+strconv.Itoa(gameID)+"/additions", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGameSeries returns games in the same series.
func (c *Client) GetGameSeries(ctx context.Context, gameID, page, pageSize int) (*PagedResult[GameListItem], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[GameListItem]
	if err := c.get(ctx, "games/"+strconv.Itoa(gameID)+"/game-series", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlatforms returns a paginated list of platforms.
func (c *Client) GetPlatforms(ctx context.Context, page, pageSize int) (*PagedResult[Platform], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[Platform]
	if err := c.get(ctx, "platforms", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlatform returns a single platform by ID.
func (c *Client) GetPlatform(ctx context.Context, id int) (*Platform, error) {
	var p Platform
	if err := c.get(ctx, "platforms/"+strconv.Itoa(id), nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetGenres returns all genres.
func (c *Client) GetGenres(ctx context.Context) (*PagedResult[Genre], error) {
	var result PagedResult[Genre]
	if err := c.get(ctx, "genres", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPublishers returns a paginated list of publishers.
func (c *Client) GetPublishers(ctx context.Context, page, pageSize int) (*PagedResult[Publisher], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[Publisher]
	if err := c.get(ctx, "publishers", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDevelopers returns a paginated list of developers.
func (c *Client) GetDevelopers(ctx context.Context, page, pageSize int) (*PagedResult[Developer], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[Developer]
	if err := c.get(ctx, "developers", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTags returns a paginated list of tags.
func (c *Client) GetTags(ctx context.Context, page, pageSize int) (*PagedResult[Tag], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	var result PagedResult[Tag]
	if err := c.get(ctx, "tags", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStores returns all stores.
func (c *Client) GetStores(ctx context.Context) (*PagedResult[Store], error) {
	var result PagedResult[Store]
	if err := c.get(ctx, "stores", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
