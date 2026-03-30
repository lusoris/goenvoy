package anidb

import "encoding/xml"

// Title represents a title with language and type metadata.
type Title struct {
	Lang string `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:",chardata"`
}

// Rating holds a numeric rating value with vote or count metadata.
type Rating struct {
	Count int    `xml:"count,attr"`
	Votes int    `xml:"votes,attr"`
	Value string `xml:",chardata"`
}

// Ratings contains the aggregated ratings for an anime.
type Ratings struct {
	Permanent       Rating `xml:"permanent"`
	Temporary       Rating `xml:"temporary"`
	Review          Rating `xml:"review"`
	Recommendations string `xml:"recommendations"`
}

// Anime holds the full details returned by the anime data command.
type Anime struct {
	XMLName         xml.Name              `xml:"anime"`
	ID              int                   `xml:"id,attr"`
	Restricted      bool                  `xml:"restricted,attr"`
	Type            string                `xml:"type"`
	EpisodeCount    int                   `xml:"episodecount"`
	StartDate       string                `xml:"startdate"`
	EndDate         string                `xml:"enddate"`
	Titles          []Title               `xml:"titles>title"`
	RelatedAnime    []RelatedAnime        `xml:"relatedanime>anime"`
	SimilarAnime    []SimilarAnime        `xml:"similaranime>anime"`
	Recommendations []AnimeRecommendation `xml:"recommendations>recommendation"`
	URL             string                `xml:"url"`
	Creators        []Creator             `xml:"creators>name"`
	Description     string                `xml:"description"`
	Ratings         Ratings               `xml:"ratings"`
	Picture         string                `xml:"picture"`
	Resources       []Resource            `xml:"resources>resource"`
	Tags            []Tag                 `xml:"tags>tag"`
	Characters      []Character           `xml:"characters>character"`
	Episodes        []Episode             `xml:"episodes>episode"`
}

// RelatedAnime represents a related anime entry.
type RelatedAnime struct {
	ID   int    `xml:"id,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:",chardata"`
}

// SimilarAnime represents a similar anime entry with approval stats.
type SimilarAnime struct {
	ID       int    `xml:"id,attr"`
	Approval int    `xml:"approval,attr"`
	Total    int    `xml:"total,attr"`
	Name     string `xml:",chardata"`
}

// AnimeRecommendation is a user recommendation in the full anime response.
type AnimeRecommendation struct {
	Type string `xml:"type,attr"`
	UID  int    `xml:"uid,attr"`
	Text string `xml:",chardata"`
}

// Creator represents a staff member credited on an anime.
type Creator struct {
	ID   int    `xml:"id,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:",chardata"`
}

// Resource links an anime to external databases and websites.
type Resource struct {
	Type             int              `xml:"type,attr"`
	ExternalEntities []ExternalEntity `xml:"externalentity"`
}

// ExternalEntity holds identifiers and URLs for an external resource.
type ExternalEntity struct {
	Identifiers []string `xml:"identifier"`
	URLs        []string `xml:"url"`
}

// Tag represents an anime tag with hierarchy and spoiler metadata.
type Tag struct {
	ID            int    `xml:"id,attr"`
	ParentID      int    `xml:"parentid,attr"`
	InfoBox       bool   `xml:"infobox,attr"`
	Weight        int    `xml:"weight,attr"`
	LocalSpoiler  bool   `xml:"localspoiler,attr"`
	GlobalSpoiler bool   `xml:"globalspoiler,attr"`
	Verified      bool   `xml:"verified,attr"`
	Update        string `xml:"update,attr"`
	Name          string `xml:"name"`
	Description   string `xml:"description"`
	PicURL        string `xml:"picurl"`
}

// Character represents an anime character.
type Character struct {
	ID            int           `xml:"id,attr"`
	Type          string        `xml:"type,attr"`
	Update        string        `xml:"update,attr"`
	Rating        Rating        `xml:"rating"`
	Name          string        `xml:"name"`
	Gender        string        `xml:"gender"`
	CharacterType CharacterType `xml:"charactertype"`
	Description   string        `xml:"description"`
	Picture       string        `xml:"picture"`
	Seiyuu        *Seiyuu       `xml:"seiyuu"`
}

// CharacterType classifies a character entry.
type CharacterType struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:",chardata"`
}

// Seiyuu identifies the voice actor for a character.
type Seiyuu struct {
	ID      int    `xml:"id,attr"`
	Picture string `xml:"picture,attr"`
	Name    string `xml:",chardata"`
}

// Episode represents an anime episode.
type Episode struct {
	ID      int     `xml:"id,attr"`
	Update  string  `xml:"update,attr"`
	EpNo    EpNo    `xml:"epno"`
	Length  int     `xml:"length"`
	AirDate string  `xml:"airdate"`
	Rating  Rating  `xml:"rating"`
	Titles  []Title `xml:"title"`
}

// EpNo holds an episode number with its type classification.
type EpNo struct {
	Type  int    `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// AnimeRef is a lightweight anime reference used in hot anime,
// random recommendation, and main page responses.
type AnimeRef struct {
	ID           int     `xml:"id,attr"`
	Restricted   bool    `xml:"restricted,attr"`
	Type         string  `xml:"type"`
	EpisodeCount int     `xml:"episodecount"`
	StartDate    string  `xml:"startdate"`
	EndDate      string  `xml:"enddate"`
	Title        Title   `xml:"title"`
	Picture      string  `xml:"picture"`
	Ratings      Ratings `xml:"ratings"`
}

// RecommendationEntry wraps the anime within a random recommendation.
type RecommendationEntry struct {
	Anime AnimeRef `xml:"anime"`
}

// SimilarPair contains a source and target anime in a similarity match.
type SimilarPair struct {
	Source SimilarRef `xml:"source"`
	Target SimilarRef `xml:"target"`
}

// SimilarRef is a lightweight anime reference used in similar pairs.
type SimilarRef struct {
	AID        int    `xml:"aid,attr"`
	Restricted bool   `xml:"restricted,attr"`
	Title      Title  `xml:"title"`
	Picture    string `xml:"picture"`
}

// MainPage holds the combined response from the main data command.
type MainPage struct {
	XMLName              xml.Name              `xml:"main"`
	HotAnime             []AnimeRef            `xml:"hotanime>anime"`
	RandomSimilar        []SimilarPair         `xml:"randomsimilar>similar"`
	RandomRecommendation []RecommendationEntry `xml:"randomrecommendation>recommendation"`
}
