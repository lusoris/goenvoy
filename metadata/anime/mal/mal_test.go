package mal_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/anime/mal"
)

func testClient(t *testing.T, handler http.HandlerFunc) *mal.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return mal.New("test-client-id", metadata.WithBaseURL(srv.URL))
}

const animeJSON = `{
	"id": 30230,
	"title": "Diamond no Ace: Second Season",
	"main_picture": {
		"medium": "https://cdn.mal.net/images/anime/30230.jpg",
		"large": "https://cdn.mal.net/images/anime/30230l.jpg"
	},
	"alternative_titles": {
		"synonyms": ["Ace of Diamond: Second Season"],
		"en": "Ace of the Diamond: Second Season",
		"ja": "\u30c0\u30a4\u30e4\u306eA"
	},
	"start_date": "2015-04-06",
	"end_date": "2016-03-28",
	"synopsis": "The story continues.",
	"mean": 8.13,
	"rank": 350,
	"popularity": 1200,
	"num_list_users": 80000,
	"num_scoring_users": 35000,
	"nsfw": "white",
	"media_type": "tv",
	"status": "finished_airing",
	"genres": [
		{"id": 1, "name": "Sports"},
		{"id": 27, "name": "Shounen"}
	],
	"num_episodes": 51,
	"start_season": {"year": 2015, "season": "spring"},
	"broadcast": {"day_of_the_week": "monday", "start_time": "18:00"},
	"source": "manga",
	"average_episode_duration": 1440,
	"rating": "pg_13",
	"pictures": [
		{"medium": "https://cdn.mal.net/p1.jpg", "large": "https://cdn.mal.net/p1l.jpg"}
	],
	"background": "Based on the manga.",
	"related_anime": [
		{
			"node": {"id": 18689, "title": "Diamond no Ace"},
			"relation_type": "prequel",
			"relation_type_formatted": "Prequel"
		}
	],
	"related_manga": [],
	"recommendations": [
		{
			"node": {"id": 15, "title": "Major S1"},
			"num_recommendations": 12
		}
	],
	"studios": [
		{"id": 11, "name": "Madhouse"},
		{"id": 10, "name": "Production I.G"}
	],
	"statistics": {
		"status": {
			"watching": "5000",
			"completed": "30000",
			"on_hold": "8000",
			"dropped": "2000",
			"plan_to_watch": "10000"
		},
		"num_list_users": 55000
	}
}`

func TestGetAnimeMetadata(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(animeJSON))
	})

	a, err := c.GetAnime(context.Background(), 30230, []string{"title", "synopsis"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID != 30230 {
		t.Errorf("ID = %d, want 30230", a.ID)
	}
	if a.Title != "Diamond no Ace: Second Season" {
		t.Errorf("Title = %q, want %q", a.Title, "Diamond no Ace: Second Season")
	}
	if a.Mean != 8.13 {
		t.Errorf("Mean = %v, want 8.13", a.Mean)
	}
	if a.Rank != 350 {
		t.Errorf("Rank = %d, want 350", a.Rank)
	}
	if a.MediaType != "tv" {
		t.Errorf("MediaType = %q, want %q", a.MediaType, "tv")
	}
}

func TestGetAnimeRelations(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(animeJSON))
	})

	a, err := c.GetAnime(context.Background(), 30230, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.RelatedAnime) != 1 {
		t.Fatalf("RelatedAnime len = %d, want 1", len(a.RelatedAnime))
	}
	if a.RelatedAnime[0].Node.ID != 18689 {
		t.Errorf("RelatedAnime[0].Node.ID = %d, want 18689", a.RelatedAnime[0].Node.ID)
	}
	if a.RelatedAnime[0].RelationType != "prequel" {
		t.Errorf("RelationType = %q, want %q", a.RelatedAnime[0].RelationType, "prequel")
	}
	if len(a.Recommendations) != 1 {
		t.Fatalf("Recommendations len = %d, want 1", len(a.Recommendations))
	}
	if a.Recommendations[0].NumRecommendations != 12 {
		t.Errorf("NumRecommendations = %d, want 12", a.Recommendations[0].NumRecommendations)
	}
}

