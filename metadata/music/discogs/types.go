package discogs

// Release represents a Discogs release (album/single/EP).
type Release struct {
	ID                int         `json:"id"`
	Title             string      `json:"title"`
	Year              int         `json:"year,omitempty"`
	ResourceURL       string      `json:"resource_url,omitempty"`
	URI               string      `json:"uri,omitempty"`
	Artists           []ArtistRef `json:"artists,omitempty"`
	ArtistsSort       string      `json:"artists_sort,omitempty"`
	Labels            []LabelRef  `json:"labels,omitempty"`
	Formats           []Format    `json:"formats,omitempty"`
	DataQuality       string      `json:"data_quality,omitempty"`
	Community         *Community  `json:"community,omitempty"`
	DateAdded         string      `json:"date_added,omitempty"`
	DateChanged       string      `json:"date_changed,omitempty"`
	NumForSale        int         `json:"num_for_sale,omitempty"`
	LowestPrice       float64     `json:"lowest_price,omitempty"`
	MasterID          int         `json:"master_id,omitempty"`
	MasterURL         string      `json:"master_url,omitempty"`
	Country           string      `json:"country,omitempty"`
	Released          string      `json:"released,omitempty"`
	Notes             string      `json:"notes,omitempty"`
	ReleasedFormatted string      `json:"released_formatted,omitempty"`
	Genres            []string    `json:"genres,omitempty"`
	Styles            []string    `json:"styles,omitempty"`
	Tracklist         []Track     `json:"tracklist,omitempty"`
	ExtraArtists      []ArtistRef `json:"extraartists,omitempty"`
	Images            []Image     `json:"images,omitempty"`
	Thumb             string      `json:"thumb,omitempty"`
	Videos            []Video     `json:"videos,omitempty"`
}

// Artist represents a Discogs artist.
type Artist struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	RealName       string   `json:"realname,omitempty"`
	ResourceURL    string   `json:"resource_url,omitempty"`
	URI            string   `json:"uri,omitempty"`
	Profile        string   `json:"profile,omitempty"`
	URLs           []string `json:"urls,omitempty"`
	NameVariations []string `json:"namevariations,omitempty"`
	Members        []Member `json:"members,omitempty"`
	Groups         []Member `json:"groups,omitempty"`
	Images         []Image  `json:"images,omitempty"`
	DataQuality    string   `json:"data_quality,omitempty"`
}

// ArtistRef is a reference to an artist within a release.
type ArtistRef struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ANV         string `json:"anv,omitempty"`
	Join        string `json:"join,omitempty"`
	Role        string `json:"role,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
	Tracks      string `json:"tracks,omitempty"`
}

// Label represents a Discogs label.
type Label struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	ResourceURL string     `json:"resource_url,omitempty"`
	URI         string     `json:"uri,omitempty"`
	Profile     string     `json:"profile,omitempty"`
	ContactInfo string     `json:"contact_info,omitempty"`
	URLs        []string   `json:"urls,omitempty"`
	Images      []Image    `json:"images,omitempty"`
	DataQuality string     `json:"data_quality,omitempty"`
	Sublabels   []LabelRef `json:"sublabels,omitempty"`
	ParentLabel *LabelRef  `json:"parent_label,omitempty"`
}

// LabelRef is a reference to a label.
type LabelRef struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CatNo       string `json:"catno,omitempty"`
	EntityType  string `json:"entity_type,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
}

