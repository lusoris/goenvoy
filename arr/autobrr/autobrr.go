package autobrr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultTimeout = 30 * time.Second

// Option configures a [Client].
type Option func(*Client)

// WithHTTPClient sets a custom [http.Client].
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithTimeout overrides the default HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(cl *Client) { cl.httpClient.Timeout = d }
}

// Client is an autobrr API client.
type Client struct {
	rawBaseURL string
	apiKey     string
	httpClient *http.Client
}

// New creates an autobrr [Client] for the instance at baseURL with the given API key.
func New(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{
		rawBaseURL: baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	RawBody    string `json:"-"`
}

func (e *APIError) Error() string {
	if e.RawBody != "" {
		return fmt.Sprintf("autobrr: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("autobrr: HTTP %d", e.StatusCode)
}

func (c *Client) do(ctx context.Context, method, path string, reqBody any) ([]byte, error) {
	u, err := url.Parse(c.rawBaseURL + "/api" + path)
	if err != nil {
		return nil, fmt.Errorf("autobrr: parse URL: %w", err)
	}

	var body io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("autobrr: encode request: %w", err)
		}
		body = bytes.NewReader(b)
	} else {
		body = http.NoBody
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("autobrr: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Token", c.apiKey)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("autobrr: %s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("autobrr: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, RawBody: string(data)}
	}

	return data, nil
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

func decode[T any](data []byte) (T, error) {
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("autobrr: decode response: %w", err)
	}
	return out, nil
}

// Health checks.

// Liveness checks if the application is running.
func (c *Client) Liveness(ctx context.Context) error {
	_, err := c.get(ctx, "/healthz/liveness")
	return err
}

// Readiness checks if the application and dependencies are ready.
func (c *Client) Readiness(ctx context.Context) error {
	_, err := c.get(ctx, "/healthz/readiness")
	return err
}

// Filters.

// GetFilters returns all filters.
func (c *Client) GetFilters(ctx context.Context) ([]Filter, error) {
	data, err := c.get(ctx, "/filters")
	if err != nil {
		return nil, err
	}
	var filters []Filter
	if err := json.Unmarshal(data, &filters); err != nil {
		return nil, fmt.Errorf("autobrr: decode filters: %w", err)
	}
	return filters, nil
}

// SetFilterEnabled enables or disables a filter.
func (c *Client) SetFilterEnabled(ctx context.Context, filterID int, enabled bool) error {
	_, err := c.do(ctx, http.MethodPut, "/filters/"+strconv.Itoa(filterID)+"/enabled", map[string]bool{"enabled": enabled})
	return err
}

// Indexers.

// GetIndexers returns all indexers.
func (c *Client) GetIndexers(ctx context.Context) ([]Indexer, error) {
	data, err := c.get(ctx, "/indexer")
	if err != nil {
		return nil, err
	}
	var indexers []Indexer
	if err := json.Unmarshal(data, &indexers); err != nil {
		return nil, fmt.Errorf("autobrr: decode indexers: %w", err)
	}
	return indexers, nil
}

// SetIndexerEnabled enables or disables an indexer.
func (c *Client) SetIndexerEnabled(ctx context.Context, indexerID int, enabled bool) error {
	_, err := c.do(ctx, http.MethodPatch, "/indexer/"+strconv.Itoa(indexerID)+"/enabled", map[string]bool{"enabled": enabled})
	return err
}

// IRC networks.

// GetIRCNetworks returns all IRC networks.
func (c *Client) GetIRCNetworks(ctx context.Context) ([]IRCNetwork, error) {
	data, err := c.get(ctx, "/irc")
	if err != nil {
		return nil, err
	}
	var networks []IRCNetwork
	if err := json.Unmarshal(data, &networks); err != nil {
		return nil, fmt.Errorf("autobrr: decode IRC networks: %w", err)
	}
	return networks, nil
}

// RestartIRCNetwork restarts a specific IRC network.
func (c *Client) RestartIRCNetwork(ctx context.Context, networkID int) error {
	_, err := c.get(ctx, "/irc/network/"+strconv.Itoa(networkID)+"/restart")
	return err
}

// Feeds.

// GetFeeds returns all feeds.
func (c *Client) GetFeeds(ctx context.Context) ([]Feed, error) {
	data, err := c.get(ctx, "/feeds")
	if err != nil {
		return nil, err
	}
	var feeds []Feed
	if err := json.Unmarshal(data, &feeds); err != nil {
		return nil, fmt.Errorf("autobrr: decode feeds: %w", err)
	}
	return feeds, nil
}

// SetFeedEnabled enables or disables a feed.
func (c *Client) SetFeedEnabled(ctx context.Context, feedID int, enabled bool) error {
	_, err := c.do(ctx, http.MethodPatch, "/feeds/"+strconv.Itoa(feedID)+"/enabled", map[string]bool{"enabled": enabled})
	return err
}

// Download clients.

// GetDownloadClients returns all download clients.
func (c *Client) GetDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	data, err := c.get(ctx, "/download_clients")
	if err != nil {
		return nil, err
	}
	var clients []DownloadClient
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("autobrr: decode download clients: %w", err)
	}
	return clients, nil
}

