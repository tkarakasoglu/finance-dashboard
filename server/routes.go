package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", Healthz)
	mux.HandleFunc("GET /symbols", Symbols)

	// Proxy TradingView symbol search to avoid browser CORS issues (primary path for UI autocomplete).
	mux.HandleFunc("GET /symbol-search", SymbolSearch)
	mux.HandleFunc("GET /api/symbol-search", SymbolSearch)

	// Exchange positions (used by auto charts mode).
	mux.HandleFunc("POST /positions", Positions)
	mux.HandleFunc("POST /api/positions", Positions)

	mux.HandleFunc("GET /", Index)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

func SymbolSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("text")
	if q == "" {
		q = r.URL.Query().Get("q")
	}
	if q == "" {
		http.Error(w, "missing query parameter: text", http.StatusBadRequest)
		return
	}

	u := url.URL{
		Scheme: "https",
		Host:   "symbol-search.tradingview.com",
		Path:   "/symbol_search/",
	}
	qs := u.Query()
	qs.Set("text", q)
	qs.Set("lang", "en")
	// hl=1 wraps matches in <em>…</em> inside JSON fields; we need plain tickers for the widget.
	qs.Set("hl", "0")
	u.RawQuery = qs.Encode()

	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, u.String(), nil)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// Browser-like headers: bare clients often get 403 from TradingView edges.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.tradingview.com")
	req.Header.Set("Referer", "https://www.tradingview.com/")

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "upstream request failed", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	w.Header().Set("Content-Type", "application/json")
	// Pass through upstream status/body for easier debugging (and to avoid hiding rate-limit/captcha responses).
	w.WriteHeader(res.StatusCode)
	_, _ = w.Write(body)
}

func Symbols(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		q = r.URL.Query().Get("text")
	}
	source := r.URL.Query().Get("source")
	if source == "" {
		source = "all"
	}
	limit := 20
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}

	st := globalSymbols.status(source)
	results := globalSymbols.search(source, q, limit)

	type item struct {
		Full     string `json:"full"`
		Exchange string `json:"exchange"`
		Symbol   string `json:"symbol"`
	}
	out := make([]item, 0, len(results))
	for _, full := range results {
		ex, sym, ok := strings.Cut(full, ":")
		if !ok {
			ex, sym = "", full
		}
		out = append(out, item{Full: full, Exchange: ex, Symbol: sym})
	}

	perSource := map[string]int{}
	for _, src := range []string{"america", "global", "crypto", "forex", "futures", "bonds"} {
		perSource[src] = globalSymbols.status(src).Count
	}

	resp := map[string]any{
		"source": source,
		"q":      q,
		// warming: no scanner-backed symbols yet and a background refresh is running
		"warming": !globalSymbols.HasScannerData() && st.Refreshing,
		// scannerReady: at least one TradingView scanner list was loaded
		"scannerReady": globalSymbols.HasScannerData(),
		"perSource":    perSource,
		"status":       st,
		"count":        len(out),
		"data":         out,
	}
	// Expensive full-list stats; opt-in only (e.g. debugging).
	if strings.TrimSpace(q) != "" && r.URL.Query().Get("stats") == "1" {
		msPer, msTotal := globalSymbols.MatchStats(q)
		resp["matchStats"] = map[string]any{
			"perSource":         msPer,
			"mergedUniqueTotal": msTotal,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
