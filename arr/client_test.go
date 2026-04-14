package arr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/arr/v2"
)

func TestNewBaseClient(t *testing.T) {
	t.Parallel()

	t.Run("valid URL", func(t *testing.T) {
		t.Parallel()
		c, err := arr.NewBaseClient("http://localhost:8989", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("trailing slash stripped", func(t *testing.T) {
		t.Parallel()
		_, err := arr.NewBaseClient("http://localhost:8989/", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		custom := &http.Client{}
		_, err := arr.NewBaseClient(
			"http://localhost:8989",
			"test-key",
			arr.WithHTTPClient(custom),
			arr.WithUserAgent("custom/1.0"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestBaseClient_Get(t *testing.T) {
	t.Parallel()

	want := arr.StatusResponse{
		AppName: "Sonarr",
		Version: "4.0.0",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("X-Api-Key") != "test-key" {
			t.Errorf("missing or wrong API key header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := arr.NewBaseClient(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got arr.StatusResponse
	if err := c.Get(context.Background(), "/api/v3/system/status", &got); err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got.AppName != want.AppName {
		t.Errorf("AppName = %q, want %q", got.AppName, want.AppName)
	}
	if got.Version != want.Version {
		t.Errorf("Version = %q, want %q", got.Version, want.Version)
	}
}

func TestBaseClient_Post(t *testing.T) {
	t.Parallel()

	wantCmd := arr.CommandResponse{
		ID:   1,
		Name: "RefreshSeries",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}

		var cmd arr.CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if cmd.Name != "RefreshSeries" {
			t.Errorf("command name = %q, want RefreshSeries", cmd.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wantCmd)
	}))
	defer srv.Close()

	c, err := arr.NewBaseClient(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got arr.CommandResponse
	err = c.Post(
		context.Background(),
		"/api/v3/command",
		arr.CommandRequest{Name: "RefreshSeries"},
		&got,
	)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}

	if got.Name != wantCmd.Name {
		t.Errorf("Name = %q, want %q", got.Name, wantCmd.Name)
	}
}

func TestBaseClient_ErrorResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer srv.Close()

	c, err := arr.NewBaseClient(srv.URL, "bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var dst arr.StatusResponse
	err = c.Get(context.Background(), "/api/v3/system/status", &dst)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	var apiErr *arr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *arr.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestBaseClient_Delete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, err := arr.NewBaseClient(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.Delete(context.Background(), "/api/v3/series/1", nil, nil); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