// MasterRelease represents a master release (canonical version).
type MasterRelease struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	MainRelease int         `json:"main_release,omitempty"`
	Year        int         `json:"year,omitempty"`
	ResourceURL string      `json:"resource_url,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Artists     []ArtistRef `json:"artists,omitempty"`
	Genres      []string    `json:"genres,omitempty"`
	Styles      []string    `json:"styles,omitempty"`
	Tracklist   []Track     `json:"tracklist,omitempty"`
	Images      []Image     `json:"images,omitempty"`
	Videos      []Video     `json:"videos,omitempty"`
	DataQuality string      `json:"data_quality,omitempty"`
}

// Track represents a track on a release.
type Track struct {
	Position     string      `json:"position,omitempty"`
	Type         string      `json:"type_,omitempty"`
	Title        string      `json:"title"`
	Duration     string      `json:"duration,omitempty"`
	ExtraArtists []ArtistRef `json:"extraartists,omitempty"`
}

// Format describes the physical format of a release.
type Format struct {
	Name         string   `json:"name"`
	Qty          string   `json:"qty,omitempty"`
	Text         string   `json:"text,omitempty"`
	Descriptions []string `json:"descriptions,omitempty"`
}

// Image represents an image resource.
type Image struct {
	Type        string `json:"type"`
	URI         string `json:"uri,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
	URI150      string `json:"uri150,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
}

// Video represents a video resource.
type Video struct {
	URI         string `json:"uri"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	Embed       bool   `json:"embed,omitempty"`
}

// Community contains community data for a release.
type Community struct {
	Have        int     `json:"have"`
	Want        int     `json:"want"`
	Rating      *Rating `json:"rating,omitempty"`
	Status      string  `json:"status,omitempty"`
	DataQuality string  `json:"data_quality,omitempty"`
}

// Rating represents community rating.
type Rating struct {
	Count   int     `json:"count"`
	Average float64 `json:"average"`
}

