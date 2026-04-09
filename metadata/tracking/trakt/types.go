package trakt

// IDs contains external identifiers for a media item.
type IDs struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	IMDb  string `json:"imdb,omitempty"`
	TMDb  int    `json:"tmdb,omitempty"`
	TVDb  int    `json:"tvdb,omitempty"`
}

// Movie represents a movie object returned by the API.
type Movie struct {
	Title                 string   `json:"title"`
	Year                  int      `json:"year"`
	IDs                   IDs      `json:"ids"`
	Tagline               string   `json:"tagline,omitempty"`
	Overview              string   `json:"overview,omitempty"`
	Released              string   `json:"released,omitempty"`
	Runtime               int      `json:"runtime,omitempty"`
	Country               string   `json:"country,omitempty"`
	Trailer               string   `json:"trailer,omitempty"`
	Homepage              string   `json:"homepage,omitempty"`
	Status                string   `json:"status,omitempty"`
	Rating                float64  `json:"rating,omitempty"`
	Votes                 int      `json:"votes,omitempty"`
	CommentCount          int      `json:"comment_count,omitempty"`
	Language              string   `json:"language,omitempty"`
	AvailableTranslations []string `json:"available_translations,omitempty"`
	Genres                []string `json:"genres,omitempty"`
	Certification         string   `json:"certification,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

// Show represents a TV show object returned by the API.
type Show struct {
	Title                 string   `json:"title"`
	Year                  int      `json:"year"`
	IDs                   IDs      `json:"ids"`
	Overview              string   `json:"overview,omitempty"`
	FirstAired            string   `json:"first_aired,omitempty"`
	Runtime               int      `json:"runtime,omitempty"`
	Certification         string   `json:"certification,omitempty"`
	Network               string   `json:"network,omitempty"`
	Country               string   `json:"country,omitempty"`
	Trailer               string   `json:"trailer,omitempty"`
	Homepage              string   `json:"homepage,omitempty"`
	Status                string   `json:"status,omitempty"`
	Rating                float64  `json:"rating,omitempty"`
	Votes                 int      `json:"votes,omitempty"`
	CommentCount          int      `json:"comment_count,omitempty"`
	Language              string   `json:"language,omitempty"`
	AvailableTranslations []string `json:"available_translations,omitempty"`
	Genres                []string `json:"genres,omitempty"`
	AiredEpisodes         int      `json:"aired_episodes,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

// Season represents a TV season.
type Season struct {
	Number        int     `json:"number"`
	IDs           IDs     `json:"ids"`
	Title         string  `json:"title,omitempty"`
	Overview      string  `json:"overview,omitempty"`
	Rating        float64 `json:"rating,omitempty"`
	Votes         int     `json:"votes,omitempty"`
	EpisodeCount  int     `json:"episode_count,omitempty"`
	AiredEpisodes int     `json:"aired_episodes,omitempty"`
	Network       string  `json:"network,omitempty"`
	FirstAired    string  `json:"first_aired,omitempty"`
}

// Episode represents a TV episode.
type Episode struct {
	Season                int      `json:"season"`
	Number                int      `json:"number"`
	Title                 string   `json:"title"`
	IDs                   IDs      `json:"ids"`
	Overview              string   `json:"overview,omitempty"`
	Rating                float64  `json:"rating,omitempty"`
	Votes                 int      `json:"votes,omitempty"`
	CommentCount          int      `json:"comment_count,omitempty"`
	FirstAired            string   `json:"first_aired,omitempty"`
	Runtime               int      `json:"runtime,omitempty"`
	AvailableTranslations []string `json:"available_translations,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

// Person represents a person (actor, director, etc.).
type Person struct {
	Name               string `json:"name"`
	IDs                IDs    `json:"ids"`
	Biography          string `json:"biography,omitempty"`
	Birthday           string `json:"birthday,omitempty"`
	Death              string `json:"death,omitempty"`
	Birthplace         string `json:"birthplace,omitempty"`
	Homepage           string `json:"homepage,omitempty"`
	Gender             string `json:"gender,omitempty"`
	KnownForDepartment string `json:"known_for_department,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
}

// TrendingMovie is a movie with its trending watcher count.
type TrendingMovie struct {
	Watchers int   `json:"watchers"`
	Movie    Movie `json:"movie"`
}

// TrendingShow is a show with its trending watcher count.
type TrendingShow struct {
	Watchers int  `json:"watchers"`
	Show     Show `json:"show"`
}

// PlayedMovie is a movie with its play/watch/collect count.
type PlayedMovie struct {
	WatcherCount   int   `json:"watcher_count"`
	PlayCount      int   `json:"play_count"`
	CollectedCount int   `json:"collected_count"`
	Movie          Movie `json:"movie"`
}

// PlayedShow is a show with its play/watch/collect count.
type PlayedShow struct {
	WatcherCount   int  `json:"watcher_count"`
	PlayCount      int  `json:"play_count"`
	CollectedCount int  `json:"collected_count"`
	Show           Show `json:"show"`
}

// AnticipatedMovie is a movie with its list count.
type AnticipatedMovie struct {
	ListCount int   `json:"list_count"`
	Movie     Movie `json:"movie"`
}

// AnticipatedShow is a show with its list count.
type AnticipatedShow struct {
	ListCount int  `json:"list_count"`
	Show      Show `json:"show"`
}

// BoxOfficeMovie is a movie with its revenue.
type BoxOfficeMovie struct {
	Revenue int   `json:"revenue"`
	Movie   Movie `json:"movie"`
}

// MovieTranslation represents a movie translation.
type MovieTranslation struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Tagline  string `json:"tagline"`
	Language string `json:"language"`
	Country  string `json:"country"`
}

