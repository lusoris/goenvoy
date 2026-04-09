package letterboxd

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error returned by the Letterboxd API.
type APIError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type,omitempty"`
	Message    string `json:"message,omitempty"`
	RawBody    string `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("letterboxd: HTTP %d: %s", e.StatusCode, e.Message)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("letterboxd: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("letterboxd: HTTP %d", e.StatusCode)
}

// Image represents a Letterboxd image with multiple sizes.
type Image struct {
	Sizes []ImageSize `json:"sizes,omitempty"`
}

// ImageSize represents an image at a specific resolution.
type ImageSize struct {
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	URL    string `json:"url,omitempty"`
}

// Link represents an external link.
type Link struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Genre represents a film genre.
type Genre struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Country represents a country.
type Country struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// Language represents a language.
type Language struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// FilmService represents a streaming service.
type FilmService struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Tag represents a user-created tag.
type Tag struct {
	Code       string `json:"code,omitempty"`
	DisplayTag string `json:"displayTag,omitempty"`
	Count      int    `json:"count,omitempty"`
}

// Films.

// Film represents full details for a Letterboxd film.
type Film struct {
	ID                 string         `json:"id,omitempty"`
	Name               string         `json:"name,omitempty"`
	OriginalName       string         `json:"originalName,omitempty"`
	AlternativeNames   []string       `json:"alternativeNames,omitempty"`
	ReleaseYear        int            `json:"releaseYear,omitempty"`
	Tagline            string         `json:"tagline,omitempty"`
	Description        string         `json:"description,omitempty"`
	RunTime            int            `json:"runTime,omitempty"`
	Poster             *Image         `json:"poster,omitempty"`
	Backdrop           *Image         `json:"backdrop,omitempty"`
	BackdropFocalPoint *FocalPoint    `json:"backdropFocalPoint,omitempty"`
	Trailer            *Trailer       `json:"trailer,omitempty"`
	Genres             []Genre        `json:"genres,omitempty"`
	Countries          []Country      `json:"countries,omitempty"`
	Languages          []Language     `json:"languages,omitempty"`
	Contributions      []Contribution `json:"contributions,omitempty"`
	Links              []Link         `json:"links,omitempty"`
	FilmCollectionID   string         `json:"filmCollectionId,omitempty"`
	Adult              bool           `json:"adult,omitempty"`
}

// FocalPoint represents a focal point for backdrop cropping.
type FocalPoint struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

// Trailer represents a film trailer.
type Trailer struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Contribution represents a person's contribution to a film.
type Contribution struct {
	Type         string            `json:"type,omitempty"`
	Contributors []ContributorInfo `json:"contributors,omitempty"`
}

// ContributorInfo represents a contributor within a film's credits.
type ContributorInfo struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	CharacterName string `json:"characterName,omitempty"`
}

// FilmSummary represents a compact film summary.
type FilmSummary struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	OriginalName     string   `json:"originalName,omitempty"`
	AlternativeNames []string `json:"alternativeNames,omitempty"`
	ReleaseYear      int      `json:"releaseYear,omitempty"`
	Poster           *Image   `json:"poster,omitempty"`
	Adult            bool     `json:"adult,omitempty"`
	Links            []Link   `json:"links,omitempty"`
}

// FilmStatistics represents statistical data for a film.
type FilmStatistics struct {
	Film    *FilmIdentifier    `json:"film,omitempty"`
	Counts  *FilmStatCounts    `json:"counts,omitempty"`
	Rating  float64            `json:"rating,omitempty"`
	Ratings []RatingsHistogram `json:"ratingsHistogram,omitempty"`
}

// FilmIdentifier references a film by ID.
type FilmIdentifier struct {
	ID string `json:"id,omitempty"`
}

// FilmStatCounts holds numeric counts for a film.
type FilmStatCounts struct {
	Watches int `json:"watches,omitempty"`
	Likes   int `json:"likes,omitempty"`
	Ratings int `json:"ratings,omitempty"`
	Fans    int `json:"fans,omitempty"`
	Lists   int `json:"lists,omitempty"`
	Reviews int `json:"reviews,omitempty"`
}

// RatingsHistogram represents a rating count at a specific value.
type RatingsHistogram struct {
	Rating float64 `json:"rating,omitempty"`
	Count  int     `json:"count,omitempty"`
}

// FilmRelationship represents a member's relationship with a film.
type FilmRelationship struct {
	Watched      bool     `json:"watched,omitempty"`
	Liked        bool     `json:"liked,omitempty"`
	InWatchlist  bool     `json:"inWatchlist,omitempty"`
	Rating       float64  `json:"rating,omitempty"`
	Reviews      []string `json:"reviews,omitempty"`
	DiaryEntries []string `json:"diaryEntries,omitempty"`
}

// FilmRelationshipUpdateRequest is the request body for updating a film relationship.
type FilmRelationshipUpdateRequest struct {
	Watched     *bool    `json:"watched,omitempty"`
	Liked       *bool    `json:"liked,omitempty"`
	InWatchlist *bool    `json:"inWatchlist,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
}

