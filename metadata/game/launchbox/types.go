package launchbox

import "encoding/xml"

// Game represents a game entry in the LaunchBox database.
type Game struct {
	XMLName             xml.Name `xml:"Game" json:"-"`
	DatabaseID          int      `xml:"DatabaseID" json:"databaseId"`
	Name                string   `xml:"Name" json:"name"`
	ReleaseDate         string   `xml:"ReleaseDate" json:"releaseDate,omitempty"`
	Overview            string   `xml:"Overview" json:"overview,omitempty"`
	MaxPlayers          int      `xml:"MaxPlayers" json:"maxPlayers,omitempty"`
	VideoURL            string   `xml:"VideoURL" json:"videoUrl,omitempty"`
	DatabaseName        string   `xml:"DatabaseName" json:"databaseName,omitempty"`
	Platform            string   `xml:"Platform" json:"platform,omitempty"`
	Developer           string   `xml:"Developer" json:"developer,omitempty"`
	Publisher           string   `xml:"Publisher" json:"publisher,omitempty"`
	Genres              string   `xml:"Genres" json:"genres,omitempty"`
	CommunityRating     float64  `xml:"CommunityRating" json:"communityRating,omitempty"`
	WikipediaURL        string   `xml:"WikipediaURL" json:"wikipediaUrl,omitempty"`
	CommunityStarRating float64  `xml:"CommunityStarRating" json:"communityStarRating,omitempty"`
}

// GameAlternateName represents an alternate name for a game.
type GameAlternateName struct {
	XMLName         xml.Name `xml:"GameAlternateName" json:"-"`
	DatabaseID      int      `xml:"DatabaseID" json:"databaseId"`
	AlternateNameID int      `xml:"AlternateNameID" json:"alternateNameId"`
	Name            string   `xml:"AlternateName" json:"name"`
	Region          string   `xml:"Region" json:"region,omitempty"`
}

// GameImage represents an image for a game.
type GameImage struct {
	XMLName    xml.Name `xml:"GameImage" json:"-"`
	DatabaseID int      `xml:"DatabaseID" json:"databaseId"`
	FileName   string   `xml:"FileName" json:"fileName"`
	Type       string   `xml:"Type" json:"type"`
	Region     string   `xml:"Region" json:"region,omitempty"`
	CRC32      string   `xml:"CRC32" json:"crc32,omitempty"`
}

// Platform represents a platform in the LaunchBox database.
type Platform struct {
	XMLName      xml.Name `xml:"Platform" json:"-"`
	Name         string   `xml:"Name" json:"name"`
	Emulated     bool     `xml:"Emulated" json:"emulated"`
	ReleaseDate  string   `xml:"ReleaseDate" json:"releaseDate,omitempty"`
	Developer    string   `xml:"Developer" json:"developer,omitempty"`
	Manufacturer string   `xml:"Manufacturer" json:"manufacturer,omitempty"`
	Category     string   `xml:"Category" json:"category,omitempty"`
}

// metadataXML is the top-level container for Metadata.xml.
type metadataXML struct {
	XMLName        xml.Name            `xml:"LaunchBox"`
	Games          []Game              `xml:"Game"`
	AlternateNames []GameAlternateName `xml:"GameAlternateName"`
	Images         []GameImage         `xml:"GameImage"`
}

// platformsXML is the top-level container for Platforms.xml.
type platformsXML struct {
	XMLName   xml.Name   `xml:"LaunchBox"`
	Platforms []Platform `xml:"Platform"`
}
