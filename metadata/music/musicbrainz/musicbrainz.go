package musicbrainz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://musicbrainz.org/ws/2"

// Client is a MusicBrainz API client.
type Client struct {
	*metadata.BaseClient
}

// New creates a MusicBrainz [Client].
func New(opts ...metadata.Option) *Client {
	opts = append([]metadata.Option{metadata.WithUserAgent("goenvoy/0.0.1 (https://github.com/lusoris/goenvoy)")}, opts...)
	bc := metadata.NewBaseClient(defaultBaseURL, "musicbrainz", opts...)
	return &Client{BaseClient: bc}
}


// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error"`
	// RawBody holds the raw response body when the error could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("musicbrainz: HTTP %d: %s", e.StatusCode, e.Message)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("musicbrainz: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("musicbrainz: HTTP %d", e.StatusCode)
}

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	u, err := url.Parse(c.BaseURL() + path)
	if err != nil {
		return fmt.Errorf("musicbrainz: parse URL: %w", err)
	}
	if params == nil {
		params = url.Values{}
	}
	params.Set("fmt", "json")
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("musicbrainz: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("musicbrainz: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("musicbrainz: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.RawBody = string(body)
		}
		return apiErr
	}

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("musicbrainz: decode response: %w", err)
		}
	}
	return nil
}

func incParams(inc []string) url.Values {
	p := url.Values{}
	if len(inc) > 0 {
		var b strings.Builder
		for i, s := range inc {
			if i > 0 {
				b.WriteByte('+')
			}
			b.WriteString(s)
		}
		p.Set("inc", b.String())
	}
	return p
}

func browseParams(limit, offset int) url.Values {
	p := url.Values{}
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		p.Set("offset", strconv.Itoa(offset))
	}
	return p
}

// Lookup methods — retrieve a single entity by MBID.