// FilmRelationshipUpdateResponse is the response after updating a film relationship.
type FilmRelationshipUpdateResponse struct {
	Data     *FilmRelationship `json:"data,omitempty"`
	Messages []string          `json:"messages,omitempty"`
}

// FilmMemberRelationshipItem represents a member's relationship with a film in a list.
type FilmMemberRelationshipItem struct {
	Member       *MemberSummary    `json:"member,omitempty"`
	Relationship *FilmRelationship `json:"relationship,omitempty"`
}

// FilmMembersResponse is the response for film member relationships.
type FilmMembersResponse struct {
	Cursor string                       `json:"cursor,omitempty"`
	Items  []FilmMemberRelationshipItem `json:"items,omitempty"`
}

// FilmFriendsResponse is the response for friends' relationships with a film.
type FilmFriendsResponse struct {
	Items []FilmMemberRelationshipItem `json:"items,omitempty"`
}

// FilmsResponse is the response for cursored film lists.
type FilmsResponse struct {
	Cursor string        `json:"cursor,omitempty"`
	Items  []FilmSummary `json:"items,omitempty"`
}

// GenresResponse is the response for genres listing.
type GenresResponse struct {
	Items []Genre `json:"items,omitempty"`
}

// FilmServicesResponse is the response for streaming services listing.
type FilmServicesResponse struct {
	Items []FilmService `json:"items,omitempty"`
}

// CountriesResponse is the response for countries listing.
type CountriesResponse struct {
	Items []Country `json:"items,omitempty"`
}

// LanguagesResponse is the response for languages listing.
type LanguagesResponse struct {
	Items []Language `json:"items,omitempty"`
}

// ReportRequest is the request body for reporting content.
type ReportRequest struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// Film collections.

// FilmCollection represents a collection of related films.
type FilmCollection struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Films       []FilmSummary `json:"films,omitempty"`
	Links       []Link        `json:"links,omitempty"`
	Description string        `json:"description,omitempty"`
}

// FilmCollectionSummary represents a compact film collection.
type FilmCollectionSummary struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	FilmCount int    `json:"filmCount,omitempty"`
}

// FilmCollectionsResponse is the response for cursored film collection lists.
type FilmCollectionsResponse struct {
	Cursor string                  `json:"cursor,omitempty"`
	Items  []FilmCollectionSummary `json:"items,omitempty"`
}

// Contributors.

// Contributor represents a film contributor (actor, director, etc.).
type Contributor struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Statistics *ContributorStatistics `json:"statistics,omitempty"`
	Links      []Link                 `json:"links,omitempty"`
}

// ContributorStatistics holds statistics for a contributor.
type ContributorStatistics struct {
	Contributions int `json:"contributions,omitempty"`
}

// ContributionItem represents a film in a contributor's filmography.
type ContributionItem struct {
	Type          string       `json:"type,omitempty"`
	Film          *FilmSummary `json:"film,omitempty"`
	CharacterName string       `json:"characterName,omitempty"`
}

// ContributionsResponse is the response for contributor filmography.
type ContributionsResponse struct {
	Cursor string             `json:"cursor,omitempty"`
	Items  []ContributionItem `json:"items,omitempty"`
}

// Lists.