// Notifications.

// GetNotifications returns all notification agents.
func (c *Client) GetNotifications(ctx context.Context) ([]Notification, error) {
	data, err := c.get(ctx, "/notification")
	if err != nil {
		return nil, err
	}
	var notifs []Notification
	if err := json.Unmarshal(data, &notifs); err != nil {
		return nil, fmt.Errorf("autobrr: decode notifications: %w", err)
	}
	return notifs, nil
}

// Config.

// GetConfig returns the autobrr configuration.
func (c *Client) GetConfig(ctx context.Context) (*Config, error) {
	data, err := c.get(ctx, "/config")
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("autobrr: decode config: %w", err)
	}
	return &cfg, nil
}

// API keys.

// GetAPIKeys returns all API keys.
func (c *Client) GetAPIKeys(ctx context.Context) ([]APIKey, error) {
	data, err := c.get(ctx, "/keys")
	if err != nil {
		return nil, err
	}
	return decode[[]APIKey](data)
}

// CreateAPIKey creates a new API key.
func (c *Client) CreateAPIKey(ctx context.Context, key APIKey) error {
	_, err := c.do(ctx, http.MethodPost, "/keys", key)
	return err
}

// DeleteAPIKey deletes an API key by its key string.
func (c *Client) DeleteAPIKey(ctx context.Context, key string) error {
	_, err := c.do(ctx, http.MethodDelete, "/keys/"+key, nil)
	return err
}

// Filters (full CRUD).

// GetFilter returns a single filter by ID.
func (c *Client) GetFilter(ctx context.Context, filterID int) (*Filter, error) {
	data, err := c.get(ctx, "/filters/"+strconv.Itoa(filterID))
	if err != nil {
		return nil, err
	}
	out, err := decode[Filter](data)
	return &out, err
}

// CreateFilter creates a new filter.
func (c *Client) CreateFilter(ctx context.Context, filter *Filter) (*Filter, error) {
	data, err := c.do(ctx, http.MethodPost, "/filters", filter)
	if err != nil {
		return nil, err
	}
	out, err := decode[Filter](data)
	return &out, err
}

// UpdateFilter updates an existing filter.
func (c *Client) UpdateFilter(ctx context.Context, filter *Filter) (*Filter, error) {
	data, err := c.do(ctx, http.MethodPut, "/filters/"+strconv.Itoa(filter.ID), filter)
	if err != nil {
		return nil, err
	}
	out, err := decode[Filter](data)
	return &out, err
}

