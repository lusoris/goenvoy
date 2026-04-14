package discogs

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

const defaultBaseURL = "https://api.discogs.com"

// Client is a Discogs API client.
type Client struct {
	*metadata.BaseClient
	token string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("discogs: %s: %s", e.Status, e.Body)
}

// New creates a new Discogs client.
//
// The token is a Discogs personal access token used for authentication.
func New(token string, opts ...metadata.Option) *Client {
	opts = append([]metadata.Option{metadata.WithUserAgent("goenvoy/1.0")}, opts...)
	bc := metadata.NewBaseClient(defaultBaseURL, "discogs", opts...)
	c := &Client{BaseClient: bc, token: token}
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Discogs token="+token)
	})
	return c
}

func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	u := c.BaseURL() + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("discogs: build request: %w", err)
	}
	req.Header.Set("Authorization", "Discogs token="+c.token)
	req.Header.Set("User-Agent", c.UserAgent())
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("discogs: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("discogs: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("discogs: decode response: %w", err)
	}
	return nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, payload, v any) error {
	u := c.BaseURL() + path

	var bodyReader io.Reader = http.NoBody
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("discogs: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return fmt.Errorf("discogs: build request: %w", err)
	}
	req.Header.Set("Authorization", "Discogs token="+c.token)
	req.Header.Set("User-Agent", c.UserAgent())
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("discogs: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("discogs: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	if v != nil && len(body) > 0 {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("discogs: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) post(ctx context.Context, path string, payload, v any) error {
	return c.doJSON(ctx, http.MethodPost, path, payload, v)
}

func (c *Client) put(ctx context.Context, path string, payload, v any) error {
	return c.doJSON(ctx, http.MethodPut, path, payload, v)
}

func (c *Client) del(ctx context.Context, path string) error {
	return c.doJSON(ctx, http.MethodDelete, path, nil, nil)
}

// GetRelease returns a release by ID.
func (c *Client) GetRelease(ctx context.Context, id int) (*Release, error) {
	var r Release
	if err := c.get(ctx, fmt.Sprintf("/releases/%d", id), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetArtist returns an artist by ID.
func (c *Client) GetArtist(ctx context.Context, id int) (*Artist, error) {
	var a Artist
	if err := c.get(ctx, fmt.Sprintf("/artists/%d", id), nil, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// GetArtistReleases returns releases for an artist.
func (c *Client) GetArtistReleases(ctx context.Context, id, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, fmt.Sprintf("/artists/%d/releases", id), params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// GetLabel returns a label by ID.
func (c *Client) GetLabel(ctx context.Context, id int) (*Label, error) {
	var l Label
	if err := c.get(ctx, fmt.Sprintf("/labels/%d", id), nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetLabelReleases returns releases for a label.
func (c *Client) GetLabelReleases(ctx context.Context, id, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, fmt.Sprintf("/labels/%d/releases", id), params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// GetMasterRelease returns a master release by ID.
func (c *Client) GetMasterRelease(ctx context.Context, id int) (*MasterRelease, error) {
	var m MasterRelease
	if err := c.get(ctx, fmt.Sprintf("/masters/%d", id), nil, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMasterVersions returns all versions of a master release.
func (c *Client) GetMasterVersions(ctx context.Context, id, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, fmt.Sprintf("/masters/%d/versions", id), params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// Search performs a database search.
func (c *Client) Search(ctx context.Context, query, searchType string, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	if query != "" {
		params.Set("q", query)
	}
	if searchType != "" {
		params.Set("type", searchType)
	}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, "/database/search", params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// Database: Release Ratings.

// GetReleaseRatingByUser returns the rating for a release by a specific user.
func (c *Client) GetReleaseRatingByUser(ctx context.Context, releaseID int, username string) (*ReleaseRating, error) {
	var rr ReleaseRating
	if err := c.get(ctx, fmt.Sprintf("/releases/%d/rating/%s", releaseID, url.PathEscape(username)), nil, &rr); err != nil {
		return nil, err
	}
	return &rr, nil
}

// UpdateReleaseRating sets/updates a release rating for the authenticated user.
func (c *Client) UpdateReleaseRating(ctx context.Context, releaseID int, username string, rating int) (*ReleaseRating, error) {
	var rr ReleaseRating
	if err := c.put(ctx, fmt.Sprintf("/releases/%d/rating/%s", releaseID, url.PathEscape(username)), map[string]int{"rating": rating}, &rr); err != nil {
		return nil, err
	}
	return &rr, nil
}

// DeleteReleaseRating removes a release rating for the authenticated user.
func (c *Client) DeleteReleaseRating(ctx context.Context, releaseID int, username string) error {
	return c.del(ctx, fmt.Sprintf("/releases/%d/rating/%s", releaseID, url.PathEscape(username)))
}

// GetCommunityReleaseRating returns the community rating for a release.
func (c *Client) GetCommunityReleaseRating(ctx context.Context, releaseID int) (*CommunityRating, error) {
	var cr CommunityRating
	if err := c.get(ctx, fmt.Sprintf("/releases/%d/rating", releaseID), nil, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

// GetReleaseStats returns marketplace stats for a release.
func (c *Client) GetReleaseStats(ctx context.Context, releaseID int) (*ReleaseStats, error) {
	var rs ReleaseStats
	if err := c.get(ctx, fmt.Sprintf("/releases/%d/stats", releaseID), nil, &rs); err != nil {
		return nil, err
	}
	return &rs, nil
}

// User Identity.

// GetIdentity returns the identity of the authenticated user.
func (c *Client) GetIdentity(ctx context.Context) (*Identity, error) {
	var id Identity
	if err := c.get(ctx, "/oauth/identity", nil, &id); err != nil {
		return nil, err
	}
	return &id, nil
}

// GetProfile returns a user's profile.
func (c *Client) GetProfile(ctx context.Context, username string) (*Profile, error) {
	var p Profile
	if err := c.get(ctx, "/users/"+url.PathEscape(username), nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// EditProfile updates the authenticated user's profile.
func (c *Client) EditProfile(ctx context.Context, username string, updates *ProfileUpdate) (*Profile, error) {
	var p Profile
	if err := c.post(ctx, "/users/"+url.PathEscape(username), updates, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetUserSubmissions returns submissions by a user.
func (c *Client) GetUserSubmissions(ctx context.Context, username string, page, perPage int) (*SubmissionsResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var sr SubmissionsResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/submissions", params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// GetUserContributions returns contributions by a user.
func (c *Client) GetUserContributions(ctx context.Context, username string, page, perPage int) (*ContributionsResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var cr ContributionsResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/contributions", params, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

// User Collection.

// GetCollectionFolders returns collection folders for a user.
func (c *Client) GetCollectionFolders(ctx context.Context, username string) (*CollectionFoldersResponse, error) {
	var cf CollectionFoldersResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/collection/folders", nil, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// CreateCollectionFolder creates a new folder in the user's collection.
func (c *Client) CreateCollectionFolder(ctx context.Context, username, name string) (*CollectionFolder, error) {
	var f CollectionFolder
	if err := c.post(ctx, "/users/"+url.PathEscape(username)+"/collection/folders", map[string]string{"name": name}, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// GetCollectionFolder returns a specific folder from the user's collection.
func (c *Client) GetCollectionFolder(ctx context.Context, username string, folderID int) (*CollectionFolder, error) {
	var f CollectionFolder
	if err := c.get(ctx, fmt.Sprintf("/users/%s/collection/folders/%d", url.PathEscape(username), folderID), nil, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// EditCollectionFolder renames a collection folder.
func (c *Client) EditCollectionFolder(ctx context.Context, username string, folderID int, name string) (*CollectionFolder, error) {
	var f CollectionFolder
	if err := c.post(ctx, fmt.Sprintf("/users/%s/collection/folders/%d", url.PathEscape(username), folderID), map[string]string{"name": name}, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// DeleteCollectionFolder deletes a collection folder.
func (c *Client) DeleteCollectionFolder(ctx context.Context, username string, folderID int) error {
	return c.del(ctx, fmt.Sprintf("/users/%s/collection/folders/%d", url.PathEscape(username), folderID))
}

// GetCollectionItemsByRelease returns instances of a release in the user's collection.
func (c *Client) GetCollectionItemsByRelease(ctx context.Context, username string, releaseID int) (*CollectionItemsResponse, error) {
	var ci CollectionItemsResponse
	if err := c.get(ctx, fmt.Sprintf("/users/%s/collection/releases/%d", url.PathEscape(username), releaseID), nil, &ci); err != nil {
		return nil, err
	}
	return &ci, nil
}

// GetCollectionItemsByFolder returns items in a collection folder.
func (c *Client) GetCollectionItemsByFolder(ctx context.Context, username string, folderID, page, perPage int) (*CollectionItemsResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var ci CollectionItemsResponse
	if err := c.get(ctx, fmt.Sprintf("/users/%s/collection/folders/%d/releases", url.PathEscape(username), folderID), params, &ci); err != nil {
		return nil, err
	}
	return &ci, nil
}

// AddToCollectionFolder adds a release to a collection folder.
func (c *Client) AddToCollectionFolder(ctx context.Context, username string, folderID, releaseID int) (*CollectionItem, error) {
	var ci CollectionItem
	if err := c.post(ctx, fmt.Sprintf("/users/%s/collection/folders/%d/releases/%d", url.PathEscape(username), folderID, releaseID), nil, &ci); err != nil {
		return nil, err
	}
	return &ci, nil
}

// ChangeRatingOfRelease changes the rating of a release instance in the collection.
func (c *Client) ChangeRatingOfRelease(ctx context.Context, username string, folderID, releaseID, instanceID, rating int) error {
	return c.post(ctx, fmt.Sprintf("/users/%s/collection/folders/%d/releases/%d/instances/%d", url.PathEscape(username), folderID, releaseID, instanceID), map[string]int{"rating": rating}, nil)
}

// DeleteInstanceFromFolder removes a release instance from a folder.
func (c *Client) DeleteInstanceFromFolder(ctx context.Context, username string, folderID, releaseID, instanceID int) error {
	return c.del(ctx, fmt.Sprintf("/users/%s/collection/folders/%d/releases/%d/instances/%d", url.PathEscape(username), folderID, releaseID, instanceID))
}

// GetCollectionCustomFields returns the custom fields for a user's collection.
func (c *Client) GetCollectionCustomFields(ctx context.Context, username string) (*CustomFieldsResponse, error) {
	var cf CustomFieldsResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/collection/fields", nil, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// EditFieldsInstance edits a custom field on a collection item instance.
func (c *Client) EditFieldsInstance(ctx context.Context, username string, folderID, releaseID, instanceID, fieldID int, value string) error {
	return c.post(ctx, fmt.Sprintf("/users/%s/collection/folders/%d/releases/%d/instances/%d/fields/%d", url.PathEscape(username), folderID, releaseID, instanceID, fieldID), map[string]string{"value": value}, nil)
}

// GetCollectionValue returns the estimated value of the user's collection.
func (c *Client) GetCollectionValue(ctx context.Context, username string) (*CollectionValue, error) {
	var cv CollectionValue
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/collection/value", nil, &cv); err != nil {
		return nil, err
	}
	return &cv, nil
}

// User Wantlist.

// GetWantlist returns the user's wantlist.
func (c *Client) GetWantlist(ctx context.Context, username string, page, perPage int) (*WantlistResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var wl WantlistResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/wants", params, &wl); err != nil {
		return nil, err
	}
	return &wl, nil
}

// AddToWantlist adds a release to the authenticated user's wantlist.
func (c *Client) AddToWantlist(ctx context.Context, username string, releaseID int, notes string, rating int) (*WantlistItem, error) {
	payload := map[string]any{}
	if notes != "" {
		payload["notes"] = notes
	}
	if rating > 0 {
		payload["rating"] = rating
	}
	var wi WantlistItem
	if err := c.put(ctx, fmt.Sprintf("/users/%s/wants/%d", url.PathEscape(username), releaseID), payload, &wi); err != nil {
		return nil, err
	}
	return &wi, nil
}

// EditWantlistItem edits a release in the user's wantlist.
func (c *Client) EditWantlistItem(ctx context.Context, username string, releaseID int, notes string, rating int) (*WantlistItem, error) {
	payload := map[string]any{}
	if notes != "" {
		payload["notes"] = notes
	}
	if rating > 0 {
		payload["rating"] = rating
	}
	var wi WantlistItem
	if err := c.post(ctx, fmt.Sprintf("/users/%s/wants/%d", url.PathEscape(username), releaseID), payload, &wi); err != nil {
		return nil, err
	}
	return &wi, nil
}

// DeleteFromWantlist removes a release from the user's wantlist.
func (c *Client) DeleteFromWantlist(ctx context.Context, username string, releaseID int) error {
	return c.del(ctx, fmt.Sprintf("/users/%s/wants/%d", url.PathEscape(username), releaseID))
}

// User Lists.

// GetUserLists returns lists created by a user.
func (c *Client) GetUserLists(ctx context.Context, username string, page, perPage int) (*UserListsResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var ul UserListsResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/lists", params, &ul); err != nil {
		return nil, err
	}
	return &ul, nil
}

// GetList returns a specific list by ID.
func (c *Client) GetList(ctx context.Context, listID int) (*UserList, error) {
	var l UserList
	if err := c.get(ctx, fmt.Sprintf("/lists/%d", listID), nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// Marketplace.

// GetInventory returns a user's marketplace inventory.
func (c *Client) GetInventory(ctx context.Context, username string, page, perPage int) (*InventoryResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var ir InventoryResponse
	if err := c.get(ctx, "/users/"+url.PathEscape(username)+"/inventory", params, &ir); err != nil {
		return nil, err
	}
	return &ir, nil
}

// GetListing returns a marketplace listing.
func (c *Client) GetListing(ctx context.Context, listingID int) (*Listing, error) {
	var l Listing
	if err := c.get(ctx, fmt.Sprintf("/marketplace/listings/%d", listingID), nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// CreateListing creates a new marketplace listing.
func (c *Client) CreateListing(ctx context.Context, listing *NewListing) (*Listing, error) {
	var l Listing
	if err := c.post(ctx, "/marketplace/listings", listing, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// EditListing edits a marketplace listing.
func (c *Client) EditListing(ctx context.Context, listingID int, listing *NewListing) error {
	return c.post(ctx, fmt.Sprintf("/marketplace/listings/%d", listingID), listing, nil)
}

// DeleteListing deletes a marketplace listing.
func (c *Client) DeleteListing(ctx context.Context, listingID int) error {
	return c.del(ctx, fmt.Sprintf("/marketplace/listings/%d", listingID))
}

// GetOrder returns a marketplace order.
func (c *Client) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	var o Order
	if err := c.get(ctx, "/marketplace/orders/"+url.PathEscape(orderID), nil, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// EditOrder updates the status of a marketplace order.
func (c *Client) EditOrder(ctx context.Context, orderID string, update *OrderUpdate) (*Order, error) {
	var o Order
	if err := c.post(ctx, "/marketplace/orders/"+url.PathEscape(orderID), update, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// ListOrders returns the authenticated user's orders.
func (c *Client) ListOrders(ctx context.Context, status string, page, perPage int) (*OrdersResponse, error) {
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var or OrdersResponse
	if err := c.get(ctx, "/marketplace/orders", params, &or); err != nil {
		return nil, err
	}
	return &or, nil
}

// GetOrderMessages returns messages for an order.
func (c *Client) GetOrderMessages(ctx context.Context, orderID string, page, perPage int) (*OrderMessagesResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var om OrderMessagesResponse
	if err := c.get(ctx, "/marketplace/orders/"+url.PathEscape(orderID)+"/messages", params, &om); err != nil {
		return nil, err
	}
	return &om, nil
}

// AddOrderMessage adds a message to an order.
func (c *Client) AddOrderMessage(ctx context.Context, orderID, message, status string) (*OrderMessage, error) {
	payload := map[string]string{}
	if message != "" {
		payload["message"] = message
	}
	if status != "" {
		payload["status"] = status
	}
	var om OrderMessage
	if err := c.post(ctx, "/marketplace/orders/"+url.PathEscape(orderID)+"/messages", payload, &om); err != nil {
		return nil, err
	}
	return &om, nil
}

// GetFee returns the marketplace fee for a given price.
func (c *Client) GetFee(ctx context.Context, price float64) (*Fee, error) {
	var f Fee
	if err := c.get(ctx, fmt.Sprintf("/marketplace/fee/%.2f", price), nil, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// GetFeeWithCurrency returns the marketplace fee for a given price and currency.
func (c *Client) GetFeeWithCurrency(ctx context.Context, price float64, currency string) (*Fee, error) {
	var f Fee
	if err := c.get(ctx, fmt.Sprintf("/marketplace/fee/%.2f/%s", price, url.PathEscape(currency)), nil, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// GetPriceSuggestions returns price suggestions for a release.
func (c *Client) GetPriceSuggestions(ctx context.Context, releaseID int) (*PriceSuggestions, error) {
	var ps PriceSuggestions
	if err := c.get(ctx, fmt.Sprintf("/marketplace/price_suggestions/%d", releaseID), nil, &ps); err != nil {
		return nil, err
	}
	return &ps, nil
}

// GetMarketplaceReleaseStats returns marketplace statistics for a release.
func (c *Client) GetMarketplaceReleaseStats(ctx context.Context, releaseID int) (*MarketplaceReleaseStats, error) {
	var mrs MarketplaceReleaseStats
	if err := c.get(ctx, fmt.Sprintf("/marketplace/stats/%d", releaseID), nil, &mrs); err != nil {
		return nil, err
	}
	return &mrs, nil
}

// Inventory Export.

// ExportInventory starts an inventory export.
func (c *Client) ExportInventory(ctx context.Context) error {
	return c.post(ctx, "/inventory/export", nil, nil)
}

// GetRecentExports returns recent inventory exports.
func (c *Client) GetRecentExports(ctx context.Context, page, perPage int) (*ExportsResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	var er ExportsResponse
	if err := c.get(ctx, "/inventory/export", params, &er); err != nil {
		return nil, err
	}
	return &er, nil
}

// GetExport returns a specific inventory export.
func (c *Client) GetExport(ctx context.Context, exportID int) (*Export, error) {
	var e Export
	if err := c.get(ctx, fmt.Sprintf("/inventory/export/%d", exportID), nil, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// DownloadExport returns the download URL for an inventory export.
func (c *Client) DownloadExport(ctx context.Context, exportID int) ([]byte, error) {
	u := c.BaseURL() + fmt.Sprintf("/inventory/export/%d/download", exportID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("discogs: build request: %w", err)
	}
	req.Header.Set("Authorization", "Discogs token="+c.token)
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("discogs: GET export: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("discogs: read export: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}
	return body, nil
}