func TestGetAnimeSubtypes(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(animeJSON))
	})

	a, err := c.GetAnime(context.Background(), 30230, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.Genres) != 2 {
		t.Fatalf("Genres len = %d, want 2", len(a.Genres))
	}
	if a.Genres[0].Name != "Sports" {
		t.Errorf("Genres[0].Name = %q, want %q", a.Genres[0].Name, "Sports")
	}
	if len(a.Studios) != 2 {
		t.Fatalf("Studios len = %d, want 2", len(a.Studios))
	}
	if a.Studios[0].Name != "Madhouse" {
		t.Errorf("Studios[0].Name = %q, want %q", a.Studios[0].Name, "Madhouse")
	}
	if a.StartSeason == nil || a.StartSeason.Year != 2015 {
		t.Errorf("StartSeason.Year want 2015, got %v", a.StartSeason)
	}
	if a.Broadcast == nil || a.Broadcast.DayOfTheWeek != "monday" {
		t.Errorf("Broadcast.DayOfTheWeek want monday, got %v", a.Broadcast)
	}
	if a.Statistics == nil || a.Statistics.NumListUsers != 55000 {
		t.Errorf("Statistics.NumListUsers want 55000, got %v", a.Statistics)
	}
}

func TestSearchAnime(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "one piece" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"node": {"id": 21, "title": "One Piece"}},
				{"node": {"id": 22, "title": "One Piece Film"}}
			],
			"paging": {"next": "https://api.myanimelist.net/v2/anime?offset=2"}
		}`))
	})

	anime, pg, err := c.SearchAnime(context.Background(), "one piece", nil, 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(anime) != 2 {
		t.Fatalf("len = %d, want 2", len(anime))
	}
	if anime[0].Title != "One Piece" {
		t.Errorf("anime[0].Title = %q, want %q", anime[0].Title, "One Piece")
	}
	if pg.Next == "" {
		t.Error("expected paging.Next to be non-empty")
	}
}

func TestAnimeRanking(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ranking_type") != "airing" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"node": {"id": 40028, "title": "Shingeki no Kyojin"}, "ranking": {"rank": 1}},
				{"node": {"id": 42203, "title": "Re:Zero S2P2"}, "ranking": {"rank": 2}}
			],
			"paging": {}
		}`))
	})

	ranked, pg, err := c.AnimeRanking(context.Background(), "airing", nil, 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ranked) != 2 {
		t.Fatalf("len = %d, want 2", len(ranked))
	}
	if ranked[0].Ranking.Rank != 1 {
		t.Errorf("ranked[0].Ranking.Rank = %d, want 1", ranked[0].Ranking.Rank)
	}
	if ranked[1].Anime.ID != 42203 {
		t.Errorf("ranked[1].Anime.ID = %d, want 42203", ranked[1].Anime.ID)
	}
	if pg == nil {
		t.Fatal("paging should not be nil")
	}
}

