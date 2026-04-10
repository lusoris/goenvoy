package kitsu

// Titles holds localized titles keyed by locale code (e.g. "en", "en_jp", "ja_jp").
type Titles map[string]string

// ImageSet holds image URLs at multiple sizes plus optional dimension metadata.
type ImageSet struct {
	Tiny     string        `json:"tiny"`
	Large    string        `json:"large"`
	Small    string        `json:"small"`
	Medium   string        `json:"medium"`
	Original string        `json:"original"`
	Meta     *ImageSetMeta `json:"meta,omitempty"`
}

// ImageSetMeta holds dimensional metadata for an [ImageSet].
type ImageSetMeta struct {
	Dimensions map[string]ImageDimension `json:"dimensions,omitempty"`
}

// ImageDimension holds width and height for a single image size.
type ImageDimension struct {
	Width  *int `json:"width,omitempty"`
	Height *int `json:"height,omitempty"`
}

// Thumbnail holds a single thumbnail image.
type Thumbnail struct {
	Original string        `json:"original"`
	Meta     *ImageSetMeta `json:"meta,omitempty"`
}

// Anime represents a Kitsu anime resource.
type Anime struct {
	ID                  string            `json:"id"`
	CreatedAt           string            `json:"createdAt"`
	UpdatedAt           string            `json:"updatedAt"`
	Slug                string            `json:"slug"`
	Synopsis            string            `json:"synopsis"`
	Description         string            `json:"description"`
	CoverImageTopOffset int               `json:"coverImageTopOffset"`
	Titles              Titles            `json:"titles"`
	CanonicalTitle      string            `json:"canonicalTitle"`
	AbbreviatedTitles   []string          `json:"abbreviatedTitles"`
	AverageRating       string            `json:"averageRating"`
	RatingFrequencies   map[string]string `json:"ratingFrequencies"`
	UserCount           int               `json:"userCount"`
	FavoritesCount      int               `json:"favoritesCount"`
	StartDate           string            `json:"startDate"`
	EndDate             string            `json:"endDate"`
	NextRelease         *string           `json:"nextRelease"`
	PopularityRank      int               `json:"popularityRank"`
	RatingRank          int               `json:"ratingRank"`
	AgeRating           string            `json:"ageRating"`
	AgeRatingGuide      string            `json:"ageRatingGuide"`
	Subtype             string            `json:"subtype"`
	Status              string            `json:"status"`
	TBA                 string            `json:"tba"`
	PosterImage         *ImageSet         `json:"posterImage"`
	CoverImage          *ImageSet         `json:"coverImage"`
	EpisodeCount        *int              `json:"episodeCount"`
	EpisodeLength       *int              `json:"episodeLength"`
	TotalLength         *int              `json:"totalLength"`
	YoutubeVideoID      string            `json:"youtubeVideoId"` //nolint:tagliatelle // API uses camelCase "Id".
	ShowType            string            `json:"showType"`
	NSFW                bool              `json:"nsfw"`
}

// Manga represents a Kitsu manga resource.
type Manga struct {
	ID                  string            `json:"id"`
	CreatedAt           string            `json:"createdAt"`
	UpdatedAt           string            `json:"updatedAt"`
	Slug                string            `json:"slug"`
	Synopsis            string            `json:"synopsis"`
	Description         string            `json:"description"`
	CoverImageTopOffset int               `json:"coverImageTopOffset"`
	Titles              Titles            `json:"titles"`
	CanonicalTitle      string            `json:"canonicalTitle"`
	AbbreviatedTitles   []string          `json:"abbreviatedTitles"`
	AverageRating       string            `json:"averageRating"`
	RatingFrequencies   map[string]string `json:"ratingFrequencies"`
	UserCount           int               `json:"userCount"`
	FavoritesCount      int               `json:"favoritesCount"`
	StartDate           string            `json:"startDate"`
	EndDate             string            `json:"endDate"`
	NextRelease         *string           `json:"nextRelease"`
	PopularityRank      int               `json:"popularityRank"`
	RatingRank          int               `json:"ratingRank"`
	AgeRating           *string           `json:"ageRating"`
	AgeRatingGuide      *string           `json:"ageRatingGuide"`
	Subtype             string            `json:"subtype"`
	Status              string            `json:"status"`
	TBA                 string            `json:"tba"`
	PosterImage         *ImageSet         `json:"posterImage"`
	CoverImage          *ImageSet         `json:"coverImage"`
	ChapterCount        *int              `json:"chapterCount"`
	VolumeCount         int               `json:"volumeCount"`
	Serialization       string            `json:"serialization"`
	MangaType           string            `json:"mangaType"`
}

