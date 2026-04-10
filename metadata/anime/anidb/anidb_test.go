package anidb

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lusoris/goenvoy/metadata"
)

func setup(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New("testclient", 1, metadata.WithBaseURL(ts.URL))
}

func serveXML(w http.ResponseWriter, data string) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	io.WriteString(w, data)
}

func getTestAnime(t *testing.T) *Anime {
	t.Helper()
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("aid") != "1" {
			http.Error(w, "bad aid", http.StatusBadRequest)
			return
		}
		serveXML(w, `<?xml version="1.0"?>
<anime id="1" restricted="false">
  <type>TV Series</type>
  <episodecount>13</episodecount>
  <startdate>1999-01-03</startdate>
  <enddate>1999-03-28</enddate>
  <titles>
    <title xml:lang="x-jat" type="main">Seikai no Monshou</title>
    <title xml:lang="en" type="official">Crest of the Stars</title>
  </titles>
  <relatedanime>
    <anime id="4" type="Sequel">Seikai no Senki</anime>
  </relatedanime>
  <similaranime>
    <anime id="584" approval="75" total="89">Ginga Eiyuu Densetsu</anime>
  </similaranime>
  <recommendations>
    <recommendation type="Must See" uid="125868">A must see anime.</recommendation>
  </recommendations>
  <url>http://www.example.com/</url>
  <creators>
    <name id="4303" type="Music">Hattori Katsuhisa</name>
  </creators>
  <description>Test description.</description>
  <ratings>
    <permanent count="4430">8.16</permanent>
    <temporary count="4460">8.23</temporary>
    <review count="12">8.70</review>
  </ratings>
  <picture>440.jpg</picture>
  <resources>
    <resource type="1">
      <externalentity>
        <identifier>14</identifier>
      </externalentity>
    </resource>
  </resources>
  <tags>
    <tag id="36" parentid="2607" weight="300" localspoiler="false" globalspoiler="false" verified="true" update="2018-01-21">
      <name>military</name>
      <description>Armed forces.</description>
    </tag>
  </tags>
  <characters>
    <character id="28" type="main character in" update="2012-07-25">
      <rating votes="1196">9.15</rating>
      <name>Lafiel</name>
      <gender>female</gender>
      <charactertype id="1">Character</charactertype>
      <description>Main protagonist.</description>
      <picture>14304.jpg</picture>
      <seiyuu id="12" picture="184301.jpg">Kawasumi Ayako</seiyuu>
    </character>
  </characters>
  <episodes>
    <episode id="1" update="2011-07-01">
      <epno type="1">1</epno>
      <length>25</length>
      <airdate>1999-01-03</airdate>
      <rating votes="28">3.31</rating>
      <title xml:lang="ja">`+"\u4fb5\u7565"+`</title>
      <title xml:lang="en">Invasion</title>
    </episode>
  </episodes>
</anime>`)
	})
	anime, err := c.GetAnime(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAnime() error: %v", err)
	}
	return anime
}

func TestGetAnimeMetadata(t *testing.T) {
	anime := getTestAnime(t)
	if anime.ID != 1 {
		t.Errorf("ID = %d, want 1", anime.ID)
	}
	if anime.Restricted {
		t.Error("expected Restricted to be false")
	}
	if anime.Type != "TV Series" {
		t.Errorf("Type = %q, want %q", anime.Type, "TV Series")
	}
	if anime.EpisodeCount != 13 {
		t.Errorf("EpisodeCount = %d, want 13", anime.EpisodeCount)
	}
	if anime.StartDate != "1999-01-03" {
		t.Errorf("StartDate = %q", anime.StartDate)
	}
	if anime.Description != "Test description." {
		t.Errorf("Description = %q", anime.Description)
	}
	if anime.Picture != "440.jpg" {
		t.Errorf("Picture = %q", anime.Picture)
	}
}

func TestGetAnimeTitles(t *testing.T) {
	anime := getTestAnime(t)
	if len(anime.Titles) != 2 {
		t.Fatalf("Titles count = %d, want 2", len(anime.Titles))
	}
	if anime.Titles[0].Lang != "x-jat" {
		t.Errorf("Title[0].Lang = %q, want %q", anime.Titles[0].Lang, "x-jat")
	}
	if anime.Titles[0].Type != "main" {
		t.Errorf("Title[0].Type = %q, want %q", anime.Titles[0].Type, "main")
	}
	if anime.Titles[0].Name != "Seikai no Monshou" {
		t.Errorf("Title[0].Name = %q, want %q", anime.Titles[0].Name, "Seikai no Monshou")
	}
}

