package discogs

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lusoris/goenvoy/metadata"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New("test-token", metadata.WithBaseURL(ts.URL))
}

func TestGetRelease(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Discogs token=test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(Release{ID: 249504, Title: "Nevermind", Year: 1991})
	})

	rel, err := c.GetRelease(context.Background(), 249504)
	if err != nil {
		t.Fatal(err)
	}
	if rel.Title != "Nevermind" || rel.Year != 1991 {
		t.Fatalf("unexpected release: %+v", rel)
	}
}

func TestGetArtist(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Artist{ID: 125246, Name: "Nirvana"})
	})

	a, err := c.GetArtist(context.Background(), 125246)
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "Nirvana" {
		t.Fatalf("unexpected artist: %+v", a)
	}
}

func TestGetArtistReleases(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(SearchResponse{
			Pagination: Pagination{Page: 1, Pages: 1, Items: 1},
			Results:    []SearchResult{{ID: 249504, Title: "Nevermind", Type: "release"}},
		})
	})

	sr, err := c.GetArtistReleases(context.Background(), 125246, 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(sr.Results) != 1 || sr.Results[0].Title != "Nevermind" {
		t.Fatalf("unexpected results: %+v", sr)
	}
}

func TestGetLabel(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Label{ID: 1, Name: "Planet E"})
	})

	l, err := c.GetLabel(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if l.Name != "Planet E" {
		t.Fatalf("unexpected label: %+v", l)
	}
}

func TestGetMasterRelease(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(MasterRelease{ID: 1000, Title: "Nevermind", Year: 1991})
	})

	m, err := c.GetMasterRelease(context.Background(), 1000)
	if err != nil {
		t.Fatal(err)
	}
	if m.Title != "Nevermind" {
		t.Fatalf("unexpected master: %+v", m)
	}
}

func TestSearch(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(SearchResponse{
			Pagination: Pagination{Page: 1, Pages: 1, Items: 1},
			Results:    []SearchResult{{ID: 249504, Title: "Nevermind", Type: "release"}},
		})
	})

	sr, err := c.Search(context.Background(), "nevermind", "release", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(sr.Results) != 1 || sr.Results[0].Title != "Nevermind" {
		t.Fatalf("unexpected results: %+v", sr)
	}
}

func TestAPIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	_, err := c.GetRelease(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := New("token", metadata.WithHTTPClient(custom))
	if c.HTTPClient() != custom {
		t.Fatal("custom HTTP client not set")
	}
}

func TestWithUserAgent(t *testing.T) {
	c := New("token", metadata.WithUserAgent("myapp/2.0"))
	if c.UserAgent() != "myapp/2.0" {
		t.Fatal("user agent not set")
	}
}

// Release Ratings.

func TestGetReleaseRatingByUser(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/releases/249504/rating/testuser" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ReleaseRating{Username: "testuser", ReleaseID: 249504, Rating: 5})
	})
	rr, err := c.GetReleaseRatingByUser(context.Background(), 249504, "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if rr.Rating != 5 {
		t.Fatalf("Rating = %d, want 5", rr.Rating)
	}
}

func TestUpdateReleaseRating(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(ReleaseRating{Username: "testuser", ReleaseID: 249504, Rating: 4})
	})
	rr, err := c.UpdateReleaseRating(context.Background(), 249504, "testuser", 4)
	if err != nil {
		t.Fatal(err)
	}
	if rr.Rating != 4 {
		t.Fatalf("Rating = %d, want 4", rr.Rating)
	}
}

func TestDeleteReleaseRating(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteReleaseRating(context.Background(), 249504, "testuser"); err != nil {
		t.Fatal(err)
	}
}

func TestGetCommunityReleaseRating(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(CommunityRating{ReleaseID: 249504, Rating: Rating{Count: 100, Average: 4.5}})
	})
	cr, err := c.GetCommunityReleaseRating(context.Background(), 249504)
	if err != nil {
		t.Fatal(err)
	}
	if cr.Rating.Average != 4.5 {
		t.Fatalf("Average = %f, want 4.5", cr.Rating.Average)
	}
}

func TestGetReleaseStats(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(ReleaseStats{NumHave: 1000, NumWant: 500})
	})
	rs, err := c.GetReleaseStats(context.Background(), 249504)
	if err != nil {
		t.Fatal(err)
	}
	if rs.NumHave != 1000 {
		t.Fatalf("NumHave = %d, want 1000", rs.NumHave)
	}
}

// User Identity.

