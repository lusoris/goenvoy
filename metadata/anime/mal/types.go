package mal

import "time"

// Anime represents a MyAnimeList anime entry.
type Anime struct {
	ID                     int               `json:"id"`
	Title                  string            `json:"title"`
	MainPicture            Picture           `json:"main_picture"`
	AlternativeTitles      AlternativeTitles `json:"alternative_titles"`
	StartDate              string            `json:"start_date"`
	EndDate                string            `json:"end_date"`
	Synopsis               string            `json:"synopsis"`
	Mean                   float64           `json:"mean"`
	Rank                   int               `json:"rank"`
	Popularity             int               `json:"popularity"`
	NumListUsers           int               `json:"num_list_users"`
	NumScoringUsers        int               `json:"num_scoring_users"`
	NSFW                   string            `json:"nsfw"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
	MediaType              string            `json:"media_type"`
	Status                 string            `json:"status"`
	Genres                 []Genre           `json:"genres"`
	NumEpisodes            int               `json:"num_episodes"`
	StartSeason            *StartSeason      `json:"start_season"`
	Broadcast              *Broadcast        `json:"broadcast"`
	Source                 string            `json:"source"`
	AverageEpisodeDuration int               `json:"average_episode_duration"`
	Rating                 string            `json:"rating"`
	Pictures               []Picture         `json:"pictures"`
	Background             string            `json:"background"`
	RelatedAnime           []RelatedAnime    `json:"related_anime"`
	RelatedManga           []RelatedManga    `json:"related_manga"`
	Recommendations        []Recommendation  `json:"recommendations"`
	Studios                []Studio          `json:"studios"`
	Statistics             *AnimeStatistics  `json:"statistics"`
}

// Manga represents a MyAnimeList manga entry.
type Manga struct {
	ID                int               `json:"id"`
	Title             string            `json:"title"`
	MainPicture       Picture           `json:"main_picture"`
	AlternativeTitles AlternativeTitles `json:"alternative_titles"`
	StartDate         string            `json:"start_date"`
	EndDate           string            `json:"end_date"`
	Synopsis          string            `json:"synopsis"`
	Mean              float64           `json:"mean"`
	Rank              int               `json:"rank"`
	Popularity        int               `json:"popularity"`
	NumListUsers      int               `json:"num_list_users"`
	NumScoringUsers   int               `json:"num_scoring_users"`
	NSFW              string            `json:"nsfw"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	MediaType         string            `json:"media_type"`
	Status            string            `json:"status"`
	Genres            []Genre           `json:"genres"`
	NumVolumes        int               `json:"num_volumes"`
	NumChapters       int               `json:"num_chapters"`
	Authors           []MangaAuthor     `json:"authors"`
	Pictures          []Picture         `json:"pictures"`
	Background        string            `json:"background"`
	RelatedAnime      []RelatedAnime    `json:"related_anime"`
	RelatedManga      []RelatedManga    `json:"related_manga"`
	Recommendations   []Recommendation  `json:"recommendations"`
	Serialization     []Serialization   `json:"serialization"`
}

// Picture holds medium and large image URLs.
type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// AlternativeTitles holds synonym and localized titles.
type AlternativeTitles struct {
	Synonyms []string `json:"synonyms"`
	En       string   `json:"en"`
	Ja       string   `json:"ja"`
}

// Genre is a content genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Studio is an animation studio.
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StartSeason is the premiere season.
type StartSeason struct {
	Year   int    `json:"year"`
	Season string `json:"season"`
}

// Broadcast is the day and time of broadcast.
type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week"`
	StartTime    string `json:"start_time"`
}

