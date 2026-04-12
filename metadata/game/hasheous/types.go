package hasheous

// HashLookupRequest is the request body for POST /Lookup/ByHash.
type HashLookupRequest struct {
	MD5    string `json:"mD5,omitempty"`
	SHA1   string `json:"shA1,omitempty"`
	SHA256 string `json:"shA256,omitempty"`
	CRC    string `json:"crc,omitempty"`
}

// HashLookup is the response from a hash lookup.
type HashLookup struct {
	ID         int64                        `json:"id"`
	Name       string                       `json:"name"`
	Platform   *MiniDataObject              `json:"platform"`
	Publisher  *MiniDataObject              `json:"publisher"`
	Signature  *SignatureResult             `json:"signature"`
	Signatures map[string][]SignatureResult `json:"signatures"`
	Metadata   []MetadataItem               `json:"metadata"`
	Attributes []AttributeItem              `json:"attributes"`
}

// MiniDataObject holds a name and optional metadata.
type MiniDataObject struct {
	Name     string         `json:"name"`
	Metadata []MetadataItem `json:"metadata"`
}

// SignatureResult holds game and ROM signature data.
type SignatureResult struct {
	Game *GameSignature `json:"game"`
	ROM  *ROMSignature  `json:"rom"`
}

// GameSignature identifies a game from a DAT signature source.
type GameSignature struct {
	ID            string            `json:"id"`
	CloneOfID     string            `json:"cloneOfId"`
	GameID        string            `json:"gameId"`
	Category      string            `json:"category"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Year          string            `json:"year"`
	Publisher     string            `json:"publisher"`
	System        string            `json:"system"`
	SystemVariant string            `json:"systemVariant"`
	Video         string            `json:"video"`
	CountryString string            `json:"countryString"`
	Country       map[string]string `json:"country"`
	LanguageStr   string            `json:"languageString"`
	Language      map[string]string `json:"language"`
	Copyright     string            `json:"copyright"`
}

// ROMSignature identifies a ROM from a DAT signature source.
type ROMSignature struct {
	Name   string `json:"name"`
	Size   *int64 `json:"size"`
	CRC    string `json:"crc"`
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
	Status string `json:"status"`
}

// MetadataItem references a metadata object.
type MetadataItem struct {
	ObjectType  string `json:"objectType"`
	ID          string `json:"id"`
	ImmutableID string `json:"immutableId"`
	Status      string `json:"status"`
}

// AttributeItem holds a compiled attribute.
type AttributeItem struct {
	AttributeType         string `json:"attributeType"`
	AttributeName         string `json:"attributeName"`
	AttributeRelationType string `json:"attributeRelationType"`
	Value                 any    `json:"value"`
	Link                  string `json:"link"`
}
