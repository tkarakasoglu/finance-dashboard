package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type scannerResponse struct {
	TotalCount int `json:"totalCount"`
	Data       []struct {
		S string `json:"s"`
	} `json:"data"`
}

// SymbolStore holds scanner-backed symbol lists per source, refreshed together on a single hourly schedule.
type SymbolStore struct {
	mu           sync.Mutex
	bySource     map[string][]string
	mergedSorted []string // sorted deduped union; rebuilt only when scanner data changes (not on every HTTP request)
	fetchedAt    time.Time
	refreshing   bool
}

var knownSources = []string{"america", "global", "crypto", "forex", "futures", "bonds"}

var scannerEndpoints = map[string]string{
	"america": "https://scanner.tradingview.com/america/scan",
	"global":  "https://scanner.tradingview.com/global/scan",
	"crypto":  "https://scanner.tradingview.com/crypto/scan",
	"forex":   "https://scanner.tradingview.com/forex/scan",
	"futures": "https://scanner.tradingview.com/futures/scan",
	"bonds":   "https://scanner.tradingview.com/bonds/scan",
}

const symbolCacheTTL = time.Hour

type SymbolStoreStatus struct {
	Count      int       `json:"count"`
	FetchedAt  time.Time `json:"fetchedAt"`
	Refreshing bool      `json:"refreshing"`
}

func (s *SymbolStore) status(source string) SymbolStoreStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bySource == nil {
		s.bySource = map[string][]string{}
	}
	if source == "all" {
		return SymbolStoreStatus{
			Count:      len(s.mergedAllUnlocked()),
			FetchedAt:  s.fetchedAt,
			Refreshing: s.refreshing,
		}
	}
	return SymbolStoreStatus{
		Count:      len(s.bySource[source]),
		FetchedAt:  s.fetchedAt,
		Refreshing: s.refreshing,
	}
}

// mergedAllUnlocked returns the deduped merged list (sorted). Caller must hold s.mu.
func (s *SymbolStore) mergedAllUnlocked() []string {
	if s.mergedSorted != nil {
		return s.mergedSorted
	}
	s.rebuildMergedCacheLocked()
	return s.mergedSorted
}

func (s *SymbolStore) rebuildMergedCacheLocked() {
	seen := make(map[string]struct{}, 50000)
	var merged []string
	for _, src := range knownSources {
		for _, sym := range s.bySource[src] {
			if sym == "" {
				continue
			}
			if _, ok := seen[sym]; ok {
				continue
			}
			seen[sym] = struct{}{}
			merged = append(merged, sym)
		}
	}
	sort.Strings(merged)
	if merged == nil {
		merged = []string{}
	}
	s.mergedSorted = merged
}

func (s *SymbolStore) mergedAll() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mergedAllUnlocked()
}

func (s *SymbolStore) cacheEmptyLocked() bool {
	for _, src := range knownSources {
		if len(s.bySource[src]) > 0 {
			return false
		}
	}
	return true
}

func (s *SymbolStore) cacheStaleLocked() bool {
	if s.cacheEmptyLocked() {
		return true
	}
	return time.Since(s.fetchedAt) > symbolCacheTTL
}

// ensureFreshAsync runs a full multi-source refresh in the background when the merged cache is empty or older than one hour.
func (s *SymbolStore) ensureFreshAsync() {
	s.mu.Lock()
	if s.bySource == nil {
		s.bySource = map[string][]string{}
	}
	if !s.cacheStaleLocked() || s.refreshing {
		s.mu.Unlock()
		return
	}
	s.refreshing = true
	s.mu.Unlock()

	go func() {
		defer func() {
			s.mu.Lock()
			s.refreshing = false
			s.mu.Unlock()
		}()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		s.refreshAll(ctx)
	}()
}

// refreshAll fetches all scanner sources in parallel and updates bySource. Successful responses replace that source;
// on error, the previous slice for that source is kept. fetchedAt is updated if at least one source returned data.
func (s *SymbolStore) refreshAll(ctx context.Context) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	success := make(map[string][]string, len(knownSources))
	for _, src := range knownSources {
		src := src
		ep := scannerEndpoints[src]
		wg.Add(1)
		go func() {
			defer wg.Done()
			syms, err := fetchAllSymbolsFromEndpoint(ctx, ep)
			if err != nil || len(syms) == 0 {
				return
			}
			mu.Lock()
			success[src] = syms
			mu.Unlock()
		}()
	}
	wg.Wait()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bySource == nil {
		s.bySource = map[string][]string{}
	}
	any := false
	for src, syms := range success {
		s.bySource[src] = syms
		any = true
	}
	if any {
		s.fetchedAt = time.Now()
		s.rebuildMergedCacheLocked()
	}
}