// List represents full details for a Letterboxd list.
type List struct {
	ID                  string         `json:"id,omitempty"`
	Name                string         `json:"name,omitempty"`
	Description         string         `json:"description,omitempty"`
	DescriptionLBML     string         `json:"descriptionLbml,omitempty"`
	Tags                []Tag          `json:"tags,omitempty"`
	Published           bool           `json:"published,omitempty"`
	Ranked              bool           `json:"ranked,omitempty"`
	HasEntriesWithNotes bool           `json:"hasEntriesWithNotes,omitempty"`
	FilmCount           int            `json:"filmCount,omitempty"`
	Owner               *MemberSummary `json:"owner,omitempty"`
	Cloneable           bool           `json:"cloneable,omitempty"`
	WhenCreated         string         `json:"whenCreated,omitempty"`
	WhenPublished       string         `json:"whenPublished,omitempty"`
	PreviewEntries      []ListEntry    `json:"previewEntries,omitempty"`
	Links               []Link         `json:"links,omitempty"`
	BackdropImage       *Image         `json:"backdrop,omitempty"`
}

// ListEntry represents an entry (film) in a list.
type ListEntry struct {
	Rank      int          `json:"rank,omitempty"`
	Film      *FilmSummary `json:"film,omitempty"`
	Notes     string       `json:"notes,omitempty"`
	NotesLBML string       `json:"notesLbml,omitempty"`
	EntryID   string       `json:"entryId,omitempty"`
}

// ListEntriesResponse is the response for cursored list entries.
type ListEntriesResponse struct {
	Cursor string      `json:"cursor,omitempty"`
	Items  []ListEntry `json:"items,omitempty"`
}

// ListStatistics represents statistical data for a list.
type ListStatistics struct {
	List   *ListIdentifier `json:"list,omitempty"`
	Counts *ListStatCounts `json:"counts,omitempty"`
}

// ListIdentifier references a list by ID.
type ListIdentifier struct {
	ID string `json:"id,omitempty"`
}

// ListStatCounts holds numeric counts for a list.
type ListStatCounts struct {
	Comments int `json:"comments,omitempty"`
	Likes    int `json:"likes,omitempty"`
}

// ListRelationship represents a member's relationship with a list.
type ListRelationship struct {
	Liked      bool `json:"liked,omitempty"`
	Subscribed bool `json:"subscribed,omitempty"`
}

// ListRelationshipUpdateRequest is the request body for updating a list relationship.
type ListRelationshipUpdateRequest struct {
	Liked      *bool `json:"liked,omitempty"`
	Subscribed *bool `json:"subscribed,omitempty"`
}

// ListRelationshipUpdateResponse is the response after updating a list relationship.
type ListRelationshipUpdateResponse struct {
	Data     *ListRelationship `json:"data,omitempty"`
	Messages []string          `json:"messages,omitempty"`
}

// ListSummary represents a compact list summary.
type ListSummary struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	FilmCount      int            `json:"filmCount,omitempty"`
	Published      bool           `json:"published,omitempty"`
	Ranked         bool           `json:"ranked,omitempty"`
	Owner          *MemberSummary `json:"owner,omitempty"`
	PreviewEntries []ListEntry    `json:"previewEntries,omitempty"`
	Description    string         `json:"description,omitempty"`
}

// ListsResponse is the response for cursored list listings.
type ListsResponse struct {
	Cursor string        `json:"cursor,omitempty"`
	Items  []ListSummary `json:"items,omitempty"`
}

// ListTopic represents a featured list topic.
type ListTopic struct {
	Name  string        `json:"name,omitempty"`
	Lists []ListSummary `json:"items,omitempty"`
}

// ListTopicsResponse is the response for list topics.
type ListTopicsResponse struct {
	Items []ListTopic `json:"items,omitempty"`
}

// ListCreationRequest is the request body for creating a new list.
type ListCreationRequest struct {
	Name            string             `json:"name"`
	Description     string             `json:"description,omitempty"`
	DescriptionLBML string             `json:"descriptionLbml,omitempty"`
	Ranked          bool               `json:"ranked,omitempty"`
	Published       bool               `json:"published,omitempty"`
	Tags            []string           `json:"tags,omitempty"`
	Entries         []ListEntryRequest `json:"entries,omitempty"`
}

// ListEntryRequest represents an entry to add to a list.
type ListEntryRequest struct {
	Film      string `json:"film"`
	Rank      int    `json:"rank,omitempty"`
	Notes     string `json:"notes,omitempty"`
	NotesLBML string `json:"notesLbml,omitempty"`
}