// Member represents a group member or group association.
type Member struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Active      bool   `json:"active,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
}

// SearchResponse is the response from a database search.
type SearchResponse struct {
	Pagination Pagination     `json:"pagination"`
	Results    []SearchResult `json:"results"`
}

// Pagination contains pagination info.
type Pagination struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"per_page"`
	Items   int `json:"items"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Thumb       string   `json:"thumb,omitempty"`
	CoverImage  string   `json:"cover_image,omitempty"`
	ResourceURL string   `json:"resource_url,omitempty"`
	URI         string   `json:"uri,omitempty"`
	Country     string   `json:"country,omitempty"`
	Year        string   `json:"year,omitempty"`
	Genre       []string `json:"genre,omitempty"`
	Style       []string `json:"style,omitempty"`
	Format      []string `json:"format,omitempty"`
	Label       []string `json:"label,omitempty"`
	Barcode     []string `json:"barcode,omitempty"`
	CatNo       string   `json:"catno,omitempty"`
}

// Release Ratings.

// ReleaseRating represents a user's rating for a release.
type ReleaseRating struct {
	Username  string `json:"username"`
	ReleaseID int    `json:"release_id"`
	Rating    int    `json:"rating"`
}

// CommunityRating represents the community rating for a release.
type CommunityRating struct {
	ReleaseID int    `json:"release_id"`
	Rating    Rating `json:"rating"`
}

// ReleaseStats contains marketplace statistics for a release.
type ReleaseStats struct {
	NumHave int `json:"num_have"`
	NumWant int `json:"num_want"`
}

// User Identity.

// Identity represents the identity of the authenticated user.
type Identity struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	ResourceURL  string `json:"resource_url"`
	ConsumerName string `json:"consumer_name"`
}

// Profile represents a Discogs user profile.
type Profile struct {
	ID                   int     `json:"id"`
	Username             string  `json:"username"`
	Name                 string  `json:"name,omitempty"`
	Email                string  `json:"email,omitempty"`
	ResourceURL          string  `json:"resource_url,omitempty"`
	URI                  string  `json:"uri,omitempty"`
	HomePage             string  `json:"home_page,omitempty"`
	Location             string  `json:"location,omitempty"`
	Profile              string  `json:"profile,omitempty"`
	Registered           string  `json:"registered,omitempty"`
	NumLists             int     `json:"num_lists,omitempty"`
	NumForSale           int     `json:"num_for_sale,omitempty"`
	NumCollection        int     `json:"num_collection,omitempty"`
	NumWantlist          int     `json:"num_wantlist,omitempty"`
	NumPending           int     `json:"num_pending,omitempty"`
	ReleasesContributed  int     `json:"releases_contributed,omitempty"`
	Rank                 int     `json:"rank,omitempty"`
	ReleasesRated        int     `json:"releases_rated,omitempty"`
	RatingAvg            float64 `json:"rating_avg,omitempty"`
	InventoryURL         string  `json:"inventory_url,omitempty"`
	CollectionFoldersURL string  `json:"collection_folders_url,omitempty"`
	WantlistURL          string  `json:"wantlist_url,omitempty"`
	AvatarURL            string  `json:"avatar_url,omitempty"`
	BannerURL            string  `json:"banner_url,omitempty"`
	BuyerRating          float64 `json:"buyer_rating,omitempty"`
	BuyerRatingStars     float64 `json:"buyer_rating_stars,omitempty"`
	BuyerNumRatings      int     `json:"buyer_num_ratings,omitempty"`
	SellerRating         float64 `json:"seller_rating,omitempty"`
	SellerRatingStars    float64 `json:"seller_rating_stars,omitempty"`
	SellerNumRatings     int     `json:"seller_num_ratings,omitempty"`
	CurrAbbr             string  `json:"curr_abbr,omitempty"`
}

// ProfileUpdate contains fields that can be updated on a user profile.
type ProfileUpdate struct {
	Name     string `json:"name,omitempty"`
	HomePage string `json:"home_page,omitempty"`
	Location string `json:"location,omitempty"`
	Profile  string `json:"profile,omitempty"`
	CurrAbbr string `json:"curr_abbr,omitempty"`
}

// SubmissionsResponse contains user submissions.
type SubmissionsResponse struct {
	Pagination  Pagination `json:"pagination"`
	Submissions struct {
		Artists  []Artist  `json:"artists,omitempty"`
		Labels   []Label   `json:"labels,omitempty"`
		Releases []Release `json:"releases,omitempty"`
	} `json:"submissions"`
}

// ContributionsResponse contains user contributions.
type ContributionsResponse struct {
	Pagination    Pagination     `json:"pagination"`
	Contributions []SearchResult `json:"contributions"`
}

// User Collection.

// CollectionFoldersResponse contains collection folders.
type CollectionFoldersResponse struct {
	Folders []CollectionFolder `json:"folders"`
}

// CollectionFolder represents a collection folder.
type CollectionFolder struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Count       int    `json:"count"`
	ResourceURL string `json:"resource_url,omitempty"`
}

// CollectionItemsResponse contains collection items.
type CollectionItemsResponse struct {
	Pagination Pagination       `json:"pagination"`
	Releases   []CollectionItem `json:"releases"`
}

// CollectionItem represents an item in a collection.
type CollectionItem struct {
	ID         int         `json:"id"`
	InstanceID int         `json:"instance_id"`
	FolderID   int         `json:"folder_id"`
	Rating     int         `json:"rating"`
	DateAdded  string      `json:"date_added,omitempty"`
	BasicInfo  *BasicInfo  `json:"basic_information,omitempty"`
	Notes      []FieldNote `json:"notes,omitempty"`
}

// BasicInfo is basic release info within a collection item.
type BasicInfo struct {
	ID         int         `json:"id"`
	Title      string      `json:"title"`
	Year       int         `json:"year,omitempty"`
	Thumb      string      `json:"thumb,omitempty"`
	CoverImage string      `json:"cover_image,omitempty"`
	Artists    []ArtistRef `json:"artists,omitempty"`
	Labels     []LabelRef  `json:"labels,omitempty"`
	Formats    []Format    `json:"formats,omitempty"`
	Genres     []string    `json:"genres,omitempty"`
	Styles     []string    `json:"styles,omitempty"`
}

// FieldNote represents a custom field value on a collection item.
type FieldNote struct {
	FieldID int    `json:"field_id"`
	Value   string `json:"value"`
}

// CustomFieldsResponse contains custom fields.
type CustomFieldsResponse struct {
	Fields []CustomField `json:"fields"`
}

// CustomField represents a custom field definition.
type CustomField struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Public   bool     `json:"public"`
	Position int      `json:"position"`
	Options  []string `json:"options,omitempty"`
	Lines    int      `json:"lines,omitempty"`
}

// CollectionValue represents the estimated value of a collection.
type CollectionValue struct {
	Maximum string `json:"maximum,omitempty"`
	Median  string `json:"median,omitempty"`
	Minimum string `json:"minimum,omitempty"`
}

// User Wantlist.

// WantlistResponse contains wantlist items.
type WantlistResponse struct {
	Pagination Pagination     `json:"pagination"`
	Wants      []WantlistItem `json:"wants"`
}

// WantlistItem represents an item in the wantlist.
type WantlistItem struct {
	ID          int        `json:"id"`
	Rating      int        `json:"rating"`
	Notes       string     `json:"notes,omitempty"`
	DateAdded   string     `json:"date_added,omitempty"`
	ResourceURL string     `json:"resource_url,omitempty"`
	BasicInfo   *BasicInfo `json:"basic_information,omitempty"`
}

// User Lists.

// UserListsResponse contains user lists.
type UserListsResponse struct {
	Pagination Pagination     `json:"pagination"`
	Lists      []UserListMeta `json:"lists"`
}

// UserListMeta is summary metadata for a user list.
type UserListMeta struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Public      bool   `json:"public"`
	DateAdded   string `json:"date_added,omitempty"`
	DateChanged string `json:"date_changed,omitempty"`
	URI         string `json:"uri,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
}

