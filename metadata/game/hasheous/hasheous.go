package hasheous

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://hasheous.org/api/v1"

// Client is a Hasheous API client.
type Client struct {
	*metadata.BaseClient
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("hasheous: %s: %s", e.Status, e.Body)
}

// New creates a Hasheous [Client]. No authentication is required for lookups.
func New(opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "hasheous", opts...)
	return &Client{BaseClient: bc}
}

func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	u := c.BaseURL() + "/" + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("hasheous: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("hasheous: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("hasheous: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("hasheous: decode response: %w", err)
	}
	return nil
}

func (c *Client) post(ctx context.Context, path string, params url.Values, body, v any) error {
	u := c.BaseURL() + "/" + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("hasheous: marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("hasheous: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("hasheous: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("hasheous: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("hasheous: decode response: %w", err)
	}
	return nil
}

// LookupByHash performs a hash lookup using one or more hashes.
// The returnAllSources parameter controls whether all DAT sources are returned.
func (c *Client) LookupByHash(ctx context.Context, req *HashLookupRequest, returnAllSources bool) (*HashLookup, error) {
	params := url.Values{}
	params.Set("returnAllSources", strconv.FormatBool(returnAllSources))
	params.Set("returnFields", "All")

	var result HashLookup
	if err := c.post(ctx, "Lookup/ByHash", params, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LookupByMD5 looks up a ROM by its MD5 hash.
func (c *Client) LookupByMD5(ctx context.Context, md5 string) (*HashLookup, error) {
	var result HashLookup
	if err := c.get(ctx, "Lookup/ByHash/md5/"+url.PathEscape(md5), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LookupBySHA1 looks up a ROM by its SHA1 hash.
func (c *Client) LookupBySHA1(ctx context.Context, sha1 string) (*HashLookup, error) {
	var result HashLookup
	if err := c.get(ctx, "Lookup/ByHash/sha1/"+url.PathEscape(sha1), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LookupBySHA256 looks up a ROM by its SHA256 hash.
func (c *Client) LookupBySHA256(ctx context.Context, sha256 string) (*HashLookup, error) {
	var result HashLookup
	if err := c.get(ctx, "Lookup/ByHash/sha256/"+url.PathEscape(sha256), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LookupByCRC looks up a ROM by its CRC hash.
func (c *Client) LookupByCRC(ctx context.Context, crc string) (*HashLookup, error) {
	var result HashLookup
	if err := c.get(ctx, "Lookup/ByHash/crc/"+url.PathEscape(crc), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlatforms returns a paginated list of platforms in the database.
func (c *Client) GetPlatforms(ctx context.Context, page, pageSize int) ([]string, error) {
	params := url.Values{}
	params.Set("PageNumber", strconv.Itoa(page))
	params.Set("PageSize", strconv.Itoa(pageSize))

	var result []string
	if err := c.get(ctx, "Lookup/Platforms", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}