// ListCreateResponse is the response after creating a list.
type ListCreateResponse struct {
	Data     *List    `json:"data,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

// ListUpdateRequest is the request body for updating a list.
type ListUpdateRequest struct {
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description,omitempty"`
	DescriptionLBML string             `json:"descriptionLbml,omitempty"`
	Ranked          *bool              `json:"ranked,omitempty"`
	Published       *bool              `json:"published,omitempty"`
	Tags            []string           `json:"tags,omitempty"`
	Entries         []ListEntryRequest `json:"entries,omitempty"`
}

// ListUpdateResponse is the response after updating a list.
type ListUpdateResponse struct {
	Data     *List    `json:"data,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

// ListAddEntriesRequest is the request body for adding films to lists.
type ListAddEntriesRequest struct {
	Lists []string           `json:"lists,omitempty"`
	Films []ListEntryRequest `json:"films,omitempty"`
}

// Comments.

// Comment represents a comment on a list, review, or story.
type Comment struct {
	ID                    string         `json:"id,omitempty"`
	Member                *MemberSummary `json:"member,omitempty"`
	Comment               string         `json:"comment,omitempty"`
	CommentLBML           string         `json:"commentLbml,omitempty"`
	WhenCreated           string         `json:"whenCreated,omitempty"`
	WhenUpdated           string         `json:"whenUpdated,omitempty"`
	Edited                bool           `json:"edited,omitempty"`
	Deleted               bool           `json:"deleted,omitempty"`
	Blocked               bool           `json:"blocked,omitempty"`
	RemovedByAdmin        bool           `json:"removedByAdmin,omitempty"`
	RemovedByContentOwner bool           `json:"removedByContentOwner,omitempty"`
}

// CommentsResponse is the response for cursored comment lists.
type CommentsResponse struct {
	Cursor string    `json:"cursor,omitempty"`
	Items  []Comment `json:"items,omitempty"`
}

// CommentCreationRequest is the request body for creating a comment.
type CommentCreationRequest struct {
	Comment string `json:"comment"`
}

// CommentUpdateRequest is the request body for updating a comment.
type CommentUpdateRequest struct {
	Comment string `json:"comment"`
}

// Log entries.

// LogEntry represents a diary entry or review.
type LogEntry struct {
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Film          *FilmSummary   `json:"film,omitempty"`
	Owner         *MemberSummary `json:"owner,omitempty"`
	DiaryDetails  *DiaryDetails  `json:"diaryDetails,omitempty"`
	Review        *Review        `json:"review,omitempty"`
	Tags          []Tag          `json:"tags2,omitempty"`
	WhenCreated   string         `json:"whenCreated,omitempty"`
	WhenUpdated   string         `json:"whenUpdated,omitempty"`
	Rating        float64        `json:"rating,omitempty"`
	Like          bool           `json:"like,omitempty"`
	Commentable   bool           `json:"commentable,omitempty"`
	Links         []Link         `json:"links,omitempty"`
	BackdropImage *Image         `json:"backdrop,omitempty"`
	Adult         bool           `json:"adult,omitempty"`
}

// DiaryDetails represents diary-specific fields of a log entry.
type DiaryDetails struct {
	DiaryDate string `json:"diaryDate,omitempty"`
	Rewatch   bool   `json:"rewatch,omitempty"`
}

// Review represents review text in a log entry.
type Review struct {
	Text     string `json:"text,omitempty"`
	TextLBML string `json:"lbml,omitempty"`
	Spoilers bool   `json:"containsSpoilers,omitempty"`
}

// LogEntriesResponse is the response for cursored log entry lists.
type LogEntriesResponse struct {
	Cursor string     `json:"cursor,omitempty"`
	Items  []LogEntry `json:"items,omitempty"`
}

// LogEntryStatistics represents statistical data for a log entry.
type LogEntryStatistics struct {
	LogEntry *LogEntryIdentifier `json:"logEntry,omitempty"`
	Counts   *LogEntryStatCounts `json:"counts,omitempty"`
}

// LogEntryIdentifier references a log entry by ID.
type LogEntryIdentifier struct {
	ID string `json:"id,omitempty"`
}

// LogEntryStatCounts holds numeric counts for a log entry.
type LogEntryStatCounts struct {
	Comments int `json:"comments,omitempty"`
	Likes    int `json:"likes,omitempty"`
}

// LogEntryRelationship represents a member's relationship with a log entry.
type LogEntryRelationship struct {
	Liked      bool `json:"liked,omitempty"`
	Subscribed bool `json:"subscribed,omitempty"`
}

// LogEntryRelationshipUpdateRequest is the request body for updating a log entry relationship.
type LogEntryRelationshipUpdateRequest struct {
	Liked      *bool `json:"liked,omitempty"`
	Subscribed *bool `json:"subscribed,omitempty"`
}

// LogEntryRelationshipUpdateResponse is the response after updating a log entry relationship.
type LogEntryRelationshipUpdateResponse struct {
	Data     *LogEntryRelationship `json:"data,omitempty"`
	Messages []string              `json:"messages,omitempty"`
}

// LogEntryCreationRequest is the request body for creating a log entry.
type LogEntryCreationRequest struct {
	FilmID    string   `json:"filmId"`
	DiaryDate string   `json:"diaryDate,omitempty"`
	Rewatch   bool     `json:"rewatch,omitempty"`
	Rating    float64  `json:"rating,omitempty"`
	Like      bool     `json:"like,omitempty"`
	Review    string   `json:"review,omitempty"`
	Spoilers  bool     `json:"containsSpoilers,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

// LogEntryUpdateRequest is the request body for updating a log entry.
type LogEntryUpdateRequest struct {
	DiaryDate string   `json:"diaryDate,omitempty"`
	Rewatch   *bool    `json:"rewatch,omitempty"`
	Rating    *float64 `json:"rating,omitempty"`
	Like      *bool    `json:"like,omitempty"`
	Review    string   `json:"review,omitempty"`
	Spoilers  *bool    `json:"containsSpoilers,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

// Members.

// Member represents full details for a Letterboxd member.
type Member struct {
	ID               string        `json:"id,omitempty"`
	Username         string        `json:"username,omitempty"`
	GivenName        string        `json:"givenName,omitempty"`
	FamilyName       string        `json:"familyName,omitempty"`
	DisplayName      string        `json:"displayName,omitempty"`
	ShortName        string        `json:"shortName,omitempty"`
	Pronoun          *Pronoun      `json:"pronoun,omitempty"`
	Bio              string        `json:"bio,omitempty"`
	BioLBML          string        `json:"bioLbml,omitempty"`
	Location         string        `json:"location,omitempty"`
	Website          string        `json:"website,omitempty"`
	Backdrop         *Image        `json:"backdrop,omitempty"`
	Avatar           *Image        `json:"avatar,omitempty"`
	MemberStatus     string        `json:"memberStatus,omitempty"`
	HideAdsInContent bool          `json:"hideAdsInContent,omitempty"`
	AccountStatus    string        `json:"accountStatus,omitempty"`
	FavoriteFilms    []FilmSummary `json:"favoriteFilms,omitempty"`
	Links            []Link        `json:"links,omitempty"`
	TwitterUsername  string        `json:"twitterUsername,omitempty"`
}

// MemberSummary represents a compact member summary.
type MemberSummary struct {
	ID          string `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	ShortName   string `json:"shortName,omitempty"`
	Avatar      *Image `json:"avatar,omitempty"`
}

// MemberStatistics represents statistical data for a member.
type MemberStatistics struct {
	Member           *MemberIdentifier  `json:"member,omitempty"`
	Counts           *MemberStatCounts  `json:"counts,omitempty"`
	RatingsHistogram []RatingsHistogram `json:"ratingsHistogram,omitempty"`
	YearsInReview    []YearInReview     `json:"yearsInReview,omitempty"`
}

// MemberIdentifier references a member by ID.
type MemberIdentifier struct {
	ID string `json:"id,omitempty"`
}

// MemberStatCounts holds numeric counts for a member.
type MemberStatCounts struct {
	FilmLikes            int `json:"filmLikes,omitempty"`
	ListLikes            int `json:"listLikes,omitempty"`
	ReviewLikes          int `json:"reviewLikes,omitempty"`
	Watches              int `json:"watches,omitempty"`
	Ratings              int `json:"ratings,omitempty"`
	Reviews              int `json:"reviews,omitempty"`
	DiaryEntries         int `json:"diaryEntries,omitempty"`
	DiaryEntriesThisYear int `json:"diaryEntriesThisYear,omitempty"`
	FilmsInDiaryThisYear int `json:"filmsInDiaryThisYear,omitempty"`
	Watchlist            int `json:"watchlist,omitempty"`
	Lists                int `json:"lists,omitempty"`
	UnpublishedLists     int `json:"unpublishedLists,omitempty"`
	Followers            int `json:"followers,omitempty"`
	Following            int `json:"following,omitempty"`
	Blocked              int `json:"blocked,omitempty"`
}

// YearInReview represents a member's year-in-review summary.
type YearInReview struct {
	Year      int `json:"year,omitempty"`
	FilmCount int `json:"filmCount,omitempty"`
}

// MemberRelationship represents a member's relationship with another member.
type MemberRelationship struct {
	Following  bool `json:"following,omitempty"`
	FollowedBy bool `json:"followedBy,omitempty"`
	Blocking   bool `json:"blocking,omitempty"`
}

// MemberRelationshipUpdateRequest is the request body for updating a member relationship.
type MemberRelationshipUpdateRequest struct {
	Following *bool `json:"following,omitempty"`
	Blocking  *bool `json:"blocking,omitempty"`
}

// MemberRelationshipUpdateResponse is the response after updating a member relationship.
type MemberRelationshipUpdateResponse struct {
	Data     *MemberRelationship `json:"data,omitempty"`
	Messages []string            `json:"messages,omitempty"`
}

// MembersResponse is the response for cursored member lists.
type MembersResponse struct {
	Cursor string          `json:"cursor,omitempty"`
	Items  []MemberSummary `json:"items,omitempty"`
}

// Pronoun represents a pronoun option.
type Pronoun struct {
	ID                string `json:"id,omitempty"`
	Label             string `json:"label,omitempty"`
	SubjectPronoun    string `json:"subjectPronoun,omitempty"`
	ObjectPronoun     string `json:"objectPronoun,omitempty"`
	PossessivePronoun string `json:"possessivePronoun,omitempty"`
	Reflexive         string `json:"reflexive,omitempty"`
}

// PronounsResponse is the response for pronouns listing.
type PronounsResponse struct {
	Items []Pronoun `json:"items,omitempty"`
}

// MemberAccount represents the authenticated member's account details.
type MemberAccount struct {
	Member                 *Member `json:"member,omitempty"`
	EmailAddress           string  `json:"emailAddress,omitempty"`
	EmailAddressValidated  bool    `json:"emailAddressValidated,omitempty"`
	PrivateAccount         bool    `json:"privateAccount,omitempty"`
	IncludeInPeopleSection bool    `json:"includeInPeopleSection,omitempty"`
	CanComment             bool    `json:"canComment,omitempty"`
	CanCloneLists          bool    `json:"canCloneLists,omitempty"`
	SuspendedComment       bool    `json:"suspendedComment,omitempty"`
	SuspendedLists         bool    `json:"suspendedLists,omitempty"`
	SuspendedReviews       bool    `json:"suspendedReviews,omitempty"`
}

// MemberSettingsUpdateRequest is the request body for updating account settings.
type MemberSettingsUpdateRequest struct {
	Bio           string   `json:"bio,omitempty"`
	BioLBML       string   `json:"bioLbml,omitempty"`
	Location      string   `json:"location,omitempty"`
	Website       string   `json:"website,omitempty"`
	GivenName     string   `json:"givenName,omitempty"`
	FamilyName    string   `json:"familyName,omitempty"`
	FavoriteFilms []string `json:"favoriteFilms,omitempty"`
}

// MemberSettingsUpdateResponse is the response after updating account settings.
type MemberSettingsUpdateResponse struct {
	Data     *MemberAccount `json:"data,omitempty"`
	Messages []string       `json:"messages,omitempty"`
}

// Activity.

// ActivityItem represents an item in a member's activity feed.
type ActivityItem struct {
	Type        string          `json:"type,omitempty"`
	Member      *MemberSummary  `json:"member,omitempty"`
	Film        *FilmSummary    `json:"film,omitempty"`
	List        *ListSummary    `json:"list,omitempty"`
	LogEntry    *LogEntry       `json:"logEntry,omitempty"`
	Review      *LogEntry       `json:"review,omitempty"`
	Story       *StorySummary   `json:"story,omitempty"`
	WhenCreated string          `json:"whenCreated,omitempty"`
	RawData     json.RawMessage `json:"-"`
}

// ActivityResponse is the response for a member's activity feed.
type ActivityResponse struct {
	Cursor string         `json:"cursor,omitempty"`
	Items  []ActivityItem `json:"items,omitempty"`
}

// Tags.

// TagsResponse is the response for tag listings.
type TagsResponse struct {
	Items []Tag `json:"items,omitempty"`
}

// Stories.

// Story represents full details for a Letterboxd story.
type Story struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Author      *MemberSummary `json:"author,omitempty"`
	Body        string         `json:"body,omitempty"`
	BodyLBML    string         `json:"bodyLbml,omitempty"`
	WhenCreated string         `json:"whenCreated,omitempty"`
	WhenUpdated string         `json:"whenUpdated,omitempty"`
	Image       *Image         `json:"image,omitempty"`
	VideoURL    string         `json:"videoUrl,omitempty"`
	Source      string         `json:"source,omitempty"`
	Links       []Link         `json:"links,omitempty"`
}

