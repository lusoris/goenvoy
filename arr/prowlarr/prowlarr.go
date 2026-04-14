package prowlarr

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/golusoris/goenvoy/arr/v2"
)

// Client is a Prowlarr API client.
type Client struct {
	base *arr.BaseClient
}

// New creates a Prowlarr [Client] for the instance at baseURL.
func New(baseURL, apiKey string, opts ...arr.Option) (*Client, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{base: base}, nil
}

// GetIndexers returns all configured indexers.
func (c *Client) GetIndexers(ctx context.Context) ([]Indexer, error) {
	var out []Indexer
	if err := c.base.Get(ctx, "/api/v1/indexer", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetIndexer returns a single indexer by its database ID.
func (c *Client) GetIndexer(ctx context.Context, id int) (*Indexer, error) {
	var out Indexer
	path := fmt.Sprintf("/api/v1/indexer/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddIndexer adds a new indexer to Prowlarr.
func (c *Client) AddIndexer(ctx context.Context, indexer *Indexer) (*Indexer, error) {
	var out Indexer
	if err := c.base.Post(ctx, "/api/v1/indexer", indexer, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexer updates an existing indexer.
func (c *Client) UpdateIndexer(ctx context.Context, indexer *Indexer) (*Indexer, error) {
	var out Indexer
	path := fmt.Sprintf("/api/v1/indexer/%d", indexer.ID)
	if err := c.base.Put(ctx, path, indexer, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteIndexer removes an indexer by its database ID.
func (c *Client) DeleteIndexer(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/indexer/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetIndexerCategories returns the global list of Newznab/Torznab categories.
func (c *Client) GetIndexerCategories(ctx context.Context) ([]IndexerCategory, error) {
	var out []IndexerCategory
	if err := c.base.Get(ctx, "/api/v1/indexer/categories", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetApplications returns all connected PVR applications.
func (c *Client) GetApplications(ctx context.Context) ([]Application, error) {
	var out []Application
	if err := c.base.Get(ctx, "/api/v1/applications", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetApplication returns a single application by its database ID.
func (c *Client) GetApplication(ctx context.Context, id int) (*Application, error) {
	var out Application
	path := fmt.Sprintf("/api/v1/applications/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddApplication adds a new application to Prowlarr.
func (c *Client) AddApplication(ctx context.Context, app *Application) (*Application, error) {
	var out Application
	if err := c.base.Post(ctx, "/api/v1/applications", app, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateApplication updates an existing application.
func (c *Client) UpdateApplication(ctx context.Context, app *Application) (*Application, error) {
	var out Application
	path := fmt.Sprintf("/api/v1/applications/%d", app.ID)
	if err := c.base.Put(ctx, path, app, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteApplication removes a connected application by its ID.
func (c *Client) DeleteApplication(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/applications/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetAppProfiles returns all app profiles.
func (c *Client) GetAppProfiles(ctx context.Context) ([]AppProfile, error) {
	var out []AppProfile
	if err := c.base.Get(ctx, "/api/v1/appprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAppProfile returns a single app profile by its ID.
func (c *Client) GetAppProfile(ctx context.Context, id int) (*AppProfile, error) {
	var out AppProfile
	path := fmt.Sprintf("/api/v1/appprofile/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddAppProfile creates a new app profile.
func (c *Client) AddAppProfile(ctx context.Context, profile *AppProfile) (*AppProfile, error) {
	var out AppProfile
	if err := c.base.Post(ctx, "/api/v1/appprofile", profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAppProfile updates an existing app profile.
func (c *Client) UpdateAppProfile(ctx context.Context, profile *AppProfile) (*AppProfile, error) {
	var out AppProfile
	path := fmt.Sprintf("/api/v1/appprofile/%d", profile.ID)
	if err := c.base.Put(ctx, path, profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAppProfile removes an app profile by its ID.
func (c *Client) DeleteAppProfile(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/appprofile/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// SearchOptions configures a search request.
type SearchOptions struct {
	Query      string
	Type       string // "search", "tvsearch", "movie", "music", "book"
	IndexerIDs []int
	Categories []int
	Limit      int
	Offset     int
}

// Search performs a search across configured indexers.
func (c *Client) Search(ctx context.Context, opts *SearchOptions) ([]Release, error) {
	var out []Release
	params := url.Values{}
	if opts.Query != "" {
		params.Set("query", opts.Query)
	}
	if opts.Type != "" {
		params.Set("type", opts.Type)
	}
	if len(opts.IndexerIDs) > 0 {
		for _, id := range opts.IndexerIDs {
			params.Add("indexerIds", strconv.Itoa(id))
		}
	}
	if len(opts.Categories) > 0 {
		for _, cat := range opts.Categories {
			params.Add("categories", strconv.Itoa(cat))
		}
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	path := "/api/v1/search"
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GrabRelease sends a release to the download client.
func (c *Client) GrabRelease(ctx context.Context, release *Release) (*Release, error) {
	var out Release
	if err := c.base.Post(ctx, "/api/v1/search", release, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDownloadClients returns all configured download clients.
func (c *Client) GetDownloadClients(ctx context.Context) ([]DownloadClientResource, error) {
	var out []DownloadClientResource
	if err := c.base.Get(ctx, "/api/v1/downloadclient", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendCommand triggers a named command.
func (c *Client) SendCommand(ctx context.Context, cmd arr.CommandRequest) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Post(ctx, "/api/v1/command", cmd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCommands returns all currently queued or running commands.
func (c *Client) GetCommands(ctx context.Context) ([]arr.CommandResponse, error) {
	var out []arr.CommandResponse
	if err := c.base.Get(ctx, "/api/v1/command", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCommand returns the status of a single command by its ID.
func (c *Client) GetCommand(ctx context.Context, id int) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	path := fmt.Sprintf("/api/v1/command/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHealth returns current health check results.
func (c *Client) GetHealth(ctx context.Context) ([]arr.HealthCheck, error) {
	var out []arr.HealthCheck
	if err := c.base.Get(ctx, "/api/v1/health", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSystemStatus returns Prowlarr system information.
func (c *Client) GetSystemStatus(ctx context.Context) (*arr.StatusResponse, error) {
	var out arr.StatusResponse
	if err := c.base.Get(ctx, "/api/v1/system/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTags returns all tags.
func (c *Client) GetTags(ctx context.Context) ([]arr.Tag, error) {
	var out []arr.Tag
	if err := c.base.Get(ctx, "/api/v1/tag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTag creates a new tag and returns it with its assigned ID.
func (c *Client) CreateTag(ctx context.Context, label string) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Post(ctx, "/api/v1/tag", arr.Tag{Label: label}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHistory returns the indexer history with pagination.
func (c *Client) GetHistory(ctx context.Context, page, pageSize int) (*arr.PagingResource[HistoryRecord], error) {
	var out arr.PagingResource[HistoryRecord]
	path := fmt.Sprintf("/api/v1/history?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHistorySince returns history events since a specific date (RFC 3339 timestamp).
func (c *Client) GetHistorySince(ctx context.Context, date string) ([]HistoryRecord, error) {
	var out []HistoryRecord
	path := "/api/v1/history/since?date=" + url.QueryEscape(date)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistoryByIndexer returns history events for a specific indexer.
func (c *Client) GetHistoryByIndexer(ctx context.Context, indexerID, limit int) ([]HistoryRecord, error) {
	var out []HistoryRecord
	path := fmt.Sprintf("/api/v1/history/indexer?indexerId=%d&limit=%d", indexerID, limit)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetIndexerStats returns aggregated indexer statistics for the given date range.
func (c *Client) GetIndexerStats(ctx context.Context, startDate, endDate string) (*IndexerStats, error) {
	var out IndexerStats
	params := []string{}
	if startDate != "" {
		params = append(params, "startDate="+url.QueryEscape(startDate))
	}
	if endDate != "" {
		params = append(params, "endDate="+url.QueryEscape(endDate))
	}
	path := "/api/v1/indexerstats"
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIndexerStatuses returns the status of all indexers (including disabled ones).
func (c *Client) GetIndexerStatuses(ctx context.Context) ([]IndexerStatus, error) {
	var out []IndexerStatus
	if err := c.base.Get(ctx, "/api/v1/indexerstatus", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Indexer Extended ----------.

// GetIndexerSchema returns the schema for all indexer types.
func (c *Client) GetIndexerSchema(ctx context.Context) ([]Indexer, error) {
	var out []Indexer
	if err := c.base.Get(ctx, "/api/v1/indexer/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestIndexer tests an indexer configuration.
func (c *Client) TestIndexer(ctx context.Context, idx *Indexer) error {
	return c.base.Post(ctx, "/api/v1/indexer/test", idx, nil)
}

// TestAllIndexers tests all configured indexers.
func (c *Client) TestAllIndexers(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/indexer/testall", nil, nil)
}

// BulkUpdateIndexers updates multiple indexers at once.
func (c *Client) BulkUpdateIndexers(ctx context.Context, bulk *IndexerBulkResource) ([]Indexer, error) {
	var out []Indexer
	if err := c.base.Put(ctx, "/api/v1/indexer/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteIndexers deletes multiple indexers at once.
func (c *Client) BulkDeleteIndexers(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v1/indexer/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// IndexerAction triggers a named action on an indexer provider.
func (c *Client) IndexerAction(ctx context.Context, name string, body *Indexer) error {
	path := "/api/v1/indexer/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Indexer Proxies ----------.

// GetIndexerProxies returns all configured indexer proxies.
func (c *Client) GetIndexerProxies(ctx context.Context) ([]IndexerProxyResource, error) {
	var out []IndexerProxyResource
	if err := c.base.Get(ctx, "/api/v1/indexerproxy", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetIndexerProxy returns a single indexer proxy by its ID.
func (c *Client) GetIndexerProxy(ctx context.Context, id int) (*IndexerProxyResource, error) {
	var out IndexerProxyResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/indexerproxy/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateIndexerProxy creates a new indexer proxy.
func (c *Client) CreateIndexerProxy(ctx context.Context, proxy *IndexerProxyResource) (*IndexerProxyResource, error) {
	var out IndexerProxyResource
	if err := c.base.Post(ctx, "/api/v1/indexerproxy", proxy, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexerProxy updates an existing indexer proxy.
func (c *Client) UpdateIndexerProxy(ctx context.Context, proxy *IndexerProxyResource) (*IndexerProxyResource, error) {
	var out IndexerProxyResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/indexerproxy/%d", proxy.ID), proxy, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteIndexerProxy deletes an indexer proxy by ID.
func (c *Client) DeleteIndexerProxy(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/indexerproxy/%d", id), nil, nil)
}

// GetIndexerProxySchema returns the schema for all indexer proxy types.
func (c *Client) GetIndexerProxySchema(ctx context.Context) ([]IndexerProxyResource, error) {
	var out []IndexerProxyResource
	if err := c.base.Get(ctx, "/api/v1/indexerproxy/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestIndexerProxy tests an indexer proxy configuration.
func (c *Client) TestIndexerProxy(ctx context.Context, proxy *IndexerProxyResource) error {
	return c.base.Post(ctx, "/api/v1/indexerproxy/test", proxy, nil)
}

// TestAllIndexerProxies tests all configured indexer proxies.
func (c *Client) TestAllIndexerProxies(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/indexerproxy/testall", nil, nil)
}

// IndexerProxyAction triggers a named action on an indexer proxy provider.
func (c *Client) IndexerProxyAction(ctx context.Context, name string, body *IndexerProxyResource) error {
	path := "/api/v1/indexerproxy/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Applications Extended ----------.

// GetApplicationSchema returns the schema for all application types.
func (c *Client) GetApplicationSchema(ctx context.Context) ([]Application, error) {
	var out []Application
	if err := c.base.Get(ctx, "/api/v1/applications/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestApplication tests an application configuration.
func (c *Client) TestApplication(ctx context.Context, app *Application) error {
	return c.base.Post(ctx, "/api/v1/applications/test", app, nil)
}

// TestAllApplications tests all configured applications.
func (c *Client) TestAllApplications(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/applications/testall", nil, nil)
}

// BulkUpdateApplications updates multiple applications at once.
func (c *Client) BulkUpdateApplications(ctx context.Context, bulk *ApplicationBulkResource) ([]Application, error) {
	var out []Application
	if err := c.base.Put(ctx, "/api/v1/applications/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteApplications deletes multiple applications at once.
func (c *Client) BulkDeleteApplications(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v1/applications/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// ApplicationAction triggers a named action on an application provider.
func (c *Client) ApplicationAction(ctx context.Context, name string, body *Application) error {
	path := "/api/v1/applications/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- App Profile Extended ----------.

// GetAppProfileSchema returns the app profile schema.
func (c *Client) GetAppProfileSchema(ctx context.Context) (*AppProfile, error) {
	var out AppProfile
	if err := c.base.Get(ctx, "/api/v1/appprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Download Clients Extended ----------.

// GetDownloadClient returns a single download client by its ID.
func (c *Client) GetDownloadClient(ctx context.Context, id int) (*DownloadClientResource, error) {
	var out DownloadClientResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/downloadclient/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDownloadClient creates a new download client.
func (c *Client) CreateDownloadClient(ctx context.Context, dc *DownloadClientResource) (*DownloadClientResource, error) {
	var out DownloadClientResource
	if err := c.base.Post(ctx, "/api/v1/downloadclient", dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClient updates an existing download client.
func (c *Client) UpdateDownloadClient(ctx context.Context, dc *DownloadClientResource) (*DownloadClientResource, error) {
	var out DownloadClientResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/downloadclient/%d", dc.ID), dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDownloadClient deletes a download client by ID.
func (c *Client) DeleteDownloadClient(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/downloadclient/%d", id), nil, nil)
}

// GetDownloadClientSchema returns the schema for all download client types.
func (c *Client) GetDownloadClientSchema(ctx context.Context) ([]DownloadClientResource, error) {
	var out []DownloadClientResource
	if err := c.base.Get(ctx, "/api/v1/downloadclient/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestDownloadClient tests a download client configuration.
func (c *Client) TestDownloadClient(ctx context.Context, dc *DownloadClientResource) error {
	return c.base.Post(ctx, "/api/v1/downloadclient/test", dc, nil)
}

// TestAllDownloadClients tests all configured download clients.
func (c *Client) TestAllDownloadClients(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/downloadclient/testall", nil, nil)
}

// BulkUpdateDownloadClients updates multiple download clients at once.
func (c *Client) BulkUpdateDownloadClients(ctx context.Context, bulk *arr.ProviderBulkResource) ([]DownloadClientResource, error) {
	var out []DownloadClientResource
	if err := c.base.Put(ctx, "/api/v1/downloadclient/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteDownloadClients deletes multiple download clients at once.
func (c *Client) BulkDeleteDownloadClients(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v1/downloadclient/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// DownloadClientAction triggers a named action on a download client provider.
func (c *Client) DownloadClientAction(ctx context.Context, name string, body *DownloadClientResource) error {
	path := "/api/v1/downloadclient/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Notifications ----------.

// GetNotifications returns all configured notifications.
func (c *Client) GetNotifications(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/notification", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetNotification returns a single notification by its ID.
func (c *Client) GetNotification(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/notification/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateNotification creates a new notification.
func (c *Client) CreateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v1/notification", n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNotification updates an existing notification.
func (c *Client) UpdateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/notification/%d", n.ID), n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNotification deletes a notification by ID.
func (c *Client) DeleteNotification(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/notification/%d", id), nil, nil)
}

// GetNotificationSchema returns the schema for all notification types.
func (c *Client) GetNotificationSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/notification/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestNotification tests a notification configuration.
func (c *Client) TestNotification(ctx context.Context, n *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v1/notification/test", n, nil)
}

// TestAllNotifications tests all configured notifications.
func (c *Client) TestAllNotifications(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/notification/testall", nil, nil)
}

// NotificationAction triggers a named action on a notification provider.
func (c *Client) NotificationAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v1/notification/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Config Endpoints ----------.

// GetDownloadClientConfig returns the download client configuration.
func (c *Client) GetDownloadClientConfig(ctx context.Context) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/downloadclient", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDownloadClientConfigByID returns the download client config by its ID.
func (c *Client) GetDownloadClientConfigByID(ctx context.Context, id int) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/downloadclient/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClientConfig updates the download client configuration.
func (c *Client) UpdateDownloadClientConfig(ctx context.Context, config *arr.DownloadClientConfigResource) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/downloadclient/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHostConfig returns the host configuration.
func (c *Client) GetHostConfig(ctx context.Context) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/host", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHostConfigByID returns the host config by its ID.
func (c *Client) GetHostConfigByID(ctx context.Context, id int) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/host/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateHostConfig updates the host configuration.
func (c *Client) UpdateHostConfig(ctx context.Context, config *arr.HostConfigResource) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/host/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUIConfig returns the UI configuration.
func (c *Client) GetUIConfig(ctx context.Context) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/ui", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUIConfigByID returns the UI config by its ID.
func (c *Client) GetUIConfigByID(ctx context.Context, id int) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/ui/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateUIConfig updates the UI configuration.
func (c *Client) UpdateUIConfig(ctx context.Context, config *arr.UIConfigResource) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/ui/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Development Config ----------.

// GetDevelopmentConfig returns the development configuration.
func (c *Client) GetDevelopmentConfig(ctx context.Context) (*DevelopmentConfigResource, error) {
	var out DevelopmentConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/development", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDevelopmentConfigByID returns the development config by its ID.
func (c *Client) GetDevelopmentConfigByID(ctx context.Context, id int) (*DevelopmentConfigResource, error) {
	var out DevelopmentConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/development/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDevelopmentConfig updates the development configuration.
func (c *Client) UpdateDevelopmentConfig(ctx context.Context, config *DevelopmentConfigResource) (*DevelopmentConfigResource, error) {
	var out DevelopmentConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/development/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Custom Filters ----------.

// GetCustomFilters returns all custom filters.
func (c *Client) GetCustomFilters(ctx context.Context) ([]arr.CustomFilterResource, error) {
	var out []arr.CustomFilterResource
	if err := c.base.Get(ctx, "/api/v1/customfilter", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCustomFilter returns a single custom filter by its ID.
func (c *Client) GetCustomFilter(ctx context.Context, id int) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/customfilter/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCustomFilter creates a new custom filter.
func (c *Client) CreateCustomFilter(ctx context.Context, filter *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Post(ctx, "/api/v1/customfilter", filter, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFilter updates an existing custom filter.
func (c *Client) UpdateCustomFilter(ctx context.Context, filter *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/customfilter/%d", filter.ID), filter, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFilter deletes a custom filter by ID.
func (c *Client) DeleteCustomFilter(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/customfilter/%d", id), nil, nil)
}

// ---------- Tags Extended ----------.

// GetTag returns a single tag by its ID.
func (c *Client) GetTag(ctx context.Context, id int) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/tag/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTag updates an existing tag.
func (c *Client) UpdateTag(ctx context.Context, tag *arr.Tag) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/tag/%d", tag.ID), tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTag deletes a tag by ID.
func (c *Client) DeleteTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/tag/%d", id), nil, nil)
}

// GetTagDetails returns all tags with details about their usage.
func (c *Client) GetTagDetails(ctx context.Context) ([]arr.TagDetail, error) {
	var out []arr.TagDetail
	if err := c.base.Get(ctx, "/api/v1/tag/detail", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTagDetail returns a single tag detail by its ID.
func (c *Client) GetTagDetail(ctx context.Context, id int) (*arr.TagDetail, error) {
	var out arr.TagDetail
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/tag/detail/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Backups ----------.

// GetBackups returns a list of all available backups.
func (c *Client) GetBackups(ctx context.Context) ([]arr.Backup, error) {
	var out []arr.Backup
	if err := c.base.Get(ctx, "/api/v1/system/backup", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBackup deletes a backup by ID.
func (c *Client) DeleteBackup(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/system/backup/%d", id), nil, nil)
}

// RestoreBackup triggers a restore from a backup by ID.
func (c *Client) RestoreBackup(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v1/system/backup/restore/%d", id), nil, nil)
}

// ---------- Logs ----------.

// GetLogs returns log entries with pagination.
func (c *Client) GetLogs(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.LogRecord], error) {
	var out arr.PagingResource[arr.LogRecord]
	path := fmt.Sprintf("/api/v1/log?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLogFiles returns a list of log files.
func (c *Client) GetLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v1/log/file", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLogFileContent returns the content of a specific log file by filename.
func (c *Client) GetLogFileContent(ctx context.Context, filename string) (string, error) {
	path := "/api/v1/log/file/" + url.PathEscape(filename)
	b, err := c.base.GetRaw(ctx, path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GetUpdateLogFiles returns a list of update log files.
func (c *Client) GetUpdateLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v1/log/file/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetUpdateLogFileContent returns the content of a specific update log file.
func (c *Client) GetUpdateLogFileContent(ctx context.Context, filename string) (string, error) {
	path := "/api/v1/log/file/update/" + url.PathEscape(filename)
	b, err := c.base.GetRaw(ctx, path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------- System ----------.

// GetTasks returns all scheduled tasks.
func (c *Client) GetTasks(ctx context.Context) ([]arr.TaskResource, error) {
	var out []arr.TaskResource
	if err := c.base.Get(ctx, "/api/v1/system/task", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTask returns a single task by its ID.
func (c *Client) GetTask(ctx context.Context, id int) (*arr.TaskResource, error) {
	var out arr.TaskResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/system/task/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUpdates returns available application updates.
func (c *Client) GetUpdates(ctx context.Context) ([]arr.UpdateResource, error) {
	var out []arr.UpdateResource
	if err := c.base.Get(ctx, "/api/v1/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSystemRoutes returns all API routes.
func (c *Client) GetSystemRoutes(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v1/system/routes", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSystemRoutesDuplicate returns duplicate API routes.
func (c *Client) GetSystemRoutesDuplicate(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v1/system/routes/duplicate", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Shutdown sends a shutdown command to Prowlarr.
func (c *Client) Shutdown(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/system/shutdown", nil, nil)
}

// Restart sends a restart command to Prowlarr.
func (c *Client) Restart(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/system/restart", nil, nil)
}

// DeleteCommand deletes a command by ID.
func (c *Client) DeleteCommand(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/command/%d", id), nil, nil)
}

// ---------- Localization ----------.

// GetLocalization returns the localization strings.
func (c *Client) GetLocalization(ctx context.Context) (map[string]string, error) {
	var out map[string]string
	if err := c.base.Get(ctx, "/api/v1/localization", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLocalizationOptions returns available localization language options.
func (c *Client) GetLocalizationOptions(ctx context.Context) ([]LocalizationOption, error) {
	var out []LocalizationOption
	if err := c.base.Get(ctx, "/api/v1/localization/options", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Ping ----------.

// Ping checks connectivity to the Prowlarr instance.
func (c *Client) Ping(ctx context.Context) error {
	return c.base.Get(ctx, "/ping", nil)
}

// ---------- File System ----------.

// BrowseFileSystem returns directory/file listings for the given path.
func (c *Client) BrowseFileSystem(ctx context.Context, path string) (map[string]any, error) {
	var out map[string]any
	reqPath := "/api/v1/filesystem?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFileSystemType returns the filesystem type for the given path.
func (c *Client) GetFileSystemType(ctx context.Context, path string) (string, error) {
	var out string
	reqPath := "/api/v1/filesystem/type?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return "", err
	}
	return out, nil
}

// ---------- Search Extended ----------.

// GrabReleasesBulk sends multiple releases to the download client.
func (c *Client) GrabReleasesBulk(ctx context.Context, releases []Release) (*Release, error) {
	var out Release
	if err := c.base.Post(ctx, "/api/v1/search/bulk", releases, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Newznab ----------.

// GetIndexerNewznab returns the Newznab/Torznab XML feed for an indexer.
func (c *Client) GetIndexerNewznab(ctx context.Context, id int) (string, error) {
	var out string
	path := fmt.Sprintf("/api/v1/indexer/%d/newznab", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return "", err
	}
	return out, nil
}

// DownloadIndexerRelease downloads a release from an indexer by ID.
func (c *Client) DownloadIndexerRelease(ctx context.Context, id int) (string, error) {
	var out string
	path := fmt.Sprintf("/api/v1/indexer/%d/download", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return "", err
	}
	return out, nil
}

// ---------- HEAD Ping ----------.

// HeadPing performs a lightweight HEAD request to /ping.
func (c *Client) HeadPing(ctx context.Context) error {
	return c.base.Head(ctx, "/ping")
}

// ---------- Backup Upload ----------.

// UploadBackup uploads a backup file via multipart form POST.
func (c *Client) UploadBackup(ctx context.Context, fileName string, data io.Reader) error {
	return c.base.Upload(ctx, "/api/v1/system/backup/upload", "file", fileName, data)
}