// DuplicateFilter duplicates a filter by ID.
func (c *Client) DuplicateFilter(ctx context.Context, filterID int) (*Filter, error) {
	data, err := c.get(ctx, "/filters/"+strconv.Itoa(filterID)+"/duplicate")
	if err != nil {
		return nil, err
	}
	out, err := decode[Filter](data)
	return &out, err
}

// DeleteFilter deletes a filter by ID.
func (c *Client) DeleteFilter(ctx context.Context, filterID int) error {
	_, err := c.do(ctx, http.MethodDelete, "/filters/"+strconv.Itoa(filterID), nil)
	return err
}

// GetFilterNotifications returns notification settings for a filter.
func (c *Client) GetFilterNotifications(ctx context.Context, filterID int) ([]FilterNotification, error) {
	data, err := c.get(ctx, "/filters/"+strconv.Itoa(filterID)+"/notifications")
	if err != nil {
		return nil, err
	}
	return decode[[]FilterNotification](data)
}

// UpdateFilterNotifications replaces notification settings for a filter.
func (c *Client) UpdateFilterNotifications(ctx context.Context, filterID int, notifs []FilterNotification) error {
	_, err := c.do(ctx, http.MethodPut, "/filters/"+strconv.Itoa(filterID)+"/notifications", notifs)
	return err
}

// Indexers (full CRUD).

// GetIndexerSchema returns all possible indexer definitions.
func (c *Client) GetIndexerSchema(ctx context.Context) ([]IndexerDefinition, error) {
	data, err := c.get(ctx, "/indexer/schema")
	if err != nil {
		return nil, err
	}
	return decode[[]IndexerDefinition](data)
}

// GetIndexerOptions returns indexer options for enabled indexers.
func (c *Client) GetIndexerOptions(ctx context.Context) ([]Indexer, error) {
	data, err := c.get(ctx, "/indexer/options")
	if err != nil {
		return nil, err
	}
	return decode[[]Indexer](data)
}

// CreateIndexer creates a new indexer.
func (c *Client) CreateIndexer(ctx context.Context, indexer Indexer) (*Indexer, error) {
	data, err := c.do(ctx, http.MethodPost, "/indexer", indexer)
	if err != nil {
		return nil, err
	}
	out, err := decode[Indexer](data)
	return &out, err
}

// UpdateIndexer updates an existing indexer.
func (c *Client) UpdateIndexer(ctx context.Context, indexer Indexer) error {
	_, err := c.do(ctx, http.MethodPut, "/indexer/"+strconv.Itoa(indexer.ID), indexer)
	return err
}

// DeleteIndexer deletes an indexer by ID.
func (c *Client) DeleteIndexer(ctx context.Context, indexerID int) error {
	_, err := c.do(ctx, http.MethodDelete, "/indexer/"+strconv.Itoa(indexerID), nil)
	return err
}

// TestIndexerAPI tests an indexer's API connectivity.
func (c *Client) TestIndexerAPI(ctx context.Context, indexerID int) error {
	_, err := c.do(ctx, http.MethodPost, "/indexer/"+strconv.Itoa(indexerID)+"/api/test", nil)
	return err
}

// IRC networks (full CRUD).

// CreateIRCNetwork creates a new IRC network.
func (c *Client) CreateIRCNetwork(ctx context.Context, network *IRCNetwork) error {
	_, err := c.do(ctx, http.MethodPost, "/irc", network)
	return err
}

// UpdateIRCNetwork updates an existing IRC network.
func (c *Client) UpdateIRCNetwork(ctx context.Context, network *IRCNetwork) error {
	_, err := c.do(ctx, http.MethodPut, "/irc/network/"+strconv.Itoa(network.ID), network)
	return err
}

// DeleteIRCNetwork deletes an IRC network by ID.
func (c *Client) DeleteIRCNetwork(ctx context.Context, networkID int) error {
	_, err := c.do(ctx, http.MethodDelete, "/irc/network/"+strconv.Itoa(networkID), nil)
	return err
}