// StorySummary represents a compact story summary.
type StorySummary struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Author      *MemberSummary `json:"author,omitempty"`
	WhenCreated string         `json:"whenCreated,omitempty"`
	Image       *Image         `json:"image,omitempty"`
}

// StoryStatistics represents statistical data for a story.
type StoryStatistics struct {
	Story  *StoryIdentifier `json:"story,omitempty"`
	Counts *StoryStatCounts `json:"counts,omitempty"`
}

// StoryIdentifier references a story by ID.
type StoryIdentifier struct {
	ID string `json:"id,omitempty"`
}

// StoryStatCounts holds numeric counts for a story.
type StoryStatCounts struct {
	Comments int `json:"comments,omitempty"`
	Likes    int `json:"likes,omitempty"`
}

// StoryRelationship represents a member's relationship with a story.
type StoryRelationship struct {
	Liked      bool `json:"liked,omitempty"`
	Subscribed bool `json:"subscribed,omitempty"`
}

// StoryRelationshipUpdateRequest is the request body for updating a story relationship.
type StoryRelationshipUpdateRequest struct {
	Liked      *bool `json:"liked,omitempty"`
	Subscribed *bool `json:"subscribed,omitempty"`
}

// StoryRelationshipUpdateResponse is the response after updating a story relationship.
type StoryRelationshipUpdateResponse struct {
	Data     *StoryRelationship `json:"data,omitempty"`
	Messages []string           `json:"messages,omitempty"`
}