func TestGetAnimeRelations(t *testing.T) {
	anime := getTestAnime(t)
	if len(anime.RelatedAnime) != 1 {
		t.Fatalf("RelatedAnime count = %d, want 1", len(anime.RelatedAnime))
	}
	if anime.RelatedAnime[0].ID != 4 {
		t.Errorf("RelatedAnime[0].ID = %d, want 4", anime.RelatedAnime[0].ID)
	}
	if anime.RelatedAnime[0].Type != "Sequel" {
		t.Errorf("RelatedAnime[0].Type = %q, want %q", anime.RelatedAnime[0].Type, "Sequel")
	}
	if len(anime.SimilarAnime) != 1 {
		t.Fatalf("SimilarAnime count = %d, want 1", len(anime.SimilarAnime))
	}
	if anime.SimilarAnime[0].Approval != 75 {
		t.Errorf("SimilarAnime[0].Approval = %d, want 75", anime.SimilarAnime[0].Approval)
	}
	if len(anime.Recommendations) != 1 {
		t.Fatalf("Recommendations count = %d, want 1", len(anime.Recommendations))
	}
	if anime.Recommendations[0].Type != "Must See" {
		t.Errorf("Recommendations[0].Type = %q, want %q", anime.Recommendations[0].Type, "Must See")
	}
}

func TestGetAnimeRatings(t *testing.T) {
	anime := getTestAnime(t)
	if anime.Ratings.Permanent.Value != "8.16" {
		t.Errorf("Permanent.Value = %q, want %q", anime.Ratings.Permanent.Value, "8.16")
	}
	if anime.Ratings.Permanent.Count != 4430 {
		t.Errorf("Permanent.Count = %d, want 4430", anime.Ratings.Permanent.Count)
	}
	if len(anime.Creators) != 1 {
		t.Fatalf("Creators count = %d, want 1", len(anime.Creators))
	}
	if anime.Creators[0].Name != "Hattori Katsuhisa" {
		t.Errorf("Creator name = %q", anime.Creators[0].Name)
	}
	if len(anime.Resources) != 1 {
		t.Fatalf("Resources count = %d, want 1", len(anime.Resources))
	}
	if len(anime.Resources[0].ExternalEntities) != 1 {
		t.Fatalf("ExternalEntities count = %d", len(anime.Resources[0].ExternalEntities))
	}
	if anime.Resources[0].ExternalEntities[0].Identifiers[0] != "14" {
		t.Errorf("Resource identifier = %q", anime.Resources[0].ExternalEntities[0].Identifiers[0])
	}
}

func TestGetAnimeTagsAndCharacters(t *testing.T) {
	anime := getTestAnime(t)
	if len(anime.Tags) != 1 {
		t.Fatalf("Tags count = %d, want 1", len(anime.Tags))
	}
	if anime.Tags[0].Name != "military" {
		t.Errorf("Tag name = %q, want %q", anime.Tags[0].Name, "military")
	}
	if anime.Tags[0].Weight != 300 {
		t.Errorf("Tag weight = %d, want 300", anime.Tags[0].Weight)
	}
	if len(anime.Characters) != 1 {
		t.Fatalf("Characters count = %d, want 1", len(anime.Characters))
	}
	if anime.Characters[0].Name != "Lafiel" {
		t.Errorf("Character name = %q, want %q", anime.Characters[0].Name, "Lafiel")
	}
	if anime.Characters[0].Rating.Votes != 1196 {
		t.Errorf("Character rating votes = %d, want 1196", anime.Characters[0].Rating.Votes)
	}
	if anime.Characters[0].Seiyuu == nil {
		t.Fatal("Seiyuu is nil")
	}
	if anime.Characters[0].Seiyuu.Name != "Kawasumi Ayako" {
		t.Errorf("Seiyuu name = %q, want %q", anime.Characters[0].Seiyuu.Name, "Kawasumi Ayako")
	}
}

