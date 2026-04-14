package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golusoris/goenvoy/metadata"
)

const (
	defaultBaseURL = "https://graphql.anilist.co"
)

// GraphQL field fragments reused across queries.
//
//nolint:misspell // API uses British "favourites".
const mediaFields = `
	id idMal
	title { romaji english native userPreferred }
	type format status
	description(asHtml: false)
	startDate { year month day }
	endDate { year month day }
	season seasonYear
	episodes duration chapters volumes
	countryOfOrigin isLicensed source
	coverImage { extraLarge large medium color }
	bannerImage
	genres synonyms
	averageScore meanScore popularity favourites
	isAdult siteUrl
	tags { id name rank isGeneralSpoiler isMediaSpoiler }
	nextAiringEpisode { airingAt timeUntilAiring episode }
`

//nolint:misspell // API uses British "favourites".
const mediaSearchFields = `
	id idMal
	title { romaji english native userPreferred }
	type format status
	coverImage { extraLarge large medium color }
	bannerImage
	genres synonyms
	episodes duration chapters volumes
	averageScore meanScore popularity favourites
	isAdult siteUrl
`

//nolint:misspell // API uses British "favourites".
const characterFields = `
	id
	name { first middle last full native alternative alternativeSpoiler userPreferred }
	image { large medium }
	description gender
	dateOfBirth { year month day }
	age siteUrl favourites
`

//nolint:misspell // API uses British "favourites".
const staffFields = `
	id
	name { first middle last full native alternative userPreferred }
	image { large medium }
	description gender
	dateOfBirth { year month day }
	dateOfDeath { year month day }
	age homeTown siteUrl favourites
`

const userFields = `
	id name about
	avatar { large medium }
	bannerImage siteUrl createdAt
`

const pageInfoFields = `pageInfo { total currentPage lastPage hasNextPage perPage }`

// Pre-built query strings.
const (
	queryGetMedia        = `query ($id: Int) { Media(id: $id) {` + mediaFields + `} }`
	queryGetMediaByMalID = `query ($idMal: Int, $type: MediaType) { Media(idMal: $idMal, type: $type) {` + mediaFields + `} }`
	querySearchMedia     = `query ($search: String, $type: MediaType, $page: Int, $perPage: Int) { Page(page: $page, perPage: $perPage) { ` + pageInfoFields + ` media(search: $search, type: $type) {` + mediaSearchFields + `} } }`
	queryGetCharacter    = `query ($id: Int) { Character(id: $id) {` + characterFields + `} }`
	querySearchChars     = `query ($search: String, $page: Int, $perPage: Int) { Page(page: $page, perPage: $perPage) { ` + pageInfoFields + ` characters(search: $search) {` + characterFields + `} } }`
	queryGetStaff        = `query ($id: Int) { Staff(id: $id) {` + staffFields + `} }`
	querySearchStaff     = `query ($search: String, $page: Int, $perPage: Int) { Page(page: $page, perPage: $perPage) { ` + pageInfoFields + ` staff(search: $search) {` + staffFields + `} } }`
	queryGetUser         = `query ($id: Int) { User(id: $id) {` + userFields + `} }`
	queryGetUserByName   = `query ($name: String) { User(name: $name) {` + userFields + `} }`
	queryGetGenres       = `query { GenreCollection }`
	queryGetTags         = `query { MediaTagCollection { id name description category isGeneralSpoiler isMediaSpoiler } }`
)

// Client is an AniList GraphQL API client.
type Client struct {
	*metadata.BaseClient
	accessToken string
}

// New creates an AniList [Client].
func New(opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "anilist", opts...)
	return &Client{BaseClient: bc}
}

// NewWithToken creates an AniList [Client] with an OAuth2 access token.
func NewWithToken(accessToken string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "anilist", opts...)
	c := &Client{BaseClient: bc, accessToken: accessToken}
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	})
	return c
}

type graphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
}