func TestGetIdentity(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Identity{ID: 1, Username: "testuser"})
	})
	id, err := c.GetIdentity(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if id.Username != "testuser" {
		t.Fatalf("Username = %q", id.Username)
	}
}

func TestGetProfile(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/testuser" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Profile{ID: 1, Username: "testuser", NumCollection: 500})
	})
	p, err := c.GetProfile(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if p.NumCollection != 500 {
		t.Fatalf("NumCollection = %d, want 500", p.NumCollection)
	}
}

func TestEditProfile(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(Profile{ID: 1, Username: "testuser", Location: "NYC"})
	})
	p, err := c.EditProfile(context.Background(), "testuser", &ProfileUpdate{Location: "NYC"})
	if err != nil {
		t.Fatal(err)
	}
	if p.Location != "NYC" {
		t.Fatalf("Location = %q", p.Location)
	}
}

func TestGetUserSubmissions(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(SubmissionsResponse{
			Pagination: Pagination{Page: 1, Pages: 1, Items: 1},
		})
	})
	sr, err := c.GetUserSubmissions(context.Background(), "testuser", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if sr.Pagination.Page != 1 {
		t.Fatalf("Page = %d", sr.Pagination.Page)
	}
}

func TestGetUserContributions(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(ContributionsResponse{
			Pagination:    Pagination{Page: 1, Pages: 1, Items: 1},
			Contributions: []SearchResult{{ID: 1, Title: "Test"}},
		})
	})
	cr, err := c.GetUserContributions(context.Background(), "testuser", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(cr.Contributions) != 1 {
		t.Fatalf("len(Contributions) = %d", len(cr.Contributions))
	}
}

// User Collection.

func TestGetCollectionFolders(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(CollectionFoldersResponse{
			Folders: []CollectionFolder{{ID: 0, Name: "All", Count: 100}},
		})
	})
	cf, err := c.GetCollectionFolders(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(cf.Folders) != 1 || cf.Folders[0].Name != "All" {
		t.Fatalf("unexpected folders: %+v", cf.Folders)
	}
}

func TestCreateCollectionFolder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CollectionFolder{ID: 5, Name: "Jazz"})
	})
	f, err := c.CreateCollectionFolder(context.Background(), "testuser", "Jazz")
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "Jazz" {
		t.Fatalf("Name = %q", f.Name)
	}
}

func TestGetCollectionFolder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/testuser/collection/folders/1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(CollectionFolder{ID: 1, Name: "Rock", Count: 50})
	})
	f, err := c.GetCollectionFolder(context.Background(), "testuser", 1)
	if err != nil {
		t.Fatal(err)
	}
	if f.Count != 50 {
		t.Fatalf("Count = %d", f.Count)
	}
}

func TestEditCollectionFolder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(CollectionFolder{ID: 1, Name: "Classic Rock"})
	})
	f, err := c.EditCollectionFolder(context.Background(), "testuser", 1, "Classic Rock")
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "Classic Rock" {
		t.Fatalf("Name = %q", f.Name)
	}
}

func TestDeleteCollectionFolder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteCollectionFolder(context.Background(), "testuser", 1); err != nil {
		t.Fatal(err)
	}
}

func TestGetCollectionItemsByRelease(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(CollectionItemsResponse{
			Pagination: Pagination{Items: 1},
			Releases:   []CollectionItem{{ID: 249504, InstanceID: 1, Rating: 5}},
		})
	})
	ci, err := c.GetCollectionItemsByRelease(context.Background(), "testuser", 249504)
	if err != nil {
		t.Fatal(err)
	}
	if len(ci.Releases) != 1 {
		t.Fatalf("len(Releases) = %d", len(ci.Releases))
	}
}

func TestGetCollectionItemsByFolder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(CollectionItemsResponse{
			Pagination: Pagination{Items: 2},
			Releases:   []CollectionItem{{ID: 1}, {ID: 2}},
		})
	})
	ci, err := c.GetCollectionItemsByFolder(context.Background(), "testuser", 0, 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(ci.Releases) != 2 {
		t.Fatalf("len(Releases) = %d", len(ci.Releases))
	}
}

func TestAddToCollectionFolder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CollectionItem{ID: 249504, InstanceID: 42, FolderID: 1})
	})
	ci, err := c.AddToCollectionFolder(context.Background(), "testuser", 1, 249504)
	if err != nil {
		t.Fatal(err)
	}
	if ci.InstanceID != 42 {
		t.Fatalf("InstanceID = %d", ci.InstanceID)
	}
}

