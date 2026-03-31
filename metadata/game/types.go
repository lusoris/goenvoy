package game

// Game represents a video game from a metadata provider.
type Game struct {
	// ID is the provider-specific identifier.
	ID int `json:"id"`
	// Name is the game title.
	Name string `json:"name"`
	// Summary is a short description of the game.
	Summary string `json:"summary"`
	// Slug is the URL-friendly identifier.
	Slug string `json:"slug"`
}

// Platform represents a gaming platform.
type Platform struct {
	// ID is the provider-specific identifier.
	ID int `json:"id"`
	// Name is the platform name (e.g. "PlayStation 5").
	Name string `json:"name"`
	// Slug is the URL-friendly identifier.
	Slug string `json:"slug"`
}

// Genre represents a game genre.
type Genre struct {
	// ID is the provider-specific identifier.
	ID int `json:"id"`
	// Name is the genre name (e.g. "Role-playing (RPG)").
	Name string `json:"name"`
	// Slug is the URL-friendly identifier.
	Slug string `json:"slug"`
}