// UserList represents a full user list with items.
type UserList struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Public      bool           `json:"public"`
	DateAdded   string         `json:"date_added,omitempty"`
	DateChanged string         `json:"date_changed,omitempty"`
	URI         string         `json:"uri,omitempty"`
	ResourceURL string         `json:"resource_url,omitempty"`
	Items       []UserListItem `json:"items,omitempty"`
}

// UserListItem represents an item in a user list.
type UserListItem struct {
	ID           int    `json:"id"`
	Type         string `json:"type"`
	Comment      string `json:"comment,omitempty"`
	URI          string `json:"uri,omitempty"`
	ResourceURL  string `json:"resource_url,omitempty"`
	DisplayTitle string `json:"display_title,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
}

// Marketplace.

// InventoryResponse contains inventory listings.
type InventoryResponse struct {
	Pagination Pagination `json:"pagination"`
	Listings   []Listing  `json:"listings"`
}

// Listing represents a marketplace listing.
type Listing struct {
	ID              int           `json:"id"`
	ResourceURL     string        `json:"resource_url,omitempty"`
	URI             string        `json:"uri,omitempty"`
	Status          string        `json:"status"`
	Condition       string        `json:"condition"`
	SleeveCondition string        `json:"sleeve_condition,omitempty"`
	Comments        string        `json:"comments,omitempty"`
	ShipsFrom       string        `json:"ships_from,omitempty"`
	Posted          string        `json:"posted,omitempty"`
	AllowOffers     bool          `json:"allow_offers"`
	FormatQuantity  int           `json:"format_quantity,omitempty"`
	ExternalID      string        `json:"external_id,omitempty"`
	Location        string        `json:"location,omitempty"`
	Weight          float64       `json:"weight,omitempty"`
	Price           *ListingPrice `json:"price,omitempty"`
	Release         *BasicInfo    `json:"release,omitempty"`
	Seller          *Seller       `json:"seller,omitempty"`
	OriginalPrice   *ListingPrice `json:"original_price,omitempty"`
}

// ListingPrice represents a listing price.
type ListingPrice struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// Seller represents a marketplace seller.
type Seller struct {
	ID          int     `json:"id"`
	Username    string  `json:"username"`
	ResourceURL string  `json:"resource_url,omitempty"`
	Rating      float64 `json:"rating,omitempty"`
	Stars       float64 `json:"stars,omitempty"`
	AvatarURL   string  `json:"avatar_url,omitempty"`
}

// NewListing contains fields for creating/editing a marketplace listing.
type NewListing struct {
	ReleaseID       int     `json:"release_id"`
	Condition       string  `json:"condition"`
	SleeveCondition string  `json:"sleeve_condition,omitempty"`
	Price           float64 `json:"price"`
	Comments        string  `json:"comments,omitempty"`
	AllowOffers     bool    `json:"allow_offers,omitempty"`
	Status          string  `json:"status,omitempty"`
	ExternalID      string  `json:"external_id,omitempty"`
	Location        string  `json:"location,omitempty"`
	Weight          float64 `json:"weight,omitempty"`
	FormatQuantity  int     `json:"format_quantity,omitempty"`
}

// Order represents a marketplace order.
type Order struct {
	ID                     string        `json:"id"`
	ResourceURL            string        `json:"resource_url,omitempty"`
	MessagesURL            string        `json:"messages_url,omitempty"`
	URI                    string        `json:"uri,omitempty"`
	Status                 string        `json:"status"`
	NextStatus             []string      `json:"next_status,omitempty"`
	Fee                    *ListingPrice `json:"fee,omitempty"`
	Created                string        `json:"created,omitempty"`
	LastActivity           string        `json:"last_activity,omitempty"`
	Buyer                  *Seller       `json:"buyer,omitempty"`
	Seller                 *Seller       `json:"seller,omitempty"`
	Items                  []OrderItem   `json:"items,omitempty"`
	Shipping               *ListingPrice `json:"shipping,omitempty"`
	ShippingAddress        string        `json:"shipping_address,omitempty"`
	AdditionalInstructions string        `json:"additional_instructions,omitempty"`
	Total                  *ListingPrice `json:"total,omitempty"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	ID      int           `json:"id"`
	Price   *ListingPrice `json:"price,omitempty"`
	Release *BasicInfo    `json:"release,omitempty"`
}

