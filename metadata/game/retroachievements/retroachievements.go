package retroachievements

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

const defaultBaseURL = "https://retroachievements.org/API"

// Client is a RetroAchievements API client.
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
	return fmt.Sprintf("retroachievements: %s: %s", e.Status, e.Body)
}

// New creates a RetroAchievements [Client] with the given web API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "retroachievements", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}

func (c *Client) get(ctx context.Context, endpoint string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("y", c.apiKey)

	u := c.BaseURL() + "/" + endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("retroachievements: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("retroachievements: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("retroachievements: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("retroachievements: decode response: %w", err)
	}
	return nil
}

// GetGame returns basic game information by ID.
func (c *Client) GetGame(ctx context.Context, gameID int) (*Game, error) {
	params := url.Values{}
	params.Set("i", strconv.Itoa(gameID))
	var game Game
	if err := c.get(ctx, "API_GetGame.php", params, &game); err != nil {
		return nil, err
	}
	return &game, nil
}

// GetGameExtended returns extended game information including achievements.
func (c *Client) GetGameExtended(ctx context.Context, gameID int) (*GameExtended, error) {
	params := url.Values{}
	params.Set("i", strconv.Itoa(gameID))
	var game GameExtended
	if err := c.get(ctx, "API_GetGameExtended.php", params, &game); err != nil {
		return nil, err
	}
	return &game, nil
}

// GetGameHashes returns ROM hashes associated with a game.
func (c *Client) GetGameHashes(ctx context.Context, gameID int) (*HashResult, error) {
	params := url.Values{}
	params.Set("i", strconv.Itoa(gameID))
	var result HashResult
	if err := c.get(ctx, "API_GetGameHashes.php", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetConsoleIDs returns the list of all consoles/platforms.
// If onlyActive is true, only active consoles are returned.
// If onlyGameSystems is true, only game systems (not hubs) are returned.
func (c *Client) GetConsoleIDs(ctx context.Context, onlyActive, onlyGameSystems bool) ([]Console, error) {
	params := url.Values{}
	if onlyActive {
		params.Set("a", "1")
	}
	if onlyGameSystems {
		params.Set("g", "1")
	}
	var consoles []Console
	if err := c.get(ctx, "API_GetConsoleIDs.php", params, &consoles); err != nil {
		return nil, err
	}
	return consoles, nil
}