// SendIRCCommand sends a command to an IRC network channel.
func (c *Client) SendIRCCommand(ctx context.Context, cmd SendIRCCmdRequest) error {
	_, err := c.do(ctx, http.MethodPost, "/irc/network/"+strconv.Itoa(cmd.NetworkID)+"/cmd", cmd)
	return err
}

// ReprocessAnnounce reprocesses an announce message on an IRC channel.
func (c *Client) ReprocessAnnounce(ctx context.Context, networkID int, channel, msg string) error {
	path := "/irc/network/" + strconv.Itoa(networkID) + "/channel/" + url.PathEscape(channel) + "/announce/process"
	_, err := c.do(ctx, http.MethodPost, path, map[string]string{"msg": msg})
	return err
}

// Feeds (full CRUD).

// CreateFeed creates a new feed.
func (c *Client) CreateFeed(ctx context.Context, feed Feed) error {
	_, err := c.do(ctx, http.MethodPost, "/feeds", feed)
	return err
}

// UpdateFeed updates an existing feed.
func (c *Client) UpdateFeed(ctx context.Context, feed Feed) error {
	_, err := c.do(ctx, http.MethodPut, "/feeds/"+strconv.Itoa(feed.ID), feed)
	return err
}

// DeleteFeed deletes a feed by ID.
func (c *Client) DeleteFeed(ctx context.Context, feedID int) error {
	_, err := c.do(ctx, http.MethodDelete, "/feeds/"+strconv.Itoa(feedID), nil)
	return err
}

// DeleteFeedCache clears the cache for a feed.
func (c *Client) DeleteFeedCache(ctx context.Context, feedID int) error {
	_, err := c.do(ctx, http.MethodDelete, "/feeds/"+strconv.Itoa(feedID)+"/cache", nil)
	return err
}

// ForceRunFeed forces an immediate run of a feed.
func (c *Client) ForceRunFeed(ctx context.Context, feedID int) error {
	_, err := c.do(ctx, http.MethodPost, "/feeds/"+strconv.Itoa(feedID)+"/forcerun", nil)
	return err
}

// TestFeed tests a feed configuration.
func (c *Client) TestFeed(ctx context.Context, feed Feed) error {
	_, err := c.do(ctx, http.MethodPost, "/feeds/test", feed)
	return err
}

// GetFeedCaps returns capabilities for a feed by ID.
func (c *Client) GetFeedCaps(ctx context.Context, feedID int) (json.RawMessage, error) {
	data, err := c.get(ctx, "/feeds/"+strconv.Itoa(feedID)+"/caps")
	return json.RawMessage(data), err
}

// Download clients (full CRUD).

// CreateDownloadClient creates a new download client.
func (c *Client) CreateDownloadClient(ctx context.Context, dc *DownloadClient) error {
	_, err := c.do(ctx, http.MethodPost, "/download_clients", dc)
	return err
}

// UpdateDownloadClient updates an existing download client.
func (c *Client) UpdateDownloadClient(ctx context.Context, dc *DownloadClient) error {
	_, err := c.do(ctx, http.MethodPut, "/download_clients", dc)
	return err
}

// DeleteDownloadClient deletes a download client by ID.
func (c *Client) DeleteDownloadClient(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodDelete, "/download_clients/"+strconv.Itoa(id), nil)
	return err
}

// TestDownloadClient tests a download client configuration.
func (c *Client) TestDownloadClient(ctx context.Context, dc *DownloadClient) error {
	_, err := c.do(ctx, http.MethodPost, "/download_clients/test", dc)
	return err
}

// GetDownloadClientArrTags returns tags from an *arr download client.
func (c *Client) GetDownloadClientArrTags(ctx context.Context, clientID int) ([]ArrTag, error) {
	data, err := c.get(ctx, "/download_clients/"+strconv.Itoa(clientID)+"/arr/tags")
	if err != nil {
		return nil, err
	}
	return decode[[]ArrTag](data)
}