func TestSeasonalAnime(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/anime/season/2024/winter") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [{"node": {"id": 100, "title": "Winter Show"}}],
			"paging": {}
		}`))
	})

	anime, _, err := c.SeasonalAnime(context.Background(), 2024, "winter", nil, "anime_score", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(anime) != 1 {
		t.Fatalf("len = %d, want 1", len(anime))
	}
	if anime[0].Title != "Winter Show" {
		t.Errorf("Title = %q, want %q", anime[0].Title, "Winter Show")
	}
}

const mangaJSON = `{
	"id": 2,
	"title": "Berserk",
	"main_picture": {
		"medium": "https://cdn.mal.net/images/manga/2.jpg",
		"large": "https://cdn.mal.net/images/manga/2l.jpg"
	},
	"alternative_titles": {
		"synonyms": [],
		"en": "Berserk",
		"ja": "\u30d9\u30eb\u30bb\u30eb\u30af"
	},
	"start_date": "1989-08-25",
	"synopsis": "Guts searches for revenge.",
	"mean": 9.47,
	"rank": 1,
	"popularity": 2,
	"media_type": "manga",
	"status": "currently_publishing",
	"genres": [
		{"id": 1, "name": "Action"},
		{"id": 8, "name": "Drama"},
		{"id": 14, "name": "Horror"}
	],
	"num_volumes": 41,
	"num_chapters": 380,
	"authors": [
		{
			"node": {"id": 1868, "first_name": "Kentarou", "last_name": "Miura"},
			"role": "Story & Art"
		}
	],
	"serialization": [
		{"node": {"id": 2, "name": "Young Animal"}, "role": ""}
	]
}`

func TestGetManga(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mangaJSON))
	})

	m, err := c.GetManga(context.Background(), 2, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.ID != 2 {
		t.Errorf("ID = %d, want 2", m.ID)
	}
	if m.Title != "Berserk" {
		t.Errorf("Title = %q, want %q", m.Title, "Berserk")
	}
	if m.Mean != 9.47 {
		t.Errorf("Mean = %v, want 9.47", m.Mean)
	}
	if m.NumVolumes != 41 {
		t.Errorf("NumVolumes = %d, want 41", m.NumVolumes)
	}
	if m.NumChapters != 380 {
		t.Errorf("NumChapters = %d, want 380", m.NumChapters)
	}
}

func TestGetMangaAuthors(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mangaJSON))
	})

	m, err := c.GetManga(context.Background(), 2, []string{"authors"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Authors) != 1 {
		t.Fatalf("Authors len = %d, want 1", len(m.Authors))
	}
	if m.Authors[0].Node.LastName != "Miura" {
		t.Errorf("Author LastName = %q, want %q", m.Authors[0].Node.LastName, "Miura")
	}
	if m.Authors[0].Role != "Story & Art" {
		t.Errorf("Author Role = %q, want %q", m.Authors[0].Role, "Story & Art")
	}
	if len(m.Serialization) != 1 {
		t.Fatalf("Serialization len = %d, want 1", len(m.Serialization))
	}
	if m.Serialization[0].Node.Name != "Young Animal" {
		t.Errorf("Serialization = %q, want %q", m.Serialization[0].Node.Name, "Young Animal")
	}
}

func TestSearchManga(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [{"node": {"id": 2, "title": "Berserk"}}],
			"paging": {"previous": "https://api.myanimelist.net/v2/manga?offset=0"}
		}`))
	})

	manga, pg, err := c.SearchManga(context.Background(), "berserk", nil, 1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(manga) != 1 {
		t.Fatalf("len = %d, want 1", len(manga))
	}
	if manga[0].Title != "Berserk" {
		t.Errorf("Title = %q, want %q", manga[0].Title, "Berserk")
	}
	if pg.Previous == "" {
		t.Error("expected paging.Previous to be non-empty")
	}
}

func TestMangaRanking(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"node": {"id": 2, "title": "Berserk"}, "ranking": {"rank": 1}},
				{"node": {"id": 13, "title": "One Piece"}, "ranking": {"rank": 2}}
			],
			"paging": {}
		}`))
	})

	ranked, _, err := c.MangaRanking(context.Background(), "all", []string{"rank"}, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ranked) != 2 {
		t.Fatalf("len = %d, want 2", len(ranked))
	}
	if ranked[0].Manga.Title != "Berserk" {
		t.Errorf("ranked[0].Manga.Title = %q, want %q", ranked[0].Manga.Title, "Berserk")
	}
	if ranked[1].Ranking.Rank != 2 {
		t.Errorf("ranked[1].Ranking.Rank = %d, want 2", ranked[1].Ranking.Rank)
	}
}

func TestForumBoards(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"categories": [
				{
					"title": "MyAnimeList",
					"boards": [
						{
							"id": 5,
							"title": "Updates & Announcements",
							"description": "Official site news",
							"subboards": [{"id": 2, "title": "MAL Guidelines"}]
						}
					]
				}
			]
		}`))
	})

	cats, err := c.ForumBoards(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 1 {
		t.Fatalf("categories len = %d, want 1", len(cats))
	}
	if cats[0].Title != "MyAnimeList" {
		t.Errorf("category title = %q, want %q", cats[0].Title, "MyAnimeList")
	}
	if len(cats[0].Boards) != 1 {
		t.Fatalf("boards len = %d, want 1", len(cats[0].Boards))
	}
	if cats[0].Boards[0].ID != 5 {
		t.Errorf("board ID = %d, want 5", cats[0].Boards[0].ID)
	}
	if len(cats[0].Boards[0].Subboards) != 1 {
		t.Fatalf("subboards len = %d, want 1", len(cats[0].Boards[0].Subboards))
	}
	if cats[0].Boards[0].Subboards[0].Title != "MAL Guidelines" {
		t.Errorf("subboard title = %q, want %q", cats[0].Boards[0].Subboards[0].Title, "MAL Guidelines")
	}
}

