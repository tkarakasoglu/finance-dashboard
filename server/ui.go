package server

// indexHTML is intentionally a single-file UI to keep the project minimal:
// Go serves this HTML, and the browser renders the dashboard + embeds TradingView widgets.
const indexHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>finance-dashboard</title>
    <style>
      :root {
        --bg: #0b0f14;
        --panel: #0f1621;
        --panel2: #0c121b;
        --text: #e7eef8;
        --muted: #9db0c8;
        --border: rgba(231, 238, 248, 0.12);
        --accent: #66a3ff;
        --shadow: 0 10px 30px rgba(0, 0, 0, 0.35);
        --radius: 12px;
      }
      * { box-sizing: border-box; }
      html, body { height: 100%; }
      body {
        margin: 0;
        background: radial-gradient(1200px 800px at 20% 0%, rgba(102, 163, 255, 0.14), transparent 55%),
                    radial-gradient(900px 600px at 100% 10%, rgba(115, 255, 227, 0.08), transparent 45%),
                    var(--bg);
        color: var(--text);
        font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial, "Apple Color Emoji", "Segoe UI Emoji";
      }
      a { color: inherit; text-decoration: none; }

      .app {
        height: 100%;
        display: grid;
        grid-template-columns: 260px 1fr;
      }
      .app.sidebar-hidden {
        grid-template-columns: 1fr;
      }
      .app.sidebar-hidden .sidebar {
        display: none;
      }

      .sidebar {
        border-right: 1px solid var(--border);
        background: linear-gradient(180deg, rgba(15, 22, 33, 0.88), rgba(12, 18, 27, 0.88));
        backdrop-filter: blur(10px);
        padding: 18px 14px;
      }

      .brand {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px 10px 14px;
        margin-bottom: 10px;
      }
      .logo {
        width: 34px;
        height: 34px;
        border-radius: 10px;
        background: linear-gradient(135deg, rgba(102,163,255,1), rgba(115,255,227,0.65));
        box-shadow: var(--shadow);
      }
      .brand h1 {
        font-size: 14px;
        line-height: 1.2;
        margin: 0;
        letter-spacing: 0.3px;
      }
      .brand p {
        margin: 2px 0 0;
        color: var(--muted);
        font-size: 12px;
      }

      .nav {
        display: flex;
        flex-direction: column;
        gap: 6px;
        margin-top: 10px;
      }
      .nav-item {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px 10px;
        border-radius: 10px;
        color: var(--muted);
        border: 1px solid transparent;
      }
      .nav-item[data-active="true"] {
        color: var(--text);
        background: rgba(102, 163, 255, 0.10);
        border-color: rgba(102, 163, 255, 0.24);
      }
      .nav-dot {
        width: 8px;
        height: 8px;
        border-radius: 999px;
        background: rgba(231, 238, 248, 0.35);
      }
      .nav-item[data-active="true"] .nav-dot {
        background: var(--accent);
      }

      .content {
        padding: 18px 18px 22px;
        overflow-x: hidden;
        overflow-y: auto;
      }
      .topbar {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 10px;
        margin-bottom: 14px;
      }
      .topbar-head {
        display: flex;
        align-items: flex-start;
        gap: 10px;
        flex: 1;
        min-width: 0;
      }
      .sidebar-toggle {
        flex-shrink: 0;
        width: 36px;
        height: 36px;
        margin-top: 2px;
        border-radius: 10px;
        border: 1px solid rgba(102, 163, 255, 0.28);
        background: rgba(102, 163, 255, 0.12);
        color: var(--text);
        cursor: pointer;
        font-size: 16px;
        line-height: 1;
        display: inline-flex;
        align-items: center;
        justify-content: center;
      }
      .sidebar-toggle:hover {
        background: rgba(102, 163, 255, 0.18);
      }
      .title {
        font-size: 16px;
        margin: 0;
        letter-spacing: 0.2px;
      }
      .subtitle {
        margin: 4px 0 0;
        color: var(--muted);
        font-size: 12px;
      }

      .topbar-actions {
        flex-shrink: 0;
        display: flex;
        align-items: center;
        gap: 10px;
      }
      .topbar-default-interval {
        font-size: 13px;
        color: var(--muted);
        white-space: nowrap;
        line-height: 1.25;
        letter-spacing: 0.2px;
      }
      .settings-btn {
        width: 40px;
        height: 40px;
        margin-top: 0;
        border-radius: 10px;
        border: 1px solid rgba(102, 163, 255, 0.28);
        background: rgba(102, 163, 255, 0.12);
        color: var(--text);
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        padding: 0;
      }
      .settings-btn:hover {
        background: rgba(102, 163, 255, 0.18);
      }
      .settings-btn svg {
        width: 20px;
        height: 20px;
        display: block;
      }

      .settings-dialog {
        border: 1px solid var(--border);
        border-radius: var(--radius);
        background: linear-gradient(180deg, rgba(15, 22, 33, 0.98), rgba(12, 18, 27, 0.98));
        color: var(--text);
        padding: 0;
        max-width: min(420px, calc(100vw - 32px));
        box-shadow: var(--shadow);
      }
      .settings-dialog::backdrop {
        background: rgba(0, 0, 0, 0.55);
      }
      .settings-dialog-inner {
        padding: 16px 16px 14px;
      }
      .settings-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        margin-bottom: 14px;
      }
      .settings-title {
        margin: 0;
        font-size: 15px;
        font-weight: 600;
        letter-spacing: 0.2px;
      }
      .settings-close {
        width: 32px;
        height: 32px;
        border-radius: 8px;
        border: 1px solid transparent;
        background: transparent;
        color: var(--muted);
        cursor: pointer;
        font-size: 22px;
        line-height: 1;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .settings-close:hover {
        color: var(--text);
        background: rgba(231, 238, 248, 0.06);
      }
      .settings-body {
        display: grid;
        gap: 12px;
      }
      .settings-footer {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 16px;
        padding-top: 14px;
        border-top: 1px solid var(--border);
      }
      .settings-footer button {
        height: 34px;
        border-radius: 10px;
        border: 1px solid var(--border);
        padding: 0 14px;
        cursor: pointer;
        font-size: 13px;
      }
      .settings-footer .btn-secondary {
        background: rgba(12, 18, 27, 0.7);
        color: var(--text);
      }
      .settings-footer .btn-primary {
        background: rgba(102, 163, 255, 0.12);
        border-color: rgba(102, 163, 255, 0.28);
        color: var(--text);
      }
      .settings-footer .btn-primary:hover {
        background: rgba(102, 163, 255, 0.18);
      }

      .controls {
        display: flex;
        gap: 10px;
        align-items: end;
        flex-wrap: wrap;
        background: rgba(15, 22, 33, 0.55);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 12px;
        backdrop-filter: blur(10px);
      }
      .control {
        display: grid;
        gap: 6px;
      }
      .control label {
        font-size: 11px;
        color: var(--muted);
      }
      .control input, .control button {
        height: 34px;
        border-radius: 10px;
        border: 1px solid var(--border);
        background: rgba(12, 18, 27, 0.7);
        color: var(--text);
        padding: 0 10px;
        outline: none;
      }
      .control select {
        height: 34px;
        border-radius: 10px;
        border: 1px solid var(--border);
        background: rgba(12, 18, 27, 0.7);
        color: var(--text);
        padding: 0 10px;
        outline: none;
        appearance: none;
      }
      .control input { width: 120px; }
      .control button {
        cursor: pointer;
        background: rgba(102, 163, 255, 0.12);
        border-color: rgba(102, 163, 255, 0.28);
      }
      .control button:hover {
        background: rgba(102, 163, 255, 0.18);
      }

      /* Logical order uses CSS order property so DnD can swap without moving DOM (moving iframes reloads them). */
      .grid {
        margin-top: 12px;
        display: grid;
        gap: 12px;
        --chart-tile-min: 260px;
      }
      .tile {
        min-height: var(--chart-tile-min, 260px);
        border-radius: var(--radius);
        border: 1px solid var(--border);
        /* allow symbol dropdown to extend outside the tile */
        overflow: visible;
        background: rgba(15, 22, 33, 0.35);
        box-shadow: 0 1px 0 rgba(255,255,255,0.05) inset;
      }
      .tile.dragging {
        opacity: 0.6;
      }
      .tile.drop-target {
        outline: 2px solid rgba(102, 163, 255, 0.55);
        outline-offset: 2px;
      }
      .tile-head {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 8px;
        padding: 10px 10px;
        border-bottom: 1px solid rgba(231, 238, 248, 0.08);
        background: rgba(12, 18, 27, 0.45);
      }
      .tile-drag-handle {
        flex-shrink: 0;
        width: 28px;
        height: 32px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--muted);
        border: 1px solid transparent;
        cursor: grab;
        user-select: none;
        -webkit-user-select: none;
      }
      .tile-drag-handle:hover {
        color: var(--text);
        border-color: rgba(231, 238, 248, 0.12);
        background: rgba(231, 238, 248, 0.04);
      }
      .tile-drag-handle:active {
        cursor: grabbing;
      }
      .tile-drag-grip {
        display: block;
        width: 10px;
        height: 14px;
        opacity: 0.85;
        background: radial-gradient(circle, currentColor 1px, transparent 1.3px);
        background-size: 4px 4px;
      }
      .tile-head .left {
        display: flex;
        align-items: center;
        gap: 8px;
        min-width: 0;
        flex: 1;
      }
      .tile-menu-wrap {
        position: relative;
        flex-shrink: 0;
      }
      .tile-menu-btn {
        width: 32px;
        height: 32px;
        padding: 0;
        border-radius: 8px;
        border: 1px solid var(--border);
        background: rgba(12, 18, 27, 0.7);
        color: var(--muted);
        cursor: pointer;
        font-size: 18px;
        line-height: 1;
      }
      .tile-menu-btn:hover {
        color: var(--text);
        border-color: rgba(231, 238, 248, 0.22);
      }
      .tile-menu-panel {
        position: absolute;
        right: 0;
        top: calc(100% + 4px);
        z-index: 50;
        min-width: 132px;
        border-radius: 10px;
        border: 1px solid var(--border);
        background: rgba(12, 18, 27, 0.98);
        backdrop-filter: blur(8px);
        padding: 4px;
        box-shadow: var(--shadow);
        display: none;
      }
      .tile-menu-panel[data-open="true"] {
        display: block;
      }
      .tile-menu-item {
        display: block;
        width: 100%;
        text-align: left;
        padding: 8px 10px;
        border: none;
        border-radius: 8px;
        background: transparent;
        color: var(--text);
        font-size: 13px;
        cursor: pointer;
      }
      .tile-menu-item:hover {
        background: rgba(255, 80, 80, 0.12);
        color: #ff8a8a;
      }
      .add-chart-footer {
        margin-top: 14px;
        display: flex;
        justify-content: center;
      }
      .add-chart-btn {
        min-width: 140px;
        height: 44px;
        border-radius: 12px;
        border: 1px dashed rgba(102, 163, 255, 0.35);
        background: rgba(102, 163, 255, 0.08);
        color: var(--accent);
        font-size: 22px;
        line-height: 1;
        cursor: pointer;
      }
      .add-chart-btn:hover {
        background: rgba(102, 163, 255, 0.14);
        border-color: rgba(102, 163, 255, 0.5);
      }
      .pill {
        font-size: 11px;
        color: var(--muted);
        border: 1px solid rgba(231, 238, 248, 0.14);
        padding: 4px 8px;
        border-radius: 999px;
        white-space: nowrap;
      }
      .symbol-input {
        height: 28px;
        width: 140px;
        border-radius: 10px;
        border: 1px solid rgba(231, 238, 248, 0.14);
        background: rgba(11, 15, 20, 0.5);
        color: var(--text);
        padding: 0 10px;
        outline: none;
      }
      .symbol-search {
        position: relative;
      }
      .symbol-search .symbol-input {
        width: 220px;
      }
      /* Position (top/left/width) set in JS — fixed to escape .content overflow clipping */
      .suggestions {
        position: fixed;
        z-index: 2147483000;
        max-height: min(260px, 40vh);
        overflow: auto;
        border-radius: 12px;
        border: 1px solid rgba(231, 238, 248, 0.14);
        background: rgba(12, 18, 27, 0.98);
        backdrop-filter: blur(8px);
        box-shadow: var(--shadow);
      }
      .suggestions .suggestion-msg {
        padding: 12px 12px;
        font-size: 12px;
        color: var(--muted);
      }
      .suggestions .suggestion-footer {
        padding: 8px 10px;
        font-size: 11px;
        line-height: 1.35;
        color: var(--muted);
        border-top: 1px solid rgba(231, 238, 248, 0.08);
        flex-shrink: 0;
      }
      .suggestion {
        display: grid;
        gap: 2px;
        padding: 10px 10px;
        cursor: pointer;
        border-bottom: 1px solid rgba(231, 238, 248, 0.06);
      }
      .suggestion:last-child {
        border-bottom: none;
      }
      .suggestion:hover {
        background: rgba(102, 163, 255, 0.10);
      }
      .suggestion .top {
        display: flex;
        gap: 8px;
        align-items: baseline;
        min-width: 0;
      }
      .suggestion .sym {
        font-size: 12px;
        white-space: nowrap;
      }
      .suggestion .ex {
        font-size: 11px;
        color: var(--muted);
        white-space: nowrap;
      }
      .suggestion .desc {
        font-size: 11px;
        color: var(--muted);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      .interval-select {
        height: 28px;
        border-radius: 10px;
        border: 1px solid rgba(231, 238, 248, 0.14);
        background: rgba(11, 15, 20, 0.5);
        color: var(--text);
        padding: 0 8px;
        outline: none;
        appearance: none;
      }
      .tile-body {
        height: 100%;
        overflow: hidden;
        border-radius: 0 0 var(--radius) var(--radius);
      }

      .empty {
        border: 1px dashed rgba(231, 238, 248, 0.18);
        border-radius: var(--radius);
        padding: 16px;
        color: var(--muted);
        margin-top: 12px;
        background: rgba(15, 22, 33, 0.35);
      }

      @media (max-width: 980px) {
        .app { grid-template-columns: 1fr; }
        .sidebar { position: sticky; top: 0; z-index: 5; border-right: none; border-bottom: 1px solid var(--border); }
      }
    </style>
  </head>
  <body>
    <div class="app" id="appRoot">
      <aside class="sidebar" id="sidebar">
        <div class="brand">
          <div class="logo" aria-hidden="true"></div>
          <div>
            <h1>finance-dashboard</h1>
            <p>minimal • multichart</p>
          </div>
        </div>

        <nav class="nav" id="nav">
          <a class="nav-item" href="#multichart" data-page="multichart">
            <span class="nav-dot" aria-hidden="true"></span>
            <span>Multichart</span>
          </a>
        </nav>
      </aside>

      <main class="content">
        <div class="topbar">
          <div class="topbar-head">
            <button type="button" class="sidebar-toggle" id="sidebarToggle" aria-controls="sidebar" aria-expanded="true" title="Hide sidebar">◀</button>
            <div>
              <h2 class="title" id="pageTitle">Multichart</h2>
              <div class="subtitle" id="pageSubtitle">TradingView chart tile grid</div>
            </div>
          </div>

          <div class="topbar-actions" id="topbarActions">
            <span class="topbar-default-interval" id="topbarDefaultInterval" title="Default chart interval"></span>
            <button type="button" class="settings-btn" id="settingsOpen" aria-haspopup="dialog" aria-controls="settingsDialog" title="Chart settings">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z" />
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
              </svg>
            </button>
          </div>
        </div>

        <dialog id="settingsDialog" class="settings-dialog" aria-labelledby="settingsTitle">
          <div class="settings-dialog-inner">
            <div class="settings-header">
              <h3 class="settings-title" id="settingsTitle">Chart settings</h3>
              <button type="button" class="settings-close" id="settingsClose" aria-label="Close">×</button>
            </div>
            <div class="settings-body">
              <div class="control">
                <label for="cols">Columns</label>
                <input id="cols" type="number" min="1" max="6" step="1" />
              </div>
              <div class="control">
                <label for="defaultInterval">Default interval</label>
                <select id="defaultInterval"></select>
              </div>
              <div class="control">
                <label for="chartHeight">Chart height (px)</label>
                <input id="chartHeight" type="number" min="180" max="900" step="10" />
              </div>
              <div class="control">
                <label for="chartTheme">Chart color scheme</label>
                <select id="chartTheme"></select>
              </div>
              <div class="control">
                <label for="barColorPreset">Bar / candle colors</label>
                <select id="barColorPreset"></select>
              </div>
            </div>
            <div class="settings-footer">
              <button type="button" class="btn-secondary" id="settingsCancel">Cancel</button>
              <button type="button" class="btn-primary" id="apply">Apply</button>
            </div>
          </div>
        </dialog>

        <div id="root"></div>
      </main>
    </div>

    <script>
      (function () {
        const STORAGE_KEY = "finance-dashboard:multichart:v2";
        const SIDEBAR_STORAGE_KEY = "finance-dashboard:sidebar:collapsed";

        function loadSidebarCollapsed() {
          try {
            return sessionStorage.getItem(SIDEBAR_STORAGE_KEY) === "1";
          } catch (_) {
            return false;
          }
        }

        function saveSidebarCollapsed(collapsed) {
          try {
            sessionStorage.setItem(SIDEBAR_STORAGE_KEY, collapsed ? "1" : "0");
          } catch (_) {}
        }

        function applySidebarCollapsed(collapsed) {
          const app = document.getElementById("appRoot");
          const btn = document.getElementById("sidebarToggle");
          const aside = document.getElementById("sidebar");
          if (!app || !btn) return;
          app.classList.toggle("sidebar-hidden", collapsed);
          btn.setAttribute("aria-expanded", collapsed ? "false" : "true");
          btn.title = collapsed ? "Show sidebar" : "Hide sidebar";
          btn.textContent = collapsed ? "▶" : "◀";
          if (aside) aside.setAttribute("aria-hidden", collapsed ? "true" : "false");
        }
        // TradingView widget intervals:
        // - Minutes as numbers (e.g. "1", "15", "60")
        // - Days/Weeks/Months as "D", "W", "M"
        // - Some multi-day/month variants are also accepted in many contexts (e.g. "2D", "3M")
        const INTERVAL_OPTIONS = [
          // Minutes
          "1", "3", "5", "15", "30", "45",
          // Hours (in minutes)
          "60", "120", "180", "240",
          // Days / Weeks / Months
          "D", "2D", "3D",
          "W",
          "M", "3M", "6M", "12M"
        ];

        function intervalLabel(v) {
          if (v === "D") return "1D";
          if (v === "2D") return "2D";
          if (v === "3D") return "3D";
          if (v === "W") return "1W";
          if (v === "M") return "1M";
          if (v === "3M") return "3M";
          if (v === "6M") return "6M";
          if (v === "12M") return "12M";
          const n = Number.parseInt(String(v), 10);
          if (!Number.isNaN(n)) {
            if (n % 60 === 0 && n >= 60) return (n / 60) + "hr";
            return n + "m";
          }
          return String(v);
        }

        function clampInt(n, min, max, fallback) {
          const x = Number.parseInt(String(n), 10);
          if (Number.isNaN(x)) return fallback;
          return Math.max(min, Math.min(max, x));
        }

        function loadState() {
          try {
            const raw = sessionStorage.getItem(STORAGE_KEY);
            if (raw) return JSON.parse(raw);
            const legacy = sessionStorage.getItem("finance-dashboard:multichart:v1");
            if (legacy) return JSON.parse(legacy);
            return null;
          } catch (_) {
            return null;
          }
        }

        function saveState(state) {
          try {
            sessionStorage.setItem(STORAGE_KEY, JSON.stringify(state));
          } catch (_) {}
        }

        function debounce(fn, delayMs) {
          let t = null;
          return function (...args) {
            if (t) clearTimeout(t);
            t = setTimeout(() => fn.apply(this, args), delayMs);
          };
        }

        const searchCache = new Map();

        function stripHtmlTags(s) {
          return String(s || "").replace(/<[^>]*>/g, "");
        }

        /** Map TradingView symbol_search row → { full, exchange, symbol } for the chart widget. */
        function mapTvSymbolRow(row) {
          if (!row || typeof row !== "object") return null;
          const full = stripHtmlTags(row.full_name || "").trim();
          const symbol = stripHtmlTags(row.symbol || "").trim();
          const exchange = stripHtmlTags(row.exchange || "").trim();
          let ex = exchange;
          let sym = symbol;
          if (full.includes(":")) {
            const parts = full.split(":");
            if (!ex) ex = parts[0] || "";
            if (!sym) sym = parts.slice(1).join(":") || "";
          }
          const mergedFull = full || (ex && sym ? ex + ":" + sym : sym);
          if (!mergedFull) return null;
          return { full: mergedFull, exchange: ex, symbol: sym || mergedFull };
        }

        /**
         * Autocomplete uses TradingView's symbol_search (server-proxied). That is indexed/ranked upstream —
         * not a linear scan over a bulk scanner export.
         */
        async function searchSymbols(text, opts) {
          const q = String(text || "").trim();
          if (q.length < 2) return { items: [] };
          const key = "plain:" + q.toLowerCase();
          const bypassCache = !!(opts && opts.bypassCache);
          if (!bypassCache && searchCache.has(key)) return searchCache.get(key);

          const url = "/api/symbol-search?" + new URLSearchParams({ text: q, lang: "en", hl: "0" }).toString();
          const init = { method: "GET", cache: "no-store" };
          if (opts && opts.signal) init.signal = opts.signal;
          const res = await fetch(url, init);
          if (!res.ok) return { items: [] };
          const raw = await res.json();
          const rows = Array.isArray(raw) ? raw : Array.isArray(raw && raw.symbols) ? raw.symbols : [];
          const items = [];
          for (let i = 0; i < rows.length && items.length < 40; i++) {
            const m = mapTvSymbolRow(rows[i]);
            if (m) items.push(m);
          }
          const result = { items: items };
          if (items.length > 0) searchCache.set(key, result);
          return result;
        }

        function defaultSymbols(count) {
          const base = [
            "BINANCE:BTCUSDT",
            "BINANCE:ETHUSDT",
            "BINANCE:SOLUSDT",
            "BINANCE:BNBUSDT",
            "BINANCE:XRPUSDT",
            "BINANCE:ADAUSDT",
            "NASDAQ:AAPL",
            "NASDAQ:MSFT",
            "NASDAQ:NVDA",
            "FOREXCOM:EURUSD",
            "FOREXCOM:USDTRY",
            "TVC:GOLD"
          ];
          const out = [];
          for (let i = 0; i < count; i++) out.push(base[i % base.length]);
          return out;
        }

        function normalizeIntervalsForSymbols(state) {
          const n = state.symbols.length;
          let cur = Array.isArray(state.intervals) ? state.intervals.slice(0, n) : [];
          while (cur.length < n) cur.push(null);
          if (cur.length > n) cur = cur.slice(0, n);
          state.intervals = cur.map((v) => (v === null || v === "" ? null : String(v)));
          return state;
        }

        /** Chart count = symbols.length; columns only control grid layout. Migrates v1 rows×cols. */
        const CHART_THEMES = [
          { value: "dark", label: "Dark" },
          { value: "light", label: "Light" }
        ];

        const BAR_COLOR_PRESETS = [
          { value: "default", label: "Default (theme)" },
          { value: "tradingview", label: "Classic green / red" },
          { value: "ocean", label: "Ocean teal / coral" },
          { value: "cyber", label: "Cyber cyan / magenta" },
          { value: "mono", label: "Monochrome" },
          { value: "high_contrast", label: "High contrast" }
        ];

        function normalizeBarColorPreset(v) {
          const s = String(v || "default");
          if (BAR_COLOR_PRESETS.some((p) => p.value === s)) return s;
          return "default";
        }

        /**
         * TradingView chart overrides (candles + bars). Omitted when preset is "default".
         * @see https://www.tradingview.com/charting-library-docs/latest/customization/overrides/chart-overrides/
         */
        function barColorOverridesForPreset(preset) {
          if (preset === "default") return undefined;
          const palettes = {
            tradingview: {
              up: "#089981",
              down: "#F23645",
              borderUp: "#089981",
              borderDown: "#F23645",
              wickUp: "#737375",
              wickDown: "#737375"
            },
            ocean: { up: "#26A69A", down: "#EF5350" },
            cyber: { up: "#00E5FF", down: "#FF00AA" },
            mono: { up: "#B0BEC5", down: "#546E7A" },
            high_contrast: { up: "#00E676", down: "#FF1744" }
          };
          const c = palettes[preset];
          if (!c) return undefined;
          const up = c.up;
          const down = c.down;
          const bu = c.borderUp || up;
          const bd = c.borderDown || down;
          const wu = c.wickUp || up;
          const wd = c.wickDown || down;
          return {
            "mainSeriesProperties.candleStyle.upColor": up,
            "mainSeriesProperties.candleStyle.downColor": down,
            "mainSeriesProperties.candleStyle.borderUpColor": bu,
            "mainSeriesProperties.candleStyle.borderDownColor": bd,
            "mainSeriesProperties.candleStyle.wickUpColor": wu,
            "mainSeriesProperties.candleStyle.wickDownColor": wd,
            "mainSeriesProperties.barStyle.upColor": up,
            "mainSeriesProperties.barStyle.downColor": down
          };
        }

        function buildMultichartState(stored) {
          const cols = clampInt(stored?.cols, 1, 6, 3);
          const defaultInterval = (stored?.defaultInterval ? String(stored.defaultInterval) : (stored?.interval ? String(stored.interval) : "60")).trim() || "60";
          const chartHeightPx = clampInt(stored?.chartHeightPx, 180, 900, 260);
          const chartTheme = stored?.chartTheme === "light" ? "light" : "dark";
          const barColorPreset = normalizeBarColorPreset(stored?.barColorPreset);
          let symbols = Array.isArray(stored?.symbols) ? [...stored.symbols] : [];
          let intervals = Array.isArray(stored?.intervals) ? [...stored.intervals] : [];
          if (stored && typeof stored.rows === "number" && stored.rows > 0) {
            const cap = stored.rows * cols;
            while (symbols.length < cap) symbols.push(defaultSymbols(cap)[symbols.length]);
            if (symbols.length > cap) symbols = symbols.slice(0, cap);
          }
          if (symbols.length === 0 && stored == null) {
            symbols = [defaultSymbols(1)[0]];
            intervals = [null];
          }
          if (symbols.length > 0) {
            while (intervals.length < symbols.length) intervals.push(null);
            intervals = intervals.slice(0, symbols.length);
          }
          return normalizeIntervalsForSymbols({ cols, defaultInterval, chartHeightPx, chartTheme, barColorPreset, symbols, intervals });
        }

        function buildBarColorPresetOptions(selectEl, value) {
          selectEl.innerHTML = "";
          BAR_COLOR_PRESETS.forEach((p) => {
            const opt = document.createElement("option");
            opt.value = p.value;
            opt.textContent = p.label;
            selectEl.appendChild(opt);
          });
          selectEl.value = normalizeBarColorPreset(value);
        }

        function buildChartThemeOptions(selectEl, value) {
          selectEl.innerHTML = "";
          CHART_THEMES.forEach((t) => {
            const opt = document.createElement("option");
            opt.value = t.value;
            opt.textContent = t.label;
            selectEl.appendChild(opt);
          });
          selectEl.value = value === "light" ? "light" : "dark";
        }

        function syncSettingsFormFromAppState() {
          if (!appState) return;
          const colsEl = document.getElementById("cols");
          const defaultIntervalEl = document.getElementById("defaultInterval");
          const chartHeightEl = document.getElementById("chartHeight");
          const chartThemeEl = document.getElementById("chartTheme");
          const barColorPresetEl = document.getElementById("barColorPreset");
          if (colsEl) colsEl.value = String(appState.cols);
          if (defaultIntervalEl) buildIntervalOptions(defaultIntervalEl, String(appState.defaultInterval), false);
          if (chartHeightEl) chartHeightEl.value = String(appState.chartHeightPx);
          if (chartThemeEl) buildChartThemeOptions(chartThemeEl, appState.chartTheme);
          if (barColorPresetEl) buildBarColorPresetOptions(barColorPresetEl, appState.barColorPreset);
        }

        function buildIntervalOptions(selectEl, value, includeDefaultOption) {
          selectEl.innerHTML = "";
          if (includeDefaultOption) {
            const opt = document.createElement("option");
            opt.value = "";
            opt.textContent = "Default";
            selectEl.appendChild(opt);
          }
          INTERVAL_OPTIONS.forEach((v) => {
            const opt = document.createElement("option");
            opt.value = String(v);
            opt.textContent = intervalLabel(v);
            selectEl.appendChild(opt);
          });
          selectEl.value = value ?? "";
        }

        function effectiveInterval(state, index) {
          const v = state.intervals?.[index];
          return (v && String(v).trim()) ? String(v).trim() : state.defaultInterval;
        }

        function updateTopbarDefaultIntervalDisplay(src) {
          const el = document.getElementById("topbarDefaultInterval");
          if (!el) return;
          const st = src || appState;
          if (!st) return;
          const v = String(st.defaultInterval || "60").trim() || "60";
          el.textContent = "Default · " + intervalLabel(v);
        }

        function activePageFromHash() {
          const h = (location.hash || "#multichart").replace("#", "");
          return h || "multichart";
        }

        function setActiveNav(page) {
          const nav = document.getElementById("nav");
          nav.querySelectorAll(".nav-item").forEach((el) => {
            el.dataset.active = String(el.dataset.page === page);
          });
        }

        function swap(arr, i, j) {
          const tmp = arr[i];
          arr[i] = arr[j];
          arr[j] = tmp;
        }

        function loadTradingViewScript() {
          return new Promise((resolve, reject) => {
            if (window.TradingView && window.TradingView.widget) return resolve();
            const existing = document.querySelector('script[data-tv="true"]');
            if (existing) {
              existing.addEventListener("load", () => resolve());
              existing.addEventListener("error", () => reject(new Error("TradingView script failed")));
              return;
            }
            const s = document.createElement("script");
            s.src = "https://s3.tradingview.com/tv.js";
            s.async = true;
            s.dataset.tv = "true";
            s.onload = () => resolve();
            s.onerror = () => reject(new Error("TradingView script failed"));
            document.head.appendChild(s);
          });
        }

        let appState = null;
        let gridEl = null;
        let gridFooterEl = null;
        let tileRefs = [];
        let mounted = [];
        let tvPromise = null;

        function syncTileLayoutMeta() {
          tileRefs.forEach((r, i) => {
            r.tile.dataset.index = String(i);
            r.tile.style.order = String(i);
          });
        }

        function deleteChartAt(index) {
          if (!appState || index < 0 || index >= appState.symbols.length) return;
          appState.symbols.splice(index, 1);
          appState.intervals.splice(index, 1);
          saveState(appState);
          if (!gridEl || !tileRefs[index]) {
            initOrUpdateMultichart(true);
            return;
          }
          const tr = tileRefs[index];
          if (tr.sugg && tr.sugg.parentNode) tr.sugg.remove();
          gridEl.removeChild(tr.tile);
          tileRefs.splice(index, 1);
          mounted.splice(index, 1);
          syncTileLayoutMeta();
          for (let j = index; j < tileRefs.length; j++) {
            refreshTile(j, false);
          }
        }

        function addChart() {
          if (!appState) return;
          appState.symbols.push(defaultSymbols(1)[0]);
          appState.intervals.push(null);
          saveState(appState);
          const before = tileRefs.length;
          ensureGrid();
          for (let i = 0; i < before; i++) {
            refreshTile(i, false);
          }
          for (let i = before; i < tileRefs.length; i++) {
            refreshTile(i, true);
          }
        }

        function ensureTradingView() {
          if (window.TradingView && window.TradingView.widget) return Promise.resolve(true);
          if (!tvPromise) tvPromise = loadTradingViewScript().then(() => true).catch(() => false);
          return tvPromise;
        }

        function desiredConfigAt(index) {
          const symbol = (appState.symbols[index] || "").trim();
          const interval = effectiveInterval(appState, index);
          return { symbol, interval };
        }

        function maybeMountWidget(index, force) {
          const ref = tileRefs[index];
          if (!ref) return;

          const { symbol, interval } = desiredConfigAt(index);
          const theme = appState.chartTheme === "light" ? "light" : "dark";
          const barPreset = normalizeBarColorPreset(appState.barColorPreset);
          const prev = mounted[index];
          if (!force && prev && prev.symbol === symbol && prev.interval === interval && prev.theme === theme && prev.barPreset === barPreset) return;

          mounted[index] = { symbol, interval, theme, barPreset };

          // If TV isn't ready, show a minimal placeholder and bail
          if (!(window.TradingView && window.TradingView.widget)) {
            ref.body.innerHTML = "";
            const empty = document.createElement("div");
            empty.className = "empty";
            empty.textContent = "Loading TradingView…";
            ref.body.appendChild(empty);
            return;
          }

          ref.body.innerHTML = "";
          if (!symbol) {
            const empty = document.createElement("div");
            empty.className = "empty";
            empty.textContent = "Symbol is empty. Example: BINANCE:BTCUSDT";
            ref.body.appendChild(empty);
            return;
          }

          const id = "tv_" + index + "_" + Math.random().toString(16).slice(2);
          const host = document.createElement("div");
          host.id = id;
          host.style.height = "100%";
          ref.body.appendChild(host);

          const barOv = barColorOverridesForPreset(barPreset);
          const wopts = {
            autosize: true,
            symbol: symbol,
            interval: interval,
            timezone: "Etc/UTC",
            theme: theme,
            style: "1",
            locale: "en",
            enable_publishing: false,
            allow_symbol_change: false,
            hide_top_toolbar: true,
            hide_legend: true,
            container_id: id
          };
          if (barOv) wopts.overrides = barOv;

          // eslint-disable-next-line no-new
          new TradingView.widget(wopts);
        }

        /**
         * @param inputsOnly If true, only sync search input / interval UI (no TradingView remount).
         *   Use after drag-and-drop: DOM nodes are not moved (iframes would reload); only refs + CSS order change.
         */
        function refreshTile(index, forceWidget, inputsOnly) {
          const ref = tileRefs[index];
          if (!ref) return;

          ref.input.value = appState.symbols[index] || "";
          buildIntervalOptions(ref.intervalSel, appState.intervals[index] ?? "", true);
          if (inputsOnly) return;
          maybeMountWidget(index, !!forceWidget);
        }

        function clearDnDHighlights() {
          tileRefs.forEach((r) => r.tile.classList.remove("drop-target"));
        }

        function swapTiles(from, to) {
          if (from === to) return;
          swap(appState.symbols, from, to);
          swap(appState.intervals, from, to);
          saveState(appState);
          const tmpRef = tileRefs[from];
          tileRefs[from] = tileRefs[to];
          tileRefs[to] = tmpRef;
          swap(mounted, from, to);
          syncTileLayoutMeta();
          refreshTile(from, false, true);
          refreshTile(to, false, true);
        }

        function createTile(index) {
          const tile = document.createElement("div");
          tile.className = "tile";
          tile.dataset.index = String(index);

          const head = document.createElement("div");
          head.className = "tile-head";

          const dragHandle = document.createElement("div");
          dragHandle.className = "tile-drag-handle";
          dragHandle.draggable = true;
          dragHandle.title = "Drag to reorder";
          dragHandle.setAttribute("aria-label", "Drag to reorder");
          const dragGrip = document.createElement("span");
          dragGrip.className = "tile-drag-grip";
          dragGrip.setAttribute("aria-hidden", "true");
          dragHandle.appendChild(dragGrip);

          const left = document.createElement("div");
          left.className = "left";

          const searchWrap = document.createElement("div");
          searchWrap.className = "symbol-search";

          const input = document.createElement("input");
          input.className = "symbol-input";
          input.placeholder = "Search symbol…";
          input.autocomplete = "off";
          input.spellcheck = false;

          const sugg = document.createElement("div");
          sugg.className = "suggestions";
          sugg.style.display = "none";

          let unbindSuggPosition = null;
          function positionSuggDropdown() {
            const r = input.getBoundingClientRect();
            const w = Math.max(r.width, 220);
            sugg.style.left = Math.max(8, r.left) + "px";
            sugg.style.top = (r.bottom + 4) + "px";
            sugg.style.width = w + "px";
          }
          function bindSuggReposition() {
            if (unbindSuggPosition) unbindSuggPosition();
            const fn = () => positionSuggDropdown();
            window.addEventListener("scroll", fn, true);
            window.addEventListener("resize", fn);
            unbindSuggPosition = function () {
              window.removeEventListener("scroll", fn, true);
              window.removeEventListener("resize", fn);
              unbindSuggPosition = null;
            };
          }

          function closeSuggestions() {
            if (unbindSuggPosition) {
              unbindSuggPosition();
            }
            sugg.style.display = "none";
            sugg.innerHTML = "";
            if (sugg.parentNode === document.body) {
              searchWrap.appendChild(sugg);
            }
          }

          function openSuggestionsPanel() {
            if (sugg.parentNode !== document.body) {
              document.body.appendChild(sugg);
            }
            positionSuggDropdown();
            sugg.style.display = "block";
            bindSuggReposition();
          }

          function setSymbolAndRefresh(i, fullSymbol) {
            const next = String(fullSymbol || "").trim();
            const prev = (appState.symbols[i] || "").trim();
            if (!next || next === prev) return;
            appState.symbols[i] = next;
            saveState(appState);
            refreshTile(i, true);
          }

          input.addEventListener("change", () => {
            const i = Number.parseInt(tile.dataset.index, 10);
            // Allow manual paste/enter of full symbol (e.g. BINANCE:BTCUSDT)
            setSymbolAndRefresh(i, input.value.trim());
          });

          let symSearchAbort = null;
          const runSearch = debounce(async () => {
            const q = input.value.trim();
            if (q.length < 2) {
              if (symSearchAbort) symSearchAbort.abort();
              closeSuggestions();
              return;
            }
            if (symSearchAbort) symSearchAbort.abort();
            symSearchAbort = new AbortController();
            const signal = symSearchAbort.signal;
            let items = [];
            try {
              const r = await searchSymbols(q, { signal });
              if (signal.aborted) return;
              items = r.items || [];
            } catch (e) {
              if (e && e.name === "AbortError") return;
              items = [];
            }
            if (!items.length) {
              closeSuggestions();
              return;
            }

            sugg.innerHTML = "";
            items.forEach((it) => {
              const row = document.createElement("div");
              row.className = "suggestion";
              row.tabIndex = 0;

              const top = document.createElement("div");
              top.className = "top";
              const sym = document.createElement("div");
              sym.className = "sym";
              sym.textContent = it.symbol || it.full;
              const ex = document.createElement("div");
              ex.className = "ex";
              ex.textContent = it.exchange || "";
              top.appendChild(sym);
              top.appendChild(ex);

              const desc = document.createElement("div");
              desc.className = "desc";
              desc.textContent = it.full || "";

              row.appendChild(top);
              row.appendChild(desc);

              const choose = () => {
                const i = Number.parseInt(tile.dataset.index, 10);
                input.value = it.full;
                closeSuggestions();
                setSymbolAndRefresh(i, it.full);
              };
              row.addEventListener("mousedown", (e) => { e.preventDefault(); choose(); });
              row.addEventListener("keydown", (e) => { if (e.key === "Enter") choose(); });
              sugg.appendChild(row);
            });
            openSuggestionsPanel();
          }, 200);

          input.addEventListener("input", () => runSearch());
          input.addEventListener("focus", () => runSearch());
          input.addEventListener("keydown", (e) => { if (e.key === "Escape") closeSuggestions(); });
          input.addEventListener("blur", () => setTimeout(closeSuggestions, 350));

          searchWrap.appendChild(input);
          searchWrap.appendChild(sugg);
          left.appendChild(searchWrap);

          const right = document.createElement("div");
          right.style.display = "flex";
          right.style.alignItems = "center";
          right.style.gap = "8px";

          const intervalSel = document.createElement("select");
          intervalSel.className = "interval-select";
          intervalSel.title = "Interval (per chart)";
          intervalSel.addEventListener("change", () => {
            const i = Number.parseInt(tile.dataset.index, 10);
            const v = intervalSel.value;
            const next = v === "" ? null : v;
            const prev = appState.intervals[i] == null ? null : String(appState.intervals[i]);
            if (next === prev) return;
            appState.intervals[i] = next;
            saveState(appState);
            refreshTile(i, true);
          });

          const menuWrap = document.createElement("div");
          menuWrap.className = "tile-menu-wrap";
          const menuBtn = document.createElement("button");
          menuBtn.type = "button";
          menuBtn.className = "tile-menu-btn";
          menuBtn.textContent = "⋯";
          menuBtn.title = "Chart options";
          menuBtn.setAttribute("aria-label", "Chart options");
          menuBtn.setAttribute("aria-haspopup", "true");
          menuBtn.setAttribute("aria-expanded", "false");
          const menuPanel = document.createElement("div");
          menuPanel.className = "tile-menu-panel";
          menuPanel.setAttribute("role", "menu");
          const delItem = document.createElement("button");
          delItem.type = "button";
          delItem.className = "tile-menu-item";
          delItem.textContent = "Delete";
          delItem.setAttribute("role", "menuitem");
          delItem.addEventListener("click", (e) => {
            e.stopPropagation();
            menuPanel.dataset.open = "false";
            menuBtn.setAttribute("aria-expanded", "false");
            const i = Number.parseInt(tile.dataset.index, 10);
            deleteChartAt(i);
          });
          menuBtn.addEventListener("click", (e) => {
            e.stopPropagation();
            const open = menuPanel.dataset.open !== "true";
            document.querySelectorAll(".tile-menu-panel").forEach((p) => {
              p.dataset.open = "false";
            });
            document.querySelectorAll(".tile-menu-btn").forEach((b) => b.setAttribute("aria-expanded", "false"));
            if (open) {
              menuPanel.dataset.open = "true";
              menuBtn.setAttribute("aria-expanded", "true");
            }
          });
          menuPanel.appendChild(delItem);
          menuWrap.appendChild(menuBtn);
          menuWrap.appendChild(menuPanel);

          right.appendChild(intervalSel);
          right.appendChild(menuWrap);

          head.appendChild(dragHandle);
          head.appendChild(left);
          head.appendChild(right);

          const body = document.createElement("div");
          body.className = "tile-body";
          body.style.height = "calc(100% - 49px)";

          dragHandle.addEventListener("dragstart", (e) => {
            const i = Number.parseInt(tile.dataset.index, 10);
            e.dataTransfer.setData("text/plain", String(i));
            e.dataTransfer.effectAllowed = "move";
            tile.classList.add("dragging");
          });
          dragHandle.addEventListener("dragend", () => {
            tile.classList.remove("dragging");
            clearDnDHighlights();
          });
          head.addEventListener("dragover", (e) => {
            e.preventDefault();
            e.dataTransfer.dropEffect = "move";
            tile.classList.add("drop-target");
          });
          head.addEventListener("dragleave", () => {
            tile.classList.remove("drop-target");
          });
          head.addEventListener("drop", (e) => {
            e.preventDefault();
            const fromRaw = e.dataTransfer.getData("text/plain");
            const from = Number.parseInt(fromRaw, 10);
            const to = Number.parseInt(tile.dataset.index, 10);
            tile.classList.remove("drop-target");
            if (Number.isNaN(from) || Number.isNaN(to) || from === to) return;
            swapTiles(from, to);
          });

          tile.appendChild(head);
          tile.appendChild(body);

          return { tile, head, input, intervalSel, body, sugg };
        }

        function ensureGrid() {
          const root = document.getElementById("root");
          if (!gridEl) {
            root.innerHTML = "";
            gridEl = document.createElement("div");
            gridEl.className = "grid";
            root.appendChild(gridEl);
            gridFooterEl = document.createElement("div");
            gridFooterEl.className = "add-chart-footer";
            const addBtn = document.createElement("button");
            addBtn.type = "button";
            addBtn.className = "add-chart-btn";
            addBtn.textContent = "+";
            addBtn.title = "Add chart";
            addBtn.setAttribute("aria-label", "Add chart");
            addBtn.addEventListener("click", () => addChart());
            gridFooterEl.appendChild(addBtn);
            root.appendChild(gridFooterEl);
          }

          gridEl.style.gridTemplateColumns = "repeat(" + appState.cols + ", minmax(0, 1fr))";
          gridEl.style.setProperty("--chart-tile-min", String(appState.chartHeightPx) + "px");

          const want = appState.symbols.length;
          const have = tileRefs.length;

          if (want < have) {
            for (let i = have - 1; i >= want; i--) {
              const tr = tileRefs[i];
              if (tr.sugg && tr.sugg.parentNode) {
                tr.sugg.parentNode.removeChild(tr.sugg);
              }
              gridEl.removeChild(tr.tile);
              tileRefs.pop();
              mounted.pop();
            }
          } else if (want > have) {
            for (let i = have; i < want; i++) {
              const ref = createTile(i);
              tileRefs.push(ref);
              mounted.push(null);
              gridEl.appendChild(ref.tile);
            }
          }

          syncTileLayoutMeta();
        }

        function initOrUpdateMultichart(forceAllWidgets) {
          ensureGrid();
          tileRefs.forEach((_, i) => refreshTile(i, !!forceAllWidgets));
        }

        function render() {
          const page = activePageFromHash();
          setActiveNav(page);

          const title = document.getElementById("pageTitle");
          const subtitle = document.getElementById("pageSubtitle");
          const topbarActions = document.getElementById("topbarActions");
          const root = document.getElementById("root");

          if (page !== "multichart") {
            title.textContent = "finance-dashboard";
            subtitle.textContent = "Minimal dashboard shell";
            if (topbarActions) topbarActions.style.display = "none";
            root.innerHTML = '<div class="empty">Select Multichart from the sidebar.</div>';
            gridEl = null;
            gridFooterEl = null;
            tileRefs = [];
            mounted = [];
            return;
          }

          if (topbarActions) topbarActions.style.display = "flex";
          title.textContent = "Multichart";
          subtitle.textContent = "TradingView chart tile grid";

          const stored = loadState();
          const state = buildMultichartState(stored);

          const settingsDialog = document.getElementById("settingsDialog");
          const settingsOpen = document.getElementById("settingsOpen");
          const colsEl = document.getElementById("cols");
          const defaultIntervalEl = document.getElementById("defaultInterval");
          const chartHeightEl = document.getElementById("chartHeight");
          const chartThemeEl = document.getElementById("chartTheme");
          const barColorPresetEl = document.getElementById("barColorPreset");

          colsEl.value = String(state.cols);
          buildIntervalOptions(defaultIntervalEl, String(state.defaultInterval), false);
          chartHeightEl.value = String(state.chartHeightPx);
          buildChartThemeOptions(chartThemeEl, state.chartTheme);
          buildBarColorPresetOptions(barColorPresetEl, state.barColorPreset);
          updateTopbarDefaultIntervalDisplay(state);

          if (settingsOpen && settingsDialog) {
            settingsOpen.onclick = () => {
              syncSettingsFormFromAppState();
              settingsDialog.showModal();
              settingsOpen.setAttribute("aria-expanded", "true");
            };
            settingsDialog.onclose = () => {
              settingsOpen.setAttribute("aria-expanded", "false");
            };
          }
          const settingsClose = document.getElementById("settingsClose");
          const settingsCancel = document.getElementById("settingsCancel");
          if (settingsClose && settingsDialog) settingsClose.onclick = () => settingsDialog.close();
          if (settingsCancel && settingsDialog) settingsCancel.onclick = () => settingsDialog.close();

          const apply = document.getElementById("apply");
          apply.onclick = () => {
            const prevDefault = state.defaultInterval;
            state.cols = clampInt(colsEl.value, 1, 6, state.cols);
            state.defaultInterval = String(defaultIntervalEl.value || state.defaultInterval).trim() || "60";
            state.chartHeightPx = clampInt(chartHeightEl.value, 180, 900, state.chartHeightPx);
            state.chartTheme = chartThemeEl.value === "light" ? "light" : "dark";
            state.barColorPreset = normalizeBarColorPreset(barColorPresetEl.value);
            saveState(state);

            appState = state;
            updateTopbarDefaultIntervalDisplay();

            initOrUpdateMultichart(false);

            if (prevDefault !== state.defaultInterval) initOrUpdateMultichart(false);

            ensureTradingView().then((ok) => { if (ok) initOrUpdateMultichart(false); });
            if (settingsDialog) settingsDialog.close();
          };

          // Initial paint
          saveState(state);
          appState = state;
          ensureTradingView().then((ok) => {
            if (!ok) {
              const root = document.getElementById("root");
              root.innerHTML = '<div class="empty">TradingView script failed to load (network/policy).</div>';
              return;
            }
            initOrUpdateMultichart(true);
          });
        }

        document.addEventListener("click", () => {
          document.querySelectorAll(".tile-menu-panel").forEach((p) => { p.dataset.open = "false"; });
          document.querySelectorAll(".tile-menu-btn").forEach((b) => b.setAttribute("aria-expanded", "false"));
        });

        window.addEventListener("hashchange", render);

        const sidebarToggle = document.getElementById("sidebarToggle");
        if (sidebarToggle) {
          sidebarToggle.addEventListener("click", () => {
            const app = document.getElementById("appRoot");
            const next = !(app && app.classList.contains("sidebar-hidden"));
            saveSidebarCollapsed(next);
            applySidebarCollapsed(next);
          });
        }
        applySidebarCollapsed(loadSidebarCollapsed());

        render();
      })();
    </script>
  </body>
</html>`;

