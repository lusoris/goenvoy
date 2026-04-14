package rtorrent_test

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/downloadclient/rtorrent"
)

type xmlMethodCall struct {
	XMLName    xml.Name `xml:"methodCall"`
	MethodName string   `xml:"methodName"`
}

func xmlStringResponse(s string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<methodResponse><params><param><value><string>%s</string></value></param></params></methodResponse>`, s)
}

func xmlIntResponse(n int) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<methodResponse><params><param><value><i4>%d</i4></value></param></params></methodResponse>`, n)
}

func xmlVoidResponse() string {
	return `<?xml version="1.0"?>
<methodResponse><params><param><value><i4>0</i4></value></param></params></methodResponse>`
}

func xmlFaultResponse(msg string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<methodResponse><fault><value><string>%s</string></value></fault></methodResponse>`, msg)
}

func xmlMulticallResponse() string {
	return `<?xml version="1.0"?>
<methodResponse><params><param><value><array><data>
<value><array><data>
  <value><string>ABC123</string></value>
  <value><string>Ubuntu 24.04</string></value>
  <value><i8>4000000000</i8></value>
  <value><i8>4000000000</i8></value>
  <value><i8>0</i8></value>
  <value><i8>50000</i8></value>
  <value><i8>8000000000</i8></value>
  <value><i8>2000</i8></value>
  <value><i4>1</i4></value>
  <value><i4>1</i4></value>
  <value><i4>1</i4></value>
  <value><i4>0</i4></value>
  <value><string>/downloads/Ubuntu</string></value>
  <value><string>/downloads</string></value>
  <value><string>linux</string></value>
  <value><i8>1700000000</i8></value>
  <value><string></string></value>
</data></array></value>
<value><array><data>
  <value><string>DEF456</string></value>
  <value><string>Fedora 40</string></value>
  <value><i8>2000000000</i8></value>
  <value><i8>1000000000</i8></value>
  <value><i8>5000000</i8></value>
  <value><i8>10000</i8></value>
  <value><i8>500000000</i8></value>
  <value><i8>500</i8></value>
  <value><i4>1</i4></value>
  <value><i4>1</i4></value>
  <value><i4>0</i4></value>
  <value><i4>0</i4></value>
  <value><string>/downloads/Fedora</string></value>
  <value><string>/downloads</string></value>
  <value><string>linux</string></value>
  <value><i8>1700000001</i8></value>
  <value><string></string></value>
</data></array></value>
</data></array></value></param></params></methodResponse>`
}

func newXMLRPCServer(t *testing.T, wantMethod, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("HTTP method = %q, want POST", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var call xmlMethodCall
		if err := xml.Unmarshal(body, &call); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if call.MethodName != wantMethod {
			t.Errorf("method = %q, want %q", call.MethodName, wantMethod)
		}
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(response))
	}))
}

func TestGetTorrents(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.multicall2", xmlMulticallResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	torrents, err := c.GetTorrents(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(torrents) != 2 {
		t.Fatalf("len = %d, want 2", len(torrents))
	}
	if torrents[0].Hash != "ABC123" {
		t.Errorf("Hash = %q, want %q", torrents[0].Hash, "ABC123")
	}
	if torrents[0].Name != "Ubuntu 24.04" {
		t.Errorf("Name = %q, want %q", torrents[0].Name, "Ubuntu 24.04")
	}
	if !torrents[0].IsComplete {
		t.Error("IsComplete = false, want true")
	}
	if torrents[1].Hash != "DEF456" {
		t.Errorf("Hash = %q, want %q", torrents[1].Hash, "DEF456")
	}
}

func TestAddTorrentURL(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "load.start", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.AddTorrentURL(context.Background(), "magnet:?xt=urn:btih:abc123"); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveTorrent(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.erase", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.RemoveTorrent(context.Background(), "ABC123"); err != nil {
		t.Fatal(err)
	}
}

func TestStartTorrent(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.start", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.StartTorrent(context.Background(), "ABC123"); err != nil {
		t.Fatal(err)
	}
}

func TestStopTorrent(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.stop", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.StopTorrent(context.Background(), "ABC123"); err != nil {
		t.Fatal(err)
	}
}

func TestPauseTorrent(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.pause", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.PauseTorrent(context.Background(), "DEF456"); err != nil {
		t.Fatal(err)
	}
}

func TestResumeTorrent(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.resume", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.ResumeTorrent(context.Background(), "DEF456"); err != nil {
		t.Fatal(err)
	}
}

func TestRecheckTorrent(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.check_hash", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.RecheckTorrent(context.Background(), "ABC123"); err != nil {
		t.Fatal(err)
	}
}

func TestSetLabel(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.custom1.set", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.SetLabel(context.Background(), "ABC123", "movies"); err != nil {
		t.Fatal(err)
	}
}

func TestGetSystemInfo(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "text/xml")
		switch {
		case strings.Contains(string(body), "system.client_version"):
			_, _ = w.Write([]byte(xmlStringResponse("0.9.8")))
		case strings.Contains(string(body), "system.library_version"):
			_, _ = w.Write([]byte(xmlStringResponse("0.13.8")))
		case strings.Contains(string(body), "system.hostname"):
			_, _ = w.Write([]byte(xmlStringResponse("seedbox")))
		case strings.Contains(string(body), "system.pid"):
			_, _ = w.Write([]byte(xmlIntResponse(1234)))
		}
	}))
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	info, err := c.GetSystemInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.ClientVersion != "0.9.8" {
		t.Errorf("ClientVersion = %q, want %q", info.ClientVersion, "0.9.8")
	}
	if info.Hostname != "seedbox" {
		t.Errorf("Hostname = %q, want %q", info.Hostname, "seedbox")
	}
	if info.PID != 1234 {
		t.Errorf("PID = %d, want 1234", info.PID)
	}
}

func TestGetDownloadRate(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "throttle.global_down.rate", xmlIntResponse(5000000))
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	rate, err := c.GetDownloadRate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rate != 5000000 {
		t.Errorf("rate = %d, want 5000000", rate)
	}
}

func TestGetUploadRate(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "throttle.global_up.rate", xmlIntResponse(1000000))
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	rate, err := c.GetUploadRate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rate != 1000000 {
		t.Errorf("rate = %d, want 1000000", rate)
	}
}

func TestSetDownloadLimit(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "throttle.global_down.max_rate.set", xmlVoidResponse())
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	if err := c.SetDownloadLimit(context.Background(), 5000000); err != nil {
		t.Fatal(err)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := newXMLRPCServer(t, "d.start", xmlFaultResponse("Could not find info-hash."))
	defer ts.Close()

	c := rtorrent.New(ts.URL)
	err := c.StartTorrent(context.Background(), "INVALID")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *rtorrent.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestWithAuth(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != "admin" || p != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(xmlStringResponse("0.9.8")))
	}))
	defer ts.Close()

	c := rtorrent.New(ts.URL, rtorrent.WithAuth("admin", "secret"))
	// GetTorrents will work if auth succeeds (even though response is wrong format)
	// but we test that auth is sent. Use a system call that returns a string.
	info, err := c.GetSystemInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.ClientVersion != "0.9.8" {
		t.Errorf("ClientVersion = %q, want %q", info.ClientVersion, "0.9.8")
	}
}
