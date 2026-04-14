package screenscraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.screenscraper.fr/api2"

// Client is a Screenscraper API v2 client.
type Client struct {
	*metadata.BaseClient
	devID        string
	devPassword  string
	softName     string
	userID       string
	userPassword string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("screenscraper: %s: %s", e.Status, e.Body)
}

// ClientOption configures the Screenscraper client.
type ClientOption func(*Client)

// WithUser sets the user credentials for the client.
func WithUser(userID, userPassword string) ClientOption {
	return func(c *Client) {
		c.userID = userID
		c.userPassword = userPassword
	}
}

// New creates a Screenscraper [Client] with the given developer credentials.
func New(devID, devPassword, softName string, opts ...any) *Client {
	var metaOpts []metadata.Option
	var clientOpts []ClientOption
	for _, o := range opts {
		switch v := o.(type) {
		case metadata.Option:
			metaOpts = append(metaOpts, v)
		case ClientOption:
			clientOpts = append(clientOpts, v)
		}
	}
	bc := metadata.NewBaseClient(defaultBaseURL, "screenscraper", metaOpts...)
	c := &Client{
		BaseClient:  bc,
		devID:       devID,
		devPassword: devPassword,
		softName:    softName,
	}
	for _, o := range clientOpts {
		o(c)
	}
	return c
}

func (c *Client) authParams() url.Values {
	params := url.Values{}
	params.Set("devid", c.devID)
	params.Set("devpassword", c.devPassword)
	params.Set("softname", c.softName)
	params.Set("output", "json")
	if c.userID != "" {
		params.Set("ssid", c.userID)
		params.Set("sspassword", c.userPassword)
	}
	return params
}

func (c *Client) get(ctx context.Context, endpoint string, params url.Values, v any) error {
	auth := c.authParams()
	for k, vs := range params {
		for _, val := range vs {
			auth.Add(k, val)
		}
	}

	u := c.BaseURL() + "/" + endpoint + "?" + auth.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("screenscraper: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("screenscraper: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("screenscraper: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("screenscraper: decode response: %w", err)
	}
	return nil
}

// GameInfoOptions holds parameters for looking up a game.
type GameInfoOptions struct {
	CRC      string
	MD5      string
	SHA1     string
	SystemID string
	ROMType  string
	ROMName  string
	ROMSize  string
}

func (o *GameInfoOptions) params() url.Values {
	params := url.Values{}
	if o.CRC != "" {
		params.Set("crc", o.CRC)
	}
	if o.MD5 != "" {
		params.Set("md5", o.MD5)
	}
	if o.SHA1 != "" {
		params.Set("sha1", o.SHA1)
	}
	if o.SystemID != "" {
		params.Set("systemeid", o.SystemID)
	}
	if o.ROMType != "" {
		params.Set("romtype", o.ROMType)
	}
	if o.ROMName != "" {
		params.Set("romnom", o.ROMName)
	}
	if o.ROMSize != "" {
		params.Set("romtaille", o.ROMSize)
	}
	return params
}

// GetGameInfo looks up game information by ROM hash or name.
func (c *Client) GetGameInfo(ctx context.Context, opts *GameInfoOptions) (*GameInfoResponse, error) {
	var result GameInfoResponse
	if err := c.get(ctx, "jeuInfos.php", opts.params(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchGames searches for games by name, optionally filtered by system.
func (c *Client) SearchGames(ctx context.Context, query, systemID string) (*GameInfoResponse, error) {
	params := url.Values{}
	params.Set("recherche", query)
	if systemID != "" {
		params.Set("systemeid", systemID)
	}
	var result GameInfoResponse
	if err := c.get(ctx, "jeuRecherche.php", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSystems returns the list of all known systems/platforms.
func (c *Client) GetSystems(ctx context.Context) (*SystemsResponse, error) {
	var result SystemsResponse
	if err := c.get(ctx, "systemesListe.php", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGenres returns the list of all genres.
func (c *Client) GetGenres(ctx context.Context) (*GenresResponse, error) {
	var result GenresResponse
	if err := c.get(ctx, "genresListe.php", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserInfo returns information about the authenticated user.
func (c *Client) GetUserInfo(ctx context.Context) (*UserInfoResponse, error) {
	var result UserInfoResponse
	if err := c.get(ctx, "ssuserInfos.php", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInfraInfo returns API infrastructure information.
func (c *Client) GetInfraInfo(ctx context.Context) (*InfraInfoResponse, error) {
	var result InfraInfoResponse
	if err := c.get(ctx, "ssinfraInfos.php", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