// OrderUpdate contains fields for updating an order.
type OrderUpdate struct {
	Status   string  `json:"status,omitempty"`
	Shipping float64 `json:"shipping,omitempty"`
}

// OrdersResponse contains a list of orders.
type OrdersResponse struct {
	Pagination Pagination `json:"pagination"`
	Orders     []Order    `json:"orders"`
}

// OrderMessagesResponse contains order messages.
type OrderMessagesResponse struct {
	Pagination Pagination     `json:"pagination"`
	Messages   []OrderMessage `json:"messages"`
}

// OrderMessage represents a message in an order.
type OrderMessage struct {
	Subject   string  `json:"subject,omitempty"`
	Message   string  `json:"message,omitempty"`
	From      *Seller `json:"from,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
	StatusID  int     `json:"status_id,omitempty"`
	Order     *Order  `json:"order,omitempty"`
}

// Fee represents a marketplace fee.
type Fee struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// PriceSuggestions contains price suggestions for different conditions.
type PriceSuggestions struct {
	VeryGood     *SuggestedPrice `json:"Very Good (VG),omitempty"`
	GoodPlus     *SuggestedPrice `json:"Good Plus (G+),omitempty"`
	NearMint     *SuggestedPrice `json:"Near Mint (NM or M-),omitempty"`
	Mint         *SuggestedPrice `json:"Mint (M),omitempty"`
	VeryGoodPlus *SuggestedPrice `json:"Very Good Plus (VG+),omitempty"`
	Fair         *SuggestedPrice `json:"Fair (F),omitempty"`
	Good         *SuggestedPrice `json:"Good (G),omitempty"`
	Poor         *SuggestedPrice `json:"Poor (P),omitempty"`
}

// SuggestedPrice represents a suggested price for a condition.
type SuggestedPrice struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// MarketplaceReleaseStats contains marketplace statistics for a release.
type MarketplaceReleaseStats struct {
	LowestPrice *ListingPrice `json:"lowest_price,omitempty"`
	NumForSale  int           `json:"num_for_sale"`
	Blocked     bool          `json:"blocked,omitempty"`
}

// Inventory Export.

// ExportsResponse contains recent exports.
type ExportsResponse struct {
	Pagination Pagination `json:"pagination"`
	Items      []Export   `json:"items"`
}

// Export represents an inventory export.
type Export struct {
	ID          int    `json:"id"`
	Status      string `json:"status"`
	CreatedTS   string `json:"created_ts,omitempty"`
	FinishedTS  string `json:"finished_ts,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
	Filename    string `json:"filename,omitempty"`
}
