package stashbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/lusoris/goenvoy/metadata"
)


// GraphQL field fragments.
const performerFields = `
	id name disambiguation gender
	birthdate deathdate height weight
	hair_color eye_color ethnicity breast_type
	country measurements career_start_year career_end_year
	aliases
	urls { url site { id name url } }
	images { id url width height }
	tattoos { location description }
	piercings { location description }
	is_favorite deleted created updated
`

const sceneFields = `
	id title details release_date production_date
	duration director code
	studio { id name }
	performers { performer { id name } as }
	tags { id name }
	images { id url width height }
	urls { url site { id name url } }
	fingerprints { hash algorithm duration submissions }
	deleted created updated
`

const studioFields = `
	id name aliases
	parent { id name }
	child_studios { id name }
	urls { url site { id name url } }
	images { id url width height }
	deleted created updated
`

const tagFields = `
	id name description aliases
	category { id name description }
	deleted created updated
`

const siteFields = `id name description url regex valid_types icon`

// Pre-built queries.
const (
	queryFindPerformer    = `query ($id: ID!) { findPerformer(id: $id) {` + performerFields + `} }`
	queryQueryPerformers  = `query ($input: PerformerQueryInput!) { queryPerformers(input: $input) { count performers {` + performerFields + `} } }`
	querySearchPerformers = `query ($term: String!, $limit: Int) { searchPerformer(term: $term, limit: $limit) {` + performerFields + `} }`
	queryFindScene        = `query ($id: ID!) { findScene(id: $id) {` + sceneFields + `} }`
	queryQueryScenes      = `query ($input: SceneQueryInput!) { queryScenes(input: $input) { count scenes {` + sceneFields + `} } }`
	querySearchScenes     = `query ($term: String!) { searchScene(term: $term) {` + sceneFields + `} }`
	queryFindFingerprints = `query ($fingerprints: [[FingerprintQueryInput!]!]!) { findScenesBySceneFingerprints(fingerprints: $fingerprints) {` + sceneFields + `} }`
	queryFindStudio       = `query ($id: ID!) { findStudio(id: $id) {` + studioFields + `} }`
	queryQueryStudios     = `query ($input: StudioQueryInput!) { queryStudios(input: $input) { count studios {` + studioFields + `} } }`
	querySearchStudios    = `query ($term: String!) { searchStudio(term: $term) {` + studioFields + `} }`
	queryFindTag          = `query ($id: ID!) { findTag(id: $id) {` + tagFields + `} }`
	queryQueryTags        = `query ($input: TagQueryInput!) { queryTags(input: $input) { count tags {` + tagFields + `} } }`
	querySearchTags       = `query ($term: String!) { searchTag(term: $term) {` + tagFields + `} }`
	queryListSites        = `query { querySites {` + siteFields + `} }`
	queryGetConfig        = `query { getConfig { host require_invite require_activation vote_promotion_threshold vote_application_threshold guidelines_url } }`
)

// Client is a StashBox GraphQL API client.
type Client struct {
	*metadata.BaseClient
	apiKey string
}

// New creates a StashBox [Client].
// endpoint is the full GraphQL URL (e.g., "https://stashdb.org/graphql").
func New(endpoint, apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(endpoint, "stashbox", opts...)
	c := &Client{BaseClient: bc, apiKey: apiKey}
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("ApiKey", apiKey)
	})
	return c
}


// GraphQLError represents an error returned by the GraphQL API.
type GraphQLError struct {
	Message string `json:"message"`
}

func (e *GraphQLError) Error() string {
	return "stashbox: " + e.Message
}

type graphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
}

// APIError is returned when the server responds with a non-2xx status.
type APIError struct {
	StatusCode int
	RawBody    string
}