// ShowTranslation represents a show translation.
type ShowTranslation struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Language string `json:"language"`
	Country  string `json:"country"`
}

// Ratings contains rating information for a media item.
type Ratings struct {
	Rating       float64      `json:"rating"`
	Votes        int          `json:"votes"`
	Distribution Distribution `json:"distribution"`
}

// Distribution maps rating values (1-10) to their counts.
type Distribution struct {
	One   int `json:"1"`
	Two   int `json:"2"`
	Three int `json:"3"`
	Four  int `json:"4"`
	Five  int `json:"5"`
	Six   int `json:"6"`
	Seven int `json:"7"`
	Eight int `json:"8"`
	Nine  int `json:"9"`
	Ten   int `json:"10"`
}

// Stats contains statistics for a media item.
type Stats struct {
	Watchers        int `json:"watchers"`
	Plays           int `json:"plays"`
	Collectors      int `json:"collectors"`
	Comments        int `json:"comments"`
	Lists           int `json:"lists"`
	Votes           int `json:"votes"`
	Favorited       int `json:"favorited"`
	Recommendations int `json:"recommendations"`
}

// CastMember represents a cast credit.
type CastMember struct {
	Characters []string `json:"characters"`
	Person     Person   `json:"person"`
}

// CrewMember represents a crew credit.
type CrewMember struct {
	Jobs   []string `json:"jobs"`
	Person Person   `json:"person"`
}

// People contains cast and crew for a media item.
type People struct {
	Cast []CastMember `json:"cast"`
	Crew *Crew        `json:"crew,omitempty"`
}

// Crew groups crew members by department.
type Crew struct {
	Production       []CrewMember `json:"production,omitempty"`
	Art              []CrewMember `json:"art,omitempty"`
	Crew             []CrewMember `json:"crew,omitempty"`
	CostumeAndMakeUp []CrewMember `json:"costume & make-up,omitempty"`
	Directing        []CrewMember `json:"directing,omitempty"`
	Writing          []CrewMember `json:"writing,omitempty"`
	Sound            []CrewMember `json:"sound,omitempty"`
	Camera           []CrewMember `json:"camera,omitempty"`
	VisualEffects    []CrewMember `json:"visual effects,omitempty"`
	Lighting         []CrewMember `json:"lighting,omitempty"`
	Editing          []CrewMember `json:"editing,omitempty"`
}

// Studio represents a production studio.
type Studio struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	IDs     IDs    `json:"ids"`
}