// RelatedAnime is an anime related to the queried entry.
type RelatedAnime struct {
	Node                  Anime  `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

// RelatedManga is a manga related to the queried entry.
type RelatedManga struct {
	Node                  Manga  `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

// Recommendation is a recommended anime or manga entry.
type Recommendation struct {
	Node               Anime `json:"node"`
	NumRecommendations int   `json:"num_recommendations"`
}

// AnimeStatistics contains aggregate list counts for an anime.
type AnimeStatistics struct {
	Status       AnimeStatusCounts `json:"status"`
	NumListUsers int               `json:"num_list_users"`
}

// AnimeStatusCounts contains per-status user counts.
type AnimeStatusCounts struct {
	Watching    string `json:"watching"`
	Completed   string `json:"completed"`
	OnHold      string `json:"on_hold"`
	Dropped     string `json:"dropped"`
	PlanToWatch string `json:"plan_to_watch"`
}

// MangaAuthor is an author of a manga.
type MangaAuthor struct {
	Node Person `json:"node"`
	Role string `json:"role"`
}

// Person represents a person (manga author, etc.).
type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Magazine is a publication outlet.
type Magazine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Serialization is a serialization entry for a manga.
type Serialization struct {
	Node Magazine `json:"node"`
	Role string   `json:"role"`
}

// Ranking holds the ranking position in a ranked list.
type Ranking struct {
	Rank int `json:"rank"`
}

// Paging contains pagination URLs.
type Paging struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

// ForumCategory is a top-level forum category.
type ForumCategory struct {
	Title  string       `json:"title"`
	Boards []ForumBoard `json:"boards"`
}

// ForumBoard is a forum board.
type ForumBoard struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Subboards   []ForumSubboard `json:"subboards"`
}

// ForumSubboard is a forum sub-board.
type ForumSubboard struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// ForumTopic is a forum topic summary.
type ForumTopic struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         ForumUser `json:"created_by"`
	NumberOfPosts     int       `json:"number_of_posts"`
	LastPostCreatedAt time.Time `json:"last_post_created_at"`
	LastPostCreatedBy ForumUser `json:"last_post_created_by"`
	IsLocked          bool      `json:"is_locked"`
}

// ForumTopicDetail holds posts and an optional poll.
type ForumTopicDetail struct {
	Title string      `json:"title"`
	Posts []ForumPost `json:"posts"`
	Poll  *ForumPoll  `json:"poll"`
}

// ForumPost is a single forum post.
type ForumPost struct {
	ID        int       `json:"id"`
	Number    int       `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy ForumUser `json:"created_by"`
	Body      string    `json:"body"`
	Signature string    `json:"signature"`
}

// ForumUser is a user reference in forum data.
type ForumUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ForumPoll is an optional poll attached to a topic.
type ForumPoll struct {
	ID       int               `json:"id"`
	Question string            `json:"question"`
	Closed   bool              `json:"closed"`
	Options  []ForumPollOption `json:"options"`
}

// ForumPollOption is a poll choice.
type ForumPollOption struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

// animeListItem wraps an anime node inside a paginated list.
type animeListItem struct {
	Node Anime `json:"node"`
}

// animeListResponse is the top-level response for anime list/search endpoints.
type animeListResponse struct {
	Data   []animeListItem `json:"data"`
	Paging Paging          `json:"paging"`
}

// animeRankingItem wraps an anime node with ranking info.
type animeRankingItem struct {
	Node    Anime   `json:"node"`
	Ranking Ranking `json:"ranking"`
}

// animeRankingResponse is the top-level response for anime ranking endpoints.
type animeRankingResponse struct {
	Data   []animeRankingItem `json:"data"`
	Paging Paging             `json:"paging"`
}

// mangaListItem wraps a manga node inside a paginated list.
type mangaListItem struct {
	Node Manga `json:"node"`
}

// mangaListResponse is the top-level response for manga list/search endpoints.
type mangaListResponse struct {
	Data   []mangaListItem `json:"data"`
	Paging Paging          `json:"paging"`
}

// mangaRankingItem wraps a manga node with ranking info.
type mangaRankingItem struct {
	Node    Manga   `json:"node"`
	Ranking Ranking `json:"ranking"`
}

// mangaRankingResponse is the top-level response for manga ranking endpoints.
type mangaRankingResponse struct {
	Data   []mangaRankingItem `json:"data"`
	Paging Paging             `json:"paging"`
}

// forumBoardsResponse wraps the top-level forum boards response.
type forumBoardsResponse struct {
	Categories []ForumCategory `json:"categories"`
}

// forumTopicsResponse wraps the paginated forum topics response.
type forumTopicsResponse struct {
	Data   []ForumTopic `json:"data"`
	Paging Paging       `json:"paging"`
}

// forumTopicDetailResponse wraps the paginated topic detail response.
type forumTopicDetailResponse struct {
	Data   ForumTopicDetail `json:"data"`
	Paging Paging           `json:"paging"`
}
