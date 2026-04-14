package letterboxd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.letterboxd.com/api/v0"

// TokenCallback is called whenever a new token is acquired or refreshed.
// This allows consumers to persist or log token events.
type TokenCallback func(accessToken string, expiresIn int)

// Client is a Letterboxd v0 API client.
type Client struct {
	*metadata.BaseClient
	accessToken  string
	clientID     string
	clientSecret string
	onToken      TokenCallback
	mu           sync.Mutex
	tokenExpiry  time.Time
}

// New creates a Letterboxd [Client] using the given OAuth2 Bearer token.
func New(accessToken string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "letterboxd", opts...)
	return &Client{BaseClient: bc, accessToken: accessToken}
}

// NewWithClientCredentials creates a Letterboxd [Client] that automatically
// acquires and refreshes an OAuth2 access token using the client_credentials grant.
func NewWithClientCredentials(clientID, clientSecret string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "letterboxd", opts...)
	return &Client{BaseClient: bc, clientID: clientID, clientSecret: clientSecret}
}

// SetTokenCallback sets a callback for token lifecycle events.
func (c *Client) SetTokenCallback(cb TokenCallback) { c.onToken = cb }

// tokenResponse is the OAuth2 token endpoint response.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// ensureToken acquires or refreshes the OAuth2 token if client credentials are configured.
func (c *Client) ensureToken(ctx context.Context) error {
	if c.clientID == "" {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return nil
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.BaseURL()+"/auth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("letterboxd: create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("letterboxd: token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("letterboxd: read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("letterboxd: token request HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tok tokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return fmt.Errorf("letterboxd: decode token response: %w", err)
	}

	c.accessToken = tok.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn)*time.Second - 30*time.Second)

	if c.onToken != nil {
		c.onToken(tok.AccessToken, tok.ExpiresIn)
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string, dst any) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}

	c.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	})

	body, status, err := c.DoRaw(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		apiErr := &APIError{StatusCode: status}
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.RawBody = string(body)
		}
		return apiErr
	}

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("letterboxd: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, payload, dst any) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}

	var bodyReader io.Reader = http.NoBody
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("letterboxd: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	c.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
	})

	body, status, err := c.DoRaw(ctx, method, path, bodyReader)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		apiErr := &APIError{StatusCode: status}
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.RawBody = string(body)
		}
		return apiErr
	}

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("letterboxd: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) post(ctx context.Context, path string, payload, dst any) error {
	return c.doJSON(ctx, http.MethodPost, path, payload, dst)
}

func (c *Client) patch(ctx context.Context, path string, payload, dst any) error {
	return c.doJSON(ctx, http.MethodPatch, path, payload, dst)
}

func (c *Client) del(ctx context.Context, path string) error {
	return c.doJSON(ctx, http.MethodDelete, path, nil, nil)
}