// SearchResult is a single result from the search endpoint.
type SearchResult struct {
	Type    string   `json:"type"`
	Score   float64  `json:"score"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
	Person  *Person  `json:"person,omitempty"`
}

// CalendarMovie is a movie in a calendar list.
type CalendarMovie struct {
	Released string `json:"released"`
	Movie    Movie  `json:"movie"`
}

// CalendarShow is a show episode in a calendar list.
type CalendarShow struct {
	FirstAired string  `json:"first_aired"`
	Episode    Episode `json:"episode"`
	Show       Show    `json:"show"`
}

// Genre represents a content genre.
type Genre struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Certification represents a content certification.
type Certification struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// Alias represents a title alias.
type Alias struct {
	Title   string `json:"title"`
	Country string `json:"country"`
}

// MovieRelease represents a movie release date and certification.
type MovieRelease struct {
	Country       string `json:"country"`
	Certification string `json:"certification"`
	ReleaseDate   string `json:"release_date"`
	ReleaseType   string `json:"release_type"`
	Note          string `json:"note,omitempty"`
}

// Country represents a country.
type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Language represents a language.
type Language struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Network represents a TV network.
type Network struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	IDs     IDs    `json:"ids"`
}

// OAuth2 types.

// DeviceCode holds the response from the device code request.
type DeviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// Token holds OAuth2 access and refresh tokens.
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}

// SyncItems is the request/response body for sync operations (watchlist, collection, history, ratings).
type SyncItems struct {
	Movies   []SyncMovie   `json:"movies,omitempty"`
	Shows    []SyncShow    `json:"shows,omitempty"`
	Episodes []SyncEpisode `json:"episodes,omitempty"`
	Seasons  []SyncSeason  `json:"seasons,omitempty"`
}

// SyncMovie is a movie reference for sync operations.
type SyncMovie struct {
	IDs       IDs    `json:"ids"`
	Title     string `json:"title,omitempty"`
	Year      int    `json:"year,omitempty"`
	Rating    int    `json:"rating,omitempty"`
	RatedAt   string `json:"rated_at,omitempty"`
	WatchedAt string `json:"watched_at,omitempty"`
}

// SyncShow is a show reference for sync operations.
type SyncShow struct {
	IDs       IDs          `json:"ids"`
	Title     string       `json:"title,omitempty"`
	Year      int          `json:"year,omitempty"`
	Rating    int          `json:"rating,omitempty"`
	RatedAt   string       `json:"rated_at,omitempty"`
	WatchedAt string       `json:"watched_at,omitempty"`
	Seasons   []SyncSeason `json:"seasons,omitempty"`
}

// SyncSeason is a season reference for sync operations.
type SyncSeason struct {
	Number   int           `json:"number"`
	Episodes []SyncEpisode `json:"episodes,omitempty"`
	Rating   int           `json:"rating,omitempty"`
	RatedAt  string        `json:"rated_at,omitempty"`
}

// SyncEpisode is an episode reference for sync operations.
type SyncEpisode struct {
	IDs       IDs    `json:"ids"`
	Rating    int    `json:"rating,omitempty"`
	RatedAt   string `json:"rated_at,omitempty"`
	WatchedAt string `json:"watched_at,omitempty"`
}

// SyncResponse is the response from sync add/remove operations.
type SyncResponse struct {
	Added    *SyncCount `json:"added,omitempty"`
	Deleted  *SyncCount `json:"deleted,omitempty"`
	Existing *SyncCount `json:"existing,omitempty"`
	NotFound *SyncItems `json:"not_found,omitempty"`
}

// SyncCount holds counts per media type from a sync operation.
type SyncCount struct {
	Movies   int `json:"movies"`
	Shows    int `json:"shows"`
	Seasons  int `json:"seasons"`
	Episodes int `json:"episodes"`
}

// WatchlistItem is an item on a user's watchlist.
type WatchlistItem struct {
	Rank     int      `json:"rank"`
	ListedAt string   `json:"listed_at"`
	Type     string   `json:"type"`
	Movie    *Movie   `json:"movie,omitempty"`
	Show     *Show    `json:"show,omitempty"`
	Episode  *Episode `json:"episode,omitempty"`
	Season   *Season  `json:"season,omitempty"`
}

// CollectionItem is a collected media item.
type CollectionItem struct {
	CollectedAt string   `json:"collected_at"`
	UpdatedAt   string   `json:"updated_at"`
	Movie       *Movie   `json:"movie,omitempty"`
	Show        *Show    `json:"show,omitempty"`
	Seasons     []Season `json:"seasons,omitempty"`
}

// HistoryItem is a watched history entry.
type HistoryItem struct {
	ID        int64    `json:"id"`
	WatchedAt string   `json:"watched_at"`
	Action    string   `json:"action"`
	Type      string   `json:"type"`
	Movie     *Movie   `json:"movie,omitempty"`
	Show      *Show    `json:"show,omitempty"`
	Episode   *Episode `json:"episode,omitempty"`
}

// RatedItem is a rated media item.
type RatedItem struct {
	RatedAt string   `json:"rated_at"`
	Rating  int      `json:"rating"`
	Type    string   `json:"type"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
	Season  *Season  `json:"season,omitempty"`
}