func TestChangeRatingOfRelease(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.ChangeRatingOfRelease(context.Background(), "testuser", 1, 249504, 42, 5); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteInstanceFromFolder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteInstanceFromFolder(context.Background(), "testuser", 1, 249504, 42); err != nil {
		t.Fatal(err)
	}
}

func TestGetCollectionCustomFields(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(CustomFieldsResponse{
			Fields: []CustomField{{ID: 1, Name: "Media Condition", Type: "dropdown"}},
		})
	})
	cf, err := c.GetCollectionCustomFields(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(cf.Fields) != 1 {
		t.Fatalf("len(Fields) = %d", len(cf.Fields))
	}
}

func TestEditFieldsInstance(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.EditFieldsInstance(context.Background(), "testuser", 1, 249504, 42, 1, "Mint"); err != nil {
		t.Fatal(err)
	}
}

func TestGetCollectionValue(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(CollectionValue{Maximum: "$1000", Median: "$500", Minimum: "$100"})
	})
	cv, err := c.GetCollectionValue(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if cv.Median != "$500" {
		t.Fatalf("Median = %q", cv.Median)
	}
}

// User Wantlist.

func TestGetWantlist(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(WantlistResponse{
			Pagination: Pagination{Items: 1},
			Wants:      []WantlistItem{{ID: 1, Rating: 4}},
		})
	})
	wl, err := c.GetWantlist(context.Background(), "testuser", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(wl.Wants) != 1 {
		t.Fatalf("len(Wants) = %d", len(wl.Wants))
	}
}

func TestAddToWantlist(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(WantlistItem{ID: 249504, Rating: 5, Notes: "want it"})
	})
	wi, err := c.AddToWantlist(context.Background(), "testuser", 249504, "want it", 5)
	if err != nil {
		t.Fatal(err)
	}
	if wi.Notes != "want it" {
		t.Fatalf("Notes = %q", wi.Notes)
	}
}

func TestEditWantlistItem(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(WantlistItem{ID: 249504, Rating: 3})
	})
	wi, err := c.EditWantlistItem(context.Background(), "testuser", 249504, "", 3)
	if err != nil {
		t.Fatal(err)
	}
	if wi.Rating != 3 {
		t.Fatalf("Rating = %d", wi.Rating)
	}
}

func TestDeleteFromWantlist(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteFromWantlist(context.Background(), "testuser", 249504); err != nil {
		t.Fatal(err)
	}
}

// User Lists.

func TestGetUserLists(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(UserListsResponse{
			Pagination: Pagination{Items: 1},
			Lists:      []UserListMeta{{ID: 1, Name: "Favorites"}},
		})
	})
	ul, err := c.GetUserLists(context.Background(), "testuser", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(ul.Lists) != 1 {
		t.Fatalf("len(Lists) = %d", len(ul.Lists))
	}
}

func TestGetList(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lists/123" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(UserList{ID: 123, Name: "Favorites", Items: []UserListItem{{ID: 1, Type: "release"}}})
	})
	l, err := c.GetList(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if l.Name != "Favorites" {
		t.Fatalf("Name = %q", l.Name)
	}
	if len(l.Items) != 1 {
		t.Fatalf("len(Items) = %d", len(l.Items))
	}
}

// Marketplace.

func TestGetInventory(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(InventoryResponse{
			Pagination: Pagination{Items: 1},
			Listings:   []Listing{{ID: 1, Status: "For Sale", Condition: "Mint (M)"}},
		})
	})
	ir, err := c.GetInventory(context.Background(), "testuser", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(ir.Listings) != 1 {
		t.Fatalf("len(Listings) = %d", len(ir.Listings))
	}
}

func TestGetListing(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Listing{ID: 1, Status: "For Sale", Condition: "Near Mint (NM or M-)"})
	})
	l, err := c.GetListing(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if l.Condition != "Near Mint (NM or M-)" {
		t.Fatalf("Condition = %q", l.Condition)
	}
}

func TestCreateListing(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Listing{ID: 42, Status: "For Sale", Condition: "Mint (M)"})
	})
	l, err := c.CreateListing(context.Background(), &NewListing{
		ReleaseID: 249504,
		Condition: "Mint (M)",
		Price:     25.00,
	})
	if err != nil {
		t.Fatal(err)
	}
	if l.ID != 42 {
		t.Fatalf("ID = %d", l.ID)
	}
}

func TestEditListing(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.EditListing(context.Background(), 42, &NewListing{Price: 20.00}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteListing(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteListing(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
}

func TestGetOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Order{ID: "1-1", Status: "Payment Received"})
	})
	o, err := c.GetOrder(context.Background(), "1-1")
	if err != nil {
		t.Fatal(err)
	}
	if o.Status != "Payment Received" {
		t.Fatalf("Status = %q", o.Status)
	}
}

func TestEditOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(Order{ID: "1-1", Status: "Shipped"})
	})
	o, err := c.EditOrder(context.Background(), "1-1", &OrderUpdate{Status: "Shipped"})
	if err != nil {
		t.Fatal(err)
	}
	if o.Status != "Shipped" {
		t.Fatalf("Status = %q", o.Status)
	}
}

func TestListOrders(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(OrdersResponse{
			Pagination: Pagination{Items: 1},
			Orders:     []Order{{ID: "1-1", Status: "Payment Received"}},
		})
	})
	or, err := c.ListOrders(context.Background(), "", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(or.Orders) != 1 {
		t.Fatalf("len(Orders) = %d", len(or.Orders))
	}
}

func TestGetOrderMessages(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(OrderMessagesResponse{
			Pagination: Pagination{Items: 1},
			Messages:   []OrderMessage{{Message: "Thanks!"}},
		})
	})
	om, err := c.GetOrderMessages(context.Background(), "1-1", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(om.Messages) != 1 {
		t.Fatalf("len(Messages) = %d", len(om.Messages))
	}
}

func TestAddOrderMessage(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(OrderMessage{Message: "Shipped!"})
	})
	om, err := c.AddOrderMessage(context.Background(), "1-1", "Shipped!", "Shipped")
	if err != nil {
		t.Fatal(err)
	}
	if om.Message != "Shipped!" {
		t.Fatalf("Message = %q", om.Message)
	}
}

func TestGetFee(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/marketplace/fee/") {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Fee{Value: 2.50, Currency: "USD"})
	})
	f, err := c.GetFee(context.Background(), 25.00)
	if err != nil {
		t.Fatal(err)
	}
	if f.Value != 2.50 {
		t.Fatalf("Value = %f", f.Value)
	}
}

func TestGetFeeWithCurrency(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "EUR") {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Fee{Value: 2.00, Currency: "EUR"})
	})
	f, err := c.GetFeeWithCurrency(context.Background(), 20.00, "EUR")
	if err != nil {
		t.Fatal(err)
	}
	if f.Currency != "EUR" {
		t.Fatalf("Currency = %q", f.Currency)
	}
}

func TestGetPriceSuggestions(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(PriceSuggestions{
			Mint:     &SuggestedPrice{Value: 50.00, Currency: "USD"},
			NearMint: &SuggestedPrice{Value: 40.00, Currency: "USD"},
		})
	})
	ps, err := c.GetPriceSuggestions(context.Background(), 249504)
	if err != nil {
		t.Fatal(err)
	}
	if ps.Mint == nil || ps.Mint.Value != 50.00 {
		t.Fatalf("Mint = %+v", ps.Mint)
	}
}

func TestGetMarketplaceReleaseStats(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(MarketplaceReleaseStats{
			NumForSale:  10,
			LowestPrice: &ListingPrice{Value: 15.00, Currency: "USD"},
		})
	})
	mrs, err := c.GetMarketplaceReleaseStats(context.Background(), 249504)
	if err != nil {
		t.Fatal(err)
	}
	if mrs.NumForSale != 10 {
		t.Fatalf("NumForSale = %d", mrs.NumForSale)
	}
}

// Inventory Export.

func TestExportInventory(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.ExportInventory(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetRecentExports(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(ExportsResponse{
			Pagination: Pagination{Items: 1},
			Items:      []Export{{ID: 1, Status: "success"}},
		})
	})
	er, err := c.GetRecentExports(context.Background(), 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(er.Items) != 1 {
		t.Fatalf("len(Items) = %d", len(er.Items))
	}
}

func TestGetExport(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Export{ID: 1, Status: "success", Filename: "export.csv"})
	})
	e, err := c.GetExport(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if e.Filename != "export.csv" {
		t.Fatalf("Filename = %q", e.Filename)
	}
}

func TestDownloadExport(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("csv,data,here"))
	})
	data, err := c.DownloadExport(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "csv,data,here" {
		t.Fatalf("data = %q", string(data))
	}
}

// doJSON / post / put / del errors.

func TestPostError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"forbidden"}`))
	})
	_, err := c.CreateListing(context.Background(), &NewListing{ReleaseID: 1, Condition: "M", Price: 10})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("StatusCode = %d", apiErr.StatusCode)
	}
}

func TestPutError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	})
	_, err := c.AddToWantlist(context.Background(), "testuser", 1, "", 0)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestDeleteError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})
	err := c.DeleteListing(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
}
