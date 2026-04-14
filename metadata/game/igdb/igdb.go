package igdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.igdb.com/v4"

// Client is an IGDB v4 API client.
type Client struct {
	*metadata.BaseClient
	clientID    string
	accessToken string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("igdb: %s: %s", e.Status, e.Body)
}

// New creates an IGDB [Client] with the given Twitch client ID and access token.
func New(clientID, accessToken string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "igdb", opts...)
	c := &Client{BaseClient: bc, clientID: clientID, accessToken: accessToken}
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Client-Id", clientID)
		req.Header.Set("Authorization", "Bearer "+accessToken)
	})
	return c
}

// query sends a POST request to the given endpoint with an APICalypse query body.
func (c *Client) query(ctx context.Context, endpoint, body string, v any) error {
	u := c.BaseURL() + "/" + endpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("igdb: create request: %w", err)
	}

	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("igdb: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("igdb: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("igdb: decode response: %w", err)
	}
	return nil
}

// idsToString converts a slice of ints to a comma-separated string like "1,2,3".
func idsToString(ids []int) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.Itoa(id)
	}
	return strings.Join(parts, ",")
}

// SearchGames searches for games by name.
func (c *Client) SearchGames(ctx context.Context, query string, limit int) ([]Game, error) {
	body := fmt.Sprintf(`search %q; fields *; limit %d;`, query, limit)
	var games []Game
	if err := c.query(ctx, "games", body, &games); err != nil {
		return nil, err
	}
	return games, nil
}

// GetGame returns a single game by ID.
func (c *Client) GetGame(ctx context.Context, id int) (*Game, error) {
	body := fmt.Sprintf("fields *; where id = %d;", id)
	var games []Game
	if err := c.query(ctx, "games", body, &games); err != nil {
		return nil, err
	}
	if len(games) == 0 {
		return nil, nil //nolint:nilnil // documented "not found" — caller checks nil
	}
	return &games[0], nil
}

// GetGames returns multiple games by their IDs.
func (c *Client) GetGames(ctx context.Context, ids []int) ([]Game, error) {
	body := fmt.Sprintf("fields *; where id = (%s);", idsToString(ids))
	var games []Game
	if err := c.query(ctx, "games", body, &games); err != nil {
		return nil, err
	}
	return games, nil
}

// GetPopularGames returns popular games sorted by total rating.
func (c *Client) GetPopularGames(ctx context.Context, limit int) ([]Game, error) {
	body := fmt.Sprintf("fields *; sort total_rating desc; where total_rating_count > 5; limit %d;", limit)
	var games []Game
	if err := c.query(ctx, "games", body, &games); err != nil {
		return nil, err
	}
	return games, nil
}

// GetPlatform returns a single platform by ID.
func (c *Client) GetPlatform(ctx context.Context, id int) (*Platform, error) {
	body := fmt.Sprintf("fields *; where id = %d;", id)
	var platforms []Platform
	if err := c.query(ctx, "platforms", body, &platforms); err != nil {
		return nil, err
	}
	if len(platforms) == 0 {
		return nil, nil //nolint:nilnil // documented "not found" — caller checks nil
	}
	return &platforms[0], nil
}

// GetPlatforms returns a paginated list of platforms.
func (c *Client) GetPlatforms(ctx context.Context, limit, offset int) ([]Platform, error) {
	body := fmt.Sprintf("fields *; limit %d; offset %d;", limit, offset)
	var platforms []Platform
	if err := c.query(ctx, "platforms", body, &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

// GetGenre returns a single genre by ID.
func (c *Client) GetGenre(ctx context.Context, id int) (*Genre, error) {
	body := fmt.Sprintf("fields *; where id = %d;", id)
	var genres []Genre
	if err := c.query(ctx, "genres", body, &genres); err != nil {
		return nil, err
	}
	if len(genres) == 0 {
		return nil, nil //nolint:nilnil // documented "not found" — caller checks nil
	}
	return &genres[0], nil
}

// GetGenres returns a paginated list of genres.
func (c *Client) GetGenres(ctx context.Context, limit, offset int) ([]Genre, error) {
	body := fmt.Sprintf("fields *; limit %d; offset %d;", limit, offset)
	var genres []Genre
	if err := c.query(ctx, "genres", body, &genres); err != nil {
		return nil, err
	}
	return genres, nil
}

// GetCompany returns a single company by ID.
func (c *Client) GetCompany(ctx context.Context, id int) (*Company, error) {
	body := fmt.Sprintf("fields *; where id = %d;", id)
	var companies []Company
	if err := c.query(ctx, "companies", body, &companies); err != nil {
		return nil, err
	}
	if len(companies) == 0 {
		return nil, nil //nolint:nilnil // documented "not found" — caller checks nil
	}
	return &companies[0], nil
}

// SearchCompanies searches for companies by name.
func (c *Client) SearchCompanies(ctx context.Context, query string, limit int) ([]Company, error) {
	body := fmt.Sprintf(`search %q; fields *; limit %d;`, query, limit)
	var companies []Company
	if err := c.query(ctx, "companies", body, &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

// GetGameCovers returns cover images for a game.
func (c *Client) GetGameCovers(ctx context.Context, gameID int) ([]Cover, error) {
	body := fmt.Sprintf("fields *; where game = %d;", gameID)
	var covers []Cover
	if err := c.query(ctx, "covers", body, &covers); err != nil {
		return nil, err
	}
	return covers, nil
}

// GetGameScreenshots returns screenshots for a game.
func (c *Client) GetGameScreenshots(ctx context.Context, gameID int) ([]Screenshot, error) {
	body := fmt.Sprintf("fields *; where game = %d;", gameID)
	var screenshots []Screenshot
	if err := c.query(ctx, "screenshots", body, &screenshots); err != nil {
		return nil, err
	}
	return screenshots, nil
}
