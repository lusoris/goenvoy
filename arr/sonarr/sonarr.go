package sonarr

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/golusoris/goenvoy/arr/v2"
)

// Client is a Sonarr API client.
type Client struct {
	base *arr.BaseClient
}

// New creates a Sonarr [Client] for the instance at baseURL.
func New(baseURL, apiKey string, opts ...arr.Option) (*Client, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{base: base}, nil
}

// GetAllSeries returns every series configured in Sonarr.
func (c *Client) GetAllSeries(ctx context.Context) ([]Series, error) {
	var out []Series
	if err := c.base.Get(ctx, "/api/v3/series", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSeries returns a single series by its database ID.
func (c *Client) GetSeries(ctx context.Context, id int) (*Series, error) {
	var out Series
	path := fmt.Sprintf("/api/v3/series/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddSeries adds a new series to Sonarr.
func (c *Client) AddSeries(ctx context.Context, series *Series) (*Series, error) {
	var out Series
	if err := c.base.Post(ctx, "/api/v3/series", series, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSeries updates an existing series. Set moveFiles to true to relocate
// files when the series path changes.
func (c *Client) UpdateSeries(ctx context.Context, series *Series, moveFiles bool) (*Series, error) {
	var out Series
	path := fmt.Sprintf("/api/v3/series/%d?moveFiles=%t", series.ID, moveFiles)
	if err := c.base.Put(ctx, path, series, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteSeries removes a series. Set deleteFiles to true to also delete
// downloaded episode files from disk.
func (c *Client) DeleteSeries(ctx context.Context, id int, deleteFiles, addImportListExclusion bool) error {
	path := fmt.Sprintf("/api/v3/series/%d?deleteFiles=%t&addImportListExclusion=%t", id, deleteFiles, addImportListExclusion)
	return c.base.Delete(ctx, path, nil, nil)
}

// LookupSeries searches for a series by term (title or TVDB ID slug).
func (c *Client) LookupSeries(ctx context.Context, term string) ([]Series, error) {
	var out []Series
	path := "/api/v3/series/lookup?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEpisodes returns all episodes for the given series.
func (c *Client) GetEpisodes(ctx context.Context, seriesID int) ([]Episode, error) {
	var out []Episode
	path := fmt.Sprintf("/api/v3/episode?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEpisode returns a single episode by its database ID.
func (c *Client) GetEpisode(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	path := fmt.Sprintf("/api/v3/episode/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateEpisode updates the metadata for an episode (typically the monitored flag).
func (c *Client) UpdateEpisode(ctx context.Context, episode *Episode) (*Episode, error) {
	var out Episode
	path := fmt.Sprintf("/api/v3/episode/%d", episode.ID)
	if err := c.base.Put(ctx, path, episode, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MonitorEpisodes sets the monitored flag for a batch of episode IDs.
func (c *Client) MonitorEpisodes(ctx context.Context, episodeIDs []int, monitored bool) error {
	body := EpisodesMonitoredResource{EpisodeIDs: episodeIDs, Monitored: monitored}
	return c.base.Put(ctx, "/api/v3/episode/monitor", body, nil)
}

// GetEpisodeFiles returns all episode files for the given series.
func (c *Client) GetEpisodeFiles(ctx context.Context, seriesID int) ([]EpisodeFile, error) {
	var out []EpisodeFile
	path := fmt.Sprintf("/api/v3/episodefile?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEpisodeFile returns a single episode file by its database ID.
func (c *Client) GetEpisodeFile(ctx context.Context, id int) (*EpisodeFile, error) {
	var out EpisodeFile
	path := fmt.Sprintf("/api/v3/episodefile/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteEpisodeFile removes a single episode file by its database ID.
func (c *Client) DeleteEpisodeFile(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v3/episodefile/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// DeleteEpisodeFiles removes multiple episode files by their IDs.
func (c *Client) DeleteEpisodeFiles(ctx context.Context, ids []int) error {
	body := EpisodeFileListResource{EpisodeFileIDs: ids}
	return c.base.Delete(ctx, "/api/v3/episodefile/bulk", &body, nil)
}

// SendCommand triggers a named command (e.g. "RefreshSeries", "EpisodeSearch").
func (c *Client) SendCommand(ctx context.Context, cmd arr.CommandRequest) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Post(ctx, "/api/v3/command", cmd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCommands returns all currently queued or running commands.
func (c *Client) GetCommands(ctx context.Context) ([]arr.CommandResponse, error) {
	var out []arr.CommandResponse
	if err := c.base.Get(ctx, "/api/v3/command", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCommand returns the status of a single command by its ID.
func (c *Client) GetCommand(ctx context.Context, id int) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	path := fmt.Sprintf("/api/v3/command/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCalendar returns episodes airing between start and end (RFC 3339 timestamps).
func (c *Client) GetCalendar(ctx context.Context, start, end string, unmonitored bool) ([]Episode, error) {
	var out []Episode
	path := fmt.Sprintf("/api/v3/calendar?start=%s&end=%s&unmonitored=%t",
		url.QueryEscape(start), url.QueryEscape(end), unmonitored)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Parse parses a release title and returns the extracted information.
func (c *Client) Parse(ctx context.Context, title string) (*ParseResult, error) {
	var out ParseResult
	path := "/api/v3/parse?title=" + url.QueryEscape(title)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSystemStatus returns Sonarr system information.
func (c *Client) GetSystemStatus(ctx context.Context) (*arr.StatusResponse, error) {
	var out arr.StatusResponse
	if err := c.base.Get(ctx, "/api/v3/system/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHealth returns current health check results.
func (c *Client) GetHealth(ctx context.Context) ([]arr.HealthCheck, error) {
	var out []arr.HealthCheck
	if err := c.base.Get(ctx, "/api/v3/health", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDiskSpace returns disk usage information for configured paths.
func (c *Client) GetDiskSpace(ctx context.Context) ([]arr.DiskSpace, error) {
	var out []arr.DiskSpace
	if err := c.base.Get(ctx, "/api/v3/diskspace", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueue returns the current download queue with pagination.
func (c *Client) GetQueue(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.QueueRecord], error) {
	var out arr.PagingResource[arr.QueueRecord]
	path := fmt.Sprintf("/api/v3/queue?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQueueItem removes an item from the download queue.
func (c *Client) DeleteQueueItem(ctx context.Context, id int, removeFromClient, blocklist bool) error {
	path := fmt.Sprintf("/api/v3/queue/%d?removeFromClient=%t&blocklist=%t", id, removeFromClient, blocklist)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetQualityProfiles returns all configured quality profiles.
func (c *Client) GetQualityProfiles(ctx context.Context) ([]arr.QualityProfile, error) {
	var out []arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v3/qualityprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTags returns all tags.
func (c *Client) GetTags(ctx context.Context) ([]arr.Tag, error) {
	var out []arr.Tag
	if err := c.base.Get(ctx, "/api/v3/tag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTag creates a new tag and returns it with its assigned ID.
func (c *Client) CreateTag(ctx context.Context, label string) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Post(ctx, "/api/v3/tag", arr.Tag{Label: label}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRootFolders returns all configured root folders.
func (c *Client) GetRootFolders(ctx context.Context) ([]arr.RootFolder, error) {
	var out []arr.RootFolder
	if err := c.base.Get(ctx, "/api/v3/rootfolder", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistory returns the download history with pagination.
func (c *Client) GetHistory(ctx context.Context, page, pageSize int) (*arr.PagingResource[HistoryRecord], error) {
	var out arr.PagingResource[HistoryRecord]
	path := fmt.Sprintf("/api/v3/history?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSeasonPass updates monitored status for multiple series and seasons.
func (c *Client) UpdateSeasonPass(ctx context.Context, pass SeasonPassResource) error {
	return c.base.Post(ctx, "/api/v3/seasonpass", pass, nil)
}

// ---------- Series Editor ----------.

// EditSeries applies bulk edits to multiple series.
func (c *Client) EditSeries(ctx context.Context, editor *SeriesEditorResource) ([]Series, error) {
	var out []Series
	if err := c.base.Put(ctx, "/api/v3/series/editor", editor, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteManySeries deletes multiple series at once.
func (c *Client) DeleteManySeries(ctx context.Context, editor *SeriesEditorResource) error {
	return c.base.Delete(ctx, "/api/v3/series/editor", editor, nil)
}

// ImportSeries imports one or more series in bulk.
func (c *Client) ImportSeries(ctx context.Context, series []Series) ([]Series, error) {
	var out []Series
	if err := c.base.Post(ctx, "/api/v3/series/import", series, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Wanted ----------.

// GetWantedMissing returns episodes that are monitored but missing files.
func (c *Client) GetWantedMissing(ctx context.Context, page, pageSize int) (*arr.PagingResource[Episode], error) {
	var out arr.PagingResource[Episode]
	path := fmt.Sprintf("/api/v3/wanted/missing?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedCutoff returns episodes that have not met the quality cutoff.
func (c *Client) GetWantedCutoff(ctx context.Context, page, pageSize int) (*arr.PagingResource[Episode], error) {
	var out arr.PagingResource[Episode]
	path := fmt.Sprintf("/api/v3/wanted/cutoff?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Blocklist ----------.

// GetBlocklist returns the blocklisted releases with pagination.
func (c *Client) GetBlocklist(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.BlocklistResource], error) {
	var out arr.PagingResource[arr.BlocklistResource]
	path := fmt.Sprintf("/api/v3/blocklist?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBlocklistItem removes a single blocklist entry.
func (c *Client) DeleteBlocklistItem(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/blocklist/%d", id), nil, nil)
}

// DeleteBlocklistBulk removes multiple blocklist entries at once.
func (c *Client) DeleteBlocklistBulk(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/blocklist/bulk", &arr.BlocklistBulkResource{IDs: ids}, nil)
}

// ---------- Custom Filters ----------.

// GetCustomFilters returns all custom filters.
func (c *Client) GetCustomFilters(ctx context.Context) ([]arr.CustomFilterResource, error) {
	var out []arr.CustomFilterResource
	if err := c.base.Get(ctx, "/api/v3/customfilter", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCustomFilter returns a single custom filter by ID.
func (c *Client) GetCustomFilter(ctx context.Context, id int) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/customfilter/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCustomFilter creates a new custom filter.
func (c *Client) CreateCustomFilter(ctx context.Context, filter *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Post(ctx, "/api/v3/customfilter", filter, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFilter updates an existing custom filter.
func (c *Client) UpdateCustomFilter(ctx context.Context, filter *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/customfilter/%d", filter.ID), filter, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFilter removes a custom filter.
func (c *Client) DeleteCustomFilter(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/customfilter/%d", id), nil, nil)
}

// ---------- Custom Formats ----------.

// GetCustomFormats returns all custom formats.
func (c *Client) GetCustomFormats(ctx context.Context) ([]arr.CustomFormatResource, error) {
	var out []arr.CustomFormatResource
	if err := c.base.Get(ctx, "/api/v3/customformat", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCustomFormat returns a single custom format by ID.
func (c *Client) GetCustomFormat(ctx context.Context, id int) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/customformat/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCustomFormat creates a new custom format.
func (c *Client) CreateCustomFormat(ctx context.Context, cf *arr.CustomFormatResource) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Post(ctx, "/api/v3/customformat", cf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFormat updates an existing custom format.
func (c *Client) UpdateCustomFormat(ctx context.Context, cf *arr.CustomFormatResource) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/customformat/%d", cf.ID), cf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFormat removes a custom format.
func (c *Client) DeleteCustomFormat(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/customformat/%d", id), nil, nil)
}

// GetCustomFormatSchema returns the schema for custom format specifications.
func (c *Client) GetCustomFormatSchema(ctx context.Context) ([]arr.CustomFormatSpecification, error) {
	var out []arr.CustomFormatSpecification
	if err := c.base.Get(ctx, "/api/v3/customformat/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Delay Profiles ----------.

// GetDelayProfiles returns all delay profiles.
func (c *Client) GetDelayProfiles(ctx context.Context) ([]arr.DelayProfileResource, error) {
	var out []arr.DelayProfileResource
	if err := c.base.Get(ctx, "/api/v3/delayprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDelayProfile returns a single delay profile by ID.
func (c *Client) GetDelayProfile(ctx context.Context, id int) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/delayprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDelayProfile creates a new delay profile.
func (c *Client) CreateDelayProfile(ctx context.Context, dp *arr.DelayProfileResource) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Post(ctx, "/api/v3/delayprofile", dp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDelayProfile updates an existing delay profile.
func (c *Client) UpdateDelayProfile(ctx context.Context, dp *arr.DelayProfileResource) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/delayprofile/%d", dp.ID), dp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDelayProfile removes a delay profile.
func (c *Client) DeleteDelayProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/delayprofile/%d", id), nil, nil)
}

// ReorderDelayProfile moves a delay profile to a new position.
func (c *Client) ReorderDelayProfile(ctx context.Context, id, afterID int) ([]arr.DelayProfileResource, error) {
	var out []arr.DelayProfileResource
	path := fmt.Sprintf("/api/v3/delayprofile/reorder/%d?after=%d", id, afterID)
	if err := c.base.Put(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Notifications ----------.

// GetNotifications returns all notification configurations.
func (c *Client) GetNotifications(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/notification", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetNotification returns a single notification by ID.
func (c *Client) GetNotification(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/notification/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateNotification creates a new notification.
func (c *Client) CreateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/notification", n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNotification updates an existing notification.
func (c *Client) UpdateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/notification/%d", n.ID), n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNotification removes a notification.
func (c *Client) DeleteNotification(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/notification/%d", id), nil, nil)
}

// GetNotificationSchema returns available notification implementations.
func (c *Client) GetNotificationSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/notification/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestNotification sends a test for a notification configuration.
func (c *Client) TestNotification(ctx context.Context, n *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/notification/test", n, nil)
}

// ---------- Download Clients ----------.

// GetDownloadClients returns all download client configurations.
func (c *Client) GetDownloadClients(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/downloadclient", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDownloadClient returns a single download client by ID.
func (c *Client) GetDownloadClient(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/downloadclient/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDownloadClient creates a new download client.
func (c *Client) CreateDownloadClient(ctx context.Context, dc *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/downloadclient", dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClient updates an existing download client.
func (c *Client) UpdateDownloadClient(ctx context.Context, dc *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/downloadclient/%d", dc.ID), dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDownloadClient removes a download client.
func (c *Client) DeleteDownloadClient(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/downloadclient/%d", id), nil, nil)
}

// GetDownloadClientSchema returns available download client implementations.
func (c *Client) GetDownloadClientSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/downloadclient/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestDownloadClient sends a test for a download client configuration.
func (c *Client) TestDownloadClient(ctx context.Context, dc *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/test", dc, nil)
}

// GetDownloadClientConfig returns the download client global configuration.
func (c *Client) GetDownloadClientConfig(ctx context.Context) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/downloadclient", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClientConfig updates the download client global configuration.
func (c *Client) UpdateDownloadClientConfig(ctx context.Context, cfg *arr.DownloadClientConfigResource) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/downloadclient/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Indexers ----------.

// GetIndexers returns all indexer configurations.
func (c *Client) GetIndexers(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/indexer", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetIndexer returns a single indexer by ID.
func (c *Client) GetIndexer(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/indexer/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateIndexer creates a new indexer.
func (c *Client) CreateIndexer(ctx context.Context, idx *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/indexer", idx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexer updates an existing indexer.
func (c *Client) UpdateIndexer(ctx context.Context, idx *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/indexer/%d", idx.ID), idx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteIndexer removes an indexer.
func (c *Client) DeleteIndexer(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/indexer/%d", id), nil, nil)
}

// GetIndexerSchema returns available indexer implementations.
func (c *Client) GetIndexerSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/indexer/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestIndexer sends a test for an indexer configuration.
func (c *Client) TestIndexer(ctx context.Context, idx *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/indexer/test", idx, nil)
}

// GetIndexerConfig returns the indexer global configuration.
func (c *Client) GetIndexerConfig(ctx context.Context) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/indexer", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexerConfig updates the indexer global configuration.
func (c *Client) UpdateIndexerConfig(ctx context.Context, cfg *arr.IndexerConfigResource) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/indexer/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIndexerFlags returns the list of indexer flags.
func (c *Client) GetIndexerFlags(ctx context.Context) ([]arr.IndexerFlagResource, error) {
	var out []arr.IndexerFlagResource
	if err := c.base.Get(ctx, "/api/v3/indexerflag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Import Lists ----------.

// GetImportLists returns all import list configurations.
func (c *Client) GetImportLists(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/importlist", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetImportList returns a single import list by ID.
func (c *Client) GetImportList(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/importlist/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateImportList creates a new import list.
func (c *Client) CreateImportList(ctx context.Context, il *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/importlist", il, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportList updates an existing import list.
func (c *Client) UpdateImportList(ctx context.Context, il *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/importlist/%d", il.ID), il, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportList removes an import list.
func (c *Client) DeleteImportList(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/importlist/%d", id), nil, nil)
}

// GetImportListSchema returns available import list implementations.
func (c *Client) GetImportListSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/importlist/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestImportList sends a test for an import list configuration.
func (c *Client) TestImportList(ctx context.Context, il *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/importlist/test", il, nil)
}

// ---------- Import List Exclusions ----------.

// GetImportListExclusions returns all import list exclusions.
func (c *Client) GetImportListExclusions(ctx context.Context) ([]arr.ImportListExclusionResource, error) {
	var out []arr.ImportListExclusionResource
	if err := c.base.Get(ctx, "/api/v3/importlistexclusion", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetImportListExclusion returns a single import list exclusion by ID.
func (c *Client) GetImportListExclusion(ctx context.Context, id int) (*arr.ImportListExclusionResource, error) {
	var out arr.ImportListExclusionResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/importlistexclusion/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateImportListExclusion creates a new import list exclusion.
func (c *Client) CreateImportListExclusion(ctx context.Context, excl *arr.ImportListExclusionResource) (*arr.ImportListExclusionResource, error) {
	var out arr.ImportListExclusionResource
	if err := c.base.Post(ctx, "/api/v3/importlistexclusion", excl, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportListExclusion updates an existing import list exclusion.
func (c *Client) UpdateImportListExclusion(ctx context.Context, excl *arr.ImportListExclusionResource) (*arr.ImportListExclusionResource, error) {
	var out arr.ImportListExclusionResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/importlistexclusion/%d", excl.ID), excl, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportListExclusion removes an import list exclusion.
func (c *Client) DeleteImportListExclusion(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/importlistexclusion/%d", id), nil, nil)
}

// ---------- Metadata ----------.

// GetMetadataConsumers returns all metadata consumer configurations.
func (c *Client) GetMetadataConsumers(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/metadata", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMetadataConsumer returns a single metadata consumer by ID.
func (c *Client) GetMetadataConsumer(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/metadata/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateMetadataConsumer creates a new metadata consumer.
func (c *Client) CreateMetadataConsumer(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/metadata", m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMetadataConsumer updates an existing metadata consumer.
func (c *Client) UpdateMetadataConsumer(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/metadata/%d", m.ID), m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMetadataConsumer removes a metadata consumer.
func (c *Client) DeleteMetadataConsumer(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/metadata/%d", id), nil, nil)
}

// GetMetadataSchema returns available metadata consumer implementations.
func (c *Client) GetMetadataSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/metadata/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestMetadataConsumer sends a test for a metadata consumer configuration.
func (c *Client) TestMetadataConsumer(ctx context.Context, m *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/metadata/test", m, nil)
}

// ---------- Auto Tagging ----------.

// GetAutoTagging returns all auto-tag rules.
func (c *Client) GetAutoTagging(ctx context.Context) ([]arr.AutoTaggingResource, error) {
	var out []arr.AutoTaggingResource
	if err := c.base.Get(ctx, "/api/v3/autotagging", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAutoTag returns a single auto-tag rule by ID.
func (c *Client) GetAutoTag(ctx context.Context, id int) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/autotagging/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateAutoTag creates a new auto-tag rule.
func (c *Client) CreateAutoTag(ctx context.Context, at *arr.AutoTaggingResource) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Post(ctx, "/api/v3/autotagging", at, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAutoTag updates an existing auto-tag rule.
func (c *Client) UpdateAutoTag(ctx context.Context, at *arr.AutoTaggingResource) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/autotagging/%d", at.ID), at, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAutoTag removes an auto-tag rule.
func (c *Client) DeleteAutoTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/autotagging/%d", id), nil, nil)
}

// GetAutoTagSchema returns available auto-tag specification implementations.
func (c *Client) GetAutoTagSchema(ctx context.Context) ([]arr.AutoTaggingSpecification, error) {
	var out []arr.AutoTaggingSpecification
	if err := c.base.Get(ctx, "/api/v3/autotagging/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Backup ----------.

// GetBackups returns all available backup files.
func (c *Client) GetBackups(ctx context.Context) ([]arr.Backup, error) {
	var out []arr.Backup
	if err := c.base.Get(ctx, "/api/v3/system/backup", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBackup removes a backup file by ID.
func (c *Client) DeleteBackup(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/system/backup/%d", id), nil, nil)
}

// RestoreBackup triggers a restore from a backup by ID.
func (c *Client) RestoreBackup(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/system/backup/restore/%d", id), nil, nil)
}

// ---------- Quality Profiles (CRUD) ----------.

// GetQualityProfile returns a single quality profile by ID.
func (c *Client) GetQualityProfile(ctx context.Context, id int) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/qualityprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateQualityProfile creates a new quality profile.
func (c *Client) CreateQualityProfile(ctx context.Context, qp *arr.QualityProfile) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Post(ctx, "/api/v3/qualityprofile", qp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateQualityProfile updates an existing quality profile.
func (c *Client) UpdateQualityProfile(ctx context.Context, qp *arr.QualityProfile) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/qualityprofile/%d", qp.ID), qp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQualityProfile removes a quality profile.
func (c *Client) DeleteQualityProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/qualityprofile/%d", id), nil, nil)
}

// ---------- Quality Definitions ----------.

// GetQualityDefinitions returns all quality definitions with size limits.
func (c *Client) GetQualityDefinitions(ctx context.Context) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Get(ctx, "/api/v3/qualitydefinition", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQualityDefinition returns a single quality definition by ID.
func (c *Client) GetQualityDefinition(ctx context.Context, id int) (*arr.QualityDefinitionResource, error) {
	var out arr.QualityDefinitionResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/qualitydefinition/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateQualityDefinition updates a quality definition.
func (c *Client) UpdateQualityDefinition(ctx context.Context, qd *arr.QualityDefinitionResource) (*arr.QualityDefinitionResource, error) {
	var out arr.QualityDefinitionResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/qualitydefinition/%d", qd.ID), qd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Tags (full CRUD) ----------.

// GetTag returns a single tag by ID.
func (c *Client) GetTag(ctx context.Context, id int) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/tag/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTag updates an existing tag.
func (c *Client) UpdateTag(ctx context.Context, tag *arr.Tag) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/tag/%d", tag.ID), tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTag removes a tag.
func (c *Client) DeleteTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/tag/%d", id), nil, nil)
}

// GetTagDetails returns all tags with details about which resources use them.
func (c *Client) GetTagDetails(ctx context.Context) ([]arr.TagDetail, error) {
	var out []arr.TagDetail
	if err := c.base.Get(ctx, "/api/v3/tag/detail", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTagDetail returns a single tag detail by ID.
func (c *Client) GetTagDetail(ctx context.Context, id int) (*arr.TagDetail, error) {
	var out arr.TagDetail
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/tag/detail/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Root Folders (full CRUD) ----------.

// GetRootFolder returns a single root folder by ID.
func (c *Client) GetRootFolder(ctx context.Context, id int) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/rootfolder/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRootFolder adds a new root folder.
func (c *Client) CreateRootFolder(ctx context.Context, rf *arr.RootFolder) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Post(ctx, "/api/v3/rootfolder", rf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRootFolder removes a root folder.
func (c *Client) DeleteRootFolder(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/rootfolder/%d", id), nil, nil)
}

// ---------- Release Profiles ----------.

// GetReleaseProfiles returns all release profiles.
func (c *Client) GetReleaseProfiles(ctx context.Context) ([]arr.ReleaseProfileResource, error) {
	var out []arr.ReleaseProfileResource
	if err := c.base.Get(ctx, "/api/v3/releaseprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetReleaseProfile returns a single release profile by ID.
func (c *Client) GetReleaseProfile(ctx context.Context, id int) (*arr.ReleaseProfileResource, error) {
	var out arr.ReleaseProfileResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/releaseprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateReleaseProfile creates a new release profile.
func (c *Client) CreateReleaseProfile(ctx context.Context, rp *arr.ReleaseProfileResource) (*arr.ReleaseProfileResource, error) {
	var out arr.ReleaseProfileResource
	if err := c.base.Post(ctx, "/api/v3/releaseprofile", rp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateReleaseProfile updates an existing release profile.
func (c *Client) UpdateReleaseProfile(ctx context.Context, rp *arr.ReleaseProfileResource) (*arr.ReleaseProfileResource, error) {
	var out arr.ReleaseProfileResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/releaseprofile/%d", rp.ID), rp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteReleaseProfile removes a release profile.
func (c *Client) DeleteReleaseProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/releaseprofile/%d", id), nil, nil)
}

// ---------- Remote Path Mappings ----------.

// GetRemotePathMappings returns all remote path mappings.
func (c *Client) GetRemotePathMappings(ctx context.Context) ([]arr.RemotePathMappingResource, error) {
	var out []arr.RemotePathMappingResource
	if err := c.base.Get(ctx, "/api/v3/remotepathmapping", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetRemotePathMapping returns a single remote path mapping by ID.
func (c *Client) GetRemotePathMapping(ctx context.Context, id int) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/remotepathmapping/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRemotePathMapping creates a new remote path mapping.
func (c *Client) CreateRemotePathMapping(ctx context.Context, rpm *arr.RemotePathMappingResource) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Post(ctx, "/api/v3/remotepathmapping", rpm, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateRemotePathMapping updates an existing remote path mapping.
func (c *Client) UpdateRemotePathMapping(ctx context.Context, rpm *arr.RemotePathMappingResource) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/remotepathmapping/%d", rpm.ID), rpm, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRemotePathMapping removes a remote path mapping.
func (c *Client) DeleteRemotePathMapping(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/remotepathmapping/%d", id), nil, nil)
}

// ---------- Releases ----------.

// SearchReleases searches for releases matching the given episode IDs.
func (c *Client) SearchReleases(ctx context.Context, episodeID int) ([]arr.ReleaseResource, error) {
	var out []arr.ReleaseResource
	path := fmt.Sprintf("/api/v3/release?episodeId=%d", episodeID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PushRelease manually pushes a download for a release.
func (c *Client) PushRelease(ctx context.Context, release *arr.ReleasePushResource) error {
	return c.base.Post(ctx, "/api/v3/release/push", release, nil)
}

// GrabRelease grabs a release by its GUID.
func (c *Client) GrabRelease(ctx context.Context, guid string, indexerID int) error {
	body := map[string]any{"guid": guid, "indexerId": indexerID}
	return c.base.Post(ctx, "/api/v3/release", body, nil)
}

// ---------- Rename ----------.

// GetRenameList returns proposed renames for a series.
func (c *Client) GetRenameList(ctx context.Context, seriesID int) ([]arr.RenameEpisodeResource, error) {
	var out []arr.RenameEpisodeResource
	path := fmt.Sprintf("/api/v3/rename?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Manual Import ----------.

// GetManualImport returns files available for manual import.
func (c *Client) GetManualImport(ctx context.Context, folder, downloadID string) ([]arr.ManualImportResource, error) {
	var out []arr.ManualImportResource
	path := fmt.Sprintf("/api/v3/manualimport?folder=%s&downloadId=%s",
		url.QueryEscape(folder), url.QueryEscape(downloadID))
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessManualImport confirms and processes a manual import.
func (c *Client) ProcessManualImport(ctx context.Context, imports []arr.ManualImportReprocessResource) error {
	return c.base.Post(ctx, "/api/v3/manualimport", imports, nil)
}

// ---------- Logs ----------.

// GetLogs returns log entries with pagination.
func (c *Client) GetLogs(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.LogRecord], error) {
	var out arr.PagingResource[arr.LogRecord]
	path := fmt.Sprintf("/api/v3/log?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLogFiles returns the list of available log files.
func (c *Client) GetLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v3/log/file", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetUpdateLogFiles returns the list of available update log files.
func (c *Client) GetUpdateLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v3/log/file/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Config: Naming ----------.

// GetNamingConfig returns the file naming configuration.
func (c *Client) GetNamingConfig(ctx context.Context) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/naming", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNamingConfig updates the file naming configuration.
func (c *Client) UpdateNamingConfig(ctx context.Context, cfg *arr.NamingConfigResource) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/naming/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Config: Host ----------.

// GetHostConfig returns the host configuration.
func (c *Client) GetHostConfig(ctx context.Context) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/host", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateHostConfig updates the host configuration.
func (c *Client) UpdateHostConfig(ctx context.Context, cfg *arr.HostConfigResource) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/host/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Config: UI ----------.

// GetUIConfig returns the UI configuration.
func (c *Client) GetUIConfig(ctx context.Context) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/ui", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateUIConfig updates the UI configuration.
func (c *Client) UpdateUIConfig(ctx context.Context, cfg *arr.UIConfigResource) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/ui/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Config: Media Management ----------.

// GetMediaManagementConfig returns the media management configuration.
func (c *Client) GetMediaManagementConfig(ctx context.Context) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/mediamanagement", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMediaManagementConfig updates the media management configuration.
func (c *Client) UpdateMediaManagementConfig(ctx context.Context, cfg *arr.MediaManagementConfigResource) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/mediamanagement/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Queue extras ----------.

// DeleteQueueItems removes multiple items from the download queue.
func (c *Client) DeleteQueueItems(ctx context.Context, ids []int, removeFromClient, blocklist bool) error {
	path := fmt.Sprintf("/api/v3/queue/bulk?removeFromClient=%t&blocklist=%t", removeFromClient, blocklist)
	return c.base.Delete(ctx, path, &arr.QueueBulkResource{IDs: ids}, nil)
}

// GrabQueueItem forces a grab for a queued item.
func (c *Client) GrabQueueItem(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/queue/grab/%d", id), nil, nil)
}

// GrabQueueItems forces grabs for multiple queued items.
func (c *Client) GrabQueueItems(ctx context.Context, ids []int) error {
	return c.base.Post(ctx, "/api/v3/queue/grab/bulk", &arr.QueueBulkResource{IDs: ids}, nil)
}

// GetQueueDetails returns all queue items without pagination.
func (c *Client) GetQueueDetails(ctx context.Context) ([]arr.QueueRecord, error) {
	var out []arr.QueueRecord
	if err := c.base.Get(ctx, "/api/v3/queue/details", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueueStatus returns the overall queue status counts.
func (c *Client) GetQueueStatus(ctx context.Context) (*arr.QueueStatusResource, error) {
	var out arr.QueueStatusResource
	if err := c.base.Get(ctx, "/api/v3/queue/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- History extras ----------.

// GetHistorySeries returns the history for a specific series.
func (c *Client) GetHistorySeries(ctx context.Context, seriesID int) ([]HistoryRecord, error) {
	var out []HistoryRecord
	path := fmt.Sprintf("/api/v3/history/series?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistorySince returns history records since a given date.
func (c *Client) GetHistorySince(ctx context.Context, date string) ([]HistoryRecord, error) {
	var out []HistoryRecord
	path := "/api/v3/history/since?date=" + url.QueryEscape(date)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MarkHistoryFailed marks a history record as failed to trigger re-download.
func (c *Client) MarkHistoryFailed(ctx context.Context, historyID int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/history/failed/%d", historyID), nil, nil)
}

// ---------- Languages ----------.

// GetLanguages returns all available languages.
func (c *Client) GetLanguages(ctx context.Context) ([]Language, error) {
	var out []Language
	if err := c.base.Get(ctx, "/api/v3/language", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- System extras ----------.

// GetSystemRoutes returns all registered API routes.
func (c *Client) GetSystemRoutes(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v3/system/routes", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSystemRoutesDuplicate returns duplicate API routes.
func (c *Client) GetSystemRoutesDuplicate(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v3/system/routes/duplicate", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Shutdown sends a shutdown command to Sonarr.
func (c *Client) Shutdown(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/system/shutdown", nil, nil)
}

// Restart sends a restart command to Sonarr.
func (c *Client) Restart(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/system/restart", nil, nil)
}

// ---------- Tasks ----------.

// GetTasks returns all scheduled tasks.
func (c *Client) GetTasks(ctx context.Context) ([]arr.TaskResource, error) {
	var out []arr.TaskResource
	if err := c.base.Get(ctx, "/api/v3/system/task", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTask returns a single scheduled task by ID.
func (c *Client) GetTask(ctx context.Context, id int) (*arr.TaskResource, error) {
	var out arr.TaskResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/system/task/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Updates ----------.

// GetUpdates returns available application updates.
func (c *Client) GetUpdates(ctx context.Context) ([]arr.UpdateResource, error) {
	var out []arr.UpdateResource
	if err := c.base.Get(ctx, "/api/v3/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Command ----------.

// DeleteCommand cancels/deletes a pending command by ID.
func (c *Client) DeleteCommand(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/command/%d", id), nil, nil)
}

// ---------- Episode File ----------.

// UpdateEpisodeFile updates an individual episode file's metadata
// (quality, language, etc.).
func (c *Client) UpdateEpisodeFile(ctx context.Context, ef *EpisodeFile) (*EpisodeFile, error) {
	var out EpisodeFile
	path := fmt.Sprintf("/api/v3/episodefile/%d", ef.ID)
	if err := c.base.Put(ctx, path, ef, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditEpisodeFiles performs a bulk update of episode file metadata
// (quality, language, release group).
func (c *Client) EditEpisodeFiles(ctx context.Context, editor *EpisodeFileEditorResource) error {
	return c.base.Put(ctx, "/api/v3/episodefile/editor", editor, nil)
}

// ---------- Custom Format Bulk ----------.

// UpdateCustomFormatsBulk performs a bulk update of custom formats.
func (c *Client) UpdateCustomFormatsBulk(ctx context.Context, body *arr.CustomFormatBulkResource) ([]arr.CustomFormatResource, error) {
	var out []arr.CustomFormatResource
	if err := c.base.Put(ctx, "/api/v3/customformat/bulk", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteCustomFormatsBulk bulk-deletes custom formats by IDs.
func (c *Client) DeleteCustomFormatsBulk(ctx context.Context, ids []int) error {
	body := &arr.CustomFormatBulkResource{IDs: ids}
	return c.base.Delete(ctx, "/api/v3/customformat/bulk", body, nil)
}

// ---------- Download Client Bulk ----------.

// UpdateDownloadClientsBulk performs a bulk update of download clients.
func (c *Client) UpdateDownloadClientsBulk(ctx context.Context, body *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/downloadclient/bulk", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteDownloadClientsBulk bulk-deletes download clients by IDs.
func (c *Client) DeleteDownloadClientsBulk(ctx context.Context, ids []int) error {
	body := &arr.ProviderBulkResource{IDs: ids}
	return c.base.Delete(ctx, "/api/v3/downloadclient/bulk", body, nil)
}

// TestAllDownloadClients tests all configured download clients.
func (c *Client) TestAllDownloadClients(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/testall", nil, nil)
}

// ---------- Indexer Bulk ----------.

// UpdateIndexersBulk performs a bulk update of indexers.
func (c *Client) UpdateIndexersBulk(ctx context.Context, body *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/indexer/bulk", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteIndexersBulk bulk-deletes indexers by IDs.
func (c *Client) DeleteIndexersBulk(ctx context.Context, ids []int) error {
	body := &arr.ProviderBulkResource{IDs: ids}
	return c.base.Delete(ctx, "/api/v3/indexer/bulk", body, nil)
}

// TestAllIndexers tests all configured indexers.
func (c *Client) TestAllIndexers(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/indexer/testall", nil, nil)
}

// ---------- Import List Bulk ----------.

// UpdateImportListsBulk performs a bulk update of import lists.
func (c *Client) UpdateImportListsBulk(ctx context.Context, body *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/importlist/bulk", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteImportListsBulk bulk-deletes import lists by IDs.
func (c *Client) DeleteImportListsBulk(ctx context.Context, ids []int) error {
	body := &arr.ProviderBulkResource{IDs: ids}
	return c.base.Delete(ctx, "/api/v3/importlist/bulk", body, nil)
}

// TestAllImportLists tests all configured import lists.
func (c *Client) TestAllImportLists(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/importlist/testall", nil, nil)
}

// ---------- Import List Config ----------.

// GetImportListConfig returns the global import list configuration.
func (c *Client) GetImportListConfig(ctx context.Context) (*ImportListConfigResource, error) {
	var out ImportListConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/importlist", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportListConfig updates the global import list configuration.
func (c *Client) UpdateImportListConfig(ctx context.Context, cfg *ImportListConfigResource) (*ImportListConfigResource, error) {
	var out ImportListConfigResource
	path := fmt.Sprintf("/api/v3/config/importlist/%d", cfg.ID)
	if err := c.base.Put(ctx, path, cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Import List Exclusion Bulk ----------.

// GetImportListExclusionsPaged returns a paginated list of import list exclusions.
func (c *Client) GetImportListExclusionsPaged(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.ImportListExclusionResource], error) {
	var out arr.PagingResource[arr.ImportListExclusionResource]
	path := fmt.Sprintf("/api/v3/importlistexclusion/paged?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportListExclusionsBulk bulk-deletes import list exclusions by IDs.
func (c *Client) DeleteImportListExclusionsBulk(ctx context.Context, ids []int) error {
	body := struct {
		IDs []int `json:"ids"`
	}{IDs: ids}
	return c.base.Delete(ctx, "/api/v3/importlistexclusion/bulk", body, nil)
}

// ---------- Notification / Metadata TestAll ----------.

// TestAllNotifications tests all configured notifications.
func (c *Client) TestAllNotifications(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/notification/testall", nil, nil)
}

// TestAllMetadataConsumers tests all configured metadata consumers.
func (c *Client) TestAllMetadataConsumers(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/metadata/testall", nil, nil)
}

// ---------- Language ----------.

// GetLanguage returns a single language by ID.
func (c *Client) GetLanguage(ctx context.Context, id int) (*Language, error) {
	var out Language
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/language/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Localization ----------.

// GetLocalization returns the localization strings for the current locale.
func (c *Client) GetLocalization(ctx context.Context) (*LocalizationResource, error) {
	var out LocalizationResource
	if err := c.base.Get(ctx, "/api/v3/localization", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Log File ----------.

// GetLogFileContent returns the content of a specific log file by filename.
func (c *Client) GetLogFileContent(ctx context.Context, filename string) (string, error) {
	path := "/api/v3/log/file/" + url.PathEscape(filename)
	b, err := c.base.GetRaw(ctx, path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------- Quality Definition Bulk ----------.

// UpdateQualityDefinitions performs a bulk update of quality definitions.
func (c *Client) UpdateQualityDefinitions(ctx context.Context, defs []arr.QualityDefinitionResource) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Put(ctx, "/api/v3/qualitydefinition/update", defs, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Quality Profile Schema ----------.

// GetQualityProfileSchema returns the schema for quality profiles.
func (c *Client) GetQualityProfileSchema(ctx context.Context) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v3/qualityprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Root Folder ----------.

// UpdateRootFolder updates an existing root folder.
func (c *Client) UpdateRootFolder(ctx context.Context, rf *arr.RootFolder) (*arr.RootFolder, error) {
	var out arr.RootFolder
	path := fmt.Sprintf("/api/v3/rootfolder/%d", rf.ID)
	if err := c.base.Put(ctx, path, rf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- FileSystem ----------.

// BrowseFileSystem returns directories and files at the given path.
func (c *Client) BrowseFileSystem(ctx context.Context, path string, includeFiles bool) (*FileSystemResource, error) {
	var out FileSystemResource
	endpoint := fmt.Sprintf("/api/v3/filesystem?path=%s&includeFiles=%t", url.QueryEscape(path), includeFiles)
	if err := c.base.Get(ctx, endpoint, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFileSystemType returns the filesystem type (e.g. local, network) for a path.
func (c *Client) GetFileSystemType(ctx context.Context, path string) (string, error) {
	var out string
	endpoint := "/api/v3/filesystem/type?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, endpoint, &out); err != nil {
		return "", err
	}
	return out, nil
}

// GetFileSystemMediaFiles returns media files at the given path.
func (c *Client) GetFileSystemMediaFiles(ctx context.Context, path string) ([]string, error) {
	var out []string
	endpoint := "/api/v3/filesystem/mediafiles?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, endpoint, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Ping ----------.

// Ping checks connectivity to the Sonarr instance.
func (c *Client) Ping(ctx context.Context) error {
	return c.base.Get(ctx, "/ping", nil)
}

// ---------- Series Folder ----------.

// GetSeriesFolder returns folder information for a series.
func (c *Client) GetSeriesFolder(ctx context.Context, seriesID int) error {
	return c.base.Get(ctx, fmt.Sprintf("/api/v3/series/%d/folder", seriesID), nil)
}

// ---------- Calendar By ID ----------.

// GetCalendarByID returns a single calendar entry by its ID.
func (c *Client) GetCalendarByID(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/calendar/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Wanted: Cutoff By ID ----------.

// GetWantedCutoffByID returns a single wanted cutoff record by its ID.
func (c *Client) GetWantedCutoffByID(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/wanted/cutoff/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Wanted: Missing By ID ----------.

// GetWantedMissingByID returns a single wanted missing record by its ID.
func (c *Client) GetWantedMissingByID(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/wanted/missing/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Config By-ID Gets ----------.

// GetDownloadClientConfigByID returns the download client config by its ID.
func (c *Client) GetDownloadClientConfigByID(ctx context.Context, id int) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/config/downloadclient/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHostConfigByID returns the host config by its ID.
func (c *Client) GetHostConfigByID(ctx context.Context, id int) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/config/host/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetImportListConfigByID returns the import list config by its ID.
func (c *Client) GetImportListConfigByID(ctx context.Context, id int) (*ImportListConfigResource, error) {
	var out ImportListConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/config/importlist/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIndexerConfigByID returns the indexer config by its ID.
func (c *Client) GetIndexerConfigByID(ctx context.Context, id int) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/config/indexer/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMediaManagementConfigByID returns the media management config by its ID.
func (c *Client) GetMediaManagementConfigByID(ctx context.Context, id int) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/config/mediamanagement/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNamingConfigByID returns the naming config by its ID.
func (c *Client) GetNamingConfigByID(ctx context.Context, id int) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/config/naming/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUIConfigByID returns the UI config by its ID.
func (c *Client) GetUIConfigByID(ctx context.Context, id int) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/config/ui/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Provider Action endpoints ----------.

// DownloadClientAction triggers a named action on a download client provider.
func (c *Client) DownloadClientAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v3/downloadclient/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ImportListAction triggers a named action on an import list provider.
func (c *Client) ImportListAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v3/importlist/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// IndexerAction triggers a named action on an indexer provider.
func (c *Client) IndexerAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v3/indexer/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// MetadataAction triggers a named action on a metadata provider.
func (c *Client) MetadataAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v3/metadata/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// NotificationAction triggers a named action on a notification provider.
func (c *Client) NotificationAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v3/notification/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Language Profile ----------.

// GetLanguageProfiles returns all language profiles.
func (c *Client) GetLanguageProfiles(ctx context.Context) ([]LanguageProfileResource, error) {
	var out []LanguageProfileResource
	if err := c.base.Get(ctx, "/api/v3/languageprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLanguageProfile returns a single language profile by ID.
func (c *Client) GetLanguageProfile(ctx context.Context, id int) (*LanguageProfileResource, error) {
	var out LanguageProfileResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/languageprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateLanguageProfile creates a new language profile.
func (c *Client) CreateLanguageProfile(ctx context.Context, profile *LanguageProfileResource) (*LanguageProfileResource, error) {
	var out LanguageProfileResource
	if err := c.base.Post(ctx, "/api/v3/languageprofile", profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateLanguageProfile updates an existing language profile.
func (c *Client) UpdateLanguageProfile(ctx context.Context, profile *LanguageProfileResource) (*LanguageProfileResource, error) {
	var out LanguageProfileResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/languageprofile/%d", profile.ID), profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteLanguageProfile deletes a language profile by ID.
func (c *Client) DeleteLanguageProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/languageprofile/%d", id), nil, nil)
}

// GetLanguageProfileSchema returns the available language profile schema.
func (c *Client) GetLanguageProfileSchema(ctx context.Context) (*LanguageProfileResource, error) {
	var out LanguageProfileResource
	if err := c.base.Get(ctx, "/api/v3/languageprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Localization extras ----------.

// GetLocalizationByID returns localization strings by localization ID.
func (c *Client) GetLocalizationByID(ctx context.Context, id int) (*LocalizationResource, error) {
	var out LocalizationResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/localization/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLocalizationLanguages returns the list of available localization languages.
func (c *Client) GetLocalizationLanguages(ctx context.Context) ([]LocalizationLanguageResource, error) {
	var out []LocalizationLanguageResource
	if err := c.base.Get(ctx, "/api/v3/localization/language", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Naming Examples ----------.

// GetNamingExamples returns naming format examples based on the current naming config.
func (c *Client) GetNamingExamples(ctx context.Context) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/naming/examples", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Quality Definition Limits ----------.

// GetQualityDefinitionLimits returns the min/max limits for quality definitions.
func (c *Client) GetQualityDefinitionLimits(ctx context.Context) (*QualityDefinitionLimitsResource, error) {
	var out QualityDefinitionLimitsResource
	if err := c.base.Get(ctx, "/api/v3/qualitydefinition/limits", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Episode File Bulk Update ----------.

// UpdateEpisodeFilesBulk performs a bulk update of episode file properties.
func (c *Client) UpdateEpisodeFilesBulk(ctx context.Context, editor *EpisodeFileEditorResource) ([]EpisodeFile, error) {
	var out []EpisodeFile
	if err := c.base.Put(ctx, "/api/v3/episodefile/bulk", editor, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Update Log File Content ----------.

// GetUpdateLogFileContent returns the content of a specific update log file.
func (c *Client) GetUpdateLogFileContent(ctx context.Context, filename string) (string, error) {
	path := "/api/v3/log/file/update/" + url.PathEscape(filename)
	b, err := c.base.GetRaw(ctx, path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------- HEAD Ping ----------.

// HeadPing performs a lightweight HEAD request to /ping.
func (c *Client) HeadPing(ctx context.Context) error {
	return c.base.Head(ctx, "/ping")
}

// ---------- Backup Upload ----------.

// UploadBackup uploads a backup file via multipart form POST.
func (c *Client) UploadBackup(ctx context.Context, fileName string, data io.Reader) error {
	return c.base.Upload(ctx, "/api/v3/system/backup/upload", "file", fileName, data)
}