// UserProfile contains user profile information.
type UserProfile struct {
	Username string `json:"username"`
	Private  bool   `json:"private"`
	Name     string `json:"name"`
	VIP      bool   `json:"vip"`
	IDs      IDs    `json:"ids"`
	JoinedAt string `json:"joined_at"`
	Location string `json:"location,omitempty"`
	About    string `json:"about,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Age      int    `json:"age,omitempty"`
}

// UserStats contains viewing statistics for a user.
type UserStats struct {
	Movies   UserMovieStats   `json:"movies"`
	Shows    UserShowStats    `json:"shows"`
	Seasons  UserSeasonStats  `json:"seasons"`
	Episodes UserEpisodeStats `json:"episodes"`
	Network  UserNetworkStats `json:"network"`
	Ratings  UserRatingStats  `json:"ratings"`
}

// UserMovieStats contains movie-specific stats.
type UserMovieStats struct {
	Plays     int `json:"plays"`
	Watched   int `json:"watched"`
	Minutes   int `json:"minutes"`
	Collected int `json:"collected"`
	Ratings   int `json:"ratings"`
	Comments  int `json:"comments"`
}

// UserShowStats contains show-specific stats.
type UserShowStats struct {
	Watched   int `json:"watched"`
	Collected int `json:"collected"`
	Ratings   int `json:"ratings"`
	Comments  int `json:"comments"`
}

// UserSeasonStats contains season-specific stats.
type UserSeasonStats struct {
	Ratings  int `json:"ratings"`
	Comments int `json:"comments"`
}

// UserEpisodeStats contains episode-specific stats.
type UserEpisodeStats struct {
	Plays     int `json:"plays"`
	Watched   int `json:"watched"`
	Minutes   int `json:"minutes"`
	Collected int `json:"collected"`
	Ratings   int `json:"ratings"`
	Comments  int `json:"comments"`
}

// UserNetworkStats contains social network stats.
type UserNetworkStats struct {
	Friends   int `json:"friends"`
	Followers int `json:"followers"`
	Following int `json:"following"`
}

// UserRatingStats contains rating distribution stats.
type UserRatingStats struct {
	Total        int          `json:"total"`
	Distribution Distribution `json:"distribution"`
}

// UserList is a custom list created by a user.
type UserList struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Privacy        string `json:"privacy,omitempty"`
	DisplayNumbers bool   `json:"display_numbers,omitempty"`
	AllowComments  bool   `json:"allow_comments,omitempty"`
	SortBy         string `json:"sort_by,omitempty"`
	SortHow        string `json:"sort_how,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	ItemCount      int    `json:"item_count,omitempty"`
	Likes          int    `json:"likes,omitempty"`
	IDs            IDs    `json:"ids,omitempty"`
}

// ListItem is an item in a user list.
type ListItem struct {
	Rank     int      `json:"rank"`
	ListedAt string   `json:"listed_at"`
	Type     string   `json:"type"`
	Movie    *Movie   `json:"movie,omitempty"`
	Show     *Show    `json:"show,omitempty"`
	Episode  *Episode `json:"episode,omitempty"`
	Season   *Season  `json:"season,omitempty"`
	Person   *Person  `json:"person,omitempty"`
}

// ScrobbleRequest is the body for scrobble start/pause/stop.
type ScrobbleRequest struct {
	Movie    *SyncMovie   `json:"movie,omitempty"`
	Show     *SyncShow    `json:"show,omitempty"`
	Episode  *SyncEpisode `json:"episode,omitempty"`
	Progress float64      `json:"progress"`
}