func TestForumTopicDetail(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/forum/topic/481") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"title": "Best anime of 2024",
				"posts": [
					{
						"id": 1001,
						"number": 1,
						"created_at": "2024-01-15T10:30:00+00:00",
						"created_by": {"id": 42, "name": "user42"},
						"body": "I think it is...",
						"signature": ""
					}
				],
				"poll": {
					"id": 50,
					"question": "What is the best?",
					"closed": false,
					"options": [
						{"id": 1, "text": "Show A", "votes": 100},
						{"id": 2, "text": "Show B", "votes": 80}
					]
				}
			},
			"paging": {"next": "https://api.myanimelist.net/v2/forum/topic/481?offset=1"}
		}`))
	})

	detail, pg, err := c.ForumTopicDetail(context.Background(), 481, 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.Title != "Best anime of 2024" {
		t.Errorf("Title = %q, want %q", detail.Title, "Best anime of 2024")
	}
	if len(detail.Posts) != 1 {
		t.Fatalf("Posts len = %d, want 1", len(detail.Posts))
	}
	if detail.Posts[0].CreatedBy.Name != "user42" {
		t.Errorf("Post CreatedBy = %q, want %q", detail.Posts[0].CreatedBy.Name, "user42")
	}
	if detail.Poll == nil {
		t.Fatal("expected poll, got nil")
	}
	if len(detail.Poll.Options) != 2 {
		t.Fatalf("Poll options len = %d, want 2", len(detail.Poll.Options))
	}
	if detail.Poll.Options[0].Votes != 100 {
		t.Errorf("Option[0].Votes = %d, want 100", detail.Poll.Options[0].Votes)
	}
	if pg.Next == "" {
		t.Error("expected paging.Next to be non-empty")
	}
}

func TestForumTopics(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "love" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"id": 999,
					"title": "Love is War discussion",
					"created_at": "2024-03-01T08:00:00+00:00",
					"created_by": {"id": 7, "name": "fan7"},
					"number_of_posts": 42,
					"last_post_created_at": "2024-03-02T12:00:00+00:00",
					"last_post_created_by": {"id": 8, "name": "fan8"},
					"is_locked": false
				}
			],
			"paging": {}
		}`))
	})

	topics, _, err := c.ForumTopics(context.Background(), "love", 0, 2, 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(topics) != 1 {
		t.Fatalf("len = %d, want 1", len(topics))
	}
	if topics[0].Title != "Love is War discussion" {
		t.Errorf("Title = %q, want %q", topics[0].Title, "Love is War discussion")
	}
	if topics[0].NumberOfPosts != 42 {
		t.Errorf("NumberOfPosts = %d, want 42", topics[0].NumberOfPosts)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "bad_request", "message": "Invalid parameters"}`))
	})

	_, err := c.GetAnime(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *mal.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *mal.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if apiErr.Err != "bad_request" {
		t.Errorf("Err = %q, want %q", apiErr.Err, "bad_request")
	}
	if apiErr.Message != "Invalid parameters" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "Invalid parameters")
	}
}

func TestAPIErrorNoMessage(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error": "forbidden"}`))
	})

	_, err := c.GetAnime(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *mal.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *mal.APIError, got %T", err)
	}
	if !strings.Contains(apiErr.Error(), "forbidden") {
		t.Errorf("error string %q should contain %q", apiErr.Error(), "forbidden")
	}
	// Message should not appear when empty.
	if strings.Contains(apiErr.Error(), ": :") {
		t.Errorf("error string should not double-colon: %q", apiErr.Error())
	}
}

