package radarr

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lusoris/goenvoy/arr"
)

// Client is a Radarr API client.
type Client struct {
	base *arr.BaseClient
}

// New creates a Radarr [Client] for the instance at baseURL.
func New(baseURL, apiKey string, opts ...arr.Option) (*Client, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{base: base}, nil
}

// GetAllMovies returns every movie configured in Radarr.
func (c *Client) GetAllMovies(ctx context.Context) ([]Movie, error) {
	var out []Movie
	if err := c.base.Get(ctx, "/api/v3/movie", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovie returns a single movie by its database ID.
func (c *Client) GetMovie(ctx context.Context, id int) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/movie/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddMovie adds a new movie to Radarr.
func (c *Client) AddMovie(ctx context.Context, movie *Movie) (*Movie, error) {
	var out Movie
	if err := c.base.Post(ctx, "/api/v3/movie", movie, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMovie updates an existing movie. Set moveFiles to true to relocate
// files when the movie path changes.
func (c *Client) UpdateMovie(ctx context.Context, movie *Movie, moveFiles bool) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/movie/%d?moveFiles=%t", movie.ID, moveFiles)
	if err := c.base.Put(ctx, path, movie, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMovie removes a movie. Set deleteFiles to true to also delete
// downloaded movie files from disk.
func (c *Client) DeleteMovie(ctx context.Context, id int, deleteFiles, addImportExclusion bool) error {
	path := fmt.Sprintf("/api/v3/movie/%d?deleteFiles=%t&addImportExclusion=%t", id, deleteFiles, addImportExclusion)
	return c.base.Delete(ctx, path, nil)
}

// LookupMovie searches for a movie by term (title).
func (c *Client) LookupMovie(ctx context.Context, term string) ([]Movie, error) {
	var out []Movie
	path := "/api/v3/movie/lookup?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LookupMovieByTmdbID looks up a movie by its TMDb ID.
func (c *Client) LookupMovieByTmdbID(ctx context.Context, tmdbID int) (*Movie, error) {
	var out Movie
	path := fmt.Sprintf("/api/v3/movie/lookup/tmdb?tmdbId=%d", tmdbID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupMovieByImdbID looks up a movie by its IMDb ID.
func (c *Client) LookupMovieByImdbID(ctx context.Context, imdbID string) (*Movie, error) {
	var out Movie
	path := "/api/v3/movie/lookup/imdb?imdbId=" + url.QueryEscape(imdbID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieFiles returns all movie files for the given movie IDs.
func (c *Client) GetMovieFiles(ctx context.Context, movieID int) ([]MovieFile, error) {
	var out []MovieFile
	path := fmt.Sprintf("/api/v3/moviefile?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovieFile returns a single movie file by its database ID.
func (c *Client) GetMovieFile(ctx context.Context, id int) (*MovieFile, error) {
	var out MovieFile
	path := fmt.Sprintf("/api/v3/moviefile/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMovieFile removes a single movie file by its database ID.
func (c *Client) DeleteMovieFile(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v3/moviefile/%d", id)
	return c.base.Delete(ctx, path, nil)
}

// DeleteMovieFiles removes multiple movie files by their IDs.
func (c *Client) DeleteMovieFiles(ctx context.Context, ids []int) error {
	body := MovieFileListResource{MovieFileIDs: ids}
	return c.base.Delete(ctx, "/api/v3/moviefile/bulk", &body)
}

// GetCollections returns all movie collections.
func (c *Client) GetCollections(ctx context.Context) ([]Collection, error) {
	var out []Collection
	if err := c.base.Get(ctx, "/api/v3/collection", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCollection returns a single collection by its database ID.
func (c *Client) GetCollection(ctx context.Context, id int) (*Collection, error) {
	var out Collection
	path := fmt.Sprintf("/api/v3/collection/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCollection updates an existing collection.
func (c *Client) UpdateCollection(ctx context.Context, collection *Collection) (*Collection, error) {
	var out Collection
	path := fmt.Sprintf("/api/v3/collection/%d", collection.ID)
	if err := c.base.Put(ctx, path, collection, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCredits returns cast and crew credits for a movie.
func (c *Client) GetCredits(ctx context.Context, movieID int) ([]Credit, error) {
	var out []Credit
	path := fmt.Sprintf("/api/v3/credit?movieId=%d", movieID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCalendar returns movies with releases between start and end (RFC 3339 timestamps).
func (c *Client) GetCalendar(ctx context.Context, start, end string, unmonitored bool) ([]Movie, error) {
	var out []Movie
	path := fmt.Sprintf("/api/v3/calendar?start=%s&end=%s&unmonitored=%t",
		url.QueryEscape(start), url.QueryEscape(end), unmonitored)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendCommand triggers a named command (e.g. "RefreshMovie", "MoviesSearch").
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

// Parse parses a release title and returns the extracted information.
func (c *Client) Parse(ctx context.Context, title string) (*ParseResult, error) {
	var out ParseResult
	path := "/api/v3/parse?title=" + url.QueryEscape(title)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSystemStatus returns Radarr system information.
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
	return c.base.Delete(ctx, path, nil)
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

// GetImportListExclusions returns all import list exclusions.
func (c *Client) GetImportListExclusions(ctx context.Context) ([]ImportListExclusion, error) {
	var out []ImportListExclusion
	if err := c.base.Get(ctx, "/api/v3/exclusions", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EditMovies performs a batch update on multiple movies.
func (c *Client) EditMovies(ctx context.Context, editor *MovieEditorResource) error {
	return c.base.Put(ctx, "/api/v3/movie/editor", editor, nil)
}

// DeleteMovies performs a batch delete of multiple movies.
func (c *Client) DeleteMovies(ctx context.Context, editor *MovieEditorResource) error {
	return c.base.Delete(ctx, "/api/v3/movie/editor", editor)
}
