package mediaserver_test

import (
	"encoding/json"
	"testing"

	"github.com/golusoris/goenvoy/mediaserver"
)

func TestServerInfoJSON(t *testing.T) {
	t.Parallel()

	si := mediaserver.ServerInfo{
		Name:      "My Server",
		Version:   "10.8.0",
		MachineID: "abc-123",
		Platform:  "linux",
	}

	data, err := json.Marshal(si)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got mediaserver.ServerInfo
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != si {
		t.Errorf("got %+v, want %+v", got, si)
	}
}

func TestLibraryJSON(t *testing.T) {
	t.Parallel()

	lib := mediaserver.Library{
		ID:        "1",
		Name:      "Movies",
		Type:      "movie",
		ItemCount: 500,
	}

	data, err := json.Marshal(lib)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got mediaserver.Library
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != lib {
		t.Errorf("got %+v, want %+v", got, lib)
	}
}

func TestSessionJSON(t *testing.T) {
	t.Parallel()

	sess := mediaserver.Session{
		ID:       "sess-1",
		UserName: "admin",
		Title:    "Inception",
		State:    "playing",
	}

	data, err := json.Marshal(sess)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got mediaserver.Session
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != sess {
		t.Errorf("got %+v, want %+v", got, sess)
	}
}
