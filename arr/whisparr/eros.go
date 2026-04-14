package whisparr

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/golusoris/goenvoy/arr/v2"
)

// ClientV3 is a Whisparr v3 (Radarr-based) API client.
type ClientV3 struct {
	base *arr.BaseClient
}

// NewV3 creates a Whisparr v3 [ClientV3] for the instance at baseURL.
func NewV3(baseURL, apiKey string, opts ...arr.Option) (*ClientV3, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientV3{base: base}, nil
}

// GetAllMovies returns every movie/scene configured in the instance.
func (c *ClientV3) GetAllMovies(ctx context.Context) ([]Movie, error) {
	var out []Movie
	if err := c.base.Get(ctx, "/api/v3/movie", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovie returns a single movie/scene by its database ID.
func (c *ClientV3) GetMovie(ctx context.Context, id int) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/movie/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddMovie adds a new movie/scene.
func (c *ClientV3) AddMovie(ctx context.Context, movie *Movie) (*Movie, error) {
	var out Movie
	if err := c.base.Post(ctx, "/api/v3/movie", movie, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMovie updates an existing movie/scene. Set moveFiles to true to
// relocate files when the path changes.
func (c *ClientV3) UpdateMovie(ctx context.Context, movie *Movie, moveFiles bool) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/movie/%d?moveFiles=%t", movie.ID, moveFiles)
	if err := c.base.Put(ctx, path, movie, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMovie removes a movie/scene by ID.
func (c *ClientV3) DeleteMovie(ctx context.Context, id int, deleteFiles, addImportExclusion bool) error {
	path := fmt.Sprintf("/api/v3/movie/%d?deleteFiles=%t&addImportExclusion=%t", id, deleteFiles, addImportExclusion)
	return c.base.Delete(ctx, path, nil, nil)
}

// LookupMovie searches for a movie by term.
func (c *ClientV3) LookupMovie(ctx context.Context, term string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/lookup/movie?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LookupScene searches for a scene by term.
func (c *ClientV3) LookupScene(ctx context.Context, term string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/lookup/scene?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMoviesByPerformer returns all movies associated with a performer foreign ID.
func (c *ClientV3) GetMoviesByPerformer(ctx context.Context, foreignID string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/movie/listbyperformerforeignid?performerForeignId=" + url.QueryEscape(foreignID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMoviesByStudio returns all movies associated with a studio foreign ID.
func (c *ClientV3) GetMoviesByStudio(ctx context.Context, foreignID string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/movie/listbystudioforeignid?studioForeignId=" + url.QueryEscape(foreignID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovieFile returns a single movie file by ID.
func (c *ClientV3) GetMovieFile(ctx context.Context, id int) (*MovieFile, error) {
	var out MovieFile
	path := fmt.Sprintf("/api/v3/moviefile/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMovieFile deletes a movie file by ID.
func (c *ClientV3) DeleteMovieFile(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v3/moviefile/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// EditMovies applies bulk edits to multiple movies.
func (c *ClientV3) EditMovies(ctx context.Context, editor *MovieEditorResource) error {
	return c.base.Put(ctx, "/api/v3/movie/editor", editor, nil)
}

// DeleteMovies deletes multiple movies according to the editor payload.
func (c *ClientV3) DeleteMovies(ctx context.Context, editor *MovieEditorResource) error {
	return c.base.Delete(ctx, "/api/v3/movie/editor", editor, nil)
}

// GetPerformers returns all performers.
func (c *ClientV3) GetPerformers(ctx context.Context) ([]Performer, error) {
	var out []Performer
	if err := c.base.Get(ctx, "/api/v3/performer", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPerformer returns a single performer by ID.
func (c *ClientV3) GetPerformer(ctx context.Context, id int) (*Performer, error) {
	var out Performer
	path := fmt.Sprintf("/api/v3/performer/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddPerformer adds a new performer to the instance.
func (c *ClientV3) AddPerformer(ctx context.Context, performer *Performer) (*Performer, error) {
	var out Performer
	if err := c.base.Post(ctx, "/api/v3/performer", performer, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdatePerformer updates an existing performer.
func (c *ClientV3) UpdatePerformer(ctx context.Context, performer *Performer) (*Performer, error) {
	var out Performer
	path := fmt.Sprintf("/api/v3/performer/%d", performer.ID)
	if err := c.base.Put(ctx, path, performer, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeletePerformer removes a performer by ID.
func (c *ClientV3) DeletePerformer(ctx context.Context, id int, deleteFiles bool) error {
	path := fmt.Sprintf("/api/v3/performer/%d?deleteFiles=%t", id, deleteFiles)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetStudios returns all studios.
func (c *ClientV3) GetStudios(ctx context.Context) ([]Studio, error) {
	var out []Studio
	if err := c.base.Get(ctx, "/api/v3/studio", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetStudio returns a single studio by ID.
func (c *ClientV3) GetStudio(ctx context.Context, id int) (*Studio, error) {
	var out Studio
	path := fmt.Sprintf("/api/v3/studio/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddStudio adds a new studio to the instance.
func (c *ClientV3) AddStudio(ctx context.Context, studio *Studio) (*Studio, error) {
	var out Studio
	if err := c.base.Post(ctx, "/api/v3/studio", studio, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateStudio updates an existing studio.
func (c *ClientV3) UpdateStudio(ctx context.Context, studio *Studio) (*Studio, error) {
	var out Studio
	path := fmt.Sprintf("/api/v3/studio/%d", studio.ID)
	if err := c.base.Put(ctx, path, studio, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteStudio removes a studio by ID.
func (c *ClientV3) DeleteStudio(ctx context.Context, id int, deleteFiles bool) error {
	path := fmt.Sprintf("/api/v3/studio/%d?deleteFiles=%t", id, deleteFiles)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetCredits returns all credits for a movie/scene.
func (c *ClientV3) GetCredits(ctx context.Context, movieID int) ([]Credit, error) {
	var out []Credit
	path := fmt.Sprintf("/api/v3/credit?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCalendar returns movies/scenes releasing between start and end dates.
func (c *ClientV3) GetCalendar(ctx context.Context, start, end string, unmonitored bool) ([]Movie, error) {
	var out []Movie
	path := fmt.Sprintf("/api/v3/calendar?start=%s&end=%s&unmonitored=%t",
		url.QueryEscape(start), url.QueryEscape(end), unmonitored)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendCommand sends a command to the instance.
func (c *ClientV3) SendCommand(ctx context.Context, cmd arr.CommandRequest) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Post(ctx, "/api/v3/command", cmd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Parse parses a title string and returns matched movie info.
func (c *ClientV3) Parse(ctx context.Context, title string) (*ParseResultV3, error) {
	var out ParseResultV3
	path := "/api/v3/parse?title=" + url.QueryEscape(title)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSystemStatus returns system information.
func (c *ClientV3) GetSystemStatus(ctx context.Context) (*arr.StatusResponse, error) {
	var out arr.StatusResponse
	if err := c.base.Get(ctx, "/api/v3/system/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHealth returns a list of health check results.
func (c *ClientV3) GetHealth(ctx context.Context) ([]arr.HealthCheck, error) {
	var out []arr.HealthCheck
	if err := c.base.Get(ctx, "/api/v3/health", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDiskSpace returns disk space information for all root folders.
func (c *ClientV3) GetDiskSpace(ctx context.Context) ([]arr.DiskSpace, error) {
	var out []arr.DiskSpace
	if err := c.base.Get(ctx, "/api/v3/diskspace", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueue returns the download queue (paged).
func (c *ClientV3) GetQueue(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.QueueRecord], error) {
	var out arr.PagingResource[arr.QueueRecord]
	path := fmt.Sprintf("/api/v3/queue?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQualityProfiles returns all quality profiles.
func (c *ClientV3) GetQualityProfiles(ctx context.Context) ([]arr.QualityProfile, error) {
	var out []arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v3/qualityprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTags returns all tags.
func (c *ClientV3) GetTags(ctx context.Context) ([]arr.Tag, error) {
	var out []arr.Tag
	if err := c.base.Get(ctx, "/api/v3/tag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTag creates a new tag with the given label.
func (c *ClientV3) CreateTag(ctx context.Context, label string) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Post(ctx, "/api/v3/tag", arr.Tag{Label: label}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRootFolders returns all configured root folders.
func (c *ClientV3) GetRootFolders(ctx context.Context) ([]arr.RootFolder, error) {
	var out []arr.RootFolder
	if err := c.base.Get(ctx, "/api/v3/rootfolder", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistory returns history records (paged).
func (c *ClientV3) GetHistory(ctx context.Context, page, pageSize int) (*arr.PagingResource[HistoryRecordV3], error) {
	var out arr.PagingResource[HistoryRecordV3]
	path := fmt.Sprintf("/api/v3/history?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetImportExclusions returns all import exclusions.
func (c *ClientV3) GetImportExclusions(ctx context.Context) ([]ImportExclusion, error) {
	var out []ImportExclusion
	if err := c.base.Get(ctx, "/api/v3/exclusions", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Alternative Titles ----------.

// GetAlternativeTitles returns alternative titles for a movie.
func (c *ClientV3) GetAlternativeTitles(ctx context.Context, movieID int) ([]AlternativeTitleResource, error) {
	var out []AlternativeTitleResource
	path := fmt.Sprintf("/api/v3/alttitle?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAlternativeTitle returns a single alternative title by ID.
func (c *ClientV3) GetAlternativeTitle(ctx context.Context, id int) (*AlternativeTitleResource, error) {
	var out AlternativeTitleResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/alttitle/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- AutoTagging ----------.

// GetAutoTags returns all auto tagging rules.
func (c *ClientV3) GetAutoTags(ctx context.Context) ([]arr.AutoTaggingResource, error) {
	var out []arr.AutoTaggingResource
	if err := c.base.Get(ctx, "/api/v3/autotagging", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAutoTag returns a single auto tag by ID.
func (c *ClientV3) GetAutoTag(ctx context.Context, id int) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/autotagging/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateAutoTag creates a new auto tagging rule.
func (c *ClientV3) CreateAutoTag(ctx context.Context, tag *arr.AutoTaggingResource) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Post(ctx, "/api/v3/autotagging", tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAutoTag updates an existing auto tagging rule.
func (c *ClientV3) UpdateAutoTag(ctx context.Context, tag *arr.AutoTaggingResource) (*arr.AutoTaggingResource, error) {
	var out arr.AutoTaggingResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/autotagging/%d", tag.ID), tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAutoTag deletes an auto tagging rule by ID.
func (c *ClientV3) DeleteAutoTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/autotagging/%d", id), nil, nil)
}

// GetAutoTagSchema returns the auto tagging schema.
func (c *ClientV3) GetAutoTagSchema(ctx context.Context) ([]arr.AutoTaggingSpecification, error) {
	var out []arr.AutoTaggingSpecification
	if err := c.base.Get(ctx, "/api/v3/autotagging/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Backup ----------.

// GetBackups returns all available backups.
func (c *ClientV3) GetBackups(ctx context.Context) ([]arr.Backup, error) {
	var out []arr.Backup
	if err := c.base.Get(ctx, "/api/v3/system/backup", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBackup deletes a backup by ID.
func (c *ClientV3) DeleteBackup(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/system/backup/%d", id), nil, nil)
}

// RestoreBackup restores from a backup by ID.
func (c *ClientV3) RestoreBackup(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/system/backup/restore/%d", id), nil, nil)
}

// ---------- Blocklist ----------.

// GetBlocklist returns the blocklist (paged).
func (c *ClientV3) GetBlocklist(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.BlocklistResource], error) {
	var out arr.PagingResource[arr.BlocklistResource]
	path := fmt.Sprintf("/api/v3/blocklist?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBlocklistItem deletes a single blocklist entry.
func (c *ClientV3) DeleteBlocklistItem(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/blocklist/%d", id), nil, nil)
}

// BulkDeleteBlocklist deletes multiple blocklist entries.
func (c *ClientV3) BulkDeleteBlocklist(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/blocklist/bulk", arr.BlocklistBulkResource{IDs: ids}, nil)
}

// GetBlocklistMovie returns blocklist entries for a specific movie.
func (c *ClientV3) GetBlocklistMovie(ctx context.Context, movieID int) ([]arr.BlocklistResource, error) {
	var out []arr.BlocklistResource
	path := fmt.Sprintf("/api/v3/blocklist/movie?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Calendar Extended ----------.

// GetCalendarByID returns a single calendar entry by ID.
func (c *ClientV3) GetCalendarByID(ctx context.Context, id int) (*Movie, error) {
	var out Movie
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/calendar/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Command Extended ----------.

// GetCommands returns all current commands.
func (c *ClientV3) GetCommands(ctx context.Context) ([]arr.CommandResponse, error) {
	var out []arr.CommandResponse
	if err := c.base.Get(ctx, "/api/v3/command", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCommand returns a command by ID.
func (c *ClientV3) GetCommand(ctx context.Context, id int) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/command/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCommand cancels a command by ID.
func (c *ClientV3) DeleteCommand(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/command/%d", id), nil, nil)
}

// ---------- Credit Extended ----------.

// GetCredit returns a single credit by ID.
func (c *ClientV3) GetCredit(ctx context.Context, id int) (*Credit, error) {
	var out Credit
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/credit/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Custom Filters ----------.

// GetCustomFilters returns all custom filters.
func (c *ClientV3) GetCustomFilters(ctx context.Context) ([]arr.CustomFilterResource, error) {
	var out []arr.CustomFilterResource
	if err := c.base.Get(ctx, "/api/v3/customfilter", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCustomFilter returns a single custom filter by ID.
func (c *ClientV3) GetCustomFilter(ctx context.Context, id int) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/customfilter/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCustomFilter creates a new custom filter.
func (c *ClientV3) CreateCustomFilter(ctx context.Context, f *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Post(ctx, "/api/v3/customfilter", f, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFilter updates an existing custom filter.
func (c *ClientV3) UpdateCustomFilter(ctx context.Context, f *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/customfilter/%d", f.ID), f, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFilter deletes a custom filter by ID.
func (c *ClientV3) DeleteCustomFilter(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/customfilter/%d", id), nil, nil)
}

// ---------- Custom Formats ----------.

// GetCustomFormats returns all custom formats.
func (c *ClientV3) GetCustomFormats(ctx context.Context) ([]arr.CustomFormatResource, error) {
	var out []arr.CustomFormatResource
	if err := c.base.Get(ctx, "/api/v3/customformat", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCustomFormat returns a custom format by ID.
func (c *ClientV3) GetCustomFormat(ctx context.Context, id int) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/customformat/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCustomFormat creates a new custom format.
func (c *ClientV3) CreateCustomFormat(ctx context.Context, cf *arr.CustomFormatResource) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Post(ctx, "/api/v3/customformat", cf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFormat updates an existing custom format.
func (c *ClientV3) UpdateCustomFormat(ctx context.Context, cf *arr.CustomFormatResource) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/customformat/%d", cf.ID), cf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFormat deletes a custom format by ID.
func (c *ClientV3) DeleteCustomFormat(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/customformat/%d", id), nil, nil)
}

// GetCustomFormatSchema returns the custom format schema.
func (c *ClientV3) GetCustomFormatSchema(ctx context.Context) ([]arr.CustomFormatSpecification, error) {
	var out []arr.CustomFormatSpecification
	if err := c.base.Get(ctx, "/api/v3/customformat/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Delay Profiles ----------.

// GetDelayProfiles returns all delay profiles.
func (c *ClientV3) GetDelayProfiles(ctx context.Context) ([]arr.DelayProfileResource, error) {
	var out []arr.DelayProfileResource
	if err := c.base.Get(ctx, "/api/v3/delayprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDelayProfile returns a single delay profile by ID.
func (c *ClientV3) GetDelayProfile(ctx context.Context, id int) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/delayprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDelayProfile creates a new delay profile.
func (c *ClientV3) CreateDelayProfile(ctx context.Context, dp *arr.DelayProfileResource) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Post(ctx, "/api/v3/delayprofile", dp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDelayProfile updates an existing delay profile.
func (c *ClientV3) UpdateDelayProfile(ctx context.Context, dp *arr.DelayProfileResource) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/delayprofile/%d", dp.ID), dp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDelayProfile deletes a delay profile by ID.
func (c *ClientV3) DeleteDelayProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/delayprofile/%d", id), nil, nil)
}

// ReorderDelayProfile moves a delay profile to a new position.
func (c *ClientV3) ReorderDelayProfile(ctx context.Context, id, afterID int) ([]arr.DelayProfileResource, error) {
	var out []arr.DelayProfileResource
	path := fmt.Sprintf("/api/v3/delayprofile/reorder/%d?after=%d", id, afterID)
	if err := c.base.Put(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Download Clients ----------.

// GetDownloadClients returns all configured download clients.
func (c *ClientV3) GetDownloadClients(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/downloadclient", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDownloadClient returns a single download client by ID.
func (c *ClientV3) GetDownloadClient(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/downloadclient/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDownloadClient creates a new download client.
func (c *ClientV3) CreateDownloadClient(ctx context.Context, dc *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/downloadclient", dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClient updates an existing download client.
func (c *ClientV3) UpdateDownloadClient(ctx context.Context, dc *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/downloadclient/%d", dc.ID), dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDownloadClient deletes a download client by ID.
func (c *ClientV3) DeleteDownloadClient(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/downloadclient/%d", id), nil, nil)
}

// GetDownloadClientSchema returns the schema for all download client types.
func (c *ClientV3) GetDownloadClientSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/downloadclient/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestDownloadClient tests a download client configuration.
func (c *ClientV3) TestDownloadClient(ctx context.Context, dc *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/test", dc, nil)
}

// TestAllDownloadClients tests all configured download clients.
func (c *ClientV3) TestAllDownloadClients(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/testall", nil, nil)
}

// BulkUpdateDownloadClients updates multiple download clients.
func (c *ClientV3) BulkUpdateDownloadClients(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/downloadclient/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteDownloadClients deletes multiple download clients.
func (c *ClientV3) BulkDeleteDownloadClients(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/downloadclient/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// DownloadClientAction triggers a named action on a download client.
func (c *ClientV3) DownloadClientAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/downloadclient/action/"+url.PathEscape(name), body, nil)
}

// ---------- Download Client Config ----------.

// GetDownloadClientConfig returns the download client configuration.
func (c *ClientV3) GetDownloadClientConfig(ctx context.Context) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/downloadclient", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClientConfig updates the download client configuration.
func (c *ClientV3) UpdateDownloadClientConfig(ctx context.Context, cfg *arr.DownloadClientConfigResource) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/downloadclient/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Extra Files ----------.

// GetExtraFiles returns extra files for a movie.
func (c *ClientV3) GetExtraFiles(ctx context.Context, movieID int) ([]ExtraFileResource, error) {
	var out []ExtraFileResource
	path := fmt.Sprintf("/api/v3/extrafile?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- File System ----------.

// BrowseFileSystem returns a directory listing.
func (c *ClientV3) BrowseFileSystem(ctx context.Context, path string) (map[string]any, error) {
	var out map[string]any
	reqPath := "/api/v3/filesystem?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFileSystemType returns the filesystem type for a path.
func (c *ClientV3) GetFileSystemType(ctx context.Context, path string) (string, error) {
	var out string
	reqPath := "/api/v3/filesystem/type?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return "", err
	}
	return out, nil
}

// GetFileSystemMediaFiles returns media files in a path.
func (c *ClientV3) GetFileSystemMediaFiles(ctx context.Context, path string) ([]map[string]any, error) {
	var out []map[string]any
	reqPath := "/api/v3/filesystem/mediafiles?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Health Extended ----------.

// GetHealthByID returns a single health check by ID.
func (c *ClientV3) GetHealthByID(ctx context.Context, id int) (*arr.HealthCheck, error) {
	var out arr.HealthCheck
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/health/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- History Extended ----------.

// GetHistorySince returns history since a specific date.
func (c *ClientV3) GetHistorySince(ctx context.Context, date string) ([]HistoryRecordV3, error) {
	var out []HistoryRecordV3
	path := "/api/v3/history/since?date=" + url.QueryEscape(date)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistoryMovie returns history for a specific movie.
func (c *ClientV3) GetHistoryMovie(ctx context.Context, movieID int) ([]HistoryRecordV3, error) {
	var out []HistoryRecordV3
	path := fmt.Sprintf("/api/v3/history/movie?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MarkHistoryFailed marks a history item as failed.
func (c *ClientV3) MarkHistoryFailed(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/history/failed/%d", id), nil, nil)
}

// ---------- Host Config ----------.

// GetHostConfig returns the host configuration.
func (c *ClientV3) GetHostConfig(ctx context.Context) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/host", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateHostConfig updates the host configuration.
func (c *ClientV3) UpdateHostConfig(ctx context.Context, cfg *arr.HostConfigResource) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/host/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Import Lists ----------.

// GetImportLists returns all configured import lists.
func (c *ClientV3) GetImportLists(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/importlist", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetImportList returns a single import list by ID.
func (c *ClientV3) GetImportList(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/importlist/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateImportList creates a new import list.
func (c *ClientV3) CreateImportList(ctx context.Context, il *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/importlist", il, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportList updates an existing import list.
func (c *ClientV3) UpdateImportList(ctx context.Context, il *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/importlist/%d", il.ID), il, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportList deletes an import list by ID.
func (c *ClientV3) DeleteImportList(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/importlist/%d", id), nil, nil)
}

// GetImportListSchema returns the schema for all import list types.
func (c *ClientV3) GetImportListSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/importlist/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestImportList tests an import list configuration.
func (c *ClientV3) TestImportList(ctx context.Context, il *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/importlist/test", il, nil)
}

// TestAllImportLists tests all configured import lists.
func (c *ClientV3) TestAllImportLists(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/importlist/testall", nil, nil)
}

// BulkUpdateImportLists updates multiple import lists.
func (c *ClientV3) BulkUpdateImportLists(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/importlist/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteImportLists deletes multiple import lists.
func (c *ClientV3) BulkDeleteImportLists(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/importlist/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// ImportListAction triggers a named action on an import list.
func (c *ClientV3) ImportListAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/importlist/action/"+url.PathEscape(name), body, nil)
}

// ---------- Import List Movies ----------.

// GetImportListMovies returns movies available from import lists.
func (c *ClientV3) GetImportListMovies(ctx context.Context) ([]Movie, error) {
	var out []Movie
	if err := c.base.Get(ctx, "/api/v3/importlist/movie", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Import Exclusions Extended ----------.

// GetImportExclusion returns a single import exclusion by ID.
func (c *ClientV3) GetImportExclusion(ctx context.Context, id int) (*ImportExclusion, error) {
	var out ImportExclusion
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/exclusions/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateImportExclusion creates a new import exclusion.
func (c *ClientV3) CreateImportExclusion(ctx context.Context, ex *ImportExclusion) (*ImportExclusion, error) {
	var out ImportExclusion
	if err := c.base.Post(ctx, "/api/v3/exclusions", ex, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportExclusion updates an existing import exclusion.
func (c *ClientV3) UpdateImportExclusion(ctx context.Context, ex *ImportExclusion) (*ImportExclusion, error) {
	var out ImportExclusion
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/exclusions/%d", ex.ID), ex, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportExclusion deletes an import exclusion by ID.
func (c *ClientV3) DeleteImportExclusion(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/exclusions/%d", id), nil, nil)
}

// BulkCreateImportExclusions creates multiple import exclusions.
func (c *ClientV3) BulkCreateImportExclusions(ctx context.Context, exs []ImportExclusion) ([]ImportExclusion, error) {
	var out []ImportExclusion
	if err := c.base.Post(ctx, "/api/v3/exclusions/bulk", exs, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteImportExclusions deletes multiple import exclusions.
func (c *ClientV3) BulkDeleteImportExclusions(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/exclusions/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// ---------- Indexers ----------.

// GetIndexers returns all configured indexers.
func (c *ClientV3) GetIndexers(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/indexer", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetIndexer returns a single indexer by ID.
func (c *ClientV3) GetIndexer(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/indexer/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateIndexer creates a new indexer.
func (c *ClientV3) CreateIndexer(ctx context.Context, idx *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/indexer", idx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexer updates an existing indexer.
func (c *ClientV3) UpdateIndexer(ctx context.Context, idx *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/indexer/%d", idx.ID), idx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteIndexer deletes an indexer by ID.
func (c *ClientV3) DeleteIndexer(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/indexer/%d", id), nil, nil)
}

// GetIndexerSchema returns the schema for all indexer types.
func (c *ClientV3) GetIndexerSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/indexer/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestIndexer tests an indexer configuration.
func (c *ClientV3) TestIndexer(ctx context.Context, idx *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/indexer/test", idx, nil)
}

// TestAllIndexers tests all configured indexers.
func (c *ClientV3) TestAllIndexers(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/indexer/testall", nil, nil)
}

// BulkUpdateIndexers updates multiple indexers.
func (c *ClientV3) BulkUpdateIndexers(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v3/indexer/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteIndexers deletes multiple indexers.
func (c *ClientV3) BulkDeleteIndexers(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/indexer/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// IndexerAction triggers a named action on an indexer.
func (c *ClientV3) IndexerAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/indexer/action/"+url.PathEscape(name), body, nil)
}

// ---------- Indexer Config ----------.

// GetIndexerConfig returns the indexer configuration.
func (c *ClientV3) GetIndexerConfig(ctx context.Context) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/indexer", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexerConfig updates the indexer configuration.
func (c *ClientV3) UpdateIndexerConfig(ctx context.Context, cfg *arr.IndexerConfigResource) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/indexer/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Indexer Flags ----------.

// GetIndexerFlags returns all available indexer flags.
func (c *ClientV3) GetIndexerFlags(ctx context.Context) ([]arr.IndexerFlagResource, error) {
	var out []arr.IndexerFlagResource
	if err := c.base.Get(ctx, "/api/v3/indexerflag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Localization ----------.

// GetLocalization returns localization strings.
func (c *ClientV3) GetLocalization(ctx context.Context) (map[string]string, error) {
	var out map[string]string
	if err := c.base.Get(ctx, "/api/v3/localization", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLocalizationLanguage returns the current localization language.
func (c *ClientV3) GetLocalizationLanguage(ctx context.Context) (*LocalizationLanguageResource, error) {
	var out LocalizationLanguageResource
	if err := c.base.Get(ctx, "/api/v3/localization/language", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Logs ----------.

// GetLogs returns log entries (paged).
func (c *ClientV3) GetLogs(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.LogRecord], error) {
	var out arr.PagingResource[arr.LogRecord]
	path := fmt.Sprintf("/api/v3/log?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLogFiles returns available log files.
func (c *ClientV3) GetLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v3/log/file", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLogFileContent returns the content of a specific log file.
func (c *ClientV3) GetLogFileContent(ctx context.Context, filename string) (string, error) {
	b, err := c.base.GetRaw(ctx, "/api/v3/log/file/"+url.PathEscape(filename))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GetUpdateLogFiles returns available update log files.
func (c *ClientV3) GetUpdateLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v3/log/file/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetUpdateLogFileContent returns the content of a specific update log file.
func (c *ClientV3) GetUpdateLogFileContent(ctx context.Context, filename string) (string, error) {
	b, err := c.base.GetRaw(ctx, "/api/v3/log/file/update/"+url.PathEscape(filename))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------- Manual Import ----------.

// GetManualImport returns items available for manual import.
func (c *ClientV3) GetManualImport(ctx context.Context, folder string) ([]arr.ManualImportResource, error) {
	var out []arr.ManualImportResource
	path := "/api/v3/manualimport?folder=" + url.QueryEscape(folder)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessManualImport triggers processing of manual import selections.
func (c *ClientV3) ProcessManualImport(ctx context.Context, items []arr.ManualImportReprocessResource) error {
	return c.base.Post(ctx, "/api/v3/manualimport", items, nil)
}

// ---------- Media Management Config ----------.

// GetMediaManagementConfig returns the media management configuration.
func (c *ClientV3) GetMediaManagementConfig(ctx context.Context) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/mediamanagement", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMediaManagementConfig updates the media management configuration.
func (c *ClientV3) UpdateMediaManagementConfig(ctx context.Context, cfg *arr.MediaManagementConfigResource) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/mediamanagement/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Metadata Consumers ----------.

// GetMetadata returns all metadata consumers.
func (c *ClientV3) GetMetadata(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/metadata", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMetadataByID returns a metadata consumer by ID.
func (c *ClientV3) GetMetadataByID(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/metadata/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateMetadata creates a new metadata consumer.
func (c *ClientV3) CreateMetadata(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/metadata", m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMetadata updates an existing metadata consumer.
func (c *ClientV3) UpdateMetadata(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/metadata/%d", m.ID), m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMetadata deletes a metadata consumer by ID.
func (c *ClientV3) DeleteMetadata(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/metadata/%d", id), nil, nil)
}

// GetMetadataSchema returns the schema for all metadata consumer types.
func (c *ClientV3) GetMetadataSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/metadata/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestMetadata tests a metadata consumer configuration.
func (c *ClientV3) TestMetadata(ctx context.Context, m *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/metadata/test", m, nil)
}

// TestAllMetadata tests all configured metadata consumers.
func (c *ClientV3) TestAllMetadata(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/metadata/testall", nil, nil)
}

// MetadataAction triggers a named action on a metadata consumer.
func (c *ClientV3) MetadataAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/metadata/action/"+url.PathEscape(name), body, nil)
}

// ---------- Movie Extended ----------.

// GetMovieList returns all movies matching a list of IDs.
func (c *ClientV3) GetMovieList(ctx context.Context, movieIDs []int) ([]Movie, error) {
	var out []Movie
	vals := make(url.Values)
	for _, id := range movieIDs {
		vals.Add("movieIds", strconv.Itoa(id))
	}
	path := "/api/v3/movie?" + vals.Encode()
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ImportMovie imports a movie.
func (c *ClientV3) ImportMovie(ctx context.Context, movies []Movie) error {
	return c.base.Post(ctx, "/api/v3/movie/import", movies, nil)
}

// ---------- Movie File Extended ----------.

// GetMovieFiles returns all movie files.
func (c *ClientV3) GetMovieFiles(ctx context.Context, movieID int) ([]MovieFile, error) {
	var out []MovieFile
	path := fmt.Sprintf("/api/v3/moviefile?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateMovieFile updates a movie file.
func (c *ClientV3) UpdateMovieFile(ctx context.Context, mf *MovieFile) (*MovieFile, error) {
	var out MovieFile
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/moviefile/%d", mf.ID), mf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditMovieFiles bulk edits movie files.
func (c *ClientV3) EditMovieFiles(ctx context.Context, editor *MovieFileEditorResource) error {
	return c.base.Put(ctx, "/api/v3/moviefile/editor", editor, nil)
}

// BulkDeleteMovieFiles deletes multiple movie files.
func (c *ClientV3) BulkDeleteMovieFiles(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v3/moviefile/bulk", struct {
		MovieFileIDs []int `json:"movieFileIds"`
	}{MovieFileIDs: ids}, nil)
}

// ---------- Movie Lookup Extended ----------.

// LookupMovieByTMDB searches for a movie by TMDb ID.
func (c *ClientV3) LookupMovieByTMDB(ctx context.Context, tmdbID int) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/lookup/movie/tmdb?tmdbId=%d", tmdbID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupMovieByIMDB searches for a movie by IMDb ID.
func (c *ClientV3) LookupMovieByIMDB(ctx context.Context, imdbID string) (*Movie, error) {
	var out Movie
	path := "/api/v3/lookup/movie/imdb?imdbId=" + url.QueryEscape(imdbID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Naming Config ----------.

// GetNamingConfig returns the naming configuration.
func (c *ClientV3) GetNamingConfig(ctx context.Context) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/naming", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNamingConfig updates the naming configuration.
func (c *ClientV3) UpdateNamingConfig(ctx context.Context, cfg *arr.NamingConfigResource) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/naming/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNamingExamples returns naming examples based on the current config.
func (c *ClientV3) GetNamingExamples(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := c.base.Get(ctx, "/api/v3/config/naming/examples", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Notifications ----------.

// GetNotifications returns all configured notifications.
func (c *ClientV3) GetNotifications(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/notification", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetNotification returns a single notification by ID.
func (c *ClientV3) GetNotification(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/notification/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateNotification creates a new notification.
func (c *ClientV3) CreateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v3/notification", n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNotification updates an existing notification.
func (c *ClientV3) UpdateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/notification/%d", n.ID), n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNotification deletes a notification by ID.
func (c *ClientV3) DeleteNotification(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/notification/%d", id), nil, nil)
}

// GetNotificationSchema returns the schema for all notification types.
func (c *ClientV3) GetNotificationSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v3/notification/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestNotification tests a notification configuration.
func (c *ClientV3) TestNotification(ctx context.Context, n *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/notification/test", n, nil)
}

// TestAllNotifications tests all configured notifications.
func (c *ClientV3) TestAllNotifications(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/notification/testall", nil, nil)
}

// NotificationAction triggers a named action on a notification.
func (c *ClientV3) NotificationAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v3/notification/action/"+url.PathEscape(name), body, nil)
}

// ---------- Performer Editor ----------.

// EditPerformers applies bulk edits to multiple performers.
func (c *ClientV3) EditPerformers(ctx context.Context, editor *PerformerEditorResource) error {
	return c.base.Put(ctx, "/api/v3/performer/editor", editor, nil)
}

// DeletePerformersBulk deletes multiple performers according to the editor payload.
func (c *ClientV3) DeletePerformersBulk(ctx context.Context, editor *PerformerEditorResource) error {
	return c.base.Delete(ctx, "/api/v3/performer/editor", editor, nil)
}

// ---------- Studio Editor ----------.

// EditStudios applies bulk edits to multiple studios.
func (c *ClientV3) EditStudios(ctx context.Context, editor *StudioEditorResource) error {
	return c.base.Put(ctx, "/api/v3/studio/editor", editor, nil)
}

// DeleteStudiosBulk deletes multiple studios according to the editor payload.
func (c *ClientV3) DeleteStudiosBulk(ctx context.Context, editor *StudioEditorResource) error {
	return c.base.Delete(ctx, "/api/v3/studio/editor", editor, nil)
}

// ---------- Quality Definitions ----------.

// GetQualityDefinitions returns all quality definitions.
func (c *ClientV3) GetQualityDefinitions(ctx context.Context) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Get(ctx, "/api/v3/qualitydefinition", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQualityDefinition returns a quality definition by ID.
func (c *ClientV3) GetQualityDefinition(ctx context.Context, id int) (*arr.QualityDefinitionResource, error) {
	var out arr.QualityDefinitionResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/qualitydefinition/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateQualityDefinition updates a quality definition.
func (c *ClientV3) UpdateQualityDefinition(ctx context.Context, qd *arr.QualityDefinitionResource) (*arr.QualityDefinitionResource, error) {
	var out arr.QualityDefinitionResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/qualitydefinition/%d", qd.ID), qd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BulkUpdateQualityDefinitions updates multiple quality definitions.
func (c *ClientV3) BulkUpdateQualityDefinitions(ctx context.Context, defs []arr.QualityDefinitionResource) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Put(ctx, "/api/v3/qualitydefinition/update", defs, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQualityDefinitionLimits returns the quality definition limits.
func (c *ClientV3) GetQualityDefinitionLimits(ctx context.Context) (*QualityDefinitionLimitsResource, error) {
	var out QualityDefinitionLimitsResource
	if err := c.base.Get(ctx, "/api/v3/qualitydefinition/limits", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Quality Profile Extended ----------.

// GetQualityProfile returns a quality profile by ID.
func (c *ClientV3) GetQualityProfile(ctx context.Context, id int) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/qualityprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateQualityProfile creates a new quality profile.
func (c *ClientV3) CreateQualityProfile(ctx context.Context, qp *arr.QualityProfile) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Post(ctx, "/api/v3/qualityprofile", qp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateQualityProfile updates an existing quality profile.
func (c *ClientV3) UpdateQualityProfile(ctx context.Context, qp *arr.QualityProfile) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/qualityprofile/%d", qp.ID), qp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQualityProfile deletes a quality profile by ID.
func (c *ClientV3) DeleteQualityProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/qualityprofile/%d", id), nil, nil)
}

// GetQualityProfileSchema returns the quality profile schema.
func (c *ClientV3) GetQualityProfileSchema(ctx context.Context) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v3/qualityprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Queue Extended ----------.

// DeleteQueueItem removes an item from the download queue.
func (c *ClientV3) DeleteQueueItem(ctx context.Context, id int, removeFromClient, blocklist bool) error {
	path := fmt.Sprintf("/api/v3/queue/%d?removeFromClient=%t&blocklist=%t", id, removeFromClient, blocklist)
	return c.base.Delete(ctx, path, nil, nil)
}

// BulkDeleteQueue removes multiple items from the download queue.
func (c *ClientV3) BulkDeleteQueue(ctx context.Context, bulk *arr.QueueBulkResource) error {
	return c.base.Delete(ctx, "/api/v3/queue/bulk", bulk, nil)
}

// GrabQueueItem sends a queue item to the download client.
func (c *ClientV3) GrabQueueItem(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v3/queue/grab/%d", id), nil, nil)
}

// BulkGrabQueue grabs multiple items from the queue.
func (c *ClientV3) BulkGrabQueue(ctx context.Context, ids []int) error {
	return c.base.Post(ctx, "/api/v3/queue/grab/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// GetQueueDetails returns detailed queue information.
func (c *ClientV3) GetQueueDetails(ctx context.Context) ([]arr.QueueRecord, error) {
	var out []arr.QueueRecord
	if err := c.base.Get(ctx, "/api/v3/queue/details", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueueDetailsByMovieID returns detailed queue info for specific movie IDs.
func (c *ClientV3) GetQueueDetailsByMovieID(ctx context.Context, movieIDs []int) ([]arr.QueueRecord, error) {
	var out []arr.QueueRecord
	vals := make(url.Values)
	for _, id := range movieIDs {
		vals.Add("movieIds", strconv.Itoa(id))
	}
	path := "/api/v3/queue/details?" + vals.Encode()
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueueStatus returns the queue status summary.
func (c *ClientV3) GetQueueStatus(ctx context.Context) (*arr.QueueStatusResource, error) {
	var out arr.QueueStatusResource
	if err := c.base.Get(ctx, "/api/v3/queue/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Releases ----------.

// SearchReleases searches for available releases.
func (c *ClientV3) SearchReleases(ctx context.Context, movieID int) ([]arr.ReleaseResource, error) {
	var out []arr.ReleaseResource
	path := fmt.Sprintf("/api/v3/release?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GrabRelease sends a release to the download client.
func (c *ClientV3) GrabRelease(ctx context.Context, release *arr.ReleaseResource) (*arr.ReleaseResource, error) {
	var out arr.ReleaseResource
	if err := c.base.Post(ctx, "/api/v3/release", release, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PushRelease pushes a release for processing.
func (c *ClientV3) PushRelease(ctx context.Context, push *arr.ReleasePushResource) ([]arr.ReleaseResource, error) {
	var out []arr.ReleaseResource
	if err := c.base.Post(ctx, "/api/v3/release/push", push, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Remote Path Mappings ----------.

// GetRemotePathMappings returns all remote path mappings.
func (c *ClientV3) GetRemotePathMappings(ctx context.Context) ([]arr.RemotePathMappingResource, error) {
	var out []arr.RemotePathMappingResource
	if err := c.base.Get(ctx, "/api/v3/remotepathmapping", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetRemotePathMapping returns a remote path mapping by ID.
func (c *ClientV3) GetRemotePathMapping(ctx context.Context, id int) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/remotepathmapping/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRemotePathMapping creates a new remote path mapping.
func (c *ClientV3) CreateRemotePathMapping(ctx context.Context, rpm *arr.RemotePathMappingResource) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Post(ctx, "/api/v3/remotepathmapping", rpm, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateRemotePathMapping updates an existing remote path mapping.
func (c *ClientV3) UpdateRemotePathMapping(ctx context.Context, rpm *arr.RemotePathMappingResource) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/remotepathmapping/%d", rpm.ID), rpm, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRemotePathMapping deletes a remote path mapping by ID.
func (c *ClientV3) DeleteRemotePathMapping(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/remotepathmapping/%d", id), nil, nil)
}

// ---------- Rename ----------.

// GetRenamePreview returns a rename preview for a movie.
func (c *ClientV3) GetRenamePreview(ctx context.Context, movieID int) ([]arr.RenameEpisodeResource, error) {
	var out []arr.RenameEpisodeResource
	path := fmt.Sprintf("/api/v3/rename?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Root Folders Extended ----------.

// GetRootFolder returns a root folder by ID.
func (c *ClientV3) GetRootFolder(ctx context.Context, id int) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/rootfolder/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRootFolder creates a new root folder.
func (c *ClientV3) CreateRootFolder(ctx context.Context, rf *arr.RootFolder) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Post(ctx, "/api/v3/rootfolder", rf, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRootFolder deletes a root folder by ID.
func (c *ClientV3) DeleteRootFolder(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/rootfolder/%d", id), nil, nil)
}

// ---------- System Extended ----------.

// GetSystemRoutes returns all API routes.
func (c *ClientV3) GetSystemRoutes(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v3/system/routes", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSystemRoutesDuplicate returns duplicate API routes.
func (c *ClientV3) GetSystemRoutesDuplicate(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v3/system/routes/duplicate", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Shutdown sends a shutdown command.
func (c *ClientV3) Shutdown(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/system/shutdown", nil, nil)
}

// Restart sends a restart command.
func (c *ClientV3) Restart(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v3/system/restart", nil, nil)
}

// ---------- Tags Extended ----------.

// GetTag returns a single tag by ID.
func (c *ClientV3) GetTag(ctx context.Context, id int) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/tag/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTag updates an existing tag.
func (c *ClientV3) UpdateTag(ctx context.Context, tag *arr.Tag) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/tag/%d", tag.ID), tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTag deletes a tag by ID.
func (c *ClientV3) DeleteTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v3/tag/%d", id), nil, nil)
}

// GetTagDetails returns all tags with usage details.
func (c *ClientV3) GetTagDetails(ctx context.Context) ([]arr.TagDetail, error) {
	var out []arr.TagDetail
	if err := c.base.Get(ctx, "/api/v3/tag/detail", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTagDetail returns a single tag detail by ID.
func (c *ClientV3) GetTagDetail(ctx context.Context, id int) (*arr.TagDetail, error) {
	var out arr.TagDetail
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/tag/detail/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Tasks ----------.

// GetTasks returns all scheduled tasks.
func (c *ClientV3) GetTasks(ctx context.Context) ([]arr.TaskResource, error) {
	var out []arr.TaskResource
	if err := c.base.Get(ctx, "/api/v3/system/task", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTask returns a scheduled task by ID.
func (c *ClientV3) GetTask(ctx context.Context, id int) (*arr.TaskResource, error) {
	var out arr.TaskResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/system/task/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- UI Config ----------.

// GetUIConfig returns the UI configuration.
func (c *ClientV3) GetUIConfig(ctx context.Context) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Get(ctx, "/api/v3/config/ui", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateUIConfig updates the UI configuration.
func (c *ClientV3) UpdateUIConfig(ctx context.Context, cfg *arr.UIConfigResource) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v3/config/ui/%d", cfg.ID), cfg, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Updates ----------.

// GetUpdates returns available application updates.
func (c *ClientV3) GetUpdates(ctx context.Context) ([]arr.UpdateResource, error) {
	var out []arr.UpdateResource
	if err := c.base.Get(ctx, "/api/v3/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Wanted ----------.

// GetWantedMissing returns missing movies (paged).
func (c *ClientV3) GetWantedMissing(ctx context.Context, page, pageSize int) (*arr.PagingResource[Movie], error) {
	var out arr.PagingResource[Movie]
	path := fmt.Sprintf("/api/v3/wanted/missing?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedMissingByID returns a single missing movie by ID.
func (c *ClientV3) GetWantedMissingByID(ctx context.Context, id int) (*Movie, error) {
	var out Movie
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/wanted/missing/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedCutoff returns cutoff-unmet movies (paged).
func (c *ClientV3) GetWantedCutoff(ctx context.Context, page, pageSize int) (*arr.PagingResource[Movie], error) {
	var out arr.PagingResource[Movie]
	path := fmt.Sprintf("/api/v3/wanted/cutoff?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedCutoffByID returns a single cutoff-unmet movie by ID.
func (c *ClientV3) GetWantedCutoffByID(ctx context.Context, id int) (*Movie, error) {
	var out Movie
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v3/wanted/cutoff/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Ping ----------.

// Ping checks if the Whisparr v3 instance is reachable.
func (c *ClientV3) Ping(ctx context.Context) error {
	return c.base.Get(ctx, "/ping", nil)
}

// ---------- HEAD Ping ----------.

// HeadPing performs a lightweight HEAD request to /ping.
func (c *ClientV3) HeadPing(ctx context.Context) error {
	return c.base.Head(ctx, "/ping")
}

// ---------- Backup Upload ----------.

// UploadBackup uploads a backup file via multipart form POST.
func (c *ClientV3) UploadBackup(ctx context.Context, fileName string, data io.Reader) error {
	return c.base.Upload(ctx, "/api/v3/system/backup/upload", "file", fileName, data)
}