func (s *SymbolStore) search(source, q string, limit int) []string {
	q = strings.TrimSpace(strings.ToLower(q))
	if q == "" {
		return nil
	}

	var syms []string
	if source == "all" {
		syms = s.mergedAll()
	} else {
		s.mu.Lock()
		if s.bySource != nil {
			syms = append([]string(nil), s.bySource[source]...)
		}
		s.mu.Unlock()
	}

	type hit struct {
		sym   string
		score int
	}
	hits := make([]hit, 0, limit)
	for _, sym := range syms {
		low := strings.ToLower(sym)
		parts := strings.SplitN(low, ":", 2)
		ticker := low
		if len(parts) == 2 {
			ticker = parts[1]
		}

		score := -1
		switch {
		case low == q || ticker == q:
			score = 100
		case strings.HasPrefix(ticker, q):
			score = 90
		case strings.HasPrefix(low, q):
			score = 80
		case strings.Contains(ticker, q):
			score = 60
		case strings.Contains(low, q):
			score = 50
		}
		if score < 0 {
			continue
		}
		hits = append(hits, hit{sym: sym, score: score})
	}

	sort.Slice(hits, func(i, j int) bool {
		if hits[i].score != hits[j].score {
			return hits[i].score > hits[j].score
		}
		return hits[i].sym < hits[j].sym
	})

	if limit <= 0 {
		limit = 20
	}
	if len(hits) > limit {
		hits = hits[:limit]
	}
	out := make([]string, 0, len(hits))
	for _, h := range hits {
		out = append(out, h.sym)
	}
	return out
}

var globalSymbols = &SymbolStore{}

func init() {
	// Bulk scanner lists are for optional /symbols API only; refresh on a timer so HTTP handlers stay light.
	go func() {
		time.Sleep(200 * time.Millisecond)
		globalSymbols.ensureFreshAsync()
		ticker := time.NewTicker(25 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			globalSymbols.ensureFreshAsync()
		}
	}()
}

// MatchStats counts substring matches per source and across the deduped merged list.
func (s *SymbolStore) MatchStats(q string) (perSource map[string]int, mergedUniqueTotal int) {
	q = strings.TrimSpace(strings.ToLower(q))
	perSource = make(map[string]int)
	if q == "" {
		return perSource, 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, src := range knownSources {
		n := 0
		for _, sym := range s.bySource[src] {
			if strings.Contains(strings.ToLower(sym), q) {
				n++
			}
		}
		perSource[src] = n
	}

	for _, sym := range s.mergedAllUnlocked() {
		if strings.Contains(strings.ToLower(sym), q) {
			mergedUniqueTotal++
		}
	}
	return perSource, mergedUniqueTotal
}

// HasScannerData reports whether at least one TradingView scanner source has been loaded into cache.
func (s *SymbolStore) HasScannerData() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bySource == nil {
		return false
	}
	for _, src := range knownSources {
		if len(s.bySource[src]) > 0 {
			return true
		}
	}
	return false
}

func fetchAllSymbolsFromEndpoint(ctx context.Context, endpoint string) ([]string, error) {
	type payload struct {
		Columns []string `json:"columns"`
		Range   []int    `json:"range"`
	}

	const pageSize = 1000
	const maxPages = 40

	client := &http.Client{Timeout: 15 * time.Second}
	seen := make(map[string]struct{}, 20000)
	var out []string

	for page := 0; page < maxPages; page++ {
		p := payload{
			Columns: []string{},
			Range:   []int{page * pageSize, (page + 1) * pageSize},
		}
		b, _ := json.Marshal(p)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Origin", "https://scanner.tradingview.com")
		req.Header.Set("Referer", "https://scanner.tradingview.com/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		body, _ := ioReadAllLimit(res.Body, 10<<20)
		_ = res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, errors.New("scanner returned non-2xx")
		}

		var sr scannerResponse
		if err := json.Unmarshal(body, &sr); err != nil {
			return nil, err
		}
		for _, row := range sr.Data {
			if row.S == "" {
				continue
			}
			if _, ok := seen[row.S]; ok {
				continue
			}
			seen[row.S] = struct{}{}
			out = append(out, row.S)
		}

		if len(sr.Data) == 0 {
			break
		}
		if sr.TotalCount > 0 && len(out) >= sr.TotalCount {
			break
		}
	}

	sort.Strings(out)
	return out, nil
}

func ioReadAllLimit(r io.Reader, max int64) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(io.LimitReader(r, max)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