func TestGetAnimeEpisodes(t *testing.T) {
	anime := getTestAnime(t)
	if len(anime.Episodes) != 1 {
		t.Fatalf("Episodes count = %d, want 1", len(anime.Episodes))
	}
	if anime.Episodes[0].EpNo.Value != "1" {
		t.Errorf("EpNo = %q, want %q", anime.Episodes[0].EpNo.Value, "1")
	}
	if anime.Episodes[0].Length != 25 {
		t.Errorf("Episode length = %d, want 25", anime.Episodes[0].Length)
	}
	if len(anime.Episodes[0].Titles) != 2 {
		t.Errorf("Episode titles count = %d, want 2", len(anime.Episodes[0].Titles))
	}
}

func TestHotAnime(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		serveXML(w, `<hotanime>
  <anime id="8556" restricted="false">
    <episodecount>12</episodecount>
    <startdate>2012-01-10</startdate>
    <title xml:lang="x-jat" type="main">Another</title>
    <ratings>
      <permanent count="248">6.61</permanent>
      <temporary count="261">7.88</temporary>
    </ratings>
    <picture>79963.jpg</picture>
  </anime>
  <anime id="1234" restricted="true">
    <episodecount>24</episodecount>
    <startdate>2011-04-03</startdate>
    <title xml:lang="en" type="official">Test Anime</title>
    <ratings>
      <permanent count="100">7.00</permanent>
    </ratings>
    <picture>test.jpg</picture>
  </anime>
</hotanime>`)
	})

	entries, err := c.HotAnime(context.Background())
	if err != nil {
		t.Fatalf("HotAnime() error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}
	if entries[0].ID != 8556 {
		t.Errorf("entries[0].ID = %d, want 8556", entries[0].ID)
	}
	if entries[0].Title.Name != "Another" {
		t.Errorf("title = %q, want %q", entries[0].Title.Name, "Another")
	}
	if entries[0].Ratings.Permanent.Count != 248 {
		t.Errorf("permanent count = %d, want 248", entries[0].Ratings.Permanent.Count)
	}
	if !entries[1].Restricted {
		t.Error("entries[1].Restricted should be true")
	}
}

func TestRandomRecommendation(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		serveXML(w, `<randomrecommendation>
  <recommendation>
    <anime id="7899" restricted="false">
      <type>TV Series</type>
      <episodecount>13</episodecount>
      <startdate>2010-10-04</startdate>
      <enddate>2010-12-27</enddate>
      <title xml:lang="x-jat" type="main">Arakawa Under the Bridge 2</title>
      <picture>54734.jpg</picture>
      <ratings>
        <permanent count="781">6.23</permanent>
        <recommendations>8</recommendations>
      </ratings>
    </anime>
  </recommendation>
</randomrecommendation>`)
	})

	entries, err := c.RandomRecommendation(context.Background())
	if err != nil {
		t.Fatalf("RandomRecommendation() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	if entries[0].Anime.ID != 7899 {
		t.Errorf("anime ID = %d, want 7899", entries[0].Anime.ID)
	}
	if entries[0].Anime.Type != "TV Series" {
		t.Errorf("anime Type = %q, want %q", entries[0].Anime.Type, "TV Series")
	}
	if entries[0].Anime.Ratings.Recommendations != "8" {
		t.Errorf("recommendations = %q, want %q", entries[0].Anime.Ratings.Recommendations, "8")
	}
	if entries[0].Anime.Picture != "54734.jpg" {
		t.Errorf("picture = %q, want %q", entries[0].Anime.Picture, "54734.jpg")
	}
}

func TestRandomSimilar(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		serveXML(w, `<randomsimilar>
  <similar>
    <source aid="7056" restricted="false">
      <title xml:lang="x-jat" type="main">Aoi Bungaku Series</title>
      <picture>37124.jpg</picture>
    </source>
    <target aid="3654" restricted="false">
      <title xml:lang="x-jat" type="main">Ayakashi</title>
      <picture>44085.jpg</picture>
    </target>
  </similar>
</randomsimilar>`)
	})

	pairs, err := c.RandomSimilar(context.Background())
	if err != nil {
		t.Fatalf("RandomSimilar() error: %v", err)
	}
	if len(pairs) != 1 {
		t.Fatalf("got %d pairs, want 1", len(pairs))
	}
	if pairs[0].Source.AID != 7056 {
		t.Errorf("source AID = %d, want 7056", pairs[0].Source.AID)
	}
	if pairs[0].Source.Title.Name != "Aoi Bungaku Series" {
		t.Errorf("source title = %q", pairs[0].Source.Title.Name)
	}
	if pairs[0].Target.AID != 3654 {
		t.Errorf("target AID = %d, want 3654", pairs[0].Target.AID)
	}
	if pairs[0].Target.Picture != "44085.jpg" {
		t.Errorf("target picture = %q", pairs[0].Target.Picture)
	}
}

