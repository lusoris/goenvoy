package whisparr

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/golusoris/goenvoy/arr/v2"
)

// Client is a Whisparr v2 (Sonarr-based) API client.
type Client struct {
	base *arr.BaseClient
}

// New creates a Whisparr v2 [Client] for the instance at baseURL.
func New(baseURL, apiKey string, opts ...arr.Option) (*Client, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{base: base}, nil
}

// GetAllSeries returns every series configured in the instance.
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

// AddSeries adds a new series.
func (c *Client) AddSeries(ctx context.Context, series *Series) (*Series, error) {
	var out Series
	if err := c.base.Post(ctx, "/api/v3/series", series, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSeries updates an existing series. Set moveFiles to true to relocate
// files when the path changes.
func (c *Client) UpdateSeries(ctx context.Context, series *Series, moveFiles bool) (*Series, error) {
	var out Series
	path := fmt.Sprintf("/api/v3/series/%d?moveFiles=%t", series.ID, moveFiles)
	if err := c.base.Put(ctx, path, series, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteSeries removes a series by ID.
func (c *Client) DeleteSeries(ctx context.Context, id int, deleteFiles, addImportListExclusion bool) error {
	path := fmt.Sprintf("/api/v3/series/%d?deleteFiles=%t&addImportListExclusion=%t", id, deleteFiles, addImportListExclusion)
	return c.base.Delete(ctx, path, nil, nil)
}

// LookupSeries searches for a series by term.
func (c *Client) LookupSeries(ctx context.Context, term string) ([]Series, error) {
	var out []Series
	path := "/api/v3/series/lookup?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEpisodes returns all episodes for a series.
func (c *Client) GetEpisodes(ctx context.Context, seriesID int) ([]Episode, error) {
	var out []Episode
	path := fmt.Sprintf("/api/v3/episode?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEpisode returns a single episode by ID.
func (c *Client) GetEpisode(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	path := fmt.Sprintf("/api/v3/episode/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEpisodeFiles returns all episode files for a series.
func (c *Client) GetEpisodeFiles(ctx context.Context, seriesID int) ([]EpisodeFile, error) {
	var out []EpisodeFile
	path := fmt.Sprintf("/api/v3/episodefile?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteEpisodeFile deletes an episode file by ID.
func (c *Client) DeleteEpisodeFile(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v3/episodefile/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetCalendar returns episodes airing between start and end dates.
func (c *Client) GetCalendar(ctx context.Context, start, end string, unmonitored bool) ([]Episode, error) {
	var out []Episode
	path := fmt.Sprintf("/api/v3/calendar?start=%s&end=%s&unmonitored=%t",
		url.QueryEscape(start), url.QueryEscape(end), unmonitored)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendCommand sends a command to the instance.
func (c *Client) SendCommand(ctx context.Context, cmd arr.CommandRequest) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Post(ctx, "/api/v3/command", cmd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Parse parses a title string and returns matched series and episodes.
func (c *Client) Parse(ctx context.Context, title string) (*V2ParseResult, error) {
	var out V2ParseResult
	path := "/api/v3/parse?title=" + url.QueryEscape(title)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSystemStatus returns system information.
func (c *Client) GetSystemStatus(ctx context.Context) (*arr.StatusResponse, error) {
	var out arr.StatusResponse
	if err := c.base.Get(ctx, "/api/v3/system/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHealth returns a list of health check results.
func (c *Client) GetHealth(ctx context.Context) ([]arr.HealthCheck, error) {
	var out []arr.HealthCheck
	if err := c.base.Get(ctx, "/api/v3/health", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDiskSpace returns disk space information for all root folders.
func (c *Client) GetDiskSpace(ctx context.Context) ([]arr.DiskSpace, error) {
	var out []arr.DiskSpace
	if err := c.base.Get(ctx, "/api/v3/diskspace", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueue returns the download queue (paged).
func (c *Client) GetQueue(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.QueueRecord], error) {
	var out arr.PagingResource[arr.QueueRecord]
	path := fmt.Sprintf("/api/v3/queue?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQualityProfiles returns all quality profiles.
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

// CreateTag creates a new tag with the given label.
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

// GetHistory returns history records (paged).
func (c *Client) GetHistory(ctx context.Context, page, pageSize int) (*arr.PagingResource[V2HistoryRecord], error) {
	var out arr.PagingResource[V2HistoryRecord]
	path := fmt.Sprintf("/api/v3/history?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSeasonPass updates monitoring for multiple series at once.
func (c *Client) UpdateSeasonPass(ctx context.Context, pass SeasonPassResource) error {
	return c.base.Post(ctx, "/api/v3/seasonpass", pass, nil)
}

// ---------- AutoTagging ----------.

// GetAutoTags returns all auto tagging rules.
func (c *Client) GetAutoTags(ctx context.Context) ([]arr.AutoTaggingResource, error) {
	var out []arr.AutoTaggingResource
	if err := c.base.Get(ctx, "/api/v3/autotagging", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAutoTag returns a single auto tag by ID.
func (c *Client) GetAutoTag(ctx context.Context, id int) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/autotagging/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateAutoTag creates a new auto tagging rule.
func (c *Client) CreateAutoTag(ctx context.Context, tag *arr.AutoTaggingResource) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Post(ctx, "/api/v3/autotagging", tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAutoTag updates an existing auto tagging rule.
func (c *Client) UpdateAutoTag(ctx context.Context, tag *arr.AutoTaggingResource) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/autotagging/%d", tag.ID), tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAutoTag deletes an auto tagging rule by ID.
func (c *Client) DeleteAutoTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/autotagging/%d", id), nil, nil)
}

// GetAutoTagSchema returns the auto tagging schema.
func (c *Client) GetAutoTagSchema(ctx context.Context) ([]arr.AutoTaggingSpecification, error) {
	var out []arr.AutoTaggingSpecification
	if err := c.base.Get(ctx, "/api/v3/autotagging/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Backup ----------.

// GetBackups returns all available backups.
func (c *Client) GetBackups(ctx context.Context) ([]arr.Backup, error) {
	var out []arr.Backup
	if err := c.base.Get(ctx, "/api/v3/system/backup", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBackup deletes a backup by ID.
func (c *Client) DeleteBackup(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/system/backup/%d", id), nil, nil)
}

// RestoreBackup restores from a backup by ID.
func (c *Client) RestoreBackup(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/system/backup/restore/%d", id), nil, nil)
}

// ---------- Blocklist ----------.

// GetBlocklist returns the blocklist (paged).
func (c *Client) GetBlocklist(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.BlocklistResource], error) {
	var out arr.PagingResource[arr.BlocklistResource]
	path := fmt.Sprintf("/api/v3/blocklist?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBlocklistItem deletes a single blocklist entry.
func (c *Client) DeleteBlocklistItem(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/blocklist/%d", id), nil, nil)
}

// BulkDeleteBlocklist deletes multiple blocklist entries.
func (c *Client) BulkDeleteBlocklist(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/blocklist/bulk", arr.BlocklistBulkResource{IDs: ids}, nil)
}

// ---------- Calendar Extended ----------.

// GetCalendarByID returns a single calendar entry by ID.
func (c *Client) GetCalendarByID(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/calendar/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Command Extended ----------.

// GetCommands returns all current commands.
func (c *Client) GetCommands(ctx context.Context) ([]arr.CommandResponse, error) {
	var out []arr.CommandResponse
	if err := c.base.Get(ctx, "/api/v3/command", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCommand returns a command by ID.
func (c *Client) GetCommand(ctx context.Context, id int) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/command/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCommand cancels a command by ID.
func (c *Client) DeleteCommand(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/command/%d", id), nil, nil)
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
func (c *Client) CreateCustomFilter(ctx context.Context, f *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Post(ctx, "/api/v3/customfilter", f, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFilter updates an existing custom filter.
func (c *Client) UpdateCustomFilter(ctx context.Context, f *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/customfilter/%d", f.ID), f, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFilter deletes a custom filter by ID.
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

// GetCustomFormat returns a custom format by ID.
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

// DeleteCustomFormat deletes a custom format by ID.
func (c *Client) DeleteCustomFormat(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/customformat/%d", id), nil, nil)
}

// GetCustomFormatSchema returns the custom format schema.
func (c *Client) GetCustomFormatSchema(ctx context.Context) ([]arr.CustomFormatSpecification, error) {
	var out []arr.CustomFormatSpecification
	if err := c.base.Get(ctx, "/api/v3/customformat/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Wanted ----------.

// GetWantedMissing returns missing episodes (paged).
func (c *Client) GetWantedMissing(ctx context.Context, page, pageSize int) (*arr.PagingResource[Episode], error) {
	var out arr.PagingResource[Episode]
	path := fmt.Sprintf("/api/v3/wanted/missing?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedMissingByID returns a single missing episode by ID.
func (c *Client) GetWantedMissingByID(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/wanted/missing/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedCutoff returns cutoff unmet episodes (paged).
func (c *Client) GetWantedCutoff(ctx context.Context, page, pageSize int) (*arr.PagingResource[Episode], error) {
	var out arr.PagingResource[Episode]
	path := fmt.Sprintf("/api/v3/wanted/cutoff?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedCutoffByID returns a single cutoff-unmet episode by ID.
func (c *Client) GetWantedCutoffByID(ctx context.Context, id int) (*Episode, error) {
	var out Episode
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/wanted/cutoff/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
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

// DeleteDelayProfile deletes a delay profile by ID.
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

// ---------- Download Clients ----------.

// GetDownloadClients returns all configured download clients.
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

// DeleteDownloadClient deletes a download client by ID.
func (c *Client) DeleteDownloadClient(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/downloadclient/%d", id), nil, nil)
}

// GetDownloadClientSchema returns the schema for all download client types.
func (c *Client) GetDownloadClientSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/downloadclient/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestDownloadClient tests a download client configuration.
func (c *Client) TestDownloadClient(ctx context.Context, dc *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/test", dc, nil)
}

// TestAllDownloadClients tests all configured download clients.
func (c *Client) TestAllDownloadClients(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/testall", nil, nil)
}

// BulkUpdateDownloadClients updates multiple download clients.
func (c *Client) BulkUpdateDownloadClients(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/downloadclient/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteDownloadClients deletes multiple download clients.
func (c *Client) BulkDeleteDownloadClients(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/downloadclient/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// DownloadClientAction triggers a named action on a download client.
func (c *Client) DownloadClientAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/action/"+url.PathEscape(name), body, nil)
}

// ---------- Download Client Config ----------.

// GetDownloadClientConfig returns the download client configuration.
func (c *Client) GetDownloadClientConfig(ctx context.Context) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/downloadclient", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClientConfig updates the download client configuration.
func (c *Client) UpdateDownloadClientConfig(ctx context.Context, cfg *arr.DownloadClientConfigResource) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/downloadclient/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Episode Extended ----------.

// UpdateEpisode updates an episode (e.g. monitored status).
func (c *Client) UpdateEpisode(ctx context.Context, ep *Episode) (*Episode, error) {
	var out Episode
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/episode/%d", ep.ID), ep, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MonitorEpisodes sets the monitored status for multiple episodes.
func (c *Client) MonitorEpisodes(ctx context.Context, req *EpisodesMonitoredResource) ([]Episode, error) {
	var out []Episode
	if err := c.base.Put(ctx, "/api/v3/episode/monitor", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Episode Files Extended ----------.

// GetEpisodeFile returns a single episode file by ID.
func (c *Client) GetEpisodeFile(ctx context.Context, id int) (*EpisodeFile, error) {
	var out EpisodeFile
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/episodefile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateEpisodeFile updates an episode file (quality, language).
func (c *Client) UpdateEpisodeFile(ctx context.Context, ef *EpisodeFile) (*EpisodeFile, error) {
	var out EpisodeFile
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/episodefile/%d", ef.ID), ef, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditEpisodeFiles bulk edits episode files.
func (c *Client) EditEpisodeFiles(ctx context.Context, editor *EpisodeFileEditorResource) error {
	return c.base.Put(ctx, "/api/v3/episodefile/editor", editor, nil)
}

// BulkDeleteEpisodeFiles deletes multiple episode files.
func (c *Client) BulkDeleteEpisodeFiles(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/episodefile/bulk", struct {
		EpisodeFileIDs []int `json:"episodeFileIds"`
	}{EpisodeFileIDs: ids}, nil)
}

// BulkUpdateEpisodeFiles updates multiple episode files.
func (c *Client) BulkUpdateEpisodeFiles(ctx context.Context, editor *EpisodeFileEditorResource) ([]EpisodeFile, error) {
	var out []EpisodeFile
	if err := c.base.Put(ctx, "/api/v3/episodefile/bulk", editor, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- File System ----------.

// BrowseFileSystem returns a directory listing.
func (c *Client) BrowseFileSystem(ctx context.Context, path string) (map[string]any, error) {
	var out map[string]any
	reqPath := "/api/v3/filesystem?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFileSystemType returns the filesystem type for a path.
func (c *Client) GetFileSystemType(ctx context.Context, path string) (string, error) {
	var out string
	reqPath := "/api/v3/filesystem/type?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return "", err
	}
	return out, nil
}

// GetFileSystemMediaFiles returns media files in a path.
func (c *Client) GetFileSystemMediaFiles(ctx context.Context, path string) ([]map[string]any, error) {
	var out []map[string]any
	reqPath := "/api/v3/filesystem/mediafiles?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Host Config ----------.

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

// ---------- History Extended ----------.

// GetHistorySince returns history since a specific date.
func (c *Client) GetHistorySince(ctx context.Context, date string) ([]V2HistoryRecord, error) {
	var out []V2HistoryRecord
	path := "/api/v3/history/since?date=" + url.QueryEscape(date)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistorySeries returns history for a specific series.
func (c *Client) GetHistorySeries(ctx context.Context, seriesID int) ([]V2HistoryRecord, error) {
	var out []V2HistoryRecord
	path := fmt.Sprintf("/api/v3/history/series?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MarkHistoryFailed marks a history item as failed.
func (c *Client) MarkHistoryFailed(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/history/failed/%d", id), nil, nil)
}

// ---------- Import Lists ----------.

// GetImportLists returns all configured import lists.
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

// DeleteImportList deletes an import list by ID.
func (c *Client) DeleteImportList(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/importlist/%d", id), nil, nil)
}

// GetImportListSchema returns the schema for all import list types.
func (c *Client) GetImportListSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/importlist/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestImportList tests an import list configuration.
func (c *Client) TestImportList(ctx context.Context, il *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/importlist/test", il, nil)
}

// TestAllImportLists tests all configured import lists.
func (c *Client) TestAllImportLists(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/importlist/testall", nil, nil)
}

// BulkUpdateImportLists updates multiple import lists.
func (c *Client) BulkUpdateImportLists(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/importlist/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteImportLists deletes multiple import lists.
func (c *Client) BulkDeleteImportLists(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/importlist/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// ImportListAction triggers a named action on an import list.
func (c *Client) ImportListAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/importlist/action/"+url.PathEscape(name), body, nil)
}

// ---------- Import List Config ----------.

// GetImportListConfig returns the import list configuration.
func (c *Client) GetImportListConfig(ctx context.Context) (*ImportListConfigResource, error) {
	var out ImportListConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/importlist", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportListConfig updates the import list configuration.
func (c *Client) UpdateImportListConfig(ctx context.Context, cfg *ImportListConfigResource) (*ImportListConfigResource, error) {
	var out ImportListConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/importlist/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
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
func (c *Client) CreateImportListExclusion(ctx context.Context, ex *arr.ImportListExclusionResource) (*arr.ImportListExclusionResource, error) {
	var out arr.ImportListExclusionResource
	if err := c.base.Post(ctx, "/api/v3/importlistexclusion", ex, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportListExclusion updates an existing import list exclusion.
func (c *Client) UpdateImportListExclusion(ctx context.Context, ex *arr.ImportListExclusionResource) (*arr.ImportListExclusionResource, error) {
	var out arr.ImportListExclusionResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/importlistexclusion/%d", ex.ID), ex, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportListExclusion deletes an import list exclusion by ID.
func (c *Client) DeleteImportListExclusion(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/importlistexclusion/%d", id), nil, nil)
}

// ---------- Indexers ----------.

// GetIndexers returns all configured indexers.
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

// DeleteIndexer deletes an indexer by ID.
func (c *Client) DeleteIndexer(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/indexer/%d", id), nil, nil)
}

// GetIndexerSchema returns the schema for all indexer types.
func (c *Client) GetIndexerSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/indexer/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestIndexer tests an indexer configuration.
func (c *Client) TestIndexer(ctx context.Context, idx *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/indexer/test", idx, nil)
}

// TestAllIndexers tests all configured indexers.
func (c *Client) TestAllIndexers(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/indexer/testall", nil, nil)
}

// BulkUpdateIndexers updates multiple indexers.
func (c *Client) BulkUpdateIndexers(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/indexer/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteIndexers deletes multiple indexers.
func (c *Client) BulkDeleteIndexers(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/indexer/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// IndexerAction triggers a named action on an indexer.
func (c *Client) IndexerAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/indexer/action/"+url.PathEscape(name), body, nil)
}

// ---------- Indexer Config ----------.

// GetIndexerConfig returns the indexer configuration.
func (c *Client) GetIndexerConfig(ctx context.Context) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/indexer", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexerConfig updates the indexer configuration.
func (c *Client) UpdateIndexerConfig(ctx context.Context, cfg *arr.IndexerConfigResource) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/indexer/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Languages ----------.

// GetLanguages returns all available languages.
func (c *Client) GetLanguages(ctx context.Context) ([]arr.LanguageResource, error) {
	var out []arr.LanguageResource
	if err := c.base.Get(ctx, "/api/v3/language", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLanguage returns a language by ID.
func (c *Client) GetLanguage(ctx context.Context, id int) (*arr.LanguageResource, error) {
	var out arr.LanguageResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/language/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Language Profiles ----------.

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
func (c *Client) CreateLanguageProfile(ctx context.Context, lp *LanguageProfileResource) (*LanguageProfileResource, error) {
	var out LanguageProfileResource
	if err := c.base.Post(ctx, "/api/v3/languageprofile", lp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateLanguageProfile updates an existing language profile.
func (c *Client) UpdateLanguageProfile(ctx context.Context, lp *LanguageProfileResource) (*LanguageProfileResource, error) {
	var out LanguageProfileResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/languageprofile/%d", lp.ID), lp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteLanguageProfile deletes a language profile by ID.
func (c *Client) DeleteLanguageProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/languageprofile/%d", id), nil, nil)
}

// GetLanguageProfileSchema returns the language profile schema.
func (c *Client) GetLanguageProfileSchema(ctx context.Context) (*LanguageProfileResource, error) {
	var out LanguageProfileResource
	if err := c.base.Get(ctx, "/api/v3/languageprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Localization ----------.

// GetLocalization returns localization strings.
func (c *Client) GetLocalization(ctx context.Context) (map[string]string, error) {
	var out map[string]string
	if err := c.base.Get(ctx, "/api/v3/localization", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLocalizationLanguage returns the current localization language.
func (c *Client) GetLocalizationLanguage(ctx context.Context) (*LocalizationLanguageResource, error) {
	var out LocalizationLanguageResource
	if err := c.base.Get(ctx, "/api/v3/localization/language", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Logs ----------.

// GetLogs returns log entries (paged).
func (c *Client) GetLogs(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.LogRecord], error) {
	var out arr.PagingResource[arr.LogRecord]
	path := fmt.Sprintf("/api/v3/log?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLogFiles returns available log files.
func (c *Client) GetLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v3/log/file", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLogFileContent returns the content of a specific log file.
func (c *Client) GetLogFileContent(ctx context.Context, filename string) (string, error) {
	b, err := c.base.GetRaw(ctx, "/api/v3/log/file/"+url.PathEscape(filename))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GetUpdateLogFiles returns available update log files.
func (c *Client) GetUpdateLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v3/log/file/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetUpdateLogFileContent returns the content of a specific update log file.
func (c *Client) GetUpdateLogFileContent(ctx context.Context, filename string) (string, error) {
	b, err := c.base.GetRaw(ctx, "/api/v3/log/file/update/"+url.PathEscape(filename))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------- Manual Import ----------.

// GetManualImport returns items available for manual import.
func (c *Client) GetManualImport(ctx context.Context, folder string) ([]arr.ManualImportResource, error) {
	var out []arr.ManualImportResource
	path := "/api/v3/manualimport?folder=" + url.QueryEscape(folder)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessManualImport triggers processing of manual import selections.
func (c *Client) ProcessManualImport(ctx context.Context, items []arr.ManualImportReprocessResource) error {
	return c.base.Post(ctx, "/api/v3/manualimport", items, nil)
}

// ---------- Media Management Config ----------.

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

// ---------- Metadata Consumers ----------.

// GetMetadata returns all metadata consumers.
func (c *Client) GetMetadata(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/metadata", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMetadataByID returns a metadata consumer by ID.
func (c *Client) GetMetadataByID(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/metadata/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateMetadata creates a new metadata consumer.
func (c *Client) CreateMetadata(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/metadata", m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMetadata updates an existing metadata consumer.
func (c *Client) UpdateMetadata(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/metadata/%d", m.ID), m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMetadata deletes a metadata consumer by ID.
func (c *Client) DeleteMetadata(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/metadata/%d", id), nil, nil)
}

// GetMetadataSchema returns the schema for all metadata consumer types.
func (c *Client) GetMetadataSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/metadata/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestMetadata tests a metadata consumer configuration.
func (c *Client) TestMetadata(ctx context.Context, m *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/metadata/test", m, nil)
}

// TestAllMetadata tests all configured metadata consumers.
func (c *Client) TestAllMetadata(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/metadata/testall", nil, nil)
}

// MetadataAction triggers a named action on a metadata consumer.
func (c *Client) MetadataAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/metadata/action/"+url.PathEscape(name), body, nil)
}

// ---------- Naming Config ----------.

// GetNamingConfig returns the naming configuration.
func (c *Client) GetNamingConfig(ctx context.Context) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/naming", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNamingConfig updates the naming configuration.
func (c *Client) UpdateNamingConfig(ctx context.Context, cfg *arr.NamingConfigResource) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/naming/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNamingExamples returns naming examples based on the current config.
func (c *Client) GetNamingExamples(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := c.base.Get(ctx, "/api/v3/config/naming/examples", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Notifications ----------.

// GetNotifications returns all configured notifications.
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

// DeleteNotification deletes a notification by ID.
func (c *Client) DeleteNotification(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/notification/%d", id), nil, nil)
}

// GetNotificationSchema returns the schema for all notification types.
func (c *Client) GetNotificationSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/notification/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestNotification tests a notification configuration.
func (c *Client) TestNotification(ctx context.Context, n *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/notification/test", n, nil)
}

// TestAllNotifications tests all configured notifications.
func (c *Client) TestAllNotifications(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/notification/testall", nil, nil)
}

// NotificationAction triggers a named action on a notification.
func (c *Client) NotificationAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/notification/action/"+url.PathEscape(name), body, nil)
}

// ---------- Quality Definitions ----------.

// GetQualityDefinitions returns all quality definitions.
func (c *Client) GetQualityDefinitions(ctx context.Context) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Get(ctx, "/api/v3/qualitydefinition", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQualityDefinition returns a quality definition by ID.
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

// BulkUpdateQualityDefinitions updates multiple quality definitions.
func (c *Client) BulkUpdateQualityDefinitions(ctx context.Context, defs []arr.QualityDefinitionResource) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Put(ctx, "/api/v3/qualitydefinition/update", defs, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQualityDefinitionLimits returns the quality definition limits.
func (c *Client) GetQualityDefinitionLimits(ctx context.Context) (*QualityDefinitionLimitsResource, error) {
	var out QualityDefinitionLimitsResource
	if err := c.base.Get(ctx, "/api/v3/qualitydefinition/limits", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Quality Profile Extended ----------.

// GetQualityProfile returns a quality profile by ID.
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

// DeleteQualityProfile deletes a quality profile by ID.
func (c *Client) DeleteQualityProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/qualityprofile/%d", id), nil, nil)
}

// GetQualityProfileSchema returns the quality profile schema.
func (c *Client) GetQualityProfileSchema(ctx context.Context) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v3/qualityprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Queue Extended ----------.

// DeleteQueueItem removes an item from the download queue.
func (c *Client) DeleteQueueItem(ctx context.Context, id int, removeFromClient, blocklist bool) error {
	path := fmt.Sprintf("/api/v3/queue/%d?removeFromClient=%t&blocklist=%t", id, removeFromClient, blocklist)
	return c.base.Delete(ctx, path, nil, nil)
}

// BulkDeleteQueue removes multiple items from the download queue.
func (c *Client) BulkDeleteQueue(ctx context.Context, bulk *arr.QueueBulkResource) error {
	return c.base.Delete(ctx, "/api/v3/queue/bulk", bulk, nil)
}

// GrabQueueItem sends a queue item to the download client.
func (c *Client) GrabQueueItem(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/queue/grab/%d", id), nil, nil)
}

// BulkGrabQueue grabs multiple items from the queue.
func (c *Client) BulkGrabQueue(ctx context.Context, ids []int) error {
	return c.base.Post(ctx, "/api/v3/queue/grab/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// GetQueueDetails returns detailed queue information.
func (c *Client) GetQueueDetails(ctx context.Context) ([]arr.QueueRecord, error) {
	var out []arr.QueueRecord
	if err := c.base.Get(ctx, "/api/v3/queue/details", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueueStatus returns the queue status summary.
func (c *Client) GetQueueStatus(ctx context.Context) (*arr.QueueStatusResource, error) {
	var out arr.QueueStatusResource
	if err := c.base.Get(ctx, "/api/v3/queue/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Releases ----------.

// SearchReleases searches for available releases.
func (c *Client) SearchReleases(ctx context.Context, episodeID int) ([]arr.ReleaseResource, error) {
	var out []arr.ReleaseResource
	path := fmt.Sprintf("/api/v3/release?episodeId=%d", episodeID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GrabRelease sends a release to the download client.
func (c *Client) GrabRelease(ctx context.Context, release *arr.ReleaseResource) (*arr.ReleaseResource, error) {
	var out arr.ReleaseResource
	if err := c.base.Post(ctx, "/api/v3/release", release, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PushRelease pushes a release for processing.
func (c *Client) PushRelease(ctx context.Context, push *arr.ReleasePushResource) ([]arr.ReleaseResource, error) {
	var out []arr.ReleaseResource
	if err := c.base.Post(ctx, "/api/v3/release/push", push, &out); err != nil {
		return nil, err
	}
	return out, nil
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

// GetReleaseProfile returns a release profile by ID.
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

// DeleteReleaseProfile deletes a release profile by ID.
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

// GetRemotePathMapping returns a remote path mapping by ID.
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

// DeleteRemotePathMapping deletes a remote path mapping by ID.
func (c *Client) DeleteRemotePathMapping(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/remotepathmapping/%d", id), nil, nil)
}

// ---------- Rename ----------.

// GetRenamePreview returns a rename preview for a series.
func (c *Client) GetRenamePreview(ctx context.Context, seriesID int) ([]arr.RenameEpisodeResource, error) {
	var out []arr.RenameEpisodeResource
	path := fmt.Sprintf("/api/v3/rename?seriesId=%d", seriesID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Root Folder Extended ----------.

// GetRootFolder returns a root folder by ID.
func (c *Client) GetRootFolder(ctx context.Context, id int) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/rootfolder/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRootFolder creates a new root folder.
func (c *Client) CreateRootFolder(ctx context.Context, rf *arr.RootFolder) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Post(ctx, "/api/v3/rootfolder", rf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRootFolder deletes a root folder by ID.
func (c *Client) DeleteRootFolder(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/rootfolder/%d", id), nil, nil)
}

// ---------- Series Editor ----------.

// EditSeries applies bulk edits to multiple series.
func (c *Client) EditSeries(ctx context.Context, editor *SeriesEditorResource) error {
	return c.base.Put(ctx, "/api/v3/series/editor", editor, nil)
}

// DeleteSeriesBulk deletes multiple series according to the editor payload.
func (c *Client) DeleteSeriesBulk(ctx context.Context, editor *SeriesEditorResource) error {
	return c.base.Delete(ctx, "/api/v3/series/editor", editor, nil)
}

// ---------- Series Import ----------.

// ImportSeries imports one or more series.
func (c *Client) ImportSeries(ctx context.Context, series []Series) error {
	return c.base.Post(ctx, "/api/v3/series/import", series, nil)
}

// ---------- System Extended ----------.

// GetSystemRoutes returns all API routes.
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

// Shutdown sends a shutdown command.
func (c *Client) Shutdown(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/system/shutdown", nil, nil)
}

// Restart sends a restart command.
func (c *Client) Restart(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/system/restart", nil, nil)
}

// ---------- Tags Extended ----------.

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

// DeleteTag deletes a tag by ID.
func (c *Client) DeleteTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/tag/%d", id), nil, nil)
}

// GetTagDetails returns all tags with usage details.
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

// ---------- Tasks ----------.

// GetTasks returns all scheduled tasks.
func (c *Client) GetTasks(ctx context.Context) ([]arr.TaskResource, error) {
	var out []arr.TaskResource
	if err := c.base.Get(ctx, "/api/v3/system/task", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTask returns a scheduled task by ID.
func (c *Client) GetTask(ctx context.Context, id int) (*arr.TaskResource, error) {
	var out arr.TaskResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/system/task/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- UI Config ----------.

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

// ---------- Updates ----------.

// GetUpdates returns available application updates.
func (c *Client) GetUpdates(ctx context.Context) ([]arr.UpdateResource, error) {
	var out []arr.UpdateResource
	if err := c.base.Get(ctx, "/api/v3/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Ping ----------.

// Ping checks if the Whisparr instance is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.base.Get(ctx, "/ping", nil)
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
