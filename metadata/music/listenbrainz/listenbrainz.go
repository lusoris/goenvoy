package listenbrainz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.listenbrainz.org"

// Client is a ListenBrainz API client.
type Client struct {
	*metadata.BaseClient
	token string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("listenbrainz: %s: %s", e.Status, e.Body)
}

// New creates a ListenBrainz [Client]. A token is only required for write
// operations and can be set via [NewWithToken].
func New(opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "listenbrainz", opts...)
	return &Client{BaseClient: bc}
}

// NewWithToken creates a ListenBrainz [Client] with a user token for authenticated operations.
func NewWithToken(token string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "listenbrainz", opts...)
	c := &Client{BaseClient: bc, token: token}
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Token "+token)
	})
	return c
}

func (c *Client) get(ctx context.Context, path string, v any) error {
	u := c.BaseURL() + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("listenbrainz: create request: %w", err)
	}

	return c.do(req, v)
}

func (c *Client) post(ctx context.Context, path string, body, v any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("listenbrainz: marshal request: %w", err)
	}

	u := c.BaseURL() + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("listenbrainz: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if c.token != "" {
		req.Header.Set("Authorization", "Token "+c.token)
	}

	return c.do(req, v)
}

func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.HTTPClient().Do(req) //nolint:gosec // user-configured base URL is intentional for API clients
	if err != nil {
		return fmt.Errorf("listenbrainz: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("listenbrainz: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("listenbrainz: decode response: %w", err)
		}
		return nil
	}

	return nil
}

// SubmitListens submits one or more listens to ListenBrainz. listenType must
// be "single", "playing_now", or "import". A token must be set via [WithToken].
func (c *Client) SubmitListens(ctx context.Context, listenType string, listens []Listen) error {
	payload := make([]listenSubmit, len(listens))
	for i, l := range listens {
		payload[i] = listenSubmit(l)
	}
	body := submitListensPayload{
		ListenType: listenType,
		Payload:    payload,
	}
	return c.post(ctx, "/1/submit-listens", body, nil)
}

// GetUserListens returns the most recent listens for a user.
func (c *Client) GetUserListens(ctx context.Context, userName string, count int) (*ListensResponse, error) {
	path := "/1/user/" + url.PathEscape(userName) + "/listens?count=" + strconv.Itoa(count)
	var resp ListensResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetListenCount returns the total number of listens for a user.
func (c *Client) GetListenCount(ctx context.Context, userName string) (int64, error) {
	path := "/1/user/" + url.PathEscape(userName) + "/listen-count"
	var resp listenCountResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return 0, err
	}
	return resp.Payload.Count, nil
}

// GetPlayingNow returns the track currently being played by a user, if any.
func (c *Client) GetPlayingNow(ctx context.Context, userName string) (*PlayingNow, error) {
	path := "/1/user/" + url.PathEscape(userName) + "/playing-now"
	var resp PlayingNow
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUserTopArtists returns top artist statistics for a user.
func (c *Client) GetUserTopArtists(ctx context.Context, userName, timeRange string, count int) (*ArtistStatsResponse, error) {
	path := "/1/stats/user/" + url.PathEscape(userName) + "/artists?" + statsQuery(timeRange, count)
	var resp ArtistStatsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUserTopReleases returns top release statistics for a user.
func (c *Client) GetUserTopReleases(ctx context.Context, userName, timeRange string, count int) (*ReleaseStatsResponse, error) {
	path := "/1/stats/user/" + url.PathEscape(userName) + "/releases?" + statsQuery(timeRange, count)
	var resp ReleaseStatsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUserTopRecordings returns top recording statistics for a user.
func (c *Client) GetUserTopRecordings(ctx context.Context, userName, timeRange string, count int) (*RecordingStatsResponse, error) {
	path := "/1/stats/user/" + url.PathEscape(userName) + "/recordings?" + statsQuery(timeRange, count)
	var resp RecordingStatsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetListeningActivity returns a user's listening activity over time.
func (c *Client) GetListeningActivity(ctx context.Context, userName, timeRange string) (*ActivityResponse, error) {
	path := "/1/stats/user/" + url.PathEscape(userName) + "/listening-activity?range=" + url.QueryEscape(timeRange)
	var resp ActivityResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDailyActivity returns a user's daily listening activity broken down by hour.
func (c *Client) GetDailyActivity(ctx context.Context, userName, timeRange string) (*DailyActivityResponse, error) {
	path := "/1/stats/user/" + url.PathEscape(userName) + "/daily-activity?range=" + url.QueryEscape(timeRange)
	var resp DailyActivityResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSitewideArtists returns sitewide top artist statistics.
func (c *Client) GetSitewideArtists(ctx context.Context, timeRange string, count int) (*ArtistStatsResponse, error) {
	path := "/1/stats/sitewide/artists?" + statsQuery(timeRange, count)
	var resp ArtistStatsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSimilarUsers returns users with similar listening habits.
func (c *Client) GetSimilarUsers(ctx context.Context, userName string) ([]SimilarUser, error) {
	path := "/1/user/" + url.PathEscape(userName) + "/similar-users"
	var resp similarUsersResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// GetLatestImport returns the Unix timestamp of the user's latest import.
func (c *Client) GetLatestImport(ctx context.Context, userName string) (int64, error) {
	path := "/1/latest/import?user_name=" + url.QueryEscape(userName)
	var resp latestImportResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return 0, err
	}
	return resp.LatestImport, nil
}

func statsQuery(timeRange string, count int) string {
	return "count=" + strconv.Itoa(count) + "&range=" + url.QueryEscape(timeRange)
}