func TestMainPageEndpoint(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		serveXML(w, `<main>
  <hotanime>
    <anime id="100" restricted="false">
      <episodecount>12</episodecount>
      <title xml:lang="en" type="main">Hot Show</title>
      <picture>hot.jpg</picture>
      <ratings><permanent count="50">7.50</permanent></ratings>
    </anime>
  </hotanime>
  <randomsimilar>
    <similar>
      <source aid="200" restricted="false">
        <title xml:lang="en" type="main">Source</title>
        <picture>src.jpg</picture>
      </source>
      <target aid="300" restricted="false">
        <title xml:lang="en" type="main">Target</title>
        <picture>tgt.jpg</picture>
      </target>
    </similar>
  </randomsimilar>
  <randomrecommendation>
    <recommendation>
      <anime id="400" restricted="false">
        <type>OVA</type>
        <title xml:lang="en" type="main">Recommended</title>
        <picture>rec.jpg</picture>
        <ratings><permanent count="10">8.00</permanent></ratings>
      </anime>
    </recommendation>
  </randomrecommendation>
</main>`)
	})

	page, err := c.MainPage(context.Background())
	if err != nil {
		t.Fatalf("MainPage() error: %v", err)
	}
	if len(page.HotAnime) != 1 || page.HotAnime[0].ID != 100 {
		t.Errorf("HotAnime unexpected: %+v", page.HotAnime)
	}
	if len(page.RandomSimilar) != 1 || page.RandomSimilar[0].Source.AID != 200 {
		t.Errorf("RandomSimilar unexpected: %+v", page.RandomSimilar)
	}
	if len(page.RandomRecommendation) != 1 || page.RandomRecommendation[0].Anime.ID != 400 {
		t.Errorf("RandomRecommendation unexpected: %+v", page.RandomRecommendation)
	}
	if page.RandomRecommendation[0].Anime.Type != "OVA" {
		t.Errorf("recommendation type = %q, want %q", page.RandomRecommendation[0].Anime.Type, "OVA")
	}
}

func TestAPIError(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		serveXML(w, `<error code="302">client version missing or invalid</error>`)
	})

	_, err := c.GetAnime(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
	if apiErr.Code != "302" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "302")
	}
	if !strings.Contains(apiErr.Message, "client version") {
		t.Errorf("Message = %q, expected to contain 'client version'", apiErr.Message)
	}
}

func TestHTTPError(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := c.HotAnime(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
	if apiErr.Code != "500" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "500")
	}
}

func TestRequestParams(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("client") != "testclient" {
			t.Errorf("client = %q, want %q", q.Get("client"), "testclient")
		}
		if q.Get("clientver") != "1" {
			t.Errorf("clientver = %q, want %q", q.Get("clientver"), "1")
		}
		if q.Get("protover") != "1" {
			t.Errorf("protover = %q, want %q", q.Get("protover"), "1")
		}
		if q.Get("request") != "anime" {
			t.Errorf("request = %q, want %q", q.Get("request"), "anime")
		}
		if q.Get("aid") != "42" {
			t.Errorf("aid = %q, want %q", q.Get("aid"), "42")
		}
		serveXML(w, `<anime id="42"><type>TV</type></anime>`)
	})

	_, err := c.GetAnime(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		serveXML(w, `<hotanime></hotanime>`)
	}))
	t.Cleanup(ts.Close)

	c := New("testclient", 2, metadata.WithBaseURL(ts.URL), metadata.WithUserAgent("custom/1.0"))
	_, _ = c.HotAnime(context.Background())
	if gotUA != "custom/1.0" {
		t.Errorf("User-Agent = %q, want %q", gotUA, "custom/1.0")
	}
}

func TestBannedError(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		serveXML(w, `<error>Banned</error>`)
	})

	_, err := c.RandomSimilar(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
	if !strings.Contains(apiErr.Message, "Banned") {
		t.Errorf("Message = %q, expected 'Banned'", apiErr.Message)
	}
}
