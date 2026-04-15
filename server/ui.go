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
      .tabsbar-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px 16px;
        flex-wrap: wrap;
        margin-bottom: 12px;
        min-width: 0;
      }
      .tabsbar-row .tabsbar {
        display: flex;
        align-items: center;
        gap: 8px;
        flex-wrap: wrap;
        flex: 1 1 auto;
        min-width: 0;
        margin-bottom: 0;
      }
      #root {
        width: 100%;
        min-width: 0;
        box-sizing: border-box;
      }
      #tabWorkspacesRoot {
        width: 100%;
        min-width: 0;
        box-sizing: border-box;
      }
      .tab {
        height: 32px;
        border-radius: 999px;
        border: 1px solid rgba(231, 238, 248, 0.14);
        background: rgba(12, 18, 27, 0.55);
        color: var(--muted);
        padding: 0 12px;
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        gap: 8px;
        user-select: none;
        -webkit-user-select: none;
      }
      .tab[data-active="true"] {
        color: var(--text);
        background: rgba(102, 163, 255, 0.12);
        border-color: rgba(102, 163, 255, 0.28);
      }
      .tab .x {
        width: 18px;
        height: 18px;
        border-radius: 999px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        color: rgba(231, 238, 248, 0.5);
        border: 1px solid transparent;
      }
      .tab .x:hover {
        color: var(--text);
        background: rgba(231, 238, 248, 0.06);
        border-color: rgba(231, 238, 248, 0.14);
      }
      .tab .tab-gear {
        width: 22px;
        height: 22px;
        margin: 0 0 0 2px;
        padding: 0;
        border: none;
        border-radius: 6px;
        background: transparent;
        color: var(--muted);
        cursor: pointer;
        flex-shrink: 0;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        line-height: 1;
      }
      .tab .tab-gear:hover {
        color: var(--accent);
        background: rgba(102, 163, 255, 0.12);
      }
      .tab[data-active="true"] .tab-gear:hover {
        color: var(--text);
      }
      .tab-add {
        width: 34px;
        height: 32px;
        border-radius: 999px;
        border: 1px dashed rgba(102, 163, 255, 0.35);
        background: rgba(102, 163, 255, 0.06);
        color: var(--accent);
        cursor: pointer;
        font-size: 18px;
        line-height: 1;
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
      .control input#tabSettingsName {
        width: 100%;
        max-width: 320px;
        box-sizing: border-box;
      }
      .control button {
        cursor: pointer;
        background: rgba(102, 163, 255, 0.12);
        border-color: rgba(102, 163, 255, 0.28);
      }
      .control button:hover {
        background: rgba(102, 163, 255, 0.18);
      }

      .tab-workspace {
        display: flex;
        flex-direction: column;
        gap: 16px;
        min-width: 0;
        width: 100%;
        box-sizing: border-box;
      }
      .tab-workspace .grid {
        margin-top: 0;
        width: 100%;
        min-width: 0;
      }
      .tab-workspace[data-auto-summary="0"] .grid {
        margin-top: 14px;
      }

      .auto-sort-bar {
        display: none;
      }
      .auto-sort-bar.auto-sort-bar--inline {
        flex: 0 0 auto;
        display: none;
        flex-wrap: nowrap;
        align-items: center;
        gap: 8px;
        min-width: 0;
        max-width: 100%;
        box-sizing: border-box;
        margin: 0 0 0 auto;
        padding: 4px 2px 4px 14px;
        border: none;
        border-radius: 0;
        border-left: 1px solid rgba(231, 238, 248, 0.1);
        background: transparent;
      }
      .auto-sort-bar.auto-sort-bar--inline[data-visible="true"] {
        display: inline-flex;
      }
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-label {
        font-size: 10px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--muted);
        opacity: 0.85;
        margin: 0;
        white-space: nowrap;
      }
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-controls {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        padding: 2px 6px 2px 8px;
        border-radius: 999px;
        background: rgba(12, 18, 27, 0.5);
        border: 1px solid rgba(231, 238, 248, 0.1);
      }
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-key,
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-dir-select {
        height: 26px;
        border-radius: 999px;
        border: 1px solid rgba(231, 238, 248, 0.1);
        background: rgba(8, 12, 18, 0.65);
        color: var(--text);
        padding: 0 9px;
        outline: none;
        appearance: none;
        font-size: 11px;
        font-weight: 500;
        cursor: pointer;
      }
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-key:hover,
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-dir-select:hover {
        border-color: rgba(102, 163, 255, 0.25);
        background: rgba(12, 18, 27, 0.85);
      }
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-key:focus,
      .auto-sort-bar.auto-sort-bar--inline .auto-sort-dir-select:focus {
        border-color: rgba(102, 163, 255, 0.4);
      }
      .tile.is-auto .tile-drag-handle {
        display: none !important;
      }

      .auto-summary-bar {
        display: none;
        width: 100%;
        min-width: 0;
        box-sizing: border-box;
      }
      .auto-summary-bar[data-visible="true"] {
        display: grid;
        grid-template-columns: repeat(6, minmax(0, 1fr));
        gap: 12px;
        align-items: stretch;
        /* Clear separation from chart grid below (in addition to flex gap). */
        margin-bottom: 16px;
      }
      @media (max-width: 1100px) {
        .auto-summary-bar[data-visible="true"] {
          grid-template-columns: repeat(3, minmax(0, 1fr));
        }
      }
      @media (max-width: 900px) {
        .auto-summary-bar[data-visible="true"] {
          grid-template-columns: repeat(2, minmax(0, 1fr));
        }
      }
      @media (max-width: 420px) {
        .auto-summary-bar[data-visible="true"] {
          grid-template-columns: minmax(0, 1fr);
        }
      }
      .auto-summary-card {
        min-width: 0;
        padding: 12px 14px;
        border-radius: var(--radius);
        border: 1px solid rgba(231, 238, 248, 0.12);
        background: rgba(12, 18, 27, 0.55);
        display: flex;
        flex-direction: column;
        justify-content: center;
        gap: 6px;
        box-sizing: border-box;
      }
      .auto-summary-label {
        font-size: 10px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.06em;
        color: var(--muted);
        line-height: 1.25;
      }
      .auto-summary-value {
        font-size: 16px;
        font-weight: 700;
        font-variant-numeric: tabular-nums;
        color: var(--text);
        line-height: 1.25;
        word-break: break-word;
        min-width: 0;
      }
      .auto-summary-card--pnl.pnl-pos {
        border-color: rgba(38, 194, 129, 0.45);
        background: linear-gradient(145deg, rgba(38, 194, 129, 0.12), rgba(10, 14, 18, 0.88));
      }
      .auto-summary-card--pnl.pnl-neg {
        border-color: rgba(242, 54, 69, 0.42);
        background: linear-gradient(145deg, rgba(242, 54, 69, 0.1), rgba(14, 10, 12, 0.9));
      }
      .auto-summary-card--pnl .auto-summary-value.pnl-pos {
        color: #3dd68c;
      }
      .auto-summary-card--pnl .auto-summary-value.pnl-neg {
        color: #ff6b7a;
      }

      /* Logical order uses CSS order property so DnD can swap without moving DOM (moving iframes reloads them). */
      .grid {
        display: grid;
        gap: 12px;
        --chart-tile-min: 260px;
      }
      .tile {
        display: flex;
        flex-direction: column;
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
      .tile-position-strip {
        display: none;
        align-items: center;
        justify-content: space-between;
        gap: 10px 12px;
        padding: 7px 10px;
        font-size: 11px;
        line-height: 1.35;
        border-bottom: 1px solid rgba(231, 238, 248, 0.06);
        background: rgba(8, 12, 18, 0.55);
        color: var(--muted);
        width: 100%;
        min-width: 0;
        box-sizing: border-box;
      }
      .tile-position-strip[data-visible="true"] {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
      }
      .pos-metrics-row {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 0 2px;
        row-gap: 4px;
        flex: 1;
        min-width: 0;
        opacity: 0.94;
      }
      .pos-stat {
        display: inline-flex;
        align-items: baseline;
        gap: 5px;
        padding: 0 10px;
        border-left: 1px solid rgba(231, 238, 248, 0.1);
        max-width: 100%;
        box-sizing: border-box;
      }
      .pos-stat:first-child {
        padding-left: 0;
        border-left: none;
      }
      .pos-stat-label {
        flex-shrink: 0;
        font-size: 8px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.06em;
        color: var(--muted);
        opacity: 0.82;
      }
      .pos-stat-value {
        font-size: 11px;
        font-weight: 600;
        font-variant-numeric: tabular-nums;
        color: var(--text);
        letter-spacing: -0.02em;
        min-width: 0;
        word-break: break-word;
        opacity: 0.95;
      }
      .pos-stat-value.pos-stat-value--contract {
        font-size: 10px;
        font-weight: 600;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        letter-spacing: 0.01em;
      }
      .pos-stat-value.pos-stat-value--emph {
        font-weight: 700;
        color: #9ec8ff;
      }
      .pos-stat--long .pos-stat-value {
        color: #4ce49a;
      }
      .pos-stat--short .pos-stat-value {
        color: #ff8f9a;
      }
      .pos-pnl-block {
        flex-shrink: 0;
        align-self: stretch;
        margin-left: auto;
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        justify-content: center;
        gap: 3px;
        padding: 6px 10px 6px 12px;
        border-left: 1px solid rgba(231, 238, 248, 0.14);
        border-right: none;
        min-width: 7.5rem;
        max-width: 100%;
        box-sizing: border-box;
        border-radius: 6px;
        background: rgba(12, 18, 28, 0.75);
        border-top: 1px solid rgba(231, 238, 248, 0.08);
        border-bottom: 1px solid rgba(231, 238, 248, 0.06);
        box-shadow: 0 1px 8px rgba(0, 0, 0, 0.22);
        z-index: 1;
        text-align: right;
      }
      .pos-pnl-block.pnl-pos {
        border-left-color: rgba(38, 194, 129, 0.45);
        background: linear-gradient(135deg, rgba(38, 194, 129, 0.18), rgba(10, 16, 14, 0.88));
        box-shadow: 0 0 14px rgba(38, 194, 129, 0.08), 0 1px 8px rgba(0, 0, 0, 0.2);
      }
      .pos-pnl-block.pnl-neg {
        border-left-color: rgba(242, 54, 69, 0.38);
        background: linear-gradient(135deg, rgba(242, 54, 69, 0.14), rgba(16, 10, 12, 0.9));
        box-shadow: 0 0 14px rgba(242, 54, 69, 0.06), 0 1px 8px rgba(0, 0, 0, 0.2);
      }
      .pos-pnl-block.pos-pnl-block-empty {
        opacity: 0.88;
        background: rgba(12, 18, 28, 0.55);
        border-left-color: rgba(231, 238, 248, 0.1);
        box-shadow: none;
      }
      .pos-pnl-label {
        font-size: 9px;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--muted);
        opacity: 0.95;
      }
      .pos-pnl-value {
        font-size: 15px;
        font-weight: 800;
        font-variant-numeric: tabular-nums;
        line-height: 1.15;
        letter-spacing: -0.03em;
        color: var(--text);
      }
      .pos-pnl-value.pnl-pos {
        color: #3dd68c;
      }
      .pos-pnl-value.pnl-neg {
        color: #ff6b7a;
      }
      .pos-pnl-value.pos-pnl-value-muted {
        font-size: 14px;
        font-weight: 600;
        color: var(--muted);
      }
      .pos-pnl-value .pos-pnl-pct {
        font-size: 12px;
        font-weight: 700;
        opacity: 0.92;
      }
      .auto-summary-value .auto-summary-pct {
        font-size: 13px;
        font-weight: 600;
        opacity: 0.88;
      }
      .tile-position-strip .pos-chip {
        display: inline-flex;
        align-items: center;
        padding: 4px 0;
        color: var(--muted);
        font-size: 11px;
        white-space: nowrap;
      }
      .tile-body {
        flex: 1;
        min-height: 0;
        height: auto;
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

      #tabFlowOverlay {
        margin-top: 12px;
      }

      .type-picker {
        margin-top: 12px;
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 22px 18px;
        background: rgba(15, 22, 33, 0.45);
      }
      .type-picker h3 {
        margin: 0 0 6px;
        font-size: 15px;
        font-weight: 600;
        letter-spacing: 0.2px;
      }
      .type-picker p {
        margin: 0 0 16px;
        font-size: 12px;
        color: var(--muted);
        line-height: 1.45;
      }
      .type-picker-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
        gap: 12px;
      }
      .type-picker-card {
        border-radius: 12px;
        border: 1px solid rgba(102, 163, 255, 0.28);
        background: rgba(102, 163, 255, 0.08);
        padding: 16px 14px;
        cursor: pointer;
        text-align: left;
        color: var(--text);
        font: inherit;
        transition: background 0.15s ease, border-color 0.15s ease;
      }
      .type-picker-card:hover {
        background: rgba(102, 163, 255, 0.14);
        border-color: rgba(102, 163, 255, 0.45);
      }
      .type-picker-card strong {
        display: block;
        font-size: 14px;
        margin-bottom: 6px;
      }
      .type-picker-card span {
        font-size: 12px;
        color: var(--muted);
        line-height: 1.4;
      }
      .type-picker-back {
        margin-bottom: 12px;
        height: 32px;
        padding: 0 12px;
        border-radius: 10px;
        border: 1px solid var(--border);
        background: rgba(12, 18, 27, 0.7);
        color: var(--muted);
        cursor: pointer;
        font: inherit;
      }
      .type-picker-back:hover {
        color: var(--text);
        border-color: rgba(231, 238, 248, 0.22);
      }
      .exchange-picker-card {
        text-align: left;
      }
      .exchange-picker-card .ex-meta {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 8px;
        margin-top: 8px;
      }
      .exchange-picker-card .ex-badge {
        font-size: 10px;
        padding: 3px 8px;
        border-radius: 999px;
        border: 1px solid rgba(231, 238, 248, 0.14);
        color: var(--muted);
        white-space: nowrap;
      }
      .exchange-picker-card .ex-badge.ok {
        border-color: rgba(115, 255, 227, 0.35);
        color: #9ee8dc;
      }
      .exchange-picker-card.missing-keys {
        opacity: 0.65;
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

        <div class="tabsbar-row">
          <div class="tabsbar" id="tabsbar" aria-label="Tabs"></div>
          <div class="auto-sort-bar auto-sort-bar--inline" id="tabsbarSort" data-visible="false" aria-label="Chart order"></div>
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
              <div class="control">
                <label>Exchange API keys (for Auto · Positions)</label>
                <div style="display:grid; gap:10px;">
                  <div style="display:grid; gap:6px;">
                    <div class="pill">Binance</div>
                    <input id="binanceApiKey" placeholder="API key" autocomplete="off" spellcheck="false" />
                    <input id="binanceApiSecret" placeholder="API secret" autocomplete="off" spellcheck="false" />
                  </div>
                  <div style="display:grid; gap:6px;">
                    <div class="pill">Bybit</div>
                    <input id="bybitApiKey" placeholder="API key" autocomplete="off" spellcheck="false" />
                    <input id="bybitApiSecret" placeholder="API secret" autocomplete="off" spellcheck="false" />
                  </div>
                  <div style="display:grid; gap:6px;">
                    <div class="pill">OKX</div>
                    <input id="okxApiKey" placeholder="API key" autocomplete="off" spellcheck="false" />
                    <input id="okxApiSecret" placeholder="API secret" autocomplete="off" spellcheck="false" />
                  </div>
                  <div style="display:grid; gap:6px;">
                    <div class="pill">Pionex</div>
                    <input id="pionexApiKey" placeholder="API key" autocomplete="off" spellcheck="false" />
                    <input id="pionexApiSecret" placeholder="API secret" autocomplete="off" spellcheck="false" />
                  </div>
                </div>
              </div>
            </div>
            <div class="settings-footer">
              <button type="button" class="btn-secondary" id="settingsCancel">Cancel</button>
              <button type="button" class="btn-primary" id="apply">Apply</button>
            </div>
          </div>
        </dialog>

        <dialog id="tabSettingsDialog" class="settings-dialog" aria-labelledby="tabSettingsTitle">
          <div class="settings-dialog-inner">
            <div class="settings-header">
              <h3 class="settings-title" id="tabSettingsTitle">Tab settings</h3>
              <button type="button" class="settings-close" id="tabSettingsClose" aria-label="Close">×</button>
            </div>
            <div class="settings-body">
              <div class="control">
                <label for="tabSettingsName">Tab name</label>
                <input id="tabSettingsName" type="text" autocomplete="off" spellcheck="false" placeholder="Tab name" />
              </div>
              <div class="control">
                <label for="tabSettingsMode">Chart type</label>
                <select id="tabSettingsMode">
                  <option value="unset">Not set</option>
                  <option value="manual">Manual</option>
                  <option value="auto_positions">Auto · Positions</option>
                </select>
              </div>
              <div class="control" id="tabSettingsExchangeWrap" style="display:none">
                <label for="tabSettingsExchange">Exchange</label>
                <select id="tabSettingsExchange">
                  <option value="">Choose exchange…</option>
                  <option value="binance">Binance</option>
                  <option value="bybit">Bybit</option>
                  <option value="okx">OKX</option>
                  <option value="pionex">Pionex</option>
                </select>
              </div>
            </div>
            <div class="settings-footer">
              <button type="button" class="btn-secondary" id="tabSettingsCancel">Cancel</button>
              <button type="button" class="btn-primary" id="tabSettingsApply">Apply</button>
            </div>
          </div>
        </dialog>

        <dialog id="exchangeKeysDialog" class="settings-dialog" aria-labelledby="exchangeKeysModalTitle">
          <div class="settings-dialog-inner">
            <div class="settings-header">
              <h3 class="settings-title" id="exchangeKeysModalTitle">API keys</h3>
              <button type="button" class="settings-close" id="exchangeKeysModalClose" aria-label="Close">×</button>
            </div>
            <div class="settings-body">
              <p style="margin:0 0 12px;font-size:12px;color:var(--muted);line-height:1.45;">Keys are stored in this browser session only.</p>
              <div class="control">
                <label for="exchangeKeysModalApiKey">API key</label>
                <input id="exchangeKeysModalApiKey" type="password" autocomplete="off" spellcheck="false" style="width:100%;max-width:none;" />
              </div>
              <div class="control">
                <label for="exchangeKeysModalApiSecret">API secret</label>
                <input id="exchangeKeysModalApiSecret" type="password" autocomplete="off" spellcheck="false" style="width:100%;max-width:none;" />
              </div>
            </div>
            <div class="settings-footer">
              <button type="button" class="btn-secondary" id="exchangeKeysModalCancel">Cancel</button>
              <button type="button" class="btn-primary" id="exchangeKeysModalSave">Save &amp; continue</button>
            </div>
          </div>
        </dialog>

        <div id="root"></div>
      </main>
    </div>

    <script>
      (function () {
        const STORAGE_KEY = "finance-dashboard:multichart:v3";
        const SIDEBAR_STORAGE_KEY = "finance-dashboard:sidebar:collapsed";
        const EXCHANGE_KEYS_STORAGE_KEY = "finance-dashboard:exchange-keys:v1";

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
            const legacyV2 = sessionStorage.getItem("finance-dashboard:multichart:v2");
            if (legacyV2) return JSON.parse(legacyV2);
            const legacyV1 = sessionStorage.getItem("finance-dashboard:multichart:v1");
            if (legacyV1) return JSON.parse(legacyV1);
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

        function loadExchangeKeys() {
          try {
            const raw = sessionStorage.getItem(EXCHANGE_KEYS_STORAGE_KEY);
            if (!raw) return null;
            const v = JSON.parse(raw);
            if (!v || typeof v !== "object") return null;
            return v;
          } catch (_) {
            return null;
          }
        }

        function saveExchangeKeys(keys) {
          try {
            sessionStorage.setItem(EXCHANGE_KEYS_STORAGE_KEY, JSON.stringify(keys || {}));
          } catch (_) {}
        }

        /** Tab "mode" is unset until the user picks a type in the workspace. */
        function normalizeTabMode(v) {
          if (v === undefined || v === null) return "manual";
          const s = String(v).trim();
          if (s === "unset") return "unset";
          if (s === "auto_positions") return "auto_positions";
          return "manual";
        }

        function tabNeedsType(tab) {
          return normalizeTabMode(tab && tab.mode) === "unset";
        }

        function isAutoMode(tabOrState) {
          const st = tabOrState || activeTab();
          if (!st) return false;
          const m = st.mode != null ? st.mode : (st.tab && st.tab.mode);
          return normalizeTabMode(m) === "auto_positions";
        }

        /** Auto tab: exchange not chosen yet (after picking Auto, before picking API). */
        function autoNeedsExchange(tab) {
          const t = tab || activeTab();
          if (!t || !isAutoMode(t)) return false;
          return String(t.autoExchange || "").trim() === "";
        }

        const AUTO_EXCHANGES = [
          { id: "binance", label: "Binance" },
          { id: "bybit", label: "Bybit" },
          { id: "okx", label: "OKX" },
          { id: "pionex", label: "Pionex" }
        ];

        function uid() {
          return Math.random().toString(16).slice(2) + Math.random().toString(16).slice(2);
        }

        function normalizeAutoSortKey(v) {
          const s = String(v || "pnl").toLowerCase();
          if (s === "symbol" || s === "pnl" || s === "size" || s === "entry") return s;
          return "pnl";
        }

        function normalizeAutoSortDir(v) {
          return String(v || "").toLowerCase() === "asc" ? "asc" : "desc";
        }

        /** Same symbols (multiset), same count — used to avoid re-sort + chart remount on every poll. */
        function sameSymbolMultiset(a, b) {
          if (!a || !b || a.length !== b.length) return false;
          const sa = a.map((x) => String(x || "").trim()).filter(Boolean).sort();
          const sb = b.map((x) => String(x || "").trim()).filter(Boolean).sort();
          if (sa.length !== sb.length) return false;
          for (let i = 0; i < sa.length; i++) {
            if (sa[i] !== sb[i]) return false;
          }
          return true;
        }

        /**
         * Reorder auto tab data by saved sort key / direction (uses position meta).
         * @returns {number[]|null} permutation: new row k shows data from old row perm[k]; null if no reorder.
         */
        function sortAutoTabCharts(tab) {
          if (!tab || !isAutoMode(tab)) return null;
          const n = tab.symbols.length;
          if (n <= 1) return null;
          const key = normalizeAutoSortKey(tab.autoSortKey);
          const dir = normalizeAutoSortDir(tab.autoSortDir) === "asc" ? 1 : -1;
          const meta = Array.isArray(tab.autoPositionMeta) ? tab.autoPositionMeta : [];
          const sym = tab.symbols;
          const ints = tab.intervals;
          const indices = [...Array(n).keys()];
          function cmp(a, b) {
            const ma = meta[a];
            const mb = meta[b];
            let va;
            let vb;
            let diff = 0;
            if (key === "symbol") {
              va = String(sym[a] || "");
              vb = String(sym[b] || "");
              diff = va.localeCompare(vb);
            } else if (key === "pnl") {
              va = Number.parseFloat(String(ma && ma.unrealizedPnL || "").replace(/,/g, ""));
              vb = Number.parseFloat(String(mb && mb.unrealizedPnL || "").replace(/,/g, ""));
              if (!Number.isFinite(va)) va = 0;
              if (!Number.isFinite(vb)) vb = 0;
              diff = va - vb;
            } else if (key === "size") {
              va = Math.abs(Number.parseFloat(String(ma && ma.size || "").replace(/,/g, "")));
              vb = Math.abs(Number.parseFloat(String(mb && mb.size || "").replace(/,/g, "")));
              if (!Number.isFinite(va)) va = 0;
              if (!Number.isFinite(vb)) vb = 0;
              diff = va - vb;
            } else if (key === "entry") {
              va = Number.parseFloat(String(ma && ma.avgPrice || "").replace(/,/g, ""));
              vb = Number.parseFloat(String(mb && mb.avgPrice || "").replace(/,/g, ""));
              if (!Number.isFinite(va)) va = 0;
              if (!Number.isFinite(vb)) vb = 0;
              diff = va - vb;
            }
            if (diff !== 0) return diff * dir;
            return (a - b) * dir;
          }
          indices.sort(cmp);
          tab.symbols = indices.map((i) => sym[i]);
          tab.intervals = indices.map((i) => (ints[i] === undefined ? null : ints[i]));
          tab.autoPositionMeta = indices.map((i) => meta[i] || null);
          return indices;
        }

        /**
         * Keep TradingView iframes on the same DOM nodes; swap tile refs so row k shows the chart that was at old row perm[k].
         */
        function applyAutoSortPermutationToTiles(perm) {
          if (!perm || perm.length <= 1) return;
          activateWorkspace(appState.activeTabId);
          const n = perm.length;
          if (tileRefs.length !== n || mounted.length !== n) return;
          const oldT = tileRefs.slice();
          const oldM = mounted.slice();
          for (let k = 0; k < n; k++) {
            tileRefs[k] = oldT[perm[k]];
            mounted[k] = oldM[perm[k]];
          }
          syncTileLayoutMeta();
          for (let i = 0; i < n; i++) {
            refreshTile(i, false, false, true);
          }
        }

        function normalizeTab(t) {
          const tab = t && typeof t === "object" ? t : {};
          const id = String(tab.id || "").trim() || uid();
          const name = String(tab.name || "Tab").trim() || "Tab";
          const mode = normalizeTabMode(tab.mode);
          if (mode === "unset") {
            return normalizeIntervalsForSymbols({ id, name, mode: "unset", symbols: [], intervals: [] });
          }
          if (mode === "auto_positions") {
            const autoSortKey = normalizeAutoSortKey(tab.autoSortKey);
            const autoSortDir = normalizeAutoSortDir(tab.autoSortDir);
            const autoEx = String(tab.autoExchange || "").trim().toLowerCase();
            let symbols = Array.isArray(tab.symbols) ? [...tab.symbols] : [];
            let intervals = Array.isArray(tab.intervals) ? [...tab.intervals] : [];
            let autoPositionMeta = Array.isArray(tab.autoPositionMeta) ? [...tab.autoPositionMeta] : [];
            if (!autoEx) {
              if (symbols.length === 0) {
                return normalizeIntervalsForSymbols({
                  id,
                  name,
                  mode: "auto_positions",
                  autoExchange: "",
                  symbols: [],
                  intervals: [],
                  autoPositionMeta: [],
                  autoSortKey,
                  autoSortDir
                });
              }
              while (intervals.length < symbols.length) intervals.push(null);
              if (intervals.length > symbols.length) intervals = intervals.slice(0, symbols.length);
              while (autoPositionMeta.length < symbols.length) autoPositionMeta.push(null);
              if (autoPositionMeta.length > symbols.length) autoPositionMeta = autoPositionMeta.slice(0, symbols.length);
              return normalizeIntervalsForSymbols({
                id,
                name,
                mode: "auto_positions",
                autoExchange: "",
                symbols,
                intervals,
                autoPositionMeta,
                autoSortKey,
                autoSortDir
              });
            }
            if (symbols.length === 0) {
              symbols = [""];
              intervals = [null];
            }
            while (intervals.length < symbols.length) intervals.push(null);
            if (intervals.length > symbols.length) intervals = intervals.slice(0, symbols.length);
            while (autoPositionMeta.length < symbols.length) autoPositionMeta.push(null);
            if (autoPositionMeta.length > symbols.length) autoPositionMeta = autoPositionMeta.slice(0, symbols.length);
            return normalizeIntervalsForSymbols({
              id,
              name,
              mode: "auto_positions",
              autoExchange: autoEx,
              symbols,
              intervals,
              autoPositionMeta,
              autoSortKey,
              autoSortDir
            });
          }
          let symbols = Array.isArray(tab.symbols) ? [...tab.symbols] : [];
          let intervals = Array.isArray(tab.intervals) ? [...tab.intervals] : [];
          if (symbols.length === 0) {
            symbols = [defaultSymbols(1)[0]];
            intervals = [null];
          }
          while (intervals.length < symbols.length) intervals.push(null);
          if (intervals.length > symbols.length) intervals = intervals.slice(0, symbols.length);
          return normalizeIntervalsForSymbols({ id, name, mode: "manual", symbols, intervals });
        }

        function buildTabbedState(stored) {
          // Global layout/theme settings
          const cols = clampInt(stored?.cols, 1, 6, 3);
          const defaultInterval = (stored?.defaultInterval ? String(stored.defaultInterval) : (stored?.interval ? String(stored.interval) : "60")).trim() || "60";
          const chartHeightPx = clampInt(stored?.chartHeightPx, 180, 900, 260);
          const chartTheme = stored?.chartTheme === "light" ? "light" : "dark";
          const barColorPreset = normalizeBarColorPreset(stored?.barColorPreset);

          // Tabs + migration from older single-tab shapes (v1/v2)
          let tabs = [];
          let activeTabId = String(stored?.activeTabId || "").trim();

          if (Array.isArray(stored?.tabs) && stored.tabs.length) {
            tabs = stored.tabs.map(normalizeTab);
            if (!activeTabId) activeTabId = String(tabs[0].id);
          } else {
            // v1/v2 were single workspace: { mode, symbols, intervals, rows? ... }
            const legacyMode = stored && stored.mode != null ? normalizeTabMode(stored.mode) : "manual";
            let symbols = Array.isArray(stored?.symbols) ? [...stored.symbols] : [];
            let intervals = Array.isArray(stored?.intervals) ? [...stored.intervals] : [];
            if (stored && typeof stored.rows === "number" && stored.rows > 0) {
              const cap = stored.rows * cols;
              while (symbols.length < cap) symbols.push(defaultSymbols(cap)[symbols.length]);
              if (symbols.length > cap) symbols = symbols.slice(0, cap);
            }
            if (symbols.length === 0) {
              symbols = [defaultSymbols(1)[0]];
              intervals = [null];
            }
            while (intervals.length < symbols.length) intervals.push(null);
            if (intervals.length > symbols.length) intervals = intervals.slice(0, symbols.length);
            const single = normalizeTab({ id: uid(), name: "Main", mode: legacyMode, symbols, intervals });
            tabs = [single];
            activeTabId = single.id;
          }

          // Ensure active tab exists
          if (!tabs.some((t) => String(t.id) === activeTabId)) activeTabId = String(tabs[0].id);

          return { cols, defaultInterval, chartHeightPx, chartTheme, barColorPreset, activeTabId, tabs };
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
          if (Array.isArray(state.autoPositionMeta)) {
            let apm = state.autoPositionMeta.slice(0, n);
            while (apm.length < n) apm.push(null);
            state.autoPositionMeta = apm;
          }
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

        function activeTab(state) {
          const st = state || appState;
          if (!st || !Array.isArray(st.tabs) || st.tabs.length === 0) return null;
          const id = String(st.activeTabId || "");
          return st.tabs.find((t) => String(t.id) === id) || st.tabs[0];
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
          const def = String(state && state.defaultInterval != null ? state.defaultInterval : "60").trim() || "60";
          return (v && String(v).trim()) ? String(v).trim() : def;
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
        /** After sort: reorder tile DOM refs to match tab order without remounting TradingView (see applyAutoSortPermutationToTiles). */
        let pendingAutoSortPermute = null;
        let tvPromise = null;
        let exchangeKeysModalPendingId = null;
        let tabSettingsTargetId = null;
        /** Per-tab chart DOM so switching tabs does not destroy TradingView iframes. */
        let tabWorkspaces = Object.create(null);
        /** Interval for Auto · Positions: refresh symbols + PnL strip from the exchange API. */
        let autoPositionsPollTimer = null;
        const AUTO_POSITIONS_POLL_MS = 5000;

        /** Pionex public futures WebSocket: live mark price (INDEX); positions still from REST poll. */
        const PIONEX_WS_PUB = "wss://ws.pionex.com/wsPub";
        let pionexIndexWs = null;
        let pionexIndexSubscribed = new Set();
        let pionexIndexReconnectTimer = null;

        function rawSymbolToPionexWsSymbol(raw) {
          let s = String(raw || "").trim();
          if (!s) return "";
          if (s.includes(":")) s = s.split(":").pop().trim();
          s = s.toUpperCase().replace(/-/g, "_");
          if (/^[A-Z0-9]+_USDT_PERP$/.test(s)) return s;
          const compact = s.replace(/_/g, "");
          const m = /^([A-Z0-9]+)USDT(PERP)?$/i.exec(compact);
          if (m) return m[1] + "_USDT_PERP";
          return "";
        }

        function shouldPionexMarkStreamActive() {
          if (!appState || document.hidden) return false;
          if (activePageFromHash() !== "multichart") return false;
          const tab = activeTab();
          if (!tab || !isAutoMode(tab) || autoNeedsExchange(tab)) return false;
          if (String(tab.autoExchange || "").trim().toLowerCase() !== "pionex") return false;
          return true;
        }

        function getDesiredPionexWsSymbols() {
          const out = new Set();
          const tab = activeTab();
          if (!tab) return out;
          const meta = Array.isArray(tab.autoPositionMeta) ? tab.autoPositionMeta : [];
          meta.forEach((m) => {
            if (!m || typeof m !== "object") return;
            const ws = rawSymbolToPionexWsSymbol(m.rawSymbol);
            if (ws) out.add(ws);
          });
          return out;
        }

        function stopPionexMarkStream() {
          if (pionexIndexReconnectTimer != null) {
            clearTimeout(pionexIndexReconnectTimer);
            pionexIndexReconnectTimer = null;
          }
          const old = pionexIndexWs;
          pionexIndexWs = null;
          pionexIndexSubscribed.clear();
          if (old) {
            try {
              old.onclose = null;
              old.close();
            } catch (_) {}
          }
        }

        function applyPionexIndexSubscriptions(ws, desired) {
          const toUnsub = [...pionexIndexSubscribed].filter((s) => !desired.has(s));
          const toSub = [...desired].filter((s) => !pionexIndexSubscribed.has(s));
          toUnsub.forEach((s) => {
            try {
              ws.send(JSON.stringify({ op: "UNSUBSCRIBE", topic: "INDEX", symbol: s }));
            } catch (_) {}
            pionexIndexSubscribed.delete(s);
          });
          toSub.forEach((s) => {
            try {
              ws.send(JSON.stringify({ op: "SUBSCRIBE", topic: "INDEX", symbol: s }));
            } catch (_) {}
            pionexIndexSubscribed.add(s);
          });
        }

        function applyPionexMarkPrice(wsSym, markStr) {
          const tab = activeTab();
          if (!tab || !isAutoMode(tab) || String(tab.autoExchange || "").trim().toLowerCase() !== "pionex") return;
          const meta = tab.autoPositionMeta;
          if (!Array.isArray(meta)) return;
          let changed = false;
          meta.forEach((m) => {
            if (!m || typeof m !== "object") return;
            const w = rawSymbolToPionexWsSymbol(m.rawSymbol);
            if (w !== wsSym) return;
            const next = String(markStr || "").trim();
            if (m.markPrice === next) return;
            m.markPrice = next;
            changed = true;
          });
          if (!changed) return;
          activateWorkspace(appState.activeTabId);
          tileRefs.forEach((ref, i) => syncPositionStrip(ref, i));
          syncAutoSummaryBar();
        }

        function handlePionexIndexWsMessage(ev) {
          let msg;
          try {
            msg = JSON.parse(ev.data);
          } catch (_) {
            return;
          }
          if (!msg || typeof msg !== "object") return;
          if (msg.op === "PING") {
            const ts = msg.timestamp;
            const w = pionexIndexWs;
            if (w && w.readyState === WebSocket.OPEN) {
              try {
                w.send(JSON.stringify({ op: "PONG", timestamp: ts }));
              } catch (_) {}
            }
            return;
          }
          const topic = String(msg.topic || msg.type || "").toUpperCase();
          if (topic !== "INDEX") return;
          const arr = msg.data;
          if (!Array.isArray(arr) || !arr.length) return;
          const row = arr[0];
          if (!row || typeof row !== "object") return;
          const sym = String(msg.symbol || row.symbol || "").trim();
          const mark = row.markPrice;
          if (!sym || mark == null || String(mark).trim() === "") return;
          applyPionexMarkPrice(sym, String(mark));
        }

        function syncPionexMarkStream() {
          if (!shouldPionexMarkStreamActive()) {
            stopPionexMarkStream();
            return;
          }
          const desired = getDesiredPionexWsSymbols();
          if (desired.size === 0) {
            stopPionexMarkStream();
            return;
          }
          if (pionexIndexWs && pionexIndexWs.readyState === WebSocket.OPEN) {
            applyPionexIndexSubscriptions(pionexIndexWs, desired);
            return;
          }
          if (pionexIndexWs && pionexIndexWs.readyState === WebSocket.CONNECTING) {
            return;
          }
          if (pionexIndexWs) {
            stopPionexMarkStream();
          }
          const ws = new WebSocket(PIONEX_WS_PUB);
          pionexIndexWs = ws;
          ws.onopen = () => {
            if (pionexIndexWs !== ws) return;
            pionexIndexSubscribed.clear();
            const d = getDesiredPionexWsSymbols();
            d.forEach((sym) => {
              try {
                ws.send(JSON.stringify({ op: "SUBSCRIBE", topic: "INDEX", symbol: sym }));
                pionexIndexSubscribed.add(sym);
              } catch (_) {}
            });
          };
          ws.onmessage = handlePionexIndexWsMessage;
          ws.onerror = () => {};
          ws.onclose = () => {
            if (pionexIndexWs !== ws) return;
            pionexIndexWs = null;
            pionexIndexSubscribed.clear();
            if (pionexIndexReconnectTimer != null) {
              clearTimeout(pionexIndexReconnectTimer);
              pionexIndexReconnectTimer = null;
            }
            if (!shouldPionexMarkStreamActive()) return;
            if (getDesiredPionexWsSymbols().size === 0) return;
            pionexIndexReconnectTimer = setTimeout(() => {
              pionexIndexReconnectTimer = null;
              syncPionexMarkStream();
            }, 2000);
          };
        }

        function ensureMultichartStructure() {
          const root = document.getElementById("root");
          if (!root) return;
          if (document.getElementById("tabWorkspacesRoot")) return;
          const twr = document.createElement("div");
          twr.id = "tabWorkspacesRoot";
          const overlay = document.createElement("div");
          overlay.id = "tabFlowOverlay";
          overlay.style.display = "none";
          root.innerHTML = "";
          root.appendChild(twr);
          root.appendChild(overlay);
        }

        function activateWorkspace(tid) {
          const id = String(tid || "");
          if (!id) return false;
          const ws = tabWorkspaces[id];
          if (!ws || !ws.gridEl) return false;
          gridEl = ws.gridEl;
          gridFooterEl = ws.gridFooterEl;
          tileRefs = ws.tileRefs;
          mounted = ws.mounted;
          return true;
        }

        function updateTabWorkspaceVisibility() {
          const tid = String(appState.activeTabId || "");
          const twr = document.getElementById("tabWorkspacesRoot");
          if (!twr) return;
          twr.querySelectorAll(".tab-workspace").forEach((el) => {
            el.style.display = el.dataset.tabId === tid ? "block" : "none";
          });
        }

        function showTabFlowOverlay() {
          ensureMultichartStructure();
          const twr = document.getElementById("tabWorkspacesRoot");
          const overlay = document.getElementById("tabFlowOverlay");
          if (twr) twr.style.display = "none";
          if (overlay) {
            overlay.style.display = "block";
            overlay.innerHTML = "";
          }
        }

        function hideTabFlowOverlay() {
          const twr = document.getElementById("tabWorkspacesRoot");
          const overlay = document.getElementById("tabFlowOverlay");
          if (overlay) {
            overlay.style.display = "none";
            overlay.innerHTML = "";
          }
          if (twr) twr.style.display = "block";
        }

        function destroyTabWorkspace(tabId) {
          const id = String(tabId || "");
          const ws = tabWorkspaces[id];
          if (ws && ws.shellEl && ws.shellEl.parentNode) {
            ws.shellEl.parentNode.removeChild(ws.shellEl);
          }
          delete tabWorkspaces[id];
        }

        function syncTileLayoutMeta() {
          tileRefs.forEach((r, i) => {
            r.tile.dataset.index = String(i);
            r.tile.style.order = String(i);
          });
        }

        function deleteChartAt(index) {
          activateWorkspace(appState.activeTabId);
          const tab = activeTab();
          if (!appState || !tab || index < 0 || index >= tab.symbols.length) return;
          if (isAutoMode(tab)) return;
          tab.symbols.splice(index, 1);
          tab.intervals.splice(index, 1);
          saveState(appState);
          if (!gridEl || !tileRefs[index]) {
            initOrUpdateMultichart(false);
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
          activateWorkspace(appState.activeTabId);
          const tab = activeTab();
          if (!appState || !tab) return;
          if (isAutoMode(tab)) return;
          tab.symbols.push(defaultSymbols(1)[0]);
          tab.intervals.push(null);
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
          const tab = activeTab();
          const symbol = ((tab && tab.symbols[index]) || "").trim();
          const interval = effectiveInterval({ defaultInterval: appState.defaultInterval, intervals: tab ? tab.intervals : [] }, index);
          return { symbol, interval };
        }

        function maybeMountWidget(index, force) {
          const ref = tileRefs[index];
          if (!ref) return;

          const { symbol, interval: intervalRaw } = desiredConfigAt(index);
          const interval = String(intervalRaw != null ? intervalRaw : "").trim() || String(appState.defaultInterval || "60").trim() || "60";
          const theme = appState.chartTheme === "light" ? "light" : "dark";
          const barPreset = normalizeBarColorPreset(appState.barColorPreset);
          const prev = mounted[index];
          const sameInt = prev && String(prev.interval) === String(interval);
          if (!force && prev && prev.symbol === symbol && sameInt && prev.theme === theme && prev.barPreset === barPreset) return;

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
         * @param skipChart If true, update head and position strip only (after tile row permute; iframe already correct).
         */
        function refreshTile(index, forceWidget, inputsOnly, skipChart) {
          const ref = tileRefs[index];
          if (!ref) return;
          const tab = activeTab();
          ref.tile.classList.toggle("is-auto", !!(tab && isAutoMode(tab)));
          if (ref.dragHandle) {
            const auto = !!(tab && isAutoMode(tab));
            ref.dragHandle.draggable = !auto;
            ref.dragHandle.style.display = auto ? "none" : "";
            ref.dragHandle.setAttribute("aria-hidden", auto ? "true" : "false");
          }

          ref.input.value = (tab && tab.symbols[index]) || "";
          ref.input.disabled = isAutoMode(tab);
          buildIntervalOptions(ref.intervalSel, (tab && tab.intervals[index]) ?? "", true);
          if (inputsOnly) return;
          syncPositionStrip(ref, index);
          if (skipChart) return;
          maybeMountWidget(index, !!forceWidget);
        }

        function clearDnDHighlights() {
          tileRefs.forEach((r) => r.tile.classList.remove("drop-target"));
        }

        function swapTiles(from, to) {
          if (from === to) return;
          activateWorkspace(appState.activeTabId);
          const tab = activeTab();
          if (!tab) return;
          if (isAutoMode(tab)) return;
          swap(tab.symbols, from, to);
          swap(tab.intervals, from, to);
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
            const tab = activeTab();
            if (!tab) return;
            if (isAutoMode(tab)) return;
            const next = String(fullSymbol || "").trim();
            const prev = (tab.symbols[i] || "").trim();
            if (!next || next === prev) return;
            tab.symbols[i] = next;
            saveState(appState);
            refreshTile(i, true);
          }

          input.addEventListener("change", () => {
            if (isAutoMode()) return;
            const i = Number.parseInt(tile.dataset.index, 10);
            // Allow manual paste/enter of full symbol (e.g. BINANCE:BTCUSDT)
            setSymbolAndRefresh(i, input.value.trim());
          });

          let symSearchAbort = null;
          const runSearch = debounce(async () => {
            if (isAutoMode()) return;
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
            const tab = activeTab();
            if (!tab) return;
            const i = Number.parseInt(tile.dataset.index, 10);
            const v = intervalSel.value;
            const next = v === "" ? null : v;
            const prev = tab.intervals[i] == null ? null : String(tab.intervals[i]);
            if (next === prev) return;
            tab.intervals[i] = next;
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

          // Auto mode: no add/remove; keep per-chart interval but lock symbol + hide delete menu.
          if (isAutoMode()) {
            input.disabled = true;
            input.placeholder = "Auto (positions)";
            menuWrap.style.display = "none";
          }

          head.appendChild(dragHandle);
          head.appendChild(left);
          head.appendChild(right);

          const positionStrip = document.createElement("div");
          positionStrip.className = "tile-position-strip";
          positionStrip.dataset.visible = "false";

          const body = document.createElement("div");
          body.className = "tile-body";

          dragHandle.addEventListener("dragstart", (e) => {
            if (isAutoMode(activeTab())) {
              e.preventDefault();
              return;
            }
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
            if (isAutoMode(activeTab())) return;
            e.preventDefault();
            e.dataTransfer.dropEffect = "move";
            tile.classList.add("drop-target");
          });
          head.addEventListener("dragleave", () => {
            if (isAutoMode(activeTab())) return;
            tile.classList.remove("drop-target");
          });
          head.addEventListener("drop", (e) => {
            if (isAutoMode(activeTab())) {
              e.preventDefault();
              return;
            }
            e.preventDefault();
            const fromRaw = e.dataTransfer.getData("text/plain");
            const from = Number.parseInt(fromRaw, 10);
            const to = Number.parseInt(tile.dataset.index, 10);
            tile.classList.remove("drop-target");
            if (Number.isNaN(from) || Number.isNaN(to) || from === to) return;
            swapTiles(from, to);
          });

          tile.appendChild(head);
          tile.appendChild(positionStrip);
          tile.appendChild(body);

          return { tile, head, input, intervalSel, body, sugg, positionStrip, dragHandle };
        }

        function applyTabTypeChoice(next) {
          const tab = activeTab();
          if (!tab || !tabNeedsType(tab)) return;
          if (next === "manual") {
            tab.mode = "manual";
            delete tab.autoExchange;
            delete tab.autoPositionMeta;
            tab.symbols = [defaultSymbols(1)[0]];
            tab.intervals = [null];
            saveState(appState);
            syncSettingsFormFromAppState();
            ensureTradingView().then((ok) => {
              if (!ok) return;
              initOrUpdateMultichart(true);
            });
          } else if (next === "auto_positions") {
            tab.mode = "auto_positions";
            tab.autoExchange = "";
            tab.symbols = [];
            tab.intervals = [];
            tab.autoPositionMeta = [];
            tab.autoSortKey = "pnl";
            tab.autoSortDir = "desc";
            saveState(appState);
            syncSettingsFormFromAppState();
            renderAutoExchangePicker();
          } else {
            return;
          }
        }

        function hasExchangeKeys(exchangeId) {
          const keys = loadExchangeKeys() || {};
          const c = keys[exchangeId];
          return !!(c && String(c.apiKey || "").trim() && String(c.apiSecret || "").trim());
        }

        function syncExchangeKeyInputsToSettingsForm() {
          const saved = loadExchangeKeys() || {};
          [
            ["binance", "binanceApiKey", "binanceApiSecret"],
            ["bybit", "bybitApiKey", "bybitApiSecret"],
            ["okx", "okxApiKey", "okxApiSecret"],
            ["pionex", "pionexApiKey", "pionexApiSecret"]
          ].forEach(([k, kid, sid]) => {
            const elK = document.getElementById(kid);
            const elS = document.getElementById(sid);
            const c = saved[k] || {};
            if (elK) elK.value = c.apiKey || "";
            if (elS) elS.value = c.apiSecret || "";
          });
        }

        function openExchangeKeysModal(exchangeId) {
          const id = String(exchangeId || "").trim().toLowerCase();
          const def = AUTO_EXCHANGES.find((e) => e.id === id);
          const title = document.getElementById("exchangeKeysModalTitle");
          const dialog = document.getElementById("exchangeKeysDialog");
          const keyEl = document.getElementById("exchangeKeysModalApiKey");
          const secEl = document.getElementById("exchangeKeysModalApiSecret");
          if (!dialog || !keyEl || !secEl) return;
          exchangeKeysModalPendingId = id;
          if (title) title.textContent = "API keys — " + (def ? def.label : id);
          const cur = (loadExchangeKeys() || {})[id] || {};
          keyEl.value = cur.apiKey || "";
          secEl.value = cur.apiSecret || "";
          dialog.showModal();
          setTimeout(() => keyEl.focus(), 0);
        }

        function saveExchangeKeysFromModal() {
          const id = exchangeKeysModalPendingId;
          if (!id) return;
          const keyEl = document.getElementById("exchangeKeysModalApiKey");
          const secEl = document.getElementById("exchangeKeysModalApiSecret");
          const dialog = document.getElementById("exchangeKeysDialog");
          const k = keyEl && String(keyEl.value || "").trim();
          const s = secEl && String(secEl.value || "").trim();
          if (!k || !s) {
            alert("Enter both API key and secret.");
            return;
          }
          const all = Object.assign({}, loadExchangeKeys() || {});
          all[id] = { apiKey: k, apiSecret: s };
          saveExchangeKeys(all);
          exchangeKeysModalPendingId = null;
          if (dialog) dialog.close();
          syncExchangeKeyInputsToSettingsForm();
          confirmAutoExchange(id);
        }

        function confirmAutoExchange(exchangeId) {
          const id = String(exchangeId || "").trim().toLowerCase();
          if (!id) return;
          if (!hasExchangeKeys(id)) {
            openExchangeKeysModal(id);
            return;
          }
          const tab = activeTab();
          if (!tab) return;
          tab.mode = "auto_positions";
          tab.autoExchange = id;
          tab.symbols = [""];
          tab.intervals = [null];
          saveState(appState);
          syncSettingsFormFromAppState();
          ensureTradingView().then((ok) => {
            if (!ok) return;
            fetchPositionsForExchange(id)
              .then((pos) => {
                const onlyMeta = updateAutoSymbolsFromPositions(pos);
                refreshAutoPositionUIAfterFetch(onlyMeta, true, false);
              })
              .catch(() => {
                initOrUpdateMultichart(true);
              });
          });
        }

        function renderAutoExchangePicker() {
          showTabFlowOverlay();
          const overlay = document.getElementById("tabFlowOverlay");
          if (!overlay) return;
          const tab = activeTab();
          const wrap = document.createElement("div");
          wrap.className = "type-picker";

          const back = document.createElement("button");
          back.type = "button";
          back.className = "type-picker-back";
          back.textContent = "← Back to chart type";
          back.addEventListener("click", () => {
            const t = activeTab();
            if (!t) return;
            if (t.mode === "auto_positions" && !String(t.autoExchange || "").trim()) {
              t.mode = "unset";
              t.symbols = [];
              t.intervals = [];
              delete t.autoExchange;
              saveState(appState);
              syncSettingsFormFromAppState();
              renderTabTypePicker();
            }
          });

          const h = document.createElement("h3");
          h.textContent = "Choose exchange API";
          const p = document.createElement("p");
          p.textContent = "Pick where to load open positions from. Keys are stored in Settings.";

          const grid = document.createElement("div");
          grid.className = "type-picker-grid";

          AUTO_EXCHANGES.forEach((ex) => {
            const ok = hasExchangeKeys(ex.id);
            const b = document.createElement("button");
            b.type = "button";
            b.className = "type-picker-card exchange-picker-card" + (ok ? "" : " missing-keys");
            const strong = document.createElement("strong");
            strong.textContent = ex.label;
            const span = document.createElement("span");
            span.textContent = "Open positions from this account.";
            const meta = document.createElement("div");
            meta.className = "ex-meta";
            const badge = document.createElement("span");
            badge.className = "ex-badge" + (ok ? " ok" : "");
            badge.textContent = ok ? "Keys configured" : "Keys missing";
            meta.appendChild(badge);
            b.appendChild(strong);
            b.appendChild(span);
            b.appendChild(meta);
            b.addEventListener("click", () => confirmAutoExchange(ex.id));
            grid.appendChild(b);
          });

          wrap.appendChild(back);
          wrap.appendChild(h);
          wrap.appendChild(p);
          wrap.appendChild(grid);
          overlay.appendChild(wrap);
        }

        function renderTabTypePicker() {
          showTabFlowOverlay();
          const overlay = document.getElementById("tabFlowOverlay");
          if (!overlay) return;
          const tab = activeTab();
          const wrap = document.createElement("div");
          wrap.className = "type-picker";
          const h = document.createElement("h3");
          h.textContent = "Choose chart type";
          const p = document.createElement("p");
          p.textContent = tab && tab.name
            ? "Tab \"" + tab.name + "\": choose how this workspace should work."
            : "Choose how this workspace should work.";
          const grid = document.createElement("div");
          grid.className = "type-picker-grid";

          const b1 = document.createElement("button");
          b1.type = "button";
          b1.className = "type-picker-card";
          b1.innerHTML = "<strong>Manual</strong><span>Add your own symbols and manage the chart grid.</span>";
          b1.addEventListener("click", () => applyTabTypeChoice("manual"));

          const b2 = document.createElement("button");
          b2.type = "button";
          b2.className = "type-picker-card";
          b2.innerHTML = "<strong>Auto · Positions</strong><span>Charts follow open positions (configure exchange API keys in settings).</span>";
          b2.addEventListener("click", () => applyTabTypeChoice("auto_positions"));

          grid.appendChild(b1);
          grid.appendChild(b2);
          wrap.appendChild(h);
          wrap.appendChild(p);
          wrap.appendChild(grid);
          overlay.appendChild(wrap);
        }

        function ensureGrid() {
          const tid = String(appState.activeTabId || "");
          ensureMultichartStructure();
          const twr = document.getElementById("tabWorkspacesRoot");
          if (!twr) return;

          let ws = tabWorkspaces[tid];
          if (!ws) {
            tabWorkspaces[tid] = { shellEl: null, autoSummaryEl: null, gridEl: null, gridFooterEl: null, tileRefs: [], mounted: [] };
            ws = tabWorkspaces[tid];
          }

          if (!ws.gridEl) {
            const shell = document.createElement("div");
            shell.className = "tab-workspace";
            shell.dataset.tabId = tid;
            shell.dataset.autoSummary = "0";
            shell.style.display = "none";

            ws.autoSummaryEl = document.createElement("div");
            ws.autoSummaryEl.className = "auto-summary-bar";
            ws.autoSummaryEl.dataset.visible = "false";

            ws.gridEl = document.createElement("div");
            ws.gridEl.className = "grid";
            ws.gridFooterEl = document.createElement("div");
            ws.gridFooterEl.className = "add-chart-footer";
            const addBtn = document.createElement("button");
            addBtn.type = "button";
            addBtn.className = "add-chart-btn";
            addBtn.textContent = "+";
            addBtn.title = "Add chart";
            addBtn.setAttribute("aria-label", "Add chart");
            addBtn.addEventListener("click", () => addChart());
            ws.gridFooterEl.appendChild(addBtn);

            shell.appendChild(ws.autoSummaryEl);
            shell.appendChild(ws.gridEl);
            shell.appendChild(ws.gridFooterEl);
            twr.appendChild(shell);
            ws.shellEl = shell;
            ws.tileRefs = [];
            ws.mounted = [];
          } else if (ws.shellEl && ws.gridEl && !ws.autoSummaryEl) {
            const bar = document.createElement("div");
            bar.className = "auto-summary-bar";
            bar.dataset.visible = "false";
            ws.shellEl.insertBefore(bar, ws.gridEl);
            ws.autoSummaryEl = bar;
            if (ws.shellEl.dataset.autoSummary == null) ws.shellEl.dataset.autoSummary = "0";
          }
          gridEl = ws.gridEl;
          gridFooterEl = ws.gridFooterEl;
          tileRefs = ws.tileRefs;
          mounted = ws.mounted;

          if (gridFooterEl) {
            gridFooterEl.style.display = isAutoMode() ? "none" : "flex";
          }

          gridEl.style.gridTemplateColumns = "repeat(" + appState.cols + ", minmax(0, 1fr))";
          gridEl.style.setProperty("--chart-tile-min", String(appState.chartHeightPx) + "px");

          const tab = activeTab();
          const want = tab ? tab.symbols.length : 0;
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

        /**
         * @param skipPollResync When true, do not restart the auto-positions poll timer (used after a poll tick to avoid re-entrancy).
         */
        function initOrUpdateMultichart(forceAllWidgets, skipPollResync) {
          ensureMultichartStructure();
          const tab = activeTab();
          if (!tab || tabNeedsType(tab)) {
            renderTabTypePicker();
            syncAutoPositionsPolling();
            syncAutoSortBar();
            return;
          }
          if (isAutoMode(tab) && autoNeedsExchange(tab)) {
            renderAutoExchangePicker();
            syncAutoPositionsPolling();
            syncAutoSortBar();
            return;
          }
          hideTabFlowOverlay();
          ensureGrid();
          updateTabWorkspaceVisibility();
          activateWorkspace(appState.activeTabId);
          let appliedAutoSortPermute = false;
          if (pendingAutoSortPermute && isAutoMode(tab)) {
            const p = pendingAutoSortPermute;
            pendingAutoSortPermute = null;
            const ok =
              p &&
              p.perm &&
              p.perm.length === tileRefs.length &&
              String(appState.activeTabId) === String(p.tabId);
            if (ok) {
              applyAutoSortPermutationToTiles(p.perm);
              appliedAutoSortPermute = true;
            }
          } else if (pendingAutoSortPermute) {
            pendingAutoSortPermute = null;
          }
          if (!appliedAutoSortPermute) {
            tileRefs.forEach((_, i) => refreshTile(i, !!forceAllWidgets));
          }
          syncAutoSummaryBar();
          syncAutoSortBar();
          if (!skipPollResync) syncAutoPositionsPolling();
        }

        async function fetchPositionsForExchange(exchangeId) {
          const id = String(exchangeId || "").trim().toLowerCase();
          const all = loadExchangeKeys() || {};
          const cred = all[id];
          if (!cred || !String(cred.apiKey || "").trim() || !String(cred.apiSecret || "").trim()) {
            throw new Error("missing exchange credentials");
          }
          const res = await fetch("/api/positions", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ exchanges: { [id]: cred } })
          });
          if (!res.ok) throw new Error("positions request failed");
          const data = await res.json();
          if (data && Array.isArray(data.errors) && data.errors.length) {
            console.warn("positions API:", data.errors.join("; "));
          }
          return data;
        }

        function stopAutoPositionsPolling() {
          if (autoPositionsPollTimer != null) {
            clearInterval(autoPositionsPollTimer);
            autoPositionsPollTimer = null;
          }
          stopPionexMarkStream();
        }

        function runAutoPositionsPollTick() {
          if (!appState || document.hidden) return;
          if (activePageFromHash() !== "multichart") return;
          const tab = activeTab();
          if (!tab || !isAutoMode(tab) || autoNeedsExchange(tab)) return;
          const ex = String(tab.autoExchange || "").trim().toLowerCase();
          if (!ex) return;
          const tid = String(tab.id);
          fetchPositionsForExchange(ex)
            .then((pos) => {
              if (!appState || String(appState.activeTabId) !== tid) return;
              const cur = activeTab();
              if (!cur || String(cur.id) !== tid || !isAutoMode(cur)) return;
              if (String(cur.autoExchange || "").trim().toLowerCase() !== ex) return;
              const onlyMeta = updateAutoSymbolsFromPositions(pos);
              refreshAutoPositionUIAfterFetch(onlyMeta, false, true);
            })
            .catch(() => {});
        }

        function syncAutoPositionsPolling() {
          stopAutoPositionsPolling();
          if (!appState || document.hidden) return;
          if (activePageFromHash() !== "multichart") return;
          const tab = activeTab();
          if (!tab || tabNeedsType(tab) || !isAutoMode(tab) || autoNeedsExchange(tab)) return;
          const ex = String(tab.autoExchange || "").trim().toLowerCase();
          if (!ex) return;
          runAutoPositionsPollTick();
          autoPositionsPollTimer = setInterval(runAutoPositionsPollTick, AUTO_POSITIONS_POLL_MS);
        }

        /**
         * Map any exchange futures/perp position ticker to a TradingView Binance USDⓈ-M **perpetual** chart.
         * Spot-style BINANCE:BTCUSDT is wrong for futures positions — use BINANCE:BTCUSDT.P (TradingView convention).
         */
        function positionSymbolToBinanceTV(p) {
          let s = String(p && p.symbol || "").trim();
          if (!s) return "";
          if (s.includes(":")) {
            s = s.split(":").pop().trim();
          }
          s = s.toUpperCase();
          s = s.replace(/\.P$/i, "");
          s = s.replace(/\.PERP$/i, "");
          s = s.replace(/_PERP$/i, "");
          s = s.replace(/-PERP$/i, "");
          s = s.replace(/-SWAP$/i, "");
          s = s.replace(/_PERPETUAL$/i, "");
          s = s.replace(/-PERPETUAL$/i, "");
          s = s.replace(/_/g, "");
          s = s.replace(/-/g, "");
          if (!s) return "";
          return "BINANCE:" + s + ".P";
        }

        function snapPositionMetaFromRow(p) {
          if (!p || typeof p !== "object") return null;
          return {
            exchange: String(p.exchange || "").trim(),
            rawSymbol: String(p.symbol || "").trim(),
            side: String(p.side || "").trim(),
            size: String(p.size || "").trim(),
            avgPrice: String(p.avgPrice || "").trim(),
            markPrice: String(p.markPrice || "").trim(),
            unrealizedPnL: String(p.unrealizedPnL || "").trim(),
            leverage: String(p.leverage || "").trim(),
            liquidationPrice: String(p.liquidationPrice || "").trim()
          };
        }

        function formatNumDisplay(s) {
          const t = String(s || "").trim();
          if (!t) return "—";
          const n = Number.parseFloat(t.replace(/,/g, ""));
          if (Number.isFinite(n)) {
            const abs = Math.abs(n);
            const dec = abs >= 1000 || (abs >= 1 && abs < 100) ? 2 : abs >= 0.01 ? 4 : 6;
            return n.toLocaleString(undefined, { maximumFractionDigits: dec, minimumFractionDigits: 0 });
          }
          return t;
        }

        /** Unrealized PnL only: at most one digit after the decimal separator. */
        function formatPnlDisplay(s) {
          const t = String(s || "").trim();
          if (!t) return "—";
          const n = Number.parseFloat(t.replace(/,/g, ""));
          if (Number.isFinite(n)) {
            return n.toLocaleString(undefined, { maximumFractionDigits: 1, minimumFractionDigits: 0 });
          }
          return t;
        }

        function parseMetaNum(s) {
          const n = Number.parseFloat(String(s || "").replace(/,/g, "").trim());
          return Number.isFinite(n) ? n : NaN;
        }

        /**
         * Denominator for uPnL %: imputed initial margin (|size|×avg÷leverage) when leverage is known,
         * else position notional |size|×avg (PnL vs notional).
         */
        function estimatePnlPercentDenominator(m) {
          if (!m || typeof m !== "object") return null;
          const sz = Math.abs(parseMetaNum(m.size));
          const avg = parseMetaNum(m.avgPrice);
          if (!Number.isFinite(sz) || !Number.isFinite(avg) || sz <= 0 || avg <= 0) return null;
          const notional = sz * avg;
          const lev = parseMetaNum(m.leverage);
          if (Number.isFinite(lev) && lev > 0) return notional / lev;
          return notional;
        }

        function unrealizedPnlPercentFromMeta(m) {
          const pnl = parseMetaNum(m.unrealizedPnL);
          const denom = estimatePnlPercentDenominator(m);
          if (!Number.isFinite(pnl) || denom == null || !(denom > 0)) return null;
          return (pnl / denom) * 100;
        }

        function formatPnlPercentDisplay(pct) {
          if (!Number.isFinite(pct)) return "";
          const sign = pct < 0 ? "-" : pct > 0 ? "+" : "";
          const abs = Math.abs(pct);
          const dec = abs >= 100 ? 1 : 2;
          return sign + abs.toLocaleString(undefined, { maximumFractionDigits: dec, minimumFractionDigits: 0 }) + "%";
        }

        function formatSideDisplay(s) {
          const t = String(s || "").trim();
          if (!t) return "";
          const u = t.toUpperCase();
          if (u === "LONG" || u === "SHORT") return u;
          if (u.includes("LONG")) return "LONG";
          if (u.includes("SHORT")) return "SHORT";
          return t;
        }

        /** @returns {"long"|"short"|null} */
        function sideToneFromPositionSide(s) {
          const u = String(s || "").toUpperCase();
          if (u.includes("LONG")) return "long";
          if (u.includes("SHORT")) return "short";
          return null;
        }

        function formatLeverageDisplay(s) {
          const t = String(s || "").trim();
          if (!t) return "";
          const n = Number.parseFloat(String(t).replace(/,/g, "").replace(/×/g, ""));
          if (Number.isFinite(n)) {
            return n.toLocaleString(undefined, { maximumFractionDigits: 2 }) + "×";
          }
          return t;
        }

        function syncAutoSummaryBar() {
          if (!appState) return;
          const tid = String(appState.activeTabId || "");
          const ws = tabWorkspaces[tid];
          if (!ws || !ws.autoSummaryEl || !ws.shellEl) return;
          const tab = activeTab();
          const bar = ws.autoSummaryEl;
          if (!tab || !isAutoMode(tab) || autoNeedsExchange(tab)) {
            bar.dataset.visible = "false";
            bar.innerHTML = "";
            ws.shellEl.dataset.autoSummary = "0";
            return;
          }
          const meta = Array.isArray(tab.autoPositionMeta) ? tab.autoPositionMeta : [];
          let totalPnL = 0;
          let hasPnL = false;
          let counted = 0;
          let totalMargin = 0;
          let hasMargin = false;
          let totalNotional = 0;
          let hasNotional = false;
          meta.forEach((m) => {
            if (!m || typeof m !== "object") return;
            counted++;
            const pnl = Number.parseFloat(String(m.unrealizedPnL || "").replace(/,/g, ""));
            if (Number.isFinite(pnl)) {
              totalPnL += pnl;
              hasPnL = true;
            }
            const sz = Math.abs(parseMetaNum(m.size));
            const avg = parseMetaNum(m.avgPrice);
            if (Number.isFinite(sz) && Number.isFinite(avg) && sz > 0 && avg > 0) {
              totalNotional += sz * avg;
              hasNotional = true;
            }
            const d = estimatePnlPercentDenominator(m);
            if (d != null && d > 0) {
              totalMargin += d;
              hasMargin = true;
            }
          });
          if (counted === 0) {
            counted = (tab.symbols || []).filter((s) => String(s || "").trim()).length;
          }
          const exId = String(tab.autoExchange || "").trim().toLowerCase();
          const def = AUTO_EXCHANGES.find((e) => e.id === exId);
          const exLabel = def ? def.label : (exId ? exId.charAt(0).toUpperCase() + exId.slice(1) : "—");

          bar.dataset.visible = "true";
          bar.innerHTML = "";
          ws.shellEl.dataset.autoSummary = "1";

          function addCard(label, valueText, opts) {
            opts = opts || {};
            const card = document.createElement("div");
            let cls = "auto-summary-card";
            if (opts.pnl) cls += " auto-summary-card--pnl";
            if (opts.pnlTone) cls += " " + opts.pnlTone;
            card.className = cls;
            const lab = document.createElement("span");
            lab.className = "auto-summary-label";
            lab.textContent = label;
            const val = document.createElement("span");
            val.className = "auto-summary-value" + (opts.pnlValClass ? " " + opts.pnlValClass : "");
            if (opts.pnlParts) {
              val.appendChild(document.createTextNode(opts.pnlParts.main));
              if (opts.pnlParts.pct) {
                const pctSpan = document.createElement("span");
                pctSpan.className = "auto-summary-pct";
                pctSpan.textContent = opts.pnlParts.pct;
                val.appendChild(pctSpan);
              }
            } else {
              val.textContent = valueText;
            }
            card.appendChild(lab);
            card.appendChild(val);
            bar.appendChild(card);
          }

          addCard("Exchange", exLabel);
          addCard("Positions", String(counted));
          addCard("Total margin", hasMargin ? formatNumDisplay(String(totalMargin)) : "—");
          addCard("Notional", hasNotional ? formatNumDisplay(String(totalNotional)) : "—");
          if (hasMargin) {
            const bal = totalMargin + (hasPnL ? totalPnL : 0);
            addCard("Total Balance", formatNumDisplay(String(bal)));
          } else {
            addCard("Total Balance", "—");
          }

          if (hasPnL) {
            const pnlClass = totalPnL >= 0 ? "pnl-pos" : "pnl-neg";
            const sign = totalPnL >= 0 ? "+" : "";
            const mainAmt = sign + formatPnlDisplay(String(totalPnL));
            let totalPct = null;
            if (hasMargin && totalMargin > 0) {
              const tp = (totalPnL / totalMargin) * 100;
              if (Number.isFinite(tp)) totalPct = tp;
            }
            if (totalPct != null) {
              addCard("Total UPNL", "", {
                pnl: true,
                pnlTone: pnlClass,
                pnlValClass: pnlClass,
                pnlParts: { main: mainAmt, pct: " (" + formatPnlPercentDisplay(totalPct) + ")" }
              });
            } else {
              addCard("Total UPNL", mainAmt, { pnl: true, pnlTone: pnlClass, pnlValClass: pnlClass });
            }
          } else {
            addCard("Total UPNL", "—");
          }
        }

        function syncAutoSortBar() {
          if (!appState) return;
          const bar = document.getElementById("tabsbarSort");
          if (!bar) return;
          const tab = activeTab();
          if (activePageFromHash() !== "multichart" || !tab || !isAutoMode(tab) || autoNeedsExchange(tab)) {
            bar.dataset.visible = "false";
            bar.innerHTML = "";
            delete bar.dataset.built;
            return;
          }
          bar.dataset.visible = "true";
          if (bar.dataset.built !== "true") {
            bar.dataset.built = "true";
            const lab = document.createElement("span");
            lab.className = "auto-sort-label";
            lab.textContent = "Sort";
            const keySel = document.createElement("select");
            keySel.className = "auto-sort-key";
            keySel.title = "Sort charts by";
            [["symbol", "Symbol"], ["pnl", "PnL"], ["size", "Size"], ["entry", "Entry price"]].forEach(([v, lbl]) => {
              const o = document.createElement("option");
              o.value = v;
              o.textContent = lbl;
              keySel.appendChild(o);
            });
            const dirSel = document.createElement("select");
            dirSel.className = "auto-sort-dir-select";
            dirSel.title = "Sort direction";
            [["asc", "Ascending"], ["desc", "Descending"]].forEach(([v, lbl]) => {
              const o = document.createElement("option");
              o.value = v;
              o.textContent = lbl;
              dirSel.appendChild(o);
            });
            function applySortFromUI() {
              const cur = activeTab();
              if (!cur || !isAutoMode(cur)) return;
              cur.autoSortKey = normalizeAutoSortKey(keySel.value);
              cur.autoSortDir = normalizeAutoSortDir(dirSel.value);
              const perm = sortAutoTabCharts(cur);
              saveState(appState);
              if (perm && perm.length === tileRefs.length) {
                applyAutoSortPermutationToTiles(perm);
              } else {
                syncTileLayoutMeta();
                tileRefs.forEach((_, i) => refreshTile(i, false));
              }
            }
            keySel.addEventListener("change", applySortFromUI);
            dirSel.addEventListener("change", applySortFromUI);
            const controls = document.createElement("span");
            controls.className = "auto-sort-controls";
            controls.appendChild(keySel);
            controls.appendChild(dirSel);
            bar.appendChild(lab);
            bar.appendChild(controls);
            bar._autoSortKeySel = keySel;
            bar._autoSortDirSel = dirSel;
          }
          const keySel = bar._autoSortKeySel;
          const dirSel = bar._autoSortDirSel;
          if (keySel && dirSel) {
            const wantK = normalizeAutoSortKey(tab.autoSortKey);
            const wantD = normalizeAutoSortDir(tab.autoSortDir);
            if (keySel.value !== wantK) keySel.value = wantK;
            if (dirSel.value !== wantD) dirSel.value = wantD;
          }
        }

        function syncPositionStrip(ref, index) {
          const tab = activeTab();
          if (!ref.positionStrip) return;
          if (!tab || !isAutoMode(tab)) {
            ref.positionStrip.dataset.visible = "false";
            ref.positionStrip.innerHTML = "";
            return;
          }
          ref.positionStrip.dataset.visible = "true";
          const m = Array.isArray(tab.autoPositionMeta) ? tab.autoPositionMeta[index] : null;
          ref.positionStrip.innerHTML = "";
          if (!m || typeof m !== "object") {
            const span = document.createElement("span");
            span.className = "pos-chip";
            span.textContent = "Loading position…";
            ref.positionStrip.appendChild(span);
            return;
          }
          const detailKeys = ["rawSymbol", "side", "size", "avgPrice", "markPrice", "leverage", "liquidationPrice"];
          const hasAny = detailKeys.some((k) => String(m[k] || "").trim()) || String(m.unrealizedPnL || "").trim();
          if (!hasAny) {
            const span = document.createElement("span");
            span.className = "pos-chip";
            span.textContent = "No details from exchange";
            ref.positionStrip.appendChild(span);
            return;
          }
          const mainEl = document.createElement("div");
          mainEl.className = "pos-metrics-row";
          function addStat(target, label, valueText, opts) {
            opts = opts || {};
            const v = String(valueText || "").trim();
            if (!v && !opts.showEmpty) return;
            const wrap = document.createElement("span");
            wrap.className = "pos-stat";
            const tone = opts.sideTone;
            if (tone === "long") wrap.classList.add("pos-stat--long");
            else if (tone === "short") wrap.classList.add("pos-stat--short");
            const lab = document.createElement("span");
            lab.className = "pos-stat-label";
            lab.textContent = label;
            const val = document.createElement("span");
            val.className = "pos-stat-value";
            if (opts.valueClass) val.classList.add(opts.valueClass);
            val.textContent = v || "—";
            wrap.appendChild(lab);
            wrap.appendChild(val);
            target.appendChild(wrap);
          }
          const symRaw = String(m.rawSymbol || "").trim();
          if (symRaw) {
            addStat(mainEl, "Contract", symRaw, { valueClass: "pos-stat-value--contract" });
          }
          const sideRaw = String(m.side || "").trim();
          if (sideRaw) {
            addStat(mainEl, "Side", formatSideDisplay(sideRaw), { sideTone: sideToneFromPositionSide(sideRaw) });
          }
          const sizeRaw = String(m.size || "").trim();
          if (sizeRaw) {
            addStat(mainEl, "Size", formatNumDisplay(sizeRaw));
          }
          const entryRaw = String(m.avgPrice || "").trim();
          if (entryRaw) {
            addStat(mainEl, "Entry", formatNumDisplay(entryRaw), { valueClass: "pos-stat-value--emph" });
          }
          const markRaw = String(m.markPrice || "").trim();
          if (markRaw) {
            addStat(mainEl, "Mark", formatNumDisplay(markRaw), { valueClass: "pos-stat-value--emph" });
          }
          const levRaw = String(m.leverage || "").trim();
          if (levRaw) {
            addStat(mainEl, "Lev", formatLeverageDisplay(levRaw));
          }
          const liqRaw = String(m.liquidationPrice || "").trim();
          if (liqRaw) {
            addStat(mainEl, "Liq", formatNumDisplay(liqRaw));
          }

          const pnlEl = document.createElement("div");
          pnlEl.className = "pos-pnl-block";
          const pnlLbl = document.createElement("span");
          pnlLbl.className = "pos-pnl-label";
          pnlLbl.textContent = "uPnL";
          const pnlVal = document.createElement("span");
          pnlVal.className = "pos-pnl-value";
          const pnlRaw = String(m.unrealizedPnL || "").trim();
          if (pnlRaw) {
            const n = Number.parseFloat(pnlRaw.replace(/,/g, ""));
            const pct = unrealizedPnlPercentFromMeta(m);
            if (Number.isFinite(n)) {
              pnlEl.classList.add(n >= 0 ? "pnl-pos" : "pnl-neg");
              pnlVal.classList.add(n >= 0 ? "pnl-pos" : "pnl-neg");
              const mainText = (n >= 0 ? "+" : "") + formatPnlDisplay(pnlRaw);
              pnlVal.textContent = "";
              pnlVal.appendChild(document.createTextNode(mainText));
              if (pct != null) {
                pnlVal.title =
                  "uPnL as % of estimated margin (|size|×avg price÷leverage), or of notional if leverage is missing";
                const pctSpan = document.createElement("span");
                pctSpan.className = "pos-pnl-pct";
                pctSpan.textContent = " (" + formatPnlPercentDisplay(pct) + ")";
                pnlVal.appendChild(pctSpan);
              }
            } else {
              pnlVal.textContent = formatPnlDisplay(pnlRaw);
            }
          } else {
            pnlEl.classList.add("pos-pnl-block-empty");
            pnlVal.classList.add("pos-pnl-value-muted");
            pnlVal.textContent = "—";
          }
          pnlEl.appendChild(pnlLbl);
          pnlEl.appendChild(pnlVal);

          ref.positionStrip.appendChild(mainEl);
          ref.positionStrip.appendChild(pnlEl);
          if (!mainEl.children.length && !pnlRaw) {
            const span = document.createElement("span");
            span.className = "pos-chip";
            span.textContent = "—";
            mainEl.appendChild(span);
          }
        }

        /**
         * @returns {boolean} true if charts/rows unchanged — only strips + summary need refresh (no initOrUpdateMultichart).
         */
        function updateAutoSymbolsFromPositions(pos) {
          const rows = Array.isArray(pos && pos.positions) ? pos.positions : [];
          const seen = new Set();
          const out = [];
          const meta = [];
          rows.forEach((p) => {
            const full = positionSymbolToBinanceTV(p);
            if (!full || seen.has(full)) return;
            seen.add(full);
            out.push(full);
            meta.push(snapPositionMetaFromRow(p));
          });
          const tab = activeTab();
          if (!tab) return false;
          if (out.length === 0) {
            tab.symbols = [""];
            tab.intervals = [null];
            tab.autoPositionMeta = [null];
            sortAutoTabCharts(tab);
            saveState(appState);
            return false;
          }
          const prevList = (tab.symbols || []).map((s) => String(s || "").trim()).filter(Boolean);
          const prevTileCount = tileRefs.length;
          const outTrim = out.map((s) => String(s || "").trim());
          const symBefore = JSON.stringify(tab.symbols);
          const intBefore = JSON.stringify(tab.intervals);

          if (prevList.length === out.length && prevList.length > 0 && sameSymbolMultiset(prevList, outTrim)) {
            const metaBySym = new Map();
            out.forEach((s, i) => {
              metaBySym.set(String(s).trim(), { full: s, row: meta[i] });
            });
            const order = prevList.filter((ps) => metaBySym.has(ps));
            if (order.length === out.length) {
              const nextSym = [];
              const nextMeta = [];
              const nextInt = [];
              for (let oi = 0; oi < order.length; oi++) {
                const ps = order[oi];
                const entry = metaBySym.get(ps);
                nextSym.push(entry.full);
                nextMeta.push(entry.row);
                const j = tab.symbols.findIndex((x) => String(x || "").trim() === ps);
                nextInt.push(j >= 0 && Array.isArray(tab.intervals) && tab.intervals[j] !== undefined ? tab.intervals[j] : null);
              }
              tab.symbols = nextSym;
              while (nextInt.length < nextSym.length) nextInt.push(null);
              tab.intervals = nextInt.slice(0, nextSym.length);
              tab.autoPositionMeta = nextMeta;
              while (tab.autoPositionMeta.length < tab.symbols.length) tab.autoPositionMeta.push(null);
              if (tab.autoPositionMeta.length > tab.symbols.length) tab.autoPositionMeta = tab.autoPositionMeta.slice(0, tab.symbols.length);
              pendingAutoSortPermute = null;
              saveState(appState);
              return true;
            }
          }

          tab.symbols = out;
          while (tab.intervals.length < out.length) tab.intervals.push(null);
          if (tab.intervals.length > out.length) tab.intervals = tab.intervals.slice(0, out.length);
          tab.autoPositionMeta = meta;
          while (tab.autoPositionMeta.length < tab.symbols.length) tab.autoPositionMeta.push(null);
          if (tab.autoPositionMeta.length > tab.symbols.length) tab.autoPositionMeta = tab.autoPositionMeta.slice(0, tab.symbols.length);
          const perm = sortAutoTabCharts(tab);
          const symAfter = JSON.stringify(tab.symbols);
          const intAfter = JSON.stringify(tab.intervals);
          const onlyMetaChanged = symBefore === symAfter && intBefore === intAfter;
          const canPermute =
            perm &&
            sameSymbolMultiset(prevList, outTrim) &&
            prevTileCount === perm.length;
          pendingAutoSortPermute = canPermute ? { tabId: String(tab.id), perm } : null;
          saveState(appState);
          return onlyMetaChanged;
        }

        function refreshAutoPositionUIAfterFetch(onlyMeta, forceAllWidgets, skipPollResync) {
          const tab = activeTab();
          const gridOk =
            tab &&
            tileRefs.length > 0 &&
            tileRefs.length === (tab.symbols || []).length &&
            (tab.symbols || []).length > 0;
          if (onlyMeta && gridOk) {
            activateWorkspace(appState.activeTabId);
            tileRefs.forEach((ref, i) => syncPositionStrip(ref, i));
            syncAutoSummaryBar();
            syncAutoSortBar();
            syncPionexMarkStream();
            return;
          }
          initOrUpdateMultichart(!!forceAllWidgets, !!skipPollResync);
          syncPionexMarkStream();
        }

        function getTabById(tabId) {
          if (!appState || !Array.isArray(appState.tabs)) return null;
          const id = String(tabId || "");
          return appState.tabs.find((x) => String(x.id) === id) || null;
        }

        function syncTabSettingsExchangeVisibility() {
          const modeEl = document.getElementById("tabSettingsMode");
          const wrap = document.getElementById("tabSettingsExchangeWrap");
          if (!modeEl || !wrap) return;
          wrap.style.display = modeEl.value === "auto_positions" ? "block" : "none";
        }

        function populateTabSettingsDialog(tab) {
          const nameEl = document.getElementById("tabSettingsName");
          const modeEl = document.getElementById("tabSettingsMode");
          const exEl = document.getElementById("tabSettingsExchange");
          if (nameEl) nameEl.value = String(tab.name || "Tab");
          if (modeEl) {
            if (tabNeedsType(tab)) modeEl.value = "unset";
            else modeEl.value = isAutoMode(tab) ? "auto_positions" : "manual";
          }
          if (exEl) exEl.value = String(tab.autoExchange || "").trim().toLowerCase();
          syncTabSettingsExchangeVisibility();
        }

        function openTabSettingsDialog(tabId) {
          const tab = getTabById(tabId);
          if (!tab || !appState) return;
          tabSettingsTargetId = String(tabId);
          populateTabSettingsDialog(tab);
          const dlg = document.getElementById("tabSettingsDialog");
          if (dlg) dlg.showModal();
        }

        function applyTabSettingsFromDialog() {
          const tab = getTabById(tabSettingsTargetId);
          if (!tab || !appState) return;
          const nameEl = document.getElementById("tabSettingsName");
          const modeEl = document.getElementById("tabSettingsMode");
          const exEl = document.getElementById("tabSettingsExchange");
          if (nameEl) tab.name = String(nameEl.value || "").trim() || "Tab";
          const next = modeEl ? modeEl.value : "manual";
          if (next === "unset") {
            if (normalizeTabMode(tab.mode) !== "unset") {
              tab.mode = "unset";
              tab.symbols = [];
              tab.intervals = [];
              delete tab.autoExchange;
              delete tab.autoPositionMeta;
              delete tab.autoSortKey;
              delete tab.autoSortDir;
            }
          } else if (next === "manual") {
            if (normalizeTabMode(tab.mode) !== "manual") {
              tab.mode = "manual";
              delete tab.autoExchange;
              delete tab.autoPositionMeta;
              delete tab.autoSortKey;
              delete tab.autoSortDir;
              if (!tab.symbols || tab.symbols.length === 0) {
                tab.symbols = [defaultSymbols(1)[0]];
                tab.intervals = [null];
              }
            }
          } else if (next === "auto_positions") {
            const wasAuto = normalizeTabMode(tab.mode) === "auto_positions";
            if (!wasAuto) {
              tab.mode = "auto_positions";
              tab.symbols = [];
              tab.intervals = [];
              tab.autoPositionMeta = [];
              tab.autoSortKey = "pnl";
              tab.autoSortDir = "desc";
            }
            tab.autoExchange = String(exEl && exEl.value || "").trim().toLowerCase();
          }
          saveState(appState);
          const dlg = document.getElementById("tabSettingsDialog");
          if (dlg) dlg.close();
          tabSettingsTargetId = null;
          renderTabsbar();
          const wasActive = String(tab.id) === String(appState.activeTabId);
          if (!wasActive) return;
          initOrUpdateMultichart(false);
          ensureTradingView().then(async (ok) => {
            if (!ok) return;
            const t = activeTab();
            if (t && isAutoMode(t) && !autoNeedsExchange(t)) {
              try {
                const pos = await fetchPositionsForExchange(t.autoExchange);
                const onlyMeta = updateAutoSymbolsFromPositions(pos);
                refreshAutoPositionUIAfterFetch(onlyMeta, false, false);
              } catch (_) {
                initOrUpdateMultichart(false);
              }
            } else {
              initOrUpdateMultichart(false);
            }
          });
        }

        function renderTabsbar() {
          const el = document.getElementById("tabsbar");
          if (!el || !appState) return;
          el.innerHTML = "";
          const tabs = Array.isArray(appState.tabs) ? appState.tabs : [];
          tabs.forEach((t, idx) => {
            const btn = document.createElement("button");
            btn.type = "button";
            btn.className = "tab";
            btn.dataset.active = String(String(t.id) === String(appState.activeTabId));
            btn.title = "Switch tab";

            const label = document.createElement("span");
            label.textContent = String(t.name || "Tab");
            btn.appendChild(label);

            const gear = document.createElement("span");
            gear.className = "tab-gear";
            gear.setAttribute("role", "button");
            gear.setAttribute("tabindex", "0");
            gear.title = "Tab settings";
            gear.setAttribute("aria-label", "Tab settings");
            gear.textContent = "\u2699";
            gear.addEventListener("click", (e) => {
              e.stopPropagation();
              e.preventDefault();
              openTabSettingsDialog(t.id);
            });
            gear.addEventListener("keydown", (e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                e.stopPropagation();
                openTabSettingsDialog(t.id);
              }
            });
            btn.appendChild(gear);

            btn.addEventListener("click", () => {
              appState.activeTabId = String(t.id);
              saveState(appState);
              syncSettingsFormFromAppState();
              renderTabsbar();
              const cur = activeTab();
              if (cur && isAutoMode(cur) && !autoNeedsExchange(cur)) {
                fetchPositionsForExchange(cur.autoExchange)
                  .then((pos) => {
                    const onlyMeta = updateAutoSymbolsFromPositions(pos);
                    refreshAutoPositionUIAfterFetch(onlyMeta, false, false);
                  })
                  .catch(() => initOrUpdateMultichart(false));
              } else {
                initOrUpdateMultichart(false);
              }
            });

            btn.addEventListener("dblclick", (e) => {
              e.preventDefault();
              const next = prompt("Tab name", String(t.name || "Tab"));
              if (next == null) return;
              t.name = String(next).trim() || "Tab";
              saveState(appState);
              renderTabsbar();
            });

            if (tabs.length > 1 && idx !== 0) {
              const x = document.createElement("span");
              x.className = "x";
              x.textContent = "×";
              x.title = "Close tab";
              x.addEventListener("click", (e) => {
                e.stopPropagation();
                const ok = confirm('Close tab "' + String(t.name || "Tab") + '"?');
                if (!ok) return;
                const closedId = String(t.id);
                const closeIdx = appState.tabs.findIndex((tt) => String(tt.id) === String(t.id));
                if (closeIdx >= 0) appState.tabs.splice(closeIdx, 1);
                destroyTabWorkspace(closedId);
                if (!appState.tabs.length) appState.tabs.push(normalizeTab({ name: "Main" }));
                if (!appState.tabs.some((tt) => String(tt.id) === String(appState.activeTabId))) {
                  appState.activeTabId = String(appState.tabs[0].id);
                }
                saveState(appState);
                syncSettingsFormFromAppState();
                renderTabsbar();
                initOrUpdateMultichart(false);
              });
              btn.appendChild(x);
            }

            el.appendChild(btn);
          });

          const add = document.createElement("button");
          add.type = "button";
          add.className = "tab-add";
          add.textContent = "+";
          add.title = "Add tab";
          add.addEventListener("click", () => {
            const name = prompt("New tab name", "New tab");
            if (name == null) return;
            const tab = normalizeTab({ name: String(name).trim() || "New tab", mode: "unset", symbols: [], intervals: [] });
            appState.tabs.push(tab);
            appState.activeTabId = tab.id;
            saveState(appState);
            syncSettingsFormFromAppState();
            renderTabsbar();
            initOrUpdateMultichart(true);
          });
          el.appendChild(add);
          syncAutoSortBar();
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
            stopAutoPositionsPolling();
            const sortBar = document.getElementById("tabsbarSort");
            if (sortBar) {
              sortBar.dataset.visible = "false";
              sortBar.innerHTML = "";
              delete sortBar.dataset.built;
            }
            root.innerHTML = '<div class="empty">Select Multichart from the sidebar.</div>';
            tabWorkspaces = Object.create(null);
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
          const state = buildTabbedState(stored);

          const settingsDialog = document.getElementById("settingsDialog");
          const settingsOpen = document.getElementById("settingsOpen");
          const colsEl = document.getElementById("cols");
          const defaultIntervalEl = document.getElementById("defaultInterval");
          const chartHeightEl = document.getElementById("chartHeight");
          const chartThemeEl = document.getElementById("chartTheme");
          const barColorPresetEl = document.getElementById("barColorPreset");
          const binanceApiKeyEl = document.getElementById("binanceApiKey");
          const binanceApiSecretEl = document.getElementById("binanceApiSecret");
          const bybitApiKeyEl = document.getElementById("bybitApiKey");
          const bybitApiSecretEl = document.getElementById("bybitApiSecret");
          const okxApiKeyEl = document.getElementById("okxApiKey");
          const okxApiSecretEl = document.getElementById("okxApiSecret");
          const pionexApiKeyEl = document.getElementById("pionexApiKey");
          const pionexApiSecretEl = document.getElementById("pionexApiSecret");

          colsEl.value = String(state.cols);
          buildIntervalOptions(defaultIntervalEl, String(state.defaultInterval), false);
          chartHeightEl.value = String(state.chartHeightPx);
          buildChartThemeOptions(chartThemeEl, state.chartTheme);
          buildBarColorPresetOptions(barColorPresetEl, state.barColorPreset);
          updateTopbarDefaultIntervalDisplay(state);

          const savedKeys = loadExchangeKeys() || {};
          if (binanceApiKeyEl) binanceApiKeyEl.value = savedKeys.binance?.apiKey || "";
          if (binanceApiSecretEl) binanceApiSecretEl.value = savedKeys.binance?.apiSecret || "";
          if (bybitApiKeyEl) bybitApiKeyEl.value = savedKeys.bybit?.apiKey || "";
          if (bybitApiSecretEl) bybitApiSecretEl.value = savedKeys.bybit?.apiSecret || "";
          if (okxApiKeyEl) okxApiKeyEl.value = savedKeys.okx?.apiKey || "";
          if (okxApiSecretEl) okxApiSecretEl.value = savedKeys.okx?.apiSecret || "";
          if (pionexApiKeyEl) pionexApiKeyEl.value = savedKeys.pionex?.apiKey || "";
          if (pionexApiSecretEl) pionexApiSecretEl.value = savedKeys.pionex?.apiSecret || "";

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

          const tabSettingsDialog = document.getElementById("tabSettingsDialog");
          const tabSettingsClose = document.getElementById("tabSettingsClose");
          const tabSettingsCancel = document.getElementById("tabSettingsCancel");
          const tabSettingsApply = document.getElementById("tabSettingsApply");
          const tabSettingsMode = document.getElementById("tabSettingsMode");
          if (tabSettingsClose && tabSettingsDialog) tabSettingsClose.onclick = () => tabSettingsDialog.close();
          if (tabSettingsCancel && tabSettingsDialog) tabSettingsCancel.onclick = () => tabSettingsDialog.close();
          if (tabSettingsApply) tabSettingsApply.onclick = () => applyTabSettingsFromDialog();
          if (tabSettingsMode && !tabSettingsMode.dataset.wired) {
            tabSettingsMode.dataset.wired = "1";
            tabSettingsMode.addEventListener("change", syncTabSettingsExchangeVisibility);
          }
          if (tabSettingsDialog) {
            tabSettingsDialog.addEventListener("close", () => { tabSettingsTargetId = null; });
          }

          const apply = document.getElementById("apply");
          apply.onclick = () => {
            const prevDefault = state.defaultInterval;
            state.cols = clampInt(colsEl.value, 1, 6, state.cols);
            state.defaultInterval = String(defaultIntervalEl.value || state.defaultInterval).trim() || "60";
            state.chartHeightPx = clampInt(chartHeightEl.value, 180, 900, state.chartHeightPx);
            state.chartTheme = chartThemeEl.value === "light" ? "light" : "dark";
            state.barColorPreset = normalizeBarColorPreset(barColorPresetEl.value);
            saveState(state);

            saveExchangeKeys({
              binance: { apiKey: binanceApiKeyEl?.value || "", apiSecret: binanceApiSecretEl?.value || "" },
              bybit: { apiKey: bybitApiKeyEl?.value || "", apiSecret: bybitApiSecretEl?.value || "" },
              okx: { apiKey: okxApiKeyEl?.value || "", apiSecret: okxApiSecretEl?.value || "" },
              pionex: { apiKey: pionexApiKeyEl?.value || "", apiSecret: pionexApiSecretEl?.value || "" }
            });

            appState = state;
            updateTopbarDefaultIntervalDisplay();

            initOrUpdateMultichart(false);

            if (prevDefault !== state.defaultInterval) initOrUpdateMultichart(false);

            ensureTradingView().then(async (ok) => {
              if (!ok) return;
              const t = activeTab();
              if (t && isAutoMode(t) && !autoNeedsExchange(t)) {
                try {
                  const pos = await fetchPositionsForExchange(t.autoExchange);
                  const onlyMeta = updateAutoSymbolsFromPositions(pos);
                  refreshAutoPositionUIAfterFetch(onlyMeta, false, false);
                } catch (_) {
                  initOrUpdateMultichart(false);
                }
              } else {
                initOrUpdateMultichart(false);
              }
            });
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
            renderTabsbar();
            const cur0 = activeTab();
            if (cur0 && isAutoMode(cur0) && !autoNeedsExchange(cur0)) {
              fetchPositionsForExchange(cur0.autoExchange)
                .then((pos) => {
                  const onlyMeta = updateAutoSymbolsFromPositions(pos);
                  refreshAutoPositionUIAfterFetch(onlyMeta, true, false);
                })
                .catch(() => initOrUpdateMultichart(true));
            } else {
              initOrUpdateMultichart(true);
            }
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

        (function wireExchangeKeysModal() {
          const dlg = document.getElementById("exchangeKeysDialog");
          const saveBtn = document.getElementById("exchangeKeysModalSave");
          const cancelBtn = document.getElementById("exchangeKeysModalCancel");
          const closeBtn = document.getElementById("exchangeKeysModalClose");
          function dismiss() {
            exchangeKeysModalPendingId = null;
            if (dlg) dlg.close();
          }
          if (saveBtn) saveBtn.onclick = () => saveExchangeKeysFromModal();
          if (cancelBtn) cancelBtn.onclick = () => dismiss();
          if (closeBtn) closeBtn.onclick = () => dismiss();
          if (dlg) dlg.addEventListener("close", () => { exchangeKeysModalPendingId = null; });
        })();

        document.addEventListener("visibilitychange", () => {
          if (document.hidden) stopAutoPositionsPolling();
          else syncAutoPositionsPolling();
        });

        render();
      })();
    </script>
  </body>
</html>`;

