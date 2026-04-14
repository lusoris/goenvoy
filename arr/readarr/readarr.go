package readarr

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/golusoris/goenvoy/arr/v2"
)

// Client is a Readarr API client.
type Client struct {
	base *arr.BaseClient
}

// New creates a Readarr [Client] for the instance at baseURL.
func New(baseURL, apiKey string, opts ...arr.Option) (*Client, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{base: base}, nil
}

// GetAllAuthors returns every author configured in Readarr.
func (c *Client) GetAllAuthors(ctx context.Context) ([]Author, error) {
	var out []Author
	if err := c.base.Get(ctx, "/api/v1/author", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAuthor returns a single author by its database ID.
func (c *Client) GetAuthor(ctx context.Context, id int) (*Author, error) {
	var out Author
	path := fmt.Sprintf("/api/v1/author/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddAuthor adds a new author to Readarr.
func (c *Client) AddAuthor(ctx context.Context, author *Author) (*Author, error) {
	var out Author
	if err := c.base.Post(ctx, "/api/v1/author", author, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateAuthor updates an existing author. Set moveFiles to true to relocate
// files when the author path changes.
func (c *Client) UpdateAuthor(ctx context.Context, author *Author, moveFiles bool) (*Author, error) {
	var out Author
	path := fmt.Sprintf("/api/v1/author/%d?moveFiles=%t", author.ID, moveFiles)
	if err := c.base.Put(ctx, path, author, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAuthor removes an author. Set deleteFiles to true to also delete
// downloaded files from disk.
func (c *Client) DeleteAuthor(ctx context.Context, id int, deleteFiles, addImportListExclusion bool) error {
	path := fmt.Sprintf("/api/v1/author/%d?deleteFiles=%t&addImportListExclusion=%t", id, deleteFiles, addImportListExclusion)
	return c.base.Delete(ctx, path, nil, nil)
}

// LookupAuthor searches for an author by name.
func (c *Client) LookupAuthor(ctx context.Context, term string) ([]Author, error) {
	var out []Author
	path := "/api/v1/author/lookup?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetBooks returns books for the given author.
func (c *Client) GetBooks(ctx context.Context, authorID int) ([]Book, error) {
	var out []Book
	path := fmt.Sprintf("/api/v1/book?authorId=%d", authorID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetBook returns a single book by its database ID.
func (c *Client) GetBook(ctx context.Context, id int) (*Book, error) {
	var out Book
	path := fmt.Sprintf("/api/v1/book/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddBook adds a new book to Readarr.
func (c *Client) AddBook(ctx context.Context, book *Book) (*Book, error) {
	var out Book
	if err := c.base.Post(ctx, "/api/v1/book", book, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateBook updates an existing book.
func (c *Client) UpdateBook(ctx context.Context, book *Book) (*Book, error) {
	var out Book
	path := fmt.Sprintf("/api/v1/book/%d", book.ID)
	if err := c.base.Put(ctx, path, book, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBook removes a book.
func (c *Client) DeleteBook(ctx context.Context, id int, deleteFiles, addImportListExclusion bool) error {
	path := fmt.Sprintf("/api/v1/book/%d?deleteFiles=%t&addImportListExclusion=%t", id, deleteFiles, addImportListExclusion)
	return c.base.Delete(ctx, path, nil, nil)
}

// LookupBook searches for a book by term.
func (c *Client) LookupBook(ctx context.Context, term string) ([]Book, error) {
	var out []Book
	path := "/api/v1/book/lookup?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MonitorBooks sets the monitored status for a list of books.
func (c *Client) MonitorBooks(ctx context.Context, req *BooksMonitoredResource) error {
	return c.base.Put(ctx, "/api/v1/book/monitor", req, nil)
}

// GetBookFiles returns all book files for the given author.
func (c *Client) GetBookFiles(ctx context.Context, authorID int) ([]BookFile, error) {
	var out []BookFile
	path := fmt.Sprintf("/api/v1/bookfile?authorId=%d", authorID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetBookFile returns a single book file by its database ID.
func (c *Client) GetBookFile(ctx context.Context, id int) (*BookFile, error) {
	var out BookFile
	path := fmt.Sprintf("/api/v1/bookfile/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBookFile removes a single book file by its database ID.
func (c *Client) DeleteBookFile(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/v1/bookfile/%d", id)
	return c.base.Delete(ctx, path, nil, nil)
}

// DeleteBookFiles removes multiple book files by their IDs.
func (c *Client) DeleteBookFiles(ctx context.Context, ids []int) error {
	body := BookFileListResource{BookFileIDs: ids}
	return c.base.Delete(ctx, "/api/v1/bookfile/bulk", &body, nil)
}

// GetEditions returns all editions for the given book IDs.
func (c *Client) GetEditions(ctx context.Context, bookID int) ([]Edition, error) {
	var out []Edition
	path := fmt.Sprintf("/api/v1/edition?bookId=%d", bookID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCalendar returns books with releases between start and end (RFC 3339 timestamps).
func (c *Client) GetCalendar(ctx context.Context, start, end string, unmonitored bool) ([]Book, error) {
	var out []Book
	path := fmt.Sprintf("/api/v1/calendar?start=%s&end=%s&unmonitored=%t",
		url.QueryEscape(start), url.QueryEscape(end), unmonitored)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendCommand triggers a named command (e.g. "RefreshAuthor", "BookSearch").
func (c *Client) SendCommand(ctx context.Context, cmd arr.CommandRequest) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	if err := c.base.Post(ctx, "/api/v1/command", cmd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCommands returns all currently queued or running commands.
func (c *Client) GetCommands(ctx context.Context) ([]arr.CommandResponse, error) {
	var out []arr.CommandResponse
	if err := c.base.Get(ctx, "/api/v1/command", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCommand returns the status of a single command by its ID.
func (c *Client) GetCommand(ctx context.Context, id int) (*arr.CommandResponse, error) {
	var out arr.CommandResponse
	path := fmt.Sprintf("/api/v1/command/%d", id)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Parse parses a release title and returns the extracted information.
func (c *Client) Parse(ctx context.Context, title string) (*ParseResult, error) {
	var out ParseResult
	path := "/api/v1/parse?title=" + url.QueryEscape(title)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSystemStatus returns Readarr system information.
func (c *Client) GetSystemStatus(ctx context.Context) (*arr.StatusResponse, error) {
	var out arr.StatusResponse
	if err := c.base.Get(ctx, "/api/v1/system/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHealth returns current health check results.
func (c *Client) GetHealth(ctx context.Context) ([]arr.HealthCheck, error) {
	var out []arr.HealthCheck
	if err := c.base.Get(ctx, "/api/v1/health", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDiskSpace returns disk usage information for configured paths.
func (c *Client) GetDiskSpace(ctx context.Context) ([]arr.DiskSpace, error) {
	var out []arr.DiskSpace
	if err := c.base.Get(ctx, "/api/v1/diskspace", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueue returns the current download queue with pagination.
func (c *Client) GetQueue(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.QueueRecord], error) {
	var out arr.PagingResource[arr.QueueRecord]
	path := fmt.Sprintf("/api/v1/queue?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQueueItem removes an item from the download queue.
func (c *Client) DeleteQueueItem(ctx context.Context, id int, removeFromClient, blocklist bool) error {
	path := fmt.Sprintf("/api/v1/queue/%d?removeFromClient=%t&blocklist=%t", id, removeFromClient, blocklist)
	return c.base.Delete(ctx, path, nil, nil)
}

// GetQualityProfiles returns all configured quality profiles.
func (c *Client) GetQualityProfiles(ctx context.Context) ([]arr.QualityProfile, error) {
	var out []arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v1/qualityprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMetadataProfiles returns all configured metadata profiles.
func (c *Client) GetMetadataProfiles(ctx context.Context) ([]MetadataProfile, error) {
	var out []MetadataProfile
	if err := c.base.Get(ctx, "/api/v1/metadataprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTags returns all tags.
func (c *Client) GetTags(ctx context.Context) ([]arr.Tag, error) {
	var out []arr.Tag
	if err := c.base.Get(ctx, "/api/v1/tag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTag creates a new tag and returns it with its assigned ID.
func (c *Client) CreateTag(ctx context.Context, label string) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Post(ctx, "/api/v1/tag", arr.Tag{Label: label}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRootFolders returns all configured root folders.
func (c *Client) GetRootFolders(ctx context.Context) ([]arr.RootFolder, error) {
	var out []arr.RootFolder
	if err := c.base.Get(ctx, "/api/v1/rootfolder", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistory returns the download history with pagination.
func (c *Client) GetHistory(ctx context.Context, page, pageSize int) (*arr.PagingResource[HistoryRecord], error) {
	var out arr.PagingResource[HistoryRecord]
	path := fmt.Sprintf("/api/v1/history?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedMissing returns books with missing files (paginated).
func (c *Client) GetWantedMissing(ctx context.Context, page, pageSize int) (*arr.PagingResource[Book], error) {
	var out arr.PagingResource[Book]
	path := fmt.Sprintf("/api/v1/wanted/missing?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedCutoff returns books not meeting quality cutoff (paginated).
func (c *Client) GetWantedCutoff(ctx context.Context, page, pageSize int) (*arr.PagingResource[Book], error) {
	var out arr.PagingResource[Book]
	path := fmt.Sprintf("/api/v1/wanted/cutoff?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetImportListExclusions returns all import list exclusions.
func (c *Client) GetImportListExclusions(ctx context.Context) ([]ImportListExclusion, error) {
	var out []ImportListExclusion
	if err := c.base.Get(ctx, "/api/v1/importlistexclusion", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSeries returns book series for the given author.
func (c *Client) GetSeries(ctx context.Context, authorID int) ([]Series, error) {
	var out []Series
	path := fmt.Sprintf("/api/v1/series?authorId=%d", authorID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EditAuthors performs a batch update on multiple authors.
func (c *Client) EditAuthors(ctx context.Context, editor *AuthorEditorResource) error {
	return c.base.Put(ctx, "/api/v1/author/editor", editor, nil)
}

// DeleteAuthors performs a batch delete of multiple authors.
func (c *Client) DeleteAuthors(ctx context.Context, editor *AuthorEditorResource) error {
	return c.base.Delete(ctx, "/api/v1/author/editor", editor, nil)
}

// EditBooks performs a batch update on multiple books.
func (c *Client) EditBooks(ctx context.Context, editor *BookEditorResource) error {
	return c.base.Put(ctx, "/api/v1/book/editor", editor, nil)
}

// DeleteBooks performs a batch delete of multiple books.
func (c *Client) DeleteBooks(ctx context.Context, editor *BookEditorResource) error {
	return c.base.Delete(ctx, "/api/v1/book/editor", editor, nil)
}

// ---------- Notifications ----------.

// GetNotifications returns all configured notifications.
func (c *Client) GetNotifications(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/notification", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetNotification returns a single notification by its ID.
func (c *Client) GetNotification(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/notification/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateNotification creates a new notification.
func (c *Client) CreateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v1/notification", n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNotification updates an existing notification.
func (c *Client) UpdateNotification(ctx context.Context, n *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/notification/%d", n.ID), n, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNotification deletes a notification by ID.
func (c *Client) DeleteNotification(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/notification/%d", id), nil, nil)
}

// GetNotificationSchema returns the schema for all notification types.
func (c *Client) GetNotificationSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/notification/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestNotification tests a notification configuration.
func (c *Client) TestNotification(ctx context.Context, n *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v1/notification/test", n, nil)
}

// TestAllNotifications tests all configured notifications.
func (c *Client) TestAllNotifications(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/notification/testall", nil, nil)
}

// NotificationAction triggers a named action on a notification provider.
func (c *Client) NotificationAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v1/notification/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Download Clients ----------.

// GetDownloadClients returns all configured download clients.
func (c *Client) GetDownloadClients(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/downloadclient", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDownloadClient returns a single download client by its ID.
func (c *Client) GetDownloadClient(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/downloadclient/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDownloadClient creates a new download client.
func (c *Client) CreateDownloadClient(ctx context.Context, dc *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v1/downloadclient", dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClient updates an existing download client.
func (c *Client) UpdateDownloadClient(ctx context.Context, dc *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/downloadclient/%d", dc.ID), dc, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDownloadClient deletes a download client by ID.
func (c *Client) DeleteDownloadClient(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/downloadclient/%d", id), nil, nil)
}

// GetDownloadClientSchema returns the schema for all download client types.
func (c *Client) GetDownloadClientSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/downloadclient/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestDownloadClient tests a download client configuration.
func (c *Client) TestDownloadClient(ctx context.Context, dc *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v1/downloadclient/test", dc, nil)
}

// TestAllDownloadClients tests all configured download clients.
func (c *Client) TestAllDownloadClients(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/downloadclient/testall", nil, nil)
}

// BulkUpdateDownloadClients updates multiple download clients at once.
func (c *Client) BulkUpdateDownloadClients(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v1/downloadclient/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteDownloadClients deletes multiple download clients at once.
func (c *Client) BulkDeleteDownloadClients(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v1/downloadclient/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// DownloadClientAction triggers a named action on a download client provider.
func (c *Client) DownloadClientAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v1/downloadclient/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Indexers ----------.

// GetIndexers returns all configured indexers.
func (c *Client) GetIndexers(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/indexer", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetIndexer returns a single indexer by its ID.
func (c *Client) GetIndexer(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/indexer/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateIndexer creates a new indexer.
func (c *Client) CreateIndexer(ctx context.Context, idx *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v1/indexer", idx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexer updates an existing indexer.
func (c *Client) UpdateIndexer(ctx context.Context, idx *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/indexer/%d", idx.ID), idx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteIndexer deletes an indexer by ID.
func (c *Client) DeleteIndexer(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/indexer/%d", id), nil, nil)
}

// GetIndexerSchema returns the schema for all indexer types.
func (c *Client) GetIndexerSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/indexer/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestIndexer tests an indexer configuration.
func (c *Client) TestIndexer(ctx context.Context, idx *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v1/indexer/test", idx, nil)
}

// TestAllIndexers tests all configured indexers.
func (c *Client) TestAllIndexers(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/indexer/testall", nil, nil)
}

// BulkUpdateIndexers updates multiple indexers at once.
func (c *Client) BulkUpdateIndexers(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v1/indexer/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteIndexers deletes multiple indexers at once.
func (c *Client) BulkDeleteIndexers(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v1/indexer/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// IndexerAction triggers a named action on an indexer provider.
func (c *Client) IndexerAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v1/indexer/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Import Lists ----------.

// GetImportLists returns all configured import lists.
func (c *Client) GetImportLists(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/importlist", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetImportList returns a single import list by its ID.
func (c *Client) GetImportList(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/importlist/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateImportList creates a new import list.
func (c *Client) CreateImportList(ctx context.Context, il *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v1/importlist", il, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportList updates an existing import list.
func (c *Client) UpdateImportList(ctx context.Context, il *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/importlist/%d", il.ID), il, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportList deletes an import list by ID.
func (c *Client) DeleteImportList(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/importlist/%d", id), nil, nil)
}

// GetImportListSchema returns the schema for all import list types.
func (c *Client) GetImportListSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/importlist/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestImportList tests an import list configuration.
func (c *Client) TestImportList(ctx context.Context, il *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v1/importlist/test", il, nil)
}

// TestAllImportLists tests all configured import lists.
func (c *Client) TestAllImportLists(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/importlist/testall", nil, nil)
}

// BulkUpdateImportLists updates multiple import lists at once.
func (c *Client) BulkUpdateImportLists(ctx context.Context, bulk *arr.ProviderBulkResource) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Put(ctx, "/api/v1/importlist/bulk", bulk, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkDeleteImportLists deletes multiple import lists at once.
func (c *Client) BulkDeleteImportLists(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v1/importlist/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// ImportListAction triggers a named action on an import list provider.
func (c *Client) ImportListAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v1/importlist/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Metadata Consumers ----------.

// GetMetadataConsumers returns all configured metadata consumers.
func (c *Client) GetMetadataConsumers(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/metadata", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMetadataConsumer returns a single metadata consumer by its ID.
func (c *Client) GetMetadataConsumer(ctx context.Context, id int) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/metadata/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateMetadataConsumer creates a new metadata consumer.
func (c *Client) CreateMetadataConsumer(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Post(ctx, "/api/v1/metadata", m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMetadataConsumer updates an existing metadata consumer.
func (c *Client) UpdateMetadataConsumer(ctx context.Context, m *arr.ProviderResource) (*arr.ProviderResource, error) {
	var out arr.ProviderResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/metadata/%d", m.ID), m, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMetadataConsumer deletes a metadata consumer by ID.
func (c *Client) DeleteMetadataConsumer(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/metadata/%d", id), nil, nil)
}

// GetMetadataConsumerSchema returns the schema for all metadata consumer types.
func (c *Client) GetMetadataConsumerSchema(ctx context.Context) ([]arr.ProviderResource, error) {
	var out []arr.ProviderResource
	if err := c.base.Get(ctx, "/api/v1/metadata/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestMetadataConsumer tests a metadata consumer configuration.
func (c *Client) TestMetadataConsumer(ctx context.Context, m *arr.ProviderResource) error {
	return c.base.Post(ctx, "/api/v1/metadata/test", m, nil)
}

// TestAllMetadataConsumers tests all configured metadata consumers.
func (c *Client) TestAllMetadataConsumers(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/metadata/testall", nil, nil)
}

// MetadataConsumerAction triggers a named action on a metadata consumer provider.
func (c *Client) MetadataConsumerAction(ctx context.Context, name string, body *arr.ProviderResource) error {
	path := "/api/v1/metadata/action/" + url.PathEscape(name)
	return c.base.Post(ctx, path, body, nil)
}

// ---------- Config Endpoints ----------.

// GetDownloadClientConfig returns the download client configuration.
func (c *Client) GetDownloadClientConfig(ctx context.Context) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/downloadclient", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDownloadClientConfig updates the download client configuration.
func (c *Client) UpdateDownloadClientConfig(ctx context.Context, config *arr.DownloadClientConfigResource) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/downloadclient/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDownloadClientConfigByID returns the download client config by its ID.
func (c *Client) GetDownloadClientConfigByID(ctx context.Context, id int) (*arr.DownloadClientConfigResource, error) {
	var out arr.DownloadClientConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/downloadclient/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIndexerConfig returns the indexer configuration.
func (c *Client) GetIndexerConfig(ctx context.Context) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/indexer", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateIndexerConfig updates the indexer configuration.
func (c *Client) UpdateIndexerConfig(ctx context.Context, config *arr.IndexerConfigResource) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/indexer/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIndexerConfigByID returns the indexer config by its ID.
func (c *Client) GetIndexerConfigByID(ctx context.Context, id int) (*arr.IndexerConfigResource, error) {
	var out arr.IndexerConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/indexer/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNamingConfig returns the naming configuration.
func (c *Client) GetNamingConfig(ctx context.Context) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/naming", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateNamingConfig updates the naming configuration.
func (c *Client) UpdateNamingConfig(ctx context.Context, config *arr.NamingConfigResource) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/naming/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNamingConfigByID returns the naming config by its ID.
func (c *Client) GetNamingConfigByID(ctx context.Context, id int) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/naming/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNamingExamples returns naming format examples based on the current naming config.
func (c *Client) GetNamingExamples(ctx context.Context) (*arr.NamingConfigResource, error) {
	var out arr.NamingConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/naming/examples", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHostConfig returns the host configuration.
func (c *Client) GetHostConfig(ctx context.Context) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/host", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateHostConfig updates the host configuration.
func (c *Client) UpdateHostConfig(ctx context.Context, config *arr.HostConfigResource) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/host/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHostConfigByID returns the host config by its ID.
func (c *Client) GetHostConfigByID(ctx context.Context, id int) (*arr.HostConfigResource, error) {
	var out arr.HostConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/host/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUIConfig returns the UI configuration.
func (c *Client) GetUIConfig(ctx context.Context) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/ui", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateUIConfig updates the UI configuration.
func (c *Client) UpdateUIConfig(ctx context.Context, config *arr.UIConfigResource) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/ui/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUIConfigByID returns the UI config by its ID.
func (c *Client) GetUIConfigByID(ctx context.Context, id int) (*arr.UIConfigResource, error) {
	var out arr.UIConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/ui/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMediaManagementConfig returns the media management configuration.
func (c *Client) GetMediaManagementConfig(ctx context.Context) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/mediamanagement", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMediaManagementConfig updates the media management configuration.
func (c *Client) UpdateMediaManagementConfig(ctx context.Context, config *arr.MediaManagementConfigResource) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/mediamanagement/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMediaManagementConfigByID returns the media management config by its ID.
func (c *Client) GetMediaManagementConfigByID(ctx context.Context, id int) (*arr.MediaManagementConfigResource, error) {
	var out arr.MediaManagementConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/mediamanagement/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Development Config ----------.

// GetDevelopmentConfig returns the development configuration.
func (c *Client) GetDevelopmentConfig(ctx context.Context) (*DevelopmentConfigResource, error) {
	var out DevelopmentConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/development", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDevelopmentConfigByID returns the development config by its ID.
func (c *Client) GetDevelopmentConfigByID(ctx context.Context, id int) (*DevelopmentConfigResource, error) {
	var out DevelopmentConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/development/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDevelopmentConfig updates the development configuration.
func (c *Client) UpdateDevelopmentConfig(ctx context.Context, config *DevelopmentConfigResource) (*DevelopmentConfigResource, error) {
	var out DevelopmentConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/development/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Metadata Provider Config ----------.

// GetMetadataProviderConfig returns the metadata provider configuration.
func (c *Client) GetMetadataProviderConfig(ctx context.Context) (*MetadataProviderConfigResource, error) {
	var out MetadataProviderConfigResource
	if err := c.base.Get(ctx, "/api/v1/config/metadataprovider", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMetadataProviderConfigByID returns the metadata provider config by its ID.
func (c *Client) GetMetadataProviderConfigByID(ctx context.Context, id int) (*MetadataProviderConfigResource, error) {
	var out MetadataProviderConfigResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/config/metadataprovider/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMetadataProviderConfig updates the metadata provider configuration.
func (c *Client) UpdateMetadataProviderConfig(ctx context.Context, config *MetadataProviderConfigResource) (*MetadataProviderConfigResource, error) {
	var out MetadataProviderConfigResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/config/metadataprovider/%d", config.ID), config, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Quality Profiles ----------.

// GetQualityProfile returns a single quality profile by its ID.
func (c *Client) GetQualityProfile(ctx context.Context, id int) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/qualityprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateQualityProfile creates a new quality profile.
func (c *Client) CreateQualityProfile(ctx context.Context, profile *arr.QualityProfile) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Post(ctx, "/api/v1/qualityprofile", profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateQualityProfile updates an existing quality profile.
func (c *Client) UpdateQualityProfile(ctx context.Context, profile *arr.QualityProfile) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/qualityprofile/%d", profile.ID), profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQualityProfile deletes a quality profile by ID.
func (c *Client) DeleteQualityProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/qualityprofile/%d", id), nil, nil)
}

// GetQualityProfileSchema returns the quality profile schema.
func (c *Client) GetQualityProfileSchema(ctx context.Context) (*arr.QualityProfile, error) {
	var out arr.QualityProfile
	if err := c.base.Get(ctx, "/api/v1/qualityprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Quality Definitions ----------.

// GetQualityDefinitions returns all quality definitions.
func (c *Client) GetQualityDefinitions(ctx context.Context) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Get(ctx, "/api/v1/qualitydefinition", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQualityDefinition returns a single quality definition by its ID.
func (c *Client) GetQualityDefinition(ctx context.Context, id int) (*arr.QualityDefinitionResource, error) {
	var out arr.QualityDefinitionResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/qualitydefinition/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateQualityDefinition updates a single quality definition.
func (c *Client) UpdateQualityDefinition(ctx context.Context, def *arr.QualityDefinitionResource) (*arr.QualityDefinitionResource, error) {
	var out arr.QualityDefinitionResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/qualitydefinition/%d", def.ID), def, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BulkUpdateQualityDefinitions updates multiple quality definitions at once.
func (c *Client) BulkUpdateQualityDefinitions(ctx context.Context, defs []arr.QualityDefinitionResource) ([]arr.QualityDefinitionResource, error) {
	var out []arr.QualityDefinitionResource
	if err := c.base.Put(ctx, "/api/v1/qualitydefinition/update", defs, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Metadata Profiles extended ----------.

// GetMetadataProfile returns a single metadata profile by its ID.
func (c *Client) GetMetadataProfile(ctx context.Context, id int) (*MetadataProfile, error) {
	var out MetadataProfile
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/metadataprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateMetadataProfile creates a new metadata profile.
func (c *Client) CreateMetadataProfile(ctx context.Context, profile *MetadataProfile) (*MetadataProfile, error) {
	var out MetadataProfile
	if err := c.base.Post(ctx, "/api/v1/metadataprofile", profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMetadataProfile updates an existing metadata profile.
func (c *Client) UpdateMetadataProfile(ctx context.Context, profile *MetadataProfile) (*MetadataProfile, error) {
	var out MetadataProfile
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/metadataprofile/%d", profile.ID), profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMetadataProfile deletes a metadata profile by ID.
func (c *Client) DeleteMetadataProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/metadataprofile/%d", id), nil, nil)
}

// GetMetadataProfileSchema returns the metadata profile schema.
func (c *Client) GetMetadataProfileSchema(ctx context.Context) (*MetadataProfile, error) {
	var out MetadataProfile
	if err := c.base.Get(ctx, "/api/v1/metadataprofile/schema", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Tags extended ----------.

// GetTag returns a single tag by its ID.
func (c *Client) GetTag(ctx context.Context, id int) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/tag/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTag updates an existing tag.
func (c *Client) UpdateTag(ctx context.Context, tag *arr.Tag) (*arr.Tag, error) {
	var out arr.Tag
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/tag/%d", tag.ID), tag, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTag deletes a tag by ID.
func (c *Client) DeleteTag(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/tag/%d", id), nil, nil)
}

// GetTagDetails returns all tags with details about their usage.
func (c *Client) GetTagDetails(ctx context.Context) ([]arr.TagDetail, error) {
	var out []arr.TagDetail
	if err := c.base.Get(ctx, "/api/v1/tag/detail", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTagDetail returns a single tag detail by its ID.
func (c *Client) GetTagDetail(ctx context.Context, id int) (*arr.TagDetail, error) {
	var out arr.TagDetail
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/tag/detail/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Root Folders extended ----------.

// GetRootFolder returns a single root folder by its ID.
func (c *Client) GetRootFolder(ctx context.Context, id int) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/rootfolder/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRootFolder creates a new root folder.
func (c *Client) CreateRootFolder(ctx context.Context, folder *arr.RootFolder) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Post(ctx, "/api/v1/rootfolder", folder, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateRootFolder updates an existing root folder.
func (c *Client) UpdateRootFolder(ctx context.Context, folder *arr.RootFolder) (*arr.RootFolder, error) {
	var out arr.RootFolder
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/rootfolder/%d", folder.ID), folder, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRootFolder deletes a root folder by ID.
func (c *Client) DeleteRootFolder(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/rootfolder/%d", id), nil, nil)
}

// ---------- Custom Filters ----------.

// GetCustomFilters returns all custom filters.
func (c *Client) GetCustomFilters(ctx context.Context) ([]arr.CustomFilterResource, error) {
	var out []arr.CustomFilterResource
	if err := c.base.Get(ctx, "/api/v1/customfilter", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCustomFilter returns a single custom filter by its ID.
func (c *Client) GetCustomFilter(ctx context.Context, id int) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/customfilter/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCustomFilter creates a new custom filter.
func (c *Client) CreateCustomFilter(ctx context.Context, filter *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Post(ctx, "/api/v1/customfilter", filter, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFilter updates an existing custom filter.
func (c *Client) UpdateCustomFilter(ctx context.Context, filter *arr.CustomFilterResource) (*arr.CustomFilterResource, error) {
	var out arr.CustomFilterResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/customfilter/%d", filter.ID), filter, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFilter deletes a custom filter by ID.
func (c *Client) DeleteCustomFilter(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/customfilter/%d", id), nil, nil)
}

// ---------- Custom Formats ----------.

// GetCustomFormats returns all custom formats.
func (c *Client) GetCustomFormats(ctx context.Context) ([]arr.CustomFormatResource, error) {
	var out []arr.CustomFormatResource
	if err := c.base.Get(ctx, "/api/v1/customformat", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCustomFormat returns a single custom format by its ID.
func (c *Client) GetCustomFormat(ctx context.Context, id int) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/customformat/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCustomFormat creates a new custom format.
func (c *Client) CreateCustomFormat(ctx context.Context, format *arr.CustomFormatResource) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Post(ctx, "/api/v1/customformat", format, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomFormat updates an existing custom format.
func (c *Client) UpdateCustomFormat(ctx context.Context, format *arr.CustomFormatResource) (*arr.CustomFormatResource, error) {
	var out arr.CustomFormatResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/customformat/%d", format.ID), format, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomFormat deletes a custom format by ID.
func (c *Client) DeleteCustomFormat(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/customformat/%d", id), nil, nil)
}

// GetCustomFormatSchema returns the custom format schema.
func (c *Client) GetCustomFormatSchema(ctx context.Context) ([]arr.CustomFormatResource, error) {
	var out []arr.CustomFormatResource
	if err := c.base.Get(ctx, "/api/v1/customformat/schema", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Delay Profiles ----------.

// GetDelayProfiles returns all delay profiles.
func (c *Client) GetDelayProfiles(ctx context.Context) ([]arr.DelayProfileResource, error) {
	var out []arr.DelayProfileResource
	if err := c.base.Get(ctx, "/api/v1/delayprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDelayProfile returns a single delay profile by its ID.
func (c *Client) GetDelayProfile(ctx context.Context, id int) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/delayprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateDelayProfile creates a new delay profile.
func (c *Client) CreateDelayProfile(ctx context.Context, profile *arr.DelayProfileResource) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Post(ctx, "/api/v1/delayprofile", profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateDelayProfile updates an existing delay profile.
func (c *Client) UpdateDelayProfile(ctx context.Context, profile *arr.DelayProfileResource) (*arr.DelayProfileResource, error) {
	var out arr.DelayProfileResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/delayprofile/%d", profile.ID), profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteDelayProfile deletes a delay profile by ID.
func (c *Client) DeleteDelayProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/delayprofile/%d", id), nil, nil)
}

// ReorderDelayProfile changes the order of a delay profile.
func (c *Client) ReorderDelayProfile(ctx context.Context, id, afterID int) ([]arr.DelayProfileResource, error) {
	var out []arr.DelayProfileResource
	path := fmt.Sprintf("/api/v1/delayprofile/reorder/%d?after=%d", id, afterID)
	if err := c.base.Put(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Release Profiles ----------.

// GetReleaseProfiles returns all release profiles.
func (c *Client) GetReleaseProfiles(ctx context.Context) ([]arr.ReleaseProfileResource, error) {
	var out []arr.ReleaseProfileResource
	if err := c.base.Get(ctx, "/api/v1/releaseprofile", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetReleaseProfile returns a single release profile by its ID.
func (c *Client) GetReleaseProfile(ctx context.Context, id int) (*arr.ReleaseProfileResource, error) {
	var out arr.ReleaseProfileResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/releaseprofile/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateReleaseProfile creates a new release profile.
func (c *Client) CreateReleaseProfile(ctx context.Context, profile *arr.ReleaseProfileResource) (*arr.ReleaseProfileResource, error) {
	var out arr.ReleaseProfileResource
	if err := c.base.Post(ctx, "/api/v1/releaseprofile", profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateReleaseProfile updates an existing release profile.
func (c *Client) UpdateReleaseProfile(ctx context.Context, profile *arr.ReleaseProfileResource) (*arr.ReleaseProfileResource, error) {
	var out arr.ReleaseProfileResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/releaseprofile/%d", profile.ID), profile, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteReleaseProfile deletes a release profile by ID.
func (c *Client) DeleteReleaseProfile(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/releaseprofile/%d", id), nil, nil)
}

// ---------- Remote Path Mappings ----------.

// GetRemotePathMappings returns all remote path mappings.
func (c *Client) GetRemotePathMappings(ctx context.Context) ([]arr.RemotePathMappingResource, error) {
	var out []arr.RemotePathMappingResource
	if err := c.base.Get(ctx, "/api/v1/remotepathmapping", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetRemotePathMapping returns a single remote path mapping by its ID.
func (c *Client) GetRemotePathMapping(ctx context.Context, id int) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/remotepathmapping/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRemotePathMapping creates a new remote path mapping.
func (c *Client) CreateRemotePathMapping(ctx context.Context, mapping *arr.RemotePathMappingResource) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Post(ctx, "/api/v1/remotepathmapping", mapping, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateRemotePathMapping updates an existing remote path mapping.
func (c *Client) UpdateRemotePathMapping(ctx context.Context, mapping *arr.RemotePathMappingResource) (*arr.RemotePathMappingResource, error) {
	var out arr.RemotePathMappingResource
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/remotepathmapping/%d", mapping.ID), mapping, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRemotePathMapping deletes a remote path mapping by ID.
func (c *Client) DeleteRemotePathMapping(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/remotepathmapping/%d", id), nil, nil)
}

// ---------- Import List Exclusions extended ----------.

// GetImportListExclusion returns a single import list exclusion by its ID.
func (c *Client) GetImportListExclusion(ctx context.Context, id int) (*ImportListExclusion, error) {
	var out ImportListExclusion
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/importlistexclusion/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateImportListExclusion creates a new import list exclusion.
func (c *Client) CreateImportListExclusion(ctx context.Context, exclusion *ImportListExclusion) (*ImportListExclusion, error) {
	var out ImportListExclusion
	if err := c.base.Post(ctx, "/api/v1/importlistexclusion", exclusion, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateImportListExclusion updates an existing import list exclusion.
func (c *Client) UpdateImportListExclusion(ctx context.Context, exclusion *ImportListExclusion) (*ImportListExclusion, error) {
	var out ImportListExclusion
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/importlistexclusion/%d", exclusion.ID), exclusion, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteImportListExclusion deletes an import list exclusion by ID.
func (c *Client) DeleteImportListExclusion(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/importlistexclusion/%d", id), nil, nil)
}

// ---------- Blocklist ----------.

// GetBlocklist returns the blocklist with pagination.
func (c *Client) GetBlocklist(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.BlocklistResource], error) {
	var out arr.PagingResource[arr.BlocklistResource]
	path := fmt.Sprintf("/api/v1/blocklist?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBlocklistItem deletes a single blocklist item by ID.
func (c *Client) DeleteBlocklistItem(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/blocklist/%d", id), nil, nil)
}

// BulkDeleteBlocklist deletes multiple blocklist items.
func (c *Client) BulkDeleteBlocklist(ctx context.Context, ids []int) error {
	return c.base.Delete(ctx, "/api/v1/blocklist/bulk", &arr.BlocklistBulkResource{IDs: ids}, nil)
}

// ---------- Queue Extended ----------.

// BulkDeleteQueue deletes multiple items from the download queue.
func (c *Client) BulkDeleteQueue(ctx context.Context, bulk *arr.QueueBulkResource, removeFromClient, blocklist bool) error {
	path := fmt.Sprintf("/api/v1/queue/bulk?removeFromClient=%t&blocklist=%t", removeFromClient, blocklist)
	return c.base.Delete(ctx, path, bulk, nil)
}

// GrabQueueItem grabs a pending release from the queue by its ID.
func (c *Client) GrabQueueItem(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v1/queue/grab/%d", id), nil, nil)
}

// GrabQueueItemsBulk grabs multiple pending releases from the queue.
func (c *Client) GrabQueueItemsBulk(ctx context.Context, ids []int) error {
	return c.base.Post(ctx, "/api/v1/queue/grab/bulk", struct {
		IDs []int `json:"ids"`
	}{IDs: ids}, nil)
}

// GetQueueDetails returns detailed information about all items in the queue.
func (c *Client) GetQueueDetails(ctx context.Context) ([]arr.QueueRecord, error) {
	var out []arr.QueueRecord
	if err := c.base.Get(ctx, "/api/v1/queue/details", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetQueueStatus returns the overall status of the download queue.
func (c *Client) GetQueueStatus(ctx context.Context) (*arr.QueueStatusResource, error) {
	var out arr.QueueStatusResource
	if err := c.base.Get(ctx, "/api/v1/queue/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- History Extended ----------.

// GetHistoryByAuthor returns history for a specific author.
func (c *Client) GetHistoryByAuthor(ctx context.Context, authorID int) ([]HistoryRecord, error) {
	var out []HistoryRecord
	path := fmt.Sprintf("/api/v1/history/author?authorId=%d", authorID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistorySince returns history since a given date.
func (c *Client) GetHistorySince(ctx context.Context, date string) ([]HistoryRecord, error) {
	var out []HistoryRecord
	path := "/api/v1/history/since?date=" + url.QueryEscape(date)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MarkHistoryFailed marks a history item as failed.
func (c *Client) MarkHistoryFailed(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v1/history/failed/%d", id), nil, nil)
}

// ---------- Releases ----------.

// SearchReleases searches for releases using the configured indexers.
func (c *Client) SearchReleases(ctx context.Context, bookID int) ([]arr.ReleaseResource, error) {
	var out []arr.ReleaseResource
	path := fmt.Sprintf("/api/v1/release?bookId=%d", bookID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GrabRelease grabs a release for download.
func (c *Client) GrabRelease(ctx context.Context, release *arr.ReleaseResource) (*arr.ReleaseResource, error) {
	var out arr.ReleaseResource
	if err := c.base.Post(ctx, "/api/v1/release", release, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PushRelease pushes a release to the download client.
func (c *Client) PushRelease(ctx context.Context, release *arr.ReleasePushResource) ([]arr.ReleaseResource, error) {
	var out []arr.ReleaseResource
	if err := c.base.Post(ctx, "/api/v1/release/push", release, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Rename ----------.

// GetRenamePreview returns a preview of book file renames for an author.
func (c *Client) GetRenamePreview(ctx context.Context, authorID, bookID int) ([]RenameBookResource, error) {
	var out []RenameBookResource
	path := fmt.Sprintf("/api/v1/rename?authorId=%d&bookId=%d", authorID, bookID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Retag ----------.

// GetRetagPreview returns a preview of book file retags for an author.
func (c *Client) GetRetagPreview(ctx context.Context, authorID, bookID int) ([]RetagBookResource, error) {
	var out []RetagBookResource
	path := fmt.Sprintf("/api/v1/retag?authorId=%d&bookId=%d", authorID, bookID)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Manual Import ----------.

// GetManualImport returns a list of potential imports for the given path.
func (c *Client) GetManualImport(ctx context.Context, folder string) ([]arr.ManualImportResource, error) {
	var out []arr.ManualImportResource
	path := "/api/v1/manualimport?folder=" + url.QueryEscape(folder)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ReprocessManualImport reprocesses manual imports.
func (c *Client) ReprocessManualImport(ctx context.Context, items []arr.ManualImportReprocessResource) ([]arr.ManualImportResource, error) {
	var out []arr.ManualImportResource
	if err := c.base.Post(ctx, "/api/v1/manualimport", items, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Backups ----------.

// GetBackups returns a list of all available backups.
func (c *Client) GetBackups(ctx context.Context) ([]arr.Backup, error) {
	var out []arr.Backup
	if err := c.base.Get(ctx, "/api/v1/system/backup", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBackup deletes a backup by ID.
func (c *Client) DeleteBackup(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/system/backup/%d", id), nil, nil)
}

// RestoreBackup triggers a restore from a backup by ID.
func (c *Client) RestoreBackup(ctx context.Context, id int) error {
	return c.base.Post(ctx, fmt.Sprintf("/api/v1/system/backup/restore/%d", id), nil, nil)
}

// ---------- Logs ----------.

// GetLogs returns log entries with pagination.
func (c *Client) GetLogs(ctx context.Context, page, pageSize int) (*arr.PagingResource[arr.LogRecord], error) {
	var out arr.PagingResource[arr.LogRecord]
	path := fmt.Sprintf("/api/v1/log?page=%d&pageSize=%d", page, pageSize)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLogFiles returns a list of log files.
func (c *Client) GetLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v1/log/file", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLogFileContent returns the content of a specific log file by filename.
func (c *Client) GetLogFileContent(ctx context.Context, filename string) (string, error) {
	path := "/api/v1/log/file/" + url.PathEscape(filename)
	b, err := c.base.GetRaw(ctx, path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GetUpdateLogFiles returns a list of update log files.
func (c *Client) GetUpdateLogFiles(ctx context.Context) ([]arr.LogFileResource, error) {
	var out []arr.LogFileResource
	if err := c.base.Get(ctx, "/api/v1/log/file/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetUpdateLogFileContent returns the content of a specific update log file.
func (c *Client) GetUpdateLogFileContent(ctx context.Context, filename string) (string, error) {
	path := "/api/v1/log/file/update/" + url.PathEscape(filename)
	b, err := c.base.GetRaw(ctx, path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------- System ----------.

// GetTasks returns all scheduled tasks.
func (c *Client) GetTasks(ctx context.Context) ([]arr.TaskResource, error) {
	var out []arr.TaskResource
	if err := c.base.Get(ctx, "/api/v1/system/task", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTask returns a single task by its ID.
func (c *Client) GetTask(ctx context.Context, id int) (*arr.TaskResource, error) {
	var out arr.TaskResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/system/task/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUpdates returns available application updates.
func (c *Client) GetUpdates(ctx context.Context) ([]arr.UpdateResource, error) {
	var out []arr.UpdateResource
	if err := c.base.Get(ctx, "/api/v1/update", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSystemRoutes returns all API routes.
func (c *Client) GetSystemRoutes(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v1/system/routes", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSystemRoutesDuplicate returns duplicate API routes.
func (c *Client) GetSystemRoutesDuplicate(ctx context.Context) ([]arr.SystemRouteResource, error) {
	var out []arr.SystemRouteResource
	if err := c.base.Get(ctx, "/api/v1/system/routes/duplicate", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Shutdown sends a shutdown command to Readarr.
func (c *Client) Shutdown(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/system/shutdown", nil, nil)
}

// Restart sends a restart command to Readarr.
func (c *Client) Restart(ctx context.Context) error {
	return c.base.Post(ctx, "/api/v1/system/restart", nil, nil)
}

// DeleteCommand deletes a command by ID.
func (c *Client) DeleteCommand(ctx context.Context, id int) error {
	return c.base.Delete(ctx, fmt.Sprintf("/api/v1/command/%d", id), nil, nil)
}

// ---------- Languages ----------.

// GetLanguages returns all languages.
func (c *Client) GetLanguages(ctx context.Context) ([]arr.LanguageResource, error) {
	var out []arr.LanguageResource
	if err := c.base.Get(ctx, "/api/v1/language", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLanguage returns a single language by its ID.
func (c *Client) GetLanguage(ctx context.Context, id int) (*arr.LanguageResource, error) {
	var out arr.LanguageResource
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/language/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Localization ----------.

// GetLocalization returns the localization strings.
func (c *Client) GetLocalization(ctx context.Context) (*LocalizationResource, error) {
	var out LocalizationResource
	if err := c.base.Get(ctx, "/api/v1/localization", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Ping ----------.

// Ping checks connectivity to the Readarr instance.
func (c *Client) Ping(ctx context.Context) error {
	return c.base.Get(ctx, "/ping", nil)
}

// ---------- Indexer Flags ----------.

// GetIndexerFlags returns the list of indexer flags.
func (c *Client) GetIndexerFlags(ctx context.Context) ([]arr.IndexerFlagResource, error) {
	var out []arr.IndexerFlagResource
	if err := c.base.Get(ctx, "/api/v1/indexerflag", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- File System ----------.

// BrowseFileSystem returns directory/file listings for the given path.
func (c *Client) BrowseFileSystem(ctx context.Context, path string) (*FileSystemResource, error) {
	var out FileSystemResource
	reqPath := "/api/v1/filesystem?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFileSystemType returns the filesystem type for the given path.
func (c *Client) GetFileSystemType(ctx context.Context, path string) (string, error) {
	var out string
	reqPath := "/api/v1/filesystem/type?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return "", err
	}
	return out, nil
}

// GetFileSystemMediaFiles returns media files for the given path.
func (c *Client) GetFileSystemMediaFiles(ctx context.Context, path string) ([]FileSystemEntry, error) {
	var out []FileSystemEntry
	reqPath := "/api/v1/filesystem/mediafiles?path=" + url.QueryEscape(path)
	if err := c.base.Get(ctx, reqPath, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Book File Update ----------.

// UpdateBookFile updates a single book file.
func (c *Client) UpdateBookFile(ctx context.Context, file *BookFile) (*BookFile, error) {
	var out BookFile
	if err := c.base.Put(ctx, fmt.Sprintf("/api/v1/bookfile/%d", file.ID), file, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditBookFilesBulk updates multiple book files at once.
func (c *Client) EditBookFilesBulk(ctx context.Context, editor *BookFileListResource) ([]BookFile, error) {
	var out []BookFile
	if err := c.base.Put(ctx, "/api/v1/bookfile/editor", editor, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- Bookshelf ----------.

// Bookshelf performs batch monitoring changes on authors and books.
func (c *Client) Bookshelf(ctx context.Context, shelf *BookshelfResource) error {
	return c.base.Post(ctx, "/api/v1/bookshelf", shelf, nil)
}

// ---------- Calendar By ID ----------.

// GetCalendarByID returns a single calendar entry by its ID.
func (c *Client) GetCalendarByID(ctx context.Context, id int) (*Book, error) {
	var out Book
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/calendar/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Wanted By ID ----------.

// GetWantedMissingByID returns a single wanted missing record by its ID.
func (c *Client) GetWantedMissingByID(ctx context.Context, id int) (*Book, error) {
	var out Book
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/wanted/missing/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedCutoffByID returns a single wanted cutoff record by its ID.
func (c *Client) GetWantedCutoffByID(ctx context.Context, id int) (*Book, error) {
	var out Book
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/wanted/cutoff/%d", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Book Overview ----------.

// GetBookOverview returns an overview for a specific book.
func (c *Client) GetBookOverview(ctx context.Context, id int) (*Book, error) {
	var out Book
	if err := c.base.Get(ctx, fmt.Sprintf("/api/v1/book/%d/overview", id), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Search ----------.

// Search searches for authors and books by term.
func (c *Client) Search(ctx context.Context, term string) ([]Author, error) {
	var out []Author
	path := "/api/v1/search?term=" + url.QueryEscape(term)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- HEAD Ping ----------.

// HeadPing performs a lightweight HEAD request to /ping.
func (c *Client) HeadPing(ctx context.Context) error {
	return c.base.Head(ctx, "/ping")
}

// ---------- Backup Upload ----------.

// UploadBackup uploads a backup file via multipart form POST.
func (c *Client) UploadBackup(ctx context.Context, fileName string, data io.Reader) error {
	return c.base.Upload(ctx, "/api/v1/system/backup/upload", "file", fileName, data)
}
