package anidb

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const titleDumpURL = "http://anidb.net/api/anime-titles.dat.gz"

// TitleEntry represents a single entry from the AniDB anime titles dump.
type TitleEntry struct {
	AID   int    // Anime ID
	Type  string // Title type: "1" (primary), "2" (synonym), "3" (short), "4" (official)
	Lang  string // Language code (e.g. "en", "ja", "x-jat")
	Title string // The title text
}

// TitleMatch is a search result from the title dump with a match quality score.
type TitleMatch struct {
	AID   int     // Anime ID
	Type  string  // Title type
	Lang  string  // Language code
	Title string  // Matched title
	Score float64 // Match quality: 1.0 = exact, 0.75 = prefix, 0.5 = contains
}

// titleStore holds the in-memory title dump.
type titleStore struct {
	mu      sync.RWMutex
	entries []TitleEntry
}

// LoadTitleDump downloads and parses the AniDB anime titles dump.
// The dump is cached in memory. Call this before using SearchByTitle.
func (c *Client) LoadTitleDump(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, titleDumpURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("anidb: create title dump request: %w", err)
	}
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("anidb: download title dump: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("anidb: title dump HTTP %d", resp.StatusCode)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("anidb: decompress title dump: %w", err)
	}
	defer gz.Close()

	entries, err := ParseTitleDump(gz)
	if err != nil {
		return err
	}

	c.titles.mu.Lock()
	c.titles.entries = entries
	c.titles.mu.Unlock()

	return nil
}

// SearchByTitle searches the loaded title dump for anime matching the query.
// Results are ranked: exact match (1.0) > prefix match (0.75) > contains match (0.5).
// Returns up to limit results sorted by score descending.
// LoadTitleDump must be called first; returns an error if the dump is not loaded.
func (c *Client) SearchByTitle(query string, limit int) ([]TitleMatch, error) {
	c.titles.mu.RLock()
	entries := c.titles.entries
	c.titles.mu.RUnlock()

	if entries == nil {
		return nil, errors.New("anidb: title dump not loaded, call LoadTitleDump first")
	}

	lowerQuery := strings.ToLower(query)
	// Track best match per AID to avoid duplicates.
	best := make(map[int]TitleMatch)

	for _, e := range entries {
		lowerTitle := strings.ToLower(e.Title)
		var score float64
		switch {
		case lowerTitle == lowerQuery:
			score = 1.0
		case strings.HasPrefix(lowerTitle, lowerQuery):
			score = 0.75
		case strings.Contains(lowerTitle, lowerQuery):
			score = 0.5
		default:
			continue
		}

		if existing, ok := best[e.AID]; !ok || score > existing.Score {
			best[e.AID] = TitleMatch{
				AID:   e.AID,
				Type:  e.Type,
				Lang:  e.Lang,
				Title: e.Title,
				Score: score,
			}
		}
	}

	results := make([]TitleMatch, 0, len(best))
	for _, m := range best {
		results = append(results, m)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// ParseTitleDump parses an AniDB anime titles dump from the given reader.
// This is exported for consumers who want to provide their own dump source.
func ParseTitleDump(r io.Reader) ([]TitleEntry, error) {
	var entries []TitleEntry
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		aid, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		entries = append(entries, TitleEntry{
			AID:   aid,
			Type:  parts[1],
			Lang:  parts[2],
			Title: parts[3],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("anidb: parse title dump: %w", err)
	}

	return entries, nil
}
