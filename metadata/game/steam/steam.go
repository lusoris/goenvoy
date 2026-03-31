package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultStoreURL  = "https://store.steampowered.com/api"
	defaultWebAPIURL = "https://api.steampowered.com"
)

// Client is a Steam API client.
type Client struct {
	storeURL  string
	webAPIURL string
	apiKey    string
	http      *http.Client
}

// Option configures a [Client].
type Option func(*Client)

// WithHTTPClient sets a custom [http.Client].
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.http = c }
}

// WithStoreURL sets a custom Steam Store API base URL (useful for testing).
func WithStoreURL(u string) Option {
	return func(cl *Client) { cl.storeURL = strings.TrimRight(u, "/") }
}

// WithWebAPIURL sets a custom Steam Web API base URL (useful for testing).
func WithWebAPIURL(u string) Option {
	return func(cl *Client) { cl.webAPIURL = strings.TrimRight(u, "/") }
}

// WithAPIKey sets the Steam Web API key for endpoints that require authentication.
func WithAPIKey(key string) Option {
	return func(cl *Client) { cl.apiKey = key }
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("steam: %s: %s", e.Status, e.Body)
}

// New creates a Steam API [Client].
func New(opts ...Option) *Client {
	c := &Client{
		storeURL:  defaultStoreURL,
		webAPIURL: defaultWebAPIURL,
		http:      http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) get(ctx context.Context, u string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("steam: create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("steam: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("steam: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	return json.Unmarshal(body, v)
}

func (c *Client) webAPIKeyParam() string {
	if c.apiKey != "" {
		return "&key=" + c.apiKey
	}
	return ""
}

// GetAppDetails returns detailed information about a single Steam app.
func (c *Client) GetAppDetails(ctx context.Context, appID int) (*AppDetails, error) {
	u := c.storeURL + "/appdetails?appids=" + strconv.Itoa(appID)

	var wrapper map[string]appDetailsWrapper
	if err := c.get(ctx, u, &wrapper); err != nil {
		return nil, err
	}

	entry, ok := wrapper[strconv.Itoa(appID)]
	if !ok || !entry.Success {
		return nil, fmt.Errorf("steam: app %d not found or request unsuccessful", appID)
	}

	return &entry.Data, nil
}

// GetMultipleAppDetails returns details for multiple Steam apps at once.
func (c *Client) GetMultipleAppDetails(ctx context.Context, appIDs []int) (map[int]*AppDetails, error) {
	ids := make([]string, len(appIDs))
	for i, id := range appIDs {
		ids[i] = strconv.Itoa(id)
	}

	u := c.storeURL + "/appdetails?appids=" + strings.Join(ids, ",")

	var wrapper map[string]appDetailsWrapper
	if err := c.get(ctx, u, &wrapper); err != nil {
		return nil, err
	}

	result := make(map[int]*AppDetails, len(appIDs))
	for _, id := range appIDs {
		entry, ok := wrapper[strconv.Itoa(id)]
		if ok && entry.Success {
			details := entry.Data
			result[id] = &details
		}
	}

	return result, nil
}

// GetFeatured returns currently featured games on the Steam store.
func (c *Client) GetFeatured(ctx context.Context) (*FeaturedResponse, error) {
	u := c.storeURL + "/featured"
	var resp FeaturedResponse
	if err := c.get(ctx, u, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFeaturedCategories returns featured store categories.
func (c *Client) GetFeaturedCategories(ctx context.Context) (*FeaturedCategories, error) {
	u := c.storeURL + "/featuredcategories"
	var resp FeaturedCategories
	if err := c.get(ctx, u, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAppList returns the complete list of all Steam applications.
func (c *Client) GetAppList(ctx context.Context) ([]AppListEntry, error) {
	u := c.webAPIURL + "/ISteamApps/GetAppList/v2/"
	if c.apiKey != "" {
		u += "?key=" + c.apiKey
	}

	var resp appListResponse
	if err := c.get(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.AppList.Apps, nil
}

// GetCurrentPlayers returns the number of current players for a Steam app.
func (c *Client) GetCurrentPlayers(ctx context.Context, appID int) (int, error) {
	u := c.webAPIURL + "/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=" + strconv.Itoa(appID) + c.webAPIKeyParam()

	var resp currentPlayersResponse
	if err := c.get(ctx, u, &resp); err != nil {
		return 0, err
	}

	return resp.Response.PlayerCount, nil
}

// GetAppNews returns news articles for a Steam app.
func (c *Client) GetAppNews(ctx context.Context, appID, count, maxLength int) ([]NewsItem, error) {
	u := c.webAPIURL + "/ISteamNews/GetNewsForApp/v2/?appid=" + strconv.Itoa(appID) +
		"&count=" + strconv.Itoa(count) +
		"&maxlength=" + strconv.Itoa(maxLength) +
		c.webAPIKeyParam()

	var resp appNewsResponse
	if err := c.get(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.AppNews.NewsItems, nil
}

// GetGlobalAchievements returns global achievement unlock percentages for a Steam app.
func (c *Client) GetGlobalAchievements(ctx context.Context, appID int) ([]Achievement, error) {
	u := c.webAPIURL + "/ISteamUserStats/GetGlobalAchievementPercentagesForApp/v2/?gameid=" + strconv.Itoa(appID) + c.webAPIKeyParam()

	var resp achievementsResponse
	if err := c.get(ctx, u, &resp); err != nil {
		return nil, err
	}

	return resp.AchievementPercentages.Achievements, nil
}