// ScrobbleResponse is the response from a scrobble operation.
type ScrobbleResponse struct {
	ID      int64    `json:"id"`
	Action  string   `json:"action"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
}

// CheckinRequest is the body for a checkin.
type CheckinRequest struct {
	Movie   *SyncMovie   `json:"movie,omitempty"`
	Show    *SyncShow    `json:"show,omitempty"`
	Episode *SyncEpisode `json:"episode,omitempty"`
	Message string       `json:"message,omitempty"`
}

// CheckinResponse is returned from a successful checkin.
type CheckinResponse struct {
	ID        int64    `json:"id"`
	WatchedAt string   `json:"watched_at"`
	Movie     *Movie   `json:"movie,omitempty"`
	Show      *Show    `json:"show,omitempty"`
	Episode   *Episode `json:"episode,omitempty"`
}

// UpdatedMovie is a movie with its update timestamp.
type UpdatedMovie struct {
	UpdatedAt string `json:"updated_at"`
	Movie     Movie  `json:"movie"`
}

// UpdatedShow is a show with its update timestamp.
type UpdatedShow struct {
	UpdatedAt string `json:"updated_at"`
	Show      Show   `json:"show"`
}

// Comment represents a comment/review on a media item.
type Comment struct {
	ID         int          `json:"id,omitempty"`
	ParentID   int          `json:"parent_id,omitempty"`
	Comment    string       `json:"comment"`
	Spoiler    bool         `json:"spoiler,omitempty"`
	Review     bool         `json:"review,omitempty"`
	Likes      int          `json:"likes,omitempty"`
	Replies    int          `json:"replies,omitempty"`
	UserRating int          `json:"user_rating,omitempty"`
	CreatedAt  string       `json:"created_at,omitempty"`
	UpdatedAt  string       `json:"updated_at,omitempty"`
	User       *UserProfile `json:"user,omitempty"`
}

// CommentItem wraps a comment with its attached media.
type CommentItem struct {
	Type    string   `json:"type"`
	Comment Comment  `json:"comment"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Season  *Season  `json:"season,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
}

// Note represents a private note on a media item.
type Note struct {
	ID        int    `json:"id,omitempty"`
	Notes     string `json:"notes"`
	Privacy   string `json:"privacy,omitempty"`
	Spoiler   bool   `json:"spoiler,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// NoteItem wraps a note with its attached media.
type NoteItem struct {
	Type    string   `json:"type"`
	Note    Note     `json:"note"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Season  *Season  `json:"season,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
}

// NoteRequest is the body for creating/updating a note.
type NoteRequest struct {
	Movie   *SyncMovie   `json:"movie,omitempty"`
	Show    *SyncShow    `json:"show,omitempty"`
	Episode *SyncEpisode `json:"episode,omitempty"`
	Notes   string       `json:"notes"`
	Privacy string       `json:"privacy,omitempty"`
	Spoiler bool         `json:"spoiler,omitempty"`
}

// CalendarDVDMovie is a movie released on DVD/Blu-ray.
type CalendarDVDMovie struct {
	Released string `json:"released"`
	Movie    Movie  `json:"movie"`
}

// WatchingItem represents an item currently being watched.
type WatchingItem struct {
	ExpiresAt string   `json:"expires_at"`
	StartedAt string   `json:"started_at"`
	Action    string   `json:"action"`
	Type      string   `json:"type"`
	Movie     *Movie   `json:"movie,omitempty"`
	Show      *Show    `json:"show,omitempty"`
	Episode   *Episode `json:"episode,omitempty"`
}

// RelatedMovie is a movie related to another movie.
type RelatedMovie struct {
	Movie
}

// RelatedShow is a show related to another show.
type RelatedShow struct {
	Show
}

// EpisodeTranslation represents an episode translation.
type EpisodeTranslation struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Language string `json:"language"`
	Country  string `json:"country"`
}

// PersonMovieCredits contains movie credits for a person.
type PersonMovieCredits struct {
	Cast []PersonMovieCast `json:"cast"`
	Crew *PersonMovieCrew  `json:"crew,omitempty"`
}

// PersonMovieCast is a movie cast credit for a person.
type PersonMovieCast struct {
	Characters []string `json:"characters"`
	Movie      Movie    `json:"movie"`
}

// PersonMovieCrew groups movie crew credits by department.
type PersonMovieCrew struct {
	Production []PersonMovieCrewMember `json:"production,omitempty"`
	Directing  []PersonMovieCrewMember `json:"directing,omitempty"`
	Writing    []PersonMovieCrewMember `json:"writing,omitempty"`
	Sound      []PersonMovieCrewMember `json:"sound,omitempty"`
	Camera     []PersonMovieCrewMember `json:"camera,omitempty"`
	Art        []PersonMovieCrewMember `json:"art,omitempty"`
}

