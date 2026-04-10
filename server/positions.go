package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type positionsRequest struct {
	Exchanges map[string]exchangeCredentials `json:"exchanges"`
}

type exchangeCredentials struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

type position struct {
	Exchange         string `json:"exchange"`
	Symbol           string `json:"symbol"`
	Side             string `json:"side,omitempty"`
	Size             string `json:"size,omitempty"`
	AvgPrice         string `json:"avgPrice,omitempty"`
	MarkPrice        string `json:"markPrice,omitempty"`
	UnrealizedPnL    string `json:"unrealizedPnL,omitempty"`
	Leverage         string `json:"leverage,omitempty"`
	LiquidationPrice string `json:"liquidationPrice,omitempty"`
}

type positionsResponse struct {
	Positions []position `json:"positions"`
	Errors    []string   `json:"errors,omitempty"`
	UpdatedAt string     `json:"updatedAt"`
}

// Positions returns open positions from supported exchanges (see handlers in this package).
func Positions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req positionsRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	var errs []string
	var out []position

	for name, cred := range req.Exchanges {
		n := strings.TrimSpace(name)
		if n == "" {
			continue
		}
		if strings.TrimSpace(cred.APIKey) == "" && strings.TrimSpace(cred.APISecret) == "" {
			continue
		}
		if strings.TrimSpace(cred.APIKey) == "" || strings.TrimSpace(cred.APISecret) == "" {
			errs = append(errs, n+": apiKey and apiSecret must both be set")
			continue
		}

		switch strings.ToLower(n) {
		case "pionex":
			pp, err := fetchPionexPositions(cred.APIKey, cred.APISecret)
			if err != nil {
				errs = append(errs, "pionex: "+err.Error())
				continue
			}
			out = append(out, pp...)
		default:
			errs = append(errs, n+": positions not implemented yet")
		}
	}

	_ = json.NewEncoder(w).Encode(positionsResponse{
		Positions: out,
		Errors:    errs,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}
