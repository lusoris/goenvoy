package steamgriddb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://www.steamgriddb.com/api/v2"

// Client is a SteamGridDB API client.
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
	return fmt.Sprintf("steamgriddb: %s: %s", e.Status, e.Body)
}

// New creates a SteamGridDB [Client] with the given API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "steamgriddb", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}

func (c *Client) get(ctx context.Context, endpoint string, params url.Values, v any) error {
	u := c.BaseURL() + "/" + endpoint
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("steamgriddb: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("steamgriddb: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("steamgriddb: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("steamgriddb: decode response: %w", err)
	}
	return nil
}

// GetGameByID returns a game by its SteamGridDB ID.
func (c *Client) GetGameByID(ctx context.Context, id int) (*Game, error) {
	var resp response[Game]
	if err := c.get(ctx, "games/id/"+strconv.Itoa(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetGameByPlatformID returns a game by its platform-specific ID.
// Platform can be "steam", "origin", "egs", "bnet", "uplay", "flashpoint", "eshop".
func (c *Client) GetGameByPlatformID(ctx context.Context, platform string, platformID int) (*Game, error) {
	var resp response[Game]
	if err := c.get(ctx, "games/"+url.PathEscape(platform)+"/"+strconv.Itoa(platformID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SearchGames searches for games by name.
func (c *Client) SearchGames(ctx context.Context, term string) ([]SearchResult, error) {
	var resp response[[]SearchResult]
	if err := c.get(ctx, "search/autocomplete/"+url.PathEscape(term), nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ImageOptions holds optional query parameters for image endpoints.
type ImageOptions struct {
	Styles     []string
	Dimensions []string
	Mimes      []string
	Types      []string
	NSFW       *bool
	Humor      *bool
	Epilepsy   *bool
	Limit      int
	Page       int
}

func (o *ImageOptions) params() url.Values {
	if o == nil {
		return nil
	}
	params := url.Values{}
	if len(o.Styles) > 0 {
		params.Set("styles", strings.Join(o.Styles, ","))
	}
	if len(o.Dimensions) > 0 {
		params.Set("dimensions", strings.Join(o.Dimensions, ","))
	}
	if len(o.Mimes) > 0 {
		params.Set("mimes", strings.Join(o.Mimes, ","))
	}
	if len(o.Types) > 0 {
		params.Set("types", strings.Join(o.Types, ","))
	}
	if o.NSFW != nil {
		params.Set("nsfw", boolStr(*o.NSFW))
	}
	if o.Humor != nil {
		params.Set("humor", boolStr(*o.Humor))
	}
	if o.Epilepsy != nil {
		params.Set("epilepsy", boolStr(*o.Epilepsy))
	}
	if o.Limit > 0 {
		params.Set("limit", strconv.Itoa(o.Limit))
	}
	if o.Page > 0 {
		params.Set("page", strconv.Itoa(o.Page))
	}
	return params
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// GetGrids returns grid images for a game.
func (c *Client) GetGrids(ctx context.Context, gameID int, opts *ImageOptions) ([]Image, error) {
	var resp response[[]Image]
	if err := c.get(ctx, "grids/game/"+strconv.Itoa(gameID), opts.params(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetHeroes returns hero images for a game.
func (c *Client) GetHeroes(ctx context.Context, gameID int, opts *ImageOptions) ([]Image, error) {
	var resp response[[]Image]
	if err := c.get(ctx, "heroes/game/"+strconv.Itoa(gameID), opts.params(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetLogos returns logo images for a game.
func (c *Client) GetLogos(ctx context.Context, gameID int, opts *ImageOptions) ([]Image, error) {
	var resp response[[]Image]
	if err := c.get(ctx, "logos/game/"+strconv.Itoa(gameID), opts.params(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetIcons returns icon images for a game.
func (c *Client) GetIcons(ctx context.Context, gameID int, opts *ImageOptions) ([]Image, error) {
	var resp response[[]Image]
	if err := c.get(ctx, "icons/game/"+strconv.Itoa(gameID), opts.params(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