// LookupArtist looks up an artist by its MBID.
// Use inc to request subqueries (e.g., "recordings", "releases", "release-groups", "works").
func (c *Client) LookupArtist(ctx context.Context, mbid string, inc []string) (*Artist, error) {
	var out Artist
	if err := c.get(ctx, "/artist/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupRelease looks up a release by its MBID.
func (c *Client) LookupRelease(ctx context.Context, mbid string, inc []string) (*Release, error) {
	var out Release
	if err := c.get(ctx, "/release/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupReleaseGroup looks up a release group by its MBID.
func (c *Client) LookupReleaseGroup(ctx context.Context, mbid string, inc []string) (*ReleaseGroup, error) {
	var out ReleaseGroup
	if err := c.get(ctx, "/release-group/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupRecording looks up a recording by its MBID.
func (c *Client) LookupRecording(ctx context.Context, mbid string, inc []string) (*Recording, error) {
	var out Recording
	if err := c.get(ctx, "/recording/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupLabel looks up a label by its MBID.
func (c *Client) LookupLabel(ctx context.Context, mbid string, inc []string) (*Label, error) {
	var out Label
	if err := c.get(ctx, "/label/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupWork looks up a work by its MBID.
func (c *Client) LookupWork(ctx context.Context, mbid string, inc []string) (*Work, error) {
	var out Work
	if err := c.get(ctx, "/work/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupArea looks up an area by its MBID.
func (c *Client) LookupArea(ctx context.Context, mbid string, inc []string) (*Area, error) {
	var out Area
	if err := c.get(ctx, "/area/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupEvent looks up an event by its MBID.
func (c *Client) LookupEvent(ctx context.Context, mbid string, inc []string) (*Event, error) {
	var out Event
	if err := c.get(ctx, "/event/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupGenre looks up a genre by its MBID.
func (c *Client) LookupGenre(ctx context.Context, mbid string) (*Genre, error) {
	var out Genre
	if err := c.get(ctx, "/genre/"+url.PathEscape(mbid), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupInstrument looks up an instrument by its MBID.
func (c *Client) LookupInstrument(ctx context.Context, mbid string, inc []string) (*Instrument, error) {
	var out Instrument
	if err := c.get(ctx, "/instrument/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupPlace looks up a place by its MBID.
func (c *Client) LookupPlace(ctx context.Context, mbid string, inc []string) (*Place, error) {
	var out Place
	if err := c.get(ctx, "/place/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupSeries looks up a series by its MBID.
func (c *Client) LookupSeries(ctx context.Context, mbid string, inc []string) (*Series, error) {
	var out Series
	if err := c.get(ctx, "/series/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LookupURL looks up a URL entity by its MBID.
func (c *Client) LookupURL(ctx context.Context, mbid string, inc []string) (*URLEntity, error) {
	var out URLEntity
	if err := c.get(ctx, "/url/"+url.PathEscape(mbid), incParams(inc), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Search methods — Lucene-style query search.

// SearchResult holds paginated search results.
type SearchResult[T any] struct {
	Created  string
	Count    int
	Offset   int
	Entities []T
}

// artistSearchResponse matches the JSON response for artist search.
type artistSearchResponse struct {
	Created string   `json:"created"`
	Count   int      `json:"count"`
	Offset  int      `json:"offset"`
	Artists []Artist `json:"artists"`
}

// SearchArtists searches for artists matching the Lucene query.
func (c *Client) SearchArtists(ctx context.Context, query string, limit, offset int) (*SearchResult[Artist], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp artistSearchResponse
	if err := c.get(ctx, "/artist", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Artist]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Artists}, nil
}

type releaseSearchResponse struct {
	Created  string    `json:"created"`
	Count    int       `json:"count"`
	Offset   int       `json:"offset"`
	Releases []Release `json:"releases"`
}

// SearchReleases searches for releases matching the Lucene query.
func (c *Client) SearchReleases(ctx context.Context, query string, limit, offset int) (*SearchResult[Release], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp releaseSearchResponse
	if err := c.get(ctx, "/release", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Release]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Releases}, nil
}

type releaseGroupSearchResponse struct {
	Created       string         `json:"created"`
	Count         int            `json:"count"`
	Offset        int            `json:"offset"`
	ReleaseGroups []ReleaseGroup `json:"release-groups"`
}

// SearchReleaseGroups searches for release groups matching the Lucene query.
func (c *Client) SearchReleaseGroups(ctx context.Context, query string, limit, offset int) (*SearchResult[ReleaseGroup], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp releaseGroupSearchResponse
	if err := c.get(ctx, "/release-group", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[ReleaseGroup]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.ReleaseGroups}, nil
}

type recordingSearchResponse struct {
	Created    string      `json:"created"`
	Count      int         `json:"count"`
	Offset     int         `json:"offset"`
	Recordings []Recording `json:"recordings"`
}

// SearchRecordings searches for recordings matching the Lucene query.
func (c *Client) SearchRecordings(ctx context.Context, query string, limit, offset int) (*SearchResult[Recording], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp recordingSearchResponse
	if err := c.get(ctx, "/recording", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Recording]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Recordings}, nil
}

type labelSearchResponse struct {
	Created string  `json:"created"`
	Count   int     `json:"count"`
	Offset  int     `json:"offset"`
	Labels  []Label `json:"labels"`
}

// SearchLabels searches for labels matching the Lucene query.
func (c *Client) SearchLabels(ctx context.Context, query string, limit, offset int) (*SearchResult[Label], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp labelSearchResponse
	if err := c.get(ctx, "/label", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Label]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Labels}, nil
}

type workSearchResponse struct {
	Created string `json:"created"`
	Count   int    `json:"count"`
	Offset  int    `json:"offset"`
	Works   []Work `json:"works"`
}

// SearchWorks searches for works matching the Lucene query.
func (c *Client) SearchWorks(ctx context.Context, query string, limit, offset int) (*SearchResult[Work], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp workSearchResponse
	if err := c.get(ctx, "/work", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Work]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Works}, nil
}

type areaSearchResponse struct {
	Created string `json:"created"`
	Count   int    `json:"count"`
	Offset  int    `json:"offset"`
	Areas   []Area `json:"areas"`
}

// SearchAreas searches for areas matching the Lucene query.
func (c *Client) SearchAreas(ctx context.Context, query string, limit, offset int) (*SearchResult[Area], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp areaSearchResponse
	if err := c.get(ctx, "/area", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Area]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Areas}, nil
}

type eventSearchResponse struct {
	Created string  `json:"created"`
	Count   int     `json:"count"`
	Offset  int     `json:"offset"`
	Events  []Event `json:"events"`
}

// SearchEvents searches for events matching the Lucene query.
func (c *Client) SearchEvents(ctx context.Context, query string, limit, offset int) (*SearchResult[Event], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp eventSearchResponse
	if err := c.get(ctx, "/event", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Event]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Events}, nil
}

type instrumentSearchResponse struct {
	Created     string       `json:"created"`
	Count       int          `json:"count"`
	Offset      int          `json:"offset"`
	Instruments []Instrument `json:"instruments"`
}

// SearchInstruments searches for instruments matching the Lucene query.
func (c *Client) SearchInstruments(ctx context.Context, query string, limit, offset int) (*SearchResult[Instrument], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp instrumentSearchResponse
	if err := c.get(ctx, "/instrument", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Instrument]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Instruments}, nil
}

type seriesSearchResponse struct {
	Created string   `json:"created"`
	Count   int      `json:"count"`
	Offset  int      `json:"offset"`
	Series  []Series `json:"series"`
}

// SearchSeries searches for series matching the Lucene query.
func (c *Client) SearchSeries(ctx context.Context, query string, limit, offset int) (*SearchResult[Series], error) {
	p := browseParams(limit, offset)
	p.Set("query", query)
	var resp seriesSearchResponse
	if err := c.get(ctx, "/series", p, &resp); err != nil {
		return nil, err
	}
	return &SearchResult[Series]{Created: resp.Created, Count: resp.Count, Offset: resp.Offset, Entities: resp.Series}, nil
}

// ListGenres returns all available genres.
func (c *Client) ListGenres(ctx context.Context) ([]Genre, error) {
	var out []Genre
	if err := c.get(ctx, "/genre/all", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Browse methods — list entities linked to another entity.

type browseArtistsResponse struct {
	ArtistCount  int      `json:"artist-count"`
	ArtistOffset int      `json:"artist-offset"`
	Artists      []Artist `json:"artists"`
}

// BrowseArtistsByReleaseGroup browses artists linked to a release group.
func (c *Client) BrowseArtistsByReleaseGroup(ctx context.Context, rgMBID string, limit, offset int) (*BrowseResult[Artist], error) {
	p := browseParams(limit, offset)
	p.Set("release-group", rgMBID)
	var resp browseArtistsResponse
	if err := c.get(ctx, "/artist", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[Artist]{Entities: resp.Artists, Count: resp.ArtistCount, Offset: resp.ArtistOffset}, nil
}

// BrowseArtistsByRecording browses artists linked to a recording.
func (c *Client) BrowseArtistsByRecording(ctx context.Context, recordingMBID string, limit, offset int) (*BrowseResult[Artist], error) {
	p := browseParams(limit, offset)
	p.Set("recording", recordingMBID)
	var resp browseArtistsResponse
	if err := c.get(ctx, "/artist", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[Artist]{Entities: resp.Artists, Count: resp.ArtistCount, Offset: resp.ArtistOffset}, nil
}

type browseReleasesResponse struct {
	ReleaseCount  int       `json:"release-count"`
	ReleaseOffset int       `json:"release-offset"`
	Releases      []Release `json:"releases"`
}

// BrowseReleasesByArtist browses releases by an artist.
func (c *Client) BrowseReleasesByArtist(ctx context.Context, artistMBID string, limit, offset int) (*BrowseResult[Release], error) {
	p := browseParams(limit, offset)
	p.Set("artist", artistMBID)
	var resp browseReleasesResponse
	if err := c.get(ctx, "/release", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[Release]{Entities: resp.Releases, Count: resp.ReleaseCount, Offset: resp.ReleaseOffset}, nil
}

// BrowseReleasesByLabel browses releases from a label.
func (c *Client) BrowseReleasesByLabel(ctx context.Context, labelMBID string, limit, offset int) (*BrowseResult[Release], error) {
	p := browseParams(limit, offset)
	p.Set("label", labelMBID)
	var resp browseReleasesResponse
	if err := c.get(ctx, "/release", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[Release]{Entities: resp.Releases, Count: resp.ReleaseCount, Offset: resp.ReleaseOffset}, nil
}

type browseReleaseGroupsResponse struct {
	ReleaseGroupCount  int            `json:"release-group-count"`
	ReleaseGroupOffset int            `json:"release-group-offset"`
	ReleaseGroups      []ReleaseGroup `json:"release-groups"`
}

// BrowseReleaseGroupsByArtist browses release groups by an artist.
func (c *Client) BrowseReleaseGroupsByArtist(ctx context.Context, artistMBID string, limit, offset int) (*BrowseResult[ReleaseGroup], error) {
	p := browseParams(limit, offset)
	p.Set("artist", artistMBID)
	var resp browseReleaseGroupsResponse
	if err := c.get(ctx, "/release-group", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[ReleaseGroup]{Entities: resp.ReleaseGroups, Count: resp.ReleaseGroupCount, Offset: resp.ReleaseGroupOffset}, nil
}

type browseRecordingsResponse struct {
	RecordingCount  int         `json:"recording-count"`
	RecordingOffset int         `json:"recording-offset"`
	Recordings      []Recording `json:"recordings"`
}

// BrowseRecordingsByArtist browses recordings by an artist.
func (c *Client) BrowseRecordingsByArtist(ctx context.Context, artistMBID string, limit, offset int) (*BrowseResult[Recording], error) {
	p := browseParams(limit, offset)
	p.Set("artist", artistMBID)
	var resp browseRecordingsResponse
	if err := c.get(ctx, "/recording", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[Recording]{Entities: resp.Recordings, Count: resp.RecordingCount, Offset: resp.RecordingOffset}, nil
}

// BrowseRecordingsByRelease browses recordings on a release.
func (c *Client) BrowseRecordingsByRelease(ctx context.Context, releaseMBID string, limit, offset int) (*BrowseResult[Recording], error) {
	p := browseParams(limit, offset)
	p.Set("release", releaseMBID)
	var resp browseRecordingsResponse
	if err := c.get(ctx, "/recording", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[Recording]{Entities: resp.Recordings, Count: resp.RecordingCount, Offset: resp.RecordingOffset}, nil
}

type browseWorksResponse struct {
	WorkCount  int    `json:"work-count"`
	WorkOffset int    `json:"work-offset"`
	Works      []Work `json:"works"`
}

// BrowseWorksByArtist browses works by an artist.
func (c *Client) BrowseWorksByArtist(ctx context.Context, artistMBID string, limit, offset int) (*BrowseResult[Work], error) {
	p := browseParams(limit, offset)
	p.Set("artist", artistMBID)
	var resp browseWorksResponse
	if err := c.get(ctx, "/work", p, &resp); err != nil {
		return nil, err
	}
	return &BrowseResult[Work]{Entities: resp.Works, Count: resp.WorkCount, Offset: resp.WorkOffset}, nil
}

// Lookup by external IDs.

// LookupByISRC looks up recordings by ISRC code.
func (c *Client) LookupByISRC(ctx context.Context, isrc string, inc []string) ([]Recording, error) {
	type isrcResponse struct {
		Recordings []Recording `json:"recordings"`
	}
	var resp isrcResponse
	if err := c.get(ctx, "/isrc/"+url.PathEscape(isrc), incParams(inc), &resp); err != nil {
		return nil, err
	}
	return resp.Recordings, nil
}

// LookupByDiscID looks up releases by disc ID.
func (c *Client) LookupByDiscID(ctx context.Context, discID string, inc []string) ([]Release, error) {
	type discIDResponse struct {
		Releases []Release `json:"releases"`
	}
	var resp discIDResponse
	if err := c.get(ctx, "/discid/"+url.PathEscape(discID), incParams(inc), &resp); err != nil {
		return nil, err
	}
	return resp.Releases, nil
}
