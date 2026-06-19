# Changelog

## [v2.0.0] - 2026-06-18

### Renamed
- Project renamed from **Uptime Kuma Probe Extension** to **KProbe**.
- GitHub repository moved from `K-cermak/Uptime-Kuma-Probe` to `K-cermak/KProbe`.

### Changed
- The `kprobe cron` command has been renamed to `kprobe scan` to better reflect its purpose.
- Replaced the simple keyword matching system with a new **Variables + Expressions** engine for HTTP scans:
  - `var_bool:<name>` – searches for a keyword in the response body.
  - `var_json_bool:<name>` – extracts a boolean value from the JSON response using a JSON path.
  - `var_json_number:<name>` – extracts a numeric value from the JSON response using a JSON path.
  - `expression=<expr>` – evaluates a logical expression combining variables (supports `&&`, `||`, `!`, `>`, `<`, `>=`, `<=`, `==`, `!=`).
- The CLI version is now **dynamically injected at build time** instead of being hardcoded.
- Better texts in the CLI.
- Improved texts in README and FAQ.

### Added
- New configurable key: `max_http_body_size` – sets the maximum HTTP response body size to read during scans (in MB, default 10).
- New configurable key: `output_http_info` – toggles detailed HTTP response output during `kprobe test http` commands.
- Installation is now done via `.deb` / `.rpm` packages built with nFPM, replacing the old `install.sh` script.
- Added build for Arm.

### Removed
- Removed the `keyword` field from scan configuration (replaced by Variables + Expressions).
- Removed the old `install.sh` and `uninstall.sh` scripts (replaced by system package management).

---

## [v1.0.1] - 2025-02-09
### Changed
- Fixed an incorrect URL in the web editor.

## [v1.0.0] - 2025-02-08
### App Release
- Initial public release of Uptime Kuma Probe Extension.