// Episode represents a Kitsu episode resource.
type Episode struct {
	ID             string     `json:"id"`
	CreatedAt      string     `json:"createdAt"`
	UpdatedAt      string     `json:"updatedAt"`
	Synopsis       string     `json:"synopsis"`
	Description    string     `json:"description"`
	Titles         Titles     `json:"titles"`
	CanonicalTitle string     `json:"canonicalTitle"`
	SeasonNumber   *int       `json:"seasonNumber"`
	Number         *int       `json:"number"`
	RelativeNumber *int       `json:"relativeNumber"`
	Airdate        string     `json:"airdate"`
	Length         *int       `json:"length"`
	Thumbnail      *Thumbnail `json:"thumbnail"`
}

// Character represents a Kitsu character resource.
type Character struct {
	ID            string    `json:"id"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
	Slug          string    `json:"slug"`
	Names         Titles    `json:"names"`
	CanonicalName string    `json:"canonicalName"`
	OtherNames    []string  `json:"otherNames"`
	Name          string    `json:"name"`
	MalID         int       `json:"malId"` //nolint:tagliatelle // API uses camelCase "Id".
	Description   string    `json:"description"`
	Image         *ImageSet `json:"image"`
}

// Category represents a Kitsu category resource.
type Category struct {
	ID              string `json:"id"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	TotalMediaCount int    `json:"totalMediaCount"`
	Slug            string `json:"slug"`
	NSFW            bool   `json:"nsfw"`
	ChildCount      int    `json:"childCount"`
}

// User represents a Kitsu user resource.
type User struct {
	ID                  string    `json:"id"`
	CreatedAt           string    `json:"createdAt"`
	UpdatedAt           string    `json:"updatedAt"`
	Name                string    `json:"name"`
	PastNames           []string  `json:"pastNames"`
	Slug                string    `json:"slug"`
	About               string    `json:"about"`
	Location            string    `json:"location"`
	WaifuOrHusbando     string    `json:"waifuOrHusbando"`
	FollowersCount      int       `json:"followersCount"`
	FollowingCount      int       `json:"followingCount"`
	LifeSpentOnAnime    int       `json:"lifeSpentOnAnime"`
	Birthday            *string   `json:"birthday"`
	Gender              *string   `json:"gender"`
	CommentsCount       int       `json:"commentsCount"`
	FavoritesCount      int       `json:"favoritesCount"`
	LikesGivenCount     int       `json:"likesGivenCount"`
	ReviewsCount        int       `json:"reviewsCount"`
	LikesReceivedCount  int       `json:"likesReceivedCount"`
	PostsCount          int       `json:"postsCount"`
	RatingsCount        int       `json:"ratingsCount"`
	MediaReactionsCount int       `json:"mediaReactionsCount"`
	Title               *string   `json:"title"`
	ProfileCompleted    bool      `json:"profileCompleted"`
	FeedCompleted       bool      `json:"feedCompleted"`
	Website             string    `json:"website"`
	Avatar              *ImageSet `json:"avatar"`
	CoverImage          *ImageSet `json:"coverImage"`
	Status              string    `json:"status"`
}

// PageLinks holds pagination URLs returned by the Kitsu API.
type PageLinks struct {
	First string `json:"first"`
	Next  string `json:"next"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
}

// Token holds OAuth2 access and refresh tokens from Kitsu.
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}

// Mapping holds an external ID mapping for an anime.
type Mapping struct {
	ID           string `json:"id"`
	ExternalSite string `json:"externalSite"`
	ExternalID   string `json:"externalId"` //nolint:tagliatelle // API uses camelCase "Id".
}