// PersonMovieCrewMember is a crew credit on a movie.
type PersonMovieCrewMember struct {
	Jobs  []string `json:"jobs"`
	Movie Movie    `json:"movie"`
}

// PersonShowCredits contains show credits for a person.
type PersonShowCredits struct {
	Cast []PersonShowCast `json:"cast"`
	Crew *PersonShowCrew  `json:"crew,omitempty"`
}

// PersonShowCast is a show cast credit for a person.
type PersonShowCast struct {
	Characters []string `json:"characters"`
	Show       Show     `json:"show"`
}

// PersonShowCrew groups show crew credits by department.
type PersonShowCrew struct {
	Production []PersonShowCrewMember `json:"production,omitempty"`
	Directing  []PersonShowCrewMember `json:"directing,omitempty"`
	Writing    []PersonShowCrewMember `json:"writing,omitempty"`
	Sound      []PersonShowCrewMember `json:"sound,omitempty"`
	Camera     []PersonShowCrewMember `json:"camera,omitempty"`
	Art        []PersonShowCrewMember `json:"art,omitempty"`
}

// PersonShowCrewMember is a crew credit on a show.
type PersonShowCrewMember struct {
	Jobs []string `json:"jobs"`
	Show Show     `json:"show"`
}

// PlaybackProgress represents a playback progress item.
type PlaybackProgress struct {
	ID       int64    `json:"id"`
	Progress float64  `json:"progress"`
	PausedAt string   `json:"paused_at"`
	Type     string   `json:"type"`
	Movie    *Movie   `json:"movie,omitempty"`
	Show     *Show    `json:"show,omitempty"`
	Episode  *Episode `json:"episode,omitempty"`
}

// LastActivities represents the last time each section was updated.
type LastActivities struct {
	All       string              `json:"all"`
	Movies    LastActivityTimes   `json:"movies"`
	Episodes  LastActivityTimes   `json:"episodes"`
	Shows     LastActivityShow    `json:"shows"`
	Seasons   LastActivitySeason  `json:"seasons"`
	Comments  LastActivityComment `json:"comments"`
	Lists     LastActivityList    `json:"lists"`
	Watchlist LastActivityTimes   `json:"watchlist"`
	Favorites LastActivityTimes   `json:"favorites"`
	Account   LastActivityAccount `json:"account"`
}

// LastActivityTimes holds timestamps for common activity types.
type LastActivityTimes struct {
	WatchedAt     string `json:"watched_at,omitempty"`
	CollectedAt   string `json:"collected_at,omitempty"`
	RatedAt       string `json:"rated_at,omitempty"`
	WatchlistedAt string `json:"watchlisted_at,omitempty"`
	FavoritedAt   string `json:"favorited_at,omitempty"`
	CommentedAt   string `json:"commented_at,omitempty"`
	PausedAt      string `json:"paused_at,omitempty"`
	HiddenAt      string `json:"hidden_at,omitempty"`
}

// LastActivityShow holds timestamps specific to shows.
type LastActivityShow struct {
	RatedAt       string `json:"rated_at,omitempty"`
	WatchlistedAt string `json:"watchlisted_at,omitempty"`
	FavoritedAt   string `json:"favorited_at,omitempty"`
	CommentedAt   string `json:"commented_at,omitempty"`
	HiddenAt      string `json:"hidden_at,omitempty"`
}

// LastActivitySeason holds timestamps specific to seasons.
type LastActivitySeason struct {
	RatedAt       string `json:"rated_at,omitempty"`
	WatchlistedAt string `json:"watchlisted_at,omitempty"`
	CommentedAt   string `json:"commented_at,omitempty"`
	HiddenAt      string `json:"hidden_at,omitempty"`
}

// LastActivityComment holds timestamps for comment activities.
type LastActivityComment struct {
	LikedAt string `json:"liked_at,omitempty"`
}

