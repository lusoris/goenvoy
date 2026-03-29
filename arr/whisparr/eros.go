package whisparr

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lusoris/goenvoy/arr"
)

// ErosClient is a Whisparr v3/Eros (Radarr-based) API client.
type ErosClient struct {
	base *arr.BaseClient
}

// NewEros creates a Whisparr v3/Eros [ErosClient] for the instance at baseURL.
func NewEros(baseURL, apiKey string, opts ...arr.Option) (*ErosClient, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &ErosClient{base: base}, nil
}

// GetAllMovies returns every movie/scene configured in the instance.
func (c *ErosClient) GetAllMovies(ctx context.Context) ([]Movie, error) {
	var out []Movie
	if err := c.base.Get(ctx, "/api/v3/movie", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovie returns a single movie/scene by its database ID.
func (c *ErosClient) GetMovie(ctx context.Context, id int) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/movie/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddMovie adds a new movie/scene.
func (c *ErosClient) AddMovie(ctx context.Context, movie *Movie) (*Movie, error) {
	var out Movie
	if err := c.base.Post(ctx, "/api/v3/movie", movie, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMovie updates an existing movie/scene. Set moveFiles to true to
// relocate files when the path changes.
func (c *ErosClient) UpdateMovie(ctx context.Context, movie *Movie, moveFiles bool) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/movie/%d?moveFiles=%t", movie.ID, moveFiles)
	if err := c.base.Put(ctx, path, movie, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMovie removes a movie/scene by ID.
func (c *ErosClient) DeleteMovie(ctx context.Context, id int, deleteFiles, addImportExclusion bool) error {
	path := fmt.Sprintf("/api/v3/movie/%d?deleteFiles=%t&addImportExclusion=%t", id, deleteFiles, addImportExclusion)
	return c.base.Delete(ctx, path, nil)
}

// LookupMovie searches for a movie by term.
func (c *ErosClient) LookupMovie(ctx context.Context, term string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/lookup/movie?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LookupScene searches for a scene by term.
func (c *ErosClient) LookupScene(ctx context.Context, term string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/lookup/scene?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMoviesByPerformer returns all movies associated with a performer foreign ID.
func (c *ErosClient) GetMoviesByPerformer(ctx context.Context, foreignID string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/movie/listbyperformerforeignid?performerForeignId=" + url.QueryEscape(foreignID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMoviesByStudio returns all movies associated with a studio foreign ID.
func (c *ErosClient) GetMoviesByStudio(ctx context.Context, foreignID string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/movie/listbystudioforeignid?studioForeignId=" + url.QueryEscape(foreignID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovieFile returns a single movie file by ID.
func (c *ErosClient) GetMovieFile(ctx context.Context, id int) (*MovieFile, error) {
	var out MovieFile
	path := fmt.Sprintf("/api/v3/moviefile/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMovieFile deletes a movie file by ID.
func (c *ErosClient) DeleteMovieFile(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v3/moviefile/%d", id)
	return c.base.Delete(ctx, path, nil)
}

// EditMovies applies bulk edits to multiple movies.
func (c *ErosClient) EditMovies(ctx context.Context, editor *MovieEditorResource) error {
	return c.base.Put(ctx, "/api/v3/movie/editor", editor, nil)
}

// DeleteMovies deletes multiple movies according to the editor payload.
func (c *ErosClient) DeleteMovies(ctx context.Context, editor *MovieEditorResource) error {
	return c.base.Delete(ctx, "/api/v3/movie/editor", editor)
}

// GetPerformers returns all performers.
func (c *ErosClient) GetPerformers(ctx context.Context) ([]Performer, error) {
	var out []Performer
	if err := c.base.Get(ctx, "/api/v3/performer", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPerformer returns a single performer by ID.
func (c *ErosClient) GetPerformer(ctx context.Context, id int) (*Performer, error) {
	var out Performer
	path := fmt.Sprintf("/api/v3/performer/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddPerformer adds a new performer to the instance.
func (c *ErosClient) AddPerformer(ctx context.Context, performer *Performer) (*Performer, error) {
	var out Performer
	if err := c.base.Post(ctx, "/api/v3/performer", performer, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdatePerformer updates an existing performer.
func (c *ErosClient) UpdatePerformer(ctx context.Context, performer *Performer) (*Performer, error) {
	var out Performer
	path := fmt.Sprintf("/api/v3/performer/%d", performer.ID)
	if err := c.base.Put(ctx, path, performer, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeletePerformer removes a performer by ID.
func (c *ErosClient) DeletePerformer(ctx context.Context, id int, deleteFiles bool) error {
	path := fmt.Sprintf("/api/v3/performer/%d?deleteFiles=%t", id, deleteFiles)
	return c.base.Delete(ctx, path, nil)
}

// GetStudios returns all studios.
func (c *ErosClient) GetStudios(ctx context.Context) ([]Studio, error) {
	var out []Studio
	if err := c.base.Get(ctx, "/api/v3/studio", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetStudio returns a single studio by ID.
func (c *ErosClient) GetStudio(ctx context.Context, id int) (*Studio, error) {
	var out Studio
	path := fmt.Sprintf("/api/v3/studio/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddStudio adds a new studio to the instance.
func (c *ErosClient) AddStudio(ctx context.Context, studio *Studio) (*Studio, error) {
	var out Studio
	if err := c.base.Post(ctx, "/api/v3/studio", studio, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateStudio updates an existing studio.
func (c *ErosClient) UpdateStudio(ctx context.Context, studio *Studio) (*Studio, error) {
	var out Studio
	path := fmt.Sprintf("/api/v3/studio/%d", studio.ID)
	if err := c.base.Put(ctx, path, studio, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteStudio removes a studio by ID.
func (c *ErosClient) DeleteStudio(ctx context.Context, id int, deleteFiles bool) error {
	path := fmt.Sprintf("/api/v3/studio/%d?deleteFiles=%t", id, deleteFiles)
	return c.base.Delete(ctx, path, nil)
}

// GetCredits returns all credits for a movie/scene.
func (c *ErosClient) GetCredits(ctx context.Context, movieID int) ([]Credit, error) {
	var out []Credit
	path := fmt.Sprintf("/api/v3/credit?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCalendar returns movies/scenes releasing between start and end dates.
func (c *ErosClient) GetCalendar(ctx context.Context, start, end string, unmonitored bool) ([]Movie, error) {
	var out []Movie
	path := fmt.Sprintf("/api/v3/calendar?start=%s&end=%s&unmonitored=%t",
		url.QueryEscape(start), url.QueryEscape(end), unmonitored)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendCommand sends a command to the instance.
func (c *ErosClient) SendCommand(ctx context.Context, cmd arr.CommandRequest) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Post(ctx, "/api/v3/command", cmd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Parse parses a title string and returns matched movie info.
func (c *ErosClient) Parse(ctx context.Context, title string) (*ErosParseResult, error) {
	var out ErosParseResult
	path := "/api/v3/parse?title=" + url.QueryEscape(title)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSystemStatus returns system information.
func (c *ErosClient) GetSystemStatus(ctx context.Context) (*arr.StatusResponse, error) {
	var out arr.StatusResponse
	if err := c.base.Get(ctx, "/api/v3/system/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHealth returns a list of health check results.
func (c *ErosClient) GetHealth(ctx context.Context) ([]arr.HealthCheck, error) {
	var out []arr.HealthCheck
	if err := c.base.Get(ctx, "/api/v3/health", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDiskSpace returns disk space information for all root folders.
func (c *ErosClient) GetDiskSpace(ctx context.Context) ([]arr.DiskSpace, error) {
	var out []arr.DiskSpace
	if err := c.base.Get(ctx, "/api/v3/diskspace", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueue returns the download queue (paged).
func (c *ErosClient) GetQueue(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.QueueRecord], error) {
	var out arr.PagingResource[arr.QueueRecord]
	path := fmt.Sprintf("/api/v3/queue?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQualityProfiles returns all quality profiles.
func (c *ErosClient) GetQualityProfiles(ctx context.Context) ([]arr.QualityProfile, error) {
	var out []arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v3/qualityprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTags returns all tags.
func (c *ErosClient) GetTags(ctx context.Context) ([]arr.Tag, error) {
	var out []arr.Tag
	if err := c.base.Get(ctx, "/api/v3/tag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTag creates a new tag with the given label.
func (c *ErosClient) CreateTag(ctx context.Context, label string) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Post(ctx, "/api/v3/tag", arr.Tag{Label: label}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRootFolders returns all configured root folders.
func (c *ErosClient) GetRootFolders(ctx context.Context) ([]arr.RootFolder, error) {
	var out []arr.RootFolder
	if err := c.base.Get(ctx, "/api/v3/rootfolder", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistory returns history records (paged).
func (c *ErosClient) GetHistory(ctx context.Context, page, pageSize int) (*arr.PagingResource[ErosHistoryRecord], error) {
	var out arr.PagingResource[ErosHistoryRecord]
	path := fmt.Sprintf("/api/v3/history?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetImportExclusions returns all import exclusions.
func (c *ErosClient) GetImportExclusions(ctx context.Context) ([]ImportExclusion, error) {
	var out []ImportExclusion
	if err := c.base.Get(ctx, "/api/v3/exclusions", &out); err != nil {
		return nil, err
	}
	return out, nil
}
