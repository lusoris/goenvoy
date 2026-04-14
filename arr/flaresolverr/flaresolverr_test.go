package flaresolverr_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/arr/flaresolverr"
)

func newTestServer(t *testing.T, wantCmd string, response any) *flaresolverr.Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1" {
			t.Errorf("path = %q, want /v1", r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		defer func() { _ = r.Body.Close() }()

		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got, ok := req["cmd"].(string); !ok || got != wantCmd {
			t.Errorf("cmd = %q, want %q", got, wantCmd)
		}

		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
	t.Cleanup(ts.Close)
	return flaresolverr.New(ts.URL)
}

func TestGet(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "request.get", map[string]any{
		"status":  "ok",
		"message": "Challenge solved!",
		"solution": map[string]any{
			"url":       "https://example.com",
			"status":    200,
			"response":  "<html>OK</html>",
			"userAgent": "Mozilla/5.0",
		},
		"startTimestamp": 1700000000000,
		"endTimestamp":   1700000005000,
		"version":        "3.3.21",
	})

	resp, err := c.Get(context.Background(), "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" {
		t.Errorf("Status = %q, want ok", resp.Status)
	}
	if resp.Solution == nil {
		t.Fatal("Solution is nil")
	}
	if resp.Solution.URL != "https://example.com" {
		t.Errorf("Solution.URL = %q, want https://example.com", resp.Solution.URL)
	}
	if resp.Solution.Response != "<html>OK</html>" {
		t.Errorf("Solution.Response = %q, want <html>OK</html>", resp.Solution.Response)
	}
}

func TestGetWithOptions(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)

		if req["session"] != "my-session" {
			t.Errorf("session = %v, want my-session", req["session"])
		}
		if req["maxTimeout"] != float64(60000) {
			t.Errorf("maxTimeout = %v, want 60000", req["maxTimeout"])
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer ts.Close()

	c := flaresolverr.New(ts.URL)
	resp, err := c.Get(context.Background(), "https://example.com", &flaresolverr.RequestOptions{
		Session:    "my-session",
		MaxTimeout: 60000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" {
		t.Errorf("Status = %q, want ok", resp.Status)
	}
}

func TestPost(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)

		if req["cmd"] != "request.post" {
			t.Errorf("cmd = %v, want request.post", req["cmd"])
		}
		if req["postData"] != "user=test" {
			t.Errorf("postData = %v, want user=test", req["postData"])
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "solution": map[string]any{"url": "https://example.com/login"}})
	}))
	defer ts.Close()

	c := flaresolverr.New(ts.URL)
	resp, err := c.Post(context.Background(), "https://example.com/login", "user=test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" {
		t.Errorf("Status = %q, want ok", resp.Status)
	}
	if resp.Solution.URL != "https://example.com/login" {
		t.Errorf("Solution.URL = %q, want https://example.com/login", resp.Solution.URL)
	}
}

func TestCreateSession(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "sessions.create", map[string]any{"status": "ok", "message": "Session created"})

	resp, err := c.CreateSession(context.Background(), "test-session", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Session created" {
		t.Errorf("Message = %q, want Session created", resp.Message)
	}
}

func TestCreateSessionWithProxy(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)

		proxy, ok := req["proxy"].(map[string]any)
		if !ok {
			t.Fatal("missing proxy")
		}
		if proxy["url"] != "http://proxy:8080" {
			t.Errorf("proxy.url = %v, want http://proxy:8080", proxy["url"])
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer ts.Close()

	c := flaresolverr.New(ts.URL)
	_, err := c.CreateSession(context.Background(), "s1", &flaresolverr.Proxy{URL: "http://proxy:8080"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListSessions(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "sessions.list", map[string]any{"status": "ok", "message": "[]"})

	resp, err := c.ListSessions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "ok" {
		t.Errorf("Status = %q, want ok", resp.Status)
	}
}

func TestDestroySession(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "sessions.destroy", map[string]any{"status": "ok", "message": "Session destroyed"})

	resp, err := c.DestroySession(context.Background(), "test-session")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Session destroyed" {
		t.Errorf("Message = %q, want Session destroyed", resp.Message)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer ts.Close()

	c := flaresolverr.New(ts.URL)
	_, err := c.Get(context.Background(), "https://example.com", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *flaresolverr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	called := false
	custom := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			called = true
			return http.DefaultTransport.RoundTrip(r)
		}),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer ts.Close()

	c := flaresolverr.New(ts.URL, flaresolverr.WithHTTPClient(custom))
	_, _ = c.Get(context.Background(), "https://example.com", nil)
	if !called {
		t.Error("custom HTTP client was not used")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
