package gotify_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/notification/gotify"
)

func newServer(t *testing.T, wantMethod, wantPath string, resp any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != wantMethod {
			t.Errorf("method = %q, want %q", r.Method, wantMethod)
		}
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
}

func TestGotifyKeyHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Gotify-Key"); got != "test-token" {
			t.Errorf("X-Gotify-Key = %q, want test-token", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(gotify.Health{Health: "green"})
	}))
	defer ts.Close()

	c := gotify.New(ts.URL, "test-token")
	_, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateMessage(t *testing.T) {
	ts := newServer(t, http.MethodPost, "/message", gotify.Message{
		Id: 1, Title: "Test", Message: "Hello", Priority: 5,
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	msg, err := c.CreateMessage(context.Background(), "Test", "Hello", 5)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Title != "Test" {
		t.Errorf("Title = %q, want Test", msg.Title)
	}
	if msg.Priority != 5 {
		t.Errorf("Priority = %d, want 5", msg.Priority)
	}
}

func TestGetMessages(t *testing.T) {
	ts := newServer(t, http.MethodGet, "/message", gotify.PagedMessages{
		Messages: []gotify.Message{{Id: 1, Title: "Msg1"}},
		Paging:   &gotify.Paging{Size: 1, Limit: 100},
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	out, err := c.GetMessages(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Messages) != 1 {
		t.Fatalf("len = %d, want 1", len(out.Messages))
	}
	if out.Messages[0].Title != "Msg1" {
		t.Errorf("Title = %q, want Msg1", out.Messages[0].Title)
	}
}

func TestDeleteMessages(t *testing.T) {
	ts := newServer(t, http.MethodDelete, "/message", nil)
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	if err := c.DeleteMessages(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteMessage(t *testing.T) {
	ts := newServer(t, http.MethodDelete, "/message/42", nil)
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	if err := c.DeleteMessage(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
}

func TestGetApplications(t *testing.T) {
	ts := newServer(t, http.MethodGet, "/application", []gotify.Application{
		{Id: 1, Name: "MyApp", Token: "app-tok"},
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	apps, err := c.GetApplications(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 1 {
		t.Fatalf("len = %d, want 1", len(apps))
	}
	if apps[0].Name != "MyApp" {
		t.Errorf("Name = %q, want MyApp", apps[0].Name)
	}
}

func TestCreateApplication(t *testing.T) {
	ts := newServer(t, http.MethodPost, "/application", gotify.Application{
		Id: 2, Name: "NewApp", Description: "desc",
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	app, err := c.CreateApplication(context.Background(), "NewApp", "desc")
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "NewApp" {
		t.Errorf("Name = %q, want NewApp", app.Name)
	}
}

func TestDeleteApplication(t *testing.T) {
	ts := newServer(t, http.MethodDelete, "/application/2", nil)
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	if err := c.DeleteApplication(context.Background(), 2); err != nil {
		t.Fatal(err)
	}
}

func TestGetClients(t *testing.T) {
	ts := newServer(t, http.MethodGet, "/client", []gotify.ClientInfo{
		{Id: 1, Name: "Desktop", Token: "cli-tok"},
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	clients, err := c.GetClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(clients) != 1 {
		t.Fatalf("len = %d, want 1", len(clients))
	}
	if clients[0].Name != "Desktop" {
		t.Errorf("Name = %q, want Desktop", clients[0].Name)
	}
}

func TestCreateClient(t *testing.T) {
	ts := newServer(t, http.MethodPost, "/client", gotify.ClientInfo{
		Id: 3, Name: "Mobile",
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	cli, err := c.CreateClient(context.Background(), "Mobile")
	if err != nil {
		t.Fatal(err)
	}
	if cli.Name != "Mobile" {
		t.Errorf("Name = %q, want Mobile", cli.Name)
	}
}

func TestDeleteClient(t *testing.T) {
	ts := newServer(t, http.MethodDelete, "/client/3", nil)
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	if err := c.DeleteClient(context.Background(), 3); err != nil {
		t.Fatal(err)
	}
}

func TestGetCurrentUser(t *testing.T) {
	ts := newServer(t, http.MethodGet, "/current/user", gotify.User{
		Id: 1, Name: "admin", Admin: true,
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	u, err := c.GetCurrentUser(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "admin" {
		t.Errorf("Name = %q, want admin", u.Name)
	}
	if !u.Admin {
		t.Error("Admin = false, want true")
	}
}

func TestGetHealth(t *testing.T) {
	ts := newServer(t, http.MethodGet, "/health", gotify.Health{
		Health: "green", Database: "ok",
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	h, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if h.Health != "green" {
		t.Errorf("Health = %q, want green", h.Health)
	}
}

func TestGetVersion(t *testing.T) {
	ts := newServer(t, http.MethodGet, "/version", gotify.VersionInfo{
		Version: "2.4.0", Commit: "abc123", BuildDate: "2025-01-01",
	})
	defer ts.Close()

	c := gotify.New(ts.URL, "token")
	v, err := c.GetVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v.Version != "2.4.0" {
		t.Errorf("Version = %q, want 2.4.0", v.Version)
	}
}

func TestAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden"))
	}))
	defer ts.Close()

	c := gotify.New(ts.URL, "bad-token")
	_, err := c.GetHealth(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *gotify.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusForbidden)
	}
}

func TestWithHTTPClient(t *testing.T) {
	ts := newServer(t, http.MethodGet, "/health", gotify.Health{Health: "green"})
	defer ts.Close()

	custom := &http.Client{}
	c := gotify.New(ts.URL, "token", gotify.WithHTTPClient(custom))
	_, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