// Notifications (full CRUD).

// CreateNotification creates a new notification agent.
func (c *Client) CreateNotification(ctx context.Context, notif *Notification) error {
	_, err := c.do(ctx, http.MethodPost, "/notification", notif)
	return err
}

// UpdateNotification updates an existing notification agent.
func (c *Client) UpdateNotification(ctx context.Context, notif *Notification) error {
	_, err := c.do(ctx, http.MethodPut, "/notification/"+strconv.Itoa(notif.ID), notif)
	return err
}

// DeleteNotification deletes a notification agent by ID.
func (c *Client) DeleteNotification(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodDelete, "/notification/"+strconv.Itoa(id), nil)
	return err
}

// TestNotification tests a notification configuration.
func (c *Client) TestNotification(ctx context.Context, notif *Notification) error {
	_, err := c.do(ctx, http.MethodPost, "/notification/test", notif)
	return err
}

// Config (update).

// UpdateConfig patches the autobrr configuration.
func (c *Client) UpdateConfig(ctx context.Context, cfg ConfigUpdate) error {
	_, err := c.do(ctx, http.MethodPatch, "/config", cfg)
	return err
}

// Actions.

// CreateAction creates a new action.
func (c *Client) CreateAction(ctx context.Context, action *Action) error {
	_, err := c.do(ctx, http.MethodPost, "/actions", action)
	return err
}

// UpdateAction updates an existing action.
func (c *Client) UpdateAction(ctx context.Context, action *Action) error {
	_, err := c.do(ctx, http.MethodPut, "/actions/"+strconv.Itoa(action.ID), action)
	return err
}

// DeleteAction deletes an action by ID.
func (c *Client) DeleteAction(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodDelete, "/actions/"+strconv.Itoa(id), nil)
	return err
}

// ToggleActionEnabled toggles whether an action is enabled.
func (c *Client) ToggleActionEnabled(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodPatch, "/actions/"+strconv.Itoa(id)+"/toggleEnabled", nil)
	return err
}

// Lists.

// GetLists returns all lists.
func (c *Client) GetLists(ctx context.Context) ([]List, error) {
	data, err := c.get(ctx, "/lists")
	if err != nil {
		return nil, err
	}
	return decode[[]List](data)
}

// GetList returns a single list by ID.
func (c *Client) GetList(ctx context.Context, id int) (*List, error) {
	data, err := c.get(ctx, "/lists/"+strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	out, err := decode[List](data)
	return &out, err
}

// CreateList creates a new list.
func (c *Client) CreateList(ctx context.Context, list *List) error {
	_, err := c.do(ctx, http.MethodPost, "/lists", list)
	return err
}

// UpdateList updates an existing list.
func (c *Client) UpdateList(ctx context.Context, list *List) error {
	_, err := c.do(ctx, http.MethodPut, "/lists/"+strconv.Itoa(list.ID), list)
	return err
}

// DeleteList deletes a list by ID.
func (c *Client) DeleteList(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodDelete, "/lists/"+strconv.Itoa(id), nil)
	return err
}

// RefreshList triggers an immediate refresh of a list.
func (c *Client) RefreshList(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodPost, "/lists/"+strconv.Itoa(id)+"/refresh", nil)
	return err
}

// RefreshAllLists refreshes all lists.
func (c *Client) RefreshAllLists(ctx context.Context) error {
	_, err := c.do(ctx, http.MethodPost, "/lists/refresh", nil)
	return err
}

// TestList tests a list configuration.
func (c *Client) TestList(ctx context.Context, list *List) error {
	_, err := c.do(ctx, http.MethodPost, "/lists/test", list)
	return err
}

// Proxies.

// GetProxies returns all proxies.
func (c *Client) GetProxies(ctx context.Context) ([]Proxy, error) {
	data, err := c.get(ctx, "/proxy")
	if err != nil {
		return nil, err
	}
	return decode[[]Proxy](data)
}