// StoriesResponse is the response for cursored story lists.
type StoriesResponse struct {
	Cursor string         `json:"cursor,omitempty"`
	Items  []StorySummary `json:"items,omitempty"`
}

// Search.

// SearchItem represents an item in search results.
type SearchItem struct {
	Type        string           `json:"type,omitempty"`
	Score       float64          `json:"score,omitempty"`
	Film        *FilmSummary     `json:"film,omitempty"`
	List        *ListSummary     `json:"listSummary,omitempty"`
	LogEntry    *LogEntry        `json:"reviewSummary,omitempty"`
	Contributor *ContributorInfo `json:"contributorSummary,omitempty"`
	Member      *MemberSummary   `json:"memberSummary,omitempty"`
	Story       *StorySummary    `json:"storySummary,omitempty"`
	Tag         *Tag             `json:"tag,omitempty"`
}

// SearchResponse is the response for search queries.
type SearchResponse struct {
	Cursor       string       `json:"cursor,omitempty"`
	Items        []SearchItem `json:"items,omitempty"`
	SearchMethod string       `json:"searchMethod,omitempty"`
}

// News.

// NewsItem represents a news item.
type NewsItem struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	URL         string `json:"url,omitempty"`
	Image       *Image `json:"image,omitempty"`
	WhenCreated string `json:"whenCreated,omitempty"`
}

// NewsResponse is the response for cursored news items.
type NewsResponse struct {
	Cursor string     `json:"cursor,omitempty"`
	Items  []NewsItem `json:"items,omitempty"`
}

// Auth.

// UsernameCheckResponse is the response for the username availability check.
type UsernameCheckResponse struct {
	Result string `json:"result,omitempty"`
}
