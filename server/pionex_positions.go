package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const pionexAPIBase = "https://api.pionex.com"

// pionexSignGET builds PIONEX-SIGNATURE for private GET requests.
// See: https://github.com/pionex-official/pionex-open-api (Futures + Spot use the same scheme).
func pionexSignGET(secret, path string, extra url.Values) (queryString string, signature string, err error) {
	if extra == nil {
		extra = url.Values{}
	}
	ts := time.Now().UnixMilli()
	extra.Set("timestamp", strconv.FormatInt(ts, 10))

	keys := make([]string, 0, len(extra))
	for k := range extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+extra.Get(k))
	}
	q := strings.Join(parts, "&")
	pathURL := path + "?" + q
	payload := "GET" + pathURL
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return q, hex.EncodeToString(mac.Sum(nil)), nil
}

func pionexNonZeroSize(s string) bool {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return false
	}
	return f != 0
}

type pionexPositionsResp struct {
	Result    bool   `json:"result"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Data      *struct {
		Positions []struct {
			Symbol           string `json:"symbol"`
			PositionSide     string `json:"positionSide"`
			NetSize          string `json:"netSize"`
			AvgPrice         string `json:"avgPrice"`
			MarkPrice        string `json:"markPrice"`
			UnrealizedPnL    string `json:"unrealizedPnL"`
			Leverage         string `json:"leverage"`
			LiquidationPrice string `json:"liquidationPrice"`
		} `json:"positions"`
	} `json:"data"`
}

func fetchPionexPositions(apiKey, apiSecret string) ([]position, error) {
	apiKey = strings.TrimSpace(apiKey)
	apiSecret = strings.TrimSpace(apiSecret)
	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing api key or secret")
	}

	path := "/uapi/v1/account/positions"
	q, sig, err := pionexSignGET(apiSecret, path, nil)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(pionexAPIBase + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PIONEX-KEY", apiKey)
	req.Header.Set("PIONEX-SIGNATURE", sig)

	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var raw pionexPositionsResp
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}
	if !raw.Result {
		msg := strings.TrimSpace(raw.Message)
		if msg == "" {
			msg = raw.Code
		}
		if msg == "" {
			msg = "unknown error"
		}
		return nil, fmt.Errorf("%s", msg)
	}
	if raw.Data == nil {
		return []position{}, nil
	}

	out := make([]position, 0, len(raw.Data.Positions))
	for _, p := range raw.Data.Positions {
		if !pionexNonZeroSize(p.NetSize) {
			continue
		}
		sym := strings.TrimSpace(p.Symbol)
		if sym == "" {
			continue
		}
		out = append(out, position{
			Exchange:         "PIONEX",
			Symbol:           sym,
			Side:             strings.TrimSpace(p.PositionSide),
			Size:             strings.TrimSpace(p.NetSize),
			AvgPrice:         strings.TrimSpace(p.AvgPrice),
			MarkPrice:        strings.TrimSpace(p.MarkPrice),
			UnrealizedPnL:    strings.TrimSpace(p.UnrealizedPnL),
			Leverage:         strings.TrimSpace(p.Leverage),
			LiquidationPrice: strings.TrimSpace(p.LiquidationPrice),
		})
	}
	return out, nil
}
