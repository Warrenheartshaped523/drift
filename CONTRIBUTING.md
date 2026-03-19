# Contributing to drift

Thank you for considering a contribution to drift. Here's everything you need to get started.

---

## Table of contents

- [Development setup](#development-setup)
- [Project structure](#project-structure)
- [Adding a new scene](#adding-a-new-scene)
- [Adding a new theme](#adding-a-new-theme)
- [Code style](#code-style)
- [Testing](#testing)
- [Submitting a pull request](#submitting-a-pull-request)
- [Reporting bugs](#reporting-bugs)

---

## Development setup

```bash
# Clone and enter the repo
git clone https://github.com/phlx0/drift
cd drift

# Install dependencies and tidy go.sum
make setup

# Build
make build

# Run
./drift

# Run a specific scene for quick iteration
./drift --scene rain --theme dracula
```

You will need **Go 1.23 or later**.
The only external runtime dependency is `tcell/v2` — no C library required.

---

## Project structure

```
drift/
├── main.go                    Entry point, ldflags injection
├── cmd/drift/
│   ├── root.go                CLI commands (cobra)
│   └── shell_snippets.go      Shell integration strings
├── internal/
│   ├── config/config.go       TOML config loading
│   ├── engine/engine.go       Render loop, scene lifecycle
│   └── scene/
│       ├── scene.go           Scene interface, Theme type, helpers
│       ├── constellation.go   Drifting stars with connection lines
│       ├── rain.go            Falling character rain
│       ├── particles.go       Flow-field particle system
│       └── waveform.go        Braille sine wave layers
└── .github/workflows/         CI and release automation
```

---

## Adding a new scene

1. Create `internal/scene/myscene.go` with `package scene`.

2. Implement the `Scene` interface:

   ```go
   type Scene interface {
       Name()   string
       Init(w, h int, t Theme)
       Update(dt float64)
       Draw(screen tcell.Screen)
       Resize(w, h int)
   }
   ```

3. Register it in `All()` inside `internal/scene/scene.go`:

   ```go
   func All() []Scene {
       return []Scene{
           NewConstellation(),
           NewRain(),
           NewParticles(),
           NewWaveform(),
           NewMyScene(), // add here
       }
   }
   ```

4. Add a config struct to `internal/config/config.go` if the scene has tunable knobs, and wire it to `SceneConfig`.

5. Test it:
   ```bash
   go build . && ./drift --scene myscene
   ```

### Scene guidelines

- **`Init` must be idempotent** — it is called again on every `Resize`.
- **`Draw` must not call `screen.Show()`** — the engine flushes once per frame.
- **Delta time**: `Update(dt float64)` receives seconds since last frame, capped at 100 ms by the engine.  Use `dt` for all time-based motion.
- **Respect terminal color** — always use `tcell.StyleDefault` as the base and only override the foreground.  Never hardcode a background color.
- **Handle all terminal sizes gracefully**, including very narrow (< 40 columns) or very short (< 10 rows) terminals.

---

## Adding a new theme

Open `internal/scene/scene.go` and add an entry to the `Themes` map:

```go
"mytheme": {
    Name: "mytheme",
    Palette: []RGBColor{
        {R, G, B},
        {R, G, B},
        {R, G, B},
        {R, G, B},
    },
    Dim: []RGBColor{
        // Darker / more muted versions of each Palette color.
        {R, G, B},
        {R, G, B},
        {R, G, B},
        {R, G, B},
    },
    Bright: RGBColor{R, G, B}, // near-white highlight
},
```

Run `./drift list themes` to confirm it appears.

---

## Code style

- Standard `gofmt` / `goimports` formatting.
- No external linters beyond `go vet` are required, but PRs must pass the CI lint step.
- Keep files focused.  If a scene file grows beyond ~300 lines, consider splitting helpers.
- Exported symbols need doc comments; unexported helpers are optional.

---

## Testing

```bash
make test       # unit tests with race detector
go vet ./...    # static analysis
```

Because most of the interesting code is pixel-level rendering, visual smoke tests are done manually:

```bash
./drift --scene <name> --theme <name>
```

For automated tests, prefer testing pure functions (math helpers, config parsing, theme lookups) rather than trying to mock `tcell.Screen`.

---

## Submitting a pull request

1. Fork the repo and create a branch off `main`.
2. Keep commits focused — one logical change per commit.
3. Run `make test` and `go vet ./...` before opening the PR.
4. Fill in the PR description template.
5. Screenshots or terminal recordings of new scenes / visual changes are very welcome.

We review PRs as time allows.  Patience is appreciated.

---

## Reporting bugs

Open an issue and include:

- Your OS and terminal emulator
- `drift version` output
- The theme and scene that triggered the bug
- What you expected vs. what happened
- A screenshot if the issue is visual

---

Made with care for the terminal community. ♥
