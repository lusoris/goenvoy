package googlebooks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://www.googleapis.com/books/v1"

// Client is a Google Books API client.
type Client struct {
	*metadata.BaseClient
	apiKey string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("googlebooks: %s: %s", e.Status, e.Body)
}

// New creates a Google Books [Client] with the given API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "googlebooks", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}


func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("key", c.apiKey)

	u := c.BaseURL() + path + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("googlebooks: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("googlebooks: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("googlebooks: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	return json.Unmarshal(body, v)
}

func (c *Client) post(ctx context.Context, path string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("key", c.apiKey)

	u := c.BaseURL() + path + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("googlebooks: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("googlebooks: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("googlebooks: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}
	return nil
}

// Search searches for volumes by query string.
func (c *Client) Search(ctx context.Context, query string) (*VolumesResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	var resp VolumesResponse
	if err := c.get(ctx, "/volumes", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchWithParams searches for volumes with detailed parameters.
func (c *Client) SearchWithParams(ctx context.Context, p *SearchParams) (*VolumesResponse, error) {
	params := url.Values{}
	if p.Query != "" {
		params.Set("q", p.Query)
	}
	if p.StartIndex > 0 {
		params.Set("startIndex", strconv.Itoa(p.StartIndex))
	}
	if p.MaxResults > 0 {
		params.Set("maxResults", strconv.Itoa(p.MaxResults))
	}
	if p.PrintType != "" {
		params.Set("printType", p.PrintType)
	}
	if p.OrderBy != "" {
		params.Set("orderBy", p.OrderBy)
	}
	if p.Filter != "" {
		params.Set("filter", p.Filter)
	}
	if p.LangRestrict != "" {
		params.Set("langRestrict", p.LangRestrict)
	}
	if p.Projection != "" {
		params.Set("projection", p.Projection)
	}
	var resp VolumesResponse
	if err := c.get(ctx, "/volumes", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVolume returns a single volume by ID.
func (c *Client) GetVolume(ctx context.Context, id string) (*Volume, error) {
	var resp Volume
	if err := c.get(ctx, "/volumes/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUserBookshelves returns public bookshelves for a user.
func (c *Client) GetUserBookshelves(ctx context.Context, userID string) (*BookshelvesResponse, error) {
	var resp BookshelvesResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(userID)+"/bookshelves", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUserBookshelf returns a single public bookshelf for a user.
func (c *Client) GetUserBookshelf(ctx context.Context, userID string, shelfID int) (*Bookshelf, error) {
	var resp Bookshelf
	if err := c.get(ctx, "/users/"+url.PathEscape(userID)+"/bookshelves/"+strconv.Itoa(shelfID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUserBookshelfVolumes returns volumes on a public bookshelf.
func (c *Client) GetUserBookshelfVolumes(ctx context.Context, userID string, shelfID, startIndex, maxResults int) (*VolumesResponse, error) {
	params := url.Values{}
	if startIndex > 0 {
		params.Set("startIndex", strconv.Itoa(startIndex))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	var resp VolumesResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(userID)+"/bookshelves/"+strconv.Itoa(shelfID)+"/volumes", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVolumeAssociatedList returns volumes associated with the given volume.
// association can be "end-of-sample" or "end-of-volume".
func (c *Client) GetVolumeAssociatedList(ctx context.Context, volumeID, association string) (*VolumesResponse, error) {
	params := url.Values{}
	if association != "" {
		params.Set("association", association)
	}
	var resp VolumesResponse
	if err := c.get(ctx, "/volumes/"+url.PathEscape(volumeID)+"/associated", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMyLibraryBookshelves returns bookshelves for the authenticated user.
// Requires OAuth2 access token set via WithHTTPClient.
func (c *Client) GetMyLibraryBookshelves(ctx context.Context) (*BookshelvesResponse, error) {
	var resp BookshelvesResponse
	if err := c.get(ctx, "/mylibrary/bookshelves", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMyLibraryBookshelf returns a single bookshelf for the authenticated user.
func (c *Client) GetMyLibraryBookshelf(ctx context.Context, shelfID int) (*Bookshelf, error) {
	var resp Bookshelf
	if err := c.get(ctx, "/mylibrary/bookshelves/"+strconv.Itoa(shelfID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMyLibraryBookshelfVolumes returns volumes on an authenticated user's bookshelf.
func (c *Client) GetMyLibraryBookshelfVolumes(ctx context.Context, shelfID, startIndex, maxResults int) (*VolumesResponse, error) {
	params := url.Values{}
	if startIndex > 0 {
		params.Set("startIndex", strconv.Itoa(startIndex))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	var resp VolumesResponse
	if err := c.get(ctx, "/mylibrary/bookshelves/"+strconv.Itoa(shelfID)+"/volumes", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddVolumeToBookshelf adds a volume to an authenticated user's bookshelf.
func (c *Client) AddVolumeToBookshelf(ctx context.Context, shelfID int, volumeID string) error {
	params := url.Values{}
	params.Set("volumeId", volumeID)
	return c.post(ctx, "/mylibrary/bookshelves/"+strconv.Itoa(shelfID)+"/addVolume", params)
}

// RemoveVolumeFromBookshelf removes a volume from an authenticated user's bookshelf.
func (c *Client) RemoveVolumeFromBookshelf(ctx context.Context, shelfID int, volumeID string) error {
	params := url.Values{}
	params.Set("volumeId", volumeID)
	return c.post(ctx, "/mylibrary/bookshelves/"+strconv.Itoa(shelfID)+"/removeVolume", params)
}

// ClearBookshelf removes all volumes from an authenticated user's bookshelf.
func (c *Client) ClearBookshelf(ctx context.Context, shelfID int) error {
	return c.post(ctx, "/mylibrary/bookshelves/"+strconv.Itoa(shelfID)+"/clearVolumes", nil)
}

// GetMyLibraryAnnotations returns annotations (highlights, notes) for the authenticated user.
func (c *Client) GetMyLibraryAnnotations(ctx context.Context, volumeID string, startIndex, maxResults int) (*AnnotationsResponse, error) {
	params := url.Values{}
	if volumeID != "" {
		params.Set("volumeId", volumeID)
	}
	if startIndex > 0 {
		params.Set("startIndex", strconv.Itoa(startIndex))
	}
	if maxResults > 0 {
		params.Set("maxResults", strconv.Itoa(maxResults))
	}
	var resp AnnotationsResponse
	if err := c.get(ctx, "/mylibrary/annotations", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMyLibraryReadingPositions returns reading positions for the authenticated user.
func (c *Client) GetMyLibraryReadingPositions(ctx context.Context, volumeID string) (*ReadingPosition, error) {
	var resp ReadingPosition
	if err := c.get(ctx, "/mylibrary/readingpositions/"+url.PathEscape(volumeID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSeries returns information about a series by series ID.
func (c *Client) GetSeries(ctx context.Context, seriesID string) (*SeriesResponse, error) {
	params := url.Values{}
	params.Set("series_id", seriesID)
	var resp SeriesResponse
	if err := c.get(ctx, "/series/get", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSeriesMembers returns volumes belonging to a series.
func (c *Client) GetSeriesMembers(ctx context.Context, seriesID string) (*VolumesResponse, error) {
	params := url.Values{}
	params.Set("series_id", seriesID)
	var resp VolumesResponse
	if err := c.get(ctx, "/series/membership/get", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