// LastActivityList holds timestamps for list activities.
type LastActivityList struct {
	LikedAt     string `json:"liked_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	CommentedAt string `json:"commented_at,omitempty"`
}

// LastActivityAccount holds timestamps for account activities.
type LastActivityAccount struct {
	SettingsAt  string `json:"settings_at,omitempty"`
	FollowedAt  string `json:"followed_at,omitempty"`
	FollowingAt string `json:"following_at,omitempty"`
	PendingAt   string `json:"pending_at,omitempty"`
}

// WatchedItem represents a watched movie or show.
type WatchedItem struct {
	Plays         int             `json:"plays"`
	LastWatchedAt string          `json:"last_watched_at"`
	LastUpdatedAt string          `json:"last_updated_at"`
	Movie         *Movie          `json:"movie,omitempty"`
	Show          *Show           `json:"show,omitempty"`
	Seasons       []WatchedSeason `json:"seasons,omitempty"`
}

// WatchedSeason is a season in a watched show.
type WatchedSeason struct {
	Number   int              `json:"number"`
	Episodes []WatchedEpisode `json:"episodes"`
}

// WatchedEpisode is an episode within a watched season.
type WatchedEpisode struct {
	Number        int    `json:"number"`
	Plays         int    `json:"plays"`
	LastWatchedAt string `json:"last_watched_at"`
}

// UserFollower represents a follower or following user.
type UserFollower struct {
	FollowedAt string      `json:"followed_at"`
	User       UserProfile `json:"user"`
}

// FollowRequest represents a pending follow request.
type FollowRequest struct {
	ID          int         `json:"id"`
	RequestedAt string      `json:"requested_at"`
	User        UserProfile `json:"user"`
}

// FavoritesItem represents a favorite item.
type FavoritesItem struct {
	Rank     int    `json:"rank"`
	ListedAt string `json:"listed_at"`
	Type     string `json:"type"`
	Movie    *Movie `json:"movie,omitempty"`
	Show     *Show  `json:"show,omitempty"`
}

// HideRequest contains items to hide from recommendations.
type HideRequest struct {
	Movies []SyncMovie `json:"movies,omitempty"`
	Shows  []SyncShow  `json:"shows,omitempty"`
}

// CommentRequest is the body for creating/updating a comment.
type CommentRequest struct {
	Comment string       `json:"comment"`
	Spoiler bool         `json:"spoiler,omitempty"`
	Review  bool         `json:"review,omitempty"`
	Movie   *SyncMovie   `json:"movie,omitempty"`
	Show    *SyncShow    `json:"show,omitempty"`
	Episode *SyncEpisode `json:"episode,omitempty"`
	Season  *SyncSeason  `json:"season,omitempty"`
}

// WatchedProgress represents the watch progress for a show.
type WatchedProgress struct {
	Aired         int              `json:"aired"`
	Completed     int              `json:"completed"`
	LastWatchedAt string           `json:"last_watched_at,omitempty"`
	ResetAt       string           `json:"reset_at,omitempty"`
	Seasons       []SeasonProgress `json:"seasons"`
	NextEpisode   *Episode         `json:"next_episode,omitempty"`
	LastEpisode   *Episode         `json:"last_episode,omitempty"`
}

// SeasonProgress represents watch progress for a season.
type SeasonProgress struct {
	Number    int               `json:"number"`
	Title     string            `json:"title,omitempty"`
	Aired     int               `json:"aired"`
	Completed int               `json:"completed"`
	Episodes  []EpisodeProgress `json:"episodes"`
}

// EpisodeProgress represents watch progress for an episode.
type EpisodeProgress struct {
	Number        int    `json:"number"`
	Completed     bool   `json:"completed"`
	LastWatchedAt string `json:"last_watched_at,omitempty"`
}

// CollectionProgress represents the collection progress for a show.
type CollectionProgress struct {
	Aired           int                        `json:"aired"`
	Completed       int                        `json:"completed"`
	LastCollectedAt string                     `json:"last_collected_at,omitempty"`
	Seasons         []CollectionSeasonProgress `json:"seasons"`
}

// CollectionSeasonProgress is collection progress for a season.
type CollectionSeasonProgress struct {
	Number    int                         `json:"number"`
	Aired     int                         `json:"aired"`
	Completed int                         `json:"completed"`
	Episodes  []CollectionEpisodeProgress `json:"episodes"`
}

// CollectionEpisodeProgress is collection progress for an episode.
type CollectionEpisodeProgress struct {
	Number      int    `json:"number"`
	Completed   bool   `json:"completed"`
	CollectedAt string `json:"collected_at,omitempty"`
}