func cursorParams(cursor string, perPage int) string {
	v := url.Values{}
	if cursor != "" {
		v.Set("cursor", cursor)
	}
	if perPage > 0 {
		v.Set("perPage", strconv.Itoa(perPage))
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}

// Films.

// GetFilm returns details for a film by its Letterboxd ID, TMDB ID (tmdb:123), or IMDB ID (imdb:tt123).
func (c *Client) GetFilm(ctx context.Context, id string) (*Film, error) {
	var f Film
	if err := c.get(ctx, "/film/"+url.PathEscape(id), &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// GetFilmStatistics returns statistical data for a film.
func (c *Client) GetFilmStatistics(ctx context.Context, id string) (*FilmStatistics, error) {
	var s FilmStatistics
	if err := c.get(ctx, "/film/"+url.PathEscape(id)+"/statistics", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetFilmRelationship returns the authenticated member's relationship with a film.
func (c *Client) GetFilmRelationship(ctx context.Context, id string) (*FilmRelationship, error) {
	var r FilmRelationship
	if err := c.get(ctx, "/film/"+url.PathEscape(id)+"/me", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateFilmRelationship updates the authenticated member's relationship with a film (watched, liked, watchlist, rating).
func (c *Client) UpdateFilmRelationship(ctx context.Context, id string, req FilmRelationshipUpdateRequest) (*FilmRelationshipUpdateResponse, error) {
	var resp FilmRelationshipUpdateResponse
	if err := c.patch(ctx, "/film/"+url.PathEscape(id)+"/me", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFilmMemberRelationship returns a specific member's relationship with a film.
func (c *Client) GetFilmMemberRelationship(ctx context.Context, filmID, memberID string) (*FilmRelationship, error) {
	var r FilmRelationship
	if err := c.get(ctx, "/film/"+url.PathEscape(filmID)+"/member/"+url.PathEscape(memberID), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetFilmMembers returns a cursored list of member relationships with a film.
func (c *Client) GetFilmMembers(ctx context.Context, id, cursor string, perPage int) (*FilmMembersResponse, error) {
	var r FilmMembersResponse
	if err := c.get(ctx, "/film/"+url.PathEscape(id)+"/members"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetFilmFriends returns the authenticated member's friends' relationships with a film.
func (c *Client) GetFilmFriends(ctx context.Context, id string) (*FilmFriendsResponse, error) {
	var r FilmFriendsResponse
	if err := c.get(ctx, "/film/"+url.PathEscape(id)+"/friends", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ReportFilm reports a film to the moderators.
func (c *Client) ReportFilm(ctx context.Context, id string, req ReportRequest) error {
	return c.post(ctx, "/film/"+url.PathEscape(id)+"/report", req, nil)
}

// GetFilms returns a cursored list of films with optional filtering.
func (c *Client) GetFilms(ctx context.Context, cursor string, perPage int) (*FilmsResponse, error) {
	var r FilmsResponse
	if err := c.get(ctx, "/films"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetGenres returns all supported film genres.
func (c *Client) GetGenres(ctx context.Context) (*GenresResponse, error) {
	var r GenresResponse
	if err := c.get(ctx, "/films/genres", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetFilmServices returns all supported streaming services.
func (c *Client) GetFilmServices(ctx context.Context) (*FilmServicesResponse, error) {
	var r FilmServicesResponse
	if err := c.get(ctx, "/films/film-services", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetCountries returns all supported film countries.
func (c *Client) GetCountries(ctx context.Context) (*CountriesResponse, error) {
	var r CountriesResponse
	if err := c.get(ctx, "/films/countries", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetLanguages returns all supported film languages.
func (c *Client) GetLanguages(ctx context.Context) (*LanguagesResponse, error) {
	var r LanguagesResponse
	if err := c.get(ctx, "/films/languages", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Film collections.

// GetFilmCollection returns details for a film collection.
func (c *Client) GetFilmCollection(ctx context.Context, id string) (*FilmCollection, error) {
	var fc FilmCollection
	if err := c.get(ctx, "/film-collection/"+url.PathEscape(id), &fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

// GetFilmCollections returns a cursored list of film collections.
func (c *Client) GetFilmCollections(ctx context.Context, cursor string, perPage int) (*FilmCollectionsResponse, error) {
	var r FilmCollectionsResponse
	if err := c.get(ctx, "/film-collections"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Contributors.

// GetContributor returns details for a film contributor (actor, director, etc.).
func (c *Client) GetContributor(ctx context.Context, id string) (*Contributor, error) {
	var ct Contributor
	if err := c.get(ctx, "/contributor/"+url.PathEscape(id), &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

// GetContributions returns a cursored list of a contributor's film contributions.
func (c *Client) GetContributions(ctx context.Context, id, cursor string, perPage int) (*ContributionsResponse, error) {
	var r ContributionsResponse
	if err := c.get(ctx, "/contributor/"+url.PathEscape(id)+"/contributions"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Lists.

// GetList returns details for a list.
func (c *Client) GetList(ctx context.Context, id string) (*List, error) {
	var l List
	if err := c.get(ctx, "/list/"+url.PathEscape(id), &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetListEntries returns a cursored list of entries (films) in a list.
func (c *Client) GetListEntries(ctx context.Context, id, cursor string, perPage int) (*ListEntriesResponse, error) {
	var r ListEntriesResponse
	if err := c.get(ctx, "/list/"+url.PathEscape(id)+"/entries"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetListComments returns a cursored list of comments on a list.
func (c *Client) GetListComments(ctx context.Context, id, cursor string, perPage int) (*CommentsResponse, error) {
	var r CommentsResponse
	if err := c.get(ctx, "/list/"+url.PathEscape(id)+"/comments"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetListStatistics returns statistical data for a list.
func (c *Client) GetListStatistics(ctx context.Context, id string) (*ListStatistics, error) {
	var s ListStatistics
	if err := c.get(ctx, "/list/"+url.PathEscape(id)+"/statistics", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetListRelationship returns the authenticated member's relationship with a list.
func (c *Client) GetListRelationship(ctx context.Context, id string) (*ListRelationship, error) {
	var r ListRelationship
	if err := c.get(ctx, "/list/"+url.PathEscape(id)+"/me", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateListRelationship updates the authenticated member's relationship with a list (like/unlike).
func (c *Client) UpdateListRelationship(ctx context.Context, id string, req ListRelationshipUpdateRequest) (*ListRelationshipUpdateResponse, error) {
	var resp ListRelationshipUpdateResponse
	if err := c.patch(ctx, "/list/"+url.PathEscape(id)+"/me", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLists returns a cursored list of lists.
func (c *Client) GetLists(ctx context.Context, cursor string, perPage int) (*ListsResponse, error) {
	var r ListsResponse
	if err := c.get(ctx, "/lists"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetListTopics returns featured list topics.
func (c *Client) GetListTopics(ctx context.Context) (*ListTopicsResponse, error) {
	var r ListTopicsResponse
	if err := c.get(ctx, "/lists/topics", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateList creates a new list.
func (c *Client) CreateList(ctx context.Context, req *ListCreationRequest) (*ListCreateResponse, error) {
	var resp ListCreateResponse
	if err := c.post(ctx, "/lists", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateList updates an existing list.
func (c *Client) UpdateList(ctx context.Context, id string, req *ListUpdateRequest) (*ListUpdateResponse, error) {
	var resp ListUpdateResponse
	if err := c.patch(ctx, "/list/"+url.PathEscape(id), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteList deletes a list.
func (c *Client) DeleteList(ctx context.Context, id string) error {
	return c.del(ctx, "/list/"+url.PathEscape(id))
}

// CreateListComment creates a comment on a list.
func (c *Client) CreateListComment(ctx context.Context, id string, req CommentCreationRequest) (*Comment, error) {
	var resp Comment
	if err := c.post(ctx, "/list/"+url.PathEscape(id)+"/comments", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReportList reports a list to the moderators.
func (c *Client) ReportList(ctx context.Context, id string, req ReportRequest) error {
	return c.post(ctx, "/list/"+url.PathEscape(id)+"/report", req, nil)
}

// ForgetList removes a shared list from the authenticated member's account.
func (c *Client) ForgetList(ctx context.Context, id string) error {
	return c.post(ctx, "/list/"+url.PathEscape(id)+"/forget", nil, nil)
}

// AddFilmsToLists adds films to one or more lists.
func (c *Client) AddFilmsToLists(ctx context.Context, req ListAddEntriesRequest) error {
	return c.patch(ctx, "/lists", req, nil)
}

// Log entries.

// GetLogEntry returns details for a log entry (diary entry or review).
func (c *Client) GetLogEntry(ctx context.Context, id string) (*LogEntry, error) {
	var e LogEntry
	if err := c.get(ctx, "/log-entry/"+url.PathEscape(id), &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// GetLogEntries returns a cursored list of log entries.
func (c *Client) GetLogEntries(ctx context.Context, cursor string, perPage int) (*LogEntriesResponse, error) {
	var r LogEntriesResponse
	if err := c.get(ctx, "/log-entries"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetLogEntryComments returns a cursored list of comments on a log entry.
func (c *Client) GetLogEntryComments(ctx context.Context, id, cursor string, perPage int) (*CommentsResponse, error) {
	var r CommentsResponse
	if err := c.get(ctx, "/log-entry/"+url.PathEscape(id)+"/comments"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetLogEntryStatistics returns statistical data for a log entry.
func (c *Client) GetLogEntryStatistics(ctx context.Context, id string) (*LogEntryStatistics, error) {
	var s LogEntryStatistics
	if err := c.get(ctx, "/log-entry/"+url.PathEscape(id)+"/statistics", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetLogEntryRelationship returns the authenticated member's relationship with a log entry.
func (c *Client) GetLogEntryRelationship(ctx context.Context, id string) (*LogEntryRelationship, error) {
	var r LogEntryRelationship
	if err := c.get(ctx, "/log-entry/"+url.PathEscape(id)+"/me", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateLogEntryRelationship updates the authenticated member's relationship with a log entry (like/subscribe).
func (c *Client) UpdateLogEntryRelationship(ctx context.Context, id string, req LogEntryRelationshipUpdateRequest) (*LogEntryRelationshipUpdateResponse, error) {
	var resp LogEntryRelationshipUpdateResponse
	if err := c.patch(ctx, "/log-entry/"+url.PathEscape(id)+"/me", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLogEntryMembers returns member relationships with the film from a log entry.
func (c *Client) GetLogEntryMembers(ctx context.Context, id, cursor string, perPage int) (*FilmMembersResponse, error) {
	var r FilmMembersResponse
	if err := c.get(ctx, "/log-entry/"+url.PathEscape(id)+"/members"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateLogEntry creates a new log entry (diary entry and/or review).
func (c *Client) CreateLogEntry(ctx context.Context, req *LogEntryCreationRequest) (*LogEntry, error) {
	var e LogEntry
	if err := c.post(ctx, "/log-entries", req, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateLogEntry updates an existing log entry.
func (c *Client) UpdateLogEntry(ctx context.Context, id string, req *LogEntryUpdateRequest) (*LogEntry, error) {
	var e LogEntry
	if err := c.patch(ctx, "/log-entry/"+url.PathEscape(id), req, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// DeleteLogEntry deletes a log entry.
func (c *Client) DeleteLogEntry(ctx context.Context, id string) error {
	return c.del(ctx, "/log-entry/"+url.PathEscape(id))
}

// CreateLogEntryComment creates a comment on a log entry.
func (c *Client) CreateLogEntryComment(ctx context.Context, id string, req CommentCreationRequest) (*Comment, error) {
	var resp Comment
	if err := c.post(ctx, "/log-entry/"+url.PathEscape(id)+"/comments", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReportLogEntry reports a log entry to the moderators.
func (c *Client) ReportLogEntry(ctx context.Context, id string, req ReportRequest) error {
	return c.post(ctx, "/log-entry/"+url.PathEscape(id)+"/report", req, nil)
}

// Members.

// GetMember returns details for a member.
func (c *Client) GetMember(ctx context.Context, id string) (*Member, error) {
	var m Member
	if err := c.get(ctx, "/member/"+url.PathEscape(id), &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMemberStatistics returns statistical data for a member.
func (c *Client) GetMemberStatistics(ctx context.Context, id string) (*MemberStatistics, error) {
	var s MemberStatistics
	if err := c.get(ctx, "/member/"+url.PathEscape(id)+"/statistics", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetMemberRelationship returns the authenticated member's relationship with another member.
func (c *Client) GetMemberRelationship(ctx context.Context, id string) (*MemberRelationship, error) {
	var r MemberRelationship
	if err := c.get(ctx, "/member/"+url.PathEscape(id)+"/me", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateMemberRelationship updates the authenticated member's relationship with another member (follow/block).
func (c *Client) UpdateMemberRelationship(ctx context.Context, id string, req MemberRelationshipUpdateRequest) (*MemberRelationshipUpdateResponse, error) {
	var resp MemberRelationshipUpdateResponse
	if err := c.patch(ctx, "/member/"+url.PathEscape(id)+"/me", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMemberActivity returns a cursored activity feed for a member.
func (c *Client) GetMemberActivity(ctx context.Context, id, cursor string, perPage int) (*ActivityResponse, error) {
	var r ActivityResponse
	if err := c.get(ctx, "/member/"+url.PathEscape(id)+"/activity"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetMemberWatchlist returns a cursored list of films on a member's watchlist.
func (c *Client) GetMemberWatchlist(ctx context.Context, id, cursor string, perPage int) (*FilmsResponse, error) {
	var r FilmsResponse
	if err := c.get(ctx, "/member/"+url.PathEscape(id)+"/watchlist"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetMemberLogEntryTags returns a member's log entry tags for autocomplete.
func (c *Client) GetMemberLogEntryTags(ctx context.Context, id string) (*TagsResponse, error) {
	var r TagsResponse
	if err := c.get(ctx, "/member/"+url.PathEscape(id)+"/log-entry-tags", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetMemberListTags returns a member's list tags for autocomplete.
func (c *Client) GetMemberListTags(ctx context.Context, id string) (*TagsResponse, error) {
	var r TagsResponse
	if err := c.get(ctx, "/member/"+url.PathEscape(id)+"/list-tags-2", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetMembers returns a cursored list of members.
func (c *Client) GetMembers(ctx context.Context, cursor string, perPage int) (*MembersResponse, error) {
	var r MembersResponse
	if err := c.get(ctx, "/members"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetPronouns returns all supported pronouns.
func (c *Client) GetPronouns(ctx context.Context) (*PronounsResponse, error) {
	var r PronounsResponse
	if err := c.get(ctx, "/members/pronouns", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ReportMember reports a member to the moderators.
func (c *Client) ReportMember(ctx context.Context, id string, req ReportRequest) error {
	return c.post(ctx, "/member/"+url.PathEscape(id)+"/report", req, nil)
}

// Me.

// GetMe returns the authenticated member's account details.
func (c *Client) GetMe(ctx context.Context) (*MemberAccount, error) {
	var a MemberAccount
	if err := c.get(ctx, "/me", &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateMe updates the authenticated member's account settings.
func (c *Client) UpdateMe(ctx context.Context, req *MemberSettingsUpdateRequest) (*MemberSettingsUpdateResponse, error) {
	var resp MemberSettingsUpdateResponse
	if err := c.patch(ctx, "/me", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Comments.

// UpdateComment updates a comment's message.
func (c *Client) UpdateComment(ctx context.Context, id string, req CommentUpdateRequest) (*Comment, error) {
	var resp Comment
	if err := c.patch(ctx, "/comment/"+url.PathEscape(id), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteComment deletes a comment.
func (c *Client) DeleteComment(ctx context.Context, id string) error {
	return c.del(ctx, "/comment/"+url.PathEscape(id))
}

// ReportComment reports a comment to the moderators.
func (c *Client) ReportComment(ctx context.Context, id string, req ReportRequest) error {
	return c.post(ctx, "/comment/"+url.PathEscape(id)+"/report", req, nil)
}

// Stories.

// GetStory returns details for a story.
func (c *Client) GetStory(ctx context.Context, id string) (*Story, error) {
	var s Story
	if err := c.get(ctx, "/story/"+url.PathEscape(id), &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStories returns a cursored list of stories.
func (c *Client) GetStories(ctx context.Context, cursor string, perPage int) (*StoriesResponse, error) {
	var r StoriesResponse
	if err := c.get(ctx, "/stories"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetStoryComments returns a cursored list of comments on a story.
func (c *Client) GetStoryComments(ctx context.Context, id, cursor string, perPage int) (*CommentsResponse, error) {
	var r CommentsResponse
	if err := c.get(ctx, "/story/"+url.PathEscape(id)+"/comments"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetStoryStatistics returns statistical data for a story.
func (c *Client) GetStoryStatistics(ctx context.Context, id string) (*StoryStatistics, error) {
	var s StoryStatistics
	if err := c.get(ctx, "/story/"+url.PathEscape(id)+"/statistics", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStoryRelationship returns the authenticated member's relationship with a story.
func (c *Client) GetStoryRelationship(ctx context.Context, id string) (*StoryRelationship, error) {
	var r StoryRelationship
	if err := c.get(ctx, "/story/"+url.PathEscape(id)+"/me", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateStoryRelationship updates the authenticated member's relationship with a story (like/subscribe).
func (c *Client) UpdateStoryRelationship(ctx context.Context, id string, req StoryRelationshipUpdateRequest) (*StoryRelationshipUpdateResponse, error) {
	var resp StoryRelationshipUpdateResponse
	if err := c.patch(ctx, "/story/"+url.PathEscape(id)+"/me", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateStoryComment creates a comment on a story.
func (c *Client) CreateStoryComment(ctx context.Context, id string, req CommentCreationRequest) (*Comment, error) {
	var resp Comment
	if err := c.post(ctx, "/story/"+url.PathEscape(id)+"/comments", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search.

// Search performs a universal search across films, members, lists, contributors, reviews, and more.
func (c *Client) Search(ctx context.Context, query, cursor string, perPage int) (*SearchResponse, error) {
	v := url.Values{}
	v.Set("input", query)
	if cursor != "" {
		v.Set("cursor", cursor)
	}
	if perPage > 0 {
		v.Set("perPage", strconv.Itoa(perPage))
	}
	var r SearchResponse
	if err := c.get(ctx, "/search?"+v.Encode(), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// News.

// GetNews returns a cursored list of recent news items.
func (c *Client) GetNews(ctx context.Context, cursor string, perPage int) (*NewsResponse, error) {
	var r NewsResponse
	if err := c.get(ctx, "/news"+cursorParams(cursor, perPage), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Auth helpers.

// CheckUsername checks if a username is available.
func (c *Client) CheckUsername(ctx context.Context, username string) (*UsernameCheckResponse, error) {
	v := url.Values{}
	v.Set("username", username)
	var r UsernameCheckResponse
	if err := c.get(ctx, "/auth/username-check?"+v.Encode(), &r); err != nil {
		return nil, err
	}
	return &r, nil
}