// GetProxy returns a single proxy by ID.
func (c *Client) GetProxy(ctx context.Context, id int) (*Proxy, error) {
	data, err := c.get(ctx, "/proxy/"+strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	out, err := decode[Proxy](data)
	return &out, err
}

// CreateProxy creates a new proxy.
func (c *Client) CreateProxy(ctx context.Context, proxy *Proxy) error {
	_, err := c.do(ctx, http.MethodPost, "/proxy", proxy)
	return err
}

// UpdateProxy updates an existing proxy.
func (c *Client) UpdateProxy(ctx context.Context, proxy *Proxy) error {
	_, err := c.do(ctx, http.MethodPut, "/proxy/"+strconv.Itoa(proxy.ID), proxy)
	return err
}

// DeleteProxy deletes a proxy by ID.
func (c *Client) DeleteProxy(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodDelete, "/proxy/"+strconv.Itoa(id), nil)
	return err
}

// TestProxy tests a proxy configuration.
func (c *Client) TestProxy(ctx context.Context, proxy *Proxy) error {
	_, err := c.do(ctx, http.MethodPost, "/proxy/test", proxy)
	return err
}

// Releases.

// GetReleases returns releases with pagination.
func (c *Client) GetReleases(ctx context.Context, offset, limit int) (*ReleaseFindResponse, error) {
	path := fmt.Sprintf("/release?offset=%d&limit=%d", offset, limit)
	data, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	out, err := decode[ReleaseFindResponse](data)
	return &out, err
}

// GetRecentReleases returns the most recent releases.
func (c *Client) GetRecentReleases(ctx context.Context) (*ReleaseFindResponse, error) {
	data, err := c.get(ctx, "/release/recent")
	if err != nil {
		return nil, err
	}
	out, err := decode[ReleaseFindResponse](data)
	return &out, err
}

// GetReleaseStats returns aggregate release statistics.
func (c *Client) GetReleaseStats(ctx context.Context) (*ReleaseStats, error) {
	data, err := c.get(ctx, "/release/stats")
	if err != nil {
		return nil, err
	}
	out, err := decode[ReleaseStats](data)
	return &out, err
}

// GetReleaseIndexers returns the list of indexers that have releases.
func (c *Client) GetReleaseIndexers(ctx context.Context) ([]string, error) {
	data, err := c.get(ctx, "/release/indexers")
	if err != nil {
		return nil, err
	}
	return decode[[]string](data)
}

// DeleteReleases bulk-deletes releases matching the given parameters.
func (c *Client) DeleteReleases(ctx context.Context, params ReleaseDeleteParams) error {
	path := "/release?"
	v := url.Values{}
	if params.OlderThan > 0 {
		v.Set("olderThan", strconv.Itoa(params.OlderThan))
	}
	for _, idx := range params.Indexers {
		v.Add("indexer", idx)
	}
	for _, rs := range params.ReleaseStatuses {
		v.Add("releaseStatus", rs)
	}
	_, err := c.do(ctx, http.MethodDelete, path+v.Encode(), nil)
	return err
}

// ReplayReleaseAction retries an action on a release.
func (c *Client) ReplayReleaseAction(ctx context.Context, releaseID, actionID int) error {
	path := fmt.Sprintf("/release/%d/actions/%d/retry", releaseID, actionID)
	_, err := c.do(ctx, http.MethodPost, path, nil)
	return err
}

// Release cleanup jobs.

// GetReleaseCleanupJobs returns all release cleanup jobs.
func (c *Client) GetReleaseCleanupJobs(ctx context.Context) ([]ReleaseCleanupJob, error) {
	data, err := c.get(ctx, "/release/cleanup-jobs")
	if err != nil {
		return nil, err
	}
	return decode[[]ReleaseCleanupJob](data)
}

