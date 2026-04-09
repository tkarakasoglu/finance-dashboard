# finance-dashboard

**Vibe-coded:** this repo was built in a loose, exploratory way—quick iterations, minimal ceremony, “make it work” energy rather than a formal product or architecture review. Expect pragmatic shortcuts and embedded UI, not a blueprint for enterprise production.

A minimal **Go** web app that serves a browser-based finance dashboard. The UI is a single embedded HTML page with a dark layout and a **TradingView** multi-chart grid. The same HTTP handlers run as a long-lived server locally and as **Vercel** serverless functions in production.

## What it does

- **Dashboard** — Open `/` in a browser to use the chart workspace (symbol search, intervals, and a tile grid backed by TradingView’s client-side widgets).
- **Symbol search proxy** — `GET /symbol-search` (and `GET /api/symbol-search` locally) forwards queries to TradingView’s symbol search API with browser-like headers so the UI can autocomplete without browser **CORS** blocking direct calls.
- **Scanner-backed symbols** — `GET /symbols` exposes searchable lists aggregated from TradingView scanner endpoints (e.g. US, global, crypto, forex, futures, bonds). Lists are refreshed on a schedule (about hourly) and merged/deduplicated for fast lookup.

TradingView is a third-party service; charts and search depend on their scripts and APIs being available and allowed by your network.

## Requirements

- [Go](https://go.dev/dl/) **1.23+**

## Quick start

```bash
make run
```

Then open:

| URL | Description |
|-----|-------------|
| [http://localhost:8080/](http://localhost:8080/) | Dashboard |
| [http://localhost:8080/healthz](http://localhost:8080/healthz) | JSON health check |
| [http://localhost:8080/symbols?q=](http://localhost:8080/symbols?q=) | Symbol list search (optional `source`, `limit`, etc.) |
| [http://localhost:8080/symbol-search?text=](http://localhost:8080/symbol-search?text=) | Proxied symbol search JSON |

To load config from a file:

```bash
cp .env.example .env
export $(cat .env | xargs) && make run
```

## Configuration

| Variable | Default | Meaning |
|----------|---------|---------|
| `HTTP_ADDR` | `:8080` | Listen address for `cmd/server` |

## Project layout

- **`cmd/server`** — Standard HTTP server with graceful shutdown on `SIGINT` / `SIGTERM`.
- **`server`** — Route registration, embedded dashboard HTML (`ui.go`), symbol store, and TradingView proxy logic.
- **`api/`** — Thin Vercel handlers that delegate to the same `server` package (`index`, `healthz`, `symbol-search`, `symbols`).

## Vercel

The repo includes `vercel.json` so paths like `/`, `/healthz`, and `/symbol-search` are rewritten to the matching `/api/*` serverless routes. After deploy, use the same paths on your deployment host (e.g. `https://<project>.vercel.app/healthz`).

## Development

```bash
make build    # build binary to bin/server
make test     # go test ./...
make tidy     # go mod tidy
make lint     # go vet ./...
```

## Disclaimer

This project is for learning and personal tooling. Market data and charts are provided by **TradingView** and related services under their terms. This repository is not affiliated with TradingView.