func TestHTTPError(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := c.GetAnime(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *mal.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *mal.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestRequestHeaders(t *testing.T) {
	t.Parallel()
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Mal-Client-Id") != "test-client-id" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": 1, "title": "Test"}`))
	})

	a, err := c.GetAnime(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Title != "Test" {
		t.Errorf("Title = %q, want %q", a.Title, "Test")
	}
}

func TestWithUserAgent(t *testing.T) {
	t.Parallel()
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": 1}`))
	}))
	t.Cleanup(srv.Close)

	c := mal.New("id", metadata.WithBaseURL(srv.URL), metadata.WithUserAgent("custom/1.0"))
	_, err := c.GetAnime(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotUA != "custom/1.0" {
		t.Errorf("User-Agent = %q, want %q", gotUA, "custom/1.0")
	}
}

func TestFieldsParam(t *testing.T) {
	t.Parallel()
	var gotFields string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotFields = r.URL.Query().Get("fields")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": 1}`))
	})

	_, err := c.GetAnime(context.Background(), 1, []string{"synopsis", "genres", "studios"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFields != "synopsis,genres,studios" {
		t.Errorf("fields = %q, want %q", gotFields, "synopsis,genres,studios")
	}
}

// OAuth2 tests.

func TestGeneratePKCE(t *testing.T) {
	t.Parallel()
	pkce, err := mal.GeneratePKCE()
	if err != nil {
		t.Fatal(err)
	}
	if pkce.CodeVerifier == "" {
		t.Fatal("CodeVerifier is empty")
	}
	if pkce.CodeChallenge == "" {
		t.Fatal("CodeChallenge is empty")
	}
	// PKCE verifier should be base64url-encoded 64 bytes ≈ 86 chars.
	if len(pkce.CodeVerifier) < 43 {
		t.Errorf("CodeVerifier too short: %d", len(pkce.CodeVerifier))
	}
	// Two calls should produce different verifiers.
	pkce2, _ := mal.GeneratePKCE()
	if pkce.CodeVerifier == pkce2.CodeVerifier {
		t.Error("two PKCE calls produced the same verifier")
	}
}

func TestAuthorizationURL(t *testing.T) {
	t.Parallel()
	authSrv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	t.Cleanup(authSrv.Close)

	c := mal.New("my-client-id")
	c.SetAuthURL(authSrv.URL)
	pkce := &mal.PKCEChallenge{
		CodeVerifier:  "test-verifier",
		CodeChallenge: "test-challenge",
	}
	u := c.AuthorizationURL("mystate", pkce)
	if !strings.Contains(u, "client_id=my-client-id") {
		t.Errorf("URL missing client_id: %s", u)
	}
	if !strings.Contains(u, "code_challenge=test-challenge") {
		t.Errorf("URL missing code_challenge: %s", u)
	}
	if !strings.Contains(u, "code_challenge_method=S256") {
		t.Errorf("URL missing code_challenge_method: %s", u)
	}
	if !strings.Contains(u, "state=mystate") {
		t.Errorf("URL missing state: %s", u)
	}
	if !strings.Contains(u, "response_type=code") {
		t.Errorf("URL missing response_type: %s", u)
	}
}

func TestExchangeCode(t *testing.T) {
	t.Parallel()
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/token" {
			t.Errorf("path = %q, want /token", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.PostForm.Get("grant_type") != "authorization_code" {
			t.Errorf("grant_type = %q", r.PostForm.Get("grant_type"))
		}
		if r.PostForm.Get("code") != "auth-code-123" {
			t.Errorf("code = %q", r.PostForm.Get("code"))
		}
		if r.PostForm.Get("code_verifier") != "test-verifier" {
			t.Errorf("code_verifier = %q", r.PostForm.Get("code_verifier"))
		}
		if r.PostForm.Get("client_id") != "cid" {
			t.Errorf("client_id = %q", r.PostForm.Get("client_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"mal-access","token_type":"Bearer","expires_in":2592000,"refresh_token":"mal-refresh"}`))
	}))
	t.Cleanup(authSrv.Close)

	var callbackToken mal.Token
	c := mal.New("cid")
	c.SetAuthURL(authSrv.URL)
	c.SetTokenCallback(func(tok mal.Token) { callbackToken = tok })
	pkce := &mal.PKCEChallenge{CodeVerifier: "test-verifier", CodeChallenge: "test-challenge"}
	tok, err := c.ExchangeCode(context.Background(), "auth-code-123", pkce, "")
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "mal-access" {
		t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "mal-access")
	}
	if tok.RefreshToken != "mal-refresh" {
		t.Errorf("RefreshToken = %q, want %q", tok.RefreshToken, "mal-refresh")
	}
	if callbackToken.AccessToken != "mal-access" {
		t.Errorf("callback AccessToken = %q, want %q", callbackToken.AccessToken, "mal-access")
	}
}

func TestRefreshToken(t *testing.T) {
	t.Parallel()
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.PostForm.Get("grant_type") != "refresh_token" {
			t.Errorf("grant_type = %q, want refresh_token", r.PostForm.Get("grant_type"))
		}
		if r.PostForm.Get("refresh_token") != "old-rt" {
			t.Errorf("refresh_token = %q, want old-rt", r.PostForm.Get("refresh_token"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh"}`))
	}))
	t.Cleanup(authSrv.Close)

	c := mal.New("cid")
	c.SetAuthURL(authSrv.URL)
	c.SetRefreshToken("old-rt")
	tok, err := c.RefreshToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "new-access")
	}
}

