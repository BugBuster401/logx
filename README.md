# logx

Minimal structured logging for Go, based on `log/slog`.

This module provides:
- `logx` — a small, implementation-agnostic logging interface
- `slogx` — a `logx.Logger` implementation based on Go `log/slog`
- `loki` — an optional remote client for sending logs to Grafana Loki

Detailed documentation is provided via GoDoc inside each package.

---

## Installation

```bash
go get github.com/BugBuster401/logx