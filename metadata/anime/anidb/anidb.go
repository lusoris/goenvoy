package anidb

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "http://api.anidb.net:9001/httpapi"

// Client is an AniDB HTTP API client.
type Client struct {
	*metadata.BaseClient
	clientName string
	clientVer  int
	titles     titleStore
}

// New creates an AniDB [Client] using the given registered client name and version.
func New(clientName string, clientVer int, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "anidb", opts...)
	return &Client{BaseClient: bc, clientName: clientName, clientVer: clientVer}
}

// APIError is returned when the AniDB API responds with an error.
type APIError struct {
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return "anidb: " + e.Code + ": " + e.Message
}

// xmlErrorResponse maps the <error> root element in error responses.
type xmlErrorResponse struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",chardata"`
}

func (c *Client) get(ctx context.Context, request string, params url.Values, dst any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("client", c.clientName)
	params.Set("clientver", strconv.Itoa(c.clientVer))
	params.Set("protover", "1")
	params.Set("request", request)

	u := c.BaseURL() + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("anidb: create request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("anidb: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("anidb: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{Code: strconv.Itoa(resp.StatusCode), Message: http.StatusText(resp.StatusCode)}
	}

	var xmlErr xmlErrorResponse
	if xml.Unmarshal(body, &xmlErr) == nil {
		return &APIError{Code: xmlErr.Code, Message: xmlErr.Message}
	}

	if dst != nil {
		if err := xml.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("anidb: decode response: %w", err)
		}
	}

	return nil
}

// hotAnimeResponse wraps the XML root for hot anime responses.
type hotAnimeResponse struct {
	XMLName xml.Name   `xml:"hotanime"`
	Entries []AnimeRef `xml:"anime"`
}

// randomRecommendationResponse wraps the XML root for random recommendation responses.
type randomRecommendationResponse struct {
	XMLName xml.Name              `xml:"randomrecommendation"`
	Entries []RecommendationEntry `xml:"recommendation"`
}

// randomSimilarResponse wraps the XML root for random similar responses.
type randomSimilarResponse struct {
	XMLName xml.Name      `xml:"randomsimilar"`
	Entries []SimilarPair `xml:"similar"`
}

// GetAnime retrieves full details for an anime by its AniDB ID.
func (c *Client) GetAnime(ctx context.Context, aid int) (*Anime, error) {
	params := url.Values{}
	params.Set("aid", strconv.Itoa(aid))
	var out Anime
	if err := c.get(ctx, "anime", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HotAnime returns the list of currently hot anime.
func (c *Client) HotAnime(ctx context.Context) ([]AnimeRef, error) {
	var resp hotAnimeResponse
	if err := c.get(ctx, "hotanime", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// RandomRecommendation returns a list of random anime recommendations.
func (c *Client) RandomRecommendation(ctx context.Context) ([]RecommendationEntry, error) {
	var resp randomRecommendationResponse
	if err := c.get(ctx, "randomrecommendation", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// RandomSimilar returns a list of random similar anime pairs.
func (c *Client) RandomSimilar(ctx context.Context) ([]SimilarPair, error) {
	var resp randomSimilarResponse
	if err := c.get(ctx, "randomsimilar", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// MainPage returns combined hot anime, random similar, and random recommendation data.
func (c *Client) MainPage(ctx context.Context) (*MainPage, error) {
	var out MainPage
	if err := c.get(ctx, "main", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
