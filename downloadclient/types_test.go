package downloadclient_test

import (
	"encoding/json"
	"testing"

	"github.com/golusoris/goenvoy/downloadclient"
)

func TestTransferStatusJSON(t *testing.T) {
	t.Parallel()

	ts := downloadclient.TransferStatus{
		ID:              "abc123",
		Name:            "ubuntu.iso",
		State:           downloadclient.TransferStateDownloading,
		Progress:        0.42,
		SizeBytes:       1024000,
		DownloadedBytes: 430080,
		DownloadRate:    50000,
		UploadRate:      10000,
		SavePath:        "/downloads",
		Category:        "iso",
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got downloadclient.TransferStatus
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != ts {
		t.Errorf("got %+v, want %+v", got, ts)
	}
}

func TestClientInfoJSON(t *testing.T) {
	t.Parallel()

	ci := downloadclient.ClientInfo{
		Name:    "qBittorrent",
		Version: "4.6.0",
	}

	data, err := json.Marshal(ci)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got downloadclient.ClientInfo
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != ci {
		t.Errorf("got %+v, want %+v", got, ci)
	}
}

func TestTransferStateConstants(t *testing.T) {
	t.Parallel()

	states := []downloadclient.TransferState{
		downloadclient.TransferStateDownloading,
		downloadclient.TransferStatePaused,
		downloadclient.TransferStateSeeding,
		downloadclient.TransferStateCompleted,
		downloadclient.TransferStateError,
		downloadclient.TransferStateQueued,
	}

	seen := make(map[downloadclient.TransferState]bool)
	for _, s := range states {
		if s == "" {
			t.Error("empty transfer state constant")
		}
		if seen[s] {
			t.Errorf("duplicate state: %s", s)
		}
		seen[s] = true
	}
}
