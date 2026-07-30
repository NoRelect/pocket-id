package mds

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-webauthn/webauthn/metadata"
	"github.com/google/uuid"
)

const (
	BlobURL  = "https://mds3.fidoalliance.org/"
	CacheTTL = 72 * time.Hour

	maxBlobSize = 50 * 1024 * 1024

	DefaultCachePath = "data/mds-cache.jwt"
)

type Entry = metadata.Entry

func IsCompromised(e Entry) bool {
	for _, r := range e.StatusReports {
		if metadata.IsUndesiredAuthenticatorStatus(r.Status) {
			return true
		}
	}
	return false
}

var fidoCertStatusOrder = []metadata.AuthenticatorStatus{
	metadata.FidoCertified,
	metadata.FidoCertifiedL1,
	metadata.FidoCertifiedL1plus,
	metadata.FidoCertifiedL2,
	metadata.FidoCertifiedL2plus,
	metadata.FidoCertifiedL3,
	metadata.FidoCertifiedL3plus,
}

func CertificationLevelRank(level string) int {
	return fidoCertStatusRank(metadata.AuthenticatorStatus(level))
}

func fidoCertStatusRank(status metadata.AuthenticatorStatus) int {
	for i, s := range fidoCertStatusOrder {
		if s == status {
			return i + 1
		}
	}
	return 0
}

func HighestCertificationLevel(e Entry) string {
	best := ""
	bestRank := 0
	for _, r := range e.StatusReports {
		if rank := fidoCertStatusRank(r.Status); rank > bestRank {
			best = string(r.Status)
			bestRank = rank
		}
	}
	return best
}

type Service struct {
	httpClient *http.Client
	blobURL    string
	cachePath  string
	decoder    *metadata.Decoder

	mu         sync.RWMutex
	entries    map[string]Entry
	fetchedAt  time.Time
	nextUpdate time.Time

	enforce atomic.Bool
}

func NewService(httpClient *http.Client) *Service {
	decoder, err := metadata.NewDecoder(metadata.WithIgnoreEntryParsingErrors())
	if err != nil {
		panic(fmt.Sprintf("mds: failed to construct metadata decoder: %v", err))
	}

	return &Service{
		httpClient: httpClient,
		blobURL:    BlobURL,
		cachePath:  DefaultCachePath,
		decoder:    decoder,
	}
}

func (s *Service) SetEnforce(enforce bool) {
	s.enforce.Store(enforce)
}

func (s *Service) LoadFromDisk() error {
	return s.loadFromDisk()
}

func (s *Service) loadFromDisk() error {
	if s.cachePath == "" {
		return nil
	}

	raw, err := os.ReadFile(s.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	md, err := s.decode(raw)
	if err != nil {
		slog.Warn("Ignoring FIDO MDS cache file that failed verification",
			slog.String("path", s.cachePath), slog.Any("error", err))
		return nil
	}

	entries := entriesByAAGUID(md)

	fetchedAt := time.Now()
	if info, statErr := os.Stat(s.cachePath); statErr == nil {
		fetchedAt = info.ModTime()
	}

	s.mu.Lock()
	s.entries = entries
	s.fetchedAt = fetchedAt
	s.nextUpdate = md.Parsed.NextUpdate
	s.mu.Unlock()

	slog.Info("Loaded FIDO MDS cache from disk",
		slog.String("path", s.cachePath),
		slog.Int("entries", len(entries)),
		slog.Time("fetchedAt", fetchedAt),
	)
	return nil
}

func (s *Service) saveToDisk(rawJWT []byte) error {
	if s.cachePath == "" {
		return nil
	}

	dir := filepath.Dir(s.cachePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create MDS cache directory: %w", err)
	}

	tmpPath := s.cachePath + ".tmp"
	if err := os.WriteFile(tmpPath, rawJWT, 0600); err != nil {
		return fmt.Errorf("failed to write MDS cache file: %w", err)
	}
	if err := os.Rename(tmpPath, s.cachePath); err != nil {
		return fmt.Errorf("failed to finalize MDS cache file: %w", err)
	}
	return nil
}

func (s *Service) decode(raw []byte) (*metadata.Metadata, error) {
	payload, err := s.decoder.DecodeBytes(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to verify MDS BLOB signature/chain: %w", err)
	}

	md, err := s.decoder.Parse(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MDS BLOB payload: %w", err)
	}

	return md, nil
}

func entriesByAAGUID(md *metadata.Metadata) map[string]Entry {
	entries := make(map[string]Entry, len(md.Parsed.Entries))
	for _, e := range md.Parsed.Entries {
		if e.AaGUID != uuid.Nil {
			entries[strings.ToLower(e.AaGUID.String())] = e
		}
	}
	return entries
}

func (s *Service) Lookup(aaguid string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[strings.ToLower(aaguid)]
	return e, ok
}

func (s *Service) FetchedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fetchedAt
}

func (s *Service) Refresh(ctx context.Context) error {
	slog.InfoContext(ctx, "Refreshing FIDO MDS cache")

	raw, err := s.fetchBlob(ctx)
	if err != nil {
		return fmt.Errorf("mds: fetch failed: %w", err)
	}

	md, err := s.decode(raw)
	if err != nil {
		return fmt.Errorf("mds: verification failed: %w", err)
	}

	entries := entriesByAAGUID(md)

	nextUpdate := md.Parsed.NextUpdate
	if nextUpdate.IsZero() {
		nextUpdate = time.Now().Add(CacheTTL)
	}

	s.mu.Lock()
	s.entries = entries
	s.fetchedAt = time.Now()
	s.nextUpdate = nextUpdate
	s.mu.Unlock()

	if err := s.saveToDisk(raw); err != nil {
		slog.ErrorContext(ctx, "Failed to persist FIDO MDS cache to disk", slog.Any("error", err))
	}

	slog.InfoContext(ctx, "FIDO MDS cache refreshed",
		slog.Int("entries", len(entries)),
		slog.Time("nextUpdate", nextUpdate),
	)
	return nil
}

func (s *Service) ListEntries() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.entries) == 0 {
		return nil
	}
	result := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		result = append(result, e)
	}
	return result
}

func (s *Service) SetEntriesForTest(entries map[string]Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = entries
	s.fetchedAt = time.Now()
}

func (s *Service) fetchBlob(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.blobURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBlobSize))
	if err != nil {
		return nil, err
	}
	return body, nil
}