func (e *APIError) Error() string {
	if e.RawBody != "" {
		return fmt.Sprintf("stashbox: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("stashbox: HTTP %d", e.StatusCode)
}

// Query sends a raw GraphQL query and unmarshals the data field into dst.
func (c *Client) Query(ctx context.Context, query string, variables map[string]any, dst any) error {
	payload, err := json.Marshal(graphQLRequest{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("stashbox: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("stashbox: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ApiKey", c.apiKey)
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("stashbox: execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("stashbox: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, RawBody: string(body)}
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return fmt.Errorf("stashbox: decode response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return &gqlResp.Errors[0]
	}

	if dst != nil {
		if err := json.Unmarshal(gqlResp.Data, dst); err != nil {
			return fmt.Errorf("stashbox: decode data: %w", err)
		}
	}
	return nil
}

// FindPerformer returns a performer by UUID.
func (c *Client) FindPerformer(ctx context.Context, id string) (*Performer, error) {
	var out struct {
		FindPerformer *Performer `json:"findPerformer"`
	}
	if err := c.Query(ctx, queryFindPerformer, map[string]any{"id": id}, &out); err != nil {
		return nil, err
	}
	return out.FindPerformer, nil
}

// QueryPerformers returns a paginated list of performers.
func (c *Client) QueryPerformers(ctx context.Context, input *QueryInput) ([]Performer, int, error) {
	var out struct {
		QueryPerformers struct {
			Count      int         `json:"count"`
			Performers []Performer `json:"performers"`
		} `json:"queryPerformers"`
	}
	vars := map[string]any{"input": input}
	if err := c.Query(ctx, queryQueryPerformers, vars, &out); err != nil {
		return nil, 0, err
	}
	return out.QueryPerformers.Performers, out.QueryPerformers.Count, nil
}

// SearchPerformers searches performers by name.
func (c *Client) SearchPerformers(ctx context.Context, term string, limit int) ([]Performer, error) {
	var out struct {
		SearchPerformer []Performer `json:"searchPerformer"`
	}
	vars := map[string]any{"term": term, "limit": limit}
	if err := c.Query(ctx, querySearchPerformers, vars, &out); err != nil {
		return nil, err
	}
	return out.SearchPerformer, nil
}

// FindScene returns a scene by UUID.
func (c *Client) FindScene(ctx context.Context, id string) (*Scene, error) {
	var out struct {
		FindScene *Scene `json:"findScene"`
	}
	if err := c.Query(ctx, queryFindScene, map[string]any{"id": id}, &out); err != nil {
		return nil, err
	}
	return out.FindScene, nil
}

// QueryScenes returns a paginated list of scenes.
func (c *Client) QueryScenes(ctx context.Context, input *QueryInput) ([]Scene, int, error) {
	var out struct {
		QueryScenes struct {
			Count  int     `json:"count"`
			Scenes []Scene `json:"scenes"`
		} `json:"queryScenes"`
	}
	vars := map[string]any{"input": input}
	if err := c.Query(ctx, queryQueryScenes, vars, &out); err != nil {
		return nil, 0, err
	}
	return out.QueryScenes.Scenes, out.QueryScenes.Count, nil
}

// SearchScenes searches scenes by text.
func (c *Client) SearchScenes(ctx context.Context, term string) ([]Scene, error) {
	var out struct {
		SearchScene []Scene `json:"searchScene"`
	}
	if err := c.Query(ctx, querySearchScenes, map[string]any{"term": term}, &out); err != nil {
		return nil, err
	}
	return out.SearchScene, nil
}

// FindScenesByFingerprints finds scenes by file fingerprints.
// Each inner slice represents fingerprints from a single file.
func (c *Client) FindScenesByFingerprints(ctx context.Context, fingerprints [][]FingerprintInput) ([][]Scene, error) {
	var out struct {
		FindScenesBySceneFingerprints [][]Scene `json:"findScenesBySceneFingerprints"`
	}
	vars := map[string]any{"fingerprints": fingerprints}
	if err := c.Query(ctx, queryFindFingerprints, vars, &out); err != nil {
		return nil, err
	}
	return out.FindScenesBySceneFingerprints, nil
}

// FindStudio returns a studio by UUID.
func (c *Client) FindStudio(ctx context.Context, id string) (*Studio, error) {
	var out struct {
		FindStudio *Studio `json:"findStudio"`
	}
	if err := c.Query(ctx, queryFindStudio, map[string]any{"id": id}, &out); err != nil {
		return nil, err
	}
	return out.FindStudio, nil
}

// QueryStudios returns a paginated list of studios.
func (c *Client) QueryStudios(ctx context.Context, input *QueryInput) ([]Studio, int, error) {
	var out struct {
		QueryStudios struct {
			Count   int      `json:"count"`
			Studios []Studio `json:"studios"`
		} `json:"queryStudios"`
	}
	vars := map[string]any{"input": input}
	if err := c.Query(ctx, queryQueryStudios, vars, &out); err != nil {
		return nil, 0, err
	}
	return out.QueryStudios.Studios, out.QueryStudios.Count, nil
}

// SearchStudios searches studios by name.
func (c *Client) SearchStudios(ctx context.Context, term string) ([]Studio, error) {
	var out struct {
		SearchStudio []Studio `json:"searchStudio"`
	}
	if err := c.Query(ctx, querySearchStudios, map[string]any{"term": term}, &out); err != nil {
		return nil, err
	}
	return out.SearchStudio, nil
}

// FindTag returns a tag by UUID.
func (c *Client) FindTag(ctx context.Context, id string) (*Tag, error) {
	var out struct {
		FindTag *Tag `json:"findTag"`
	}
	if err := c.Query(ctx, queryFindTag, map[string]any{"id": id}, &out); err != nil {
		return nil, err
	}
	return out.FindTag, nil
}

// QueryTags returns a paginated list of tags.
func (c *Client) QueryTags(ctx context.Context, input *QueryInput) ([]Tag, int, error) {
	var out struct {
		QueryTags struct {
			Count int   `json:"count"`
			Tags  []Tag `json:"tags"`
		} `json:"queryTags"`
	}
	vars := map[string]any{"input": input}
	if err := c.Query(ctx, queryQueryTags, vars, &out); err != nil {
		return nil, 0, err
	}
	return out.QueryTags.Tags, out.QueryTags.Count, nil
}

// SearchTags searches tags by name.
func (c *Client) SearchTags(ctx context.Context, term string) ([]Tag, error) {
	var out struct {
		SearchTag []Tag `json:"searchTag"`
	}
	if err := c.Query(ctx, querySearchTags, map[string]any{"term": term}, &out); err != nil {
		return nil, err
	}
	return out.SearchTag, nil
}

// ListSites returns all known content sites.
func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	var out struct {
		QuerySites []Site `json:"querySites"`
	}
	if err := c.Query(ctx, queryListSites, nil, &out); err != nil {
		return nil, err
	}
	return out.QuerySites, nil
}

// Config holds server configuration.
type Config struct {
	Host                     string `json:"host"`
	RequireInvite            bool   `json:"require_invite"`
	RequireActivation        bool   `json:"require_activation"`
	VotePromotionThreshold   *int   `json:"vote_promotion_threshold,omitempty"`
	VoteApplicationThreshold *int   `json:"vote_application_threshold,omitempty"`
	GuidelinesURL            string `json:"guidelines_url,omitempty"`
}

// GetConfig returns the server configuration.
func (c *Client) GetConfig(ctx context.Context) (*Config, error) {
	var out struct {
		GetConfig *Config `json:"getConfig"`
	}
	if err := c.Query(ctx, queryGetConfig, nil, &out); err != nil {
		return nil, err
	}
	return out.GetConfig, nil
}
