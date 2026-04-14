package jellyfin_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/mediaserver/jellyfin/v2"
)

func newTestServer(t *testing.T, wantPath string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestAuthenticateByName(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/Users/AuthenticateByName" {
			t.Errorf("path = %q, want /Users/AuthenticateByName", r.URL.Path)
		}
		var body struct {
			Username string `json:"Username"`
			Pw       string `json:"Pw"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Username != "admin" {
			t.Errorf("username = %q, want admin", body.Username)
		}
		if body.Pw != "secret" {
			t.Errorf("password = %q, want secret", body.Pw)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jellyfin.AuthenticationResult{
			AccessToken: "test-token-123",
			ServerID:    "server-1",
		})
	}))
	defer ts.Close()

	c := jellyfin.New(ts.URL)
	if err := c.AuthenticateByName(context.Background(), "admin", "secret"); err != nil {
		t.Fatal(err)
	}
}

func TestAuthenticationError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Invalid username or password"})
	}))
	defer ts.Close()

	c := jellyfin.New(ts.URL)
	err := c.AuthenticateByName(context.Background(), "bad", "creds")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *jellyfin.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *jellyfin.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
}

func TestGetPublicSystemInfo(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/System/Info/Public", jellyfin.PublicSystemInfo{
		ServerName:   "My Jellyfin Server",
		Version:      "10.8.13",
		ID:           "abc123",
		LocalAddress: "192.168.1.100",
	})
	defer ts.Close()

	c := jellyfin.New(ts.URL)
	info, err := c.GetPublicSystemInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.ServerName != "My Jellyfin Server" {
		t.Errorf("ServerName = %q, want My Jellyfin Server", info.ServerName)
	}
	if info.Version != "10.8.13" {
		t.Errorf("Version = %q, want 10.8.13", info.Version)
	}
}

func TestGetSystemInfo(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/System/Info", jellyfin.SystemInfo{
		ServerName:                 "My Server",
		Version:                    "10.8.13",
		ID:                         "abc123",
		OperatingSystem:            "Linux",
		OperatingSystemDisplayName: "Debian GNU/Linux 12",
	})
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	info, err := c.GetSystemInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.OperatingSystem != "Linux" {
		t.Errorf("OS = %q, want Linux", info.OperatingSystem)
	}
}

func TestPing(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/System/Ping", "")
	defer ts.Close()

	c := jellyfin.New(ts.URL)
	if err := c.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/Users", []jellyfin.UserDto{
		{Name: "admin", ID: "user-1", HasPassword: true},
		{Name: "guest", ID: "user-2", HasPassword: false},
	})
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	users, err := c.GetUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("len = %d, want 2", len(users))
	}
	if users[0].Name != "admin" {
		t.Errorf("Name = %q, want admin", users[0].Name)
	}
}

func TestGetCurrentUser(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/Users/Me", jellyfin.UserDto{
		Name:        "admin",
		ID:          "user-1",
		ServerID:    "server-1",
		HasPassword: true,
	})
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	user, err := c.GetCurrentUser(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != "admin" {
		t.Errorf("Name = %q, want admin", user.Name)
	}
}

func TestGetSessions(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/Sessions", []jellyfin.SessionInfoDto{
		{
			ID:         "session-1",
			UserID:     "user-1",
			UserName:   "admin",
			DeviceName: "Chrome",
			Client:     "Jellyfin Web",
			IsActive:   true,
		},
	})
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	sessions, err := c.GetSessions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len = %d, want 1", len(sessions))
	}
	if sessions[0].UserName != "admin" {
		t.Errorf("UserName = %q, want admin", sessions[0].UserName)
	}
}

func TestGetItems(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/Items", jellyfin.ItemsResult{
		Items: []jellyfin.BaseItemDto{
			{ID: "item-1", Name: "Inception", Type: "Movie"},
			{ID: "item-2", Name: "Interstellar", Type: "Movie"},
		},
		TotalRecordCount: 2,
		StartIndex:       0,
	})
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	result, err := c.GetItems(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRecordCount != 2 {
		t.Errorf("TotalRecordCount = %d, want 2", result.TotalRecordCount)
	}
	if len(result.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(result.Items))
	}
	if result.Items[0].Name != "Inception" {
		t.Errorf("Name = %q, want Inception", result.Items[0].Name)
	}
}

func TestGetItemsByParent(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Items" {
			t.Errorf("path = %q, want /Items", r.URL.Path)
		}
		if got := r.URL.Query().Get("ParentId"); got != "lib-1" {
			t.Errorf("ParentId = %q, want lib-1", got)
		}
		if got := r.URL.Query().Get("StartIndex"); got != "0" {
			t.Errorf("StartIndex = %q, want 0", got)
		}
		if got := r.URL.Query().Get("Limit"); got != "50" {
			t.Errorf("Limit = %q, want 50", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jellyfin.ItemsResult{
			Items:            []jellyfin.BaseItemDto{{ID: "item-1", Name: "Movie 1"}},
			TotalRecordCount: 100,
			StartIndex:       0,
		})
	}))
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	result, err := c.GetItemsByParent(context.Background(), "lib-1", 0, 50)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRecordCount != 100 {
		t.Errorf("TotalRecordCount = %d, want 100", result.TotalRecordCount)
	}
}

func TestGetItem(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/Items/item-1", jellyfin.BaseItemDto{
		ID:           "item-1",
		Name:         "The Matrix",
		Type:         "Movie",
		RunTimeTicks: 8208000000,
		Genres:       []string{"Action", "Sci-Fi"},
	})
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	item, err := c.GetItem(context.Background(), "item-1")
	if err != nil {
		t.Fatal(err)
	}
	if item.Name != "The Matrix" {
		t.Errorf("Name = %q, want The Matrix", item.Name)
	}
	if len(item.Genres) != 2 {
		t.Errorf("len(Genres) = %d, want 2", len(item.Genres))
	}
}

func TestGetUserViews(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/UserViews" {
			t.Errorf("path = %q, want /UserViews", r.URL.Path)
		}
		if got := r.URL.Query().Get("userId"); got != "user-1" {
			t.Errorf("userId = %q, want user-1", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jellyfin.ItemsResult{
			Items: []jellyfin.BaseItemDto{
				{ID: "lib-1", Name: "Movies", Type: "CollectionFolder", IsFolder: true},
				{ID: "lib-2", Name: "TV Shows", Type: "CollectionFolder", IsFolder: true},
			},
			TotalRecordCount: 2,
		})
	}))
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	result, err := c.GetUserViews(context.Background(), "user-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("len = %d, want 2", len(result.Items))
	}
	if result.Items[0].Name != "Movies" {
		t.Errorf("Name = %q, want Movies", result.Items[0].Name)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("{\"message\":\"Access denied\"}"))
	}))
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("bad"))
	_, err := c.GetUsers(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *jellyfin.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *jellyfin.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want 403", apiErr.StatusCode)
	}
	if apiErr.Message != "Access denied" {
		t.Errorf("Message = %q, want Access denied", apiErr.Message)
	}
}

func TestAPIErrorMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  jellyfin.APIError
		want string
	}{
		{"with message", jellyfin.APIError{StatusCode: 401, Message: "Unauthorized"}, "jellyfin: HTTP 401: Unauthorized"},
		{"raw body", jellyfin.APIError{StatusCode: 502, RawBody: "gateway error"}, "jellyfin: HTTP 502: gateway error"},
		{"code only", jellyfin.APIError{StatusCode: 500}, "jellyfin: HTTP 500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/Users", []jellyfin.UserDto{})
	defer ts.Close()

	c := jellyfin.New(ts.URL, jellyfin.WithAccessToken("token"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GetUsers(ctx)
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
}

func TestWithOptions(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.Contains(auth, `DeviceId="my-device-123"`) {
			t.Errorf("Authorization = %q, want DeviceId my-device-123", auth)
		}
		if !strings.Contains(auth, `Token="token"`) {
			t.Errorf("Authorization = %q, want Token", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]jellyfin.UserDto{})
	}))
	defer ts.Close()

	c := jellyfin.New(ts.URL,
		jellyfin.WithAccessToken("token"),
		jellyfin.WithDeviceID("my-device-123"),
	)
	_, err := c.GetUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
