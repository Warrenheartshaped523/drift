# Changelog

All notable changes to drift will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
drift uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.1.0] — 2026-03-19

First public release.

### Scenes

- **constellation** — stars drift slowly across the screen, connecting nearby neighbours with dotted lines; brightness twinkles per star
- **rain** — columns of half-width katakana characters and digits fall at varying speeds with bright heads and fading trails
- **particles** — a sinusoidal flow field drives 120 glyphs across the screen, leaving ghost trails as they move
- **waveform** — three layered sine waves rendered with Unicode braille characters for sub-character precision; amplitudes breathe in and out independently

### Themes

Seven built-in color themes matched to popular terminal colorschemes: `cosmic`, `nord`, `dracula`, `catppuccin`, `gruvbox`, `forest`, `mono`

### Shell integration

Idle detection via native shell mechanisms — no background daemons:

- **zsh** — TMOUT + TRAPALRM
- **bash** — PROMPT_COMMAND with a background timer
- **fish** — `fish_prompt` / `fish_preexec` event hooks

Activate with `eval "$(drift shell-init zsh)"` (or bash/fish).

### CLI

- `drift --scene <name>` — lock to a specific scene
- `drift --theme <name>` — override the color theme
- `drift --duration <n>` — seconds per scene when cycling, 0 = no cycling
- `drift list scenes` — list available scenes
- `drift list themes` — list themes with live color swatches
- `drift config --init` — write default config to `~/.config/drift/config.toml`
- `drift shell-init zsh|bash|fish` — print shell integration snippet

### Distribution

- Single static binary, no CGO, no runtime dependencies
- Pre-built releases for macOS and Linux (amd64 + arm64)
- goreleaser pipeline with SHA-256 checksums

[0.1.0]: https://github.com/phlx0/drift/releases/tag/v0.1.0
