package opensubtitles

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Api-Key") != "test-api-key" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(ts.Close)

	return New("test-api-key", metadata.WithBaseURL(ts.URL))
}

func TestSearch(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(SearchResponse{
			TotalCount: 1,
			TotalPages: 1,
			Page:       1,
			PerPage:    60,
			Data: []Subtitle{{
				ID:   "1",
				Type: "subtitle",
				Attributes: SubtitleAttributes{
					Language: "en",
					Release:  "Inception.2010.BluRay",
				},
			}},
		})
	})

	resp, err := c.Search(context.Background(), &SearchParams{Query: "inception"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Attributes.Language != "en" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestDownload(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(DownloadResponse{
			Link:      "https://dl.opensubtitles.com/file/123",
			FileName:  "subtitle.srt",
			Remaining: 99,
		})
	})
	c.token = "test-token"

	resp, err := c.Download(context.Background(), DownloadRequest{FileID: 123})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Link == "" || resp.Remaining != 99 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(LoginResponse{
			Token:  "jwt-token-123",
			Status: 200,
		})
	})

	resp, err := c.Login(context.Background(), "user", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Token != "jwt-token-123" {
		t.Fatalf("unexpected token: %s", resp.Token)
	}
	if c.token != "jwt-token-123" {
		t.Fatal("token not stored on client")
	}
}

func TestSearchFeatures(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(FeaturesResponse{
			TotalCount: 1,
			Data: []Feature{{
				ID:   "f1",
				Type: "movie",
				Attributes: FeatureAttributes{
					Title: "Inception",
					Year:  "2010",
				},
			}},
		})
	})

	resp, err := c.SearchFeatures(context.Background(), "inception")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Attributes.Title != "Inception" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetLanguages(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []Language{{LanguageCode: "en", LanguageName: "English"}},
		})
	})

	langs, err := c.GetLanguages(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(langs) != 1 || langs[0].LanguageCode != "en" {
		t.Fatalf("unexpected languages: %+v", langs)
	}
}

func TestGetFormats(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []SubtitleFormat{{FormatName: "SubRip", Extension: "srt"}},
		})
	})

	fmts, err := c.GetFormats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(fmts) != 1 || fmts[0].Extension != "srt" {
		t.Fatalf("unexpected formats: %+v", fmts)
	}
}

func TestPopular(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(FeaturesResponse{
			TotalCount: 1,
			Data: []Feature{{
				ID:         "f1",
				Type:       "movie",
				Attributes: FeatureAttributes{Title: "Dune: Part Two"},
			}},
		})
	})

	resp, err := c.Popular(context.Background(), "en")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "rate limited", http.StatusTooManyRequests)
	})

	_, err := c.Search(context.Background(), &SearchParams{Query: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	custom := &http.Client{}
	c := New("key", metadata.WithHTTPClient(custom))
	if c.HTTPClient() != custom {
		t.Fatal("custom HTTP client not set")
	}
}

func TestWithUserAgent(t *testing.T) {
	t.Parallel()

	c := New("key", metadata.WithUserAgent("myapp v2.0"))
	if c.UserAgent() != "myapp v2.0" {
		t.Fatal("user agent not set")
	}
}