func TestRefreshTokenMissing(t *testing.T) {
	t.Parallel()
	c := mal.New("cid")
	_, err := c.RefreshToken(context.Background())
	if err == nil {
		t.Fatal("expected error when no refresh token set")
	}
}

func TestBearerTokenOverClientID(t *testing.T) {
	t.Parallel()
	var gotAuth, gotClientID string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotClientID = r.Header.Get("X-Mal-Client-Id")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": 1}`))
	}))
	t.Cleanup(srv.Close)

	c := mal.New("cid", metadata.WithBaseURL(srv.URL))
	c.SetAccessToken("my-tok")
	_, err := c.GetAnime(context.Background(), 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer my-tok" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer my-tok")
	}
	if gotClientID != "" {
		t.Errorf("X-MAL-CLIENT-ID = %q, want empty when token present", gotClientID)
	}
}

func TestClientIDFallback(t *testing.T) {
	t.Parallel()
	var gotAuth, gotClientID string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotClientID = r.Header.Get("X-Mal-Client-Id")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": 1}`))
	}))
	t.Cleanup(srv.Close)

	c := mal.New("cid", metadata.WithBaseURL(srv.URL))
	_, err := c.GetAnime(context.Background(), 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "" {
		t.Errorf("Authorization = %q, want empty when no token", gotAuth)
	}
	if gotClientID != "cid" {
		t.Errorf("X-MAL-CLIENT-ID = %q, want %q", gotClientID, "cid")
	}
}

func TestPaginationParams(t *testing.T) {
	t.Parallel()
	var gotLimit, gotOffset string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotLimit = r.URL.Query().Get("limit")
		gotOffset = r.URL.Query().Get("offset")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data": [], "paging": {}}`))
	})

	_, _, err := c.SearchAnime(context.Background(), "test", nil, 25, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotLimit != "25" {
		t.Errorf("limit = %q, want %q", gotLimit, "25")
	}
	if gotOffset != "50" {
		t.Errorf("offset = %q, want %q", gotOffset, "50")
	}
}