// Query sends a raw GraphQL query and unmarshals the data field into dst.
func (c *Client) Query(ctx context.Context, query string, variables map[string]any, dst any) error {
	payload, err := json.Marshal(graphQLRequest{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("anilist: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("anilist: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("anilist: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("anilist: read response: %w", err)
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		if resp.StatusCode >= 400 {
			return &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status}
		}
		return fmt.Errorf("anilist: decode response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return &APIError{Errors: gqlResp.Errors}
	}

	if dst != nil {
		if err := json.Unmarshal(gqlResp.Data, dst); err != nil {
			return fmt.Errorf("anilist: decode data: %w", err)
		}
	}

	return nil
}

// GetMedia returns a media entry by its AniList ID.
func (c *Client) GetMedia(ctx context.Context, id int) (*Media, error) {
	var resp struct {
		Media Media `json:"Media"`
	}
	if err := c.Query(ctx, queryGetMedia, map[string]any{"id": id}, &resp); err != nil {
		return nil, err
	}
	return &resp.Media, nil
}

// GetMediaByMalID returns a media entry by its MyAnimeList ID and type.
func (c *Client) GetMediaByMalID(ctx context.Context, malID int, mediaType MediaType) (*Media, error) {
	var resp struct {
		Media Media `json:"Media"`
	}
	vars := map[string]any{"idMal": malID, "type": mediaType}
	if err := c.Query(ctx, queryGetMediaByMalID, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Media, nil
}

// SearchMedia searches for media entries with pagination.
func (c *Client) SearchMedia(ctx context.Context, search string, mediaType MediaType, page, perPage int) (*MediaPage, error) {
	var resp struct {
		Page MediaPage `json:"Page"`
	}
	vars := map[string]any{
		"search":  search,
		"type":    mediaType,
		"page":    page,
		"perPage": perPage,
	}
	if err := c.Query(ctx, querySearchMedia, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Page, nil
}

// GetCharacter returns a character by its AniList ID.
func (c *Client) GetCharacter(ctx context.Context, id int) (*Character, error) {
	var resp struct {
		Character Character `json:"Character"`
	}
	if err := c.Query(ctx, queryGetCharacter, map[string]any{"id": id}, &resp); err != nil {
		return nil, err
	}
	return &resp.Character, nil
}

// SearchCharacters searches for characters with pagination.
func (c *Client) SearchCharacters(ctx context.Context, search string, page, perPage int) (*CharacterPage, error) {
	var resp struct {
		Page CharacterPage `json:"Page"`
	}
	vars := map[string]any{"search": search, "page": page, "perPage": perPage}
	if err := c.Query(ctx, querySearchChars, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Page, nil
}

// GetStaff returns a staff member by their AniList ID.
func (c *Client) GetStaff(ctx context.Context, id int) (*Staff, error) {
	var resp struct {
		Staff Staff `json:"Staff"`
	}
	if err := c.Query(ctx, queryGetStaff, map[string]any{"id": id}, &resp); err != nil {
		return nil, err
	}
	return &resp.Staff, nil
}

// SearchStaff searches for staff members with pagination.
func (c *Client) SearchStaff(ctx context.Context, search string, page, perPage int) (*StaffPage, error) {
	var resp struct {
		Page StaffPage `json:"Page"`
	}
	vars := map[string]any{"search": search, "page": page, "perPage": perPage}
	if err := c.Query(ctx, querySearchStaff, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Page, nil
}

// GetUser returns a user by their AniList ID.
func (c *Client) GetUser(ctx context.Context, id int) (*User, error) {
	var resp struct {
		User User `json:"User"`
	}
	if err := c.Query(ctx, queryGetUser, map[string]any{"id": id}, &resp); err != nil {
		return nil, err
	}
	return &resp.User, nil
}

// GetUserByName returns a user by their username.
func (c *Client) GetUserByName(ctx context.Context, name string) (*User, error) {
	var resp struct {
		User User `json:"User"`
	}
	if err := c.Query(ctx, queryGetUserByName, map[string]any{"name": name}, &resp); err != nil {
		return nil, err
	}
	return &resp.User, nil
}

// GetGenres returns all available media genres.
func (c *Client) GetGenres(ctx context.Context) ([]string, error) {
	var resp struct {
		GenreCollection []string `json:"GenreCollection"`
	}
	if err := c.Query(ctx, queryGetGenres, nil, &resp); err != nil {
		return nil, err
	}
	return resp.GenreCollection, nil
}

// GetTags returns all available media tags.
func (c *Client) GetTags(ctx context.Context) ([]MediaTag, error) {
	var resp struct {
		MediaTagCollection []MediaTag `json:"MediaTagCollection"`
	}
	if err := c.Query(ctx, queryGetTags, nil, &resp); err != nil {
		return nil, err
	}
	return resp.MediaTagCollection, nil
}
