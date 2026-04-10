package stashbox_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/adult/stashbox"
	"github.com/lusoris/goenvoy/metadata"
)

type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

func newGQLServer(t *testing.T, wantKey, dataField string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("ApiKey"); got != wantKey {
			t.Errorf("ApiKey = %q, want %q", got, wantKey)
		}
		var req gqlRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		data, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{` + `"` + dataField + `":` + string(data) + `}}`))
	}))
}

func TestFindPerformer(t *testing.T) {
	p := stashbox.Performer{
		ID:      "abc-123",
		Name:    "Test Performer",
		Gender:  "FEMALE",
		Country: "US",
		Aliases: []string{"Alias1"},
	}
	ts := newGQLServer(t, "key-1", "findPerformer", p)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-1")
	result, err := c.FindPerformer(context.Background(), "abc-123")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Test Performer" {
		t.Errorf("Name = %q, want %q", result.Name, "Test Performer")
	}
	if result.Gender != "FEMALE" {
		t.Errorf("Gender = %q, want %q", result.Gender, "FEMALE")
	}
}

func TestQueryPerformers(t *testing.T) {
	resp := struct {
		Count      int                  `json:"count"`
		Performers []stashbox.Performer `json:"performers"`
	}{
		Count:      1,
		Performers: []stashbox.Performer{{ID: "p1", Name: "Performer 1"}},
	}
	ts := newGQLServer(t, "key-2", "queryPerformers", resp)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-2")
	performers, count, err := c.QueryPerformers(context.Background(), &stashbox.QueryInput{Page: 1, PerPage: 25})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
	if len(performers) != 1 {
		t.Fatalf("len = %d, want 1", len(performers))
	}
}

func TestSearchPerformers(t *testing.T) {
	performers := []stashbox.Performer{{ID: "sp1", Name: "Searched Performer"}}
	ts := newGQLServer(t, "key-3", "searchPerformer", performers)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-3")
	result, err := c.SearchPerformers(context.Background(), "test", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Name != "Searched Performer" {
		t.Errorf("Name = %q, want %q", result[0].Name, "Searched Performer")
	}
}

func TestFindScene(t *testing.T) {
	s := stashbox.Scene{
		ID:       "scene-1",
		Title:    "Test Scene",
		Duration: intPtr(3600),
		Studio:   &stashbox.Studio{ID: "studio-1", Name: "Test Studio"},
		Tags:     []stashbox.Tag{{ID: "tag-1", Name: "Tag1"}},
	}
	ts := newGQLServer(t, "key-4", "findScene", s)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-4")
	result, err := c.FindScene(context.Background(), "scene-1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Title != "Test Scene" {
		t.Errorf("Title = %q, want %q", result.Title, "Test Scene")
	}
	if result.Studio == nil || result.Studio.Name != "Test Studio" {
		t.Error("Studio not parsed correctly")
	}
}

func TestQueryScenes(t *testing.T) {
	resp := struct {
		Count  int              `json:"count"`
		Scenes []stashbox.Scene `json:"scenes"`
	}{
		Count:  2,
		Scenes: []stashbox.Scene{{ID: "s1", Title: "Scene 1"}, {ID: "s2", Title: "Scene 2"}},
	}
	ts := newGQLServer(t, "key-5", "queryScenes", resp)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-5")
	scenes, count, err := c.QueryScenes(context.Background(), &stashbox.QueryInput{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	if len(scenes) != 2 {
		t.Fatalf("len = %d, want 2", len(scenes))
	}
}

func TestSearchScenes(t *testing.T) {
	scenes := []stashbox.Scene{{ID: "ss1", Title: "Found Scene"}}
	ts := newGQLServer(t, "key-6", "searchScene", scenes)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-6")
	result, err := c.SearchScenes(context.Background(), "test query")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestFindScenesByFingerprints(t *testing.T) {
	scenes := [][]stashbox.Scene{
		{{ID: "fp-1", Title: "Fingerprint Match"}},
	}
	ts := newGQLServer(t, "key-7", "findScenesBySceneFingerprints", scenes)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-7")
	fp := [][]stashbox.FingerprintInput{{{Hash: "abc123", Algorithm: "MD5"}}}
	result, err := c.FindScenesByFingerprints(context.Background(), fp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || len(result[0]) != 1 {
		t.Fatalf("unexpected result shape")
	}
	if result[0][0].Title != "Fingerprint Match" {
		t.Errorf("Title = %q, want %q", result[0][0].Title, "Fingerprint Match")
	}
}

func TestFindStudio(t *testing.T) {
	studio := stashbox.Studio{
		ID:     "studio-1",
		Name:   "Test Studio",
		Parent: &stashbox.Studio{ID: "parent-1", Name: "Parent Studio"},
	}
	ts := newGQLServer(t, "key-8", "findStudio", studio)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-8")
	result, err := c.FindStudio(context.Background(), "studio-1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Test Studio" {
		t.Errorf("Name = %q, want %q", result.Name, "Test Studio")
	}
	if result.Parent == nil || result.Parent.Name != "Parent Studio" {
		t.Error("Parent not parsed correctly")
	}
}

func TestQueryStudios(t *testing.T) {
	resp := struct {
		Count   int               `json:"count"`
		Studios []stashbox.Studio `json:"studios"`
	}{
		Count:   1,
		Studios: []stashbox.Studio{{ID: "qs1", Name: "Queried Studio"}},
	}
	ts := newGQLServer(t, "key-9", "queryStudios", resp)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-9")
	studios, count, err := c.QueryStudios(context.Background(), &stashbox.QueryInput{Page: 1, PerPage: 25})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
	if len(studios) != 1 {
		t.Fatalf("len = %d, want 1", len(studios))
	}
}

func TestSearchStudios(t *testing.T) {
	studios := []stashbox.Studio{{ID: "ss1", Name: "Searched Studio"}}
	ts := newGQLServer(t, "key-10", "searchStudio", studios)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-10")
	result, err := c.SearchStudios(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestFindTag(t *testing.T) {
	tag := stashbox.Tag{
		ID:          "tag-1",
		Name:        "Test Tag",
		Description: "A tag for testing",
		Category:    &stashbox.TagCategory{ID: "cat-1", Name: "Category"},
	}
	ts := newGQLServer(t, "key-11", "findTag", tag)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-11")
	result, err := c.FindTag(context.Background(), "tag-1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Test Tag" {
		t.Errorf("Name = %q, want %q", result.Name, "Test Tag")
	}
	if result.Category == nil || result.Category.Name != "Category" {
		t.Error("Category not parsed correctly")
	}
}

func TestQueryTags(t *testing.T) {
	resp := struct {
		Count int            `json:"count"`
		Tags  []stashbox.Tag `json:"tags"`
	}{
		Count: 3,
		Tags:  []stashbox.Tag{{ID: "t1", Name: "Tag1"}, {ID: "t2", Name: "Tag2"}, {ID: "t3", Name: "Tag3"}},
	}
	ts := newGQLServer(t, "key-12", "queryTags", resp)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-12")
	tags, count, err := c.QueryTags(context.Background(), &stashbox.QueryInput{Page: 1, PerPage: 50})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
	if len(tags) != 3 {
		t.Fatalf("len = %d, want 3", len(tags))
	}
}

func TestSearchTags(t *testing.T) {
	tags := []stashbox.Tag{{ID: "st1", Name: "Searched Tag"}}
	ts := newGQLServer(t, "key-13", "searchTag", tags)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-13")
	result, err := c.SearchTags(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestListSites(t *testing.T) {
	sites := []stashbox.Site{
		{ID: "site-1", Name: "Test Site", URL: "https://example.com"},
		{ID: "site-2", Name: "Other Site", URL: "https://other.com", ValidTypes: []string{"SCENE"}},
	}
	ts := newGQLServer(t, "key-14", "querySites", sites)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-14")
	result, err := c.ListSites(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0].Name != "Test Site" {
		t.Errorf("Name = %q, want %q", result[0].Name, "Test Site")
	}
}

func TestGetConfig(t *testing.T) {
	cfg := stashbox.Config{
		Host:          "https://stashdb.org",
		RequireInvite: true,
		GuidelinesURL: "https://guidelines.stashdb.org",
	}
	ts := newGQLServer(t, "key-15", "getConfig", cfg)
	defer ts.Close()

	c := stashbox.New(ts.URL, "key-15")
	result, err := c.GetConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Host != "https://stashdb.org" {
		t.Errorf("Host = %q, want %q", result.Host, "https://stashdb.org")
	}
	if !result.RequireInvite {
		t.Error("RequireInvite should be true")
	}
}

func TestGraphQLError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"performer not found"}]}`))
	}))
	defer ts.Close()

	c := stashbox.New(ts.URL, "key")
	_, err := c.FindPerformer(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	var gqlErr *stashbox.GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("expected *GraphQLError, got %T", err)
	}
	if gqlErr.Message != "performer not found" {
		t.Errorf("Message = %q, want %q", gqlErr.Message, "performer not found")
	}
}

func TestAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized"))
	}))
	defer ts.Close()

	c := stashbox.New(ts.URL, "bad-key")
	_, err := c.FindPerformer(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *stashbox.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestWithCustomHTTPClient(t *testing.T) {
	p := stashbox.Performer{ID: "p1", Name: "Custom"}
	ts := newGQLServer(t, "custom-key", "findPerformer", p)
	defer ts.Close()

	custom := &http.Client{}
	c := stashbox.New(ts.URL, "custom-key", metadata.WithHTTPClient(custom))
	result, err := c.FindPerformer(context.Background(), "p1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Custom" {
		t.Errorf("Name = %q, want %q", result.Name, "Custom")
	}
}

func intPtr(v int) *int { return &v }