// GetReleaseCleanupJob returns a cleanup job by ID.
func (c *Client) GetReleaseCleanupJob(ctx context.Context, id int) (*ReleaseCleanupJob, error) {
	data, err := c.get(ctx, "/release/cleanup-jobs/"+strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	out, err := decode[ReleaseCleanupJob](data)
	return &out, err
}

// CreateReleaseCleanupJob creates a new cleanup job.
func (c *Client) CreateReleaseCleanupJob(ctx context.Context, job *ReleaseCleanupJob) error {
	_, err := c.do(ctx, http.MethodPost, "/release/cleanup-jobs", job)
	return err
}

// UpdateReleaseCleanupJob updates an existing cleanup job.
func (c *Client) UpdateReleaseCleanupJob(ctx context.Context, job *ReleaseCleanupJob) error {
	_, err := c.do(ctx, http.MethodPut, "/release/cleanup-jobs/"+strconv.Itoa(job.ID), job)
	return err
}

// DeleteReleaseCleanupJob deletes a cleanup job by ID.
func (c *Client) DeleteReleaseCleanupJob(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodDelete, "/release/cleanup-jobs/"+strconv.Itoa(id), nil)
	return err
}

// ToggleReleaseCleanupJobEnabled toggles whether a cleanup job is enabled.
func (c *Client) ToggleReleaseCleanupJobEnabled(ctx context.Context, id int, enabled bool) error {
	_, err := c.do(ctx, http.MethodPatch, "/release/cleanup-jobs/"+strconv.Itoa(id)+"/enabled", map[string]bool{"enabled": enabled})
	return err
}

// ForceRunReleaseCleanupJob triggers immediate execution of a cleanup job.
func (c *Client) ForceRunReleaseCleanupJob(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodPost, "/release/cleanup-jobs/"+strconv.Itoa(id)+"/run", nil)
	return err
}

// Release duplicate profiles.

// GetReleaseDuplicateProfiles returns all duplicate release profiles.
func (c *Client) GetReleaseDuplicateProfiles(ctx context.Context) ([]ReleaseProfileDuplicate, error) {
	data, err := c.get(ctx, "/release/profiles/duplicate")
	if err != nil {
		return nil, err
	}
	return decode[[]ReleaseProfileDuplicate](data)
}

// CreateReleaseDuplicateProfile creates a duplicate release profile.
func (c *Client) CreateReleaseDuplicateProfile(ctx context.Context, profile ReleaseProfileDuplicate) error {
	_, err := c.do(ctx, http.MethodPost, "/release/profiles/duplicate", profile)
	return err
}

// DeleteReleaseDuplicateProfile deletes a duplicate release profile by ID.
func (c *Client) DeleteReleaseDuplicateProfile(ctx context.Context, id int) error {
	_, err := c.do(ctx, http.MethodDelete, "/release/profiles/duplicate/"+strconv.Itoa(id), nil)
	return err
}

// Logs.

// GetLogFiles returns a list of available log files.
func (c *Client) GetLogFiles(ctx context.Context) (*LogFileResponse, error) {
	data, err := c.get(ctx, "/logs/files")
	if err != nil {
		return nil, err
	}
	out, err := decode[LogFileResponse](data)
	return &out, err
}

// GetLogFile returns the contents of a specific log file.
func (c *Client) GetLogFile(ctx context.Context, filename string) ([]byte, error) {
	return c.get(ctx, "/logs/files/"+url.PathEscape(filename))
}

// Updates.

// CheckForUpdates checks if a newer version is available.
func (c *Client) CheckForUpdates(ctx context.Context) error {
	_, err := c.get(ctx, "/updates/check")
	return err
}

// GetLatestRelease returns the latest GitHub release.
func (c *Client) GetLatestRelease(ctx context.Context) (*GithubRelease, error) {
	data, err := c.get(ctx, "/updates/latest")
	if err != nil {
		return nil, err
	}
	out, err := decode[GithubRelease](data)
	return &out, err
}
